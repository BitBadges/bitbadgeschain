package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

func TestV35SetsBlockMaxGasAndMinGasPrice(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Seed the pre-upgrade values so the test proves the handler changes them.
	params, err := app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	params.Block.MaxGas = -1
	require.NoError(t, app.ConsensusParamsKeeper.ParamsStore.Set(ctx, params))

	fm := app.FeeMarketKeeper.GetParams(ctx)
	fm.MinGasPrice = sdkmath.LegacyZeroDec()
	require.NoError(t, app.FeeMarketKeeper.SetParams(ctx, fm))

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35.Keepers{
		ConsensusParams: app.ConsensusParamsKeeper,
		FeeMarket:       app.FeeMarketKeeper,
		IBCRateLimit:    app.IBCRateLimitKeeper,
		Tokenization:    app.TokenizationKeeper,
	}))

	params, err = app.ConsensusParamsKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, v35.BlockMaxGas, params.Block.MaxGas)

	want := sdkmath.LegacyMustNewDecFromStr(v35.MinGasPrice)
	require.True(t, app.FeeMarketKeeper.GetParams(ctx).MinGasPrice.Equal(want))
}
