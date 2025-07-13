# Overview

In the transferability page, we mainly talked about how we match and check an approval. Here, we will talk about the **approvalCriteria.**

These are the additional options or restrictions you can set which decide whether a transfer is approved or not (e.g. how much can be transferred? how many times? etc). To be approved, it must satisfy all the options / restrictions set (everything talked about in the last page AND the **approvalCriteria**).

```typescript
export interface iApprovalCriteria<T extends NumberType> {
    /** The $BADGE transfers to be executed upon every approval. */
    coinTransfers?: iCoinTransfer<T>[];
    /** The list of merkle challenges that need valid proofs to be approved. */
    merkleChallenges?: iMerkleChallenge<T>[];
    /** The predetermined balances for each transfer. */
    predeterminedBalances?: iPredeterminedBalances<T>;
    /** The maximum approved amounts for this approval. */
    approvalAmounts?: iApprovalAmounts<T>;
    /** The max num transfers for this approval. */
    maxNumTransfers?: iMaxNumTransfers<T>;
    /** Whether the to address must equal the initiatedBy address. */
    requireToEqualsInitiatedBy?: boolean;
    /** Whether the from address must equal the initiatedBy address. */
    requireFromEqualsInitiatedBy?: boolean;
    /** Whether the to address must not equal the initiatedBy address. */
    requireToDoesNotEqualInitiatedBy?: boolean;
    /** Whether the from address must not equal the initiatedBy address. */
    requireFromDoesNotEqualInitiatedBy?: boolean;
    /** Whether this approval overrides the from address's approved outgoing transfers. */
    overridesFromOutgoingApprovals?: boolean;
    /** Whether this approval overrides the to address's approved incoming transfers. */
    overridesToIncomingApprovals?: boolean;
}
```

**Extending Functionality**

These are just the native options provided in the interface for convenience and consistency, but you can always implement your own approvals via custom logic with a CosmWASM smart contract.

Or, you can leverage the native features for custom implementations:

-   Use Merkle challenges as a commit-reveal claim code mechanism. This is what we do for BitBadges claims. Then, you can assign claim codes how you wish.&#x20;
-   Create custom, general-purpose zero-knowledge circuits that need to be satisfied. Implement your custom logic within in the circuit (if it lends itself to this approach).

#### Tracker IDs

The approval interface utilizes trackers behind the scenes for certain fields which are identified by IDs (**amountTrackerId, challengeTrackerId)**. Trackers are increment only and immutable in storage and referenced by an ID consisting of **approvalId-trackerId** plus other identifying details. Because the **approvalId** is in there, this enforces that all trackers are scoped to a specific approvalId. However, since they are ID based and increment only, it is important to be careful to not use previous IDs that have state. See best practices below.

**Best Practices - Creating / Updating / Deleting**

A really important part in creating / editing / deleting approvals is to keep track of the trackers' state. You do not want to use a tracker with prior history when creating an approval which should start from scratch. This messes up the expected behavior of the approval. You can't just simply delete the tracker because it is increment only and immutable.

To combat this, we recommend the following best practices:

-   When creating an approval for the first time (which is expected to start completely from scratch), always use unique, unused IDs for **approvalId**, **amountTrackerId** and **challengeTrackerId**. If these are unique, the tracker IDs will also be unique which means they will have no prior history.
-   When you want entirely new state (whole approval), change the **approvalId** to something unique. This makes all tracker IDs different. This is the recommended option (if applicable) because you know for certain that all trackers have been reset.
-   When you want just an individual tracker's state to be reset, you need to change the individual tracker ID to a aunique unused vaue.

#### Scoped vs Cross-Approval Logic

All approvals' state is expected to be scoped, but sometimes, you may want to implement cross-approval logic (e.g. do not allow double dipping between two approvals). This, unfortunately, is out of scope for the native interface at the moment. Consider workarounds and careful design decisions.&#x20;

However, if you do need more advanced functionality, you may have to go up a level and use CosmWASM or alternative solutions. If you run into this problem, let us know, and we can recommend what to do.&#x20;
