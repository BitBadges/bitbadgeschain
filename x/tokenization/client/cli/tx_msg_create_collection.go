package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMsgCreateCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collection [tx-json-or-file]",
		Short: "Broadcast message createCollection",
		Long:  "Accepts JSON either inline or from a file path. If the argument is a valid file path, it will read the JSON from that file.",
		Args:  cobra.ExactArgs(1), // Accept exactly one argument (the JSON string or file path)
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON, err := ReadJSONFromFileOrString(args[0])
			if err != nil {
				return err
			}

			var txData types.MsgCreateCollection
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
