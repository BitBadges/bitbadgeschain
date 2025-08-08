# GetCollection

Retrieves complete information about a collection.

## Proto Definition

```protobuf
message QueryGetCollectionRequest {
  string collectionId = 1; // ID of collection to retrieve
}

message QueryGetCollectionResponse {
  BadgeCollection collection = 1;
}

message BadgeCollection {
  string collectionId = 1; // Unique identifier for this collection
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 2; // Collection metadata over time
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 3; // Token metadata over time
  string balancesType = 4; // Type of balances ("Standard", "Off-Chain - Indexed", etc.)
  repeated OffChainBalancesMetadataTimeline offChainBalancesMetadataTimeline = 5; // Off-chain balance metadata
  repeated CustomDataTimeline customDataTimeline = 7; // Arbitrary custom data over time
  repeated ManagerTimeline managerTimeline = 8; // Manager address over time
  CollectionPermissions collectionPermissions = 9; // Collection permissions
  repeated CollectionApproval collectionApprovals = 10; // Collection-level approvals
  repeated StandardsTimeline standardsTimeline = 11; // Standards over time
  repeated IsArchivedTimeline isArchivedTimeline = 12; // Archive status over time
  UserBalanceStore defaultBalances = 13; // Default balance store for users
  string createdBy = 14; // Creator of the collection
  repeated UintRange validBadgeIds = 15; // Valid token ID ranges
  string mintEscrowAddress = 16; // Generated escrow address for the collection
}

// See all the proto definitions [here](https://github.com/bitbadges/bitbadgeschain/tree/master/proto/badges)
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-collection [id]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_collection/1"
```

### Response Example

```json
{
    "collection": {
        "collectionId": "1"
        // ...
    }
}
```
