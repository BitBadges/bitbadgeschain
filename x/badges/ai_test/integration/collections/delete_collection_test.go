package collections_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type DeleteCollectionTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestDeleteCollectionSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(DeleteCollectionTestSuite))
}

func (suite *DeleteCollectionTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestDeleteCollection_PurgesAllState tests that deleting a collection purges all related state
func (suite *DeleteCollectionTestSuite) TestDeleteCollection_PurgesAllState() {
	// Setup: Create collection with balances, approvals, and challenges
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Add collection approvals - need mint approval first
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:               collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval, approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens to create balances
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(collectionId, suite.Alice, mintBalances)
	suite.MintBadges(collectionId, suite.Bob, mintBalances)

	// Set user approvals to create approval trackers
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Verify balances exist before deletion
	aliceBalanceBefore := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().NotNil(aliceBalanceBefore)
	suite.Require().Greater(len(aliceBalanceBefore.Balances), 0, "Alice should have balances before deletion")

	bobBalanceBefore := suite.GetBalance(collectionId, suite.Bob)
	suite.Require().NotNil(bobBalanceBefore)
	suite.Require().Greater(len(bobBalanceBefore.Balances), 0, "Bob should have balances before deletion")

	// Delete the collection
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}
	_, err = suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().NoError(err, "collection deletion should succeed")

	// Verify collection is deleted
	_, found := suite.Keeper.GetCollectionFromStore(suite.Ctx, collectionId)
	suite.Require().False(found, "collection should not exist after deletion")

	// Verify balances are purged - check that no balances remain for this collection
	allBalances, _, allCollectionIds := suite.Keeper.GetUserBalancesFromStore(suite.Ctx)
	for i, id := range allCollectionIds {
		suite.Require().NotEqual(collectionId, id, "No balances should remain for deleted collection. Found balance at index %d", i)
		_ = allBalances[i] // Use the balance to avoid unused variable
	}
}

// TestDeleteCollection_NoOrphanedState tests that no orphaned state remains after deletion
func (suite *DeleteCollectionTestSuite) TestDeleteCollection_NoOrphanedState() {
	// Create and populate a collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Set up mint approval first
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
	updateMintMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:               collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMintMsg)
	suite.Require().NoError(err)

	// Add multiple balances
	for i := 0; i < 5; i++ {
		mintBalances := []*types.Balance{
			testutil.GenerateSimpleBalance(10, uint64(i+1)),
		}
		suite.MintBadges(collectionId, suite.Alice, mintBalances)
	}

	// Add approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Delete collection
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}
	_, err = suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().NoError(err)

	// Verify no balances remain for this collection
	// Iterate all balances and ensure none belong to deleted collection
	allBalances, allAddresses, allCollectionIds := suite.Keeper.GetUserBalancesFromStore(suite.Ctx)
	for i, id := range allCollectionIds {
		if id.Equal(collectionId) {
			suite.T().Errorf("Found orphaned balance for deleted collection %s at address %s", collectionId, allAddresses[i])
		}
		_ = allBalances[i] // Use the balance to avoid unused variable
	}
}

