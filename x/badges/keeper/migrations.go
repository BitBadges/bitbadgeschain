package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"

	v3types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v2types "github.com/bitbadges/bitbadgeschain/x/badges/types/v2"
)

// MigrateBadgesKeeper migrates the badges keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	// Get all collections
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	// Unchanged: params, nextCollectionId, challengeTrackers, approvalTrackers

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

	return nil
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// First unmarshal into v2 type
		var v2Collection v2types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v2Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v2Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v3 type
		var v3Collection v3types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v3Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		for _, approval := range v3Collection.CollectionApprovals {
			correspondingV2Approval := &v2types.CollectionApproval{}
			for _, v2Approval := range v2Collection.CollectionApprovals {
				if v2Approval.ApprovalId == approval.ApprovalId {
					correspondingV2Approval = v2Approval
					break
				}
			}

			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}

			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &v3types.RecurringOwnershipTimes{
				StartTime:          sdkmath.NewUint(0),
				IntervalLength:     sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp = false
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp = correspondingV2Approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow
		}

		for _, approval := range v3Collection.DefaultBalances.IncomingApprovals {
			correspondingV2Approval := &v2types.UserIncomingApproval{}
			for _, v2Approval := range v2Collection.DefaultBalances.IncomingApprovals {
				if v2Approval.ApprovalId == approval.ApprovalId {
					correspondingV2Approval = v2Approval
					break
				}
			}

			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &v3types.RecurringOwnershipTimes{
				StartTime:          sdkmath.NewUint(0),
				IntervalLength:     sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp = false
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp = correspondingV2Approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow
		}

		for _, approval := range v3Collection.DefaultBalances.OutgoingApprovals {
			correspondingV2Approval := &v2types.UserOutgoingApproval{}
			for _, v2Approval := range v2Collection.DefaultBalances.OutgoingApprovals {
				if v2Approval.ApprovalId == approval.ApprovalId {
					correspondingV2Approval = v2Approval
					break
				}
			}
			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &v3types.RecurringOwnershipTimes{
				StartTime:          sdkmath.NewUint(0),
				IntervalLength:     sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp = false
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp = correspondingV2Approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v3Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance v2types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into v3 type
		var v3Balance v3types.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &v3Balance); err != nil {
			return err
		}

		for _, approval := range v3Balance.IncomingApprovals {
			correspondingV2Approval := &v2types.UserIncomingApproval{}
			for _, v2Approval := range UserBalance.IncomingApprovals {
				if v2Approval.ApprovalId == approval.ApprovalId {
					correspondingV2Approval = v2Approval
					break
				}
			}
			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &v3types.RecurringOwnershipTimes{
				StartTime:          sdkmath.NewUint(0),
				IntervalLength:     sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp = false
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp = correspondingV2Approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow
		}

		for _, approval := range v3Balance.OutgoingApprovals {
			correspondingV2Approval := &v2types.UserOutgoingApproval{}
			for _, v2Approval := range UserBalance.OutgoingApprovals {
				if v2Approval.ApprovalId == approval.ApprovalId {
					correspondingV2Approval = v2Approval
					break
				}
			}
			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals = &v3types.ResetTimeIntervals{
				StartTime:      sdkmath.NewUint(0),
				IntervalLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &v3types.RecurringOwnershipTimes{
				StartTime:          sdkmath.NewUint(0),
				IntervalLength:     sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			}
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp = false
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp = correspondingV2Approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v3Balance))
	}
	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v2types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v3 type
	// 	var v3AddressList v3types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v3AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v3AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ApprovalTracker v2types.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

		// Convert to JSON
		jsonBytes, err := json.Marshal(ApprovalTracker)
		if err != nil {
			return err
		}

		// Unmarshal into v3 type
		var v3ApprovalTracker v3types.ApprovalTracker
		if err := json.Unmarshal(jsonBytes, &v3ApprovalTracker); err != nil {
			return err
		}

		wctx := sdk.UnwrapSDKContext(ctx)
		nowUnixMilli := wctx.BlockTime().UnixMilli()
		v3ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v3ApprovalTracker))
	}
	return nil
}
