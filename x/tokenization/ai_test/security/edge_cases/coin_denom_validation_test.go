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

// TestCoinDenomValidation_MultiCoinUnauthorizedDenom tests MED-012 fix:
// Validates that ALL coins in a transfer are checked, not just the first.
// This prevents unauthorized denoms from bypassing validation in multi-coin transfers.
func (suite *CoinDenomValidationTestSuite) TestCoinDenomValidation_MultiCoinUnauthorizedDenom() {
	wctx := sdk.WrapSDKContext(suite.Ctx)
	collectionId := suite.CollectionId

	// Set allowed denoms to only include "ubadge"
	params := suite.Keeper.GetParams(suite.Ctx)
	params.AllowedDenoms = []string{"ubadge"}
	err := suite.Keeper.SetParams(suite.Ctx, params)
	suite.Require().NoError(err, "should set params")

	// Mint badges to Bob
	suite.MintBadges(collectionId, suite.Bob, []*tokenizationtypes.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: getFullUintRanges(),
		},
	})

	// Create an approval with coin transfer that has multiple coins:
	// - First coin: "ubadge" (allowed)
	// - Second coin: "uatom" (NOT allowed - should cause failure)
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
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"}, // Allowed denom
								{Amount: sdkmath.NewInt(50), Denom: "uatom"},   // Unauthorized denom - should fail
							},
						},
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(wctx, updateMsg)
	suite.Require().NoError(err, "should create approval with multi-coin transfer")

	// Set incoming approval for Alice (needed for transfer to succeed up to coin validation)
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

	// Attempt to execute transfer - should FAIL because "uatom" is not in allowed denoms
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
						Version:         actualVersion, // Use the actual version from the created approval
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(wctx, transferMsg)
	suite.Require().Error(err, "should fail when unauthorized denom is in multi-coin transfer")
	suite.Require().Contains(err.Error(), "denom uatom is not allowed", "error should mention unauthorized denom")
}

// TestCoinDenomValidation_AllCoinsValidated tests that all coins are validated even if first is allowed
func (suite *CoinDenomValidationTestSuite) TestCoinDenomValidation_AllCoinsValidated() {
	wctx := sdk.WrapSDKContext(suite.Ctx)
	collectionId := suite.CollectionId

	// Set allowed denoms to only include "ubadge"
	params := suite.Keeper.GetParams(suite.Ctx)
	params.AllowedDenoms = []string{"ubadge"}
	err := suite.Keeper.SetParams(suite.Ctx, params)
	suite.Require().NoError(err, "should set params")

	// Mint badges to Bob
	suite.MintBadges(collectionId, suite.Bob, []*tokenizationtypes.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: getFullUintRanges(),
		},
	})

	// Test case: First coin is unauthorized, second is authorized
	// Should fail on first coin check
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
								{Amount: sdkmath.NewInt(50), Denom: "uatom"},   // Unauthorized denom first
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"}, // Allowed denom second
							},
						},
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(wctx, updateMsg)
	suite.Require().NoError(err, "should create approval")

	// Set incoming approval for Alice (needed for transfer to succeed up to coin validation)
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
		if approval.ApprovalId == "reverse-order-test" {
			actualVersion = approval.Version
			break
		}
	}

	// Attempt transfer - should fail on first coin
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
						Version:         actualVersion, // Use the actual version from the created approval
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(wctx, transferMsg)
	suite.Require().Error(err, "should fail when unauthorized denom is present")
	suite.Require().Contains(err.Error(), "denom uatom is not allowed", "error should mention unauthorized denom")
}
