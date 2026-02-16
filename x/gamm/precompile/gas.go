package gamm

const (
	// Base gas costs for transactions
	// IMPORTANT: These values represent the MINIMUM gas required for precompile operations.
	// The cosmos/evm RunNativeAction creates a Cosmos gas meter with limit = contract.Gas.
	// These base values should account for:
	//   - BalancerGasFeeForSwap = 10,000 (per swap)
	//   - UniversalRemoveOverlaps = 500 (per call, can be many during token transfers)
	//   - Bank module operations (transfers, balance checks)
	//   - State reads/writes (KV store operations)
	//
	// If these values are too low, gas estimation will succeed but actual execution
	// may fail with "out of gas" errors that produce missing revert data.
	GasJoinPoolBase                         = 150_000 // Join requires pool lookup + token transfers + share minting
	GasExitPoolBase                         = 150_000 // Exit requires pool lookup + share burning + token transfers
	GasSwapExactAmountInBase                = 150_000 // Swap requires: balancer computation (10k) + token transfers + permission checks
	GasCreatePoolBase                       = 200_000 // Pool creation is complex: multiple state writes + token deposits
	GasSwapExactAmountInWithIBCTransferBase = 200_000 // Swap + IBC transfer overhead

	// Gas costs per element for dynamic calculations
	GasPerRoute     = 50_000 // Each additional route adds another swap operation
	GasPerCoin      = 10_000 // Each coin involves balance lookups and potential permission checks
	GasPerAffiliate = 5_000  // Affiliate processing
	GasPerMemoByte  = 10

	// Gas costs for queries (read-only operations are cheaper)
	GasGetPoolBase                     = 10_000
	GasGetPoolsBase                    = 20_000
	GasGetPoolTypeBase                 = 5_000
	GasCalcJoinPoolNoSwapSharesBase    = 20_000
	GasCalcExitPoolCoinsFromSharesBase = 20_000
	GasCalcJoinPoolSharesBase          = 20_000
	GasGetPoolParamsBase               = 10_000
	GasGetTotalSharesBase              = 10_000
	GasGetTotalLiquidityBase           = 20_000
	GasEstimateSwapExactAmountInBase   = 50_000 // Estimation runs full swap logic without state changes
	GasEstimateSwapExactAmountOutBase  = 50_000
)

// CalculateDynamicGas calculates dynamic gas based on input complexity
func CalculateDynamicGas(baseGas uint64, numRoutes, numCoins, numAffiliates int) uint64 {
	gas := baseGas
	gas += uint64(numRoutes) * GasPerRoute
	gas += uint64(numCoins) * GasPerCoin
	gas += uint64(numAffiliates) * GasPerAffiliate
	return gas
}
