package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting token: %s")

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: collection.CollectionApprovals,
	})
	suite.Require().Nil(err, "Error updating collection approvals")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection.CollectionApprovals,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error creating token: %s")
}

func (suite *TestSuite) TestNewBadgesNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating token: %s")
}

func (suite *TestSuite) TestNewBadgeBadgeNotExists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						TokenIds: []*types.UintRange{
							{
								Start: sdkmath.NewUint(2),
								End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
							},
						},
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating token: %s")
}

func (suite *TestSuite) TestNewBadgesNotAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
				{
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TokenIds:                  GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating token: %s")
}

func (suite *TestSuite) TestNewBadgesPermissionIsApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
				{
					PermanentlyPermittedTimes: GetFullUintRanges(),
					TokenIds:                  GetOneUintRange(),
				},
				{
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					TokenIds:                  GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error creating token: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetTwoUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
	})
	suite.Require().Error(err, "Error creating token: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetOneUintRange(),
			},
		},
	})
	suite.Require().Error(err, "Error creating token: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetOneUintRange(),
			},
		},
	})
	suite.Require().Nil(err, "Error creating token: %s")
}
