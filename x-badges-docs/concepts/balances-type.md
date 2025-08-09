# Balances Type

The balances type determines how token ownership and transfers are managed within a collection. This is permanent upon genesis and cannot be changed.

> **Important**: Always use `"Standard"` for new collections. Other balance types are supported behind the scenes for legacy purposes only, but "Standard" should always be used for on-chain balances.

## Standard Balances Type

```json
"balancesType": "Standard"
```

With standard balances, everything is facilitated on-chain in a decentralized manner. All balances are stored on the blockchain, and everything is facilitated through on-chain transfers and approvals.

## Key Features

### On-Chain Storage

-   All balances are stored directly on the blockchain
-   Provides complete transparency and decentralization
-   No reliance on external systems for balance verification

### Transfer Requirements

All transfers require:

-   **Sufficient balances** in the sender's account
-   **Valid approvals** for the collection, sender, and recipient where necessary
-   **Three-tier approval system** verification (collection, sender, recipient)

### Approval Structure

All transfers must specify:

-   **Collection approvals** - Managed by the collection manager
-   **Sender's outgoing approvals** - User-controlled outgoing approval permissions (if not forcefully overridden by collection approvals)
-   **Recipient's incoming approvals** - User-controlled incoming approval permissions (if not forcefully overridden by collection approvals)

The collection approvals are managed by the manager and can optionally override the user-level approvals.

### Mint Address Behavior

The "Mint" address has special properties:

-   **Unlimited balances** - Can mint any amount of tokens
-   **Send-only** - Can only send tokens, not receive them
-   **Circulating supply control** - All circulating tokens originate from Mint transfers. Thus, the circulating supply is controlled by the collection approvals from the Mint address.
-   **Non-Controllable** - The Mint address cannot set its own approvals, so all approvals must be set by the collection manager and forcefully override the Mint address's approvals.

### Circulating Supply Management

The circulating supply is controlled by:

-   **Transfers from Mint address** - Initial distribution mechanism
-   **Collection approvals** - Manager-controlled transfer rules
-   **User approvals** - Individual transfer permissions
-   **Collection permissions** - Controls approval updatability
