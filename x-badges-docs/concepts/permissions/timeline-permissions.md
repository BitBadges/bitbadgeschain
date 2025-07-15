# Timeline Permissions

Timeline permissions control when timeline-based fields can be updated, such as collection metadata and badge metadata.

## High-Level Logic

### Basic Timeline Permissions

```
For each timeline update request:
    Check if timeline time matches any timelineTimes criteria
        → If no match: ALLOW (neutral state)
        → If match: Check if current time is in permanentlyPermittedTimes
            → If yes: ALLOW
            → If no: Check if current time is in permanentlyForbiddenTimes
                → If yes: DENY
                → If no: ALLOW (neutral state)
```

### Badge-Specific Timeline Permissions

```
For each badge timeline update request:
    Check if timeline time AND badge ID match criteria
        → If no match: ALLOW (neutral state)
        → If match: Check if current time is in permanentlyPermittedTimes
            → If yes: ALLOW
            → If no: Check if current time is in permanentlyForbiddenTimes
                → If yes: DENY
                → If no: ALLOW (neutral state)
```

**English**:

-   **Basic**: "For these permission execution times, the (timelineTime -> timelineValue) pairs can be updated"
-   **Badge-Specific**: "For these permission execution times, the (badgeId, timelineTime -> timelineValue) pairs can be updated"

## Timeline vs Execution Times

-   **Timeline Times**: Which timeline values can be updated?
-   **Execution Times**: When the permission can be executed?

These may not align. For example, you might forbid updating timeline values for Jan 2024 during 2023.

## Overview

```
Timeline Update
    ↓
Timeline Time Match
    ↓
Time Permission Check
    ↓
Execute/Deny
```

## Types

### Basic Timeline Permissions

Control collection-level timeline updates:

```typescript
interface TimedUpdatePermission {
    timelineTimes: UintRange[];
    permanentlyPermittedTimes: UintRange[];
    permanentlyForbiddenTimes: UintRange[];
}
```

**Available Actions:**

-   `canArchiveCollection`
-   `canUpdateOffChainBalancesMetadata`
-   `canUpdateStandards`
-   `canUpdateCustomData`
-   `canUpdateManager`
-   `canUpdateCollectionMetadata`

### Badge-Specific Timeline Permissions

Control badge metadata timeline updates:

```typescript
interface TimedUpdateWithBadgeIdsPermission {
    timelineTimes: UintRange[];
    badgeIds: UintRange[];
    permanentlyPermittedTimes: UintRange[];
    permanentlyForbiddenTimes: UintRange[];
}
```

**Available Actions:**

-   `canUpdateBadgeMetadata`

## Examples

### Lock Collection Metadata Forever

```json
{
    "canUpdateCollectionMetadata": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Lock Specific Timeline Period

```json
{
    "canUpdateCollectionMetadata": [
        {
            "timelineTimes": [{ "start": "1000", "end": "2000" }],
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Lock Badge Metadata for Existing Badges

```json
{
    "canUpdateBadgeMetadata": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "badgeIds": [{ "start": "1", "end": "100" }],
            "permanentlyPermittedTimes": [],
            "permanentlyForbiddenTimes": [
                { "start": "1", "end": "18446744073709551615" }
            ]
        }
    ]
}
```

### Allow Updates Only During Specific Period

```json
{
    "canUpdateCollectionMetadata": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "permanentlyPermittedTimes": [
                { "start": "1704067200000", "end": "1735689600000" }
            ],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```
