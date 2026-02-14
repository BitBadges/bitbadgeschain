package gamm

const (
	// Base gas costs for transactions
	GasJoinPoolBase                      = 40_000
	GasExitPoolBase                      = 40_000
	GasSwapExactAmountInBase             = 50_000
	GasCreatePoolBase                    = 60_000 // Higher gas for pool creation due to complexity
	GasSwapExactAmountInWithIBCTransferBase = 60_000

	// Gas costs per element for dynamic calculations
	GasPerRoute     = 5_000
	GasPerCoin      = 2_000
	GasPerAffiliate = 3_000
	GasPerMemoByte  = 10

	// Gas costs for queries (lower since they're read-only)
	GasGetPoolBase                    = 3_000
	GasGetPoolsBase                   = 5_000
	GasGetPoolTypeBase                = 2_000
	GasCalcJoinPoolNoSwapSharesBase   = 5_000
	GasCalcExitPoolCoinsFromSharesBase = 5_000
	GasCalcJoinPoolSharesBase         = 5_000
	GasGetPoolParamsBase              = 3_000
	GasGetTotalSharesBase             = 3_000
	GasGetTotalLiquidityBase          = 5_000
	GasEstimateSwapExactAmountInBase  = 10_000
	GasEstimateSwapExactAmountOutBase = 10_000
)

// CalculateDynamicGas calculates dynamic gas based on input complexity
func CalculateDynamicGas(baseGas uint64, numRoutes, numCoins, numAffiliates int) uint64 {
	gas := baseGas
	gas += uint64(numRoutes) * GasPerRoute
	gas += uint64(numCoins) * GasPerCoin
	gas += uint64(numAffiliates) * GasPerAffiliate
	return gas
}

