package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSetValidTokenIds_ValidateBasicCalled(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first using UpdateCollection
	createMsg := &types.MsgUniversalUpdateCollection{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(0), // New collection
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		UpdateValidTokenIds:         true,
		UpdateCollectionPermissions: true,
		CollectionPermissions:       &types.CollectionPermissions{},
	}

	res, err := UpdateCollectionWithRes(suite, wctx, createMsg)
	require.NoError(t, err)
	collectionId := res.CollectionId

	// Test that ValidateBasic is called by passing an invalid message
	// Invalid creator address should fail ValidateBasic
	invalidMsg := &types.MsgSetValidTokenIds{
		Creator:      "invalid-address", // Invalid address
		CollectionId: collectionId,
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{},
	}

	_, err = suite.msgServer.SetValidTokenIds(wctx, invalidMsg)
	require.Error(t, err, "Invalid message should fail validation")
	// Verify the error is from ValidateBasic - should contain validation error
	require.Contains(t, err.Error(), "invalid", "Error should indicate validation failure")
}
