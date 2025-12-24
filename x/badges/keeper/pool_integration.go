package keeper

import (
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBalancesToTransferWithAlias calculates the balances to transfer for a given denom and amount
func GetBalancesToTransferWithAlias(collection *badgestypes.TokenCollection, denom string, amount sdkmath.Uint) ([]*badgestypes.Balance, error) {
	path, err := GetCorrespondingAliasPath(collection, denom)
	if err != nil {
		return nil, err
	}

	conversionAmount := path.Amount
	if conversionAmount.IsZero() || conversionAmount.IsNil() {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidRequest, "conversion amount is zero")
	}

	// Throw if not evenly divisible
	if !amount.Mod(conversionAmount).IsZero() {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidRequest, "amount is not evenly divisible by path amount")
	}

	multiplierToUse := amount.Quo(conversionAmount)
	balancesToTransfer := badgestypes.DeepCopyBalances(path.Balances)
	for _, balance := range balancesToTransfer {
		balance.Amount = balance.Amount.Mul(multiplierToUse)
	}

	return balancesToTransfer, nil
}

func (k Keeper) SetAllAutoApprovalFlagsForAddressUnsafe(ctx sdk.Context, collection *badgestypes.TokenCollection, address string) error {
	badgesMsgServer := NewMsgServerImpl(k)
	currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address)

	alreadyAutoApprovedAllIncomingTransfers := currBalances.AutoApproveAllIncomingTransfers
	alreadyAutoApprovedSelfInitiatedOutgoingTransfers := currBalances.AutoApproveSelfInitiatedOutgoingTransfers
	alreadyAutoApprovedSelfInitiatedIncomingTransfers := currBalances.AutoApproveSelfInitiatedIncomingTransfers

	autoApprovedAll := alreadyAutoApprovedAllIncomingTransfers && alreadyAutoApprovedSelfInitiatedOutgoingTransfers && alreadyAutoApprovedSelfInitiatedIncomingTransfers

	if !autoApprovedAll {
		// We override all approvals to be default allowed
		// Incoming - All, no matter what
		// Outgoing - Self-initiated
		//
		// This should cover the transfer to this address (rare edge case where default opt-in only)
		updateApprovalsMsg := &badgestypes.MsgUpdateUserApprovals{
			Creator:                               address,
			CollectionId:                          collection.CollectionId,
			UpdateAutoApproveAllIncomingTransfers: true,
			AutoApproveAllIncomingTransfers:       true,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers:       true,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedIncomingTransfers:       true,
		}
		_, err := badgesMsgServer.UpdateUserApprovals(ctx, updateApprovalsMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

// sendNativeTokensToAddressWithPoolApprovals sends native badges tokens to an address
func (k Keeper) sendNativeTokensToAddressWithPoolApprovals(ctx sdk.Context, poolAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransferWithAlias(collection, denom, amount)
	if err != nil {
		return err
	}

	err = k.SetAllAutoApprovalFlagsForAddressUnsafe(ctx, collection, poolAddress)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := NewMsgServerImpl(k)

	// Important: We should only allow auto-scanned approvals here
	//            Anything prioritized is potentially unsafe if we are using an IBC hook (where we cannot trust the sender)
	msg := &badgestypes.MsgTransferTokens{
		Creator:      poolAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        poolAddress,
				ToAddresses: []string{toAddress},
				Balances:    balancesToTransfer,
			},
		},
	}

	_, err = badgesMsgServer.TransferTokens(ctx, msg)
	return err
}

