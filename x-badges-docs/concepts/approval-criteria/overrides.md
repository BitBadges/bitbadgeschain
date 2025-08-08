# Override User Level Approvals

Collection-level approvals can override user-level approvals to force transfers.

## Interface

```typescript
interface ApprovalCriteria<T extends NumberType> {
    overridesFromOutgoingApprovals?: boolean;
    overridesToIncomingApprovals?: boolean;
}
```

## How It Works

-   **`overridesFromOutgoingApprovals: true`**: Skip sender's outgoing approvals
-   **`overridesToIncomingApprovals: true`**: Skip recipient's incoming approvals

This enables forced transfers without user consent.

## Use Cases

-   **Force Revoke**: Remove tokens from users
-   **Freeze Tokens**: Prevent transfers regardless of user settings
-   **Emergency Actions**: Administrative control over transfers

## Mint Address Requirement

**CRITICAL**: Mint address approvals must always override outgoing approvals:

```json
{
    "fromListId": "Mint",
    "approvalCriteria": {
        "overridesFromOutgoingApprovals": true
    }
}
```

The Mint address has no user-level approvals, so overrides are required for functionality.
