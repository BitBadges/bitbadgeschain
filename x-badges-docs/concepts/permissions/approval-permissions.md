# Approval Permissions

Approval permissions control when transfer approvals can be updated, allowing you to freeze specific transfer rules.

## High-Level Logic

```
For each approval update request:
    Check if approval criteria match (from, to, initiatedBy, transferTimes, tokenIds, ownershipTimes, approvalId)
        → If no match: ALLOW (neutral state)
        → If match: Check if current time is in permanentlyPermittedTimes
            → If yes: ALLOW
            → If no: Check if current time is in permanentlyForbiddenTimes
                → If yes: DENY
                → If no: ALLOW (neutral state)
```

**English**: "For these permission execution times, the approvals matching to these criteria can be updated"

## Overview

```
Approval Update
    ↓
Transfer Match
    ↓
Approval ID Match
    ↓
Time Permission Check
    ↓
Execute/Deny
```

## Interface

```typescript
interface ApprovalPermission {
    fromListId: string;
    toListId: string;
    initiatedByListId: string;
    transferTimes: UintRange[];
    tokenIds: UintRange[];
    ownershipTimes: UintRange[];
    approvalId: string;

    permanentlyPermittedTimes: UintRange[];
    permanentlyForbiddenTimes: UintRange[];
}
```

## Available Actions

| Action                         | Scope      | Description                         |
| ------------------------------ | ---------- | ----------------------------------- |
| `canUpdateCollectionApprovals` | Collection | Control collection-level approvals  |
| `canUpdateIncomingApprovals`   | User       | Control incoming transfer approvals |
| `canUpdateOutgoingApprovals`   | User       | Control outgoing transfer approvals |

**Note**: For user approvals, `fromListId` and `toListId` are automatically set:

-   **Incoming**: `toListId` is hardcoded to the user's address
-   **Outgoing**: `fromListId` is hardcoded to the user's address

## Key Concepts

### Approval Tuple

An approval tuple consists of: `(from, to, initiatedBy, tokenIds, transferTimes, ownershipTimes, approvalId)`

### Brute Force Pattern

To lock specific criteria, specify the target and set all other criteria to maximum ranges:

```json
{
    "fromListId": "All",
    "toListId": "All",
    "initiatedByListId": "All",
    "tokenIds": [{ "start": "1", "end": "10" }],
    "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
    "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
    "approvalId": "All",
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [
        { "start": "1", "end": "18446744073709551615" }
    ]
}
```

## Examples

### Lock Specific ID Range

```json
{
    "canUpdateCollectionApprovals": [
        {
            "fromListId": "All",
            "toListId": "All",
            "initiatedByListId": "All",
            "tokenIds": [{ "start": "1", "end": "100" }],
            "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "approvalId": "All",
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Lock Specific Approval ID

```json
{
    "canUpdateCollectionApprovals": [
        {
            "fromListId": "All",
            "toListId": "All",
            "initiatedByListId": "All",
            "tokenIds": [{ "start": "1", "end": "18446744073709551615" }],
            "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "approvalId": "specific-approval-id",
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

## Protection Strategies

### 1. Specific Approval Lock

Lock a specific approval by its unique ID:

```json
"approvalId": "unique-approval-id"
```

### 2. Range Lock with Overlap Protection

Lock a token range AND all overlapping approvals:

```json
// Lock token range
{
    "tokenIds": [{ "start": "1", "end": "10" }],
    "approvalId": "All"
}

// Lock overlapping approval
{
    "approvalId": "overlapping-approval-id"
}
```

### 3. Complete Freeze

Lock all approvals for a collection:

```json
{
    "fromListId": "All",
    "toListId": "All",
    "initiatedByListId": "All",
    "tokenIds": [{ "start": "1", "end": "18446744073709551615" }],
    "transferTimes": [{ "start": "1", "end": "18446744073709551615" }],
    "ownershipTimes": [{ "start": "1", "end": "18446744073709551615" }],
    "approvalId": "All",
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [
        { "start": "1", "end": "18446744073709551615" }
    ]
}
```
