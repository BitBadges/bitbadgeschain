**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetManager

Sets the manager timeline and update permissions for a collection. This is a convenience message that focuses specifically on manager management.

## Overview

This message allows you to:

-   Set who manages the collection at different times
-   Configure permissions to update the manager in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update the manager timeline according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetManager {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetManager";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New manager timeline to set.
  repeated ManagerTimeline managerTimeline = 3;

  // Permission to update manager timeline
  repeated TimedUpdatePermission canUpdateManager = 4;
}

message MsgSetManagerResponse {
  // ID of the collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-manager '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "managerTimeline": [
        {
            "manager": "bb1def456...",
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canUpdateManager": [
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
