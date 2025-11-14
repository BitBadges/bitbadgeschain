package keeper

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func CheckStartsWithBadges(denom string) bool {
	return strings.HasPrefix(denom, "badges:") || strings.HasPrefix(denom, "badgeslp:")
}

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

func ParseDenomCollectionId(denom string) (uint64, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return 0, err
	}

	// this is equivalent to split(':')[1]
	return strconv.ParseUint(parts[1], 10, 64)
}

func ParseDenomPath(denom string) (string, error) {
	parts, err := GetPartsFromDenom(denom)
	if err != nil {
		return "", err
	}
	// this is equivalent to split(':')[1]
	return parts[2], nil
}

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

func (k Keeper) ParseCollectionFromDenom(ctx sdk.Context, denom string) (*badgestypes.TokenCollection, error) {
	collectionId, err := ParseDenomCollectionId(denom)
	if err != nil {
		return nil, err
	}

	collection, found := k.badgesKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return nil, sdkerrors.Wrapf(badgestypes.ErrInvalidCollectionID, "collection %s not found", sdkmath.NewUint(collectionId).String())
	}

	return collection, nil
}

func (k Keeper) SendNativeTokensToPool(ctx sdk.Context, recipientAddress string, poolAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(k.badgesKeeper)

	currBalances, _ := k.badgesKeeper.GetBalanceOrApplyDefault(ctx, collection, poolAddress)

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
			Creator:                               poolAddress,
			CollectionId:                          collection.CollectionId,
			UpdateAutoApproveAllIncomingTransfers: true,
			AutoApproveAllIncomingTransfers:       true,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers:       true,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedIncomingTransfers:       true,
		}
		_, err = badgesMsgServer.UpdateUserApprovals(ctx, updateApprovalsMsg)
		if err != nil {
			return err
		}
	}

	msg := &badgestypes.MsgTransferTokens{
		Creator:      recipientAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        recipientAddress,
				ToAddresses: []string{poolAddress},
				Balances:    balancesToTransfer,
			},
		},
	}

	_, err = badgesMsgServer.TransferTokens(ctx, msg)
	return err
}

func (k Keeper) SendNativeTokensFromPool(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	collection, err := k.ParseCollectionFromDenom(ctx, denom)
	if err != nil {
		return err
	}

	balancesToTransfer, err := GetBalancesToTransfer(collection, denom, amount)
	if err != nil {
		return err
	}

	// Create and execute MsgTransferTokens to ensure proper event handling and validation
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(k.badgesKeeper)

	// Just for sanity checks, we override all approvals to be default allowed
	// Incoming - All, no matter what
	// Outgoing - Self-initiated
	updateApprovalsMsg := &badgestypes.MsgUpdateUserApprovals{
		Creator:                               poolAddress,
		CollectionId:                          collection.CollectionId,
		UpdateAutoApproveAllIncomingTransfers: true,
		AutoApproveAllIncomingTransfers:       true,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers:       true,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedIncomingTransfers:       true,

		//One-time outgoing approval for the pool to send tokens to the recipient
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

	msg := &badgestypes.MsgTransferTokens{
		Creator:      recipientAddress,
		CollectionId: collection.CollectionId,
		Transfers: []*badgestypes.Transfer{
			{
				From:        poolAddress,
				ToAddresses: []string{recipientAddress},
				Balances:    balancesToTransfer,
				PrioritizedApprovals: []*badgestypes.ApprovalIdentifierDetails{
					{
						ApprovalId:      "one-time-outgoing",
						ApprovalLevel:   "outgoing",
						ApproverAddress: poolAddress,
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

	// We then make sure that the pool no longer has the one-time outgoing approval
	// This is needed as opposed to auto-deletion because technically the approval might not
	// be used if there is some forceful override (thus never deletes and we have a dangling approval)
	updateApprovalsMsg2 := &badgestypes.MsgUpdateUserApprovals{
		Creator:      poolAddress,
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

// IMPORTANT: Should ONLY be called when to address is a pool address
func (k Keeper) SendCoinsToPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, wrap it
	for _, coin := range coins {
		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeTokensToPool(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IMPORTANT: Should ONLY be called when from address is a pool address
func (k Keeper) SendCoinsFromPoolWithUnwrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// if denom is a badges denom, unwrap it
	for _, coin := range coins {
		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeTokensFromPool(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}

		} else {
			// otherwise, send the coins normally
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Used for taker fees
func (k Keeper) FundCommunityPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		moduleAddress := authtypes.NewModuleAddress(distrtypes.ModuleName).String()

		if k.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.SendNativeTokensToPool(ctx, from.String(), moduleAddress, coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			err := k.communityPoolKeeper.FundCommunityPool(ctx, sdk.NewCoins(coin), from)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
