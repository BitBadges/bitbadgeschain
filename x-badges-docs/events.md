# Events

The badges module emits events for all message operations to enable blockchain monitoring and external application integration.

## Event Categories

### Standard Message Events

All message handlers emit `sdk.EventTypeMessage` events with message-specific attributes.

### Indexer Events

Duplicate events with type "indexer" for external application consumption.

### Transfer Events

Detailed events for transfer operations including approval usage and challenge tracking.

### IBC Events

Events for cross-chain operations with acknowledgment handling.

## Standard Message Events

### Collection Management

#### CreateCollection

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "create_collection"
  - msg: string (JSON-encoded message)
```

#### UpdateCollection

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "update_collection"
  - msg: string (JSON-encoded message)
```

#### UniversalUpdateCollection

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "universal_update_collection"
  - msg: string (JSON-encoded message)
```

#### DeleteCollection

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "delete_collection"
  - msg: string (JSON-encoded message)
```

### Badge Transfers

#### TransferBadges

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (initiator address)
  - action: "transfer_badges"
  - msg: string (JSON-encoded message)
  - transfer: string (JSON transfer details)
  - from: string (sender address)
  - to: string (recipient address)
  - initiatedBy: string (initiator address)
  - coinTransfers: string (JSON coin transfer details)
  - approvalsUsed: string (JSON approval usage details)
  - balances: string (JSON balance details)
```

### User Approvals

#### UpdateUserApprovals

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (user address)
  - action: "update_user_approvals"
  - msg: string (JSON-encoded message)
```

### Address Lists

#### CreateAddressLists

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "create_address_lists"
  - msg: string (JSON-encoded message)
```

### Dynamic Stores

#### CreateDynamicStore

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (creator address)
  - action: "create_dynamic_store"
  - store_id: string
  - msg: string (JSON-encoded message)
```

#### UpdateDynamicStore

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (updater address)
  - action: "update_dynamic_store"
  - store_id: string
  - msg: string (JSON-encoded message)
```

#### DeleteDynamicStore

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (deleter address)
  - action: "delete_dynamic_store"
  - store_id: string
  - msg: string (JSON-encoded message)
```

#### SetDynamicStoreValue

```
Type: "message"
Attributes:
  - module: "badges"
  - sender: string (setter address)
  - action: "set_dynamic_store_value"
  - store_id: string
  - address: string (target address)
  - value: string ("true"/"false")
  - msg: string (JSON-encoded message)
```

## Transfer Events

### Approval Usage

```
Type: "usedApprovalDetails"
Attributes:
  - collectionId: string
  - approverAddress: string
  - approvalId: string
  - amountTrackerId: string
  - approvalLevel: string
  - trackerType: string
  - address: string
  - amounts: string (JSON array)
  - numTransfers: string
  - lastUpdatedAt: string
```

### Challenge Events

```
Type: "challenge{approvalId}{challengeId}{leafIndex}{approverAddress}{approvalLevel}{newNumUsed}"
Attributes:
  - challengeTrackerId: string
  - approvalId: string
  - leafIndex: string
  - approverAddress: string
  - approvalLevel: string
  - numUsed: string
```

### Dynamic Approval Events

```
Type: "approval{collectionId}{approverAddress}{approvalId}{amountsTrackerId}{approvalLevel}{trackerType}{address}"
Attributes:
  - amountTrackerId: string
  - approvalId: string
  - approverAddress: string
  - approvalLevel: string
  - trackerType: string
  - approvedAddress: string
  - amounts: string (JSON array)
  - numTransfers: string
  - lastUpdatedAt: string
```

## IBC Events

### Packet Events

```
Type: "timeout" (for timeouts)
Attributes:
  - acknowledgement: string
  - success: string ("true"/"false")
  - error: string (if applicable)
```
