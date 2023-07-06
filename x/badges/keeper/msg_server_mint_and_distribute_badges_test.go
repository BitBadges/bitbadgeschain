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

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetFullIdRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From: bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						OwnershipTimes: GetOneIdRange(),
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
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetFullIdRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From: bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						OwnershipTimes: GetOneIdRange(),
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
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetFullIdRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From: bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: []*types.IdRange{
							{
								Start: sdkmath.NewUint(2),
								End: sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
							},
						},
						OwnershipTimes: GetOneIdRange(),
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

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanCreateMoreBadges: []*types.BalancesActionPermission{
				{
					DefaultValues: &types.BalancesActionDefaultValues{
						ForbiddenTimes: GetFullIdRanges(),
						BadgeIds: GetFullIdRanges(),
						OwnershipTimes: GetFullIdRanges(),
					},
					Combinations: []*types.BalancesActionCombination{{
					}},
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
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetFullIdRanges(),
			},
		},
		Transfers: []*types.Transfer{
			{
				From: bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						OwnershipTimes: GetOneIdRange(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestNewBadgesPermissionIsAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateCollectionPermissions(suite, wctx, &types.MsgUpdateCollectionPermissions{
		Creator: bob,
		CollectionId: sdkmath.NewUint(1),
		Permissions: &types.CollectionPermissions{
			CanCreateMoreBadges: []*types.BalancesActionPermission{
				{
					DefaultValues: &types.BalancesActionDefaultValues{
						PermittedTimes: GetFullIdRanges(),
						BadgeIds: GetOneIdRange(),
						OwnershipTimes: GetOneIdRange(),
					},
					Combinations: []*types.BalancesActionCombination{{

					}},
				},
				{
					DefaultValues: &types.BalancesActionDefaultValues{
						ForbiddenTimes: GetFullIdRanges(),
						BadgeIds: GetFullIdRanges(),
						OwnershipTimes: GetFullIdRanges(),
					},
					Combinations: []*types.BalancesActionCombination{{
						
					}},
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
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetFullIdRanges(),
				OwnershipTimes: GetOneIdRange(),
			},
		},
	})
	suite.Require().Error(err, "Error creating badge: %s")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				OwnershipTimes: GetOneIdRange(),
			},
		},
	})
	suite.Require().Nil(err, "Error creating badge: %s")
}