# User Royalties

Apply percentage-based royalties to transfers.

## Interface

```typescript
interface UserRoyalties {
    percentage: string; // 1 to 10000 represents basis points (0.01% to 100%)
    payoutAddress: string; // Address to receive the royalties
}
```

## How It Works

User royalties automatically deduct a percentage from transfers and send it to a specified payout address:

-   **Percentage**: Expressed in basis points (1 = 0.01%, 100 = 1%, 10000 = 100%)
-   **Payout**: Automatically sent to the specified address on each transfer
-   **Deduction**: Applied to the transfer amount before the transfer is processed

## Usage Examples

### 5% Royalty

```json
{
    "userRoyalties": {
        "percentage": "500", // 500 basis points = 5%
        "payoutAddress": "bb1creator..."
    }
}
```

**Result**: 5% of each transfer amount is sent to the creator's address.

### 2.5% Royalty

```json
{
    "userRoyalties": {
        "percentage": "250", // 250 basis points = 2.5%
        "payoutAddress": "bb1artist..."
    }
}
```

**Result**: 2.5% of each transfer amount is sent to the artist's address.

## Edge Case: One per Transfer

Currently, we only support one specific royalty percentage applied per transfer. If a transfer matches to different approvals with multiple royalties, the transfer may fail.
