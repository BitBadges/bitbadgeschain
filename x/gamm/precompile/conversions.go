package gamm

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// ConvertCoinFromEVM converts a Solidity Coin struct to sdk.Coin
func ConvertCoinFromEVM(coin struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
},
) (sdk.Coin, error) {
	if err := ValidateCoin(coin, "coin"); err != nil {
		return sdk.Coin{}, err
	}
	amount := sdkmath.NewIntFromBigInt(coin.Amount)
	return sdk.NewCoin(coin.Denom, amount), nil
}

// ConvertCoinToEVM converts sdk.Coin to Solidity Coin struct
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

// ConvertCoinsFromEVM converts an array of Solidity Coin structs to sdk.Coins
func ConvertCoinsFromEVM(coins []struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
},
) (sdk.Coins, error) {
	if err := ValidateCoins(coins, "coins"); err != nil {
		return nil, err
	}
	sdkCoins := make(sdk.Coins, len(coins))
	for i, coin := range coins {
		sdkCoin, err := ConvertCoinFromEVM(coin)
		if err != nil {
			return nil, fmt.Errorf("coins[%d]: %w", i, err)
		}
		sdkCoins[i] = sdkCoin
	}
	return sdkCoins.Sort(), nil
}

// ConvertCoinsFromEVMAllowZero converts an array of Solidity Coin structs to sdk.Coins, allowing zero amounts
// This is used for tokenOutMins in exit pool where zero means no minimum
func ConvertCoinsFromEVMAllowZero(coins []struct {
	Denom  string   `json:"denom"`
	Amount *big.Int `json:"amount"`
},
) (sdk.Coins, error) {
	if err := ValidateCoinsAllowZero(coins, "coins"); err != nil {
		return nil, err
	}
	sdkCoins := make(sdk.Coins, 0, len(coins))
	for i, coin := range coins {
		// Skip zero amounts (they mean no minimum)
		if coin.Amount.Sign() == 0 {
			continue
		}
		sdkCoin, err := ConvertCoinFromEVM(coin)
		if err != nil {
			return nil, fmt.Errorf("coins[%d]: %w", i, err)
		}
		sdkCoins = append(sdkCoins, sdkCoin)
	}
	return sdkCoins.Sort(), nil
}

// ConvertCoinsToEVM converts sdk.Coins to an array of Solidity Coin structs
// Returns structs with json tags to match ABI definition
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

// ConvertSwapRouteFromEVM converts a Solidity SwapAmountInRoute struct to poolmanagertypes.SwapAmountInRoute
func ConvertSwapRouteFromEVM(route struct {
	PoolId        uint64 `json:"poolId"`
	TokenOutDenom string `json:"tokenOutDenom"`
},
) poolmanagertypes.SwapAmountInRoute {
	return poolmanagertypes.SwapAmountInRoute{
		PoolId:        route.PoolId,
		TokenOutDenom: route.TokenOutDenom,
	}
}

// ConvertSwapRoutesFromEVM converts an array of Solidity SwapAmountInRoute structs
func ConvertSwapRoutesFromEVM(routes []struct {
	PoolId        uint64 `json:"poolId"`
	TokenOutDenom string `json:"tokenOutDenom"`
},
) ([]poolmanagertypes.SwapAmountInRoute, error) {
	if err := ValidateRoutes(routes, "routes"); err != nil {
		return nil, err
	}
	result := make([]poolmanagertypes.SwapAmountInRoute, len(routes))
	for i, route := range routes {
		result[i] = ConvertSwapRouteFromEVM(route)
	}
	return result, nil
}

// ConvertAffiliateFromEVM converts a Solidity Affiliate struct to poolmanagertypes.Affiliate
func ConvertAffiliateFromEVM(affiliate struct {
	Address        common.Address `json:"address"`
	BasisPointsFee *big.Int       `json:"basisPointsFee"`
},
) (poolmanagertypes.Affiliate, error) {
	if err := ValidateAddress(affiliate.Address, "affiliate.address"); err != nil {
		return poolmanagertypes.Affiliate{}, err
	}
	if affiliate.BasisPointsFee == nil {
		return poolmanagertypes.Affiliate{}, ErrInvalidInput("affiliate.basisPointsFee cannot be nil")
	}
	// Basis points should be between 0 and 10000
	maxBasisPoints := big.NewInt(10000)
	if affiliate.BasisPointsFee.Cmp(maxBasisPoints) > 0 {
		return poolmanagertypes.Affiliate{}, ErrInvalidInput(fmt.Sprintf("affiliate.basisPointsFee (%s) exceeds maximum (10000)", affiliate.BasisPointsFee.String()))
	}

	cosmosAddr := sdk.AccAddress(affiliate.Address.Bytes()).String()
	return poolmanagertypes.Affiliate{
		Address:        cosmosAddr,
		BasisPointsFee: fmt.Sprintf("%d", affiliate.BasisPointsFee.Uint64()),
	}, nil
}

