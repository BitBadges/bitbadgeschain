package cli

import (
	"encoding/json"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func CmdPurgeApprovals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge-approvals [collection-id] [purge-expired] [approver-address] [purge-counterparty-approvals] [approvals-to-purge-json-or-file]",
		Short: "Broadcast message PurgeApprovals",
		Long: `Broadcast message PurgeApprovals with specific approvals to purge.

The approvals-to-purge-json-or-file parameter is REQUIRED and must be a JSON array of approval identifier details.
Accepts JSON either inline or from a file path. If the argument is a valid file path, it will read the JSON from that file.
Example inline: '[{"approvalId":"approval1","approvalLevel":"collection","approverAddress":"","version":"1"},{"approvalId":"approval2","approvalLevel":"incoming","approverAddress":"cosmos1...","version":"1"}]'

Rules:
- For self-purge (creator purging own approvals): purgeExpired must be true, purgeCounterpartyApprovals must be false
- For other-purge (creator purging someone else's approvals): can set either purgeExpired or purgeCounterpartyApprovals
- approvalsToPurge cannot be empty - you must specify exactly which approvals to purge`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCollectionId := args[0]
			argPurgeExpired := args[1]
			argApproverAddress := args[2]
			argPurgeCounterpartyApprovals := args[3]
			argApprovalsToPurgeJson := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId, err := strconv.ParseUint(argCollectionId, 10, 64)
			if err != nil {
				return err
			}

			purgeExpired, err := strconv.ParseBool(argPurgeExpired)
			if err != nil {
				return err
			}

			purgeCounterpartyApprovals, err := strconv.ParseBool(argPurgeCounterpartyApprovals)
			if err != nil {
				return err
			}

			// Parse the JSON for specific approvals to purge (now required)
			var approvalsToPurge []*types.ApprovalIdentifierDetails
			if argApprovalsToPurgeJson == "" {
				return errorsmod.Wrap(types.ErrApprovalsToPurgeEmpty, "")
			}

			// Support file or inline JSON
			jsonBytes, err := ReadJSONBytesFromFileOrString(argApprovalsToPurgeJson)
			if err != nil {
				return err
			}

			err = json.Unmarshal(jsonBytes, &approvalsToPurge)
			if err != nil {
				return err
			}

			if len(approvalsToPurge) == 0 {
				return errorsmod.Wrap(types.ErrApprovalsToPurgeEmpty, "")
			}

			msg := types.NewMsgPurgeApprovals(
				clientCtx.GetFromAddress().String(),
				sdkmath.NewUint(collectionId),
				purgeExpired,
				argApproverAddress,
				purgeCounterpartyApprovals,
				approvalsToPurge,
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
