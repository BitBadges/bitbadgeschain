package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

// DeriveCouncilAddress derives a unique module account address for a council.
func DeriveCouncilAddress(councilId uint64) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("council/%d", councilId))))
}

// councilKey returns the store key for a council by ID.
func councilKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(types.CouncilKeyPrefix, bz...)
}

// nextProposalIdKey returns the store key for the next proposal ID counter for a council.
func nextProposalIdKey(councilId uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, councilId)
	return append(types.NextProposalIdKeyPrefix, bz...)
}

// GetCouncil retrieves a council by ID. Returns (council, found).
func (k Keeper) GetCouncil(ctx sdk.Context, id uint64) (types.Council, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(councilKey(id))
	if bz == nil {
		return types.Council{}, false
	}
	var council types.Council
	if err := json.Unmarshal(bz, &council); err != nil {
		panic(fmt.Sprintf("failed to unmarshal council %d: %v", id, err))
	}
	return council, true
}

// SetCouncil stores a council.
func (k Keeper) SetCouncil(ctx sdk.Context, council types.Council) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(council)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal council %d: %v", council.Id, err))
	}
	store.Set(councilKey(council.Id), bz)
}

// GetNextCouncilId returns the next available council ID and increments the counter.
func (k Keeper) GetNextCouncilId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextCouncilIdKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}

	// Increment and store
	next := make([]byte, 8)
	binary.BigEndian.PutUint64(next, id+1)
	store.Set(types.NextCouncilIdKey, next)

	return id
}

// GetNextProposalId returns the next available proposal ID for a council and increments the counter.
func (k Keeper) GetNextProposalId(ctx sdk.Context, councilId uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nextProposalIdKey(councilId))
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}

	next := make([]byte, 8)
	binary.BigEndian.PutUint64(next, id+1)
	store.Set(nextProposalIdKey(councilId), next)

	return id
}
