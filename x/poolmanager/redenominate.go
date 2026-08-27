package poolmanager

import (
	"fmt"
	"strings"

	"github.com/cosmos/gogoproto/proto"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// RedenominationResult reports what RedenominateTakerFeeState moved.
type RedenominationResult struct {
	// Trackers counts the per-denom entries in the three taker-fee ProtoRev
	// coin arrays (stakers, community pool, burn) that were re-keyed and
	// rescaled.
	Trackers int
	// SkimAccruals counts the taker-fee-share accrual entries re-keyed. An
	// entry is keyed by *two* denoms — the share denom and the charged denom —
	// and either or both can be the one being redenominated.
	SkimAccruals int
	// ShareAgreements counts the TakerFeeShareAgreement records re-keyed.
	ShareAgreements int
}

// RedenominateTakerFeeState moves every piece of taker-fee accounting that
// names a denom from the legacy denom to the new one, scaling the amounts.
//
// This lives in the module rather than in the upgrade handler because all three
// pieces need raw store access that the keeper deliberately does not export:
// the public API can only *increase* a tracker, never delete or set one, so a
// re-key is not expressible through it.
//
// Three pieces, each with the denom in the key *and* an amount in the value:
//
//  1. The ProtoRev taker-fee trackers (keys.go KeyTakerFee{Stakers,CommunityPool,Burn}ProtoRevArray).
//     Key is prefix+denom, value is a marshalled Int. Written on every swap.
//     Orphaned, the BADGE fees accrued since the last accounting height sit
//     under a denom nothing reads and epoch distribution pays out zero.
//
//  2. The taker-fee-share skim accruals (TakerFeeSkimAccrualPrefix).
//     Key is prefix|shareDenom|chargedDenom, value a marshalled sdk.IntProto.
//     Either denom position can be the redenominated one; the *amount* is an
//     amount of the charged denom, so it scales only when the charged denom is
//     the one moving.
//
//  3. The TakerFeeShareAgreement records (KeyTakerFeeShare). Rename only — the
//     value is a percentage, not an amount. But the keeper caches this map in
//     memory at construction, so re-keying the store without invalidating the
//     cache leaves the agreement resolvable under the retired denom and
//     invisible under the live one for the rest of the process's life.
//
// The receiver is a value. That is safe for the cache repair below only because
// the repair edits the map's contents rather than reassigning the field — see
// the comment there.
func (k Keeper) RedenominateTakerFeeState(
	ctx sdk.Context,
	legacyDenom, newDenom string,
	factor osmomath.Int,
) (RedenominationResult, error) {
	var res RedenominationResult
	store := ctx.KVStore(k.storeKey)

	// --- 1. ProtoRev taker-fee trackers ---
	//
	// Key is the prefix followed by the raw denom bytes, value a marshalled
	// Int (see osmoutils.IncreaseCoinByDenomFromPrefix).
	for _, prefix := range [][]byte{
		types.KeyTakerFeeStakersProtoRevArray,
		types.KeyTakerFeeCommunityPoolProtoRevArray,
		types.KeyTakerFeeBurnProtoRevArray,
	} {
		oldKey := append(append([]byte(nil), prefix...), []byte(legacyDenom)...)
		bz := store.Get(oldKey)
		if bz == nil {
			continue
		}

		var amount osmomath.Int
		if err := amount.Unmarshal(bz); err != nil {
			return res, fmt.Errorf("unmarshalling taker fee tracker for %s: %w", legacyDenom, err)
		}

		newKey := append(append([]byte(nil), prefix...), []byte(newDenom)...)
		// A tracker may already exist under the new denom on a re-run; add
		// rather than overwrite so nothing is lost either way.
		if existing := store.Get(newKey); existing != nil {
			var already osmomath.Int
			if err := already.Unmarshal(existing); err != nil {
				return res, fmt.Errorf("unmarshalling taker fee tracker for %s: %w", newDenom, err)
			}
			amount = amount.Mul(factor).Add(already)
		} else {
			amount = amount.Mul(factor)
		}

		scaled, err := amount.Marshal()
		if err != nil {
			return res, fmt.Errorf("marshalling taker fee tracker for %s: %w", newDenom, err)
		}
		store.Delete(oldKey)
		store.Set(newKey, scaled)
		res.Trackers++
	}

	// --- 2. Taker-fee-share skim accruals ---
	type accrualRewrite struct {
		oldKey, newKey []byte
		value          osmomath.Int
	}
	var accruals []accrualRewrite

	iter := storetypes.KVStorePrefixIterator(store, types.TakerFeeSkimAccrualPrefix)
	for ; iter.Valid(); iter.Next() {
		key := append([]byte(nil), iter.Key()...)
		parts := strings.Split(string(key), types.KeySeparator)
		if len(parts) != 3 {
			ctx.Logger().Error("v35: skipping taker fee skim accrual with an unrecognised key shape",
				"key", string(key))
			continue
		}
		shareDenom, chargedDenom := parts[1], parts[2]
		if shareDenom != legacyDenom && chargedDenom != legacyDenom {
			continue
		}

		var accrued sdk.IntProto
		if err := proto.Unmarshal(iter.Value(), &accrued); err != nil {
			iter.Close()
			return res, fmt.Errorf("unmarshalling skim accrual %s: %w", string(key), err)
		}

		amount := accrued.Int
		// The stored value is an amount *of the charged denom*. It moves only
		// when the charged denom is the one being redenominated; when only the
		// share denom moves this is a pure re-key.
		if chargedDenom == legacyDenom {
			amount = amount.Mul(factor)
			chargedDenom = newDenom
		}
		if shareDenom == legacyDenom {
			shareDenom = newDenom
		}

		accruals = append(accruals, accrualRewrite{
			oldKey: key,
			newKey: types.KeyTakerFeeShareDenomAccrualForTakerFeeChargedDenom(shareDenom, chargedDenom),
			value:  amount,
		})
	}
	iter.Close()

	for _, a := range accruals {
		bz, err := proto.Marshal(&sdk.IntProto{Int: a.value})
		if err != nil {
			return res, fmt.Errorf("marshalling skim accrual %s: %w", string(a.newKey), err)
		}
		store.Delete(a.oldKey)
		store.Set(a.newKey, bz)
		res.SkimAccruals++
	}

	// --- 3. TakerFeeShareAgreement records ---
	oldAgreementKey := types.FormatTakerFeeShareAgreementKey(legacyDenom)
	if bz := store.Get(oldAgreementKey); bz != nil {
		var agreement types.TakerFeeShareAgreement
		if err := proto.Unmarshal(bz, &agreement); err != nil {
			return res, fmt.Errorf("unmarshalling taker fee share agreement for %s: %w", legacyDenom, err)
		}
		agreement.Denom = newDenom
		reencoded, err := proto.Marshal(&agreement)
		if err != nil {
			return res, fmt.Errorf("marshalling taker fee share agreement for %s: %w", newDenom, err)
		}
		store.Delete(oldAgreementKey)
		store.Set(types.FormatTakerFeeShareAgreementKey(newDenom), reencoded)
		res.ShareAgreements++

		// Repair the in-memory cache in place. getTakerFeeShareAgreementFromDenom
		// reads *only* the cache, so a store that has been re-keyed while the
		// cache still holds the retired denom leaves every BADGE swap resolving
		// no share agreement for as long as the process lives.
		//
		// Mutating the map's contents rather than reassigning the field is
		// deliberate on two counts. The keeper is held by value in several
		// places, so a reassignment would only be seen by this copy, while a map
		// is a reference and an in-place edit is seen by all of them. And
		// BeginBlock repopulates the cache from the store only while it is
		// empty, so seeding an entry into an empty cache would suppress that
		// and leave the rest of the map unpopulated — hence the guard below.
		if len(k.cachedTakerFeeShareAgreementMap) > 0 {
			delete(k.cachedTakerFeeShareAgreementMap, legacyDenom)
			k.cachedTakerFeeShareAgreementMap[newDenom] = agreement
		}
	}

	ctx.Logger().Info(
		"v35: redenominated poolmanager taker fee state",
		"from", legacyDenom,
		"to", newDenom,
		"trackers", res.Trackers,
		"skim_accruals", res.SkimAccruals,
		"share_agreements", res.ShareAgreements,
	)

	return res, nil
}
