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

func CmdUpdateDisallowedTransfers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-disallowed-transfers [collection-id] [disallowed-transfers]",
		Short: "Broadcast message updateDisallowedTransfers",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			argCollectionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			var argDisallowedTransfers []*types.TransferMapping
			if err := json.Unmarshal([]byte(args[1]), &argDisallowedTransfers); err != nil {
				return err
			}
			
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDisallowedTransfers(
				clientCtx.GetFromAddress().String(),
				argCollectionId,
				argDisallowedTransfers,
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
