package keeper_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type BalanceManagementTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestBalanceManagementSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(BalanceManagementTestSuite))
}

func (suite *BalanceManagementTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestBalanceManagement_GetBalanceOrApplyDefault tests default balance application
func (suite *BalanceManagementTestSuite) TestBalanceManagement_GetBalanceOrApplyDefault() {
	// Get balance for address that doesn't exist - should apply default
	collection := suite.GetCollection(suite.CollectionId)
	balance, appliedDefault := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)

	suite.Require().True(appliedDefault, "default should be applied for new address")
	suite.Require().NotNil(balance)
	suite.Require().NotNil(balance.UserPermissions)
}

// TestBalanceManagement_SetBalanceForAddress tests setting balance for address
func (suite *BalanceManagementTestSuite) TestBalanceManagement_SetBalanceForAddress() {
	collection := suite.GetCollection(suite.CollectionId)

	// Create custom balance
	customBalance := &types.UserBalanceStore{
		Balances: []*types.Balance{
			{
				Amount: sdkmath.NewUint(100),
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
			},
		},
		UserPermissions: &types.UserPermissions{},
	}

	// Set balance
	err := suite.Keeper.SetBalanceForAddress(suite.Ctx, collection, suite.Alice, customBalance)
	suite.Require().NoError(err)

	// Get balance and verify
	balance, appliedDefault := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().False(appliedDefault, "default should not be applied for existing balance")
	suite.Require().Equal(1, len(balance.Balances))
	suite.Require().True(balance.Balances[0].Amount.Equal(sdkmath.NewUint(100)))
}

// TestBalanceManagement_DefaultInheritance tests that defaults are inherited correctly
func (suite *BalanceManagementTestSuite) TestBalanceManagement_DefaultInheritance() {
	// Create collection with specific default balances
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Get balance for new address - should inherit defaults
	collection := suite.GetCollection(collectionId)
	balance, appliedDefault := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)

	suite.Require().True(appliedDefault, "default should be applied")
	suite.Require().NotNil(balance)
	// Balance should match collection defaults
	suite.Require().NotNil(collection.DefaultBalances)
}

// TestBalanceManagement_SpecialAddresses tests balance handling for special addresses
func (suite *BalanceManagementTestSuite) TestBalanceManagement_SpecialAddresses() {
	collection := suite.GetCollection(suite.CollectionId)

	// Test Mint address
	mintBalance, appliedDefault := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, types.MintAddress)
	suite.Require().False(appliedDefault, "Mint address should not apply default")
	suite.Require().NotNil(mintBalance)
	suite.Require().Equal(0, len(mintBalance.Balances), "Mint address should have empty balances (unlimited)")

	// Test Total address
	totalBalance, appliedDefault := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, types.TotalAddress)
	suite.Require().False(appliedDefault, "Total address should not apply default")
	suite.Require().NotNil(totalBalance)
	suite.Require().Equal(0, len(totalBalance.Balances), "Total address should have empty balances (unlimited)")
}

// TestBalanceManagement_BalanceUpdate tests balance updates during transfers
func (suite *BalanceManagementTestSuite) TestBalanceManagement_BalanceUpdate() {
	// Setup approvals for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
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

	// Get balance before transfer
	balanceBefore := suite.GetBalance(suite.CollectionId, suite.Alice)

	// Transfer tokens
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(50, 1),
			}),
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err)

	// Get balance after transfer
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)

	// Balance should have changed
	suite.Require().NotEqual(balanceBefore.Balances, balanceAfter.Balances, "balance should change after transfer")
}

