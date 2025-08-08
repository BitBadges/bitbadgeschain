# Badge Ownership

Require specific badge holdings from the specified party as a prerequisite for transfer approval. This approval criteria checks on-chain badge balances to ensure the specified party owns required badges before allowing transfers.

## Overview

Badge ownership requirements enable gating mechanisms where users must possess specific badges to access certain transfers. This creates dependency relationships between collections and enables sophisticated access control systems.

**Key Benefits**:

-   **Access Control**: Gate transfers based on badge ownership
-   **Collection Dependencies**: Create relationships between different badge collections
-   **On-Chain Verification**: Automatic balance checking without external data
-   **Flexible Requirements**: Support for amount ranges, time-based ownership, and multiple badge types
-   **Party-Specific Checks**: Check ownership for initiator, sender, or recipient

## Interface

```typescript
interface MustOwnBadges<T extends NumberType> {
    collectionId: T;
    amountRange: UintRange<T>; // Min/max amount expected
    ownershipTimes: UintRange<T>[];
    badgeIds: UintRange<T>[];

    overrideWithCurrentTime: boolean; // Use current block time. Overrides ownershipTimes with [{ start: currentTime, end: currentTime }]
    mustSatisfyForAllAssets: boolean; // All vs one badge requirement
    ownershipCheckParty: string; // Which party to check ownership for: "initiator", "sender", "recipient" (default: "initiator" if empty)
}
```

## Field Descriptions

### collectionId

-   **Type**: `T` (NumberType)
-   **Description**: The ID of the collection containing the required badges
-   **Example**: `"1"` for collection ID 1

### amountRange

-   **Type**: `UintRange<T>`
-   **Description**: Minimum and maximum amount of badges the user must own
-   **Format**: `{ start: "minAmount", end: "maxAmount" }`
-   **Example**: `{ start: "1", end: "10" }` requires 1-10 badges (amounts)

### ownershipTimes

-   **Type**: `UintRange<T>[]`
-   **Description**: Time ranges when the user must have owned the badges (UNIX milliseconds)
-   **Example**: `[{ start: "1691931600000", end: "1723554000000" }]` for Aug 13, 2023 - Aug 13, 2024

### badgeIds

-   **Type**: `UintRange<T>[]`
-   **Description**: Specific badge IDs that must be owned
-   **Example**: `[{ start: "1", end: "100" }]` for badge IDs 1-100

### overrideWithCurrentTime

-   **Type**: `boolean`
-   **Description**: When true, ignores `ownershipTimes` and uses current block time
-   **Behavior**: Sets ownership time to `[{ start: currentTime, end: currentTime }]`
-   **Use Case**: Require current ownership only, not historical ownership

### mustSatisfyForAllAssets

-   **Type**: `boolean`
-   **Description**: Controls whether all specified badge requirements must be met or just one
-   **True**: User must own ALL specified badge combinations
-   **False**: User must own AT LEAST ONE of the specified badge combinations

### ownershipCheckParty

-   **Type**: `string`
-   **Description**: Specifies which party of the transfer to check ownership for
-   **Options**:
    -   `"initiator"` (default): Check ownership for the address that initiated the transfer
    -   `"sender"`: Check ownership for the address sending the badges
    -   `"recipient"`: Check ownership for the address receiving the badges
-   **Default**: `"initiator"` (if empty or not specified)
-   **Example**: `"sender"` to require the sender to own specific badges before allowing the transfer

## Example

Require users to own specific badges to access premium features or exclusive transfers.

```json
{
    "mustOwnBadges": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "badgeIds": [{ "start": "1", "end": "1" }],
            "overrideWithCurrentTime": false,
            "mustSatisfyForAllAssets": true,
            "ownershipCheckParty": "initiator"
        }
    ]
}
```

### Party-Specific Examples

#### Check Initiator Ownership (Default)

```json
{
    "mustOwnBadges": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "badgeIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "initiator"
        }
    ]
}
```

#### Check Sender Ownership

```json
{
    "mustOwnBadges": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "badgeIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "sender"
        }
    ]
}
```

#### Check Recipient Ownership

```json
{
    "mustOwnBadges": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "badgeIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "recipient"
        }
    ]
}
```
