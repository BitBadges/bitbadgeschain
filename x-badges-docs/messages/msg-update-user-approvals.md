# MsgUpdateUserApprovals

Updates a user's approval settings for token transfers.

## Collection ID Auto-Lookup

If you specify `collectionId` as `"0"`, it will automatically lookup the latest collection ID created. This can be used if you are creating a collection and do not know the official collection ID yet but want to perform a multi-message transaction.

## Update Flag Pattern

This message uses an update flag + value pattern for selective updates. Each updatable field has a corresponding boolean flag (e.g., `updateOutgoingApprovals`, `updateIncomingApprovals`, `updateAutoApproveSelfInitiatedOutgoingTransfers`).

-   **If update flag is `true`**: The corresponding value field is processed and the user's settings are updated with the new value
-   **If update flag is `false`**: The corresponding value field is completely ignored, regardless of what data is provided

This allows you to update only specific approval settings without affecting others, and you can safely leave unused value fields empty or with placeholder data.

## Authorization & Permissions

Users can only update their own approvals. Updates must be performed according to the permissions set (i.e. the `userPermissions` previously set for that user).

**Note**: Typically, user permissions are almost always permanently allowed/set to enabled. These permissions only need to be customized in advanced cases where fine-grained control over user approval updates is required.

## Proto Definition

```protobuf
message MsgUpdateUserApprovals {
  string creator = 1; // User updating their approval settings
  string collectionId = 2; // Target collection for approval updates
  bool updateOutgoingApprovals = 3;
  repeated UserOutgoingApproval outgoingApprovals = 4;
  bool updateIncomingApprovals = 5;
  repeated UserIncomingApproval incomingApprovals = 6;
  bool updateAutoApproveSelfInitiatedOutgoingTransfers = 7;
  bool autoApproveSelfInitiatedOutgoingTransfers = 8;
  bool updateAutoApproveSelfInitiatedIncomingTransfers = 9;
  bool autoApproveSelfInitiatedIncomingTransfers = 10;
  bool updateAutoApproveAllIncomingTransfers = 11;
  bool autoApproveAllIncomingTransfers = 12;
  bool updateUserPermissions = 13;
  UserPermissions userPermissions = 14;
}

message MsgUpdateUserApprovalsResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges update-user-approved-transfers '[tx-json]' --from user-key
```

### JSON Example

For complete transaction examples, see [MsgUpdateUserApprovals Examples](../examples/txs/msgupdate-user-approvals/).

```json
{
    "creator": "bb1user123...",
    "collectionId": "1",

    "updateOutgoingApprovals": false,
    "outgoingApprovals": [],

    "updateIncomingApprovals": false,
    "incomingApprovals": [],

    "updateAutoApproveSelfInitiatedOutgoingTransfers": true,
    "autoApproveSelfInitiatedOutgoingTransfers": true,

    "updateAutoApproveSelfInitiatedIncomingTransfers": false,
    "autoApproveSelfInitiatedIncomingTransfers": true,

    "updateAutoApproveAllIncomingTransfers": false,
    "autoApproveAllIncomingTransfers": false,

    "updateUserPermissions": false,
    "userPermissions": {
        "canUpdateOutgoingApprovals": [],
        "canUpdateIncomingApprovals": [],
        "canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [],
        "canUpdateAutoApproveSelfInitiatedIncomingTransfers": [],
        "canUpdateAutoApproveAllIncomingTransfers": []
    }
}
```
