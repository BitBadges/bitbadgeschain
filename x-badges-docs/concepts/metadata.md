# Metadata

BitBadges allows defining metadata for both collections and individual tokens using timeline-based metadata fields. This enables rich, dynamic content that can change over time while maintaining on-chain verifiability.

## Metadata Timelines

### Collection Metadata Timeline

The `collectionMetadataTimeline` defines metadata for the entire collection over time.

```json
"collectionMetadataTimeline": [
  {
    "timelineTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "collectionMetadata": {
      "uri": "ipfs://Qmf8xxN2fwXGgouue3qsJtN8ZRSsnoHxM9mGcynTPhh6Ub",
      "customData": ""
    }
  }
]
```

### Token Metadata Timeline

The `badgeMetadataTimeline` defines metadata for individual tokens over time. The order of `badgeMetadata` entries matters, as it uses a **first-match approach** via linear scan for specific token IDs. BitBadges uses the `{id}` placeholder in the token metadata URI to replace with the actual token ID.

```json
"badgeMetadataTimeline": [
  {
    "timelineTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "badgeMetadata": [
      {
        "uri": "ipfs://Qmf8xxN2fwXGgouue3qsJtN8ZRSsnoHxM9mGcynTPhh6Ub/{id}",
        "badgeIds": [
          {
            "start": "1",
            "end": "10000000000000"
          }
        ],
        "customData": ""
      }
    ]
  }
]
```

## Metadata Interface

The BitBadges API, Indexer, and Site expect metadata to follow this format by default:

```typescript
export interface Metadata<T extends NumberType> {
    name: string;
    description: string;
    image: string;
    video?: string;
    category?: string;
    externalUrl?: string;
    tags?: string[];
    socials?: {
        [key: string]: string;
    };
}
```

## Key Features

### Dynamic Token ID Replacement

-   If the token metadata URI includes `"{id}"`, it's replaced with the actual token ID
-   Example: `"...abc.com/metadata/{id}"` becomes `"...abc.com/metadata/1"` for token ID 1
-   Enables efficient metadata generation for large collections

### First-Match Token Metadata

-   Token metadata entries are evaluated in order
-   First matching entry for a token ID is used
-   Allows specific overrides before general rules

## Permission Control

Metadata updates are controlled by collection permissions:

### Collection Metadata Permission

```json
"canUpdateCollectionMetadata": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

### Token Metadata Permission

```json
"canUpdateBadgeMetadata": [
  {
    "badgeIds": [{"start": "1", "end": "100"}],
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyPermittedTimes": [{"start": "1", "end": "18446744073709551615"}],
    "permanentlyForbiddenTimes": []
  }
]
```

### Timeline Times vs Permission Times Within Permissions

As explained in [Permissions](permissions/), the `timelineTimes` field is used to define the timeline times that can be updated or not. The `permanentlyPermittedTimes` and `permanentlyForbiddenTimes` fields are used to define the times when the permission is enabled or disabled.
