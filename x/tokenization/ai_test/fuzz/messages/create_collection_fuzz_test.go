package messages_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// Note: Fuzz tests use different structure than regular test suites

// FuzzCreateCollection_ValidInputs fuzzes CreateCollection with valid inputs
func FuzzCreateCollection_ValidInputs(f *testing.F) {
	// Seed with valid inputs - use valid bech32 addresses
	f.Add("bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx", uint64(1), uint64(100), "https://example.com/metadata")
	f.Add("bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx", uint64(1), uint64(1000), "")
	f.Add("bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx", uint64(100), uint64(200), "https://test.com")

	f.Fuzz(func(t *testing.T, creator string, tokenIdStart, tokenIdEnd uint64, uri string) {
		// Skip invalid inputs
		if tokenIdStart > tokenIdEnd {
			return
		}
		if tokenIdStart == 0 {
			return
		}

		// Setup
		k, ctx := keepertest.TokenizationKeeper(t)
		msgServer := keeper.NewMsgServerImpl(k)

		// Create message with fuzzed inputs
		msg := &types.MsgCreateCollection{
			Creator: creator,
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
			},
			ValidTokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(tokenIdStart), End: sdkmath.NewUint(tokenIdEnd)},
			},
			CollectionPermissions: &types.CollectionPermissions{},
			Manager: creator,
			CollectionMetadata: &types.CollectionMetadata{
				Uri:        uri,
				CustomData: "",
			},
			TokenMetadata: []*types.TokenMetadata{},
			CustomData:    "",
			CollectionApprovals: []*types.CollectionApproval{},
			Standards:      []string{},
			IsArchived:    false,
		}

		// Attempt to create collection - should either succeed or fail with validation error
		_, err := msgServer.CreateCollection(sdk.WrapSDKContext(ctx), msg)
		if err != nil {
			// Validation errors are expected for invalid inputs
			// We just want to ensure no panics occur
			return
		}
	})
}

// TestCreateCollectionFuzz_RandomInputs tests CreateCollection with random inputs
func TestCreateCollectionFuzz_RandomInputs(t *testing.T) {
	rand.Seed(42) // Fixed seed for reproducibility

	k, ctx := keepertest.TokenizationKeeper(t)
	msgServer := keeper.NewMsgServerImpl(k)
	manager := "bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx"

	for i := 0; i < 100; i++ {
		// Generate random valid inputs
		tokenIdStart := uint64(rand.Intn(1000) + 1)
		tokenIdEnd := tokenIdStart + uint64(rand.Intn(1000))

		msg := &types.MsgCreateCollection{
			Creator: manager,
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
			},
			ValidTokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(tokenIdStart), End: sdkmath.NewUint(tokenIdEnd)},
			},
			CollectionPermissions: &types.CollectionPermissions{},
			Manager: manager,
			CollectionMetadata: &types.CollectionMetadata{
				Uri:        "",
				CustomData: "",
			},
			TokenMetadata: []*types.TokenMetadata{},
			CustomData: "",
			CollectionApprovals: []*types.CollectionApproval{},
			Standards: []string{},
			IsArchived: false,
		}

		_, err := msgServer.CreateCollection(sdk.WrapSDKContext(ctx), msg)
		if err != nil {
			// Validation errors are acceptable in fuzz testing
			// We just want to ensure no panics occur
			return
		}
	}
}

