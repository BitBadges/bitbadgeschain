package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// newHomeCmd builds a command carrying a --home flag with the given default,
// mirroring how the root command supplies one to `pre-upgrade`.
func newHomeCmd(defaultHome string) *cobra.Command {
	cmd := &cobra.Command{Use: "pre-upgrade"}
	cmd.Flags().String("home", defaultHome, "directory for config and data")
	return cmd
}

// TestResolveHomePrefersDaemonHomeOverFlagDefault is the single most important
// property of this hook, and it is the one with no visible symptom when broken.
//
// cosmovisor invokes the hook as exec.Command(bin, "pre-upgrade") with no
// arguments whatsoever. --home is therefore never passed, and the root command's
// persistent flag supplies its compiled-in *default* rather than the home the
// node actually runs from. An operator whose node lives anywhere else — a
// non-default --home, a systemd unit with its own path — would have had a
// pristine default-path config.toml rewritten while their real one was left at
// "flood", and the hook would have exited 0 reporting success.
//
// The node then dies at the upgrade height with a config error, which is exactly
// the outcome the hook exists to prevent, and the logs say the pre-upgrade
// succeeded.
//
// DAEMON_HOME is the value to trust because cosmovisor always exports it.
func TestResolveHomePrefersDaemonHomeOverFlagDefault(t *testing.T) {
	daemonHome := t.TempDir()
	flagDefault := t.TempDir()

	t.Setenv("DAEMON_HOME", daemonHome)

	// No --home passed, exactly as cosmovisor invokes it.
	got := resolveHome(newHomeCmd(flagDefault))

	if got != daemonHome {
		t.Fatalf("resolveHome() = %q, want DAEMON_HOME %q.\n"+
			"Falling back to the --home default (%q) rewrites the wrong config.toml "+
			"and still reports success, which is the silent failure this hook must not have.",
			got, daemonHome, flagDefault)
	}
}

// TestResolveHomePrefersExplicitFlagOverDaemonHome pins the other direction: an
// operator running the hook by hand with --home means it, even inside an
// environment where cosmovisor has exported DAEMON_HOME.
func TestResolveHomePrefersExplicitFlagOverDaemonHome(t *testing.T) {
	daemonHome := t.TempDir()
	explicit := t.TempDir()

	t.Setenv("DAEMON_HOME", daemonHome)

	cmd := newHomeCmd(t.TempDir())
	if err := cmd.Flags().Set("home", explicit); err != nil {
		t.Fatalf("setting --home: %v", err)
	}

	if got := resolveHome(cmd); got != explicit {
		t.Fatalf("resolveHome() = %q, want the explicitly passed --home %q", got, explicit)
	}
}

// TestResolveHomeFallsBackToFlagWhenDaemonHomeIsAbsent covers running the hook
// outside cosmovisor, and the case where a systemd unit exports DAEMON_HOME as
// an empty string. Those are one test, not two: os.Getenv returns "" for both
// unset and empty, so a separate "unset" case would execute an identical path
// and only look like extra coverage.
//
// Treating "" as a home would send the hook to /config/config.toml, fail to
// find it, and exit 30 — aborting an upgrade that had nothing wrong with it.
func TestResolveHomeFallsBackToFlagWhenDaemonHomeIsAbsent(t *testing.T) {
	t.Setenv("DAEMON_HOME", "")

	flagDefault := t.TempDir()
	got := resolveHome(newHomeCmd(flagDefault))

	if got == "" {
		t.Fatal("resolveHome() = \"\", an absent DAEMON_HOME must not be treated as a home")
	}
	if got != flagDefault {
		t.Fatalf("resolveHome() = %q, want the --home value %q", got, flagDefault)
	}
}
