package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateDynamicStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-dynamic-store [store-id] [default-value]",
		Short: "Broadcast message updateDynamicStore",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId := types.NewUintFromString(args[0])
			defaultValue, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDynamicStore(
				clientCtx.GetFromAddress().String(),
				argStoreId,
				defaultValue,
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
