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

func GetAddressMapping(suite *TestSuite, ctx context.Context, mappingId string) (*types.AddressMapping, error) {
	res, err := suite.app.BadgesKeeper.GetAddressMapping(ctx, &types.QueryGetAddressMappingRequest{
		MappingId: mappingId,
	})
	if err != nil {
		return &types.AddressMapping{}, err
	}

	return res.Mapping, nil
}

func GetNumUsedForChallenge(suite *TestSuite, ctx context.Context, challengeId string, level string, leafIndex sdkmath.Uint, collectionId sdkmath.Uint) (sdkmath.Uint, error) {
	res, err := suite.app.BadgesKeeper.GetNumUsedForChallenge(ctx, &types.QueryGetNumUsedForChallengeRequest{
		ChallengeId:  challengeId,
		Level:        level,
		LeafIndex:    leafIndex,
		CollectionId: collectionId,
	})
	if err != nil {
		return sdkmath.Uint{}, err
	}

	return res.NumUsed, nil
}

func GetApprovalsTracker(suite *TestSuite, ctx context.Context, collectionId sdkmath.Uint, address string, trackerId string, level string, depth string) (*types.ApprovalsTracker, error) {
	res, err := suite.app.BadgesKeeper.GetApprovalsTracker(ctx, &types.QueryGetApprovalsTrackerRequest{
		CollectionId: sdkmath.Uint(collectionId),
		Address:      address,
		TrackerId:    trackerId,
		Level:        level,
		Depth:        depth,
	})
	if err != nil {
		return &types.ApprovalsTracker{}, err
	}

	return res.Tracker, nil
}

// /* Msg helpers */

func UpdateCollection(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UpdateCollection(ctx, msg)
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

func TransferBadges(suite *TestSuite, ctx context.Context, msg *types.MsgTransferBadges) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.TransferBadges(ctx, msg)
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

func CreateAddressMappings(suite *TestSuite, ctx context.Context, msg *types.MsgCreateAddressMappings) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.CreateAddressMappings(ctx, msg)
	return err
}


/** Legacy casts for test compatibility */

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []*types.MsgNewCollection) error {
	for _, collectionToCreate := range collectionsToCreate {
		balancesType := ""
		if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(1)) {
			balancesType = "Standard"
		} else if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(2)) {
			balancesType = "Off-Chain"
		} else if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(3)) {
			balancesType = "Inherited"
		}

		collectionRes, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
			CollectionId: sdkmath.NewUint(0),
			Creator:      bob,
			BalancesType: balancesType,
			CollectionPermissions: collectionToCreate.Permissions,
			CollectionApprovedTransfersTimeline: collectionToCreate.CollectionApprovedTransfersTimeline,
			DefaultApprovedOutgoingTransfersTimeline: collectionToCreate.DefaultApprovedOutgoingTransfersTimeline,
			DefaultApprovedIncomingTransfersTimeline: collectionToCreate.DefaultApprovedIncomingTransfersTimeline,
			// ManagerTimeline: []*types.ManagerTimeline{
			// 	{
			// 		Manager: collectionToCreate.Creator,
			// 		TimelineTimes: GetFullUintRanges(),
			// 	},
			// },
			CollectionMetadataTimeline: collectionToCreate.CollectionMetadataTimeline,
			BadgeMetadataTimeline: collectionToCreate.BadgeMetadataTimeline,
			OffChainBalancesMetadataTimeline: collectionToCreate.OffChainBalancesMetadataTimeline,
			InheritedBalancesTimeline: collectionToCreate.InheritedBalancesTimeline,
			CustomDataTimeline: collectionToCreate.CustomDataTimeline,
			ContractAddressTimeline: collectionToCreate.ContractAddressTimeline,
			StandardsTimeline: collectionToCreate.StandardsTimeline,
			BadgesToCreate: collectionToCreate.BadgesToCreate,
			// IsArchivedTimeline: collectionToCreate.IsArchivedTimeline,


			// ManagerTimeline: collectionToCreate.ManagerTimeline,
			UpdateCollectionPermissions: true,
			// UpdateManagerTimeline: true,
			UpdateCollectionMetadataTimeline: true,
			UpdateBadgeMetadataTimeline: true,
			UpdateOffChainBalancesMetadataTimeline: true,
			UpdateInheritedBalancesTimeline: true,
			UpdateCustomDataTimeline: true,
			UpdateContractAddressTimeline: true,
			UpdateCollectionApprovedTransfersTimeline: true,
			UpdateStandardsTimeline: true,
			// UpdateIsArchivedTimeline: true,
		})
		if err != nil {
			return err
		}

		_, err = suite.msgServer.TransferBadges(ctx, &types.MsgTransferBadges{
			Creator: bob,
			CollectionId: collectionRes.CollectionId,
			Transfers: collectionToCreate.Transfers,
		})
		if err != nil {
			return err
		}

		_, err = suite.msgServer.CreateAddressMappings(ctx, &types.MsgCreateAddressMappings{
			Creator: bob,
			AddressMappings: collectionToCreate.AddressMappings,
		})
		if err != nil {
			return err
		}
	
	}
	return nil
}

