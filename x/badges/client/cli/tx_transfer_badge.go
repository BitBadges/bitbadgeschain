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

func CmdTransferBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-badge [from] [to] [amount] [collection-id] [badge-id-start] [badge-id-end] [expiry-time] [cant-cancel-before-time]",
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

			argAmount, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			argCollectionId, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argBadgeIdRanges, err := GetIdRanges(args[4], args[5])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferBadge(
				clientCtx.GetFromAddress().String(),
				argCollectionId,
				argFrom,
				[]*types.Transfers{
					{
						ToAddresses: argToAddressesUint64,
						Balances: []*types.Balance{
							{
								Balance: argAmount,
								BadgeIds: argBadgeIdRanges,
							},
						},
					},
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
