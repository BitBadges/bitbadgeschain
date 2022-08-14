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

func CmdRevokeBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-badge [addresses] [amounts] [badge-id] [subbadge-start-id] [subbadge-end-id]",
		Short: "Broadcast message revoke-badge",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAddressesUInt64, err := GetIdArrFromString(args[0])
			if err != nil {
				return err
			}

			argAmountsUInt64, err := GetIdArrFromString(args[1])
			if err != nil {
				return err
			}

			argBadgeId, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			subbadgeRanges, err := GetIdRanges(args[3], args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevokeBadge(
				clientCtx.GetFromAddress().String(),
				argAddressesUInt64,
				argAmountsUInt64,
				argBadgeId,
				subbadgeRanges,
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
