package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestArchiveCollection() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	_, err = GetCollection(suite, wctx, sdk.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")

	err = ArchiveCollection(suite, wctx, &types.MsgArchiveCollection{
		Creator:       bob,
		CollectionId:  sdk.NewUint(1),
		IsArchivedTimeline: []*types.IsArchivedTimeline{
			{
				IsArchived: true,
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error archiving collection: %s")

	//Still should be able to get collection
	_, err = GetCollection(suite, wctx, sdk.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")

	err = UpdateManager(suite, wctx, &types.MsgUpdateManager{
		Creator:       bob,
		CollectionId:  sdk.NewUint(1),
		ManagerTimeline:       []*types.ManagerTimeline{
			{
				Manager: alice,
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Error(err, "Error updating manager: %s")

	err = ArchiveCollection(suite, wctx, &types.MsgArchiveCollection{
		Creator:       bob,
		CollectionId:  sdk.NewUint(1),
		IsArchivedTimeline: []*types.IsArchivedTimeline{
			{
				IsArchived: false,
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error archiving collection: %s")

	err = UpdateManager(suite, wctx, &types.MsgUpdateManager{
		Creator:       bob,
		CollectionId:  sdk.NewUint(1),
		ManagerTimeline:       []*types.ManagerTimeline{
			{
				Manager: alice,
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error updating manager: %s")
}