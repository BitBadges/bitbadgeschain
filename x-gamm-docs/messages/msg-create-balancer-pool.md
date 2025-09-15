# MsgCreateBalancerPool

Creates a new balancer pool.

The poolId will be assigned at execution time and is obtainable in the transaction response. The pool creator must provide initial liquidity and set pool parameters.

## Pool Creation Properties

The creation transaction for a balancer pool is unique in several ways:

-   Initial liquidity must be provided by the creator
-   Pool parameters like swap fee and exit fee are set
-   Token weights are configured for the pool assets
-   A dedicated module account is created for the pool

## Proto Definition

```protobuf
// ===================== MsgCreatePool
message MsgCreateBalancerPool {
  option (amino.name) = "gamm/create-balancer-pool";
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];

  gamm.poolmodels.balancer.PoolParams pool_params = 2
      [ (gogoproto.moretags) = "yaml:\"pool_params\"" ];

  repeated gamm.poolmodels.balancer.PoolAsset pool_assets = 3
      [ (gogoproto.nullable) = false ];
}

// Returns the poolID
message MsgCreateBalancerPoolResponse {
  uint64 pool_id = 1 [ (gogoproto.customname) = "PoolID" ];
}
```

### JSON Example

```json
{
    "sender": "bb1abc123...",
    "pool_params": {
        "swap_fee": "0.003", // Note: Depending on serialization format this may be something like "300000000" for compatibility with sdk.Int format
        "exit_fee": "0.000"
    },
    "pool_assets": [
        {
            "token": {
                "denom": "ubadge",
                "amount": "1000000"
            },
            "weight": "50"
        },
        {
            "token": {
                "denom": "badgeslp:21:utoken",
                "amount": "5000000"
            },
            "weight": "50"
        }
    ]
}
```
