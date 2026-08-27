package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// ---------------------------------------------------------------------------
// PowerReduction: the premise the whole redenomination rests on
// ---------------------------------------------------------------------------

// TestV35PowerReductionScalesByTheSameFactorAsTheBondDenom pins the premise that
// makes moving PowerReduction safe.
//
// Consensus power is tokens/PowerReduction. The redenomination multiplies the
// numerator by 10^9, so PowerReduction has to move by exactly 10^9 too — no
// more, no less — or every validator's power changes. Nothing else in the tree
// asserts the relationship between the two constants, which is how a reviewer
// came to believe the two had been moved independently.
//
// The second half of the test is the concrete consequence, checked against a
// real mainnet figure rather than a round number, so the assertion cannot be
// satisfied by an arithmetic identity that happens to hold for powers of ten.
func TestV35PowerReductionScalesByTheSameFactorAsTheBondDenom(t *testing.T) {
	// The bond denom moves by ConversionFactor. PowerReduction must move by the
	// same factor off the SDK default it started at.
	require.Equal(t,
		v35.LegacyPowerReduction.Mul(v35.ConversionFactor).String(),
		appparams.PowerReduction.String(),
		"PowerReduction must scale by exactly the same factor as the denom it counts")

	// And it must actually be installed, not merely declared: the packages that
	// call InitSDKConfigWithoutSeal without building the app rely on this.
	require.Equal(t, appparams.PowerReduction.String(), sdk.DefaultPowerReduction.String(),
		"the declared PowerReduction must be the one the SDK is using")

	// bitbadges-1's largest validator at the time this was written.
	const mainnetValidatorUbadge = int64(307_572_905_653_630)

	before := sdkmath.NewInt(mainnetValidatorUbadge).Quo(v35.LegacyPowerReduction)
	after := sdkmath.NewInt(mainnetValidatorUbadge).Mul(v35.ConversionFactor).Quo(appparams.PowerReduction)

	require.Equal(t, before.String(), after.String(),
		"a mainnet-sized validator's consensus power must be identical across the upgrade")
	require.Equal(t, "307572905", after.String(),
		"and it must be the power the live chain actually reports")

	// The SDK's own conversion, so this is not just the test doing the same
	// division twice with different constants.
	val := stakingtypes.Validator{
		Status: stakingtypes.Bonded,
		Tokens: sdkmath.NewInt(mainnetValidatorUbadge).Mul(v35.ConversionFactor),
	}
	require.Equal(t, before.Int64(), val.ConsensusPower(sdk.DefaultPowerReduction),
		"the SDK's own power calculation must agree")
}

// The guard that makes the dev-versus-mainnet divergence safe.
//
// Mainnet bonds ubadge, which this upgrade redenominates, so PowerReduction
// moving with it is correct. The local dev chain bonds ustake, which this
// upgrade does NOT touch — running the same binary there would divide every
// validator's power by 10^9 and collapse the set to zero power with no error
// anywhere. So the handler refuses before it writes anything.
func TestV35RefusesToRunWhenTheBondDenomIsNotRedenominated(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	// The dev chain's bond denom, which this upgrade leaves alone.
	params, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	params.BondDenom = "ustake"
	require.NoError(t, app.StakingKeeper.SetParams(ctx, params))

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 1_000)

	err = v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app))
	require.Error(t, err)
	require.Contains(t, err.Error(), "ustake")
	require.Contains(t, err.Error(), "PowerReduction")

	// It must refuse *before* touching state, or a partly-migrated chain is left
	// behind for the operator to unpick.
	require.Equal(t, sdkmath.NewInt(1_000), bk.GetBalance(ctx, alice, legacyDenom).Amount,
		"the handler must not have written anything before refusing")
	require.True(t, bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).IsZero())
}

// The mainnet configuration is accepted. Without this, the test above passes
// just as well against a guard that refuses everything.
func TestV35AcceptsTheMainnetBondDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))
}

// ---------------------------------------------------------------------------
// x/precisebank: the dust DoS
// ---------------------------------------------------------------------------

