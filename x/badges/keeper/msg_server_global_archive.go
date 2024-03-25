package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) GlobalArchive(goCtx context.Context, msg *types.MsgGlobalArchive) (*types.MsgGlobalArchiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	approvedCreators := []string{}
	approvedCreators = append(approvedCreators, "cosmos1kfr2xajdvs46h0ttqadu50nhu8x4v0tcfn4p0x")

	admin := false
	for _, creator := range approvedCreators {
		if creator == msg.Creator {
			admin = true
			break
		}
	}

	if !admin {
		return nil, types.ErrUnauthorized
	}

	err := k.SetGlobalArchiveInStore(ctx, msg.Archive)
	if err != nil {
		return nil, err
	}

	return &types.MsgGlobalArchiveResponse{}, nil
}
