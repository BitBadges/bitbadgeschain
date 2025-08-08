# MsgUpdateCollection

Updates an existing collection's properties.

## Update Flag Pattern

This message uses an update flag + value pattern for selective updates. Each updatable field has a corresponding boolean flag (e.g., `updateValidBadgeIds`, `updateCollectionPermissions`).

-   **If update flag is `true`**: The corresponding value field is processed and the collection is updated with the new value
-   **If update flag is `false`**: The corresponding value field is completely ignored, regardless of what data is provided

This allows you to update only specific fields without affecting others, and you can safely leave unused value fields empty or with placeholder data.

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. All updates must obey the previously set permissions - meaning the permission settings that were in effect _before_ this message was started.

**Important**: If you update the permissions in the current message, those new permissions are applied last and will not be applicable until the following transaction. This prevents circumventing permission restrictions within the same transaction.

## Proto Definition

```protobuf
message MsgUpdateCollection {
  string creator = 1; // Address updating collection (must be manager)
  string collectionId = 2; // ID of collection to update
  bool updateValidBadgeIds = 3;
  repeated UintRange validBadgeIds = 4;
  bool updateCollectionPermissions = 7;
  CollectionPermissions collectionPermissions = 8;
  bool updateManagerTimeline = 9;
  repeated ManagerTimeline managerTimeline = 10;
  bool updateCollectionMetadataTimeline = 11;
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 12;
  bool updateBadgeMetadataTimeline = 13;
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 14;
  bool updateOffChainBalancesMetadataTimeline = 15;
  repeated OffChainBalancesMetadataTimeline offChainBalancesMetadataTimeline = 16;
  bool updateCustomDataTimeline = 17;
  repeated CustomDataTimeline customDataTimeline = 18;
  bool updateCollectionApprovals = 21;
  repeated CollectionApproval collectionApprovals = 22;
  bool updateStandardsTimeline = 23;
  repeated StandardsTimeline standardsTimeline = 24;
  bool updateIsArchivedTimeline = 27;
  repeated IsArchivedTimeline isArchivedTimeline = 28;
  repeated cosmos.base.v1beta1.Coin mintEscrowCoinsToTransfer = 29;
  repeated CosmosCoinWrapperPathAddObject cosmosCoinWrapperPathsToAdd = 30;
  CollectionInvariants invariants = 31;
}

message MsgUpdateCollectionResponse {
  string collectionId = 1; // ID of updated collection
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges update-collection '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "updateValidBadgeIds": true,
    "validBadgeIds": [{ "start": "1", "end": "200" }],
    "updateCollectionPermissions": false,
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
    "updateManagerTimeline": false,
    "managerTimeline": [],
    "updateCollectionMetadataTimeline": false,
    "collectionMetadataTimeline": [],
    "updateBadgeMetadataTimeline": false,
    "badgeMetadataTimeline": [],
    "updateOffChainBalancesMetadataTimeline": false,
    "offChainBalancesMetadataTimeline": [],
    "updateCustomDataTimeline": false,
    "customDataTimeline": [],
    "updateCollectionApprovals": false,
    "collectionApprovals": [],
    "updateStandardsTimeline": false,
    "standardsTimeline": [],
    "updateIsArchivedTimeline": false,
    "isArchivedTimeline": [],
    "mintEscrowCoinsToTransfer": [],
    "cosmosCoinWrapperPathsToAdd": [],
    "invariants": {
        "noCustomOwnershipTimes": false
    }
}
```
