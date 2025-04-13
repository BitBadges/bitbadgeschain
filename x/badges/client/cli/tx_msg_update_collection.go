package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMsgUpdateCollection() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "update-collection [tx-json]",
		Short: "Broadcast message updateCollection",
		Args:  cobra.ExactArgs(1), // Accept exactly one argument (the JSON string)
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgUpdateCollection
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			txData.Creator = clientCtx.GetFromAddress().String()

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
