# Manager

The manager is the central authority for a collection, controlling all administrative operations and having exclusive rights to perform updates, deletions, and other management tasks.

## Manager Timeline

### Structure

```json
"managerTimeline": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "manager": "bb1alice..."
  }
]
```

### Time-Based Manager Changes

Managers can be scheduled to change automatically:

```json
"managerTimeline": [
  {
    "timelineTimes": [{"start": "1", "end": "1672531199000"}],
    "manager": "bb1alice..."
  },
  {
    "timelineTimes": [{"start": "1672531200000", "end": "18446744073709551615"}],
    "manager": "bb1bob..."
  }
]
```

This transfers management from Alice to Bob on January 1, 2023.

## Manager Permissions

The manager role can be granted various permissions, allowing for flexible administration of the collection. These permissions include:

### Core Administrative Permissions

1. **Collection Deletion** - The ability to permanently remove the collection from the system
2. **Collection Archiving** - Archive a collection, making it read-only and rejecting all transactions until unarchived
3. **Core Collection Updates** - Modifying essential details such as metadata URLs and collection standards
4. **Manager Role Transfer** - The ability to pass the manager role to another address
5. **Badge Creation** - Permission to mint additional badges within the collection
6. **Custom Permissions** - Collection-specific permissions depending on setup

### Metadata Management

-   **Collection Metadata Updates** - Modify collection-level metadata and URIs
-   **Badge Metadata Updates** - Update individual badge metadata (with badge-specific permissions)
-   **Timeline Management** - Schedule metadata changes over time

### Transferability Control

-   **Approval Settings** - Modify the collection's approval settings that determine how badges can be transferred
-   **Transfer Rules** - Update transferability conditions and restrictions
-   **Permission Updates** - Configure transferability permissions

### Off-Chain Management

-   **Off-chain Balance Management** - For collections using off-chain balance storage, managers can update these balances
-   **External Integrations** - Manager role can extend to off-chain functionalities and custom utilities

### User-Level Operation Limits

The manager cannot directly:

-   Modify user balances (must follow approval system)
-   Access user private keys or personal data

## Fine-Grained Permission Customizability

One of the key features of the manager role in BitBadges is the ability to customize permissions at a granular level. This allows for precise control over the collection's management.

Permissions can be customized based on various factors:

### Permission Dimensions

-   **Badge Specificity** - Which particular badges within the collection can be affected
-   **Time Constraints** - When can certain actions be performed
-   **Value Limitations** - What specific values or ranges are allowed for updates
-   **Conditional Triggers** - Under what circumstances can certain permissions be exercised

### Permission States

Each permission can exist in one of three states:

1. **Forbidden + Permanently Frozen**

    - The permission is permanently disallowed
    - This state cannot be changed, ensuring certain actions remain off-limits indefinitely

2. **Permitted + Not Frozen**

    - The permission is currently allowed
    - This state can be changed to either of the other two states, offering flexibility in management

3. **Permitted + Permanently Frozen**
    - The permission is permanently allowed
    - Like the first state, this cannot be changed, ensuring certain capabilities always remain available

**Note**: There is no "Forbidden + Not Frozen" state because such a state could theoretically be updated to "Permitted" at any time and then immediately executed, effectively making it a "Permitted" state.

## Permission Control Examples

### Manager Updates

Manager updates are controlled by the `canUpdateManager` permission:

```json
"canUpdateManager": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

## Usage Examples

### Setting Initial Manager

During collection creation:

```json
{
    "creator": "bb1alice...",
    "managerTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "manager": "bb1alice..."
        }
    ],
    "collectionPermissions": {
        "canUpdateManager": [
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

### Decentralized Management Transition

```json
{
    "managerTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "1672531199000" }],
            "manager": "bb1alice..."
        },
        {
            "timelineTimes": [
                { "start": "1672531200000", "end": "18446744073709551615" }
            ],
            "manager": "bb1qqqq...."
        }
    ],
    "collectionPermissions": {
        "canUpdateManager": [
            {
                "permanentlyPermittedTimes": [],
                "permanentlyForbiddenTimes": [
                    { "start": "1672531200000", "end": "18446744073709551615" }
                ]
            }
        ]
    }
}
```

This transitions to a burn address manager and locks management permanently, creating a decentralized collection.
