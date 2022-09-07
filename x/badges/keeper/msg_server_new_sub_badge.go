package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:            msg.Creator,
		BadgeId:            msg.BadgeId,
		MustBeManager:      true,
		CanCreateSubbadges: true,
	})
	if err != nil {
		return nil, err
	}

	originalSubassetId := badge.NextSubassetId

	managerBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	managerBalanceInfo, found := k.GetUserBalanceFromStore(ctx, managerBalanceKey)
	if !found {
		managerBalanceInfo = types.UserBalanceInfo{}
	}

	badge, managerBalanceInfo, err = CreateSubassets(badge, managerBalanceInfo, msg.Supplys, msg.AmountsToCreate)
	if err != nil {
		return nil, err
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	if err := k.SetUserBalanceInStore(ctx, managerBalanceKey, GetBalanceInfoToInsertToStorage(managerBalanceInfo)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedSubBadges"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(badge.Id)),
			sdk.NewAttribute("FirstId", fmt.Sprint(originalSubassetId)),
			sdk.NewAttribute("LastId", fmt.Sprint(badge.NextSubassetId-1)),
		),
	)

	return &types.MsgNewSubBadgeResponse{
		NextSubassetId: badge.NextSubassetId,
	}, nil
}
