package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDeleteDynamicStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-dynamic-store [store-id]",
		Short: "Broadcast message deleteDynamicStore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId := types.NewUintFromString(args[0])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteDynamicStore(
				clientCtx.GetFromAddress().String(),
				argStoreId,
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
