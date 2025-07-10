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
	v9types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	// v7types "github.com/bitbadges/bitbadgeschain/x/badges/types/v7"
)

// MigrateBadgesKeeper migrates the badges keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	// Get all collections
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

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

func MigrateIncomingApprovals(incomingApprovals []*v9types.UserIncomingApproval) []*v9types.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v9types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*v9types.UserOutgoingApproval) []*v9types.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v9types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*v9types.CollectionApproval) []*v9types.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v9types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return collectionApprovals
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// TODO: CHANGE THIS BACK
		var v7Collection v9types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v7Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v7Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v7 type
		var v9Collection v9types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v9Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		// v9Collection.CollectionApprovals = MigrateApprovals(v9Collection.CollectionApprovals)
		// v9Collection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(v9Collection.DefaultBalances.IncomingApprovals)
		// v9Collection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(v9Collection.DefaultBalances.OutgoingApprovals)

		// From cosmos SDK x/group moduleAdd commentMore actions
		// Generate account address of collection
		var accountAddr sdk.AccAddress
		// loop here in the rare case where a ADR-028-derived address creates a
		// collision with an existing address.
		for {
			derivationKey := make([]byte, 8)
			binary.BigEndian.PutUint64(derivationKey, v9Collection.CollectionId.Uint64())

			ac, err := authtypes.NewModuleCredential(v9types.ModuleName, AccountGenerationPrefix, derivationKey)
			if err != nil {
				return err
			}
			//generate the address from the credential
			accountAddr = sdk.AccAddress(ac.Address())

			break
		}
		v9Collection.MintEscrowAddress = accountAddr.String()

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v9Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var UserBalance v7types.UserBalanceStore
	// 	k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(UserBalance)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7Balance v9types.UserBalanceStore
	// 	if err := json.Unmarshal(jsonBytes, &v7Balance); err != nil {
	// 		return err
	// 	}

	// 	v7Balance.IncomingApprovals = MigrateIncomingApprovals(v7Balance.IncomingApprovals)
	// 	v7Balance.OutgoingApprovals = MigrateOutgoingApprovals(v7Balance.OutgoingApprovals)

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7Balance))
	// }

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v7types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7AddressList v9types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v7AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker v7types.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7ApprovalTracker v9types.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &v7ApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	v7ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7ApprovalTracker))
	// }
	return nil
}
