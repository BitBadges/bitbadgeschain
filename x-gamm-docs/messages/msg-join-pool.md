# MsgJoinPool

Joins an existing pool by providing liquidity. Users can join a pool by providing tokens proportional to the current pool composition. In return, they receive LP tokens representing their share of the pool.

## Join Pool Properties

When joining a pool:

-   Tokens must be provided in the correct proportions
-   LP tokens are minted to the user
-   Pool liquidity increases
-   User becomes eligible for trading fees

## Proto Definition

```protobuf
// ===================== MsgJoinPool
// This is really MsgJoinPoolNoSwap
message MsgJoinPool {
  option (amino.name) = "gamm/join-pool";
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  uint64 pool_id = 2 [ (gogoproto.moretags) = "yaml:\"pool_id\"" ];
  string share_out_amount = 3 [

    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.moretags) = "yaml:\"pool_amount_out\"",
    (gogoproto.nullable) = false
  ];
  repeated cosmos.base.v1beta1.Coin token_in_maxs = 4 [
    (gogoproto.moretags) = "yaml:\"token_in_max_amounts\"",
    (gogoproto.nullable) = false
  ];
}

message MsgJoinPoolResponse {
  string share_out_amount = 1 [

    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.moretags) = "yaml:\"share_out_amount\"",
    (gogoproto.nullable) = false
  ];
  repeated cosmos.base.v1beta1.Coin token_in = 2 [
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
    "share_out_amount": "1000000",
    "token_in_maxs": [
        {
            "denom": "uatom",
            "amount": "100000"
        },
        {
            "denom": "uosmo",
            "amount": "500000"
        }
    ]
}
```

## Token Proportions

The tokens provided must be in the same proportion as the current pool composition. If not, the transaction will fail or tokens will be returned.

## Slippage Protection

The `token_in_maxs` field provides slippage protection by setting maximum amounts for each token that can be used in the join operation.

## LP Token Minting

Upon successful join, LP tokens are minted to the user's address. These tokens represent ownership of the pool and can be used for:

-   Earning trading fees
-   Governance participation
-   Staking in yield farming programs
