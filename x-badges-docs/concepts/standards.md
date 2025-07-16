# Standards

Standards are informational tags that provide guidance on how to interpret and implement collection features. The collection interface is very feature-rich, and oftentimes you may need certain features to be implemented in a certain way, avoid certain features, etc. That is what standards are for.

## Timeline Implementation

```json
"standardsTimeline": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "standards": ["transferable", "text-only-metadata", "non-fungible", "attendance-format"]
  }
]
```

## Important Notes

-   **No blockchain validation** - Standards are purely informational
-   **Multiple standards allowed** - As long as they are compatible
-   **Application responsibility** - Queriers must verify compliance

## BitBadges Site Standards

The BitBadges site recognizes specific standards that collections can implement to ensure compatibility with various features and integrations:

### 1. Tradable Standard

Collections marked with the **Tradable** standard are marked as tradable on the BitBadges site. We will track orderbook, volume, price, and other metrics for these collections. Also, the interface will be optimized for trading.

-   **Requirements**:
    -   Ensure the collection lends itself to user-to-user trading
    -   Must be able to support bids / offers / listings / collection offers. These are standardized approvals that follow specific rules.

### 2. NFT Standard

Collections marked with the **NFT** standard are expected to be non-fungible tokens with supply = 1 for every badge ID.

-   **Requirements**:
    -   Each badge ID must have supply = 1 and full ownership times
    -   No fungible badge IDs allowed
    -   Maintains uniqueness across all badge IDs in the collection

### 3. Cosmos Wrappable Standard

Collections marked with the **Cosmos Wrappable** standard can be wrapped into Cosmos SDK coin denominations.

-   **Requirements**:
    -   Must have at least one wrapper path defined
    -   Should support bidirectional wrapping/unwrapping
    -   Refer to the [Cosmos Wrapper documentation](../cosmos-wrapper-paths.md) for detailed implementation guidelines

### 4. Subscriptions Standard

Collections marked with the **Subscriptions** standard are designed for recurring content delivery and subscription-based systems.

-   **Requirements**:
    -   Must support time-based ownership periods for subscription-like behavior
    -   Must be able to handle recurring badge issuance and expiration
    -   Should support dynamic content updates based on subscription status
-   **Implementation**: See [Subscriptions Protocol](protocols/subscriptions-protocol.md) for detailed implementation requirements and validation logic

### 5. Quests Standard

Collections marked with the **Quests** standard are designed for achievement-based systems and quest completion tracking.

-   **Requirements**:
    -   Should implement quest completion tracking and reward distribution
    -   Must support achievement-based badge issuance
    -   Should handle quest progression and milestone tracking

### Using BitBadges Standards

To implement these standards on your collection, add them to your `standardsTimeline`:

```json
{
    "standardsTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "standards": [
                "Tradable",
                "NFT",
                "Cosmos Wrappable",
                "Subscriptions",
                "Quests"
            ]
        }
    ]
}
```

**Note**: These standards are informational and do not enforce blockchain-level validation. Applications and platforms are responsible for verifying compliance with the specified standards.
