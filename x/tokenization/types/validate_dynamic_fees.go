package types

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// ValidateDynamicFeeSchedule performs basic validation on a DynamicFeeSchedule.
func ValidateDynamicFeeSchedule(schedule *DynamicFeeSchedule) error {
	if schedule == nil {
		return nil
	}

	if len(schedule.Tiers) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "dynamic fee schedule must have at least one tier")
	}

	if schedule.FeeRecipient == "" {
		return sdkerrors.Wrap(ErrInvalidAddress, "dynamic fee schedule fee recipient is empty")
	}

	if err := ValidateAddress(schedule.FeeRecipient, false); err != nil {
		return sdkerrors.Wrap(ErrInvalidAddress, "invalid dynamic fee schedule fee recipient address")
	}

	if schedule.FeeDenom == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "dynamic fee schedule fee denom is empty")
	}

	for i, tier := range schedule.Tiers {
		if tier.BasisPoints > 10000 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "tier %d: basis points %d exceeds 10000 (100%%)", i, tier.BasisPoints)
		}

		if tier.MinAmount.GT(tier.MaxAmount) {
			return sdkerrors.Wrapf(ErrInvalidRequest, "tier %d: min amount %s > max amount %s", i, tier.MinAmount.String(), tier.MaxAmount.String())
		}

		// Check for overlapping tiers
		for j := i + 1; j < len(schedule.Tiers); j++ {
			otherTier := schedule.Tiers[j]
			// Tiers overlap if: tier.Min <= other.Max AND other.Min <= tier.Max
			if tier.MinAmount.LTE(otherTier.MaxAmount) && otherTier.MinAmount.LTE(tier.MaxAmount) {
				return sdkerrors.Wrapf(ErrInvalidRequest,
					"tiers %d and %d overlap: [%s, %s] and [%s, %s]",
					i, j,
					tier.MinAmount.String(), tier.MaxAmount.String(),
					otherTier.MinAmount.String(), otherTier.MaxAmount.String(),
				)
			}
		}
	}

	return nil
}

// ValidateCoinPerTokenMultiplier performs basic validation on a CoinPerTokenMultiplier.
func ValidateCoinPerTokenMultiplier(m *CoinPerTokenMultiplier) error {
	if m == nil {
		return nil
	}

	if len(m.CoinAmountPerToken) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "coin per token multiplier must have at least one coin")
	}

	for i, coin := range m.CoinAmountPerToken {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "coin per token multiplier coin %d has empty denom", i)
		}
		if coin.Amount.IsNil() || coin.Amount.IsZero() {
			return sdkerrors.Wrapf(ErrInvalidRequest, "coin per token multiplier coin %d has zero or nil amount", i)
		}
		if coin.Amount.IsNegative() {
			return sdkerrors.Wrapf(ErrInvalidRequest, "coin per token multiplier coin %d has negative amount", i)
		}
	}

	return nil
}

// ValidateTimeBasedRefundFormula performs basic validation on a TimeBasedRefundFormula.
func ValidateTimeBasedRefundFormula(formula *TimeBasedRefundFormula) error {
	if formula == nil {
		return nil
	}

	if len(formula.BaseRefundAmount) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "time-based refund formula must have at least one base refund coin")
	}

	if formula.TotalDuration.IsZero() {
		return sdkerrors.Wrap(ErrInvalidRequest, "time-based refund formula total duration cannot be zero")
	}

	if formula.EndTime.IsZero() {
		return sdkerrors.Wrap(ErrInvalidRequest, "time-based refund formula end time cannot be zero")
	}

	for i, coin := range formula.BaseRefundAmount {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "time-based refund coin %d has empty denom", i)
		}
		if coin.Amount.IsNil() || coin.Amount.IsZero() {
			return sdkerrors.Wrapf(ErrInvalidRequest, "time-based refund coin %d has zero or nil amount", i)
		}
		if coin.Amount.IsNegative() {
			return sdkerrors.Wrapf(ErrInvalidRequest, "time-based refund coin %d has negative amount", i)
		}
	}

	// EndTime should be > TotalDuration start (basic sanity)
	// EndTime - TotalDuration = subscription start time
	if formula.EndTime.LT(formula.TotalDuration) {
		return sdkerrors.Wrap(ErrInvalidRequest,
			fmt.Sprintf("time-based refund formula end time (%s) must be >= total duration (%s)",
				formula.EndTime.String(), formula.TotalDuration.String()))
	}

	// Sanity check: TotalDuration should not be absurdly large
	// (100 years in milliseconds = ~3.15e12, use 10 trillion as upper bound)
	maxDuration := sdkmath.NewUint(10_000_000_000_000) // ~317 years
	if formula.TotalDuration.GT(maxDuration) {
		return sdkerrors.Wrap(ErrInvalidRequest,
			fmt.Sprintf("time-based refund formula total duration %s exceeds maximum", formula.TotalDuration.String()))
	}

	return nil
}
