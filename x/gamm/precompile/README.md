# Gamm Precompile

This package implements a precompiled contract for the BitBadges gamm module, enabling Solidity smart contracts to interact with liquidity pools.

## Precompile Address

```
0x0000000000000000000000000000000000001002
```

## Files

- `precompile.go`: Main precompile implementation
- `handlers.go`: Transaction and query handler methods
- `validation.go`: Input validation helpers
- `errors.go`: Structured error types
- `events.go`: Event emission helpers
- `gas.go`: Dynamic gas calculation
- `security.go`: Security utilities
- `conversions.go`: Type conversion helpers (EVM <-> Cosmos)
- `metrics.go`: Usage metrics tracking
- `abi.json`: ABI definition for the precompile
- `precompile_test.go`: Unit tests

## Usage

### In Solidity

```solidity
import "./interfaces/IGammPrecompile.sol";

IGammPrecompile precompile = IGammPrecompile(0x0000000000000000000000000000000000001002);

// Join a pool
Coin[] memory tokenInMaxs = new Coin[](2);
tokenInMaxs[0] = Coin({denom: "uatom", amount: 1000000});
tokenInMaxs[1] = Coin({denom: "uosmo", amount: 2000000});

(uint256 shareOutAmount, Coin[] memory tokenIn) = precompile.joinPool(
    poolId,
    desiredShareAmount,
    tokenInMaxs
);

// Exit a pool
Coin[] memory tokenOutMins = new Coin[](2);
tokenOutMins[0] = Coin({denom: "uatom", amount: 500000});
tokenOutMins[1] = Coin({denom: "uosmo", amount: 1000000});

Coin[] memory tokenOut = precompile.exitPool(
    poolId,
    shareInAmount,
    tokenOutMins
);

// Swap exact amount in
SwapAmountInRoute[] memory routes = new SwapAmountInRoute[](1);
routes[0] = SwapAmountInRoute({
    poolId: poolId,
    tokenOutDenom: "uosmo"
});

Coin memory tokenIn = Coin({denom: "uatom", amount: 1000000});
Affiliate[] memory affiliates = new Affiliate[](0);

uint256 tokenOutAmount = precompile.swapExactAmountIn(
    routes,
    tokenIn,
    minTokenOutAmount,
    affiliates
);
```

## Transaction Methods

### joinPool
Join a liquidity pool by providing tokens.

**Parameters:**
- `poolId` (uint64): The pool ID to join
- `shareOutAmount` (uint256): Desired amount of LP shares to receive
- `tokenInMaxs` (Coin[]): Maximum amounts of each token to provide

**Returns:**
- `shareOutAmount` (uint256): Actual shares received
- `tokenIn` (Coin[]): Actual tokens provided

### exitPool
Exit a liquidity pool by burning shares.

**Parameters:**
- `poolId` (uint64): The pool ID to exit
- `shareInAmount` (uint256): Amount of LP shares to burn
- `tokenOutMins` (Coin[]): Minimum amounts of each token to receive

**Returns:**
- `tokenOut` (Coin[]): Tokens received

### swapExactAmountIn
Swap tokens with exact input amount.

**Parameters:**
- `routes` (SwapAmountInRoute[]): Swap route through pools
- `tokenIn` (Coin): Input token
- `tokenOutMinAmount` (uint256): Minimum output amount
- `affiliates` (Affiliate[]): Optional affiliate fee recipients

**Returns:**
- `tokenOutAmount` (uint256): Output token amount

### swapExactAmountInWithIBCTransfer
Swap tokens and transfer via IBC.

**Parameters:**
- `routes` (SwapAmountInRoute[]): Swap route through pools
- `tokenIn` (Coin): Input token
- `tokenOutMinAmount` (uint256): Minimum output amount
- `ibcTransferInfo` (IBCTransferInfo): IBC transfer details
- `affiliates` (Affiliate[]): Optional affiliate fee recipients

**Returns:**
- `tokenOutAmount` (uint256): Output token amount

## Query Methods

### getPool
Query pool data by ID.

**Parameters:**
- `poolId` (uint64): The pool ID

**Returns:**
- `pool` (bytes): Protobuf-encoded pool data

### getPools
Query all pools with pagination.

**Parameters:**
- `offset` (uint256): Pagination offset
- `limit` (uint256): Pagination limit (max 1000)

**Returns:**
- `pools` (bytes): Protobuf-encoded pools data

### getPoolType
Query pool type by ID.

**Parameters:**
- `poolId` (uint64): The pool ID

**Returns:**
- `poolType` (string): Pool type name

### calcJoinPoolNoSwapShares
Calculate shares for joining pool without swap.

**Parameters:**
- `poolId` (uint64): The pool ID
- `tokensIn` (Coin[]): Tokens to provide

**Returns:**
- `tokensOut` (Coin[]): Tokens that will be used
- `sharesOut` (uint256): Shares that will be received

### calcExitPoolCoinsFromShares
Calculate tokens received for exiting pool.

**Parameters:**
- `poolId` (uint64): The pool ID
- `shareInAmount` (uint256): Shares to burn

**Returns:**
- `tokensOut` (Coin[]): Tokens that will be received

