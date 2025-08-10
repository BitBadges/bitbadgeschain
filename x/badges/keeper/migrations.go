package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v11"
)

// MigrateBadgesKeeper migrates the badges keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	// Get all collections
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	currParams := k.GetParams(ctx)
	currParams.AffiliatePercentage = sdkmath.NewUint(5000)

	err := k.SetParams(ctx, currParams)
	if err != nil {
		return err
	}

	if err := MigrateCollections(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateBalances(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateAddressLists(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateApprovalTrackers(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateDynamicStores(ctx, store, k); err != nil {
		return err
	}

	return nil
}

func MigrateTokenIdsActionPermission(tokenIdsActionPermission *oldtypes.BadgeIdsActionPermission) *newtypes.TokenIdsActionPermission {
	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(tokenIdsActionPermission)
	if err != nil {
		return nil
	}

	var newTokenIdsActionPermission newtypes.TokenIdsActionPermission
	if err := json.Unmarshal(jsonBytes, &newTokenIdsActionPermission); err != nil {
		return nil
	}

	newTokenIdsActionPermission.TokenIds = convertUintRanges(tokenIdsActionPermission.BadgeIds)

	return &newTokenIdsActionPermission
}

func MigrateTimedUpdateWithTokenIdsPermission(timedUpdateWithTokenIdsPermission *oldtypes.TimedUpdateWithBadgeIdsPermission) *newtypes.TimedUpdateWithTokenIdsPermission {

	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(timedUpdateWithTokenIdsPermission)
	if err != nil {
		return nil
	}

	var newTimedUpdateWithTokenIdsPermission newtypes.TimedUpdateWithTokenIdsPermission
	if err := json.Unmarshal(jsonBytes, &newTimedUpdateWithTokenIdsPermission); err != nil {
		return nil
	}

	newTimedUpdateWithTokenIdsPermission.TokenIds = convertUintRanges(timedUpdateWithTokenIdsPermission.BadgeIds)

	return &newTimedUpdateWithTokenIdsPermission
}

func MigrateCollectionApprovalPermission(collectionApprovalPermission *oldtypes.CollectionApprovalPermission) *newtypes.CollectionApprovalPermission {
	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(collectionApprovalPermission)
	if err != nil {
		return nil
	}

	var newCollectionApprovalPermission newtypes.CollectionApprovalPermission
	if err := json.Unmarshal(jsonBytes, &newCollectionApprovalPermission); err != nil {
		return nil
	}

	newCollectionApprovalPermission.TokenIds = convertUintRanges(collectionApprovalPermission.BadgeIds)

	return &newCollectionApprovalPermission
}

func MigrateUserOutgoingApprovalPermission(userOutgoingApprovalPermission *oldtypes.UserOutgoingApprovalPermission) *newtypes.UserOutgoingApprovalPermission {

	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(userOutgoingApprovalPermission)
	if err != nil {
		return nil
	}

	var newUserOutgoingApprovalPermission newtypes.UserOutgoingApprovalPermission
	if err := json.Unmarshal(jsonBytes, &newUserOutgoingApprovalPermission); err != nil {
		return nil
	}

	return &newUserOutgoingApprovalPermission
}

func MigrateUserIncomingApprovalPermission(userIncomingApprovalPermission *oldtypes.UserIncomingApprovalPermission) *newtypes.UserIncomingApprovalPermission {

	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(userIncomingApprovalPermission)
	if err != nil {
		return nil
	}

	var newUserIncomingApprovalPermission newtypes.UserIncomingApprovalPermission
	if err := json.Unmarshal(jsonBytes, &newUserIncomingApprovalPermission); err != nil {
		return nil
	}

	return &newUserIncomingApprovalPermission
}

func MigrateUserPermissions(userPermissions *oldtypes.UserPermissions) *newtypes.UserPermissions {
	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(userPermissions)
	if err != nil {
		return nil
	}

	var newUserPermissions newtypes.UserPermissions
	if err := json.Unmarshal(jsonBytes, &newUserPermissions); err != nil {
		return nil
	}

	for _, canUpdateOutgoingApprovals := range userPermissions.CanUpdateOutgoingApprovals {
		newUserPermissions.CanUpdateOutgoingApprovals = append(newUserPermissions.CanUpdateOutgoingApprovals, MigrateUserOutgoingApprovalPermission(canUpdateOutgoingApprovals))
	}

	for _, canUpdateIncomingApprovals := range userPermissions.CanUpdateIncomingApprovals {
		newUserPermissions.CanUpdateIncomingApprovals = append(newUserPermissions.CanUpdateIncomingApprovals, MigrateUserIncomingApprovalPermission(canUpdateIncomingApprovals))
	}

	return &newUserPermissions
}

func MigrateCollectionPermissions(collectionPermissions *oldtypes.CollectionPermissions) *newtypes.CollectionPermissions {
	//try to marshal as much as possible
	jsonBytes, err := json.Marshal(collectionPermissions)
	if err != nil {
		return nil
	}

	var newCollectionPermissions newtypes.CollectionPermissions
	if err := json.Unmarshal(jsonBytes, &newCollectionPermissions); err != nil {
		return nil
	}

	for _, canUpdateValidTokenIds := range collectionPermissions.CanUpdateValidBadgeIds {
		newCollectionPermissions.CanUpdateValidTokenIds = append(newCollectionPermissions.CanUpdateValidTokenIds, MigrateTokenIdsActionPermission(canUpdateValidTokenIds))
	}

	for _, canUpdateTokenMetadata := range collectionPermissions.CanUpdateBadgeMetadata {
		newCollectionPermissions.CanUpdateTokenMetadata = append(newCollectionPermissions.CanUpdateTokenMetadata, MigrateTimedUpdateWithTokenIdsPermission(canUpdateTokenMetadata))
	}

	for _, canUpdateCollectionApprovals := range collectionPermissions.CanUpdateCollectionApprovals {
		newCollectionPermissions.CanUpdateCollectionApprovals = append(newCollectionPermissions.CanUpdateCollectionApprovals, MigrateCollectionApprovalPermission(canUpdateCollectionApprovals))
	}

	return &newCollectionPermissions
}

func MigrateTokenMetadata(tokenMetadata *oldtypes.BadgeMetadata) *newtypes.TokenMetadata {
	return &newtypes.TokenMetadata{
		TokenIds:   convertUintRanges(tokenMetadata.BadgeIds),
		Uri:        tokenMetadata.Uri,
		CustomData: tokenMetadata.CustomData,
	}
}

func MigrateTokenMetadataTimeline(tokenMetadataTimeline []*oldtypes.BadgeMetadataTimeline) []*newtypes.TokenMetadataTimeline {

	if tokenMetadataTimeline == nil {
		return nil
	}

	newTokenMetadataTimeline := make([]*newtypes.TokenMetadataTimeline, len(tokenMetadataTimeline))
	for i, tokenMetadata := range tokenMetadataTimeline {
		newTokenMetadataArray := make([]*newtypes.TokenMetadata, len(tokenMetadata.BadgeMetadata))
		for j, metadata := range tokenMetadata.BadgeMetadata {
			newTokenMetadataArray[j] = MigrateTokenMetadata(metadata)
		}

		newTokenMetadataTimeline[i] = &newtypes.TokenMetadataTimeline{
			TimelineTimes: convertUintRanges(tokenMetadata.TimelineTimes),
			TokenMetadata: newTokenMetadataArray,
		}
	}

	return newTokenMetadataTimeline
}

func MigrateBalancesType(balances []*oldtypes.Balance) []*newtypes.Balance {
	if balances == nil {
		return nil
	}

	newBalances := make([]*newtypes.Balance, len(balances))
	for i, balance := range balances {
		newBalances[i] = &newtypes.Balance{
			Amount:         balance.Amount,
			TokenIds:       convertUintRanges(balance.BadgeIds),
			OwnershipTimes: convertUintRanges(balance.OwnershipTimes),
		}
	}

	return newBalances
}

func MigrateCosmosCoinWrapperPath(cosmosCoinWrapperPath *oldtypes.CosmosCoinWrapperPath) *newtypes.CosmosCoinWrapperPath {
	//marshal as much as possible
	jsonBytes, err := json.Marshal(cosmosCoinWrapperPath)
	if err != nil {
		return nil
	}

	var newCosmosCoinWrapperPath newtypes.CosmosCoinWrapperPath
	if err := json.Unmarshal(jsonBytes, &newCosmosCoinWrapperPath); err != nil {
		return nil
	}

	newCosmosCoinWrapperPath.Balances = MigrateBalancesType(cosmosCoinWrapperPath.Balances)
	
	// Set allowOverrideWithAnyValidToken to false for existing paths during migration
	newCosmosCoinWrapperPath.AllowOverrideWithAnyValidToken = false

	return &newCosmosCoinWrapperPath
}

func MigrateMustOwnTokens(mustOwnTokens []*oldtypes.MustOwnBadges) []*newtypes.MustOwnTokens {
	if mustOwnTokens == nil {
		return nil
	}

	newMustOwnTokens := make([]*newtypes.MustOwnTokens, len(mustOwnTokens))
	for i, mustOwnToken := range mustOwnTokens {
		// marshal as much as possible
		jsonBytes, err := json.Marshal(mustOwnToken)
		if err != nil {
			return nil
		}

		var newMustOwnToken newtypes.MustOwnTokens
		if err := json.Unmarshal(jsonBytes, &newMustOwnToken); err != nil {
			return nil
		}

		newMustOwnToken.TokenIds = convertUintRanges(mustOwnToken.BadgeIds)
		newMustOwnTokens[i] = &newMustOwnToken
	}

	return newMustOwnTokens
}

func MigrateApprovalCriteria(approvalCriteria *oldtypes.ApprovalCriteria) *newtypes.ApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	// marshal as much as possible
	jsonBytes, err := json.Marshal(approvalCriteria)
	if err != nil {
		return nil
	}

	var newApprovalCriteria newtypes.ApprovalCriteria
	if err := json.Unmarshal(jsonBytes, &newApprovalCriteria); err != nil {
		return nil
	}

	newApprovalCriteria.MustOwnTokens = MigrateMustOwnTokens(approvalCriteria.MustOwnBadges)
	if newApprovalCriteria.PredeterminedBalances != nil {
		if newApprovalCriteria.PredeterminedBalances.IncrementedBalances != nil {
			newApprovalCriteria.PredeterminedBalances.IncrementedBalances.IncrementTokenIdsBy = approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy
			newApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidToken = approvalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidBadge
			newApprovalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances = MigrateBalancesType(approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances)

		}

		for i, manualBalance := range approvalCriteria.PredeterminedBalances.ManualBalances {
			newApprovalCriteria.PredeterminedBalances.ManualBalances[i].Balances = MigrateBalancesType(manualBalance.Balances)
		}
	}

	return &newApprovalCriteria
}

func MigrateIncomingApprovalCriteria(approvalCriteria *oldtypes.IncomingApprovalCriteria) *newtypes.IncomingApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	// marshal as much as possible
	jsonBytes, err := json.Marshal(approvalCriteria)
	if err != nil {
		return nil
	}

	var newApprovalCriteria newtypes.IncomingApprovalCriteria
	if err := json.Unmarshal(jsonBytes, &newApprovalCriteria); err != nil {
		return nil
	}

	newApprovalCriteria.MustOwnTokens = MigrateMustOwnTokens(approvalCriteria.MustOwnBadges)

	return &newApprovalCriteria
}

func MigrateOutgoingApprovalCriteria(approvalCriteria *oldtypes.OutgoingApprovalCriteria) *newtypes.OutgoingApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	// marshal as much as possible
	jsonBytes, err := json.Marshal(approvalCriteria)
	if err != nil {
		return nil
	}

	var newApprovalCriteria newtypes.OutgoingApprovalCriteria
	if err := json.Unmarshal(jsonBytes, &newApprovalCriteria); err != nil {
		return nil
	}

	newApprovalCriteria.MustOwnTokens = MigrateMustOwnTokens(approvalCriteria.MustOwnBadges)

	return &newApprovalCriteria
}

