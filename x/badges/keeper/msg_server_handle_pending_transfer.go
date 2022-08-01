package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) HandlePendingTransfer(goCtx context.Context, msg *types.MsgHandlePendingTransfer) (*types.MsgHandlePendingTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	// Verify that the badge and subbadge exist and are valid
	err := k.AssertBadgeAndSubBadgeExists(ctx, msg.BadgeId, msg.SubbadgeId)
	if err != nil {
		return nil, err
	}

	// Handle the transfers
	balance_key := GetBalanceKey(creator_account_num, msg.BadgeId, msg.SubbadgeId)
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	} else {
		//In the future, we can make this a binary search since this is all sorted
		for _, pending_transfer := range badgeBalanceInfo.Pending {
			if pending_transfer.ThisPendingNonce == msg.ThisNonce {
				//We have four scenarios
				if pending_transfer.SendRequest && msg.Accept {
					return nil, ErrCantAcceptOwnTransferRequest
				}
				sentAndWantToCancel := pending_transfer.SendRequest && !msg.Accept
				receivedAndWantToAccept := !pending_transfer.SendRequest && msg.Accept
				received := !pending_transfer.SendRequest

				//Get the other party's account number, nonce, and information
				other_account_num := uint64(0)
				outgoing := false
				if pending_transfer.From == creator_account_num {
					outgoing = true
					other_account_num = pending_transfer.To
				} else if pending_transfer.To == creator_account_num {
					outgoing = false
					other_account_num = pending_transfer.From
				} else {
					panic("This pending transfer should always have the creator address as 'To' or 'From'")
				}
				other_balance_key := GetBalanceKey(other_account_num, msg.BadgeId, msg.SubbadgeId)
				other_nonce := pending_transfer.OtherPendingNonce

				//Remove from both parties' pending because it will always be removed (whether accepted or not)
				if err := k.RemovePending(ctx, balance_key, msg.ThisNonce, other_nonce); err != nil {
					return nil, err
				}

				if err = k.RemovePending(ctx, other_balance_key, other_nonce, msg.ThisNonce); err != nil {
					return nil, err
				}


				if sentAndWantToCancel {
					// If an outgoing transfer, we need to revert the balance and approval back to what it was
					if creator_account_num == pending_transfer.From {
						if err := k.AddToBadgeBalance(ctx, balance_key, pending_transfer.Amount); err != nil {
							return nil, err
						}

						
						//If it was sent via an approval, we need to add the approval back
						if pending_transfer.ApprovedBy != pending_transfer.From {
							if err = k.AddBalanceToApproval(ctx, balance_key, pending_transfer.Amount, pending_transfer.ApprovedBy); err != nil {
								return nil, err
							}
						}
					}
				} else if received {
					to_balance_key := balance_key
					from_balance_key := other_balance_key
					if outgoing {
						to_balance_key =  other_balance_key
						from_balance_key =  balance_key
					}

					if receivedAndWantToAccept {
						// We alteady removed from "From" balance so all we have to do now is add to "To" balance
						if err := k.AddToBadgeBalance(ctx, to_balance_key, pending_transfer.Amount); err != nil {
							return nil, err
						}
					} else {
						//If an incoming transfer and you want to reject it, we need to revert the balances and approvals back to what it was
						if !outgoing {
							if err := k.AddToBadgeBalance(ctx, from_balance_key, pending_transfer.Amount); err != nil {
								return nil, err
							}

							//If it was sent via an approval, we need to add the approval back
							if pending_transfer.ApprovedBy != pending_transfer.From {
								if err = k.AddBalanceToApproval(ctx, from_balance_key, pending_transfer.Amount, pending_transfer.ApprovedBy); err != nil {
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
}
