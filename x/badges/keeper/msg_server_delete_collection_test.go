package keeper_test

import (
	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestDeleteCollection() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	_, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")

	err = DeleteCollection(suite, wctx, &types.MsgDeleteCollection{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error archiving collection: %s")

	//Still should be able to get collection
	_, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Error(err, "Error getting badge: %s")

	err = ArchiveCollection(suite, wctx, &types.MsgArchiveCollection{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		IsArchivedTimeline: []*types.IsArchivedTimeline{
			{
				IsArchived:    true,
				TimelineTimes: GetFullUintRanges(),
			},
		},
	})
	suite.Require().Error(err, "Error getting badge: %s")
}
