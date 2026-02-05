# Tokenization Precompile API Reference

## Overview

The Tokenization Precompile provides a Solidity interface to interact with the BitBadges tokenization module from EVM-compatible smart contracts. The precompile is available at address `0x0000000000000000000000000000000000001001`.

## Precompile Address

```
0x0000000000000000000000000000000000001001
```

## Data Types

### UintRange

A range structure used for token IDs, ownership times, and transfer times:

```solidity
struct UintRange {
    uint256 start;
    uint256 end;
}
```

**Constraints:**
- `start` must be <= `end`
- Both `start` and `end` must be non-negative
- Maximum array size: 100 ranges

### UserIncomingApproval

Structure for incoming approval settings:

```solidity
struct UserIncomingApproval {
    string approvalId;
    string fromListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
}
```

### UserOutgoingApproval

Structure for outgoing approval settings:

```solidity
struct UserOutgoingApproval {
    string approvalId;
    string toListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
}
```

## Transaction Methods

### transferTokens

Transfers tokens from the caller (`msg.sender`) to one or more recipient addresses.

**Signature:**
```solidity
function transferTokens(
    uint256 collectionId,
    address[] calldata toAddresses,
    uint256 amount,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external returns (bool);
```

**Parameters:**
- `collectionId` (uint256): The collection ID to transfer from
- `toAddresses` (address[]): Array of recipient EVM addresses (max 100 addresses)
- `amount` (uint256): Amount to transfer to each recipient (must be > 0)
- `tokenIds` (UintRange[]): Array of token ID ranges to transfer (max 100 ranges)
- `ownershipTimes` (UintRange[]): Array of ownership time ranges to transfer (max 100 ranges)

**Returns:**
- `bool`: `true` if transfer succeeded

**Gas Cost:**
- Base: 30,000 gas
- Per recipient: 5,000 gas
- Per token ID range: 1,000 gas
- Per ownership time range: 1,000 gas

**Example:**
```solidity
IBadgesPrecompile precompile = IBadgesPrecompile(0x0000000000000000000000000000000000001001);

UintRange[] memory tokenIds = new UintRange[](1);
tokenIds[0] = UintRange({start: 1, end: 10});

UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange({start: 0, end: type(uint256).max});

address[] memory recipients = new address[](1);
recipients[0] = 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0;

bool success = precompile.transferTokens(
    collectionId,
    recipients,
    10,
    tokenIds,
    ownershipTimes
);
require(success, "Transfer failed");
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters (zero addresses, invalid ranges, etc.)
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `4` (ErrorCodeTransferFailed): Transfer operation failed (insufficient balance, approval issues, etc.)

---

### setIncomingApproval

Sets an incoming approval for the caller, allowing specified addresses to transfer tokens to the caller.

**Signature:**
```solidity
function setIncomingApproval(
    uint256 collectionId,
    UserIncomingApproval calldata approval
) external returns (bool);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approval` (UserIncomingApproval): Approval configuration struct

**Returns:**
- `bool`: `true` if approval was set successfully

**Gas Cost:**
- Base: 20,000 gas
- Dynamic: Based on number of ranges in approval struct

**Example:**
```solidity
UserIncomingApproval memory approval = UserIncomingApproval({
    approvalId: "my_incoming_approval",
    fromListId: "All",
    initiatedByListId: "All",
    transferTimes: new UintRange[](0), // Empty means all times
    tokenIds: new UintRange[](0), // Empty means all token IDs
    ownershipTimes: new UintRange[](0), // Empty means all ownership times
    uri: "",
    customData: ""
});

bool success = precompile.setIncomingApproval(collectionId, approval);
require(success, "Set incoming approval failed");
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `5` (ErrorCodeApprovalFailed): Approval operation failed

---

### setOutgoingApproval

Sets an outgoing approval for the caller, allowing the caller to transfer tokens to specified addresses.

**Signature:**
```solidity
function setOutgoingApproval(
    uint256 collectionId,
    UserOutgoingApproval calldata approval
) external returns (bool);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approval` (UserOutgoingApproval): Approval configuration struct

**Returns:**
- `bool`: `true` if approval was set successfully

**Gas Cost:**
- Base: 20,000 gas
- Dynamic: Based on number of ranges in approval struct