// One ubadge sent to the precisebank reserve before the upgrade height used to
// halt the chain.
//
// The reserve address is a plain module address that was never on the bank's
// blocked list, so anybody could put coins on it — and on mainnet somebody did:
// it held 3 ubadge of dust at the time this was written. Step 2 of the handler
// turns that into 3*10^9 abadge of surplus, the old two-sided reserve == owed
// check fails, CustomUpgradeHandlerLogic returns an error, and the upgrade
// BeginBlocker panics every node. Cost to the attacker: 10^-9 BADGE plus gas.
func TestV35PreciseBankSurplusDustDoesNotHaltTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 4)

	// One whole ubadge of legitimate fractional balance, correctly backed...
	seedPreciseBankFractionalBalances(t, app, ctx, bk, map[string]int64{
		alice.String(): 1_000_000_000,
	}, 0)
	// ...plus the attacker's single ubadge of dust on top of it.
	fundPreciseBankReserve(t, ctx, bk, 1)

	supplyBefore := bk.GetSupply(ctx, legacyDenom).Amount

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)),
		"a surplus on the reserve must not be able to halt the upgrade")

	factor := v35.ConversionFactor

	// The legitimate holder is paid in full: 4 ubadge of bank balance plus the
	// 10^9 fractional units, which are one abadge each after the redenomination.
	require.Equal(t, sdkmath.NewInt(4).Mul(factor).AddRaw(1_000_000_000),
		bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount,
		"the surplus must not change what an honest holder is owed")

	reserveAddr := authtypes.NewModuleAddress(v35.PreciseBankStoreKey)
	require.True(t, bk.GetBalance(ctx, reserveAddr, appparams.BaseCoinUnit).IsZero(),
		"the reserve must be drained, surplus included")

	// The surplus is burned, so supply drops by exactly the dust and balances
	// still sum to supply. Anything else would leave supply inflated against
	// coins nobody can reach.
	summed := sdkmath.ZeroInt()
	bk.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(coin.Amount)
		}
		return false
	})
	newSupply := bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount
	require.Equal(t, newSupply, summed, "balances must sum to supply")
	require.Equal(t, supplyBefore.Mul(factor).Sub(sdkmath.NewInt(1).Mul(factor)), newSupply,
		"supply must fall by exactly the burned surplus and nothing else")
}

// The other side of the one-sided check: a reserve that is genuinely short must
// still refuse. Widening the check must not have widened it into "never fails".
func TestV35PreciseBankStillRefusesWhenTheReserveIsShort(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 1)

	// Two whole ubadge claimed, one ubadge behind it: short by exactly one.
	seedPreciseBankFractionalBalances(t, app, ctx, bk, map[string]int64{
		alice.String(): 2_000_000_000,
	}, 0)
	drainPreciseBankReserve(t, ctx, bk, 1)

	err := v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app))
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not back it")
}

// Defence in depth for the same DoS: after v35 nothing can put coins on the
// reserve address in the first place.
func TestV35PreciseBankReserveRejectsIncomingCoins(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	reserveAddr := authtypes.NewModuleAddress(v35.PreciseBankStoreKey)
	require.True(t, bk.BlockedAddr(reserveAddr),
		"the retired precisebank reserve must be a blocked address")

	attacker := randAddr()
	fundLegacy(t, ctx, bk, attacker, 10)

	// Through the msg server, which is the path an attacker actually has and the
	// only one that consults the blocked list — keeper.SendCoins deliberately
	// does not, because module-internal moves have to be able to reach these
	// accounts.
	msgServer := bankkeeper.NewMsgServerImpl(app.BankKeeper)
	_, err := msgServer.Send(ctx, banktypes.NewMsgSend(
		attacker, reserveAddr,
		sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1))),
	))
	require.Error(t, err, "a send to the retired reserve must be refused")
	require.Contains(t, err.Error(), "not allowed to receive funds")

	require.True(t, bk.GetBalance(ctx, reserveAddr, legacyDenom).IsZero(),
		"and nothing may have landed there")
}

// drainPreciseBankReserve moves ubadge off the reserve, to build the
// under-backed state the migration must refuse.
//
// The coins go to an ordinary account rather than being burned, so balances
// still sum to supply on the way in. Otherwise RedenominateBank's own
// unaccounted-supply check fires first and the test passes for the wrong reason.
func drainPreciseBankReserve(t *testing.T, ctx sdk.Context, bk bankkeeper.BaseKeeper, ubadge int64) {
	t.Helper()
	reserveAddr := authtypes.NewModuleAddress(v35.PreciseBankStoreKey)
	sink := randAddr()

	remaining := bk.GetBalance(ctx, reserveAddr, legacyDenom).Amount.SubRaw(ubadge)
	require.False(t, remaining.IsNegative())
	require.NoError(t, bk.UncheckedSetBalance(ctx, reserveAddr, sdk.NewCoin(legacyDenom, remaining)))
	require.NoError(t, bk.UncheckedSetBalance(ctx, sink, sdk.NewCoin(legacyDenom, sdkmath.NewInt(ubadge))))
}

// ---------------------------------------------------------------------------
// x/ibc-rate-limit: caps that fail closed
// ---------------------------------------------------------------------------

