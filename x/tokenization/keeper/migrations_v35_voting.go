package keeper

import (
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
)

func collectV35VotingBalanceCollisions(store storetypes.KVStore) map[string]bool {
	collisions := map[string]bool{}
	rewrites := collectV35KeyRewrites(store, UserBalanceKey, func(suffix string) (string, bool) {
		return rewriteV35DelimitedKey(suffix, -1)
	})
	for _, rw := range rewrites {
		if store.Has(rw.newKey) {
			collisions[string(rw.newKey[len(UserBalanceKey):])] = true
		}
	}
	return collisions
}

func (k Keeper) migrateV35VotingApproverKeys(store storetypes.KVStore, balanceCollisions map[string]bool) {
	canonicalScopes := map[string]bool{}
	collect := func(prefix []byte, vote bool) []v35KeyRewrite {
		return collectV35KeyRewrites(store, prefix, func(suffix string) (string, bool) {
			canonical, changed := rewriteV35DelimitedKey(suffix, 1)
			if !changed {
				scope := suffix
				if vote {
					last := strings.LastIndex(scope, BalanceKeyDelimiter)
					if last < 0 {
						return suffix, false
					}
					scope = scope[:last]
				}
				canonicalScopes[scope] = true
			}
			return canonical, changed
		})
	}
	votes := collect(VotingTrackerKey, true)
	trackers := collect(VotingChallengeTrackerKey, false)
	resetScopes := map[string]bool{}
	needsReset := func(scope string) bool {
		parts := strings.SplitN(scope, BalanceKeyDelimiter, 3)
		return canonicalScopes[scope] || (len(parts) == 3 && balanceCollisions[parts[0]+BalanceKeyDelimiter+parts[1]])
	}
	for _, rw := range votes {
		suffix := string(rw.newKey[len(VotingTrackerKey):])
		last := strings.LastIndex(suffix, BalanceKeyDelimiter)
		if last >= 0 && needsReset(suffix[:last]) {
			resetScopes[suffix[:last]] = true
		}
		if !store.Has(rw.newKey) {
			store.Set(rw.newKey, store.Get(rw.oldKey))
		}
		store.Delete(rw.oldKey)
	}
	for _, rw := range trackers {
		scope := string(rw.newKey[len(VotingChallengeTrackerKey):])
		if needsReset(scope) {
			resetScopes[scope] = true
		}
		if !store.Has(rw.newKey) {
			store.Set(rw.newKey, store.Get(rw.oldKey))
		}
		store.Delete(rw.oldKey)
	}
	// A fresh vote must restart the delay after combining approver scopes.
	scopes := make([]string, 0, len(resetScopes))
	for scope := range resetScopes {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	for _, scope := range scopes {
		key := votingChallengeTrackerStoreKey(scope)
		store.Set(key, k.cdc.MustMarshal(&types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.ZeroUint()}))
	}
}
