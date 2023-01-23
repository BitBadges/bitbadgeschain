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

func CmdSetApproval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-approval [amount] [address] [badge-id] [badge-id-start] [badge-id-end] [expiry-time]",
		Short: "Broadcast message setApproval",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAddress, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argBadgeId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argBadgeIdRanges, err := GetIdRanges(args[3], args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetApproval(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				argAddress,
				[]*types.Balance{
					{
						Balance: argAmount,
						BadgeIds: argBadgeIdRanges,
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
