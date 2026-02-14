# Gamm Precompile - Frontend Integration Guide

## Overview

The Gamm Precompile enables Solidity smart contracts to interact with liquidity pools through a standardized EVM interface. It provides both transaction methods (state-changing operations) and query methods (read-only operations) for managing liquidity pools, swaps, and pool information.

**Key Features:**
- ✅ Join liquidity pools by providing tokens
- ✅ Exit liquidity pools by burning shares
- ✅ Swap tokens with exact input amounts
- ✅ Swap tokens with IBC transfers
- ✅ Query pool information, liquidity, and calculations
- ✅ All operations use JSON-based protobuf format

## Precompile Address

```
0x0000000000000000000000000000000000001002
```

**Address Space Convention:**
- **0x0800-0x0806**: Reserved for default Cosmos precompiles
- **0x1001+**: Reserved for custom BitBadges precompiles
  - 0x1001: Tokenization precompile
  - 0x1002: Gamm precompile ✅ (This document)
  - 0x1003: SendManager precompile
  - 0x1004+: Available for future precompiles

## Registration & Enablement

The gamm precompile is **registered** in `app/evm.go` and **enabled** in the genesis configuration (`config.yml`). It's also enabled via the upgrade handler for existing chains.

## ABI Interface

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./types/GammTypes.sol";

/// @title IGammPrecompile
/// @notice Interface for the BitBadges gamm precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000001002
///      All methods use JSON string parameters matching protobuf JSON format.
///      The caller address (sender) is automatically set from msg.sender.
///      Use helper libraries to construct JSON strings from Solidity types.
interface IGammPrecompile {
    // ============ Transactions ============
    
    /// @notice Join a liquidity pool by providing tokens
    /// @param msgJson JSON string matching MsgJoinPool protobuf JSON format
    /// @return shareOutAmount The amount of pool shares received
    /// @return tokenIn The actual tokens provided (may be less than tokenInMaxs)
    function joinPool(string memory msgJson)
        external
        returns (uint256 shareOutAmount, GammTypes.Coin[] memory tokenIn);

    /// @notice Exit a liquidity pool by burning shares
    /// @param msgJson JSON string matching MsgExitPool protobuf JSON format
    /// @return tokenOut The tokens received from exiting the pool
    function exitPool(string memory msgJson)
        external
        returns (GammTypes.Coin[] memory tokenOut);

    /// @notice Swap tokens with exact input amount
    /// @param msgJson JSON string matching MsgSwapExactAmountIn protobuf JSON format
    /// @return tokenOutAmount The amount of output tokens received
    function swapExactAmountIn(string memory msgJson)
        external
        returns (uint256 tokenOutAmount);

    /// @notice Swap tokens with exact input amount and transfer via IBC
    /// @param msgJson JSON string matching MsgSwapExactAmountInWithIBCTransfer protobuf JSON format
    /// @return tokenOutAmount The amount of output tokens received
    function swapExactAmountInWithIBCTransfer(string memory msgJson)
        external
        returns (uint256 tokenOutAmount);

    // ============ Queries ============

    /// @notice Get pool data by ID
    /// @param msgJson JSON string matching QueryPoolRequest protobuf JSON format
    /// @return pool The pool data as protobuf-encoded bytes
    function getPool(string memory msgJson)
        external
        view
        returns (bytes memory pool);

    /// @notice Get all pools with pagination
    /// @param msgJson JSON string matching QueryPoolsRequest protobuf JSON format
    /// @return pools The pools data as protobuf-encoded bytes
    function getPools(string memory msgJson)
        external
        view
        returns (bytes memory pools);

    /// @notice Get pool type by ID
    /// @param msgJson JSON string matching QueryPoolTypeRequest protobuf JSON format
    /// @return poolType The pool type string
    function getPoolType(string memory msgJson)
        external
        view
        returns (string memory poolType);

