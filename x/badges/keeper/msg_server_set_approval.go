package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Assumes a BadgeBalance exists for the given badge and subbadge and msg.Creator
func (k msgServer) SetApproval(goCtx context.Context, msg *types.MsgSetApproval) (*types.MsgSetApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	//Can't send to same address
	if msg.Address == creator_account_num {
		return nil, ErrSenderAndReceiverSame
	}

	// Verify that the from and to addresses are registered
	account_nums := []uint64{}
	account_nums = append(account_nums, msg.Address)
	err := k.AssertAccountNumbersAreValid(ctx, account_nums)
	if err != nil {
		return nil, err
	}

	// Verify that the badge and subbadge exist and are valid
	err = k.AssertBadgeAndSubBadgeExists(ctx, msg.BadgeId, msg.SubbadgeId)
	if err != nil {
		return nil, err
	}

	// Set the approval to msg.Amount (no math; just set to whatever the user requests)
	balance_key := GetBalanceKey(creator_account_num, msg.BadgeId, msg.SubbadgeId)
	err = k.Keeper.SetApproval(ctx, balance_key, msg.Amount, msg.Address)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetApprovalResponse{}, nil
}
