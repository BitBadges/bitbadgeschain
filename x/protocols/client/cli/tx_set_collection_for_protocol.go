package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSetCollectionForProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-collection-for-protocol [name] [collectionId]",
		Short: "Broadcast message setCollectionForProtocol",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetCollectionForProtocol(
				clientCtx.GetFromAddress().String(),
				args[0],
				sdkmath.NewUintFromString(args[1]),
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
