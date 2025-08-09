**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetValidBadgeIds

Sets the valid token IDs and update permissions for a collection. This is a convenience message that focuses specifically on token ID management.

## Overview

This message allows you to:

-   Set which token IDs are valid for the collection
-   Configure permissions to update the valid token IDs in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update valid token IDs according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetValidBadgeIds {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetValidBadgeIds";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New token IDs to add to this collection
  repeated UintRange validBadgeIds = 3;

  // Permission to update valid token IDs
  repeated BadgeIdsActionPermission canUpdateValidBadgeIds = 4;
}

message MsgSetValidBadgeIdsResponse {
  // ID of the collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-valid-badge-ids '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "validBadgeIds": [
        { "start": "1", "end": "100" },
        { "start": "200", "end": "300" }
    ],
    "canUpdateValidBadgeIds": [
        {
            "badgeIds": [{ "start": "1", "end": "50" }],
            "permanentlyPermittedTimes": [{ "start": "1000", "end": "2000" }],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```

## Related Messages

-   [MsgUniversalUpdateCollection](./msg-universal-update-collection.md) - Full collection update with all fields
-   [MsgUpdateCollection](./msg-update-collection.md) - Legacy update message
