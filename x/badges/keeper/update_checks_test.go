package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestCheckTimedUpdatePermission() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes: GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{
							
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				ContractAddress: "0x123",
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckTimedUpdatePermissionDefaultAllowed() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				ContractAddress: "0x123",
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckTimedUpdatePermissionInvalidTimes() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes: GetOneUintRange(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{
							
						},
					},
				},
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						ForbiddenTimes: GetFullUintRanges(),
						TimelineTimes: GetTwoUintRanges(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{
							
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				ContractAddress: "0x123",
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes: GetOneUintRange(),
				ContractAddress: "0x123",
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckTimedUpdateWithBadgeIdsPermission() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes: GetFullUintRanges(),
						BadgeIds: GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{
							
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckTimedUpdateWithBadgeIdsPermissionDefaultAllowed() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckTimedUpdateWithBadgeIdsPermissionInvalidTimes() {	
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes: GetOneUintRange(),
						BadgeIds: GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{
							
						},
					},
				},
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						ForbiddenTimes: GetFullUintRanges(),
						TimelineTimes: GetTwoUintRanges(),
						BadgeIds: GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{
							
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetOneUintRange(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}


//TODO: Collection/UserApprovedTransfer Updates check