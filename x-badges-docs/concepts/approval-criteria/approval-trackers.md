# Approval Trackers

```typescript
export interface iApprovalAmounts<T extends NumberType> {
  ...
  
  /** The ID of the approval tracker. This is the key used to track tallies. */
  amountTrackerId: string;
}
```

```typescript
export interface iMaxNumTransfers<T extends NumberType> {
  ...
  
  /** The ID of the approval tracker. This is the key used to track tallies. */
  amountTrackerId: string;
}
```

## Approval Trackers

Approval or amount trackers track how many badges have been transferred and how many transfers have occurred. This is done via an incrementing tally system with a threshold.&#x20;

Take the following approval tracker

1. You are approved for x10 of badge IDs 1-10. For simplicity, lets say the tracker ID is "xyz".
2. You transfer x5 of badge IDs 1-10 -> "xyz" tally goes from x0/10 -> x5/10
3. You transfer another x5 -> "xyz" tally goes to x10/10
4. You transfer another x1 -> exceed threshold so transfer fails and will also fail for all subsequent approvals that match to "xyz".
5. However, if subsequent transfers match to tracker "abc", the tally starts from zero again because it is a different tracker w/o any history.

### **How are approval trackers identified?**

Above, we used "xyz" for simplicity, but the identifier of each approval tracker actually consists of **amountTrackerId** along with other identifying details.

Note that if multiple approvals specify the same **amountTrackerId,** the SAME tracker will be incremented when DIFFERENT approvals are used. This is because the tracker identifier will be the same and thus increment the same tracker. However, all tracker IDs specify the **approvalId** which must be unique. Thus, all trackers are only ever scoped to a single aproval.

```
ID: collectionId-approvalLevel-approverAddress-approvalId-amountTrackerId-trackerType-approvedAddress
```

```typescript
export interface ApprovalTrackerIdDetails<T extends NumberType> {
  collectionId: T
  approvalLevel: "collection" | "incoming" | "outgoing" | ""
  approvalId: string
  approverAddress: string
  amountTrackerId: string
  trackerType: "overall" | "to" | "from" | "initiatedBy" | ""
  approvedAddress: string
}
```

The **trackerType** corresponds to what type of tracker it is. For example, should we increment every time this approval is used? per unique recipient? sender? initiator?

If "overall", this is applicable to any transfer and will increment everytime the approval is used. This creates a single universal tally. **approvedAddress** will be empty.&#x20;

If "to", "from", or "initiatedBy", the **approvedAddress** is the sender, recipient, or initiator of the transfer, respectively. Note since **approvedAddress** and **trackerType** are part of the approval tracker's identifier, this creates unique individual tallies (trackers) per address.

For example, these correspond to different trackers because the **approvedAddress** is different. Thus, Alice's transfers will be tracked separately from Bob's.

`1-collection- -approvalId-uniqueID-initiatedBy-alice`

`1-collection- -approvalId-uniqueID-initiatedBy-bob`

**Handling Multiple Trackers**

Trackers are ID-based, and thus, multiple trackers can be created. Take note of what makes up the ID. The collection ID, approval level, approver address, and more are all considered. If one changes or is different, the whole ID is different and will correspond to a new tracker.

**Increment Only**

Trackers are increment only and immutable in storage. To start an approval tally from scratch, you will need to map the approval to a new unused ID. This can be done simply by editing **amountTrackerId** (because this changes the whole ID) or restructuring to change one of the other fields that make up the overall ID.

IMPORTANT: Because of the immutable nature, be careful to not revert to a previously used ID unintentionally because the starting point will be the previous tally (not starting from scratch).

### **What is tracked?**

```typescript
export interface ApprovalTrackerInfoBase<T extends NumberType> extends ApprovalTrackerIdDetails<T> {
  numTransfers: T;
  amounts: Balance<T>[];
  lastUpdatedAt: UnixMilliTimestamp<T>;
}
```

Each transfer that maps to the tracker increments **numTransfers** by 1, and each badge transferred increments the **amounts** in the interface (if tracked).

Example:

`ID: 1-collection- -uniqueID-initiatedBy-alice`

```json
{
    "numTransfers": 10,
    "amounts": [{ 
        "amount": 10n, 
        "badgeIds": [{ start: 1n, end: 1n }], 
        "ownershipTimes":  [{ start: 1n, end: 100000000000n }], 
    }]
}
```

**As-Needed Basis**

We increment on an as-needed basis. Meaning, if there is no need to increment the tally (unlimited limit and/or not restrictions), we **do not increment** for efficiency purposes. For example, if we only have requirements for **numTransfers** but do not need the **amounts**, we do not increment the amounts.

**Different Tracker IDs - Amounts vs Transfers**

It is possible to have different tracker IDs for the number of transfers and amounts (as seen at the top of the page). However, typically, these will be the same for simplicity.

### **Resets**

For both tracker types, we also allow occasional resets to zero. If it is the first update of the interval, we will reset all tracker progress. This is useful for recurring subscriptions (one transfer per month), for example.

If they are left as 0, there is no reset.&#x20;

Whether an update is needed or not is calculated using the lastUpdatedAt field.

```typescript
/**
 * @category Interfaces
 */
export interface iResetTimeIntervals<T extends NumberType> {
  /** The start time of the first interval. */
  startTime: T;
  /** The length of the interval. */
  intervalLength: T;
}

```
