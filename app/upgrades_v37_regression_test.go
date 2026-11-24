package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v37 "github.com/bitbadges/bitbadgeschain/app/upgrades/v37"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// TestV37RedenominateBankPreservesPreExistingAbadgeBalances verifies that if
// any account (including the staging evm module account) already holds abadge
// prior to the upgrade, RedenominateBank accumulates converted ubadge onto
// those pre-existing balances rather than overwriting them, preserving total
// supply conservation exactly.
func TestV37RedenominateBankPreservesPreExistingAbadgeBalances(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	alice := randAddr()
	bob := randAddr()
	evmModAddr := authtypes.NewModuleAddress(evmtypes.ModuleName)

	// Mint pre-existing abadge on alice, bob, and the evm module account.
	alicePreAbadge := sdkmath.NewInt(123_456)
	bobPreAbadge := sdkmath.NewInt(777)
	modPreAbadge := sdkmath.NewInt(999)
	mintTo := func(addr sdk.AccAddress, coin sdk.Coin) {
		coins := sdk.NewCoins(coin)
		require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", coins))
		if addr.Equals(evmModAddr) {
			require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", evmtypes.ModuleName, coins))
			return
		}
		require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, coins))
	}

	mintTo(alice, sdk.NewCoin(appparams.BaseCoinUnit, alicePreAbadge))
	mintTo(bob, sdk.NewCoin(appparams.BaseCoinUnit, bobPreAbadge))
	mintTo(evmModAddr, sdk.NewCoin(appparams.BaseCoinUnit, modPreAbadge))

	// Mint legacy ubadge on alice and the evm module account (bob holds only abadge).
	aliceLegacy := sdkmath.NewInt(50)
	modLegacy := sdkmath.NewInt(10)
	mintTo(alice, sdk.NewCoin(appparams.LegacyBaseCoinUnit, aliceLegacy))
	mintTo(evmModAddr, sdk.NewCoin(appparams.LegacyBaseCoinUnit, modLegacy))

	preNewSupply := app.BankKeeper.GetSupply(ctx, appparams.BaseCoinUnit).Amount
	legacySupply := app.BankKeeper.GetSupply(ctx, appparams.LegacyBaseCoinUnit).Amount
	require.Equal(t, aliceLegacy.Add(modLegacy), legacySupply)

	res, err := v37.RedenominateBank(ctx, app.V37Keepers().Bank)
	require.NoError(t, err)
	require.Equal(t, legacySupply, res.LegacySupply)
	require.Equal(t, legacySupply.Mul(v37.ConversionFactor), res.NewSupply)

	// Legacy supply must be completely retired.
	require.True(t, app.BankKeeper.GetSupply(ctx, appparams.LegacyBaseCoinUnit).Amount.IsZero())

	// Every account retains pre-existing abadge PLUS converted legacy balance.
	wantAlice := alicePreAbadge.Add(aliceLegacy.Mul(v37.ConversionFactor))
	wantMod := modPreAbadge.Add(modLegacy.Mul(v37.ConversionFactor))
	require.Equal(t, wantAlice.String(), app.BankKeeper.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, bobPreAbadge.String(), app.BankKeeper.GetBalance(ctx, bob, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, wantMod.String(), app.BankKeeper.GetBalance(ctx, evmModAddr, appparams.BaseCoinUnit).Amount.String())

	// Sum of all abadge balances must match BankKeeper total supply exactly.
	wantSupply := preNewSupply.Add(legacySupply.Mul(v37.ConversionFactor))
	gotSupply := app.BankKeeper.GetSupply(ctx, appparams.BaseCoinUnit).Amount
	require.Equal(t, wantSupply.String(), gotSupply.String())

	summed := sdkmath.ZeroInt()
	app.BankKeeper.IterateAllBalances(ctx, func(_ sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(c.Amount)
		}
		return false
	})
	require.Equal(t, gotSupply.String(), summed.String(), "sum of all abadge balances must equal total supply")
}

