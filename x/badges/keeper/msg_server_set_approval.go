package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets approval to msg.Amount (no math involved)
func (k msgServer) SetApproval(goCtx context.Context, msg *types.MsgSetApproval) (*types.MsgSetApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, _, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{ msg.Address }, msg.BadgeId,	msg.SubbadgeId, false,
	)
	if err != nil {
		return nil, err
	}

	if CreatorAccountNum == msg.Address {
		return nil, ErrSenderAndReceiverSame // Can't approve yourself
	}

	BalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, msg.SubbadgeId)	

	if err := k.Keeper.SetApproval(ctx, BalanceKey, msg.Amount, msg.Address); err != nil {
		return nil, err
	}

	return &types.MsgSetApprovalResponse{}, nil
}
