# Subscriptions Protocol

The Subscriptions Protocol enables collections to implement subscription-based token ownership with recurring payments, time-limited access, and a tipping system for automatic renewal. This protocol standardizes how subscription tokens are created, distributed, and renewed.

## Protocol Overview

Subscription collections allow users to pay a recurring fee + tip to maintain ownership of tokens for specific time periods. The protocol ensures predictable behavior for subscription management across different applications.

## Protocol Requirements

### Standards Declaration

Collections must include "Subscriptions" in their standards timeline:

```json
{
    "standardsTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "standards": ["Subscriptions"]
        }
    ]
}
```

### Token ID Configuration

-   **Single Token ID**: Only one token ID range (1-1) is permitted
-   **Valid Token IDs**: Must match exactly with subscription approval token IDs

### Collection Approvals

Must contain at least one subscription faucet approval with the following characteristics:

#### From Address

-   **fromListId**: Must be "Mint" (minting operation)

#### Approval Criteria Requirements

##### Coin Transfers

-   **Required**: At least one coin transfer specification
-   **Single Denom**: All coin transfers must use the same denomination
-   **No Address Override**: `overrideFromWithApproverAddress` and `overrideToWithInitiator` must be false

##### Predetermined Balances

Must use `incrementedBalances` with specific configuration:

```json
{
    "predeterminedBalances": {
        "incrementedBalances": {
            "startBalances": [
                {
                    "amount": "1",
                    "badgeIds": [{ "start": "1", "end": "1" }],
                    "ownershipTimes": [{ "start": "0", "end": "0" }]
                }
            ],
            "incrementBadgeIdsBy": "0",
            "incrementOwnershipTimesBy": "0",
            "durationFromTimestamp": "2592000000", // 30 days in milliseconds
            "allowOverrideTimestamp": true,
            "recurringOwnershipTimes": {
                "startTime": "0",
                "intervalLength": "0",
                "chargePeriodLength": "0"
            }
        }
    }
}
```

**Key Requirements:**

-   **Amount**: Must be exactly 1
-   **Token IDs**: Single token ID (1-1)
-   **Duration**: Must be greater than 0 (subscription period length)
-   **Override Timestamp**: Must be true for faucet functionality
-   **No Increments**: Token ID and ownership time increments must be 0
-   **No Recurring**: Recurring ownership times must be all 0

##### Restrictions

-   **No Merkle Challenges**: Cannot include merkle challenges
-   **No Token Requirements**: Cannot include mustOwnBadges requirements
-   **No Address Restrictions**: Cannot require from/to equals initiated by
-   **No Override Approvals**: Cannot override incoming approvals

## User Subscription Management

Users manage their subscriptions through incoming approvals that complement the collection's faucet approval:

### User Incoming Approval Requirements

#### Basic Configuration

-   **fromListId**: Must be "Mint"
-   **Token IDs**: Must match subscription approval token IDs exactly
-   **Single Token**: Only one token ID range permitted

#### Coin Transfer Configuration

```json
{
    "coinTransfers": [
        {
            "coins": [
                {
                    "denom": "ubadge", // Must match subscription denom
                    "amount": "100000" // Must be >= subscription amount
                }
            ],
            "overrideFromWithApproverAddress": true,
            "overrideToWithInitiator": true
        }
    ]
}
```

#### Predetermined Balances for Renewals

```json
{
    "predeterminedBalances": {
        "incrementedBalances": {
            "startBalances": [
                {
                    "amount": "1",
                    "badgeIds": [{ "start": "1", "end": "1" }]
                }
            ],
            "incrementBadgeIdsBy": "0",
            "incrementOwnershipTimesBy": "0",
            "durationFromTimestamp": "0",
            "allowOverrideTimestamp": false,
            "recurringOwnershipTimes": {
                "startTime": "1672531200000", // Current subscription start
                "intervalLength": "2592000000", // 30 days
                "chargePeriodLength": "604800000" // 7 days max charge period
            }
        }
    }
}
```

#### Transfer Limits

```json
{
    "maxNumTransfers": {
        "overallMaxNumTransfers": "1",
        "resetTimeIntervals": {
            "startTime": "1672531200000",
            "intervalLength": "2592000000" // Same as subscription interval
        }
    }
}
```

## Protocol Validation Logic

### Collection Validation

