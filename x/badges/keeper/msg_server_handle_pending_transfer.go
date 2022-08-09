package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
	}

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}

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
		1. Creator wants to cancel an outgoing transfer they sent -> Revert and remove from both
		2. Creator wants to accept an outgoing transfer they sent -> Error (can't accept a transfer you sent)

		3. Creator wants to cancel a request for transfer they sent -> Nothing additional and remove fom both
		4. Creator wants to accept a request for transfer they sent -> Can accept if other party approved it

		5. Creator wants to reject a request for transfer they received -> Nothing Additional
		6. Creator wants to accept a request for transfer they received -> Accept forcefully or mark as approved

		7. Creator wants to reject an incoming transfer they received -> Revert
		8. Creator wants to accept an incoming transfer they received -> Simple Accept

		Initial Pending Transfer From A -> Escrow
			B Accepts -> Simple Accept, remove from both pending
			B Rejects -> Remove from Your Pending -> Original has to reclaim, remove from their pending to revert
			A cancels -> Revert and remove from both pending

		Initial Transfer Request From A
			B Accepts It in Full Forcefully, remove from both pending
			B Marks As Approved -> A Can Forcefully Accept Only If Other Party Approved, remove from both pending
			B Rejects -> Remove from Their Pending -> A has to remove from their pending transfer
			A cancels -> Remove from both pending
	*/

	CreatorBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.BadgeId)
	creatorBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, CreatorBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	newCreatorBadgeBalanceInfo := creatorBadgeBalanceInfo

	balanceInfoCache := make(map[uint64]types.BadgeBalanceInfo, 0)
	balanceInfoUpdated := make(map[uint64]bool, 0)

	updated := false
	for _, nonceRange := range msg.NonceRanges {
		//TODO: validate nonce ranges
		//In the future, we can make this a binary search since this is all sorted by the nonces (append-only)
		for idx, CurrPendingTransfer := range creatorBadgeBalanceInfo.Pending {
			if CurrPendingTransfer.ThisPendingNonce <= nonceRange.End && CurrPendingTransfer.ThisPendingNonce >= nonceRange.Start {
				updated = true
				expired := false
				if CurrPendingTransfer.ExpirationTime != 0 {
					expired = uint64(ctx.BlockTime().Unix()) > CurrPendingTransfer.ExpirationTime
				}

				sentAndWantToAccept := CurrPendingTransfer.SendRequest && msg.Accept
				sentAndWantToCancel := CurrPendingTransfer.SendRequest && !msg.Accept
				receivedAndWantToAccept := !CurrPendingTransfer.SendRequest && msg.Accept
				receivedAndWantToReject := !CurrPendingTransfer.SendRequest && !msg.Accept
				outgoingTransfer := CurrPendingTransfer.From == CreatorAccountNum

				// Handle incoming / outgoing transfers
				simpleAddToYourBalance := receivedAndWantToAccept && !outgoingTransfer      //Remove from both
				needToRevertYourOutgoingTransfer := sentAndWantToCancel && outgoingTransfer //Remove from your pending and theirs if they haven't already

				// Handle rejections (throw expensive revert computation back at the proposing party)
				simpleRemoveFromYourPending := receivedAndWantToReject //Remove from only your pending

				// Handle transfer requests
				acceptAndMarkAsApproved := !msg.ForcefulAccept && receivedAndWantToAccept && outgoingTransfer //Don't remove at all
				forcefulAcceptAndTransfer := msg.ForcefulAccept && receivedAndWantToAccept && outgoingTransfer
				finalizeOnlyAfterApproved := sentAndWantToAccept && !outgoingTransfer       //Remove from both pending, full transfer, needs to be approved
				needToRevertYourTransferRequest := sentAndWantToCancel && !outgoingTransfer //Remove from both pendings if they haven't already

				// Can't accept your own outgoing transfer
				if sentAndWantToAccept && outgoingTransfer {
					return nil, ErrCantAcceptOwnTransferRequest
				}

				//Get basic info
				OtherPartyAccountNum := CurrPendingTransfer.From
				if outgoingTransfer {
					OtherPartyAccountNum = CurrPendingTransfer.To
				}
				OtherPartyBalanceKey := ConstructBalanceKey(OtherPartyAccountNum, msg.BadgeId)
				OtherPartyNonce := CurrPendingTransfer.OtherPendingNonce
				otherPartyBalanceInfo, ok := balanceInfoCache[OtherPartyAccountNum]

				// If its only a simple reject or accept, we don't need the other party's balances. Let the proposing party do the writes and reads
				if !ok && (!simpleRemoveFromYourPending || !acceptAndMarkAsApproved) {
					otherPartyBalanceInfo, found = k.GetBadgeBalanceFromStore(ctx, OtherPartyBalanceKey)
					balanceInfoCache[OtherPartyAccountNum] = otherPartyBalanceInfo

					if !found {
						if needToRevertYourOutgoingTransfer || needToRevertYourTransferRequest {

						} else {
							return nil, ErrBadgeBalanceNotExists
						}
					}
				}

				for i := CurrPendingTransfer.SubbadgeRange.Start; i <= CurrPendingTransfer.SubbadgeRange.End; i++ {
					if simpleAddToYourBalance {
						if expired {
							return nil, ErrPendingTransferExpired
						}
						newCreatorBadgeBalanceInfo, err = k.AddToBadgeBalance(ctx, newCreatorBadgeBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}
					} else if acceptAndMarkAsApproved {
						if expired {
							return nil, ErrPendingTransferExpired
						}
						newCreatorBadgeBalanceInfo.Pending[idx].MarkedAsApproved = true
					} else if needToRevertYourOutgoingTransfer {
						// Depending on if it is outgoing or not determines which party's balances to revert and add approvals back to
						FromInfo := newCreatorBadgeBalanceInfo
						if !outgoingTransfer {
							FromInfo = otherPartyBalanceInfo
						}

						FromInfo, err := k.AddToBadgeBalance(ctx, FromInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}

						//If it was sent via an approval, we need to add the approval back
						if CurrPendingTransfer.ApprovedBy != CurrPendingTransfer.From {
							FromInfo, err = k.AddBalanceToApproval(ctx, FromInfo, CurrPendingTransfer.Amount, CurrPendingTransfer.ApprovedBy, types.NumberRange{Start: i, End: i})
							if err != nil {
								return nil, err
							}
						}

						if outgoingTransfer {
							newCreatorBadgeBalanceInfo = FromInfo
						} else {
							otherPartyBalanceInfo = FromInfo
						}
					} else if finalizeOnlyAfterApproved {
						if expired {
							return nil, ErrPendingTransferExpired
						}

						approved := false
						for _, pending_info := range otherPartyBalanceInfo.Pending {
							if pending_info.ThisPendingNonce == OtherPartyNonce && pending_info.OtherPendingNonce == CurrPendingTransfer.ThisPendingNonce && pending_info.MarkedAsApproved {
								approved = true
							}
						}

						if !approved {
							return nil, ErrNotApproved
						}

						permissions := types.GetPermissions(badge.PermissionFlags)
						can_transfer := AccountNotFrozen(badge, permissions, CurrPendingTransfer.From)
						if !can_transfer {
							return nil, ErrAddressFrozen
						}

						otherPartyBalanceInfo, err = k.RemoveFromBadgeBalance(ctx, otherPartyBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}

						newCreatorBadgeBalanceInfo, err = k.AddToBadgeBalance(ctx, newCreatorBadgeBalanceInfo, i, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}
					} else if forcefulAcceptAndTransfer {
						if expired {
							return nil, ErrPendingTransferExpired
						}
						
						newCreatorBadgeBalanceInfo, err = k.HandlePreTransfer(ctx, newCreatorBadgeBalanceInfo, badge, msg.BadgeId, i, CreatorAccountNum, OtherPartyAccountNum, CreatorAccountNum, CurrPendingTransfer.Amount)
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
					}
				}

				//Remove in all situations from your balances if it is nto an approval mark
				if !acceptAndMarkAsApproved {
					//We already handled cases 2, 4, where we try and accept own request so all will end up with removing from both parties' pending requests whether accepting or rejecting
					newCreatorBadgeBalanceInfo, err = k.RemovePending(ctx, newCreatorBadgeBalanceInfo, CurrPendingTransfer.ThisPendingNonce, OtherPartyNonce)
					if err != nil {
						return nil, err
					}

					if !simpleRemoveFromYourPending {
						otherPartyBalanceInfo, err = k.RemovePending(ctx, otherPartyBalanceInfo, OtherPartyNonce, CurrPendingTransfer.ThisPendingNonce)

						if err != nil {
							if err == ErrPendingNotFound && (needToRevertYourOutgoingTransfer || needToRevertYourTransferRequest) {
								//This is okay. Other party already removed it
								balanceInfoUpdated[OtherPartyAccountNum] = false
							} else {
								return nil, err
							}
						} else {
							balanceInfoUpdated[OtherPartyAccountNum] = true
						}
					}
				}

				balanceInfoCache[OtherPartyAccountNum] = otherPartyBalanceInfo
			}
		}
	}

	if updated {
		err = k.SetBadgeBalanceInStore(ctx, CreatorBalanceKey, newCreatorBadgeBalanceInfo)
		if err != nil {
			return nil, err
		}

		for key := range balanceInfoUpdated {
			err = k.SetBadgeBalanceInStore(ctx, ConstructBalanceKey(key, msg.BadgeId), balanceInfoCache[key])
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