    /// @notice Calculate shares for joining pool without swap
    /// @param msgJson JSON string matching QueryCalcJoinPoolNoSwapSharesRequest protobuf JSON format
    /// @return tokensOut The tokens that would be provided
    /// @return sharesOut The shares that would be received
    function calcJoinPoolNoSwapShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory tokensOut, uint256 sharesOut);

    /// @notice Calculate tokens received for exiting pool
    /// @param msgJson JSON string matching QueryCalcExitPoolCoinsFromSharesRequest protobuf JSON format
    /// @return tokensOut The tokens that would be received
    function calcExitPoolCoinsFromShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory tokensOut);

    /// @notice Calculate shares for joining pool (with swap)
    /// @param msgJson JSON string matching QueryCalcJoinPoolSharesRequest protobuf JSON format
    /// @return shareOutAmount The shares that would be received
    /// @return tokensOut The tokens that would be provided
    function calcJoinPoolShares(string memory msgJson)
        external
        view
        returns (uint256 shareOutAmount, GammTypes.Coin[] memory tokensOut);

    /// @notice Get pool parameters
    /// @param msgJson JSON string matching QueryPoolParamsRequest protobuf JSON format
    /// @return params The pool parameters as protobuf-encoded bytes
    function getPoolParams(string memory msgJson)
        external
        view
        returns (bytes memory params);

    /// @notice Get total shares for a pool
    /// @param msgJson JSON string matching QueryTotalSharesRequest protobuf JSON format
    /// @return totalShares The total shares as a Coin struct
    function getTotalShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin memory totalShares);

    /// @notice Get total liquidity for a pool
    /// @param msgJson JSON string matching QueryTotalLiquidityRequest protobuf JSON format
    /// @return liquidity The total liquidity as an array of Coin structs
    function getTotalLiquidity(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory liquidity);
}
```

## Transaction Methods

### 1. joinPool

Join a liquidity pool by providing tokens. Returns the shares received and actual tokens used.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1,
  "share_out_amount": "1000000",
  "token_in_maxs": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    },
    {
      "denom": "ustake",
      "amount": "2000000000"
    }
  ]
}
```

**Important Notes:**
- `sender` is **automatically set** from `msg.sender` (cannot be spoofed)
- `share_out_amount` is the desired number of pool shares to receive
- `token_in_maxs` is an array of maximum tokens you're willing to provide
- Actual tokens used may be less than `token_in_maxs`

**TypeScript Example:**
```typescript
import { ethers } from "ethers";

const GAMM_PRECOMPILE_ADDRESS = "0x0000000000000000000000000000000000001002";

const gammABI = [
  "function joinPool(string memory msgJson) external returns (uint256 shareOutAmount, tuple(string denom, uint256 amount)[] tokenIn)"
];

const gamm = new ethers.Contract(GAMM_PRECOMPILE_ADDRESS, gammABI, signer);

const msgJson = JSON.stringify({
      pool_id: 1,
  share_out_amount: "1000000",
  token_in_maxs: [
    { denom: "ubadge", amount: "1000000000" },
    { denom: "ustake", amount: "2000000000" }
  ]
});

const tx = await gamm.joinPool(msgJson);
const receipt = await tx.wait();

// Decode return values from events or logs
const [shareOutAmount, tokenIn] = await gamm.joinPool.staticCall(msgJson);
console.log("Shares received:", shareOutAmount.toString());
console.log("Tokens used:", tokenIn);
```

### 2. exitPool

Exit a liquidity pool by burning shares. Returns the tokens received.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1,
  "share_in_amount": "500000",
  "token_out_mins": [
    {
      "denom": "ubadge",
      "amount": "500000000"
    },
    {
      "denom": "ustake",
      "amount": "1000000000"
    }
  ]
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
      pool_id: 1,
  share_in_amount: "500000",
  token_out_mins: [
    { denom: "ubadge", amount: "500000000" },
    { denom: "ustake", amount: "1000000000" }
  ]
});

const tx = await gamm.exitPool(msgJson);
const receipt = await tx.wait();

// Get return value
const tokenOut = await gamm.exitPool.staticCall(msgJson);
console.log("Tokens received:", tokenOut);
```

### 3. swapExactAmountIn

Swap tokens with exact input amount. Returns the output token amount.

**Protobuf JSON Format:**
```json
{
  "routes": [
    {
      "pool_id": 1,
      "token_out_denom": "ustake"
    }
  ],
  "token_in": {
    "denom": "ubadge",
    "amount": "1000000000"
  },
  "token_out_min_amount": "2000000000",
  "affiliates": []
}
```

**Multi-Hop Swap Example:**
```json
{
  "routes": [
    {
      "pool_id": 1,
      "token_out_denom": "uatom"
    },
    {
      "pool_id": 2,
      "token_out_denom": "ustake"
    }
  ],
  "token_in": {
    "denom": "ubadge",
    "amount": "1000000000"
  },
  "token_out_min_amount": "1800000000",
  "affiliates": [
    {
      "address": "bb1...",
      "basis_points_fee": "100"
    }
  ]
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
  routes: [
    { pool_id: 1, token_out_denom: "ustake" }
  ],
  token_in: {
    denom: "ubadge",
    amount: "1000000000"
  },
  token_out_min_amount: "2000000000",
  affiliates: []
});

