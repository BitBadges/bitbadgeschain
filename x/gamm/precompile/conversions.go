package gamm

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EVMCoin represents a coin for EVM return values
type EVMCoin struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
}

// ConvertCoinToEVM converts sdk.Coin to Solidity Coin struct
// This is used for response packing in transaction handlers
func ConvertCoinToEVM(coin sdk.Coin) EVMCoin {
	return EVMCoin{
		Denom:  coin.Denom,
		Amount: coin.Amount.BigInt(),
	}
}

// ConvertCoinToEVMSafe converts sdk.Coin to Solidity Coin struct with overflow check.
// Returns error if the amount exceeds uint256 max.
func ConvertCoinToEVMSafe(coin sdk.Coin, fieldName string) (EVMCoin, error) {
	amount := coin.Amount.BigInt()
	if err := CheckUint256Overflow(amount, fmt.Sprintf("%s.amount", fieldName)); err != nil {
		return EVMCoin{}, err
	}
	return EVMCoin{
		Denom:  coin.Denom,
		Amount: amount,
	}, nil
}

// ConvertCoinsToEVM converts sdk.Coins to an array of Solidity Coin structs
// Returns structs with json tags to match ABI definition
// This is used for response packing in transaction handlers
func ConvertCoinsToEVM(coins sdk.Coins) []EVMCoin {
	result := make([]EVMCoin, len(coins))
	for i, coin := range coins {
		result[i] = ConvertCoinToEVM(coin)
	}
	return result
}

// ConvertCoinsToEVMSafe converts sdk.Coins to an array of Solidity Coin structs with overflow checks.
// Returns error if any amount exceeds uint256 max.
func ConvertCoinsToEVMSafe(coins sdk.Coins, fieldName string) ([]EVMCoin, error) {
	result := make([]EVMCoin, len(coins))
	for i, coin := range coins {
		coinEVM, err := ConvertCoinToEVMSafe(coin, fmt.Sprintf("%s[%d]", fieldName, i))
		if err != nil {
			return nil, err
		}
		result[i] = coinEVM
	}
	return result, nil
}
