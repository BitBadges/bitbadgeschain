# Predetermined Balances

Predetermined balances are a new way of having fine-grained control over the amounts that are approved with each transfer. In a typical tally-based system where you approve X amount to be transferred, you have no control over the combination of amounts that will add up to X. For example, if you approve x100, you can't control whether the transfers are \[x1, x1, x98] or \[x100] or another combination.

Predetermined balances let you explicitly define the amounts that must be transferred and the order of the transfers. For example, you can enforce x1 of badge ID 1 has to be transferred before x1 of badge ID 2, and so on.

Although this can be used in tandem with approval amounts, either one or the other is usually used because they both specify amount restrictions.

**TLDR; The transfer will fail if the balances are not EXACTLY as defined in the predetermined balances.**

```typescript
export interface PredeterminedBalances<T extends NumberType> {
    manualBalances: ManualBalances<T>[];
    incrementedBalances: IncrementedBalances<T>;
    orderCalculationMethod: PredeterminedOrderCalculationMethod;
}
```

## **Defining Balances**

There are two ways to define the balances. Both can not be used together.

-   **Manual Balances:** Simply define an array of balances manually. Each element corresponds to a different set of balances for a unique transfer.
-   ```json
    "manualBalances": [
      {
        "amount": "1",
        "badgeIds": [
          {
            "start": "1",
            "end": "1"
          }
        ],
        "ownershipTimes": [
          {
            "start": "1691978400000",
            "end": "1723514400000"
          }
        ]
      },
      {...},
      {...},
    ]
    ```
-   **Incremented Balances:** Define starting balances and then define how much to increment or calculate the IDs and times by after each transfer. There are different approaches here (incompatible with each other).
    -   Increments: You can enforce x1 of badge ID 1 has to be transferred before x1 of badge ID 2, and so on. This is typically used for minting badges. You can also customize the ownership times to increment by a certain amount. Or, have them dynamically overriden to be the current time + a interval length (now + 1 month, now + 1 year, etc).
    -   Duration From Timestamp: If enabled, this will dynamically calculate the ownership times from a timestamp (default: transfer time) + a set duration of time. All ownership times will be overwritten. If the override timestamp is allowed, users can specify a custom timestamp to start from in MsgTransferBadges precalculationOptions.
    -   Recurring Ownership Times: Recurring ownership times are similar to the above, but they define set intervals + charge periods that are approved. For example, you could approve the ownership times for the 1st to the 30th of the month which repeats indefinitely. The charge period is how long before the next interval starts, the approval can be used. For example, allow this approval to be charged up to 7 days in advance of the next interval.
    -   Allowing Badge Override: If enabled, users can specify a custom badge ID in MsgTransferBadges precalculationOptions. This will override all the badge IDs in the starting balances to this specified badge ID. This is useful for collection offers. The specified badge IDs must only be a single badge ID and must be a valid badge ID in the collection.
-   ```json
    "incrementedBalances": {
      "startBalances": [
        {
          "amount": "1",
          "badgeIds": [
            {
              "start": "1",
              "end": "1"
            }
          ],
          "ownershipTimes": [
            {
              "start": "1691978400000",
              "end": "1723514400000"
            }
          ]
        }
      ],
      "incrementBadgeIdsBy": "1",
      "incrementOwnershipTimesBy": "0",
      "durationFromTimestamp": "0", // UNIX milliseconds
      "allowOverrideTimestamp": false,
      "allowOverrideWithAnyValidBadge": false,
      "recurringOwnershipTimes": {
        "startTime": "0",
        "intervalLength": "0",
        "chargePeriodLength": "0"
      }
    }
    ```

## **Precalculating Balances**

Predetermined balances can quickly change, such as in between the time a transaction is broadcasted and confirmed. For example, other users' mints get processed, and thus, the badge IDs one should receive changes. This creates a problem because you can't manually specify balances because that results in race conditions and failed transfers / claims.

