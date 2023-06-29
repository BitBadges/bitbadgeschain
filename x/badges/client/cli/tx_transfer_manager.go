package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateManager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-manager [collection-id] [address]",
		Short: "Broadcast message updateManager",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
return nil
			// argBadgeId := types.NewUintFromString(args[0])
			// if err != nil {
			// 	return err
			// }
			// argAddress, err := cast.ToStringE(args[1])
			// if err != nil {
			// 	return err
			// }

			// clientCtx, err := client.GetClientTxContext(cmd)
			// if err != nil {
			// 	return err
			// }

			// msg := types.NewMsgUpdateManager(
			// 	clientCtx.GetFromAddress().String(),
			// 	argBadgeId,
			// 	argAddress,
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
