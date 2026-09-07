package tokenization

import (
	"fmt"
	"strings"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestJSONBudgetBoundsNestedFieldsBeforeSemanticValidation(t *testing.T) {
	for _, payload := range []string{
		`{"aliasPathsToAdd":[` + strings.Repeat(`{},`, 100) + `{}]}`,
		`{"collectionPermissions":{"canUpdateManager":[{"permanentlyPermittedTimes":[` + strings.Repeat(`{},`, 100) + `{}]}]}}`,
		`{"defaultBalances":{"balances":[` + strings.Repeat(`{},`, 100) + `{}]}}`,
		`{"standards":["` + strings.Repeat("x", MaxMetadataLength+1) + `"]}`,
		strings.Repeat("[", 34) + "0" + strings.Repeat("]", 34),
		`{} {}`,
	} {
		ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
		require.Error(t, meterJSONInput(ctx, payload))
	}
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
	require.NoError(t, meterJSONInput(ctx, fmt.Sprintf(`{"addresses":[%s"bb1last"]}`, strings.Repeat(`"bb1address",`, 999))))
	require.Positive(t, ctx.GasMeter().GasConsumed())
}

func TestJSONBudgetChargesBeforeRunningOutOfGas(t *testing.T) {
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewGasMeter(GasPerApprovalField))
	require.Panics(t, func() { _ = meterJSONInput(ctx, `{"standards":["a","b"]}`) })
}

func BenchmarkV35JSONBudget(b *testing.B) {
	for _, size := range []int{100, 10000} {
		b.Run(fmt.Sprintf("string-%d", size), func(b *testing.B) {
			payload := `{"standards":["` + strings.Repeat("a", size) + `"]}`
			b.ReportAllocs()
			var gas uint64
			for i := 0; i < b.N; i++ {
				ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
				if err := meterJSONInput(ctx, payload); err != nil {
					b.Fatal(err)
				}
				gas = ctx.GasMeter().GasConsumed()
			}
			b.ReportMetric(float64(gas), "gas/op")
		})
	}
}
