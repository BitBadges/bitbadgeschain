package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
)

func TestCoinPullAllowance_FullSubscriptionLifecycle(t *testing.T) {
	const (
		subscriber = "bb1subscriber"
		merchant   = "bb1merchant"
		keeperBot  = "bb1keeperbot"
		dayMs      = uint64(86_400_000)
		periodMs   = 30 * dayMs
		startMs    = uint64(1_700_000_000_000)
	)

	terms := types.SubscriptionTerms{
		CollectionId:    42,
		ApprovalId:      "pro-monthly",
		ApprovalVersion: 1,
		PriceDenom:      "uusdc",
		PriceAmount:     sdkmath.NewUint(20_000_000), // 20 USDC
		PayoutAddress:   merchant,
		PeriodLengthMs:  periodMs,
	}

	reg := types.NewPullAllowanceRegistry()
	reg.SetBalance(subscriber, "uusdc", sdkmath.NewUint(50_000_000)) // 50 USDC

	// 1. Initial subscriber consent grant.
	allowance, err := reg.GrantAllowance(
		"sub-allow-1",
		subscriber,
		subscriber,
		terms,
		sdkmath.NewUint(20_500_000), // 20 USDC + up to 0.5 USDC keeper fee per period
		sdkmath.NewUint(500_000),
		startMs,
		startMs+12*periodMs,
	)
	require.NoError(t, err)
	require.Equal(t, terms.Hash(), allowance.TermsHash)

	// 2. First periodic capped pull (Period 0) succeeds.
	res0, err := reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+1000, sdkmath.NewUint(100_000))
	require.NoError(t, err)
	require.Equal(t, types.PullStatusSucceeded, res0.Status)
	require.Equal(t, uint64(0), res0.PeriodIndex)
	require.Equal(t, sdkmath.NewUint(20_000_000), reg.GetBalance(merchant, "uusdc"))
	require.Equal(t, sdkmath.NewUint(100_000), reg.GetBalance(keeperBot, "uusdc"))
	require.Equal(t, sdkmath.NewUint(29_900_000), reg.GetBalance(subscriber, "uusdc"))

	// Second pull in the SAME period is rejected by hard per-period cap.
	_, err = reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+5*dayMs, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrPeriodCapExceeded)

	// Period 1 rollover automatically resets period quota and succeeds.
	res1, err := reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+periodMs+1000, sdkmath.ZeroUint())
	require.NoError(t, err)
	require.Equal(t, types.PullStatusSucceeded, res1.Status)
	require.Equal(t, uint64(1), res1.PeriodIndex)
	require.Equal(t, sdkmath.NewUint(9_900_000), reg.GetBalance(subscriber, "uusdc"))

	// 3. Insufficient balance in Period 2 triggers graceful failure without bricking state.
	res2Fail, err := reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+2*periodMs+1000, sdkmath.ZeroUint())
	require.NoError(t, err)
	require.Equal(t, types.PullStatusInsufficientBalance, res2Fail.Status)
	require.Equal(t, sdkmath.NewUint(9_900_000), reg.GetBalance(subscriber, "uusdc"), "balance must remain untouched")
	require.Equal(t, sdkmath.NewUint(20_500_000), res2Fail.RemainingPeriodQuota, "failed pull must not consume period quota")

	snap, ok := reg.GetAllowance("sub-allow-1")
	require.True(t, ok)
	require.Equal(t, uint64(1), snap.FailedBalanceChecks)
	require.False(t, snap.Revoked)

	// After subscriber tops up, retry in Period 2 succeeds cleanly.
	reg.SetBalance(subscriber, "uusdc", sdkmath.NewUint(30_000_000))
	res2Ok, err := reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+2*periodMs+2*dayMs, sdkmath.ZeroUint())
	require.NoError(t, err)
	require.Equal(t, types.PullStatusSucceeded, res2Ok.Status)
	require.Equal(t, uint64(2), res2Ok.PeriodIndex)

	// 4. Instant mid-period revocation: non-subscriber is rejected, subscriber succeeds.
	require.ErrorIs(t, reg.RevokeAllowance("sub-allow-1", merchant, startMs+2*periodMs+3*dayMs), types.ErrUnauthorizedRevocation)
	require.NoError(t, reg.RevokeAllowance("sub-allow-1", subscriber, startMs+2*periodMs+3*dayMs))

	// Subsequent pull in Period 3 immediately fails with ErrAllowanceRevoked.
	_, err = reg.ExecutePull("sub-allow-1", terms, keeperBot, startMs+3*periodMs+1000, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrAllowanceRevoked)
}

