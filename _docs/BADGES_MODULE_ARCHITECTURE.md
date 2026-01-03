# Module Architecture Overview

This document provides a high-level architectural overview of the BitBadges `x/badges` module, focusing on design patterns, data flow, and integration points.

## Table of Contents

-   [Architectural Principles](#architectural-principles)
-   [Core Components](#core-components)
-   [Data Flow Patterns](#data-flow-patterns)
-   [State Management Architecture](#state-management-architecture)
-   [Cross-Chain Integration Architecture](#cross-chain-integration-architecture)
-   [Scalability Considerations](#scalability-considerations)
-   [Security Architecture](#security-architecture)

## Architectural Principles

### 1. Timeline-Driven Design

The module is built around the concept that most properties can change over time, providing:

-   **Immutable History**: Once a timeline point has passed, it cannot be modified
-   **Scheduled Changes**: Future property changes can be pre-programmed
-   **Predictable Behavior**: Users can rely on announced future changes

```
Timeline Architecture:
┌─────────────────────────────────────────────────┐
│ Property Timeline                               │
├─────────────────────────────────────────────────┤
│ Time: 0────1000────2000────3000────∞           │
│ Value: A    B       C       D       D           │
└─────────────────────────────────────────────────┘
```

### 2. Three-Tier Approval System

All token transfers must satisfy three independent approval levels:

```
Transfer Approval Flow:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Collection      │    │ Outgoing        │    │ Incoming        │
│ Approval        │───▶│ Approval        │───▶│ Approval        │
│ (Manager Sets)  │    │ (Sender Sets)   │    │ (Recipient Sets)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
    ✅ ALLOW                 ✅ ALLOW                 ✅ ALLOW
                                    │
                                    ▼
                              Transfer Succeeds
```

### 3. Range-Based Efficiency

Everything is expressed as ranges (UintRange) for computational efficiency:

-   Token IDs: `[1-1000]` instead of 1000 individual IDs
-   Time ranges: `[start_time-end_time]` for temporal validity
-   Amount ranges: `[min_amount-max_amount]` for quantities

### 4. Default Inheritance Pattern

Users inherit collection defaults until they explicitly override them:

-   Minimizes storage requirements
-   Provides consistent default behavior
-   Lazy initialization for state efficiency

## Core Components

### 1. Collection Manager

**Responsibility**: Manages the lifecycle and properties of collections

```go
type CollectionManager struct {
    // Timeline-based properties
    ManagerTimeline        []ManagerTimeline
    MetadataTimeline       []CollectionMetadataTimeline
    PermissionsTimeline    []CollectionPermissions

    // Static properties
    CollectionId           Uint
    BalanceType           BalanceType
    ValidTokenIds         []UintRange
}
```

**Key Functions**:

-   Collection creation and updates
-   Timeline property management
-   Permission enforcement
-   Archive status control

### 2. Balance Engine

**Responsibility**: Tracks token ownership with temporal and quantity precision

```go
type BalanceEngine struct {
    // User balances with full temporal context
    UserBalances map[string]UserBalanceStore

    // Collection defaults
    DefaultBalances UserBalanceStore
}
```

**Key Functions**:

-   Balance queries with time filtering
-   Transfer execution and validation
-   Default inheritance resolution
-   Balance versioning for approval invalidation

### 3. Approval Processor

**Responsibility**: Evaluates complex approval criteria for transfers

```go
type ApprovalProcessor struct {
    // Three-tier approval system
    CollectionApprovals   []CollectionApproval
    OutgoingApprovals     []UserOutgoingApproval
    IncomingApprovals     []UserIncomingApproval

    // Usage tracking
    ApprovalTrackers      map[string]ApprovalTracker
    ChallengeTrackers     map[string]ChallengeTracker
}
```

**Key Functions**:

-   Multi-level approval validation
-   Merkle challenge verification
-   Usage limit enforcement
-   Approval criteria evaluation

### 4. Address List Manager

**Responsibility**: Manages reusable address collections

```go
type AddressListManager struct {
    AddressLists map[Uint]AddressList
}
```

**Key Functions**:

-   Address list creation and management
-   Whitelist/blacklist logic evaluation
-   Address pattern matching (wildcards, etc.)

### 5. Dynamic Store Engine

**Global Kill Switch Feature**: All dynamic stores include a `globalEnabled` field that acts as a global kill switch. When `globalEnabled = false`, all approvals using that store via `DynamicStoreChallenge` will fail immediately, regardless of per-address values. This enables quick halting of approvals (e.g., when a 2FA protocol is compromised). The field defaults to `true` on creation and can be toggled via `MsgUpdateDynamicStore`.

**Responsibility**: Provides flexible key-value storage for custom logic

```go
type DynamicStoreEngine struct {
    Stores map[Uint]DynamicStore
    Values map[string]bool  // storeId+address -> boolean
}
```

**Key Functions**:

-   Store creation and management
-   Value setting and retrieval
-   Integration with approval criteria

## Data Flow Patterns

### 1. Collection Creation Flow

```
User Request
     │
     ▼
┌─────────────────┐
│ Message         │
│ Validation      │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Permission      │
│ Check           │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Timeline        │
│ Validation      │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ State Update    │
│ (Atomic)        │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Event Emission  │
└─────────────────┘
```

### 2. Transfer Execution Flow

```
Transfer Request
     │
     ▼
┌─────────────────┐
│ Basic           │
│ Validation      │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Collection      │
│ Approval Check  │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Outgoing        │
│ Approval Check  │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Incoming        │
│ Approval Check  │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Balance Update  │
│ (Atomic)        │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Tracker Update  │
└─────────────────┘
```

### 3. Query Resolution Flow

```
Query Request
     │
     ▼
┌─────────────────┐
│ Parameter       │
│ Validation      │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ State Lookup    │
└─────────────────┘
     │
     ▼
┌─────────────────┐    ┌─────────────────┐
│ Found?          │───▶│ Return Data     │
└─────────────────┘    └─────────────────┘
     │
     ▼
┌─────────────────┐
│ Check Defaults  │
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Return Defaults │
└─────────────────┘
```

## State Management Architecture

### 1. Layered Storage Model

```
Application Layer
├── Collection Management
├── Balance Tracking
├── Approval Processing
└── Dynamic Stores

Keeper Layer
├── State Validation
├── Permission Enforcement
├── Event Emission
└── Cross-cutting Concerns

Storage Layer
├── KV Store (Primary)
├── Index Management
├── Cache Layer
└── Migration Support
```

### 2. Key-Value Store Organization

```
Prefix-Based Key Structure:
┌─────────────────────────────────────────────────┐
│ 0x01 │ Collections                               │
├─────────────────────────────────────────────────┤
│ 0x02 │ User Balances                             │
├─────────────────────────────────────────────────┤
│ 0x03 │ Address Lists                             │
├─────────────────────────────────────────────────┤
│ 0x04 │ Approval Trackers                         │
├─────────────────────────────────────────────────┤
│ 0x07 │ Dynamic Stores                            │
├─────────────────────────────────────────────────┤
│ 0x08 │ Dynamic Store Values                      │
└─────────────────────────────────────────────────┘
```

### 3. Composite Key Construction

```go
// Example: User Balance Key
userBalanceKey := concat(
    UserBalancePrefix,     // 0x02
    collectionId,          // "1"
    address,              // "bb1abc123..."
)

// Example: Approval Tracker Key
approvalTrackerKey := concat(
    ApprovalTrackerPrefix, // 0x04
    collectionId,          // "1"
    approvalLevel,         // "outgoing"
    trackerType,           // "amount"
    approverAddress,       // "bb1abc123..."
    trackerId,            // "tracker1"
)
```

### 4. Atomic Operation Patterns

All state modifications are atomic within transaction boundaries:

```go
func (k Keeper) ExecuteTransfer(ctx sdk.Context, transfer Transfer) error {
    // Phase 1: Validation (read-only)
    if err := k.ValidateTransfer(ctx, transfer); err != nil {
        return err
    }

    // Phase 2: State modifications (atomic)
    return k.executeInTransaction(ctx, func() error {
        // Update sender balances
        if err := k.UpdateSenderBalance(ctx, transfer); err != nil {
            return err
        }

        // Update recipient balances
        if err := k.UpdateRecipientBalance(ctx, transfer); err != nil {
            return err
        }

        // Update approval trackers
        if err := k.UpdateApprovalTrackers(ctx, transfer); err != nil {
            return err
        }

        // All succeed or all fail
        return nil
    })
}
```

## Cross-Chain Integration Architecture

### 1. Multi-Chain Signature Support

```
Signature Verification Layer:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Ethereum        │    │ Bitcoin         │    │ Solana          │
│ (EIP712)        │    │ (JSON Schema)   │    │ (JSON Schema)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│ Unified Verification Engine                                     │
├─────────────────────────────────────────────────────────────────┤
│ • Schema validation                                             │
│ • Signature verification                                        │
│ • Address format handling                                       │
│ • Cross-chain message routing                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2. IBC Integration

```
IBC Packet Flow:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Source Chain    │    │ IBC Relay       │    │ Destination     │
│ (BitBadges)     │───▶│ Infrastructure  │───▶│ Chain          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                                               │
        ▼                                               ▼
┌─────────────────┐                         ┌─────────────────┐
│ Token Lock      │                         │ Token Unlock    │
│ • Escrow tokens │                         │ • Mint/Release  │
│ • Create packet │                         │ • Verify packet │
│ • Send via IBC  │                         │ • Update state  │
└─────────────────┘                         └─────────────────┘
```

### 3. WASM Smart Contract Integration

```
WASM Binding Architecture:
┌─────────────────────────────────────────────────────────────────┐
│ Smart Contract Layer                                            │
├─────────────────────────────────────────────────────────────────┤
│ Custom Query Bindings    │ Custom Message Bindings             │
│ • Get collection info    │ • Create collection                  │
│ • Query balances         │ • Transfer tokens                    │
│ • Check approvals        │ • Update approvals                   │
│ • Dynamic store queries  │ • Set dynamic store values           │
└─────────────────────────────────────────────────────────────────┘
        │                               │
        ▼                               ▼
┌─────────────────┐         ┌─────────────────┐
│ Query Handler   │         │ Message Router  │
│ (Read-only)     │         │ (State Changes) │
└─────────────────┘         └─────────────────┘
```

## Scalability Considerations

### 1. Range-Based Computation

Instead of iterating individual elements:

-   Token IDs stored as ranges: O(1) vs O(n)
-   Time ranges for temporal queries: O(log n) vs O(n)
-   Amount ranges for balance operations: O(1) vs O(n)

### 2. Lazy State Initialization

```
User Balance Lifecycle:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ No explicit     │    │ First           │    │ Explicit        │
│ balance         │───▶│ modification    │───▶│ balance store   │
│ (inherits       │    │ (create store)  │    │ (independent)   │
│ defaults)       │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 3. Efficient Query Patterns

-   **Direct lookups**: O(1) for collections, balances by ID
-   **Prefix iteration**: O(k) where k = number of results
-   **Binary search**: O(log n) for range queries
-   **Pagination support**: Prevents large result set issues

### 4. State Pruning Strategy

```
State Retention Policy:
┌─────────────────────────────────────────────────────────────────┐
│ Current State: Always retained                                  │
├─────────────────────────────────────────────────────────────────┤
│ Recent History: Configurable retention period                   │
├─────────────────────────────────────────────────────────────────┤
│ Archived Data: Optional long-term storage                       │
├─────────────────────────────────────────────────────────────────┤
│ Pruned Data: Removed to reduce storage burden                   │
└─────────────────────────────────────────────────────────────────┘
```

## Security Architecture

### 1. Permission Model

```
Permission Hierarchy:
┌─────────────────────────────────────────────────────────────────┐
│ Module Authority (Governance)                                   │
├─────────────────────────────────────────────────────────────────┤
│ Collection Manager                                              │
│ • Collection-level permissions                                  │
│ • Timeline modifications                                        │
│ • Approval configurations                                       │
├─────────────────────────────────────────────────────────────────┤
│ Individual Users                                                │
│ • Personal approval settings                                    │
│ • Token transfers                                               │
│ • Dynamic store interactions                                    │
└─────────────────────────────────────────────────────────────────┘
```

### 2. Validation Layers

```
Validation Stack:
┌─────────────────┐
│ Message Format  │ ← Basic format and required fields
├─────────────────┤
│ Business Rules  │ ← Timeline, range, and logic validation
├─────────────────┤
│ Permission      │ ← Authorization and access control
├─────────────────┤
│ State           │ ← Consistency and integrity checks
└─────────────────┘
```

### 3. Replay Attack Prevention

-   **Approval versioning**: Increments on approval changes
-   **Merkle challenge tracking**: Per-leaf usage counters
-   **Transaction nonces**: Standard Cosmos SDK protection
-   **Timeline immutability**: Prevents historical manipulation

### 4. Economic Security

-   **Gas metering**: All operations have appropriate gas costs
-   **Spam prevention**: Creation costs and rate limiting
-   **Resource bounds**: Practical limits on collection sizes
-   **Fee mechanisms**: Optional protocol fees for sustainability

This architecture provides a robust foundation for the complex token management system while maintaining performance, security, and extensibility for future enhancements.
