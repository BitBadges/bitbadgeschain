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

func CmdClaimBadge() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "claim-badge [claim-id] [collection-id] [whitelist-proof] [code-proof]",
		Short: "Broadcast message claimBadge",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			argClaimId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argCollectionId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argWhitelistProofJson, err := parseJson(args[2])
			if err != nil {
				return err
			}

			argWhitelistProof := &types.ClaimProof{
				Leaf: argWhitelistProofJson["leaf"].(string),
				Aunts: argWhitelistProofJson["aunts"].([]*types.ClaimProofItem),
			}

			argCodeProofJson, err := parseJson(args[3])
			if err != nil {
				return err
			}

			argCodeProof := &types.ClaimProof{
				Leaf: argCodeProofJson["leaf"].(string),
				Aunts: argCodeProofJson["aunts"].([]*types.ClaimProofItem),
			}

			msg := types.NewMsgClaimBadge(
				clientCtx.GetFromAddress().String(),
				argClaimId,
				argCollectionId,
				argWhitelistProof,
				argCodeProof,
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