func MintAndDistributeBadges(suite *TestSuite, ctx context.Context, msg *types.MsgMintAndDistributeBadges) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		BadgesToCreate: msg.BadgesToCreate,
		InheritedBalancesTimeline: msg.InheritedBalancesTimeline,
		UpdateInheritedBalancesTimeline: true,
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		UpdateCollectionMetadataTimeline: true,
		BadgeMetadataTimeline: msg.BadgeMetadataTimeline,
		UpdateBadgeMetadataTimeline: true,
		OffChainBalancesMetadataTimeline: msg.OffChainBalancesMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		CollectionApprovedTransfersTimeline: msg.CollectionApprovedTransfersTimeline,
		UpdateCollectionApprovedTransfersTimeline: true,
	})
	if err != nil {
		return err 
	}

	_, err = suite.msgServer.TransferBadges(ctx, &types.MsgTransferBadges{
		Creator: bob,
		CollectionId: msg.CollectionId,
		Transfers: msg.Transfers,
	})
	return err
}

func UpdateCollectionApprovedTransfers(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollectionApprovedTransfers) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		CollectionApprovedTransfersTimeline: msg.CollectionApprovedTransfersTimeline,
		UpdateCollectionApprovedTransfersTimeline: true,
	})
	return err
}

func ArchiveCollection(suite *TestSuite, ctx context.Context, msg *types.MsgArchiveCollection) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		IsArchivedTimeline: msg.IsArchivedTimeline,
		UpdateIsArchivedTimeline: true,
	})
	return err
}

func UpdateManager(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateManager) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		ManagerTimeline: msg.ManagerTimeline,
		UpdateManagerTimeline: true,
	})
	return err
}

func UpdateMetadata(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateMetadata) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		UpdateCollectionMetadataTimeline: true,
		BadgeMetadataTimeline: msg.BadgeMetadataTimeline,
		UpdateBadgeMetadataTimeline: true,
		OffChainBalancesMetadataTimeline: msg.OffChainBalancesMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		StandardsTimeline: msg.StandardsTimeline,
		UpdateStandardsTimeline: true,
		ContractAddressTimeline: msg.ContractAddressTimeline,
		UpdateContractAddressTimeline: true,
		CustomDataTimeline: msg.CustomDataTimeline,
		UpdateCustomDataTimeline: true,
	})
	if err != nil {
		return err 
	}

	return nil
}

func UpdateCollectionPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollectionPermissions) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		CollectionPermissions: msg.Permissions,
		UpdateCollectionPermissions: true,
	})
	return err
}


func UpdateUserPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateUserPermissions) error {
	_, err := suite.msgServer.UpdateUserApprovedTransfers(ctx, &types.MsgUpdateUserApprovedTransfers{
		Creator: bob,
		CollectionId: msg.CollectionId,
		Permissions: msg.Permissions,
		UpdateApprovedTransfersUserPermissions: true,
	})
	return err
}
