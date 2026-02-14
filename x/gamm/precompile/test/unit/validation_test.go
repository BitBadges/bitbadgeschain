package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

func TestValidatePoolId(t *testing.T) {
	tests := []struct {
		name   string
		poolId uint64
		wantErr bool
	}{
		{"valid pool ID", 1, false},
		{"zero pool ID", 0, true},
		{"large pool ID", 1000000, false},
		{"max uint64", ^uint64(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidatePoolId(tt.poolId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateShareAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   *big.Int
		fieldName string
		wantErr  bool
	}{
		{"valid amount", big.NewInt(1000), "shareAmount", false},
		{"zero amount", big.NewInt(0), "shareAmount", true}, // Zero is not allowed (Sign() <= 0)
		{"nil amount", nil, "shareAmount", true},
		{"negative amount", big.NewInt(-1), "shareAmount", true},
		{"large amount", big.NewInt(1e18), "shareAmount", false},
		{"max uint256", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)), "shareAmount", true}, // Exceeds int256 max (2^255-1)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateShareAmount(tt.amount, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCoin(t *testing.T) {
	tests := []struct {
		name     string
		coin     struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}
		fieldName string
		wantErr  bool
	}{
		{
			name: "valid coin",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "uatom", Amount: big.NewInt(1000)},
			fieldName: "coin",
			wantErr:   false,
		},
		{
			name: "empty denom",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "", Amount: big.NewInt(1000)},
			fieldName: "coin",
			wantErr:   true,
		},
		{
			name: "nil amount",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "uatom", Amount: nil},
			fieldName: "coin",
			wantErr:   true,
		},
		{
			name: "negative amount",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "uatom", Amount: big.NewInt(-1)},
			fieldName: "coin",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateCoin(tt.coin, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCoins(t *testing.T) {
	validCoin := struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}{Denom: "uatom", Amount: big.NewInt(1000)}

	tests := []struct {
		name     string
		coins    []struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}
		fieldName string
		wantErr  bool
	}{
		{
			name:      "valid coins",
			coins:     []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{validCoin},
			fieldName: "coins",
			wantErr:   false,
		},
		{
			name:      "empty array",
			coins:     []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{},
			fieldName: "coins",
			wantErr:   true, // Empty array is not allowed by ValidateArraySize
		},
		{
			name: "exceeds max size",
			coins: func() []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` } {
				coins := make([]struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }, gamm.MaxCoins+1)
				for i := range coins {
					coins[i] = validCoin
				}
				return coins
			}(),
			fieldName: "coins",
			wantErr:   true,
		},
		{
			name: "invalid coin in array",
			coins: []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{
				{Denom: "", Amount: big.NewInt(1000)}, // Empty denom
			},
			fieldName: "coins",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateCoins(tt.coins, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRoutes(t *testing.T) {
	validRoute := struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}{PoolId: 1, TokenOutDenom: "uosmo"}

	tests := []struct {
		name     string
		routes   []struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}
		fieldName string
		wantErr  bool
	}{
		{
			name:      "valid routes",
			routes:    []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{validRoute},
			fieldName: "routes",
			wantErr:   false,
		},
		{
			name:      "empty array",
			routes:    []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{},
			fieldName: "routes",
			wantErr:   true, // Routes cannot be empty
		},
		{
			name: "exceeds max size",
			routes: func() []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` } {
				routes := make([]struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }, gamm.MaxRoutes+1)
				for i := range routes {
					routes[i] = validRoute
				}
				return routes
			}(),
			fieldName: "routes",
			wantErr:   true,
		},
		{
			name: "zero pool ID",
			routes: []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{
				{PoolId: 0, TokenOutDenom: "uosmo"},
			},
			fieldName: "routes",
			wantErr:   true,
		},
		{
			name: "empty token out denom",
			routes: []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{
				{PoolId: 1, TokenOutDenom: ""},
			},
			fieldName: "routes",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to the correct type for validation
			routes := make([]struct {
				PoolId        uint64 `json:"poolId"`
				TokenOutDenom string `json:"tokenOutDenom"`
			}, len(tt.routes))
			copy(routes, tt.routes)
			
			err := gamm.ValidateRoutes(routes, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateAffiliates(t *testing.T) {
	validAffiliate := struct {
		Address        common.Address `json:"address"`
		BasisPointsFee *big.Int       `json:"basisPointsFee"`
	}{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: big.NewInt(100)}

	tests := []struct {
		name      string
		affiliates []struct {
			Address        common.Address `json:"address"`
			BasisPointsFee *big.Int       `json:"basisPointsFee"`
		}
		fieldName string
		wantErr  bool
	}{
		{
			name:      "valid affiliates",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{validAffiliate},
			fieldName: "affiliates",
			wantErr:   false,
		},
		{
			name:      "empty array",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{},
			fieldName: "affiliates",
			wantErr:   false, // Empty array is valid
		},
		{
			name: "exceeds max size",
			affiliates: func() []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` } {
				affiliates := make([]struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }, gamm.MaxAffiliates+1)
				for i := range affiliates {
					affiliates[i] = validAffiliate
				}
				return affiliates
			}(),
			fieldName: "affiliates",
			wantErr:   true,
		},
		{
			name: "nil basis points fee",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{
				{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: nil},
			},
			fieldName: "affiliates",
			wantErr:   true,
		},
		{
			name: "negative basis points fee",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{
				{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: big.NewInt(-1)},
			},
			fieldName: "affiliates",
			wantErr:   true,
		},
		{
			name: "basis points fee exceeds 10000",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{
				{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: big.NewInt(10001)},
			},
			fieldName: "affiliates",
			wantErr:   true,
		},
		{
			name: "zero address",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{
				{Address: common.Address{}, BasisPointsFee: big.NewInt(100)},
			},
			fieldName: "affiliates",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateAffiliates(tt.affiliates, tt.fieldName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateIBCTransferInfo(t *testing.T) {
	// Use a valid future timestamp (nanoseconds since epoch)
	// Year 2286, should be far enough in the future for tests
	futureTimestamp := uint64(9999999999)
	// Current timestamp (past)
	currentTimestamp := uint64(1000000000)

	tests := []struct {
		name     string
		ibcInfo  struct {
			SourceChannel    string `json:"sourceChannel"`
			Receiver         string `json:"receiver"`
			Memo             string `json:"memo"`
			TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
		}
		ctxTimestamp uint64
		wantErr      bool
	}{
		{
			name: "valid IBC info",
			ibcInfo: struct {
				SourceChannel    string `json:"sourceChannel"`
				Receiver         string `json:"receiver"`
				Memo             string `json:"memo"`
				TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
			}{
				SourceChannel:    "channel-0",
				Receiver:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
				Memo:             "test memo",
				TimeoutTimestamp: futureTimestamp,
			},
			ctxTimestamp: currentTimestamp,
			wantErr:      false,
		},
		{
			name: "empty source channel",
			ibcInfo: struct {
				SourceChannel    string `json:"sourceChannel"`
				Receiver         string `json:"receiver"`
				Memo             string `json:"memo"`
				TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
			}{
				SourceChannel:    "",
				Receiver:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
				Memo:             "test memo",
				TimeoutTimestamp: futureTimestamp,
			},
			ctxTimestamp: currentTimestamp,
			wantErr:      true,
		},
		{
			name: "empty receiver",
			ibcInfo: struct {
				SourceChannel    string `json:"sourceChannel"`
				Receiver         string `json:"receiver"`
				Memo             string `json:"memo"`
				TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
			}{
				SourceChannel:    "channel-0",
				Receiver:         "",
				Memo:             "test memo",
				TimeoutTimestamp: futureTimestamp,
			},
			ctxTimestamp: currentTimestamp,
			wantErr:      true,
		},
		{
			name: "memo too long",
			ibcInfo: struct {
				SourceChannel    string `json:"sourceChannel"`
				Receiver         string `json:"receiver"`
				Memo             string `json:"memo"`
				TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
			}{
				SourceChannel:    "channel-0",
				Receiver:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
				Memo:             string(make([]byte, gamm.MaxMemoLength+1)), // Use MaxMemoLength for IBC memo validation
				TimeoutTimestamp: futureTimestamp,
			},
			ctxTimestamp: currentTimestamp,
			wantErr:      true,
		},
		{
			name: "timeout in past",
			ibcInfo: struct {
				SourceChannel    string `json:"sourceChannel"`
				Receiver         string `json:"receiver"`
				Memo             string `json:"memo"`
				TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
			}{
				SourceChannel:    "channel-0",
				Receiver:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
				Memo:             "test memo",
				TimeoutTimestamp: currentTimestamp - 1, // Past timestamp
			},
			ctxTimestamp: currentTimestamp,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidateIBCTransferInfo(tt.ibcInfo, tt.ctxTimestamp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		offset  *big.Int
		limit   *big.Int
		wantErr bool
	}{
		{"valid pagination", big.NewInt(0), big.NewInt(10), false},
		{"nil offset", nil, big.NewInt(10), true},
		{"nil limit", big.NewInt(0), nil, true},
		{"negative offset", big.NewInt(-1), big.NewInt(10), true},
		{"negative limit", big.NewInt(0), big.NewInt(-1), true},
		{"limit exceeds max", big.NewInt(0), big.NewInt(1001), true}, // ValidatePagination uses 1000 as max
		{"max limit", big.NewInt(0), big.NewInt(int64(gamm.MaxPaginationLimit)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gamm.ValidatePagination(tt.offset, tt.limit)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

