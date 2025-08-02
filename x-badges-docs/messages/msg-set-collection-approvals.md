**Disclaimer:**  
This message is a streamlined alternative to [MsgUpdateCollection](./msg-update-collection.md). If you need to update many fields at once, we recommend using MsgUpdateCollection instead.

# MsgSetCollectionApprovals

Sets the collection approvals and update permissions for a badge collection. This is a convenience message that focuses specifically on collection approvals management.

## Overview

This message allows you to:

-   Set collection approvals for the collection
-   Configure permissions to update the collection approvals in the future

## Authorization & Permissions

Updates can only be performed by the **current manager** of the collection. The manager must have permission to update collection approvals according to the collection's current permission settings.

## Proto Definition

```protobuf
message MsgSetCollectionApprovals {
  option (cosmos.msg.v1.signer) = "creator";
  option (amino.name) = "badges/SetCollectionApprovals";

  // Address of the creator.
  string creator = 1;

  // ID of the collection.
  string collectionId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // New collection approvals to set.
  repeated CollectionApproval collectionApprovals = 3;

  // Permission to update collection approvals
  repeated CollectionApprovalPermission canUpdateCollectionApprovals = 4;
}

message MsgSetCollectionApprovalsResponse {
  // ID of the badge collection.
  string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-collection-approvals '[tx-json]' --from manager-key
```

### JSON Example

```json
{
    "creator": "bb1abc123...",
    "collectionId": "1",
    "collectionApprovals": [
        {
            "fromListId": "list1",
            "toListId": "list2",
            "initiatedByListId": "list3",
            "transferTimes": [{ "start": "1000", "end": "2000" }],
            "badgeIds": [{ "start": "1", "end": "10" }],
            "ownershipTimes": [{ "start": "1", "end": "100" }],
            "approvalId": "approval1",
            "approvalCriteria": {
                "mustOwnBadges": [],
                "merkleChallenges": [],
                "ethSignatureChallenges": [],
                "coinTransfers": [],
                "predeterminedBalances": null,
                "approvalAmounts": null,
                "autoDeletionOptions": null,
                "maxNumTransfers": null,
                "dynamicStoreChallenges": []
            }
        }
    ],
    "canUpdateCollectionApprovals": [
        {
            "fromListId": "list1",
            "toListId": "list2",
            "initiatedByListId": "list3",
            "transferTimes": [{ "start": "1000", "end": "2000" }],
            "badgeIds": [{ "start": "1", "end": "10" }],
            "ownershipTimes": [{ "start": "1", "end": "100" }],
            "approvalId": "approval1",
            "permanentlyPermittedTimes": [{ "start": "1000", "end": "2000" }],
            "permanentlyForbiddenTimes": []
        }
    ]
}
```

## Related Messages

-   [MsgUniversalUpdateCollection](./msg-universal-update-collection.md) - Full collection update with all fields
-   [MsgUpdateCollection](./msg-update-collection.md) - Legacy update message
