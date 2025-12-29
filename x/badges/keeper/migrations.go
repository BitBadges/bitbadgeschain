package keeper

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v21"
)

// MigrateBadgesKeeper migrates the tokens keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	if err := MigratePools(ctx, k); err != nil {
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

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Set mustPrioritize to false for migrated data
		// This ensures existing approvals continue to work without requiring explicit prioritization
		approval.ApprovalCriteria.MustPrioritize = false

		// Initialize votingChallenges to empty array if nil
		if approval.ApprovalCriteria.VotingChallenges == nil {
			approval.ApprovalCriteria.VotingChallenges = []*newtypes.VotingChallenge{}
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Set mustPrioritize to false for migrated data
		// This ensures existing approvals continue to work without requiring explicit prioritization
		approval.ApprovalCriteria.MustPrioritize = false

		// Initialize votingChallenges to empty array if nil
		if approval.ApprovalCriteria.VotingChallenges == nil {
			approval.ApprovalCriteria.VotingChallenges = []*newtypes.VotingChallenge{}
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Set mustPrioritize to false for migrated data
		// This ensures existing approvals continue to work without requiring explicit prioritization
		approval.ApprovalCriteria.MustPrioritize = false

		// Initialize votingChallenges to empty array if nil
		if approval.ApprovalCriteria.VotingChallenges == nil {
			approval.ApprovalCriteria.VotingChallenges = []*newtypes.VotingChallenge{}
		}
	}

	return collectionApprovals
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

// convertActionPermissions converts old ActionPermission to new ActionPermission
func convertActionPermissions(oldPerms []*oldtypes.ActionPermission) []*newtypes.ActionPermission {
	if oldPerms == nil {
		return nil
	}
	newPerms := make([]*newtypes.ActionPermission, len(oldPerms))
	for i, oldPerm := range oldPerms {
		newPerms[i] = &newtypes.ActionPermission{
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		}
	}
	return newPerms
}

// convertTokenIdsActionPermissions converts old TokenIdsActionPermission to new TokenIdsActionPermission
func convertTokenIdsActionPermissions(oldPerms []*oldtypes.TokenIdsActionPermission) []*newtypes.TokenIdsActionPermission {
	if oldPerms == nil {
		return nil
	}
	newPerms := make([]*newtypes.TokenIdsActionPermission, len(oldPerms))
	for i, oldPerm := range oldPerms {
		newPerms[i] = &newtypes.TokenIdsActionPermission{
			TokenIds:                  convertUintRanges(oldPerm.TokenIds),
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		}
	}
	return newPerms
}

// convertCollectionApprovalPermissions converts old CollectionApprovalPermission to new CollectionApprovalPermission
func convertCollectionApprovalPermissions(oldPerms []*oldtypes.CollectionApprovalPermission) []*newtypes.CollectionApprovalPermission {
	if oldPerms == nil {
		return nil
	}
	newPerms := make([]*newtypes.CollectionApprovalPermission, len(oldPerms))
	for i, oldPerm := range oldPerms {
		newPerms[i] = &newtypes.CollectionApprovalPermission{
			FromListId:                oldPerm.FromListId,
			ToListId:                  oldPerm.ToListId,
			InitiatedByListId:         oldPerm.InitiatedByListId,
			TransferTimes:             convertUintRanges(oldPerm.TransferTimes),
			TokenIds:                  convertUintRanges(oldPerm.TokenIds),
			OwnershipTimes:            convertUintRanges(oldPerm.OwnershipTimes),
			ApprovalId:                oldPerm.ApprovalId,
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		}
	}
	return newPerms
}

// MigratePools iterates through all existing pools and sets their addresses as reserved protocol addresses
// and caches them in the pool address cache
func MigratePools(ctx sdk.Context, k Keeper) error {
	// Iterate through pool IDs from 1 to a reasonable upper bound
	// We check up to 10000 pools - if there are more, they will be handled when created
	// maxPoolId := uint64(10000)
	// for poolId := uint64(1); poolId < maxPoolId; poolId++ {
	// 	pool, err := k.gammKeeper.GetPool(ctx, poolId)
	// 	if err != nil {
	// 		// Pool doesn't exist, continue to next ID
	// 		continue
	// 	}

	// 	// Get pool address
	// 	poolAddress := pool.GetAddress().String()

	// 	// Set pool address as reserved protocol address
	// 	if err := k.SetReservedProtocolAddressInStore(ctx, poolAddress, true); err != nil {
	// 		// Log error but continue - don't fail migration for individual pools
	// 		ctx.Logger().Error(fmt.Sprintf("Failed to set pool %d address as reserved protocol: %v", poolId, err))
	// 		continue
	// 	}

	// 	// Cache the pool address -> pool ID mapping
	// 	k.SetPoolAddressInCache(ctx, poolAddress, poolId)
	// }

	return nil
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()

	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	for ; iterator.Valid(); iterator.Next() {
		var oldCollection oldtypes.TokenCollection
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

		// Migrate timeline fields to simple fields by extracting current value at block time
		// Manager
		newCollection.Manager = getCurrentManagerFromTimeline(blockTime, oldCollection.ManagerTimeline)

		// CollectionMetadata
		newCollection.CollectionMetadata = getCurrentCollectionMetadataFromTimeline(blockTime, oldCollection.CollectionMetadataTimeline)

		// TokenMetadata
		newCollection.TokenMetadata = getCurrentTokenMetadataFromTimeline(blockTime, oldCollection.TokenMetadataTimeline)

		// CustomData
		newCollection.CustomData = getCurrentCustomDataFromTimeline(blockTime, oldCollection.CustomDataTimeline)

		// Standards
		newCollection.Standards = getCurrentStandardsFromTimeline(blockTime, oldCollection.StandardsTimeline)

		// IsArchived
		newCollection.IsArchived = getCurrentIsArchivedFromTimeline(blockTime, oldCollection.IsArchivedTimeline)

		// Migrate permissions - only keep those where current time is in timeline times
		// Work with old types first to access TimelineTimes
		if oldCollection.CollectionPermissions != nil {
			newCollection.CollectionPermissions = MigrateCollectionPermissions(blockTime, oldCollection.CollectionPermissions)
		}

		newCollection.CollectionApprovals = MigrateApprovals(newCollection.CollectionApprovals)
		newCollection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(newCollection.DefaultBalances.IncomingApprovals)
		newCollection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(newCollection.DefaultBalances.OutgoingApprovals)

		// Migrate cosmosCoinWrapperPaths from old format to new format
		// Old format has balances and allowCosmosWrapping, new format has conversion
		// If allowCosmosWrapping is false, convert to aliasPaths instead
		newCollection.CosmosCoinWrapperPaths, newCollection.AliasPaths = MigrateCosmosCoinWrapperPaths(oldCollection.CosmosCoinWrapperPaths)

		// Migrate CollectionInvariants.CosmosCoinBackedPath from old format to new format
		// Old format: has ibcDenom, balances, and ibcAmount
		// New format: has conversion (with sideA containing amount+denom, and sideB containing balances)
		if oldCollection.Invariants != nil && oldCollection.Invariants.CosmosCoinBackedPath != nil {
			// Ensure newCollection.Invariants is initialized
			if newCollection.Invariants == nil {
				newCollection.Invariants = &newtypes.CollectionInvariants{}
			}
			// Migrate the CosmosCoinBackedPath
			newCollection.Invariants.CosmosCoinBackedPath = MigrateCosmosCoinBackedPath(oldCollection.Invariants.CosmosCoinBackedPath)
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &newCollection, true); err != nil {
			return err
		}
	}

	return nil
}

// MigrateCosmosCoinWrapperPaths migrates old format cosmosCoinWrapperPaths to new format
// Old format: has balances and allowCosmosWrapping
// New format: has conversion (with sideA.amount and sideB balances)
// If allowCosmosWrapping is false, convert to aliasPaths instead
func MigrateCosmosCoinWrapperPaths(oldPaths []*oldtypes.CosmosCoinWrapperPath) ([]*newtypes.CosmosCoinWrapperPath, []*newtypes.AliasPath) {
	var newWrapperPaths []*newtypes.CosmosCoinWrapperPath
	var newAliasPaths []*newtypes.AliasPath

	for _, oldPath := range oldPaths {
		if oldPath == nil {
			continue
		}

		// Convert balances to new format
		// Old format has balances directly, new format has conversion.sideB
		var sideBBalances []*newtypes.Balance
		if oldPath.Balances != nil {
			sideBBalances = convertBalances(oldPath.Balances)
		}

		// Set amount to 1 if not already set (default to 1)
		// Note: In old format, there was no explicit amount field, so we default to 1
		amount := newtypes.Uint(sdkmath.NewUint(1))

		// Create conversion
		conversion := &newtypes.ConversionWithoutDenom{
			SideA: &newtypes.ConversionSideA{
				Amount: amount,
			},
			SideB: sideBBalances,
		}

		// Migrate DenomUnits to include metadata field
		var newDenomUnits []*newtypes.DenomUnit
		if oldPath.DenomUnits != nil {
			newDenomUnits = convertDenomUnits(oldPath.DenomUnits)
		}

		// If allowCosmosWrapping is false, convert to aliasPath
		if !oldPath.AllowCosmosWrapping {
			aliasPath := &newtypes.AliasPath{
				Denom:      oldPath.Denom,
				Conversion: conversion,
				Symbol:     oldPath.Symbol,
				DenomUnits: newDenomUnits,
				Metadata:   nil, // Old format didn't have metadata, set to nil
			}
			newAliasPaths = append(newAliasPaths, aliasPath)
		} else {
			// Convert to new CosmosCoinWrapperPath format
			wrapperPath := &newtypes.CosmosCoinWrapperPath{
				Address:                        oldPath.Address,
				Denom:                          oldPath.Denom,
				Conversion:                     conversion,
				Symbol:                         oldPath.Symbol,
				DenomUnits:                     newDenomUnits,
				AllowOverrideWithAnyValidToken: oldPath.AllowOverrideWithAnyValidToken,
				Metadata:                       nil, // Old format didn't have metadata, set to nil
			}
			newWrapperPaths = append(newWrapperPaths, wrapperPath)
		}
	}

	return newWrapperPaths, newAliasPaths
}

// convertBalances converts old Balance format to new Balance format
func convertBalances(oldBalances []*oldtypes.Balance) []*newtypes.Balance {
	if oldBalances == nil {
		return nil
	}
	var newBalances []*newtypes.Balance
	for _, oldBalance := range oldBalances {
		if oldBalance == nil {
			continue
		}
		newBalances = append(newBalances, &newtypes.Balance{
			Amount:         newtypes.Uint(oldBalance.Amount),
			OwnershipTimes: convertUintRanges(oldBalance.OwnershipTimes),
			TokenIds:       convertUintRanges(oldBalance.TokenIds),
		})
	}
	return newBalances
}

// convertDenomUnits converts old DenomUnit format to new DenomUnit format (adds metadata field)
func convertDenomUnits(oldUnits []*oldtypes.DenomUnit) []*newtypes.DenomUnit {
	if oldUnits == nil {
		return nil
	}
	newUnits := make([]*newtypes.DenomUnit, len(oldUnits))
	for i, oldUnit := range oldUnits {
		if oldUnit == nil {
			continue
		}
		newUnits[i] = &newtypes.DenomUnit{
			Decimals:         newtypes.Uint(oldUnit.Decimals),
			Symbol:           oldUnit.Symbol,
			IsDefaultDisplay: oldUnit.IsDefaultDisplay,
			Metadata:         nil, // Old format didn't have metadata, set to nil
		}
	}
	return newUnits
}

// MigrateCosmosCoinBackedPath migrates old format CosmosCoinBackedPath to new format
// Old format: has address, ibcDenom, balances, and ibcAmount
// New format: has address and conversion (with sideA containing amount+denom, and sideB containing balances)
func MigrateCosmosCoinBackedPath(oldPath *oldtypes.CosmosCoinBackedPath) *newtypes.CosmosCoinBackedPath {
	if oldPath == nil {
		return nil
	}

	// Convert balances to new format
	var sideBBalances []*newtypes.Balance
	if oldPath.Balances != nil {
		sideBBalances = convertBalances(oldPath.Balances)
	}

	// Convert ibcAmount and ibcDenom to conversion.sideA
	// Old format has ibcAmount (Uint) and ibcDenom (string)
	// New format has conversion.sideA (ConversionSideAWithDenom) with amount and denom
	sideA := &newtypes.ConversionSideAWithDenom{
		Amount: newtypes.Uint(oldPath.IbcAmount),
		Denom:  oldPath.IbcDenom,
	}

	// Create conversion
	conversion := &newtypes.Conversion{
		SideA: sideA,
		SideB: sideBBalances,
	}

	// Create new CosmosCoinBackedPath
	return &newtypes.CosmosCoinBackedPath{
		Address:    oldPath.Address,
		Conversion: conversion,
	}
}

// Helper functions for specific timeline types
func getCurrentManagerFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.ManagerTimeline) string {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				return timelineVal.Manager
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		return timeline[0].Manager
	}
	return ""
}

func getCurrentCustomDataFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.CustomDataTimeline) string {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				return timelineVal.CustomData
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		return timeline[0].CustomData
	}
	return ""
}

func getCurrentStandardsFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.StandardsTimeline) []string {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				return timelineVal.Standards
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		return timeline[0].Standards
	}
	return nil
}

func getCurrentIsArchivedFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.IsArchivedTimeline) bool {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				return timelineVal.IsArchived
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		return timeline[0].IsArchived
	}
	return false
}
func getCurrentCollectionMetadataFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.CollectionMetadataTimeline) *newtypes.CollectionMetadata {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				return &newtypes.CollectionMetadata{
					Uri:        timelineVal.CollectionMetadata.Uri,
					CustomData: timelineVal.CollectionMetadata.CustomData,
				}
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		return &newtypes.CollectionMetadata{
			Uri:        timeline[0].CollectionMetadata.Uri,
			CustomData: timeline[0].CollectionMetadata.CustomData,
		}
	}
	return nil
}

func getCurrentTokenMetadataFromTimeline(blockTime sdkmath.Uint, timeline []*oldtypes.TokenMetadataTimeline) []*newtypes.TokenMetadata {
	for _, timelineVal := range timeline {
		for _, timeRange := range timelineVal.TimelineTimes {
			start := newtypes.Uint(timeRange.Start)
			end := newtypes.Uint(timeRange.End)
			if blockTime.GTE(start) && blockTime.LTE(end) {
				// Convert old TokenMetadata to new TokenMetadata
				result := make([]*newtypes.TokenMetadata, len(timelineVal.TokenMetadata))
				for i, tm := range timelineVal.TokenMetadata {
					result[i] = &newtypes.TokenMetadata{
						Uri:      tm.Uri,
						TokenIds: convertUintRanges(tm.TokenIds),
					}
				}
				return result
			}
		}
	}
	// Return first value if no match found (fallback)
	if len(timeline) > 0 {
		result := make([]*newtypes.TokenMetadata, len(timeline[0].TokenMetadata))
		for i, tm := range timeline[0].TokenMetadata {
			result[i] = &newtypes.TokenMetadata{
				Uri:      tm.Uri,
				TokenIds: convertUintRanges(tm.TokenIds),
			}
		}
		return result
	}
	return nil
}

