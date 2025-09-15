package keeper

import (
	"fmt"
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	types "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func CheckStartsWithBadges(denom string) bool {
	return strings.HasPrefix(denom, "badges:")
}

func (k Keeper) CheckIsBadgesWrappedDenom(ctx sdk.Context, denom string) bool {
	if !CheckStartsWithBadges(denom) {
		return false
	}

	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return false
	}

	path, err := GetCorrespondingPath(collection, denom)
	if err != nil {
		return false
	}

	// This is a little bit of an edge case
	// It is possible to have a badges: denom that is not the auto-converted denom
	// If this flag is true, we assume that they have to be wrapped first
	//
	// Ex: chaosnet denomination (badges:49:chaosnet)
	if path.AllowCosmosWrapping {
		return false
	}

	return true
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
	baseDenom, err := ParseDenomPath(denom)
	if err != nil {
		return nil, err
	}

	// This is okay because we don't allow numeric chars in denoms
	numericStr := ""
	for _, char := range baseDenom {
		if char >= '0' && char <= '9' {
			numericStr += string(char)
		}
	}

	cosmosPaths := collection.CosmosCoinWrapperPaths
	for _, path := range cosmosPaths {
		if path.AllowOverrideWithAnyValidToken {
			// 1. Replace the {id} placeholder with the actual denom
			// 2. Convert all balance.badgeIds to the actual badge ID
			if numericStr == "" {
				continue
			}

			idFromDenom := sdkmath.NewUintFromString(numericStr)
			path.Denom = strings.ReplaceAll(path.Denom, "{id}", idFromDenom.String())
			path.Balances = badgestypes.DeepCopyBalances(path.Balances)
			for _, balance := range path.Balances {
				balance.BadgeIds = []*badgestypes.UintRange{
					{Start: idFromDenom, End: idFromDenom},
				}
			}
		}

		if path.Denom == baseDenom {
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
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidCollectionID, "collection %s not found", sdkmath.NewUint(collectionId).String())
	}

	return collection, nil
}

func (k Keeper) SendNativeBadgesToPool(ctx sdk.Context, recipientAddress string, poolAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferBadges to ensure proper event handling and validation
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(k.badgesKeeper)

	msg := &badgestypes.MsgTransferBadges{
		Creator:      recipientAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        recipientAddress,
				ToAddresses: []string{poolAddress},
				Balances:    balancesToTransfer,
			},
		},
	}

	_, err = badgesMsgServer.TransferBadges(ctx, msg)
	return err
}

func (k Keeper) SendNativeBadgesFromPool(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferBadges to ensure proper event handling and validation
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(k.badgesKeeper)

	msg := &badgestypes.MsgTransferBadges{
		Creator:      poolAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        poolAddress,
				ToAddresses: []string{recipientAddress},
				Balances:    balancesToTransfer,
			},
		},
	}

	_, err = badgesMsgServer.TransferBadges(ctx, msg)
	return err
}

// IMPORTANT: Should ONLY be called when to address is a pool address
func (k Keeper) SendCoinsToPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if k.CheckIsBadgesWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeBadgesToPool(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

			//Mint corresponding coins for compatibillity
			err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				return err
			}

			// Send to the pool address
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		} else {
			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IMPORTANT: Should ONLY be called when from address is a pool address
func (k Keeper) SendCoinsFromPoolWithUnwrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, unwrap it
	for _, coin := range coins {
		if k.CheckIsBadgesWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeBadgesFromPool(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

			// Send coins to module from pool
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				return err
			}

			//Burn corresponding coins from the pool for compatibillity
			err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		} else {

			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Used for taker fees
func (k Keeper) FundCommunityPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		moduleAddress := authtypes.NewModuleAddress(distrtypes.ModuleName).String()

		if k.CheckIsBadgesWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeBadgesToPool(ctx, from.String(), moduleAddress, coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

		} else {
			err := k.communityPoolKeeper.FundCommunityPool(ctx, sdk.NewCoins(coin), from)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
