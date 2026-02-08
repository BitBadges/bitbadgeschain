package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func TestConvertShareAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  *big.Int
		wantErr bool
	}{
		{"valid amount", big.NewInt(1000), false},
		{"zero amount", big.NewInt(0), true}, // Zero is not allowed for share amounts
		{"nil amount", nil, true},
		{"negative amount", big.NewInt(-1), true},
		{"large amount", big.NewInt(1e18), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertShareAmount(tt.amount)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, sdkmath.NewIntFromBigInt(tt.amount), result)
			}
		})
	}
}

func TestConvertCoinFromEVM(t *testing.T) {
	tests := []struct {
		name    string
		coin    struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}
		wantErr bool
	}{
		{
			name: "valid coin",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "uatom", Amount: big.NewInt(1000)},
			wantErr: false,
		},
		{
			name: "empty denom",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "", Amount: big.NewInt(1000)},
			wantErr: true,
		},
		{
			name: "nil amount",
			coin: struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{Denom: "uatom", Amount: nil},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertCoinFromEVM(tt.coin)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.coin.Denom, result.Denom)
				require.Equal(t, sdkmath.NewIntFromBigInt(tt.coin.Amount), result.Amount)
			}
		})
	}
}

func TestConvertCoinsFromEVM(t *testing.T) {
	validCoin := struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}{Denom: "uatom", Amount: big.NewInt(1000)}

	tests := []struct {
		name    string
		coins   []struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}
		wantErr bool
	}{
		{
			name:    "valid coins",
			coins:   []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{validCoin},
			wantErr: false,
		},
		{
			name:    "empty array",
			coins:   []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{},
			wantErr: true, // Empty array is not allowed by ValidateArraySize
		},
		{
			name: "invalid coin in array",
			coins: []struct{ Denom string `json:"denom"`; Amount *big.Int `json:"amount"` }{
				{Denom: "", Amount: big.NewInt(1000)}, // Empty denom
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertCoinsFromEVM(tt.coins)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, result, len(tt.coins))
			}
		})
	}
}

func TestConvertCoinToEVM(t *testing.T) {
	coin := sdk.NewCoin("uatom", sdkmath.NewInt(1000))
	result := gamm.ConvertCoinToEVM(coin)

	require.Equal(t, coin.Denom, result.Denom)
	require.Equal(t, coin.Amount.BigInt(), result.Amount)
}

func TestConvertCoinsToEVM(t *testing.T) {
	coins := sdk.NewCoins(
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		sdk.NewCoin("uosmo", sdkmath.NewInt(2000)),
	)

	result := gamm.ConvertCoinsToEVM(coins)

	require.Len(t, result, len(coins))
	require.Equal(t, coins[0].Denom, result[0].Denom)
	require.Equal(t, coins[0].Amount.BigInt(), result[0].Amount)
}

func TestConvertSwapRoutesFromEVM(t *testing.T) {
	validRoute := struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}{PoolId: 1, TokenOutDenom: "uosmo"}

	tests := []struct {
		name    string
		routes  []struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}
		wantErr bool
	}{
		{
			name:    "valid routes",
			routes:  []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{validRoute},
			wantErr: false,
		},
		{
			name:    "empty array",
			routes:  []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{},
			wantErr: true, // Routes cannot be empty
		},
		{
			name: "zero pool ID",
			routes: []struct{ PoolId uint64 `json:"poolId"`; TokenOutDenom string `json:"tokenOutDenom"` }{
				{PoolId: 0, TokenOutDenom: "uosmo"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertSwapRoutesFromEVM(tt.routes)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, result, len(tt.routes))
				require.Equal(t, poolmanagertypes.SwapAmountInRoute{
					PoolId:        tt.routes[0].PoolId,
					TokenOutDenom: tt.routes[0].TokenOutDenom,
				}, result[0])
			}
		})
	}
}

func TestConvertAffiliatesFromEVM(t *testing.T) {
	validAffiliate := struct {
		Address        common.Address `json:"address"`
		BasisPointsFee *big.Int       `json:"basisPointsFee"`
	}{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: big.NewInt(100)}

	tests := []struct {
		name       string
		affiliates []struct {
			Address        common.Address `json:"address"`
			BasisPointsFee *big.Int       `json:"basisPointsFee"`
		}
		wantErr bool
	}{
		{
			name:       "valid affiliates",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{validAffiliate},
			wantErr:    false,
		},
		{
			name:       "empty array",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{},
			wantErr:    false, // Empty array is valid
		},
		{
			name: "invalid affiliate",
			affiliates: []struct{ Address common.Address `json:"address"`; BasisPointsFee *big.Int `json:"basisPointsFee"` }{
				{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), BasisPointsFee: nil},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertAffiliatesFromEVM(tt.affiliates)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, result, len(tt.affiliates))
			}
		})
	}
}

func TestConvertIBCTransferInfoFromEVM(t *testing.T) {
	// Use a valid future timestamp (nanoseconds since epoch)
	futureTimestamp := uint64(9999999999) // Year 2286, should be far enough in the future for tests

	tests := []struct {
		name     string
		ibcInfo  struct {
			SourceChannel    string `json:"sourceChannel"`
			Receiver         string `json:"receiver"`
			Memo             string `json:"memo"`
			TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
		}
		wantErr bool
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
			wantErr: false,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gamm.ConvertIBCTransferInfoFromEVM(tt.ibcInfo)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.ibcInfo.SourceChannel, result.SourceChannel)
				require.Equal(t, tt.ibcInfo.Receiver, result.Receiver)
				require.Equal(t, tt.ibcInfo.Memo, result.Memo)
				require.Equal(t, tt.ibcInfo.TimeoutTimestamp, result.TimeoutTimestamp)
			}
		})
	}
}

