package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

func TestVerifyCaller_ValidAddress(t *testing.T) {
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")
	err := gamm.VerifyCaller(caller)
	require.NoError(t, err)
}

func TestVerifyCaller_ZeroAddress(t *testing.T) {
	caller := common.Address{}
	err := gamm.VerifyCaller(caller)
	require.Error(t, err)
	// Just check that error is not nil (validation moved to ValidateBasic)
	require.NotNil(t, err)
}

func TestCheckOverflow_ValidValues(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{"zero", big.NewInt(0)},
		{"one", big.NewInt(1)},
		{"max_int256", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))}, // 2^255-1, max for int256
		{"large_value", big.NewInt(1e18)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.CheckOverflow(tt.value, "testField")
			require.NoError(t, err)
		})
	}
}

func TestCheckOverflow_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{"nil", nil},
		{"negative", big.NewInt(-1)},
		{"exceeds_uint256", new(big.Int).Lsh(big.NewInt(1), 256)}, // 2^256
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.CheckOverflow(tt.value, "testField")
			require.Error(t, err)
		})
	}
}

func TestValidateArraySize_ValidSizes(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		maxSize int
		wantErr bool
	}{
		{"zero", 0, 10, true}, // Zero is not allowed
		{"one", 1, 10, false},
		{"max", 10, 10, false},
		{"small", 5, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateArraySize(tt.size, tt.maxSize, "testField")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateArraySize_InvalidSizes(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		maxSize int
	}{
		{"exceeds_max", 11, 10},
		{"way_over", 100, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateArraySize(tt.size, tt.maxSize, "testField")
			require.Error(t, err)
			// Just check that error is not nil (validation moved to ValidateBasic)
			require.NotNil(t, err)
		})
	}
}

func TestValidateArraySizeAllowEmpty(t *testing.T) {
	// Empty array should be allowed
	err := gamm.ValidateArraySizeAllowEmpty(0, 10, "testField")
	require.NoError(t, err, "Empty array should be allowed by ValidateArraySizeAllowEmpty")

	// Non-empty array within bounds should be allowed
	err = gamm.ValidateArraySizeAllowEmpty(5, 10, "testField")
	require.NoError(t, err, "Array within bounds should be allowed")

	// Array at max should be allowed
	err = gamm.ValidateArraySizeAllowEmpty(10, 10, "testField")
	require.NoError(t, err, "Array at max should be allowed")

	// Array exceeding max should be rejected
	err = gamm.ValidateArraySizeAllowEmpty(11, 10, "testField")
	require.Error(t, err, "Array exceeding max should be rejected")
}

func TestGetCallerAddress_ValidCaller(t *testing.T) {
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// Test VerifyCaller directly
	err := gamm.VerifyCaller(caller)
	require.NoError(t, err)
}

func TestGetCallerAddress_ZeroAddress(t *testing.T) {
	caller := common.Address{} // Zero address

	// Test VerifyCaller directly
	err := gamm.VerifyCaller(caller)
	require.Error(t, err)
	// Just check that error is not nil (validation moved to ValidateBasic)
	require.NotNil(t, err)
}

func TestMaxArraySizes(t *testing.T) {
	// Verify max sizes are defined
	require.Equal(t, 10, gamm.MaxRoutes)
	require.Equal(t, 10, gamm.MaxAffiliates)
	require.Equal(t, 20, gamm.MaxCoins)
	require.Equal(t, 256, gamm.MaxMemoLength)
	require.Equal(t, 100, gamm.MaxPaginationLimit)
	require.Equal(t, uint64(10), uint64(gamm.GasPerMemoByte))
}

