package keeper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

func canonicalV35Address(address *string) {
	*address, _ = canonicalBech32(*address)
}

func (k Keeper) canonicalV35List(ctx sdk.Context, list *string) {
	_, cleaned, _ := parseInversionPattern(*list, true)
	if _, found := k.GetAddressListFromStore(ctx, cleaned); found {
		return
	}
	body := cleaned
	prefix := ""
	if strings.HasPrefix(body, "AllWithout") {
		prefix = "AllWithout"
		body = strings.TrimPrefix(body, prefix)
	}
	addresses := strings.Split(body, ":")
	for i := range addresses {
		canonicalV35Address(&addresses[i])
		if types.ValidateAddress(addresses[i], true) != nil {
			return
		}
	}
	replacement := prefix + strings.Join(addresses, ":")
	if strings.HasPrefix(*list, "!(") && strings.HasSuffix(*list, ")") {
		*list = "!(" + replacement + ")"
	} else if strings.HasPrefix(*list, "!") {
		*list = "!" + replacement
	} else {
		*list = replacement
	}
}

type v35AddressCriteria interface {
	GetMustOwnTokens() []*types.MustOwnTokens
	GetDynamicStoreChallenges() []*types.DynamicStoreChallenge
	GetCoinTransfers() []*types.CoinTransfer
	GetVotingChallenges() []*types.VotingChallenge
	GetEvmQueryChallenges() []*types.EVMQueryChallenge
}

func canonicalV35Criteria(criteria v35AddressCriteria) error {
	for _, c := range criteria.GetMustOwnTokens() {
		if c != nil {
			canonicalV35Address(&c.OwnershipCheckParty)
		}
	}
	for _, c := range criteria.GetDynamicStoreChallenges() {
		if c != nil {
			canonicalV35Address(&c.OwnershipCheckParty)
		}
	}
	for _, c := range criteria.GetCoinTransfers() {
		if c != nil {
			canonicalV35Address(&c.To)
		}
	}
	for _, c := range criteria.GetEvmQueryChallenges() {
		if c != nil {
			canonicalV35Address(&c.ContractAddress)
		}
	}
	for _, c := range criteria.GetVotingChallenges() {
		if c == nil {
			continue
		}
		seen := map[string]bool{}
		for _, voter := range c.Voters {
			if voter == nil {
				continue
			}
			canonicalV35Address(&voter.Address)
			if seen[voter.Address] {
				return fmt.Errorf("duplicate canonical voter in proposal %s", c.ProposalId)
			}
			seen[voter.Address] = true
		}
	}
	return nil
}

func (k Keeper) canonicalV35UserBalance(ctx sdk.Context, balance *types.UserBalanceStore) error {
	if balance == nil {
		return nil
	}
	for _, a := range balance.IncomingApprovals {
		if a == nil {
			continue
		}
		k.canonicalV35List(ctx, &a.FromListId)
		k.canonicalV35List(ctx, &a.InitiatedByListId)
		if err := canonicalV35Criteria(a.ApprovalCriteria); err != nil {
			return err
		}
	}
	for _, a := range balance.OutgoingApprovals {
		if a == nil {
			continue
		}
		k.canonicalV35List(ctx, &a.ToListId)
		k.canonicalV35List(ctx, &a.InitiatedByListId)
		if err := canonicalV35Criteria(a.ApprovalCriteria); err != nil {
			return err
		}
	}
	if p := balance.UserPermissions; p != nil {
		for _, a := range p.CanUpdateIncomingApprovals {
			if a == nil {
				continue
			}
			k.canonicalV35List(ctx, &a.FromListId)
			k.canonicalV35List(ctx, &a.InitiatedByListId)
		}
		for _, a := range p.CanUpdateOutgoingApprovals {
			if a == nil {
				continue
			}
			k.canonicalV35List(ctx, &a.ToListId)
			k.canonicalV35List(ctx, &a.InitiatedByListId)
		}
	}
	return nil
}

