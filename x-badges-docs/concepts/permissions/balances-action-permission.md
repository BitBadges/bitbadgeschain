# Badge IDs Action Permission

The BadgeIdsAction permission denotes for what (badge ID, ownership times), can an action be executed? For example, can I create more of badge ID 1-10?

This permission refers to the UPDATABILITY of the balances and has no bearing on what the circulating supplys are currently set to.

```json
"collectionPermissions": {
    "canUpdateValidBadgeIds": [...],
    ...
}
```

```typescript
export interface BadgeIdsActionPermission<T extends NumberType> {
  badgeIds: UintRange<T>[];
  
  permanentlyPermittedTimes: UintRange<T>[];
  permanentlyForbiddenTimes: UintRange<T>[];
}
```

{% content-ref url="../balances-transfers/creating-badges.md" %}
[Creating Badges](../balances-transfers/creating-badges.md)
{% endcontent-ref %}

```json
"canUpdateValidBadgeIds": [
  {
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
    ]
  }
]
```
