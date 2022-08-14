package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdPruneBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune-balances [badge-ids] [addresses]",
		Short: "Broadcast message pruneBalances",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeIdsUint64, err := GetIdArrFromString(args[0])
			if err != nil {
				return err
			}

			argAddressesUint64, err := GetIdArrFromString(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgPruneBalances(
				clientCtx.GetFromAddress().String(),
				argBadgeIdsUint64,
				argAddressesUint64,
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
