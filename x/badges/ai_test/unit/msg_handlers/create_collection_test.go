package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type CreateCollectionTestSuite struct {
	testutil.AITestSuite
}

func TestCreateCollectionSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CreateCollectionTestSuite))
}

// TestCreateCollection_ValidInput tests creating a collection with valid inputs
func (suite *CreateCollectionTestSuite) TestCreateCollection_ValidInput() {
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			testutil.GenerateUintRange(1, 100),
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: suite.Manager,
		CollectionMetadata: testutil.GenerateCollectionMetadata("https://example.com/metadata", ""),
		TokenMetadata: []*types.TokenMetadata{},
		CustomData: "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards: []string{},
		IsArchived: false,
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().True(resp.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0")

	// Verify collection exists
	collection := suite.GetCollection(resp.CollectionId)
	suite.Require().Equal(suite.Manager, collection.Manager)
	suite.Require().Equal("https://example.com/metadata", collection.CollectionMetadata.Uri)
}

// TestCreateCollection_InvalidCreator tests creating a collection with invalid creator
func (suite *CreateCollectionTestSuite) TestCreateCollection_InvalidCreator() {
	msg := &types.MsgCreateCollection{
		Creator: "", // Invalid empty creator
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			testutil.GenerateUintRange(1, 100),
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: "",
		CollectionMetadata: testutil.GenerateCollectionMetadata("", ""),
		TokenMetadata: []*types.TokenMetadata{},
		CustomData: "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards: []string{},
		IsArchived: false,
	}

	_, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid creator")
}

// TestCreateCollection_InvalidTokenIds tests creating a collection with invalid token ID ranges
func (suite *CreateCollectionTestSuite) TestCreateCollection_InvalidTokenIds() {
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(1)}, // Invalid: start > end
		},
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
	suite.Require().Error(err, "should fail with invalid token ID ranges")
}

// TestCreateCollection_WithApprovals tests creating a collection with approval settings
func (suite *CreateCollectionTestSuite) TestCreateCollection_WithApprovals() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "All", "All"),
	}

	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			testutil.GenerateUintRange(1, 100),
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: suite.Manager,
		CollectionMetadata: testutil.GenerateCollectionMetadata("", ""),
		TokenMetadata: []*types.TokenMetadata{},
		CustomData: "",
		CollectionApprovals: approvals,
		Standards: []string{},
		IsArchived: false,
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	// Verify approvals are set
	collection := suite.GetCollection(resp.CollectionId)
	suite.Require().Equal(1, len(collection.CollectionApprovals))
	suite.Require().Equal("approval1", collection.CollectionApprovals[0].ApprovalId)
}

// TestCreateCollection_ManagerPermission tests that creator becomes manager
func (suite *CreateCollectionTestSuite) TestCreateCollection_ManagerPermission() {
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{}, // Empty balances - zero amounts are not allowed
		},
		ValidTokenIds: []*types.UintRange{
			testutil.GenerateUintRange(1, 100),
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager: suite.Manager,
		CollectionMetadata: testutil.GenerateCollectionMetadata("", ""),
		TokenMetadata: []*types.TokenMetadata{},
		CustomData: "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards: []string{},
		IsArchived: false,
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	collection := suite.GetCollection(resp.CollectionId)
	suite.Require().Equal(suite.Manager, collection.Manager, "creator should be set as manager")
	suite.Require().Equal(suite.Manager, collection.CreatedBy, "createdBy should match creator")
}

