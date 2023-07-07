package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateCollectionApprovedTransfers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-allowed-transfers [collection-id] [allowed-transfers]",
		Short: "Broadcast message UpdateCollectionApprovedTransfers",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return nil

			// argCollectionId := types.NewUintFromString(args[0])
			// if err != nil {
			// 	return err
			// }

			// var argApprovedTransfers []*types.CollectionApprovedTransfer
			// if err := json.Unmarshal([]byte(args[1]), &argApprovedTransfers); err != nil {
			// 	return err
			// }

			// clientCtx, err := client.GetClientTxContext(cmd)
			// if err != nil {
			// 	return err
			// }

			// msg := types.NewMsgUpdateCollectionApprovedTransfers(
			// 	clientCtx.GetFromAddress().String(),
			// 	argCollectionId,
			// 	argApprovedTransfers,
			// )

			// if err := msg.ValidateBasic(); err != nil {
			// 	return err
			// }

			// return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
