package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckDynamicStoreChallenges validates and processes dynamic store challenges for an approval
// It checks if the initiator has sufficient remaining uses for each challenge and decrements the usage count
// Returns error if validation fails, nil on success
func (k Keeper) CheckDynamicStoreChallenges(
	ctx sdk.Context,
	challenges []*types.DynamicStoreChallenge,
	initiatedBy string,
	isPrioritizedApproval bool,
	addPotentialError func(bool, string),
	simulation bool,
) error {
	for _, challenge := range challenges {
		storeId := challenge.StoreId

		// Get the current value for the initiator
		dynamicStoreValue, found := k.GetDynamicStoreValueFromStore(ctx, storeId, initiatedBy)

		var val sdkmath.Uint
		if found {
			val = dynamicStoreValue.Value
		} else {
			// If no specific value found, get the default value from the store
			dynamicStore, foundStore := k.GetDynamicStoreFromStore(ctx, storeId)
			if !foundStore {
				errorMsg := fmt.Sprintf("dynamic store not found for storeId %s", storeId.String())
				addPotentialError(isPrioritizedApproval, errorMsg)
				return sdkerrors.New("dynamic_store_not_found", 1, errorMsg)
			}
			val = dynamicStore.DefaultValue
		}

		// Check if the initiator has remaining uses
		if val.Equal(sdkmath.NewUint(0)) {
			errorMsg := fmt.Sprintf("initiator has no remaining uses for dynamic store challenge storeId %s", storeId.String())
			addPotentialError(isPrioritizedApproval, errorMsg)
			return sdkerrors.New("no_remaining_uses", 1, errorMsg)
		}

		// Decrement the usage count only if not simulating
		if !simulation {
			newValue := val.SubUint64(1)
			if found {
				// Update existing value
				if err := k.SetDynamicStoreValueInStore(ctx, storeId, initiatedBy, newValue); err != nil {
					return sdkerrors.Wrapf(err, "failed to decrement dynamic store value for storeId %s", storeId.String())
				}
			} else {
				// Create new value with decremented default
				if err := k.SetDynamicStoreValueInStore(ctx, storeId, initiatedBy, newValue); err != nil {
					return sdkerrors.Wrapf(err, "failed to create dynamic store value for storeId %s", storeId.String())
				}
			}
		}
	}

	return nil
}

// SimulateDynamicStoreChallenges is a wrapper around CheckDynamicStoreChallenges for simulation
func (k Keeper) SimulateDynamicStoreChallenges(
	ctx sdk.Context,
	challenges []*types.DynamicStoreChallenge,
	initiatedBy string,
	isPrioritizedApproval bool,
	addPotentialError func(bool, string),
) error {
	return k.CheckDynamicStoreChallenges(ctx, challenges, initiatedBy, isPrioritizedApproval, addPotentialError, true)
}

// ExecuteDynamicStoreChallenges is a wrapper around CheckDynamicStoreChallenges for execution
func (k Keeper) ExecuteDynamicStoreChallenges(
	ctx sdk.Context,
	challenges []*types.DynamicStoreChallenge,
	initiatedBy string,
	isPrioritizedApproval bool,
	addPotentialError func(bool, string),
) error {
	return k.CheckDynamicStoreChallenges(ctx, challenges, initiatedBy, isPrioritizedApproval, addPotentialError, false)
}
