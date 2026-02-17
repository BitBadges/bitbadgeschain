package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type PoolIntegrationTestSuite struct {
	TestSuite
}

func TestPoolIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PoolIntegrationTestSuite))
}

func (suite *PoolIntegrationTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

// Helper to create a collection with alias paths for pool integration testing
func (suite *PoolIntegrationTestSuite) createCollectionWithAliasPath(creator string, denom string) (*types.TokenCollection, string, error) {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(creator)
	collectionsToCreate[0].AliasPathsToAdd = []*types.AliasPathAddObject{
		{
			Denom: denom,
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(1),
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
			Symbol:     "TEST",
			DenomUnits: []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: denom, IsDefaultDisplay: true}},
		},
	}

	// Add collection approvals for transfers to/from wrapper paths
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	})

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	if err != nil {
		return nil, "", err
	}

	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	if err != nil {
		return nil, "", err
	}

	// Generate the wrapper denom (format: badgeslp:collectionId:denom)
	wrapperDenom := keeper.AliasDenomPrefix + collection.CollectionId.String() + ":" + denom

	return collection, wrapperDenom, nil
}

// =============================================================================
// GetBalancesToTransferWithAlias Edge Case Tests
// =============================================================================

// TestGetBalancesToTransferWithAlias_NonDivisibleAmount tests non-evenly-divisible amount error
func (suite *PoolIntegrationTestSuite) TestGetBalancesToTransferWithAlias_NonDivisibleAmount() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection where conversion requires amount of 10
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].AliasPathsToAdd = []*types.AliasPathAddObject{
		{
			Denom: "divisibletest",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(10), // Requires amounts divisible by 10
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
			Symbol:     "DIV",
			DenomUnits: []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "divisibletest", IsDefaultDisplay: true}},
		},
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	})

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	wrapperDenom := keeper.AliasDenomPrefix + collection.CollectionId.String() + ":divisibletest"

	// Try to get balances for non-divisible amount (should fail)
	_, err = keeper.GetBalancesToTransferWithAlias(collection, wrapperDenom, sdkmath.NewUint(7))
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not evenly divisible")

	// Divisible amount should work
	balances, err := keeper.GetBalancesToTransferWithAlias(collection, wrapperDenom, sdkmath.NewUint(10))
	suite.Require().NoError(err)
	suite.Require().Len(balances, 1)
	suite.Require().Equal(sdkmath.NewUint(1), balances[0].Amount) // 10 / 10 = 1
}

// =============================================================================
// SendNativeTokensViaAliasDenom Tests
// =============================================================================

// TestSendNativeTokensViaAliasDenom_Success tests successful alias denom routing
func (suite *PoolIntegrationTestSuite) TestSendNativeTokensViaAliasDenom_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collection, wrapperDenom, err := suite.createCollectionWithAliasPath(bob, "aliastest")
	suite.Require().NoError(err)

	// Mint tokens to bob
	err = TransferTokens(&suite.TestSuite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err)

	// Execute via alias denom
	err = suite.app.TokenizationKeeper.SendNativeTokensViaAliasDenom(
		suite.ctx,
		bob,
		alice,
		wrapperDenom,
		sdkmath.NewUint(1),
	)
	suite.Require().NoError(err)

	// Verify alice received tokens
	aliceBalance, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, alice)
	suite.Require().NoError(err)
	suite.Require().True(len(aliceBalance.Balances) > 0, "Alice should have received tokens")
}

// =============================================================================
// FundCommunityPoolViaAliasDenom Tests
// =============================================================================

// TestFundCommunityPoolViaAliasDenom_Success tests community pool funding via alias denom
func (suite *PoolIntegrationTestSuite) TestFundCommunityPoolViaAliasDenom_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collection, wrapperDenom, err := suite.createCollectionWithAliasPath(bob, "communitytest")
	suite.Require().NoError(err)

	// Mint tokens to bob
	err = TransferTokens(&suite.TestSuite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err)

	// Use signer as a "community pool" address for testing (valid bech32 address)
	communityPoolAddr := signer

	// Fund the community pool via alias denom
	err = suite.app.TokenizationKeeper.FundCommunityPoolViaAliasDenom(
		suite.ctx,
		bob,
		communityPoolAddr,
		wrapperDenom,
		sdkmath.NewUint(1),
	)
	suite.Require().NoError(err)

	// Verify community pool received tokens
	poolBalance, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, communityPoolAddr)
	suite.Require().NoError(err)
	suite.Require().True(len(poolBalance.Balances) > 0, "Community pool should have received tokens")
}

