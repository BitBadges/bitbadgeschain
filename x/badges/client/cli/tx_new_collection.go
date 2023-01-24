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

func CmdNewCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-collection [collectionUri] [badgeUri] [permissions] [bytes-string] [subasset-supplys] [subasset-amounts] [standard]",
		Short: "Broadcast message newCollection",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionUri := args[0]
			argBadgeUri := args[1]
			_ = argBadgeUri

			argBytesStr := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permissions, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			argSupplysUInt64, err := GetIdArrFromString(args[4])
			if err != nil {
				return err
			}

			argAmountsUInt64, err := GetIdArrFromString(args[5])
			if err != nil {
				return err
			}

			argStandard, err := cast.ToUint64E(args[6])
			if err != nil {
				return err
			}

			argBadgeSupplys := make([]*types.BadgeSupplyAndAmount, len(argSupplysUInt64))
			for i := 0; i < len(argSupplysUInt64); i++ {
				argBadgeSupplys[i] = &types.BadgeSupplyAndAmount{
					Supply: argSupplysUInt64[i],
					Amount: argAmountsUInt64[i],
				}
			}

			msg := types.NewMsgNewCollection(
				clientCtx.GetFromAddress().String(),
				argStandard,
				argBadgeSupplys,
				argCollectionUri,
				argBadgeUri,
				permissions,
				[]*types.TransferMapping{},
				[]*types.TransferMapping{},
				argBytesStr,
				[]*types.Transfers{}, //TODO: add whitelist capabilities to CLI
				[]*types.Claim{},
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
