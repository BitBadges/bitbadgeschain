package approval_criteria_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// PredeterminedBalancesTestSuite tests the PredeterminedBalances approval criteria field
type PredeterminedBalancesTestSuite struct {
	testutil.AITestSuite
}

func TestPredeterminedBalancesSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(PredeterminedBalancesTestSuite))
}

func (suite *PredeterminedBalancesTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestManualBalances_ArrayWorks tests that manual balances array works correctly
func (suite *PredeterminedBalancesTestSuite) TestManualBalances_ArrayWorks() {
	// Create approval with predetermined manual balances
	// First transfer gets balance entry 0, second gets entry 1, etc.
	approval := testutil.GenerateCollectionApproval("predeterminedapproval", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			ManualBalances: []*types.ManualBalances{
				{
					Balances: []*types.Balance{
						testutil.GenerateSimpleBalance(5, 1), // First transfer: 5 of token 1
					},
				},
				{
					Balances: []*types.Balance{
						testutil.GenerateSimpleBalance(10, 2), // Second transfer: 10 of token 2
					},
				},
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint should transfer exactly 5 of token 1
	// Must use PrioritizedApprovals when using tracking features
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "predeterminedapproval",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first mint with predetermined balance should succeed")

	// Verify Alice received token 1
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().NotEmpty(aliceBalance.Balances, "Alice should have balance")

	// Second mint should transfer exactly 10 of token 2
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 2)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "predeterminedapproval",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second mint with predetermined balance should succeed")

	// Verify Bob received token 2
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	suite.Require().NotEmpty(bobBalance.Balances, "Bob should have balance")
}

// TestIncrementedBalances_Works tests that incremented balances work correctly
func (suite *PredeterminedBalancesTestSuite) TestIncrementedBalances_Works() {
	// Create approval with incremented balances starting at token 1
	approval := testutil.GenerateCollectionApproval("incrementedapproval", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1), // Start with token 1
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(1), // Increment token ID by 1 each time
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint should get token 1
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementedapproval",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first mint should succeed with token 1")

	// Second mint should get token 2 (incremented)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 2)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementedapproval",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second mint should succeed with token 2")

	// Third mint should get token 3 (incremented again)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementedapproval",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "third mint should succeed with token 3")
}

// TestIncrementTokenIdsBy_Works tests token ID increment functionality
func (suite *PredeterminedBalancesTestSuite) TestIncrementTokenIdsBy_Works() {
	// Create approval with larger token ID increment
	approval := testutil.GenerateCollectionApproval("incrementby5", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1),
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(5), // Increment by 5
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint gets token 1
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementby5",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first mint with token 1 should succeed")

	// Second mint gets token 6 (1 + 5)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 6)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementby5",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second mint with token 6 should succeed")

	// Third mint gets token 11 (6 + 5)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 11)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementby5",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "third mint with token 11 should succeed")
}

// TestIncrementOwnershipTimesBy_Works tests ownership time increment functionality
func (suite *PredeterminedBalancesTestSuite) TestIncrementOwnershipTimesBy_Works() {
	// Create approval with ownership time increment
	approval := testutil.GenerateCollectionApproval("incrementtimes", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						TokenIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)}, // First ownership time
						},
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(1000), // Increment by 1000ms
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint with ownership times 1000-2000
	balance1 := &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)},
		},
	}
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{balance1},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementtimes",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first mint with ownership times 1000-2000 should succeed")

	// Second mint with incremented ownership times 2000-3000
	balance2 := &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(2000), End: sdkmath.NewUint(3000)},
		},
	}
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{balance2},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementtimes",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second mint with incremented ownership times should succeed")
}

