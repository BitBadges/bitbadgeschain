package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdHandlePendingTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle-pending-transfer [badge-id] [starting-pending-ids] [ending-pending-ids] [action]",
		Short: "Broadcast message handlePendingTransfer",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argStartingNonces := args[1]
			argEndingNonces := args[2]
			argNonceRanges, err := GetIdRanges(argStartingNonces, argEndingNonces)
			if err != nil {
				return err
			}

			argAction, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgHandlePendingTransfer(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				argNonceRanges,
				[]uint64{ argAction },
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
