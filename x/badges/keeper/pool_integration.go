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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// CheckStartsWithBadges checks if a denom starts with "badges:" or "badgeslp:"
func CheckStartsWithBadges(denom string) bool {
	return strings.HasPrefix(denom, "badges:") || strings.HasPrefix(denom, "badgeslp:")
}

// GetPartsFromDenom parses a badges denom into its parts
func GetPartsFromDenom(denom string) ([]string, error) {
	if !CheckStartsWithBadges(denom) {
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

// CheckIsWrappedDenom checks if a denom is a wrapped badges denom
func (k Keeper) CheckIsWrappedDenom(ctx sdk.Context, denom string) bool {
	if !CheckStartsWithBadges(denom) {
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

// SendCoinsToPoolWithAliasRouting sends coins to a pool, wrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when to address is a pool address
// bankKeeper is required for sending non-badges coins
func (k Keeper) SendCoinsToPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.sendNativeTokensToAddressWithPoolApprovals(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return customhookstypes.WrapErr(&ctx, err, "failed to send coins to pool: %s",
					coin.Denom)
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
		if k.CheckIsWrappedDenom(ctx, coin.Denom) {

			err := k.sendNativeTokensFromAddressWithPoolApprovals(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

		} else {
			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return customhookstypes.WrapErr(&ctx, err, "failed to send native tokens from pool: %s",
					coin.Denom)
			}
		}
	}

	return nil
}

// FundCommunityPoolWithAliasRouting funds the community pool, wrapping badges denoms if needed.
// Used for taker fees
// bankKeeper and communityPoolKeeper are required for non-badges coins
func (k Keeper) FundCommunityPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		moduleAddress := authtypes.NewModuleAddress(distrtypes.ModuleName).String()

		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			collection, err := k.ParseCollectionFromDenom(ctx, coin.Denom)
			if err != nil {
				return err
			}

			// To accept incoming transfers (if disallowed by default)
			err = k.SetAllAutoApprovalFlagsForAddressUnsafe(ctx, collection, moduleAddress)
			if err != nil {
				return err
			}

			// No approvals to be set here
			err = k.SendNativeTokensViaAliasDenom(ctx, from.String(), moduleAddress, coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			err := k.distributionKeeper.FundCommunityPool(ctx, sdk.NewCoins(coin), from)
			if err != nil {
				return customhookstypes.WrapErr(&ctx, err, "failed to fund community pool: %s",
					coin.Denom)
			}
		}
	}

	return nil
}

// SendCoinWithAliasRouting sends a coin using the appropriate routing (wrapped badges or bank)
// For wrapped badges denoms, uses sendNativeTokensToAddressWithPoolApprovals which sets up recipient approvals
// For regular denoms, uses bank keeper SendCoins
func (k Keeper) SendCoinWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coin *sdk.Coin,
) error {
	// Check if this is a wrapped badges denom
	if k.CheckIsWrappedDenom(ctx, coin.Denom) {
		amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
		return k.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
	} else {
		return k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
	}
}

func (k Keeper) SendCoinsWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	// Check if this is a wrapped badges denom
	for _, coin := range coins {
		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