// ConvertAffiliatesFromEVM converts an array of Solidity Affiliate structs
func ConvertAffiliatesFromEVM(affiliates []struct {
	Address        common.Address `json:"address"`
	BasisPointsFee *big.Int       `json:"basisPointsFee"`
},
) ([]poolmanagertypes.Affiliate, error) {
	if len(affiliates) == 0 {
		return []poolmanagertypes.Affiliate{}, nil
	}
	if err := ValidateAffiliates(affiliates, "affiliates"); err != nil {
		return nil, err
	}
	result := make([]poolmanagertypes.Affiliate, len(affiliates))
	for i, affiliate := range affiliates {
		converted, err := ConvertAffiliateFromEVM(affiliate)
		if err != nil {
			return nil, fmt.Errorf("affiliates[%d]: %w", i, err)
		}
		result[i] = converted
	}
	return result, nil
}

// ConvertIBCTransferInfoFromEVM converts a Solidity IBCTransferInfo struct to gammtypes.IBCTransferInfo
func ConvertIBCTransferInfoFromEVM(ibcInfo struct {
	SourceChannel    string `json:"sourceChannel"`
	Receiver         string `json:"receiver"`
	Memo             string `json:"memo"`
	TimeoutTimestamp uint64 `json:"timeoutTimestamp"`
},
) (gammtypes.IBCTransferInfo, error) {
	// Note: timestamp validation is done in the handler where we have access to ctx
	// This conversion function just converts the struct format
	_ = ibcInfo // ibcInfo is already validated in the handler
	return gammtypes.IBCTransferInfo{
		SourceChannel:    ibcInfo.SourceChannel,
		Receiver:         ibcInfo.Receiver,
		Memo:             ibcInfo.Memo,
		TimeoutTimestamp: ibcInfo.TimeoutTimestamp,
	}, nil
}

// ConvertShareAmount converts *big.Int to sdkmath.Int for share amounts
func ConvertShareAmount(amount *big.Int) (sdkmath.Int, error) {
	if err := ValidateShareAmount(amount, "shareAmount"); err != nil {
		return sdkmath.Int{}, err
	}
	return sdkmath.NewIntFromBigInt(amount), nil
}

// ConvertCoinsToABIFormat converts an array of Coin structs to ABI-compatible format for packing
// The ABI library expects []interface{} where each element is []interface{} containing [denom, amount]
// Accepts structs with or without json tags
func ConvertCoinsToABIFormat(coins interface{}) []interface{} {
	rv := reflect.ValueOf(coins)
	if rv.Kind() != reflect.Slice {
		return nil
	}
	
	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		// Extract denom and amount fields
		denom := elem.FieldByName("Denom").Interface()
		amount := elem.FieldByName("Amount").Interface()
		// Each tuple element must be []interface{} with fields in ABI order
		result[i] = []interface{}{
			denom,
			amount,
		}
	}
	return result
}

// ConvertRoutesToABIFormat converts an array of SwapAmountInRoute structs to ABI-compatible format
func ConvertRoutesToABIFormat(routes []struct {
	PoolId        uint64 `json:"poolId"`
	TokenOutDenom string `json:"tokenOutDenom"`
},
) []interface{} {
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		// Each tuple element must be []interface{} with fields in ABI order
		result[i] = []interface{}{
			route.PoolId,
			route.TokenOutDenom,
		}
	}
	return result
}

// ConvertAffiliatesToABIFormat converts an array of Affiliate structs to ABI-compatible format
func ConvertAffiliatesToABIFormat(affiliates []struct {
	Address        common.Address `json:"address"`
	BasisPointsFee *big.Int       `json:"basisPointsFee"`
},
) []interface{} {
	result := make([]interface{}, len(affiliates))
	for i, affiliate := range affiliates {
		// Each tuple element must be []interface{} with fields in ABI order
		result[i] = []interface{}{
			affiliate.Address,
			affiliate.BasisPointsFee,
		}
	}
	return result
}
