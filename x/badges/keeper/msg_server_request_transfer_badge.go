package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		SubbadgeRangesToValidate: msg.SubbadgeRanges,
		AccountsThatCantEqualCreator: []uint64{msg.From},
	}

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}

	FromBalanceKey := ConstructBalanceKey(msg.From, msg.BadgeId)
	ToBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)

	fromBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, FromBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	toBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, ToBalanceKey)
	if !found {
		toBadgeBalanceInfo = types.BadgeBalanceInfo{}
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		fromBadgeBalanceInfo, toBadgeBalanceInfo, err = k.AppendPendingTransferForBothParties(ctx, fromBadgeBalanceInfo, toBadgeBalanceInfo, *subbadgeRange, CreatorAccountNum, msg.From, msg.Amount, CreatorAccountNum, false, msg.ExpirationTime)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetBadgeBalanceInStore(ctx, FromBalanceKey, fromBadgeBalanceInfo); err != nil {
		return nil, err
	}

	if err := k.SetBadgeBalanceInStore(ctx, ToBalanceKey, toBadgeBalanceInfo); err != nil {
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
