package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdNewCollection())
	cmd.AddCommand(CmdMintBadge())
	cmd.AddCommand(CmdTransferBadge())
	cmd.AddCommand(CmdSetApproval())
	cmd.AddCommand(CmdUpdateDisallowedTransfers())
	cmd.AddCommand(CmdUpdateUris())
	cmd.AddCommand(CmdUpdatePermissions())
	cmd.AddCommand(CmdTransferManager())
	cmd.AddCommand(CmdRequestTransferManager())
	cmd.AddCommand(CmdUpdateBytes())
	cmd.AddCommand(CmdRegisterAddresses())
	// this line is used by starport scaffolding # 1

	return cmd
}
