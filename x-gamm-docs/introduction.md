# Introduction

The `x/gamm` module implements the Generalized Automated Market Maker (GAMM) functionality for the BitBadges blockchain. This module was forked from Osmosis's `x/gamm` module and provides the core infrastructure for decentralized exchange (DEX) operations, including liquidity pool creation, token swaps, and yield farming.

## Fork from Osmosis

This module maintains compatibility with Osmosis's GAMM functionality while adding specialized support for BitBadges tokens. The key modifications include:

-   **Interface Simplification**: Removed unused parameters and streamlined type definitions
-   **Badge Token Support**: Integrated native badge token handling with automatic conversion
-   **Enhanced Compatibility**: Ensures seamless operation with existing DeFi infrastructure

## Key Concepts

### Automated Market Maker (AMM)

An AMM is a decentralized exchange protocol that uses mathematical formulas to determine token prices and facilitate trades without the need for traditional orderbooks.

### Liquidity Pools

Liquidity pools are smart contracts that hold pairs of tokens and allow users to trade between them. Each pool has:

-   **Pool Assets**: The tokens held in the pool
-   **Pool Shares**: LP tokens representing ownership of the pool
-   **Swap Fee**: Fee charged on each trade
-   **Exit Fee**: Fee charged when exiting the pool

### Pool Types

The GAMM module currently supports:

-   **Balancer Pools**: Standard AMM pools with configurable weights
