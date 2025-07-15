package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

func CreateDynamicStore(suite *TestSuite, ctx context.Context, msg *types.MsgCreateDynamicStore) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.CreateDynamicStore(ctx, msg)
	return err
}

func UpdateDynamicStore(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateDynamicStore) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UpdateDynamicStore(ctx, msg)
	return err
}

func DeleteDynamicStore(suite *TestSuite, ctx context.Context, msg *types.MsgDeleteDynamicStore) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.DeleteDynamicStore(ctx, msg)
	return err
}

func SetDynamicStoreValue(suite *TestSuite, ctx context.Context, msg *types.MsgSetDynamicStoreValue) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.SetDynamicStoreValue(ctx, msg)
	return err
}

/** Legacy casts for test compatibility */

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []*types.MsgNewCollection) error {
	for _, collectionToCreate := range collectionsToCreate {
		balancesType := ""
		if sdkmath.NewUint(1).Equal(sdkmath.NewUint(collectionToCreate.BalancesType.Uint64())) {
			balancesType = "Standard"
		} else if sdkmath.NewUint(2).Equal(sdkmath.NewUint(collectionToCreate.BalancesType.Uint64())) {
			balancesType = "Off-Chain - Indexed"
		} else if sdkmath.NewUint(3).Equal(sdkmath.NewUint(collectionToCreate.BalancesType.Uint64())) {
			balancesType = "Inherited"
		} else if sdkmath.NewUint(4).Equal(sdkmath.NewUint(collectionToCreate.BalancesType.Uint64())) {
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
				OutgoingApprovals: []*types.UserOutgoingApproval{},
				IncomingApprovals: []*types.UserIncomingApproval{},
				AutoApproveSelfInitiatedOutgoingTransfers: true,
				AutoApproveSelfInitiatedIncomingTransfers: true,
				UserPermissions: nil,
			},
			CollectionMetadataTimeline:       collectionToCreate.CollectionMetadataTimeline,
			BadgeMetadataTimeline:            collectionToCreate.BadgeMetadataTimeline,
			OffChainBalancesMetadataTimeline: collectionToCreate.OffChainBalancesMetadataTimeline,
			// InheritedCollectionId: collectionToCreate.InheritedCollectionId,
			CustomDataTimeline:          collectionToCreate.CustomDataTimeline,
			StandardsTimeline:           collectionToCreate.StandardsTimeline,
			CosmosCoinWrapperPathsToAdd: collectionToCreate.CosmosCoinWrapperPathsToAdd,
			ValidBadgeIds:               allBadgeIds,
			// IsArchivedTimeline: collectionToCreate.IsArchivedTimeline,

			// ManagerTimeline: collectionToCreate.ManagerTimeline,
			UpdateValidBadgeIds:         true,
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

		//Update for bob
		err = UpdateUserApprovals(suite, ctx, &types.MsgUpdateUserApprovals{
			Creator:                 bob,
			CollectionId:            collectionRes.CollectionId,
			OutgoingApprovals:       collectionToCreate.DefaultOutgoingApprovals,
			IncomingApprovals:       collectionToCreate.DefaultIncomingApprovals,
			UpdateOutgoingApprovals: true,
			UpdateIncomingApprovals: true,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
			AutoApproveSelfInitiatedIncomingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
		})
		if err != nil {
			return err
		}

		//Update for alice
		err = UpdateUserApprovals(suite, ctx, &types.MsgUpdateUserApprovals{
			Creator:                 alice,
			CollectionId:            collectionRes.CollectionId,
			OutgoingApprovals:       collectionToCreate.DefaultOutgoingApprovals,
			IncomingApprovals:       collectionToCreate.DefaultIncomingApprovals,
			UpdateOutgoingApprovals: true,
			UpdateIncomingApprovals: true,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
			AutoApproveSelfInitiatedIncomingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
		})
		if err != nil {
			return err
		}

		//Update for charlie
		err = UpdateUserApprovals(suite, ctx, &types.MsgUpdateUserApprovals{
			Creator:                 charlie,
			CollectionId:            collectionRes.CollectionId,
			OutgoingApprovals:       collectionToCreate.DefaultOutgoingApprovals,
			IncomingApprovals:       collectionToCreate.DefaultIncomingApprovals,
			UpdateOutgoingApprovals: true,
			UpdateIncomingApprovals: true,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
			AutoApproveSelfInitiatedIncomingTransfers:       !collectionToCreate.DefaultDisapproveSelfInitiated,
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
		UpdateValidBadgeIds:                    true,
		ValidBadgeIds:                          allBadgeIds,
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

	newTransfers := []*types.Transfer{}
	for _, transfer := range msg.Transfers {
		newTransfer := transfer
		newTransfer.PrioritizedApprovals = GetDefaultPrioritizedApprovals(sdk.UnwrapSDKContext(ctx), suite.app.BadgesKeeper, msg.CollectionId)
		newTransfers = append(newTransfers, newTransfer)
	}

	if len(msg.Transfers) > 0 {
		_, err = suite.msgServer.TransferBadges(ctx, &types.MsgTransferBadges{
			Creator:      bob,
			CollectionId: msg.CollectionId,
			Transfers:    newTransfers,
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

// Helper functions for UniversalUpdateCollection subsets using native messages

func SetValidBadgeIds(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, validBadgeIds []*types.UintRange) error {
	msg := &types.MsgSetValidBadgeIds{
		Creator:       creator,
		CollectionId:  collectionId,
		ValidBadgeIds: validBadgeIds,
	}
	_, err := suite.msgServer.SetValidBadgeIds(ctx, msg)
	return err
}

func SetManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, managerTimeline []*types.ManagerTimeline) error {
	msg := &types.MsgSetManager{
		Creator:         creator,
		CollectionId:    collectionId,
		ManagerTimeline: managerTimeline,
	}
	_, err := suite.msgServer.SetManager(ctx, msg)
	return err
}

func SetCollectionMetadata(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, collectionMetadataTimeline []*types.CollectionMetadataTimeline) error {
	msg := &types.MsgSetCollectionMetadata{
		Creator:                    creator,
		CollectionId:               collectionId,
		CollectionMetadataTimeline: collectionMetadataTimeline,
	}
	_, err := suite.msgServer.SetCollectionMetadata(ctx, msg)
	return err
}

func SetBadgeMetadata(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, badgeMetadataTimeline []*types.BadgeMetadataTimeline) error {
	msg := &types.MsgSetBadgeMetadata{
		Creator:               creator,
		CollectionId:          collectionId,
		BadgeMetadataTimeline: badgeMetadataTimeline,
	}
	_, err := suite.msgServer.SetBadgeMetadata(ctx, msg)
	return err
}

func SetCustomData(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, customDataTimeline []*types.CustomDataTimeline) error {
	msg := &types.MsgSetCustomData{
		Creator:            creator,
		CollectionId:       collectionId,
		CustomDataTimeline: customDataTimeline,
	}
	_, err := suite.msgServer.SetCustomData(ctx, msg)
	return err
}

func SetStandards(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, standardsTimeline []*types.StandardsTimeline) error {
	msg := &types.MsgSetStandards{
		Creator:           creator,
		CollectionId:      collectionId,
		StandardsTimeline: standardsTimeline,
	}
	_, err := suite.msgServer.SetStandards(ctx, msg)
	return err
}

func SetCollectionApprovals(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, collectionApprovals []*types.CollectionApproval) error {
	msg := &types.MsgSetCollectionApprovals{
		Creator:             creator,
		CollectionId:        collectionId,
		CollectionApprovals: collectionApprovals,
	}
	_, err := suite.msgServer.SetCollectionApprovals(ctx, msg)
	return err
}

func SetIsArchived(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, isArchivedTimeline []*types.IsArchivedTimeline) error {
	msg := &types.MsgSetIsArchived{
		Creator:            creator,
		CollectionId:       collectionId,
		IsArchivedTimeline: isArchivedTimeline,
	}
	_, err := suite.msgServer.SetIsArchived(ctx, msg)
	return err
}
