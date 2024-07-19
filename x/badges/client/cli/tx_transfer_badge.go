package cli

import (
	"strconv"

	"bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

/*
your-cli-command transfer-badges '{
  "creator": "your-creator-address",
  "collectionId": "your-collection-id",
  "transfers": [
    {
      "from": "from-address",
      "toAddresses": ["to-address1", "to-address2"],
      "balances": [
				...
      ],
      "precalculateBalancesFromApproval": {...}, // Populate with details
      "merkleProofs": [...], // Populate with proofs
      "memo": "Transfer memo"
    }
  ]
}'
*/

func CmdTransferBadges() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-badges [tx-json]",
		Short: "Broadcast message transferBadges",
		Args:  cobra.ExactArgs(1), // Accept exactly one argument (the JSON string)
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgTransferBadges
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
