package keeper

import (
	"bitbadgeschain/x/offers/types"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
)

/****************************************PROPOSALS****************************************/

func (k Keeper) SetProposalInStore(ctx sdk.Context, proposal *types.Proposal) error {
	marshaled_badge, err := k.cdc.Marshal(proposal)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Proposal failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(proposalKey(proposal.Id), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the proposalId.
func (k Keeper) GetProposalFromStore(ctx sdk.Context, proposalId sdkmath.Uint) (*types.Proposal, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_proposal := store.Get(proposalKey(proposalId))

	var proposal types.Proposal
	if len(marshaled_proposal) == 0 {
		return &proposal, false
	}
	k.cdc.MustUnmarshal(marshaled_proposal, &proposal)
	return &proposal, true
}

// GetProposalsFromStore defines a method for returning all badges information by key.
func (k Keeper) GetProposalsFromStore(ctx sdk.Context) (proposals []*types.Proposal) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, ProposalKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		proposals = append(proposals, &proposal)
	}
	return
}

// StoreHasProposalID determines whether the specified proposalId exists
func (k Keeper) StoreHasProposalID(ctx sdk.Context, proposalId sdkmath.Uint) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(proposalKey(proposalId))
}

// DeleteProposalFromStore deletes a badge from the store.
func (k Keeper) DeleteProposalFromStore(ctx sdk.Context, proposalId sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(proposalKey(proposalId))
}

/****************************************NEXT COLLECTION ID****************************************/

// Gets the next proposal ID.
func (k Keeper) GetNextProposalId(ctx sdk.Context) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	nextProposalId := store.Get(nextProposalIdKey())
	nextProposalIdStr := string((nextProposalId))
	nextID := types.NewUintFromString(nextProposalIdStr)
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextProposalId(ctx sdk.Context, nextID sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(nextProposalIdKey(), []byte(nextID.String()))
}

// Increments the next proposal ID by 1.
func (k Keeper) IncrementNextProposalId(ctx sdk.Context) {
	nextID := k.GetNextProposalId(ctx)
	k.SetNextProposalId(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}
