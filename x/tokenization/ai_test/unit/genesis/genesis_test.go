package genesis_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/module"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GenesisTestSuite tests genesis export/import functionality
type GenesisTestSuite struct {
	testutil.AITestSuite
}

func TestGenesisTestSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestGenesis_ExportImportRoundtrip tests that genesis export then import produces identical state
func (suite *GenesisTestSuite) TestGenesis_ExportImportRoundtrip() {
	// Create multiple collections
	collection1Id := suite.CreateTestCollection(suite.Manager)
	collection2Id := suite.CreateTestCollection(suite.Manager)
	collection3Id := suite.CreateTestCollection(suite.Manager)

	// Setup mint approvals and mint tokens to various addresses
	suite.SetupMintApproval(collection1Id)
	suite.SetupMintApproval(collection2Id)
	suite.SetupMintApproval(collection3Id)

	suite.MintTokens(collection1Id, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collection1Id, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(50, 2)})
	suite.MintTokens(collection2Id, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})
	suite.MintTokens(collection3Id, suite.Charlie, []*types.Balance{testutil.GenerateSimpleBalance(75, 1)})

	// Get the state before export
	collection1Before := suite.GetCollection(collection1Id)
	collection2Before := suite.GetCollection(collection2Id)
	collection3Before := suite.GetCollection(collection3Id)

	aliceBalance1Before := suite.GetBalance(collection1Id, suite.Alice)
	bobBalance1Before := suite.GetBalance(collection1Id, suite.Bob)
	aliceBalance2Before := suite.GetBalance(collection2Id, suite.Alice)
	charlieBalance3Before := suite.GetBalance(collection3Id, suite.Charlie)

	nextCollectionIdBefore := suite.Keeper.GetNextCollectionId(suite.Ctx)

	// Export genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)

	// Verify export has data
	suite.Require().NotNil(genesis)
	suite.Require().True(len(genesis.Collections) >= 3, "should have at least 3 collections")
	suite.Require().True(len(genesis.Balances) > 0, "should have balances")
	suite.Require().Equal(len(genesis.Balances), len(genesis.BalanceStoreKeys), "balances and keys should match")

	// Now re-initialize with the exported genesis
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get the state after import
	collection1After := suite.GetCollection(collection1Id)
	collection2After := suite.GetCollection(collection2Id)
	collection3After := suite.GetCollection(collection3Id)

	aliceBalance1After := suite.GetBalance(collection1Id, suite.Alice)
	bobBalance1After := suite.GetBalance(collection1Id, suite.Bob)
	aliceBalance2After := suite.GetBalance(collection2Id, suite.Alice)
	charlieBalance3After := suite.GetBalance(collection3Id, suite.Charlie)

	nextCollectionIdAfter := suite.Keeper.GetNextCollectionId(suite.Ctx)

	// Verify collections are preserved
	suite.Require().Equal(collection1Before.CollectionId, collection1After.CollectionId)
	suite.Require().Equal(collection1Before.Manager, collection1After.Manager)
	suite.Require().Equal(collection2Before.CollectionId, collection2After.CollectionId)
	suite.Require().Equal(collection3Before.CollectionId, collection3After.CollectionId)

	// Verify balances are preserved
	suite.Require().Equal(len(aliceBalance1Before.Balances), len(aliceBalance1After.Balances))
	suite.Require().Equal(len(bobBalance1Before.Balances), len(bobBalance1After.Balances))
	suite.Require().Equal(len(aliceBalance2Before.Balances), len(aliceBalance2After.Balances))
	suite.Require().Equal(len(charlieBalance3Before.Balances), len(charlieBalance3After.Balances))

	// Verify amounts are preserved
	if len(aliceBalance1Before.Balances) > 0 && len(aliceBalance1After.Balances) > 0 {
		suite.Require().Equal(aliceBalance1Before.Balances[0].Amount, aliceBalance1After.Balances[0].Amount)
	}

	// Verify nextCollectionId is preserved
	suite.Require().Equal(nextCollectionIdBefore, nextCollectionIdAfter)
}

// TestGenesis_AllCollectionsPreserved tests that all collections are preserved through genesis
func (suite *GenesisTestSuite) TestGenesis_AllCollectionsPreserved() {
	// Create several collections with different configurations
	numCollections := 5
	collectionIds := make([]sdkmath.Uint, numCollections)

	for i := 0; i < numCollections; i++ {
		collectionIds[i] = suite.CreateTestCollection(suite.Manager)
	}

	// Export genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)

	// Verify all collections are in genesis
	suite.Require().True(len(genesis.Collections) >= numCollections,
		"genesis should have at least %d collections, got %d", numCollections, len(genesis.Collections))

	// Verify each collection ID exists in genesis
	for _, collectionId := range collectionIds {
		found := false
		for _, collection := range genesis.Collections {
			if collection.CollectionId.Equal(collectionId) {
				found = true
				break
			}
		}
		suite.Require().True(found, "collection %s should be in genesis", collectionId)
	}
}

