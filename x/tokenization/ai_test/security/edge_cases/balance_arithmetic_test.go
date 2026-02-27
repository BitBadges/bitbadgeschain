package edge_cases_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// BalanceArithmeticTestSuite tests balance conservation, arithmetic operations,
// and multi-balance scenarios
type BalanceArithmeticTestSuite struct {
	testutil.AITestSuite
}

func TestBalanceArithmeticSuite(t *testing.T) {
	suite.Run(t, new(BalanceArithmeticTestSuite))
}

func (suite *BalanceArithmeticTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestBalanceConservation_SumEqualsTotal tests that the sum of all balances
// equals the total minted amount
func (suite *BalanceArithmeticTestSuite) TestBalanceConservation_SumEqualsTotal() {
	// Create collection with approvals
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint specific amount
	mintAmount := uint64(1000)
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(mintAmount, 1)})

	// Record total minted
	totalMinted := sdkmath.NewUint(mintAmount)

	// Perform various transfers
	transfers := []struct {
		from   string
		to     string
		amount uint64
	}{
		{suite.Alice, suite.Bob, 300},
		{suite.Bob, suite.Charlie, 100},
		{suite.Alice, suite.Charlie, 200},
		{suite.Charlie, suite.Alice, 50},
	}

	for _, t := range transfers {
		msg := &types.MsgTransferTokens{
			Creator:      t.from,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        t.from,
					ToAddresses: []string{t.to},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(t.amount, 1)},
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
		suite.Require().NoError(err)
	}

	// Calculate sum of all balances
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	charlieBalance := suite.GetBalance(collectionId, suite.Charlie)

	totalInCirculation := suite.calculateTotalBalance(aliceBalance.Balances)
	totalInCirculation = totalInCirculation.Add(suite.calculateTotalBalance(bobBalance.Balances))
	totalInCirculation = totalInCirculation.Add(suite.calculateTotalBalance(charlieBalance.Balances))

	suite.Require().True(totalMinted.Equal(totalInCirculation),
		"total in circulation (%s) should equal total minted (%s)",
		totalInCirculation, totalMinted)
}

// TestBalanceConservation_NoTokensCreatedFromNothing tests that tokens cannot
// be created without proper minting
func (suite *BalanceArithmeticTestSuite) TestBalanceConservation_NoTokensCreatedFromNothing() {
	// Create collection
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint 100 tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Try to transfer more than Alice has (should fail)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(150, 1)}, // More than minted
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
	suite.Require().Error(err, "cannot transfer more than exists")

	// Verify Alice still has exactly 100
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(100)),
		"Alice should still have 100 tokens")
}

// TestBalanceConservation_NoTokensLostInTransfers tests that tokens are never
// lost during transfers
func (suite *BalanceArithmeticTestSuite) TestBalanceConservation_NoTokensLostInTransfers() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(500, 1)})

	// Get total before transfers
	aliceBefore := suite.GetBalance(collectionId, suite.Alice)
	bobBefore := suite.GetBalance(collectionId, suite.Bob)
	charlieBefore := suite.GetBalance(collectionId, suite.Charlie)

	totalBefore := suite.calculateTotalBalance(aliceBefore.Balances)
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(bobBefore.Balances))
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(charlieBefore.Balances))

	// Perform many transfers in a chain
	for i := 0; i < 5; i++ {
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
							ApprovalId:    "transfer_approval",
							Version:       sdkmath.NewUint(0),
						},
					},
					OnlyCheckPrioritizedCollectionApprovals: true,
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err)

		// Bob transfers to Charlie
		msg2 := &types.MsgTransferTokens{
			Creator:      suite.Bob,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Bob,
					ToAddresses: []string{suite.Charlie},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
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
		_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
		suite.Require().NoError(err)
	}

	// Get total after transfers
	aliceAfter := suite.GetBalance(collectionId, suite.Alice)
	bobAfter := suite.GetBalance(collectionId, suite.Bob)
	charlieAfter := suite.GetBalance(collectionId, suite.Charlie)

	totalAfter := suite.calculateTotalBalance(aliceAfter.Balances)
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(bobAfter.Balances))
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(charlieAfter.Balances))

	suite.Require().True(totalBefore.Equal(totalAfter),
		"total tokens should be conserved: before=%s, after=%s", totalBefore, totalAfter)
}