func (k Keeper) canonicalV35Collection(ctx sdk.Context, collection *types.TokenCollection) error {
	canonicalV35Address(&collection.Manager)
	canonicalV35Address(&collection.CreatedBy)
	canonicalV35Address(&collection.MintEscrowAddress)
	for _, path := range collection.CosmosCoinWrapperPaths {
		if path != nil {
			canonicalV35Address(&path.Address)
		}
	}
	if inv := collection.Invariants; inv != nil {
		if inv.CosmosCoinBackedPath != nil {
			canonicalV35Address(&inv.CosmosCoinBackedPath.Address)
		}
		for _, challenge := range inv.EvmQueryChallenges {
			if challenge != nil {
				canonicalV35Address(&challenge.ContractAddress)
			}
		}
	}
	for _, a := range collection.CollectionApprovals {
		if a == nil {
			continue
		}
		k.canonicalV35List(ctx, &a.FromListId)
		k.canonicalV35List(ctx, &a.ToListId)
		k.canonicalV35List(ctx, &a.InitiatedByListId)
		if err := canonicalV35Criteria(a.ApprovalCriteria); err != nil {
			return err
		}
		if settings := a.ApprovalCriteria.GetUserApprovalSettings(); settings != nil && settings.UserRoyalties != nil {
			canonicalV35Address(&settings.UserRoyalties.PayoutAddress)
		}
	}
	if p := collection.CollectionPermissions; p != nil {
		for _, a := range p.CanUpdateCollectionApprovals {
			if a == nil {
				continue
			}
			k.canonicalV35List(ctx, &a.FromListId)
			k.canonicalV35List(ctx, &a.ToListId)
			k.canonicalV35List(ctx, &a.InitiatedByListId)
		}
	}
	return k.canonicalV35UserBalance(ctx, collection.DefaultBalances)
}

func (k Keeper) migrateV35AddressValues(ctx sdk.Context, store storetypes.KVStore) error {
	for _, keyPrefix := range [][]byte{CollectionKey, UserBalanceKey, DynamicStoreKey, VotingTrackerKey} {
		iterator := storetypes.KVStorePrefixIterator(store, keyPrefix)
		type update struct{ key, value []byte }
		updates := []update{}
		for ; iterator.Valid(); iterator.Next() {
			var value proto.Message
			var err error
			switch {
			case bytes.Equal(keyPrefix, CollectionKey):
				collection := new(types.TokenCollection)
				k.cdc.MustUnmarshal(iterator.Value(), collection)
				err = k.canonicalV35Collection(ctx, collection)
				value = collection
			case bytes.Equal(keyPrefix, UserBalanceKey):
				balance := new(types.UserBalanceStore)
				k.cdc.MustUnmarshal(iterator.Value(), balance)
				err = k.canonicalV35UserBalance(ctx, balance)
				value = balance
			case bytes.Equal(keyPrefix, DynamicStoreKey):
				dynamicStore := new(types.DynamicStore)
				k.cdc.MustUnmarshal(iterator.Value(), dynamicStore)
				canonicalV35Address(&dynamicStore.CreatedBy)
				value = dynamicStore
			case bytes.Equal(keyPrefix, VotingTrackerKey):
				vote := new(types.VoteProof)
				k.cdc.MustUnmarshal(iterator.Value(), vote)
				canonicalV35Address(&vote.Voter)
				value = vote
			}
			if err != nil {
				iterator.Close()
				return err
			}
			updated := k.cdc.MustMarshal(value)
			if !bytes.Equal(updated, iterator.Value()) {
				updates = append(updates, update{append([]byte(nil), iterator.Key()...), updated})
			}
		}
		iterator.Close()
		for _, update := range updates {
			store.Set(update.key, update.value)
		}
	}
	return nil
}

func (k Keeper) migrateV35VotingEntries(store storetypes.KVStore) error {
	for _, keyPrefix := range [][]byte{VotingTrackerKey, VotingChallengeTrackerKey} {
		positions := []int{1}
		if bytes.Equal(keyPrefix, VotingTrackerKey) {
			positions = append(positions, -1)
		}
		rewrites := collectV35KeyRewrites(store, keyPrefix, func(suffix string) (string, bool) {
			return rewriteV35DelimitedKey(suffix, positions...)
		})
		for _, rw := range rewrites {
			value := store.Get(rw.oldKey)
			if bytes.Equal(keyPrefix, VotingTrackerKey) {
				vote := new(types.VoteProof)
				k.cdc.MustUnmarshal(value, vote)
				canonicalV35Address(&vote.Voter)
				value = k.cdc.MustMarshal(vote)
			}
			if existing := store.Get(rw.newKey); existing != nil && !bytes.Equal(value, existing) {
				return fmt.Errorf("conflicting canonical voting entry %x", rw.newKey)
			}
			store.Set(rw.newKey, value)
			store.Delete(rw.oldKey)
		}
	}
	return nil
}
