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

func CmdMintBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-sub-badge [id] [supplys] [amounts] [new-collection-uri] [new-badge-uri]",
		Short: "Creates a subasset of the badge ID. Must be executed by the manager. CLI args delimited by commas.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argSupplysUint64, err := GetIdArrFromString(args[1])
			if err != nil {
				return err
			}

			argAmountsUint64, err := GetIdArrFromString(args[2])
			if err != nil {
				return err
			}

			argCollectionUri, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}

			// argBadgeUrsi, err := cast.ToStringE(args[4])
			// if err != nil {
			// 	return err
			// }

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argBadgeSupplys := make([]*types.BadgeSupplyAndAmount, len(argSupplysUint64))
			for i := 0; i < len(argSupplysUint64); i++ {
				argBadgeSupplys[i] = &types.BadgeSupplyAndAmount{
					Supply: argSupplysUint64[i],
					Amount: argAmountsUint64[i],
				}
			}

			msg := types.NewMsgMintBadge(
				clientCtx.GetFromAddress().String(),
				argId,
				argBadgeSupplys,
				[]*types.Transfers{}, //TODO:
				[]*types.Claim{},
				argCollectionUri,
				[]*types.BadgeUri{},
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
