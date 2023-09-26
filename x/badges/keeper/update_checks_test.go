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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes:  GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes:   GetFullUintRanges(),
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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes:   GetFullUintRanges(),
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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateContractAddress: []*types.TimedUpdatePermission{
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes:  GetOneUintRange(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{},
					},
				},
				{
					DefaultValues: &types.TimedUpdateDefaultValues{
						ForbiddenTimes: GetFullUintRanges(),
						TimelineTimes:  GetTwoUintRanges(),
					},
					Combinations: []*types.TimedUpdateCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes:   GetFullUintRanges(),
				ContractAddress: "0x123",
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ContractAddressTimeline: []*types.ContractAddressTimeline{
			{
				TimelineTimes:   GetOneUintRange(),
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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes:  GetFullUintRanges(),
						BadgeIds:       GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri:      "https://example.com",
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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri:      "https://example.com",
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
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						PermittedTimes: GetFullUintRanges(),
						TimelineTimes:  GetOneUintRange(),
						BadgeIds:       GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{},
					},
				},
				{
					DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
						ForbiddenTimes: GetFullUintRanges(),
						TimelineTimes:  GetTwoUintRanges(),
						BadgeIds:       GetFullUintRanges(),
					},
					Combinations: []*types.TimedUpdateWithBadgeIdsCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri:      "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
			{
				TimelineTimes: GetOneUintRange(),
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri:      "https://example.com",
						BadgeIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckCollectionApprovedTransferUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{
				{
					DefaultValues: &types.CollectionApprovedTransferDefaultValues{
						FromMappingId:        alice,
						ToMappingId:          "AllWithoutMint",
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						ApprovalTrackerId: 		"All",
						ChallengeTrackerId:	  "All",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 	GetFullUintRanges(),
					},

					Combinations: []*types.CollectionApprovedTransferCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalTrackerId:                      "test",
				ApprovalId: 													  "test",
				ApprovalDetails: &types.ApprovalDetails{
					RequireToEqualsInitiatedBy:             true,
					
					MaxNumTransfers: 												&types.MaxNumTransfers{},
					ApprovalAmounts: 												&types.ApprovalAmounts{},
					OverridesFromApprovedOutgoingTransfers: true,
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												&types.ApprovalDetails{},
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")


	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												&types.ApprovalDetails{},
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: []*types.IsCollectionTransferAllowed{{
						IsApproved: false,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",

				ApprovalDetails: &types.ApprovalDetails{},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))


	approvedTransfers := []*types.CollectionApprovedTransfer{
		{
			FromMappingId:                          bob,
			ToMappingId:                            "AllWithoutMint",
			InitiatedByMappingId:                   "AllWithoutMint",
			BadgeIds:                               GetFullUintRanges(),
			TransferTimes:                          GetFullUintRanges(),
			OwnershipTimes: 		 										GetFullUintRanges(),
			ApprovalId: "test",
			ApprovalTrackerId:                      "test",
			
			ApprovalDetails: &types.ApprovalDetails{
				
					RequireToEqualsInitiatedBy:             true,
					MaxNumTransfers: 												&types.MaxNumTransfers{},
					ApprovalAmounts: 												&types.ApprovalAmounts{},
					OverridesFromApprovedOutgoingTransfers: true,
				
			},
		},
	}
	approvedTransfers = append(approvedTransfers, collection.CollectionApprovedTransfers...)


	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: approvedTransfers,
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}


func (suite *TestSuite) TestCheckCollectionApprovedTransferUpdateApprovalTrackerIds() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{
				{
					DefaultValues: &types.CollectionApprovedTransferDefaultValues{
						FromMappingId:        alice,
						ToMappingId:          "AllWithoutMint",
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						ApprovalTrackerId: 		"All",
						ChallengeTrackerId:	  "All",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 	GetFullUintRanges(),
					},

					Combinations: []*types.CollectionApprovedTransferCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 													"test",
				ApprovalTrackerId: 										"something that is not the same",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 													"test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckCollectionApprovedTransferUpdateApprovalTrackerIdsSpecificIdLocked() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{
				{
					DefaultValues: &types.CollectionApprovedTransferDefaultValues{
						FromMappingId:        alice,
						ToMappingId:          "AllWithoutMint",
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						ApprovalTrackerId: 		"test",
						ChallengeTrackerId:	  "All",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 	GetFullUintRanges(),
					},

					Combinations: []*types.CollectionApprovedTransferCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 													"test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: false,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 													"test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: false,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 														"test",
				ApprovalTrackerId: 											"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: 													"test",
				ApprovalTrackerId: 										"asdffdafs",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: 										"test",
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}


func (suite *TestSuite) TestCheckUserApprovedTransferUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultApprovedIncomingTransfers = []*types.UserApprovedIncomingTransfer{}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfers = []*types.UserApprovedOutgoingTransfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateUserPermissions(suite, wctx, &types.MsgUpdateUserPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.UserPermissions{
			CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransferPermission{
				{
					DefaultValues: &types.UserApprovedOutgoingTransferDefaultValues{
						ToMappingId:          alice,
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 GetFullUintRanges(),
						ApprovalTrackerId: "All",
						ChallengeTrackerId: "All",
					},

					Combinations: []*types.UserApprovedOutgoingTransferCombination{
						{},
					},
				},
			},
			CanUpdateApprovedIncomingTransfers: []*types.UserApprovedIncomingTransferPermission{
				{
					DefaultValues: &types.UserApprovedIncomingTransferDefaultValues{
						FromMappingId:        alice,
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 GetFullUintRanges(),
						ApprovalTrackerId: "All",
						ChallengeTrackerId: "All",
					},
					Combinations: []*types.UserApprovedIncomingTransferCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateUserApprovedTransfers(suite, wctx, &types.MsgUpdateUserApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		UpdateApprovedOutgoingTransfers: true,
		UpdateApprovedIncomingTransfers: true,
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:                alice,
				InitiatedByMappingId:       "AllWithoutMint",
				BadgeIds:                   GetFullUintRanges(),
				TransferTimes:              GetFullUintRanges(),
				OwnershipTimes: 		 GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId:                  "test",
				ApprovalDetails: 					&types.OutgoingApprovalDetails{
					
						RequireToEqualsInitiatedBy: true,
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: &types.ApprovalAmounts{},
					
				},
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{{
					IsApproved: true,
				}},

			},
		},
		ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
			{
				FromMappingId:        alice,
				InitiatedByMappingId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes: 		 GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId:                 "test",
				AllowedCombinations: []*types.IsUserIncomingTransferAllowed{{
					IsApproved: true,
				}},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateUserApprovedTransfers(suite, wctx, &types.MsgUpdateUserApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		UpdateApprovedOutgoingTransfers: true,
		UpdateApprovedIncomingTransfers: true,
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:                bob,
				InitiatedByMappingId:       "AllWithoutMint",
				BadgeIds:                   GetFullUintRanges(),
				TransferTimes:              GetFullUintRanges(),
				OwnershipTimes: 		 GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId:                  "test",
				ApprovalDetails: &types.OutgoingApprovalDetails{
					
						RequireToEqualsInitiatedBy: true,
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: 												&types.ApprovalAmounts{},
					
				},
			},
		},
		ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
			{
				FromMappingId:                bob,
				InitiatedByMappingId:         "AllWithoutMint",
				BadgeIds:                     GetFullUintRanges(),
				TransferTimes:                GetFullUintRanges(),
				OwnershipTimes: 		 GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId:                 "test",
				ApprovalDetails: &types.IncomingApprovalDetails{
					
						RequireFromEqualsInitiatedBy: true,

						MaxNumTransfers: 												&types.MaxNumTransfers{
							
						},
						ApprovalAmounts: 												&types.ApprovalAmounts{
							
						},
					
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentBaseCaseNoSplit() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{
				{
					DefaultValues: &types.CollectionApprovedTransferDefaultValues{
						FromMappingId:        alice,
						ToMappingId:          "AllWithoutMint",
						ForbiddenTimes:       GetFullUintRanges(),
						InitiatedByMappingId: "AllWithoutMint",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 	GetFullUintRanges(),
						ApprovalTrackerId: 		"All",
						ChallengeTrackerId: 	"All",
					},
					Combinations: []*types.CollectionApprovedTransferCombination{
						{},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	// 	ApprovalDetails: &types.ApprovalDetails{
	// 		{
	// 			MerkleChallenges:                []*types.MerkleChallenge{},
	// 			ApprovalTrackerId:                 "test",
	// 			MaxNumTransfers: &types.MaxNumTransfers{
	// 				OverallMaxNumTransfers: sdkmath.NewUint(1000),
	// 			},
	// 			ApprovalAmounts: &types.ApprovalAmounts{
	// 				PerFromAddressApprovalAmount: sdkmath.NewUint(1), //potentially unlimited
	// 			},
	// 		},
	// 	},
	// }},

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalent() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	// 	ApprovalDetails: &types.ApprovalDetails{
	// 		{
	// 			MerkleChallenges:                []*types.MerkleChallenge{},
	// 			ApprovalTrackerId:                 "test",
	// 			MaxNumTransfers: &types.MaxNumTransfers{
	// 				OverallMaxNumTransfers: sdkmath.NewUint(1000),
	// 			},
	// 			ApprovalAmounts: &types.ApprovalAmounts{
	// 				PerFromAddressApprovalAmount: sdkmath.NewUint(1), //potentially unlimited
	// 			},
	// 		},
	// 	},
	// }},

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badges")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparateBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	// 	ApprovalDetails: &types.ApprovalDetails{
	// 		{
	// 			MerkleChallenges:                []*types.MerkleChallenge{},
	// 			ApprovalTrackerId:                 "test",
	// 			MaxNumTransfers: &types.MaxNumTransfers{
	// 				OverallMaxNumTransfers: sdkmath.NewUint(1000),
	// 			},
	// 			ApprovalAmounts: &types.ApprovalAmounts{
	// 				PerFromAddressApprovalAmount: sdkmath.NewUint(1), //potentially unlimited
	// 			},
	// 		},
	// 	},
	// }},

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparateBalancesTwoTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)


	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          "!" + bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	// 	ApprovalDetails: &types.ApprovalDetails{
	// 		{
	// 			MerkleChallenges:                []*types.MerkleChallenge{},
	// 			ApprovalTrackerId:                 "test",
	// 			MaxNumTransfers: &types.MaxNumTransfers{
	// 				OverallMaxNumTransfers: sdkmath.NewUint(1000),
	// 			},
	// 			ApprovalAmounts: &types.ApprovalAmounts{
	// 				PerFromAddressApprovalAmount: sdkmath.NewUint(1), //potentially unlimited
	// 			},
	// 		},
	// 	},
	// }},

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badges")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparatePredeterminedBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)
	newApprovalDetails.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy: sdk.NewUint(0),
			IncrementOwnershipTimesBy: sdk.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetBottomHalfUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetTopHalfUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	//Not exactly the predetermined balances, but the same number of transfers
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(2),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badges")
}

func (suite *TestSuite) TestSplitPredetrminedBalancesEquivalentButNotSameTransferBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)
	newApprovalDetails.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy: sdk.NewUint(0),
			IncrementOwnershipTimesBy: sdk.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetBottomHalfUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetTopHalfUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")
}

func (suite *TestSuite) TestGetMaxPossible() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultApprovedIncomingTransfers = []*types.UserApprovedIncomingTransfer{
		{
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "AllWithoutMint",
			TransferTimes:        GetFullUintRanges(),
			OwnershipTimes: 			GetFullUintRanges(),
			ApprovalId: "test",
			BadgeIds:             GetFullUintRanges(),
			AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
				{
					IsApproved: true,
				},
			},
		},
	}
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount: sdkmath.NewUint(20),
			BadgeIds: GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].Transfers[0].Balances[0].Amount = sdkmath.NewUint(20)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
			collectionsToCreate[0].CollectionApprovedTransfers[0],
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: "adsfhjals",
				ApprovalDetails: 												&types.ApprovalDetails{
					
						
						ApprovalAmounts: &types.ApprovalAmounts{
							OverallApprovalAmount: sdk.NewUint(10),
						},
						
						MaxNumTransfers: &types.MaxNumTransfers{},
					},
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
			{
				FromMappingId:                          bob,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalId: "test",
				ApprovalTrackerId: "adsfhjaladsfasdf",
				ApprovalDetails: 												&types.ApprovalDetails{
					
					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdk.NewUint(10),
					},
					MaxNumTransfers: &types.MaxNumTransfers{},
				},
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{
					IsApproved: true,
				}},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:        bob,
		CollectionId:  sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From: 			bob,
				ToAddresses: []string{alice},
				Balances: 	[]*types.Balance{
					{
						Amount: sdkmath.NewUint(20),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")
}

//TODO: Equality checks