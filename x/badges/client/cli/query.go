package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group badges queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGetBalance())
	cmd.AddCommand(CmdGetCollection())
	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdGetAddressMapping())
	cmd.AddCommand(CmdGetApprovalsTrackers())
	cmd.AddCommand(CmdGetNumUsedForMerkleChallenge())
	// this line is used by starport scaffolding # 1

	return cmd
}

// // Queries an approvals tracker by ID.
// rpc GetApprovalsTracker(QueryGetApprovalsTrackerRequest) returns (QueryGetApprovalsTrackerResponse) {
// 	option (google.api.http).get = "/bitbadges/bitbadgeschain/badges/get_approvals_tracker/{collectionId}/{approvalLevel}/{approverAddress}/{amountTrackerId}/{trackerType}/{approvedAddress}";
// }

// // Queries the number of times a given leaf has been used for a given merkle challenge.
// rpc GetNumUsedForMerkleChallenge(QueryGetNumUsedForMerkleChallengeRequest) returns (QueryGetNumUsedForMerkleChallengeResponse) {
// 	option (google.api.http).get = "/bitbadges/bitbadgeschain/badges/get_num_used_for_challenge/{collectionId}/{approvalLevel}/{approverAddress}/{challengeTrackerId}/{leafIndex}";
// }
