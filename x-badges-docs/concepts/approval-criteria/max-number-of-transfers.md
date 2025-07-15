# Max Number of Transfers

Limit the number of transfers that can occur using approval trackers.

See [Approval Trackers](./approval-trackers.md) for more information on how trackers work.

## Interface

```typescript
interface ApprovalCriteria<T extends NumberType> {
    maxNumTransfers?: MaxNumTransfers<T>;
}
```

## How It Works

Similar to approval amounts, specify maximum transfers on:

-   **Overall**: Universal limit for all transfers
-   **Per Sender**: Limit per unique sender address
-   **Per Recipient**: Limit per unique recipient address
-   **Per Initiator**: Limit per unique initiator address

"0" means unlimited and not tracked. "N" means max N transfers allowed.

## Example

```json
{
    "maxNumTransfers": {
        "overallMaxNumTransfers": "0",
        "perFromAddressMaxNumTransfers": "0",
        "perToAddressMaxNumTransfers": "0",
        "perInitiatedByAddressMaxNumTransfers": "1",
        "amountTrackerId": "uniqueID"
    }
}
```

Alice can initiate 1 transfer, then no more. Bob can still transfer (different tracker).

```typescript
{
    "fullTrackerId": {
        "numTransfers": 1,
        "amounts": [],
        "lastUpdatedAt": 1691978400000
    }
}
```

## As-Needed Basis

We track on an as-needed basis, meaning if we do not have requirements that use the number of transfers, we will not increment the tracker.

Edge Case: In Predetermined Balances, you may need the number of transfers for determining the balances to assign to each transfer (e.g. transfer #10 -> badge ID 10). In this case, we do need to track the number of transfers. This is all facilitated via the same tracker, so even if you have "0" or unlimited set for the corresponding value in maxNumTransfers, the tracker may be incremented behind the scenes. Consider this when editing / creating approvals. You do not want to use a tracker that has prior history when you expect it to start from scratch.
