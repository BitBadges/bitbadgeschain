package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, Badge, Permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{msg.To, msg.From}, msg.BadgeId, msg.SubbadgeId, false,
	)
	if err != nil {
		return nil, err
	}

	FromBalanceKey := GetBalanceKey(msg.From, msg.BadgeId, msg.SubbadgeId)
	ToBalanceKey := GetBalanceKey(msg.To, msg.BadgeId, msg.SubbadgeId)

	// Checks and handles if this account can transfer or is approved to transfer
	err = k.HandlePreTransfer(ctx, Badge, msg.BadgeId, msg.SubbadgeId, msg.From, msg.To, CreatorAccountNum, msg.Amount)
	if err != nil {
		return nil, err
	}

	//We will always remove from "From" balance for both forceful (transfer it) and pending (put it in escrow)
	err = k.RemoveFromBadgeBalance(ctx, FromBalanceKey, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Handle the transfer forcefully (no pending) if forceful transfers is set or "burning" (sending to manager address)
	// Else, handle it by adding a pending transfer

	// TODO: support forceful transfers when sending to reserved address numbers such as ETH NULL address
	var reservedAddress = []uint64{ }
	sendingToReservedAddress := false
	for _, address := range reservedAddress {
		if address == msg.To {
			sendingToReservedAddress = true
			break
		}
	}

	
	if sendingToReservedAddress || Permissions.ForcefulTransfers() || Badge.Manager == msg.To {
		err := k.AddToBadgeBalance(ctx, ToBalanceKey, msg.Amount)
		if err != nil {
			return nil, err
		}
	} else {
		err = k.AddToBothPendingBadgeBalances(ctx, msg.BadgeId, msg.SubbadgeId, msg.To, msg.From, msg.Amount, CreatorAccountNum, true)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgTransferBadgeResponse{}, nil
}
