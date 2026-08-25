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
