package keeper_test

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

/* Query helpers */

func GetCollection(suite *TestSuite, ctx context.Context, id sdkmath.Uint) (*types.BadgeCollection, error) {
	res, err := suite.app.BadgesKeeper.GetCollection(ctx, &types.QueryGetCollectionRequest{CollectionId: sdkmath.Uint(id).String()})
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	return res.Collection, nil
}

func GetUserBalance(suite *TestSuite, ctx context.Context, collectionId sdkmath.Uint, address string) (*types.UserBalanceStore, error) {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		CollectionId: collectionId.String(),
		Address:      address,
	})
	if err != nil {
		return &types.UserBalanceStore{}, err
	}

	return res.Balance, nil
}

func GetAddressList(suite *TestSuite, ctx context.Context, listId string) (*types.AddressList, error) {
	res, err := suite.app.BadgesKeeper.GetAddressList(ctx, &types.QueryGetAddressListRequest{
		ListId: listId,
	})
	if err != nil {
		return &types.AddressList{}, err
	}

	return res.List, nil
}

// func GetChallengeTracker(suite *TestSuite, ctx context.Context, challengeId string, level string, leafIndex sdkmath.Uint, collectionId sdkmath.Uint) (sdkmath.Uint, error) {
// 	res, err := suite.app.BadgesKeeper.GetChallengeTracker(ctx, &types.QueryGetChallengeTrackerRequest{
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

// func GetApprovalTracker(suite *TestSuite, ctx context.Context, collectionId sdkmath.Uint, address string, amountTrackerId string, level string, trackerType string) (*types.ApprovalTracker, error) {
// 	res, err := suite.app.BadgesKeeper.GetApprovalTracker(ctx, &types.QueryGetApprovalTrackerRequest{
// 		CollectionId: sdkmath.Uint(collectionId),
// 		Address:      address,
// 		AmountTrackerId:    amountTrackerId,
// 		Level:        level,
// 		Depth:        depth,
// 	})
// 	if err != nil {
// 		return &types.ApprovalTracker{}, err
// 	}

// 	return res.Tracker, nil
// }

// /* Msg helpers */

func UpdateCollection(suite *TestSuite, ctx context.Context, msg *types.MsgUniversalUpdateCollection) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UniversalUpdateCollection(ctx, msg)
	return err
}

func UpdateCollectionWithRes(suite *TestSuite, ctx context.Context, msg *types.MsgUniversalUpdateCollection) (*types.MsgUniversalUpdateCollectionResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	res, err := suite.msgServer.UniversalUpdateCollection(ctx, msg)
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

func CreateAddressLists(suite *TestSuite, ctx context.Context, msg *types.MsgCreateAddressLists) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.CreateAddressLists(ctx, msg)
	return err
}

