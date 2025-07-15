# Action Permissions

Action permissions are the simplest type - they only control when an action can be executed based on time.

## High-Level Logic

```
For each action request:
    Check if current time is in permanentlyPermittedTimes
        → If yes: ALLOW
        → If no: Check if current time is in permanentlyForbiddenTimes
            → If yes: DENY
            → If no: ALLOW (neutral state)
```

**English**: "For these times, this action can be performed" or "For these times, this action is blocked"

## Overview

```
Action Request
    ↓
Time Check
    ↓
┌─────────────────┬─────────────────┐
│ Permitted Times │ Forbidden Times │
└─────────────────┴─────────────────┘
    ↓
Execute Action    Deny Action
```

## Interface

```typescript
interface ActionPermission {
    permanentlyPermittedTimes: UintRange[];
    permanentlyForbiddenTimes: UintRange[];
}
```

## Collection Actions

| Action                | Description              | Use Case       |
| --------------------- | ------------------------ | -------------- |
| `canDeleteCollection` | Delete entire collection | Permanent lock |

## User Actions

| Action                                               | Description                     | Use Case         |
| ---------------------------------------------------- | ------------------------------- | ---------------- |
| `canUpdateAutoApproveSelfInitiatedOutgoingTransfers` | Auto-approve outgoing transfers | User convenience |
| `canUpdateAutoApproveSelfInitiatedIncomingTransfers` | Auto-approve incoming transfers | User convenience |
| `canUpdateAutoApproveAllIncomingTransfers`           | Auto-approve all incoming       | User convenience |

## Examples

### Lock Collection Deletion Forever

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

### Allow Collection Deletion Only During Specific Period

```json
{
    "canDeleteCollection": [
        {
            "permanentlyPermittedTimes": [
                { "start": "1704067200000", "end": "1735689600000" }
            ],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```

### Default Behavior (No Restrictions)

```json
{
    "canDeleteCollection": []
}
```