// TestV37MixedDenomDistributionRewardsConserved verifies that a validator with
// multiple delegators holding BOTH native (ubadge) and foreign (ibc/USDC)
// distribution rewards settles pre-upgrade foreign rewards losslessly at full
// 9-decimal precision AND rescales historical foreign ratios by 10^-9 so post-
// upgrade withdrawals neither inflate by 10^9 nor clamp-drain other delegators.
func TestV37MixedDenomDistributionRewardsConserved(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Seed pre-upgrade 9-decimal bond denom and staking/distribution state.
	seedLegacyBondDenom(t, app, ctx)

	const foreignDenom = "ibc/E1116484B327AEE59CDC3DA73D319834781A13DB2A7DFC1F38A30CD45ABF58B8"

	valAddr := sdk.ValAddress(randAddr())
	del1 := randAddr()
	del2 := randAddr()

	// Validator with 1000 ubadge bonded (del1 = 250, del2 = 750).
	valTokens := sdkmath.NewInt(1_000)
	del1Stake := sdkmath.NewInt(250)
	del2Stake := sdkmath.NewInt(750)

	validator := stakingtypes.Validator{
		OperatorAddress: valAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          valTokens,
		DelegatorShares: sdkmath.LegacyNewDecFromInt(valTokens),
		Commission:      stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyOneDec(), sdkmath.LegacyOneDec()),
	}
	require.NoError(t, app.StakingKeeper.SetValidator(ctx, validator))
	require.NoError(t, app.StakingKeeper.SetDelegation(ctx, stakingtypes.NewDelegation(del1.String(), valAddr.String(), sdkmath.LegacyNewDecFromInt(del1Stake))))
	require.NoError(t, app.StakingKeeper.SetDelegation(ctx, stakingtypes.NewDelegation(del2.String(), valAddr.String(), sdkmath.LegacyNewDecFromInt(del2Stake))))

	// Fund bonded pool for staking consistency.
	bondedCoins := sdk.NewCoins(sdk.NewCoin(appparams.LegacyBaseCoinUnit, valTokens))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", bondedCoins))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", stakingtypes.BondedPoolName, bondedCoins))

	// Initialize distribution state for validator across period 1 -> period 2:
	// Period 1 baseline ratio = 0, Period 2 accrues 400 ubadge + 400 foreignDenom.
	// Expected payouts: del1 gets 100 ubadge + 100 foreign, del2 gets 300 ubadge + 300 foreign.
	nativeReward := sdkmath.NewInt(400)
	foreignReward := sdkmath.NewInt(400)
	distrCoins := sdk.NewCoins(
		sdk.NewCoin(appparams.LegacyBaseCoinUnit, nativeReward),
		sdk.NewCoin(foreignDenom, foreignReward),
	)
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", distrCoins))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", distrtypes.ModuleName, distrCoins))

	require.NoError(t, app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 1, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 2)))
	ratio := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(appparams.LegacyBaseCoinUnit, sdkmath.LegacyNewDec(400).QuoInt(valTokens)),
		sdk.NewDecCoinFromDec(foreignDenom, sdkmath.LegacyNewDec(400).QuoInt(valTokens)),
	)
	require.NoError(t, app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 2, distrtypes.NewValidatorHistoricalRewards(ratio, 1)))
	require.NoError(t, app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 3)))
	require.NoError(t, app.DistrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, distrtypes.ValidatorOutstandingRewards{
		Rewards: sdk.NewDecCoinsFromCoins(distrCoins...),
	}))
	require.NoError(t, app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, valAddr, distrtypes.InitialValidatorAccumulatedCommission()))
	require.NoError(t, app.DistrKeeper.SetDelegatorStartingInfo(ctx, valAddr, del1, distrtypes.NewDelegatorStartingInfo(1, sdkmath.LegacyNewDecFromInt(del1Stake), 1)))
	require.NoError(t, app.DistrKeeper.SetDelegatorStartingInfo(ctx, valAddr, del2, distrtypes.NewDelegatorStartingInfo(1, sdkmath.LegacyNewDecFromInt(del2Stake), 1)))

	// Run the full v37 upgrade handler.
	require.NoError(t, v37.CustomUpgradeHandlerLogic(ctx, app.V37Keepers()))

	// Verify pre-upgrade foreign rewards were settled losslessly AND pre-upgrade
	// native rewards were paid out & converted x10^9 into abadge.
	require.Equal(t, "100", app.BankKeeper.GetBalance(ctx, del1, foreignDenom).Amount.String())
	require.Equal(t, "300", app.BankKeeper.GetBalance(ctx, del2, foreignDenom).Amount.String())
	require.Equal(t, sdkmath.NewInt(100).Mul(v37.ConversionFactor).String(), app.BankKeeper.GetBalance(ctx, del1, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, sdkmath.NewInt(300).Mul(v37.ConversionFactor).String(), app.BankKeeper.GetBalance(ctx, del2, appparams.BaseCoinUnit).Amount.String())

	// Verify any remaining historical foreign ratio was divided by 10^9 (not left at old scale).
	app.DistrKeeper.IterateValidatorHistoricalRewards(ctx, func(v sdk.ValAddress, _ uint64, hr distrtypes.ValidatorHistoricalRewards) bool {
		if !v.Equals(valAddr) {
			return false
		}
		for _, c := range hr.CumulativeRewardRatio {
			require.NotEqual(t, appparams.LegacyBaseCoinUnit, c.Denom)
			if c.Denom == foreignDenom {
				// 0.4 / 10^9 = 4e-10
				expected := sdkmath.LegacyNewDec(400).QuoInt(valTokens).Quo(sdkmath.LegacyNewDecFromInt(v37.ConversionFactor))
				require.True(t, c.Amount.Equal(expected), "foreign ratio must scale by 10^-9: got %s, want %s", c.Amount, expected)
			}
		}
		return false
	})

	// Now accrue NEW post-upgrade mixed rewards (800 foreign + 800e9 abadge) and
	// verify both delegators withdraw their exact 25% / 75% shares with zero clamp!
	postNative := sdkmath.NewInt(800).Mul(v37.ConversionFactor)
	postForeign := sdkmath.NewInt(800)
	postCoins := sdk.NewCoins(
		sdk.NewCoin(appparams.BaseCoinUnit, postNative),
		sdk.NewCoin(foreignDenom, postForeign),
	)
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", postCoins))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, "mint", distrtypes.ModuleName, postCoins))

	cur, err := app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr)
	require.NoError(t, err)
	cur.Rewards = cur.Rewards.Add(sdk.NewDecCoinsFromCoins(postCoins...)...)
	require.NoError(t, app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddr, cur))

	out, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddr)
	require.NoError(t, err)
	out.Rewards = out.Rewards.Add(sdk.NewDecCoinsFromCoins(postCoins...)...)
	require.NoError(t, app.DistrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, out))

	// Advance block height so x/distribution allows post-upgrade reward
	// withdrawal (startingInfo.Height == currentHeight returns zero).
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	w1, err := app.DistrKeeper.WithdrawDelegationRewards(ctx, del1, valAddr)
	require.NoError(t, err)
	w2, err := app.DistrKeeper.WithdrawDelegationRewards(ctx, del2, valAddr)
	require.NoError(t, err)

	require.Equal(t, "200", w1.AmountOf(foreignDenom).String())
	require.Equal(t, "600", w2.AmountOf(foreignDenom).String())
	require.Equal(t, sdkmath.NewInt(200).Mul(v37.ConversionFactor).String(), w1.AmountOf(appparams.BaseCoinUnit).String())
	require.Equal(t, sdkmath.NewInt(600).Mul(v37.ConversionFactor).String(), w2.AmountOf(appparams.BaseCoinUnit).String())
}

