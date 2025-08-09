package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
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
					TimelineTimes:             GetFullUintRanges(),
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
					TimelineTimes:             GetOneUintRange(),
				},
				{

					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TimelineTimes:             GetTwoUintRanges(),
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

func (suite *TestSuite) TestCheckIfTimedUpdateWithTokenIdsPermissionPermits() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TimedUpdateWithTokenIdsPermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:             GetFullUintRanges(),
					TokenIds:                  GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		TokenMetadataTimeline: []*types.TokenMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri:      "https://example.com",
						TokenIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckIfTimedUpdateWithTokenIdsPermissionPermitsDefaultAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TimedUpdateWithTokenIdsPermission{},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		TokenMetadataTimeline: []*types.TokenMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri:      "https://example.com",
						TokenIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating metadata")
}

func (suite *TestSuite) TestCheckIfTimedUpdateWithTokenIdsPermissionPermitsInvalidTimes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TimedUpdateWithTokenIdsPermission{
				{

					PermanentlyPermittedTimes: GetFullUintRanges(),
					TimelineTimes:             GetOneUintRange(),
					TokenIds:                  GetFullUintRanges(),
				},
				{

					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TimelineTimes:             GetTwoUintRanges(),
					TokenIds:                  GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		TokenMetadataTimeline: []*types.TokenMetadataTimeline{
			{
				TimelineTimes: GetFullUintRanges(),
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri:      "https://example.com",
						TokenIds: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error updating metadata")

	err = UpdateMetadata(suite, wctx, &types.MsgUpdateMetadata{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		TokenMetadataTimeline: []*types.TokenMetadataTimeline{
			{
				TimelineTimes: GetOneUintRange(),
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri:      "https://example.com",
						TokenIds: GetFullUintRanges(),
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

					FromListId:                alice,
					ToListId:                  "AllWithoutMint",
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					ApprovalId:                "All",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testasdfas",
				ApprovalCriteria:  &types.ApprovalCriteria{},
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testasdfasdfasfd",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  &types.ApprovalCriteria{},
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testafdasdf",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testasdfasdf",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testdfgh",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "something different",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testhdfgjhdf",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",

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
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "test2",

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

					FromListId:                alice,
					ToListId:                  "AllWithoutMint",
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					ApprovalId:                "All",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "something that is not the same",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "tesfasdft",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testadsfasdf",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",

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

					FromListId:                alice,
					ToListId:                  "AllWithoutMint",
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					ApprovalId:                "test",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				InitiatedByListId: "Mint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",

				ApprovalCriteria: collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testafdsasdf",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetOneUintRange(), //different
				ApprovalId:        "test",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testasdfas",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "tesadsft",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
			},
			{
				FromListId:        "!" + alice,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testasdfasd",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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

					FromListId:                alice,
					ToListId:                  "AllWithoutMint",
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					ApprovalId:                "approvalidtotest",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetOneUintRange(), //different
				ApprovalId:        "approvalidtotest",

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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "different id",
				ApprovalCriteria:  collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria,
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

					ToListId:                  alice,
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
					ApprovalId:                "All",
				},
			},
			CanUpdateIncomingApprovals: []*types.UserIncomingApprovalPermission{
				{

					FromListId:                alice,
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
					ApprovalId:                "All",
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
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

					FromListId:                alice,
					ToListId:                  "AllWithoutMint",
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					InitiatedByListId:         "AllWithoutMint",
					TokenIds:                  GetFullUintRanges(),
					TransferTimes:             GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
					ApprovalId:                "All",
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalent() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testgfag",
				ApprovalCriteria:  newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring tokens")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparateBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testagdf",
				ApprovalCriteria:  newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparateBalancesTwoTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			collectionsToCreate[0].CollectionApprovals[0],
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  newApprovalCriteria,
			},
			{
				FromListId:        "!" + bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testadfgsd",
				ApprovalCriteria:  newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring tokens")
}

func (suite *TestSuite) TestSplittingIntoMultipleIsEquivalentSeparatePredeterminedBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	newApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
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
				TokenIds:          GetBottomHalfUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  newApprovalCriteria,
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetTopHalfUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testsgdfs",
				ApprovalCriteria:  newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")

	//Not exactly the predetermined balances, but the same number of transfers
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring tokens")
}

func (suite *TestSuite) TestSplitPredetrminedBalancesEquivalentButNotSameTransferBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	newApprovalCriteria := collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria
	newApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	newApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
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
				TokenIds:          GetBottomHalfUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria:  newApprovalCriteria,
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetTopHalfUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "testsdfgsdf",
				ApprovalCriteria:  newApprovalCriteria,
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")
}

func (suite *TestSuite) TestGetMaxPossible() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals = []*types.UserIncomingApproval{
		{
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "test",
			TokenIds:          GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(20),
			TokenIds:       GetFullUintRanges(),
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
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "test",
				ApprovalCriteria: &types.ApprovalCriteria{
					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdkmath.NewUint(10),
					},

					MaxNumTransfers: &types.MaxNumTransfers{},
				},
			},
			{
				FromListId:        bob,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TokenIds:          GetFullUintRanges(),
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				ApprovalId:        "tessdgfst",
				ApprovalCriteria: &types.ApprovalCriteria{

					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdkmath.NewUint(10),
					},
					MaxNumTransfers: &types.MaxNumTransfers{},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection approved transfers")

	//Test that the number of balances does not matter as long as they are equivalent
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(20),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring tokens")
}

//TODO: Equality checks
