# Predetermined Balances

## Overview

Predetermined balances provide fine-grained control over the exact amounts and order of transfers in an approval. Unlike traditional tally-based systems where you approve a total amount (e.g., 100 badges) without controlling the specific combinations, predetermined balances let you explicitly define:

-   **Exact amounts** that must be transferred
-   **Specific order** of transfers
-   **Precise badge IDs and ownership times** for each transfer

**Key Principle**: The transfer will fail if the balances are not EXACTLY as defined in the predetermined balances.

## Interface Definition

```typescript
export interface PredeterminedBalances<T extends NumberType> {
    manualBalances: ManualBalances<T>[];
    incrementedBalances: IncrementedBalances<T>;
    orderCalculationMethod: PredeterminedOrderCalculationMethod;
}
```

## Balance Definition Methods

There are two mutually exclusive ways to define balances:

### 1. Manual Balances

Define an array of specific balance sets manually. Each element corresponds to a different transfer.

```json
{
    "manualBalances": [
        {
            "amount": "1",
            "badgeIds": [
                {
                    "start": "1",
                    "end": "1"
                }
            ],
            "ownershipTimes": [
                {
                    "start": "1691978400000",
                    "end": "1723514400000"
                }
            ]
        },
        {
            "amount": "5",
            "badgeIds": [
                {
                    "start": "2",
                    "end": "6"
                }
            ],
            "ownershipTimes": [
                {
                    "start": "1691978400000",
                    "end": "1723514400000"
                }
            ]
        }
    ]
}
```

**Use Case**: When you need complete control over each specific transfer amount and timing.

### 2. Incremented Balances

Define starting balances and rules for subsequent transfers. Perfect for sequential minting or time-based releases or other common patterns. Note that most options are incompatible with each other.

```json
{
    "incrementedBalances": {
        "startBalances": [
            {
                "amount": "1",
                "badgeIds": [
                    {
                        "start": "1",
                        "end": "1"
                    }
                ],
                "ownershipTimes": [
                    {
                        "start": "1691978400000",
                        "end": "1723514400000"
                    }
                ]
            }
        ],
        "incrementBadgeIdsBy": "1",
        "incrementOwnershipTimesBy": "0",
        "durationFromTimestamp": "0",
        "allowOverrideTimestamp": false,
        "allowOverrideWithAnyValidBadge": false,
        "recurringOwnershipTimes": {
            "startTime": "0",
            "intervalLength": "0",
            "chargePeriodLength": "0"
        }
    }
}
```

#### Increment Options

| Field                            | Description                                          | Example                                              |
| -------------------------------- | ---------------------------------------------------- | ---------------------------------------------------- |
| `incrementBadgeIdsBy`            | Amount to increment badge IDs by after each transfer | `"1"` = next transfer gets badge ID 2, then 3, etc.  |
| `incrementOwnershipTimesBy`      | Amount to increment ownership times by               | `"86400000"` = add 1 day to ownership times          |
| `durationFromTimestamp`          | Calculate ownership times from timestamp + duration  | `"2592000000"` = 30 days from transfer time          |
| `allowOverrideTimestamp`         | Allow custom timestamp override in transfer          | `true` = users can specify custom start time         |
| `allowOverrideWithAnyValidBadge` | Allow any valid badge ID (one) override              | `true` = users can specify any single valid badge ID |
| `recurringOwnershipTimes`        | Define recurring time intervals                      | Monthly subscriptions, weekly rewards                |

#### Duration From Timestamp

Dynamically calculate ownership times from a timestamp plus a set duration. This overwrites all ownership times in the starting balances.

```json
{
    "durationFromTimestamp": "2592000000", // 30 days in milliseconds
    "allowOverrideTimestamp": true
}
```

**Behavior**:

-   **Default**: Uses transfer time as the base timestamp
-   **Override**: If `allowOverrideTimestamp` is true, users can specify a custom timestamp in `MsgTransferBadges` `precalculationOptions`
-   **Calculation**: `ownershipTime = baseTimestamp + durationFromTimestamp`
-   **Overwrite**: All ownership times in starting balances are replaced with [{ "start": baseTimestamp, "end": baseTimestamp + durationFromTimestamp }]

#### Recurring Ownership Times

Define repeating time intervals for subscriptions or periodic rewards:

```json
{
    "recurringOwnershipTimes": {
        "startTime": "1691978400000", // When intervals begin
        "intervalLength": "2592000000", // 30 days in milliseconds
        "chargePeriodLength": "604800000" // 7 days advance charging
    }
}
```

**Example**: Monthly subscription starting August 13, 2023, with 7-day advance charging period.

## Precalculating Balances

