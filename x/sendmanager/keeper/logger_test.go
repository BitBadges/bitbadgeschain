package keeper_test

import (
	"testing"

	"cosmossdk.io/log/v2"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	"github.com/stretchr/testify/require"
)

func TestKeeperLoggerSafeDefaultAndOverride(t *testing.T) {
	// Zero-value keeper must not panic when Logger() is called
	var zeroKeeper keeper.Keeper
	require.NotPanics(t, func() {
		require.NotNil(t, zeroKeeper.Logger())
	})

	// NewKeeper initializes a non-nil logger and SetLogger overrides safely
	k := keeper.NewKeeper(nil, nil, nil, nil, nil, nil)
	require.NotPanics(t, func() {
		require.NotNil(t, k.Logger())
	})
	k.SetLogger(nil)
	require.NotPanics(t, func() {
		require.NotNil(t, k.Logger())
	})
	k.SetLogger(log.NewNopLogger())
	require.NotNil(t, k.Logger())
}