func TestCoinPullAllowance_RejectsMerchantPriceAndTermsMutation(t *testing.T) {
	const (
		subscriber = "bb1subscriber"
		merchant   = "bb1merchant"
		attacker   = "bb1attacker"
		periodMs   = uint64(30 * 86_400_000)
		startMs    = uint64(1_700_000_000_000)
	)

	origTerms := types.SubscriptionTerms{
		CollectionId:    7,
		ApprovalId:      "starter",
		ApprovalVersion: 1,
		PriceDenom:      "uusdc",
		PriceAmount:     sdkmath.NewUint(10_000_000),
		PayoutAddress:   merchant,
		PeriodLengthMs:  periodMs,
	}

	reg := types.NewPullAllowanceRegistry()
	reg.SetBalance(subscriber, "uusdc", sdkmath.NewUint(100_000_000))
	_, err := reg.GrantAllowance("allow-mut", subscriber, subscriber, origTerms, sdkmath.NewUint(50_000_000), sdkmath.ZeroUint(), startMs, startMs+6*periodMs)
	require.NoError(t, err)

	// Merchant silently raises price (even though still below subscriber's max cap!).
	raisedPrice := origTerms
	raisedPrice.PriceAmount = sdkmath.NewUint(15_000_000)
	_, err = reg.ExecutePull("allow-mut", raisedPrice, "", startMs+1000, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrTermsMutated)

	// Merchant bumps approval version.
	bumpedVersion := origTerms
	bumpedVersion.ApprovalVersion = 2
	_, err = reg.ExecutePull("allow-mut", bumpedVersion, "", startMs+1000, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrTermsMutated)

	// Merchant redirects payout address.
	redirectedPayout := origTerms
	redirectedPayout.PayoutAddress = attacker
	_, err = reg.ExecutePull("allow-mut", redirectedPayout, "", startMs+1000, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrTermsMutated)

	// Subscriber balance must remain 100% untouched after all rejected attempts.
	require.Equal(t, sdkmath.NewUint(100_000_000), reg.GetBalance(subscriber, "uusdc"))
}

func TestCoinPullAllowance_ExpiryAndExecutorFeeBounds(t *testing.T) {
	const (
		subscriber = "bb1subscriber"
		merchant   = "bb1merchant"
		periodMs   = uint64(30 * 86_400_000)
		startMs    = uint64(1_700_000_000_000)
	)

	terms := types.SubscriptionTerms{
		CollectionId:    9,
		ApprovalId:      "basic",
		ApprovalVersion: 1,
		PriceDenom:      "uusdc",
		PriceAmount:     sdkmath.NewUint(5_000_000),
		PayoutAddress:   merchant,
		PeriodLengthMs:  periodMs,
	}

	reg := types.NewPullAllowanceRegistry()
	reg.SetBalance(subscriber, "uusdc", sdkmath.NewUint(100_000_000))
	_, err := reg.GrantAllowance("allow-exp", subscriber, subscriber, terms, sdkmath.NewUint(6_000_000), sdkmath.NewUint(100_000), startMs, startMs+periodMs)
	require.NoError(t, err)

	// Before start time -> rejected.
	_, err = reg.ExecutePull("allow-exp", terms, "", startMs-1, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrAllowanceNotYetActive)

	// Executor fee above cap -> rejected.
	_, err = reg.ExecutePull("allow-exp", terms, "bb1greedy", startMs+1, sdkmath.NewUint(200_000))
	require.ErrorIs(t, err, types.ErrExecutorFeeExceeded)

	// After validUntil -> rejected.
	_, err = reg.ExecutePull("allow-exp", terms, "", startMs+periodMs+1, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrAllowanceExpired)

	// Cap below terms price -> rejected at grant time.
	_, err = reg.GrantAllowance("allow-lowcap", subscriber, subscriber, terms, sdkmath.NewUint(1_000_000), sdkmath.ZeroUint(), startMs, startMs+periodMs)
	require.ErrorIs(t, err, types.ErrInvalidAllowanceConfig)

	// Monotonic timestamp guard: pull in period 0 then replay older timestamp -> rejected.
	_, err = reg.ExecutePull("allow-exp", terms, "", startMs+5000, sdkmath.ZeroUint())
	require.NoError(t, err)
	_, err = reg.ExecutePull("allow-exp", terms, "", startMs+4000, sdkmath.ZeroUint())
	require.ErrorIs(t, err, types.ErrInvalidAllowanceConfig)
}
