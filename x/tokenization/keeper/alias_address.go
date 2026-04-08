package keeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GenerateWrapperPathAddress derives the deterministic wrapper path address from a denom.
// Exposed for tests and any logic that needs the address without storing it on-chain.
func GenerateWrapperPathAddress(denom string) (string, error) {
	accountAddr, err := generatePathAddress(denom, WrapperPathGenerationPrefix)
	if err != nil {
		return "", err
	}
	return accountAddr.String(), nil
}

// MustGenerateWrapperPathAddress is a convenience helper that panics on error (suitable for tests).
func MustGenerateWrapperPathAddress(denom string) string {
	addr, err := GenerateWrapperPathAddress(denom)
	if err != nil {
		panic(fmt.Errorf("failed to generate wrapper path address: %w", err))
	}
	return addr
}

// ResolveAddressAlias resolves a known address alias to its real bb1... address using collection context.
func ResolveAddressAlias(collection *types.TokenCollection, alias string) (string, error) {
	switch {
	case alias == types.MintEscrowAlias:
		if collection.MintEscrowAddress == "" {
			return "", fmt.Errorf("collection has no mint escrow address")
		}
		return collection.MintEscrowAddress, nil

	case strings.HasPrefix(alias, types.CosmosWrapperPrefix):
		idxStr := alias[len(types.CosmosWrapperPrefix):]
		idx, err := strconv.ParseUint(idxStr, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid CosmosWrapper index: %s", idxStr)
		}
		if int(idx) >= len(collection.CosmosCoinWrapperPaths) {
			return "", fmt.Errorf("CosmosWrapper/%d out of range (collection has %d wrapper paths)", idx, len(collection.CosmosCoinWrapperPaths))
		}
		addr := collection.CosmosCoinWrapperPaths[idx].Address
		if addr == "" {
			return "", fmt.Errorf("CosmosWrapper/%d has no address", idx)
		}
		return addr, nil

	case alias == types.IBCBackingAlias:
		if collection.Invariants == nil || collection.Invariants.CosmosCoinBackedPath == nil {
			return "", fmt.Errorf("collection has no IBC backing path")
		}
		addr := collection.Invariants.CosmosCoinBackedPath.Address
		if addr == "" {
			return "", fmt.Errorf("IBC backing path has no address")
		}
		return addr, nil

	default:
		return "", fmt.Errorf("unrecognized address alias: %s", alias)
	}
}

// ResolveAddressIfAlias resolves the address if it's a known alias, otherwise returns it unchanged.
func ResolveAddressIfAlias(collection *types.TokenCollection, address string) (string, error) {
	if types.IsAddressAlias(address) {
		return ResolveAddressAlias(collection, address)
	}
	return address, nil
}

// ResolveListIdAliases resolves any address aliases embedded within a list ID string.
// List IDs can contain aliases in colon-separated format, with AllWithout prefix, or with inversion.
// Reserved keywords (All, AllWithMint, None, Mint) are passed through unchanged.
func ResolveListIdAliases(collection *types.TokenCollection, listId string) (string, error) {
	if listId == "" {
		return listId, nil
	}

	// Preserve and strip inversion prefix
	inversionPrefix := ""
	inner := listId
	if strings.HasPrefix(inner, "!(") && strings.HasSuffix(inner, ")") {
		inversionPrefix = "!("
		inner = inner[2 : len(inner)-1]
		// Will re-add ")" at end
	} else if strings.HasPrefix(inner, "!") {
		inversionPrefix = "!"
		inner = inner[1:]
	}

	// Check for reserved keywords that should not be resolved
	if inner == "All" || inner == "AllWithMint" || inner == "None" || inner == "Mint" {
		return listId, nil // Return original with inversion prefix intact
	}

	// Handle AllWithout prefix
	allWithoutPrefix := ""
	if strings.HasPrefix(inner, "AllWithout") {
		allWithoutPrefix = "AllWithout"
		inner = inner[len("AllWithout"):]
	}

	// Split by colon and resolve each segment
	segments := strings.Split(inner, ":")
	resolved := make([]string, len(segments))
	for i, seg := range segments {
		r, err := ResolveAddressIfAlias(collection, seg)
		if err != nil {
			return "", fmt.Errorf("error resolving alias in list ID segment %q: %w", seg, err)
		}
		resolved[i] = r
	}

	// Reconstruct
	result := allWithoutPrefix + strings.Join(resolved, ":")
	if inversionPrefix == "!(" {
		result = "!(" + result + ")"
	} else if inversionPrefix == "!" {
		result = "!" + result
	}

	return result, nil
}

