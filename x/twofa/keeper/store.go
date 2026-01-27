package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	twofatypes "github.com/bitbadges/bitbadgeschain/x/twofa/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetUser2FARequirementsInStore sets the 2FA requirements for a user
func (k Keeper) SetUser2FARequirementsInStore(ctx sdk.Context, address string, requirements *twofatypes.User2FARequirements) error {
	// Validate address
	if err := types.ValidateAddress(address, false); err != nil {
		return sdkerrors.Wrap(err, "invalid address")
	}

	marshaled, err := k.cdc.Marshal(requirements)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal User2FARequirements failed")
	}

	store := k.getStore(ctx)
	store.Set(user2FARequirementsStoreKey(address), marshaled)
	return nil
}

// GetUser2FARequirementsFromStore gets the 2FA requirements for a user
func (k Keeper) GetUser2FARequirementsFromStore(ctx sdk.Context, address string) (*twofatypes.User2FARequirements, bool) {
	store := k.getStore(ctx)
	marshaled := store.Get(user2FARequirementsStoreKey(address))

	if len(marshaled) == 0 {
		return &twofatypes.User2FARequirements{
			MustOwnTokens:          []*types.MustOwnTokens{},
			DynamicStoreChallenges: []*types.DynamicStoreChallenge{},
		}, false
	}

	var requirements twofatypes.User2FARequirements
	k.cdc.MustUnmarshal(marshaled, &requirements)
	return &requirements, true
}

