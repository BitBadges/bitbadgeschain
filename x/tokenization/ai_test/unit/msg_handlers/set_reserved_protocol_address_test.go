package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetReservedProtocolAddressTestSuite struct {
	testutil.AITestSuite
}

func TestSetReservedProtocolAddressSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetReservedProtocolAddressTestSuite))
}

// TestSetReservedProtocolAddress_GovernanceAuthority tests that governance authority can set reserved address
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_GovernanceAuthority() {
	// Get governance authority
	authority := suite.Keeper.GetAuthority()

	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Alice,
		IsReservedProtocol: true,
	}

	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "governance authority should be able to set reserved protocol address")

	// Verify address is now reserved
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Alice)
	suite.Require().True(isReserved, "address should be marked as reserved protocol")
}

// TestSetReservedProtocolAddress_UnsetReservedAddress tests unsetting a reserved address
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_UnsetReservedAddress() {
	authority := suite.Keeper.GetAuthority()

	// First, set the address as reserved
	setMsg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Bob,
		IsReservedProtocol: true,
	}
	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Verify it's reserved
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Bob)
	suite.Require().True(isReserved, "address should be reserved after setting")

	// Now unset the reserved status
	unsetMsg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Bob,
		IsReservedProtocol: false,
	}
	_, err = suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), unsetMsg)
	suite.Require().NoError(err, "governance authority should be able to unset reserved protocol address")

	// Verify address is no longer reserved
	isReserved = suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Bob)
	suite.Require().False(isReserved, "address should not be reserved after unsetting")
}

// TestSetReservedProtocolAddress_NonGovernanceRejected tests that non-governance authority is rejected
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_NonGovernanceRejected() {
	// Try using a non-authority address
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          suite.Alice, // Not the governance authority
		Address:            suite.Bob,
		IsReservedProtocol: true,
	}

	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-governance authority should not be able to set reserved protocol address")
	suite.Require().Contains(err.Error(), "invalid authority", "error should mention invalid authority")

	// Verify address is NOT reserved
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Bob)
	suite.Require().False(isReserved, "address should not be reserved when set by non-authority")
}

// TestSetReservedProtocolAddress_AffectsApprovalCriteria tests that reserved address affects approval criteria checks
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_AffectsApprovalCriteria() {
	authority := suite.Keeper.GetAuthority()

	// Create a collection with forceful transfer approval
	approval := testutil.GenerateCollectionApproval("forceful_transfer", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true, // This triggers reserved protocol check
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Set up mint approval and mint tokens to Alice
	suite.SetupMintApproval(collectionId)
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintTokens(collectionId, suite.Alice, mintBalances)

	// Mark Alice as reserved protocol address
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Alice,
		IsReservedProtocol: true,
	}
	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Try to forcefully transfer FROM Alice (reserved address) - should fail
	// Bob initiates the transfer from Alice's address
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Bob, // Bob initiates
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // From reserved protocol address
				ToAddresses: []string{suite.Bob},
				Balances: []*types.Balance{
					testutil.GenerateSimpleBalance(5, 1),
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "forceful_transfer",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "forceful transfer from reserved protocol address should fail")
}

// TestSetReservedProtocolAddress_QueryReservedStatus tests querying reserved status after setting
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_QueryReservedStatus() {
	authority := suite.Keeper.GetAuthority()

	// Initially, Alice should not be reserved
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Alice)
	suite.Require().False(isReserved, "Alice should not be reserved initially")

	// Set Alice as reserved
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Alice,
		IsReservedProtocol: true,
	}
	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Query again - should now be reserved
	isReserved = suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Alice)
	suite.Require().True(isReserved, "Alice should be reserved after setting")

	// Bob should still not be reserved
	isReserved = suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Bob)
	suite.Require().False(isReserved, "Bob should not be affected by Alice being reserved")
}

// TestSetReservedProtocolAddress_MultipleAddresses tests setting multiple addresses as reserved
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_MultipleAddresses() {
	authority := suite.Keeper.GetAuthority()

	addresses := []string{suite.Alice, suite.Bob, suite.Charlie}

	// Set all addresses as reserved
	for _, addr := range addresses {
		msg := &types.MsgSetReservedProtocolAddress{
			Authority:          authority,
			Address:            addr,
			IsReservedProtocol: true,
		}
		_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "should be able to set %s as reserved", addr)
	}

	// Verify all are reserved
	for _, addr := range addresses {
		isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, addr)
		suite.Require().True(isReserved, "address %s should be reserved", addr)
	}

	// Unset one address
	unsetMsg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            suite.Bob,
		IsReservedProtocol: false,
	}
	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), unsetMsg)
	suite.Require().NoError(err)

	// Verify Bob is no longer reserved, others still are
	suite.Require().True(suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Alice), "Alice should still be reserved")
	suite.Require().False(suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Bob), "Bob should not be reserved")
	suite.Require().True(suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Charlie), "Charlie should still be reserved")
}

// TestSetReservedProtocolAddress_InvalidAddress tests setting invalid address as reserved
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_InvalidAddress() {
	authority := suite.Keeper.GetAuthority()

	// Try to set invalid address as reserved
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          authority,
		Address:            "invalid_address",
		IsReservedProtocol: true,
	}

	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid address")
	suite.Require().Contains(err.Error(), "invalid address", "error should mention invalid address")
}

// TestSetReservedProtocolAddress_EmptyAuthority tests that empty authority is rejected
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_EmptyAuthority() {
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:          "",
		Address:            suite.Alice,
		IsReservedProtocol: true,
	}

	_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with empty authority")
}

// TestSetReservedProtocolAddress_IdempotentSet tests setting same address twice
func (suite *SetReservedProtocolAddressTestSuite) TestSetReservedProtocolAddress_IdempotentSet() {
	authority := suite.Keeper.GetAuthority()

	// Set Alice as reserved twice
	for i := 0; i < 2; i++ {
		msg := &types.MsgSetReservedProtocolAddress{
			Authority:          authority,
			Address:            suite.Alice,
			IsReservedProtocol: true,
		}
		_, err := suite.MsgServer.SetReservedProtocolAddress(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "setting same address as reserved should succeed (iteration %d)", i)
	}

	// Should still be reserved
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, suite.Alice)
	suite.Require().True(isReserved, "Alice should be reserved after setting twice")
}
