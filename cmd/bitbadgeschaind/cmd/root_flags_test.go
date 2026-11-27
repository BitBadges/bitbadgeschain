package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/version"
)

func TestRootCmdHasUniqueChildCommandNames(t *testing.T) {
	rootCmd := NewRootCmd()

	counts := map[string]int{}
	for _, sub := range rootCmd.Commands() {
		name := sub.Name()
		if name == "" {
			t.Errorf("rootCmd contains an unnamed child command (Use=%q)", sub.Use)
			continue
		}
		counts[name]++
	}

	for name, count := range counts {
		if count > 1 {
			t.Errorf("rootCmd registered child command %q %d times, want 1", name, count)
		}
	}

	if counts["keys"] != 1 {
		t.Fatalf("rootCmd has %d 'keys' commands, want 1", counts["keys"])
	}
}

func TestKeysCommandHasNoEmptySubcommands(t *testing.T) {
	rootCmd := NewRootCmd()

	keysCmd, _, err := rootCmd.Find([]string{"keys"})
	if err != nil || keysCmd == nil || keysCmd.Name() != "keys" {
		t.Fatalf("failed to find 'keys' command: cmd=%v err=%v", keysCmd, err)
	}

	seen := map[string]bool{}
	for _, sub := range keysCmd.Commands() {
		if sub.Name() == "" {
			t.Errorf("keys command contains an empty-name subcommand (Use=%q)", sub.Use)
		}
		seen[sub.Name()] = true
	}

	for _, want := range []string{"add", "list", "show", "delete", "unsafe-export-eth-key", "unsafe-import-eth-key"} {
		if !seen[want] {
			t.Errorf("keys command missing expected subcommand %q", want)
		}
	}

	var buf bytes.Buffer
	keysCmd.SetOut(&buf)
	keysCmd.SetErr(&buf)
	if err := keysCmd.Usage(); err != nil {
		t.Fatalf("keysCmd.Usage() failed: %v", err)
	}
	lines := strings.Split(buf.String(), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "Available Commands:" && i+1 < len(lines) {
			if strings.TrimSpace(lines[i+1]) == "" {
				t.Errorf("keys --help has a blank first row under Available Commands:\n%s", buf.String())
			}
			break
		}
	}
}

func writeStubCLI(t *testing.T, script string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "bitbadges-cli")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatalf("write stub cli: %v", err)
	}
	return p
}

func TestRootVersionFlagWithoutSDKCLI(t *testing.T) {
	t.Setenv("BITBADGES_SDK_CLI_PATH", "")
	t.Setenv("PATH", t.TempDir())

	origVer, origCommit := version.Version, version.Commit
	version.Version = "v36-test"
	version.Commit = "abc1234"
	defer func() {
		version.Version = origVer
		version.Commit = origCommit
	}()

	rootCmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs([]string{"--version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("bb --version failed: %v (stderr=%q)", err, errBuf.String())
	}

	got := strings.TrimSpace(out.String())
	want := "bitbadgeschaind v36-test (abc1234)"
	if got != want {
		t.Fatalf("bb --version output = %q, want %q", got, want)
	}
}

func TestRootVersionFlagWithSDKCLI(t *testing.T) {
	stub := writeStubCLI(t, `if [ "$1" = "--version" ]; then echo "0.46.0"; exit 0; fi; exit 1`)
	t.Setenv("BITBADGES_SDK_CLI_PATH", stub)

	origVer, origCommit := version.Version, version.Commit
	version.Version = "v36"
	version.Commit = ""
	defer func() {
		version.Version = origVer
		version.Commit = origCommit
	}()

	rootCmd := NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("bb --version failed: %v", err)
	}

	got := strings.TrimSpace(out.String())
	want := "bitbadgeschaind v36\nbitbadges-cli 0.46.0"
	if got != want {
		t.Fatalf("bb --version output = %q, want %q", got, want)
	}
}