// MigrateCollectionPermissions filters permissions to only keep those where current time is in timeline times
func MigrateCollectionPermissions(blockTime sdkmath.Uint, oldPerms *oldtypes.CollectionPermissions) *newtypes.CollectionPermissions {
	if oldPerms == nil {
		return nil
	}

	newPerms := &newtypes.CollectionPermissions{}

	// Migrate each permission type - work with old types to access TimelineTimes
	// Convert TimedUpdatePermission to ActionPermission (only keep those where current time is in timeline times)
	if oldPerms.CanUpdateManager != nil {
		newPerms.CanUpdateManager = filterTimedUpdatePermissionsToActionPermissions(blockTime, oldPerms.CanUpdateManager)
	}
	if oldPerms.CanUpdateCollectionMetadata != nil {
		newPerms.CanUpdateCollectionMetadata = filterTimedUpdatePermissionsToActionPermissions(blockTime, oldPerms.CanUpdateCollectionMetadata)
	}
	if oldPerms.CanUpdateTokenMetadata != nil {
		// Convert TimedUpdateWithTokenIdsPermission to TokenIdsActionPermission
		newPerms.CanUpdateTokenMetadata = filterTimedUpdateWithTokenIdsPermissionsToTokenIdsActionPermissions(blockTime, oldPerms.CanUpdateTokenMetadata)
	}
	if oldPerms.CanUpdateCustomData != nil {
		newPerms.CanUpdateCustomData = filterTimedUpdatePermissionsToActionPermissions(blockTime, oldPerms.CanUpdateCustomData)
	}
	if oldPerms.CanUpdateStandards != nil {
		newPerms.CanUpdateStandards = filterTimedUpdatePermissionsToActionPermissions(blockTime, oldPerms.CanUpdateStandards)
	}
	if oldPerms.CanDeleteCollection != nil {
		// ActionPermission doesn't have timeline times, just convert directly
		newPerms.CanDeleteCollection = convertActionPermissions(oldPerms.CanDeleteCollection)
	}
	if oldPerms.CanArchiveCollection != nil {
		newPerms.CanArchiveCollection = filterTimedUpdatePermissionsToActionPermissions(blockTime, oldPerms.CanArchiveCollection)
	}
	if oldPerms.CanUpdateCollectionApprovals != nil {
		// CollectionApprovalPermission doesn't have timeline times in the same way
		newPerms.CanUpdateCollectionApprovals = convertCollectionApprovalPermissions(oldPerms.CanUpdateCollectionApprovals)
	}
	if oldPerms.CanUpdateValidTokenIds != nil {
		// TokenIdsActionPermission doesn't have timeline times
		newPerms.CanUpdateValidTokenIds = convertTokenIdsActionPermissions(oldPerms.CanUpdateValidTokenIds)
	}

	return newPerms
}

