package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) SelfDestructBadge(goCtx context.Context, msg *types.MsgSelfDestructBadge) (*types.MsgSelfDestructBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	nextSubassetId := badge.NextSubassetId

	for i := uint64(0); i < nextSubassetId; i++ {
		ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, i)
		SubassetSupply := uint64(1) //Default if not found
		for _, subasset := range badge.SubassetsTotalSupply {
			if subasset.StartId <= i && subasset.EndId >= i {
				SubassetSupply = subasset.Supply
				break
			}
		}

		balanceInfo, found := k.GetBadgeBalanceFromStore(ctx, ManagerBalanceKey)
		if !found || balanceInfo.Balance < SubassetSupply {
			return nil, ErrMustOwnTotalSupplyToSelfDestruct
		}
	}

	k.DeleteBadgeFromStore(ctx, msg.BadgeId)

	return &types.MsgSelfDestructBadgeResponse{}, nil
}