// Renaming a rate-limit config's denom without rescaling its caps is strictly
// worse than doing neither.
//
// The module allows any transfer with no matching config, so a stale config
// fails open. Renaming makes it match again — and then enforces caps that are
// still in the old scale. Mainnet's ubadge configs carry supply-shift caps of
// 10^16 and 3*10^16; left alone, a 10,000,000-BADGE daily ceiling becomes a
// 0.01-BADGE one and every non-trivial BADGE IBC transfer is rejected from the
// upgrade height onward.
func TestV35RateLimitCapsScaleWithTheRenamedDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	// The shape mainnet actually runs: supply-shift caps at 10^16/HOUR and
	// 3*10^16/DAY, address caps at 10^15/HOUR and 10^16/DAY.
	const foreignDenom = "ibc/ABCDEF0123456789"
	mustInt := func(s string) sdkmath.Int {
		v, ok := sdkmath.NewIntFromString(s)
		require.True(t, ok)
		return v
	}

	require.NoError(t, app.IBCRateLimitKeeper.SetParams(ctx, ratelimittypes.Params{
		RateLimits: []ratelimittypes.RateLimitConfig{
			{
				ChannelId: "channel-0",
				Denom:     legacyDenom,
				SupplyShiftLimits: []ratelimittypes.TimeframeLimit{
					{MaxAmount: mustInt("10000000000000000"), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, TimeframeDuration: 1},
					{MaxAmount: mustInt("30000000000000000"), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, TimeframeDuration: 1},
					// A disabled limit. Zero means "off" to the keeper and must
					// stay zero, not become zero-times-something-else.
					{MaxAmount: sdkmath.ZeroInt(), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, TimeframeDuration: 1},
				},
				AddressLimits: []ratelimittypes.AddressLimit{
					{MaxTransfers: 5, MaxAmount: mustInt("1000000000000000"), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, TimeframeDuration: 1},
				},
				UniqueSenderLimits: []ratelimittypes.UniqueSenderLimit{
					// A count, not an amount. Must not move.
					{MaxUniqueSenders: 7, TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, TimeframeDuration: 1},
				},
			},
			{
				ChannelId: "channel-0",
				Denom:     foreignDenom,
				SupplyShiftLimits: []ratelimittypes.TimeframeLimit{
					{MaxAmount: mustInt("12345"), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, TimeframeDuration: 1},
				},
			},
		},
	}))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	params := app.IBCRateLimitKeeper.GetParams(ctx)
	require.Len(t, params.RateLimits, 2)

	var migrated, foreign ratelimittypes.RateLimitConfig
	for _, cfg := range params.RateLimits {
		switch cfg.Denom {
		case appparams.BaseCoinUnit:
			migrated = cfg
		case foreignDenom:
			foreign = cfg
		default:
			t.Fatalf("unexpected denom %q left in the rate limit params", cfg.Denom)
		}
	}

	factor := v35.ConversionFactor
	require.Equal(t, mustInt("10000000000000000").Mul(factor).String(),
		migrated.SupplyShiftLimits[0].MaxAmount.String(),
		"an hourly BADGE cap must buy the same amount of BADGE after the upgrade")
	require.Equal(t, mustInt("30000000000000000").Mul(factor).String(),
		migrated.SupplyShiftLimits[1].MaxAmount.String())
	require.True(t, migrated.SupplyShiftLimits[2].MaxAmount.IsZero(),
		"a disabled limit must stay disabled")
	require.Equal(t, mustInt("1000000000000000").Mul(factor).String(),
		migrated.AddressLimits[0].MaxAmount.String())
	require.EqualValues(t, 5, migrated.AddressLimits[0].MaxTransfers,
		"a transfer *count* is not an amount and must not move")
	require.EqualValues(t, 7, migrated.UniqueSenderLimits[0].MaxUniqueSenders,
		"a sender count is not an amount and must not move")

	require.Equal(t, "12345", foreign.SupplyShiftLimits[0].MaxAmount.String(),
		"a config on a foreign denom must be untouched")
}