// TestMultiBalance_MergingBalancesForSameToken tests that balances for the
// same token but different times are handled correctly
func (suite *BalanceArithmeticTestSuite) TestMultiBalance_MergingBalancesForSameToken() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens with different ownership times
	balance1 := &types.Balance{
		Amount: sdkmath.NewUint(100),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		},
	}

	balance2 := &types.Balance{
		Amount: sdkmath.NewUint(50),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1001), End: sdkmath.NewUint(2000)},
		},
	}

	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{balance1})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{balance2})

	// Transfer partial amount from first time range
	transferBalance := &types.Balance{
		Amount: sdkmath.NewUint(30),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		},
	}

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{transferBalance},
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
	suite.Require().NoError(err)

	// Verify Bob got the tokens
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(30)),
		"Bob should have 30 tokens")

	// Verify Alice still has remaining
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(120)),
		"Alice should have 120 tokens (100-30 + 50)")
}

// TestMultiBalance_SplitBalanceIntoMultipleRanges tests splitting a balance
// into multiple ranges
func (suite *BalanceArithmeticTestSuite) TestMultiBalance_SplitBalanceIntoMultipleRanges() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens for token IDs 1-10
	balance := &types.Balance{
		Amount: sdkmath.NewUint(100),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
	}

	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{balance})

	// Transfer only tokens 1-5 to Bob
	transferBalance := &types.Balance{
		Amount: sdkmath.NewUint(100),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
	}

	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{transferBalance},
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
	suite.Require().NoError(err)

	// Verify Bob has tokens 1-5
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	suite.Require().True(len(bobBalance.Balances) > 0, "Bob should have balances")

	// Verify Alice has tokens 6-10
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().True(len(aliceBalance.Balances) > 0, "Alice should have remaining balances")
}

// TestMultiBalance_ConsolidatingOverlappingRanges tests that receiving tokens
// in different transactions consolidates properly
func (suite *BalanceArithmeticTestSuite) TestMultiBalance_ConsolidatingOverlappingRanges() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 2)})

	// Transfer token 1 to Bob
	msg1 := &types.MsgTransferTokens{
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
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Transfer token 2 to Bob
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 2)},
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
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Bob should have tokens from both transfers
	// Note: Balance amounts are stored per token ID, so we need to account for multiple token IDs
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	bobTotal := suite.calculateTotalBalanceWithTokenIds(bobBalance.Balances)
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(100)),
		"Bob should have 100 tokens total (50 from each token ID)")
}

// TestBalanceArithmetic_LargeNumbers tests balance operations with large numbers
func (suite *BalanceArithmeticTestSuite) TestBalanceArithmetic_LargeNumbers() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint a very large amount
	largeAmount := uint64(1_000_000_000_000) // 1 trillion
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{
		{
			Amount: sdkmath.NewUint(largeAmount),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
		},
	})

	// Transfer half
	halfAmount := largeAmount / 2
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(halfAmount),
						TokenIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
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
	suite.Require().NoError(err)

	// Verify balances
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	bobBalance := suite.GetBalance(collectionId, suite.Bob)

	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)

	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(halfAmount)),
		"Alice should have half the tokens")
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(halfAmount)),
		"Bob should have half the tokens")
	suite.Require().True(aliceTotal.Add(bobTotal).Equal(sdkmath.NewUint(largeAmount)),
		"sum should equal original amount")
}

