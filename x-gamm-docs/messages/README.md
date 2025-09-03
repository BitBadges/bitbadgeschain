# Messages

This directory contains documentation for the core transaction messages supported by the `x/gamm` module.

## Core Operations

### Pool Management

-   [MsgCreateBalancerPool](./msg-create-balancer-pool.md) - Create new balancer pool with initial liquidity
-   [MsgJoinPool](./msg-join-pool.md) - Join existing pool with liquidity
-   [MsgExitPool](./msg-exit-pool.md) - Exit pool and receive underlying tokens

### Trading

-   [MsgSwapExactAmountIn](./msg-swap-exact-amount-in.md) - Swap exact amount of tokens in

## Core Message List

The GAMM module supports 4 core transaction messages:

1. **MsgCreateBalancerPool** - Create new balancer pool
2. **MsgJoinPool** - Join pool with proportional tokens
3. **MsgSwapExactAmountIn** - Swap exact input for minimum output
4. **MsgExitPool** - Exit pool and receive underlying tokens
