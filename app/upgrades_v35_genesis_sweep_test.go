package app

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/stretchr/testify/require"

	erc20types "github.com/cosmos/evm/x/erc20/types"

	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// TestV35NoLegacyDenomSurvivesInExportedGenesis is the backstop for the whole
// migration.
//
// Every other test asserts that a *specific* module was converted, which only
// catches the subsystems someone thought to look at. This one exports the full
// genesis after the upgrade and fails if the retired denom appears anywhere at
// all — so a store nobody remembered is caught by omission rather than by
// having been enumerated correctly.
//
// It is deliberately a blunt string search over the serialized genesis. The
// point is to not depend on knowing where to look.
func TestV35NoLegacyDenomSurvivesInExportedGenesis(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	seedLegacyChainState(t, app, ctx)
	seedLegacyDenomEverywhere(t, app, ctx, bk)

	// Sanity: the seeding actually put the legacy denom into state, otherwise
	// this test would pass vacuously.
	preExport := exportGenesisJSON(t, app, ctx)
	require.Contains(t, preExport, legacyDenom,
		"seed did not place the legacy denom in state; the sweep would pass vacuously")

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	postExport := exportGenesisJSON(t, app, ctx)

	if occurrences := findLegacyDenom(postExport); len(occurrences) > 0 {
		t.Fatalf("the retired denom %q still appears %d time(s) in exported genesis after the upgrade.\n"+
			"Each is a store the migration does not cover. Context around the first few:\n%s",
			legacyDenom, len(occurrences), strings.Join(occurrences, "\n"))
	}
}

// exportGenesisJSON serializes the app's module state the same way an export
// would, so the sweep sees whatever the modules actually persist.
func exportGenesisJSON(t *testing.T, app *App, ctx sdk.Context) string {
	t.Helper()
	genesis, err := app.ModuleManager.ExportGenesis(ctx, app.AppCodec())
	require.NoError(t, err)
	raw, err := json.Marshal(genesis)
	require.NoError(t, err)
	return string(raw)
}

// findLegacyDenom returns a short context window around each occurrence, so a
// failure names the module rather than just the count.
func findLegacyDenom(doc string) []string {
	// Word-boundary match so a denom that merely contains the string (an IBC
	// hash, say) does not produce a false positive.
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(legacyDenom) + `\b`)
	locs := re.FindAllStringIndex(doc, -1)

	out := make([]string, 0, len(locs))
	for i, loc := range locs {
		if i >= 5 {
			out = append(out, "  ... and more")
			break
		}
		start := loc[0] - 220
		if start < 0 {
			start = 0
		}
		end := loc[1] + 80
		if end > len(doc) {
			end = len(doc)
		}
		out = append(out, "  ..."+doc[start:end]+"...")
	}
	return out
}

// seedLegacyDenomEverywhere puts the legacy denom into as many distinct stores
// as can be reached through public keepers, so the sweep has something to catch
// in each. Without it the export is nearly empty and the test proves little.
func seedLegacyDenomEverywhere(t *testing.T, app *App, ctx sdk.Context, bk bankkeeper.BaseKeeper) {
	t.Helper()

	alice, bob := randAddr(), randAddr()
	fundLegacy(t, ctx, bk, alice, 5_000_000)
	fundLegacy(t, ctx, bk, bob, 7_000_000)

	legacyCoins := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1_000_000)))

	// A vesting account: locked amounts live on the account, not the balance.
	vestAddr := randAddr()
	fundLegacy(t, ctx, bk, vestAddr, 3_000_000)
	base := authtypes.NewBaseAccountWithAddress(vestAddr)
	base.AccountNumber = 9001
	// A fixed far-future unix time; ctx.BlockTime() is zero in tests, so
	// deriving from it yields a negative end time.
	const vestingEnd = int64(4_000_000_000)
	vestAcc, err := vestingtypes.NewDelayedVestingAccount(base, legacyCoins, vestingEnd)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, vestAcc)

	// A fee grant with a spend limit.
	require.NoError(t, app.FeeGrantKeeper.GrantAllowance(ctx, alice, bob, &feegrant.BasicAllowance{
		SpendLimit: legacyCoins,
	}))

	// An authz send authorization with a spend limit.
	require.NoError(t, app.AuthzKeeper.SaveGrant(ctx, bob, alice,
		banktypes.NewSendAuthorization(legacyCoins, nil), nil))

	// IBC total-escrow accounting for the native denom.
	app.TransferKeeper.SetTotalEscrowForDenom(ctx, sdk.NewCoin(legacyDenom, sdkmath.NewInt(2_500_000)))

	// An IBC rate-limit config keyed on the native denom. This one is worth
	// seeding explicitly: a stale config fails *open* — the module allows any
	// transfer with no matching config — so a missed rename silently removes
	// the throttle rather than blocking transfers.
	rlParams := app.IBCRateLimitKeeper.GetParams(ctx)
	rlParams.RateLimits = append(rlParams.RateLimits, ratelimittypes.RateLimitConfig{
		ChannelId: "channel-0",
		Denom:     legacyDenom,
	})
	require.NoError(t, app.IBCRateLimitKeeper.SetParams(ctx, rlParams))

	// An erc20 token pair mapping the native denom to a contract.
	app.ERC20Keeper.SetTokenPair(ctx, erc20types.TokenPair{
		Erc20Address:  "0x0000000000000000000000000000000000000bad",
		Denom:         legacyDenom,
		Enabled:       true,
		ContractOwner: erc20types.OWNER_MODULE,
	})
}
