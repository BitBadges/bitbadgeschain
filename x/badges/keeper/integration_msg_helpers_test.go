package keeper_test

import (
	"context"
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* Query helpers */

func GetCollection(suite *TestSuite, ctx context.Context, id sdk.Uint) (types.BadgeCollection, error) {
	res, err := suite.app.BadgesKeeper.GetCollection(ctx, &types.QueryGetCollectionRequest{CollectionId: sdk.Uint(id)})
	if err != nil {
		return types.BadgeCollection{}, err
	}

	return *res.Collection, nil
}

func GetUserBalance(suite *TestSuite, ctx context.Context, collectionId sdk.Uint, address string) (types.UserBalanceStore, error) {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		CollectionId: sdk.Uint(collectionId),
		Address: address,
	})
	if err != nil {
		return types.UserBalanceStore{}, err
	}

	return *res.Balance, nil
}

func GetClaim(suite *TestSuite, ctx context.Context, collectionId sdk.Uint, claimId sdk.Uint) (types.Claim, error) {
	res, err := suite.app.BadgesKeeper.GetClaim(ctx, &types.QueryGetClaimRequest{
		CollectionId: sdk.Uint(collectionId),
		ClaimId: sdk.Uint(claimId),
	})
	if err != nil {
		return types.Claim{}, err
	}

	return *res.Claim, nil
}

/* Msg helpers */

type CollectionsToCreate struct {
	Collection types.MsgNewCollection
	Amount     sdk.Uint
	Creator    string
}

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []CollectionsToCreate) error {
	for _, collectionToCreate := range collectionsToCreate {
		for i := 0; i < int(collectionToCreate.Amount.BigInt().Int64()); i++ {
			msg := types.NewMsgNewCollection(collectionToCreate.Creator, collectionToCreate.Collection.Standard, collectionToCreate.Collection.BadgeSupplys, collectionToCreate.Collection.CollectionUri, collectionToCreate.Collection.BadgeUris, collectionToCreate.Collection.Permissions, collectionToCreate.Collection.AllowedTransfers, collectionToCreate.Collection.ManagerApprovedTransfers, collectionToCreate.Collection.Bytes, collectionToCreate.Collection.Transfers, collectionToCreate.Collection.Claims, collectionToCreate.Collection.BalancesUri)
			_, err := suite.msgServer.NewCollection(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateBadges(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, supplysAndAmounts []*types.BadgeSupplyAndAmount, transfers []*types.Transfer, claims []*types.Claim, collectionUri string, badgeUris []*types.BadgeUri, balancesUri string) error {
	msg := types.NewMsgMintAndDistributeBadges(creator, collectionId, supplysAndAmounts, transfers, claims, collectionUri, badgeUris, balancesUri)
	_, err := suite.msgServer.MintAndDistributeBadges(ctx, msg)
	return err
}

// Note: Only supports Bob and Alice and only supports supplysAndAmounts[0]
func CreateBadgesAndMintAllToCreator(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, supplysAndAmounts []*types.BadgeSupplyAndAmount) error {
	collection, err := GetCollection(suite, ctx, collectionId)
	if err != nil {
		return err
	}

	transfers := []*types.Transfer{
		{
			ToAddresses: []string{creator},
			Balances: []*types.Balance{
				{
					Amount: supplysAndAmounts[0].Supply,
					BadgeIds: []*types.IdRange{
						{
							Start: collection.NextBadgeId,
							End:   collection.NextBadgeId.Add(supplysAndAmounts[0].Amount).SubUint64(1),
						},
					},
				},
			},
		},
	}

	msg := types.NewMsgMintAndDistributeBadges(creator, collectionId, supplysAndAmounts, transfers, []*types.Claim{}, "https://example.com",
		[]*types.BadgeUri{
			{
				Uri: "https://example.com/{id}",
				BadgeIds: []*types.IdRange{
					{
						Start: sdk.NewUint(1),
						End:   sdk.NewUint(math.MaxUint64),
					},
				},
			},
		}, "")
	_, err = suite.msgServer.MintAndDistributeBadges(ctx, msg)
	return err
}

func TransferBadge(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, from string, transfers []*types.Transfer) error {
	msg := types.NewMsgTransferBadge(creator, collectionId, from, transfers)
	_, err := suite.msgServer.TransferBadge(ctx, msg)
	return err
}

func SetApproval(suite *TestSuite, ctx context.Context, creator string, address string, collectionId sdk.Uint, balances []*types.Balance) error {
	msg := types.NewMsgSetApproval(creator, collectionId, address, balances)
	_, err := suite.msgServer.SetApproval(ctx, msg)
	return err
}

func UpdateAllowedTransfers(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, allowedTransfers []*types.TransferMapping) error {
	msg := types.NewMsgUpdateAllowedTransfers(creator, collectionId, allowedTransfers)
	_, err := suite.msgServer.UpdateAllowedTransfers(ctx, msg)
	return err
}

func RequestTransferManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, add bool) error {
	msg := types.NewMsgRequestTransferManager(creator, collectionId, add)
	_, err := suite.msgServer.RequestTransferManager(ctx, msg)
	return err
}

func TransferManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, address string) error {
	msg := types.NewMsgTransferManager(creator, collectionId, address)
	_, err := suite.msgServer.TransferManager(ctx, msg)
	return err
}

func UpdateURIs(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, collectionUri string, badgeUris []*types.BadgeUri, balancesUri string) error {
	msg := types.NewMsgUpdateUris(creator, collectionId, collectionUri, badgeUris, balancesUri)
	_, err := suite.msgServer.UpdateUris(ctx, msg)
	return err
}

func UpdatePermissions(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, permissions sdk.Uint) error {
	msg := types.NewMsgUpdatePermissions(creator, collectionId, permissions)
	_, err := suite.msgServer.UpdatePermissions(ctx, msg)
	return err
}

func UpdateBytes(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, bytes string) error {
	msg := types.NewMsgUpdateBytes(creator, collectionId, bytes)
	_, err := suite.msgServer.UpdateBytes(ctx, msg)
	return err
}

func ClaimBadge(suite *TestSuite, ctx context.Context, creator string, claimId sdk.Uint, collectionId sdk.Uint, solutions []*types.ChallengeSolution) error {
	msg := types.NewMsgClaimBadge(creator, claimId, collectionId, solutions)
	_, err := suite.msgServer.ClaimBadge(ctx, msg)
	return err
}
