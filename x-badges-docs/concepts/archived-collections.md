# Archived Collections

Collections can be archived to temporarily or permanently disable all transactions while keeping the collection data verifiable and public on-chain.

## Key Concepts

### Archive State

-   Controlled by the `isArchivedTimeline` field
-   When archived, all transactions fail until unarchived
-   Collection remains readable and verifiable on-chain
-   Does not delete the collection, only makes it read-only

### Timeline-Based Archiving

Collections can be archived for specific time periods:

-   **Temporary archiving** - Archive for maintenance or security
-   **Permanent archiving** - Sunset collections while preserving data
-   **Scheduled archiving** - Pre-planned archive periods

## Implementation

### isArchivedTimeline Structure

```json
"isArchivedTimeline": [
  {
    "timelineTimes": [{"start": "1672531200000", "end": "18446744073709551615"}],
    "isArchived": true
  }
]
```

### Permission Control

Archiving is controlled by the `canArchiveCollection` permission:

```json
"canArchiveCollection": [
  {
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

Note that the `canArchiveCollection` permission is for the updatability of the `isArchivedTimeline` field. It has no bearing on the current value of the `isArchived` field.

For example, when you permanently forbid updating the archive status forever, it could be locked as `true` forever or `false` forever.

## Transaction Behavior

### When Archived

-   **All transactions fail** - No updates, transfers, or changes allowed
-   **Read operations continue** - Queries and data access remain available
-   **Unarchiving exception** - Only unarchiving transactions can succeed

### When Unarchived

-   **Normal operations resume** - All transaction types are allowed
-   **No data loss** - All collection data remains intact
-   **Permissions apply** - Standard permission checks resume

## Archiving a Collection

### During Collection Creation

```json
{
    "creator": "bb1...",
    "isArchivedTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "isArchived": false
        }
    ],
    "collectionPermissions": {
        "canArchiveCollection": [
            {
                "permanentlyPermittedTimes": [
                    { "start": "1", "end": "18446744073709551615" }
                ],
                "permanentlyForbiddenTimes": []
            }
        ]
    }
}
```

### During Collection Updates

Use [MsgUpdateCollection](../../messages/msg-update-collection.md) to update the archive status:

```json
{
    "creator": "bb1...",
    "collectionId": "1",
    "updateIsArchivedTimeline": true,
    "isArchivedTimeline": [
        {
            "timelineTimes": [
                { "start": "1672531200000", "end": "18446744073709551615" }
            ],
            "isArchived": true
        }
    ]
}
```
