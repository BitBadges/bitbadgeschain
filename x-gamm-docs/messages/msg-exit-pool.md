# MsgExitPool

Exits an existing pool by burning LP tokens and receiving underlying tokens.

Users can exit a pool by burning their LP tokens. In return, they receive the underlying pool tokens proportional to their share.

## Exit Pool Properties

When exiting a pool:

-   LP tokens are burned from the user
-   Underlying tokens are returned proportionally
-   Pool liquidity decreases
-   Exit fees may be applied

## Proto Definition

```protobuf
// ===================== MsgExitPool
message MsgExitPool {
  option (amino.name) = "gamm/exit-pool";
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  uint64 pool_id = 2 [ (gogoproto.moretags) = "yaml:\"pool_id\"" ];
  string share_in_amount = 3 [

    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.moretags) = "yaml:\"share_in_amount\"",
    (gogoproto.nullable) = false
  ];

  repeated cosmos.base.v1beta1.Coin token_out_mins = 4 [
    (gogoproto.moretags) = "yaml:\"token_out_min_amounts\"",
    (gogoproto.nullable) = false
  ];
}

message MsgExitPoolResponse {
  repeated cosmos.base.v1beta1.Coin token_out = 1 [
    (gogoproto.moretags) = "yaml:\"token_out\"",
    (gogoproto.nullable) = false
  ];
}
```

### JSON Example

```json
{
    "sender": "bb1abc123...",
    "pool_id": "1",
    "share_in_amount": "100000",
    "token_out_mins": [
        {
            "denom": "uatom",
            "amount": "10000"
        },
        {
            "denom": "uosmo",
            "amount": "50000"
        }
    ]
}
```

## Token Proportions

The tokens received will be in the same proportion as the current pool composition. The user cannot specify which tokens to receive.

## Slippage Protection

The `token_out_mins` field provides slippage protection by setting minimum amounts for each token that must be received from the exit operation.

## Exit Fees

Some pools may charge exit fees, which are deducted from the tokens returned to the user.

## LP Token Burning

Upon successful exit, LP tokens are burned from the user's address, reducing their ownership stake in the pool.
