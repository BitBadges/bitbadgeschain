package keeper

import (
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleSpecialAddressWrapping processes cosmos coin wrapping/unwrapping for special addresses
func (k Keeper) HandleSpecialAddressWrapping(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transferBalances []*types.Balance,
	from string,
	to string,
	initiatedBy string,
) error {
	// Get denomination information
	denomInfo := &types.CosmosCoinWrapperPath{}
	isSendingToSpecialAddress := false
	isSendingFromSpecialAddress := false

	for _, path := range collection.CosmosCoinWrapperPaths {
		if path.Address == to {
			isSendingToSpecialAddress = true
			denomInfo = path
		}
		if path.Address == from {
			isSendingFromSpecialAddress = true
			denomInfo = path
		}
	}

	if !isSendingFromSpecialAddress && !isSendingToSpecialAddress {
		return nil // No special wrapping needed
	} else if isSendingToSpecialAddress && isSendingFromSpecialAddress {
		return sdkerrors.Wrapf(ErrNotImplemented, "cannot send to and from special addresses at the same time")
	}

	// Validate initiator
	if err := validateSpecialAddressTransfer(isSendingToSpecialAddress, isSendingFromSpecialAddress, from, to, initiatedBy); err != nil {
		return err
	}

	if denomInfo.Denom == "" {
		return sdkerrors.Wrapf(ErrNotImplemented, "no denom info found for %s", denomInfo.Address)
	}

	// Check if cosmos wrapping is allowed for this path
	// If explicity not allowed, throw
	if !denomInfo.AllowCosmosWrapping {
		return sdkerrors.Wrapf(ErrNotImplemented, "cosmos wrapping is not allowed for this wrapper path")
	}

	ibcDenom := denomInfo.Denom
	if len(transferBalances) == 0 || len(transferBalances[0].TokenIds) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one token ID")
	}

	firstTokenId := transferBalances[0].TokenIds[0].Start

	// Check if the denom contains the {id} placeholder
	if strings.Contains(denomInfo.Denom, "{id}") {
		if !denomInfo.AllowOverrideWithAnyValidToken {
			return sdkerrors.Wrapf(ErrInvalidConversion, "allowOverrideWithAnyValidToken is not true for this wrapper path")
		}

		//Throw if balances len != 1
		if len(transferBalances) != 1 {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine token ID for {id} placeholder replacement")
		}

		//Throw if TokenIds len != 1
		if len(transferBalances[0].TokenIds) != 1 {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine token ID for {id} placeholder replacement")
		}

		//Throw if TokenIds[0].Start != TokenIds[0].End
		if !transferBalances[0].TokenIds[0].Start.Equal(transferBalances[0].TokenIds[0].End) {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine token ID for {id} placeholder replacement")
		}

		// For {id} placeholder, we need to replace it with the actual token ID
		// Since we're in a transfer context, we can use the first token ID from the transfer
		ibcDenom = strings.ReplaceAll(denomInfo.Denom, "{id}", firstTokenId.String())
	}

	conversionBalances := types.DeepCopyBalances(denomInfo.Balances)

	// If allowOverrideWithAnyValidToken is true, allow any valid token ID
	if denomInfo.AllowOverrideWithAnyValidToken {
		maxTokenIdForCollection := sdkmath.NewUint(1)
		for _, validRange := range collection.ValidTokenIds {
			if validRange.End.GT(maxTokenIdForCollection) {
				maxTokenIdForCollection = validRange.End
			}
		}

		if firstTokenId.GT(maxTokenIdForCollection) || firstTokenId.LT(sdkmath.NewUint(1)) {
			return sdkerrors.Wrapf(ErrInvalidConversion, "token ID not in valid range for overrideWithAnyValidToken")
		}

		// Perform the override
		for _, balance := range conversionBalances {
			balance.TokenIds = []*types.UintRange{
				{
					Start: firstTokenId,
					End:   firstTokenId,
				},
			}
		}
	}

	// Calculate conversion multiplier
	multiplier, err := k.calculateConversionMultiplier(ctx, collection, transferBalances, conversionBalances)
	if err != nil {
		return err
	}

	// Construct the full IBC denomination
	ibcDenom = WrappedDenomPrefix + collection.CollectionId.String() + ":" + ibcDenom

	bankKeeper := k.bankKeeper
	amountInt := multiplier.BigInt()

	if isSendingToSpecialAddress {
		if types.IsMintAddress(from) {
			return sdkerrors.Wrapf(ErrNotImplemented, "the Mint address cannot perform wrap / unwrap actions")
		}

		userAddressAcc := sdk.MustAccAddressFromBech32(from)

		err = bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}

		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	} else if isSendingFromSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(to)

		err = bankKeeper.SendCoinsFromAccountToModule(ctx, userAddressAcc, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}

		err = bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	}

	return nil
}

