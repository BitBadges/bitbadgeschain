package edge_cases

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

var _ = testutil.GenerateCollectionApproval // Ensure testutil is imported

type CoinDenomValidationTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestCoinDenomValidationTestSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CoinDenomValidationTestSuite))
}

func (suite *CoinDenomValidationTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)

	// Set up mint approval so we can mint badges
	suite.SetupMintApproval(suite.CollectionId)

	// Set up collection approval to allow transfers between users
	// This is needed so transfers can proceed to coin validation
	wctx := sdk.WrapSDKContext(suite.Ctx)
	collectionApproval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	updateMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{collectionApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(wctx, updateMsg)
	suite.Require().NoError(err, "should set up collection approval")
}

// getFullUintRanges returns a full uint range (1 to MaxUint64)
func getFullUintRanges() []*tokenizationtypes.UintRange {
	return []*tokenizationtypes.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
	}
}

// TestCoinDenomValidation_MultiCoinInsufficientBalance tests that coin transfers
// with multiple coins fail gracefully when the sender has insufficient balance.
// Uses allowed denoms (ubadge) so that the transfer reaches the balance check.
func (suite *CoinDenomValidationTestSuite) TestCoinDenomValidation_MultiCoinInsufficientBalance() {
	wctx := sdk.WrapSDKContext(suite.Ctx)
	collectionId := suite.CollectionId

	// Mint badges to Bob
	suite.MintTokens(collectionId, suite.Bob, []*tokenizationtypes.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: getFullUintRanges(),
		},
	})

	// Create an approval with coin transfer that has multiple coins.
	// The sender (Bob) does not hold these coins, so the transfer should fail
	// with an insufficient balance error.
	updateMsg := &tokenizationtypes.MsgUpdateUserApprovals{
		Creator:                 suite.Bob,
		CollectionId:            collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*tokenizationtypes.UserOutgoingApproval{
			{
				ToListId:          suite.Alice,
				InitiatedByListId: suite.Alice,
				TransferTimes:     getFullUintRanges(),
				OwnershipTimes:    getFullUintRanges(),
				TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "multi-coin-test",
				ApprovalCriteria: &tokenizationtypes.OutgoingApprovalCriteria{
					MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*tokenizationtypes.CoinTransfer{
						{
							To: suite.Alice,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},
								{Amount: sdkmath.NewInt(50), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(wctx, updateMsg)
	suite.Require().NoError(err, "should create approval with multi-coin transfer")

	// Set incoming approval for Alice
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(wctx, setIncomingMsg)
	suite.Require().NoError(err, "should create incoming approval")

	// Get the actual version of the approval that was created
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	var actualVersion sdkmath.Uint = sdkmath.NewUint(0)
	for _, approval := range bobBalance.OutgoingApprovals {
		if approval.ApprovalId == "multi-coin-test" {
			actualVersion = approval.Version
			break
		}
	}

	// Attempt to execute transfer - should FAIL because Bob has insufficient coin balance
	transferMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Alice},
				Balances: []*tokenizationtypes.Balance{
					{
						OwnershipTimes: getFullUintRanges(),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*tokenizationtypes.ApprovalIdentifierDetails{
					{
						ApprovalId:      "multi-coin-test",
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Bob,
						Version:         actualVersion,
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(wctx, transferMsg)
	suite.Require().Error(err, "should fail when sender has insufficient coin balance")
	suite.Require().Contains(err.Error(), "insufficient", "error should mention insufficient balance")
}

// TestCoinDenomValidation_AllCoinsCheckedForBalance tests that all coins in a
// multi-coin transfer are checked. Even if the first coin has sufficient balance,
// the transfer should fail if a subsequent coin has insufficient balance.
func (suite *CoinDenomValidationTestSuite) TestCoinDenomValidation_AllCoinsCheckedForBalance() {
	wctx := sdk.WrapSDKContext(suite.Ctx)
	collectionId := suite.CollectionId

	// Mint badges to Bob
	suite.MintTokens(collectionId, suite.Bob, []*tokenizationtypes.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: getFullUintRanges(),
		},
	})

	// Coin transfer with an allowed denom Bob doesn't hold enough of
	updateMsg := &tokenizationtypes.MsgUpdateUserApprovals{
		Creator:                 suite.Bob,
		CollectionId:            collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*tokenizationtypes.UserOutgoingApproval{
			{
				ToListId:          suite.Alice,
				InitiatedByListId: suite.Alice,
				TransferTimes:     getFullUintRanges(),
				OwnershipTimes:    getFullUintRanges(),
				TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "reverse-order-test",
				ApprovalCriteria: &tokenizationtypes.OutgoingApprovalCriteria{
					MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*tokenizationtypes.CoinTransfer{
						{
							To: suite.Alice,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(50), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(wctx, updateMsg)
	suite.Require().NoError(err, "should create approval")

	// Set incoming approval for Alice
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(wctx, setIncomingMsg)
	suite.Require().NoError(err, "should create incoming approval")

	// Get the actual version of the approval
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	var actualVersion sdkmath.Uint = sdkmath.NewUint(0)
	for _, approval := range bobBalance.OutgoingApprovals {
		if approval.ApprovalId == "reverse-order-test" {
			actualVersion = approval.Version
			break
		}
	}

	// Attempt transfer - should fail because Bob has no uatom
	transferMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Alice},
				Balances: []*tokenizationtypes.Balance{
					{
						OwnershipTimes: getFullUintRanges(),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*tokenizationtypes.ApprovalIdentifierDetails{
					{
						ApprovalId:      "reverse-order-test",
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Bob,
						Version:         actualVersion,
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(wctx, transferMsg)
	suite.Require().Error(err, "should fail when sender has insufficient coin balance")
	suite.Require().Contains(err.Error(), "insufficient", "error should mention insufficient balance")
}
