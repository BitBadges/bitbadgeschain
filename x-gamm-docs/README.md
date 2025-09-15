# 📚 Overview

This directory contains comprehensive developer documentation for the BitBadges blockchain's `x/gamm` module.

This module was forked from Osmosis's `x/gamm` module with several key modifications to support BitBadges tokens and enhance compatibility.

## Key Differences from Osmosis

### 1. Interface Revamps

-   Removed `smoothWeightChangeParams` and other unused parameters
-   Updated certain type definitions for better compatibility with our codebase
-   Streamlined interfaces for improved performance
-   Remove unneeded logic like stableswap pools and governance proposal handling
-   Removed future pool governor functionality
-   Removed pool creation fee requirements

### 2. Badge Token Integration

The main difference is in the badge token handling system. With every attempted transfer of BitBadges tokens, the system:

-   **Wrapping Conversion**: Uses `cosmosCoinWrapperPaths` defined by the collection to convert badge tokens to `x/bank` denominations (using the `path.balances` array) at a set conversion rate
-   **1:1 Minting**: Behind the scenes, mints/burns 1:1 at the conversion rate into an `x/bank` denomination before / after each transfer to/from a pool
-   **Pool Compatibility**: Ensures seamless integration with existing pool infrastructure
-   **Automatic Conversion**: Handles badge token transfers to/from pools automatically

#### Conversion Example

Here's how the badge token conversion works:

**Badge Token**: `badgeslp:21:utoken`

-   Collection ID: `21`
-   Base Denom: `utoken`
-   Wrapper Path Balances Conversion Rate: `[{ amount: 1n, badgeIds: [{ start: 1n, end: 1n }], ownershipTimes: UintRangeArray.FullRanges() }]`

**Conversion Process**:

```
1 badgeslp:21:utoken = [{ amount: 1n, badgeIds: [{ start: 1n, end: 1n }], ownershipTimes: UintRangeArray.FullRanges() }]

2 badgeslp:21:utoken = [{ amount: 2n, badgeIds: [{ start: 1n, end: 1n }], ownershipTimes: UintRangeArray.FullRanges() }]
```

**Visual Flow**:

```
User wants to add 5 badgeslp:21:utoken to a pool in liquidity
    ↓
System reads cosmosCoinWrapperPaths for collection 21
    ↓
Finds path with denom "utoken" and balances array (using the wrapper path's balances array)
    ↓
Transfers the x/badges balances to the pool address with conversions applied
    ↓
Pool receives the x/badges balances and mints the corresponding badgeslp:21:utoken balances which are native x/bank denominations
    ↓
The user now has the equivalent of badgeslp:21:utoken balances in the pool (all else in x/gamm is kept the same as standard in Osmosis)
```

For the reverse, when a pool wants to send tokens to the user, it will burn the native x/bank denominations and transfer the x/badges balances back to the user address.

## Table of Contents

1. [Introduction](./introduction.md) - Overview and key concepts
2. [Messages](./messages/) - Transaction messages and handlers

## Message Reference

### Core Operations

-   [MsgCreateBalancerPool](./messages/msg-create-balancer-pool.md) - Create new balancer pool
-   [MsgJoinPool](./messages/msg-join-pool.md) - Join existing pool with liquidity
-   [MsgSwapExactAmountIn](./messages/msg-swap-exact-amount-in.md) - Swap exact amount of tokens in
-   [MsgExitPool](./messages/msg-exit-pool.md) - Exit pool and receive tokens

## Query Reference

For all GAMM queries, please refer to the [BitBadges LCD API](https://lcd.bitbadges.io/).

The LCD provides comprehensive query endpoints for:

-   Pool information and statistics
-   Trading data and spot prices
-   Module parameters
-   And more

All queries follow the standard Cosmos SDK query patterns and can be accessed via REST API or gRPC.

## Quick Links

-   [BitBadges Chain Repository](https://github.com/bitbadges/bitbadgeschain)
-   [BitBadges Documentation](https://docs.bitbadges.io)
-   [Proto Definitions](https://github.com/bitbadges/bitbadgeschain/tree/master/proto/gamm)

## Documentation Style

This documentation follows the [Cosmos SDK module documentation standards](https://docs.cosmos.network/main/building-modules/README) and is designed for developers building on or integrating with the BitBadges blockchain.
