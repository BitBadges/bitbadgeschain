package gamm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ValidatePoolId validates that a pool ID is valid (non-zero)
func ValidatePoolId(poolId uint64) error {
	if poolId == 0 {
		return ErrInvalidInput("poolId cannot be zero")
	}
	return nil
}

// ValidateShareAmount validates that a share amount is positive and not overflow
func ValidateShareAmount(amount *big.Int, fieldName string) error {
	if amount == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if amount.Sign() <= 0 {
		return ErrInvalidInput(fmt.Sprintf("%s must be greater than zero, got %s", fieldName, amount.String()))
	}
	if err := CheckOverflow(amount, fieldName); err != nil {
		return err
	}
	return nil
}

// ValidateCoin validates that a Coin struct is valid
func ValidateCoin(coin struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
}, fieldName string) error {
	if coin.Denom == "" {
		return ErrInvalidInput(fmt.Sprintf("%s.denom cannot be empty", fieldName))
	}
	if err := ValidateStringLength(coin.Denom, fmt.Sprintf("%s.denom", fieldName)); err != nil {
		return err
	}
	if err := ValidateShareAmount(coin.Amount, fmt.Sprintf("%s.amount", fieldName)); err != nil {
		return err
	}
	return nil
}

// ValidateCoins validates that a coins array is valid
func ValidateCoins(coins []struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
}, fieldName string) error {
	if err := ValidateArraySize(len(coins), MaxCoins, fieldName); err != nil {
		return err
	}
	for i, coin := range coins {
		if err := ValidateCoin(coin, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateCoinsAllowZero validates that a coins array is valid, allowing zero amounts
// This is used for tokenOutMins in exit pool where zero means no minimum
func ValidateCoinsAllowZero(coins []struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
}, fieldName string) error {
	if err := ValidateArraySize(len(coins), MaxCoins, fieldName); err != nil {
		return err
	}
	for i, coin := range coins {
		if coin.Denom == "" {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].denom cannot be empty", fieldName, i))
		}
		if err := ValidateStringLength(coin.Denom, fmt.Sprintf("%s[%d].denom", fieldName, i)); err != nil {
			return err
		}
		if coin.Amount == nil {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].amount cannot be nil", fieldName, i))
		}
		// Allow zero, but still check for overflow
		if coin.Amount.Sign() < 0 {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].amount cannot be negative, got %s", fieldName, i, coin.Amount.String()))
		}
		if err := CheckOverflow(coin.Amount, fmt.Sprintf("%s[%d].amount", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateRoutes validates that a routes array is valid
func ValidateRoutes(routes []struct {
	PoolId        uint64 `json:"poolId"`
	TokenOutDenom string `json:"tokenOutDenom"`
}, fieldName string) error {
	// Routes cannot be empty for swap operations
	if len(routes) == 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	if err := ValidateArraySize(len(routes), MaxRoutes, fieldName); err != nil {
		return err
	}
	for i, route := range routes {
		if err := ValidatePoolId(route.PoolId); err != nil {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].poolId: %v", fieldName, i, err))
		}
		if route.TokenOutDenom == "" {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].tokenOutDenom cannot be empty", fieldName, i))
		}
		if err := ValidateStringLength(route.TokenOutDenom, fmt.Sprintf("%s[%d].tokenOutDenom", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateIBCTransferInfo validates that IBCTransferInfo is valid
func ValidateIBCTransferInfo(ibcInfo struct {
	SourceChannel    string `json:"sourceChannel"`
	Receiver         string `json:"receiver"`
	Memo             string `json:"memo"`
	TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
}, ctxTimestamp uint64) error {
	if ibcInfo.SourceChannel == "" {
		return ErrInvalidInput("ibcTransferInfo.sourceChannel cannot be empty")
	}
	if err := ValidateStringLength(ibcInfo.SourceChannel, "ibcTransferInfo.sourceChannel"); err != nil {
		return err
	}
	if ibcInfo.Receiver == "" {
		return ErrInvalidInput("ibcTransferInfo.receiver cannot be empty")
	}
	if err := ValidateStringLength(ibcInfo.Receiver, "ibcTransferInfo.receiver"); err != nil {
		return err
	}
	// Validate memo length (IBC memos have stricter limits)
	if len(ibcInfo.Memo) > MaxMemoLength {
		return ErrInvalidInput(fmt.Sprintf("ibcTransferInfo.memo length (%d) exceeds maximum allowed length (%d)", len(ibcInfo.Memo), MaxMemoLength))
	}
	// Validate timeout timestamp is in the future
	if ibcInfo.TimeoutTimestamp <= ctxTimestamp {
		return ErrInvalidInput(fmt.Sprintf("ibcTransferInfo.timeoutTimestamp (%d) must be in the future (current: %d)", ibcInfo.TimeoutTimestamp, ctxTimestamp))
	}
	return nil
}

// ValidateAffiliates validates that an affiliates array is valid
func ValidateAffiliates(affiliates []struct {
	Address        common.Address `json:"address"`
	BasisPointsFee *big.Int       `json:"basisPointsFee"`
}, fieldName string) error {
	if len(affiliates) == 0 {
		// Affiliates are optional, so empty array is valid
		return nil
	}
	if len(affiliates) > MaxAffiliates {
		return ErrInvalidInput(fmt.Sprintf("%s size (%d) exceeds maximum allowed size (%d)", fieldName, len(affiliates), MaxAffiliates))
	}
	for i, affiliate := range affiliates {
		if err := ValidateAddress(affiliate.Address, fmt.Sprintf("%s[%d].address", fieldName, i)); err != nil {
			return err
		}
		if affiliate.BasisPointsFee == nil {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].basisPointsFee cannot be nil", fieldName, i))
		}
		if affiliate.BasisPointsFee.Sign() < 0 {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].basisPointsFee cannot be negative", fieldName, i))
		}
		// Basis points should be between 0 and 10000 (0% to 100%)
		maxBasisPoints := big.NewInt(10000)
		if affiliate.BasisPointsFee.Cmp(maxBasisPoints) > 0 {
			return ErrInvalidInput(fmt.Sprintf("%s[%d].basisPointsFee (%s) exceeds maximum (10000)", fieldName, i, affiliate.BasisPointsFee.String()))
		}
	}
	return nil
}

// ValidateAddress validates that an address is not zero
func ValidateAddress(addr common.Address, fieldName string) error {
	if addr == (common.Address{}) {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be zero address", fieldName))
	}
	return nil
}

// ValidateString validates that a string is not empty
func ValidateString(s string, fieldName string) error {
	if s == "" {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(offset, limit *big.Int) error {
	if offset == nil {
		return ErrInvalidInput("offset cannot be nil")
	}
	if offset.Sign() < 0 {
		return ErrInvalidInput("offset cannot be negative")
	}
	if limit == nil {
		return ErrInvalidInput("limit cannot be nil")
	}
	if limit.Sign() <= 0 {
		return ErrInvalidInput("limit must be greater than zero")
	}
	// Limit reasonable maximum for pagination
	maxLimit := big.NewInt(1000)
	if limit.Cmp(maxLimit) > 0 {
		return ErrInvalidInput(fmt.Sprintf("limit (%s) exceeds maximum allowed (1000)", limit.String()))
	}
	return nil
}

