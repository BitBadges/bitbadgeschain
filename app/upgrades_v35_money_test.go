package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	seedLegacyBondDenom(t, app, ctx)

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
	seedLegacyBondDenom(t, app, ctx)

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
	seedLegacyBondDenom(t, app, ctx)

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
	seedLegacyBondDenom(t, app, ctx)

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

// seedPreciseBankFractionalBalances writes the retired module's state directly:
// per-account fractional remainders, an unowned remainder, and the ubadge
// reserve on the module account that backs both.
//
// Written raw rather than through the module's keeper because the module is
// unwired — which is exactly the situation the migration has to cope with.
func seedPreciseBankFractionalBalances(
	t *testing.T,
	app *App,
	ctx sdk.Context,
	bk bankkeeper.BaseKeeper,
	fractional map[string]int64,
	remainder int64,
) {
	t.Helper()

	store := ctx.KVStore(app.LegacyPreciseBankKey)

	total := int64(0)
	for addrStr, amount := range fractional {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		require.NoError(t, err)
		if amount == 0 {
			// x/precisebank deletes rather than stores a zero, so a zero
			// fractional balance is simply absent. Asserting on it still
			// matters: the account must come out with its integer part intact.
			continue
		}
		bz, err := sdkmath.NewInt(amount).Marshal()
		require.NoError(t, err)
		store.Set(append(append([]byte{}, v35.PreciseBankFractionalPrefix...), addr.Bytes()...), bz)
		total += amount
	}

	if remainder > 0 {
		bz, err := sdkmath.NewInt(remainder).Marshal()
		require.NoError(t, err)
		store.Set(v35.PreciseBankRemainderKey, bz)
		total += remainder
	}

	// The reserve backs every fractional unit 1:1 in ubadge, so at genesis-time
	// validation sum(fractional) + remainder == reserve * 10^9.
	require.Zero(t, total%1_000_000_000, "seed must leave a whole number of ubadge in the reserve")
	fundPreciseBankReserve(t, ctx, bk, total/1_000_000_000)
}

// fundPreciseBankReserve puts ubadge on the retired module's reserve address.
//
// Deliberately not SendCoinsFromModuleToAccount: v35 adds "precisebank" to the
// bank's blocked-address list, so a send there is refused. That block is the
// point — nothing should be able to put coins on an address whose module is
// gone — and the seeding here has to reach around it the same way the real
// pre-upgrade chain got there, which was through x/precisebank's own keeper
// while the module still existed.
func fundPreciseBankReserve(t *testing.T, ctx sdk.Context, bk bankkeeper.BaseKeeper, ubadge int64) {
	t.Helper()
	if ubadge == 0 {
		return
	}
	reserveAddr := authtypes.NewModuleAddress(v35.PreciseBankStoreKey)
	mintAddr := authtypes.NewModuleAddress("mint")

	require.NoError(t, bk.MintCoins(ctx, "mint",
		sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(ubadge)))))

	// Hand the freshly minted coins over without a send. Supply is already
	// right; only the two balances move, so they must move together or the
	// balances stop summing to supply before the migration even starts.
	require.NoError(t, bk.UncheckedSetBalance(ctx, mintAddr,
		sdk.NewCoin(legacyDenom, bk.GetBalance(ctx, mintAddr, legacyDenom).Amount.SubRaw(ubadge))))
	require.NoError(t, bk.UncheckedSetBalance(ctx, reserveAddr,
		sdk.NewCoin(legacyDenom, bk.GetBalance(ctx, reserveAddr, legacyDenom).Amount.AddRaw(ubadge))))
}

// At 9 decimals an EVM account's balance was ubadge*10^9 + a fractional
// remainder held by x/precisebank and backed by ubadge on its reserve. Deleting
// that store as part of the upgrade destroys the fractional part of every
// account that ever received a non-integer-ubadge transfer, and strands the
// reserve on a module address whose module no longer exists.
//
// After the upgrade a fractional unit is one abadge, so the post-upgrade
// balance must be exactly ubadge*10^9 + fractional.
func TestV35PreciseBankFractionalBalancesSurviveTheUpgrade(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	// maxFractional is the largest a fractional balance can be: one below the
	// 10^9 that would have carried into an integer ubadge.
	const maxFractional = int64(999_999_999)

	full := randAddr()   // integer balance plus the maximum fractional part
	none := randAddr()   // integer balance, no fractional part at all
	dustOnly := randAddr() // fractional part only, no integer balance

	fundLegacy(t, ctx, bk, full, 5)
	fundLegacy(t, ctx, bk, none, 7)

	// The three fractional parts plus the remainder must come to a whole
	// number of ubadge, which is what the reserve holds.
	const dust = int64(3)
	const remainder = int64(999_999_998) // maxFractional + dust + remainder == 2 ubadge
	seedPreciseBankFractionalBalances(t, app, ctx, bk, map[string]int64{
		full.String():     maxFractional,
		none.String():     0,
		dustOnly.String(): dust,
	}, remainder)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor
	require.Equal(t, sdkmath.NewInt(5).Mul(factor).AddRaw(maxFractional),
		bk.GetBalance(ctx, full, appparams.BaseCoinUnit).Amount,
		"an account's fractional remainder is real value and must survive")
	require.Equal(t, sdkmath.NewInt(7).Mul(factor),
		bk.GetBalance(ctx, none, appparams.BaseCoinUnit).Amount,
		"an account with no fractional balance must be unaffected")
	require.Equal(t, sdkmath.NewInt(dust),
		bk.GetBalance(ctx, dustOnly, appparams.BaseCoinUnit).Amount,
		"an account whose whole balance was fractional must still have it")

	reserveAddr := authtypes.NewModuleAddress(v35.PreciseBankStoreKey)
	require.True(t, bk.GetBalance(ctx, reserveAddr, appparams.BaseCoinUnit).IsZero(),
		"the reserve's module no longer exists; anything left there is stranded")
	require.True(t, bk.GetBalance(ctx, reserveAddr, legacyDenom).IsZero())

	// Supply must still equal the sum of balances: the remainder was owned by
	// nobody, so it is burned rather than left inflating supply.
	summed := sdkmath.ZeroInt()
	bk.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(coin.Amount)
		}
		return false
	})
	require.Equal(t, bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount, summed,
		"balances must sum to supply after paying out the fractional balances")
}

// A reserve that does not back what the fractional balances claim means the
// module's own invariant has been broken. Paying out anyway would mint value
// out of nothing, so the upgrade must refuse rather than guess.
func TestV35PreciseBankRefusesWhenTheReserveDoesNotBackTheBalances(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyBondDenom(t, app, ctx)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 1)

	// One whole ubadge of fractional balance, with no reserve behind it.
	store := ctx.KVStore(app.LegacyPreciseBankKey)
	bz, err := sdkmath.NewInt(500_000_000).Marshal()
	require.NoError(t, err)
	store.Set(append(append([]byte{}, v35.PreciseBankFractionalPrefix...), alice.Bytes()...), bz)

	err = v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app))
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not back it")
}
