package cli

import (
	"encoding/json"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateAllowedTransfers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-allowed-transfers [collection-id] [allowed-transfers]",
		Short: "Broadcast message updateAllowedTransfers",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			argCollectionId := types.NewUintFromString(args[0])
			if err != nil {
				return err
			}

			var argAllowedTransfers []*types.TransferMapping
			if err := json.Unmarshal([]byte(args[1]), &argAllowedTransfers); err != nil {
				return err
			}
			
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAllowedTransfers(
				clientCtx.GetFromAddress().String(),
				argCollectionId,
				argAllowedTransfers,
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