// TestGenesis_AllBalancesPreserved tests that all balances are preserved through genesis
func (suite *GenesisTestSuite) TestGenesis_AllBalancesPreserved() {
	// Create collection and mint to multiple addresses
	collectionId := suite.CreateTestCollection(suite.Manager)
	suite.SetupMintApproval(collectionId)

	// Mint different amounts to different addresses
	addresses := []string{suite.Alice, suite.Bob, suite.Charlie}
	amounts := []uint64{100, 200, 300}

	for i, addr := range addresses {
		suite.MintTokens(collectionId, addr, []*types.Balance{
			testutil.GenerateSimpleBalance(amounts[i], 1),
		})
	}

	// Export genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)

	// Re-initialize
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Verify all balances
	for i, addr := range addresses {
		balance := suite.GetBalance(collectionId, addr)
		suite.Require().True(len(balance.Balances) > 0, "address %s should have balance", addr)
		suite.Require().Equal(sdkmath.NewUint(amounts[i]), balance.Balances[0].Amount,
			"address %s should have amount %d", addr, amounts[i])
	}
}

// TestGenesis_NextCollectionIdGreaterThanMax tests that nextCollectionId > max existing collectionId
func (suite *GenesisTestSuite) TestGenesis_NextCollectionIdGreaterThanMax() {
	// Create several collections
	var maxCollectionId sdkmath.Uint
	for i := 0; i < 5; i++ {
		collectionId := suite.CreateTestCollection(suite.Manager)
		if maxCollectionId.IsNil() || collectionId.GT(maxCollectionId) {
			maxCollectionId = collectionId
		}
	}

	// Export genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)

	// Verify nextCollectionId > max collection ID
	suite.Require().True(genesis.NextCollectionId.GT(maxCollectionId),
		"nextCollectionId (%s) should be greater than max existing collectionId (%s)",
		genesis.NextCollectionId, maxCollectionId)
}

// TestGenesis_CollectionMetadataPreserved tests that collection metadata is preserved
func (suite *GenesisTestSuite) TestGenesis_CollectionMetadataPreserved() {
	// Create collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Update metadata
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                  suite.Manager,
		CollectionId:             collectionId,
		UpdateCollectionMetadata: true,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/custom-metadata",
			CustomData: "custom data for testing",
		},
		UpdateCustomData: true,
		CustomData:       "collection custom data",
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Get state before
	collectionBefore := suite.GetCollection(collectionId)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get state after
	collectionAfter := suite.GetCollection(collectionId)

	// Verify metadata preserved
	suite.Require().Equal(collectionBefore.CollectionMetadata.Uri, collectionAfter.CollectionMetadata.Uri)
	suite.Require().Equal(collectionBefore.CollectionMetadata.CustomData, collectionAfter.CollectionMetadata.CustomData)
	suite.Require().Equal(collectionBefore.CustomData, collectionAfter.CustomData)
}

// TestGenesis_CollectionApprovalsPreserved tests that collection approvals are preserved
func (suite *GenesisTestSuite) TestGenesis_CollectionApprovalsPreserved() {
	// Create collection with approvals
	approval1 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	approval2 := testutil.GenerateCollectionApproval("approval2", types.MintAddress, "All")
	// When using Mint in fromListId, must set OverridesFromOutgoingApprovals = true
	approval2.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	approvals := []*types.CollectionApproval{approval1, approval2}
	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, approvals)

	// Get state before
	collectionBefore := suite.GetCollection(collectionId)
	numApprovalsBefore := len(collectionBefore.CollectionApprovals)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get state after
	collectionAfter := suite.GetCollection(collectionId)
	numApprovalsAfter := len(collectionAfter.CollectionApprovals)

	// Verify approvals preserved
	suite.Require().Equal(numApprovalsBefore, numApprovalsAfter,
		"number of approvals should be preserved")

	// Verify approval IDs preserved
	approvalIdsBefore := make(map[string]bool)
	for _, approval := range collectionBefore.CollectionApprovals {
		approvalIdsBefore[approval.ApprovalId] = true
	}

	for _, approval := range collectionAfter.CollectionApprovals {
		suite.Require().True(approvalIdsBefore[approval.ApprovalId],
			"approval ID %s should be preserved", approval.ApprovalId)
	}
}

