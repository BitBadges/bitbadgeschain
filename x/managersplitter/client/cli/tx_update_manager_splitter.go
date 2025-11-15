package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	"github.com/cosmos/cosmos-sdk/client"
)

func CmdUpdateManagerSplitter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-manager-splitter [tx-json]",
		Short: "Broadcast message updateManagerSplitter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgUpdateManagerSplitter
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			txData.Admin = clientCtx.GetFromAddress().String()

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

