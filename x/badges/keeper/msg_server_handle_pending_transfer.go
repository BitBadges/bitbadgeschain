package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, _, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{ }, msg.BadgeId, msg.SubbadgeId, false,
	)
	if err != nil {
		return nil, err
	}

	/* 
		These are the cases we have to handle:
		1. Creator wants to cancel an outgoing transfer they sent where they are the ApprovedBy 
		2. Creator wants to cancel an outgoing transfer they sent where they are not the ApprovedBy
		3. Creator wants to accept an outgoing transfer they sent where they are the ApprovedBy (can't accept own transfer)
		4. Creator wants to accept an outgoing transfer they sent where they are not the Approved By (can't accept own transfer)
		
		5. Creator wants to cancel an incoming transfer (transfer request) they sent 
		6. Creator wants to accept an incoming transfer (transfer request) they sent (can't accept own transfer)

		7. Creator wants to reject an outgoing transfer (transfer request) they received
		8. Creator wants to accept an outgoing transfer (transfer request) they received
        
		9. Creator wants to reject an incoming transfer they received where From == ApprovedBy
		10. Creator wants to reject an incoming transfer they received where From != ApprovedBy
		11. Creator wants to accept an incoming transfer they received where From == ApprovedBy
		12. Creator wants to accept an incoming transfer they received where From != ApprovedBy
	*/

	CreatorBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, msg.SubbadgeId)
	BadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, CreatorBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	} 

	//In the future, we can make this a binary search since this is all sorted by the nonces (append-only)
	for _, CurrPendingTransfer := range BadgeBalanceInfo.Pending {
		if CurrPendingTransfer.ThisPendingNonce == msg.ThisNonce {
			if CurrPendingTransfer.SendRequest && msg.Accept {
				return nil, ErrCantAcceptOwnTransferRequest //Handle cases 3, 4, 6
			}
			sentAndWantToCancel := CurrPendingTransfer.SendRequest && !msg.Accept //Cases 1, 2, 5
			receivedAndWantToAccept := !CurrPendingTransfer.SendRequest && msg.Accept // Cases 8, 11, 12
			receivedAndWantToReject := !CurrPendingTransfer.SendRequest && !msg.Accept //Cases 7, 9, 10
			outgoingTransfer := CurrPendingTransfer.From == CreatorAccountNum
			
			//Get basic info
			OtherPartyAccountNum := CurrPendingTransfer.From
			if outgoingTransfer {
				OtherPartyAccountNum = CurrPendingTransfer.To
			}
			OtherPartyBalanceKey := GetBalanceKey(OtherPartyAccountNum, msg.BadgeId, msg.SubbadgeId)
			OtherPartyNonce := CurrPendingTransfer.OtherPendingNonce

			//We already handled cases 3, 4, 6, so all will end up with removing from both parties' pending requests whether accepting or rejecting
			if err := k.RemovePending(ctx, CreatorBalanceKey, msg.ThisNonce, OtherPartyNonce); err != nil {
				return nil, err
			}
			if err = k.RemovePending(ctx, OtherPartyBalanceKey, OtherPartyNonce, msg.ThisNonce); err != nil {
				return nil, err
			}

			if sentAndWantToCancel {
				if outgoingTransfer {
					// Cases 1, 2: We need to revert balances since they are in escrow
					if err := k.AddToBadgeBalance(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount); err != nil {
						return nil, err
					}

					//If it was sent via an approval, we need to add the approval back
					if CurrPendingTransfer.ApprovedBy != CurrPendingTransfer.From {
						if err = k.AddBalanceToApproval(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount, CurrPendingTransfer.ApprovedBy); err != nil {
							return nil, err
						}
					}
				} else {
					// Case 5: Everything is already handled. We just remove from pending
				}
			} else {
				if receivedAndWantToAccept {
					if outgoingTransfer {
						//Case 8: Accepting a transfer request we received. Forcefully send the funds to the other party
						err = k.HandlePreTransfer(ctx, badge, msg.BadgeId, msg.SubbadgeId, CurrPendingTransfer.From, CurrPendingTransfer.To, CreatorAccountNum, CurrPendingTransfer.Amount)
						if err != nil {
							return nil, err
						}

						if err := k.RemoveFromBadgeBalance(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount); err != nil {
							return nil, err
						}

						// We already removed from "From" balance so all we have to do now is add to "To" balance
						if err := k.AddToBadgeBalance(ctx, OtherPartyBalanceKey, CurrPendingTransfer.Amount); err != nil {
							return nil, err
						}
					} else {
						// Case 11, 12: Accepting an incoming transfer; Balances and approvals already taken care of 
						// We already removed from "From" balance so all we have to do now is add to "To" balance
						if err := k.AddToBadgeBalance(ctx, CreatorBalanceKey, CurrPendingTransfer.Amount); err != nil {
							return nil, err
						}
					}
				} else if receivedAndWantToReject {
					if outgoingTransfer {
						// Case 7: Reject a transfer request we received. Nothing needs to be done besides remove from pending which we already did
					} else {
						// Case 9, 10: Reject an incoming transfer we received. We need to revert balances since they are in escrow
						if err := k.AddToBadgeBalance(ctx, OtherPartyBalanceKey, CurrPendingTransfer.Amount); err != nil {
							return nil, err
						}

						//If it was sent via an approval, we need to add the approval back
						if CurrPendingTransfer.ApprovedBy != CurrPendingTransfer.From {
							if err = k.AddBalanceToApproval(ctx, OtherPartyBalanceKey, CurrPendingTransfer.Amount, CurrPendingTransfer.ApprovedBy); err != nil {
								return nil, err
							}
						}
					}
				}
			}
			return &types.MsgHandlePendingTransferResponse{}, nil
		}
	}
	return nil, ErrNoPendingTransferFound

}
