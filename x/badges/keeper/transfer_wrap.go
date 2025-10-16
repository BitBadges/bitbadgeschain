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
	collection *types.BadgeCollection,
	transferBalances []*types.Balance,
	from string,
	to string,
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

	if denomInfo.Denom == "" {
		return sdkerrors.Wrapf(ErrNotImplemented, "no denom info found for %s", denomInfo.Address)
	}

	// Check if cosmos wrapping is allowed for this path
	// If explicity not allowed, throw
	if !denomInfo.AllowCosmosWrapping {
		return sdkerrors.Wrapf(ErrNotImplemented, "cosmos wrapping is not allowed for this wrapper path")
	}

	ibcDenom := denomInfo.Denom
	if len(transferBalances) == 0 || len(transferBalances[0].BadgeIds) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one badge ID")
	}

	firstTokenId := transferBalances[0].BadgeIds[0].Start

	// Check if the denom contains the {id} placeholder
	if strings.Contains(denomInfo.Denom, "{id}") {
		if !denomInfo.AllowOverrideWithAnyValidToken {
			return sdkerrors.Wrapf(ErrInvalidConversion, "allowOverrideWithAnyValidToken is not true for this wrapper path")
		}

		//Throw if balances len != 1
		if len(transferBalances) != 1 {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine badge ID for {id} placeholder replacement")
		}

		//Throw if BadgeIds len != 1
		if len(transferBalances[0].BadgeIds) != 1 {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine badge ID for {id} placeholder replacement")
		}

		//Throw if BadgeIds[0].Start != BadgeIds[0].End
		if !transferBalances[0].BadgeIds[0].Start.Equal(transferBalances[0].BadgeIds[0].End) {
			return sdkerrors.Wrapf(ErrInvalidConversion, "cannot determine badge ID for {id} placeholder replacement")
		}

		// For {id} placeholder, we need to replace it with the actual badge ID
		// Since we're in a transfer context, we can use the first token ID from the transfer
		ibcDenom = strings.ReplaceAll(denomInfo.Denom, "{id}", firstTokenId.String())
	}

	conversionBalances := types.DeepCopyBalances(denomInfo.Balances)

	// If allowOverrideWithAnyValidToken is true, allow any valid badge ID
	if denomInfo.AllowOverrideWithAnyValidToken {
		maxBadgeIdForCollection := sdkmath.NewUint(1)
		for _, validRange := range collection.ValidBadgeIds {
			if validRange.End.GT(maxBadgeIdForCollection) {
				maxBadgeIdForCollection = validRange.End
			}
		}

		if firstTokenId.GT(maxBadgeIdForCollection) || firstTokenId.LT(sdkmath.NewUint(1)) {
			return sdkerrors.Wrapf(ErrInvalidConversion, "token ID not in valid range for overrideWithAnyValidToken")
		}

		// Perform the override
		for _, balance := range conversionBalances {
			balance.BadgeIds = []*types.UintRange{
				{
					Start: firstTokenId,
					End:   firstTokenId,
				},
			}
		}
	}

	// Little hacky but we find the amount for a specific time and ID
	// Then we will check if it is evenly divisible by the number of transfer balances
	if len(transferBalances[0].OwnershipTimes) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one ownership time")
	}
	firstOwnershipTime := transferBalances[0].OwnershipTimes[0].Start
	firstAmount := transferBalances[0].Amount

	multiplier := sdkmath.NewUint(0)
	for _, balance := range conversionBalances {
		foundTokenId, err := types.SearchUintRangesForUint(firstTokenId, balance.BadgeIds)
		if err != nil {
			return err
		}
		foundOwnershipTime, err := types.SearchUintRangesForUint(firstOwnershipTime, balance.OwnershipTimes)
		if err != nil {
			return err
		}
		if foundTokenId && foundOwnershipTime {
			multiplier = firstAmount.Quo(balance.Amount)
			break
		}
	}

	if multiplier.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
	}

	conversionBalancesMultiplied := types.DeepCopyBalances(conversionBalances)
	for _, balance := range conversionBalancesMultiplied {
		balance.Amount = balance.Amount.Mul(multiplier)
	}

	transferBalancesCopy := types.DeepCopyBalances(transferBalances)
	remainingBalances, err := types.SubtractBalances(ctx, transferBalancesCopy, conversionBalancesMultiplied)
	if err != nil {
		return sdkerrors.Wrapf(err, "conversion is not evenly divisible")
	}

	if len(remainingBalances) > 0 {
		return sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
	}

	// Construct the full IBC denomination
	badgePrefix := "badges:"
	ibcDenom = badgePrefix + collection.CollectionId.String() + ":" + ibcDenom

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
