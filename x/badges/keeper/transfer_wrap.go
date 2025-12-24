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

	if !isSendingToSpecialAddress && !isSendingFromSpecialAddress {
		return nil
	}

	if denomInfo == nil || denomInfo.Denom == "" {
		return sdkerrors.Wrapf(ErrNotImplemented, "no denom info found")
	}

	ibcDenom := denomInfo.Denom
	if len(transferBalances) == 0 || len(transferBalances[0].TokenIds) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "transfer balances must contain at least one token ID")
	}

	// Validate initiator
	// Note: We previously verify isPrioritized elsewhere as well
	if err := validateSpecialAddressTransfer(isSendingToSpecialAddress, isSendingFromSpecialAddress, from, to, initiatedBy); err != nil {
		return err
	}

	firstTokenId := transferBalances[0].TokenIds[0].Start

	// Check if the denom contains the {id} placeholder and dynamic replace if so
	// We use the requested token ID from the transfer balances so must be size() == 1 and start == end
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

		maxTokenIdForCollection := sdkmath.NewUint(1)
		for _, validRange := range collection.ValidTokenIds {
			if validRange.End.GT(maxTokenIdForCollection) {
				maxTokenIdForCollection = validRange.End
			}
		}

		if firstTokenId.GT(maxTokenIdForCollection) || firstTokenId.LT(sdkmath.NewUint(1)) {
			return sdkerrors.Wrapf(ErrInvalidConversion, "token ID not in valid range for overrideWithAnyValidToken")
		}

		// For {id} placeholder, we need to replace it with the actual token ID
		// Since we're in a transfer context, we can use the first token ID from the transfer
		ibcDenom = strings.ReplaceAll(denomInfo.Denom, "{id}", firstTokenId.String())
	}

	// If allowOverrideWithAnyValidToken is true, allow any valid token ID
	conversionBalances := types.DeepCopyBalances(denomInfo.Balances)
	if denomInfo.AllowOverrideWithAnyValidToken {
		for _, balance := range conversionBalances {
			balance.TokenIds = []*types.UintRange{
				{
					Start: firstTokenId,
					End:   firstTokenId,
				},
			}
		}
	}

	// Conversion Rate = [{ amount: ibcAmount, denom: ibcDenom }] x 1 -> path.Balances x 1
	multiplier, err := k.calculateConversionMultiplier(ctx, transferBalances, conversionBalances)
	if err != nil {
		return err
	}

	// This is important. Prefixing allows only our module to control such denominations
	// and doesn't allow users to mint other coins like "ubadge" or "factory/..."
	ibcDenom = WrappedDenomPrefix + collection.CollectionId.String() + ":" + ibcDenom
	amountInt := multiplier.Mul(denomInfo.Amount).BigInt()

	if isSendingToSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(from)
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	} else if isSendingFromSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(to)
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddressAcc, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}

		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
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

	if !isSendingToSpecialAddress && !isSendingFromSpecialAddress {
		return nil
	}

	if denomInfo == nil || denomInfo.IbcDenom == "" {
		return sdkerrors.Wrapf(ErrNotImplemented, "no ibc denom info found")
	}

	// Validate initiator
	if err := validateSpecialAddressTransfer(isSendingToSpecialAddress, isSendingFromSpecialAddress, from, to, initiatedBy); err != nil {
		return err
	}

	// Conversion Rate = [{ amount: ibcAmount, denom: ibcDenom }] x 1 -> path.Balances x 1
	ibcDenom := denomInfo.IbcDenom
	conversionBalances := types.DeepCopyBalances(denomInfo.Balances)
	multiplier, err := k.calculateConversionMultiplier(ctx, transferBalances, conversionBalances)
	if err != nil {
		return err
	}

	amountInt := multiplier.Mul(denomInfo.IbcAmount).BigInt()

	if isSendingToSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(from)
		specialAddressAcc := sdk.MustAccAddressFromBech32(denomInfo.Address)
		err = k.bankKeeper.SendCoins(ctx, specialAddressAcc, userAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	} else if isSendingFromSpecialAddress {
		userAddressAcc := sdk.MustAccAddressFromBech32(to)
		specialAddressAcc := sdk.MustAccAddressFromBech32(denomInfo.Address)
		err = k.bankKeeper.SendCoins(ctx, userAddressAcc, specialAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
		if err != nil {
			return err
		}
	}

	return nil
}
