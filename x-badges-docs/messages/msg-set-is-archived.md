**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetIsArchived

Sets the isArchived timeline and update permissions for a collection. This is a convenience message that focuses specifically on archiving management.

## Overview

This message allows you to:

-   Set isArchived timeline for the collection
-   Configure permissions to archive the collection in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to archive the collection according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetIsArchived {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetIsArchived";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New isArchived timeline to set.
  repeated IsArchivedTimeline isArchivedTimeline = 3;

  // Permission to archive collection
  repeated TimedUpdatePermission canArchiveCollection = 4;
}

message MsgSetIsArchivedResponse {
  // ID of the collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-is-archived '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "isArchivedTimeline": [
        {
            "isArchived": true,
            "timelineTimes": [{ "start": "1000", "end": "2000" }]
        }
    ],
    "canArchiveCollection": [
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
