# Valid Token IDs

Valid Token IDs define the range of token IDs that exist within a collection. This is mainly informational but also may be used to enforce certain rules within the collection.

## Creating Token IDs

### During Collection Creation

Use the `validTokenIds` field in [MsgCreateCollection](../../messages/msg-create-collection.md):

```json
{
    "creator": "bb1...",
    "validTokenIds": [
        {
            "start": "1",
            "end": "100"
        }
    ],
    "collectionPermissions": {
        "canUpdateValidTokenIds": [
            // { ... }
        ]
    }
}
```

### During Collection Updates

Use the `validTokenIds` field in [MsgUpdateCollection](../../messages/msg-update-collection.md):

```json
{
    "creator": "bb1...",
    "collectionId": "1",
    "validTokenIds": [
        {
            "start": "101",
            "end": "200"
        }
    ]
}
```

## Permission Control

Updates to valid token IDs must obey the `canUpdateValidTokenIds` permission:

### Permission Structure

```json
"canUpdateValidTokenIds": [
  {
    "tokenIds": [{"start": "1", "end": "1000"}],
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

### Permission Behaviors

Note that the `canUpdateValidTokenIds` permission applies to the updatability of the `validTokenIds` field.

We find the first-match for (current time, token ID) for each token ID that is changed, and check the permission for that time. If no time matches, the permission is default enabled. See [Permissions](permissions/) for more details.

### Permission Best-Practices

Typically, the desired functionality falls into one of the following categories:

-   **Set and Lock All**: Set the valid token IDs upon genesis and lock everything from further updates
-   **Set and Lock All Current, Allow Expansion**: Set the valid token IDs upon genesis and lock the current ones from being updated, but allow expansion in the future.

## Best Practices

1. **Plan ahead**: Consider future expansion when setting initial token ID ranges
2. **Sequential additions**: Always add token IDs sequentially to maintain the no-gaps requirement
3. **Permission management**: Carefully configure `canUpdateValidTokenIds` permissions based on collection lifecycle
4. **Documentation**: Clearly document the intended use of different token ID ranges
