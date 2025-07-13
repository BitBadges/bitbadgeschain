# Dynamic Store Challenges

Dynamic Store Challenges are approval criteria that require the transfer initiator to pass boolean checks against one or more dynamic stores. The concept is simple: if a user can return `true` for all dynamic store fetches by their address (the initiator), they pass all challenges.

This is a powerful feature that can be used to implement complex approval logic. This is intended to be used in conjunctions with smart contracts allowing for more complex logic to be implemented.

Note: For fully off-chain alternatives, you may want to consider the merkleChallenges field to save on gas costs and avoid the need for transactions.

## Proto Definition

```protobuf
// DynamicStoreChallenge defines a challenge that requires the initiator to pass a dynamic store check.
message DynamicStoreChallenge {
  // The ID of the dynamic store to check.
  string storeId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## How It Works

### Challenge Logic

1. **Initiator Check**: The system checks the transfer initiator's address against the specified dynamic store
2. **Boolean Evaluation**: The dynamic store returns a boolean value for the initiator's address
3. **Pass Requirement**: The initiator must have a `true` value in the dynamic store to pass the challenge
4. **Multiple Challenges**: If multiple challenges exist, the initiator must pass ALL of them

### Challenge Evaluation

```javascript
function passesAllChallenges(initiatorAddress, challenges) {
    for (const challenge of challenges) {
        const storeValue = getDynamicStoreValue(
            challenge.storeId,
            initiatorAddress
        );
        if (!storeValue) {
            return false; // Failed challenge
        }
    }
    return true; // Passed all challenges
}
```

## Usage in Approval Criteria

Dynamic store challenges can be used in all three approval levels:

### Collection Approval Criteria

```json
{
    "dynamicStoreChallenges": [
        {
            "storeId": "1"
        },
        {
            "storeId": "5"
        }
    ]
}
```

### Outgoing Approval Criteria

```json
{
    "dynamicStoreChallenges": [
        {
            "storeId": "3"
        }
    ]
}
```

### Incoming Approval Criteria

```json
{
    "dynamicStoreChallenges": [
        {
            "storeId": "2"
        }
    ]
}
```

## Use Cases

### Membership Verification

-   **VIP Access**: Store `true` for VIP members in a dynamic store
-   **Active Users**: Track active users with periodic updates
-   **Subscription Status**: Verify current subscription status

### Governance and Voting

-   **Voting Rights**: Store voting eligibility in dynamic stores
-   **Proposal Participation**: Track participation in governance proposals
-   **Staking Requirements**: Verify staking status for transfer privileges

### Game Mechanics

-   **Achievement Unlocks**: Require specific achievements for certain transfers
-   **Level Requirements**: Check player levels or progression
-   **Quest Completion**: Verify quest completion status

### Compliance and KYC

-   **Identity Verification**: Store KYC completion status
-   **Compliance Checks**: Verify regulatory compliance
-   **Whitelist Management**: Maintain dynamic whitelists

## Example Implementation

### Setup Dynamic Store Challenge

```json
{
    "collectionApprovals": [
        {
            "approvalId": "member-only-transfers",
            "fromListId": "All",
            "toListId": "All",
            "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "badgeIds": [{ "start": "1", "end": "100" }],
            "dynamicStoreChallenges": [
                {
                    "storeId": "1" // Member status store
                },
                {
                    "storeId": "2" // Active subscription store
                }
            ]
        }
    ]
}
```

### Managing Store Values

Users and managers can update dynamic store values using the appropriate messages and queries.

## Related Operations

For more information on managing dynamic stores, see:

-   **[Create Dynamic Store](../../messages/msg-create-dynamic-store.md)** - Creating new dynamic stores
-   **[Update Dynamic Store](../../messages/msg-update-dynamic-store.md)** - Updating store properties
-   **[Delete Dynamic Store](../../messages/msg-delete-dynamic-store.md)** - Removing dynamic stores
-   **[Set Dynamic Store Value](../../messages/msg-set-dynamic-store-value.md)** - Setting address-specific values
-   **[Get Dynamic Store Value](../../queries/query-get-dynamic-store-value.md)** - Querying store values

## Important Notes

### Default Values

-   Dynamic stores have default values for uninitialized addresses
-   If an address has no explicit value, the default is used
-   Default values are set when creating the dynamic store

### Performance Considerations

-   Dynamic store lookups are efficient on-chain operations
-   Multiple challenges are evaluated sequentially
-   Consider the gas cost of multiple store lookups

### Security Implications

-   Store creators control who can update values
-   Values can be changed, affecting future challenge results
-   Design challenges with appropriate access controls

Dynamic Store Challenges provide a flexible and powerful way to implement complex approval logic while maintaining simplicity in the core protocol design.
