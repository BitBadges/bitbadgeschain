# Badge ID Action Permissions

Badge ID action permissions control which badge-specific actions can be performed based on badge IDs.

## High-Level Logic

```
For each badge action request:
    Check if badge ID matches any badgeIds criteria
        → If no match: ALLOW (neutral state)
        → If match: Check if current time is in permanentlyPermittedTimes
            → If yes: ALLOW
            → If no: Check if current time is in permanentlyForbiddenTimes
                → If yes: DENY
                → If no: ALLOW (neutral state)
```

**English**: "For these times, these badge IDs can be updated" or "For these times, these badge IDs are locked"

## Overview

```
Badge Action
    ↓
Badge ID Match
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
| `canUpdateValidBadgeIds` | Update valid badge IDs | Configuration |

## Examples

### Lock All Badge ID Updates

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

### Lock Specific Badge Range

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

### Allow Future Badge IDs Only

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
