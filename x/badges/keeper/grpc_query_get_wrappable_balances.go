package keeper

import (
	"context"
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
	if !strings.HasPrefix(req.Denom, WrappedDenomPrefix) && !strings.HasPrefix(req.Denom, AliasDenomPrefix) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "denom must start with '%s' or '%s'", WrappedDenomPrefix, AliasDenomPrefix)
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
	path, err := GetCorrespondingAliasPath(collection, req.Denom)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "wrapper path not found for denom: %s", req.Denom)
	}

	// Get user's native balances (non-wrapped)
	userBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, req.Address)
	maxWrappableAmount, err := k.calculateMaxWrappableAmount(ctx, userBalances.Balances, path)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "error calculating max wrappable amount")
	}

	return &types.QueryGetWrappableBalancesResponse{
		MaxWrappableAmount: maxWrappableAmount,
	}, nil
}

// calculateMaxWrappableAmount calculates the maximum amount that can be wrapped.
//
// The conversion rate is: 1 x { amount: path.Amount, denom } = 1 x path.Balances
//
// Algorithm:
// 1. For each balance in path.Balances, find the corresponding user balance (matching token IDs and ownership times)
// 2. Calculate how many times that path balance can fit: userBalance.Amount / pathBalance.Amount
// 3. Take the minimum across all path balances (since we need all of them to perform a conversion)
// 4. Multiply by path.Amount to get the total wrappable amount
func (k Keeper) calculateMaxWrappableAmount(ctx sdk.Context, userBalances []*types.Balance, path *types.AliasPath) (sdkmath.Uint, error) {
	if len(path.Balances) == 0 {
		return sdkmath.NewUint(0), nil
	}

	if path.Amount.IsZero() || path.Amount.IsNil() {
		return sdkmath.NewUint(0), sdkerrors.Wrapf(types.ErrInvalidRequest, "path amount is zero")
	}

	// Track the minimum number of conversions possible across all path balances
	// We need all path balances to perform a conversion, so we're limited by the scarcest one
	var minConversions *sdkmath.Uint

	// For each balance required by the path, find how many conversions are possible
	for _, pathBalance := range path.Balances {
		// Calculate how many times this path balance can fit into the user's total balance
		// If pathBalance.Amount is zero, we can't perform any conversions
		if pathBalance.Amount.IsZero() {
			return sdkmath.NewUint(0), nil
		}

		// Get user balances that match this path balance's token IDs and ownership times
		userBalancesForPath, err := types.GetBalancesForIds(ctx, pathBalance.TokenIds, pathBalance.OwnershipTimes, userBalances)
		if err != nil {
			return sdkmath.NewUint(0), err
		}

		// If multiple balances are returned, they represent different ID/time combinations.
		// We need to take the minimum amount, as each balance represents a different combination
		// and we can only use the amount available for the specific ID/time combination we need.
		if len(userBalancesForPath) == 0 {
			return sdkmath.NewUint(0), nil
		}

		var minUserAmount sdkmath.Uint = userBalancesForPath[0].Amount
		for i := 1; i < len(userBalancesForPath); i++ {
			if userBalancesForPath[i].Amount.LT(minUserAmount) {
				minUserAmount = userBalancesForPath[i].Amount
			}
		}

		conversionsForThisBalance := minUserAmount.Quo(pathBalance.Amount)

		// Update minimum conversions (first iteration or if this is smaller)
		if minConversions == nil || conversionsForThisBalance.LT(*minConversions) {
			minConversions = &conversionsForThisBalance
		}
	}

	// If we couldn't perform any conversions, return 0
	if minConversions == nil {
		return sdkmath.NewUint(0), nil
	}

	// Multiply by path.Amount to get the total wrappable amount
	// This represents how many wrapped units (denom with amount path.Amount) can be created
	return minConversions.Mul(path.Amount), nil
}
