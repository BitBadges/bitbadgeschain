# Quest Protocol

The Quest Protocol is a standardized way to create quest-based badge collections that reward users with badges and coins for completing cryptographically verified tasks. Quest collections are designed for achievement-based systems where users complete quests (tasks) and provide cryptographic proofs (Merkle proofs) to claim badges along with coin incentives.

## Protocol Requirements

### Collection Standards

-   Must include "Quests" in the `standardsTimeline` for the current time period
-   Must have exactly one valid badge ID: `{"start": "1", "end": "1"}`

### Quest Approval Requirements

-   **From List**: Must be "Mint" (minting from the mint address)
-   **Merkle Challenge**: Must have exactly one Merkle challenge with:
    -   `maxUsesPerLeaf`: 1 (one use per proof)
    -   `useCreatorAddressAsLeaf`: false (custom proof verification)
    -   Valid Merkle root hash for proof verification
-   **Coin Transfers**: Must have exactly one coin transfer with:
    -   Exactly one coin denomination and amount
    -   `overrideFromWithApproverAddress`: true (coins come from collection creator)
    -   `overrideToWithInitiator`: true (coins go to the quest completer)
-   **Max Transfers**: Must have `overallMaxNumTransfers` > 0
-   **Predetermined Balances**: Must have:
    -   Exactly one `startBalance` with amount 1 for badge ID 1
    -   `incrementBadgeIdsBy`: 0 (no badge ID incrementing)
    -   `incrementOwnershipTimesBy`: 0 (no time incrementing)
    -   `durationFromTimestamp`: 0 (no time-based duration)
    -   `allowOverrideTimestamp`: false (no timestamp overrides)
    -   All `recurringOwnershipTimes` fields set to 0 (no recurring)
-   **Additional Constraints**:
    -   `mustOwnBadges`: empty (no prerequisite badges)
    -   `requireToEqualsInitiatedBy`: false (no address matching required)

## Validation Functions

### Collection Validation

**API Documentation:** [doesCollectionFollowQuestProtocol](https://bitbadges.github.io/bitbadgesjs/functions/doesCollectionFollowQuestProtocol.html)

```typescript
export const doesCollectionFollowQuestProtocol = (
    collection?: Readonly<iCollectionDoc<bigint>>
) => {
    if (!collection) {
        return false;
    }

    // Check if "Quests" standard is active for current time
    let found = false;
    for (const standard of collection.standardsTimeline) {
        const isCurrentTime = UintRangeArray.From(
            standard.timelineTimes
        ).searchIfExists(BigInt(Date.now()));
        if (!isCurrentTime) {
            continue;
        }

        if (!standard.standards.includes('Quests')) {
            continue;
        }

        found = true;
    }

    if (!found) {
        return false;
    }

    // Assert valid badge IDs are only 1n-1n
    const badgeIds = UintRangeArray.From(collection.validBadgeIds)
        .sortAndMerge()
        .convert(BigInt);
    if (badgeIds.length !== 1 || badgeIds.size() !== 1n) {
        return false;
    }

    if (badgeIds[0].start !== 1n || badgeIds[0].end !== 1n) {
        return false;
    }

    return true;
};
```

### Approval Validation

**API Documentation:** [isQuestApproval](https://bitbadges.github.io/bitbadgesjs/functions/isQuestApproval.html)

```typescript
export const isQuestApproval = (approval: iCollectionApproval<bigint>) => {
    const approvalCriteria = approval.approvalCriteria;
    if (!approvalCriteria?.coinTransfers) {
        return false;
    }

    // Must be minting approval
    if (approval.fromListId !== 'Mint') {
        return false;
    }

    // Must have exactly one Merkle challenge
    if (
        !approvalCriteria.merkleChallenges ||
        approvalCriteria.merkleChallenges.length !== 1
    ) {
        return false;
    }

    let merkleChallenge = approvalCriteria.merkleChallenges?.[0];
    if (merkleChallenge.maxUsesPerLeaf !== 1n) {
        return false;
    }

    // Must not require owning other badges
    if (approvalCriteria.mustOwnBadges?.length) {
        return false;
    }

    // Must not use creator address as leaf
    if (merkleChallenge.useCreatorAddressAsLeaf) {
        return false;
    }

    // Must have max transfer limit
    const maxNumTransfers =
        approvalCriteria.maxNumTransfers?.overallMaxNumTransfers;
    if (!maxNumTransfers) {
        return false;
    }

    if (maxNumTransfers <= 0n) {
        return false;
    }

    // Must have exactly one coin transfer
    if (approvalCriteria.coinTransfers.length !== 1) {
        return false;
    }

    // Validate coin transfer configuration
    for (const coinTransfer of approvalCriteria.coinTransfers) {
        if (coinTransfer.coins.length !== 1) {
            return false;
        }

        if (
            !coinTransfer.overrideFromWithApproverAddress ||
            !coinTransfer.overrideToWithInitiator
        ) {
            return false;
        }
    }

    // Validate predetermined balances
    const incrementedBalances =
        approvalCriteria.predeterminedBalances?.incrementedBalances;
    if (!incrementedBalances) {
        return false;
    }

    if (incrementedBalances.startBalances.length !== 1) {
        return false;
    }

    const allBadgeIds = UintRangeArray.From(
        incrementedBalances.startBalances[0].badgeIds
    )
        .sortAndMerge()
        .convert(BigInt);
    if (allBadgeIds.length !== 1 || allBadgeIds.size() !== 1n) {
        return false;
    }

    if (allBadgeIds[0].start !== 1n || allBadgeIds[0].end !== 1n) {
        return false;
    }

    const amount = incrementedBalances.startBalances[0].amount;
    if (amount !== 1n) {
        return false;
    }

    if (incrementedBalances.incrementBadgeIdsBy !== 0n) {
        return false;
    }

    if (incrementedBalances.incrementOwnershipTimesBy !== 0n) {
        return false;
    }

    if (incrementedBalances.durationFromTimestamp !== 0n) {
        return false;
    }

    // Needs this to be false for the subscription faucet to work
    if (incrementedBalances.allowOverrideTimestamp) {
        return false;
    }

    if (incrementedBalances.recurringOwnershipTimes.startTime !== 0n) {
        return false;
    }

    if (incrementedBalances.recurringOwnershipTimes.intervalLength !== 0n) {
        return false;
    }

    if (incrementedBalances.recurringOwnershipTimes.chargePeriodLength !== 0n) {
        return false;
    }

    if (approvalCriteria.requireToEqualsInitiatedBy) {
        return false;
    }

    return true;
};
```

## Implementation Example

For a complete implementation example, see the [Quest Badge Collection Example](../../examples/txs/msgcreatecollection/quest-badge-collection.md).
