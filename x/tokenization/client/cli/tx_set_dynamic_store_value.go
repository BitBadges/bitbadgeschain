package cli

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdSetDynamicStoreValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-dynamic-store-value [store-id] [address] [value]",
		Short: "Broadcast message setDynamicStoreValue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId := types.NewUintFromString(args[0])
			argAddress := args[1]
			argValueInt, ok := sdkmath.NewIntFromString(args[2])
			if !ok || argValueInt.IsNegative() {
				return types.ErrInvalidRequest.Wrap("invalid value: must be a non-negative integer")
			}
			argValue := sdkmath.NewUintFromBigInt(argValueInt.BigInt())

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetDynamicStoreValue(
				clientCtx.GetFromAddress().String(),
				argStoreId,
				argAddress,
				argValue,
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
