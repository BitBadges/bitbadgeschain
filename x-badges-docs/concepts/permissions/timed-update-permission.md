# Timed Update Permission

```json
"collectionPermissions": {
    "canArchiveCollection": [...],
    "canUpdateOffChainBalancesMetadata": [...],
    "canUpdateStandards": [...],
    "canUpdateCustomData": [...],
    "canUpdateManager": [...],
    "canUpdateCollectionMetadata": [...],
    ...
}
```

```typescript
export interface TimedUpdatePermission<T extends NumberType> {
  timelineTimes: UintRange<T>[];
  
  permanentlyPermittedTimes: UintRange<T>[];
  permanentlyForbiddenTimes: UintRange<T>[];
}
```

**TimedUpdatePermission**s simply denote for what **timelineTimes**, can the manager update the scheduled value? These are only applicable to normal [timeline-based fields](../timelines.md) such as the collection metadata timeline. This permission refers to the UPDATABILITY of the timeline and has no bearing on what the timeline is currently set to.

The **timelineTimes** are which timeline time values can be updated. The permitted / forbidden times are when the permission can be executed (the update can take place). Note these may not be aligned. Maybe, you want to forbid updating the timeline from Jan 2024 - Dec 2024 during 2023.

**Examples**

Below, this forbids updating the entire timeline because all timelineTimes are specified.

```json
"canUpdateCollectionMetadata": [
  {
    "timelineTimes": [
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

Below, this forbids ever updating the times 1000-2000 only. All other times can still be updated.

```json
"canUpdateCollectionMetadata": [
  {
    "timelineTimes": [
      {
        "start": "1000",
        "end": "2000"
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
