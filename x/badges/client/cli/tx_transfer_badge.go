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

func CmdTransferBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-badge [from] [to] [amount] [badge-id] [subbadge-id-start] [subbadge-id-end] [expiry-time] [cant-cancel-before-time]",
		Short: "Broadcast message transferBadge",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argFrom, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argToAddressesUint64, err := GetIdArrFromString(args[1])
			if err != nil {
				return err
			}

			argAmountsUint64, err := GetIdArrFromString(args[2])
			if err != nil {
				return err
			}

			argBadgeId, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argSubbadgeRanges, err := GetIdRanges(args[4], args[5])
			if err != nil {
				return err
			}

			argExpirationTime, err := cast.ToUint64E(args[6])
			if err != nil {
				return err
			}

			argCantCancelBeforeTime, err := cast.ToUint64E(args[7])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferBadge(
				clientCtx.GetFromAddress().String(),
				argFrom,
				argToAddressesUint64,
				argAmountsUint64,
				argBadgeId,
				argSubbadgeRanges,
				argExpirationTime,
				argCantCancelBeforeTime,
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
