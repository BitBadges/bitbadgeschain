package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeConfig drops content into a temp config.toml and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("seeding config: %v", err)
	}
	return path
}

func readConfig(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading config: %v", err)
	}
	return string(raw)
}

// countMempoolTypeKeys counts uncommented `type =` keys inside [mempool].
// Two would be a TOML duplicate-key error and would stop the node just as
// surely as the wrong value, so every success case asserts exactly one.
func countMempoolTypeKeys(content string) int {
	start := strings.Index(content, "[mempool]")
	if start < 0 {
		return 0
	}
	rest := content[start+len("[mempool]"):]
	if next := strings.Index(rest, "\n["); next >= 0 {
		rest = rest[:next]
	}
	n := 0
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "type") && strings.Contains(trimmed, "=") {
			n++
		}
	}
	return n
}

func TestEnsureAppMempoolType(t *testing.T) {
	tests := []struct {
		name string
		in   string
		// wantErr: the call must fail and leave the file untouched.
		wantErr bool
		// wantChanged: the return value of the first call.
		wantChanged bool
		// wantContains: substrings the rewritten file must contain.
		wantContains []string
		// wantNotContains: substrings the rewritten file must not contain.
		wantNotContains []string
	}{
		{
			name: "no mempool section is an error",
			in: `[p2p]
seeds = "abc@1.2.3.4:26656"
`,
			wantErr: true,
		},
		{
			name: "rewrites a double-quoted flood value",
			in: `[mempool]
recheck = true
type = "flood"
size = 5000

[p2p]
seeds = "abc@1.2.3.4:26656"
`,
			wantChanged:     true,
			wantContains:    []string{`type = "app"`, "recheck = true", "size = 5000", `seeds = "abc@1.2.3.4:26656"`},
			wantNotContains: []string{`"flood"`},
		},
		{
			name: "rewrites a single-quoted flood value",
			in: `[mempool]
type = 'flood'
size = 5000
`,
			wantChanged:     true,
			wantContains:    []string{`type = "app"`},
			wantNotContains: []string{"'flood'"},
		},
		{
			name: "inserts the key when the section has none",
			in: `[mempool]
recheck = true
size = 5000

[p2p]
seeds = "abc@1.2.3.4:26656"
`,
			wantChanged:  true,
			wantContains: []string{`type = "app"`, "recheck = true", `seeds = "abc@1.2.3.4:26656"`},
		},
		{
			name: "inserts the key when the only occurrence is commented out",
			in: `[mempool]
# type = "flood"
recheck = true
`,
			wantChanged: true,
			// The commented line is a comment, not a key: leave it alone and
			// add a real one. Removing it would lose operator context.
			wantContains: []string{`type = "app"`, `# type = "flood"`},
		},
		{
			name: "inserts the key when the section is empty and last",
			in: `[p2p]
seeds = "abc@1.2.3.4:26656"

[mempool]
`,
			wantChanged:  true,
			wantContains: []string{`type = "app"`, `seeds = "abc@1.2.3.4:26656"`},
		},
		{
			name: "already app is a no-op",
			in: `[mempool]
type = "app"
size = 5000
`,
			wantChanged:  false,
			wantContains: []string{`type = "app"`, "size = 5000"},
		},
		{
			name: "leaves a type key in a different section alone",
			in: `[mempool]
recheck = true

[tx_index]
indexer = "kv"
type = "psql"
`,
			wantChanged: true,
			// The [tx_index] type must survive verbatim.
			wantContains: []string{`type = "app"`, `type = "psql"`, `indexer = "kv"`},
		},
		{
			name: "leaves a type key in a section BEFORE mempool alone",
			in: `[tx_index]
type = "psql"

[mempool]
recheck = true
`,
			wantChanged:  true,
			wantContains: []string{`type = "app"`, `type = "psql"`},
		},
		{
			name: "rewrites an indented key",
			in: `[mempool]
  type = "flood"
`,
			wantChanged:     true,
			wantContains:    []string{`type = "app"`},
			wantNotContains: []string{`"flood"`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := writeConfig(t, tc.in)

			changed, err := ensureAppMempoolType(path)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got changed=%v", changed)
				}
				if got := readConfig(t, path); got != tc.in {
					t.Fatalf("file was modified on the error path:\n%s", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}

			got := readConfig(t, path)
			for _, want := range tc.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("result missing %q:\n%s", want, got)
				}
			}
			for _, unwanted := range tc.wantNotContains {
				if strings.Contains(got, unwanted) {
					t.Errorf("result still contains %q:\n%s", unwanted, got)
				}
			}
			if n := countMempoolTypeKeys(got); n != 1 {
				t.Errorf("[mempool] has %d type keys, want exactly 1 (duplicates are a TOML error):\n%s", n, got)
			}

			// Idempotence: a second run must report no change and produce a
			// byte-identical file. cosmovisor can invoke pre-upgrade more than
			// once (retried upgrade, operator re-running it by hand).
			changed2, err := ensureAppMempoolType(path)
			if err != nil {
				t.Fatalf("second run errored: %v", err)
			}
			if changed2 {
				t.Errorf("second run reported a change; not idempotent")
			}
			if got2 := readConfig(t, path); got2 != got {
				t.Errorf("second run altered the file:\n--- first ---\n%s\n--- second ---\n%s", got, got2)
			}
		})
	}
}

func TestEnsureAppMempoolTypeMissingFile(t *testing.T) {
	_, err := ensureAppMempoolType(filepath.Join(t.TempDir(), "nope.toml"))
	if err == nil {
		t.Fatal("expected an error for a missing config.toml")
	}
}

// TestEnsureAppMempoolTypePreservesMode guards the atomic-write rewrite: a
// temp-file + rename must not silently widen or narrow the operator's
// permissions on config.toml.
func TestEnsureAppMempoolTypePreservesMode(t *testing.T) {
	path := writeConfig(t, "[mempool]\ntype = \"flood\"\n")
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	if _, err := ensureAppMempoolType(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %o, want 600", got)
	}
}

// TestEnsureAppMempoolTypeFollowsSymlink guards against replacing a symlinked
// config.toml with a regular file. Operators commonly symlink config into a
// managed directory; os.Rename onto the link path would break that.
func TestEnsureAppMempoolTypeFollowsSymlink(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "real-config.toml")
	if err := os.WriteFile(real, []byte("[mempool]\ntype = \"flood\"\n"), 0o644); err != nil {
		t.Fatalf("seeding: %v", err)
	}
	link := filepath.Join(dir, "config.toml")
	if err := os.Symlink(real, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	changed, err := ensureAppMempoolType(link)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatal("expected a change")
	}

	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("the symlink was replaced by a regular file")
	}
	if got := readConfig(t, real); !strings.Contains(got, `type = "app"`) {
		t.Fatalf("the symlink target was not updated:\n%s", got)
	}
}

// TestEnsureAppMempoolTypeLeavesNoTempFiles ensures the atomic write cleans up
// after itself and does not litter the operator's config directory.
func TestEnsureAppMempoolTypeLeavesNoTempFiles(t *testing.T) {
	path := writeConfig(t, "[mempool]\ntype = \"flood\"\n")
	if _, err := ensureAppMempoolType(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if name == "config.toml" || name == "config.toml.bak" {
			continue
		}
		t.Errorf("stray file left behind: %s", name)
	}
}
