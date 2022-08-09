package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets approval to msg.Amount (no math involved)
func (k msgServer) SetApproval(goCtx context.Context, msg *types.MsgSetApproval) (*types.MsgSetApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		AccountsToCheckIfRegistered: []uint64{msg.Address},
		SubbadgeRangesToValidate: msg.SubbadgeRanges,
		AccountsThatCantEqualCreator: []uint64{msg.Address},
	}

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}


	BalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	badgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, BalanceKey)
	if !found {
		badgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		badgeBalanceInfo, err = k.Keeper.SetApproval(ctx, badgeBalanceInfo, msg.Amount, msg.Address, *subbadgeRange, msg.ExpirationTime)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetBadgeBalanceInStore(ctx, BalanceKey, badgeBalanceInfo); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SetApproval"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeRanges", fmt.Sprint(msg.SubbadgeRanges)),
			sdk.NewAttribute("ApprovedAddress", fmt.Sprint(msg.Address)),
			sdk.NewAttribute("Amount", fmt.Sprint(msg.Amount)),
		),
	)

	return &types.MsgSetApprovalResponse{}, nil
}
