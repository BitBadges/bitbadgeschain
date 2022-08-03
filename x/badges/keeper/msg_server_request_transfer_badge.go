package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	CreatorAccountNum, _, _, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{msg.From}, msg.BadgeId, msg.SubbadgeId, false,
	)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RequestTransfer"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("From", fmt.Sprint(msg.From)),
			sdk.NewAttribute("Amount", fmt.Sprint(msg.Amount)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeId", fmt.Sprint(msg.SubbadgeId)),
		),
	)
	return &types.MsgRequestTransferBadgeResponse{}, nil
}