// =============================================================================
// Pool Auto-Approval Tests
// =============================================================================

// TestSetAllAutoApprovalFlagsForPoolAddress_Success tests that pool addresses get auto-approval flags
func (suite *PoolIntegrationTestSuite) TestSetAllAutoApprovalFlagsForPoolAddress_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collection, _, err := suite.createCollectionWithAliasPath(bob, "autoapptest")
	suite.Require().NoError(err)

	// Use charlie as a "pool address" for testing (valid bech32 address)
	poolAddr := charlie

	// Set auto-approval flags for pool address
	err = suite.app.TokenizationKeeper.SetAllAutoApprovalFlagsForPoolAddress(
		suite.ctx,
		collection,
		poolAddr,
	)
	suite.Require().NoError(err)

	// Verify the flags were set
	poolBalance, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, poolAddr)
	suite.Require().NoError(err)
	suite.Require().True(poolBalance.AutoApproveAllIncomingTransfers, "AutoApproveAllIncomingTransfers should be true")
	suite.Require().True(poolBalance.AutoApproveSelfInitiatedOutgoingTransfers, "AutoApproveSelfInitiatedOutgoingTransfers should be true")
	suite.Require().True(poolBalance.AutoApproveSelfInitiatedIncomingTransfers, "AutoApproveSelfInitiatedIncomingTransfers should be true")
}

// TestSetAllAutoApprovalFlagsForPoolAddress_DoesNotOverride tests that existing flags are not overridden
func (suite *PoolIntegrationTestSuite) TestSetAllAutoApprovalFlagsForPoolAddress_DoesNotOverride() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collection, _, err := suite.createCollectionWithAliasPath(bob, "nooverridetest")
	suite.Require().NoError(err)

	// Use signer as a "pool address" for testing (valid bech32 address)
	poolAddr := signer

	// First call - sets all flags
	err = suite.app.TokenizationKeeper.SetAllAutoApprovalFlagsForPoolAddress(
		suite.ctx,
		collection,
		poolAddr,
	)
	suite.Require().NoError(err)

	// Second call - should not error (idempotent)
	err = suite.app.TokenizationKeeper.SetAllAutoApprovalFlagsForPoolAddress(
		suite.ctx,
		collection,
		poolAddr,
	)
	suite.Require().NoError(err)

	// Verify the flags are still set
	poolBalance, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, poolAddr)
	suite.Require().NoError(err)
	suite.Require().True(poolBalance.AutoApproveAllIncomingTransfers)
}

// =============================================================================
// GetBalancesToTransferWithAlias Tests
// =============================================================================

// TestGetBalancesToTransferWithAlias_Success tests balance calculation for alias transfer
func (suite *PoolIntegrationTestSuite) TestGetBalancesToTransferWithAlias_Success() {
	collection, wrapperDenom, err := suite.createCollectionWithAliasPath(bob, "balancetest")
	suite.Require().NoError(err)

	// Get balances for amount 5
	balances, err := keeper.GetBalancesToTransferWithAlias(collection, wrapperDenom, sdkmath.NewUint(5))
	suite.Require().NoError(err)
	suite.Require().NotNil(balances)
	suite.Require().Len(balances, 1)
	suite.Require().Equal(sdkmath.NewUint(5), balances[0].Amount)
}

// TestGetBalancesToTransferWithAlias_ZeroConversion tests error on zero conversion amount
func (suite *PoolIntegrationTestSuite) TestGetBalancesToTransferWithAlias_InvalidDenom() {
	collection, _, err := suite.createCollectionWithAliasPath(bob, "invalidtest")
	suite.Require().NoError(err)

	// Try with invalid denom
	_, err = keeper.GetBalancesToTransferWithAlias(collection, "invalid-denom", sdkmath.NewUint(1))
	suite.Require().Error(err)
}

