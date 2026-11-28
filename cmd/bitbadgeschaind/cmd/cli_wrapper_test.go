package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

// TestSDKForwarderBypassesClientConfigWithMissingHome verifies BB-361:
// SDK-forwarded commands (`bb dev --help`, `bb build --help`, `bb cli --help`,
// `bb tx status --help`) must bypass Cosmos client.toml / node config loading
// so they succeed even when --home points to an empty, missing, or read-only
// directory, and must never create client.toml as a side effect.
func TestSDKForwarderBypassesClientConfigWithMissingHome(t *testing.T) {
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
			name:     "dev --help with missing unwritable home",
			args:     []string{"--home", missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "dev --help with --home= syntax",
			args:     []string{"--home=" + missingHome, "dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "dev --help without --home flag",
			args:     []string{"dev", "--help"},
			wantArgs: "dev --help",
		},
		{
			name:     "build --help with missing unwritable home",
			args:     []string{"--home", missingHome, "build", "--help"},
			wantArgs: "build --help",
		},
		{
			name:     "cli --help with missing unwritable home",
			args:     []string{"--home", missingHome, "cli", "--help"},
			wantArgs: "--help",
		},
		{
			name:     "tx status --help with missing unwritable home",
			args:     []string{"--home", missingHome, "tx", "status", "--help"},
			wantArgs: "tx status --help",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Remove(recordFile)

			rootCmd := NewRootCmd()
			rootCmd.SetArgs(tc.args)
			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("Execute(%v) failed: %v", tc.args, err)
			}

			gotBytes, err := os.ReadFile(recordFile)
			if err != nil {
				t.Fatalf("fake bitbadges-cli was not invoked: %v", err)
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
