package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdHandlePendingTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle-pending-transfer [accept] [badge-id] [subbadge-id] [pending-id]",
		Short: "Broadcast message handlePendingTransfer",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAccept, err := cast.ToBoolE(args[0])
			if err != nil {
				return err
			}
			argBadgeId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argSubbadgeId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			argPendingId, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgHandlePendingTransfer(
				clientCtx.GetFromAddress().String(),
				argAccept,
				argBadgeId,
				argSubbadgeId,
				argPendingId,
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
