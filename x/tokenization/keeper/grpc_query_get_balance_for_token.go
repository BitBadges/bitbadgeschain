package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBalanceForToken queries the balance amount for a specific token ID at a specific time.
func (k Keeper) GetBalanceForToken(goCtx context.Context, req *types.QueryGetBalanceForTokenRequest) (*types.QueryGetBalanceForTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	collectionId := sdkmath.NewUintFromString(req.CollectionId)
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s not found", req.CollectionId)
	}

	// Fetch the full balance store for this address
	balanceStore, _, err := k.GetBalanceOrApplyDefault(ctx, collection, req.Address)
	if err != nil {
		return nil, err
	}

	// Parse tokenId
	tokenId := sdkmath.NewUintFromString(req.TokenId)

	// Parse time - default to current block time if empty or "0"
	var timeVal sdkmath.Uint
	if req.Time == "" || req.Time == "0" {
		timeVal = sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	} else {
		timeVal = sdkmath.NewUintFromString(req.Time)
	}

	// Use GetBalancesForIds to look up the balance for the specific tokenId and time
	tokenIdRange := &types.UintRange{Start: tokenId, End: tokenId}
	timeRange := &types.UintRange{Start: timeVal, End: timeVal}

	fetchedBalances, err := types.GetBalancesForIds(
		ctx,
		[]*types.UintRange{tokenIdRange},
		[]*types.UintRange{timeRange},
		balanceStore.Balances,
	)
	if err != nil {
		return nil, err
	}

	if len(fetchedBalances) > 1 {
		return nil, status.Error(codes.Internal, "unexpected: GetBalancesForIds returned more than one entry for a single tokenId and time")
	}

	amount := sdkmath.NewUint(0)
	if len(fetchedBalances) == 1 {
		amount = fetchedBalances[0].Amount
	}

	return &types.QueryGetBalanceForTokenResponse{
		Balance: amount.String(),
	}, nil
}
