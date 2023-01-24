package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClaimBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-badge [claim-id]",
		Short: "Broadcast message claimBadge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// TODO:
			return types.ErrNotImplemented
			// argClaimId, err := cast.ToUint64E(args[0])
			// if err != nil {
			// 	return err
			// }

			// clientCtx, err := client.GetClientTxContext(cmd)
			// if err != nil {
			// 	return err
			// }

			// msg := types.NewMsgClaimBadge(
			// 	clientCtx.GetFromAddress().String(),
			// 	argClaimId,
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
