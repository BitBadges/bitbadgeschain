package cli

import (
	"strconv"
	"strings"

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
		Use:   "transfer-badge [from] [to] [amount] [badge-id] [subbadge-id-start] [subbadge-id-end] [expiry-time]",
		Short: "Broadcast message transferBadge",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argFrom, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argToAddresses := strings.Split(args[1], ",")

			argToAddressesUint64 := []uint64{}
			for _, toAddress := range argToAddresses {
				addressAsUint64, err := cast.ToUint64E(toAddress)
				if err != nil {
					return err
				}

				argToAddressesUint64 = append(argToAddressesUint64, addressAsUint64)
			}

			argAmounts := strings.Split(args[2], ",")

			argAmountsUint64 := []uint64{}
			for _, amount := range argAmounts {
				amountAsUint64, err := cast.ToUint64E(amount)
				if err != nil {
					return err
				}

				argAmountsUint64 = append(argAmountsUint64, amountAsUint64)
			}

			argBadgeId, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argSubbadgeIdStart, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argSubbadgeIdEnd, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			argExpirationTime, err := cast.ToUint64E(args[6])
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
				[]*types.IdRange{
					{
						Start: argSubbadgeIdStart,
						End:   argSubbadgeIdEnd,
					},
				},
				argExpirationTime,
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
