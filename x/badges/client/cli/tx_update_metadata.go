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

func CmdUpdateMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-uris [collection-id] [uri] [badge-uris] [offChainBalancesMetadata]",
		Short: "Broadcast message updateMetadata",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBadgeId := types.NewUintFromString(args[0])
			if err != nil {
				return err
			}

			argUri, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}

			var argBadgeMetadata []*types.BadgeMetadata
			if err := json.Unmarshal([]byte(args[2]), &argBadgeMetadata); err != nil {
				return err
			}

			argOffChainBalancesMetadata, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateMetadata(
				clientCtx.GetFromAddress().String(),
				argBadgeId,
				argUri,
				argBadgeMetadata,
				argOffChainBalancesMetadata,
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
