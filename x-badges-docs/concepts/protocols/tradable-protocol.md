# Tradable Protocol

The Tradable Protocol enables orderbook-style trading of badges through standardized bid and listing approvals. This protocol supports both fungible and non-fungible badge types with coin-based transactions.

## Protocol Overview

The Tradable Protocol creates a decentralized marketplace where users can:

-   **List badges for sale** at specific prices
-   **Place bids** to buy badges from other users
-   **Support both NFT and fungible badges** with flexible trading rules
-   **Execute trades** through standardized approval mechanisms

## Protocol Requirements

### Collection Standards

-   Must include "Tradable" in the `standardsTimeline` for the current time period. Note this often also goes well with "NFT" standard if your collection is NFTs.
-   Compatible with all badge types (fungible and non-fungible)
-   No restrictions on badge ID ranges or quantities

### Approval Types

#### 1. Listing Approvals (Outgoing)

Listings allow badge owners to sell their badges for coins.

**Requirements:**

-   Single transfer time range
-   Exactly one coin transfer with one coin denomination
-   Coin recipient must be the badge owner (`to` equals `fromListId`)
-   No address overrides (`overrideFromWithApproverAddress` and `overrideToWithInitiator` must be false)
-   Specific badge IDs (no `allowOverrideWithAnyValidBadge`)
-   Full ownership times
-   No Merkle challenges or prerequisite badges
-   `overallMaxNumTransfers` > 0
-   Typically, you want the denomination to match the collection's preferred denomination.

#### 2. Bid Approvals (Incoming)

Bids allow users to offer coins to purchase badges from others.

**Requirements:**

-   Single transfer time range
-   Exactly one coin transfer with one coin denomination
-   Coins come from bidder (`overrideFromWithApproverAddress` must be true)
-   Coins go to badge owner (`overrideToWithInitiator` must be true)
-   Specific badge IDs (no `allowOverrideWithAnyValidBadge` unless collection bid)
-   Full ownership times
-   No Merkle challenges or prerequisite badges
-   `overallMaxNumTransfers` > 0

#### 3. Collection Bids (Special Case)

Collection bids allow users to bid on any badge within a collection.

**Additional Requirements:**

-   Must have `allowOverrideWithAnyValidBadge` set to true
-   All other bid requirements apply

## Validation Functions

### General Orderbook Validation

**API Documentation:** [isOrderbookBidOrListingApproval](https://bitbadges.github.io/bitbadgesjs/functions/isOrderbookBidOrListingApproval.html)

```typescript
export const isOrderbookBidOrListingApproval = (
    approval: iCollectionApproval<bigint>,
    approvalLevel: 'incoming' | 'outgoing'
) => {
    return isBidOrListingApproval(approval, approvalLevel, {
        isFungibleCheck: true,
        fungibleOrNonFungibleAllowed: true,
    });
};
```

### Core Bid/Listing Validation

**API Documentation:** [isBidOrListingApproval](https://bitbadges.github.io/bitbadgesjs/functions/isBidOrListingApproval.html)

```typescript
export const isBidOrListingApproval = (
    approval: iCollectionApproval<bigint>,
    approvalLevel: 'incoming' | 'outgoing',
    options?: {
        isFungibleCheck?: boolean;
        fungibleOrNonFungibleAllowed?: boolean;
        isCollectionBid?: boolean;
    }
) => {
    const approvalCriteria = approval.approvalCriteria;
    if (approvalCriteria?.coinTransfers?.length !== 1) {
        return false;
    }

    if (approval.transferTimes.length !== 1) {
        return false;
    }

    const coinTransfer = approvalCriteria.coinTransfers[0];
    if (coinTransfer.coins.length !== 1) {
        return false;
    }

    // Validate address overrides for incoming approvals (bids)
    if (
        approvalLevel === 'incoming' &&
        !coinTransfer.overrideFromWithApproverAddress
    ) {
        return false;
    }

    if (approvalLevel === 'incoming' && !coinTransfer.overrideToWithInitiator) {
        return false;
    }

    // Validate address overrides for outgoing approvals (listings)
    if (
        approvalLevel === 'outgoing' &&
        coinTransfer.overrideFromWithApproverAddress
    ) {
        return false;
    }

    if (approvalLevel === 'outgoing' && coinTransfer.overrideToWithInitiator) {
        return false;
    }

    // For listings, recipient must be the approving user
    const to = coinTransfer.to;
    if (approvalLevel === 'outgoing' && to !== approval.fromListId) {
        return false;
    }

    const incrementedBalances =
        approvalCriteria.predeterminedBalances?.incrementedBalances;
    if (!incrementedBalances) {
        return false;
    }

    if (incrementedBalances.startBalances.length !== 1) {
        return false;
    }

    // Collection bids can accept any valid badge
    if (options?.isCollectionBid) {
        if (!incrementedBalances.allowOverrideWithAnyValidBadge) {
            return false;
        }
    } else {
        const allBadgeIds = UintRangeArray.From(
            incrementedBalances.startBalances[0].badgeIds
        )
            .sortAndMerge()
            .convert(BigInt);
        if (allBadgeIds.length !== 1 || allBadgeIds.size() !== 1n) {
            return false;
        }

        if (incrementedBalances.allowOverrideWithAnyValidBadge) {
            return false;
        }
    }

    const amount = incrementedBalances.startBalances[0].amount;
    const toCheckAmountOne =
        !options ||
        (!options.isFungibleCheck && !options.fungibleOrNonFungibleAllowed);
    if (toCheckAmountOne) {
        if (amount !== 1n) {
            return false;
        }
    }

    if (
        !UintRangeArray.From(
            incrementedBalances.startBalances[0].ownershipTimes
        ).isFull()
    ) {
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

    if (approvalCriteria.requireFromEqualsInitiatedBy) {
        return false;
    }

    if (approvalCriteria.requireToEqualsInitiatedBy) {
        return false;
    }

    if (approvalCriteria.overridesToIncomingApprovals) {
        return false;
    }

    if (approvalCriteria.merkleChallenges?.length) {
        return false;
    }

    if (approvalCriteria.mustOwnBadges?.length) {
        return false;
    }

    if (
        (approvalCriteria.maxNumTransfers?.overallMaxNumTransfers ?? 0n) === 0n
    ) {
        return false;
    }

    return true;
};
```

### Collection Bid Validation

**API Documentation:** [isCollectionBid](https://bitbadges.github.io/bitbadgesjs/functions/isCollectionBid.html)

```typescript
export const isCollectionBid = (approval: iCollectionApproval<bigint>) => {
    return isBidOrListingApproval(approval, 'incoming', {
        isCollectionBid: true,
    });
};
```

## Implementation Example

For a complete implementation example, see the [Tradable NFT Collection Example](../../examples/txs/msgcreatecollection/tradable-nft-collection.md).
