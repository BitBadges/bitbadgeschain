# Cosmos Coin Wrapper Tutorial

This tutorial walks you through setting up cosmos coin wrappers to bridge BitBadges with the broader Cosmos ecosystem. Cosmos coin wrappers automatically convert tokens to fungible Cosmos coins and vice versa.

## Prerequisites

-   Understanding of [Cosmos Wrapper Paths](../concepts/cosmos-wrapper-paths.md)
-   Basic knowledge of BitBadges collections and approvals

## Step 1: Set Up Your Cosmos Denominations

First, define your cosmos coin wrapper paths. For detailed information about available options, see [Cosmos Wrapper Paths](../concepts/cosmos-wrapper-paths.md).

```typescript
const cosmosCoinWrapperPaths = [ ... ];
```

## Step 2: Generate Your Special Address

When you create a collection with cosmos coin wrapper paths, the system automatically generates a special address for each wrapper. This address acts as the bridge between tokens and cosmos coins. This will also be available on the BitBadges site if you want to go that route.

```typescript
import { generateAliasAddressForDenom } from 'bitbadgesjs-sdk';

const denom = 'utoken1';
const wrapperAddress = generateAliasAddressForDenom(denom);
console.log('Wrapper Address:', wrapperAddress);
```

## Step 3: Set Up Approvals for Wrapping/Unwrapping

The transfers still operate under the approval / transferability system. We will use the following examples from our examples section, but you can customize as you see fit. Note the need to override the wrapper address's approvals where necessary because the wrapper address is uncontrollable.

-   [Cosmos Wrapper Approval](./approvals/cosmos-wrapper-approval.md)
-   [Cosmos Unwrapper Approval](./approvals/cosmos-unwrapper-approval.md)

```typescript
const collection = {
    ...BaseCollectionDetails,
    collectionApprovalTimeline: [
        {
            timelineTimes: FullTimeRanges,
            collectionApprovals: [
                ...otherApprovals,
                wrapperApproval,
                unwrapperApproval,
            ],
        },
    ],
};
```

## Related Concepts

-   [Cosmos Wrapper Paths](../concepts/cosmos-wrapper-paths.md)
-   [Approval System](../concepts/approval-criteria/approval-system.md)
-   [Token Collections](../concepts/badge-collections.md)
