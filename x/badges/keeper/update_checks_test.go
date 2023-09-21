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
				
				ApprovalDetails: []*types.ApprovalDetails{
					{
						RequireToEqualsInitiatedBy:             true,
						ApprovalTrackerId:                      "test",
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: 												&types.ApprovalAmounts{},
						OverridesFromApprovedOutgoingTransfers: true,
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	newApprovalDetails := collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails
	newApprovalDetails = append(newApprovalDetails, &types.ApprovalDetails{})

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
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
				ApprovalDetails: 												collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	newApprovalDetails = []*types.ApprovalDetails{}
	newApprovalDetails = append(newApprovalDetails, &types.ApprovalDetails{})
	newApprovalDetails = append(newApprovalDetails, collectionsToCreate[0].CollectionApprovedTransfers[1].ApprovalDetails...)


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
				ApprovalDetails: 												newApprovalDetails,
				AllowedCombinations: 										[]*types.IsCollectionTransferAllowed{{}},
			},
			{
				FromMappingId:                          "!" + alice,
				ToMappingId:                            "AllWithoutMint",
				InitiatedByMappingId:                   "AllWithoutMint",
				BadgeIds:                               GetFullUintRanges(),
				TransferTimes:                          GetFullUintRanges(),
				OwnershipTimes: 		 										GetFullUintRanges(),
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
				
				ApprovalDetails: []*types.ApprovalDetails{},
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
			
			ApprovalDetails: []*types.ApprovalDetails{
				{
					RequireToEqualsInitiatedBy:             true,
					ApprovalTrackerId:                      "test",
					MaxNumTransfers: 												&types.MaxNumTransfers{},
					ApprovalAmounts: 												&types.ApprovalAmounts{},
					OverridesFromApprovedOutgoingTransfers: true,
				},
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
				ApprovalDetails: 					[]*types.OutgoingApprovalDetails{
					{
						RequireToEqualsInitiatedBy: true,
						ApprovalTrackerId:                  "test",
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: &types.ApprovalAmounts{},
					},
				},
			},
		},
		ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
			{
				FromMappingId:        alice,
				InitiatedByMappingId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes: 		 GetFullUintRanges(),
				ApprovalDetails: []*types.IncomingApprovalDetails{
					{
						ApprovalTrackerId:                 "test",
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: 												&types.ApprovalAmounts{},
					},
				},
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
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						RequireToEqualsInitiatedBy: true,
						ApprovalTrackerId:                  "test",
						MaxNumTransfers: 												&types.MaxNumTransfers{},
						ApprovalAmounts: 												&types.ApprovalAmounts{},
					},
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
				ApprovalDetails: []*types.IncomingApprovalDetails{
					{
						RequireFromEqualsInitiatedBy: true,

						ApprovalTrackerId:                 "test",
						MaxNumTransfers: 												&types.MaxNumTransfers{
							
						},
						ApprovalAmounts: 												&types.ApprovalAmounts{
							
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

//TODO: Equality checks
