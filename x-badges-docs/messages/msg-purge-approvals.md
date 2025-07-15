# MsgPurgeApprovals

A message to purge specific approvals from approval lists. This is a targeted approach that requires specifying exactly which approvals to purge.

## Overview

This message allows you to purge specific approvals by their identifier.

### Usage 1: Self-Purge (Creator purging their own approvals)

-   **`purgeExpired` must be `true`**
-   **`purgeCounterpartyApprovals` must be `false`**
-   **`approvalsToPurge` must contain the specific approvals to purge**
-   Specified approvals will be purged if they are expired (no future transfer times)

### Usage 2: Other-Purge (Creator purging someone else's approvals)

-   Can set either `purgeExpired` or `purgeCounterpartyApprovals` (or both)
-   **`approvalsToPurge` must contain the specific approvals to purge**
-   Purge permissions are determined by the approval's auto-deletion options in `approvalCriteria`:
    -   `allowPurgeIfExpired`: Allows others to purge expired approvals
    -   `allowCounterpartyPurge`: Allows counterparty to purge if they are the only initiator (initiatedByList must be a whitelist with exactly one address matching the counterparty)
-   Specified approvals that match the conditions will be purged

## Fields

-   `creator`: The address submitting the transaction.
-   `collectionId`: The target collection for approval cleanup.
-   `purgeExpired`: Whether to purge expired approvals (must be true for self-purge).
-   `approverAddress`: The address whose approvals to purge. If empty, defaults to `creator`.
-   `purgeCounterpartyApprovals`: Whether to purge counterparty approvals (must be false for self-purge).
-   `approvalsToPurge`: **Required** - An array of approval identifier details specifying exactly which approvals to purge. Cannot be empty.

## ApprovalIdentifierDetails

Each approval to purge must be specified with:

```typescript
interface ApprovalIdentifierDetails {
    approvalId: string; // The ID of the approval
    approvalLevel: string; // "collection", "incoming", or "outgoing"
    approverAddress: string; // Address of the approver (empty for collection-level)
    version: string; // Version of the approval (must match or else we will not purge)
}
```

## Auto-Deletion Options

The following flags in approval criteria control purge permissions in `approvalCriteria`:

-   `allowCounterpartyPurge`: Allows the counterparty to purge the approval if they are the ONLY initiator in `initiatedByList` (must be a whitelist with exactly one address matching the counterparty).
-   `allowPurgeIfExpired`: Allows others (besides the approval owner) to call `PurgeApprovals` on their behalf for expired approvals.

## Permissions

Although user approval permissions are rarely disabled, we still check these purges obey them. If the user does not have permission to purge their own approval, the purge will fail. With counterparty purges, this can be thought of purging on behalf of the user, so the user's permissions are still checked.

## Example Usage

```bash
# [collectionId, purgeExpired, approverAddress, purgeCounterpartyApprovals, approvalsToPurge]

bitbadgeschaind tx badges purge-approvals 1 true "" false '[{"approvalId":"my-approval","approvalLevel":"outgoing","approverAddress":"bb1...","version":"0"}]' --from user-key
```

## Response

The response includes the number of approvals that were successfully purged:

```json
{
    "numPurged": "3"
}
```

## Related Messages

-   [MsgUpdateUserApprovals](./msg-update-user-approvals.md) - Full approval management
-   [MsgSetIncomingApproval](./msg-set-incoming-approval.md) - Set an incoming approval
-   [MsgDeleteIncomingApproval](./msg-delete-incoming-approval.md) - Delete a single incoming approval
-   [MsgSetOutgoingApproval](./msg-set-outgoing-approval.md) - Set a single outgoing approval
-   [MsgDeleteOutgoingApproval](./msg-delete-outgoing-approval.md) - Delete a single outgoing approval