const tx = await gamm.swapExactAmountIn(msgJson);
const receipt = await tx.wait();

// Get return value
const tokenOutAmount = await gamm.swapExactAmountIn.staticCall(msgJson);
console.log("Output amount:", tokenOutAmount.toString());
```

### 4. swapExactAmountInWithIBCTransfer

Swap tokens and transfer the output via IBC to another chain.

**Protobuf JSON Format:**
```json
{
  "routes": [
    {
      "pool_id": 1,
      "token_out_denom": "ustake"
    }
  ],
  "token_in": {
    "denom": "ubadge",
    "amount": "1000000000"
  },
  "token_out_min_amount": "2000000000",
  "ibc_transfer_info": {
    "source_channel": "channel-0",
    "receiver": "cosmos1...",
    "memo": "",
    "timeout_timestamp": "1735689600000000000"
  },
  "affiliates": []
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
  routes: [
    { pool_id: 1, token_out_denom: "ustake" }
  ],
  token_in: {
    denom: "ubadge",
    amount: "1000000000"
  },
  token_out_min_amount: "2000000000",
  ibc_transfer_info: {
    source_channel: "channel-0",
    receiver: "cosmos1...",
    memo: "",
    timeout_timestamp: "1735689600000000000" // Unix timestamp in nanoseconds
  },
  affiliates: []
});

const tx = await gamm.swapExactAmountInWithIBCTransfer(msgJson);
await tx.wait();
```

## Query Methods

### 1. getPool

Get pool data by ID. Returns protobuf-encoded bytes.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({ pool_id: 1 });
const poolBytes = await gamm.getPool(msgJson);
// Decode protobuf bytes to get pool data
```

### 2. getPools

Get all pools with pagination. Returns protobuf-encoded bytes.

**Protobuf JSON Format:**
```json
{
  "pagination": {
    "key": "",
    "offset": "0",
    "limit": "100",
    "count_total": true
  }
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
  pagination: {
    key: "",
    offset: "0",
    limit: "100",
    count_total: true
  }
});
const poolsBytes = await gamm.getPools(msgJson);
```

### 3. getPoolType

Get the type of a pool (e.g., "Balancer").

**Protobuf JSON Format:**
```json
{
  "pool_id": 1
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({ pool_id: 1 });
const poolType = await gamm.getPoolType(msgJson);
console.log("Pool type:", poolType); // e.g., "Balancer"
```

### 4. calcJoinPoolNoSwapShares

Calculate shares for joining a pool without swap.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1,
  "token_in_maxs": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    },
    {
      "denom": "ustake",
      "amount": "2000000000"
    }
  ]
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
      pool_id: 1,
  token_in_maxs: [
    { denom: "ubadge", amount: "1000000000" },
    { denom: "ustake", amount: "2000000000" }
  ]
});

const [tokensOut, sharesOut] = await gamm.calcJoinPoolNoSwapShares(msgJson);
console.log("Tokens out:", tokensOut);
console.log("Shares out:", sharesOut.toString());
```

### 5. calcExitPoolCoinsFromShares

Calculate tokens received for exiting a pool.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1,
  "share_in_amount": "500000"
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
      pool_id: 1,
  share_in_amount: "500000"
});

const tokensOut = await gamm.calcExitPoolCoinsFromShares(msgJson);
console.log("Tokens out:", tokensOut);
```

### 6. calcJoinPoolShares

Calculate shares for joining a pool (with swap).

**Protobuf JSON Format:**
```json
{
  "pool_id": 1,
  "token_in_maxs": [
    {
      "denom": "ubadge",
      "amount": "1000000000"
    }
  ]
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({
      pool_id: 1,
  token_in_maxs: [
    { denom: "ubadge", amount: "1000000000" }
  ]
});

const [shareOutAmount, tokensOut] = await gamm.calcJoinPoolShares(msgJson);
console.log("Shares out:", shareOutAmount.toString());
console.log("Tokens out:", tokensOut);
```

### 7. getPoolParams

