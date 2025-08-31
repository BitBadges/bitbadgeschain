package keeper

import (
	"fmt"
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckStartsWithBadges(denom string) bool {
	return strings.HasPrefix(denom, "badges:")
}

func GetPartsFromDenom(denom string) ([]string, error) {
	if !CheckStartsWithBadges(denom) {
		return nil, fmt.Errorf("invalid denom: %s", denom)
	}

	parts := strings.Split(denom, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid denom: %s", denom)
	}
	return parts, nil
}

func ParseDenomCollectionId(denom string) (uint64, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return 0, err
	}

	// this is equivalent to split(':')[1]
	return strconv.ParseUint(parts[1], 10, 64)
}

func ParseDenomPath(denom string) (string, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return "", err
	}
	// this is equivalent to split(':')[1]
	return parts[2], nil
}

func GetCorrespondingPath(collection *badgestypes.BadgeCollection, denom string) (*badgestypes.CosmosCoinWrapperPath, error) {
	cosmosPaths := collection.CosmosCoinWrapperPaths
	for _, path := range cosmosPaths {
		if path.Denom == denom {
			return path, nil
		}
	}

	return nil, fmt.Errorf("path not found for denom: %s", denom)
}

func GetBalancesToTransfer(collection *badgestypes.BadgeCollection, denom string, amount sdkmath.Uint) ([]*badgestypes.Balance, error) {
	path, err := GetCorrespondingPath(collection, denom)
	if err != nil {
		return nil, err
	}

	balancesToTransfer := badgestypes.DeepCopyBalances(path.Balances)
	for _, balance := range balancesToTransfer {
		balance.Amount = balance.Amount.Mul(amount)
	}

	return balancesToTransfer, nil
}

func (k Keeper) ParseCollectionFromDenom(ctx sdk.Context, denom string) (*badgestypes.BadgeCollection, error) {
	collectionId, err := ParseDenomCollectionId(denom)
	if err != nil {
		return nil, err
	}

	collection, found := k.badgesKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidCollectionID, "collection %s not found", collectionId)
	}

	return collection, nil
}

func (k Keeper) WrapBadgesToSDKDenom(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// TODO: Make this a MsgTransferBadges or handle pre / post checks manually

	err = k.badgesKeeper.HandleTransfers(ctx, collection, []*badgestypes.Transfer{
		{
			From:        poolAddress,
			ToAddresses: []string{recipientAddress},
			Balances:    balancesToTransfer,
		},
	}, poolAddress)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UnwrapSDKDenomToBadges(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// TODO: Make this a MsgTransferBadges or handle pre / post checks manually
	err = k.badgesKeeper.HandleTransfers(ctx, collection, []*badgestypes.Transfer{
		{
			From:        recipientAddress,
			ToAddresses: []string{poolAddress},
			Balances:    balancesToTransfer,
		},
	}, recipientAddress)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) SendCoinsWithWrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if CheckStartsWithBadges(coin.Denom) {
			err := k.WrapBadgesToSDKDenom(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		}

		// otherwise, send the coins normally
		err := k.bankKeeper.SendCoins(ctx, from, to, coins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) SendCoinsWitUhwrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if CheckStartsWithBadges(coin.Denom) {
			err := k.UnwrapSDKDenomToBadges(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		}

		// otherwise, send the coins normally
		err := k.bankKeeper.SendCoins(ctx, from, to, coins)
		if err != nil {
			return err
		}
	}

	return nil
}
