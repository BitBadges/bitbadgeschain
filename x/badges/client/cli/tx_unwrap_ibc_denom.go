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

/*
your-cli-command unwrap-ibc-denom '{
  "creator": "your-creator-address",
  "collectionId": "your-collection-id",
  "amount": {
    "denom": "ibc/SHA256_HASH",
    "amount": "1000000"
  },
  "overrideTokenId": "optional-token-id-for-override"
}'
*/

func CmdUnwrapIBCDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unwrap-ibc-denom [tx-json]",
		Short: "Broadcast message unwrapIBCDenom",
		Args:  cobra.ExactArgs(1), // Accept exactly one argument (the JSON string)
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgUnwrapIBCDenom
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			// Set the creator to the client context's from address if not provided
			if txData.Creator == "" {
				txData.Creator = clientCtx.GetFromAddress().String()
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
