# Max Number of Transfers

Pre-Readings (IMPORTANT): [Approval Trackers](approval-trackers.md)

```typescript
export interface ApprovalCriteria<T extends NumberType> {
  ...
  maxNumTransfers?: MaxNumTransfers<T>;
  ...
}
```

Similar to approval amounts, you can also specify the maximum number of transfers that can occur on an overall or per sender/recipient/initiatedBy basis. If tracked, the corresponding approval tracker will be incremented by 1 each transfer.&#x20;

"0" means unlimited allowed and not tracked. "N" means max N transfers allowed.

**Example**

Let's say we have the ID `1-collection- -approvalId-uniqueID-initiatedBy-alice` and the defined values below:

<pre class="language-json"><code class="lang-json">"maxNumTransfers": {
<strong>    "overallMaxNumTransfers": "0",
</strong>    "perFromAddressMaxNumTransfers": "0",
    "perToAddressMaxNumTransfers": "0",
    "perInitiatedByAddressMaxNumTransfers": "1",
    "amountTrackerId": "uniqueID",
    "resetTimeIntervals": {
      "startTime": "0",
      "intervalLength": "0"
    }
}
</code></pre>

The first transfer initiated by Alice would increment the approval tracker ID`1-collection- -approvalId-uniqueID-initiatedBy-alice` to 1/1 transfers used. Alice can no longer initiate another transfer.

However, Bob can still transfer because his ID is `1-collection- -approvalId-uniqueID-initiatedBy-bob` which is a different tracker.

**As-Needed Basis**

We track on an as-needed basis, meaning if we do not have requirements that use the number of transfers, we will not increment / track.

Edge Case: In [Predetermined Balances](max-number-of-transfers.md#predetermined-balances), you may need the number of transfers for determining the balances to assign to each transfer (e.g. transfer #10 -> badge ID 10). In this case, we do need to track the number of transfers. This is all facilitated via the same tracker, so even if you have "0" or unlimited set for the corresponding value in **maxNumTransfers**, the tracker may be incremented behind the scenes. Consider this when editing / creating approvals. You do not want to use a tracker that has prior history when you expect it to start from scratch.
