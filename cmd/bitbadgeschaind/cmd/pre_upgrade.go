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
//
//   - v34 runs CometBFT v0.39 and cosmos/evm v0.7, which enables the EVM mempool
//     and enforces via server/config.ValidateCrossConfig that config.toml's
//     mempool.type is "app". Otherwise it refuses to start.
//
//     The mempool is enabled by the start command, not by app.toml - editing
//     app.toml alone neither enables nor disables it. See initCometBFTConfig
//     in config.go for why.
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

	start, end, ok := mempoolSectionBounds(content)
	if !ok {
		return false, fmt.Errorf("no [mempool] section in %s", path)
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

// mempoolSectionBounds locates the [mempool] table in a config.toml and returns
// the byte range it occupies. It scans LINES, because the obvious substring
// search was wrong in two ways that both end with a dead node:
//
//   - strings.Index(content, "[mempool]") matches inside comments and strings.
//     An operator note such as `# see the [mempool] section below` sitting under
//     [p2p] won the race, so `type = "app"` was written into [p2p], the real
//     mempool.type stayed "flood", and the hook exited 0 printing success. The
//     node then died at the upgrade height: exactly the silent-success failure
//     this hook exists to eliminate. Stock 0.38/0.39 templates carry only one
//     literal [mempool], but operators annotate their configs.
//
//   - Finding the end with strings.Index(rest, "\n[") let a multi-line TOML
//     string containing a line that starts with "[" truncate the scan. A type
//     key below such a string fell outside the section, so a second one was
//     inserted and config.toml became a duplicate-key parse error.
//
// The section starts at its header line and runs to the start of the next
// top-level table header line, ignoring comment lines and anything inside a
// multi-line string. ok is false when there is no [mempool] header line.
func mempoolSectionBounds(content string) (start, end int, ok bool) {
	inMultiline := false
	offset := 0
	found := false

	for _, line := range strings.SplitAfter(content, "\n") {
		if line == "" {
			break
		}
		lineStart := offset
		offset += len(line)

		if inMultiline {
			inMultiline = !closesMultilineString(line)
			continue
		}

		trimmed := strings.TrimSpace(line)
		// A comment is not a section header, whatever it happens to spell.
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		if name, isHeader := tomlTableName(trimmed); isHeader {
			if found {
				return start, lineStart, true
			}
			if name == "mempool" {
				start, found = lineStart, true
			}
			continue
		}

		if opensMultilineString(line) {
			inMultiline = true
		}
	}

	if !found {
		return 0, 0, false
	}
	// [mempool] is the last table in the file.
	return start, len(content), true
}

// tomlTableName returns the table name of a header line such as `[mempool]` or
// `[mempool] # comment`. Array-of-tables headers ([[x]]) count as headers too:
// they end the preceding section just the same.
func tomlTableName(trimmed string) (string, bool) {
	if !strings.HasPrefix(trimmed, "[") {
		return "", false
	}
	closeIdx := strings.Index(trimmed, "]")
	if closeIdx < 0 {
		return "", false
	}
	return strings.Trim(trimmed[1:closeIdx], "[ \t"), true
}

// opensMultilineString reports whether line leaves a multi-line TOML string
// open. An odd number of delimiters on the line means it is still open at the
// end of the line.
func opensMultilineString(line string) bool {
	return strings.Count(line, `"""`)%2 == 1 || strings.Count(line, "'''")%2 == 1
}

// closesMultilineString reports whether line closes an already-open multi-line
// TOML string.
func closesMultilineString(line string) bool {
	return strings.Contains(line, `"""`) || strings.Contains(line, "'''")
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
	// contents are still only in the page cache, which reaches the same
	// truncated-config failure mode by another route.
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
	// fsync the directory too. The rename itself is a directory metadata change,
	// and on several filesystems that entry can still be only in the page cache
	// after os.Rename returns: a crash immediately afterwards would lose the
	// edit even though the file contents were durable. The file is never
	// corrupted either way, so this is completeness rather than a fix, and a
	// failure to sync the directory must not abort an upgrade whose config edit
	// already landed.
	if dirFile, err := os.Open(dir); err == nil {
		_ = dirFile.Sync()
		_ = dirFile.Close()
	}
	return nil
}
