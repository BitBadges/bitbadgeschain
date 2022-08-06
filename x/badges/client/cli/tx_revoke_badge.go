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

func CmdRevokeBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-badge [addresses] [amounts] [badge-id] [subbadge-start-id] [subbadge-end-id]",
		Short: "Broadcast message revoke-badge",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argAddressesStringArr := strings.Split(args[0], ",")

			argAddressesUInt64 := []uint64{}
			for _, address := range argAddressesStringArr {
				addressAsUint64, err := cast.ToUint64E(address)
				if err != nil {
					return err
				}

				argAddressesUInt64 = append(argAddressesUInt64, addressAsUint64)
			}

			argAmountsStringArr := strings.Split(args[1], ",")

			argAmountsUInt64 := []uint64{}
			for _, amount := range argAmountsStringArr {
				amountAsUint64, err := cast.ToUint64E(amount)
				if err != nil {
					return err
				}

				argAmountsUInt64 = append(argAmountsUInt64, amountAsUint64)
			}
			argBadgeId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argSubbadgeStartId, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argSubbadgeEndId, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevokeBadge(
				clientCtx.GetFromAddress().String(),
				argAddressesUInt64,
				argAmountsUInt64,
				argBadgeId,
				types.NumberRange{
					Start: argSubbadgeStartId,
					End:   argSubbadgeEndId,
				},
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
