package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RevokeBadge(goCtx context.Context, msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                      msg.Creator,
		BadgeId:                      msg.BadgeId,
		SubbadgeRangesToValidate:     msg.SubbadgeRanges,
		AccountsThatCantEqualCreator: msg.Addresses,
		MustBeManager:                true,
		CanRevoke:                    true,
	})
	if err != nil {
		return nil, err
	}

	managerBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	managerBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, managerBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	for i, revokeAddress := range msg.Addresses {
		// Note that we check for duplicates in ValidateBasic, so these addresses will be unique every time
		addressBalanceKey := ConstructBalanceKey(revokeAddress, msg.BadgeId)
		addressBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, addressBalanceKey)
		if !found {
			return nil, ErrUserBalanceNotExists
		}

		revokeAmount := msg.Amounts[i]
		addressBalanceInfo, err = SubtractBalancesForIdRanges(ctx, addressBalanceInfo, msg.SubbadgeRanges, revokeAmount)
		if err != nil {
			return nil, err
		}

		managerBalanceInfo, err = AddBalancesForIdRanges(ctx, managerBalanceInfo, msg.SubbadgeRanges, revokeAmount)
		if err != nil {
			return nil, err
		}

		err = k.SetUserBalanceInStore(ctx, addressBalanceKey, addressBalanceInfo)
		if err != nil {
			return nil, err
		}
	}

	err = k.SetUserBalanceInStore(ctx, managerBalanceKey, managerBalanceInfo)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RevokeBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeRanges", fmt.Sprint(msg.SubbadgeRanges)),
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
		),
	)

	return &types.MsgRevokeBadgeResponse{}, nil
}