// The in-flight accumulators are keyed by denom too, and hold amounts. Orphaned,
// every rate-limit window silently restarts at zero at the upgrade height, which
// briefly allows up to twice the configured cap across the boundary.
func TestV35RateLimitInFlightFlowsAreRekeyedAndScaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const channel = "channel-0"
	const foreignDenom = "ibc/ABCDEF0123456789"
	sender := randAddr().String()

	const tfType = int32(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR)
	const tfDur = int64(1)

	store := ctx.KVStore(app.GetKey(ratelimittypes.StoreKey))
	cdc := app.AppCodec()

	setFlow := func(denom string, netFlow int64) {
		flow := ratelimittypes.ChannelFlow{NetFlow: sdkmath.NewInt(netFlow)}
		store.Set(ratelimittypes.ChannelFlowKey(channel, denom, tfType, tfDur), cdc.MustMarshal(&flow))
	}
	setAddrData := func(denom string, total int64) {
		data := ratelimittypes.AddressTransferData{TransferCount: 3, TotalAmount: sdkmath.NewInt(total)}
		store.Set(ratelimittypes.AddressTransferDataKey(sender, channel, denom, tfType, tfDur), cdc.MustMarshal(&data))
	}

	setFlow(legacyDenom, 1_000_000)
	setFlow(foreignDenom, 777)
	setAddrData(legacyDenom, 500_000)
	setAddrData(foreignDenom, 42)

	// The paired window records carry no amount but are keyed the same way; a
	// window left behind makes its flow look unopened.
	window := ratelimittypes.ChannelFlowWindow{WindowStart: 12345, WindowDuration: 600}
	store.Set(ratelimittypes.ChannelFlowWindowKey(channel, legacyDenom, tfType, tfDur), cdc.MustMarshal(&window))
	store.Set(ratelimittypes.AddressTransferWindowKey(sender, channel, legacyDenom, tfType, tfDur), cdc.MustMarshal(&window))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor

	require.Nil(t, store.Get(ratelimittypes.ChannelFlowKey(channel, legacyDenom, tfType, tfDur)),
		"the retired denom's flow record must not survive")
	var flow ratelimittypes.ChannelFlow
	cdc.MustUnmarshal(store.Get(ratelimittypes.ChannelFlowKey(channel, appparams.BaseCoinUnit, tfType, tfDur)), &flow)
	require.Equal(t, sdkmath.NewInt(1_000_000).Mul(factor).String(), flow.NetFlow.String(),
		"the accumulated flow is an amount and must move with the denom")

	require.Nil(t, store.Get(ratelimittypes.AddressTransferDataKey(sender, channel, legacyDenom, tfType, tfDur)))
	var data ratelimittypes.AddressTransferData
	cdc.MustUnmarshal(store.Get(ratelimittypes.AddressTransferDataKey(sender, channel, appparams.BaseCoinUnit, tfType, tfDur)), &data)
	require.Equal(t, sdkmath.NewInt(500_000).Mul(factor).String(), data.TotalAmount.String())
	require.EqualValues(t, 3, data.TransferCount, "a transfer count must not scale")

	var movedWindow ratelimittypes.ChannelFlowWindow
	cdc.MustUnmarshal(store.Get(ratelimittypes.ChannelFlowWindowKey(channel, appparams.BaseCoinUnit, tfType, tfDur)), &movedWindow)
	require.Equal(t, window, movedWindow,
		"the window must follow its flow unchanged, or the flow looks unopened")
	require.Nil(t, store.Get(ratelimittypes.ChannelFlowWindowKey(channel, legacyDenom, tfType, tfDur)))
	require.NotNil(t, store.Get(ratelimittypes.AddressTransferWindowKey(sender, channel, appparams.BaseCoinUnit, tfType, tfDur)))

	// Foreign denoms must be left exactly where they were.
	var foreignFlow ratelimittypes.ChannelFlow
	cdc.MustUnmarshal(store.Get(ratelimittypes.ChannelFlowKey(channel, foreignDenom, tfType, tfDur)), &foreignFlow)
	require.Equal(t, "777", foreignFlow.NetFlow.String())
	var foreignData ratelimittypes.AddressTransferData
	cdc.MustUnmarshal(store.Get(ratelimittypes.AddressTransferDataKey(sender, channel, foreignDenom, tfType, tfDur)), &foreignData)
	require.Equal(t, "42", foreignData.TotalAmount.String())
}

// ---------------------------------------------------------------------------
// x/authz: StakeAuthorization
// ---------------------------------------------------------------------------

// StakeAuthorization.MaxTokens is a coin amount that fell through the switch's
// default branch, whose comment claimed no unhandled type carried one.
//
// It is a single optional sdk.Coin rather than an sdk.Coins, so the helper every
// other case uses does not apply — which is exactly why it was easy to miss.
// Left alone, a 100-BADGE delegation cap becomes 10^-7 BADGE, and because
// Accept() writes the remainder back on every partial spend, the wrong scale is
// baked into a fresh grant the first time the grant is used.
func TestV35StakeAuthorizationMaxTokensIsScaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	granter, grantee := randAddr(), randAddr()
	capped := sdk.NewCoin(legacyDenom, sdkmath.NewInt(100_000_000_000)) // 100 BADGE at 9 decimals

	// The SDK rejects a grant with neither an allow nor a deny list.
	allowed := []sdk.ValAddress{sdk.ValAddress(randAddr())}
	auth, err := stakingtypes.NewStakeAuthorization(allowed, nil,
		stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, &capped)
	require.NoError(t, err)
	require.NoError(t, app.AuthzKeeper.SaveGrant(ctx, grantee, granter, auth, nil))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	stored, _ := app.AuthzKeeper.GetAuthorization(ctx, grantee, granter,
		sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}))
	require.NotNil(t, stored, "the grant must survive the migration")

	migrated, ok := stored.(*stakingtypes.StakeAuthorization)
	require.True(t, ok)
	require.NotNil(t, migrated.MaxTokens)
	require.Equal(t, appparams.BaseCoinUnit, migrated.MaxTokens.Denom)
	require.Equal(t, capped.Amount.Mul(v35.ConversionFactor).String(), migrated.MaxTokens.Amount.String(),
		"a 100-BADGE delegation cap must still be a 100-BADGE cap")
}

// The nil case is normal — StakeAuthorization reads a nil MaxTokens as "no cap"
// — and must not be turned into a zero-coin cap, which would refuse every
// delegation. A migration that dereferences unconditionally would panic here.
func TestV35StakeAuthorizationWithNoCapIsLeftAlone(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	granter, grantee := randAddr(), randAddr()
	allowed := []sdk.ValAddress{sdk.ValAddress(randAddr())}
	auth, err := stakingtypes.NewStakeAuthorization(allowed, nil,
		stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, nil)
	require.NoError(t, err)
	require.NoError(t, app.AuthzKeeper.SaveGrant(ctx, grantee, granter, auth, nil))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	stored, _ := app.AuthzKeeper.GetAuthorization(ctx, grantee, granter,
		sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}))
	require.NotNil(t, stored)
	migrated, ok := stored.(*stakingtypes.StakeAuthorization)
	require.True(t, ok)
	require.Nil(t, migrated.MaxTokens, "an uncapped grant must stay uncapped")
}

