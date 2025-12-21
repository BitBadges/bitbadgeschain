package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckDynamicStoreChallenges validates dynamic store challenges for an approval
// It checks if the initiator has a true value for each challenge (read-only check)
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
		if challenge == nil {
			errorMsg := "challenge is nil"
			addPotentialError(isPrioritizedApproval, errorMsg)
			return sdkerrors.New("invalid_challenge", 1, errorMsg)
		}

		storeId := challenge.StoreId

		// Get the current value for the initiator
		dynamicStoreValue, found := k.GetDynamicStoreValueFromStore(ctx, storeId, initiatedBy)

		var val bool
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

		// Check if the initiator has a true value (read-only check, no updates)
		if !val {
			errorMsg := fmt.Sprintf("initiator does not have permission for dynamic store challenge storeId %s", storeId.String())
			addPotentialError(isPrioritizedApproval, errorMsg)
			return sdkerrors.New("no_permission", 1, errorMsg)
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
