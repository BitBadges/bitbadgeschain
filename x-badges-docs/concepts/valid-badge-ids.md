# Valid Token IDs

Valid Token IDs define the range of token IDs that exist within a collection. This is mainly informational but also may be used to enforce certain rules within the collection.

## Creating Token IDs

### During Collection Creation

Use the `validBadgeIds` field in [MsgCreateCollection](../../messages/msg-create-collection.md):

```json
{
    "creator": "bb1...",
    "validBadgeIds": [
        {
            "start": "1",
            "end": "100"
        }
    ],
    "collectionPermissions": {
        "canUpdateValidBadgeIds": [
            // { ... }
        ]
    }
}
```

### During Collection Updates

Use the `validBadgeIds` field in [MsgUpdateCollection](../../messages/msg-update-collection.md):

```json
{
    "creator": "bb1...",
    "collectionId": "1",
    "validBadgeIds": [
        {
            "start": "101",
            "end": "200"
        }
    ]
}
```

## Permission Control

Updates to valid token IDs must obey the `canUpdateValidBadgeIds` permission:

### Permission Structure

```json
"canUpdateValidBadgeIds": [
  {
    "badgeIds": [{"start": "1", "end": "1000"}],
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

### Permission Behaviors

Note that the `canUpdateValidBadgeIds` permission applies to the updatability of the `validBadgeIds` field.

We find the first-match for (current time, token ID) for each token ID that is changed, and check the permission for that time. If no time matches, the permission is default enabled. See [Permissions](permissions/) for more details.

### Permission Best-Practices

Typically, the desired functionality falls into one of the following categories:

-   **Set and Lock All**: Set the valid token IDs upon genesis and lock everything from further updates
-   **Set and Lock All Current, Allow Expansion**: Set the valid token IDs upon genesis and lock the current ones from being updated, but allow expansion in the future.

## Best Practices

1. **Plan ahead**: Consider future expansion when setting initial token ID ranges
2. **Sequential additions**: Always add token IDs sequentially to maintain the no-gaps requirement
3. **Permission management**: Carefully configure `canUpdateValidBadgeIds` permissions based on collection lifecycle
4. **Documentation**: Clearly document the intended use of different token ID ranges
