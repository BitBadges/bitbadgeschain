# 🗃️ State

This document describes how the badges module manages state storage, including key structures and data organization.

## Table of Contents

-   [State Store Overview](#state-store-overview)
-   [Key Generation](#key-generation)
-   [Data Structures](#data-structures)

## State Store Overview

The badges module uses the Cosmos SDK's KVStore with a prefix-based key structure to organize different types of data efficiently. All data is stored in a single store with logical separation through key prefixes.

### Store Architecture

```
KVStore (badges)
├── Collections [0x01]
├── User Balances [0x02]
├── Next Collection ID [0x03]
├── Challenge Trackers [0x04]
├── Address Lists [0x06]
├── Approval Trackers [0x07]
├── Next Address List ID [0x0A]
├── Dynamic Stores [0x0D]
├── Next Dynamic Store ID [0x0E]
└── Dynamic Store Values [0x0F]
```

## Key Generation

The badges module uses deterministic key generation with prefixes and delimiters to ensure unique storage keys.

### Key Prefixes

```go
var (
    CollectionKey           = []byte{0x01}
    UserBalanceKey          = []byte{0x02}
    NextCollectionIdKey     = []byte{0x03}
    UsedClaimChallengeKey   = []byte{0x04}
    AddressListKey          = []byte{0x06}
    ApprovalTrackerKey      = []byte{0x07}
    NextAddressListIdKey    = []byte{0x0A}
    DynamicStoreKey         = []byte{0x0D}
    NextDynamicStoreIdKey   = []byte{0x0E}
    DynamicStoreValueKey    = []byte{0x0F}

    BalanceKeyDelimiter = "-"
)
```

## Data Structures

### Collections

-   **Key**: `0x01 + collectionId`
-   **Value**: `BadgeCollection` protobuf message
-   **Description**: Complete collection metadata, permissions, and configuration

### User Balances

-   **Key**: `0x02 + collectionId + "-" + address`
-   **Value**: `UserBalanceStore` protobuf message
-   **Description**: User's badge balances and approval settings for a collection

### Address Lists

-   **Key**: `0x06 + listId`
-   **Value**: `AddressList` protobuf message
-   **Description**: Reusable address lists for access control and approvals

### Approval Trackers

-   **Key**: `0x07 + collectionId + "-" + addressForApproval + "-" + approvalId + "-" + amountTrackerId + "-" + level + "-" + trackerType + "-" + address`
-   **Value**: `ApprovalTracker` protobuf message
-   **Description**: Tracks usage of approvals for transfer limits and restrictions

### Challenge Trackers

-   **Key**: `0x04 + collectionId + "-" + addressForChallenge + "-" + approvalLevel + "-" + approvalId + "-" + challengeId + "-" + leafIndex`
-   **Value**: Usage count (uint64)
-   **Description**: Tracks merkle proof leaf usage to prevent replay attacks

### Dynamic Stores

-   **Key**: `0x0D + storeId`
-   **Value**: `DynamicStore` protobuf message
-   **Description**: Dynamic store configuration and metadata

### Dynamic Store Values

-   **Key**: `0x0F + storeId + address`
-   **Value**: Boolean value
-   **Description**: Address-specific boolean values within dynamic stores

### Counter Keys

-   **Next Collection ID**: `0x03` → Current collection counter
-   **Next Address List ID**: `0x0A` → Current address list counter
-   **Next Dynamic Store ID**: `0x0E` → Current dynamic store counter