Get pool parameters. Returns protobuf-encoded bytes.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1
}
```

### 8. getTotalShares

Get total shares for a pool.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({ pool_id: 1 });
const totalShares = await gamm.getTotalShares(msgJson);
console.log("Total shares:", totalShares);
// totalShares is { denom: string, amount: uint256 }
```

### 9. getTotalLiquidity

Get total liquidity for a pool.

**Protobuf JSON Format:**
```json
{
  "pool_id": 1
}
```

**TypeScript Example:**
```typescript
const msgJson = JSON.stringify({ pool_id: 1 });
const liquidity = await gamm.getTotalLiquidity(msgJson);
console.log("Total liquidity:", liquidity);
// liquidity is [{ denom: string, amount: uint256 }, ...]
```

## Complete Frontend Integration Example

```typescript
import { ethers } from "ethers";

// Precompile address
const GAMM_PRECOMPILE_ADDRESS = "0x0000000000000000000000000000000000001002";

// ABI (simplified - include all methods you need)
const gammABI = [
  "function joinPool(string memory msgJson) external returns (uint256 shareOutAmount, tuple(string denom, uint256 amount)[] tokenIn)",
  "function exitPool(string memory msgJson) external returns (tuple(string denom, uint256 amount)[] tokenOut)",
  "function swapExactAmountIn(string memory msgJson) external returns (uint256 tokenOutAmount)",
  "function swapExactAmountInWithIBCTransfer(string memory msgJson) external returns (uint256 tokenOutAmount)",
  "function getPool(string memory msgJson) external view returns (bytes memory pool)",
  "function getPoolType(string memory msgJson) external view returns (string memory poolType)",
  "function calcJoinPoolNoSwapShares(string memory msgJson) external view returns (tuple(string denom, uint256 amount)[] tokensOut, uint256 sharesOut)",
  "function calcExitPoolCoinsFromShares(string memory msgJson) external view returns (tuple(string denom, uint256 amount)[] tokensOut)",
  "function getTotalShares(string memory msgJson) external view returns (tuple(string denom, uint256 amount) totalShares)",
  "function getTotalLiquidity(string memory msgJson) external view returns (tuple(string denom, uint256 amount)[] liquidity)"
];

// Create contract instance
const gamm = new ethers.Contract(
  GAMM_PRECOMPILE_ADDRESS,
  gammABI,
  signer
);

// Example: Join a pool
async function joinPool(poolId: string, shareOutAmount: string, tokenInMaxs: Array<{denom: string, amount: string}>) {
  const msgJson = JSON.stringify({
    pool_id: poolId,
    share_out_amount: shareOutAmount,
    token_in_maxs: tokenInMaxs
  });

  // Estimate gas first
  const gasEstimate = await gamm.joinPool.estimateGas(msgJson);
  
  // Execute transaction
  const tx = await gamm.joinPool(msgJson, { gasLimit: gasEstimate });
  const receipt = await tx.wait();
  
  // Get return values (from events or static call)
  const [shares, tokens] = await gamm.joinPool.staticCall(msgJson);
  
  return { shares, tokens, receipt };
}

// Example: Swap tokens
async function swapTokens(
  routes: Array<{pool_id: number, token_out_denom: string}>,
  tokenIn: {denom: string, amount: string},
  tokenOutMinAmount: string
) {
  const msgJson = JSON.stringify({
    routes,
    token_in: tokenIn,
    token_out_min_amount: tokenOutMinAmount,
    affiliates: []
  });

  const tx = await gamm.swapExactAmountIn(msgJson);
  const receipt = await tx.wait();
  
  const tokenOutAmount = await gamm.swapExactAmountIn.staticCall(msgJson);
  
  return { tokenOutAmount, receipt };
}

// Example: Query pool information
async function getPoolInfo(poolId: string) {
  const msgJson = JSON.stringify({ pool_id: poolId });
  
  const [poolType, totalShares, totalLiquidity] = await Promise.all([
    gamm.getPoolType(msgJson),
    gamm.getTotalShares(msgJson),
    gamm.getTotalLiquidity(msgJson)
  ]);
  
  return { poolType, totalShares, totalLiquidity };
}
```

## Gas Costs

Base gas costs for transaction methods:
- **joinPool**: 50,000 gas
- **exitPool**: 50,000 gas
- **swapExactAmountIn**: 60,000 gas
- **swapExactAmountInWithIBCTransfer**: 80,000 gas

