package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

/* Query helpers */

func GetCollection(suite *TestSuite, ctx context.Context, id uint64) (types.BadgeCollection, error) {
	res, err := suite.app.BadgesKeeper.GetCollection(ctx, &types.QueryGetCollectionRequest{Id: uint64(id)})
	if err != nil {
		return types.BadgeCollection{}, err
	}

	return *res.Collection, nil
}

func GetUserBalance(suite *TestSuite, ctx context.Context, collectionId uint64, address uint64) (types.UserBalance, error) {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		BadgeId: uint64(collectionId),
		Address: uint64(address),
	})
	if err != nil {
		return types.UserBalance{}, err
	}

	return *res.Balance, nil
}

/* Msg helpers */

type CollectionsToCreate struct {
	Collection types.MsgNewCollection
	Amount     uint64
	Creator    string
}

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []CollectionsToCreate) error {
	for _, collectionToCreate := range collectionsToCreate {
		for i := 0; i < int(collectionToCreate.Amount); i++ {
			msg := types.NewMsgNewCollection(collectionToCreate.Creator, collectionToCreate.Collection.Standard, collectionToCreate.Collection.BadgeSupplys, collectionToCreate.Collection.CollectionUri, collectionToCreate.Collection.BadgeUri, collectionToCreate.Collection.Permissions, collectionToCreate.Collection.DisallowedTransfers, collectionToCreate.Collection.ManagerApprovedTransfers, collectionToCreate.Collection.Bytes, collectionToCreate.Collection.Transfers, collectionToCreate.Collection.Claims)
			_, err := suite.msgServer.NewCollection(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateBadges(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, supplysAndAmounts []*types.BadgeSupplyAndAmount, transfers []*types.Transfers, claims []*types.Claim) error {
	msg := types.NewMsgMintBadge(creator, collectionId, supplysAndAmounts, transfers, claims)
	_, err := suite.msgServer.MintBadge(ctx, msg)
	return err
}

//Note: Only supports Bob and Alice and only supports supplysAndAmounts[0]
func CreateBadgesAndMintAllToCreator(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, supplysAndAmounts []*types.BadgeSupplyAndAmount) error {
	creatorAccountNum := uint64(0)
	if creator == bob {
		creatorAccountNum = bobAccountNum
	}
	if creator == alice {
		creatorAccountNum = aliceAccountNum
	}

	collection, err := GetCollection(suite, ctx, collectionId)
	if err != nil {
		return err
	}

	transfers := []*types.Transfers{
		{
			ToAddresses: []uint64{creatorAccountNum},
			Balances: []*types.Balance{
				{
					Balance: supplysAndAmounts[0].Supply,
					BadgeIds: []*types.IdRange{
						{
							Start: collection.NextBadgeId,
							End:   collection.NextBadgeId + supplysAndAmounts[0].Amount - 1,
						},
					},
				},
			},
		},
	}

	msg := types.NewMsgMintBadge(creator, collectionId, supplysAndAmounts, transfers, []*types.Claim{})
	_, err = suite.msgServer.MintBadge(ctx, msg)
	return err
}

func TransferBadge(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, from uint64, transfers []*types.Transfers) error {
	msg := types.NewMsgTransferBadge(creator, collectionId, from, transfers)
	_, err := suite.msgServer.TransferBadge(ctx, msg)
	return err
}

func SetApproval(suite *TestSuite, ctx context.Context, creator string, address uint64, collectionId uint64, balances []*types.Balance) error {
	msg := types.NewMsgSetApproval(creator, collectionId, address, balances)
	_, err := suite.msgServer.SetApproval(ctx, msg)
	return err
}

func UpdateDisallowedTransfers(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, disallowedTransfers []*types.TransferMapping) error {
	msg := types.NewMsgUpdateDisallowedTransfers(creator, collectionId, disallowedTransfers)
	_, err := suite.msgServer.UpdateDisallowedTransfers(ctx, msg)
	return err
}

func RequestTransferManager(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, add bool) error {
	msg := types.NewMsgRequestTransferManager(creator, collectionId, add)
	_, err := suite.msgServer.RequestTransferManager(ctx, msg)
	return err
}

func TransferManager(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, address uint64) error {
	msg := types.NewMsgTransferManager(creator, collectionId, address)
	_, err := suite.msgServer.TransferManager(ctx, msg)
	return err
}

func UpdateURIs(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, collectionUri string, badgeUri string) error {
	msg := types.NewMsgUpdateUris(creator, collectionId, collectionUri, badgeUri)
	_, err := suite.msgServer.UpdateUris(ctx, msg)
	return err
}

func UpdatePermissions(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, permissions uint64) error {
	msg := types.NewMsgUpdatePermissions(creator, collectionId, permissions)
	_, err := suite.msgServer.UpdatePermissions(ctx, msg)
	return err
}

func UpdateBytes(suite *TestSuite, ctx context.Context, creator string, collectionId uint64, bytes string) error {
	msg := types.NewMsgUpdateBytes(creator, collectionId, bytes)
	_, err := suite.msgServer.UpdateBytes(ctx, msg)
	return err
}

func RegisterAddresses(suite *TestSuite, ctx context.Context, creator string, addresses []string) error {
	msg := types.NewMsgRegisterAddresses(creator, addresses)
	_, err := suite.msgServer.RegisterAddresses(ctx, msg)
	return err
}

func ClaimBadge(suite *TestSuite, ctx context.Context, creator string, claimId uint64, collectionId uint64, leaf string, proof *types.Proof, uri string, timeRange *types.IdRange) error {
	msg := types.NewMsgClaimBadge(creator, claimId, collectionId, leaf, proof, uri, timeRange)
	_, err := suite.msgServer.ClaimBadge(ctx, msg)
	return err
}