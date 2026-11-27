package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/cmd/bitbadgeschaind/cmd"
)

func main() {
	if extraArgs, ok := cmd.MatchRootHelpJSON(os.Args[1:]); ok {
		if err := cmd.ExecRootHelpJSON(extraArgs); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
