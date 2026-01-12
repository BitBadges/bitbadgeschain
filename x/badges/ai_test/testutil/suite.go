package testutil

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/badges/testutil/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AITestSuite provides a comprehensive test suite for AI-generated tests
type AITestSuite struct {
	suite.Suite

	Keeper      keeper.Keeper
	Ctx         sdk.Context
	MsgServer   types.MsgServer
	QueryClient types.QueryClient

	// Test addresses
	Alice   string
	Bob     string
	Charlie string
	Manager string
	Creator string

	// Test collection ID (will be set after collection creation)
	CollectionId sdkmath.Uint
}

// SetupTest initializes the test suite with a fresh keeper and context
func (suite *AITestSuite) SetupTest() {
	k, ctx := keepertest.BadgesKeeper(suite.T())

	suite.Keeper = k
	suite.Ctx = ctx.WithBlockTime(time.Now())
	suite.MsgServer = keeper.NewMsgServerImpl(k)

	// Initialize test addresses - using valid Bech32 addresses with "bb" prefix
	// These addresses are from existing tests and are known to be valid
	suite.Alice = "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	suite.Bob = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	suite.Charlie = "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	suite.Manager = "bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx"
	suite.Creator = suite.Manager // Default creator is manager

	suite.CollectionId = sdkmath.NewUint(0) // Will be set after collection creation
}

