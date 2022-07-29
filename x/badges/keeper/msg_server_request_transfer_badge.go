package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Note there is no support for approvals or transfer requests on your behalf. The "To" address will be msg.Creator
func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	// Can't request transfer to yourself
	if creator_account_num == msg.From {
		return nil, ErrSenderAndReceiverSame
	}

	// Verify that the from address is registered
	account_nums := []uint64{}
	account_nums = append(account_nums, msg.From)
	err := k.AssertAccountNumbersAreValid(ctx, account_nums)
	if err != nil {
		return nil, err
	}

	// Verify that the badge and subbadge exist and are valid
	err = k.AssertBadgeAndSubBadgeExists(ctx, msg.BadgeId, msg.SubbadgeId)
	if err != nil {
		return nil, err
	}
	
	// Add to both account's pending transfers (we handle permissions when acecepting / rejecting the transfer)
	err = k.AddToBothPendingBadgeBalances(ctx, msg.BadgeId, msg.SubbadgeId, creator_account_num, msg.From, msg.Amount, creator_account_num, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestTransferBadgeResponse{}, nil
}