// TestGenesis_InvariantsPreserved tests that collection invariants are preserved
func (suite *GenesisTestSuite) TestGenesis_InvariantsPreserved() {
	// Create collection with invariants
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "https://example.com/metadata",
		},
		Invariants: &types.InvariantsAddObject{
			NoCustomOwnershipTimes:      true,
			NoForcefulPostMintTransfers: true,
			MaxSupplyPerId:              sdkmath.NewUint(1000),
			DisablePoolCreation:         true,
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId

	// Get state before
	collectionBefore := suite.GetCollection(collectionId)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get state after
	collectionAfter := suite.GetCollection(collectionId)

	// Verify invariants preserved
	suite.Require().NotNil(collectionAfter.Invariants, "invariants should be preserved")
	suite.Require().Equal(collectionBefore.Invariants.NoCustomOwnershipTimes,
		collectionAfter.Invariants.NoCustomOwnershipTimes)
	suite.Require().Equal(collectionBefore.Invariants.NoForcefulPostMintTransfers,
		collectionAfter.Invariants.NoForcefulPostMintTransfers)
	suite.Require().Equal(collectionBefore.Invariants.MaxSupplyPerId,
		collectionAfter.Invariants.MaxSupplyPerId)
	suite.Require().Equal(collectionBefore.Invariants.DisablePoolCreation,
		collectionAfter.Invariants.DisablePoolCreation)
}

// TestGenesis_UserApprovalsPreserved tests that user-level approvals are preserved
func (suite *GenesisTestSuite) TestGenesis_UserApprovalsPreserved() {
	// Create collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Setup mint approval and mint tokens
	suite.SetupMintApproval(collectionId)
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Set user approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("user_outgoing", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err := suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Get state before
	balanceBefore := suite.GetBalance(collectionId, suite.Alice)
	numOutgoingBefore := len(balanceBefore.OutgoingApprovals)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get state after
	balanceAfter := suite.GetBalance(collectionId, suite.Alice)
	numOutgoingAfter := len(balanceAfter.OutgoingApprovals)

	// Verify user approvals preserved
	suite.Require().Equal(numOutgoingBefore, numOutgoingAfter,
		"number of user outgoing approvals should be preserved")
}

// TestGenesis_ValidTokenIdsPreserved tests that valid token IDs are preserved
func (suite *GenesisTestSuite) TestGenesis_ValidTokenIdsPreserved() {
	// Create collection with specific valid token IDs
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
			{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "https://example.com/metadata",
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId

	// Get state before
	collectionBefore := suite.GetCollection(collectionId)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get state after
	collectionAfter := suite.GetCollection(collectionId)

	// Verify valid token IDs preserved
	suite.Require().Equal(len(collectionBefore.ValidTokenIds), len(collectionAfter.ValidTokenIds),
		"number of valid token ID ranges should be preserved")

	for i, rangeBefore := range collectionBefore.ValidTokenIds {
		rangeAfter := collectionAfter.ValidTokenIds[i]
		suite.Require().Equal(rangeBefore.Start, rangeAfter.Start)
		suite.Require().Equal(rangeBefore.End, rangeAfter.End)
	}
}

// TestGenesis_AddressListsPreserved tests that address lists are preserved
func (suite *GenesisTestSuite) TestGenesis_AddressListsPreserved() {
	// Create address list (list ID must be alphanumeric only)
	err := suite.Keeper.CreateAddressList(suite.Ctx, &types.AddressList{
		ListId:    "testlist",
		Addresses: []string{suite.Alice, suite.Bob},
		Whitelist: true,
		CreatedBy: suite.Manager,
	})
	suite.Require().NoError(err)

	// Export genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)

	// Find our list in genesis
	found := false
	for _, list := range genesis.AddressLists {
		if list.ListId == "testlist" {
			found = true
			suite.Require().Equal(2, len(list.Addresses))
			suite.Require().True(list.Whitelist)
			break
		}
	}
	suite.Require().True(found, "address list should be in genesis")

	// Reimport
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Verify list exists
	list, listFound := suite.Keeper.GetAddressListFromStore(suite.Ctx, "testlist")
	suite.Require().True(listFound)
	suite.Require().Equal(2, len(list.Addresses))
}

// TestGenesis_ParamsPreserved tests that module params are preserved
func (suite *GenesisTestSuite) TestGenesis_ParamsPreserved() {
	// Get current params
	paramsBefore := suite.Keeper.GetParams(suite.Ctx)

	// Export and reimport genesis
	genesis := tokenization.ExportGenesis(suite.Ctx, *suite.Keeper)
	tokenization.InitGenesis(suite.Ctx, *suite.Keeper, *genesis)

	// Get params after
	paramsAfter := suite.Keeper.GetParams(suite.Ctx)

	// Verify params preserved
	suite.Require().Equal(paramsBefore, paramsAfter, "params should be preserved")
}

// TestGenesis_DefaultGenesisValid tests that default genesis is valid
func (suite *GenesisTestSuite) TestGenesis_DefaultGenesisValid() {
	defaultGenesis := types.DefaultGenesis()

	// Validate default genesis
	err := defaultGenesis.Validate()
	suite.Require().NoError(err, "default genesis should be valid")

	// Verify default values
	suite.Require().Equal(sdkmath.NewUint(1), defaultGenesis.NextCollectionId)
	suite.Require().Equal(sdkmath.NewUint(1), defaultGenesis.NextDynamicStoreId)
	suite.Require().Equal(types.PortID, defaultGenesis.PortId)
}
