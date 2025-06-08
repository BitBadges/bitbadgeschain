package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	storetypes "cosmossdk.io/store/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	v5types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v4types "github.com/bitbadges/bitbadgeschain/x/badges/types/v4"
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

func MigrateIncomingApprovals(incomingApprovals []*v5types.UserIncomingApproval) []*v5types.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		for _, challenge := range approval.ApprovalCriteria.MerkleChallenges {
			challenge.LeafSigner = ""
		}

		approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidBadge = false
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*v5types.UserOutgoingApproval) []*v5types.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		for _, challenge := range approval.ApprovalCriteria.MerkleChallenges {
			challenge.LeafSigner = ""
		}

		approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidBadge = false
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*v5types.CollectionApproval) []*v5types.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		for _, challenge := range approval.ApprovalCriteria.MerkleChallenges {
			challenge.LeafSigner = ""
		}

		approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidBadge = false
	}

	return collectionApprovals
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// First unmarshal into v4 type
		var v4Collection v4types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v4Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v4Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v5 type
		var v5Collection v5types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v5Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		v5Collection.CollectionApprovals = MigrateApprovals(v5Collection.CollectionApprovals)
		v5Collection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(v5Collection.DefaultBalances.IncomingApprovals)
		v5Collection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(v5Collection.DefaultBalances.OutgoingApprovals)

		// From cosmos SDK x/group moduleAdd commentMore actions
		// Generate account address of collection
		var accountAddr sdk.AccAddress
		// loop here in the rare case where a ADR-028-derived address creates a
		// collision with an existing address.
		for {
			derivationKey := make([]byte, 8)
			binary.BigEndian.PutUint64(derivationKey, v5Collection.CollectionId.Uint64())

			ac, err := authtypes.NewModuleCredential(types.ModuleName, AccountGenerationPrefix, derivationKey)
			if err != nil {
				return err
			}
			//generate the address from the credential
			accountAddr = sdk.AccAddress(ac.Address())

			break
		}

		v5Collection.MintEscrowAddress = accountAddr.String()

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v5Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance v4types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into v5 type
		var v5Balance v5types.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &v5Balance); err != nil {
			return err
		}

		v5Balance.IncomingApprovals = MigrateIncomingApprovals(v5Balance.IncomingApprovals)
		v5Balance.OutgoingApprovals = MigrateOutgoingApprovals(v5Balance.OutgoingApprovals)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v5Balance))
	}
	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v4types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v5 type
	// 	var v5AddressList v5types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v5AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v5AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker v4types.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v5 type
	// 	var v5ApprovalTracker v5types.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &v5ApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	v5ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v5ApprovalTracker))
	// }
	return nil
}