To combat this, when initiating a transfer, we allow you to specify **precalculateBalancesFromApproval** (in [MsgTransferBadges](../../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md)). Here, you define which **approvalId** you want to precalculate from, and at execution time, we calculate what the predetermined balances are and override the requested balances to transfer with them. Note this is the unique **approvalId** of the approval, not the tracker ID. Additional override options can be specified in the **precalculationOptions** field as well.

<pre class="language-typescript"><code class="lang-typescript"><strong>precalculateBalancesFromApproval: {
</strong>    approvalId: string;
    approvalLevel: string; //"collection" | "incoming" | "outgoing"
    approverAddress: string; //"" if collection-level
    version: string; //"1"
}
</code></pre>

## **Defining Order of Transfers**

Which balances to assign for a transfer is calculated by a specified order calculation method.

For manual balances, we want to determine which element index of the array is transferred (e.g. order number = 0 means the balances of manualBalances\[0] will be transferred). For incremented balances, this corresponds to how many times we should increment (e.g. order number = 5 means apply the increments to the starting balances five times).

There are five calculation methods to determine the order method.

### Defining Order by Number of Transfers

We either use a running tally of the number of transfers to calculate the order number (no previous transfers = order number 0, one previous transfer = order number 1, and so on). This can be done on an overall or per to/from/initiatedBy address basis and is incremented using an approval tracker as explained in [Max Number of Transfers](predetermined-balances.md#max-number-of-transfers).

IMPORTANT: Note the number of transfers is tracked using the same tracker as used within **maxNumTransfers**. Trackers are increment only, immutable, and incremented on an as-needed basis. Be mindful of this. If the tracker has prior history (potentially because **maxNumTransfers** was set), the order numbers will be calculated according to the prior history of this tracker. The opposite is also true. If you are tracking transfers here for predetermined balances, the **maxNumTransfers** restrictions will be calculated according to the tracker's history. Consider this when editing / creating approvals. You do not want to use a tracker that has prior history when you expect it to start from scratch.

### Reserved Order

We also support using the leaf index for the defined Merkle challenge proof (see [Merkle Challenges](predetermined-balances.md#merkle-challenges)) to calculate the order number (e.g. leftmost leaf on expected leaf layer will correspond to order number 0, next leaf will be order number 1, and so on). The leftmost leaf means the leftmost leaf of the **expectedProofLength** layer. The challenge we will use is the one with the corresponding **challengeTrackerId**.

This is used to reserve specific badges for specific users / claim codes. For example, reserve the badges corresponding to order number 10 (leaf number 10) for address xyz.eth.

```typescript
export interface PredeterminedOrderCalculationMethod {
    useOverallNumTransfers: boolean;
    usePerToAddressNumTransfers: boolean;
    usePerFromAddressNumTransfers: boolean;
    usePerInitiatedByAddressNumTransfers: boolean;
    useMerkleChallengeLeafIndex: boolean;
    challengeTrackerId: string;
}
```

**Overlap / Out of Bounds**

In the base approval interface, we specify the bounds for the approval ("Alice" can transfer the IDs 1-10 for Mon-Fri to "Bob" initiated by "Alice"). Typically, the precalculated balances should be completely within these bounds. However, the order number may eventually correspond to balances that have no overlap with these bounds or partially overlap. For example, if you approve x1 of ID 1, then x1 of ID 2 and so on up to x1 of ID 10000, eventually, the order number will be 10001 which corresponds to balances that are out of bounds.

If it is completely out of bounds (e.g. order number = 101 but approved badgeIds 1-100 with increments of 1), this is practically ignored. This is because if you try and transfer badge ID 101, it will never match to the current approval.

You should try and design your approvals for no partial overlaps. But, in rare cases, this may occur (some in bounds and some out of bounds). In this case, the overall transfer balances still must be **exactly** as defined (in bounds + out of bounds); however, we only approve the in bounds ones for the current approval. The out of bounds ones must be approved by a separate approval.
