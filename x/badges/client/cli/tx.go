package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands. For nested types, use JSON format with quotes (e.g. '{\"key\": \"value\"}') or '[...]'", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdTransferBadges())
	cmd.AddCommand(CmdDeleteCollection())
	cmd.AddCommand(CmdUpdateUserApprovedOutgoingTransfers())
	cmd.AddCommand(CmdUpdateCollection())
	cmd.AddCommand(CmdCreateAddressMappings())
	// this line is used by starport scaffolding # 1

	return cmd
}
