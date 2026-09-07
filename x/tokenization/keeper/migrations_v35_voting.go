package keeper

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/pkg/storewalk"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Voting scopes move before balances, so collision checks read both original
// balance namespaces without retaining a full-store rewrite list.
func (k Keeper) migrateV35VotingApproverKeys(ctx sdk.Context, store storetypes.KVStore) error {
	hasVotes := func(scope string) bool {
		prefix := storeKey(VotingTrackerKey, scope+BalanceKeyDelimiter)
		iterator := storetypes.KVStorePrefixIterator(store, prefix)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			if !strings.Contains(string(iterator.Key()[len(prefix):]), BalanceKeyDelimiter) {
				return true
			}
		}
		return false
	}
	moveScope := func(oldScope, newScope string) error {
		oldTracker := votingChallengeTrackerStoreKey(oldScope)
		newTracker := votingChallengeTrackerStoreKey(newScope)
		reset := store.Has(newTracker) || hasVotes(newScope)
		parts := strings.SplitN(newScope, BalanceKeyDelimiter, 3)
		if len(parts) == 3 {
			canonical := storeKey(UserBalanceKey, parts[0]+BalanceKeyDelimiter+parts[1])
			upper := storeKey(UserBalanceKey, parts[0]+BalanceKeyDelimiter+strings.ToUpper(parts[1]))
			reset = reset || (store.Has(canonical) && store.Has(upper))
		}
		prefix := storeKey(VotingTrackerKey, oldScope+BalanceKeyDelimiter)
		if err := storewalk.Prefix(ctx, store, prefix, func(key, value []byte) error {
			voter := string(key[len(prefix):])
			if strings.Contains(voter, BalanceKeyDelimiter) {
				return nil
			}
			newKey := storeKey(VotingTrackerKey, newScope+BalanceKeyDelimiter+voter)
			if !store.Has(newKey) {
				store.Set(newKey, value)
			}
			store.Delete(key)
			return nil
		}); err != nil {
			return err
		}
		if store.Has(oldTracker) {
			if !store.Has(newTracker) {
				store.Set(newTracker, store.Get(oldTracker))
			}
			store.Delete(oldTracker)
		}
		if reset {
			store.Set(newTracker, k.cdc.MustMarshal(&types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.ZeroUint()}))
		}
		return nil
	}
	for _, prefix := range [][]byte{VotingTrackerKey, VotingChallengeTrackerKey} {
		if err := storewalk.Prefix(ctx, store, prefix, func(key, value []byte) error {
			oldScope := string(key[len(prefix):])
			if prefix[0] == VotingTrackerKey[0] {
				last := strings.LastIndex(oldScope, BalanceKeyDelimiter)
				if last < 0 {
					return nil
				}
				oldScope = oldScope[:last]
			}
			newScope, changed := rewriteV35DelimitedKey(oldScope, 1)
			if !changed {
				return nil
			}
			return moveScope(oldScope, newScope)
		}); err != nil {
			return err
		}
	}
	return nil
}
