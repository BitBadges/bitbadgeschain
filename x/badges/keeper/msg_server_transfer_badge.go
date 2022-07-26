package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	accsToCheck := []uint64{msg.From}
	accsToCheck = append(accsToCheck, msg.ToAddresses...)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                     msg.Creator,
		BadgeId:                     msg.BadgeId,
		SubbadgeRangesToValidate:    msg.SubbadgeRanges,
		AccountsToCheckRegistration: accsToCheck,
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

		for _, amount := range msg.Amounts {
			for _, subbadgeRange := range msg.SubbadgeRanges {
				fromUserBalanceInfo, toUserBalanceInfo, err = HandleTransfer(badge, subbadgeRange, fromUserBalanceInfo, toUserBalanceInfo, amount, msg.From, to, CreatorAccountNum, msg.ExpirationTime, msg.CantCancelBeforeTime)
				if err != nil {
					return nil, err
				}
			}
		}

		if err := k.SetUserBalanceInStore(ctx, toBalanceKey, GetBalanceInfoToInsertToStorage(toUserBalanceInfo)); err != nil {
			return nil, err
		}
	}

	if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, GetBalanceInfoToInsertToStorage(fromUserBalanceInfo)); err != nil {
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
