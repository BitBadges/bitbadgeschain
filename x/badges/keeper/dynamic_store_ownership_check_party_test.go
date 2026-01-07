package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	approvalcriteria "github.com/bitbadges/bitbadgeschain/x/badges/approval_criteria"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestDynamicStoreOwnershipCheckParty_Initiator tests that ownership check defaults to initiator
func TestDynamicStoreOwnershipCheckParty_Initiator(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	sender := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	recipient := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme" // Use same as sender for simplicity

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set value for initiator to true
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: initiator,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that uses this dynamic store (default/empty ownershipCheckParty)
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:            createResp.StoreId,
							OwnershipCheckParty: "", // Empty should default to initiator
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Should pass because initiator has value = true
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when checking initiator and initiator has permission")

	// Explicitly set to "initiator"
	collection.CollectionApprovals[0].ApprovalCriteria.DynamicStoreChallenges[0].OwnershipCheckParty = "initiator"
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	// Should still pass
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when explicitly checking initiator and initiator has permission")
}

// TestDynamicStoreOwnershipCheckParty_Sender tests that ownership check can use sender
func TestDynamicStoreOwnershipCheckParty_Sender(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	sender := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	recipient := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme" // Use same as sender for simplicity

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set value for sender to true (not initiator)
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: sender,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that checks sender
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:            createResp.StoreId,
							OwnershipCheckParty: "sender",
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Should pass because sender has value = true
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when checking sender and sender has permission")

	// Should fail if sender doesn't have permission
	setValueMsg.Value = false
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.Error(t, err, "should fail when sender doesn't have permission")
}

// TestDynamicStoreOwnershipCheckParty_Recipient tests that ownership check can use recipient
func TestDynamicStoreOwnershipCheckParty_Recipient(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	sender := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	recipient := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme" // Use same as sender for simplicity

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set value for recipient to true
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: recipient,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that checks recipient
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:            createResp.StoreId,
							OwnershipCheckParty: "recipient",
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Should pass because recipient has value = true
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when checking recipient and recipient has permission")

	// Should fail if recipient doesn't have permission
	setValueMsg.Value = false
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.Error(t, err, "should fail when recipient doesn't have permission")
}

// TestDynamicStoreOwnershipCheckParty_HardcodedAddress tests that ownership check can use a hardcoded bb1 address
func TestDynamicStoreOwnershipCheckParty_HardcodedAddress(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	sender := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	recipient := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme" // Use same as sender for simplicity
	hardcodedAddress := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q" // Use creator as hardcoded address

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set value for hardcoded address to true
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: hardcodedAddress,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that checks hardcoded address
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:            createResp.StoreId,
							OwnershipCheckParty: hardcodedAddress,
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Should pass because hardcoded address has value = true
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when checking hardcoded address and address has permission")

	// Should fail if hardcoded address doesn't have permission
	setValueMsg.Value = false
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.Error(t, err, "should fail when hardcoded address doesn't have permission")
}

// TestDynamicStoreOwnershipCheckParty_InvalidAddressFallsBackToInitiator tests that invalid addresses fall back to initiator
func TestDynamicStoreOwnershipCheckParty_InvalidAddressFallsBackToInitiator(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	sender := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	recipient := "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme" // Use same as sender for simplicity

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set value for initiator to true (fallback target)
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: initiator,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that uses an invalid address (should fall back to initiator)
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:            createResp.StoreId,
							OwnershipCheckParty: "invalid-address-format",
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Should pass because invalid address falls back to initiator, and initiator has permission
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, recipient, sender, initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when invalid address falls back to initiator and initiator has permission")
}

