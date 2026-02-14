package types

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type mockTransferKeeper struct {
	// Map from hash to full denom path (e.g., "transfer/channel-3/uatom")
	traces map[string]string
}

func (m mockTransferKeeper) DenomPathFromHash(_ sdk.Context, denom string) (string, error) {
	// Simulate the same logic as the real DenomPathFromHash
	const hashedPrefix = "ibc/"
	if !strings.HasPrefix(denom, hashedPrefix) {
		return "", fmt.Errorf("denom %s is not an ibc/<hash> denom", denom)
	}

	hexHash := denom[len(hashedPrefix):]
	path, found := m.traces[hexHash]
	if !found {
		return "", fmt.Errorf("trace not found for hash: %s", hexHash)
	}

	return path, nil
}

func TestExpandIBCDenomToFullPath_WithStoredTrace(t *testing.T) {
	ctx := sdk.Context{} // empty context is fine for this helper
	hash := "A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701"
	// Full path format: "transfer/channel-3/uatom"
	fullPath := "transfer/channel-3/uatom"

	keeper := mockTransferKeeper{
		traces: map[string]string{
			hash: fullPath,
		},
	}

	full, err := ExpandIBCDenomToFullPath(ctx, "ibc/"+hash, keeper)
	require.NoError(t, err)
	require.Equal(t, fullPath, full)
}

func TestExpandIBCDenomToFullPath_NoStoredTrace(t *testing.T) {
	ctx := sdk.Context{}
	hash := "ABCDEF1234"

	keeper := mockTransferKeeper{
		traces: map[string]string{},
	}

	// DenomPathFromHash should return error when trace not found
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