// resolveApprovalCriteriaAliases resolves aliases in the address fields of ApprovalCriteria.
func resolveApprovalCriteriaAliases(collection *types.TokenCollection, criteria *types.ApprovalCriteria) error {
	if criteria == nil {
		return nil
	}

	// CoinTransfers
	for _, ct := range criteria.CoinTransfers {
		if ct == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, ct.To)
		if err != nil {
			return fmt.Errorf("coinTransfer.To: %w", err)
		}
		ct.To = resolved
	}

	// MustOwnTokens
	for _, mot := range criteria.MustOwnTokens {
		if mot == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, mot.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("mustOwnTokens.OwnershipCheckParty: %w", err)
		}
		mot.OwnershipCheckParty = resolved
	}

	// DynamicStoreChallenges
	for _, dsc := range criteria.DynamicStoreChallenges {
		if dsc == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, dsc.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("dynamicStoreChallenge.OwnershipCheckParty: %w", err)
		}
		dsc.OwnershipCheckParty = resolved
	}

	// UserApprovalSettings.UserRoyalties.PayoutAddress
	if criteria.UserApprovalSettings != nil && criteria.UserApprovalSettings.UserRoyalties != nil {
		resolved, err := ResolveAddressIfAlias(collection, criteria.UserApprovalSettings.UserRoyalties.PayoutAddress)
		if err != nil {
			return fmt.Errorf("userRoyalties.PayoutAddress: %w", err)
		}
		criteria.UserApprovalSettings.UserRoyalties.PayoutAddress = resolved
	}

	return nil
}

// resolveOutgoingCriteriaAliases resolves aliases in OutgoingApprovalCriteria address fields.
func resolveOutgoingCriteriaAliases(collection *types.TokenCollection, criteria *types.OutgoingApprovalCriteria) error {
	if criteria == nil {
		return nil
	}

	for _, ct := range criteria.CoinTransfers {
		if ct == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, ct.To)
		if err != nil {
			return fmt.Errorf("coinTransfer.To: %w", err)
		}
		ct.To = resolved
	}

	for _, mot := range criteria.MustOwnTokens {
		if mot == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, mot.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("mustOwnTokens.OwnershipCheckParty: %w", err)
		}
		mot.OwnershipCheckParty = resolved
	}

	for _, dsc := range criteria.DynamicStoreChallenges {
		if dsc == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, dsc.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("dynamicStoreChallenge.OwnershipCheckParty: %w", err)
		}
		dsc.OwnershipCheckParty = resolved
	}

	return nil
}

// resolveIncomingCriteriaAliases resolves aliases in IncomingApprovalCriteria address fields.
func resolveIncomingCriteriaAliases(collection *types.TokenCollection, criteria *types.IncomingApprovalCriteria) error {
	if criteria == nil {
		return nil
	}

	for _, ct := range criteria.CoinTransfers {
		if ct == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, ct.To)
		if err != nil {
			return fmt.Errorf("coinTransfer.To: %w", err)
		}
		ct.To = resolved
	}

	for _, mot := range criteria.MustOwnTokens {
		if mot == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, mot.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("mustOwnTokens.OwnershipCheckParty: %w", err)
		}
		mot.OwnershipCheckParty = resolved
	}

	for _, dsc := range criteria.DynamicStoreChallenges {
		if dsc == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, dsc.OwnershipCheckParty)
		if err != nil {
			return fmt.Errorf("dynamicStoreChallenge.OwnershipCheckParty: %w", err)
		}
		dsc.OwnershipCheckParty = resolved
	}

	return nil
}

