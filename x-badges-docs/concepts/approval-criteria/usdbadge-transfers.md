# Coin Transfers

Automatic token transfers executed on every approval use. Supports any Cosmos SDK denomination. These are triggered every time an approval is used.

To get all supported denominations, use the query parameters for the badges module.

## Interface

```typescript
interface iCoinTransfer<T extends NumberType> {
    to: string; // Recipient BitBadges address
    coins: iCosmosCoin<T>[];

    overrideFromWithApproverAddress: boolean; // By default, this is the initiator address
    overrideToWithInitiator: boolean; // By default, this is the to address specified
}

interface iCosmosCoin<T extends NumberType> {
    amount: T;
    denom: string; // Any Cosmos SDK denomination (e.g., "ubadge", "uatom", "uosmo")
}
```

## Mint Escrow Address

For collection approvals with `overrideFromWithApproverAddress: true`, the approver address is a special mint escrow address.

### Generation

```typescript
const mintEscrowAddress = generateAlias(
    'badges',
    getAliasDerivationKeysForCollection(collectionId)
);
```

### Properties

-   Longer than normal addresses
-   No private key (cannot be controlled by users)
-   Can receive Cosmos-native tokens
-   Only collection approvals can trigger transfers from it

## Example

```json
[
    {
        "to": "bb1...",
        "coins": [{ "amount": "1000000000", "denom": "ubadge" }],
        "overrideFromWithApproverAddress": false,
        "overrideToWithInitiator": false
    }
]
```
