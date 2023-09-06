package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestInheritedBalanceTypes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].BalancesType = sdkmath.NewUint(3)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline = nil
	collectionsToCreate[0].InheritedCollectionId = sdkmath.NewUint(2)
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline = nil
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = nil
	// collectionsToCreate[0].DefaultUserPermissions = nil

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Equal(collection.BalancesType, "Inherited")
	suite.Require().Equal(collection.InheritedCollectionId, sdkmath.NewUint(2))

	collectionsToCreate = GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{{ Amount: sdkmath.NewUint(1), BadgeIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges() }}
	collectionsToCreate[0].BalancesType = sdkmath.NewUint(3)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline = nil
	collectionsToCreate[0].InheritedCollectionId = sdkmath.NewUint(2)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")

	collectionsToCreate = GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].BalancesType = sdkmath.NewUint(3)
	collectionsToCreate[0].InheritedCollectionId = sdkmath.NewUint(2)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")
}