// sendNativeTokensFromAddressWithPoolApprovals sends native badges tokens from an address
func (k Keeper) sendNativeTokensFromAddressWithPoolApprovals(ctx sdk.Context, fromAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransferWithAlias(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := NewMsgServerImpl(k)

	// Just for sanity checks, we override all approvals to be default allowed
	// Incoming - All, no matter what
	// Outgoing - Self-initiated
	updateApprovalsMsg := &badgestypes.MsgUpdateUserApprovals{
		Creator:                               fromAddress,
		CollectionId:                          collection.CollectionId,
		UpdateAutoApproveAllIncomingTransfers: true,
		AutoApproveAllIncomingTransfers:       true,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers:       true,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedIncomingTransfers:       true,

		//One-time outgoing approval for the address to send tokens to the recipient
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*badgestypes.UserOutgoingApproval{
			{
				ToListId:          recipientAddress,
				InitiatedByListId: recipientAddress,
				TransferTimes:     []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				OwnershipTimes:    []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				TokenIds:          []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				Version:           sdkmath.NewUint(0),
				ApprovalId:        "one-time-outgoing",
			},
		},
	}
	_, err = badgesMsgServer.UpdateUserApprovals(ctx, updateApprovalsMsg)
	if err != nil {
		return err
	}

	// Important: We should only allow auto-scanned approvals here
	//            Anything prioritized is potentially unsafe if we are using an IBC hook (where we cannot trust the sender)
	// The one time outgoing approval is safe because it is hardcoded
	msg := &badgestypes.MsgTransferTokens{
		Creator:      recipientAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        fromAddress,
				ToAddresses: []string{recipientAddress},
				Balances:    balancesToTransfer,
				PrioritizedApprovals: []*badgestypes.ApprovalIdentifierDetails{
					{
						ApprovalId:      "one-time-outgoing",
						ApprovalLevel:   "outgoing",
						ApproverAddress: fromAddress,
						Version:         sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
	}

	_, err = badgesMsgServer.TransferTokens(ctx, msg)
	if err != nil {
		return err
	}

	// We then make sure that the address no longer has the one-time outgoing approval
	// This is needed as opposed to auto-deletion because technically the approval might not
	// be used if there is some forceful override (thus never deletes and we have a dangling approval)
	updateApprovalsMsg2 := &badgestypes.MsgUpdateUserApprovals{
		Creator:      fromAddress,
		CollectionId: collection.CollectionId,

		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*badgestypes.UserOutgoingApproval{},
	}
	_, err = badgesMsgServer.UpdateUserApprovals(ctx, updateApprovalsMsg2)
	if err != nil {
		return err
	}
	return nil
}

// sendNativeTokensToAddressWithPoolApprovals sends native badges tokens to an address
func (k Keeper) SendNativeTokensViaAliasDenom(ctx sdk.Context, recipientAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransferWithAlias(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := NewMsgServerImpl(k)

	// Important: We should only allow auto-scanned approvals here
	//            Anything prioritized is potentially unsafe if we are using an IBC hook (where we cannot trust the sender)
	msg := &badgestypes.MsgTransferTokens{
		Creator:      recipientAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        recipientAddress,
				ToAddresses: []string{toAddress},
				Balances:    balancesToTransfer,
			},
		},
	}

	_, err = badgesMsgServer.TransferTokens(ctx, msg)
	return err
}

// FundCommunityPoolViaAliasDenom funds the community pool using alias denom routing
// This handles the alias denom-specific logic (e.g., setting auto-approvals for the module address)
func (k Keeper) FundCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	// To accept incoming transfers (if disallowed by default)
	err = k.SetAllAutoApprovalFlagsForAddressUnsafe(ctx, collection, toAddress)
	if err != nil {
		return err
	}

	// No approvals to be set here
	return k.SendNativeTokensViaAliasDenom(ctx, fromAddress, toAddress, denom, amount)
}

// SpendFromCommunityPoolViaAliasDenom spends from the community pool using alias denom routing
// This handles the alias denom-specific logic (e.g., setting auto-approvals for the recipient address)
func (k Keeper) SpendFromCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	// To set outgoing transfers (if disallowed by default)
	err = k.SetAllAutoApprovalFlagsForAddressUnsafe(ctx, collection, fromAddress)
	if err != nil {
		return err
	}

	// Standard send from community pool to recipient
	return k.SendNativeTokensViaAliasDenom(ctx, fromAddress, toAddress, denom, amount)
}

// TODO: For both of these, I'd love to DRY more with sendManager. I just need to handle the pre/post approvals (and also the prioritized correctly which isn't supported natively)

// SendCoinsToPoolWithAliasRouting sends coins to a pool, wrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when to address is a pool address
// bankKeeper is required for sending non-badges coins
func (k Keeper) SendCoinsToPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if k.CheckIsAliasDenom(ctx, coin.Denom) {
			err := k.sendNativeTokensToAddressWithPoolApprovals(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			// otherwise, send the coins normally
			err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, from, to, &coin)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SendCoinsFromPoolWithAliasRouting sends coins from a pool, unwrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when from address is a pool address
// bankKeeper is required for sending non-badges coins
func (k Keeper) SendCoinsFromPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, unwrap it
	for _, coin := range coins {
		if k.CheckIsAliasDenom(ctx, coin.Denom) {

			err := k.sendNativeTokensFromAddressWithPoolApprovals(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

		} else {
			// otherwise, send the coins normally
			err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, from, to, &coin)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) GetSpendableCoinAmountBadgesLPOnly(ctx sdk.Context, address sdk.AccAddress, denom string) (sdkmath.Int, error) {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// Get the corresponding wrapper path
	path, err := GetCorrespondingAliasPath(collection, denom)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// Get user's badge balance
	userBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address.String())
	maxWrappableAmount, err := k.calculateMaxWrappableAmount(ctx, userBalances.Balances, path)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	return sdkmath.NewIntFromBigInt(maxWrappableAmount.BigInt()), nil
}
