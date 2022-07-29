package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	// Verify that the from and to addresses are registered; 
	account_nums := []uint64{}
	account_nums = append(account_nums, msg.To)
	must_approve := msg.From != creator_account_num
	if must_approve {
		account_nums = append(account_nums, msg.From)
	}
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
	permissions := types.GetPermissions(badge.PermissionFlags)
	from_balance_key := GetBalanceKey(msg.From, msg.BadgeId, msg.SubbadgeId)
	to_balance_key := GetBalanceKey(msg.To, msg.BadgeId, msg.SubbadgeId)


	// Check approvals if msg.Creator != msg.From
	if must_approve {
		err := k.RemoveBalanceFromApproval(ctx, from_balance_key, msg.Amount, creator_account_num) //if pending and cancelled, this approval will be added back
		if err != nil {
			return nil, err
		}
	}

	// TODO: Check if account is frozen


	// Handle the transfer forcefully (no pending) if forceful transfers is set or "burning" (sending to manager address)
	// Else, handle it by adding a pending transfer
	if permissions.ForcefulTransfers() || badge.Manager == msg.To {
		err := k.AddToBadgeBalance(ctx, to_balance_key, msg.Amount)
		if err != nil {
			return nil, err
		}

		err = k.RemoveFromBadgeBalance(ctx, from_balance_key, msg.Amount)
		if err != nil {
			return nil, err
		}
	} else {
		err = k.RemoveFromBadgeBalance(ctx, from_balance_key, msg.Amount) //We remove the balance while in pending (will be added  back if the transfer is rejected / cancelled)
		if err != nil {
			return nil, err
		}

		err = k.AddToBothPendingBadgeBalances(ctx, msg.BadgeId, msg.SubbadgeId,  msg.To, msg.From, msg.Amount, creator_account_num, true)

		if err != nil {
			return nil, err
		}
	}

	return &types.MsgTransferBadgeResponse{}, nil
}
