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
		Use:   "request-transfer-badge [from] [amount] [badge-id] [subbadge-id]",
		Short: "Broadcast message requestTransferBadge",
		Args:  cobra.ExactArgs(4),
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
			argSubbadgeId, err := cast.ToUint64E(args[3])
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
				argSubbadgeId,
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
