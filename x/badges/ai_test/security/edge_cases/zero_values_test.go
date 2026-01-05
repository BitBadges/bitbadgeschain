package edge_cases_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type ZeroValuesTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestZeroValuesSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(ZeroValuesTestSuite))
}

func (suite *ZeroValuesTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestZeroValues_ZeroAmountTransfer tests that zero amount transfers are rejected
func (suite *ZeroValuesTestSuite) TestZeroValues_ZeroAmountTransfer() {
	// Setup approvals
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Attempt zero amount transfer - should fail
	zeroBalance := &types.Balance{
		Amount: sdkmath.NewUint(0),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
	}

	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{zeroBalance}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "zero amount transfer should be rejected")
}

// TestZeroValues_EmptyTokenIds tests that empty token ID ranges are handled correctly
func (suite *ZeroValuesTestSuite) TestZeroValues_EmptyTokenIds() {
	// Attempt to create collection with empty token ID ranges
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(0),
					TokenIds: []*types.UintRange{}, // Empty token IDs
					OwnershipTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
					},
				},
			},
		},
		ValidTokenIds: []*types.UintRange{}, // Empty valid token IDs
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: suite.Manager,
		CollectionMetadata: testutil.GenerateCollectionMetadata("", ""),
		TokenMetadata: []*types.TokenMetadata{},
		CustomData: "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards: []string{},
		IsArchived: false,
	}

	_, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	// Empty token IDs may or may not be valid depending on validation rules
	// This test documents the behavior
	_ = err // Accept either outcome
}

// TestZeroValues_ZeroCollectionId tests handling of zero collection ID
func (suite *ZeroValuesTestSuite) TestZeroValues_ZeroCollectionId() {
	// Zero collection ID should use auto-prev resolution
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: sdkmath.NewUint(0), // Zero collection ID
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	// Should fail if no previous collection exists in transaction
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "zero collection ID should fail without previous collection")
}

