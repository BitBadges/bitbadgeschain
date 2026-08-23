package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/app"
)

// preUpgradeFatal is the exit code cosmovisor treats as "pre-upgrade failed,
// abort the upgrade". Exit code 1 means "no such command, carry on", which is
// the opposite of what we want when the edit genuinely failed, so any real
// failure must exit with this instead of returning an error to cobra.
const preUpgradeFatal = 30

var mempoolTypeRe = regexp.MustCompile(`(?m)^(\s*type\s*=\s*)"[^"]*"`)

// NewPreUpgradeCmd returns the `pre-upgrade` command that cosmovisor invokes on
// the new binary after swapping the symlink but before starting it.
//
// v34 needs it because there is no single config.toml that both binaries
// accept:
//
//   - v33 runs CometBFT v0.38, which knows only the "flood" and "nop" mempool
//     types and panics with `unknown mempool type: "app"` on anything else.
//   - v34 runs CometBFT v0.39 and cosmos/evm v0.7, which enables the EVM mempool
//     by default (app.toml mempool.max-txs defaults to 0, i.e. unbounded) and
//     enforces via server/config.ValidateCrossConfig that config.toml's
//     mempool.type is "app". Otherwise it refuses to start.
//
// So the value cannot be set ahead of the upgrade, and without this hook every
// node would crash on restart at the upgrade height and stay down until an
// operator hand-edited config.toml. Rewriting it here closes that window.
//
// Operators who do not run cosmovisor must make the same edit manually; see the
// v34 release notes.
func NewPreUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "Pre-upgrade hook invoked by cosmovisor before starting the new binary",
		Long: `Applies configuration changes that must happen between stopping the old
binary and starting the new one.

For v34 this sets config.toml's mempool.type to "app", which CometBFT v0.39
requires when the EVM mempool is enabled and which CometBFT v0.38 (v33) rejects
outright, so it cannot be set in advance.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			home := resolveHome(cmd)

			changed, err := ensureAppMempoolType(filepath.Join(home, "config", "config.toml"))
			if err != nil {
				// Do not return the error: cobra would exit 1, which cosmovisor
				// reads as "this command does not exist" and would continue into
				// a node that cannot start.
				fmt.Fprintf(os.Stderr, "pre-upgrade failed: %v\n", err)
				os.Exit(preUpgradeFatal)
			}

			if changed {
				cmd.Printf("pre-upgrade: set mempool.type = \"app\" in %s\n", home)
			} else {
				cmd.Printf("pre-upgrade: mempool.type already \"app\" in %s, nothing to do\n", home)
			}
			return nil
		},
	}

	return cmd
}

// resolveHome picks the home directory this hook should edit.
//
// Order matters. cosmovisor invokes us as `exec.Command(bin, "pre-upgrade")`
// with no arguments at all, so --home is never passed and the root command's
// persistent flag supplies its *default* rather than the node's actual home.
// Trusting that default silently rewrote the wrong config.toml while still
// reporting success. DAEMON_HOME, which cosmovisor always exports, is the
// authoritative value; an explicitly passed --home still wins for operators
// running this by hand.
func resolveHome(cmd *cobra.Command) string {
	if f := cmd.Flags().Lookup("home"); f != nil && f.Changed {
		if v := f.Value.String(); v != "" {
			return v
		}
	}
	if h := os.Getenv("DAEMON_HOME"); h != "" {
		return h
	}
	if f := cmd.Flags().Lookup("home"); f != nil {
		if v := f.Value.String(); v != "" {
			return v
		}
	}
	return app.DefaultNodeHome
}

// ensureAppMempoolType rewrites the `type` key inside config.toml's [mempool]
// section, leaving the rest of the file (comments, operator customizations)
// byte-identical. Reports whether it changed anything.
func ensureAppMempoolType(path string) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", path, err)
	}
	content := string(raw)

	start := strings.Index(content, "[mempool]")
	if start < 0 {
		return false, fmt.Errorf("no [mempool] section in %s", path)
	}

	// The section runs until the next top-level table header. [p2p] may appear
	// before [mempool] in the generated file, so scan forward rather than
	// assuming any particular ordering.
	end := len(content)
	if next := strings.Index(content[start+len("[mempool]"):], "\n["); next >= 0 {
		end = start + len("[mempool]") + next + 1
	}

	section := content[start:end]
	loc := mempoolTypeRe.FindStringSubmatchIndex(section)
	if loc == nil {
		return false, fmt.Errorf("no type key in the [mempool] section of %s", path)
	}
	if strings.Contains(section[loc[0]:loc[1]], `"app"`) {
		return false, nil
	}

	updated := section[:loc[0]] + section[loc[2]:loc[3]] + `"app"` + section[loc[1]:]
	if err := os.WriteFile(path, []byte(content[:start]+updated+content[end:]), 0o644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}
