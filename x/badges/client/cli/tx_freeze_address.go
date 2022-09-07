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

func CmdFreezeAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze-address [badge-id] [add] [start-addresses] [end-addresses]",
		Short: "Broadcast message freezeAddress",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argAdd, err := cast.ToBoolE(args[1])
			if err != nil {
				return err
			}

			argStartAddresses := args[2]
			argEndAddresses := args[3]

			addressRanges, err := GetIdRanges(argStartAddresses, argEndAddresses)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgFreezeAddress(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				argAdd,
				addressRanges,
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
