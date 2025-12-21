package keeper

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckStartsWithWrappedOrAliasDenom(denom string) bool {
	return strings.HasPrefix(denom, WrappedDenomPrefix) || strings.HasPrefix(denom, AliasDenomPrefix)
}

// GetPartsFromDenom parses a badges denom into its parts
func GetPartsFromDenom(denom string) ([]string, error) {
	if !CheckStartsWithWrappedOrAliasDenom(denom) {
		return nil, fmt.Errorf("invalid denom: %s", denom)
	}

	parts := strings.Split(denom, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid denom: %s", denom)
	}
	return parts, nil
}

// ParseDenomCollectionId extracts the collection ID from a badges denom
func ParseDenomCollectionId(denom string) (uint64, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return 0, err
	}

	// this is equivalent to split(':')[1]
	return strconv.ParseUint(parts[1], 10, 64)
}

// ParseDenomPath extracts the path from a badges denom
func ParseDenomPath(denom string) (string, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return "", err
	}
	// this is equivalent to split(':')[2]
	return parts[2], nil
}

// GetCorrespondingPath finds the CosmosCoinWrapperPath for a given denom
func GetCorrespondingPath(collection *badgestypes.TokenCollection, denom string) (*badgestypes.CosmosCoinWrapperPath, error) {
	baseDenom, err := ParseDenomPath(denom)
	if err != nil {
		return nil, err
	}

	// This is okay because we don't allow numeric chars in denoms
	numericStr := ""
	for _, char := range baseDenom {
		if char >= '0' && char <= '9' {
			numericStr += string(char)
		}
	}

	cosmosPaths := collection.CosmosCoinWrapperPaths
	for _, path := range cosmosPaths {
		if path.AllowOverrideWithAnyValidToken {
			// 1. Replace the {id} placeholder with the actual denom
			// 2. Convert all balance.tokenIds to the actual token ID
			if numericStr == "" {
				continue
			}

			idFromDenom := sdkmath.NewUintFromString(numericStr)
			path.Denom = strings.ReplaceAll(path.Denom, "{id}", idFromDenom.String())
			path.Balances = badgestypes.DeepCopyBalances(path.Balances)
			for _, balance := range path.Balances {
				balance.TokenIds = []*badgestypes.UintRange{
					{Start: idFromDenom, End: idFromDenom},
				}
			}
		}

		if path.Denom == baseDenom {
			return path, nil
		}
	}

	return nil, fmt.Errorf("path not found for denom: %s", denom)
}

// GetBalancesToTransfer calculates the balances to transfer for a given denom and amount
func GetBalancesToTransfer(collection *badgestypes.TokenCollection, denom string, amount sdkmath.Uint) ([]*badgestypes.Balance, error) {
	path, err := GetCorrespondingPath(collection, denom)
	if err != nil {
		return nil, err
	}

	balancesToTransfer := badgestypes.DeepCopyBalances(path.Balances)
	for _, balance := range balancesToTransfer {
		balance.Amount = balance.Amount.Mul(amount)
	}

	return balancesToTransfer, nil
}

// CheckIsAliasDenom checks if a denom is a wrapped badges denom
func (k Keeper) CheckIsAliasDenom(ctx sdk.Context, denom string) bool {
	if !CheckStartsWithWrappedOrAliasDenom(denom) {
		return false
	}

	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return false
	}

	path, err := GetCorrespondingPath(collection, denom)
	if err != nil {
		return false
	}

	// This is a little bit of an edge case
	// It is possible to have a badges: denom that is not the auto-converted denom
	// If this flag is true, we assume that they have to be wrapped first
	//
	// Ex: chaosnet denomination (badges:49:chaosnet)
	if path.AllowCosmosWrapping {
		return false
	}

	return true
}

// ParseCollectionFromDenom parses a collection from a badges denom
func (k Keeper) ParseCollectionFromDenom(ctx sdk.Context, denom string) (*badgestypes.TokenCollection, error) {
	collectionId, err := ParseDenomCollectionId(denom)
	if err != nil {
		return nil, err
	}

	collection, found := k.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return nil, customhookstypes.WrapErr(&ctx, badgestypes.ErrInvalidCollectionID, "collection %s not found",
			sdkmath.NewUint(collectionId).String())
	}

	return collection, nil
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

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
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

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
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

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
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

	// badgeslp: denom - calculate from badge balances
	// Parse collection from denom
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// Get the corresponding wrapper path
	path, err := GetCorrespondingPath(collection, denom)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// Get user's badge balance
	userBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address.String())

	// Use the same calculation as GetWrappableBalances
	maxWrappableAmount, err := k.calculateMaxWrappableAmount(ctx, userBalances.Balances, path.Balances)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	return sdkmath.NewIntFromBigInt(maxWrappableAmount.BigInt()), nil
}
