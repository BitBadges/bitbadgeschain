# 📨 Messages

This directory contains detailed documentation for all message types supported by the badges module.

## Message Categories

### Collection Management

-   [MsgCreateCollection](./msg-create-collection.md) - Create new badge collection
-   [MsgUpdateCollection](./msg-update-collection.md) - Update existing collection properties
-   [MsgUniversalUpdateCollection](./msg-universal-update-collection.md) - Universal create/update interface with invariants support
-   [MsgDeleteCollection](./msg-delete-collection.md) - Archive/delete collection

### Helper Collection Update Messages

-   [MsgSetValidBadgeIds](./msg-set-valid-badge-ids.md) - Update valid badge IDs and permissions
-   [MsgSetManager](./msg-set-manager.md) - Update manager timeline and permissions
-   [MsgSetCollectionMetadata](./msg-set-collection-metadata.md) - Update collection metadata timeline and permissions
-   [MsgSetBadgeMetadata](./msg-set-badge-metadata.md) - Update badge metadata timeline and permissions
-   [MsgSetCustomData](./msg-set-custom-data.md) - Update custom data timeline and permissions
-   [MsgSetStandards](./msg-set-standards.md) - Update standards timeline and permissions
-   [MsgSetCollectionApprovals](./msg-set-collection-approvals.md) - Update collection approvals and permissions
-   [MsgSetIsArchived](./msg-set-is-archived.md) - Update isArchived timeline and permissions

### Badge Transfers

-   [MsgTransferBadges](./msg-transfer-badges.md) - Transfer badges between addresses with approval validation

### User Approval Management

-   [MsgUpdateUserApprovals](./msg-update-user-approvals.md) - Update user transfer approval settings
-   [MsgSetIncomingApproval](./msg-set-incoming-approval.md) - Set a single incoming approval (helper)
-   [MsgDeleteIncomingApproval](./msg-delete-incoming-approval.md) - Delete a single incoming approval (helper)
-   [MsgSetOutgoingApproval](./msg-set-outgoing-approval.md) - Set a single outgoing approval (helper)
-   [MsgDeleteOutgoingApproval](./msg-delete-outgoing-approval.md) - Delete a single outgoing approval (helper)
-   [MsgPurgeApprovals](./msg-purge-approvals.md) - Purge expired approvals (helper)

### Address List Management

-   [MsgCreateAddressLists](./msg-create-address-lists.md) - Create reusable address lists for access control

### Dynamic Store Management

-   [MsgCreateDynamicStore](./msg-create-dynamic-store.md) - Create numeric stores for approval criteria
-   [MsgUpdateDynamicStore](./msg-update-dynamic-store.md) - Update dynamic store configuration
-   [MsgDeleteDynamicStore](./msg-delete-dynamic-store.md) - Delete dynamic store
-   [MsgSetDynamicStoreValue](./msg-set-dynamic-store-value.md) - Set individual address values in dynamic store
-   [MsgIncrementStoreValue](./msg-increment-store-value.md) - Increase values for addresses
-   [MsgDecrementStoreValue](./msg-decrement-store-value.md) - Decrease values for addresses

## Additional Message Types

The following message types exist in the protocol but may be documented separately:

-   **MsgUpdateParams** - Update module parameters via governance
