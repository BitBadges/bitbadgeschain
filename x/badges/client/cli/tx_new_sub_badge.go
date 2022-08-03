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

func CmdNewSubBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-sub-badge [id] [supplys] [amounts]",
		Short: "Creates a subasset of the badge ID. Must be executed by the manager. CLI args delimited by commas.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argSupplysStringArr := strings.Split(args[1], ",")


			argSupplysUInt64 := []uint64{}
			for _, supply := range argSupplysStringArr {
				println(supply)
				supplyAsUint64, err := cast.ToUint64E(supply)
				if err != nil {
					return err
				}

				argSupplysUInt64 = append(argSupplysUInt64, supplyAsUint64)
			}

			argAmountsStringArr := strings.Split(args[2], ",")

			argAmountsUInt64 := []uint64{}
			for _, amount := range argAmountsStringArr {
				amountAsUint64, err := cast.ToUint64E(amount)
				if err != nil {
					return err
				}

				argAmountsUInt64 = append(argAmountsUInt64, amountAsUint64)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewSubBadge(
				clientCtx.GetFromAddress().String(),
				argId,
				argSupplysUInt64,
				argAmountsUInt64,
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
