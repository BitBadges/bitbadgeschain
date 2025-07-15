package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func CmdSetOutgoingApproval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-outgoing-approval [collection-id] [approval-json]",
		Short: "Broadcast message SetOutgoingApproval",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionId := args[0]
			argApprovalJSON := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId, err := strconv.ParseUint(argCollectionId, 10, 64)
			if err != nil {
				return err
			}

			var approval types.UserOutgoingApproval
			if err := clientCtx.Codec.UnmarshalJSON([]byte(argApprovalJSON), &approval); err != nil {
				return err
			}

			msg := types.NewMsgSetOutgoingApproval(
				clientCtx.GetFromAddress().String(),
				sdkmath.NewUint(collectionId),
				&approval,
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
