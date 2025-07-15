# Approval Trackers

Track transfer amounts and counts using increment-only tallies with thresholds.

## How It Works

Trackers use an incrementing tally system with thresholds:

1. **Setup**: Approved for x10 of badge IDs 1-10 with tracker ID "xyz"
2. **Transfer x5**: Tracker "xyz" goes from 0/10 → 5/10
3. **Transfer x5**: Tracker "xyz" goes to 10/10
4. **Transfer x1**: Exceeds threshold, transfer fails

## Tracker Identification

Tracker IDs include multiple components:

```
ID: collectionId-approvalLevel-approverAddress-approvalId-amountTrackerId-trackerType-approvedAddress
```

### Tracker ID Details Interface

```typescript
interface ApprovalTrackerIdDetails<T extends NumberType> {
    collectionId: T;
    approvalLevel: 'collection' | 'incoming' | 'outgoing' | '';
    approvalId: string;
    approverAddress: string;
    amountTrackerId: string;
    trackerType: 'overall' | 'to' | 'from' | 'initiatedBy' | '';
    approvedAddress: string;
}
```

### Component Breakdown

-   **collectionId**: The collection this tracker belongs to
-   **approvalLevel**: Level of approval ("collection", "incoming", "outgoing", or empty)
-   **approvalId**: Unique identifier for the specific approval
-   **approverAddress**: Address of the approver (empty for collection-level)
-   **amountTrackerId**: User-defined tracker identifier specified in approvalAmounts or maxNumTransfers (see below)
-   **trackerType**: Type of tracking ("overall", "to", "from", "initiatedBy", or empty)
-   **approvedAddress**: Specific address being tracked (empty for "overall")

```typescript
interface iApprovalAmounts<T extends NumberType> {
    amountTrackerId: string; // Key for tracking tallies
}

interface iMaxNumTransfers<T extends NumberType> {
    amountTrackerId: string; // Key for tracking tallies
}
```

### Tracker Types

-   **"overall"**: Universal tally for any transfer (approvedAddress empty)
-   **"to"**: Per-recipient tally (approvedAddress = recipient)
-   **"from"**: Per-sender tally (approvedAddress = sender)
-   **"initiatedBy"**: Per-initiator tally (approvedAddress = initiator)

## Increment Only and Immutable

Trackers are increment only and immutable in storage. To start an approval tally from scratch, you will need to map the approval to a new unused tracker ID. This can be done simply by editing `amountTrackerId` (because this changes the whole ID) or restructuring to change one of the other fields that make up the overall ID.

**IMPORTANT**: Because of the immutable nature, be careful to not revert to a previously used ID unintentionally because the starting point will be the previous tally (not starting from scratch).

## As-Needed Basis

Only increment when necessary (e.g., if no amount restrictions, don't track amounts). Meaning, if there is no need to increment the tally (unlimited limit and/or not restrictions), we do not increment for efficiency purposes. For example, if we only have requirements for numTransfers but do not need the amounts, we do not increment the amounts.

### Example Tracker States

```json
{
    "fullTrackerId1": {
        "numTransfers": 5,
        "amounts": [
            {
                "amount": 50,
                "badgeIds": [{ "start": 1, "end": 10 }],
                "ownershipTimes": [{ "start": 1, "end": 100000000000 }]
            }
        ],
        "lastUpdatedAt": 1691978400000
    },
    "fullTrackerId2": {
        "numTransfers": 3,
        "amounts": [
            {
                "amount": 15,
                "badgeIds": [{ "start": 1, "end": 5 }],
                "ownershipTimes": [{ "start": 1, "end": 100000000000 }]
            }
        ],
        "lastUpdatedAt": 1691978400000
    }
}
```

## Periodic Resets

Trackers support periodic resets to zero using time intervals.

Leave the values at 0 to disable periodic resets.

```typescript
interface ResetTimeIntervals<T extends NumberType> {
    startTime: T; // Original start time of the first interval
    intervalLength: T; // Interval length in unix milliseconds
}
```

### How It Works

-   **First Update**: If it's the first update of the interval, all tracker progress is reset to zero
-   **Recurring**: Useful for recurring subscriptions (e.g., one transfer per month)
-   **No Reset**: Set both values to 0 for no periodic resets

### Example

```json
{
    "approvalAmounts": {
        "overallApprovalAmount": "100",
        "amountTrackerId": "monthly-tracker",
        "resetTimeIntervals": {
            "startTime": "1691978400000", // Aug 13, 2023
            "intervalLength": "2592000000" // 30 days in milliseconds
        }
    }
}
```

This creates a monthly reset cycle starting from August 13, 2023.