### calcJoinPoolShares
Calculate shares for joining pool.

**Parameters:**
- `poolId` (uint64): The pool ID
- `tokensIn` (Coin[]): Tokens to provide

**Returns:**
- `shareOutAmount` (uint256): Shares that will be received
- `tokensOut` (Coin[]): Tokens that will be used

### getPoolParams
Query pool parameters.

**Parameters:**
- `poolId` (uint64): The pool ID

**Returns:**
- `params` (bytes): Protobuf-encoded pool parameters

### getTotalShares
Query total shares for a pool.

**Parameters:**
- `poolId` (uint64): The pool ID

**Returns:**
- `totalShares` (Coin): Total shares as a Coin struct

### getTotalLiquidity
Query total liquidity across all pools.

**Parameters:** None

**Returns:**
- `liquidity` (Coin[]): Total liquidity for all pools

## ABI Notes

### Struct Definitions

The ABI includes the following struct types:

- `Coin { string denom, uint256 amount }` - Token representation
- `SwapAmountInRoute { uint64 poolId, string tokenOutDenom }` - Swap route step
- `Affiliate { address address, uint256 basisPointsFee }` - Affiliate fee recipient
- `IBCTransferInfo { string sourceChannel, string receiver, string memo, uint64 timeoutTimestamp }` - IBC transfer details

### Return Type Simplifications

Complex protobuf types are returned as bytes for full access. Simple types (uint256, string) are returned directly.

## DoS Protection Limits

| Field | Max Size |
|-------|----------|
| Routes | 10 |
| Coins | 20 |
| Affiliates | 10 |
| String Length | 10,000 chars |

## Production Deployment

### Genesis Configuration

The gamm precompile must be **enabled** in the genesis state for it to be callable. The precompile is registered during app initialization, but it must also be enabled via the EVM module's genesis state.

**Required Genesis Configuration:**

```json
{
  "evm": {
    "params": {
      "active_static_precompiles": [
        "0x0000000000000000000000000000000000001002"
      ]
    }
  }
}
```

### Verification Steps

1. **Check ABI Loading**: Verify the precompile ABI loaded successfully:
   ```go
   err := gamm.ValidatePrecompileEnabled()
   if err != nil {
       // Precompile is not properly configured
   }
   ```

2. **Verify Precompile Address**: Ensure the precompile address matches:
   ```
   0x0000000000000000000000000000000000001002
   ```

3. **Check Genesis State**: Verify the precompile address is in `active_static_precompiles` in the EVM module's genesis state.

4. **Test Precompile Access**: Query the precompile to ensure it's accessible:
   ```bash
   # Query a pool to verify precompile is working
   # This should return pool data, not an error
   ```

### Monitoring Recommendations

1. **Event Monitoring**: Monitor precompile events:
   - `precompile_join_pool`
   - `precompile_exit_pool`
   - `precompile_swap_exact_amount_in`
   - `precompile_swap_exact_amount_in_with_ibc_transfer`
   - `precompile_metrics`

2. **Error Tracking**: Monitor error codes and frequencies:
   - Track `ErrorCode` values in logs
   - Alert on high error rates
   - Monitor specific error types (e.g., `ErrorCodeSwapFailed`)

3. **Gas Usage**: Monitor gas consumption:
   - Track gas usage per method
   - Alert on unusually high gas costs
   - Monitor dynamic gas calculation accuracy

4. **Performance Metrics**: Track:
   - Method call frequency
   - Success/failure rates
   - Average gas consumption per operation type

### Troubleshooting

#### Precompile Not Callable

**Symptoms**: Calls to precompile address return errors or no data.

**Solutions**:
1. Verify precompile is enabled in genesis state (`active_static_precompiles`)
2. Check ABI loaded successfully: `gamm.GetABILoadError()`
3. Verify precompile address is correct: `0x0000000000000000000000000000000000001002`
4. Check EVM module logs for precompile registration errors

#### ABI Load Errors

**Symptoms**: Precompile returns "ABI failed to load" errors.

**Solutions**:
1. Verify `abi.json` file is present and valid JSON
2. Check file is embedded correctly (build issue)
3. Verify ABI matches the precompile implementation
4. Check application logs for ABI loading warnings

#### Snapshot Errors (Fixed)

**Symptoms**: Panics with "snapshot index 0 out of bound [0..0)" when precompile returns errors.

**Status**: This issue has been fixed in `pkg/evmcompat/atomic.go`. The fix detects EVM context and uses native `Snapshot()`/`RevertToSnapshot()` instead of `CacheContext()`, which prevents the snapshot stack corruption. See `docs/security/SNAPSHOT_CORRUPTION_REPORT.md` for technical details.

#### Gas Estimation Issues

**Symptoms**: Gas estimation is inaccurate or too low.

**Solutions**:
1. Dynamic gas calculation is implemented in `RequiredGas()` - verify it's working
2. For complex operations, gas may be underestimated if ABI parsing fails (falls back to base gas)
3. Monitor actual gas usage vs estimated gas
4. Adjust base gas costs if needed based on production metrics

#### Type Conversion Errors

