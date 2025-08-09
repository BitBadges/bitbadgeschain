**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetCustomData

Sets the custom data timeline and update permissions for a collection. This is a convenience message that focuses specifically on custom data management.

## Overview

This message allows you to:

-   Set custom data timeline for the collection
-   Configure permissions to update the custom data in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update the custom data timeline according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetCustomData {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetCustomData";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New custom data timeline to set.
  repeated CustomDataTimeline customDataTimeline = 3;

  // Permission to update custom data timeline
  repeated TimedUpdatePermission canUpdateCustomData = 4;
}

message MsgSetCustomDataResponse {
  // ID of the collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-custom-data '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "customDataTimeline": [
        {
            "customData": "{\"description\": \"My custom data\", \"version\": \"1.0\"}",
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canUpdateCustomData": [
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
