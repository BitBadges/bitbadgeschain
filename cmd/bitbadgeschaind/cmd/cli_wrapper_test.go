package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/spf13/viper"
)

// TestRootCmdKeysRegisteredExactlyOnce verifies BB-361: `keys` must appear
// exactly once in the root command tree and in `bb --help` output, and no
// top-level command name may be registered multiple times.
func TestRootCmdKeysRegisteredExactlyOnce(t *testing.T) {
	rootCmd := NewRootCmd()

	counts := make(map[string]int)
	for _, c := range rootCmd.Commands() {
		counts[c.Name()]++
	}

	if counts["keys"] != 1 {
		t.Fatalf("expected `keys` to be registered exactly once, got %d", counts["keys"])
	}

	for name, n := range counts {
		if n > 1 {
			t.Errorf("top-level command %q registered %d times", name, n)
		}
	}

	var helpOut bytes.Buffer
	rootCmd.SetOut(&helpOut)
	rootCmd.SetErr(&helpOut)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd --help failed: %v", err)
	}

	keysLines := 0
	for _, line := range strings.Split(helpOut.String(), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "keys ") || trimmed == "keys" {
			keysLines++
		}
	}
	if keysLines != 1 {
		t.Fatalf("expected `keys` to appear once in --help output, found %d occurrences:\n%s", keysLines, helpOut.String())
	}
}

// TestSDKForwarderRootFlagStrippingAndHomeBypass verifies BB-361 end-to-end
// through svrcmd.Execute (the real main.go entrypoint with --home, --log_level,
// --trace, etc.):
//   - Leading Cosmos root flags BEFORE the subcommand (`--home`, `--log_level`,
//     `--trace`) are stripped before forwarding to bitbadges-cli, and client.toml
//     is never loaded or created even when --home is missing/empty/read-only.
//   - Flags appearing AFTER the subcommand (`bb dev --home /foo`, `bb cli --home /foo`)
//     are preserved verbatim and NOT silently eaten.
func TestSDKForwarderRootFlagStrippingAndHomeBypass(t *testing.T) {
	// Create a fake bitbadges-cli script that records its invocation args to a file.
	tmpDir := t.TempDir()
	recordFile := filepath.Join(tmpDir, "invoked.txt")
	fakeCLI := filepath.Join(tmpDir, "fake-bitbadges-cli.sh")
	script := "#!/bin/sh\nprintf '%s' \"$*\" > " + recordFile + "\n"
	if err := os.WriteFile(fakeCLI, []byte(script), 0o755); err != nil {
		t.Fatalf("writing fake CLI: %v", err)
	}
	t.Setenv("BITBADGES_SDK_CLI_PATH", fakeCLI)

	// Use a read-only parent so any attempt to mkdir/write <home>/config/client.toml
	// would immediately fail with EACCES.
	roParent := filepath.Join(tmpDir, "readonly-parent")
	if err := os.MkdirAll(roParent, 0o555); err != nil {
		t.Fatalf("creating read-only parent: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(roParent, 0o755) })

	missingHome := filepath.Join(roParent, "missing-home")

	cases := []struct {
		name     string
		args     []string
		wantArgs string
	}{
		{
			name:     "bb --home /tmp/empty dev --help strips leading root --home",
			args:     []string{"--home", missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "bb --home=/tmp/empty dev --help strips leading --home= syntax",
			args:     []string{"--home=" + missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "bb --log_level info --home /tmp/empty dev --help strips all leading root flags",
			args:     []string{"--log_level", "info", "--home", missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "bb --trace --log_level=debug --home /tmp/empty dev --help strips bool and value root flags",
			args:     []string{"--trace", "--log_level=debug", "--home", missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "bb dev --home /foo preserves subcommand --home flag verbatim",
			args:     []string{"dev", "--home", "/foo"},
			wantArgs: "dev --home /foo",
		},
		{
			name:     "bb --home /tmp/empty dev --home /foo strips root --home and keeps subcommand --home",
			args:     []string{"--home", missingHome, "dev", "--home", "/foo"},
			wantArgs: "dev --home /foo",
		},
		{
			name:     "bb cli --home /foo preserves --home /foo and does not silently exit with empty args",
			args:     []string{"cli", "--home", "/foo"},
			wantArgs: "--home /foo",
		},
		{
			name:     "bb --home /tmp/empty cli --help strips leading root --home",
			args:     []string{"--home", missingHome, "cli", "--help"},
			wantArgs: "--help",
		},
		{
			name:     "bb --log_level info --home /tmp/empty cli --home /foo strips root flags and keeps subcommand flag",
			args:     []string{"--log_level", "info", "--home", missingHome, "cli", "--home", "/foo"},
			wantArgs: "--home /foo",
		},
		{
			name:     "bb dev --help without root flags works cleanly",
			args:     []string{"dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "bb --home /tmp/empty build --help works cleanly",
			args:     []string{"--home", missingHome, "build", "--help"},
			wantArgs: "build --help",
		},
		{
			name:     "bb --log_level info --home /tmp/empty tx status --help works cleanly",
			args:     []string{"--log_level", "info", "--home", missingHome, "tx", "status", "--help"},
			wantArgs: "tx status --help",
		},
		{
			name:     "bb tx status --home /foo preserves subcommand --home flag verbatim",
			args:     []string{"tx", "status", "--home", "/foo"},
			wantArgs: "tx status --home /foo",
		},
	}

	origOSArgs := os.Args
	t.Cleanup(func() {
		os.Args = origOSArgs
		viper.Reset()
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Remove(recordFile)
			viper.Reset()
			os.Args = append([]string{"bitbadgeschaind"}, tc.args...)

			rootCmd := NewRootCmd()
			if err := svrcmd.Execute(rootCmd, "", missingHome); err != nil {
				t.Fatalf("svrcmd.Execute(%v) failed: %v", tc.args, err)
			}

			gotBytes, err := os.ReadFile(recordFile)
			if err != nil {
				t.Fatalf("fake bitbadges-cli was not invoked (args=%v): %v", tc.args, err)
			}
			if got := string(gotBytes); got != tc.wantArgs {
				t.Fatalf("forwarded args = %q, want %q", got, tc.wantArgs)
			}

			if _, err := os.Stat(filepath.Join(missingHome, "config", "client.toml")); err == nil {
				t.Fatalf("SDK forwarder unexpectedly created client.toml under %s", missingHome)
			}
		})
	}
}

func TestIsSDKForwarderCmdClassification(t *testing.T) {
	rootCmd := NewRootCmd()

	forwarders := [][]string{
		{"dev"},
		{"build"},
		{"check"},
		{"cli"},
		{"portfolio"}, // deprecated forwarder
		{"tx", "status"},
		{"tx", "wait"},
	}
	for _, path := range forwarders {
		cmd, _, err := rootCmd.Find(path)
		if err != nil || cmd == nil {
			t.Fatalf("Find(%v) failed: %v", path, err)
		}
		if !isSDKForwarderCmd(cmd) {
			t.Errorf("expected %v to be classified as SDK forwarder", path)
		}
	}

	chainNatives := [][]string{
		{"start"},
		{"keys"},
		{"query"},
		{"tx"},
		{"status"},
		{"init"},
	}
	for _, path := range chainNatives {
		cmd, _, err := rootCmd.Find(path)
		if err != nil || cmd == nil {
			t.Fatalf("Find(%v) failed: %v", path, err)
		}
		if isSDKForwarderCmd(cmd) {
			t.Errorf("expected chain-native %v NOT to be classified as SDK forwarder", path)
		}
	}
}
