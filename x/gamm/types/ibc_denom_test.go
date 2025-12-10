package types

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

type mockTransferKeeper struct {
	traces map[string]transfertypes.DenomTrace
}

func (m mockTransferKeeper) DenomPathFromHash(_ sdk.Context, denom string) (string, error) {
	// Simulate the same logic as the real DenomPathFromHash
	const hashedPrefix = "ibc/"
	if !strings.HasPrefix(denom, hashedPrefix) {
		return "", fmt.Errorf("denom %s is not an ibc/<hash> denom", denom)
	}

	hexHash := denom[len(hashedPrefix):]
	trace, found := m.traces[hexHash]
	if !found {
		return "", fmt.Errorf("trace not found for hash: %s", hexHash)
	}

	return trace.GetFullDenomPath(), nil
}

func TestExpandIBCDenomToFullPath_WithStoredTrace(t *testing.T) {
	ctx := sdk.Context{} // empty context is fine for this helper
	hash := "A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701"
	trace := transfertypes.DenomTrace{
		Path:      "transfer/channel-3",
		BaseDenom: "uatom",
	}

	keeper := mockTransferKeeper{
		traces: map[string]transfertypes.DenomTrace{
			hash: trace,
		},
	}

	full, err := ExpandIBCDenomToFullPath(ctx, "ibc/"+hash, keeper)
	require.NoError(t, err)
	require.Equal(t, trace.GetFullDenomPath(), full)
}

func TestExpandIBCDenomToFullPath_NoStoredTrace(t *testing.T) {
	ctx := sdk.Context{}
	hash := "ABCDEF1234"

	keeper := mockTransferKeeper{
		traces: map[string]transfertypes.DenomTrace{},
	}

	// Fallback ParseDenomTrace should return base denom when no path
	full, err := ExpandIBCDenomToFullPath(ctx, "ibc/"+hash, keeper)
	require.Error(t, err)
	require.Empty(t, full)
}

func TestExpandIBCDenomToFullPath_NonIBCDenom(t *testing.T) {
	ctx := sdk.Context{}
	keeper := mockTransferKeeper{}

	full, err := ExpandIBCDenomToFullPath(ctx, "uatom", keeper)
	require.NoError(t, err)
	require.Equal(t, "uatom", full)
}
