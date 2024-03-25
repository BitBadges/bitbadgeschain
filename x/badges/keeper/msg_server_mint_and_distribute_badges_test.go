package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
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
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestNewBadgesNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
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
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestNewBadgeBadgeNotExists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
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
						BadgeIds: []*types.UintRange{
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
	suite.Require().Error(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestNewBadgesNotAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanCreateMoreBadges: []*types.BalancesActionPermission{
				{
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					BadgeIds:                  GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				BadgeIds:       GetOneUintRange(),
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
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetOneUintRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestNewBadgesPermissionIsApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUniversalUpdateCollectionPermissions{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanCreateMoreBadges: []*types.BalancesActionPermission{
				{
					PermanentlyPermittedTimes: GetFullUintRanges(),
					BadgeIds:                  GetOneUintRange(),
					OwnershipTimes:            GetOneUintRange(),
				},
				{
					PermanentlyForbiddenTimes: GetFullUintRanges(),
					BadgeIds:                  GetFullUintRanges(),
					OwnershipTimes:            GetFullUintRanges(),
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
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetFullUintRanges(),
				OwnershipTimes: GetOneUintRange(),
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetOneUintRange(),
			},
		},
	})
	suite.Require().Nil(err, "Error creating badge: %s")
}