// filterTimedUpdatePermissionsToActionPermissions filters permissions to only keep those where current time is in timeline times
// and converts them to ActionPermission (removing timeline times)
func filterTimedUpdatePermissionsToActionPermissions(blockTime sdkmath.Uint, oldPerms []*oldtypes.TimedUpdatePermission) []*newtypes.ActionPermission {
	if oldPerms == nil {
		return nil
	}

	filtered := []*newtypes.ActionPermission{}
	for _, oldPerm := range oldPerms {
		// Check if blockTime is in any of the timeline times
		// Only keep permissions where current time is in timeline times
		if len(oldPerm.TimelineTimes) > 0 {
			found, err := newtypes.SearchUintRangesForUint(blockTime, convertUintRanges(oldPerm.TimelineTimes))
			if err != nil || !found {
				continue // Skip this permission if blockTime is not in timeline times
			}
		}
		// Convert to ActionPermission (removing TimelineTimes)
		filtered = append(filtered, &newtypes.ActionPermission{
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		})
	}
	return filtered
}

// filterTimedUpdateWithTokenIdsPermissionsToTokenIdsActionPermissions filters permissions to only keep those where current time is in timeline times
// and converts them to TokenIdsActionPermission (removing timeline times)
func filterTimedUpdateWithTokenIdsPermissionsToTokenIdsActionPermissions(blockTime sdkmath.Uint, oldPerms []*oldtypes.TimedUpdateWithTokenIdsPermission) []*newtypes.TokenIdsActionPermission {
	if oldPerms == nil {
		return nil
	}

	filtered := []*newtypes.TokenIdsActionPermission{}
	for _, oldPerm := range oldPerms {
		// Check if blockTime is in any of the timeline times
		// Only keep permissions where current time is in timeline times
		if len(oldPerm.TimelineTimes) > 0 {
			found, err := newtypes.SearchUintRangesForUint(blockTime, convertUintRanges(oldPerm.TimelineTimes))
			if err != nil || !found {
				continue // Skip this permission if blockTime is not in timeline times
			}
		}
		// Convert to TokenIdsActionPermission (removing TimelineTimes)
		filtered = append(filtered, &newtypes.TokenIdsActionPermission{
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
			TokenIds:                  convertUintRanges(oldPerm.TokenIds),
		})
	}
	return filtered
}