// TestV37MigratesSDKv055KeyRotationFeePoolManagerKVParamsAndEVMCoinInfo verifies
// all v36/SDK v0.55 state surfaces that did not exist in the earlier v35 draft:
//   - x/staking Params.KeyRotationFee (denom + x10^9 amount)
//   - x/poolmanager KV-stored Params (CommunityPoolDenomToSwapNonWhitelistedAssetsTo + Whitelist)
//   - x/vm persisted EvmCoinInfo + process-global evmCoinInfo (abadge, 18 decimals)
func TestV37MigratesSDKv055KeyRotationFeePoolManagerKVParamsAndEVMCoinInfo(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Seed pre-upgrade 9-decimal staking params including KeyRotationFee in ubadge.
	sp, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	sp.BondDenom = appparams.LegacyBaseCoinUnit
	sp.KeyRotationFee = sdk.NewCoin(appparams.LegacyBaseCoinUnit, sdkmath.NewInt(5_000_000))
	require.NoError(t, app.StakingKeeper.SetParams(ctx, sp))

	// Seed pre-upgrade poolmanager KV params naming ubadge.
	pmParams := poolmanagertypes.DefaultParams()
	pmParams.TakerFeeParams.DefaultTakerFee = osmomath.MustNewDecFromStr("0.002")
	pmParams.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo = appparams.LegacyBaseCoinUnit
	pmParams.TakerFeeParams.CommunityPoolDenomWhitelist = []string{appparams.LegacyBaseCoinUnit, "ibc/1234"}
	app.PoolManagerKeeper.SetParams(ctx, pmParams)

	// Seed pre-upgrade 9-decimal EVM params and persisted EvmCoinInfo.
	evmParams := app.EVMKeeper.GetParams(ctx)
	evmParams.EvmDenom = appparams.LegacyBaseCoinUnit
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: appparams.BaseCoinUnit}
	require.NoError(t, app.EVMKeeper.SetParams(ctx, evmParams))
	require.NoError(t, app.EVMKeeper.SetEvmCoinInfo(ctx, evmtypes.EvmCoinInfo{
		Denom:         appparams.LegacyBaseCoinUnit,
		ExtendedDenom: appparams.BaseCoinUnit,
		DisplayDenom:  appparams.DisplayCoinUnit,
		Decimals:      appparams.LegacyBaseCoinDecimals,
	}))

	// Run the full v37 upgrade handler.
	require.NoError(t, v37.CustomUpgradeHandlerLogic(ctx, app.V37Keepers()))

	// 1. Staking params & KeyRotationFee.
	gotSP, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, appparams.BaseCoinUnit, gotSP.BondDenom)
	require.Equal(t, appparams.BaseCoinUnit, gotSP.KeyRotationFee.Denom)
	require.Equal(t, sdkmath.NewInt(5_000_000).Mul(v37.ConversionFactor).String(), gotSP.KeyRotationFee.Amount.String())
	require.NoError(t, gotSP.Validate())

	// 2. PoolManager KV params.
	gotPM := app.PoolManagerKeeper.GetParams(ctx)
	require.Equal(t, appparams.BaseCoinUnit, gotPM.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo)
	require.Equal(t, []string{appparams.BaseCoinUnit, "ibc/1234"}, gotPM.TakerFeeParams.CommunityPoolDenomWhitelist)

	// 3. EVM persisted EvmCoinInfo & process globals.
	gotCI := app.EVMKeeper.GetEvmCoinInfo(ctx)
	require.Equal(t, appparams.BaseCoinUnit, gotCI.Denom)
	require.Equal(t, appparams.BaseCoinUnit, gotCI.ExtendedDenom)
	require.Equal(t, uint32(appparams.BaseCoinDecimals), gotCI.Decimals)
	require.Equal(t, appparams.BaseCoinUnit, evmtypes.GetEVMCoinDenom())
	require.Equal(t, appparams.BaseCoinUnit, evmtypes.GetEVMCoinExtendedDenom())
	require.Equal(t, evmtypes.Decimals(appparams.BaseCoinDecimals), evmtypes.GetEVMCoinDecimals())
}
