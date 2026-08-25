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

// mempoolTypeRe matches an uncommented `type = <string>` key. Both TOML string
// forms are accepted: CometBFT writes double quotes, but a hand-edited config
// may use TOML literal (single-quoted) strings. Failing to match one would make
// us append a second `type` key and turn config.toml into a duplicate-key
// error, which stops the node just as surely as the wrong value.
var mempoolTypeRe = regexp.MustCompile(`(?m)^(\s*type\s*=\s*)("[^"]*"|'[^']*')`)

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

// ensureAppMempoolType sets the `type` key inside config.toml's [mempool]
// section to "app", leaving the rest of the file (comments, operator
// customizations) byte-identical. If the section has no uncommented type key,
// one is inserted rather than treated as an error. Reports whether it changed
// anything; the write itself is atomic (see atomicWriteFile).
//
// A missing [mempool] section IS still an error: that is not a config this hook
// can reason about, and guessing where to put the section could produce a file
// CometBFT rejects outright.
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

	var updated string
	if loc := mempoolTypeRe.FindStringSubmatchIndex(section); loc != nil {
		if strings.Contains(section[loc[0]:loc[1]], `"app"`) {
			return false, nil
		}
		updated = section[:loc[0]] + section[loc[2]:loc[3]] + `"app"` + section[loc[1]:]
	} else {
		// The section exists but carries no uncommented type key — either it was
		// never written, or the operator commented it out. Erroring here would
		// exit preUpgradeFatal and abort the upgrade, halting the node at the
		// upgrade height: exactly the downtime this hook exists to prevent.
		// Insert the key instead, immediately after the section header, and
		// leave any commented-out line in place as operator context.
		updated = insertAfterSectionHeader(section, `type = "app"`)
	}

	return true, atomicWriteFile(path, content[:start]+updated+content[end:])
}

// insertAfterSectionHeader puts line directly beneath the section's header line.
func insertAfterSectionHeader(section, line string) string {
	nl := strings.Index(section, "\n")
	if nl < 0 {
		// The header is the final line of the file and has no trailing newline.
		return section + "\n" + line + "\n"
	}
	return section[:nl+1] + line + "\n" + section[nl+1:]
}

// atomicWriteFile replaces path's contents without ever leaving a truncated
// file behind.
//
// os.WriteFile truncates and then writes. Interrupted in between — power loss,
// ENOSPC, cosmovisor killing us — it leaves the operator's config.toml empty or
// half-written, losing persistent_peers, seeds, external_address and moniker at
// the exact moment of a coordinated upgrade. Writing a sibling temp file,
// fsyncing it, and renaming makes the replacement atomic: the config is either
// wholly old or wholly new. A .bak of the original is kept alongside so an
// operator can diff or restore by hand.
func atomicWriteFile(path, content string) error {
	// Follow symlinks first. Operators commonly symlink config.toml into a
	// managed directory; renaming onto the link path would silently replace the
	// link with a regular file and detach it from whatever manages it.
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}

	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}

	if original, err := os.ReadFile(path); err == nil {
		// Best effort: a failed backup must not block the upgrade edit itself.
		_ = os.WriteFile(path+".bak", original, mode)
	}

	dir := filepath.Dir(path)
	// The temp file must live in the same directory as the target: os.Rename is
	// only atomic within a filesystem.
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file in %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op once the rename succeeds

	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return fmt.Errorf("writing %s: %w", tmpName, err)
	}
	// fsync before rename: without it the rename can land while the new file's
	// contents are still only in the page cache, which is the same truncated
	// -config failure mode by another route.
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("syncing %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", tmpName, err)
	}
	// CreateTemp makes the file 0600; restore the original permissions.
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("chmod %s: %w", tmpName, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replacing %s: %w", path, err)
	}
	return nil
}
