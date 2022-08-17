package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                      msg.Creator,
		BadgeId:                      msg.BadgeId,
		SubbadgeRangesToValidate:     msg.SubbadgeRanges,
		AccountsThatCantEqualCreator: []uint64{msg.From},
	})
	if err != nil {
		return nil, err
	}
	
	fromBalanceKey := ConstructBalanceKey(msg.From, msg.BadgeId)
	fromUserBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, fromBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	toBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	toUserBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, toBalanceKey)
	if !found {
		toUserBalanceInfo = types.UserBalanceInfo{}
	}

	//TODO: Maybe have []subbadgeRange within pending transfers instead of creating new transfers each time
	for _, subbadgeRange := range msg.SubbadgeRanges {
		fromUserBalanceInfo, toUserBalanceInfo, err = AppendPendingTransferForBothParties(ctx, fromUserBalanceInfo, toUserBalanceInfo, subbadgeRange, CreatorAccountNum, msg.From, msg.Amount, CreatorAccountNum, false, msg.ExpirationTime, msg.CantCancelBeforeTime)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, GetBalanceInfoToInsertToStorage(fromUserBalanceInfo)); err != nil {
		return nil, err
	}

	if err := k.SetUserBalanceInStore(ctx, toBalanceKey, GetBalanceInfoToInsertToStorage(toUserBalanceInfo)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RequestTransfer"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("From", fmt.Sprint(msg.From)),
			sdk.NewAttribute("Amount", fmt.Sprint(msg.Amount)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeRanges", fmt.Sprint(msg.SubbadgeRanges)),
		),
	)
	return &types.MsgRequestTransferBadgeResponse{}, nil
}