func TestVersionSubcommandUnbroken(t *testing.T) {
	origVer := version.Version
	version.Version = "v36-subcmd"
	defer func() { version.Version = origVer }()

	rootCmd := NewRootCmd()
	verCmd, _, err := rootCmd.Find([]string{"version"})
	if err != nil || verCmd == nil || verCmd.Name() != "version" {
		t.Fatalf("failed to resolve 'version' subcommand: cmd=%v err=%v", verCmd, err)
	}

	var out bytes.Buffer
	verCmd.SetOut(&out)
	verCmd.SetErr(&out)
	if err := verCmd.RunE(verCmd, nil); err != nil {
		t.Fatalf("bb version failed: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "v36-subcmd" {
		t.Fatalf("bb version output = %q, want %q", got, "v36-subcmd")
	}
}

func TestRootHelpJSONFlagMissingSDKCLIExits127(t *testing.T) {
	t.Setenv("BITBADGES_SDK_CLI_PATH", "")
	t.Setenv("PATH", t.TempDir())

	code := -1
	orig := osExit
	osExit = func(c int) { code = c }
	defer func() { osExit = orig }()

	rootCmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs([]string{"--help-json"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("bb --help-json returned unexpected error: %v", err)
	}
	if code != exitNoSDKCLI {
		t.Fatalf("exit code = %d, want %d", code, exitNoSDKCLI)
	}
	if !strings.Contains(errBuf.String(), "bitbadges-cli") {
		t.Errorf("expected install hint on stderr, got %q", errBuf.String())
	}
}

func TestRootHelpJSONFlagForwardsToSDKCLI(t *testing.T) {
	stub := writeStubCLI(t, `printf '{"args":"%s"}\n' "$*"`)
	t.Setenv("BITBADGES_SDK_CLI_PATH", stub)

	rootCmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs([]string{"--help-json", "--", "--quiet"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("bb --help-json failed: %v (stderr=%q)", err, errBuf.String())
	}

	got := strings.TrimSpace(out.String())
	want := `{"args":"--help-json --quiet"}`
	if got != want {
		t.Fatalf("bb --help-json output = %q, want %q", got, want)
	}
}

func TestMatchRootHelpJSON(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantOK    bool
		wantExtra []string
	}{
		{
			name:      "bare --help-json",
			args:      []string{"--help-json"},
			wantOK:    true,
			wantExtra: []string{},
		},
		{
			name:      "--help-json with trailing extra flags",
			args:      []string{"--help-json", "--quiet", "--compact"},
			wantOK:    true,
			wantExtra: []string{"--quiet", "--compact"},
		},
		{
			name:      "leading flag before --help-json",
			args:      []string{"--quiet", "--help-json"},
			wantOK:    true,
			wantExtra: []string{"--quiet"},
		},
		{
			name:      "--help-json with double-dash separator",
			args:      []string{"--help-json", "--", "--foo"},
			wantOK:    true,
			wantExtra: []string{"--foo"},
		},
		{
			name:   "subcommand before --help-json routes to forwarder",
			args:   []string{"build", "--help-json"},
			wantOK: false,
		},
		{
			name:   "help flag takes precedence",
			args:   []string{"--help", "--help-json"},
			wantOK: false,
		},
		{
			name:   "version flag does not match",
			args:   []string{"--version"},
			wantOK: false,
		},
		{
			name:   "unknown flag without --help-json does not match",
			args:   []string{"--unknown"},
			wantOK: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			extra, ok := MatchRootHelpJSON(tc.args)
			if ok != tc.wantOK {
				t.Fatalf("MatchRootHelpJSON(%v) ok = %v, want %v", tc.args, ok, tc.wantOK)
			}
			if tc.wantOK && !reflect.DeepEqual(extra, tc.wantExtra) {
				t.Fatalf("MatchRootHelpJSON(%v) extra = %v, want %v", tc.args, extra, tc.wantExtra)
			}
		})
	}
}

func TestRootUnknownFlagAndUnknownCommandStillError(t *testing.T) {
	t.Run("unknown flag", func(t *testing.T) {
		rootCmd := NewRootCmd()
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"--definitely-not-a-flag"})
		if err := rootCmd.Execute(); err == nil {
			t.Fatal("expected error for unknown flag, got nil")
		}
	})

	t.Run("unknown command", func(t *testing.T) {
		rootCmd := NewRootCmd()
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"definitely-not-a-command"})
		if err := rootCmd.Execute(); err == nil {
			t.Fatal("expected error for unknown command, got nil")
		}
	})
}
