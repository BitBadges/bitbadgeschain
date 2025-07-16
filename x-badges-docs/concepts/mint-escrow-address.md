# Mint Escrow Address

The Mint Escrow Address (_mintEscrowAddress_) is a special reserved address generated from the collection ID that holds Cosmos native funds on behalf of the "Mint" address for a specific collection. This address has no known private key and is not controlled by anyone. The only way to get funds out is via collection approvals from the Mint address.

See the [coin transfers](../approval-criteria/usdbadge-transfers.md) section for more details.

## Functionality

### Cosmos Native Fund Storage

The Mint Escrow Address can hold Cosmos native tokens (like "ubadge" tokens) that are associated with the Mint address for a specific collection.

See the coinTransfers section for more details. This is the only way to get funds out of the Mint Escrow Address.

## Auto-Escrow During Collection Creation

The `MsgCreateCollection` interface includes a `mintEscrowCoinsToTransfer` field of type `repeated cosmos.base.v1beta1.Coin` that allows you to automatically escrow native coins to the Mint Escrow Address during collection creation.

### Pre-Creation Escrow

-   **Unknown collection ID** - Escrow coins before knowing the final collection ID
-   **Automatic transfer** - Coins are automatically transferred to the generated Mint Escrow Address
-   **Collection initialization** - Funds are available immediately when the collection is created
-   **Single transaction** - Combine collection creation and coin escrow in one operation

### Usage

```json
{
    "creator": "cosmos1...",
    "collectionId": "0",
    "mintEscrowCoinsToTransfer": [
        {
            "denom": "ubadge",
            "amount": "1000000"
        }
    ]
    // ... other collection fields
}
```

This field is particularly useful when you need to fund the Mint Escrow Address but don't know the collection ID beforehand, since the escrow address is derived from the collection ID itself. Thus, it can be done all in one transaction.
