package gamm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

const (
	// Base gas costs for transactions
	// IMPORTANT: These values are DEDUCTED from the transaction gas before the precompile runs.
	// The actual execution gas comes from the remaining gas (contract.Gas after deduction).
	// Setting these too high causes "out of gas" errors because there's not enough remaining
	// gas for the Cosmos SDK operations (bank transfers, state updates, etc.).
	// These values should be minimal "entry fees" - the actual gas consumption happens inside
	// RunNativeAction via the Cosmos SDK gas meter.
	GasJoinPoolBase                         = 10_000
	GasExitPoolBase                         = 10_000
	GasSwapExactAmountInBase                = 10_000
	GasCreatePoolBase                       = 15_000
	GasSwapExactAmountInWithIBCTransferBase = 15_000

	// Gas costs per element for dynamic calculations
	GasPerRoute      = 5_000
	GasPerCoin       = 2_000
	GasPerAffiliate  = 3_000
	GasPerMemoByte   = 10
	GasPerInputChunk = 100 // Gas per 32-byte chunk of input, for JSON parsing

	// Gas costs for queries (lower since they're read-only)
	GasGetPoolBase                     = 3_000
	GasGetPoolsBase                    = 5_000
	GasGetPoolTypeBase                 = 2_000
	GasCalcJoinPoolNoSwapSharesBase    = 5_000
	GasCalcExitPoolCoinsFromSharesBase = 5_000
	GasCalcJoinPoolSharesBase          = 5_000
	GasGetPoolParamsBase               = 3_000
	GasGetTotalSharesBase              = 3_000
	GasGetTotalLiquidityBase           = 5_000
	GasEstimateSwapExactAmountInBase   = 10_000
	GasEstimateSwapExactAmountOutBase  = 10_000
)

// CalculateDynamicGas calculates dynamic gas based on input complexity
func CalculateDynamicGas(baseGas uint64, numRoutes, numCoins, numAffiliates int) uint64 {
	gas := baseGas
	gas += uint64(numRoutes) * GasPerRoute
	gas += uint64(numCoins) * GasPerCoin
	gas += uint64(numAffiliates) * GasPerAffiliate
	return gas
}

// MeterMessage validates the size of a parsed message's variable-length parts
// against the caps in security.go and charges their per-element gas on the
// SDK meter, which RunNativeAction bills back to the EVM.
//
// RequiredGas runs before the JSON is parsed and can only price the raw
// input; this is the other half, taken once the message shape is known and
// before the keeper does any work.
func MeterMessage(ctx sdk.Context, msg sdk.Msg) error {
	gas, err := messageGas(msg)
	if err != nil {
		return err
	}
	ctx.GasMeter().ConsumeGas(gas, "precompile message elements")
	return nil
}

func messageGas(msg sdk.Msg) (uint64, error) {
	switch m := msg.(type) {
	case *gammtypes.MsgJoinPool:
		if err := ValidateArraySizeAllowEmpty(len(m.TokenInMaxs), MaxCoins, "tokenInMaxs"); err != nil {
			return 0, err
		}
		return CalculateDynamicGas(0, 0, len(m.TokenInMaxs), 0), nil
	case *gammtypes.MsgExitPool:
		if err := ValidateArraySizeAllowEmpty(len(m.TokenOutMins), MaxCoins, "tokenOutMins"); err != nil {
			return 0, err
		}
		return CalculateDynamicGas(0, 0, len(m.TokenOutMins), 0), nil
	case *gammtypes.MsgSwapExactAmountIn:
		return swapGas(len(m.Routes), len(m.Affiliates), "")
	case *gammtypes.MsgSwapExactAmountInWithIBCTransfer:
		return swapGas(len(m.Routes), len(m.Affiliates), m.IbcTransferInfo.Memo)
	case *balancer.MsgCreateBalancerPool:
		if err := ValidateArraySizeAllowEmpty(len(m.PoolAssets), MaxCoins, "poolAssets"); err != nil {
			return 0, err
		}
		return CalculateDynamicGas(0, 0, len(m.PoolAssets), 0), nil
	}
	return 0, nil
}

func swapGas(routes, affiliates int, memo string) (uint64, error) {
	if err := ValidateArraySize(routes, MaxRoutes, "routes"); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(affiliates, MaxAffiliates, "affiliates"); err != nil {
		return 0, err
	}
	if len(memo) > MaxMemoLength {
		return 0, ErrInvalidInput(fmt.Sprintf("ibcTransferInfo.memo length (%d) exceeds maximum allowed length (%d)", len(memo), MaxMemoLength))
	}
	return CalculateDynamicGas(0, routes, 1, affiliates) + uint64(len(memo))*GasPerMemoByte, nil
}
