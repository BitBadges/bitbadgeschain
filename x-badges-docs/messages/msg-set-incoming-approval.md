# MsgSetIncomingApproval

A helper message to set a single incoming approval for token transfers. This is a developer-friendly wrapper around `MsgUpdateUserApprovals` that simplifies setting individual incoming approvals. For more information, we refer to the [MsgUpdateUserApprovals](./msg-update-user-approvals.md) documentation.

## Overview

This message allows you to set or update a single incoming approval without having to construct the full `MsgUpdateUserApprovals` message. It automatically handles version management and validation.

## Proto Definition

```protobuf
message MsgSetIncomingApproval {
  string creator = 1; // User setting the approval
  string collectionId = 2; // Target collection for approval
  UserIncomingApproval approval = 3; // The incoming approval to set
}

message MsgSetIncomingApprovalResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-incoming-approval [collection-id] '[approval-json]' --from user-key
```

## Behavior

-   **New Approval**: If the approval ID doesn't exist, a new approval is created with version 0
-   **Update Existing**: If the approval ID already exists, the approval is updated and the version is incremented
-   **No Change**: If the approval content hasn't changed, the version remains the same
-   **Validation**: The approval is validated according to the collection's permissions and user's approval update permissions

## Authorization & Permissions

Users can only set their own incoming approvals. The operation must be performed according to the permissions set (i.e. the `userPermissions` previously set for that user).

## Related Messages

-   [MsgUpdateUserApprovals](./msg-update-user-approvals.md) - Full approval management
-   [MsgDeleteIncomingApproval](./msg-delete-incoming-approval.md) - Delete an incoming approval
-   [MsgSetOutgoingApproval](./msg-set-outgoing-approval.md) - Set an outgoing approval
