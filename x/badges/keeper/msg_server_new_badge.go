package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewBadge(goCtx context.Context, msg *types.MsgNewBadge) (*types.MsgNewBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)
	
	// We shouldn't have to call UniversalValidate() because anyone can call this function

	NextBadgeId := k.GetNextBadgeId(ctx)
	k.IncrementNextBadgeId(ctx)

	badge := types.BitBadge{
		Id:                    NextBadgeId,
		Uri:                   msg.Uri,
		Manager:               CreatorAccountNum,
		Permissions:       msg.Permissions,
		SubassetUriFormat:     msg.SubassetUris,
		DefaultSubassetSupply: msg.DefaultSubassetSupply,
		// SubassetSupplys: []*types.Subasset{},
		// NextSubassetId:       0,
		// FreezeAddresses:      []uint64{},
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(NextBadgeId)),
		),
	)

	return &types.MsgNewBadgeResponse{
		Id: NextBadgeId,
	}, nil
}
