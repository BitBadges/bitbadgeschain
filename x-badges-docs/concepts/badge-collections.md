# Collections

A collection is the primary entity that defines a group of related tokens with shared properties and rules. We refer you to other pages for more details on the different concepts that make up a collection.

Note: This is what is stored on-chain in storage for a collection. You may typically interact with similar concepts but moreso in Messages and Queries format.

## Proto Definition

```protobuf
message BadgeCollection {
  // The unique identifier for this collection. This is assigned by the blockchain. First collection has ID 1.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // The metadata for the collection itself, which can vary over time.
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 2;

  // The metadata for each token in the collection, also subject to changes over time.
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 3;

  // The type of balances this collection uses ("Standard", "Off-Chain - Indexed", "Off-Chain - Non-Indexed", or "Non-Public").
  string balancesType = 4;

  // Metadata for fetching balances for collections with off-chain balances, subject to changes over time.
  repeated OffChainBalancesMetadataTimeline offChainBalancesMetadataTimeline = 5;

  // An arbitrary field that can store any data, subject to changes over time.
  repeated CustomDataTimeline customDataTimeline = 7;

  // The address of the manager of this collection, subject to changes over time.
  repeated ManagerTimeline managerTimeline = 8;

  // Permissions that define what the manager of the collection can do or not do.
  CollectionPermissions collectionPermissions = 9;

  // Transferability of the collection for collections with standard balances, subject to changes over time.
  // Overrides user approvals for a transfer if specified.
  // Transfer must satisfy both user and collection-level approvals.
  // Only applicable to on-chain balances.
  repeated CollectionApproval collectionApprovals = 10;

  // Standards that define how to interpret the fields of the collection, subject to changes over time.
  repeated StandardsTimeline standardsTimeline = 11;

  // Whether the collection is archived or not, subject to changes over time.
  // When archived, it becomes read-only, and no transactions can be processed until it is unarchived.
  repeated IsArchivedTimeline isArchivedTimeline = 12;

  // The default store of a balance for a user, upon genesis.
  UserBalanceStore defaultBalances = 13;

  // The user or entity who created the collection.
  string createdBy = 14;

  // The valid token IDs for this collection.
  repeated UintRange validBadgeIds = 15;

  // The generated address of the collection. Also used to escrow Mint balances.
  string mintEscrowAddress = 16;

  // The IBC wrapper (sdk.coin) paths for the collection.
  repeated CosmosCoinWrapperPath cosmosCoinWrapperPaths = 17;

  // Collection-level invariants that cannot be broken.
  // These are set upon genesis and cannot be modified.
  CollectionInvariants invariants = 18;
}
```
