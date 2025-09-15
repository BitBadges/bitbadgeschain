# GetWrappableBalances

Retrieves the maximum amount of badges that can be wrapped into cosmos coins for a specific denom and user address.

## Proto Definition

```protobuf
message QueryGetWrappableBalancesRequest {
  string denom = 1; // The cosmos coin denom (e.g., "badges:1:mytoken")
  string address = 2; // Address to get wrappable balances for
}

message QueryGetWrappableBalancesResponse {
  Uint maxWrappableAmount = 1; // Maximum amount that can be wrapped
}
```

## Description

This query calculates the maximum amount of badges that a user can wrap into cosmos coins for a given denom. It:

1. **Parses the denom**: Extracts the collection ID from the denom format `badges:COLL_ID:*` or `badgeslp:COLL_ID:*`
2. **Finds the wrapper path**: Locates the corresponding cosmos coin wrapper path for the denom
3. **Calculates maximum wrappable amount**: Determines the largest amount the user can wrap based on their current badge balances

The query supports both static denoms and dynamic `{id}` placeholder denoms. For dynamic denoms, it extracts numeric characters from the base denom to replace the `{id}` placeholder.

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-wrappable-balances [denom] [address]

# Example with static denom
bitbadgeschaind query badges get-wrappable-balances "badges:1:mytoken" "bb1..."

# Example with dynamic denom (where 123 is the badge ID)
bitbadgeschaind query badges get-wrappable-balances "badgeslp:1:token123" "bb1..."
```

### REST API

```bash
# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_wrappable_balances?denom=badges:1:mytoken&address=bb1..."
```

### Response Example

```json
{
    "maxWrappableAmount": "1000"
}
```

## Error Cases

-   **Invalid denom format**: Denom must start with "badges:" or "badgeslp:" and follow the format `badges:COLL_ID:*` or `badgeslp:COLL_ID:*`
-   **Collection not found**: The specified collection ID doesn't exist
-   **Wrapper path not found**: No cosmos coin wrapper path matches the given denom
-   **No balances**: User has no balances for the required badge IDs and ownership times

## Use Cases

-   **Pre-wrapping validation**: Check how much a user can wrap before attempting the wrap operation
-   **UI display**: Show users their maximum wrappable amount for different tokens
-   **Batch operations**: Calculate optimal amounts for multiple wrapping operations
-   **Integration testing**: Verify wrapper path configurations work correctly

## Related Queries

-   [GetBalance](./get-balance.md) - Get user's current badge balances
-   [GetCollection](./get-collection.md) - Get collection details including wrapper paths
