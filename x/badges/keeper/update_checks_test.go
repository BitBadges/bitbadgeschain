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
						TimelineTimes:        GetFullUintRanges(),
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
		CollectionApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
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
								ApprovalId:                              "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{},
								ApprovalAmounts: 												&types.ApprovalAmounts{},
								OverridesFromApprovedOutgoingTransfers: true,
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				CollectionApprovedTransfers: []*types.CollectionApprovedTransfer{
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
								ApprovalId:                             "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{},
								ApprovalAmounts: 												&types.ApprovalAmounts{},
								OverridesFromApprovedOutgoingTransfers: true,
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckUserApprovedTransferUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline = []*types.UserApprovedIncomingTransferTimeline{}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = []*types.UserApprovedOutgoingTransferTimeline{}

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
						TimelineTimes:        GetFullUintRanges(),
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
						TimelineTimes:        GetFullUintRanges(),
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
		UpdateApprovedOutgoingTransfersTimeline: true,
		UpdateApprovedIncomingTransfersTimeline: true,
		ApprovedOutgoingTransfersTimeline: []*types.UserApprovedOutgoingTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
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
								ApprovalId:                  "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{},
								ApprovalAmounts: &types.ApprovalAmounts{},
							},
						},
					},
				},
			},
		},
		ApprovedIncomingTransfersTimeline: []*types.UserApprovedIncomingTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
					{
						FromMappingId:        alice,
						InitiatedByMappingId: "AllWithoutMint",
						BadgeIds:             GetFullUintRanges(),
						TransferTimes:        GetFullUintRanges(),
						OwnershipTimes: 		 GetFullUintRanges(),
						ApprovalDetails: []*types.IncomingApprovalDetails{
							{
								ApprovalId:                 "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{},
								ApprovalAmounts: 												&types.ApprovalAmounts{},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateUserApprovedTransfers(suite, wctx, &types.MsgUpdateUserApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		UpdateApprovedOutgoingTransfersTimeline: true,
		UpdateApprovedIncomingTransfersTimeline: true,
		ApprovedOutgoingTransfersTimeline: []*types.UserApprovedOutgoingTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
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
								ApprovalId:                  "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{},
								ApprovalAmounts: 												&types.ApprovalAmounts{},
							},
						},
					},
				},
			},
		},
		ApprovedIncomingTransfersTimeline: []*types.UserApprovedIncomingTransferTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
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

								ApprovalId:                 "test",
								MaxNumTransfers: 												&types.MaxNumTransfers{
									
								},
								ApprovalAmounts: 												&types.ApprovalAmounts{
									
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

//TODO: Equality checks
