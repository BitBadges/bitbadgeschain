# Timed Update With Badge Ids Permission

```json
"collectionPermissions": {
    "canUpdateBadgeMetadata": [...],
    ...
}
```

```typescript
export interface TimedUpdateWithBadgeIdsPermission<T extends NumberType> {
  timelineTimes: UintRange<T>[];
  badgeIds: UintRange<T>[];
  
  permanentlyPermittedTimes: UintRange<T>[];
  permanentlyForbiddenTimes: UintRange<T>[];
}
```

**TimedUpdatePermissionWithBadgeId**s simply denote for what **timelineTimes** and **badgeIds** combinations, can the manager update the scheduled value? These are only applicable to badge ID [timeline-based fields](../timelines.md) such as the badge metadata timeline. This permission refers to the UPDATABILITY of the timeline and has no bearing on what the timeline is currently set to.

The **timelineTimes** are which timeline time values can be updated. The **badgeIds** are which badge IDs can be updated. For a pair such as (Mon-Fri, IDs 1-10), this means the values corresponding to the badge IDs at the timeline times can be updated or not. Both have to match. IDs 11+ are not handled at all in this case. Sunday is not handled. Updating IDs 1-10 on Sunday is not handled.

**Examples**

Below, this forbids updating the entire timeline because all timelineTimes and badgeIds are specified.

```json
"canUpdateBadgeMetadata": [
  {
    "timelineTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "badgeIds": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
  }
]
```

A commonly set value for this permission may look like the following. Let's say you have a collection with 100 badges but in the future, you can create new badges 101+. This permission allows you to freeze the metadata of the current badges but allow you to set the metadata for any new badges in the future.

```json
"canUpdateBadgeMetadata": [
  {
    "timelineTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "badgeIds": [
      {
        "start": "1",
        "end": "100"
      }
    ],
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
  }
]
```
