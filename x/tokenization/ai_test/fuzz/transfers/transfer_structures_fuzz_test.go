package transfers_test

import (
	"math"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

// FuzzTransfer_BalanceRanges fuzzes transfer balance ranges
func FuzzTransfer_BalanceRanges(f *testing.F) {
	// Seed with valid inputs
	f.Add(uint64(1), uint64(100), uint64(1), uint64(100), uint64(1), uint64(math.MaxUint64))
	f.Add(uint64(10), uint64(50), uint64(1), uint64(10), uint64(1), uint64(math.MaxUint64))
	f.Add(uint64(1), uint64(1), uint64(1), uint64(1), uint64(1), uint64(math.MaxUint64))

	f.Fuzz(func(t *testing.T, amount, tokenIdStart, tokenIdEnd, ownershipTimeStart, ownershipTimeEnd uint64, maxOwnershipTime uint64) {
		// Skip invalid inputs
		if tokenIdStart > tokenIdEnd {
			return
		}
		if ownershipTimeStart > ownershipTimeEnd {
			return
		}
		if ownershipTimeEnd > maxOwnershipTime {
			return
		}
		if amount == 0 {
			return
		}

		// Setup
		k, ctx := keepertest.TokenizationKeeper(t)
		msgServer := keeper.NewMsgServerImpl(k)

		// Create collection
		manager := "bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx"
		createMsg := &types.MsgCreateCollection{
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
			CustomData:    "",
			CollectionApprovals: []*types.CollectionApproval{},
			Standards:      []string{},
			IsArchived:    false,
		}

		resp, err := msgServer.CreateCollection(sdk.WrapSDKContext(ctx), createMsg)
		if err != nil {
			// Validation errors are expected for invalid inputs
			return
		}

		// Create transfer with fuzzed balance ranges
		transferBalance := &types.Balance{
			Amount: sdkmath.NewUint(amount),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(tokenIdStart), End: sdkmath.NewUint(tokenIdEnd)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(ownershipTimeStart), End: sdkmath.NewUint(ownershipTimeEnd)},
			},
		}

		transferMsg := &types.MsgTransferTokens{
			Creator:      manager,
			CollectionId: resp.CollectionId,
			Transfers: []*types.Transfer{
				{
					From:        types.MintAddress,
					ToAddresses: []string{"bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"},
					Balances:    []*types.Balance{transferBalance},
				},
			},
		}

		// Attempt transfer - should either succeed or fail with validation error
		_, err = msgServer.TransferTokens(sdk.WrapSDKContext(ctx), transferMsg)
		if err != nil {
			// Validation errors are expected for invalid inputs
			return
		}
	})
}

// TestTransferFuzz_RandomBalances tests transfers with random balance structures
func TestTransferFuzz_RandomBalances(t *testing.T) {
	rand.Seed(42) // Fixed seed for reproducibility

	k, ctx := keepertest.TokenizationKeeper(t)
	msgServer := keeper.NewMsgServerImpl(k)
	manager := "bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx"
	alice := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create collection
	createMsg := &types.MsgCreateCollection{
		Creator: manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "",
			CustomData: "",
		},
		TokenMetadata: []*types.TokenMetadata{},
		CustomData:    "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:      []string{},
		IsArchived:    false,
	}

	resp, err := msgServer.CreateCollection(sdk.WrapSDKContext(ctx), createMsg)
	require.NoError(t, err)

	// Test multiple random balance structures
	for i := 0; i < 50; i++ {
		amount := uint64(rand.Intn(100) + 1)
		tokenId := uint64(rand.Intn(100) + 1)

		transferBalance := &types.Balance{
			Amount: sdkmath.NewUint(amount),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(tokenId), End: sdkmath.NewUint(tokenId)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
		}

		transferMsg := &types.MsgTransferTokens{
			Creator:      manager,
			CollectionId: resp.CollectionId,
			Transfers: []*types.Transfer{
				{
					From:        types.MintAddress,
					ToAddresses: []string{alice},
					Balances:    []*types.Balance{transferBalance},
				},
			},
		}

		_, err = msgServer.TransferTokens(sdk.WrapSDKContext(ctx), transferMsg)
		// Some transfers may fail due to missing approvals, which is expected
		_ = err
	}
}

