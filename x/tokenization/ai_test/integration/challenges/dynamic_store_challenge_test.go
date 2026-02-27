package challenges_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type DynamicStoreChallengeTestSuite struct {
	testutil.AITestSuite
	StoreId sdkmath.Uint
}

func TestDynamicStoreChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(DynamicStoreChallengeTestSuite))
}

func (suite *DynamicStoreChallengeTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.StoreId = sdkmath.NewUint(0)
}

// createDynamicStore creates a new dynamic store and returns its ID
func (suite *DynamicStoreChallengeTestSuite) createDynamicStore(creator string, defaultValue bool) sdkmath.Uint {
	msg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: defaultValue,
		Uri:          "",
		CustomData:   "",
	}

	resp, err := suite.MsgServer.CreateDynamicStore(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "creating dynamic store should succeed")
	suite.Require().NotNil(resp, "response should not be nil")

	return resp.StoreId
}

// setDynamicStoreValue sets a value for an address in the dynamic store
func (suite *DynamicStoreChallengeTestSuite) setDynamicStoreValue(storeId sdkmath.Uint, creator, address string, value bool) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: storeId,
		Address: address,
		Value:   value,
	}

	_, err := suite.MsgServer.SetDynamicStoreValue(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting dynamic store value should succeed")
}

// updateDynamicStore updates the dynamic store configuration
func (suite *DynamicStoreChallengeTestSuite) updateDynamicStore(storeId sdkmath.Uint, creator string, defaultValue, globalEnabled bool) {
	msg := &types.MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:       storeId,
		DefaultValue:  defaultValue,
		GlobalEnabled: globalEnabled,
	}

	_, err := suite.MsgServer.UpdateDynamicStore(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating dynamic store should succeed")
}

// TestDynamicStoreChallenge_StoreValueTruePasses tests that store value=true passes (expected true)
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_StoreValueTruePasses() {
	// Create dynamic store with default value false
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Alice's value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)

	// Create challenge that checks the store
	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer initiated by Alice (who has value=true) should succeed
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when initiator has store value=true")
}

// TestDynamicStoreChallenge_StoreValueFalseFails tests that store value=false fails (expected true)
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_StoreValueFalseFails() {
	// Create dynamic store with default value false
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Alice's value to false explicitly
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer initiated by Alice (who has value=false) should fail
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when initiator has store value=false")
}

// TestDynamicStoreChallenge_DefaultValueUsed tests that defaultValue is used when address not set
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_DefaultValueUsed() {
	// Create dynamic store with default value true
	storeId := suite.createDynamicStore(suite.Manager, true)

	// Do NOT set any value for Alice - it should use default

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should succeed because default value is true
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when default value is true and address not set")
}

// TestDynamicStoreChallenge_DefaultValueFalseFails tests default value false causes failure
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_DefaultValueFalseFails() {
	// Create dynamic store with default value false
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Do NOT set any value for Alice - it should use default (false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because default value is false
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when default value is false and address not set")
}

// TestDynamicStoreChallenge_GlobalEnabledFalseFailsAll tests that globalEnabled=false fails all checks
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_GlobalEnabledFalseFailsAll() {
	// Create dynamic store with default value true
	storeId := suite.createDynamicStore(suite.Manager, true)

	// Set Alice's value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)

	// Disable the store globally
	suite.updateDynamicStore(storeId, suite.Manager, true, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because globalEnabled is false (even though Alice has value=true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when globalEnabled=false")
}

// TestDynamicStoreChallenge_OwnershipCheckPartyInitiator tests ownershipCheckParty="initiator"
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartyInitiator() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Alice (initiator) value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)
	// Set Bob (recipient) value to false
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Bob, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator", // Check the initiator
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer initiated by Alice should succeed (Alice=true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when initiator has value=true")
}

// TestDynamicStoreChallenge_OwnershipCheckPartySender tests ownershipCheckParty="sender"
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartySender() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Alice (sender in this case) value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "sender", // Check the sender (from address)
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer from Alice (sender) should succeed (Alice=true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when sender has value=true")
}

// TestDynamicStoreChallenge_OwnershipCheckPartySenderFails tests sender check fails when sender=false
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartySenderFails() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Alice (sender) value to false
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "sender",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because sender (Alice) has value=false
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when sender has value=false")
}

// TestDynamicStoreChallenge_OwnershipCheckPartyRecipient tests ownershipCheckParty="recipient"
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartyRecipient() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Bob (recipient) value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Bob, true)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "recipient", // Check the recipient (to address)
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer to Bob should succeed (Bob=true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when recipient has value=true")
}

// TestDynamicStoreChallenge_OwnershipCheckPartyRecipientFails tests recipient check fails when recipient=false
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartyRecipientFails() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Bob (recipient) value to false
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Bob, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "recipient",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because recipient (Bob) has value=false
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when recipient has value=false")
}

