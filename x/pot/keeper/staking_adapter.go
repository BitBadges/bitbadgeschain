package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// StakingValidatorSetAdapter wraps the existing StakingKeeper + SlashingKeeper
// into the abstract ValidatorSetKeeper interface.
type StakingValidatorSetAdapter struct {
	stakingKeeper  types.StakingKeeper
	slashingKeeper types.SlashingKeeper
}

// NewStakingAdapter creates a new StakingValidatorSetAdapter.
func NewStakingAdapter(sk types.StakingKeeper, slk types.SlashingKeeper) *StakingValidatorSetAdapter {
	return &StakingValidatorSetAdapter{stakingKeeper: sk, slashingKeeper: slk}
}

// IterateActiveValidators walks all bonded validators by power and converts
// each to a ValidatorInfo. The callback returns true to stop iteration.
func (a *StakingValidatorSetAdapter) IterateActiveValidators(ctx context.Context, fn func(val types.ValidatorInfo) bool) error {
	return a.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, val stakingtypes.ValidatorI) bool {
		consAddrBytes, err := val.GetConsAddr()
		if err != nil {
			return false // skip, continue iteration
		}

		// Convert valoper (bbvaloper...) to account address (bb1...).
		operatorAddr := val.GetOperator()
		valAddr, err := sdk.ValAddressFromBech32(operatorAddr)
		if err != nil {
			return false // skip
		}
		accAddr := sdk.AccAddress(valAddr).String()

		info := types.ValidatorInfo{
			ConsAddr:     sdk.ConsAddress(consAddrBytes),
			OperatorAddr: accAddr,
			Power:        0, // power is not used by the EndBlocker currently
		}
		return fn(info)
	})
}

// GetValidatorByConsAddr looks up a validator by consensus address.
func (a *StakingValidatorSetAdapter) GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (types.ValidatorInfo, error) {
	val, err := a.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	// Convert valoper to account address.
	operatorAddr := val.GetOperator()
	valAddr, valErr := sdk.ValAddressFromBech32(operatorAddr)
	accAddr := operatorAddr
	if valErr == nil {
		accAddr = sdk.AccAddress(valAddr).String()
	}

	return types.ValidatorInfo{
		ConsAddr:     consAddr,
		OperatorAddr: accAddr,
		Power:        0,
	}, nil
}

// DisableValidator jails a validator via x/staking.
func (a *StakingValidatorSetAdapter) DisableValidator(ctx context.Context, consAddr sdk.ConsAddress) error {
	return a.stakingKeeper.Jail(ctx, consAddr)
}

// EnableValidator unjails a validator via x/staking.
func (a *StakingValidatorSetAdapter) EnableValidator(ctx context.Context, consAddr sdk.ConsAddress) error {
	return a.stakingKeeper.Unjail(ctx, consAddr)
}

// IsValidatorJailed checks the staking module's jailed flag.
func (a *StakingValidatorSetAdapter) IsValidatorJailed(ctx context.Context, consAddr sdk.ConsAddress) (bool, error) {
	val, err := a.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		return false, fmt.Errorf("validator not found: %w", err)
	}
	return val.IsJailed(), nil
}

// CanSafelyEnable checks whether it is safe to unjail a validator by consulting
// the slashing module. Returns false if tombstoned or if the slashing jail
// period has not yet expired.
func (a *StakingValidatorSetAdapter) CanSafelyEnable(ctx context.Context, consAddr sdk.ConsAddress) bool {
	// Tombstoned validators are permanently banned.
	if a.slashingKeeper.IsTombstoned(ctx, consAddr) {
		return false
	}

	// Check jailed-until — respect slashing jail period.
	signingInfo, err := a.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if err != nil {
		// No signing info means validator was never slashing-jailed; safe to unjail.
		return true
	}

	// If JailedUntil is set and hasn't passed yet, don't enable.
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if !signingInfo.JailedUntil.IsZero() && !signingInfo.JailedUntil.Before(sdkCtx.BlockTime()) {
		return false
	}

	return true
}

// Compile-time check that StakingValidatorSetAdapter implements ValidatorSetKeeper.
var _ types.ValidatorSetKeeper = (*StakingValidatorSetAdapter)(nil)
