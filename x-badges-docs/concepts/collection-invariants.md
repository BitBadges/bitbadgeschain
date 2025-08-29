# Collection Invariants

Collection invariants are immutable rules that are set upon collection creation and cannot be broken or modified afterward. These invariants enforce fundamental constraints on how the collection operates, ensuring consistency and preventing certain types of restrictions.

## Overview

Invariants are collection-level properties that are set during genesis (collection creation) and remain fixed for the lifetime of the collection. Unlike permissions or other configurable settings, invariants cannot be updated or removed once established.

## Proto Definition

```protobuf
message CollectionInvariants {
  // If true, all ownership times must be full ranges [{ start: 1, end: GoMaxUInt64 }].
  // This prevents time-based restrictions on token ownership.
  bool noCustomOwnershipTimes = 1;

  // Maximum supply per token ID. If set, no balance can exceed this amount.
  // This prevents any single token ID from having more than the specified supply.
  string maxSupplyPerId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Available Invariants

### noCustomOwnershipTimes

When enabled, this invariant enforces that all ownership times throughout the collection must represent full ranges from the beginning of time (1) to the maximum possible time (18446744073709551615).

#### What it affects:

1. **Collection Approvals**: All collection-level approvals must have ownership times that are full ranges
2. **User Approvals**: All user-level approvals (incoming/outgoing) must have ownership times that are full ranges
3. **Transfer Balances**: All transfer operations must involve balances with full ownership time ranges

#### Validation Logic:

The invariant checks that ownership times are exactly:

```json
[{ "start": "1", "end": "18446744073709551615" }]
```

Any other ownership time configuration will cause validation to fail.

#### Use Cases:

-   **Preventing Time-Based Restrictions**: Ensures tokens cannot have time-limited ownership periods
-   **Simplifying Ownership Model**: Eliminates complexity around time-based approval restrictions
-   **Compliance Requirements**: Some applications may require permanent, unrestricted ownership

#### Example:

```json
// ✅ Valid - Full ownership time range
{
  "ownershipTimes": [
    {
      "start": "1",
      "end": "18446744073709551615"
    }
  ]
}

// ❌ Invalid - Restricted time range
{
  "ownershipTimes": [
    {
      "start": "1000",
      "end": "2000"
    }
  ]
}

// ❌ Invalid - Multiple ranges
{
  "ownershipTimes": [
    {
      "start": "1",
      "end": "1000"
    },
    {
      "start": "2000",
      "end": "3000"
    }
  ]
}
```

### maxSupplyPerId

When set to a non-zero value, this invariant enforces that no balance amount can exceed the specified maximum supply per token ID. This prevents supply inflation and ensures that the total supply of any individual token ID remains within the defined limits.

#### What it affects:

1. **Total Address Balances**: When setting balances for the "Total" address (which represents the total supply across all users), no individual balance amount can exceed the maximum
2. **Supply Control**: Prevents any single token ID from having more than the specified supply amount
3. **Collection Integrity**: Ensures the collection maintains its intended supply constraints

#### Validation Logic:

The invariant checks that when setting "Total" address balances, all balance amounts must be less than or equal to the specified `maxSupplyPerId`:

```go
if balance.Amount.GT(collection.Invariants.MaxSupplyPerId) {
    return error("maxSupplyPerId invariant violation")
}
```

#### Use Cases:

-   **Supply Caps**: Enforce maximum supply limits for individual token IDs
-   **Anti-Inflation**: Prevent supply manipulation through balance operations
-   **Compliance**: Meet regulatory requirements for maximum token supply
-   **Economic Control**: Maintain scarcity and value of specific token IDs

#### Example:

```json
// ✅ Valid - Balance amount within limit
{
  "invariants": {
    "maxSupplyPerId": "1000"
  },
  "balances": [
    {
      "amount": "500",
      "badgeIds": [{ "start": "1", "end": "1" }]
    }
  ]
}

// ❌ Invalid - Balance amount exceeds limit
{
  "invariants": {
    "maxSupplyPerId": "1000"
  },
  "balances": [
    {
      "amount": "1500",  // Exceeds maxSupplyPerId of 1000
      "badgeIds": [{ "start": "1", "end": "1" }]
    }
  ]
}

// ✅ Valid - Non-Total address not affected
{
  "invariants": {
    "maxSupplyPerId": "1000"
  },
  "balances": [
    {
      "amount": "2000",  // Allowed for non-Total addresses
      "badgeIds": [{ "start": "1", "end": "1" }]
    }
  ]
}
```

#### Important Notes:

-   **Only affects "Total" address**: The invariant only applies when setting balances for the "Total" address
-   **Zero value ignored**: If `maxSupplyPerId` is set to 0, the invariant is not enforced
-   **Immutable**: Once set during collection creation, this value cannot be changed
-   **Per-token ID basis**: The limit applies to each individual token ID, not the total collection supply

## Setting Invariants

Invariants can only be set during collection creation via `MsgCreateCollection` or `MsgUniversalUpdateCollection` (when creating a new collection with CollectionId = 0).

```protobuf
message MsgCreateCollection {
  // ... other fields ...
  CollectionInvariants invariants = 18;
}
```

## Validation Points

The invariants are validated at several points:

### noCustomOwnershipTimes

1. **Collection Creation**: When creating a new collection
2. **Collection Updates**: When updating collection approvals
3. **Transfer Execution**: When processing token transfers
4. **Approval Updates**: When updating user or collection approvals

### maxSupplyPerId

1. **Balance Storage**: When setting user balances in the store (specifically for "Total" address)
2. **Supply Validation**: Before any balance amount is stored that would exceed the maximum

## Error Messages

When invariants are violated, you'll receive error messages like:

### noCustomOwnershipTimes

```
noCustomOwnershipTimes invariant is enabled: ownership times must be full range [{ start: 1, end: 18446744073709551615 }]
```

### maxSupplyPerId

```
maxSupplyPerId invariant violation: balance amount 1500 exceeds maximum supply per ID 1000
```

## Related Concepts

-   [Collections](./badge-collections.md)
-   [Transferability Approvals](./transferability-approvals.md)
-   [Time Fields](./time-fields.md)
-   [UintRange](./uintrange.md)
