package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	stableswap "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/stableswap"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// fundLegacyModule mints legacy coins onto a module account without clobbering
// its module-account record, which SendCoinsFromModuleToAccount would do.
func fundLegacyModule(t *testing.T, ctx sdk.Context, app *App, bk bankkeeper.BaseKeeper, module string, amount int64) {
	t.Helper()
	require.NotNil(t, app.AccountKeeper.GetModuleAccount(ctx, module))
	coins := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(amount)))
	require.NoError(t, bk.MintCoins(ctx, "mint", coins))
	require.NoError(t, bk.SendCoinsFromModuleToModule(ctx, "mint", module, coins))
}

// seedDelegationsWithPendingRewards builds the state a real chain has when
// delegators have unclaimed staking rewards: a validator with distribution
// records initialised through the module's own hooks, delegations splitting the
// validator's shares evenly, and an allocation of reward tokens backed by coins
// on the distribution module account.
//
// It goes through the keeper's exported hooks rather than hand-writing the
// distribution store so the reward accounting is the same shape the chain
// produces, not a shape the test invented.
func seedDelegationsWithPendingRewards(
	t *testing.T,
	app *App,
	ctx sdk.Context,
	rewardAmount int64,
	numDelegators int64,
) (sdk.Context, stakingtypes.Validator, []sdk.AccAddress) {
	t.Helper()

	val := seedLegacyChainState(t, app, ctx)
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	require.NoError(t, err)

	require.NoError(t, app.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valAddr))

	sharesEach := val.DelegatorShares.QuoInt64(numDelegators)
	delAddrs := make([]sdk.AccAddress, 0, numDelegators)
	for i := int64(0); i < numDelegators; i++ {
		delAddr := randAddr()
		del := stakingtypes.NewDelegation(delAddr.String(), val.OperatorAddress, sharesEach)
		require.NoError(t, app.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr))
		require.NoError(t, app.StakingKeeper.SetDelegation(ctx, del))
		require.NoError(t, app.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr))
		delAddrs = append(delAddrs, delAddr)
	}

	// Rewards accrue at a later height than the delegations were created;
	// CalculateDelegationRewards short-circuits to zero when they match.
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	fundLegacyModule(t, ctx, app, bk, distrtypes.ModuleName, rewardAmount)

	reward := sdk.NewDecCoins(sdk.NewDecCoin(legacyDenom, sdkmath.NewInt(rewardAmount)))
	require.NoError(t, app.DistrKeeper.AllocateTokensToValidator(ctx, val, reward))

	return ctx, val, delAddrs
}

// pendingReward computes what the delegator could withdraw right now, through
// the real reward calculation rather than by reading the store.
func pendingReward(t *testing.T, app *App, ctx sdk.Context, valAddr sdk.ValAddress, delAddr sdk.AccAddress) sdk.DecCoins {
	t.Helper()

	val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
	require.NoError(t, err)
	del, err := app.StakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	require.NoError(t, err)

	endingPeriod, err := app.DistrKeeper.IncrementValidatorPeriod(ctx, val)
	require.NoError(t, err)

	rewards, err := app.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	require.NoError(t, err)
	return rewards
}

// TestV35PendingStakingRewardsScaleByExactlyTheConversionFactor is the test the
// whole distribution migration turns on.
//
// A delegator's payout is (cumulative reward ratio delta) * stake.
// CumulativeRewardRatio is coins *per token*; DelegatorStartingInfo.Stake is a
// token amount. A redenomination multiplies coins and tokens by the same
// factor, so the ratio is invariant and only the stake moves. Scaling both
// squares the factor: a 10^9 redenomination turns into a 10^18 one, and the
// pending reward exceeds the outstanding rewards backing it.
//
// Asserting on the migration's own bookkeeping would not catch that. This runs
// the SDK's own CalculateDelegationRewards on both sides of the upgrade.
func TestV35PendingStakingRewardsScaleByExactlyTheConversionFactor(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	const rewardAmount = 1_000 // ubadge
	ctx, val, delAddrs := seedDelegationsWithPendingRewards(t, app, ctx, rewardAmount, 1)
	delAddr := delAddrs[0]
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	require.NoError(t, err)

	before := pendingReward(t, app, ctx, valAddr, delAddr)
	require.Equal(t,
		sdk.NewDecCoins(sdk.NewDecCoin(legacyDenom, sdkmath.NewInt(rewardAmount))).String(),
		before.String(),
		"seed must produce exactly the allocated reward as pending")

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	after := pendingReward(t, app, ctx, valAddr, delAddr)

	want := sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(rewardAmount).Mul(v35.ConversionFactor))
	require.Equal(t, want.String(), after.AmountOf(appparams.BaseCoinUnit).String(),
		"pending rewards must scale by exactly 10^9, not 10^18")
	require.True(t, after.AmountOf(legacyDenom).IsZero(),
		"no reward may remain denominated in the retired denom")
}

