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

func CmdUpdateUris() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-uris [collection-id] [uri] [badge-uris]",
		Short: "Broadcast message updateUris",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argUri, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}

			var argBadgeUris []*types.BadgeUri
			if err := json.Unmarshal([]byte(args[2]), &argBadgeUris); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateUris(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				argUri,
				argBadgeUris,
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
