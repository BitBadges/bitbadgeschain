package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)
	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
	
	err :=  *new(error)
	/*
		Outgoing : Creator -> OtherParty
		Incoming : OtherParty -> CreatorAccountNum

		Sent vs Received: Who originally sent the pending transfer
		ApprovedBy: Who approved the sending of pending transfer

		"Requests for transfer" (5-8) currently always sent by that address (i.e. no ApprovedBy).

		Outcomes (all successful outcomes result in both pending transfers being removed from balanceInfo.Pending):
		Revert - Revert the balance and approval to the original values
		Accept Forcefully - Transfer the badge forcefully if permissions allow (i.e. no pending or else we will have an infinite loop of pending transfers)
		Simple Accept - Balance and approval already in escrow. Just need to simply add the new balance to the recipient.
		Nothing Additional - Besides removing from pending, nothing additional happens.

		These are the cases we have to handle:
		1. Creator wants to cancel an outgoing transfer they sent -> Revert
		2. Creator wants to accept an outgoing transfer they sent -> Error (can't accept a transfer you sent)

		3. Creator wants to cancel a request for transfer they sent -> Nothing Additional
		4. Creator wants to accept a request for transfer they sent -> Error (can't accept own request)

		5. Creator wants to reject a request for transfer they received -> Nothing Additional
		6. Creator wants to accept a request for transfer they received -> Accept Forcefully

		7. Creator wants to reject an incoming transfer they received -> Revert
		8. Creator wants to accept an incoming transfer they received -> Simple Accept
	*/

	CreatorBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId)
	creatorBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, CreatorBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	newCreatorBadgeBalanceInfo := creatorBadgeBalanceInfo
	
	balanceInfoCache := make(map[uint64]types.BadgeBalanceInfo, 0)

	updated := false
	
		//In the future, we can make this a binary search since this is all sorted by the nonces (append-only)
		for _, CurrPendingTransfer := range creatorBadgeBalanceInfo.Pending {
			if CurrPendingTransfer.ThisPendingNonce <= msg.NonceRanges.End && CurrPendingTransfer.ThisPendingNonce >= msg.NonceRanges.Start {
				updated = true
				if CurrPendingTransfer.SendRequest && msg.Accept {
					return nil, ErrCantAcceptOwnTransferRequest //Handle cases 2, 4
				}

				sentAndWantToCancel := CurrPendingTransfer.SendRequest && !msg.Accept      //Cases 1, 3
				receivedAndWantToAccept := !CurrPendingTransfer.SendRequest && msg.Accept  // Cases 6, 8
				receivedAndWantToReject := !CurrPendingTransfer.SendRequest && !msg.Accept //Cases 5, 7
				outgoingTransfer := CurrPendingTransfer.From == CreatorAccountNum
				balancesAreInEscrow := CurrPendingTransfer.SendRequest == outgoingTransfer

				// Cases 1, 7: Existing transfer was sent, is pending, but needs to be reversed
				needToRevertBalances := balancesAreInEscrow && ((sentAndWantToCancel && outgoingTransfer) || (receivedAndWantToReject && !outgoingTransfer))
				// Case 6: Accept a transfer / mint request from another party. Must go through all pre transfer checks. Forceful transfer (no pending)
				fullForcefulTransfer := receivedAndWantToAccept && outgoingTransfer
				// Case 8: Accepting an incoming transfer. Balances and approvals already in escrow.
				simpleAddToRecipientBalance := receivedAndWantToAccept && !outgoingTransfer
				// Cases 3 and 5: All we need to do is remove pending requests
				// Cases 2 and 4 already handled

				//Get basic info
				OtherPartyAccountNum := CurrPendingTransfer.From
				if outgoingTransfer {
					OtherPartyAccountNum = CurrPendingTransfer.To
				}
				OtherPartyBalanceKey := GetBalanceKey(OtherPartyAccountNum, msg.BadgeId)
				OtherPartyNonce := CurrPendingTransfer.OtherPendingNonce
				otherPartyBalanceInfo, ok := balanceInfoCache[OtherPartyAccountNum]
				if !ok {
					otherPartyBalanceInfo, found = k.GetBadgeBalanceFromStore(ctx, OtherPartyBalanceKey)
					balanceInfoCache[OtherPartyAccountNum] = otherPartyBalanceInfo
					if !found {
						return nil, ErrBadgeBalanceNotExists
					}
				}
				
				for i := CurrPendingTransfer.SubbadgeRange.Start; i <= CurrPendingTransfer.SubbadgeRange.End; i++ {
					if needToRevertBalances {
						// Depending on if it is outgoing or not determines which party's balances to revert and add approvals back to
						FromInfo := newCreatorBadgeBalanceInfo
						if !outgoingTransfer {
							FromInfo = otherPartyBalanceInfo
						}

						FromInfo, err := k.AddToBadgeBalance(ctx, FromInfo, i, CurrPendingTransfer.Amount); 
						if err != nil {
							return nil, err
						}

						//If it was sent via an approval, we need to add the approval back
						if CurrPendingTransfer.ApprovedBy != CurrPendingTransfer.From {
							FromInfo, err = k.AddBalanceToApproval(ctx, FromInfo, CurrPendingTransfer.Amount, CurrPendingTransfer.ApprovedBy, types.SubbadgeRange{Start: i, End: i}) 
							if err != nil {
								return nil, err
							}
						}

						if outgoingTransfer {
							newCreatorBadgeBalanceInfo = FromInfo
						} else {
							otherPartyBalanceInfo = FromInfo
						}
					} else if fullForcefulTransfer {
						
						newCreatorBadgeBalanceInfo, err = k.HandlePreTransfer(ctx, newCreatorBadgeBalanceInfo, badge, msg.BadgeId, i, CurrPendingTransfer.From, CurrPendingTransfer.To, CreatorAccountNum, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}

						newCreatorBadgeBalanceInfo, err = k.RemoveFromBadgeBalance(ctx, newCreatorBadgeBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}

						otherPartyBalanceInfo, err = k.AddToBadgeBalance(ctx, otherPartyBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}
					} else if simpleAddToRecipientBalance {
						newCreatorBadgeBalanceInfo, err = k.AddToBadgeBalance(ctx, newCreatorBadgeBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}
					}
				}

				//We already handled cases 2, 4, where we try and accept own request so all will end up with removing from both parties' pending requests whether accepting or rejecting
				newCreatorBadgeBalanceInfo, err = k.RemovePending(ctx, newCreatorBadgeBalanceInfo, CurrPendingTransfer.ThisPendingNonce, OtherPartyNonce)
				if err != nil {
					return nil, err
				}
				otherPartyBalanceInfo, err = k.RemovePending(ctx, otherPartyBalanceInfo, OtherPartyNonce, CurrPendingTransfer.ThisPendingNonce)
				if err != nil {
					return nil, err
				}

				balanceInfoCache[OtherPartyAccountNum] = otherPartyBalanceInfo
			}
		}
	

	

	if updated {
		err = k.SetBadgeBalanceInStore(ctx, CreatorBalanceKey, newCreatorBadgeBalanceInfo)
		if err != nil {
			return nil, err
		}

		for key, balanceInfo := range balanceInfoCache {
			err = k.SetBadgeBalanceInStore(ctx, GetBalanceKey(key, msg.BadgeId), balanceInfo)
			if err != nil {
				return nil, err
			}
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute(sdk.AttributeKeyAction, "HandledPendingTransfers"),
				sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
				sdk.NewAttribute("Accepted", fmt.Sprint(msg.Accept)),
				sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
				sdk.NewAttribute("NonceRanges", fmt.Sprint(msg.NonceRanges)),
						),
		)
		
		return &types.MsgHandlePendingTransferResponse{}, nil
	} else {
		return nil, ErrNoPendingTransferFound
	}
}