**Symptoms**: "invalid type" errors when calling precompile methods.

**Solutions**:
1. Verify Solidity contract uses correct struct definitions matching the ABI
2. Check ABI encoding/decoding in Solidity contract
3. Review error messages - they now include method names and actual types received
4. Ensure struct field names match exactly (case-sensitive)

## Known Limitations

1. **Pool Creation**: Pool creation is not supported through the precompile. Use native Cosmos SDK interface.

2. **Complex Pool Types**: Some pool-specific operations may require direct keeper access.

3. **ABI Load Failure**: If `abi.json` is corrupted, the precompile will be disabled but the chain will still start. The error can be checked via `gamm.GetABILoadError()`.

4. **EVM Snapshot Bug (Fixed)**: The `cosmos/evm` module's snapshot stack corruption issue has been fixed locally in `pkg/evmcompat/atomic.go`. The fix detects EVM context and uses native `Snapshot()`/`RevertToSnapshot()` instead of `CacheContext()`. See `docs/security/SNAPSHOT_CORRUPTION_REPORT.md` for details.

5. **Dynamic Gas Estimation**: While dynamic gas calculation is implemented, if ABI parsing fails during `RequiredGas()`, the method falls back to base gas. This may result in underestimation for complex operations. The EVM will still charge actual gas used during execution.

## Production Checklist

Before deploying to production:

- [ ] Precompile address is in `active_static_precompiles` in genesis state
- [ ] ABI loaded successfully (check `gamm.GetABILoadError()`)
- [ ] Precompile validation passes (`gamm.ValidatePrecompileEnabled()`)
- [ ] All required methods are accessible (test queries)
- [ ] Monitoring is set up for precompile events and errors
- [ ] Gas costs are calibrated based on testnet usage
- [ ] Error handling is tested
- [ ] Documentation is reviewed and up-to-date
- [ ] Known limitations are documented and understood

## Development

### Adding New Methods

1. Add method to `abi.json`
2. Add method constant to `precompile.go`
3. Add gas cost constant to `gas.go`
4. Add case to `RequiredGas` switch in `precompile.go`
5. Add case to `ExecuteWithMethodName` switch in `precompile.go`
6. Implement method handler in `handlers.go`
7. Add validation in `validation.go` (if needed)
8. Add conversion helpers in `conversions.go` (if needed)
9. Add event emission in `events.go` (for transactions)
10. Add tests:
    - Unit tests for validation and core logic
    - Integration tests for EVM keeper interaction
    - Benchmarks for gas calculation performance
    - Fuzz tests for validation functions

### Testing

The precompile includes comprehensive test coverage:

- **Unit Tests** (`test/unit/`): Core functionality, validation, conversions, error handling, security
- **Integration Tests** (`test/integration/`): EVM keeper integration, events, gas accuracy
- **E2E Tests** (`test/e2e/`): Complete workflows (note: some ABI packing limitations exist)
- **Advanced Tests** (`test/advanced/`): Concurrency, reentrancy, stress, multi-user workflows
- **Benchmarks** (`test/benchmarks/`): Gas calculation performance
- **Fuzz Tests** (`test/fuzz/`): Robustness testing for validation functions

Run all tests:
```bash
go test -tags=test ./x/gamm/precompile/...
```

Run benchmarks:
```bash
go test -tags=test -bench=. ./x/gamm/precompile/test/benchmarks/...
```

Run fuzz tests:
```bash
go test -tags=test -fuzz=. ./x/gamm/precompile/test/fuzz/...
```

#### Test Coverage Summary

- **Unit Tests**: Validation, conversions, error handling, security patterns, gas calculation
- **Integration Tests**: Full EVM keeper integration, event emission, gas accuracy verification
- **E2E Tests**: Complete workflows for join/exit/swap operations (some ABI tuple packing limitations)
- **Advanced Tests**: 
  - Concurrency: Parallel queries, state consistency, validation function thread-safety
  - Reentrancy: Call stack depth, nested calls, state consistency
  - Stress: Maximum routes/coins/affiliates, large amounts, complex routing
  - Multi-user: Sequential operations, concurrent queries, state isolation
- **Benchmarks**: Gas calculation performance for various input sizes
- **Fuzz Tests**: Edge case discovery for all validation functions

#### Known Limitations

1. **E2E ABI Packing**: Some E2E tests have complex ABI tuple packing issues for arrays of structs. The core functionality is fully tested through unit, integration, and advanced test suites. E2E ABI packing improvements can be addressed when better ABI handling tooling is available.

## Error Codes

- `ErrorCodeInvalidInput` (1): Invalid input parameters
- `ErrorCodePoolNotFound` (2): Pool not found
- `ErrorCodeSwapFailed` (3): Swap operation failed
- `ErrorCodeQueryFailed` (4): Query operation failed
- `ErrorCodeInternalError` (5): Internal error
- `ErrorCodeUnauthorized` (6): Unauthorized operation
- `ErrorCodeJoinPoolFailed` (7): Join pool operation failed
- `ErrorCodeExitPoolFailed` (8): Exit pool operation failed
- `ErrorCodeIBCTransferFailed` (9): IBC transfer operation failed

