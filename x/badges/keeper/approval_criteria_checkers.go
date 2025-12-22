package keeper

import (
	sdkmath "cosmossdk.io/math"
	approvalcriteria "github.com/bitbadges/bitbadgeschain/x/badges/approval_criteria"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// collectionServiceAdapter adapts the Keeper to the CollectionService interface
type collectionServiceAdapter struct {
	keeper *Keeper
}

func (a *collectionServiceAdapter) GetCollection(ctx sdk.Context, collectionId sdkmath.Uint) (*types.TokenCollection, bool) {
	return a.keeper.GetCollectionFromStore(ctx, collectionId)
}

func (a *collectionServiceAdapter) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool) {
	return a.keeper.GetBalanceOrApplyDefault(ctx, collection, userAddress)
}

// addressCheckServiceAdapter adapts the Keeper to the AddressCheckService interface
type addressCheckServiceAdapter struct {
	keeper *Keeper
}

func (a *addressCheckServiceAdapter) IsWasmContract(ctx sdk.Context, address string) (bool, error) {
	return a.keeper.IsWasmContract(ctx, address)
}

func (a *addressCheckServiceAdapter) IsLiquidityPool(ctx sdk.Context, address string) (bool, error) {
	return a.keeper.IsLiquidityPool(ctx, address)
}

func (a *addressCheckServiceAdapter) IsAddressReservedProtocol(ctx sdk.Context, address string) bool {
	return a.keeper.IsAddressReservedProtocolInStore(ctx, address)
}

// addressListServiceAdapter adapts the Keeper to the AddressListService interface
type addressListServiceAdapter struct {
	keeper *Keeper
}

func (a *addressListServiceAdapter) CheckAddresses(ctx sdk.Context, addressListId string, addressToCheck string) (bool, error) {
	return a.keeper.CheckAddresses(ctx, addressListId, addressToCheck)
}

// dynamicStoreServiceAdapter adapts the Keeper to the DynamicStoreService interface
type dynamicStoreServiceAdapter struct {
	keeper *Keeper
}

func (a *dynamicStoreServiceAdapter) GetDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) (*types.DynamicStore, bool) {
	store, found := a.keeper.GetDynamicStoreFromStore(ctx, storeId)
	if !found {
		return nil, false
	}
	return &store, true
}

func (a *dynamicStoreServiceAdapter) GetDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) (*types.DynamicStoreValue, bool) {
	value, found := a.keeper.GetDynamicStoreValueFromStore(ctx, storeId, address)
	if !found {
		return nil, false
	}
	return &value, true
}

// GetApprovalCriteriaCheckers returns all applicable checkers for the given approval
// This includes basic validation checkers (address matching, transfer times) and approval criteria checkers
func (k Keeper) GetApprovalCriteriaCheckers(approval *types.CollectionApproval) []approvalcriteria.ApprovalCriteriaChecker {
	checkers := []approvalcriteria.ApprovalCriteriaChecker{}

	if approval.ApprovalCriteria == nil {
		return checkers
	}

	approvalCriteria := approval.ApprovalCriteria

	// MustOwnTokens checker
	if len(approvalCriteria.MustOwnTokens) > 0 {
		collectionService := &collectionServiceAdapter{keeper: &k}
		checkers = append(checkers, approvalcriteria.NewMustOwnTokensChecker(collectionService))
	}

	// Address checks for sender
	if approvalCriteria.SenderChecks != nil {
		addressCheckService := &addressCheckServiceAdapter{keeper: &k}
		checkers = append(checkers, approvalcriteria.NewAddressChecksChecker(addressCheckService, approvalCriteria.SenderChecks, "sender"))
	}

	// Address checks for recipient
	if approvalCriteria.RecipientChecks != nil {
		addressCheckService := &addressCheckServiceAdapter{keeper: &k}
		checkers = append(checkers, approvalcriteria.NewAddressChecksChecker(addressCheckService, approvalCriteria.RecipientChecks, "recipient"))
	}

	// Address checks for initiator
	if approvalCriteria.InitiatorChecks != nil {
		addressCheckService := &addressCheckServiceAdapter{keeper: &k}
		checkers = append(checkers, approvalcriteria.NewAddressChecksChecker(addressCheckService, approvalCriteria.InitiatorChecks, "initiator"))
	}

	// AltTimeChecks checker
	if approvalCriteria.AltTimeChecks != nil {
		checkers = append(checkers, approvalcriteria.NewAltTimeChecksChecker(approvalCriteria.AltTimeChecks))
	}

	// Address equality checkers
	if approvalCriteria.RequireFromDoesNotEqualInitiatedBy {
		checkers = append(checkers, approvalcriteria.NewRequireFromDoesNotEqualInitiatedByChecker())
	}

	if approvalCriteria.RequireFromEqualsInitiatedBy {
		checkers = append(checkers, approvalcriteria.NewRequireFromEqualsInitiatedByChecker())
	}

	if approvalCriteria.RequireToDoesNotEqualInitiatedBy {
		checkers = append(checkers, approvalcriteria.NewRequireToDoesNotEqualInitiatedByChecker())
	}

	if approvalCriteria.RequireToEqualsInitiatedBy {
		checkers = append(checkers, approvalcriteria.NewRequireToEqualsInitiatedByChecker())
	}

	// DynamicStoreChallenges checker
	if len(approvalCriteria.DynamicStoreChallenges) > 0 {
		dynamicStoreService := &dynamicStoreServiceAdapter{keeper: &k}
		checkers = append(checkers, approvalcriteria.NewDynamicStoreChallengesChecker(dynamicStoreService))
	}

	// NoForcefulPostMintTransfers checker (always added, will check if invariant is enabled in Check method)
	checkers = append(checkers, approvalcriteria.NewNoForcefulPostMintTransfersChecker())

	// ReservedProtocolAddress checker (always added, will check conditions in Check method)
	addressCheckService := &addressCheckServiceAdapter{keeper: &k}
	checkers = append(checkers, approvalcriteria.NewReservedProtocolAddressChecker(addressCheckService))

	// Append custom checkers registered by developers
	for _, provider := range k.customCheckerProviders {
		customCheckers := provider(approval)
		if customCheckers != nil {
			checkers = append(checkers, customCheckers...)
		}
	}

	return checkers
}
