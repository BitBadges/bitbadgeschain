package tokenization

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// TestValidateAddress tests the ValidateAddress function
func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name      string
		address   common.Address
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid address",
			address:   common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"),
			fieldName: "userAddress",
			wantErr:   false,
		},
		{
			name:      "zero address",
			address:   common.Address{},
			fieldName: "userAddress",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.address, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateAddresses tests the ValidateAddresses function
func TestValidateAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addresses []common.Address
		fieldName string
		wantErr   bool
	}{
		{
			name: "valid addresses",
			addresses: []common.Address{
				common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"),
				common.HexToAddress("0x8ba1f109551bD432803012645Hac136c22C9e7"),
			},
			fieldName: "toAddresses",
			wantErr:   false,
		},
		{
			name:      "empty addresses",
			addresses: []common.Address{},
			fieldName: "toAddresses",
			wantErr:   true,
		},
		{
			name: "contains zero address",
			addresses: []common.Address{
				common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"),
				common.Address{},
			},
			fieldName: "toAddresses",
			wantErr:   true,
		},
		{
			name: "too many addresses",
			addresses: func() []common.Address {
				addrs := make([]common.Address, MaxRecipients+1)
				for i := range addrs {
					addrs[i] = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
				}
				return addrs
			}(),
			fieldName: "toAddresses",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddresses(tt.addresses, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateCollectionId tests the ValidateCollectionId function
func TestValidateCollectionId(t *testing.T) {
	tests := []struct {
		name        string
		collectionId *big.Int
		wantErr     bool
	}{
		{
			name:        "valid collection ID",
			collectionId: big.NewInt(1),
			wantErr:     false,
		},
		{
			name:        "zero collection ID",
			collectionId: big.NewInt(0),
			wantErr:     false, // 0 is valid (used for creating new collections)
		},
		{
			name:        "negative collection ID",
			collectionId: big.NewInt(-1),
			wantErr:     true,
		},
		{
			name:        "nil collection ID",
			collectionId: nil,
			wantErr:     true,
		},
		{
			name:        "large collection ID",
			collectionId: new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCollectionId(tt.collectionId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateAmount tests the ValidateAmount function
func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    *big.Int
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid amount",
			amount:    big.NewInt(100),
			fieldName: "amount",
			wantErr:   false,
		},
		{
			name:      "zero amount",
			amount:    big.NewInt(0),
			fieldName: "amount",
			wantErr:   true,
		},
		{
			name:      "negative amount",
			amount:    big.NewInt(-1),
			fieldName: "amount",
			wantErr:   true,
		},
		{
			name:      "nil amount",
			amount:    nil,
			fieldName: "amount",
			wantErr:   true,
		},
		{
			name:      "large amount",
			amount:    new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			fieldName: "amount",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmount(tt.amount, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateBigIntRanges tests the ValidateBigIntRanges function
func TestValidateBigIntRanges(t *testing.T) {
	tests := []struct {
		name      string
		ranges    []struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}
		fieldName string
		wantErr   bool
	}{
		{
			name: "valid ranges",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(10)},
				{Start: big.NewInt(20), End: big.NewInt(30)},
			},
			fieldName: "tokenIds",
			wantErr:   false,
		},
		{
			name:      "empty ranges",
			ranges:    []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{},
			fieldName: "tokenIds",
			wantErr:   true,
		},
		{
			name: "start greater than end",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(10), End: big.NewInt(5)},
			},
			fieldName: "tokenIds",
			wantErr:   true,
		},
		{
			name: "nil start",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: nil, End: big.NewInt(10)},
			},
			fieldName: "tokenIds",
			wantErr:   true,
		},
		{
			name: "nil end",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: nil},
			},
			fieldName: "tokenIds",
			wantErr:   true,
		},
		{
			name: "negative start",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(-1), End: big.NewInt(10)},
			},
			fieldName: "tokenIds",
			wantErr:   true,
		},
		{
			name: "negative end",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(-1)},
			},
			fieldName: "tokenIds",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBigIntRanges(tt.ranges, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateString tests the ValidateString function
func TestValidateString(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid string",
			str:       "test-id",
			fieldName: "approvalId",
			wantErr:   false,
		},
		{
			name:      "empty string",
			str:       "",
			fieldName: "approvalId",
			wantErr:   true,
		},
		{
			name:      "whitespace only",
			str:       "   ",
			fieldName: "approvalId",
			wantErr:   false, // ValidateString doesn't trim, so this is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.str, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestCheckOverflow tests the CheckOverflow function
func TestCheckOverflow(t *testing.T) {
	tests := []struct {
		name      string
		value     *big.Int
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid positive value",
			value:     big.NewInt(100),
			fieldName: "amount",
			wantErr:   false,
		},
		{
			name:      "zero value",
			value:     big.NewInt(0),
			fieldName: "amount",
			wantErr:   false,
		},
		{
			name:      "negative value",
			value:     big.NewInt(-1),
			fieldName: "amount",
			wantErr:   true,
		},
		{
			name:      "nil value",
			value:     nil,
			fieldName: "amount",
			wantErr:   true,
		},
		{
			name:      "large value",
			value:     new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
			fieldName: "amount",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckOverflow(tt.value, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestVerifyCaller tests the VerifyCaller function
func TestVerifyCaller(t *testing.T) {
	tests := []struct {
		name    string
		caller  common.Address
		wantErr bool
	}{
		{
			name:    "valid caller",
			caller:  common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"),
			wantErr: false,
		},
		{
			name:    "zero address caller",
			caller:  common.Address{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyCaller(tt.caller)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

