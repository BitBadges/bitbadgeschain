package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdHandlePendingTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle-pending-transfer [accept] [badge-id] [starting-pending-id] [ending-pending-id] [forceful-accept]",
		Short: "Broadcast message handlePendingTransfer",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAccept, err := cast.ToBoolE(args[0])
			if err != nil {
				return err
			}
			argBadgeId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argStartingNonce, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			argEndingNonce, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			// argStartingNonces := strings.Split(args[2], ",")

			// argStartingNoncesUint64 := []uint64{}
			// for _, nonce := range argStartingNonces {
			// 	nonceAsUint64, err := cast.ToUint64E(nonce)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	argStartingNoncesUint64 = append(argStartingNoncesUint64, nonceAsUint64)
			// }

			// argEndingNonces := strings.Split(args[3], ",")

			// argEndingNoncesUint64 := []uint64{}
			// for _, nonce := range argEndingNonces {
			// 	nonceAsUint64, err := cast.ToUint64E(nonce)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	argEndingNoncesUint64 = append(argEndingNoncesUint64, nonceAsUint64)
			// }

			// if len(argStartingNoncesUint64) != len(argEndingNoncesUint64) {
			// 	return types.ErrInvalidArgumentLengths
			// }

			// nonceRanges := []*types.NumberRange{}
			// for i := 0; i < len(argStartingNoncesUint64); i++ {
			// 	nonceRanges = append(nonceRanges, &types.NumberRange{
			// 		Start: argStartingNoncesUint64[i],
			// 		End:   argEndingNoncesUint64[i],
			// 	})
			// }

			forcefulAccept, err := cast.ToBoolE(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgHandlePendingTransfer(
				clientCtx.GetFromAddress().String(),
				argAccept,
				argBadgeId,
				[]*types.NumberRange{
					{
						Start: argStartingNonce,
						End:   argEndingNonce,
					},
				},
				forcefulAccept,
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
