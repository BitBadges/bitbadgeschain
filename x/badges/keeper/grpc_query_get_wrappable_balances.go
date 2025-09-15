package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetWrappableBalances queries the maximum wrappable amount for a given denom and user address.
func (k Keeper) GetWrappableBalances(goCtx context.Context, req *types.QueryGetWrappableBalancesRequest) (*types.QueryGetWrappableBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Parse collection ID from denom (format: badges:COLL_ID:*)
	if !strings.HasPrefix(req.Denom, "badges:") && !strings.HasPrefix(req.Denom, "badgeslp:") {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "denom must start with 'badges:' or 'badgeslp:'")
	}

	parts := strings.Split(req.Denom, ":")
	if len(parts) < 3 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid denom format, expected '*:COLL_ID:*'")
	}

	collectionIdStr := parts[1]
	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid collection ID: %s", collectionIdStr)
	}

	// Fetch the collection
	collection, found := k.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return nil, sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s not found", collectionIdStr)
	}

	// Find the corresponding cosmos wrapper path
	wrapperPath, err := k.getCorrespondingWrapperPath(collection, req.Denom)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "wrapper path not found for denom: %s", req.Denom)
	}

	// Get user's native balances (non-wrapped)
	userBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, req.Address)
	maxWrappableAmount, err := k.calculateMaxWrappableAmount(ctx, userBalances.Balances, wrapperPath.Balances)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "error calculating max wrappable amount")
	}

	return &types.QueryGetWrappableBalancesResponse{
		MaxWrappableAmount: maxWrappableAmount,
	}, nil
}

// getCorrespondingWrapperPath finds the wrapper path that matches the given denom
func (k Keeper) getCorrespondingWrapperPath(collection *types.BadgeCollection, denom string) (*types.CosmosCoinWrapperPath, error) {
	parts := strings.Split(denom, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid denom format")
	}
	baseDenom := parts[2]

	// Extract numeric string from base denom for {id} placeholder replacement
	numericStr := ""
	for _, char := range baseDenom {
		if char >= '0' && char <= '9' {
			numericStr += string(char)
		}
	}

	cosmosPaths := collection.CosmosCoinWrapperPaths
	for _, path := range cosmosPaths {
		// Handle {id} placeholder replacement
		if path.AllowOverrideWithAnyValidToken && strings.Contains(path.Denom, "{id}") {
			if numericStr == "" {
				continue // Skip if no numeric part found
			}

			idFromDenom := sdkmath.NewUintFromString(numericStr)
			replacedDenom := strings.ReplaceAll(path.Denom, "{id}", idFromDenom.String())
			if replacedDenom == baseDenom {
				return path, nil
			}
		} else if path.Denom == baseDenom {
			return path, nil
		}
	}

	return nil, fmt.Errorf("path not found for denom: %s", denom)
}

// calculateMaxWrappableAmount calculates the maximum amount that can be wrapped
// by trying different amounts and checking if the user has enough balances to cover the conversion
func (k Keeper) calculateMaxWrappableAmount(ctx sdk.Context, userBalances []*types.Balance, wrapperPathBalances []*types.Balance) (sdkmath.Uint, error) {
	if len(wrapperPathBalances) == 0 {
		return sdkmath.NewUint(0), nil
	}

	// Get all badge IDs and ownership times from wrapper path balances
	var allBadgeIds []*types.UintRange
	var allOwnershipTimes []*types.UintRange

	for _, wrapperBalance := range wrapperPathBalances {
		allBadgeIds = append(allBadgeIds, wrapperBalance.BadgeIds...)
		allOwnershipTimes = append(allOwnershipTimes, wrapperBalance.OwnershipTimes...)
	}

	allBadgeIds, err := types.SortUintRangesAndMerge(allBadgeIds, true)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	allOwnershipTimes, err = types.SortUintRangesAndMerge(allOwnershipTimes, true)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	// Get user balances for the wrapper path badge IDs and ownership times
	balancesForIds, err := types.GetBalancesForIds(ctx, allBadgeIds, allOwnershipTimes, userBalances)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	// Filter out zero balances
	balancesForIds = types.FilterZeroBalances(balancesForIds)
	if len(balancesForIds) == 0 {
		return sdkmath.NewUint(0), nil
	}

	// Get all potential amounts from user balances and sort from largest to smallest
	var potentialAmounts []sdkmath.Uint
	for _, balance := range balancesForIds {
		potentialAmounts = append(potentialAmounts, balance.Amount)
	}

	// Sort from largest to smallest
	for i := 0; i < len(potentialAmounts); i++ {
		for j := i + 1; j < len(potentialAmounts); j++ {
			if potentialAmounts[i].LT(potentialAmounts[j]) {
				potentialAmounts[i], potentialAmounts[j] = potentialAmounts[j], potentialAmounts[i]
			}
		}
	}

	alreadyChecked := make(map[string]bool)

	// Find the largest amount that is less than or equal to the conversion amount
	for _, amount := range potentialAmounts {
		if alreadyChecked[amount.String()] {
			continue
		}
		alreadyChecked[amount.String()] = true

		// Create conversion balances by multiplying wrapper path balances by the amount
		var conversionBalances []*types.Balance
		for _, wrapperBalance := range wrapperPathBalances {
			conversionBalance := &types.Balance{
				Amount:         wrapperBalance.Amount.Mul(amount),
				BadgeIds:       wrapperBalance.BadgeIds,
				OwnershipTimes: wrapperBalance.OwnershipTimes,
			}
			conversionBalances = append(conversionBalances, conversionBalance)
		}

		// Clone user balances and subtract conversion balances
		userBalancesClone := types.DeepCopyBalances(userBalances)
		userBalancesClone, err = types.SubtractBalances(ctx, conversionBalances, userBalancesClone)
		if err != nil {
			continue // Skip this amount if subtraction fails
		}

		// If no error, no underflow and they have enough
		return amount, nil
	}

	return sdkmath.NewUint(0), nil
}
