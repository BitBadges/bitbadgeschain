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

	addressesToCheck := []uint64{}
	addressesToCheck = append(addressesToCheck, msg.ToAddresses...)
	addressesToCheck = append(addressesToCheck, msg.From)
	
	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		SubbadgeRangesToValidate: msg.SubbadgeRanges,
		AccountsToCheckIfRegistered: addressesToCheck,
	}

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}

	permissions := types.GetPermissions(badge.PermissionFlags)

	FromBalanceKey := ConstructBalanceKey(msg.From, msg.BadgeId)
	fromBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, FromBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	for _, to := range msg.ToAddresses {
		ToBalanceKey := ConstructBalanceKey(to, msg.BadgeId)
		toBadgeBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, ToBalanceKey)
		if !found {
			toBadgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
		}

		for _, amount := range msg.Amounts {
			handledForcefulTransfer := false
			for _, subbadgeRange := range msg.SubbadgeRanges {
				for currSubbadgeId := subbadgeRange.Start; currSubbadgeId <= subbadgeRange.End; currSubbadgeId++ {
					// Checks and handles if this account can transfer or is approved to transfer
					fromBadgeBalanceInfo, err = k.HandlePreTransfer(ctx, fromBadgeBalanceInfo, badge, msg.BadgeId, currSubbadgeId, msg.From, to, CreatorAccountNum, amount)
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
					var reservedAddress = []uint64{}
					sendingToReservedAddress := false
					for _, address := range reservedAddress {
						if address == to {
							sendingToReservedAddress = true
							break
						}
					}

					if sendingToReservedAddress || permissions.ForcefulTransfers() || badge.Manager == to {
						handledForcefulTransfer = true
						toBadgeBalanceInfo, err = k.AddToBadgeBalance(ctx, toBadgeBalanceInfo, currSubbadgeId, amount)
						if err != nil {
							return nil, err
						}
					}
				}

				if !handledForcefulTransfer {
					fromBadgeBalanceInfo, toBadgeBalanceInfo, err = k.AddToBothPendingBadgeBalances(ctx, fromBadgeBalanceInfo, toBadgeBalanceInfo, *subbadgeRange, to, msg.From, amount, CreatorAccountNum, true, msg.ExpirationTime)
					if err != nil {
						return nil, err
					}
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
			sdk.NewAttribute("SubbadgeRanges", fmt.Sprint(msg.SubbadgeRanges)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
			sdk.NewAttribute("From", fmt.Sprint(msg.From)),
			sdk.NewAttribute("ToAddresses", fmt.Sprint(msg.ToAddresses)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
