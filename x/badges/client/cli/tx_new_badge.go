package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdNewBadge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-badge [uri] [permissions] [subasset-uris] [metadata-hash] [default-supply]",
		Short: "Broadcast message newBadge",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argUri := args[0]
			argSubassetUris := args[2]
			argMetadataHash := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permissions, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			defaultSupply, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewBadge(
				clientCtx.GetFromAddress().String(),
				argUri,
				permissions,
				argSubassetUris,
				argMetadataHash,
				defaultSupply,
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
