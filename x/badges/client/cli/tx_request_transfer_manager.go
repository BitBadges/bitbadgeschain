package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRequestUpdateManager() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-transfer-manager [badge-id] [add]",
		Short: "Broadcast message requestUpdateManager",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
return nil
		// 	argBadgeId := types.NewUintFromString(args[0])
		// 	if err != nil {
		// 		return err
		// 	}

		// 	argAdd, err := cast.ToBoolE(args[0])
		// 	if err != nil {
		// 		return err
		// 	}

		// 	clientCtx, err := client.GetClientTxContext(cmd)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	msg := types.NewMsgRequestUpdateManager(
		// 		clientCtx.GetFromAddress().String(),
		// 		argBadgeId,
		// 		argAdd,
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
