package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams implements the MsgServer interface
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// UpdateRateLimit implements the MsgServer interface
func (k msgServer) UpdateRateLimit(goCtx context.Context, req *types.MsgUpdateRateLimit) (*types.MsgUpdateRateLimitResponse, error) {
	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get current params
	params := k.GetParams(ctx)

	// Find if a rate limit with the same channel_id and denom exists
	foundIndex := -1
	for i, config := range params.RateLimits {
		if config.ChannelId == req.RateLimit.ChannelId && config.Denom == req.RateLimit.Denom {
			foundIndex = i
			break
		}
	}

	// Update or append the rate limit
	if foundIndex >= 0 {
		// Update existing rate limit
		params.RateLimits[foundIndex] = req.RateLimit
	} else {
		// Append new rate limit
		params.RateLimits = append(params.RateLimits, req.RateLimit)
	}

	// Validate and set the updated params
	if err := k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateRateLimitResponse{}, nil
}

