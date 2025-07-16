# Transferability & Approvals

## Overview

Transferability in BitBadges is controlled through a hierarchical approval system with three levels:

<figure><img src="../../../.gitbook/assets/image (33).png" alt="Approval hierarchy diagram"><figcaption>Approval hierarchy: Collection → User (Incoming/Outgoing)</figcaption></figure>

### Approval Levels

| Level          | Description                              | Fields                                    |
| -------------- | ---------------------------------------- | ----------------------------------------- |
| **Collection** | Global rules for the entire collection   | All fields                                |
| **Incoming**   | User-specific rules for receiving badges | `toList` = user's address, no overrides   |
| **Outgoing**   | User-specific rules for sending badges   | `fromList` = user's address, no overrides |

**Key Rule**: A transfer must satisfy collection-level approvals AND (unless overridden) user-level incoming/outgoing approvals.

## Approval Structure

```typescript
interface CollectionApproval<T extends NumberType> {
    // Core Fields
    toListId: string; // Who can receive?
    fromListId: string; // Who can send?
    initiatedByListId: string; // Who can initiate?
    transferTimes: UintRange<T>[]; // When can transfer happen?
    badgeIds: UintRange<T>[]; // Which badge IDs?
    ownershipTimes: UintRange<T>[]; // Which ownership times?
    approvalId: string; // Unique identifier

    // Version control (incremented on each update)
    version: T;

    // Optional Fields
    uri?: string; // Metadata link
    customData?: string; // Custom data
    approvalCriteria?: ApprovalCriteria<T>; // Additional restrictions
}
```

See [Approval Criteria](approval-criteria/) for more details on the `approvalCriteria` field.

## Approval Value vs Permission

While the value may seem similar to the approval update permissions, the permission corresponds to the **updatability** of the approvals (i.e. `canUpdateCollectionApprovals`). The approvals themselves correspond to if a transfer is currently approved or not.

## The Six Core Fields

Every approval defines **Who? When? What?** through these fields:

| Field               | Type                            | Purpose                                 | Example                                            |
| ------------------- | ------------------------------- | --------------------------------------- | -------------------------------------------------- |
| `toListId`          | Address List ID                 | Who can receive badges                  | `"All"`, `"Mint"`, `"bb1..."`                      |
| `fromListId`        | Address List ID                 | Who can send badges                     | `"Mint"`, `"!Mint"`                                |
| `initiatedByListId` | Address List ID                 | Who can initiate transfer               | `"All"`, `"bb1..."`                                |
| `transferTimes`     | UintRange[] (UNIX Milliseconds) | When transfer can occur                 | `[{start: "1691931600000", end: "1723554000000"}]` |
| `badgeIds`          | UintRange[] (Badge IDs)         | Which badge IDs                         | `[{start: "1", end: "100"}]`                       |
| `ownershipTimes`    | UintRange[] (UNIX Milliseconds) | Which ownership times to be transferred | `[{start: "1", end: "18446744073709551615"}]`      |

### Example Approval

```json
{
    "fromListId": "Mint",
    "toListId": "All",
    "initiatedByListId": "All",
    "transferTimes": [{ "start": "1691931600000", "end": "1723554000000" }],
    "badgeIds": [{ "start": "1", "end": "100" }],
    "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
    "approvalId": "mint-to-all"
}
```

**Translation**: Allow anyone to claim badges 1-100 from the Mint address between Aug 13, 2023 and Aug 13, 2024.

## Approval Matching Process

### Step-by-Step Flow (High-Level)

```
Transfer Request
    ↓
Check Collection Approvals
    ↓
Match Found? ──No──→ TRANSFER DENIED
    ↓ Yes
Check Approval Criteria
    ↓
Criteria Met? ──No──→ TRANSFER DENIED
    ↓ Yes
Check User Approvals (where necessary)
    ↓
User Approvals Met? ──No──→ TRANSFER DENIED
    ↓ Yes
TRANSFER APPROVED
```

### Matching Logic

For optimal design, you should try to design transfers such that they only use specific approvals without the need for splitting. However, if needed, we split the transfer / approvals as fine-grained as we can to make it succeed. In other words, we deduct as much as possible from each approval as we iterate.

### Prioritized Approvals

In MsgTransferBadges, you can specify which approvals to prioritize. This allows you to prioritize certain approvals over others.

```typescript
// In MsgTransferBadges
{
  prioritizedApprovals: [{
    approvalId: "approval1",
    approvalLevel: "collection",
    approverAdress: "", // ""bb1" if approvalLevel is "incoming" or "outgoing",
    version: 0,
  }],
  onlyCheckPrioritizedCollectionApprovals: true,
  // If true, the transfer will be denied if no prioritized approvals match
}
```

## Auto-Scan vs Prioritized Approvals

The transfer approval system operates in two modes to balance efficiency and precision:

### Auto-Scan Mode (Default)

By default, the system automatically scans through available approvals to find a match for the transfer. This mode:

-   **Works with**: Approvals using [Empty Approval Criteria](../examples/empty-approval-criteria.md) (no side effects). For example, when you approve all incoming transfers w/ no restrictions, this has no side effects.
-   **Behavior**: Automatically finds and uses the first matching approval
-   **Use case**: Simple transfers without custom logic or side effects
-   **No versioning required**: The system handles approval selection automatically

### Prioritized Approvals (Required for Side Effects)

**CRITICAL REQUIREMENT**: Any transfer with side effects or custom approval criteria MUST always be prioritized with proper versioning set. No exceptions.

#### Race Condition Protection

The versioning control ensures that before submitting, the user knows the exact approval they are using:

```typescript
"prioritizedApprovals": [
    {
        "approvalId": "abc123",
        "approvalLevel": "collection",
        "approverAddress": "",
        "version": "2" // Must specify exact version
  }
]
```

See [MsgTransferBadges](../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md) for the complete message structure.

## Related Topics

-   [Approval Criteria](approval-criteria/) - Additional restrictions and challenges
-   [Address Lists](../address-lists.md) - Managing address groups
-   [UintRanges](../uint-ranges.md) - Range logic implementation
-   [Permissions](permissions/) - Controlling who can update approvals