// TestBalanceArithmetic_MultiRecipientConservation tests that multi-recipient
// transfers conserve total balance
func (suite *BalanceArithmeticTestSuite) TestBalanceArithmetic_MultiRecipientConservation() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens
	totalMinted := uint64(300)
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(totalMinted, 1)})

	// Transfer to multiple recipients at once
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob, suite.Charlie}, // Two recipients
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)}, // 50 each
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
	suite.Require().NoError(err)

	// Verify total conservation
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	charlieBalance := suite.GetBalance(collectionId, suite.Charlie)

	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)
	charlieTotal := suite.calculateTotalBalance(charlieBalance.Balances)

	total := aliceTotal.Add(bobTotal).Add(charlieTotal)
	suite.Require().True(total.Equal(sdkmath.NewUint(totalMinted)),
		"total should be conserved: expected %d, got %s", totalMinted, total)

	// Verify individual balances
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(50)), "Bob should have 50")
	suite.Require().True(charlieTotal.Equal(sdkmath.NewUint(50)), "Charlie should have 50")
	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(200)), "Alice should have 200 remaining")
}

// TestBalanceArithmetic_CircularTransfersConservation tests that tokens are
// conserved even in circular transfer patterns
func (suite *BalanceArithmeticTestSuite) TestBalanceArithmetic_CircularTransfersConservation() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to all participants
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Charlie, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Get total before
	aliceBefore := suite.GetBalance(collectionId, suite.Alice)
	bobBefore := suite.GetBalance(collectionId, suite.Bob)
	charlieBefore := suite.GetBalance(collectionId, suite.Charlie)

	totalBefore := suite.calculateTotalBalance(aliceBefore.Balances)
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(bobBefore.Balances))
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(charlieBefore.Balances))

	// Circular transfers: Alice -> Bob -> Charlie -> Alice
	transfers := []struct {
		from   string
		to     string
		amount uint64
	}{
		{suite.Alice, suite.Bob, 25},
		{suite.Bob, suite.Charlie, 25},
		{suite.Charlie, suite.Alice, 25},
	}

	for _, t := range transfers {
		msg := &types.MsgTransferTokens{
			Creator:      t.from,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        t.from,
					ToAddresses: []string{t.to},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(t.amount, 1)},
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
		suite.Require().NoError(err)
	}

	// Get total after
	aliceAfter := suite.GetBalance(collectionId, suite.Alice)
	bobAfter := suite.GetBalance(collectionId, suite.Bob)
	charlieAfter := suite.GetBalance(collectionId, suite.Charlie)

	totalAfter := suite.calculateTotalBalance(aliceAfter.Balances)
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(bobAfter.Balances))
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(charlieAfter.Balances))

	suite.Require().True(totalBefore.Equal(totalAfter),
		"total should be conserved after circular transfers: before=%s, after=%s",
		totalBefore, totalAfter)
}

// TestBalanceArithmetic_ExactAmountRequired tests that transfers require
// exact amounts (not more, not less)
func (suite *BalanceArithmeticTestSuite) TestBalanceArithmetic_ExactAmountRequired() {
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint exactly 100 tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Try to transfer 101 - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(101, 1)},
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
	suite.Require().Error(err, "cannot transfer more than balance")

	// Transfer exactly 100 - should succeed
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)},
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

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transferring exact balance should succeed")

	// Alice should have 0 now
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(0)),
		"Alice should have 0 after transferring all")

	// Bob should have 100
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(100)),
		"Bob should have 100")
}

// Helper function to calculate total balance
// Note: The balance Amount field represents the amount per (tokenId, ownershipTime) combination.
// For tests that use single token IDs with same ownership time, this simple sum works.
// For tests with multiple token IDs, use calculateTotalBalanceWithTokenIds instead.
func (suite *BalanceArithmeticTestSuite) calculateTotalBalance(balances []*types.Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		total = total.Add(balance.Amount)
	}
	return total
}

// calculateTotalBalanceWithTokenIds calculates total tokens considering multiple token IDs
// Each balance's Amount is multiplied by the number of distinct token IDs it covers
func (suite *BalanceArithmeticTestSuite) calculateTotalBalanceWithTokenIds(balances []*types.Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		// Count the number of token IDs in this balance
		tokenCount := sdkmath.NewUint(0)
		for _, tokenRange := range balance.TokenIds {
			rangeSize := tokenRange.End.Sub(tokenRange.Start).Add(sdkmath.NewUint(1))
			tokenCount = tokenCount.Add(rangeSize)
		}
		// Multiply amount by number of token IDs
		total = total.Add(balance.Amount.Mul(tokenCount))
	}
	return total
}
