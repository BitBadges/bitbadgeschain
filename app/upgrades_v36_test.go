package app

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	v36 "github.com/bitbadges/bitbadgeschain/app/upgrades/v36"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestV36MigratesLegacyPoolManagerParams(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	adminAddr := sdk.AccAddress("poolmanager_admin___").String()
	whitelistAddr := sdk.AccAddress("poolmanager_wl______").String()

	expected := poolmanagertypes.Params{
		TakerFeeParams: poolmanagertypes.TakerFeeParams{
			DefaultTakerFee: osmomath.MustNewDecFromStr("0.0025"),
			OsmoTakerFeeDistribution: poolmanagertypes.TakerFeeDistributionPercentage{
				StakingRewards: osmomath.MustNewDecFromStr("0.8"),
				CommunityPool:  osmomath.MustNewDecFromStr("0.2"),
			},
			NonOsmoTakerFeeDistribution: poolmanagertypes.TakerFeeDistributionPercentage{
				StakingRewards: osmomath.MustNewDecFromStr("0.5"),
				CommunityPool:  osmomath.MustNewDecFromStr("0.5"),
			},
			AdminAddresses:                                 []string{adminAddr},
			CommunityPoolDenomToSwapNonWhitelistedAssetsTo: "ubadge",
			ReducedFeeWhitelist:                            []string{whitelistAddr},
			CommunityPoolDenomWhitelist:                    []string{"ubadge"},
		},
	}

	// Seed the legacy "params" KV store in the exact Amino JSON format x/params wrote.
	legacyKey := app.GetKey(v36.LegacyParamsStoreKey)
	require.NotNil(t, legacyKey)
	legacyStore := ctx.KVStore(legacyKey)
	writeLegacy := func(paramKey []byte, val interface{}) {
		bz, err := app.LegacyAmino().MarshalJSON(val)
		require.NoError(t, err)
		fullKey := append(append([]byte(poolmanagertypes.ModuleName), '/'), paramKey...)
		legacyStore.Set(fullKey, bz)
	}
	writeLegacy(poolmanagertypes.KeyDefaultTakerFee, expected.TakerFeeParams.DefaultTakerFee)
	writeLegacy(poolmanagertypes.KeyOsmoTakerFeeDistribution, expected.TakerFeeParams.OsmoTakerFeeDistribution)
	writeLegacy(poolmanagertypes.KeyNonOsmoTakerFeeDistribution, expected.TakerFeeParams.NonOsmoTakerFeeDistribution)
	writeLegacy(poolmanagertypes.KeyAdminAddresses, expected.TakerFeeParams.AdminAddresses)
	writeLegacy(poolmanagertypes.KeyCommunityPoolDenomToSwapNonWhitelistedAssetsTo, expected.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo)
	writeLegacy(poolmanagertypes.KeyReducedTakerFeeByWhitelist, expected.TakerFeeParams.ReducedFeeWhitelist)
	writeLegacy(poolmanagertypes.KeyCommunityPoolDenomWhitelist, expected.TakerFeeParams.CommunityPoolDenomWhitelist)

	require.NoError(t, v36.CustomUpgradeHandlerLogic(ctx, v36.Keepers{
		Account:         app.AccountKeeper,
		LegacyAmino:     app.LegacyAmino(),
		LegacyParamsKey: legacyKey,
		PoolManager:     &app.PoolManagerKeeper,
	}))

	got := app.PoolManagerKeeper.GetParams(ctx)
	require.Equal(t, expected, got)
	require.True(t, app.PoolManagerKeeper.GetDefaultTakerFee(ctx).Equal(expected.TakerFeeParams.DefaultTakerFee))
}

func TestV36PoolManagerMigrationFallsBackToDefaults(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	require.NoError(t, v36.CustomUpgradeHandlerLogic(ctx, v36.Keepers{
		Account:         app.AccountKeeper,
		LegacyAmino:     app.LegacyAmino(),
		LegacyParamsKey: app.GetKey(v36.LegacyParamsStoreKey),
		PoolManager:     &app.PoolManagerKeeper,
	}))

	require.Equal(t, poolmanagertypes.DefaultParams(), app.PoolManagerKeeper.GetParams(ctx))
}

func TestV36KeyRotationFeePoolAccountWired(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	macc := app.AccountKeeper.GetModuleAccount(ctx, stakingtypes.KeyRotationFeePoolName)
	require.NotNil(t, macc)
	require.True(t, macc.HasPermission(authtypes.Burner))
	require.True(t, app.BankKeeper.BlockedAddr(macc.GetAddress()))
}

func TestV36ConvertsExistingBaseAccountAtKeyRotationFeePool(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	addr := app.AccountKeeper.GetModuleAddress(stakingtypes.KeyRotationFeePoolName)
	require.NotNil(t, addr)
	existing := app.AccountKeeper.GetAccount(ctx, addr)
	require.NotNil(t, existing)

	// Simulate an address that existed as a plain BaseAccount prior to v36.
	base := authtypes.NewBaseAccount(addr, nil, existing.GetAccountNumber(), existing.GetSequence())
	app.AccountKeeper.SetAccount(ctx, base)

	require.NoError(t, v36.EnsureKeyRotationFeePoolAccount(ctx, app.AccountKeeper))

	macc := app.AccountKeeper.GetModuleAccount(ctx, stakingtypes.KeyRotationFeePoolName)
	require.NotNil(t, macc)
	require.True(t, macc.HasPermission(authtypes.Burner))
	require.Equal(t, existing.GetAccountNumber(), macc.GetAccountNumber())
}
