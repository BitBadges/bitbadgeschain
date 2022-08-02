package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	CreatorAccountNum, _, _, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{ msg.From }, msg.BadgeId, msg.SubbadgeId, false,
	)
	if err != nil {
		return nil, err
	}

	if CreatorAccountNum == msg.From {
		return nil, ErrSenderAndReceiverSame // Can't request yourself for transfer
	}

	// Add to both account's pending transfers (we handle permissions when acecepting / rejecting the transfer)
	err = k.AddToBothPendingBadgeBalances(ctx, msg.BadgeId, msg.SubbadgeId, CreatorAccountNum, msg.From, msg.Amount, CreatorAccountNum, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestTransferBadgeResponse{}, nil
}
