package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)


func (k msgServer) RevokeBadge(goCtx context.Context,  msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	//Can't revoke from same address
	if creator_account_num == msg.Address {
		return nil, ErrSenderAndReceiverSame
	}

	// Verify that the from and to addresses are registered; 
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

	// Verify that the permissions are valid
	badge, _ := k.GetBadgeFromStore(ctx, msg.BadgeId) //currently ignore error because above we assert that it exists
	permissions := GetPermissions(badge.PermissionFlags)
	if !permissions.can_revoke {
		return nil, ErrInvalidPermissions
	}

	if badge.Manager != creator_account_num {
		return nil, ErrSenderIsNotManager
	}

	address_balance_key := GetBalanceKey(msg.Address, msg.BadgeId, msg.SubbadgeId)
	manager_balance_key := GetBalanceKey(creator_account_num, msg.BadgeId, msg.SubbadgeId)


	err = k.RemoveFromBadgeBalance(ctx, address_balance_key, msg.Amount)
	if err != nil {
		return nil, err
	}

	err = k.AddToBadgeBalance(ctx, manager_balance_key, msg.Amount)
	if err != nil {
		return nil, err
	}

    // TODO: Handling the message
    _ = ctx

	return &types.MsgRevokeBadgeResponse{}, nil
}
