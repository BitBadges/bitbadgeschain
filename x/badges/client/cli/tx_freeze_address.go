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

func CmdFreezeAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze-address [addresses] [badge-id] [subbadge-id] [add]",
		Short: "Broadcast message freezeAddress",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAddresses := strings.Split(args[0], ",")

			argAddressesUint64 := []uint64{}
			for _, address := range argAddresses {
				addressAsUint64, err := cast.ToUint64E(address)
				if err != nil {
					return err
				}

				argAddressesUint64 = append(argAddressesUint64, addressAsUint64)
			}

			argBadgeId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argSubbadgeId, err := cast.ToUint64E(args[2])
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
				argAddressesUint64,
				argBadgeId,
				argSubbadgeId,
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
