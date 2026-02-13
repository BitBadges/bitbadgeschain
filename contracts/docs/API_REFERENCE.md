# API Reference

Complete reference for the BitBadges tokenization precompile and helper libraries.

## Table of Contents

- [Precompile Interface](#precompile-interface)
- [TokenizationWrappers](#tokenizationwrappers)
- [TokenizationHelpers](#tokenizationhelpers)
- [TokenizationBuilders](#tokenizationbuilders)
- [TokenizationJSONHelpers](#tokenizationjsonhelpers)
- [TokenizationErrors](#tokenizationerrors)

## Precompile Interface

### Address

```solidity
ITokenizationPrecompile constant TOKENIZATION = 
    ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
```

### Transaction Methods

All transaction methods use JSON string parameters and return success indicators or result values.

#### `transferTokens(string calldata msgJson) → bool`

Transfer tokens from the caller to specified addresses.

**Parameters:**
- `msgJson`: JSON string matching `MsgTransferTokens` protobuf format

**Returns:**
- `success`: True if transfer succeeded

**Example:**
```solidity
string memory json = TokenizationJSONHelpers.transferTokensJSON(
    collectionId, recipients, amount, tokenIdsJson, ownershipTimesJson
);
bool success = TOKENIZATION.transferTokens(json);
```

#### `createCollection(string calldata msgJson) → uint256`

Create a new token collection.

**Parameters:**
- `msgJson`: JSON string matching `MsgCreateCollection` protobuf format

**Returns:**
- `newCollectionId`: The newly created collection ID

**Example:**
```solidity
string memory json = builder.build(); // Using TokenizationBuilders
uint256 collectionId = TOKENIZATION.createCollection(json);
```

#### `createDynamicStore(string calldata msgJson) → uint256`

Create a new dynamic store (boolean registry).

**Parameters:**
- `msgJson`: JSON string matching `MsgCreateDynamicStore` protobuf format

**Returns:**
- `storeId`: The newly created store ID

**Example:**
```solidity
uint256 storeId = TokenizationWrappers.createDynamicStore(
    TOKENIZATION, false, "ipfs://...", "KYC Registry"
);
```

### Query Methods

Query methods return protobuf-encoded bytes (except `getBalanceAmount`, `getTotalSupply`, and `getAllReservedProtocolAddresses`).

#### `getCollection(string calldata msgJson) → bytes`

Get collection details by ID.

**Parameters:**
- `msgJson`: JSON string with `collectionId`

**Returns:**
- `collection`: Protobuf-encoded `TokenCollection` bytes

**Note:** Full decoding requires off-chain tools. Use `getBalanceAmount` for simple queries.

#### `getBalance(string calldata msgJson) → bytes`

Get user balance for a collection.

**Parameters:**
- `msgJson`: JSON string with `collectionId` and `userAddress`

**Returns:**
- `balance`: Protobuf-encoded `UserBalanceStore` bytes

#### `getBalanceAmount(string calldata msgJson) → uint256`

Get balance amount for specific token/ownership ranges.

**Parameters:**
- `msgJson`: JSON string with `collectionId`, `userAddress`, `tokenIds`, and `ownershipTimes`

**Returns:**
- `amount`: The total balance amount

**Example:**
```solidity
uint256 balance = TokenizationWrappers.getBalanceAmount(
    TOKENIZATION, collectionId, user, tokenIds, ownershipTimes
);
```

#### `getTotalSupply(string calldata msgJson) → uint256`

Get total supply for specific token/ownership ranges.

**Parameters:**
- `msgJson`: JSON string with `collectionId`, `tokenIds`, and `ownershipTimes`

**Returns:**
- `amount`: The total supply amount

## TokenizationWrappers

Type-safe wrapper functions that accept structs instead of JSON strings.

### Transaction Wrappers

#### `transferTokens(...) → bool`

Transfer tokens with typed parameters.

```solidity
function transferTokens(
    ITokenizationPrecompile precompile,
    uint256 collectionId,
    address[] memory toAddresses,
    uint256 amount,
    TokenizationTypes.UintRange[] memory tokenIds,
    TokenizationTypes.UintRange[] memory ownershipTimes
) internal returns (bool success)
```

#### `transferSingleToken(...) → bool`

Convenience wrapper for transferring a single token.

```solidity
function transferSingleToken(
    ITokenizationPrecompile precompile,
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) internal returns (bool success)
```

#### `transferTokensWithFullOwnership(...) → bool`

Transfer tokens with full ownership time range.

```solidity
function transferTokensWithFullOwnership(
    ITokenizationPrecompile precompile,
    uint256 collectionId,
    address[] memory toAddresses,
    uint256 amount,
    TokenizationTypes.UintRange[] memory tokenIds
) internal returns (bool success)
```

### Query Wrappers

#### `getBalanceAmount(...) → uint256`

Get balance amount with typed parameters.

```solidity
function getBalanceAmount(
    ITokenizationPrecompile precompile,
    uint256 collectionId,
    address userAddress,
    TokenizationTypes.UintRange[] memory tokenIds,
    TokenizationTypes.UintRange[] memory ownershipTimes
) internal view returns (uint256 amount)
```

## TokenizationHelpers

Utility functions for creating and validating tokenization types.

### Range Creation

#### `createUintRange(uint256 start, uint256 end) → UintRange`

Create a UintRange struct.

#### `createFullOwnershipTimeRange() → UintRange`

Create a full ownership time range (1 to max uint256).

#### `createSingleTokenIdRange(uint256 tokenId) → UintRange`

Create a range for a single token ID.

#### `createOwnershipTimeRange(uint256 startTime, uint256 duration) → UintRange`

Create an ownership time range from start time with duration.

### Metadata Creation

#### `createCollectionMetadata(string memory uri, string memory customData) → CollectionMetadata`

Create a CollectionMetadata struct.

#### `createTokenMetadata(string memory uri, string memory customData, UintRange[] memory tokenIds) → TokenMetadata`

Create a TokenMetadata struct.

### Validation

#### `validateUintRange(UintRange memory range) → bool`

Validate that a UintRange is valid (start <= end).

#### `validateUintRangeArray(UintRange[] memory ranges) → bool`

Validate an array of UintRanges.

## TokenizationBuilders

Fluent builder APIs for complex operations.

### CollectionBuilder

```solidity
TokenizationBuilders.CollectionBuilder memory builder = 
    TokenizationBuilders.newCollection();

builder = builder.withValidTokenIdRange(1, 1000);
builder = builder.withManager(managerAddress);
builder = builder.withMetadata(metadata);
builder = builder.withStandards(standards);

string memory json = builder.build();
uint256 collectionId = TOKENIZATION.createCollection(json);
```

**Methods:**
- `newCollection()` - Create a new builder
- `withValidTokenIds(...)` - Set valid token ID ranges
- `withManager(...)` - Set manager address
- `withMetadata(...)` - Set collection metadata
- `withDefaultBalances(...)` - Set default balances
- `withStandards(...)` - Set standards array
- `build()` - Build the JSON string

### TransferBuilder

```solidity
TokenizationBuilders.TransferBuilder memory builder = 
    TokenizationBuilders.newTransfer();

builder = builder.withRecipient(recipient);
builder = builder.withAmount(amount);
builder = builder.withTokenId(tokenId);
builder = builder.withFullOwnershipTime();

string memory json = builder.buildTransfer(collectionId);
bool success = TOKENIZATION.transferTokens(json);
```

## TokenizationJSONHelpers

Low-level JSON construction helpers.

### Common Methods

#### `transferTokensJSON(...) → string`

Construct JSON for `transferTokens`.

#### `getCollectionJSON(uint256 collectionId) → string`

Construct JSON for `getCollection`.

#### `getBalanceJSON(uint256 collectionId, address userAddress) → string`

Construct JSON for `getBalance`.

#### `uintRangeArrayToJson(uint256[] memory starts, uint256[] memory ends) → string`

Convert UintRange arrays to JSON.

## TokenizationErrors

Custom error types for better error handling.

### Common Errors

```solidity
error CollectionNotFound(uint256 collectionId);
error InsufficientBalance(uint256 collectionId, uint256 required, uint256 available);
error TransferFailed(uint256 collectionId, string reason);
error InvalidTokenId(uint256 collectionId, uint256 tokenId);
error ApprovalDenied(uint256 collectionId, string approvalId);
```

### Validation Helpers

```solidity
TokenizationErrors.requireValidCollectionId(collectionId);
TokenizationErrors.requireValidAddress(address_);
TokenizationErrors.requireNonEmptyString(str, "parameterName");
```

## Type Definitions

See `contracts/types/TokenizationTypes.sol` for complete type definitions including:
- `UintRange` - Range of IDs
- `Balance` - Token balance
- `CollectionMetadata` - Collection metadata
- `TokenMetadata` - Token metadata
- `UserBalanceStore` - User balance and approvals
- And many more...

## Return Value Decoding

Most query methods return protobuf-encoded bytes. For production use:

1. **Use direct queries** when possible:
   - `getBalanceAmount()` returns `uint256` directly
   - `getTotalSupply()` returns `uint256` directly
   - `getAllReservedProtocolAddresses()` returns `address[]` directly

2. **Decode off-chain** for complex types:
   - Use the TypeScript SDK for full decoding
   - Emit events with raw bytes for off-chain indexing

3. **See TokenizationDecoders** for limitations and recommendations

