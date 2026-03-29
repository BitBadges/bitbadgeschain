package cli

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateDynamicStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-dynamic-store [default-value]",
		Short: "Broadcast message createDynamicStore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defaultValue, ok := sdkmath.NewIntFromString(args[0])
			if !ok {
				return types.ErrInvalidRequest.Wrap("invalid default value: must be a non-negative integer")
			}
			if defaultValue.IsNegative() {
				return types.ErrInvalidRequest.Wrap("invalid default value: must be a non-negative integer")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDynamicStore(
				clientCtx.GetFromAddress().String(),
				sdkmath.NewUintFromBigInt(defaultValue.BigInt()),
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