// The default branch of the authorization switch is now a hard error, so every
// amountless type on this chain has to be named explicitly. GenericAuthorization
// is the one that would otherwise halt the upgrade on any chain with a generic
// grant — which is most of them.
func TestV35GenericAuthorizationDoesNotHaltTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	granter, grantee := randAddr(), randAddr()
	auth := authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{}))
	require.NoError(t, app.AuthzKeeper.SaveGrant(ctx, grantee, granter, auth, nil))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)),
		"a generic grant carries no amount and must not fail the upgrade")

	stored, _ := app.AuthzKeeper.GetAuthorization(ctx, grantee, granter,
		sdk.MsgTypeURL(&banktypes.MsgSend{}))
	require.NotNil(t, stored, "and it must still be there afterwards")
}

// ---------------------------------------------------------------------------
// x/poolmanager: accrued taker fees
// ---------------------------------------------------------------------------

// The taker-fee trackers put the denom in the key and an amount in the value, so
// the rename orphans accrued BADGE fees under a denom nothing reads and epoch
// distribution pays out zero.
func TestV35TakerFeeTrackersAreRekeyedAndScaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"
	pmk := app.PoolManagerKeeper

	require.NoError(t, pmk.UpdateTakerFeeTrackerForStakersByDenom(ctx, legacyDenom, osmomath.NewInt(1_000)))
	require.NoError(t, pmk.UpdateTakerFeeTrackerForCommunityPoolByDenom(ctx, legacyDenom, osmomath.NewInt(2_000)))
	require.NoError(t, pmk.UpdateTakerFeeTrackerForBurnByDenom(ctx, legacyDenom, osmomath.NewInt(3_000)))
	require.NoError(t, pmk.UpdateTakerFeeTrackerForStakersByDenom(ctx, foreignDenom, osmomath.NewInt(9)))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor

	for _, tc := range []struct {
		name string
		get  func(sdk.Context, string) (sdk.Coin, error)
		want int64
	}{
		{"stakers", pmk.GetTakerFeeTrackerForStakersByDenom, 1_000},
		{"community pool", pmk.GetTakerFeeTrackerForCommunityPoolByDenom, 2_000},
		{"burn", pmk.GetTakerFeeTrackerForBurnByDenom, 3_000},
	} {
		migrated, err := tc.get(ctx, appparams.BaseCoinUnit)
		require.NoError(t, err)
		require.Equal(t, sdkmath.NewInt(tc.want).Mul(factor).String(), migrated.Amount.String(),
			"%s tracker must carry the accrued fees over at the new scale", tc.name)

		orphaned, err := tc.get(ctx, legacyDenom)
		require.NoError(t, err)
		require.True(t, orphaned.Amount.IsZero(),
			"%s tracker must leave nothing behind under the retired denom", tc.name)
	}

	foreign, err := pmk.GetTakerFeeTrackerForStakersByDenom(ctx, foreignDenom)
	require.NoError(t, err)
	require.Equal(t, "9", foreign.Amount.String(), "a foreign denom's tracker must not move")
}

// The taker-fee-share accrual is keyed by *two* denoms and the value is an
// amount of the second one, so which of the two moved decides whether the value
// scales.
func TestV35TakerFeeSkimAccrualsAreRekeyedAndScaledOnlyOnTheChargedDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"
	pmk := app.PoolManagerKeeper

	// charged in the redenominated denom -> the amount moves
	require.NoError(t, pmk.SetTakerFeeShareDenomsToAccruedValue(ctx, foreignDenom, legacyDenom, osmomath.NewInt(100)))
	// shared by the redenominated denom, charged in a foreign one -> re-key only
	require.NoError(t, pmk.SetTakerFeeShareDenomsToAccruedValue(ctx, legacyDenom, foreignDenom, osmomath.NewInt(200)))
	// neither -> untouched
	require.NoError(t, pmk.SetTakerFeeShareDenomsToAccruedValue(ctx, foreignDenom, foreignDenom, osmomath.NewInt(300)))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor

	charged, err := pmk.GetTakerFeeShareDenomsToAccruedValue(ctx, foreignDenom, appparams.BaseCoinUnit)
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewInt(100).Mul(factor).String(), charged.String(),
		"the accrued value is an amount of the charged denom, so it moves with it")

	shared, err := pmk.GetTakerFeeShareDenomsToAccruedValue(ctx, appparams.BaseCoinUnit, foreignDenom)
	require.NoError(t, err)
	require.Equal(t, "200", shared.String(),
		"only the share denom moved; the value is a foreign amount and must not scale")

	untouched, err := pmk.GetTakerFeeShareDenomsToAccruedValue(ctx, foreignDenom, foreignDenom)
	require.NoError(t, err)
	require.Equal(t, "300", untouched.String())

	_, err = pmk.GetTakerFeeShareDenomsToAccruedValue(ctx, foreignDenom, legacyDenom)
	require.Error(t, err, "nothing must be left under the retired denom")
}

