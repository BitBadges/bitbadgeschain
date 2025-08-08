# MsgDeleteOutgoingApproval

A helper message to delete a single outgoing approval for token transfers. This is a developer-friendly wrapper around `MsgUpdateUserApprovals` that simplifies deleting individual outgoing approvals. For more information, we refer to the [MsgUpdateUserApprovals](./msg-update-user-approvals.md) documentation.

## Overview

This message allows you to delete a single outgoing approval by its ID without having to construct the full `MsgUpdateUserApprovals` message with an empty approval list.

## Proto Definition

```protobuf
message MsgDeleteOutgoingApproval {
  string creator = 1; // User deleting the approval
  string collectionId = 2; // Target collection for approval
  string approvalId = 3; // The ID of the approval to delete
}

message MsgDeleteOutgoingApprovalResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges delete-outgoing-approval [collection-id] [approval-id] --from user-key
```

### Example

```bash
bitbadgeschaind tx badges delete-outgoing-approval 1 "my-approval-1" --from user-key
```

## Behavior

-   **Approval Lookup**: The system searches for an outgoing approval with the specified `approvalId`
-   **Deletion**: If found, the approval is removed from the user's outgoing approvals list
-   **Error Handling**: If the approval ID is not found, an error is returned
-   **Validation**: The deletion is validated according to the collection's permissions and user's approval update permissions

## Authorization & Permissions

Users can only delete their own outgoing approvals. The operation must be performed according to the permissions set (i.e. the `userPermissions` previously set for that user).

## Related Messages

-   [MsgUpdateUserApprovals](./msg-update-user-approvals.md) - Full approval management
-   [MsgSetOutgoingApproval](./msg-set-outgoing-approval.md) - Set an outgoing approval
-   [MsgDeleteIncomingApproval](./msg-delete-incoming-approval.md) - Delete an incoming approval
