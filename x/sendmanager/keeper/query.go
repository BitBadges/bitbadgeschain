package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Balance queries the balance of a specific denom for an address with alias routing.
// This allows querying both standard coins and alias denoms (e.g., badgeslp:).
func (q queryServer) Balance(ctx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "invalid request")
	}

	if req.Address == "" {
		return nil, sdkerrors.Wrap(errortypes.ErrInvalidAddress, "address cannot be empty")
	}

	if req.Denom == "" {
		return nil, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "denom cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrapf(errortypes.ErrInvalidAddress, "invalid address: %s", err)
	}

	// Use sendmanager's GetBalanceWithAliasRouting which handles both standard coins and alias denoms
	balance, err := q.k.GetBalanceWithAliasRouting(sdkCtx, address, req.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryBalanceResponse{
		Balance: balance,
	}, nil
}
