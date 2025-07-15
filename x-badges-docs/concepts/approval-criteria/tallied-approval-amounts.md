# Tallied Approval Amounts

Limit transfer amounts using increment-only trackers with thresholds.

## Interface

```typescript
interface ApprovalCriteria<T extends NumberType> {
    approvalAmounts?: ApprovalAmounts<T>;
}
```

## How It Works

Specify maximum amounts that can be transferred using four tracker types:

-   **Overall** (`trackerType = "overall"`): Universal limit for all transfers
-   **Per To Address** (`trackerType = "to"`): Limit per unique recipient
-   **Per From Address** (`trackerType = "from"`): Limit per unique sender
-   **Per Initiated By Address** (`trackerType = "initiatedBy"`): Limit per unique initiator

"0" means unlimited and not tracked. "N" means max N amount allowed.

## Example

```json
{
    "approvalAmounts": {
        "overallApprovalAmount": "1000",
        "perFromAddressApprovalAmount": "0",
        "perToAddressApprovalAmount": "0",
        "perInitiatedByAddressApprovalAmount": "10",
        "amountTrackerId": "uniqueID"
    }
}
```

## Tracker Types

### Overall Tracker

-   **ID**: `1-collection- -approvalId-uniqueID-overall-`
-   **Behavior**: Increments for all transfers regardless of sender/recipient/initiator
-   **Use Case**: Global collection limits

### Per-To Address Tracker

-   **ID**: `1-collection- -approvalId-uniqueID-to-recipientAddress`
-   **Behavior**: Separate tracker for each unique recipient
-   **Use Case**: Limit how much each user can receive

### Per-From Address Tracker

-   **ID**: `1-collection- -approvalId-uniqueID-from-senderAddress`
-   **Behavior**: Separate tracker for each unique sender
-   **Use Case**: Limit how much each user can send

### Per-InitiatedBy Address Tracker

-   **ID**: `1-collection- -approvalId-uniqueID-initiatedBy-initiatorAddress`
-   **Behavior**: Separate tracker for each unique initiator
-   **Use Case**: Limit how much each user can initiate

## Detailed Example

Using the approval amounts defined above, when Alice initiates a transfer of x10 from Bob:

### Two Trackers Get Incremented

**#1) Overall Tracker**

-   **ID**: `1-collection- -approvalId-uniqueID-overall-`
-   **Before**: 0/1000
-   **After**: 10/1000
-   **Behavior**: Any subsequent transfers (from Charlie, etc.) will also increment this universal tracker

**#2) Per-Initiator Tracker**

-   **ID**: `1-collection- -approvalId-uniqueID-initiatedBy-alice`
-   **Before**: 0/10
-   **After**: 10/10 (fully used)
-   **Behavior**: Only incremented when Alice initiates. Charlie's transfers use a separate tracker: `1-collection- -approvalId-uniqueID-initiatedBy-charlie`

### Amount Tracking with Balance Type

Trackers store amounts using the balance type structure. Above, we simplified it to just the amount.

```json
{
    "amounts": [
        {
            "amount": 10n,
            "badgeIds": [{ "start": 1n, "end": 1n }],
            "ownershipTimes": [{ "start": 1n, "end": 100000000000n }]
        }
    ]
}
```

**What Gets Incremented**:

-   **Amount**: The total quantity transferred
-   **Badge IDs**: Specific badge IDs that were transferred
-   **Ownership Times**: The ownership time ranges that were transferred

### Unlimited Trackers (No Increment)

Since "to" and "from" trackers are set to "0" (unlimited), no tracking occurs for these types.

## Tracker Behavior

-   **As-Needed**: Only increment trackers when necessary (unlimited = no tracking)
-   **Separate Counts**: Each tracker type maintains independent tallies
-   **Address Scoped**: Per-address trackers create unique counters per address
-   **Balance Tracking**: Increments for specific badge IDs and ownership times transferred

## Resets and ID Changes

### Changing Tracker ID

When you update `amountTrackerId` from "uniqueID" to "uniqueID2":

```
1-collection- -approvalId-uniqueID-initiatedBy-alice
↓
1-collection- -approvalId-uniqueID2-initiatedBy-alice
```

**Result**: All tracker IDs change, so all tallies start from scratch.

### Reusing Old IDs

If you later change back to "uniqueID", the starting point will be the previous tally:

-   Alice's initiatedBy tracker: 10/10 used (not 0/10)

**Important**: Never reuse tracker IDs unless you want to continue from the previous state. They are increment-only.
