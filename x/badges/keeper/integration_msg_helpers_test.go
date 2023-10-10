package keeper_test

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
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

// func GetNumUsedForMerkleChallenge(suite *TestSuite, ctx context.Context, challengeId string, level string, leafIndex sdkmath.Uint, collectionId sdkmath.Uint) (sdkmath.Uint, error) {
// 	res, err := suite.app.BadgesKeeper.GetNumUsedForMerkleChallenge(ctx, &types.QueryGetNumUsedForMerkleChallengeRequest{
// 		ChallengeId:  challengeId,
// 		Level:        level,
// 		LeafIndex:    leafIndex,
// 		CollectionId: collectionId,
// 	})
// 	if err != nil {
// 		return sdkmath.Uint{}, err
// 	}

// 	return res.NumUsed, nil
// }

// func GetApprovalsTracker(suite *TestSuite, ctx context.Context, collectionId sdkmath.Uint, address string, amountTrackerId string, level string, trackerType string) (*types.ApprovalsTracker, error) {
// 	res, err := suite.app.BadgesKeeper.GetApprovalsTracker(ctx, &types.QueryGetApprovalsTrackerRequest{
// 		CollectionId: sdkmath.Uint(collectionId),
// 		Address:      address,
// 		AmountTrackerId:    amountTrackerId,
// 		Level:        level,
// 		Depth:        depth,
// 	})
// 	if err != nil {
// 		return &types.ApprovalsTracker{}, err
// 	}

// 	return res.Tracker, nil
// }

// /* Msg helpers */

func UpdateCollection(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UpdateCollection(ctx, msg)
	return err
}

func UpdateCollectionWithRes(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollection) (*types.MsgUpdateCollectionResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	res, err := suite.msgServer.UpdateCollection(ctx, msg)
	return res, err
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

func UpdateUserApprovals(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateUserApprovals) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UpdateUserApprovals(ctx, msg)
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
		} else {
			return sdkerrors.Wrapf(types.ErrInvalidCollectionID, "Balances type %s not supported", collectionToCreate.BalancesType)
		}

		collectionRes, err := UpdateCollectionWithRes(suite, ctx, &types.MsgUpdateCollection{
			CollectionId: sdkmath.NewUint(0),
			Creator:      bob,
			BalancesType: balancesType,
			CollectionPermissions: collectionToCreate.Permissions,
			CollectionApprovals: collectionToCreate.CollectionApprovals,
			DefaultOutgoingApprovals: collectionToCreate.DefaultOutgoingApprovals,
			DefaultIncomingApprovals: collectionToCreate.DefaultIncomingApprovals,
			// ManagerTimeline: []*types.ManagerTimeline{
			// 	{
			// 		Manager: collectionToCreate.Creator,
			// 		TimelineTimes: GetFullUintRanges(),
			// 	},
			// },
			CollectionMetadataTimeline: collectionToCreate.CollectionMetadataTimeline,
			BadgeMetadataTimeline: collectionToCreate.BadgeMetadataTimeline,
			OffChainBalancesMetadataTimeline: collectionToCreate.OffChainBalancesMetadataTimeline,
			// InheritedCollectionId: collectionToCreate.InheritedCollectionId,
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
			
			UpdateCustomDataTimeline: true,
			UpdateContractAddressTimeline: true,
			UpdateCollectionApprovals: true,
			UpdateStandardsTimeline: true,
			// UpdateIsArchivedTimeline: true,

			DefaultAutoApproveSelfInitiatedOutgoingTransfers: !collectionToCreate.DefaultDisapproveSelfInitiated,
			DefaultAutoApproveSelfInitiatedIncomingTransfers: !collectionToCreate.DefaultDisapproveSelfInitiated,
		})
		if err != nil {
			return err
		}

		if len(collectionToCreate.Transfers) > 0 {
			err = TransferBadges(suite, ctx, &types.MsgTransferBadges{
				Creator: bob,
				CollectionId: collectionRes.CollectionId,
				Transfers: collectionToCreate.Transfers,
			})
			if err != nil {
				return err
			}
		}

		err = CreateAddressMappings(suite, ctx, &types.MsgCreateAddressMappings{
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
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		UpdateCollectionMetadataTimeline: true,
		BadgeMetadataTimeline: msg.BadgeMetadataTimeline,
		UpdateBadgeMetadataTimeline: true,
		OffChainBalancesMetadataTimeline: msg.OffChainBalancesMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		CollectionApprovals: msg.CollectionApprovals,
		UpdateCollectionApprovals: true,
	})
	if err != nil {
		return err 
	}

	if len(msg.Transfers) > 0 {
		_, err = suite.msgServer.TransferBadges(ctx, &types.MsgTransferBadges{
			Creator: bob,
			CollectionId: msg.CollectionId,
			Transfers: msg.Transfers,
		})
	}
	return err
}

func UpdateCollectionApprovals(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateCollectionApprovals) error {
	_, err := suite.msgServer.UpdateCollection(ctx, &types.MsgUpdateCollection{
		Creator: bob,
		CollectionId: msg.CollectionId,
		CollectionApprovals: msg.CollectionApprovals,
		UpdateCollectionApprovals: true,
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
	_, err := suite.msgServer.UpdateUserApprovals(ctx, &types.MsgUpdateUserApprovals{
		Creator: bob,
		CollectionId: msg.CollectionId,
		UserPermissions: msg.Permissions,
		UpdateUserPermissions: true,
	})
	return err
}
