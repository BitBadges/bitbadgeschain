package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDecrementStoreValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrement-store-value [store-id] [address] [amount] [set-to-zero-on-underflow]",
		Short: "Broadcast message decrementStoreValue",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStoreId := types.NewUintFromString(args[0])
			argAddress := args[1]
			argAmount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			argSetToZeroOnUnderflow, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDecrementStoreValue(
				clientCtx.GetFromAddress().String(),
				argStoreId,
				argAddress,
				sdkmath.NewUint(argAmount),
				argSetToZeroOnUnderflow,
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