// An over-scaled pending reward does not surface as an error: cosmos-sdk
// v0.54.4 clamps the payout to the validator's outstanding rewards
// (rewardsRaw.Intersect(outstanding)) and only logs a "rounding error". With
// two delegators that clamp is theft — the first claimant drains the whole
// outstanding pot and the second is left with nothing, silently.
//
// So the assertion is not "the withdrawal succeeds" but "each delegator gets
// their own share, redenominated".
func TestV35WithdrawDelegatorRewardPaysEachDelegatorTheirShare(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	const rewardAmount = 1_000
	ctx, val, delAddrs := seedDelegationsWithPendingRewards(t, app, ctx, rewardAmount, 2)
	valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	require.NoError(t, err)

	// Materialise the accrued rewards into a historical cumulative ratio, which
	// is what any delegation change, slash or withdrawal does on a live chain.
	// Without this the ratio is still zero at the upgrade height and scaling it
	// wrongly has nothing to act on.
	require.False(t, pendingReward(t, app, ctx, valAddr, delAddrs[0]).IsZero())

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	msgServer := distrkeeper.NewMsgServerImpl(app.DistrKeeper)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)
	want := sdkmath.NewInt(rewardAmount / 2).Mul(v35.ConversionFactor)

	for i, delAddr := range delAddrs {
		_, err := msgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: val.OperatorAddress,
		})
		require.NoError(t, err, "withdrawal %d after the upgrade must not fail", i)

		got := bk.GetBalance(ctx, delAddr, appparams.BaseCoinUnit).Amount
		require.Equal(t, want, got,
			"delegator %d must receive exactly its own redenominated share, not the whole pot", i)
	}
}

// A stableswap pool prices swaps off liquidity[i] / scalingFactors[i]. Scaling
// the BADGE liquidity by 10^9 and leaving its scaling factor alone makes the
// BADGE leg look 10^9 times deeper than it is, so the pool quotes BADGE at
// roughly zero and the first swapper drains the other side.
//
// The scaling factors are also index-aligned with PoolLiquidity, and the rename
// moves BADGE's position in the sorted coin list — so the factors have to be
// permuted with it, not just scaled.
func TestV35StableswapSpotPriceIsUnchangedAcrossTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// uatom sorts after abadge and before ubadge, so the rename reorders
	// PoolLiquidity and any index-aligned array has to move with it.
	const otherDenom = "uatom"

	liquidity := sdk.NewCoins(
		sdk.NewCoin(legacyDenom, sdkmath.NewInt(2_000_000_000)),
		sdk.NewCoin(otherDenom, sdkmath.NewInt(3_000_000_000)),
	)
	require.Equal(t, otherDenom, liquidity[0].Denom, "seed assumes uatom sorts before ubadge")

	// Non-1 scaling factors, index-aligned with the sorted liquidity above.
	scalingFactors := []uint64{7, 100}

	pool, err := stableswap.NewStableswapPool(
		1,
		stableswap.PoolParams{SwapFee: osmomath.ZeroDec(), ExitFee: osmomath.ZeroDec()},
		liquidity,
		scalingFactors,
		"",
		"",
	)
	require.NoError(t, err)
	require.NoError(t, app.GammKeeper.OverwritePoolV15MigrationUnsafe(ctx, &pool))

	// Quote BADGE against the other asset. This direction multiplies by the
	// BADGE scaling factor at the end, so the post-upgrade figure is the
	// pre-upgrade one times 10^9 exactly — the unit shrank by 10^9, the price
	// per unit did not move.
	before, err := pool.SpotPrice(ctx, legacyDenom, otherDenom)
	require.NoError(t, err)
	require.False(t, before.IsZero())

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	migrated, err := app.GammKeeper.GetPool(ctx, 1)
	require.NoError(t, err)
	after, ok := migrated.(*stableswap.Pool)
	require.True(t, ok)

	// The factors must still line up with the liquidity they scale.
	require.Equal(t, len(after.PoolLiquidity), len(after.ScalingFactors))
	require.Equal(t, appparams.BaseCoinUnit, after.PoolLiquidity[0].Denom,
		"abadge sorts first, so the liquidity must have been re-sorted")
	require.Equal(t, otherDenom, after.PoolLiquidity[1].Denom)

	// The scaled reserve the AMM actually prices off — liquidity/factor — must
	// be the same rational number on both sides of the upgrade.
	require.Equal(t,
		liquidity.AmountOf(legacyDenom).Mul(sdkmath.NewIntFromUint64(after.ScalingFactors[0])),
		after.PoolLiquidity[0].Amount.Mul(sdkmath.NewIntFromUint64(scalingFactors[1])),
		"the BADGE leg's scaled reserve must not move")
	require.Equal(t, liquidity.AmountOf(otherDenom), after.PoolLiquidity[1].Amount,
		"the non-BADGE leg must not move")
	require.Equal(t, scalingFactors[0], after.ScalingFactors[1],
		"the non-BADGE scaling factor must not move")

	// The unit shrank by 10^9, so the price per unit grows by 10^9 and nothing
	// else. Compared at the pre-upgrade decimal precision, because the larger
	// descaling multiplier lets the post-upgrade figure carry more significant
	// digits than the pre-upgrade one could represent.
	afterPrice, err := after.SpotPrice(ctx, appparams.BaseCoinUnit, otherDenom)
	require.NoError(t, err)
	require.Equal(t, before.Dec().String(), afterPrice.QuoInt64(1_000_000_000).Dec().String(),
		"a redenomination must not move the price of a stableswap pool")
}

