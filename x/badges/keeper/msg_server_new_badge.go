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
		Permissions:           msg.Permissions,
		DefaultSubassetSupply: msg.DefaultSubassetSupply,
		FreezeRanges:          msg.FreezeAddressRanges,
		ArbitraryBytes:        msg.ArbitraryBytes,
		// SubassetSupplys: []*types.Subasset{},
		// NextSubassetId:       0,
	}

	if len(msg.SubassetSupplys) != 0 {
		managerBalanceInfo := types.UserBalanceInfo{}
		err := *new(error)
		badge, managerBalanceInfo, err = CreateSubassets(badge, managerBalanceInfo, msg.SubassetSupplys, msg.SubassetAmountsToCreate)
		if err != nil {
			return nil, err
		}

		if err := k.SetUserBalanceInStore(ctx, ConstructBalanceKey(CreatorAccountNum, badge.Id), GetBalanceInfoToInsertToStorage(managerBalanceInfo)); err != nil {
			return nil, err
		}
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