**Example:**
```solidity
UserOutgoingApproval memory approval = UserOutgoingApproval({
    approvalId: "my_outgoing_approval",
    toListId: "All",
    initiatedByListId: "All",
    transferTimes: new UintRange[](0),
    tokenIds: new UintRange[](0),
    ownershipTimes: new UintRange[](0),
    uri: "",
    customData: ""
});

bool success = precompile.setOutgoingApproval(collectionId, approval);
require(success, "Set outgoing approval failed");
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `5` (ErrorCodeApprovalFailed): Approval operation failed

---

## Query Methods

### getCollection

Queries collection data by ID. Returns protobuf-encoded collection data.

**Signature:**
```solidity
function getCollection(uint256 collectionId) external view returns (bytes);
```

**Parameters:**
- `collectionId` (uint256): The collection ID to query

**Returns:**
- `bytes`: Protobuf-encoded collection data

**Gas Cost:** 3,000 gas

**Example:**
```solidity
bytes memory collectionData = precompile.getCollection(collectionId);
// Decode using protobuf library in Solidity
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid collection ID
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getBalance

Queries balance data for a user address. Returns protobuf-encoded balance data.

**Signature:**
```solidity
function getBalance(
    uint256 collectionId,
    address userAddress
) external view returns (bytes);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `userAddress` (address): The user address to query (cannot be zero address)

**Returns:**
- `bytes`: Protobuf-encoded balance data

**Gas Cost:** 3,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getBalanceAmount

Gets the balance amount for a user with specific token IDs and ownership times. Returns a uint256 value directly (no protobuf decoding needed).

**Signature:**
```solidity
function getBalanceAmount(
    uint256 collectionId,
    address userAddress,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external view returns (uint256);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `userAddress` (address): The user address to query
- `tokenIds` (UintRange[]): Array of token ID ranges to query
- `ownershipTimes` (UintRange[]): Array of ownership time ranges to query

**Returns:**
- `uint256`: Total balance amount matching the specified ranges

**Gas Cost:**
- Base: 3,000 gas
- Per range: 500 gas

**Example:**
```solidity
UintRange[] memory tokenIds = new UintRange[](1);
tokenIds[0] = UintRange({start: 1, end: 10});

UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange({start: 0, end: type(uint256).max});

uint256 balance = precompile.getBalanceAmount(
    collectionId,
    userAddress,
    tokenIds,
    ownershipTimes
);
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getTotalSupply

Gets the total supply for a collection with specific token IDs and ownership times. Returns a uint256 value directly.

**Signature:**
```solidity
function getTotalSupply(
    uint256 collectionId,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external view returns (uint256);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `tokenIds` (UintRange[]): Array of token ID ranges to query
- `ownershipTimes` (UintRange[]): Array of ownership time ranges to query

**Returns:**
- `uint256`: Total supply matching the specified ranges

**Gas Cost:**
- Base: 3,000 gas
- Per range: 500 gas

**Example:**
```solidity
UintRange[] memory tokenIds = new UintRange[](1);
tokenIds[0] = UintRange({start: 1, end: 100});

UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange({start: 0, end: type(uint256).max});

uint256 totalSupply = precompile.getTotalSupply(
    collectionId,
    tokenIds,
    ownershipTimes
);
```

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `2` (ErrorCodeCollectionNotFound): Collection does not exist
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getAddressList

Queries an address list by ID. Returns protobuf-encoded address list data.

**Signature:**
```solidity
function getAddressList(string calldata listId) external view returns (bytes);
```

**Parameters:**
- `listId` (string): The address list ID

**Returns:**
- `bytes`: Protobuf-encoded address list data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid list ID (empty string)
- `6` (ErrorCodeQueryFailed): Address list not found

---

### getApprovalTracker

Queries an approval tracker. Returns protobuf-encoded tracker data.

**Signature:**
```solidity
function getApprovalTracker(
    uint256 collectionId,
    string calldata approvalLevel,
    address approverAddress,
    string calldata amountTrackerId,
    string calldata trackerType,
    address approvedAddress,
    string calldata approvalId
) external view returns (bytes);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approvalLevel` (string): Approval level (e.g., "collection", "user")
- `approverAddress` (address): Address of the approver
- `amountTrackerId` (string): Amount tracker ID
- `trackerType` (string): Tracker type (e.g., "overall", "perFromAddress")
- `approvedAddress` (address): Address that was approved
- `approvalId` (string): Approval ID

**Returns:**
- `bytes`: Protobuf-encoded approval tracker data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getChallengeTracker

Queries a challenge tracker. Returns the number of times the challenge has been used as uint256.

**Signature:**
```solidity
function getChallengeTracker(
    uint256 collectionId,
    string calldata approvalLevel,
    address approverAddress,
    string calldata challengeTrackerId,
    uint256 leafIndex,
    string calldata approvalId
) external view returns (uint256);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approvalLevel` (string): Approval level
- `approverAddress` (address): Address of the approver
- `challengeTrackerId` (string): Challenge tracker ID
- `leafIndex` (uint256): Leaf index in the merkle tree
- `approvalId` (string): Approval ID

**Returns:**
- `uint256`: Number of times the challenge has been used

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters (negative leafIndex, etc.)
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getETHSignatureTracker

Queries an ETH signature tracker. Returns the number of times the signature has been used as uint256.

**Signature:**
```solidity
function getETHSignatureTracker(
    uint256 collectionId,
    string calldata approvalLevel,
    address approverAddress,
    string calldata approvalId,
    string calldata challengeTrackerId,
    string calldata signature
) external view returns (uint256);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approvalLevel` (string): Approval level
- `approverAddress` (address): Address of the approver
- `approvalId` (string): Approval ID
- `challengeTrackerId` (string): Challenge tracker ID
- `signature` (string): The signature to check

**Returns:**
- `uint256`: Number of times the signature has been used

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getDynamicStore

Queries a dynamic store by ID. Returns protobuf-encoded store data.

**Signature:**
```solidity
function getDynamicStore(uint256 storeId) external view returns (bytes);
```

**Parameters:**
- `storeId` (uint256): The dynamic store ID

**Returns:**
- `bytes`: Protobuf-encoded dynamic store data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid store ID
- `6` (ErrorCodeQueryFailed): Dynamic store not found

---

### getDynamicStoreValue

Queries a dynamic store value for a specific user address. Returns protobuf-encoded value data.

**Signature:**
```solidity
function getDynamicStoreValue(
    uint256 storeId,
    address userAddress
) external view returns (bytes);
```

**Parameters:**
- `storeId` (uint256): The dynamic store ID
- `userAddress` (address): The user address to query

**Returns:**
- `bytes`: Protobuf-encoded store value data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getWrappableBalances

Gets wrappable balances for a user address and denom. Returns the amount as uint256.

**Signature:**
```solidity
function getWrappableBalances(
    string calldata denom,
    address userAddress
) external view returns (uint256);
```

**Parameters:**
- `denom` (string): The denomination (e.g., "stake")
- `userAddress` (address): The user address to query

**Returns:**
- `uint256`: Wrappable balance amount

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid user address
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### isAddressReservedProtocol

Checks if an address is reserved for protocol use. Returns a boolean.

**Signature:**
```solidity
function isAddressReservedProtocol(address addr) external view returns (bool);
```

**Parameters:**
- `addr` (address): The address to check

**Returns:**
- `bool`: `true` if the address is reserved protocol, `false` otherwise

**Gas Cost:** 2,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid address (zero address)

---

### getAllReservedProtocolAddresses

Gets all reserved protocol addresses. Returns an array of addresses.

**Signature:**
```solidity
function getAllReservedProtocolAddresses() external view returns (address[]);
```

**Parameters:** None

**Returns:**
- `address[]`: Array of all reserved protocol addresses

**Gas Cost:** 5,000 gas

**Error Codes:**
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getVote

Queries a vote for a specific proposal. Returns protobuf-encoded vote data.

**Signature:**
```solidity
function getVote(
    uint256 collectionId,
    string calldata approvalLevel,
    address approverAddress,
    string calldata approvalId,
    string calldata proposalId,
    address voterAddress
) external view returns (bytes);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approvalLevel` (string): Approval level
- `approverAddress` (address): Address of the approver
- `approvalId` (string): Approval ID
- `proposalId` (string): Proposal ID
- `voterAddress` (address): Address of the voter

**Returns:**
- `bytes`: Protobuf-encoded vote data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### getVotes

Queries all votes for a proposal. Returns protobuf-encoded votes data.

**Signature:**
```solidity
function getVotes(
    uint256 collectionId,
    string calldata approvalLevel,
    address approverAddress,
    string calldata approvalId,
    string calldata proposalId
) external view returns (bytes);
```

**Parameters:**
- `collectionId` (uint256): The collection ID
- `approvalLevel` (string): Approval level
- `approverAddress` (address): Address of the approver
- `approvalId` (string): Approval ID
- `proposalId` (string): Proposal ID

**Returns:**
- `bytes`: Protobuf-encoded votes data

**Gas Cost:** 5,000 gas

**Error Codes:**
- `1` (ErrorCodeInvalidInput): Invalid input parameters
- `6` (ErrorCodeQueryFailed): Query operation failed

---

### params

Queries the module parameters. Returns protobuf-encoded parameters data.

**Signature:**
```solidity
function params() external view returns (bytes);
```

**Parameters:** None

**Returns:**
- `bytes`: Protobuf-encoded module parameters

**Gas Cost:** 2,000 gas

**Error Codes:**
- `6` (ErrorCodeQueryFailed): Query operation failed

---

## Error Codes

All methods return structured errors with the following error codes:

| Code | Name | Description |
|------|------|-------------|
| 1 | ErrorCodeInvalidInput | Invalid input parameters (zero addresses, invalid ranges, negative values, etc.) |
| 2 | ErrorCodeCollectionNotFound | Collection does not exist |
| 3 | ErrorCodeBalanceNotFound | Balance not found |
| 4 | ErrorCodeTransferFailed | Transfer operation failed (insufficient balance, approval issues, etc.) |
| 5 | ErrorCodeApprovalFailed | Approval operation failed |
| 6 | ErrorCodeQueryFailed | Query operation failed |
| 7 | ErrorCodeInternalError | Internal error (marshaling, etc.) |
| 8 | ErrorCodeUnauthorized | Unauthorized operation |
| 9 | ErrorCodeCollectionArchived | Collection is archived (read-only) |

## Input Validation

All methods perform rigorous input validation:

- **Addresses**: Cannot be zero addresses
- **Collection IDs**: Must be non-negative (0 is valid for new collections)
- **Amounts**: Must be greater than zero
- **Ranges**: `start` must be <= `end`, both must be non-negative
- **Array Sizes**: 
  - Maximum 100 recipients per transfer
  - Maximum 100 token ID ranges
  - Maximum 100 ownership time ranges
- **Strings**: Cannot be empty where required

## Gas Costs Summary

### Transaction Methods
- `transferTokens`: 30,000 base + 5,000 per recipient + 1,000 per token ID range + 1,000 per ownership time range
- `setIncomingApproval`: 20,000 base + dynamic
- `setOutgoingApproval`: 20,000 base + dynamic

### Query Methods
- `getCollection`: 3,000 gas
- `getBalance`: 3,000 gas
- `getBalanceAmount`: 3,000 base + 500 per range
- `getTotalSupply`: 3,000 base + 500 per range
- `getAddressList`: 5,000 gas
- `getApprovalTracker`: 5,000 gas
- `getChallengeTracker`: 5,000 gas
- `getETHSignatureTracker`: 5,000 gas
- `getDynamicStore`: 5,000 gas
- `getDynamicStoreValue`: 5,000 gas
- `getWrappableBalances`: 5,000 gas
- `isAddressReservedProtocol`: 2,000 gas
- `getAllReservedProtocolAddresses`: 5,000 gas
- `getVote`: 5,000 gas
- `getVotes`: 5,000 gas
- `params`: 2,000 gas

## Protobuf Decoding

Most query methods return protobuf-encoded bytes. To decode these in Solidity, you'll need a protobuf decoding library. The following methods return direct values (no decoding needed):

- `getBalanceAmount`: Returns `uint256`
- `getTotalSupply`: Returns `uint256`
- `getChallengeTracker`: Returns `uint256`
- `getETHSignatureTracker`: Returns `uint256`
- `getWrappableBalances`: Returns `uint256`
- `isAddressReservedProtocol`: Returns `bool`
- `getAllReservedProtocolAddresses`: Returns `address[]`

## Best Practices

1. **Always check return values**: Transaction methods return `bool` indicating success
2. **Validate inputs before calling**: Check addresses, amounts, and ranges in your contract
3. **Handle errors gracefully**: Use try-catch blocks where appropriate
4. **Use direct return methods when possible**: Prefer `getBalanceAmount` and `getTotalSupply` over protobuf-encoded methods when you only need the amount
5. **Respect array size limits**: Keep recipient and range arrays under 100 elements
6. **Check collection existence**: Verify collection exists before attempting transfers

## Example: Complete Transfer Flow

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IBadgesPrecompile {
    struct UintRange {
        uint256 start;
        uint256 end;
    }
    
    function transferTokens(
        uint256 collectionId,
        address[] calldata toAddresses,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool);
    
    function getBalanceAmount(
        uint256 collectionId,
        address userAddress,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
}

contract MyTokenContract {
    IBadgesPrecompile constant BADGES_PRECOMPILE = IBadgesPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public collectionId;
    
    function transfer(uint256 amount, address to) external returns (bool) {
        IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
        tokenIds[0] = IBadgesPrecompile.UintRange({start: 1, end: 100});
        
        IBadgesPrecompile.UintRange[] memory ownershipTimes = new IBadgesPrecompile.UintRange[](1);
        ownershipTimes[0] = IBadgesPrecompile.UintRange({start: 0, end: type(uint256).max});
        
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        
        return BADGES_PRECOMPILE.transferTokens(
            collectionId,
            recipients,
            amount,
            tokenIds,
            ownershipTimes
        );
    }
    
    function balanceOf(address user) external view returns (uint256) {
        IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
        tokenIds[0] = IBadgesPrecompile.UintRange({start: 1, end: 100});
        
        IBadgesPrecompile.UintRange[] memory ownershipTimes = new IBadgesPrecompile.UintRange[](1);
        ownershipTimes[0] = IBadgesPrecompile.UintRange({start: 0, end: type(uint256).max});
        
        return BADGES_PRECOMPILE.getBalanceAmount(
            collectionId,
            user,
            tokenIds,
            ownershipTimes
        );
    }
}
```

