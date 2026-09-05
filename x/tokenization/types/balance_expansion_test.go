package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestBalanceExpansionChargesBeforeMaterialization(t *testing.T) {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		{Start: sdkmath.NewUint(3), End: sdkmath.NewUint(3)},
	}
	balance := &types.Balance{Amount: sdkmath.OneUint(), TokenIds: ranges, OwnershipTimes: ranges}
	for name, operation := range map[string]func(sdk.Context){
		"fetch requested": func(ctx sdk.Context) { _, _ = types.GetBalancesForIds(ctx, ranges, ranges, nil) },
		"fetch existing":  func(ctx sdk.Context) { _, _ = types.GetBalancesForIds(ctx, nil, nil, []*types.Balance{balance}) },
		"delete existing": func(ctx sdk.Context) { _, _ = types.DeleteBalances(ctx, nil, nil, []*types.Balance{balance}) },
		"delete requested": func(ctx sdk.Context) {
			_, _ = types.DeleteBalances(ctx, ranges, ranges, []*types.Balance{{Amount: sdkmath.OneUint()}})
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := testutil.DefaultContext(storetypes.NewKVStoreKey("test"), storetypes.NewTransientStoreKey("transient")).WithGasMeter(storetypes.NewGasMeter(100))
			require.PanicsWithValue(t, storetypes.ErrorOutOfGas{Descriptor: "balance range expansion"}, func() { operation(ctx) })
			require.Equal(t, uint64(4000), ctx.GasMeter().GasConsumed())
		})
	}
}