// FormatTakerFeeShareAgreementKey embeds the denom, and the keeper serves every
// lookup from an in-memory cache built from that key. Re-keying the store
// without repairing the cache leaves the agreement resolvable under the retired
// denom and invisible under the live one for the rest of the process's life.
func TestV35TakerFeeShareAgreementIsRekeyedAndTheCacheFollows(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	pmk := app.PoolManagerKeeper
	require.NoError(t, pmk.SetTakerFeeShareAgreementForDenom(ctx, poolmanagertypes.TakerFeeShareAgreement{
		Denom:       legacyDenom,
		SkimPercent: osmomath.MustNewDecFromStr("0.01"),
		SkimAddress: randAddr().String(),
	}))

	// SetTakerFeeShareAgreementForDenom seeds the cache, so this reproduces a
	// running node rather than a cold-started one.
	_, found := pmk.GetTakerFeeShareAgreementFromDenomUNSAFE(legacyDenom)
	require.True(t, found, "precondition: the cache holds the retired denom")

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	stored, found := pmk.GetTakerFeeShareAgreementFromDenomNoCache(ctx, appparams.BaseCoinUnit)
	require.True(t, found, "the agreement must be re-keyed onto the live denom")
	require.Equal(t, appparams.BaseCoinUnit, stored.Denom)
	require.Equal(t, "0.010000000000000000", stored.SkimPercent.String(),
		"a percentage is not an amount and must not scale")

	_, found = pmk.GetTakerFeeShareAgreementFromDenomNoCache(ctx, legacyDenom)
	require.False(t, found, "nothing must be left under the retired denom")

	// The cache is what every swap actually reads.
	_, found = pmk.GetTakerFeeShareAgreementFromDenomUNSAFE(appparams.BaseCoinUnit)
	require.True(t, found, "the cache must resolve the live denom")
	_, found = pmk.GetTakerFeeShareAgreementFromDenomUNSAFE(legacyDenom)
	require.False(t, found, "and must no longer resolve the retired one")
}

// ---------------------------------------------------------------------------
// x/tokenization
// ---------------------------------------------------------------------------

// A second, per-approval denom allowlist sits beside the module params one, and
// only the params one was being migrated. It fails closed: a collection pinned
// to ["ubadge"] rejects every user-level coin transfer after the rename —
// including the ones this same migration just rescaled into the live denom.
func TestV35UserApprovalSettingsAllowedDenomsAreRepointed(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{{
			ApprovalId: "pinned",
			ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
				UserApprovalSettings: &tokenizationtypes.UserApprovalSettings{
					AllowedDenoms: []string{legacyDenom, foreignDenom},
				},
			},
		}},
	}
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, collection, true))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	got, found := app.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	allowed := got.CollectionApprovals[0].ApprovalCriteria.UserApprovalSettings.AllowedDenoms
	require.Equal(t, []string{appparams.BaseCoinUnit, foreignDenom}, allowed,
		"the live denom must be allowed and the foreign one must be left alone")
}

// The one that strands real money if it is wrong.
//
// A backed collection's escrow address is a module credential over the denom
// string and nothing else — not the collection id. So redenominating the denom
// re-derives a different address, and coins already escrowed do not follow:
// unbacking would query the new address, find it empty, and the backing would be
// permanently unredeemable. Mainnet collection 73 is in exactly this position.
func TestV35BackedPathEscrowMovesToTheRederivedAddress(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	oldEscrow, err := tokenizationkeeper.DerivePathAddress(legacyDenom, tokenizationkeeper.BackedPathGenerationPrefix)
	require.NoError(t, err)
	newEscrow, err := tokenizationkeeper.DerivePathAddress(appparams.BaseCoinUnit, tokenizationkeeper.BackedPathGenerationPrefix)
	require.NoError(t, err)
	require.False(t, oldEscrow.Equals(newEscrow),
		"precondition: the derivation must actually move, or this test proves nothing")

	// Mainnet collection 73's rate: 10^9 ubadge against 10^9 tokens.
	const sideAAmount = int64(1_000_000_000)
	// The coins actually escrowed behind it.
	const escrowed = int64(4_000_000_000)

	newBacked := func(id uint64) *tokenizationtypes.TokenCollection {
		return &tokenizationtypes.TokenCollection{
			CollectionId: sdkmath.NewUint(id),
			Invariants: &tokenizationtypes.CollectionInvariants{
				CosmosCoinBackedPath: &tokenizationtypes.CosmosCoinBackedPath{
					Address: oldEscrow.String(),
					Conversion: &tokenizationtypes.Conversion{
						SideA: &tokenizationtypes.ConversionSideAWithDenom{
							Amount: sdkmath.NewUint(uint64(sideAAmount)),
							Denom:  legacyDenom,
						},
					},
				},
			},
		}
	}
	// Two collections on the same denom, which mainnet also has for its IBC
	// denoms: they share one escrow, so the balance must be moved exactly once.
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, newBacked(1), true))
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, newBacked(2), true))

	fundLegacy(t, ctx, bk, oldEscrow, escrowed)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor

	for _, id := range []uint64{1, 2} {
		got, found := app.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(id))
		require.True(t, found)
		path := got.Invariants.CosmosCoinBackedPath

		require.Equal(t, appparams.BaseCoinUnit, path.Conversion.SideA.Denom)
		require.Equal(t, sdkmath.NewUint(uint64(sideAAmount)).Mul(sdkmath.NewUintFromBigInt(factor.BigInt())).String(),
			path.Conversion.SideA.Amount.String(),
			"collection %d: the backing rate must be unchanged in real terms", id)
		require.Equal(t, newEscrow.String(), path.Address,
			"collection %d: the path must point at the address its denom now derives", id)
	}

	require.Equal(t, sdkmath.NewInt(escrowed).Mul(factor).String(),
		bk.GetBalance(ctx, newEscrow, appparams.BaseCoinUnit).Amount.String(),
		"every escrowed coin must have followed the path to the re-derived address")
	require.True(t, bk.GetBalance(ctx, newEscrow, legacyDenom).IsZero())
	require.True(t, bk.GetBalance(ctx, oldEscrow, appparams.BaseCoinUnit).IsZero(),
		"and nothing may be left stranded at the address nothing will look at again")

	require.True(t, app.AccountKeeper.HasAccount(ctx, newEscrow),
		"the re-derived escrow needs an account record, as a send into it would have created")

	// No coins were invented or destroyed on the way.
	summed := sdkmath.ZeroInt()
	bk.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(coin.Amount)
		}
		return false
	})
	require.Equal(t, bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount, summed,
		"moving an escrow must conserve value")
}

