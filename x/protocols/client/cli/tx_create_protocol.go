package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreateProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-protocol [name] [uri] [customData] [isFrozen]",
		Short: "Broadcast message createProtocol",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			//parse args[3] to bool
			isFrozen := false
			if args[3] == "true" {
				isFrozen = true
			}

			msg := types.NewMsgCreateProtocol(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				args[2],
				isFrozen,
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