// calculateConversionMultiplier calculates the multiplier for converting transfer balances to coin amounts
// This is shared between wrapper and backed path handlers
func (k Keeper) calculateConversionMultiplier(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transferBalances []*types.Balance,
	conversionBalances []*types.Balance,
) (sdkmath.Uint, error) {
	if len(transferBalances) == 0 || len(transferBalances[0].TokenIds) == 0 {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one token ID")
	}

	firstTokenId := transferBalances[0].TokenIds[0].Start

	if len(transferBalances[0].OwnershipTimes) == 0 {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one ownership time")
	}
	firstOwnershipTime := transferBalances[0].OwnershipTimes[0].Start
	firstAmount := transferBalances[0].Amount

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

// validateSpecialAddressTransfer validates that the initiator matches the expected address
// This is shared between wrapper and backed path handlers
func validateSpecialAddressTransfer(
	isSendingToSpecialAddress bool,
	isSendingFromSpecialAddress bool,
	from string,
	to string,
	initiatedBy string,
) error {
	if isSendingToSpecialAddress && initiatedBy != from {
		return sdkerrors.Wrapf(ErrNotImplemented, "initiator must be the same as the sender when sending to special addresses")
	}

	if isSendingFromSpecialAddress && initiatedBy != to {
		return sdkerrors.Wrapf(ErrNotImplemented, "initiator must be the same as the recipient when sending from special addresses")
	}

	return nil
}

// HandleSpecialAddressBacking processes cosmos coin backing/unbacking for special addresses
// This uses bankKeeper.SendCoins instead of minting/burning coins
func (k Keeper) HandleSpecialAddressBacking(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transferBalances []*types.Balance,
	from string,
	to string,
	initiatedBy string,
) error {
	// Get denomination information
	var denomInfo *types.CosmosCoinBackedPath
	isSendingToSpecialAddress := false
	isSendingFromSpecialAddress := false

	if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
		path := collection.Invariants.CosmosCoinBackedPath
		if path.Address == to {
			isSendingToSpecialAddress = true
			denomInfo = path
		}
		if path.Address == from {
			isSendingFromSpecialAddress = true
			denomInfo = path
		}
	}

	if !isSendingFromSpecialAddress && !isSendingToSpecialAddress {
		return nil // No special backing needed
	} else if isSendingToSpecialAddress && isSendingFromSpecialAddress {
		return sdkerrors.Wrapf(ErrNotImplemented, "cannot send to and from special addresses at the same time")
	}

	// Validate initiator
	if err := validateSpecialAddressTransfer(isSendingToSpecialAddress, isSendingFromSpecialAddress, from, to, initiatedBy); err != nil {
		return err
	}

	if denomInfo == nil || denomInfo.IbcDenom == "" {
		if denomInfo != nil {
			return sdkerrors.Wrapf(ErrNotImplemented, "no ibc denom info found for %s", denomInfo.Address)
		}
		return nil // No backed path configured
	}

	ibcDenom := denomInfo.IbcDenom
	conversionBalances := types.DeepCopyBalances(denomInfo.Balances)

	// Calculate conversion multiplier
	multiplier, err := k.calculateConversionMultiplier(ctx, collection, transferBalances, conversionBalances)
	if err != nil {
		return err
	}

	// ibcAmount is validated in ValidateBasic to be non-zero
	ibcAmount := denomInfo.IbcAmount

	// Calculate total IBC coin amount: multiplier * ibcAmount
	totalIbcAmount := multiplier.Mul(ibcAmount)

	bankKeeper := k.bankKeeper
	amountInt := totalIbcAmount.BigInt()

	if isSendingToSpecialAddress {
		if types.IsMintAddress(from) {
			return sdkerrors.Wrapf(ErrNotImplemented, "the Mint address cannot perform backing actions")
		}

		userAddressAcc := sdk.MustAccAddressFromBech32(from)
		specialAddressAcc := sdk.MustAccAddressFromBech32(denomInfo.Address)
		// Send coins from special address to user
		err = bankKeeper.SendCoins(ctx, specialAddressAcc, userAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	} else if isSendingFromSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(to)
		specialAddressAcc := sdk.MustAccAddressFromBech32(denomInfo.Address)

		// Send coins from user to special address
		err = bankKeeper.SendCoins(ctx, userAddressAcc, specialAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	}

	return nil
}

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
