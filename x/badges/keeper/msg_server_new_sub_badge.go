package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.Id)

	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanCreateSubbadges() {
		return nil, ErrInvalidPermissions
	}

	//Once here, we should be safe to mint
	//By default, we assume non fungible subbadge (i.e. supply == 1) so we don't store if supply == 1
	subasset_id := badge.NextSubassetId
	if msg.Supply != 1 {
		badge.SubassetsTotalSupply = append(badge.SubassetsTotalSupply, &types.Subasset{
			Id:     subasset_id,
			Supply: msg.Supply,
		})
	}
	badge.NextSubassetId += 1

	//Mint the total supply of subbadge to the manager
	ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.Id, subasset_id)
	if err := k.AddToBadgeBalance(ctx, ManagerBalanceKey, msg.Supply); err != nil {
		return nil, err
	}

	if err := k.UpdateBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	return &types.MsgNewSubBadgeResponse{
		SubassetId: subasset_id,
	}, nil
}