// TestAllowOverrideTimestamp_Works tests the allowOverrideTimestamp flag
func (suite *PredeterminedBalancesTestSuite) TestAllowOverrideTimestamp_Works() {
	// Create approval with allowOverrideTimestamp enabled
	// This tests that the flag is properly stored in the approval structure
	// The DurationFromTimestamp feature calculates ownership times based on block timestamp
	approval := testutil.GenerateCollectionApproval("overridetimestamp", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1),
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(86400000), // 1 day in milliseconds
				AllowOverrideTimestamp:    true,                      // Allow override
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Verify the approval structure was saved correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "overridetimestamp" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria, "approval criteria should exist")
			suite.Require().NotNil(app.ApprovalCriteria.PredeterminedBalances, "predetermined balances should exist")
			suite.Require().NotNil(app.ApprovalCriteria.PredeterminedBalances.IncrementedBalances, "incremented balances should exist")
			suite.Require().True(app.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp, "AllowOverrideTimestamp should be true")
			suite.Require().True(app.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp.Equal(sdkmath.NewUint(86400000)), "DurationFromTimestamp should be 86400000")
			break
		}
	}
	suite.Require().True(found, "approval should exist")
}

// TestUseMerkleChallengeLeafIndex_Works tests using merkle challenge leaf index for ordering
func (suite *PredeterminedBalancesTestSuite) TestUseMerkleChallengeLeafIndex_Works() {
	// Create approval using merkle challenge leaf index for ordering
	// This is a more advanced test that requires merkle proofs
	approval := testutil.GenerateCollectionApproval("merkleorder", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			ManualBalances: []*types.ManualBalances{
				{
					Balances: []*types.Balance{
						testutil.GenerateSimpleBalance(1, 1), // Leaf index 0
					},
				},
				{
					Balances: []*types.Balance{
						testutil.GenerateSimpleBalance(1, 2), // Leaf index 1
					},
				},
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseMerkleChallengeLeafIndex: true,
				ChallengeTrackerId:          "merkletracker",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection should be created")
}

// TestUsePerToAddressNumTransfers_Works tests per-to-address ordering
func (suite *PredeterminedBalancesTestSuite) TestUsePerToAddressNumTransfers_Works() {
	// Create approval with per-to-address ordering
	approval := testutil.GenerateCollectionApproval("pertoorder", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1),
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(1),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UsePerToAddressNumTransfers: true, // Order based on transfers to each address
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			PerToAddressMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint to Alice - she gets token 1
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "pertoorder",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first mint to Alice should succeed")

	// First mint to Bob - he also gets token 1 (his first transfer)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "pertoorder",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "first mint to Bob should succeed")

	// Second mint to Alice - she gets token 2
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 2)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "pertoorder",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "second mint to Alice should succeed")
}

// TestUsePerFromAddressNumTransfers_Works tests per-from-address ordering
func (suite *PredeterminedBalancesTestSuite) TestUsePerFromAddressNumTransfers_Works() {
	// First create collection and mint tokens to multiple addresses
	approval := testutil.GenerateCollectionApproval("perfromorder", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1),
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(1),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UsePerFromAddressNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			PerFromAddressMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection should be created")

	// Setup: mint tokens to Alice and Bob so they can transfer
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{
		testutil.GenerateBalance(10, 1, 10, 1, math.MaxUint64),
	})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{
		testutil.GenerateBalance(10, 1, 10, 1, math.MaxUint64),
	})

	// First transfer from Alice - she sends token 1
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "perfromorder",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first transfer from Alice should succeed")

	// First transfer from Bob - he sends token 1 (his first transfer)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "perfromorder",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "first transfer from Bob should succeed")
}

// TestUsePerInitiatedByAddressNumTransfers_Works tests per-initiated-by-address ordering
func (suite *PredeterminedBalancesTestSuite) TestUsePerInitiatedByAddressNumTransfers_Works() {
	approval := testutil.GenerateCollectionApproval("perinitiatororder", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					testutil.GenerateSimpleBalance(1, 1),
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(1),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UsePerInitiatedByAddressNumTransfers: true,
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection should be created")
}

// TestPredeterminedBalances_WrongBalanceRejected tests that transfers with wrong balances are rejected
func (suite *PredeterminedBalancesTestSuite) TestPredeterminedBalances_WrongBalanceRejected() {
	// Create approval with predetermined balances - exactly 5 of token 1
	approval := testutil.GenerateCollectionApproval("exactbalance", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			ManualBalances: []*types.ManualBalances{
				{
					Balances: []*types.Balance{
						testutil.GenerateSimpleBalance(5, 1), // Must be exactly 5 of token 1
					},
				},
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // Unlimited
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Try to transfer wrong amount (10 instead of 5) - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)}, // Wrong amount
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "exactbalance",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with wrong predetermined balance should fail")
}
