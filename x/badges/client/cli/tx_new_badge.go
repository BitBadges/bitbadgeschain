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

func CmdNewBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-badge [uri] [permissions] [subasset-uris] [bytes-string] [default-supply] [subasset-supplys] [subasset-amounts] [freeze-start] [freeze-end] [standard]",
		Short: "Broadcast message newBadge",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argUri := args[0]
			argSubassetUris := args[2]
			_ = argSubassetUris
			argBytesStr := []byte(args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permissions, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			defaultSupply, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			argSupplysStringArr := strings.Split(args[5], ",")

			argSupplysUInt64 := []uint64{}
			for _, supply := range argSupplysStringArr {
				println(supply)
				supplyAsUint64, err := cast.ToUint64E(supply)
				if err != nil {
					return err
				}

				argSupplysUInt64 = append(argSupplysUInt64, supplyAsUint64)
			}

			argAmountsStringArr := strings.Split(args[6], ",")

			argAmountsUInt64 := []uint64{}
			for _, amount := range argAmountsStringArr {
				amountAsUint64, err := cast.ToUint64E(amount)
				if err != nil {
					return err
				}

				argAmountsUInt64 = append(argAmountsUInt64, amountAsUint64)
			}

			argStartAddress, err := cast.ToUint64E(args[7])
			if err != nil {
				return err
			}
			argEndAddress, err := cast.ToUint64E(args[8])
			if err != nil {
				return err
			}

			argStandard, err := cast.ToUint64E(args[9])
			if err != nil {
				return err
			}

			//TODO: parse differences between uris and subasseturis


			msg := types.NewMsgNewBadge(
				clientCtx.GetFromAddress().String(),
				types.UriObject{
					Uri: []byte(argUri),
				},
				permissions,
				argBytesStr,
				defaultSupply,
				argAmountsUInt64,
				argSupplysUInt64,
				[]*types.IdRange{
					{
						Start: argStartAddress,
						End:   argEndAddress,
					},
				},
				argStandard,
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
