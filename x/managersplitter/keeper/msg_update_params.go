package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "invalid authority; expected "+k.GetAuthority()+", got "+msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	// Store params (if needed in the future)
	// For now, params are empty, so we just validate

	return &types.MsgUpdateParamsResponse{}, nil
}

