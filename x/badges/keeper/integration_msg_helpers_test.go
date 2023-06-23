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
		Address:      address,
	})
	if err != nil {
		return types.UserBalanceStore{}, err
	}

	return *res.Balance, nil
}

func GetClaim(suite *TestSuite, ctx context.Context, collectionId sdk.Uint, claimId sdk.Uint) (types.Claim, error) {
	res, err := suite.app.BadgesKeeper.GetClaim(ctx, &types.QueryGetClaimRequest{
		CollectionId: sdk.Uint(collectionId),
		ClaimId:      sdk.Uint(claimId),
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
			msg := types.NewMsgNewCollection(collectionToCreate.Creator, collectionToCreate.Collection.Standard, collectionToCreate.Collection.BadgesToCreate, collectionToCreate.Collection.CollectionMetadata, collectionToCreate.Collection.BadgeMetadata, collectionToCreate.Collection.Permissions, collectionToCreate.Collection.ApprovedTransfers, collectionToCreate.Collection.ManagerApprovedTransfers, collectionToCreate.Collection.Bytes, collectionToCreate.Collection.Transfers, collectionToCreate.Collection.Claims, collectionToCreate.Collection.OffChainBalancesMetadata)
			_, err := suite.msgServer.NewCollection(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateBadges(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, supplysAndAmounts []*types.BadgeSupplyAndAmount, transfers []*types.Transfer, claims []*types.Claim, collectionMetadata *CollectionMetadata, badgeMetadata []*types.BadgeMetadata, offChainBalancesMetadata *OffChainBalancesMetadata) error {
	msg := types.NewMsgMintAndDistributeBadges(creator, collectionId, supplysAndAmounts, transfers, claims, collectionMetadata, badgeMetadata, offChainBalancesMetadata)
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
		[]*types.BadgeMetadata{
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

func UpdateCollectionApprovedTransfers(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, approvedTransfers []*types.CollectionApprovedTransfer) error {
	msg := types.NewMsgUpdateCollectionApprovedTransfers(creator, collectionId, approvedTransfers)
	_, err := suite.msgServer.UpdateCollectionApprovedTransfers(ctx, msg)
	return err
}

func RequestUpdateManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, add bool) error {
	msg := types.NewMsgRequestUpdateManager(creator, collectionId, add)
	_, err := suite.msgServer.RequestUpdateManager(ctx, msg)
	return err
}

func UpdateManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, address string) error {
	msg := types.NewMsgUpdateManager(creator, collectionId, address)
	_, err := suite.msgServer.UpdateManager(ctx, msg)
	return err
}

func UpdateURIs(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, collectionMetadata *CollectionMetadata, badgeMetadata []*types.BadgeMetadata, offChainBalancesMetadata *OffChainBalancesMetadata) error {
	msg := types.NewMsgUpdateMetadata(creator, collectionId, collectionMetadata, badgeMetadata, offChainBalancesMetadata)
	_, err := suite.msgServer.UpdateMetadata(ctx, msg)
	return err
}

func UpdateCollectionPermissions(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, permissions sdk.Uint) error {
	msg := types.NewMsgUpdateCollectionPermissions(creator, collectionId, permissions)
	_, err := suite.msgServer.UpdateCollectionPermissions(ctx, msg)
	return err
}

func UpdateCustomData(suite *TestSuite, ctx context.Context, creator string, collectionId sdk.Uint, bytes string) error {
	msg := types.NewMsgUpdateCustomData(creator, collectionId, bytes)
	_, err := suite.msgServer.UpdateCustomData(ctx, msg)
	return err
}

func ClaimBadge(suite *TestSuite, ctx context.Context, creator string, claimId sdk.Uint, collectionId sdk.Uint, solutions []*types.ChallengeSolution) error {
	msg := types.NewMsgClaimBadge(creator, claimId, collectionId, solutions)
	_, err := suite.msgServer.ClaimBadge(ctx, msg)
	return err
}
