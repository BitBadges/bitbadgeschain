# EVM Precompile Usage Examples

This document provides practical Solidity examples for using the BitBadges tokenization precompile.

## Table of Contents

1. [Basic Setup](#basic-setup)
2. [Token Transfers](#token-transfers)
3. [Balance Queries](#balance-queries)
4. [Approval Management](#approval-management)
5. [ERC-3643 Wrapper](#erc-3643-wrapper)
6. [Advanced Patterns](#advanced-patterns)
7. [Error Handling](#error-handling)
8. [Gas Optimization](#gas-optimization)

## Basic Setup

### Interface Definition

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IBadgesPrecompile {
    struct UintRange {
        uint256 start;
        uint256 end;
    }
    
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
    
    // Transaction methods
    function transferTokens(
        uint256 collectionId,
        address[] calldata toAddresses,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool);
    
    function setIncomingApproval(
        uint256 collectionId,
        UserIncomingApproval calldata approval
    ) external returns (bool);
    
    function setOutgoingApproval(
        uint256 collectionId,
        UserOutgoingApproval calldata approval
    ) external returns (bool);
    
    // Query methods
    function getBalanceAmount(
        uint256 collectionId,
        address userAddress,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
    
    function getTotalSupply(
        uint256 collectionId,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
    
    function getCollection(uint256 collectionId) external view returns (bytes);
    function getBalance(uint256 collectionId, address userAddress) external view returns (bytes);
    function params() external view returns (bytes);
}

// Precompile address
address constant BADGES_PRECOMPILE = 0x0000000000000000000000000000000000001001;
```

### Basic Contract Template

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IBadgesPrecompile.sol";

contract BadgesContract {
    IBadgesPrecompile public constant precompile = IBadgesPrecompile(BADGES_PRECOMPILE);
    uint256 public collectionId;
    
    constructor(uint256 _collectionId) {
        collectionId = _collectionId;
    }
    
    // Helper to create full range
    function getFullRange() internal pure returns (IBadgesPrecompile.UintRange[] memory) {
        IBadgesPrecompile.UintRange[] memory ranges = new IBadgesPrecompile.UintRange[](1);
        ranges[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
        return ranges;
    }
}
```

## Token Transfers

### Simple Transfer to One Recipient

```solidity
function transfer(address to, uint256 amount) external returns (bool) {
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    
    IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.transferTokens(
        collectionId,
        recipients,
        amount,
        tokenIds,
        ownershipTimes
    );
}
```

### Transfer to Multiple Recipients

```solidity
function transferBatch(address[] calldata recipients, uint256 amount) external returns (bool) {
    require(recipients.length > 0 && recipients.length <= 100, "Invalid recipients");
    
    IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.transferTokens(
        collectionId,
        recipients,
        amount,
        tokenIds,
        ownershipTimes
    );
}
```

### Transfer Specific Token IDs

```solidity
function transferSpecificTokens(
    address to,
    uint256 amount,
    uint256 tokenIdStart,
    uint256 tokenIdEnd
) external returns (bool) {
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    
    IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
    tokenIds[0] = IBadgesPrecompile.UintRange({
        start: tokenIdStart,
        end: tokenIdEnd
    });
    
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.transferTokens(
        collectionId,
        recipients,
        amount,
        tokenIds,
        ownershipTimes
    );
}
```

### Transfer with Multiple Token ID Ranges

```solidity
function transferMultipleRanges(
    address to,
    uint256 amount,
    uint256[] calldata tokenIdStarts,
    uint256[] calldata tokenIdEnds
) external returns (bool) {
    require(tokenIdStarts.length == tokenIdEnds.length, "Mismatched arrays");
    require(tokenIdStarts.length <= 100, "Too many ranges");
    
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    
    IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](tokenIdStarts.length);
    for (uint256 i = 0; i < tokenIdStarts.length; i++) {
        require(tokenIdStarts[i] <= tokenIdEnds[i], "Invalid range");
        tokenIds[i] = IBadgesPrecompile.UintRange({
            start: tokenIdStarts[i],
            end: tokenIdEnds[i]
        });
    }
    
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.transferTokens(
        collectionId,
        recipients,
        amount,
        tokenIds,
        ownershipTimes
    );
}
```

## Balance Queries

### Get Balance for User

```solidity
function balanceOf(address user) external view returns (uint256) {
    IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.getBalanceAmount(
        collectionId,
        user,
        tokenIds,
        ownershipTimes
    );
}
```

### Get Balance for Specific Token IDs

```solidity
function balanceOfTokens(
    address user,
    uint256 tokenIdStart,
    uint256 tokenIdEnd
) external view returns (uint256) {
    IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
    tokenIds[0] = IBadgesPrecompile.UintRange({
        start: tokenIdStart,
        end: tokenIdEnd
    });
    
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.getBalanceAmount(
        collectionId,
        user,
        tokenIds,
        ownershipTimes
    );
}
```

### Get Total Supply

```solidity
function totalSupply() external view returns (uint256) {
    IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.getTotalSupply(
        collectionId,
        tokenIds,
        ownershipTimes
    );
}
```

### Get Total Supply for Specific Token IDs

```solidity
function totalSupplyForTokens(
    uint256 tokenIdStart,
    uint256 tokenIdEnd
) external view returns (uint256) {
    IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
    tokenIds[0] = IBadgesPrecompile.UintRange({
        start: tokenIdStart,
        end: tokenIdEnd
    });
    
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    return precompile.getTotalSupply(
        collectionId,
        tokenIds,
        ownershipTimes
    );
}
```

## Approval Management

### Set Incoming Approval (Allow Others to Send to You)

```solidity
function allowIncomingFromAll() external returns (bool) {
    IBadgesPrecompile.UserIncomingApproval memory approval = IBadgesPrecompile.UserIncomingApproval({
        approvalId: "incoming_all",
        fromListId: "All",
        initiatedByListId: "All",
        transferTimes: new IBadgesPrecompile.UintRange[](0), // Empty = all times
        tokenIds: new IBadgesPrecompile.UintRange[](0), // Empty = all token IDs
        ownershipTimes: new IBadgesPrecompile.UintRange[](0), // Empty = all ownership times
        uri: "",
        customData: ""
    });
    
    return precompile.setIncomingApproval(collectionId, approval);
}
```

### Set Outgoing Approval (Allow Sending to Specific Addresses)

```solidity
function allowOutgoingTo(address to) external returns (bool) {
    // Create address list first (via Cosmos SDK), then reference it
    // For this example, we'll use "All" which allows transfers to any address
    IBadgesPrecompile.UserOutgoingApproval memory approval = IBadgesPrecompile.UserOutgoingApproval({
        approvalId: "outgoing_all",
        toListId: "All",
        initiatedByListId: "All",
        transferTimes: new IBadgesPrecompile.UintRange[](0),
        tokenIds: new IBadgesPrecompile.UintRange[](0),
        ownershipTimes: new IBadgesPrecompile.UintRange[](0),
        uri: "",
        customData: ""
    });
    
    return precompile.setOutgoingApproval(collectionId, approval);
}
```

### Set Approval with Specific Ranges

```solidity
function setRestrictedApproval(
    uint256 tokenIdStart,
    uint256 tokenIdEnd,
    uint256 timeStart,
    uint256 timeEnd
) external returns (bool) {
    IBadgesPrecompile.UintRange[] memory tokenIds = new IBadgesPrecompile.UintRange[](1);
    tokenIds[0] = IBadgesPrecompile.UintRange({
        start: tokenIdStart,
        end: tokenIdEnd
    });
    
    IBadgesPrecompile.UintRange[] memory ownershipTimes = new IBadgesPrecompile.UintRange[](1);
    ownershipTimes[0] = IBadgesPrecompile.UintRange({
        start: timeStart,
        end: timeEnd
    });
    
    IBadgesPrecompile.UserOutgoingApproval memory approval = IBadgesPrecompile.UserOutgoingApproval({
        approvalId: "restricted_approval",
        toListId: "All",
        initiatedByListId: "All",
        transferTimes: new IBadgesPrecompile.UintRange[](0),
        tokenIds: tokenIds,
        ownershipTimes: ownershipTimes,
        uri: "",
        customData: ""
    });
    
    return precompile.setOutgoingApproval(collectionId, approval);
}
```

## ERC-3643 Wrapper

### Complete ERC-3643 Implementation

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IBadgesPrecompile.sol";

contract ERC3643Badges {
    IBadgesPrecompile public constant precompile = IBadgesPrecompile(BADGES_PRECOMPILE);
    
    uint256 public collectionId;
    string public name;
    string public symbol;
    uint8 public decimals;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    
    constructor(
        uint256 _collectionId,
        string memory _name,
        string memory _symbol,
        uint8 _decimals
    ) {
        collectionId = _collectionId;
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
    }
    
    function totalSupply() external view returns (uint256) {
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        return precompile.getTotalSupply(collectionId, tokenIds, ownershipTimes);
    }
    
    function balanceOf(address account) external view returns (uint256) {
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        return precompile.getBalanceAmount(collectionId, account, tokenIds, ownershipTimes);
    }
    
    function transfer(address to, uint256 amount) external returns (bool) {
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        bool success = precompile.transferTokens(
            collectionId,
            recipients,
            amount,
            tokenIds,
            ownershipTimes
        );
        
        if (success) {
            emit Transfer(msg.sender, to, amount);
        }
        
        return success;
    }
    
    function getFullRange() internal pure returns (IBadgesPrecompile.UintRange[] memory) {
        IBadgesPrecompile.UintRange[] memory ranges = new IBadgesPrecompile.UintRange[](1);
        ranges[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
        return ranges;
    }
}
```

## Advanced Patterns

### Batch Operations

```solidity
function batchTransfer(
    address[] calldata recipients,
    uint256[] calldata amounts
) external returns (bool) {
    require(recipients.length == amounts.length, "Mismatched arrays");
    require(recipients.length <= 100, "Too many recipients");
    
    IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
    IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
    
    bool allSuccess = true;
    for (uint256 i = 0; i < recipients.length; i++) {
        address[] memory singleRecipient = new address[](1);
        singleRecipient[0] = recipients[i];
        
        bool success = precompile.transferTokens(
            collectionId,
            singleRecipient,
            amounts[i],
            tokenIds,
            ownershipTimes
        );
        
        if (!success) {
            allSuccess = false;
        }
    }
    
    return allSuccess;
}
```

### Conditional Transfer

```solidity
function transferIfBalanceSufficient(
    address to,
    uint256 amount
) external returns (bool) {
    // Check balance first
    uint256 balance = balanceOf(msg.sender);
    require(balance >= amount, "Insufficient balance");
    
    // Perform transfer
    return transfer(to, amount);
}
```

### Multi-Collection Contract

```solidity
contract MultiCollectionBadges {
    IBadgesPrecompile public constant precompile = IBadgesPrecompile(BADGES_PRECOMPILE);
    
    mapping(uint256 => string) public collectionNames;
    uint256[] public collectionIds;
    
    function addCollection(uint256 collectionId, string memory name) external {
        collectionNames[collectionId] = name;
        collectionIds.push(collectionId);
    }
    
    function transferFromCollection(
        uint256 collectionId,
        address to,
        uint256 amount
    ) external returns (bool) {
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        return precompile.transferTokens(
            collectionId,
            recipients,
            amount,
            tokenIds,
            ownershipTimes
        );
    }
    
    function getFullRange() internal pure returns (IBadgesPrecompile.UintRange[] memory) {
        IBadgesPrecompile.UintRange[] memory ranges = new IBadgesPrecompile.UintRange[](1);
        ranges[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
        return ranges;
    }
}
```

## Error Handling

### Try-Catch Pattern

```solidity
function safeTransfer(address to, uint256 amount) external {
    try precompile.transferTokens(...) returns (bool success) {
        require(success, "Transfer returned false");
    } catch Error(string memory reason) {
        // Handle error with reason
        revert(reason);
    } catch (bytes memory lowLevelData) {
        // Handle low-level error
        revert("Transfer failed");
    }
}
```

### Error Code Parsing

```solidity
function parseErrorCode(bytes memory errorData) internal pure returns (uint256) {
    // Error format: "precompile error [code=X]: message"
    // This is a simplified parser - in production, use a proper library
    if (errorData.length < 20) return 0;
    
    // Look for "[code=" pattern
    for (uint256 i = 0; i < errorData.length - 10; i++) {
        if (errorData[i] == '[' && 
            errorData[i+1] == 'c' && 
            errorData[i+2] == 'o' && 
            errorData[i+3] == 'd' && 
            errorData[i+4] == 'e' && 
            errorData[i+5] == '=') {
            // Extract code number
            uint256 code = 0;
            uint256 j = i + 6;
            while (j < errorData.length && errorData[j] >= '0' && errorData[j] <= '9') {
                code = code * 10 + (uint256(uint8(errorData[j])) - 48);
                j++;
            }
            return code;
        }
    }
    return 0;
}
```

### Validation Before Transfer

```solidity
function validatedTransfer(
    address to,
    uint256 amount
) external returns (bool) {
    require(to != address(0), "Cannot transfer to zero address");
    require(amount > 0, "Amount must be greater than zero");
    
    uint256 balance = balanceOf(msg.sender);
    require(balance >= amount, "Insufficient balance");
    
    return transfer(to, amount);
}
```

## Gas Optimization

### Reuse Range Arrays

```solidity
contract OptimizedBadges {
    IBadgesPrecompile public constant precompile = IBadgesPrecompile(BADGES_PRECOMPILE);
    uint256 public collectionId;
    
    // Cache full range arrays
    IBadgesPrecompile.UintRange[] private fullTokenRange;
    IBadgesPrecompile.UintRange[] private fullOwnershipRange;
    
    constructor(uint256 _collectionId) {
        collectionId = _collectionId;
        
        // Initialize cached ranges
        fullTokenRange = new IBadgesPrecompile.UintRange[](1);
        fullTokenRange[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
        
        fullOwnershipRange = new IBadgesPrecompile.UintRange[](1);
        fullOwnershipRange[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
    }
    
    function transfer(address to, uint256 amount) external returns (bool) {
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        
        // Use cached ranges instead of creating new ones
        return precompile.transferTokens(
            collectionId,
            recipients,
            amount,
            fullTokenRange,
            fullOwnershipRange
        );
    }
}
```

### Batch Range Creation

```solidity
function createRanges(
    uint256[] calldata starts,
    uint256[] calldata ends
) internal pure returns (IBadgesPrecompile.UintRange[] memory) {
    require(starts.length == ends.length, "Mismatched arrays");
    
    IBadgesPrecompile.UintRange[] memory ranges = new IBadgesPrecompile.UintRange[](starts.length);
    for (uint256 i = 0; i < starts.length; i++) {
        require(starts[i] <= ends[i], "Invalid range");
        ranges[i] = IBadgesPrecompile.UintRange({
            start: starts[i],
            end: ends[i]
        });
    }
    return ranges;
}
```

## Complete Example: Token Marketplace

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IBadgesPrecompile.sol";

contract BadgesMarketplace {
    IBadgesPrecompile public constant precompile = IBadgesPrecompile(BADGES_PRECOMPILE);
    
    struct Listing {
        uint256 collectionId;
        address seller;
        uint256 amount;
        uint256 price;
        bool active;
    }
    
    mapping(uint256 => Listing) public listings;
    uint256 public listingCounter;
    
    event Listed(uint256 indexed listingId, address indexed seller, uint256 collectionId, uint256 amount, uint256 price);
    event Sold(uint256 indexed listingId, address indexed buyer);
    
    function list(
        uint256 collectionId,
        uint256 amount,
        uint256 price
    ) external returns (uint256) {
        // Transfer tokens to this contract
        address[] memory recipients = new address[](1);
        recipients[0] = address(this);
        
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        bool success = precompile.transferTokens(
            collectionId,
            recipients,
            amount,
            tokenIds,
            ownershipTimes
        );
        require(success, "Transfer to marketplace failed");
        
        // Create listing
        uint256 listingId = listingCounter++;
        listings[listingId] = Listing({
            collectionId: collectionId,
            seller: msg.sender,
            amount: amount,
            price: price,
            active: true
        });
        
        emit Listed(listingId, msg.sender, collectionId, amount, price);
        return listingId;
    }
    
    function buy(uint256 listingId) external payable {
        Listing storage listing = listings[listingId];
        require(listing.active, "Listing not active");
        require(msg.value >= listing.price, "Insufficient payment");
        
        // Transfer tokens to buyer
        address[] memory recipients = new address[](1);
        recipients[0] = msg.sender;
        
        IBadgesPrecompile.UintRange[] memory tokenIds = getFullRange();
        IBadgesPrecompile.UintRange[] memory ownershipTimes = getFullRange();
        
        bool success = precompile.transferTokens(
            listing.collectionId,
            recipients,
            listing.amount,
            tokenIds,
            ownershipTimes
        );
        require(success, "Transfer to buyer failed");
        
        // Transfer payment to seller
        payable(listing.seller).transfer(listing.price);
        
        // Refund excess
        if (msg.value > listing.price) {
            payable(msg.sender).transfer(msg.value - listing.price);
        }
        
        listing.active = false;
        emit Sold(listingId, msg.sender);
    }
    
    function getFullRange() internal pure returns (IBadgesPrecompile.UintRange[] memory) {
        IBadgesPrecompile.UintRange[] memory ranges = new IBadgesPrecompile.UintRange[](1);
        ranges[0] = IBadgesPrecompile.UintRange({
            start: 0,
            end: type(uint256).max
        });
        return ranges;
    }
}
```

## Best Practices Summary

1. **Always validate inputs** before calling precompile methods
2. **Use direct return methods** (`getBalanceAmount`, `getTotalSupply`) when possible
3. **Cache range arrays** to save gas on repeated calls
4. **Handle errors gracefully** with try-catch blocks
5. **Respect array size limits** (max 100 elements)
6. **Check balances** before attempting transfers
7. **Use events** to track important state changes
8. **Optimize gas** by reusing arrays and minimizing ranges

