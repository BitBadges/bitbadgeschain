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

func CmdPruneBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune-balances [badge-ids] [addresses]",
		Short: "Broadcast message pruneBalances",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeIds := strings.Split(args[0], ",")

			argBadgeIdsUint64 := []uint64{}
			for _, badgeId := range argBadgeIds {
				badgeIdAsUint64, err := cast.ToUint64E(badgeId)
				if err != nil {
					return err
				}

				argBadgeIdsUint64 = append(argBadgeIdsUint64, badgeIdAsUint64)
			}

			argAddresses := strings.Split(args[0], ",")

			argAddressesUint64 := []uint64{}
			for _, address := range argAddresses {
				addressAsUint64, err := cast.ToUint64E(address)
				if err != nil {
					return err
				}

				argAddressesUint64 = append(argAddressesUint64, addressAsUint64)
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
