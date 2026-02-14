package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands. For nested types, use JSON format with quotes (e.g. '{\"key\": \"value\"}') or '[...]'", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdTransferTokens())
	cmd.AddCommand(CmdDeleteCollection())
	cmd.AddCommand(CmdUpdateUserOutgoingApprovals())
	cmd.AddCommand(CmdUpdateCollection())
	cmd.AddCommand(CmdCreateAddressLists())
	cmd.AddCommand(CmdMsgCreateCollection())
	cmd.AddCommand(CmdMsgUpdateCollection())
	cmd.AddCommand(CmdCreateDynamicStore())
	cmd.AddCommand(CmdUpdateDynamicStore())
	cmd.AddCommand(CmdDeleteDynamicStore())
	cmd.AddCommand(CmdSetDynamicStoreValue())
	cmd.AddCommand(CmdSetIncomingApproval())
	cmd.AddCommand(CmdDeleteIncomingApproval())
	cmd.AddCommand(CmdSetOutgoingApproval())
	cmd.AddCommand(CmdDeleteOutgoingApproval())
	cmd.AddCommand(CmdPurgeApprovals())
	cmd.AddCommand(CmdSetValidTokenIds())
	cmd.AddCommand(CmdSetManager())
	cmd.AddCommand(CmdSetCollectionMetadata())
	cmd.AddCommand(CmdSetTokenMetadata())
	cmd.AddCommand(CmdSetCustomData())
	cmd.AddCommand(CmdSetStandards())
	cmd.AddCommand(CmdSetCollectionApprovals())
	cmd.AddCommand(CmdSetIsArchived())
	cmd.AddCommand(CmdCastVote())
	// this line is used by starport scaffolding # 1

	return cmd
}
