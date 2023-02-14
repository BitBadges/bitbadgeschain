package keeper

import (
	"context"
	"encoding/json"
	"fmt"

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

	CreatorAccountNum, _, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                      msg.Creator,
		CollectionId:                 msg.CollectionId,
		BadgeIdRangesToValidate:      rangesToValidate,
		AccountsThatCantEqualCreator: []uint64{msg.Address},
		AccountsToCheckRegistration:  []uint64{msg.Address},
	})
	if err != nil {
		return nil, err
	}

	creatorBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.CollectionId)
	creatorbalance, found := k.Keeper.GetUserBalanceFromStore(ctx, creatorBalanceKey)
	if !found {
		creatorbalance = types.UserBalance{}
	}

	for _, balance := range msg.Balances {
		amount := balance.Balance
		for _, badgeIdRange := range balance.BadgeIds {
			creatorbalance, err = SetApproval(creatorbalance, amount, msg.Address, badgeIdRange)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := k.SetUserBalanceInStore(ctx, creatorBalanceKey, creatorbalance); err != nil {
		return nil, err
	}

	

	userBalanceJson, err := json.Marshal(creatorbalance)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SetApproval"),
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
			sdk.NewAttribute("creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("user_balance", string(userBalanceJson)),
		),
	)

	return &types.MsgSetApprovalResponse{}, nil
}
