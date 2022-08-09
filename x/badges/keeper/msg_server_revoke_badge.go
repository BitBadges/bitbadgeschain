package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RevokeBadge(goCtx context.Context, msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		SubbadgeRangesToValidate: msg.SubbadgeRanges,
		MustBeManager: true,
		AccountsToCheckIfRegistered: msg.Addresses,
		CanRevoke: true,
	}

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}

	ManagerBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	managerBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, ManagerBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	for i, revokeAddress := range msg.Addresses {
		if revokeAddress == CreatorAccountNum {
			return nil, ErrSenderAndReceiverSame
		}

		// Note that we check for duplicates in ValidateBasic, so these addresses will be unique every time
		AddressBalanceKey := ConstructBalanceKey(revokeAddress, msg.BadgeId)
		addressBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, AddressBalanceKey)
		if !found {
			return nil, ErrBadgeBalanceNotExists
		}

		revokeAmount := msg.Amounts[i]

		for _, subbadgeRange := range msg.SubbadgeRanges {
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				addressBalanceInfo, err = k.RemoveFromBadgeBalance(ctx, addressBalanceInfo, i, revokeAmount)
				if err != nil {
					return nil, err
				}

				managerBalanceInfo, err = k.AddToBadgeBalance(ctx, managerBalanceInfo, i, revokeAmount)
				if err != nil {
					return nil, err
				}

			}
		}
		err = k.SetBadgeBalanceInStore(ctx, AddressBalanceKey, addressBalanceInfo)
		if err != nil {
			return nil, err
		}
	}

	err = k.SetBadgeBalanceInStore(ctx, ManagerBalanceKey, managerBalanceInfo)
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
