# Auto-Deletion Options

Automatically delete approvals after specific conditions are met.

## Interface

```typescript
interface AutoDeletionOptions {
    afterOneUse: boolean;
    afterOverallMaxNumTransfers: boolean;
}
```

## How It Works

Auto-deletion options allow approvals to be automatically removed when certain conditions are met:

-   **`afterOneUse`**: Delete the approval after it's used once
-   **`afterOverallMaxNumTransfers`**: Delete the approval after the overall max number of transfers threshold is met

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