// A collection backed by a foreign denom must be left entirely alone — its
// escrow derivation has not changed, so moving it would strand the coins this
// migration is trying to protect.
func TestV35BackedPathOnAForeignDenomIsUntouched(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"
	escrow, err := tokenizationkeeper.DerivePathAddress(foreignDenom, tokenizationkeeper.BackedPathGenerationPrefix)
	require.NoError(t, err)

	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		Invariants: &tokenizationtypes.CollectionInvariants{
			CosmosCoinBackedPath: &tokenizationtypes.CosmosCoinBackedPath{
				Address: escrow.String(),
				Conversion: &tokenizationtypes.Conversion{
					SideA: &tokenizationtypes.ConversionSideAWithDenom{
						Amount: sdkmath.NewUint(1_000_000),
						Denom:  foreignDenom,
					},
				},
			},
		},
	}, true))

	coins := sdk.NewCoins(sdk.NewCoin(foreignDenom, sdkmath.NewInt(5_000)))
	require.NoError(t, bk.MintCoins(ctx, "mint", coins))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", escrow, coins))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	got, found := app.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	path := got.Invariants.CosmosCoinBackedPath
	require.Equal(t, foreignDenom, path.Conversion.SideA.Denom)
	require.Equal(t, "1000000", path.Conversion.SideA.Amount.String())
	require.Equal(t, escrow.String(), path.Address)
	require.Equal(t, "5000", bk.GetBalance(ctx, escrow, foreignDenom).Amount.String())
}

// ---------------------------------------------------------------------------
// x/gamm and x/auth odds and ends
// ---------------------------------------------------------------------------

// The total-liquidity index is keyed by denom with an amount as its value, so
// the rename orphans it and the live denom reports zero liquidity.
func TestV35GammTotalLiquidityIndexIsRekeyedAndScaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"
	app.GammKeeper.RecordTotalLiquidityIncrease(ctx, sdk.NewCoins(
		sdk.NewCoin(legacyDenom, sdkmath.NewInt(1_000)),
		sdk.NewCoin(foreignDenom, sdkmath.NewInt(55)),
	))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	require.Equal(t, sdkmath.NewInt(1_000).Mul(v35.ConversionFactor).String(),
		app.GammKeeper.GetDenomLiquidity(ctx, appparams.BaseCoinUnit).String(),
		"recorded liquidity is an amount and must move with the denom")
	require.True(t, app.GammKeeper.GetDenomLiquidity(ctx, legacyDenom).IsZero(),
		"nothing must be left under the retired denom")
	require.Equal(t, "55", app.GammKeeper.GetDenomLiquidity(ctx, foreignDenom).String())

	total, err := app.GammKeeper.GetTotalLiquidity(ctx)
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewInt(1_000).Mul(v35.ConversionFactor).String(),
		total.AmountOf(appparams.BaseCoinUnit).String())
	require.True(t, total.AmountOf(legacyDenom).IsZero())
}

