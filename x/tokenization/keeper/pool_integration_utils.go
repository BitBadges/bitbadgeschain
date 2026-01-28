package keeper

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckStartsWithAliasDenom(denom string) bool {
	return strings.HasPrefix(denom, AliasDenomPrefix)
}

func CheckStartsWithWrappedOrAliasDenom(denom string) bool {
	return strings.HasPrefix(denom, WrappedDenomPrefix) || strings.HasPrefix(denom, AliasDenomPrefix)
}

// GetPartsFromDenom parses a badges denom into its parts
func GetPartsFromDenom(denom string) ([]string, error) {
	if !CheckStartsWithWrappedOrAliasDenom(denom) {
		return nil, errorsmod.Wrapf(ErrInvalidDenomFormat, "denom: %s", denom)
	}

	parts := strings.Split(denom, ":")
	if len(parts) < 3 {
		return nil, errorsmod.Wrapf(ErrInvalidDenomFormat, "denom: %s", denom)
	}
	return parts, nil
}

// ParseDenomCollectionId extracts the collection ID from a badges denom
func ParseDenomCollectionId(denom string) (uint64, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return 0, err
	}

	// this is equivalent to split(':')[1]
	return strconv.ParseUint(parts[1], 10, 64)
}

// ParseDenomPath extracts the path from a badges denom
func ParseDenomPath(denom string) (string, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return "", err
	}
	return strings.Join(parts[2:], ":"), nil
}

// GetCorrespondingAliasPath finds the AliasPath for a given denom
func GetCorrespondingAliasPath(collection *tokenizationtypes.TokenCollection, denom string) (*tokenizationtypes.AliasPath, error) {
	baseDenom, err := ParseDenomPath(denom)
	if err != nil {
		return nil, err
	}

	aliasPaths := collection.AliasPaths
	for _, path := range aliasPaths {
		if path.Denom == baseDenom {
			return path, nil
		}
	}

	return nil, errorsmod.Wrapf(ErrAliasPathNotFound, "denom: %s", denom)
}

// CheckIsAliasDenom checks if a denom is a wrapped badges denom
func (k Keeper) CheckIsAliasDenom(ctx sdk.Context, denom string) bool {
	if !CheckStartsWithWrappedOrAliasDenom(denom) {
		return false
	}

	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return false
	}

	_, err = GetCorrespondingAliasPath(collection, denom)
	return err == nil
}

// ParseCollectionFromDenom parses a collection from a badges denom
func (k Keeper) ParseCollectionFromDenom(ctx sdk.Context, denom string) (*tokenizationtypes.TokenCollection, error) {
	collectionId, err := ParseDenomCollectionId(denom)
	if err != nil {
		return nil, err
	}

	collection, found := k.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return nil, customhookstypes.WrapErr(&ctx, tokenizationtypes.ErrInvalidCollectionID, "collection %s not found",
			sdkmath.NewUint(collectionId).String())
	}

	return collection, nil
}
