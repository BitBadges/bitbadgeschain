package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	v34 "github.com/bitbadges/bitbadgeschain/app/upgrades/v34"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// The Injective-routed USDC denom has to be spendable through x/tokenization
// after the upgrade, or paid mints, prediction markets, payment requests and
// subscriptions all reject it as a disallowed denom.
func TestV34AllowsInjectiveUSDC(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	tk := *app.TokenizationKeeper

	// Simulate a chain that predates the new denom, which is what any running
	// chain looks like — DefaultParams only applies from genesis.
	params := tk.GetParams(ctx)
	params.AllowedDenoms = []string{tokenizationtypes.NativeDenom, tokenizationtypes.USDCNobleDenom}
	tk.SetParams(ctx, params)

	require.True(t, v34.AllowInjectiveUSDC(ctx, tk), "should report that it changed something")

	after := tk.GetParams(ctx).AllowedDenoms
	require.Contains(t, after, tokenizationtypes.USDCDenom, "Injective USDC must be allowed after the upgrade")
	require.Contains(t, after, tokenizationtypes.USDCNobleDenom,
		"legacy Noble USDC must stay allowed — 16 collections are backed by it and their escrow address derives from that denom string")
	require.Contains(t, after, tokenizationtypes.NativeDenom)
}

// Upgrade handlers are re-executed during replay and recovery, so a second run
// must not append a duplicate.
func TestV34AllowsInjectiveUSDCIsIdempotent(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	tk := *app.TokenizationKeeper

	// Start from pre-upgrade state. A fresh test app already carries the denom
	// via DefaultParams, so without this the first call is a no-op and the test
	// would assert nothing about idempotency.
	params := tk.GetParams(ctx)
	params.AllowedDenoms = []string{tokenizationtypes.NativeDenom, tokenizationtypes.USDCNobleDenom}
	tk.SetParams(ctx, params)

	require.True(t, v34.AllowInjectiveUSDC(ctx, tk))
	first := append([]string(nil), tk.GetParams(ctx).AllowedDenoms...)

	require.False(t, v34.AllowInjectiveUSDC(ctx, tk), "second run must find nothing to do")
	require.Equal(t, first, tk.GetParams(ctx).AllowedDenoms, "second run must not change the denom list")

	count := 0
	for _, d := range tk.GetParams(ctx).AllowedDenoms {
		if d == tokenizationtypes.USDCDenom {
			count++
		}
	}
	require.Equal(t, 1, count, "Injective USDC must appear exactly once")
}

// A fresh chain gets both denoms from DefaultParams without needing the
// upgrade at all.
func TestDefaultParamsCarryBothUSDCDenoms(t *testing.T) {
	allowed := tokenizationtypes.DefaultParams().AllowedDenoms
	require.Contains(t, allowed, tokenizationtypes.USDCDenom, "canonical USDC")
	require.Contains(t, allowed, tokenizationtypes.USDCNobleDenom, "legacy USDC.noble")
	require.Contains(t, allowed, tokenizationtypes.NativeDenom)
}

// The two USDC denoms must never collide — they are different IBC routes for
// the same underlying asset, and conflating them is the whole class of bug this
// change exists to avoid.
func TestUSDCDenomsAreDistinct(t *testing.T) {
	require.NotEqual(t, tokenizationtypes.USDCDenom, tokenizationtypes.USDCNobleDenom)
	require.Equal(t,
		"ibc/0E485657AEF4C39D551E7D53463734E4C445A96E6C814DC4C2FF0031470B40BB",
		tokenizationtypes.USDCDenom,
		"canonical USDC is sha256 of transfer/channel-40/transfer/channel-148/uusdc")
	require.Equal(t,
		"ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349",
		tokenizationtypes.USDCNobleDenom,
		"legacy USDC is sha256 of transfer/channel-2/uusdc")
}
