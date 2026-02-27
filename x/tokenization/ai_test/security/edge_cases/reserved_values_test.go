package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ReservedValuesTestSuite tests handling of reserved list IDs like
// "Mint", "!Mint", "All", "!All", "AllWithoutMint", etc.
type ReservedValuesTestSuite struct {
	testutil.AITestSuite
}

func TestReservedValuesSuite(t *testing.T) {
	suite.Run(t, new(ReservedValuesTestSuite))
}

func (suite *ReservedValuesTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestReservedListId_MintMatchesOnlyMintAddress tests that "Mint" list
// matches only the special Mint address
func (suite *ReservedValuesTestSuite) TestReservedListId_MintMatchesOnlyMintAddress() {
	// Create collection with approval that only allows transfers FROM Mint address
	mintOnlyApproval := testutil.GenerateCollectionApproval("mint_only", types.MintAddress, "All")
	mintOnlyApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{mintOnlyApproval})

	// Minting should work (from Mint address)
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Verify Alice received the tokens
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().True(len(aliceBalance.Balances) > 0, "Alice should have tokens")

	// Transfers from regular addresses should fail with this approval
	// (need different approval for regular transfers)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "mint_only",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer from Alice should fail - approval only matches Mint address")
}

// TestReservedListId_NotMintMatchesAllExceptMint tests that "!Mint" list
// matches all addresses except the Mint address
func (suite *ReservedValuesTestSuite) TestReservedListId_NotMintMatchesAllExceptMint() {
	// Create collection with approval for transfers (not from Mint)
	notMintApproval := testutil.GenerateCollectionApproval("not_mint", "!Mint", "All")
	notMintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Also need a mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		notMintApproval,
	})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer from Alice to Bob should work (Alice matches !Mint)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "not_mint",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer from Alice should succeed - !Mint matches regular addresses")

	// Verify Bob received the tokens
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	suite.Require().True(len(bobBalance.Balances) > 0, "Bob should have tokens")
}

// TestReservedListId_AllMatchesAllAddresses tests that "All" list matches
// all addresses. Note: "All" cannot be used in fromListId because the system
// requires mint approvals to use "Mint" specifically, not "All" which includes Mint.
// This test documents the proper usage of "All" in to and initiatedBy fields.
func (suite *ReservedValuesTestSuite) TestReservedListId_AllMatchesAllAddresses() {
	// Create collection with approval using AllWithoutMint for from (proper pattern)
	// and All for to and initiatedBy
	transferApproval := testutil.GenerateCollectionApproval("all_to", "AllWithoutMint", "All")
	transferApproval.InitiatedByListId = "All"
	transferApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Also need mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		transferApproval,
	})

	// Minting works with dedicated mint approval
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to anyone works (All matches all addresses in toListId)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "all_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed - All matches any recipient")

	// Transfer initiated by anyone also works (All in initiatedByListId)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager, // Different initiator than from
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(25, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "all_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer should succeed - All matches any initiator")
}

// TestReservedListId_NotAllMatchesNoAddresses tests that "!All" list
// matches no addresses (empty set)
func (suite *ReservedValuesTestSuite) TestReservedListId_NotAllMatchesNoAddresses() {
	// Create collection with approval from !All (should match nothing)
	noOneApproval := testutil.GenerateCollectionApproval("no_one", "!All", "All")
	noOneApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Need a working mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		noOneApproval,
	})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer using no_one approval should fail (!All matches nothing)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "no_one",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer should fail - !All matches no addresses")
}

// TestReservedListId_AllWithoutMintMatchesAllExceptMint tests the AllWithoutMint
// list ID which is commonly used for transfer approvals
func (suite *ReservedValuesTestSuite) TestReservedListId_AllWithoutMintMatchesAllExceptMint() {
	// Create collection with AllWithoutMint from list
	transferApproval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	transferApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Need a mint approval for minting
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		transferApproval,
	})

	// Mint tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer using transfer_approval should work (Alice matches AllWithoutMint)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "transfer_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed - AllWithoutMint matches regular addresses")
}

