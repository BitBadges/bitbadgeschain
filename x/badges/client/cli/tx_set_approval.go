package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdSetApproval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-approval [amount] [address] [badge-id] [subbadge-id-start] [subbadge-id-end] [expiry-time]",
		Short: "Broadcast message setApproval",
		Args:  cobra.ExactArgs(6),
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
			argSubbadgeIdStart, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argSubbadgeIdEnd, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argExpirationTime, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetApproval(
				clientCtx.GetFromAddress().String(),
				argAmount,
				argAddress,
				argBadgeId,
				[]*types.NumberRange{
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
