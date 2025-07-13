# $BADGE Transfers

**coinTransfers** are the $BADGE credits to be sent **every** use of the approval. There can be multiple transfers here to implement complex royalty systems, for example. Or, you can leave it blank for no $BADGE transfers. The only allowed denom is "ubadge" which has 9 decimals. 1e9 ubadge = 1 $BADGE.

Note: $BADGE refers to the native gas credits token of the blockchain, not a specific badge.

This will be executed every time the approval is used.

For collection approvals that **overrideFromWithApproverAddress**, the approver address will be overridden with a special mint escrow address. This address can be used to transfer / send $BADGE via collection approvals for the Mint address. For example, quest payouts are implemented by using this address to send $BADGE from the Mint address to the initiator.

```typescript
// To generate the mint escrow address for a collection ID
const mintEscrowAddress = generateAlias(
    'badges',
    getAliasDerivationKeysForCollection(doc.collectionId)
);
```

The Mint escrow address is slightly longer than a normal address and cannot be controlled by any end-user because it has no private key. It is technically an address, so it can receive badges, $BADGE, etc. However, the only way to trigger a $BADGE transfer from the mint escrow address is via collection approvals.

```typescript
export interface iCoinTransfer<T extends NumberType> {
    /**
     * The recipient of the coin transfer. This should be a Bech32 BitBadges address.
     */
    to: string;
    /**
     * The coins
     */
    coins: iCosmosCoin<T>[];
    /**
     * Whether or not to override the from address with the approver address.
     */
    overrideFromWithApproverAddress: boolean;
    /**
     * Whether or not to override the to address with the initiator of the transaction.
     */
    overrideToWithInitiator: boolean;
}
```

```typescript
export interface iCosmosCoin<T extends NumberType> {
    /** The amount of the coin. */
    amount: T;
    /** The denomination of the coin. */
    denom: string;
}
```
