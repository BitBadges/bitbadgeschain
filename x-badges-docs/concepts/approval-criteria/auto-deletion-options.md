# Auto-Deletion Options

Automatically delete approvals after specific conditions are met.

## Interface

```typescript
interface AutoDeletionOptions {
    afterOneUse: boolean;
    afterOverallMaxNumTransfers: boolean;
    allowCounterpartyPurge?: boolean;
    allowPurgeIfExpired?: boolean;
}
```

## How It Works

Auto-deletion options allow approvals to be automatically removed when certain conditions are met:

-   **`afterOneUse`**: Delete the approval after it's used once
-   **`afterOverallMaxNumTransfers`**: Delete the approval after the overall max number of transfers threshold is met
-   **`allowCounterpartyPurge`**: If true, allows the counterparty (the only initiator in `initiatedByList`, must be a whitelist with exactly one address) to purge the approval, even if they are not the owner. This may be used for like a rejection of the approval.
-   **`allowPurgeIfExpired`**: If true, allows others (in addition to the approval owner) to purge expired approvals on the owner's behalf. This may be used for a cleanup-like system.

## Usage Examples

### Single-Use Approval

```json
{
    "autoDeletionOptions": {
        "afterOneUse": true,
        "afterOverallMaxNumTransfers": false
    }
}
```

**Result**: Approval is deleted immediately after the first transfer.

### Limited-Use Approval

```json
{
    "maxNumTransfers": {
        "overallMaxNumTransfers": "10"
    },
    "autoDeletionOptions": {
        "afterOneUse": false,
        "afterOverallMaxNumTransfers": true
    }
}
```

**Result**: Approval is deleted after 10 transfers are completed.

### Allow Counterparty Purge

```json
{
    "autoDeletionOptions": {
        "afterOneUse": false,
        "afterOverallMaxNumTransfers": false,
        "allowCounterpartyPurge": true
    }
}
```

**Result**: The counterparty (if they are the only initiator in a whitelist) can purge this approval, even if not the owner.

### Allow Others to Purge Expired Approvals

```json
{
    "autoDeletionOptions": {
        "afterOneUse": false,
        "afterOverallMaxNumTransfers": false,
        "allowPurgeIfExpired": true
    }
}
```

**Result**: Any user can purge this approval if it is expired (no future valid transfer times), not just the owner.
