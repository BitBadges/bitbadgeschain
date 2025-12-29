package cli

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCastVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cast-vote [collection-id] [approval-level] [approver-address] [approval-id] [proposal-id] [yes-weight]",
		Short: "Broadcast message castVote",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionId := types.NewUintFromString(args[0])
			argApprovalLevel := args[1]
			argApproverAddress := args[2]
			argApprovalId := args[3]
			argProposalId := args[4]
			argYesWeight := types.NewUintFromString(args[5])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCastVote(
				clientCtx.GetFromAddress().String(),
				argCollectionId,
				argApprovalLevel,
				argApproverAddress,
				argApprovalId,
				argProposalId,
				argYesWeight,
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

