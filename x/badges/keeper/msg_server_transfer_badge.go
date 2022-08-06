package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addresses_to_check := []uint64{}
	addresses_to_check = append(addresses_to_check, msg.ToAddresses...)
	addresses_to_check = append(addresses_to_check, msg.From)
	CreatorAccountNum, Badge, Permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, addresses_to_check, msg.BadgeId, msg.NumberRange.End, false,
	)
	
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
	if err != nil {
		return nil, err
	}


	FromBalanceKey := GetBalanceKey(msg.From, msg.BadgeId)
	fromBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, FromBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	for _, to := range msg.ToAddresses {
		ToBalanceKey := GetBalanceKey(to, msg.BadgeId)
		toBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, ToBalanceKey)
		if !found {
			toBadgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
		}
		
		for _, amount := range msg.Amounts {
			handledForcefulTransfer := false
			for currSubbadgeId := msg.NumberRange.Start; currSubbadgeId <= msg.NumberRange.End; currSubbadgeId++ {
				// Checks and handles if this account can transfer or is approved to transfer
				fromBadgeBalanceInfo, err = k.HandlePreTransfer(ctx, fromBadgeBalanceInfo, Badge, msg.BadgeId, currSubbadgeId, msg.From, to, CreatorAccountNum, amount)
				if err != nil {
					return nil, err
				}

				//We will always remove from "From" balance for both forceful (transfer it) and pending (put it in escrow)
				fromBadgeBalanceInfo, err = k.RemoveFromBadgeBalance(ctx, fromBadgeBalanceInfo, currSubbadgeId, amount)
				if err != nil {
					return nil, err
				}

				// Handle the transfer forcefully (no pending) if forceful transfers is set or "burning" (sending to manager address)
				// Else, handle it by adding a pending transfer

				// TODO: support forceful transfers when sending to reserved address numbers such as ETH NULL address
				var reservedAddress = []uint64{ }
				sendingToReservedAddress := false
				for _, address := range reservedAddress {
					if address == to {
						sendingToReservedAddress = true
						break
					}
				}

				if sendingToReservedAddress || Permissions.ForcefulTransfers() || Badge.Manager == to {
					handledForcefulTransfer = true
					toBadgeBalanceInfo, err = k.AddToBadgeBalance(ctx, toBadgeBalanceInfo, currSubbadgeId, amount)
					if err != nil {
						return nil, err
					}
				}
			}
			
			if !handledForcefulTransfer {
				fromBadgeBalanceInfo, toBadgeBalanceInfo, err = k.AddToBothPendingBadgeBalances(ctx, fromBadgeBalanceInfo, toBadgeBalanceInfo, *msg.NumberRange, to, msg.From, amount, CreatorAccountNum, true)
				if err != nil {
					return nil, err
				}
			}
		}
		if err := k.SetBadgeBalanceInStore(ctx, ToBalanceKey, toBadgeBalanceInfo); err != nil {
			return nil, err
		}		
	}	


	

	

	

	

	if err := k.SetBadgeBalanceInStore(ctx, FromBalanceKey, fromBadgeBalanceInfo); err != nil {
		return nil, err
	}

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "TransferBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeId Start", fmt.Sprint(msg.NumberRange.Start)),
			sdk.NewAttribute("SubbadgeId End", fmt.Sprint(msg.NumberRange.End)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
			sdk.NewAttribute("From", fmt.Sprint(msg.From)),
			sdk.NewAttribute("ToAddresses", fmt.Sprint(msg.ToAddresses)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
