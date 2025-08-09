# Overview

Permissions control who can perform actions on collections and user balances, and when those actions can be executed.

Typically, user permissions are always set to allowed / empty. In advanced cases like escrows, user permissions can be set to forbidden to make certain on-chain actions impossible.

Collection permissions are important to be set correctly: If there is no manager, the collection permission values do not matter.

## Overview

```
Action Request
    ↓
Permission Check
    ↓
┌─────────────────┬─────────────────┐
│ Collection      │ User            │
│ Permissions     │ Permissions     │
│ (Manager Only)  │ (User Control)  │
└─────────────────┴─────────────────┘
    ↓
Permitted/Forbidden
    ↓
Execute/Deny
```

## Permission Types

| Type           | Scope            | Purpose                          |
| -------------- | ---------------- | -------------------------------- |
| **Collection** | Manager only     | Control collection-level actions |
| **User**       | Individual users | Control user-specific actions    |

## Permission States

**Note**: Once a permission is set to permanently permitted or forbidden, it cannot be changed.

| State                     | Description           | Behavior               |
| ------------------------- | --------------------- | ---------------------- |
| **Permanently Permitted** | Action ALWAYS allowed | Can be executed        |
| **Permanently Forbidden** | Action ALWAYS blocked | Cannot be executed     |
| **Neutral**               | Not specified         | **Allowed by default** |

There is no forbidden + not frozen state because theoretically, it could be updated to permitted at any time and executed (thus making it permitted).

## Time Control

All permissions support time-based control via UNIX millisecond UintRanges.

```json
{
    "permanentlyPermittedTimes": [{ "start": "1", "end": "1000" }],
    "permanentlyForbiddenTimes": [
        { "start": "1001", "end": "18446744073709551615" }
    ]
}
```

## Permission Categories

There are **five types** of permissions, each with different criteria:

-   **[Action Permissions](action-permissions.md)** - Simple time-based permissions (no criteria)
-   **[Timeline Permissions](timeline-permissions.md)** - Control timeline updates (timelineTimes)
-   **[Timeline with Token IDs](timeline-permissions.md)** - Control token-specific timeline updates (timelineTimes + tokenIds)
-   **[Token ID Action Permissions](badge-id-permissions.md)** - Control token-specific actions (tokenIds)
-   **[Approval Permissions](approval-permissions.md)** - Control approval updates (transfer criteria + approvalId)

### Correct Categorization

Based on the proto definitions:

**Action Permissions** (only time control):

-   `canDeleteCollection` (collection)
-   `canUpdateAutoApproveSelfInitiatedOutgoingTransfers` (user)
-   `canUpdateAutoApproveSelfInitiatedIncomingTransfers` (user)
-   `canUpdateAutoApproveAllIncomingTransfers` (user)

**Timeline Permissions** (timelineTimes + time control):

-   `canArchiveCollection`
-   `canUpdateOffChainBalancesMetadata`
-   `canUpdateStandards`
-   `canUpdateCustomData`
-   `canUpdateManager`
-   `canUpdateCollectionMetadata`

**Timeline with Token IDs** (timelineTimes + tokenIds + time control):

-   `canUpdateTokenMetadata`

**Token ID Action Permissions** (tokenIds + time control):

-   `canUpdateValidTokenIds`

**Approval Permissions** (transfer criteria + approvalId + time control):

-   `canUpdateCollectionApprovals` (collection)
-   `canUpdateIncomingApprovals` (user)
-   `canUpdateOutgoingApprovals` (user)

## Quick Examples

### Lock Collection Deletion

When can the collection be deleted?

```json
{
    "canDeleteCollection": [
        {
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Freeze Token Metadata

When can the token metadata be updated? And which (token IDs, timeline time) pairs does it apply to?

```json
{
    "canUpdateTokenMetadata": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "tokenIds": [{ "start": "1", "end": "100" }],

            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

## First Match Policy

Permissions are evaluated as a linear array where each element has criteria and time controls. Only the **first matching element** is applied - all subsequent matches are ignored.

### Key Rules

-   **First Match Only**: Only the first element that matches all criteria is used
-   **Deterministic State**: Each criteria combination has exactly one permission state
-   **No Overlap**: Times cannot be in both `permanentlyPermittedTimes` and `permanentlyForbiddenTimes`
-   **Order Matters**: Array order affects which permissions are applied

### Example: Timeline Permissions

```json
"canUpdateCollectionMetadata": [
    {
        "timelineTimes": [{ "start": "1", "end": "10" }],
        "permanentlyPermittedTimes": [],
        "permanentlyForbiddenTimes": [{ "start": "1", "end": "10" }]
    },
    {
        "timelineTimes": [{ "start": "1", "end": "100" }],
        "permanentlyPermittedTimes": [{ "start": "1", "end": "18446744073709551615" }],
        "permanentlyForbiddenTimes": []
    }
]
```

**Result:**

-   Timeline times 1-10: **Forbidden** (first element matches, second element does not)
-   Timeline times 11-100: **Permitted** (second element matches)

## Satisfying Criteria

All criteria in a permission element must match for it to be applied. Partial matches are ignored.

### Example: Token Metadata Permissions

```json
"canUpdateTokenMetadata": [
    {
        "timelineTimes": [{ "start": "1", "end": "10" }],
        "tokenIds": [{ "start": "1", "end": "10" }],
        "permanentlyPermittedTimes": [{ "start": "1", "end": "18446744073709551615" }],
        "permanentlyForbiddenTimes": []
    }
]
```

**This permission only covers:**

-   Timeline times 1-10 AND token IDs 1-10

**It does NOT cover:**

-   Timeline time 1 with token ID 11
-   Timeline time 11 with token ID 1
-   Timeline time 11 with token ID 11

These combinations are **unhandled** and **allowed by default** since they do not match the permission criteria.

## Brute Force Pattern

To lock specific criteria, you must specify the target and set all other criteria to maximum ranges.

### Example: Lock Token IDs 1-10

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

## Important Notes

-   **First Match Policy**: Only the first matching permission is applied
-   **Default Allow**: Unspecified permissions are allowed by default
-   **Manager Required**: Collection permissions require a manager
-   **User Control**: User permissions typically remain empty for full control
-   **Brute Force**: Use maximum ranges to ensure complete coverage
-   **Order Matters**: Array order affects permission evaluation
