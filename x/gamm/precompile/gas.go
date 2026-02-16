package gamm

const (
	// Base gas costs for transactions
	// IMPORTANT: These values are DEDUCTED from the transaction gas before the precompile runs.
	// The actual execution gas comes from the remaining gas (contract.Gas after deduction).
	// Setting these too high causes "out of gas" errors because there's not enough remaining
	// gas for the Cosmos SDK operations (bank transfers, state updates, etc.).
	// These values should be minimal "entry fees" - the actual gas consumption happens inside
	// RunNativeAction via the Cosmos SDK gas meter.
	GasJoinPoolBase                         = 20_000
	GasExitPoolBase                         = 20_000
	GasSwapExactAmountInBase                = 20_000 // Base entry fee; actual gas (balancer computation 10k, bank transfers, etc.) consumed from contract.Gas
	GasCreatePoolBase                       = 30_000 // Higher gas for pool creation due to complexity
	GasSwapExactAmountInWithIBCTransferBase = 25_000 // Higher than regular swap due to IBC transfer overhead

	// Gas costs per element for dynamic calculations
	GasPerRoute     = 5_000
	GasPerCoin      = 2_000
	GasPerAffiliate = 3_000
	GasPerMemoByte  = 10

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