// MigrateUserPermissions migrates user permissions (these don't have TimelineTimes, so we keep all)
func MigrateUserPermissions(blockTime sdkmath.Uint, oldPerms *oldtypes.UserPermissions) *newtypes.UserPermissions {
	if oldPerms == nil {
		return nil
	}

	newPerms := &newtypes.UserPermissions{}

	// These permissions don't have TimelineTimes, so we keep all and just convert
	if oldPerms.CanUpdateOutgoingApprovals != nil {
		newPerms.CanUpdateOutgoingApprovals = convertUserOutgoingApprovalPermissions(oldPerms.CanUpdateOutgoingApprovals)
	}
	if oldPerms.CanUpdateIncomingApprovals != nil {
		newPerms.CanUpdateIncomingApprovals = convertUserIncomingApprovalPermissions(oldPerms.CanUpdateIncomingApprovals)
	}
	if oldPerms.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers != nil {
		newPerms.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers = convertActionPermissions(oldPerms.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers)
	}
	if oldPerms.CanUpdateAutoApproveSelfInitiatedIncomingTransfers != nil {
		newPerms.CanUpdateAutoApproveSelfInitiatedIncomingTransfers = convertActionPermissions(oldPerms.CanUpdateAutoApproveSelfInitiatedIncomingTransfers)
	}
	if oldPerms.CanUpdateAutoApproveAllIncomingTransfers != nil {
		newPerms.CanUpdateAutoApproveAllIncomingTransfers = convertActionPermissions(oldPerms.CanUpdateAutoApproveAllIncomingTransfers)
	}

	return newPerms
}