// TestDynamicStoreChallenge_OwnershipCheckPartySpecificAddress tests ownershipCheckParty with specific address
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartySpecificAddress() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Charlie's value to true (Charlie is not involved in the transfer)
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Charlie, true)

	// Check a specific address (Charlie) instead of initiator/sender/recipient
	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: suite.Charlie, // Specific address
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should succeed because Charlie has value=true
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when specific address has value=true")
}

// TestDynamicStoreChallenge_OwnershipCheckPartySpecificAddressFails tests specific address check fails
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OwnershipCheckPartySpecificAddressFails() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Set Charlie's value to false
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Charlie, false)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: suite.Charlie,
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because Charlie has value=false
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when specific address has value=false")
}

// TestDynamicStoreChallenge_MultipleChallengesAllMustPass tests multiple challenges must all pass
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_MultipleChallengesAllMustPass() {
	// Create two stores
	storeId1 := suite.createDynamicStore(suite.Manager, false)
	storeId2 := suite.createDynamicStore(suite.Manager, false)

	// Set Alice's values to true in both stores
	suite.setDynamicStoreValue(storeId1, suite.Manager, suite.Alice, true)
	suite.setDynamicStoreValue(storeId2, suite.Manager, suite.Alice, true)

	dynamicStoreChallenge1 := &types.DynamicStoreChallenge{
		StoreId:             storeId1,
		OwnershipCheckParty: "initiator",
	}
	dynamicStoreChallenge2 := &types.DynamicStoreChallenge{
		StoreId:             storeId2,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge1, dynamicStoreChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should succeed because Alice has value=true in both stores
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when all challenges pass")
}

// TestDynamicStoreChallenge_MultipleChallengesOneFails tests that one failing challenge blocks transfer
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_MultipleChallengesOneFails() {
	// Create two stores
	storeId1 := suite.createDynamicStore(suite.Manager, false)
	storeId2 := suite.createDynamicStore(suite.Manager, false)

	// Set Alice's value to true in store1, but false in store2
	suite.setDynamicStoreValue(storeId1, suite.Manager, suite.Alice, true)
	suite.setDynamicStoreValue(storeId2, suite.Manager, suite.Alice, false)

	dynamicStoreChallenge1 := &types.DynamicStoreChallenge{
		StoreId:             storeId1,
		OwnershipCheckParty: "initiator",
	}
	dynamicStoreChallenge2 := &types.DynamicStoreChallenge{
		StoreId:             storeId2,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge1, dynamicStoreChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because Alice has value=false in store2
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when one challenge fails")
}

// TestDynamicStoreChallenge_ValueChangeAffectsTransfer tests that changing value affects future transfers
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_ValueChangeAffectsTransfer() {
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Initially set Alice's value to true
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(20, 1)})

	// First transfer should succeed (Alice=true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "first transfer should succeed")

	// Change Alice's value to false
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, false)

	// Second transfer should fail (Alice=false now)
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "second transfer should fail after value changed to false")
}

// TestDynamicStoreChallenge_GlobalEnabledToggle tests toggling globalEnabled
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_GlobalEnabledToggle() {
	storeId := suite.createDynamicStore(suite.Manager, false)
	suite.setDynamicStoreValue(storeId, suite.Manager, suite.Alice, true)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             storeId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(30, 1)})

	// First transfer should succeed (globalEnabled defaults to true)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when globalEnabled=true")

	// Disable globally
	suite.updateDynamicStore(storeId, suite.Manager, false, false)

	// Second transfer should fail
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "transfer should fail when globalEnabled=false")

	// Re-enable globally
	suite.updateDynamicStore(storeId, suite.Manager, false, true)

	// Third transfer should succeed again
	transferMsg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg3)
	suite.Require().NoError(err, "transfer should succeed after re-enabling globalEnabled")
}

// TestDynamicStoreChallenge_NonExistentStoreFails tests that non-existent store causes failure
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_NonExistentStoreFails() {
	// Use a store ID that doesn't exist
	nonExistentStoreId := sdkmath.NewUint(99999)

	dynamicStoreChallenge := &types.DynamicStoreChallenge{
		StoreId:             nonExistentStoreId,
		OwnershipCheckParty: "initiator",
	}

	approval := testutil.GenerateCollectionApproval("dynamic_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		DynamicStoreChallenges:         []*types.DynamicStoreChallenge{dynamicStoreChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Transfer should fail because store doesn't exist
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when dynamic store doesn't exist")
}

// TestDynamicStoreChallenge_OnlyCreatorCanSetValue tests that only creator can set values
func (suite *DynamicStoreChallengeTestSuite) TestDynamicStoreChallenge_OnlyCreatorCanSetValue() {
	// Manager creates the store
	storeId := suite.createDynamicStore(suite.Manager, false)

	// Alice (not creator) tries to set value
	msg := &types.MsgSetDynamicStoreValue{
		Creator: suite.Alice,
		StoreId: storeId,
		Address: suite.Bob,
		Value:   true,
	}

	_, err := suite.MsgServer.SetDynamicStoreValue(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-creator should not be able to set value")
}
