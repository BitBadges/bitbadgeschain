package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestCheckIfTimedUpdatePermissionPermits() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateStandards: []*types.TimedUpdatePermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:  GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		StandardsTimeline: []*types.StandardsTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				Standards:     []string{"0x123"},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckIfTimedUpdatePermissionPermitsDefaultAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateStandards: []*types.TimedUpdatePermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		StandardsTimeline: []*types.StandardsTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				Standards:     []string{"0x123"},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckIfTimedUpdatePermissionPermitsInvalidTimes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateStandards: []*types.TimedUpdatePermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:  GetOneUintRange(),
				},
				{

					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TimelineTimes:  GetTwoUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		StandardsTimeline: []*types.StandardsTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				Standards:     []string{"0x123"},
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		StandardsTimeline: []*types.StandardsTimeline{
			{
				TimelineTimes: GetOneUintRange(),
				Standards:     []string{"0x123"},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckIfTimedUpdateWithBadgeIdsPermissionPermits() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:  GetFullUintRanges(),
					BadgeIds:       GetFullUintRanges(),
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

func (suite *TestSuite) TestCheckIfTimedUpdateWithBadgeIdsPermissionPermitsDefaultAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
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

func (suite *TestSuite) TestCheckIfTimedUpdateWithBadgeIdsPermissionPermitsInvalidTimes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:  GetOneUintRange(),
					BadgeIds:       GetFullUintRanges(),
				},
				{

					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TimelineTimes:  GetTwoUintRanges(),
					BadgeIds:       GetFullUintRanges(),
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

func (suite *TestSuite) TestCheckCollectionApprovalUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{

					FromListId:        alice,
					ToListId:          "AllWithoutMint",
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					ApprovalId: "All",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ChallengeTrackerId:   "test",
				AmountTrackerId:      "test",
				ApprovalId:           "test",
				ApprovalCriteria: &types.ApprovalCriteria{
					RequireToEqualsInitiatedBy: true,

					MaxNumTransfers:                &types.MaxNumTransfers{},
					ApprovalAmounts:                &types.ApprovalAmounts{},
					OverridesFromOutgoingApprovals: true,
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testasdfas",
				ApprovalCriteria:     &types.ApprovalCriteria{},
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testasdfasdfasfd",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ApprovalCriteria:     &types.ApprovalCriteria{},
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testafdasdf",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testasdfasdf",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testdfgh",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
				ChallengeTrackerId:   "something different",
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testhdfgjhdf",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
				ChallengeTrackerId:   "test",
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",

				ApprovalCriteria: &types.ApprovalCriteria{},
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	approvals := []*types.CollectionApproval{
		{
			FromListId:        bob,
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			BadgeIds:             GetFullUintRanges(),
			TransferTimes:        GetFullUintRanges(),
			OwnershipTimes:       GetFullUintRanges(),
			ApprovalId:           "test2",
			ChallengeTrackerId:   "test2",
			AmountTrackerId:      "test2",

			ApprovalCriteria: &types.ApprovalCriteria{

				RequireToEqualsInitiatedBy:     true,
				MaxNumTransfers:                &types.MaxNumTransfers{},
				ApprovalAmounts:                &types.ApprovalAmounts{},
				OverridesFromOutgoingApprovals: true,
			},
		},
	}
	approvals = append(approvals, collection.CollectionApprovals...)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckCollectionApprovalUpdateAmountTrackerIds() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{

					FromListId:        alice,
					ToListId:          "AllWithoutMint",
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					ApprovalId: "All",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "something that is not the same",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "tesfasdft",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testadsfasdf",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckCollectionApprovalUpdateAmountTrackerIdsSpecificIdLocked() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",

				ApprovalCriteria: collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")


	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{

					FromListId:        alice,
					ToListId:          "AllWithoutMint",
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					ApprovalId:      		"test",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "something different",

				ApprovalCriteria: collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testafdsasdf",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
				ChallengeTrackerId:   "something different",
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testasdfas",
				AmountTrackerId:      "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
				ChallengeTrackerId:   "test",
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "tesadsft",
				AmountTrackerId:      "asdffdafs",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testasdfasd",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckCollectionApprovalUpdateAmountTrackerIdsSpecificIdLockedNoPreviousValues() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")
	
	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{

					FromListId:        alice,
					ToListId:          "AllWithoutMint",
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					ApprovalId:      		"approvalidtotest",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "approvalidtotest",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "something different",

				ApprovalCriteria: collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "different id",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "something different",

				ApprovalCriteria: collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")
}

func (suite *TestSuite) TestCheckUserApprovalUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals = []*types.UserIncomingApproval{}
	collectionsToCreate[0].DefaultOutgoingApprovals = []*types.UserOutgoingApproval{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateUserPermissions(suite, wctx, &types.MsgUpdateUserPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.UserPermissions{
			CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{
				{

					ToListId:          alice,
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
					ApprovalId: "All",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
				},
			},
			CanUpdateIncomingApprovals: []*types.UserIncomingApprovalPermission{
				{

					FromListId:        alice,
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
					ApprovalId: "All",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		UpdateIncomingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          alice,
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ChallengeTrackerId:   "test",
				AmountTrackerId:      "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{

					RequireToEqualsInitiatedBy: true,
					MaxNumTransfers:            &types.MaxNumTransfers{},
					ApprovalAmounts:            &types.ApprovalAmounts{},
				},
			},
		},
		IncomingApprovals: []*types.UserIncomingApproval{
			{
				FromListId:        alice,
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ChallengeTrackerId:   "test",
				AmountTrackerId:      "test",
			},
		},
	})
	suite.Require().Error(err, "Error updating collection approved transfers")

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		UpdateIncomingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          bob,
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ChallengeTrackerId:   "test",
				AmountTrackerId:      "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{

					RequireToEqualsInitiatedBy: true,
					MaxNumTransfers:            &types.MaxNumTransfers{},
					ApprovalAmounts:            &types.ApprovalAmounts{},
				},
			},
		},
		IncomingApprovals: []*types.UserIncomingApproval{
			{
				FromListId:        bob,
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				ChallengeTrackerId:   "test",
				AmountTrackerId:      "test",
				ApprovalCriteria: &types.IncomingApprovalCriteria{

					RequireFromEqualsInitiatedBy: true,

					MaxNumTransfers: &types.MaxNumTransfers{},
					ApprovalAmounts: &types.ApprovalAmounts{},
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

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{

					FromListId:        alice,
					ToListId:          "AllWithoutMint",
					PermanentlyForbiddenTimes:       GetFullUintRanges(),
					InitiatedByListId: "AllWithoutMint",
					BadgeIds:             GetFullUintRanges(),
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes:       GetFullUintRanges(),
					ApprovalId: "All",
					AmountTrackerId: "All",
					ChallengeTrackerId: "All",
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
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

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testgfag",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
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

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testagdf",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTopHalfUintRanges(),
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

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testadfgsd",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTopHalfUintRanges(),
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

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)
	newApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy:       sdk.NewUint(0),
			IncrementOwnershipTimesBy: sdk.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetBottomHalfUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetTopHalfUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testsgdfs",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")

	//Not exactly the predetermined balances, but the same number of transfers
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2),
						BadgeIds:       GetFullUintRanges(),
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

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdk.NewUint(1)
	newApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy:       sdk.NewUint(0),
			IncrementOwnershipTimesBy: sdk.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetBottomHalfUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetTopHalfUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "testsdfgsdf",
				AmountTrackerId:      "test",
				ChallengeTrackerId:   "test",
				ApprovalCriteria:     newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
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
	collectionsToCreate[0].DefaultIncomingApprovals = []*types.UserIncomingApproval{
		{
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:        GetFullUintRanges(),
			OwnershipTimes:       GetFullUintRanges(),
			ApprovalId:           "test",
			AmountTrackerId:      "test",
			ChallengeTrackerId:   "test",
			BadgeIds:             GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(20),
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].Transfers[0].Balances[0].Amount = sdkmath.NewUint(20)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "test",
				AmountTrackerId:      "adsfhjals",
				ChallengeTrackerId:   "test",
				ApprovalCriteria: &types.ApprovalCriteria{

					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdk.NewUint(10),
					},

					MaxNumTransfers: &types.MaxNumTransfers{},
				},
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				BadgeIds:             GetFullUintRanges(),
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes:       GetFullUintRanges(),
				ApprovalId:           "tessdgfst",
				AmountTrackerId:      "adsfhjaladsfasdf",
				ChallengeTrackerId:   "test",
				ApprovalCriteria: &types.ApprovalCriteria{

					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdk.NewUint(10),
					},
					MaxNumTransfers: &types.MaxNumTransfers{},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(20),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badges")
}

//TODO: Equality checks
