# 📚 Overview

This directory contains comprehensive developer documentation for the BitBadges blockchain's `x/badges` module.

This section is a knowledge dump for how tokens operate behind the scenes. For most use cases, you will not care about any of this as it will be handled for you via the site. And if you are self-implementing a token-gated service, you can just fetch balances and metadata from the API without worrying about the underlying details.

```typescript
const res = await BitBadgesApi.getBalanceByAddress(collectionId, address, {
    ...options,
});
console.log(res);

const res = await BitBadgesApi.getBadgeMetadata(1, 5);
```

## Table of Contents

1. [Introduction](./introduction.md) - Overview and key concepts
2. [Concepts](./02-concepts.md) - Core data structures and business logic
3. [State](./state.md) - State management and storage patterns
4. [Messages](./messages/) - Transaction messages and handlers
5. [Queries](./queries/) - Query types and endpoints
6. [Events](./events.md) - Event emissions and tracking
7. [Examples](./examples/) - Common usage patterns and building blocks

## Message Reference

### Collection Management

-   [MsgCreateCollection](./messages/msg-create-collection.md) - Create new collection
-   [MsgUpdateCollection](./messages/msg-update-collection.md) - Update existing collection
-   [MsgUniversalUpdateCollection](./messages/msg-universal-update-collection.md) - Universal create/update interface with invariants support
-   [MsgDeleteCollection](./messages/msg-delete-collection.md) - Delete collection

### Token Transfers

-   [MsgTransferBadges](./messages/msg-transfer-badges.md) - Transfer tokens between addresses

### User Approvals

-   [MsgUpdateUserApprovals](./messages/msg-update-user-approvals.md) - Update transfer approvals

### Address Lists & Dynamic Stores

-   [MsgCreateAddressLists](./messages/msg-create-address-lists.md) - Create reusable address lists
-   [MsgCreateDynamicStore](./messages/msg-create-dynamic-store.md) - Create key-value store
-   [MsgUpdateDynamicStore](./messages/msg-update-dynamic-store.md) - Update dynamic store properties
-   [MsgDeleteDynamicStore](./messages/msg-delete-dynamic-store.md) - Delete dynamic store
-   [MsgSetDynamicStoreValue](./messages/msg-set-dynamic-store-value.md) - Set address-specific store values
-   [MsgIncrementStoreValue](./messages/msg-increment-store-value.md) - Increase values for addresses
-   [MsgDecrementStoreValue](./messages/msg-decrement-store-value.md) - Decrease values for addresses
-   [More messages...](./messages/) - See full message reference

## Query Reference

### Core Queries

-   [GetCollection](./queries/get-collection.md) - Retrieve collection data
-   [GetBalance](./queries/get-balance.md) - Get user balances
-   [GetApprovalTracker](./queries/get-approval-tracker.md) - Get approval usage data
-   [GetAddressList](./queries/get-address-list.md) - Retrieve address list
-   [More queries...](./queries/) - See full query reference

## Quick Links

-   [BitBadges Chain Repository](https://github.com/bitbadges/bitbadgeschain)
-   [BitBadges Documentation](https://docs.bitbadges.io)
-   [Proto Definitions](https://github.com/bitbadges/bitbadgeschain/tree/master/proto/badges)

## Documentation Style

This documentation follows the [Cosmos SDK module documentation standards](https://docs.cosmos.network/main/building-modules/README) and is designed for developers building on or integrating with the BitBadges blockchain.
