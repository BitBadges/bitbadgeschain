package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClaimBadge() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "claim-badge [claim-id] [collection-id] [solutions]",
		Short: "Broadcast message claimBadge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			argClaimId := types.NewUintFromString(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argCollectionId := types.NewUintFromString(args[1])
			if err != nil {
				return err
			}

			argSolutionsJsonArr, err := parseJsonArr(args[2])
			if err != nil {
				return err
			}
			
			solutions := []*types.ChallengeSolution{}
			for _, solutionInterface := range argSolutionsJsonArr {
				solution := solutionInterface.(*types.ChallengeSolution)
				solutions = append(solutions, solution)
			}

			argSolutions := solutions

			msg := types.NewMsgClaimBadge(
				clientCtx.GetFromAddress().String(),
				argClaimId,
				argCollectionId,
				argSolutions,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