// x/tokenization approvals also live on TokenCollection.DefaultBalances, which
// every user inherits the first time they touch a collection. Skipping it makes
// every paid mint priced through a default approval cost 10^-9 of its intent —
// which is to say, free — for anyone who joins after the upgrade.
func TestV35DefaultBalancesCoinTransfersAreScaled(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	price := sdk.NewCoin(legacyDenom, sdkmath.NewInt(1_000_000_000)) // 1 BADGE
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		DefaultBalances: &tokenizationtypes.UserBalanceStore{
			IncomingApprovals: []*tokenizationtypes.UserIncomingApproval{{
				ApprovalId: "default-incoming",
				ApprovalCriteria: &tokenizationtypes.IncomingApprovalCriteria{
					CoinTransfers: []*tokenizationtypes.CoinTransfer{{
						To:    randAddr().String(),
						Coins: []*sdk.Coin{&price},
					}},
				},
			}},
		},
	}
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, collection, true))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	got, found := app.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	coin := got.DefaultBalances.IncomingApprovals[0].ApprovalCriteria.CoinTransfers[0].Coins[0]
	require.Equal(t, appparams.BaseCoinUnit, coin.Denom,
		"a default approval priced in the retired denom is priced in nothing")
	require.Equal(t, sdkmath.NewInt(1_000_000_000).Mul(v35.ConversionFactor), coin.Amount,
		"a paid mint must cost the same after the upgrade as before it")
}

// The rename moves BADGE's sort position, so a multi-denom coin transfer that
// was sorted before the migration is not sorted after it unless the slice is
// rebuilt. Everywhere else in this migration re-sorts (convertCoins,
// scaleDecCoins, the gamm pool assets); this path mutated in place.
func TestV35CoinTransferCoinsStaySortedByDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// ibc/... sorts after abadge and before ubadge.
	const ibcDenom = "ibc/ABCDEF0123456789"
	badge := sdk.NewCoin(legacyDenom, sdkmath.NewInt(5))
	voucher := sdk.NewCoin(ibcDenom, sdkmath.NewInt(9))

	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{{
			ApprovalId: "paid",
			ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
				CoinTransfers: []*tokenizationtypes.CoinTransfer{{
					To:    randAddr().String(),
					Coins: []*sdk.Coin{&voucher, &badge}, // sorted: ibc/... < ubadge
				}},
			},
		}},
	}
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, collection, true))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	got, found := app.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	coins := got.CollectionApprovals[0].ApprovalCriteria.CoinTransfers[0].Coins

	flat := sdk.Coins{}
	for _, c := range coins {
		flat = append(flat, *c)
	}
	require.NoError(t, flat.Validate(),
		"the migrated coins must still be a valid, sorted sdk.Coins")
	require.Equal(t, sdkmath.NewInt(5).Mul(v35.ConversionFactor), flat.AmountOf(appparams.BaseCoinUnit),
		"AmountOf binary-searches a sorted slice; an unsorted one silently returns zero")
	require.Equal(t, sdkmath.NewInt(9), flat.AmountOf(ibcDenom),
		"a foreign denom must not move")
}

// A stored taker-fee override that happens to equal the current default is a
// reachable configuration: governance can change the default after an override
// was written. Erroring there aborts the whole upgrade handler and halts the
// chain at the upgrade height with no recovery short of a new binary.
func TestV35TakerFeeOverrideEqualToDefaultDoesNotHaltTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	defaultFee := app.PoolManagerKeeper.GetDefaultTakerFee(ctx)
	app.PoolManagerKeeper.SetDenomPairTakerFee(ctx, legacyDenom, "uatom", defaultFee.Add(osmomath.NewDecWithPrec(1, 3)))

	// Now move the default onto the override's value, which is exactly the
	// state a governance parameter change leaves behind.
	pmParams := app.PoolManagerKeeper.GetParams(ctx)
	pmParams.TakerFeeParams.DefaultTakerFee = defaultFee.Add(osmomath.NewDecWithPrec(1, 3))
	app.PoolManagerKeeper.SetParams(ctx, pmParams)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)),
		"a redundant taker-fee override must not halt the chain")
}
