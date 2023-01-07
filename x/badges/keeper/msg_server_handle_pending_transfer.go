package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
	Initial Pending Transfer From A -> Balances / approvals put into escrow
		B Accepts -> Simple add escrowed balance to B, remove from both pending
		B Rejects -> Remove from B Pending
			A has to trigger the escrow revert and remove from their balance
		A cancels -> Revert escrows and remove from both pending

	Initial Transfer Request From A
		B Accepts It (Forceful) -> Forceful transfer B to A -> Remove from both pending
		B Accepts It (Approve) -> Set pending transfer as approved in B's pending list
			A can then trigger a forceful transfer B to A -> Remove from both pending
		B Rejects -> Remove from B Pending
			A has to remove from their pending
		A cancels -> Remove from both pending
*/
func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
	})
	if err != nil {
		return nil, err
	}

	//Get creator balance information
	creatorBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	creatorBalanceInfo, found := k.GetUserBalanceFromStore(ctx, creatorBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	//We will store user balances in a cache and a didUpdate accountNum -> bool map to keep track of which ones are updated to avoid unnecessary writes if we didn't update anything.
	balanceInfoCache := make(map[uint64]types.UserBalanceInfo, 0)
	didUpdateBalanceInfo := make(map[uint64]bool, 0)

	for idx, CurrPendingTransfer := range creatorBalanceInfo.Pending {
		for _, nonceRange := range msg.NonceRanges {
			action := msg.Actions[0]
			if len(msg.Actions) > 1 {
				action = msg.Actions[idx]
			}
			accept := action == 1 || action == 2
			forcefulAccept := action == 2
			
			

			if nonceRange.End == 0 {
				nonceRange.End = nonceRange.Start
			}

			if !(CurrPendingTransfer.ThisPendingNonce <= nonceRange.End && CurrPendingTransfer.ThisPendingNonce >= nonceRange.Start) {
				continue
			}

			//Handle trying to accept after expiration time
			expired := uint64(ctx.BlockTime().Unix()) > CurrPendingTransfer.ExpirationTime
			if accept && CurrPendingTransfer.ExpirationTime != 0 {
				if expired {
					return nil, ErrPendingTransferExpired
				}
			}

			//An outgoingTransfer is when the balances of the badge are being transferred from the account calling this function.
			//This doesn't depend at all on whether it was sent by this account (could be a request) or if it is being accepted or rejected.
			outgoingTransfer := CurrPendingTransfer.From == CreatorAccountNum

			//Normal transfer flags (top four) and transfer request flags (bottom five) - Only one of these nine flags will be true
			acceptOwnOutgoingTransfer := CurrPendingTransfer.Sent && accept && outgoingTransfer
			acceptIncomingTransfer := !CurrPendingTransfer.Sent && accept && !outgoingTransfer
			cancelOwnOutgoingTransfer := CurrPendingTransfer.Sent && !accept && outgoingTransfer
			rejectIncomingTransfer := !CurrPendingTransfer.Sent && !accept

			finalizeOwnTransferRequestAfterApprovedByOtherParty := CurrPendingTransfer.Sent && accept && !outgoingTransfer
			//These two are the same scenario but split into forceful and non-forceful transfers, so manager doesn't have to pay gas for every transfer request
			acceptTransferRequestButMarkAsApproved := !forcefulAccept && !CurrPendingTransfer.Sent && accept && outgoingTransfer
			acceptTransferRequestForcefully := forcefulAccept && !CurrPendingTransfer.Sent && accept && outgoingTransfer
			rejectTransferRequest := !CurrPendingTransfer.Sent && !accept
			cancelOwnTransferRequest := CurrPendingTransfer.Sent && !accept && !outgoingTransfer

			if acceptOwnOutgoingTransfer {
				return nil, ErrCantAcceptOwnTransferRequest
			}

			//Check if before cancellation deadline and throw error if cancellation attempt
			pastCancelDeadline := true //if set to 0, it is cancellable at any time
			if CurrPendingTransfer.CantCancelBeforeTime != 0 {
				pastCancelDeadline = uint64(ctx.BlockTime().Unix()) > CurrPendingTransfer.CantCancelBeforeTime
			}
			if !pastCancelDeadline && (cancelOwnOutgoingTransfer || cancelOwnTransferRequest) {
				return nil, ErrCantCancelYet
			}

			//Other helper flags
			onlyUpdatingCreatorBalance := rejectIncomingTransfer || rejectTransferRequest || acceptTransferRequestButMarkAsApproved
			otherPartyMayHaveRemovedAlready := cancelOwnOutgoingTransfer || cancelOwnTransferRequest
			needToRemoveFromThisPending := !acceptTransferRequestButMarkAsApproved

			//Get basic information about the other party and update cache if not already in cache
			otherPartyAccountNum := CurrPendingTransfer.From
			if outgoingTransfer {
				otherPartyAccountNum = CurrPendingTransfer.To
			}
			otherPartyBalanceKey := ConstructBalanceKey(otherPartyAccountNum, msg.BadgeId)
			otherPartyNonce := CurrPendingTransfer.OtherPendingNonce
			otherPartyBalanceInfo, ok := balanceInfoCache[otherPartyAccountNum]
			//If we need the other party's balances and not in cache, get it from store
			if !ok && !onlyUpdatingCreatorBalance {
				otherPartyBalanceInfo, found = k.GetUserBalanceFromStore(ctx, otherPartyBalanceKey)
				balanceInfoCache[otherPartyAccountNum] = otherPartyBalanceInfo
				if !found {
					return nil, ErrUserBalanceNotExists
				}
			}

			//If we need to do something with the balances or pending transfers, update that
			if acceptIncomingTransfer {
				//Simple add to "To" balance
				creatorBalanceInfo, err = AddBalancesForIdRanges(creatorBalanceInfo, []*types.IdRange{CurrPendingTransfer.SubbadgeRange}, CurrPendingTransfer.Amount)
			} else if acceptTransferRequestButMarkAsApproved {
				//Just mark as accepted, don't remove from pending or do anything yet. Outsource gas to requester
				creatorBalanceInfo.Pending[idx].MarkedAsAccepted = true
			} else if cancelOwnOutgoingTransfer {
				//Previously escrowed balances and approvals. Need to revert
				creatorBalanceInfo, err = RevertEscrowedBalancesAndApprovals(creatorBalanceInfo, CurrPendingTransfer.SubbadgeRange, CurrPendingTransfer.From, CurrPendingTransfer.ApprovedBy, CurrPendingTransfer.Amount)
			} else if finalizeOwnTransferRequestAfterApprovedByOtherParty {
				//If other party marked as accepted, you can go ahead and finalize the transfer forcefully.
				idx, found := SearchPendingByNonce(otherPartyBalanceInfo.Pending, otherPartyNonce)
				if found {
					if !otherPartyBalanceInfo.Pending[idx].MarkedAsAccepted {
						return nil, ErrNotApproved
					}
				} else {
					return nil, ErrNotApproved
				}

				otherPartyBalanceInfo, creatorBalanceInfo, err = ForcefulTransfer(badge, CurrPendingTransfer.SubbadgeRange, otherPartyBalanceInfo, creatorBalanceInfo, CurrPendingTransfer.Amount, CurrPendingTransfer.From, CurrPendingTransfer.To, CurrPendingTransfer.From, CurrPendingTransfer.ExpirationTime)
			} else if acceptTransferRequestForcefully {
				//Accept a transfer request forcefully
				creatorBalanceInfo, otherPartyBalanceInfo, err = ForcefulTransfer(badge, CurrPendingTransfer.SubbadgeRange, creatorBalanceInfo, otherPartyBalanceInfo, CurrPendingTransfer.Amount, CurrPendingTransfer.From, CurrPendingTransfer.To, CreatorAccountNum, CurrPendingTransfer.ExpirationTime)
			}
			if err != nil {
				return nil, err
			}

			//Remove from this party's pending and other party's pending if applicable
			if needToRemoveFromThisPending {
				creatorBalanceInfo, err = RemovePending(creatorBalanceInfo, CurrPendingTransfer.ThisPendingNonce, otherPartyNonce)
				if err != nil {
					return nil, err
				}

				//Try to remove from the other party's pending, but in some situations it may have already been removed
				if !onlyUpdatingCreatorBalance {
					otherPartyBalanceInfo, err = RemovePending(otherPartyBalanceInfo, otherPartyNonce, CurrPendingTransfer.ThisPendingNonce)
					balanceInfoCache[otherPartyAccountNum] = otherPartyBalanceInfo
					if err != nil {
						if !(err == ErrPendingNotFound && otherPartyMayHaveRemovedAlready) {
							return nil, err
						}
					} else {
						didUpdateBalanceInfo[otherPartyAccountNum] = true
					}
				}
			}
		}
	}

	err = k.SetUserBalanceInStore(ctx, creatorBalanceKey, GetBalanceInfoToInsertToStorage(creatorBalanceInfo))
	if err != nil {
		return nil, err
	}

	//For all user balances that we did update, update the store
	for key := range didUpdateBalanceInfo {
		err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(key, msg.BadgeId), GetBalanceInfoToInsertToStorage(balanceInfoCache[key]))
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "HandledPendingTransfers"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("Actions", fmt.Sprint(msg.Actions)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("NonceRanges", fmt.Sprint(msg.NonceRanges)),
		),
	)

	return &types.MsgHandlePendingTransferResponse{}, nil

}
