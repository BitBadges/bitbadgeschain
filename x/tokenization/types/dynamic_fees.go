package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeTier defines a single tier in a dynamic fee schedule.
// When the token transfer amount falls within [MinAmount, MaxAmount],
// the fee is computed as: transferAmount * BasisPoints / 10000.
type FeeTier struct {
	// Minimum token transfer amount for this tier (inclusive).
	MinAmount sdkmath.Uint
	// Maximum token transfer amount for this tier (inclusive).
	// Use math.MaxUint for unbounded upper tiers.
	MaxAmount sdkmath.Uint
	// Fee in basis points (100 = 1%, 10000 = 100%).
	BasisPoints uint64
}

// DynamicFeeSchedule defines a tiered, percentage-based fee that scales
// with the actual token transfer amount. The first matching tier wins.
type DynamicFeeSchedule struct {
	// Tiered fee brackets. First matching bracket wins.
	Tiers []FeeTier
	// Address that receives the fee.
	FeeRecipient string
	// Denom of the coin to charge (e.g. "ubadge").
	FeeDenom string
}

// CoinPerTokenMultiplier enables proportional pricing for non-backed
// collections. When set on a CoinTransfer, the coin transfer amount
// becomes CoinAmountPerToken * actualTokenTransferAmount instead of
// the fixed coins amount.
type CoinPerTokenMultiplier struct {
	// Cost per 1 token unit, expressed as sdk.Coins.
	CoinAmountPerToken []sdk.Coin
}

// TimeBasedRefundFormula enables pro-rated subscription refunds.
// When set on a CoinTransfer, the actual refund amount is computed as:
//
//	refund = BaseRefundAmount * (timeRemaining / TotalDuration)
//
// where timeRemaining = endTime - currentBlockTime.
type TimeBasedRefundFormula struct {
	// Full refund amount if cancelled immediately (time remaining == total duration).
	BaseRefundAmount []sdk.Coin
	// Total subscription duration in milliseconds.
	TotalDuration sdkmath.Uint
	// The end time of the subscription in milliseconds (Unix epoch).
	// Used to compute timeRemaining = EndTime - currentBlockTime.
	EndTime sdkmath.Uint
}

// ComputeDynamicFee calculates the fee for a given token transfer amount
// using the first matching tier in the schedule.
// Returns the fee amount and whether a matching tier was found.
func (d *DynamicFeeSchedule) ComputeDynamicFee(tokenTransferAmount sdkmath.Uint) (sdkmath.Int, bool) {
	for _, tier := range d.Tiers {
		if tokenTransferAmount.GTE(tier.MinAmount) && tokenTransferAmount.LTE(tier.MaxAmount) {
			// fee = transferAmount * basisPoints / 10000
			fee := tokenTransferAmount.Mul(sdkmath.NewUint(tier.BasisPoints)).Quo(sdkmath.NewUint(10000))
			return sdkmath.NewIntFromBigInt(fee.BigInt()), true
		}
	}
	return sdkmath.ZeroInt(), false
}

// ComputeMeteredAmount calculates the total coin amounts when using
// CoinPerTokenMultiplier. Each coin in the multiplier is scaled by
// the token transfer amount.
func (m *CoinPerTokenMultiplier) ComputeMeteredAmount(tokenTransferAmount sdkmath.Uint) []sdk.Coin {
	result := make([]sdk.Coin, 0, len(m.CoinAmountPerToken))
	for _, coin := range m.CoinAmountPerToken {
		scaledAmount := sdkmath.NewIntFromBigInt(
			coin.Amount.BigInt(),
		).Mul(sdkmath.NewIntFromBigInt(tokenTransferAmount.BigInt()))
		result = append(result, sdk.NewCoin(coin.Denom, scaledAmount))
	}
	return result
}

// ComputeRefundAmount calculates the pro-rated refund based on remaining time.
// Returns the refund coins. If the subscription has expired (currentTime >= endTime),
// all refund amounts are zero.
func (t *TimeBasedRefundFormula) ComputeRefundAmount(currentTimeMs sdkmath.Uint) []sdk.Coin {
	result := make([]sdk.Coin, 0, len(t.BaseRefundAmount))

	// If current time >= end time, subscription expired, no refund
	if currentTimeMs.GTE(t.EndTime) {
		for _, coin := range t.BaseRefundAmount {
			result = append(result, sdk.NewCoin(coin.Denom, sdkmath.ZeroInt()))
		}
		return result
	}

	// timeRemaining = endTime - currentTime
	timeRemaining := t.EndTime.Sub(currentTimeMs)

	// Cap timeRemaining at TotalDuration (should not exceed, but be safe)
	if timeRemaining.GT(t.TotalDuration) {
		timeRemaining = t.TotalDuration
	}

	for _, coin := range t.BaseRefundAmount {
		// refund = baseAmount * timeRemaining / totalDuration
		refundAmount := sdkmath.NewIntFromBigInt(coin.Amount.BigInt()).
			Mul(sdkmath.NewIntFromBigInt(timeRemaining.BigInt())).
			Quo(sdkmath.NewIntFromBigInt(t.TotalDuration.BigInt()))
		result = append(result, sdk.NewCoin(coin.Denom, refundAmount))
	}
	return result
}

// GetTotalTokenTransferAmount sums the Amount field across all balances,
// weighted by the number of (tokenId, ownershipTime) pairs each balance covers.
// This gives the total number of token units being transferred.
//
// TODO: This is a simplified calculation that sums amounts directly.
// A more precise version would multiply each balance.Amount by the number
// of individual token IDs and ownership time units it covers, but that
// requires expanding ranges which can be expensive. For now, we sum
// balance.Amount * numberOfTokenIdRanges * numberOfOwnershipTimeRanges
// as a practical approximation that works for most single-range transfers.
func GetTotalTokenTransferAmount(balances []*Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		if balance == nil || balance.Amount.IsZero() {
			continue
		}
		// For each balance entry, the amount applies to all (tokenId, ownershipTime) combinations.
		// For dynamic fee purposes, we sum the raw amounts.
		// In the common case (single token ID range, single ownership time range),
		// this is simply the transfer amount.
		total = total.Add(balance.Amount)
	}
	return total
}
