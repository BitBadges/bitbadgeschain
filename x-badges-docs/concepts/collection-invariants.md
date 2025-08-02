# Collection Invariants

Collection invariants are immutable rules that are set upon collection creation and cannot be broken or modified afterward. These invariants enforce fundamental constraints on how the collection operates, ensuring consistency and preventing certain types of restrictions.

## Overview

Invariants are collection-level properties that are set during genesis (collection creation) and remain fixed for the lifetime of the collection. Unlike permissions or other configurable settings, invariants cannot be updated or removed once established.

## Proto Definition

```protobuf
message CollectionInvariants {
  // If true, all ownership times must be full ranges [{ start: 1, end: GoMaxUInt64 }].
  // This prevents time-based restrictions on badge ownership.
  bool noCustomOwnershipTimes = 1;
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

-   **Preventing Time-Based Restrictions**: Ensures badges cannot have time-limited ownership periods
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

## Setting Invariants

Invariants can only be set during collection creation via `MsgCreateCollection` or `MsgUniversalUpdateCollection` (when creating a new collection with CollectionId = 0).

```protobuf
message MsgCreateCollection {
  // ... other fields ...
  CollectionInvariants invariants = 18;
}
```

## Validation Points

The `noCustomOwnershipTimes` invariant is validated at several points:

1. **Collection Creation**: When creating a new collection
2. **Collection Updates**: When updating collection approvals
3. **Transfer Execution**: When processing badge transfers
4. **Approval Updates**: When updating user or collection approvals

## Error Messages

When the invariant is violated, you'll receive error messages like:

```
noCustomOwnershipTimes invariant is enabled: ownership times must be full range [{ start: 1, end: 18446744073709551615 }]
```

## Related Concepts

-   [Badge Collections](./badge-collections.md)
-   [Transferability Approvals](./transferability-approvals.md)
-   [Time Fields](./time-fields.md)
-   [UintRange](./uintrange.md)
