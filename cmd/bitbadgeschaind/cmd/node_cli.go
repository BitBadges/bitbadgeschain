package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// exitNoSDKCLI is the exit status when an SDK verb is invoked but
// bitbadges-cli cannot be found. 127 matches the shell's "command not
// found", so scripts and agents can tell a missing dependency apart from a
// real command failure.
const exitNoSDKCLI = 127

// osExit is swapped out in tests.
var osExit = os.Exit

// execNodeCLI delegates a subcommand to the Node.js bitbadges-cli binary.
// It forwards all args, pipes stdio through, and preserves exit codes.
// If the CLI is not found, it prints an install hint and exits 127. This
// path only runs for SDK verbs (`bb build`, `bb check`, `bb cli ...`), so
// chain-native commands are unaffected. It used to return nil (exit 0),
// which made `bb check tx.json` look like a passing check with empty
// output.
func execNodeCLI(subcommand string, args []string) error {
	cliPath, useNpx := findNodeCLI()
	if cliPath == "" {
		fmt.Fprintf(os.Stderr, "SDK CLI not available: `bb %s` needs bitbadges-cli. Install with: npm install -g bitbadges (provides bitbadges-cli command)\n", subcommand)
		fmt.Fprintln(os.Stderr, "Or set BITBADGES_SDK_CLI_PATH environment variable.")
		osExit(exitNoSDKCLI)
		return nil
	}

	var fullArgs []string
	if useNpx {
		// npx -p bitbadges bitbadges-cli <subcommand> <args...>
		fullArgs = append([]string{"-p", "bitbadges", "bitbadges-cli", subcommand}, args...)
	} else {
		// <binary> <subcommand> <args...>
		fullArgs = append([]string{subcommand}, args...)
	}

	execCmd := exec.Command(cliPath, fullArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	if err := execCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}

// findNodeCLI locates the bitbadges-cli binary.
// Returns the path to use and whether npx is being used as the wrapper.
func findNodeCLI() (path string, useNpx bool) {
	// 1. Check BITBADGES_SDK_CLI_PATH env var for an explicit path
	if envPath := os.Getenv("BITBADGES_SDK_CLI_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, false
		}
	}

	// 2. Try bitbadges-cli directly on PATH
	if p, err := exec.LookPath("bitbadges-cli"); err == nil {
		return p, false
	}

	// 3. Fall back to npx
	if p, err := exec.LookPath("npx"); err == nil {
		return p, true
	}

	return "", false
}
