package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsBackedOrWrappingPathAddress checks if an address is a cosmos coin wrapper path address or backed path address
func (k Keeper) IsBackedOrWrappingPathAddress(ctx sdk.Context, collection *types.TokenCollection, address string) bool {
	for _, path := range collection.CosmosCoinWrapperPaths {
		if path.Address == address {
			return true
		}
	}
	if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
		if collection.Invariants.CosmosCoinBackedPath.Address == address {
			return true
		}
	}
	return false
}

// IsBackedOrWrappingPathAddress checks if an address is a cosmos coin wrapper path address or backed path address
func (k Keeper) IsSpecialBackedAddress(ctx sdk.Context, collection *types.TokenCollection, address string) bool {
	if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
		if collection.Invariants.CosmosCoinBackedPath.Address == address {
			return true
		}
	}
	return false
}

// validateSpecialAddressTransfer validates that the initiator matches the expected address
// This is shared between wrapper and backed path handlers
func validateSpecialAddressTransfer(
	isSendingToSpecialAddress bool,
	isSendingFromSpecialAddress bool,
	from string,
	to string,
	initiatedBy string,
) error {
	if !isSendingFromSpecialAddress && !isSendingToSpecialAddress {
		return nil // No special wrapping needed
	} else if isSendingToSpecialAddress && isSendingFromSpecialAddress {
		return sdkerrors.Wrapf(ErrNotImplemented, "cannot send to and from special addresses at the same time")
	}

	if isSendingToSpecialAddress && initiatedBy != from {
		return sdkerrors.Wrapf(ErrNotImplemented, "initiator must be the same as the sender when sending to special addresses")
	}

	if isSendingFromSpecialAddress && initiatedBy != to {
		return sdkerrors.Wrapf(ErrNotImplemented, "initiator must be the same as the recipient when sending from special addresses")
	}

	if isSendingToSpecialAddress && types.IsMintAddress(from) {
		return sdkerrors.Wrapf(ErrNotImplemented, "the Mint address cannot perform wrap / unwrap actions")
	}

	if isSendingFromSpecialAddress && types.IsMintAddress(to) {
		return sdkerrors.Wrapf(ErrNotImplemented, "the Mint address cannot perform wrap / unwrap actions")
	}

	return nil
}

// calculateConversionMultiplier calculates the multiplier for converting transfer balances to coin amounts
// This is shared between wrapper and backed path handlers
func (k Keeper) calculateConversionMultiplier(
	ctx sdk.Context,
	transferBalances []*types.Balance,
	conversionBalances []*types.Balance,
) (sdkmath.Uint, error) {
	if len(transferBalances) == 0 || len(transferBalances[0].TokenIds) == 0 {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one token ID")
	}

	if len(transferBalances[0].OwnershipTimes) == 0 {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one ownership time")
	}

	// Little hacky but we only allow evenly divisible conversions so we can simply use the first combination,
	// get the multiplier, and check it
	firstBalances := transferBalances[0]
	firstTokenId := firstBalances.TokenIds[0].Start
	firstOwnershipTime := firstBalances.OwnershipTimes[0].Start
	firstAmount := firstBalances.Amount

	// Find the multiplier for the conversion
	// How many of the conversion rate balances evenly fit into the transfer balances? That is multiplier
	multiplier := sdkmath.NewUint(0)
	for _, balance := range conversionBalances {
		foundTokenId, err := types.SearchUintRangesForUint(firstTokenId, balance.TokenIds)
		if err != nil {
			return sdkmath.NewUint(0), err
		}
		foundOwnershipTime, err := types.SearchUintRangesForUint(firstOwnershipTime, balance.OwnershipTimes)
		if err != nil {
			return sdkmath.NewUint(0), err
		}
		if foundTokenId && foundOwnershipTime {
			multiplier = firstAmount.Quo(balance.Amount)
			break
		}
	}

	if multiplier.IsZero() {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
	}

	conversionBalancesMultiplied := types.DeepCopyBalances(conversionBalances)
	for _, balance := range conversionBalancesMultiplied {
		balance.Amount = balance.Amount.Mul(multiplier)
	}

	// Assert that the conversion is evenly divisible for all combinations, not just our first one we checked
	transferBalancesCopy := types.DeepCopyBalances(transferBalances)
	remainingBalances, err := types.SubtractBalances(ctx, transferBalancesCopy, conversionBalancesMultiplied)
	if err != nil {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(err, "conversion is not evenly divisible")
	}

	if len(remainingBalances) > 0 {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
	}

	return multiplier, nil
}
