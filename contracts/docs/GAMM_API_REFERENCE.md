# Gamm API Reference

Complete reference for the BitBadges gamm precompile and helper libraries.

## Table of Contents

- [Precompile Interface](#precompile-interface)
- [GammWrappers](#gammwrappers)
- [GammHelpers](#gammhelpers)
- [GammBuilders](#gammbuilders)
- [GammJSONHelpers](#gammjsonhelpers)
- [GammErrors](#gammerrors)

## Precompile Interface

### Address

```solidity
IGammPrecompile constant GAMM = 
    IGammPrecompile(0x0000000000000000000000000000000000001002);
```

### Transaction Methods

All transaction methods use JSON string parameters and return result values.

#### `joinPool(string calldata msgJson) → (uint256, Coin[])`

Join a liquidity pool by providing tokens.

**Parameters:**
- `msgJson`: JSON string matching `MsgJoinPool` protobuf format

**Returns:**
- `shareOutAmount`: The amount of pool shares received
- `tokenIn`: The actual tokens provided (may be less than tokenInMaxs)

**Example:**
```solidity
string memory json = GammJSONHelpers.joinPoolJSON(
    poolId, shareOutAmount, tokenInMaxsJson
);
(uint256 shares, GammTypes.Coin[] memory tokens) = GAMM.joinPool(json);
```

#### `exitPool(string calldata msgJson) → Coin[]`

Exit a liquidity pool by burning shares.

**Parameters:**
- `msgJson`: JSON string matching `MsgExitPool` protobuf format

**Returns:**
- `tokenOut`: The tokens received from exiting the pool

**Example:**
```solidity
string memory json = GammJSONHelpers.exitPoolJSON(
    poolId, shareInAmount, tokenOutMinsJson
);
GammTypes.Coin[] memory tokens = GAMM.exitPool(json);
```

#### `swapExactAmountIn(string calldata msgJson) → uint256`

Swap tokens with exact input amount.

**Parameters:**
- `msgJson`: JSON string matching `MsgSwapExactAmountIn` protobuf format

**Returns:**
- `tokenOutAmount`: The amount of output tokens received

**Example:**
```solidity
string memory json = GammJSONHelpers.swapExactAmountInJSON(
    routesJson, tokenInJson, tokenOutMinAmount, affiliatesJson
);
uint256 tokenOut = GAMM.swapExactAmountIn(json);
```

#### `swapExactAmountInWithIBCTransfer(string calldata msgJson) → uint256`

Swap tokens with exact input amount and transfer via IBC.

**Parameters:**
- `msgJson`: JSON string matching `MsgSwapExactAmountInWithIBCTransfer` protobuf format

**Returns:**
- `tokenOutAmount`: The amount of output tokens received

### Query Methods

Query methods return typed values (Coin, Coin[], uint256, string) or protobuf-encoded bytes.

#### `getPool(string calldata msgJson) → bytes`

Get pool data by ID.

**Parameters:**
- `msgJson`: JSON string with `poolId`

**Returns:**
- `pool`: Protobuf-encoded pool bytes

**Note:** Full decoding requires off-chain tools. Use `getTotalShares` or `getTotalLiquidity` for simple queries.

#### `getTotalShares(string calldata msgJson) → Coin`

Get total shares for a pool.

**Parameters:**
- `msgJson`: JSON string with `poolId`

**Returns:**
- `totalShares`: The total shares as a Coin struct

**Example:**
```solidity
GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);
```

#### `getTotalLiquidity(string calldata msgJson) → Coin[]`

Get total liquidity for a pool.

**Parameters:**
- `msgJson`: JSON string with `poolId`

**Returns:**
- `liquidity`: The total liquidity as an array of Coin structs

**Example:**
```solidity
GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, poolId);
```

#### `calcJoinPoolNoSwapShares(string calldata msgJson) → (Coin[], uint256)`

Calculate shares for joining pool without swap.

**Parameters:**
- `msgJson`: JSON string with `poolId` and `tokenInMaxs`

**Returns:**
- `tokensOut`: The tokens that would be provided
- `sharesOut`: The shares that would be received

#### `calcExitPoolCoinsFromShares(string calldata msgJson) → Coin[]`

Calculate tokens received for exiting pool.

**Parameters:**
- `msgJson`: JSON string with `poolId` and `shareInAmount`

**Returns:**
- `tokensOut`: The tokens that would be received

#### `calcJoinPoolShares(string calldata msgJson) → (uint256, Coin[])`

Calculate shares for joining pool (with swap).

**Parameters:**
- `msgJson`: JSON string with `poolId` and `tokenInMaxs`

**Returns:**
- `shareOutAmount`: The shares that would be received
- `tokensOut`: The tokens that would be provided

## GammWrappers

Type-safe wrapper functions that accept structs instead of JSON strings.

### Transaction Wrappers

#### `joinPool(...) → (uint256, Coin[])`

Join a pool with typed parameters.

```solidity
function joinPool(
    IGammPrecompile precompile,
    uint64 poolId,
    uint256 shareOutAmount,
    GammTypes.Coin[] memory tokenInMaxs
) internal returns (uint256, GammTypes.Coin[] memory)
```

#### `exitPool(...) → Coin[]`

Exit a pool with typed parameters.

```solidity
function exitPool(
    IGammPrecompile precompile,
    uint64 poolId,
    uint256 shareInAmount,
    GammTypes.Coin[] memory tokenOutMins
) internal returns (GammTypes.Coin[] memory)
```

#### `swapExactAmountIn(...) → uint256`

Swap tokens with typed parameters.

```solidity
function swapExactAmountIn(
    IGammPrecompile precompile,
    GammTypes.SwapAmountInRoute[] memory routes,
    GammTypes.Coin memory tokenIn,
    uint256 tokenOutMinAmount,
    GammTypes.Affiliate[] memory affiliates
) internal returns (uint256)
```

### Query Wrappers

#### `getPool(...) → bytes`

Get pool data by ID.

```solidity
function getPool(
    IGammPrecompile precompile,
    uint64 poolId
) internal view returns (bytes memory)
```

#### `getTotalShares(...) → Coin`

Get total shares for a pool.

```solidity
function getTotalShares(
    IGammPrecompile precompile,
    uint64 poolId
) internal view returns (GammTypes.Coin memory)
```

#### `getTotalLiquidity(...) → Coin[]`

Get total liquidity for a pool.

```solidity
function getTotalLiquidity(
    IGammPrecompile precompile,
    uint64 poolId
) internal view returns (GammTypes.Coin[] memory)
```

## GammHelpers

Utility functions for creating gamm types.

### Type Creation

#### `createCoin(string memory denom, uint256 amount) → Coin`

Create a Coin struct.

#### `createSwapRoute(uint64 poolId, string memory tokenOutDenom) → SwapAmountInRoute`

Create a SwapAmountInRoute struct.

#### `createAffiliate(address address_, uint256 basisPointsFee) → Affiliate`

Create an Affiliate struct.

#### `createIBCTransferInfo(...) → IBCTransferInfo`

Create an IBCTransferInfo struct.

#### `createSingleHopRoute(uint64 poolId, string memory tokenOutDenom) → SwapAmountInRoute[]`

Create a single-hop swap route.

#### `createTwoHopRoute(...) → SwapAmountInRoute[]`

Create a two-hop swap route.

## GammBuilders

Fluent builder APIs for complex operations.

### JoinPoolBuilder

```solidity
GammBuilders.JoinPoolBuilder memory builder = 
    GammBuilders.newJoinPool(poolId, shareOutAmount);

builder = builder.addTokenInMax("uatom", 1000000);
builder = builder.addTokenInMax("uosmo", 2000000);

string memory json = builder.build();
(uint256 shares, GammTypes.Coin[] memory tokens) = GAMM.joinPool(json);
```

**Methods:**
- `newJoinPool(poolId, shareOutAmount)` - Create a new builder
- `addTokenInMax(denom, amount)` - Add a token input maximum
- `build()` - Build the JSON string

### ExitPoolBuilder

```solidity
GammBuilders.ExitPoolBuilder memory builder = 
    GammBuilders.newExitPool(poolId, shareInAmount);

builder = builder.addTokenOutMin("uatom", 500000);
builder = builder.addTokenOutMin("uosmo", 1000000);

string memory json = builder.build();
GammTypes.Coin[] memory tokens = GAMM.exitPool(json);
```

### SwapBuilder

```solidity
GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();

builder = builder.addRoute(poolId1, "uosmo");
builder = builder.addRoute(poolId2, "uion");
builder = builder.withTokenIn("uatom", 1000000);
builder = builder.withTokenOutMinAmount(900000);
builder = builder.addAffiliate(affiliateAddr, 100); // 1% fee

string memory json = builder.build();
uint256 tokenOut = GAMM.swapExactAmountIn(json);
```

**Methods:**
- `newSwap()` - Create a new builder
- `addRoute(poolId, tokenOutDenom)` - Add a swap route
- `withTokenIn(denom, amount)` - Set input token
- `withTokenOutMinAmount(amount)` - Set minimum output amount
- `addAffiliate(address, basisPointsFee)` - Add an affiliate
- `build()` - Build the JSON string

## GammJSONHelpers

Low-level JSON construction helpers.

### Transaction JSON Constructors

#### `joinPoolJSON(poolId, shareOutAmount, tokenInMaxsJson) → string`

Construct JSON for `joinPool`.

#### `exitPoolJSON(poolId, shareInAmount, tokenOutMinsJson) → string`

Construct JSON for `exitPool`.

#### `swapExactAmountInJSON(routesJson, tokenInJson, tokenOutMinAmount, affiliatesJson) → string`

Construct JSON for `swapExactAmountIn`.

### Query JSON Constructors

#### `getPoolJSON(poolId) → string`

Construct JSON for `getPool`.

#### `getTotalSharesJSON(poolId) → string`

Construct JSON for `getTotalShares`.

### Type to JSON Converters

#### `coinToJson(Coin memory coin) → string`

Convert a Coin to JSON.

#### `coinsToJson(Coin[] memory coins) → string`

Convert an array of Coins to JSON.

#### `swapRouteToJson(SwapAmountInRoute memory route) → string`

Convert a SwapAmountInRoute to JSON.

#### `swapRoutesToJson(SwapAmountInRoute[] memory routes) → string`

Convert an array of SwapAmountInRoute to JSON.

#### `affiliateToJson(Affiliate memory affiliate) → string`

Convert an Affiliate to JSON.

#### `affiliatesToJson(Affiliate[] memory affiliates) → string`

Convert an array of Affiliates to JSON.

## GammErrors

Custom error types for better error handling.

### Common Errors

```solidity
error InvalidPoolId(uint64 poolId);
error InvalidCoin(string message);
error InvalidRoute(string message);
error InvalidAffiliate(string message);
error PoolNotFound(uint64 poolId);
error TransactionFailed(string message);
error QueryFailed(string message);
```

### Validation Helpers

```solidity
GammErrors.requireValidPoolId(poolId);
GammErrors.requireValidCoin(denom, amount);
GammErrors.requireValidRoute(poolId, tokenOutDenom);
GammErrors.requireValidAffiliate(address_, basisPointsFee);
```

## Type Definitions

See `contracts/types/GammTypes.sol` for complete type definitions including:
- `Coin` - Token with denomination and amount
- `SwapAmountInRoute` - Single hop in a swap route
- `Affiliate` - Fee recipient for swaps
- `IBCTransferInfo` - IBC transfer information
- Message types: `MsgJoinPool`, `MsgExitPool`, `MsgSwapExactAmountIn`, etc.
- Query types: `QueryPoolRequest`, `QueryTotalSharesRequest`, etc.

