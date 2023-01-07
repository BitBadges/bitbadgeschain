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

func CmdNewBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-badge [uri] [permissions] [subasset-uris] [bytes-string] [default-supply] [subasset-supplys] [subasset-amounts] [freeze-starts] [freeze-ends] [standard]",
		Short: "Broadcast message newBadge",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argUri := args[0]
			argSubassetUris := args[2]
			_ = argSubassetUris
			argBytesStr := args[3]

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

			argSupplysUInt64, err := GetIdArrFromString(args[5])
			if err != nil {
				return err
			}

			argAmountsUInt64, err := GetIdArrFromString(args[6])
			if err != nil {
				return err
			}

			freezeRanges, err := GetIdRanges(args[7], args[8])
			if err != nil {
				return err
			}

			argStandard, err := cast.ToUint64E(args[9])
			if err != nil {
				return err
			}

			uriObject, err := GetUriObject(argUri, argSubassetUris)
			if err != nil {
				return err
			}

			argSubassetSupplysAndAmounts := make([]*types.SubassetSupplyAndAmount, len(argSupplysUInt64))
			for i := 0; i < len(argSupplysUInt64); i++ {
				argSubassetSupplysAndAmounts[i] = &types.SubassetSupplyAndAmount{
					Supply: argSupplysUInt64[i],
					Amount: argAmountsUInt64[i],
				}
			}

			msg := types.NewMsgNewBadge(
				clientCtx.GetFromAddress().String(),
				argStandard,
				defaultSupply,
				argSubassetSupplysAndAmounts,
				uriObject,
				permissions,
				freezeRanges,
				argBytesStr,
				[]*types.WhitelistMintInfo{}, //TODO: add whitelist capabilities to CLI
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
