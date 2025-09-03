# MsgSwapExactAmountIn

Swaps an exact amount of tokens in for a minimum amount of tokens out.

This message allows users to swap a specific amount of input tokens for output tokens, with slippage protection through the minimum output amount.

## Swap Properties

When executing a swap:

-   Exact input amount is specified
-   Minimum output amount provides slippage protection
-   Swap fee is deducted from the input
-   Price impact is calculated based on pool liquidity

## Proto Definition

```protobuf
// ===================== MsgSwapExactAmountIn
message MsgSwapExactAmountIn {
  option (amino.name) = "gamm/swap-exact-amount-in";
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  repeated poolmanager.v1beta1.SwapAmountInRoute routes = 2
      [ (gogoproto.nullable) = false ];
  cosmos.base.v1beta1.Coin token_in = 3 [
    (gogoproto.moretags) = "yaml:\"token_in\"",
    (gogoproto.nullable) = false
  ];
  string token_out_min_amount = 4 [

    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.moretags) = "yaml:\"token_out_min_amount\"",
    (gogoproto.nullable) = false
  ];
}

message MsgSwapExactAmountInResponse {
  string token_out_amount = 1 [

    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.moretags) = "yaml:\"token_out_amount\"",
    (gogoproto.nullable) = false
  ];
}
```

### JSON Example

```json
{
    "sender": "bb1abc123...",
    "routes": [
        {
            "pool_id": "1",
            "token_out_denom": "uosmo"
        }
    ],
    "token_in": {
        "denom": "uatom",
        "amount": "1000000"
    },
    "token_out_min_amount": "5000000"
}
```

## Multi-Hop Swaps

The `routes` field allows for multi-hop swaps through multiple pools. The swap will execute through each pool in sequence.

## Slippage Protection

The `token_out_min_amount` field ensures that the user receives at least the specified amount of output tokens, protecting against price slippage.

## Swap Fees

Each pool in the swap route charges a swap fee, which is deducted from the input amount before the swap is executed.
