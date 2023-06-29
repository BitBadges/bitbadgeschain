package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdForkCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-collection [collection-id]",
		Short: "Broadcast message forkCollection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
return nil
		// 	argCollectionId, err := cast.ToUint64E(args[0])
		// 	if err != nil {
		// 		return err
		// 	}

		// 	clientCtx, err := client.GetClientTxContext(cmd)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	msg := types.NewMsgForkCollection(
		// 		clientCtx.GetFromAddress().String(),
		// 		argCollectionId,
		// 	)
		// 	if err := msg.ValidateBasic(); err != nil {
		// 		return err
		// 	}
		// 	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