// convertUserOutgoingApprovalPermissions converts old UserOutgoingApprovalPermission to new UserOutgoingApprovalPermission
func convertUserOutgoingApprovalPermissions(oldPerms []*oldtypes.UserOutgoingApprovalPermission) []*newtypes.UserOutgoingApprovalPermission {
	if oldPerms == nil {
		return nil
	}
	newPerms := make([]*newtypes.UserOutgoingApprovalPermission, len(oldPerms))
	for i, oldPerm := range oldPerms {
		newPerms[i] = &newtypes.UserOutgoingApprovalPermission{
			ToListId:                  oldPerm.ToListId,
			InitiatedByListId:         oldPerm.InitiatedByListId,
			TransferTimes:             convertUintRanges(oldPerm.TransferTimes),
			TokenIds:                  convertUintRanges(oldPerm.TokenIds),
			OwnershipTimes:            convertUintRanges(oldPerm.OwnershipTimes),
			ApprovalId:                oldPerm.ApprovalId,
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		}
	}
	return newPerms
}

// convertUserIncomingApprovalPermissions converts old UserIncomingApprovalPermission to new UserIncomingApprovalPermission
func convertUserIncomingApprovalPermissions(oldPerms []*oldtypes.UserIncomingApprovalPermission) []*newtypes.UserIncomingApprovalPermission {
	if oldPerms == nil {
		return nil
	}
	newPerms := make([]*newtypes.UserIncomingApprovalPermission, len(oldPerms))
	for i, oldPerm := range oldPerms {
		newPerms[i] = &newtypes.UserIncomingApprovalPermission{
			FromListId:                oldPerm.FromListId,
			InitiatedByListId:         oldPerm.InitiatedByListId,
			TransferTimes:             convertUintRanges(oldPerm.TransferTimes),
			TokenIds:                  convertUintRanges(oldPerm.TokenIds),
			OwnershipTimes:            convertUintRanges(oldPerm.OwnershipTimes),
			ApprovalId:                oldPerm.ApprovalId,
			PermanentlyPermittedTimes: convertUintRanges(oldPerm.PermanentlyPermittedTimes),
			PermanentlyForbiddenTimes: convertUintRanges(oldPerm.PermanentlyForbiddenTimes),
		}
	}
	return newPerms
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	// Get block time from context (convert context.Context to sdk.Context)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockTime := sdkmath.NewUint(uint64(sdkCtx.BlockTime().UnixMilli()))

	for ; iterator.Valid(); iterator.Next() {
		var oldBalance oldtypes.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &oldBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldBalance)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newBalance newtypes.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &newBalance); err != nil {
			return err
		}

		// Migrate approvals
		newBalance.IncomingApprovals = MigrateIncomingApprovals(newBalance.IncomingApprovals)
		newBalance.OutgoingApprovals = MigrateOutgoingApprovals(newBalance.OutgoingApprovals)

		// Migrate permissions - work with old types first to access TimelineTimes if any
		if oldBalance.UserPermissions != nil {
			newBalance.UserPermissions = MigrateUserPermissions(blockTime, oldBalance.UserPermissions)
		}

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

		// Convert defaultValue from Uint to bool if needed
		// Check if the JSON contains a string defaultValue (old Uint format) or bool (new format)
		var jsonData map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &jsonData); err == nil {
			if defaultValue, exists := jsonData["defaultValue"]; exists {
				// If it's a string (old Uint format like "0" or "1"), convert to bool
				if strValue, ok := defaultValue.(string); ok {
					// Parse the string as Uint and convert: "0" -> false, anything else -> true
					// Since Uint is stored as string in JSON, check if it's "0"
					newDynamicStore.DefaultValue = strValue != "0" && strValue != ""
				} else if boolValue, ok := defaultValue.(bool); ok {
					// Already bool, use as-is
					newDynamicStore.DefaultValue = boolValue
				} else if numValue, ok := defaultValue.(float64); ok {
					// If it's a number (shouldn't happen but handle it), convert to bool
					newDynamicStore.DefaultValue = numValue != 0
				}
			}
		}

		// Save the updated dynamic store
		if err := k.SetDynamicStoreInStore(sdk.UnwrapSDKContext(ctx), newDynamicStore); err != nil {
			return err
		}
	}

	// Migrate dynamic store values
	valueIterator := storetypes.KVStorePrefixIterator(store, DynamicStoreValueKey)
	defer valueIterator.Close()
	for ; valueIterator.Valid(); valueIterator.Next() {
		var oldDynamicStoreValue oldtypes.DynamicStoreValue
		k.cdc.MustUnmarshal(valueIterator.Value(), &oldDynamicStoreValue)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldDynamicStoreValue)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newDynamicStoreValue newtypes.DynamicStoreValue
		if err := json.Unmarshal(jsonBytes, &newDynamicStoreValue); err != nil {
			return err
		}

		// Convert value from Uint to bool if needed
		// Check if the JSON contains a string value (old Uint format) or bool (new format)
		var jsonData map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &jsonData); err == nil {
			if value, exists := jsonData["value"]; exists {
				// If it's a string (old Uint format like "0" or "1"), convert to bool
				if strValue, ok := value.(string); ok {
					// Parse the string as Uint and convert: "0" -> false, anything else -> true
					// Since Uint is stored as string in JSON, check if it's "0"
					newDynamicStoreValue.Value = strValue != "0" && strValue != ""
				} else if boolValue, ok := value.(bool); ok {
					// Already bool, use as-is
					newDynamicStoreValue.Value = boolValue
				} else if numValue, ok := value.(float64); ok {
					// If it's a number (shouldn't happen but handle it), convert to bool
					newDynamicStoreValue.Value = numValue != 0
				}
			}
		}

		// Save the updated dynamic store value
		if err := k.SetDynamicStoreValueInStore(sdk.UnwrapSDKContext(ctx), newDynamicStoreValue.StoreId, newDynamicStoreValue.Address, newDynamicStoreValue.Value); err != nil {
			return err
		}
	}

	return nil
}
