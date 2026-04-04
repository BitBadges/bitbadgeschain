package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// PoAKeeper is the interface that a PoA keeper must satisfy for x/pot integration.
// Chain developers using PoA implement this interface wrapping their actual PoA keeper.
type PoAKeeper interface {
	// IterateActiveValidators walks validators with power > 0.
	// The callback receives the raw consensus address bytes, operator address string,
	// and voting power. Return true to stop iteration.
	IterateActiveValidators(ctx context.Context, fn func(consAddr []byte, operatorAddr string, power int64) bool) error

	// GetValidatorByConsAddr looks up a validator.
	GetValidatorByConsAddr(ctx context.Context, consAddr []byte) (operatorAddr string, power int64, err error)

	// SetValidatorPower updates a validator's voting power.
	SetValidatorPower(ctx context.Context, consAddr []byte, power int64) error
}

// PoAValidatorSetAdapter wraps a PoAKeeper into a ValidatorSetKeeper.
// It persists saved power via a PowerStore so that disabled validators
// can have their power restored on re-enable.
type PoAValidatorSetAdapter struct {
	poaKeeper  PoAKeeper
	powerStore types.PowerStore
}

// NewPoAAdapter creates a new PoAValidatorSetAdapter.
// The powerStore should be the x/pot Keeper (which implements types.PowerStore)
// to persist saved power across blocks.
func NewPoAAdapter(pk PoAKeeper, ps types.PowerStore) *PoAValidatorSetAdapter {
	return &PoAValidatorSetAdapter{poaKeeper: pk, powerStore: ps}
}

// IterateActiveValidators walks all PoA validators with positive power.
func (a *PoAValidatorSetAdapter) IterateActiveValidators(ctx context.Context, fn func(val types.ValidatorInfo) bool) error {
	return a.poaKeeper.IterateActiveValidators(ctx, func(consAddr []byte, operatorAddr string, power int64) bool {
		return fn(types.ValidatorInfo{
			ConsAddr:     sdk.ConsAddress(consAddr),
			OperatorAddr: operatorAddr,
			Power:        power,
		})
	})
}

// GetValidatorByConsAddr looks up a PoA validator by consensus address.
func (a *PoAValidatorSetAdapter) GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (types.ValidatorInfo, error) {
	operatorAddr, power, err := a.poaKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		return types.ValidatorInfo{}, err
	}
	return types.ValidatorInfo{
		ConsAddr:     consAddr,
		OperatorAddr: operatorAddr,
		Power:        power,
	}, nil
}

// DisableValidator saves the validator's current power and sets it to 0.
func (a *PoAValidatorSetAdapter) DisableValidator(ctx context.Context, consAddr sdk.ConsAddress) error {
	_, power, err := a.poaKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		return fmt.Errorf("cannot disable validator: %w", err)
	}

	// Save current power so it can be restored later.
	if power > 0 {
		a.powerStore.SetSavedPower(ctx, consAddr, power)
	}

	return a.poaKeeper.SetValidatorPower(ctx, consAddr, 0)
}

// EnableValidator restores the validator's saved power.
func (a *PoAValidatorSetAdapter) EnableValidator(ctx context.Context, consAddr sdk.ConsAddress) error {
	savedPower, found := a.powerStore.GetSavedPower(ctx, consAddr)
	if !found || savedPower <= 0 {
		// No saved power — use a default of 1 to re-enable.
		savedPower = 1
	}

	a.powerStore.RemoveSavedPower(ctx, consAddr)
	return a.poaKeeper.SetValidatorPower(ctx, consAddr, savedPower)
}

// IsValidatorJailed returns true if the validator's power is 0.
// In PoA, "jailed" means "disabled" (power == 0).
func (a *PoAValidatorSetAdapter) IsValidatorJailed(ctx context.Context, consAddr sdk.ConsAddress) (bool, error) {
	_, power, err := a.poaKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		return false, err
	}
	return power == 0, nil
}

// CanSafelyEnable always returns true for PoA — there is no slashing concept.
func (a *PoAValidatorSetAdapter) CanSafelyEnable(_ context.Context, _ sdk.ConsAddress) bool {
	return true
}

// Compile-time check that PoAValidatorSetAdapter implements ValidatorSetKeeper.
var _ types.ValidatorSetKeeper = (*PoAValidatorSetAdapter)(nil)
