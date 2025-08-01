**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetCollectionMetadata

Sets the collection metadata timeline and update permissions for a badge collection. This is a convenience message that focuses specifically on collection metadata management.

## Overview

This message allows you to:

-   Set collection metadata timeline for the collection
-   Configure permissions to update the collection metadata in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update the collection metadata timeline according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetCollectionMetadata {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetCollectionMetadata";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New collection metadata timeline to set.
  repeated CollectionMetadataTimeline collectionMetadataTimeline = 3;

  // Permission to update collection metadata timeline
  repeated TimedUpdatePermission canUpdateCollectionMetadata = 4;
}

message MsgSetCollectionMetadataResponse {
  // ID of the badge collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-collection-metadata '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "collectionMetadataTimeline": [
        {
            "collectionMetadata": {
                "uri": "https://example.com/collection.json",
                "customData": "{\"description\": \"My collection\"}"
            },
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canUpdateCollectionMetadata": [
        {
            "timelineTimes": [{ "start": "1000", "end": "2000" }],
            "permanentlyPermittedTimes": [{ "start": "1000", "end": "2000" }],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```

## Related Messages

-   [MsgUniversalUpdateCollection](./msg-universal-update-collection.md) - Full collection update with all fields
-   [MsgUpdateCollection](./msg-update-collection.md) - Legacy update message
