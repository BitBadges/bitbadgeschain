package app

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

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

// ibcDenom derives an IBC denom from its full trace exactly the way ibc-go's
// transfer module does: sha256 of the trace, uppercase hex, "ibc/" prefix.
//
// Tests must derive this rather than restate the literal. A constant asserted
// against a second copy of itself proves nothing, and these constants are the
// difference between a spendable balance and an unreachable one.
func ibcDenom(trace string) string {
	sum := sha256.Sum256([]byte(trace))
	return "ibc/" + strings.ToUpper(hex.EncodeToString(sum[:]))
}

// Every IBC denom constant must equal the hash of the trace its doc comment
// claims. A typo in either the constant or the documented route fails here.
func TestIBCDenomConstantsMatchTheirTraces(t *testing.T) {
	for _, tc := range []struct {
		name  string
		trace string
		got   string
	}{
		{
			name:  "canonical USDC via Injective",
			trace: "transfer/channel-40/transfer/channel-148/uusdc",
			got:   tokenizationtypes.USDCDenom,
		},
		{
			name:  "legacy USDC direct from Noble",
			trace: "transfer/channel-2/uusdc",
			got:   tokenizationtypes.USDCNobleDenom,
		},
		{
			name:  "ATOM",
			trace: "transfer/channel-3/uatom",
			got:   tokenizationtypes.ATOMDenom,
		},
		{
			name:  "OSMO",
			trace: "transfer/channel-0/uosmo",
			got:   tokenizationtypes.OSMODenom,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, ibcDenom(tc.trace), tc.got,
				"constant does not hash from its documented trace %q", tc.trace)
		})
	}
}

// The two USDC denoms must never collide — they are different IBC routes for
// the same underlying asset, and conflating them is the whole class of bug this
// change exists to avoid.
func TestUSDCDenomsAreDistinct(t *testing.T) {
	require.NotEqual(t, tokenizationtypes.USDCDenom, tokenizationtypes.USDCNobleDenom)
	require.NotEqual(t, "transfer/channel-40/transfer/channel-148/uusdc", "transfer/channel-2/uusdc")
}

// The upgrade handler itself has to reach AllowInjectiveUSDC. Calling
// AllowInjectiveUSDC directly proves the function works but not that anything
// invokes it — dropping the call site would leave every other test in this file
// green while the mainnet upgrade silently did nothing.
func TestV34UpgradeHandlerAllowsInjectiveUSDC(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	tk := *app.TokenizationKeeper

	// Pre-upgrade state: a running chain that predates the canonical denom.
	params := tk.GetParams(ctx)
	params.AllowedDenoms = []string{tokenizationtypes.NativeDenom, tokenizationtypes.USDCNobleDenom}
	require.NoError(t, tk.SetParams(ctx, params))

	handler := v34.CreateUpgradeHandler(
		app.ModuleManager,
		app.Configurator(),
		app.AccountKeeper,
		tk,
		app.PoolManagerKeeper,
		app.IBCRateLimitKeeper,
	)

	fromVM := app.ModuleManager.GetVersionMap()
	toVM, err := handler(ctx, upgradetypes.Plan{Name: v34.UpgradeName}, fromVM)
	require.NoError(t, err, "v34 upgrade handler must run cleanly")
	require.NotEmpty(t, toVM)

	after := tk.GetParams(ctx).AllowedDenoms
	require.Contains(t, after, tokenizationtypes.USDCDenom,
		"canonical USDC must be allowed after running the real upgrade handler")
	require.Contains(t, after, tokenizationtypes.USDCNobleDenom,
		"legacy Noble USDC must survive the upgrade — backed-path escrows derive from that denom string")
}

// The handler is only reachable on mainnet if it is registered under the
// upgrade name the governance plan will carry.
func TestV34UpgradeHandlerIsRegistered(t *testing.T) {
	app := Setup(false)
	require.True(t, app.UpgradeKeeper.HasHandler(v34.UpgradeName),
		"v34 upgrade handler must be registered with the upgrade keeper")
}
