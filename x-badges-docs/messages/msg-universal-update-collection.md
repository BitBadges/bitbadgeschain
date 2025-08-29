# MsgUniversalUpdateCollection

A universal message that can be used to either create a new collection or update an existing one. This message combines the functionality of both `MsgCreateCollection` and `MsgUpdateCollection` into a single interface.

## Dual Purpose

-   **Collection Creation**: When `collectionId` is set to `"0"`, this message creates a new collection
-   **Collection Update**: When `collectionId` is set to an existing collection ID, this message updates that collection

## Update Flag Pattern

This message uses an update flag + value pattern for selective updates. Each updatable field has a corresponding boolean flag (e.g., `updateValidBadgeIds`, `updateCollectionPermissions`).

-   **If update flag is `true`**: The corresponding value field is processed and the collection is updated with the new value
-   **If update flag is `false`**: The corresponding value field is completely ignored, regardless of what data is provided

## Authorization & Permissions

-   **For Collection Creation**: Can be executed by any address
-   **For Collection Updates**: Can only be executed by the **current manager** of the collection. All updates must obey the previously set permissions.

## Proto Definition

```protobuf
message MsgUniversalUpdateCollection {
  string creator = 1; // Address creating/updating collection
  string collectionId = 2; // "0" for new collection, existing ID for updates

  // Creation-only fields (only used when collectionId = "0")
  string balancesType = 3; // "Standard", "Off-Chain - Indexed", etc.
  UserBalanceStore defaultBalances = 4;

  // Updateable fields (used for both creation and updates)
  repeated UintRange validBadgeIds = 5;
  bool updateCollectionPermissions = 6;
  CollectionPermissions collectionPermissions = 7;
  bool updateManagerTimeline = 8;
  repeated ManagerTimeline managerTimeline = 9;
  bool updateCollectionMetadataTimeline = 10;
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 11;
  bool updateBadgeMetadataTimeline = 12;
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 13;
  bool updateOffChainBalancesMetadataTimeline = 14;
  repeated OffChainBalancesMetadataTimeline offChainBalancesMetadataTimeline = 15;
  bool updateCustomDataTimeline = 16;
  repeated CustomDataTimeline customDataTimeline = 17;
  bool updateCollectionApprovals = 18;
  repeated CollectionApproval collectionApprovals = 19;
  bool updateStandardsTimeline = 20;
  repeated StandardsTimeline standardsTimeline = 21;
  bool updateIsArchivedTimeline = 22;
  repeated IsArchivedTimeline isArchivedTimeline = 23;

  // Transfer fields
  repeated cosmos.base.v1beta1.Coin mintEscrowCoinsToTransfer = 24;
  repeated CosmosCoinWrapperPathAddObject cosmosCoinWrapperPathsToAdd = 25;

  // Invariants (creation-only)
  CollectionInvariants invariants = 26;
}

message MsgUniversalUpdateCollectionResponse {
  string collectionId = 1; // ID of created/updated collection
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges universal-update-collection '[tx-json]' --from creator-key
```

### JSON Example - Creating a New Collection

```json
{
    "creator": "bb1abc123...",
    "collectionId": "0",
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
    "validBadgeIds": [{ "start": "1", "end": "100" }],
    "updateCollectionPermissions": true,
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
    "updateManagerTimeline": true,
    "managerTimeline": [],
    "updateCollectionMetadataTimeline": true,
    "collectionMetadataTimeline": [],
    "updateBadgeMetadataTimeline": true,
    "badgeMetadataTimeline": [],
    "updateOffChainBalancesMetadataTimeline": true,
    "offChainBalancesMetadataTimeline": [],
    "updateCustomDataTimeline": true,
    "customDataTimeline": [],
    "updateCollectionApprovals": true,
    "collectionApprovals": [],
    "updateStandardsTimeline": true,
    "standardsTimeline": [],
    "updateIsArchivedTimeline": true,
    "isArchivedTimeline": [],
    "mintEscrowCoinsToTransfer": [],
    "cosmosCoinWrapperPathsToAdd": [],
    "invariants": {
        "noCustomOwnershipTimes": false,
        "maxSupplyPerId": "0"
    }
}
```

### JSON Example - Updating an Existing Collection

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "updateValidBadgeIds": true,
    "validBadgeIds": [{ "start": "1", "end": "200" }],
    "updateCollectionPermissions": false,
    "collectionPermissions": {},
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
    "invariants": {}
}
```

## Key Differences from Other Messages

### vs MsgCreateCollection

-   More flexible update flag pattern
-   Can be used for both creation and updates
-   Includes invariants support

### vs MsgUpdateCollection

-   Can create new collections when collectionId = "0"
-   Includes creation-only fields like `balancesType` and `defaultBalances`
-   Includes invariants support

## Invariants Support

When creating a new collection (collectionId = "0"), you can set collection invariants using the `invariants` field. Invariants cannot be modified after collection creation.

```json
{
    "invariants": {
        "noCustomOwnershipTimes": true,
        "maxSupplyPerId": "0"
    }
}
```

## Related Messages

-   [MsgCreateCollection](./msg-create-collection.md)
-   [MsgUpdateCollection](./msg-update-collection.md)
-   [Collection Invariants](../concepts/collection-invariants.md)
