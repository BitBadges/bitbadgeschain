package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, _, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{}, msg.BadgeId, msg.SubbadgeId, false,
	)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
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
		1. Creator wants to cancel an outgoing transfer they sent -> Revert
		2. Creator wants to accept an outgoing transfer they sent -> Error (can't accept a transfer you sent)

		3. Creator wants to cancel a request for transfer they sent -> Nothing Additional
		4. Creator wants to accept a request for transfer they sent -> Error (can't accept own request)

		5. Creator wants to reject a request for transfer they received -> Nothing Additional
		6. Creator wants to accept a request for transfer they received -> Accept Forcefully

		7. Creator wants to reject an incoming transfer they received -> Revert
		8. Creator wants to accept an incoming transfer they received -> Simple Accept
	*/

	CreatorBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, msg.SubbadgeId)
	BadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, CreatorBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	updated := false
	//In the future, we can make this a binary search since this is all sorted by the nonces (append-only)
	for _, CurrPendingTransfer := range BadgeBalanceInfo.Pending {
		if CurrPendingTransfer.ThisPendingNonce <= msg.EndingNonce && CurrPendingTransfer.ThisPendingNonce >= msg.StartingNonce {
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
			OtherPartyBalanceKey := GetBalanceKey(OtherPartyAccountNum, msg.BadgeId, msg.SubbadgeId)
			OtherPartyNonce := CurrPendingTransfer.OtherPendingNonce

			

			if needToRevertBalances {
				// Depending on if it is outgoing or not determines which party's balances to revert and add approvals back to
				FromKey := CreatorBalanceKey
				if !outgoingTransfer {
					FromKey = OtherPartyBalanceKey
				}

				if err := k.AddToBadgeBalance(ctx, FromKey, CurrPendingTransfer.Amount); err != nil {
					return nil, err
				}

				//If it was sent via an approval, we need to add the approval back
				if CurrPendingTransfer.ApprovedBy != CurrPendingTransfer.From {
					if err = k.AddBalanceToApproval(ctx, FromKey, CurrPendingTransfer.Amount, CurrPendingTransfer.ApprovedBy); err != nil {
						return nil, err
					}
				}
			} else if fullForcefulTransfer {
				err = k.HandlePreTransfer(ctx, badge, msg.BadgeId, msg.SubbadgeId, CurrPendingTransfer.From, CurrPendingTransfer.To, CreatorAccountNum, CurrPendingTransfer.Amount)
				if err != nil {
					return nil, err
				}

				if err := k.RemoveFromBadgeBalance(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount); err != nil {
					return nil, err
				}

				if err := k.AddToBadgeBalance(ctx, OtherPartyBalanceKey, CurrPendingTransfer.Amount); err != nil {
					return nil, err
				}
			} else if simpleAddToRecipientBalance {
				if err := k.AddToBadgeBalance(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount); err != nil {
					return nil, err
				}
			}

			//We already handled cases 2, 4, where we try and accept own request so all will end up with removing from both parties' pending requests whether accepting or rejecting
			if err := k.RemovePending(ctx, CreatorBalanceKey, CurrPendingTransfer.ThisPendingNonce, OtherPartyNonce); err != nil {
				return nil, err
			}
			if err = k.RemovePending(ctx, OtherPartyBalanceKey, OtherPartyNonce, CurrPendingTransfer.ThisPendingNonce); err != nil {
				return nil, err
			}
		}
	}
	if updated {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute(sdk.AttributeKeyAction, "HandledPendingTransfers"),
				sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
				sdk.NewAttribute("Accepted", fmt.Sprint(msg.Accept)),
				sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
				sdk.NewAttribute("SubbadgeId", fmt.Sprint(msg.SubbadgeId)),
				sdk.NewAttribute("StartingNonce", fmt.Sprint(msg.StartingNonce)),
				sdk.NewAttribute("EndingNonce", fmt.Sprint(msg.EndingNonce)),
			),
		)
		
		return &types.MsgHandlePendingTransferResponse{}, nil
	} else {
		return nil, ErrNoPendingTransferFound
	}
}
