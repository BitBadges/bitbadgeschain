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

func CmdFreezeAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze-address [start-addresses] [end-addresses] [badge-id] [add]",
		Short: "Broadcast message freezeAddress",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStartAddresses := args[0]
			argEndAddresses := args[1]

			addressRanges, err := GetIdRanges(argStartAddresses, argEndAddresses)
			if err != nil {
				return err
			}

			argBadgeId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			argAdd, err := cast.ToBoolE(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgFreezeAddress(
				clientCtx.GetFromAddress().String(),
				addressRanges,
				argBadgeId,
				argAdd,
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
