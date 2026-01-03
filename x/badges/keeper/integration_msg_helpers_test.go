package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

/* Query helpers */

func GetCollection(suite *TestSuite, ctx context.Context, id sdkmath.Uint) (*types.TokenCollection, error) {
	res, err := suite.app.BadgesKeeper.GetCollection(ctx, &types.QueryGetCollectionRequest{CollectionId: sdkmath.Uint(id).String()})
	if err != nil {
		return &types.TokenCollection{}, err
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

func TransferTokens(suite *TestSuite, ctx context.Context, msg *types.MsgTransferTokens) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.TransferTokens(ctx, msg)
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

// Helper functions to extract current values from timeline arrays (for legacy MsgNewCollection compatibility)
func getCurrentCollectionMetadataFromTimeline(timeline []*types.CollectionMetadataTimeline) *types.CollectionMetadata {
	if len(timeline) == 0 {
		return nil
	}
	// Return the first timeline entry's metadata
	return timeline[0].CollectionMetadata
}

func getCurrentTokenMetadataFromTimeline(timeline []*types.TokenMetadataTimeline) []*types.TokenMetadata {
	if len(timeline) == 0 {
		return nil
	}
	// Return the first timeline entry's token metadata
	return timeline[0].TokenMetadata
}

func getCurrentCustomDataFromTimeline(timeline []*types.CustomDataTimeline) string {
	if len(timeline) == 0 {
		return ""
	}
	// Return the first timeline entry's custom data
	return timeline[0].CustomData
}

func getCurrentStandardsFromTimeline(timeline []*types.StandardsTimeline) []string {
	if len(timeline) == 0 {
		return nil
	}
	// Return the first timeline entry's standards
	return timeline[0].Standards
}

func CreateCollections(suite *TestSuite, ctx context.Context, collectionsToCreate []*types.MsgNewCollection) error {
	for _, collectionToCreate := range collectionsToCreate {
		// All collections now use Standard balances

		allTokenIds := []*types.UintRange{}
		for _, badge := range collectionToCreate.TokensToCreate {
			allTokenIds = append(allTokenIds, badge.TokenIds...)
		}
		allTokenIds = types.SortUintRangesAndMergeAdjacentAndIntersecting(allTokenIds)

		normalizedWrapperPaths := []*types.CosmosCoinWrapperPathAddObject{}
		for _, wrapperPath := range collectionToCreate.CosmosCoinWrapperPathsToAdd {
			if wrapperPath.Conversion == nil || wrapperPath.Conversion.SideA == nil {
				wrapperPath.Conversion = &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{
						Amount: sdkmath.NewUint(1),
					},
					SideB: []*types.Balance{},
				}
			} else if wrapperPath.Conversion.SideA.Amount.IsNil() || wrapperPath.Conversion.SideA.Amount.IsZero() {
				wrapperPath.Conversion.SideA.Amount = sdkmath.NewUint(1)
			}
			normalizedWrapperPaths = append(normalizedWrapperPaths, wrapperPath)
		}

		normalizedAliasPaths := []*types.AliasPathAddObject{}
		for _, aliasPath := range collectionToCreate.AliasPathsToAdd {
			if aliasPath.Conversion == nil || aliasPath.Conversion.SideA == nil {
				aliasPath.Conversion = &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{
						Amount: sdkmath.NewUint(1),
					},
					SideB: []*types.Balance{},
				}
			} else if aliasPath.Conversion.SideA.Amount.IsNil() || aliasPath.Conversion.SideA.Amount.IsZero() {
				aliasPath.Conversion.SideA.Amount = sdkmath.NewUint(1)
			}
			normalizedAliasPaths = append(normalizedAliasPaths, aliasPath)
		}

		//For legacy purposes, we will use tokensToCreate which mints them to the mint address
		collectionRes, err := UpdateCollectionWithRes(suite, ctx, &types.MsgUniversalUpdateCollection{
			CollectionId:          sdkmath.NewUint(0),
			Creator:               bob,
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
			// Convert timeline fields to simple values (use first timeline entry or empty)
			CollectionMetadata:          getCurrentCollectionMetadataFromTimeline(collectionToCreate.CollectionMetadataTimeline),
			TokenMetadata:               getCurrentTokenMetadataFromTimeline(collectionToCreate.TokenMetadataTimeline),
			CustomData:                  getCurrentCustomDataFromTimeline(collectionToCreate.CustomDataTimeline),
			Standards:                   getCurrentStandardsFromTimeline(collectionToCreate.StandardsTimeline),
			CosmosCoinWrapperPathsToAdd: normalizedWrapperPaths,
			AliasPathsToAdd:             normalizedAliasPaths,
			ValidTokenIds:               allTokenIds,
			Invariants:                  collectionToCreate.Invariants,
			IsArchived:                  false, // Default to not archived

			Manager:                     "", // Default manager will be set to creator
			UpdateValidTokenIds:         true,
			UpdateCollectionPermissions: true,
			UpdateManager:               false, // Don't update manager, use default
			UpdateCollectionMetadata:    len(collectionToCreate.CollectionMetadataTimeline) > 0,
			UpdateTokenMetadata:         len(collectionToCreate.TokenMetadataTimeline) > 0,
			UpdateCustomData:            len(collectionToCreate.CustomDataTimeline) > 0,
			UpdateCollectionApprovals:   true,
			UpdateStandards:             len(collectionToCreate.StandardsTimeline) > 0,
			UpdateIsArchived:            false,
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
			err = TransferTokens(suite, ctx, &types.MsgTransferTokens{
				Creator:      bob,
				CollectionId: collectionRes.CollectionId,
				Transfers:    collectionToCreate.Transfers,
			})
			if err != nil {
				return err
			}
		}

		// Convert AddressList to AddressListInput (remove createdBy field)
		addressListInputs := make([]*types.AddressListInput, len(collectionToCreate.AddressLists))
		for i, addrList := range collectionToCreate.AddressLists {
			addressListInputs[i] = &types.AddressListInput{
				ListId:     addrList.ListId,
				Addresses:  addrList.Addresses,
				Whitelist:  addrList.Whitelist,
				Uri:        addrList.Uri,
				CustomData: addrList.CustomData,
			}
		}
		err = CreateAddressLists(suite, ctx, &types.MsgCreateAddressLists{
			Creator:      bob,
			AddressLists: addressListInputs,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func MintAndDistributeTokens(suite *TestSuite, ctx context.Context, msg *types.MsgMintAndDistributeTokens) error {
	allTokenIds := []*types.UintRange{}
	for _, badge := range msg.TokensToCreate {
		allTokenIds = append(allTokenIds, badge.TokenIds...)
	}
	allTokenIds = types.SortUintRangesAndMergeAdjacentAndIntersecting(allTokenIds)

	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                   bob,
		CollectionId:              msg.CollectionId,
		UpdateValidTokenIds:       true,
		ValidTokenIds:             allTokenIds,
		CollectionMetadata:        getCurrentCollectionMetadataFromTimeline(msg.CollectionMetadataTimeline),
		UpdateCollectionMetadata:  len(msg.CollectionMetadataTimeline) > 0,
		TokenMetadata:             getCurrentTokenMetadataFromTimeline(msg.TokenMetadataTimeline),
		UpdateTokenMetadata:       len(msg.TokenMetadataTimeline) > 0,
		CollectionApprovals:       msg.CollectionApprovals,
		UpdateCollectionApprovals: true,
		DefaultBalances:           &types.UserBalanceStore{},
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
		_, err = suite.msgServer.TransferTokens(ctx, &types.MsgTransferTokens{
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
	// Extract isArchived value from timeline (use first entry or default to true for archive)
	isArchived := true
	if len(msg.IsArchivedTimeline) > 0 {
		isArchived = msg.IsArchivedTimeline[0].IsArchived
	}
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:          bob,
		CollectionId:     msg.CollectionId,
		IsArchived:       isArchived,
		UpdateIsArchived: true,
	})
	return err
}

func UpdateManager(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateManager) error {
	// Extract manager value from timeline (use first entry)
	manager := ""
	if len(msg.ManagerTimeline) > 0 {
		manager = msg.ManagerTimeline[0].Manager
	}
	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:       bob,
		CollectionId:  msg.CollectionId,
		Manager:       manager,
		UpdateManager: true,
	})
	return err
}

func UpdateMetadata(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateMetadata) error {
	// Extract values from timelines (use first entry for each)
	collectionMetadata := getCurrentCollectionMetadataFromTimeline(msg.CollectionMetadataTimeline)
	tokenMetadata := getCurrentTokenMetadataFromTimeline(msg.TokenMetadataTimeline)
	standards := getCurrentStandardsFromTimeline(msg.StandardsTimeline)
	customData := getCurrentCustomDataFromTimeline(msg.CustomDataTimeline)

	_, err := suite.msgServer.UniversalUpdateCollection(ctx, &types.MsgUniversalUpdateCollection{
		Creator:                  bob,
		CollectionId:             msg.CollectionId,
		CollectionMetadata:       collectionMetadata,
		UpdateCollectionMetadata: len(msg.CollectionMetadataTimeline) > 0,
		TokenMetadata:            tokenMetadata,
		UpdateTokenMetadata:      len(msg.TokenMetadataTimeline) > 0,
		Standards:                standards,
		UpdateStandards:          len(msg.StandardsTimeline) > 0,
		CustomData:               customData,
		UpdateCustomData:         len(msg.CustomDataTimeline) > 0,
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

func SetValidTokenIds(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, validTokenIds []*types.UintRange) error {
	msg := &types.MsgSetValidTokenIds{
		Creator:       creator,
		CollectionId:  collectionId,
		ValidTokenIds: validTokenIds,
	}
	_, err := suite.msgServer.SetValidTokenIds(ctx, msg)
	return err
}

func SetManager(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, managerTimeline []*types.ManagerTimeline) error {
	// Extract manager from timeline (use first entry)
	manager := ""
	if len(managerTimeline) > 0 {
		manager = managerTimeline[0].Manager
	}
	msg := &types.MsgSetManager{
		Creator:      creator,
		CollectionId: collectionId,
		Manager:      manager,
	}
	_, err := suite.msgServer.SetManager(ctx, msg)
	return err
}

func SetCollectionMetadata(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, collectionMetadataTimeline []*types.CollectionMetadataTimeline) error {
	msg := &types.MsgSetCollectionMetadata{
		Creator:            creator,
		CollectionId:       collectionId,
		CollectionMetadata: getCurrentCollectionMetadataFromTimeline(collectionMetadataTimeline),
	}
	_, err := suite.msgServer.SetCollectionMetadata(ctx, msg)
	return err
}

func SetTokenMetadata(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, tokenMetadataTimeline []*types.TokenMetadataTimeline) error {
	msg := &types.MsgSetTokenMetadata{
		Creator:       creator,
		CollectionId:  collectionId,
		TokenMetadata: getCurrentTokenMetadataFromTimeline(tokenMetadataTimeline),
	}
	_, err := suite.msgServer.SetTokenMetadata(ctx, msg)
	return err
}

func SetCustomData(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, customDataTimeline []*types.CustomDataTimeline) error {
	msg := &types.MsgSetCustomData{
		Creator:      creator,
		CollectionId: collectionId,
		CustomData:   getCurrentCustomDataFromTimeline(customDataTimeline),
	}
	_, err := suite.msgServer.SetCustomData(ctx, msg)
	return err
}

func SetStandards(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, standardsTimeline []*types.StandardsTimeline) error {
	msg := &types.MsgSetStandards{
		Creator:      creator,
		CollectionId: collectionId,
		Standards:    getCurrentStandardsFromTimeline(standardsTimeline),
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
	// Extract isArchived from timeline (use first entry or default to false)
	isArchived := false
	if len(isArchivedTimeline) > 0 {
		isArchived = isArchivedTimeline[0].IsArchived
	}
	msg := &types.MsgSetIsArchived{
		Creator:      creator,
		CollectionId: collectionId,
		IsArchived:   isArchived,
	}
	_, err := suite.msgServer.SetIsArchived(ctx, msg)
	return err
}
