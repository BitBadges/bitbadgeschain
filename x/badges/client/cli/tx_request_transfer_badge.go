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

func CmdRequestTransferBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-transfer-badge [from] [amount] [badge-id] [subbadge-id-start] [subbadge-id-end] [expiry-time]",
		Short: "Broadcast message requestTransferBadge",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argFrom, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argAmount, err := cast.ToUint64E(args[1])
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

			msg := types.NewMsgRequestTransferBadge(
				clientCtx.GetFromAddress().String(),
				argFrom,
				argAmount,
				argBadgeId,
				[]*types.IdRange{
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
