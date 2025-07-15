**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetBadgeMetadata

Sets the badge metadata timeline and update permissions for a badge collection. This is a convenience message that focuses specifically on badge metadata management.

## Overview

This message allows you to:

-   Set badge metadata timeline for the collection
-   Configure permissions to update the badge metadata in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update the badge metadata timeline according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetBadgeMetadata {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetBadgeMetadata";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New badge metadata timeline to set.
  repeated BadgeMetadataTimeline badgeMetadataTimeline = 3;

  // Permission to update badge metadata timeline
  repeated TimedUpdateWithBadgeIdsPermission canUpdateBadgeMetadata = 4;
}

message MsgSetBadgeMetadataResponse {
  // ID of the badge collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-badge-metadata '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "badgeMetadataTimeline": [
        {
            "badgeMetadata": [
                {
                    "uri": "https://example.com/badge1.json",
                    "customData": "{\"description\": \"First badge\"}",
                    "badgeIds": [{ "start": "1", "end": "10" }]
                }
            ],
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canUpdateBadgeMetadata": [
        {
            "badgeIds": [{ "start": "1", "end": "10" }],
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
