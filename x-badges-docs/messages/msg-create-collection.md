# MsgCreateCollection

Creates a new badge collection.

The collectionId will be assigned at execution time and is obtainable in the transaction response. Subsequent updates to the collection will be through MsgUpdateCollection.

## Creation Only Properties

The creation or genesis transaction for a collection is unique in a couple ways.

There are no permissions previously set, so there are no restrictions for what can be set vs not. Subsequent updates to the collection must follow any previously set permissions.

This is the only time that you can specify `balancesType` and the `defaultBalances` information.

## Proto Definition

```protobuf
message MsgCreateCollection {
  string creator = 1; // Address creating the collection
  string balancesType = 2; // "Standard", "Off-Chain - Indexed", etc.
  UserBalanceStore defaultBalances = 4;
  repeated UintRange validBadgeIds = 5; // Badge ID ranges to include
  CollectionPermissions collectionPermissions = 6;
  repeated ManagerTimeline managerTimeline = 7;
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 8;
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 9;
  repeated OffChainBalancesMetadataTimeline offChainBalancesMetadataTimeline = 10;
  repeated CustomDataTimeline customDataTimeline = 11;
  repeated CollectionApproval collectionApprovals = 12;
  repeated StandardsTimeline standardsTimeline = 13;
  repeated IsArchivedTimeline isArchivedTimeline = 14;
  repeated cosmos.base.v1beta1.Coin mintEscrowCoinsToTransfer = 16;
  repeated CosmosCoinWrapperPathAddObject cosmosCoinWrapperPathsToAdd = 17;
}

message MsgCreateCollectionResponse {
  string collectionId = 1; // ID of the created collection
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges create-collection '[tx-json]' --from creator-key
```

### JSON Example

For complete transaction examples, see [MsgCreateCollection Examples](../examples/txs/msgcreatecollection/).

```json
{
  "creator": "bb1abc123...",
  "balancesType": "Standard",
  "defaultBalances": {
    "balances": [],
    "outgoingApprovals": [],
    "incomingApprovals": [],
    "autoApproveSelfInitiatedOutgoingTransfers": false,
    "autoApproveSelfInitiatedIncomingTransfers": true,
    "autoApproveAllIncomingTransfers": false,
    "userPermissions": {
      "canUpdateOutgoingApprovals": [],
      "canUpdateIncomingApprovals": [],
      "canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [],
      "canUpdateAutoApproveSelfInitiatedIncomingTransfers": [],
      "canUpdateAutoApproveAllIncomingTransfers": []
    }
  },
  "validBadgeIds": [{"start": "1", "end": "100"}],
  "collectionPermissions": {
    "canDeleteCollection": [],
    "canArchiveCollection": [],
    "canUpdateOffChainBalancesMetadata": [],
    "canUpdateStandards": [],
    "canUpdateCustomData": [],
    "canUpdateManager": [],
    "canUpdateCollectionMetadata": [],
    "canUpdateValidBadgeIds": [],
    "canUpdateBadgeMetadata": [],
    "canUpdateCollectionApprovals": []
  },
  "managerTimeline": [],
  "collectionMetadataTimeline": [],
  "badgeMetadataTimeline": [],
  "offChainBalancesMetadataTimeline": [],
  "customDataTimeline": [],
  "collectionApprovals": [],
  "standardsTimeline": [],
  "isArchivedTimeline": [],
  "mintEscrowCoinsToTransfer": [],
  "cosmosCoinWrapperPathsToAdd": []
}
```