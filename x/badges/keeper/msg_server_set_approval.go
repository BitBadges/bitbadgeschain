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

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                      msg.Creator,
		BadgeId:                      msg.BadgeId,
		SubbadgeRangesToValidate:     msg.SubbadgeRanges,
		AccountsThatCantEqualCreator: []uint64{msg.Address},
		AccountsToCheckRegistration:  []uint64{msg.Address},
	})
	if err != nil {
		return nil, err
	}

	creatorBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	creatorBalanceInfo, found := k.Keeper.GetUserBalanceFromStore(ctx, creatorBalanceKey)
	if !found {
		creatorBalanceInfo = types.UserBalanceInfo{}
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		creatorBalanceInfo, err = SetApproval(creatorBalanceInfo, msg.Amount, msg.Address, subbadgeRange)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetUserBalanceInStore(ctx, creatorBalanceKey, GetBalanceInfoToInsertToStorage(creatorBalanceInfo)); err != nil {
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
