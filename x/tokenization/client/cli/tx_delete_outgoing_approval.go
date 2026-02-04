package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func CmdDeleteOutgoingApproval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-outgoing-approval [collection-id] [approval-id]",
		Short: "Broadcast message DeleteOutgoingApproval",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionId := args[0]
			argApprovalId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId, err := strconv.ParseUint(argCollectionId, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteOutgoingApproval(
				clientCtx.GetFromAddress().String(),
				sdkmath.NewUint(collectionId),
				argApprovalId,
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