// resolveCollectionApprovalAliases resolves all aliases in a CollectionApproval.
func resolveCollectionApprovalAliases(collection *types.TokenCollection, approval *types.CollectionApproval) error {
	if approval == nil {
		return nil
	}

	// List IDs
	var err error
	approval.FromListId, err = ResolveListIdAliases(collection, approval.FromListId)
	if err != nil {
		return fmt.Errorf("fromListId: %w", err)
	}
	approval.ToListId, err = ResolveListIdAliases(collection, approval.ToListId)
	if err != nil {
		return fmt.Errorf("toListId: %w", err)
	}
	approval.InitiatedByListId, err = ResolveListIdAliases(collection, approval.InitiatedByListId)
	if err != nil {
		return fmt.Errorf("initiatedByListId: %w", err)
	}

	return resolveApprovalCriteriaAliases(collection, approval.ApprovalCriteria)
}

// resolveOutgoingApprovalAliases resolves all aliases in a UserOutgoingApproval.
func resolveOutgoingApprovalAliases(collection *types.TokenCollection, approval *types.UserOutgoingApproval) error {
	if approval == nil {
		return nil
	}

	var err error
	approval.ToListId, err = ResolveListIdAliases(collection, approval.ToListId)
	if err != nil {
		return fmt.Errorf("toListId: %w", err)
	}
	approval.InitiatedByListId, err = ResolveListIdAliases(collection, approval.InitiatedByListId)
	if err != nil {
		return fmt.Errorf("initiatedByListId: %w", err)
	}

	return resolveOutgoingCriteriaAliases(collection, approval.ApprovalCriteria)
}

// resolveIncomingApprovalAliases resolves all aliases in a UserIncomingApproval.
func resolveIncomingApprovalAliases(collection *types.TokenCollection, approval *types.UserIncomingApproval) error {
	if approval == nil {
		return nil
	}

	var err error
	approval.FromListId, err = ResolveListIdAliases(collection, approval.FromListId)
	if err != nil {
		return fmt.Errorf("fromListId: %w", err)
	}
	approval.InitiatedByListId, err = ResolveListIdAliases(collection, approval.InitiatedByListId)
	if err != nil {
		return fmt.Errorf("initiatedByListId: %w", err)
	}

	return resolveIncomingCriteriaAliases(collection, approval.ApprovalCriteria)
}

// resolveUserBalanceStoreAliases resolves all aliases in a UserBalanceStore (used for default balances).
func resolveUserBalanceStoreAliases(collection *types.TokenCollection, store *types.UserBalanceStore) error {
	if store == nil {
		return nil
	}

	for _, approval := range store.OutgoingApprovals {
		if err := resolveOutgoingApprovalAliases(collection, approval); err != nil {
			return fmt.Errorf("defaultBalances.outgoingApprovals: %w", err)
		}
	}
	for _, approval := range store.IncomingApprovals {
		if err := resolveIncomingApprovalAliases(collection, approval); err != nil {
			return fmt.Errorf("defaultBalances.incomingApprovals: %w", err)
		}
	}

	return nil
}

// ResolveAllMsgAliases resolves all address aliases in a MsgUniversalUpdateCollection.
// Must be called after the collection is fetched/created so addresses can be resolved.
// Modifies the msg in-place, replacing all aliases with real bb1... addresses.
func ResolveAllMsgAliases(collection *types.TokenCollection, msg *types.MsgUniversalUpdateCollection) error {
	// Collection approvals
	for _, approval := range msg.CollectionApprovals {
		if err := resolveCollectionApprovalAliases(collection, approval); err != nil {
			return fmt.Errorf("collectionApprovals: %w", err)
		}
	}

	// Default balances
	if err := resolveUserBalanceStoreAliases(collection, msg.DefaultBalances); err != nil {
		return err
	}

	return nil
}

// ResolveTransferAliases resolves all address aliases in a MsgTransferTokens.
func ResolveTransferAliases(collection *types.TokenCollection, msg *types.MsgTransferTokens) error {
	for _, transfer := range msg.Transfers {
		if transfer == nil {
			continue
		}
		resolved, err := ResolveAddressIfAlias(collection, transfer.From)
		if err != nil {
			return fmt.Errorf("transfer.From: %w", err)
		}
		transfer.From = resolved

		for i, to := range transfer.ToAddresses {
			resolved, err := ResolveAddressIfAlias(collection, to)
			if err != nil {
				return fmt.Errorf("transfer.ToAddresses[%d]: %w", i, err)
			}
			transfer.ToAddresses[i] = resolved
		}
	}
	return nil
}
