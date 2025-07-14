# 📨 Messages

This directory contains detailed documentation for all message types supported by the badges module.

## Message Categories

### Collection Management

-   [MsgCreateCollection](./msg-create-collection.md) - Create new badge collection
-   [MsgUpdateCollection](./msg-update-collection.md) - Update existing collection properties
-   [MsgDeleteCollection](./msg-delete-collection.md) - Archive/delete collection

### Badge Transfers

-   [MsgTransferBadges](./msg-transfer-badges.md) - Transfer badges between addresses with approval validation

### User Approval Management

-   [MsgUpdateUserApprovals](./msg-update-user-approvals.md) - Update user transfer approval settings

### Address List Management

-   [MsgCreateAddressLists](./msg-create-address-lists.md) - Create reusable address lists for access control

### Dynamic Store Management

-   [MsgCreateDynamicStore](./msg-create-dynamic-store.md) - Create boolean flag stores for approval criteria
-   [MsgUpdateDynamicStore](./msg-update-dynamic-store.md) - Update dynamic store configuration
-   [MsgDeleteDynamicStore](./msg-delete-dynamic-store.md) - Delete dynamic store
-   [MsgSetDynamicStoreValue](./msg-set-dynamic-store-value.md) - Set individual address values in dynamic store

## Additional Message Types

The following message types exist in the protocol but may be documented separately:

-   **MsgUniversalUpdateCollection** - Legacy unified create/update interface
-   **MsgUpdateParams** - Update module parameters via governance
