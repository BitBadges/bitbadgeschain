# Token ID Action Permissions

Token ID action permissions control which token-specific actions can be performed based on token IDs.

## High-Level Logic

```
For each token action request:
    Check if token ID matches any badgeIds criteria
        → If no match: ALLOW (neutral state)
        → If match: Check if current time is in permanentlyPermittedTimes
            → If yes: ALLOW
            → If no: Check if current time is in permanentlyForbiddenTimes
                → If yes: DENY
                → If no: ALLOW (neutral state)
```

**English**: "For these times, these token IDs can be updated" or "For these times, these token IDs are locked"

## Overview

```
Token Action
    ↓
Token ID Match
    ↓
Time Permission Check
    ↓
Execute/Deny
```

## Interface

```typescript
interface BadgeIdsActionPermission {
    badgeIds: UintRange[];
    permanentlyPermittedTimes: UintRange[];
    permanentlyForbiddenTimes: UintRange[];
}
```

## Available Actions

| Action                   | Description            | Use Case      |
| ------------------------ | ---------------------- | ------------- |
| `canUpdateValidBadgeIds` | Update valid token IDs | Configuration |

## Examples

### Lock All Token ID Updates

```json
{
    "canUpdateValidBadgeIds": [
        {
            "badgeIds": [{ "start": "1", "end": "18446744073709551615" }],
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Lock Specific ID Range

```json
{
    "canUpdateValidBadgeIds": [
        {
            "badgeIds": [{ "start": "1", "end": "100" }],
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Allow Future Token IDs Only

```json
{
    "canUpdateValidBadgeIds": [
        {
            "badgeIds": [{ "start": "101", "end": "18446744073709551615" }],
            "permanentlyPermittedTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```
