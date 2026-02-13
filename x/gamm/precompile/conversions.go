package gamm

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConvertCoinToEVM converts sdk.Coin to Solidity Coin struct
// This is used for response packing in transaction handlers
func ConvertCoinToEVM(coin sdk.Coin) struct {
	Denom  string
	Amount *big.Int
} {
	return struct {
		Denom  string
		Amount *big.Int
	}{
		Denom:  coin.Denom,
		Amount: coin.Amount.BigInt(),
	}
}

// ConvertCoinsToEVM converts sdk.Coins to an array of Solidity Coin structs
// Returns structs with json tags to match ABI definition
// This is used for response packing in transaction handlers
func ConvertCoinsToEVM(coins sdk.Coins) []struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
} {
	result := make([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}, len(coins))
	for i, coin := range coins {
		coinEVM := ConvertCoinToEVM(coin)
		result[i] = struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}{
			Denom:  coinEVM.Denom,
			Amount: coinEVM.Amount,
		}
	}
	return result
}