**API Documentation:** [doesCollectionFollowSubscriptionProtocol](https://bitbadges.github.io/bitbadgesjs/functions/doesCollectionFollowSubscriptionProtocol.html)

```typescript
function doesCollectionFollowSubscriptionProtocol(collection) {
    // Check for "Subscriptions" standard
    const hasSubscriptionStandard = collection.standardsTimeline.some(
        (standard) =>
            standard.standards.includes('Subscriptions') &&
            isCurrentTime(standard.timelineTimes)
    );

    if (!hasSubscriptionStandard) return false;

    // Find subscription faucet approvals
    const subscriptionApprovals = collection.collectionApprovals.filter(
        (approval) => isSubscriptionFaucetApproval(approval)
    );

    if (subscriptionApprovals.length < 1) return false;

    // Validate single token ID requirement
    if (collection.validBadgeIds.length !== 1) return false;

    // Ensure approval token IDs match collection token IDs
    const allApprovalBadgeIds = subscriptionApprovals
        .map((approval) => approval.badgeIds)
        .flat();

    return badgeIdsMatch(collection.validBadgeIds, allApprovalBadgeIds);
}
```

### Faucet Approval Validation

**API Documentation:** [isSubscriptionFaucetApproval](https://bitbadges.github.io/bitbadgesjs/functions/isSubscriptionFaucetApproval.html)

```typescript
function isSubscriptionFaucetApproval(approval) {
    // Must be from Mint
    if (approval.fromListId !== 'Mint') return false;

    // Must have coin transfers
    if (!approval.approvalCriteria?.coinTransfers?.length) return false;

    // Single denomination requirement
    const allDenoms = approval.approvalCriteria.coinTransfers.flatMap((ct) =>
        ct.coins.map((c) => c.denom)
    );
    if (new Set(allDenoms).size > 1) return false;

    // No address overrides in coin transfers
    for (const coinTransfer of approval.approvalCriteria.coinTransfers) {
        if (
            coinTransfer.overrideFromWithApproverAddress ||
            coinTransfer.overrideToWithInitiator
        ) {
            return false;
        }
    }

    // Validate incremented balances configuration
    const incrementedBalances =
        approval.approvalCriteria.predeterminedBalances?.incrementedBalances;
    if (!incrementedBalances) return false;

    return validateIncrementedBalances(incrementedBalances, approval.badgeIds);
}
```

### User Approval Validation

**API Documentation:** [isUserRecurringApproval](https://bitbadges.github.io/bitbadgesjs/functions/isUserRecurringApproval.html)

```typescript
function isUserRecurringApproval(userApproval, subscriptionApproval) {
    // Must be from Mint
    if (userApproval.fromListId !== 'Mint') return false;

    // Token IDs must match subscription
    if (!badgeIdsMatch(userApproval.badgeIds, subscriptionApproval.badgeIds)) {
        return false;
    }

    // Payment amount must be >= subscription amount
    const userAmount =
        userApproval.approvalCriteria?.coinTransfers?.[0]?.coins?.[0]?.amount;
    const subscriptionAmount =
        subscriptionApproval.approvalCriteria?.coinTransfers?.[0]?.coins?.[0]
            ?.amount;
    if (userAmount < subscriptionAmount) return false;

    // Validate coin transfer overrides
    const coinTransfer = userApproval.approvalCriteria.coinTransfers[0];
    if (
        !coinTransfer.overrideFromWithApproverAddress ||
        !coinTransfer.overrideToWithInitiator
    ) {
        return false;
    }

    // Validate recurring configuration
    return validateRecurringConfiguration(userApproval, subscriptionApproval);
}
```

## Usage Examples

### Basic Subscription Collection

```json
{
    "validBadgeIds": [{ "start": "1", "end": "1" }],
    "standardsTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "standards": ["Subscriptions"]
        }
    ],
    "collectionApprovals": [
        {
            "fromListId": "Mint",
            "toListId": "All",
            "initiatedByListId": "All",
            "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "badgeIds": [{ "start": "1", "end": "1" }],
            "approvalCriteria": {
                "coinTransfers": [
                    {
                        "coins": [{ "denom": "ubadge", "amount": "100000" }],
                        "overrideFromWithApproverAddress": false,
                        "overrideToWithInitiator": false
                    }
                ],
                "predeterminedBalances": {
                    "incrementedBalances": {
                        "startBalances": [
                            {
                                "amount": "1",
                                "badgeIds": [{ "start": "1", "end": "1" }]
                            }
                        ],
                        "durationFromTimestamp": "2592000000",
                        "allowOverrideTimestamp": true
                    }
                }
            }
        }
    ]
}
```

### User Subscription Setup

```json
{
    "fromListId": "Mint",
    "badgeIds": [{ "start": "1", "end": "1" }],
    "approvalCriteria": {
        "coinTransfers": [
            {
                "coins": [{ "denom": "ubadge", "amount": "100000" }],
                "overrideFromWithApproverAddress": true,
                "overrideToWithInitiator": true
            }
        ],
        "predeterminedBalances": {
            "incrementedBalances": {
                "recurringOwnershipTimes": {
                    "intervalLength": "2592000000",
                    "chargePeriodLength": "604800000"
                }
            }
        },
        "maxNumTransfers": {
            "overallMaxNumTransfers": "1",
            "resetTimeIntervals": {
                "intervalLength": "2592000000"
            }
        }
    }
}
```

## Implementation Benefits

1. **Standardization**: Predictable subscription behavior across applications
2. **Interoperability**: Common interface for subscription management
3. **Automation**: Recurring payment and renewal mechanisms
4. **Flexibility**: Configurable subscription periods and pricing
5. **Validation**: Built-in compliance checking for protocol adherence

The Subscriptions Protocol provides a robust foundation for implementing subscription-based token systems while maintaining the flexibility and security of the BitBadges approval system.