// RescaleVestingAccounts wrote back and counted every vesting account it saw,
// including ones with nothing to convert, so its Rescaled figure overstated what
// the migration did — the one number an operator would check the upgrade log
// against.
func TestV35VestingAccountWithNoLegacyCoinsIsNotCountedAsRescaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/ABCDEF0123456789"

	withBadge := randAddr()
	withoutBadge := randAddr()

	newVesting := func(addr sdk.AccAddress, coins sdk.Coins) {
		base := app.AccountKeeper.NewAccountWithAddress(ctx, addr).(*authtypes.BaseAccount)
		// An absolute future time: the test context's block time is the zero
		// Time, whose Unix() is negative and which the constructor rejects.
		acc, err := vestingtypes.NewDelayedVestingAccount(base, coins, 4_000_000_000)
		require.NoError(t, err)
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	newVesting(withBadge, sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1_000))))
	newVesting(withoutBadge, sdk.NewCoins(sdk.NewCoin(foreignDenom, sdkmath.NewInt(1_000))))

	res, err := v35.RescaleVestingAccounts(ctx, app.AccountKeeper)
	require.NoError(t, err)
	require.Equal(t, 1, res.Rescaled,
		"only the account that actually held the retired denom was rescaled")
	require.GreaterOrEqual(t, res.Scanned, 2)

	// And the account that was skipped must be genuinely unchanged.
	acc := app.AccountKeeper.GetAccount(ctx, withoutBadge).(*vestingtypes.DelayedVestingAccount)
	require.Equal(t, "1000", acc.OriginalVesting.AmountOf(foreignDenom).String(),
		"a vesting account with nothing to convert must come out byte-identical")
}

// ---------------------------------------------------------------------------
// x/gov: in-flight proposal payloads
// ---------------------------------------------------------------------------

// A proposal still in its deposit or voting period at the upgrade height carries
// an Any-packed payload written against the retired denom. Mainnet's four
// ubadge-bearing proposals are all MsgUpdateParams, and one of those passing
// after the upgrade would write the retired denom straight back into the params
// this migration just moved.
//
// The handler must not rewrite it — a governance payload is what voters approved
// — and must not reject it either, or filing one would be a free chain-halting
// DoS for the price of a minimum deposit. It flags it, and carries on.
func TestV35UndecidedProposalNamingTheRetiredDenomDoesNotHaltTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	govParams, err := app.GovKeeper.Params.Get(ctx)
	require.NoError(t, err)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1)))
	require.NoError(t, app.GovKeeper.Params.Set(ctx, govParams))

	proposer := randAddr()
	deposit := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(10)))
	require.NoError(t, bk.MintCoins(ctx, "mint", deposit))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", proposer, deposit))

	// The payload names the retired denom, the way mainnet's ubadge-bearing
	// proposals do.
	payload := banktypes.NewMsgSend(
		authtypes.NewModuleAddress("gov"), randAddr(),
		sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1_000))),
	)
	proposal, err := app.GovKeeper.SubmitProposal(ctx,
		[]sdk.Msg{payload}, "", "flagged", "flagged", proposer, false)
	require.NoError(t, err)
	_, err = app.GovKeeper.AddDeposit(ctx, proposal.Id, proposer, deposit)
	require.NoError(t, err)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)),
		"an undecided proposal naming the retired denom must not halt the upgrade")

	stored, err := app.GovKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(t, err)
	require.Len(t, stored.Messages, 1)

	// The payload is left byte-identical: the handler reports it, it does not
	// edit what voters were shown.
	var migrated banktypes.MsgSend
	require.NoError(t, app.AppCodec().Unmarshal(stored.Messages[0].Value, &migrated))
	require.Equal(t, legacyDenom, migrated.Amount[0].Denom,
		"a governance payload must not be rewritten underneath its voters")
	require.Equal(t, "1000", migrated.Amount[0].Amount.String())

	// The deposit, which is gov's own accounting rather than the voters', does
	// move.
	require.Equal(t, sdkmath.NewInt(10).Mul(v35.ConversionFactor).String(),
		sdk.Coins(stored.TotalDeposit).AmountOf(appparams.BaseCoinUnit).String())
}

// Flagging is the entire mitigation for the case above, so it has to be asserted
// directly: a "did not halt" test passes just as well against a handler that
// never looked.
func TestV35FlagsUndecidedProposalsAndOnlyUndecidedOnes(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	seedLegacyBondDenom(t, app, ctx)

	govParams, err := app.GovKeeper.Params.Get(ctx)
	require.NoError(t, err)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1)))
	require.NoError(t, app.GovKeeper.Params.Set(ctx, govParams))

	proposer := randAddr()
	funds := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(100)))
	require.NoError(t, bk.MintCoins(ctx, "mint", funds))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", proposer, funds))

	submit := func(denom string) uint64 {
		msg := banktypes.NewMsgSend(authtypes.NewModuleAddress("gov"), randAddr(),
			sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(1_000))))
		p, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msg}, "", "p", "p", proposer, false)
		require.NoError(t, err)
		return p.Id
	}

	flagged := submit(legacyDenom)
	clean := submit("ibc/ABCDEF0123456789")

	// The same payload on a proposal that is already decided must NOT be
	// flagged: it will never execute again, so it cannot affect anything.
	decided := submit(legacyDenom)
	decidedProposal, err := app.GovKeeper.Proposals.Get(ctx, decided)
	require.NoError(t, err)
	decidedProposal.Status = govv1.StatusPassed
	require.NoError(t, app.GovKeeper.SetProposal(ctx, decidedProposal))

	res, err := v35.RescaleGovDeposits(ctx, app.GovKeeper)
	require.NoError(t, err)
	require.Equal(t, 3, res.Proposals)
	require.Equal(t, 1, res.FlaggedMessages,
		"exactly the undecided proposal naming the retired denom must be flagged")

	_ = flagged
	_ = clean
}
