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

func CmdNewSubBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-sub-badge [id] [supplys] [amounts]",
		Short: "Creates a subasset of the badge ID. Must be executed by the manager. CLI args delimited by commas.",
		Args:  cobra.ExactArgs(3),
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

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewSubBadge(
				clientCtx.GetFromAddress().String(),
				argId,
				argSupplysUint64,
				argAmountsUint64,
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
