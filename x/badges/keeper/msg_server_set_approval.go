package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Sets approval to msg.Amount (no math involved)
func (k msgServer) SetApproval(goCtx context.Context, msg *types.MsgSetApproval) (*types.MsgSetApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rangesToValidate := []*types.IdRange{}
	for _, balance := range msg.Balances {
		rangesToValidate = append(rangesToValidate, balance.BadgeIds...)
	}

	_, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                      msg.Creator,
		CollectionId:                 msg.CollectionId,
		BadgeIdRangesToValidate:      rangesToValidate,
		AccountsThatCantEqualCreator: []string{msg.Address},
	})
	if err != nil {
		return nil, err
	}

	creatorBalanceKey := ConstructBalanceKey(msg.Creator, msg.CollectionId)
	creatorBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, creatorBalanceKey)
	if !found {
		creatorBalance = types.UserBalanceStore{}
	}

	for _, balance := range msg.Balances {
		creatorBalance.Approvals, err = SetApproval(creatorBalance.Approvals, balance.Amount, msg.Address, balance.BadgeIds)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetUserBalanceInStore(ctx, creatorBalanceKey, creatorBalance); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SetApproval"),
		),
	)

	return &types.MsgSetApprovalResponse{}, nil
}
