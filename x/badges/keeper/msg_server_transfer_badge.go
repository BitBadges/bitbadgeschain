package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                  msg.Creator,
		BadgeId:                  msg.BadgeId,
		SubbadgeRangesToValidate: msg.SubbadgeRanges,
	})
	if err != nil {
		return nil, err
	}

	fromBalanceKey := ConstructBalanceKey(msg.From, msg.BadgeId)
	fromUserBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, fromBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	for _, to := range msg.ToAddresses {
		toBalanceKey := ConstructBalanceKey(to, msg.BadgeId)
		toUserBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, toBalanceKey)
		if !found {
			toUserBalanceInfo = types.UserBalanceInfo{}
		}

		//TODO: Batch
		for _, amount := range msg.Amounts {
			for _, subbadgeRange := range msg.SubbadgeRanges {
				fromUserBalanceInfo, toUserBalanceInfo, err = HandleTransfer(ctx, badge, *subbadgeRange, fromUserBalanceInfo, toUserBalanceInfo, amount, msg.From, to, CreatorAccountNum, msg.ExpirationTime, msg.CantCancelBeforeTime)
				if err != nil {
					return nil, err
				}
			}
		}

		if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalanceInfo); err != nil {
			return nil, err
		}
	}

	if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalanceInfo); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "TransferBadge"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeRanges", fmt.Sprint(msg.SubbadgeRanges)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
			sdk.NewAttribute("From", fmt.Sprint(msg.From)),
			sdk.NewAttribute("To", fmt.Sprint(msg.ToAddresses)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
