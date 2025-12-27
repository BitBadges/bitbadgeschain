package keeper

import (
	"fmt"
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

	if path.Conversion == nil || path.Conversion.SideA == nil {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidRequest, "conversion or sideA is nil")
	}

	conversionAmount := path.Conversion.SideA.Amount
	if conversionAmount.IsZero() || conversionAmount.IsNil() {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidRequest, "conversion amount is zero")
	}

	// Throw if not evenly divisible
	if !amount.Mod(conversionAmount).IsZero() {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidRequest, "amount is not evenly divisible by path amount")
	}

	multiplierToUse := amount.Quo(conversionAmount)
	balancesToTransfer := badgestypes.DeepCopyBalances(path.Conversion.SideB)
	for _, balance := range balancesToTransfer {
		balance.Amount = balance.Amount.Mul(multiplierToUse)
	}

	return balancesToTransfer, nil
}

// setAutoApprovalFlagsIfNeeded sets auto-approval flags on a balance if they're not already set.
// Returns true if any flags were changed, false otherwise.
// This is a DRY helper to avoid repeating the flagsChanged pattern.
func setAutoApprovalFlagsIfNeeded(balance *badgestypes.UserBalanceStore) bool {
	flagsChanged := false

	// Set AutoApproveAllIncomingTransfers only if not already set
	if !balance.AutoApproveAllIncomingTransfers {
		balance.AutoApproveAllIncomingTransfers = true
		flagsChanged = true
	}

	// Set AutoApproveSelfInitiatedOutgoingTransfers only if not already set
	if !balance.AutoApproveSelfInitiatedOutgoingTransfers {
		balance.AutoApproveSelfInitiatedOutgoingTransfers = true
		flagsChanged = true
	}

	// Set AutoApproveSelfInitiatedIncomingTransfers only if not already set
	if !balance.AutoApproveSelfInitiatedIncomingTransfers {
		balance.AutoApproveSelfInitiatedIncomingTransfers = true
		flagsChanged = true
	}

	return flagsChanged
}

// This function follows the same pattern as setAutoApproveFlagsForPathAddress to ensure
// consistent behavior and prevent unintended overrides of user-configured settings.
func (k Keeper) SetAllAutoApprovalFlagsForPoolAddress(ctx sdk.Context, collection *badgestypes.TokenCollection, address string) error {
	currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address)

	// Set flags if needed (DRY helper)
	if setAutoApprovalFlagsIfNeeded(currBalances) {
		err := k.SetBalanceForAddress(ctx, collection, address, currBalances)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to set auto-approve flags for pool/path address: %s", address)
		}
	}

	return nil
}

// SetAllAutoApprovalFlagsForIntermediateAddress sets auto-approval flags for an intermediate address.
// This is used for IBC hooks where intermediate addresses (derived from channel and sender) need
// auto-approval to function correctly. Intermediate addresses are deterministic system addresses
// that are not pool addresses but still need auto-approval for IBC hook operations.
//
// Security:
//   - Only sets flags if they're not already set (prevents overriding existing settings)
//   - Each flag is checked individually before setting
//   - This function does not validate that the address is a pool/path address since
//     intermediate addresses are a different type of system address
func (k Keeper) SetAllAutoApprovalFlagsForIntermediateAddress(ctx sdk.Context, collection *badgestypes.TokenCollection, address string) error {
	// Get current balances
	currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address)

	// Set flags if needed (DRY helper)
	if setAutoApprovalFlagsIfNeeded(currBalances) {
		err := k.SetBalanceForAddress(ctx, collection, address, currBalances)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to set auto-approve flags for intermediate address: %s", address)
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

	err = k.SetAllAutoApprovalFlagsForPoolAddress(ctx, collection, poolAddress)
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

// SendNativeTokensFromAddressWithPoolApprovals sends native badges tokens from an address
// Security: Uses unique approval IDs and ensures cleanup even on failure to prevent approval reuse.
// This function is exported for testing purposes to verify security properties.
func (k Keeper) SendNativeTokensFromAddressWithPoolApprovals(ctx sdk.Context, fromAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransferWithAlias(collection, denom, amount)
	if err != nil {
		return err
	}

	// Generate unique approval ID to prevent reuse attacks
	// Format: "one-time-outgoing-{blockHeight}-{collectionId}-{hash}"
	// This ensures each approval is unique and cannot be reused
	blockHeight := ctx.BlockHeight()
	approvalId := fmt.Sprintf("one-time-outgoing-%d-%s-%s", blockHeight, collection.CollectionId.String(), recipientAddress)

	// Security: The version must be incremented to prevent replay attacks
	// Using version 0 would allow the approval to be reused if deletion fails
	approvalVersion := k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", fromAddress, approvalId)

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := NewMsgServerImpl(k)

	// Helper function to delete the one-time approval
	deleteOneTimeApproval := func() error {
		cleanupMsg := &badgestypes.MsgUpdateUserApprovals{
			Creator:                 fromAddress,
			CollectionId:            collection.CollectionId,
			UpdateOutgoingApprovals: true,
			OutgoingApprovals:       []*badgestypes.UserOutgoingApproval{},
		}
		_, err := badgesMsgServer.UpdateUserApprovals(ctx, cleanupMsg)
		return err
	}

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

		// One-time outgoing approval for the address to send tokens to the recipient
		// Security: Uses unique approval ID to prevent reuse attacks
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*badgestypes.UserOutgoingApproval{
			{
				ToListId:          recipientAddress,
				InitiatedByListId: recipientAddress,
				TransferTimes:     []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				OwnershipTimes:    []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				TokenIds:          []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				Version:           approvalVersion,
				ApprovalId:        approvalId,
			},
		},
	}
	_, err = badgesMsgServer.UpdateUserApprovals(ctx, updateApprovalsMsg)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to create one-time approval: %s", approvalId)
	}

	// Important: We should only allow auto-scanned approvals here
	//            Anything prioritized is potentially unsafe if we are using an IBC hook (where we cannot trust the sender)
	// The one-time outgoing approval is safe because it uses a unique ID and version
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
						ApprovalId:      approvalId,
						ApprovalLevel:   "outgoing",
						ApproverAddress: fromAddress,
						Version:         approvalVersion,
					},
				},
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
	}

	_, err = badgesMsgServer.TransferTokens(ctx, msg)
	if err != nil {
		// Transfer failed - clean up the approval before returning
		cleanupErr := deleteOneTimeApproval()
		if cleanupErr != nil {
			// Log cleanup error but return the original transfer error
			return sdkerrors.Wrapf(err, "transfer failed for one-time approval: %s (cleanup also failed: %v)", approvalId, cleanupErr)
		}
		return sdkerrors.Wrapf(err, "transfer failed for one-time approval: %s", approvalId)
	}

	// Transfer succeeded - delete the approval
	err = deleteOneTimeApproval()
	if err != nil {
		// Log error but don't fail - transfer already succeeded
		// This is a cleanup operation, so we return success even if cleanup fails
		// The approval will be cleaned up on next access or can be manually deleted
		return sdkerrors.Wrapf(err, "transfer succeeded but failed to delete one-time approval: %s", approvalId)
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
	err = k.SetAllAutoApprovalFlagsForPoolAddress(ctx, collection, toAddress)
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
	err = k.SetAllAutoApprovalFlagsForPoolAddress(ctx, collection, fromAddress)
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

			err := k.SendNativeTokensFromAddressWithPoolApprovals(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
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
