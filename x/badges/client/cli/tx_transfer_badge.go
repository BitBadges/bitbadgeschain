package cli

import (
	"encoding/json"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdTransferBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-badge [collection-id] [from] [transfers]",
		Short: "Broadcast message transferBadge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionId := types.NewUintFromString(args[0])
			if err != nil {
				return err
			}

			argFrom, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}

			var argTransfers []*types.Transfer
			err = json.Unmarshal([]byte(args[2]), &argTransfers)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferBadge(
				clientCtx.GetFromAddress().String(),
				argCollectionId,
				argFrom,
				argTransfers,
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