Base gas costs for query methods:
- **getPool**: 5,000 gas
- **getPools**: 5,000 gas
- **getPoolType**: 3,000 gas
- **calcJoinPoolNoSwapShares**: 10,000 gas
- **calcExitPoolCoinsFromShares**: 10,000 gas
- **calcJoinPoolShares**: 10,000 gas
- **getPoolParams**: 5,000 gas
- **getTotalShares**: 5,000 gas
- **getTotalLiquidity**: 5,000 gas

Additional gas may be consumed based on:
- Number of tokens in pools
- Number of routes in swaps
- Complexity of calculations

## Error Handling

The precompile uses structured error codes:

- **ErrorCodeInvalidInput (1)**: Invalid input parameters
- **ErrorCodePoolNotFound (2)**: Pool not found
- **ErrorCodeSwapFailed (3)**: Swap operation failed
- **ErrorCodeQueryFailed (4)**: Query operation failed
- **ErrorCodeInternalError (5)**: Internal error
- **ErrorCodeUnauthorized (6)**: Unauthorized operation
- **ErrorCodeJoinPoolFailed (7)**: Join pool operation failed
- **ErrorCodeExitPoolFailed (8)**: Exit pool operation failed
- **ErrorCodeIBCTransferFailed (9)**: IBC transfer operation failed

## Security Considerations

1. **Sender Authentication**: The `sender` field is **always** set from `msg.sender` and cannot be spoofed. Any `sender` in the JSON is ignored.

2. **Slippage Protection**: 
   - For `joinPool`: Use `token_in_maxs` to limit maximum tokens provided
   - For `exitPool`: Use `token_out_mins` to ensure minimum tokens received
   - For `swapExactAmountIn`: Use `token_out_min_amount` to ensure minimum output

3. **Route Validation**: Validate swap routes before execution to ensure they're correct.

4. **IBC Timeout**: When using `swapExactAmountInWithIBCTransfer`, ensure `timeout_timestamp` is set appropriately (not too short, not too long).

5. **Affiliate Fees**: Affiliate fees are calculated from `token_out_min_amount`, not the actual output. Be aware of this when setting affiliate fees.

## Common Use Cases

### 1. Add Liquidity to a Pool

```typescript
// Calculate expected shares first
const calcMsgJson = JSON.stringify({
      pool_id: 1,
  token_in_maxs: [
    { denom: "ubadge", amount: "1000000000" },
    { denom: "ustake", amount: "2000000000" }
  ]
});

const [tokensOut, sharesOut] = await gamm.calcJoinPoolNoSwapShares(calcMsgJson);
console.log("Expected shares:", sharesOut.toString());

// Join the pool
const joinMsgJson = JSON.stringify({
      pool_id: 1,
  share_out_amount: sharesOut.toString(),
  token_in_maxs: [
    { denom: "ubadge", amount: "1000000000" },
    { denom: "ustake", amount: "2000000000" }
  ]
});

const tx = await gamm.joinPool(joinMsgJson);
await tx.wait();
```

### 2. Remove Liquidity from a Pool

```typescript
// Calculate expected tokens first
const calcMsgJson = JSON.stringify({
      pool_id: 1,
  share_in_amount: "500000"
});

const tokensOut = await gamm.calcExitPoolCoinsFromShares(calcMsgJson);
console.log("Expected tokens:", tokensOut);

// Exit the pool
const exitMsgJson = JSON.stringify({
      pool_id: 1,
  share_in_amount: "500000",
  token_out_mins: tokensOut.map(t => ({
    denom: t.denom,
    amount: (BigInt(t.amount) * 95n / 100n).toString() // 5% slippage tolerance
  }))
});

const tx = await gamm.exitPool(exitMsgJson);
await tx.wait();
```

### 3. Swap Tokens (Single Hop)

```typescript
const msgJson = JSON.stringify({
  routes: [
    { pool_id: 1, token_out_denom: "ustake" }
  ],
  token_in: {
    denom: "ubadge",
    amount: "1000000000" // 1 BADGE
  },
  token_out_min_amount: "1900000000", // 5% slippage tolerance
  affiliates: []
});

const tx = await gamm.swapExactAmountIn(msgJson);
await tx.wait();
```

### 4. Swap Tokens (Multi-Hop)

