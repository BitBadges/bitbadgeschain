**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetTokenMetadata

Sets the token metadata timeline and update permissions for a collection. This is a convenience message that focuses specifically on token metadata management.

## Overview

This message allows you to:

-   Set token metadata timeline for the collection
-   Configure permissions to update the token metadata in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update the token metadata timeline according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetTokenMetadata {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetTokenMetadata";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New token metadata timeline to set.
  repeated TokenMetadataTimeline tokenMetadataTimeline = 3;

  // Permission to update token metadata timeline
  repeated TimedUpdateWithTokenIdsPermission canUpdateTokenMetadata = 4;
}

message MsgSetTokenMetadataResponse {
  // ID of the collection.
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
    "tokenMetadataTimeline": [
        {
            "tokenMetadata": [
                {
                    "uri": "https://example.com/badge1.json",
                    "customData": "{\"description\": \"First token\"}",
                    "tokenIds": [{ "start": "1", "end": "10" }]
                }
            ],
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canUpdateTokenMetadata": [
        {
            "tokenIds": [{ "start": "1", "end": "10" }],
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