### The Race Condition Problem

Predetermined balances can change rapidly between transaction broadcast and confirmation. For example:

-   Other users' mints get processed
-   Badge IDs shift due to concurrent activity
-   Manual balance specification becomes unreliable

### The Solution: Precalculation

Use `precalculateBalancesFromApproval` in [MsgTransferBadges](../../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md) to dynamically calculate balances at execution time.

```typescript
{
  precalculateBalancesFromApproval: {
    approvalId: string;           // The approval to precalculate from
    approvalLevel: string;        // "collection" | "incoming" | "outgoing"
    approverAddress: string;      // "" if collection-level
    version: string;              // Must specify exact version
  },
  precalculationOptions: {
    // Additional override options dependent on the selections
  }
}
```

## Order Calculation Methods

The system needs to determine which balance set to use for each transfer. This is controlled by the `orderCalculationMethod`.

### How Order Numbers Work

The order number determines which balances to transfer, but it works differently depending on the balance type:

#### Manual Balances

-   **Order number = 0**: Transfer `manualBalances[0]` (first element)
-   **Order number = 1**: Transfer `manualBalances[1]` (second element)
-   **Order number = 5**: Transfer `manualBalances[5]` (sixth element)

**Example**: If you have 3 manual balance sets, order numbers 0, 1, and 2 will use each set once. Order number 3 would be out of bounds.

#### Incremented Balances

-   **Order number = 0**: Use starting balances as-is (no increments)
-   **Order number = 1**: Apply increments once to starting balances
-   **Order number = 5**: Apply increments five times to starting balances

**Example**: Starting with badge ID 1, increment by 1:

-   Order 0: Badge ID 1
-   Order 1: Badge ID 2
-   Order 2: Badge ID 3
-   Order 5: Badge ID 6

### Transfer-Based Order Numbers

Track the number of transfers to determine order:

| Method                                 | Description           | Use Case                    |
| -------------------------------------- | --------------------- | --------------------------- |
| `useOverallNumTransfers`               | Global transfer count | Simple sequential transfers |
| `usePerToAddressNumTransfers`          | Per-recipient count   | User-specific limits        |
| `usePerFromAddressNumTransfers`        | Per-sender count      | Sender-specific limits      |
| `usePerInitiatedByAddressNumTransfers` | Per-initiator count   | Initiator-specific limits   |

**Important**: Uses the same tracker as [Max Number of Transfers](max-number-of-transfers.md). Trackers are:

-   Increment-only and immutable
-   Shared between predetermined balances and max transfer limits
-   Must be carefully managed to avoid conflicts

### Merkle-Based Order Numbers

Use Merkle challenge leaf indices (leftmost = 0, rightmost = numLeaves - 1) for reserved transfers:

```typescript
{
  "useMerkleChallengeLeafIndex": true,
  "challengeTrackerId": "uniqueId"
}
```

**Use Case**: Reserve specific badge IDs for specific users or claim codes.

## Order Calculation Interface

```typescript
export interface PredeterminedOrderCalculationMethod {
    useOverallNumTransfers: boolean;
    usePerToAddressNumTransfers: boolean;
    usePerFromAddressNumTransfers: boolean;
    usePerInitiatedByAddressNumTransfers: boolean;
    useMerkleChallengeLeafIndex: boolean;
    challengeTrackerId: string;
}
```

## Boundary Handling

### Understanding Bounds

Every approval defines bounds through its core fields (badgeIds, ownershipTimes, etc.). For example:

-   **Badge IDs**: 1-100
-   **Ownership Times**: Mon-Fri only
-   **Transfer Times**: Specific date range

Predetermined balances must work within these bounds, but note that order numbers can eventually exceed them.

### Boundary Scenarios

#### Complete Out-of-Bounds

**Scenario**: Order number corresponds to balances completely outside approval bounds.

**Example**:

-   Approval allows badge IDs 1-100
-   Increment by 1 for each transfer
-   Order number 101 would require badge ID 101 (out of bounds)

**Result**: Transfer is ignored because badge ID 101 never matches the approval's badge ID range.

#### Partial Overlap

**Scenario**: Order number corresponds to balances that partially overlap with approval bounds.

**Example**:

-   Approval allows badge IDs 1-100
-   Transfer requires badge IDs 95-105
-   Badge IDs 95-100 are in bounds, 101-105 are out of bounds

**Result**:

-   Only in-bounds balances (95-100) are approved by current approval
-   Out-of-bounds balances (101-105) must be approved by a separate approval
-   The complete transfer (95-105) must still be exactly as defined

**Important**: The transfer will fail unless all out-of-bounds balances are approved by other approvals.
