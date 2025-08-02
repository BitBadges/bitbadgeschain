# Approval Criteria

Additional restrictions and conditions that determine whether a transfer is approved beyond the basic approval matching.

## Interface

```typescript
export interface iApprovalCriteria<T extends NumberType> {
    /** The $BADGE transfers to be executed upon every approval. */
    coinTransfers?: iCoinTransfer<T>[];
    /** The list of merkle challenges that need valid proofs to be approved. */
    merkleChallenges?: iMerkleChallenge<T>[];
    /** The list of must own badges that need valid proofs to be approved. */
    mustOwnBadges?: iMustOwnBadge<T>[];
    /** The predetermined balances for each transfer. These allow approvals to use predetermined balance amounts rather than an incrementing tally system. */
    predeterminedBalances?: iPredeterminedBalances<T>;
    /** The maximum approved amounts for this approval. */
    approvalAmounts?: iApprovalAmounts<T>;
    /** The max num transfers for this approval. */
    maxNumTransfers?: iMaxNumTransfers<T>;
    /** Whether the approval should be deleted after one use. */
    autoDeletionOptions?: iAutoDeletionOptions;
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
    /** The royalties to apply to the transfer. */
    userRoyalties?: iUserRoyalties<T>;
    /** The list of dynamic store challenges that the initiator must pass for approval. */
    dynamicStoreChallenges?: iDynamicStoreChallenge<T>[];
    /** The list of ETH signature challenges that require valid Ethereum signatures for approval. */
    ethSignatureChallenges?: iETHSignatureChallenge<T>[];
}
```

## Core Components

-   **[Approval Trackers](approval-trackers.md)** - Tracking transfer amounts and counts
-   **[Tallied Approval Amounts](tallied-approval-amounts.md)** - Amount limits and thresholds
-   **[Max Number of Transfers](max-number-of-transfers.md)** - Transfer count limits
-   **[Predetermined Balances](predetermined-balances.md)** - Exact balance requirements
-   **[Merkle Challenges](merkle-challenges.md)** - Cryptographic proof requirements
-   **[Dynamic Store Challenges](dynamic-store-challenges.md)** - On-chain numeric checks
-   **[ETH Signature Challenges](eth-signature-challenges.md)** - Ethereum signature requirements
-   **[Badge Ownership](badge-ownership.md)** - Required badge holdings
-   **[$BADGE Transfers](usdbadge-transfers.md)** - Automatic token transfers
-   **[Overrides](overrides.md)** - Bypassing user-level approvals
-   **[Requires](requires.md)** - Address relationship restrictions
-   **[Auto-Deletion Options](auto-deletion-options.md)** - Automatic approval cleanup
-   **[User Royalties](user-royalties.md)** - Percentage-based transfer fees

## Key Concepts

### Tracker IDs

Trackers use IDs with format: `approvalId-trackerId` plus identifying details. All trackers are scoped to a specific `approvalId`.

**Important**: Trackers are increment-only and immutable. Never reuse tracker IDs with prior history.

### Best Practices - Creating / Updating / Deleting

Trackers are increment-only and immutable. Never reuse tracker IDs with prior history when creating approvals that should start from scratch.

### Extending Functionality

For custom logic beyond native options:

-   Use CosmWASM smart contracts
-   Leverage Merkle challenges for commit-reveal mechanisms

### Cross-Approval Logic

Native interface doesn't support cross-approval logic (e.g., preventing double-dipping between approvals). Consider:

-   Workarounds and careful design
-   CosmWASM for advanced functionality
-   Contact us for recommendations
