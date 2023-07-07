package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// string creator = 1;
//     uint64 collectionId = 2;
//     uint64 address = 3; //The address that are approved to transfer the balances.
//     repeated Balance balances = 4; //approval balances for every badgeId

func CmdSetApproval() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-approval [collection-id] [address] [balances]",
		Short: "Broadcast message setApproval",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return nil
			// 	argBadgeId := types.NewUintFromString(args[0])
			// 	if err != nil {
			// 		return err
			// 	}

			// 	argAddress, err := cast.ToStringE(args[1])
			// 	if err != nil {
			// 		return err
			// 	}

			// 	var argBalances []*types.Balance
			// 	if err := json.Unmarshal([]byte(args[2]), &argBalances); err != nil {
			// 		return err
			// 	}

			// 	clientCtx, err := client.GetClientTxContext(cmd)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	msg := types.NewMsgSetApproval(
			// 		clientCtx.GetFromAddress().String(),
			// 		argBadgeId,
			// 		argAddress,
			// 		argBalances,
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