func MigrateIncomingApprovals(incomingApprovals []*oldtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {

	newIncomingApprovals := make([]*newtypes.UserIncomingApproval, len(incomingApprovals))
	for i, approval := range incomingApprovals {
		// marshal as much as possible
		jsonBytes, err := json.Marshal(approval)
		if err != nil {
			return nil
		}

		var newApproval newtypes.UserIncomingApproval
		if err := json.Unmarshal(jsonBytes, &newApproval); err != nil {
			return nil
		}

		newApproval.ApprovalCriteria = MigrateIncomingApprovalCriteria(approval.ApprovalCriteria)
		newApproval.TokenIds = convertUintRanges(approval.BadgeIds)
		newIncomingApprovals[i] = &newApproval
	}

	// for _, approval := range incomingApprovals {
	// 	if approval.ApprovalCriteria == nil {
	// 		continue
	// 	}

	// 	if approval.ApprovalCriteria.AutoDeletionOptions == nil {
	// 		continue
	// 	}

	// 	approval.ApprovalCriteria.AutoDeletionOptions = &newtypes.AutoDeletionOptions{
	// 		AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
	// 		AfterOverallMaxNumTransfers: false,
	// 	}
	// }

	return newIncomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*oldtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {

	newOutgoingApprovals := make([]*newtypes.UserOutgoingApproval, len(outgoingApprovals))
	for i, approval := range outgoingApprovals {
		// marshal as much as possible
		jsonBytes, err := json.Marshal(approval)
		if err != nil {
			return nil
		}

		var newApproval newtypes.UserOutgoingApproval
		if err := json.Unmarshal(jsonBytes, &newApproval); err != nil {
			return nil
		}

		newApproval.ApprovalCriteria = MigrateOutgoingApprovalCriteria(approval.ApprovalCriteria)
		newApproval.TokenIds = convertUintRanges(approval.BadgeIds)
		newOutgoingApprovals[i] = &newApproval
	}

	return newOutgoingApprovals
}

func MigrateApprovals(collectionApprovals []*oldtypes.CollectionApproval) []*newtypes.CollectionApproval {
	newCollectionApprovals := make([]*newtypes.CollectionApproval, len(collectionApprovals))
	for i, approval := range collectionApprovals {
		// marshal as much as possible
		jsonBytes, err := json.Marshal(approval)
		if err != nil {
			return nil
		}

		var newApproval newtypes.CollectionApproval
		if err := json.Unmarshal(jsonBytes, &newApproval); err != nil {
			return nil
		}

		newApproval.TokenIds = convertUintRanges(approval.BadgeIds)
		newApproval.ApprovalCriteria = MigrateApprovalCriteria(approval.ApprovalCriteria)
		newCollectionApprovals[i] = &newApproval
	}

	return newCollectionApprovals
}

// convertUintRange converts old v9 UintRange to new UintRange
func convertUintRange(oldRange *oldtypes.UintRange) *newtypes.UintRange {
	return &newtypes.UintRange{
		Start: newtypes.Uint(oldRange.Start),
		End:   newtypes.Uint(oldRange.End),
	}
}

// convertUintRanges converts a slice of old v9 UintRange to new UintRange
func convertUintRanges(oldRanges []*oldtypes.UintRange) []*newtypes.UintRange {
	newRanges := make([]*newtypes.UintRange, len(oldRanges))
	for i, oldRange := range oldRanges {
		newRanges[i] = convertUintRange(oldRange)
	}
	return newRanges
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldCollection oldtypes.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &oldCollection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldCollection)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newCollection newtypes.TokenCollection
		if err := json.Unmarshal(jsonBytes, &newCollection); err != nil {
			return err
		}

		newCollection.CollectionPermissions = MigrateCollectionPermissions(oldCollection.CollectionPermissions)
		newCollection.TokenMetadataTimeline = MigrateTokenMetadataTimeline(oldCollection.BadgeMetadataTimeline)
		newCollection.ValidTokenIds = convertUintRanges(oldCollection.ValidBadgeIds)
		newCollection.CollectionApprovals = MigrateApprovals(oldCollection.CollectionApprovals)
		newCollection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(oldCollection.DefaultBalances.IncomingApprovals)
		newCollection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(oldCollection.DefaultBalances.OutgoingApprovals)
		newCollection.DefaultBalances.Balances = MigrateBalancesType(oldCollection.DefaultBalances.Balances)
		newCollection.DefaultBalances.UserPermissions = MigrateUserPermissions(oldCollection.DefaultBalances.UserPermissions)
		for _, cosmosCoinWrapperPath := range oldCollection.CosmosCoinWrapperPaths {
			newCollection.CosmosCoinWrapperPaths = append(newCollection.CosmosCoinWrapperPaths, MigrateCosmosCoinWrapperPath(cosmosCoinWrapperPath))
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &newCollection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance oldtypes.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var newBalance newtypes.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &newBalance); err != nil {
			return err
		}

		newBalance.Balances = MigrateBalancesType(UserBalance.Balances)
		newBalance.IncomingApprovals = MigrateIncomingApprovals(UserBalance.IncomingApprovals)
		newBalance.OutgoingApprovals = MigrateOutgoingApprovals(UserBalance.OutgoingApprovals)
		newBalance.UserPermissions = MigrateUserPermissions(UserBalance.UserPermissions)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&newBalance))
	}

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var AddressList oldtypes.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

		// Convert to JSON
		jsonBytes, err := json.Marshal(AddressList)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var oldAddressList newtypes.AddressList
		if err := json.Unmarshal(jsonBytes, &oldAddressList); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&oldAddressList))
	}

	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ApprovalTracker oldtypes.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

		// Convert to JSON
		jsonBytes, err := json.Marshal(ApprovalTracker)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var oldApprovalTracker newtypes.ApprovalTracker
		if err := json.Unmarshal(jsonBytes, &oldApprovalTracker); err != nil {
			return err
		}

		oldApprovalTracker.Amounts = MigrateBalancesType(ApprovalTracker.Amounts)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&oldApprovalTracker))
	}

	return nil
}

func MigrateDynamicStores(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// Migrate base dynamic stores
	iterator := storetypes.KVStorePrefixIterator(store, DynamicStoreKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldDynamicStore oldtypes.DynamicStore
		k.cdc.MustUnmarshal(iterator.Value(), &oldDynamicStore)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldDynamicStore)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newDynamicStore newtypes.DynamicStore
		if err := json.Unmarshal(jsonBytes, &newDynamicStore); err != nil {
			return err
		}

		// Save the updated dynamic store
		if err := k.SetDynamicStoreInStore(sdk.UnwrapSDKContext(ctx), newDynamicStore); err != nil {
			return err
		}
	}

	return nil
}