// CreateTestCollection creates a basic test collection with default settings
func (suite *AITestSuite) CreateTestCollection(creator string) sdkmath.Uint {
	msg := &types.MsgCreateCollection{
		Creator: creator,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               creator,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation should succeed")
	suite.Require().NotNil(resp, "collection creation response should not be nil")
	suite.Require().True(resp.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0, got: %s", resp.CollectionId.String())

	collectionId := resp.CollectionId
	suite.CollectionId = collectionId
	return collectionId
}

// CreateTestCollectionWithApprovals creates a test collection with approval settings
func (suite *AITestSuite) CreateTestCollectionWithApprovals(creator string, approvals []*types.CollectionApproval) sdkmath.Uint {
	msg := &types.MsgCreateCollection{
		Creator: creator,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               creator,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: approvals,
		Standards:           []string{},
		IsArchived:          false,
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation should succeed")
	suite.Require().NotNil(resp, "collection creation response should not be nil")
	suite.Require().True(resp.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0, got: %s", resp.CollectionId.String())

	collectionId := resp.CollectionId
	suite.CollectionId = collectionId
	return collectionId
}

// SetupMintApproval sets up a collection approval that allows minting from Mint address
// Preserves existing approvals by merging the mint approval with current approvals
// According to BitBadges docs: Mint approvals must have fromListId: 'Mint' and must forcefully override user-level outgoing approval
func (suite *AITestSuite) SetupMintApproval(collectionId sdkmath.Uint) {
	// Get current collection to check existing approvals
	collection := suite.GetCollection(collectionId)

	mintApprovalId := "mint_approval"

	// Check if mint approval already exists and has correct flags
	hasMintApproval := false
	needsUpdate := false
	var existingMintApproval *types.CollectionApproval
	for _, approval := range collection.CollectionApprovals {
		if approval.ApprovalId == mintApprovalId {
			hasMintApproval = true
			existingMintApproval = approval
			// Check if it has the required override flags
			if approval.ApprovalCriteria == nil ||
				!approval.ApprovalCriteria.OverridesFromOutgoingApprovals ||
				approval.FromListId != types.MintAddress {
				needsUpdate = true
			}
			break
		}
	}

	// Create or update mint approval with fromListId: 'Mint' as required by BitBadges docs
	var mintApproval *types.CollectionApproval
	if hasMintApproval && !needsUpdate {
		// Approval exists and has correct flags, no need to update
		return
	} else if hasMintApproval && needsUpdate {
		// Update existing approval to ensure it has correct flags
		mintApproval = existingMintApproval
		if mintApproval.ApprovalCriteria == nil {
			mintApproval.ApprovalCriteria = &types.ApprovalCriteria{}
		}
		mintApproval.FromListId = types.MintAddress
		mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
		mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
	} else {
		// Create new mint approval
		mintApproval = GenerateCollectionApproval(mintApprovalId, types.MintAddress, "All")
		// Mint approvals must forcefully override user-level outgoing approval because Mint cannot be managed
		mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true // Forcefully override outgoing approvals
		mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true   // Allow mint to bypass incoming approvals
	}

	// Merge with existing approvals (prepend to ensure it's checked first due to first-match policy)
	newApprovals := make([]*types.CollectionApproval, 0, len(collection.CollectionApprovals)+1)
	if hasMintApproval {
		// Replace existing approval
		for _, approval := range collection.CollectionApprovals {
			if approval.ApprovalId != mintApprovalId {
				newApprovals = append(newApprovals, approval)
			}
		}
		newApprovals = append([]*types.CollectionApproval{mintApproval}, newApprovals...)
	} else {
		// Add new approval at the beginning
		newApprovals = append(newApprovals, mintApproval)
		newApprovals = append(newApprovals, collection.CollectionApprovals...)
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       newApprovals,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "failed to set up mint approval")

	// Verify the approval was saved by refreshing the collection
	updatedCollection := suite.GetCollection(collectionId)
	found := false
	for _, approval := range updatedCollection.CollectionApprovals {
		if approval.ApprovalId == mintApprovalId {
			found = true
			suite.Require().Equal(types.MintAddress, approval.FromListId, "mint approval FromListId should be 'Mint'")
			suite.Require().NotNil(approval.ApprovalCriteria, "mint approval must have ApprovalCriteria")
			suite.Require().True(approval.ApprovalCriteria.OverridesFromOutgoingApprovals, "mint approval must override outgoing approvals")
			break
		}
	}
	suite.Require().True(found, "mint approval should be saved in collection")
}

// MintBadges mints badges to an address
// Automatically sets up mint approval if needed
func (suite *AITestSuite) MintBadges(collectionId sdkmath.Uint, to string, balances []*types.Balance) {
	// Ensure minting is allowed by setting up a collection approval for Mint
	suite.SetupMintApproval(collectionId)

	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{to},
				Balances:    balances,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
}

// AdvanceBlockTime advances the block time by the specified duration
func (suite *AITestSuite) AdvanceBlockTime(duration time.Duration) {
	newTime := suite.Ctx.BlockTime().Add(duration)
	suite.Ctx = suite.Ctx.WithBlockTime(newTime)
}

// SetBlockTime sets the block time to a specific time
func (suite *AITestSuite) SetBlockTime(t time.Time) {
	suite.Ctx = suite.Ctx.WithBlockTime(t)
}

// GetCollection retrieves a collection by ID
func (suite *AITestSuite) GetCollection(collectionId sdkmath.Uint) *types.TokenCollection {
	collection, found := suite.Keeper.GetCollectionFromStore(suite.Ctx, collectionId)
	suite.Require().True(found, "collection should exist")
	return collection
}

// GetBalance retrieves a user's balance for a collection
func (suite *AITestSuite) GetBalance(collectionId sdkmath.Uint, address string) *types.UserBalanceStore {
	collection := suite.GetCollection(collectionId)
	balance, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, address)
	return balance
}

// AssertBalance asserts that a user's balance matches expected values
func (suite *AITestSuite) AssertBalance(collectionId sdkmath.Uint, address string, expectedBalances []*types.Balance) {
	balance := suite.GetBalance(collectionId, address)
	suite.Require().Equal(len(expectedBalances), len(balance.Balances), "balance count should match")

	// Simple comparison - in real tests, you might want more sophisticated balance comparison
	for i, expected := range expectedBalances {
		if i < len(balance.Balances) {
			actual := balance.Balances[i]
			suite.Require().True(expected.Amount.Equal(actual.Amount), "amounts should match")
			suite.Require().Equal(len(expected.TokenIds), len(actual.TokenIds), "token ID ranges should match")
			suite.Require().Equal(len(expected.OwnershipTimes), len(actual.OwnershipTimes), "ownership time ranges should match")
		}
	}
}

// RunTestSuite runs a test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
