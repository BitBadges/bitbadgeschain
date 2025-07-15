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

func CmdSetValidBadgeIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-valid-badge-ids [tx-json]",
		Short: "Broadcast message setValidBadgeIds",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgSetValidBadgeIds
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
