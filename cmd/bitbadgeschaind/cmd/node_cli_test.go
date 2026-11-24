package cmd

import "testing"

// When bitbadges-cli is unresolvable, an SDK verb must not exit 0: agents
// and scripts would read `bb check tx.json` as a passing check.
func TestExecNodeCLIMissingSDKCLIExitsNonZero(t *testing.T) {
	t.Setenv("BITBADGES_SDK_CLI_PATH", "")
	t.Setenv("PATH", t.TempDir()) // no bitbadges-cli, no npx

	code := -1
	orig := osExit
	osExit = func(c int) { code = c }
	defer func() { osExit = orig }()

	_ = execNodeCLI("check", []string{"tx.json"})
	if code != exitNoSDKCLI {
		t.Fatalf("exit code = %d, want %d", code, exitNoSDKCLI)
	}
}
