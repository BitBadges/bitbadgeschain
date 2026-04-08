package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// execNodeCLI delegates a subcommand to the Node.js bitbadges-cli binary.
// It forwards all args, pipes stdio through, and preserves exit codes.
// If the CLI is not found, it prints a helpful message and returns nil
// so chain operations are unaffected.
func execNodeCLI(subcommand string, args []string) error {
	cliPath, useNpx := findNodeCLI()
	if cliPath == "" {
		fmt.Fprintln(os.Stderr, "SDK CLI not available. Install with: npm install -g bitbadgesjs-sdk")
		fmt.Fprintln(os.Stderr, "Or set BITBADGES_SDK_CLI_PATH environment variable.")
		return nil // Don't error - chain operations should still work
	}

	var fullArgs []string
	if useNpx {
		// npx bitbadges-cli <subcommand> <args...>
		fullArgs = append([]string{"bitbadges-cli", subcommand}, args...)
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
