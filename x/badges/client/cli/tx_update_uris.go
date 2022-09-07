package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateUris() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-uris [badge-id] [uri] [subasset-uri]",
		Short: "Broadcast message updateUris",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argUri := args[1]
			argSubassetUri := args[2]
			_ = argSubassetUri

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			uriObject, err := GetUriObject(argUri, argSubassetUri)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateUris(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				uriObject,
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