/** Legacy casts for test compatibility */

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []*types.MsgNewCollection) error {
	for _, collectionToCreate := range collectionsToCreate {
		balancesType := ""
		if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(1)) {
			balancesType = "Standard"
		} else if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(2)) {
			balancesType = "Off-Chain - Indexed"
		} else if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(3)) {
			balancesType = "Inherited"
		} else if collectionToCreate.BalancesType.Equal(sdkmath.NewUint(4)) {
			balancesType = "Off-Chain - Non-Indexed"
		} else {
			return sdkerrors.Wrapf(types.ErrInvalidCollectionID, "Balances type %s not supported", collectionToCreate.BalancesType)
		}

		allBadgeIds := []*types.UintRange{}
		for _, badge := range collectionToCreate.BadgesToCreate {
			allBadgeIds = append(allBadgeIds, badge.BadgeIds...)
		}
		allBadgeIds = types.SortUintRangesAndMergeAdjacentAndIntersecting(allBadgeIds)

		//For legacy purposes, we will use badgesToCreate which mints them to the mint address
		collectionRes, err := UpdateCollectionWithRes(suite, ctx, &types.MsgUniversalUpdateCollection{
			CollectionId:          sdkmath.NewUint(0),
			Creator:               bob,
			BalancesType:          balancesType,
			CollectionPermissions: collectionToCreate.Permissions,
			CollectionApprovals:   collectionToCreate.CollectionApprovals,

			DefaultBalances: &types.UserBalanceStore{
				Balances:          collectionToCreate.DefaultBalances,
				OutgoingApprovals: collectionToCreate.DefaultOutgoingApprovals,
				IncomingApprovals: collectionToCreate.DefaultIncomingApprovals,
				AutoApproveSelfInitiatedOutgoingTransfers: !collectionToCreate.DefaultDisapproveSelfInitiated,
				AutoApproveSelfInitiatedIncomingTransfers: !collectionToCreate.DefaultDisapproveSelfInitiated,
				UserPermissions: nil,
			},

			CollectionMetadataTimeline:       collectionToCreate.CollectionMetadataTimeline,
			BadgeMetadataTimeline:            collectionToCreate.BadgeMetadataTimeline,
			OffChainBalancesMetadataTimeline: collectionToCreate.OffChainBalancesMetadataTimeline,
			// InheritedCollectionId: collectionToCreate.InheritedCollectionId,
			CustomDataTimeline: collectionToCreate.CustomDataTimeline,
			StandardsTimeline:  collectionToCreate.StandardsTimeline,
			BadgeIdsToAdd:      allBadgeIds,
			// IsArchivedTimeline: collectionToCreate.IsArchivedTimeline,

			// ManagerTimeline: collectionToCreate.ManagerTimeline,
			UpdateCollectionPermissions: true,
			// UpdateManagerTimeline: true,
			UpdateCollectionMetadataTimeline:       true,
			UpdateBadgeMetadataTimeline:            true,
			UpdateOffChainBalancesMetadataTimeline: true,

			UpdateCustomDataTimeline:  true,
			UpdateCollectionApprovals: true,
			UpdateStandardsTimeline:   true,
			// UpdateIsArchivedTimeline: true,
		})
		if err != nil {
			return err
		}

		if len(collectionToCreate.Transfers) > 0 {
			err = TransferBadges(suite, ctx, &types.MsgTransferBadges{
				Creator:      bob,
				CollectionId: collectionRes.CollectionId,
				Transfers:    collectionToCreate.Transfers,
			})
			if err != nil {
				return err
			}
		}

		err = CreateAddressLists(suite, ctx, &types.MsgCreateAddressLists{
			Creator:      bob,
			AddressLists: collectionToCreate.AddressLists,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func MintAndDistributeBadges(suite *TestSuite, ctx context.Context, msg *types.MsgMintAndDistributeBadges) error {
	allBadgeIds := []*types.UintRange{}
	for _, badge := range msg.BadgesToCreate {
		allBadgeIds = append(allBadgeIds, badge.BadgeIds...)
	}
	allBadgeIds = types.SortUintRangesAndMergeAdjacentAndIntersecting(allBadgeIds)

	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                                bob,
		CollectionId:                           msg.CollectionId,
		BadgeIdsToAdd:                          allBadgeIds,
		CollectionMetadataTimeline:             msg.CollectionMetadataTimeline,
		UpdateCollectionMetadataTimeline:       true,
		BadgeMetadataTimeline:                  msg.BadgeMetadataTimeline,
		UpdateBadgeMetadataTimeline:            true,
		OffChainBalancesMetadataTimeline:       msg.OffChainBalancesMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		CollectionApprovals:                    msg.CollectionApprovals,
		UpdateCollectionApprovals:              true,
		DefaultBalances:                        &types.UserBalanceStore{},
	})
	if err != nil {
		return err
	}

	if len(msg.Transfers) > 0 {
		_, err = suite.msgServer.TransferBadges(ctx, &types.MsgTransferBadges{
			Creator:      bob,
			CollectionId: msg.CollectionId,
			Transfers:    msg.Transfers,
		})
	}
	return err
}

func UpdateCollectionApprovals(suite *TestSuite, ctx context.Context, msg *types.MsgUniversalUpdateCollectionApprovals) error {
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                   bob,
		CollectionId:              msg.CollectionId,
		CollectionApprovals:       msg.CollectionApprovals,
		UpdateCollectionApprovals: true,
	})
	return err
}

func ArchiveCollection(suite *TestSuite, ctx context.Context, msg *types.MsgArchiveCollection) error {
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                  bob,
		CollectionId:             msg.CollectionId,
		IsArchivedTimeline:       msg.IsArchivedTimeline,
		UpdateIsArchivedTimeline: true,
	})
	return err
}

func UpdateManager(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateManager) error {
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:               bob,
		CollectionId:          msg.CollectionId,
		ManagerTimeline:       msg.ManagerTimeline,
		UpdateManagerTimeline: true,
	})
	return err
}

func UpdateMetadata(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateMetadata) error {
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                                bob,
		CollectionId:                           msg.CollectionId,
		CollectionMetadataTimeline:             msg.CollectionMetadataTimeline,
		UpdateCollectionMetadataTimeline:       true,
		BadgeMetadataTimeline:                  msg.BadgeMetadataTimeline,
		UpdateBadgeMetadataTimeline:            true,
		OffChainBalancesMetadataTimeline:       msg.OffChainBalancesMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		StandardsTimeline:                      msg.StandardsTimeline,
		UpdateStandardsTimeline:                true,
		CustomDataTimeline:                     msg.CustomDataTimeline,
		UpdateCustomDataTimeline:               true,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateCollectionPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUniversalUpdateCollectionPermissions) error {
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                     bob,
		CollectionId:                msg.CollectionId,
		CollectionPermissions:       msg.Permissions,
		UpdateCollectionPermissions: true,
	})
	return err
}

func UpdateUserPermissions(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateUserPermissions) error {
	_, err := suite.msgServer.UpdateUserApprovals(ctx, &types.MsgUpdateUserApprovals{
		Creator:               bob,
		CollectionId:          msg.CollectionId,
		UserPermissions:       msg.Permissions,
		UpdateUserPermissions: true,
	})
	return err
}
