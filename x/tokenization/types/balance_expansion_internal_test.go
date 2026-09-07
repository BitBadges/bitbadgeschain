package types

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestBalanceExpansionGasArithmetic(t *testing.T) {
	ctx := testutil.DefaultContext(storetypes.NewKVStoreKey("test"), storetypes.NewTransientStoreKey("transient")).WithGasMeter(storetypes.NewInfiniteGasMeter())
	chargeBalanceExpansion(ctx, 0, 100)
	chargeBalanceExpansion(ctx, 100, 0)
	require.Zero(t, ctx.GasMeter().GasConsumed())
	chargeBalanceExpansion(ctx, 2, 3)
	require.Equal(t, uint64(6000), ctx.GasMeter().GasConsumed())
	maxInt := int(^uint(0) >> 1)
	for _, dimensions := range [][2]int{{maxInt, maxInt}, {maxInt, 1}} {
		require.PanicsWithValue(t, storetypes.ErrorGasOverflow{Descriptor: "balance range expansion"}, func() {
			chargeBalanceExpansion(ctx, dimensions[0], dimensions[1])
		})
	}
	require.Equal(t, uint64(6000), ctx.GasMeter().GasConsumed())
}
