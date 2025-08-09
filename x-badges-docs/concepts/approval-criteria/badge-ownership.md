# Token Ownership

Require specific token holdings from the initiator as a prerequisite for transfer approval. This approval criteria checks on-chain balances to ensure users own required tokens before allowing transfers.

## Overview

Token ownership requirements enable gating mechanisms where users must possess specific tokens to access certain transfers. This creates dependency relationships between collections and enables sophisticated access control systems.

**Key Benefits**:

-   **Access Control**: Gate transfers based on token ownership
-   **Collection Dependencies**: Create relationships between different collections
-   **On-Chain Verification**: Automatic balance checking without external data
-   **Flexible Requirements**: Support for amount ranges, time-based ownership, and multiple token types

## Interface

```typescript
interface MustOwnTokens<T extends NumberType> {
    collectionId: T;
    amountRange: UintRange<T>; // Min/max amount expected
    ownershipTimes: UintRange<T>[];
    tokenIds: UintRange<T>[];

    overrideWithCurrentTime: boolean; // Use current block time. Overrides ownershipTimes with [{ start: currentTime, end: currentTime }]
    mustSatisfyForAllAssets: boolean; // All vs one requirement
    ownershipCheckParty: string; // Which party to check ownership for: "initiator", "sender", "recipient" (default: "initiator" if empty)
}
```

## Field Descriptions

### collectionId

-   **Type**: `T` (NumberType)
-   **Description**: The ID of the collection containing the required tokens
-   **Example**: `"1"` for collection ID 1

### amountRange

-   **Type**: `UintRange<T>`
-   **Description**: Minimum and maximum amount of tokens the user must own
-   **Format**: `{ start: "minAmount", end: "maxAmount" }`
-   **Example**: `{ start: "1", end: "10" }` requires 1-10 tokens (amounts)

### ownershipTimes

-   **Type**: `UintRange<T>[]`
-   **Description**: Time ranges when the user must have owned the tokens (UNIX milliseconds)
-   **Example**: `[{ start: "1691931600000", end: "1723554000000" }]` for Aug 13, 2023 - Aug 13, 2024

### tokenIds

-   **Type**: `UintRange<T>[]`
-   **Description**: Specific token IDs that must be owned
-   **Example**: `[{ start: "1", end: "100" }]` for token IDs 1-100

### overrideWithCurrentTime

-   **Type**: `boolean`
-   **Description**: When true, ignores `ownershipTimes` and uses current block time
-   **Behavior**: Sets ownership time to `[{ start: currentTime, end: currentTime }]`
-   **Use Case**: Require current ownership only, not historical ownership

### mustSatisfyForAllAssets

-   **Type**: `boolean`
-   **Description**: Controls whether all specified token requirements must be met or just one
-   **True**: User must own ALL specified token combinations
-   **False**: User must own AT LEAST ONE of the specified token combinations

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

Require users to own specific tokens to access premium features or exclusive transfers.

```json
{
    "mustOwnTokens": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "tokenIds": [{ "start": "1", "end": "1" }],
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
    "mustOwnTokens": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "tokenIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "initiator"
        }
    ]
}
```

#### Check Sender Ownership

```json
{
    "mustOwnTokens": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "tokenIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "sender"
        }
    ]
}
```

#### Check Recipient Ownership

```json
{
    "mustOwnTokens": [
        {
            "collectionId": "1",
            "amountRange": { "start": "1", "end": "1" },
            "tokenIds": [{ "start": "1", "end": "1" }],
            "ownershipCheckParty": "recipient"
        }
    ]
}
```
