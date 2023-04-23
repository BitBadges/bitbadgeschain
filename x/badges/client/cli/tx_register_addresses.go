package cli

import (
	"encoding/json"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRegisterAddresses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-addresses [addresses-to-register]",
		Short: "Broadcast message registerAddresses",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var argAddressesToRegister []string
			err = json.Unmarshal([]byte(args[0]), &argAddressesToRegister)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterAddresses(
				clientCtx.GetFromAddress().String(),
				argAddressesToRegister,
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
