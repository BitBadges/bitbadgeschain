package app

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// A fresh chain gets both denoms from DefaultParams without needing any
// migration at all.
func TestDefaultParamsCarryBothUSDCDenoms(t *testing.T) {
	allowed := tokenizationtypes.DefaultParams().AllowedDenoms
	require.Contains(t, allowed, tokenizationtypes.USDCDenom, "canonical USDC")
	require.Contains(t, allowed, tokenizationtypes.USDCNobleDenom, "legacy USDC.n")
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
			name:  "canonical USDC: native Injective USDC, single hop",
			trace: "transfer/channel-40/erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a",
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

	// The trace hash is case-sensitive and Injective's bank module spells the
	// erc20 denom with a CHECKSUMMED address. A lowercase spelling silently
	// yields a different denom nobody would ever receive.
	require.NotEqual(t, tokenizationtypes.USDCDenom,
		ibcDenom("transfer/channel-40/erc20:0xa00c59ff5a080d2b954d0c75e46e22a0c371235a"),
		"lowercased erc20 address must not produce the canonical denom")

	// Nor is the canonical denom the forwarded Noble voucher (2-hop) route.
	require.NotEqual(t, tokenizationtypes.USDCDenom,
		ibcDenom("transfer/channel-40/transfer/channel-148/uusdc"),
		"canonical USDC is the single-hop native asset, not the forwarded Noble voucher")
}

// Mainnet does not get the canonical denom from DefaultParams (a running chain
// keeps whatever params are in state) and there is no upgrade migration for it:
// the cutover mechanism is a governance proposal carrying
// tokenization.MsgUpdateParams. This exercises that exact path — the same msg
// server, the same authority check — against pre-cutover state.
//
// MsgUpdateParams REPLACES the whole params object, so the proposal must carry
// the full current params with the canonical denom appended, never a partial
// object. The test mirrors that by mutating a queried copy.
func TestGovUpdateParamsAllowsInjectiveUSDC(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	tk := app.TokenizationKeeper

	// Pre-cutover mainnet state: canonical denom absent.
	params := tk.GetParams(ctx)
	params.AllowedDenoms = []string{tokenizationtypes.NativeDenom, tokenizationtypes.USDCNobleDenom}
	require.NoError(t, tk.SetParams(ctx, params))

	ms := tokenizationkeeper.NewMsgServerImpl(tk)

	// The proposal payload: full current params + the appended denom.
	proposed := tk.GetParams(ctx)
	proposed.AllowedDenoms = append(proposed.AllowedDenoms, tokenizationtypes.USDCDenom)

	// A non-authority sender must be rejected — only governance can do this.
	_, err := ms.UpdateParams(ctx, &tokenizationtypes.MsgUpdateParams{
		Authority: "invalid",
		Params:    proposed,
	})
	require.Error(t, err, "a non-governance authority must not be able to change AllowedDenoms")

	_, err = ms.UpdateParams(ctx, &tokenizationtypes.MsgUpdateParams{
		Authority: tk.GetAuthority(),
		Params:    proposed,
	})
	require.NoError(t, err, "the governance authority must be able to allow the canonical denom")

	after := tk.GetParams(ctx).AllowedDenoms
	require.Contains(t, after, tokenizationtypes.USDCDenom,
		"canonical USDC must be allowed after the gov proposal")
	require.Contains(t, after, tokenizationtypes.USDCNobleDenom,
		"legacy Noble USDC must stay allowed — 16 collections are backed by it and their escrow address derives from that denom string")
	require.Contains(t, after, tokenizationtypes.NativeDenom)
}
