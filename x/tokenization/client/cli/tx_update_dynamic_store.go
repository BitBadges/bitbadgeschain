package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUpdateDynamicStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-dynamic-store [store-id] [default-value] [global-enabled]",
		Short: "Broadcast message updateDynamicStore",
		Long: `Update a dynamic store.

Arguments:
  store-id: The ID of the dynamic store to update
  default-value: The default numeric value for uninitialized addresses (non-negative integer)
  global-enabled: The global kill switch state (true = enabled, false = disabled/halted)`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId := types.NewUintFromString(args[0])
			defaultValueInt, ok := sdkmath.NewIntFromString(args[1])
			if !ok || defaultValueInt.IsNegative() {
				return types.ErrInvalidRequest.Wrap("invalid default value: must be a non-negative integer")
			}
			defaultValue := sdkmath.NewUintFromBigInt(defaultValueInt.BigInt())
			globalEnabled, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDynamicStoreWithGlobalEnabled(
				clientCtx.GetFromAddress().String(),
				argStoreId,
				defaultValue,
				globalEnabled,
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