// TestReservedListId_UsedInFromField tests reserved list IDs in the from field
func (suite *ReservedValuesTestSuite) TestReservedListId_UsedInFromField() {
	// This is covered by other tests but let's be explicit
	// Test specific address in from field

	// Use Alice's address directly as the fromListId
	specificFromApproval := testutil.GenerateCollectionApproval("alice_only", suite.Alice, "All")
	specificFromApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		specificFromApproval,
	})

	// Mint to both Alice and Bob
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer from Alice should work
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "alice_only",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer from Alice should succeed")

	// Transfer from Bob should fail with this approval
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "alice_only",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer from Bob should fail - approval only allows Alice")
}

// TestReservedListId_UsedInToField tests reserved list IDs in the to field
func (suite *ReservedValuesTestSuite) TestReservedListId_UsedInToField() {
	// Create approval that only allows transfers TO Bob
	toBobOnlyApproval := testutil.GenerateCollectionApproval("to_bob_only", "AllWithoutMint", suite.Bob)
	toBobOnlyApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		toBobOnlyApproval,
	})

	// Mint to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to Bob should work
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "to_bob_only",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to Bob should succeed")

	// Transfer to Charlie should fail with this approval
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "to_bob_only",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer to Charlie should fail - approval only allows Bob as recipient")
}

// TestReservedListId_UsedInInitiatedByField tests reserved list IDs in the initiatedBy field
func (suite *ReservedValuesTestSuite) TestReservedListId_UsedInInitiatedByField() {
	// Create approval that only allows transfers initiated by Manager
	managerInitiatedApproval := testutil.GenerateCollectionApproval("manager_initiated", "AllWithoutMint", "All")
	managerInitiatedApproval.InitiatedByListId = suite.Manager
	managerInitiatedApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		managerInitiatedApproval,
	})

	// Mint to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer initiated by Manager should work (even if from Alice)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager, // Manager initiates
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "manager_initiated",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer initiated by Manager should succeed")

	// Transfer initiated by Alice should fail with this approval
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Alice initiates
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "manager_initiated",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer initiated by Alice should fail - approval only allows Manager to initiate")
}

// TestReservedListId_NoneMatchesEmptySet tests the "None" reserved list ID
func (suite *ReservedValuesTestSuite) TestReservedListId_NoneMatchesEmptySet() {
	// Create approval with None as the to list (should match no recipients)
	noneToApproval := testutil.GenerateCollectionApproval("none_to", "AllWithoutMint", "None")
	noneToApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		noneToApproval,
	})

	// Mint to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to anyone should fail with None as recipient list
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "none_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer should fail - None matches no addresses")
}

// TestReservedListId_ColonSeparatedAddresses tests colon-separated address lists
func (suite *ReservedValuesTestSuite) TestReservedListId_ColonSeparatedAddresses() {
	// Create approval with multiple addresses using colon separator
	multiAddressListId := suite.Bob + ":" + suite.Charlie
	multiToApproval := testutil.GenerateCollectionApproval("multi_to", "AllWithoutMint", multiAddressListId)
	multiToApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		multiToApproval,
	})

	// Mint to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to Bob should work
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(25, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "multi_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to Bob should succeed")

	// Transfer to Charlie should also work
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(25, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "multi_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer to Charlie should succeed")

	// Transfer to Manager should fail (not in the list)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Manager},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "multi_to",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().Error(err, "transfer to Manager should fail - not in Bob:Charlie list")
}

// TestReservedListId_AllWithoutSpecificAddress tests AllWithout<address> pattern
func (suite *ReservedValuesTestSuite) TestReservedListId_AllWithoutSpecificAddress() {
	// Create approval that allows transfers to everyone except Bob
	allWithoutBobListId := "AllWithout" + suite.Bob
	allWithoutBobApproval := testutil.GenerateCollectionApproval("all_without_bob", "AllWithoutMint", allWithoutBobListId)
	allWithoutBobApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		allWithoutBobApproval,
	})

	// Mint to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to Charlie should work (not Bob)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(25, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "all_without_bob",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to Charlie should succeed")

	// Transfer to Bob should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(25, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "all_without_bob",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer to Bob should fail - AllWithout<Bob> excludes Bob")
}
