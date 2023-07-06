package keeper_test

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

/* Query helpers */

func GetCollection(suite *TestSuite, ctx context.Context, id sdkmath.Uint) (*types.BadgeCollection, error) {
	res, err := suite.app.BadgesKeeper.GetCollection(ctx, &types.QueryGetCollectionRequest{CollectionId: sdkmath.Uint(id)})
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	return res.Collection, nil
}

func GetUserBalance(suite *TestSuite, ctx context.Context, collectionId sdkmath.Uint, address string) (*types.UserBalanceStore, error) {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		CollectionId: sdkmath.Uint(collectionId),
		Address:      address,
	})
	if err != nil {
		return &types.UserBalanceStore{}, err
	}

	return res.Balance, nil
}

//TODO: Queries

// /* Msg helpers */

type CollectionsToCreate struct {
	Collection types.MsgNewCollection
	Amount     sdkmath.Uint
	Creator    string
}

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []CollectionsToCreate) error {
	for _, collectionToCreate := range collectionsToCreate {
		for i := 0; i < int(collectionToCreate.Amount.BigInt().Int64()); i++ {
			err := collectionToCreate.Collection.ValidateBasic()
			if err != nil {
				return err
			}

			_, err = suite.msgServer.NewCollection(ctx, &collectionToCreate.Collection)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewCollection(suite *TestSuite, ctx context.Context, msg *types.MsgNewCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.NewCollection(ctx, msg)
	return err
}

func MintAndDistributeBadges(suite *TestSuite, ctx context.Context, msg *types.MsgMintAndDistributeBadges) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.MintAndDistributeBadges(ctx, msg)
	return err
}

func ArchiveCollection(suite *TestSuite, ctx context.Context, msg *types.MsgArchiveCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.ArchiveCollection(ctx, msg)
	return err
}

func DeleteCollection(suite *TestSuite, ctx context.Context, msg *types.MsgDeleteCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.DeleteCollection(ctx, msg)
	return err
}


func TransferBadge(suite *TestSuite, ctx context.Context, msg *types.MsgTransferBadge) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.TransferBadge(ctx, msg)
	return err
}

func UpdateManager(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateManager) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateManager(ctx, msg)
	return err
}

func UpdateMetadata(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateMetadata) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateMetadata(ctx, msg)
	return err
}

func UpdateCollectionApprovedTransfers(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollectionApprovedTransfers) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateCollectionApprovedTransfers(ctx, msg)
	return err
}

func UpdateCollectionPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollectionPermissions) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateCollectionPermissions(ctx, msg)
	return err
}

func UpdateUserApprovedTransfers(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateUserApprovedTransfers) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateUserApprovedTransfers(ctx, msg)
	return err
}

func UpdateUserPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateUserPermissions) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	_, err = suite.msgServer.UpdateUserPermissions(ctx, msg)
	return err
}