```typescript
const msgJson = JSON.stringify({
  routes: [
    { pool_id: 1, token_out_denom: "uatom" },
    { pool_id: 2, token_out_denom: "ustake" }
  ],
  token_in: {
    denom: "ubadge",
    amount: "1000000000"
  },
  token_out_min_amount: "1800000000", // Account for multiple hops
  affiliates: []
});

const tx = await gamm.swapExactAmountIn(msgJson);
await tx.wait();
```

### 5. Swap and Transfer via IBC

```typescript
const timeoutTimestamp = Math.floor(Date.now() / 1000) + 300; // 5 minutes from now
const msgJson = JSON.stringify({
  routes: [
    { pool_id: 1, token_out_denom: "ustake" }
  ],
  token_in: {
    denom: "ubadge",
    amount: "1000000000"
  },
  token_out_min_amount: "1900000000",
  ibc_transfer_info: {
    source_channel: "channel-0",
    receiver: "cosmos1...",
    memo: "",
    timeout_timestamp: (timeoutTimestamp * 1e9).toString() // Convert to nanoseconds
  },
  affiliates: []
});

const tx = await gamm.swapExactAmountInWithIBCTransfer(msgJson);
await tx.wait();
```

## Important Notes

1. **Amount Format**: All amounts are strings in the smallest unit (e.g., `"1000000000"` for 1 BADGE if BADGE has 9 decimals).

2. **Pool IDs**: Pool IDs are `uint64` values, but in JSON they're strings (e.g., `"1"` not `1`).

3. **Sender Field**: The `sender` field in all transaction messages is **automatically set** from `msg.sender` and cannot be overridden.

4. **Return Values**: Transaction methods return values that can be accessed via:
   - Static calls: `await contract.methodName.staticCall(msgJson)`
   - Events: Decode from transaction receipt events
   - Return data: Decode from transaction receipt return data

5. **Query Methods**: Query methods are `view` functions and don't consume gas when called from frontend (read-only calls).

6. **Protobuf Encoding**: Some query methods return protobuf-encoded bytes. You'll need to decode these using the appropriate protobuf library.

7. **Slippage Protection**: Always use minimum/maximum amounts to protect against slippage:
   - `token_in_maxs` for joins
   - `token_out_mins` for exits
   - `token_out_min_amount` for swaps

8. **Multi-Hop Swaps**: When using multiple routes, account for cumulative slippage across all hops.

## Helper Libraries

The codebase includes helper libraries for constructing JSON strings:

- **GammJSONHelpers**: Solidity library for building JSON strings
- **GammTypes**: Type definitions for all gamm types
- **GammWrappers**: Wrapper functions for common operations

See `contracts/libraries/` for implementation details.

## Testing

```typescript
import { ethers } from "hardhat";

describe("Gamm Precompile", () => {
  it("should join a pool", async () => {
    const gamm = await ethers.getContractAt(
      "IGammPrecompile",
      "0x0000000000000000000000000000000000001002"
    );

    const msgJson = JSON.stringify({
      pool_id: 1,
      share_out_amount: "1000000",
      token_in_maxs: [
        { denom: "ubadge", amount: "1000000000" },
        { denom: "ustake", amount: "2000000000" }
      ]
    });

    const tx = await gamm.joinPool(msgJson);
    await tx.wait();

    // Verify success
    expect(tx).to.not.be.null;
  });

  it("should swap tokens", async () => {
    const gamm = await ethers.getContractAt(
      "IGammPrecompile",
      "0x0000000000000000000000000000000000001002"
    );

    const msgJson = JSON.stringify({
      routes: [
        { pool_id: 1, token_out_denom: "ustake" }
      ],
      token_in: {
        denom: "ubadge",
        amount: "1000000000"
      },
      token_out_min_amount: "1900000000",
      affiliates: []
    });

    const tx = await gamm.swapExactAmountIn(msgJson);
    await tx.wait();

    expect(tx).to.not.be.null;
  });
});
```

## Summary

The Gamm Precompile provides a comprehensive interface for interacting with liquidity pools from EVM. It supports:

- ✅ Joining and exiting pools
- ✅ Swapping tokens (single and multi-hop)
- ✅ IBC transfers with swaps
- ✅ Querying pool information
- ✅ Calculating expected outcomes

All methods use JSON-based protobuf format, making it easy to integrate with frontend applications.

**Key Benefits:**
- Standardized EVM interface for liquidity pools
- Support for complex operations (multi-hop swaps, IBC transfers)
- Comprehensive query methods for pool information
- Slippage protection built-in
- Helper libraries for JSON construction

