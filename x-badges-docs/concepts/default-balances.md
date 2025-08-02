# Default Balances

Default balances are predefined balance stores that are automatically assigned to new users (uninitialized balance stores) when they first interact with a collection. These defaults are set during collection creation and cannot be updated after genesis.

## Structure

Default balances follow the same `UserBalanceStore` structure:

```json
{
    "defaultBalances": {
        "balances": [],
        "outgoingApprovals": [],
        "incomingApprovals": [],
        "autoApproveSelfInitiatedOutgoingTransfers": false,
        "autoApproveSelfInitiatedIncomingTransfers": true,
        "userPermissions": {
            // User permission configuration
        }
    }
}
```

## Important Limitations

### No Complex Approval Criteria

Default balances **cannot** include:

-   **Approval criteria** with complex conditions (merkle challenges, ETH signature challenges, dynamic store challenges, etc.)
-   **Coin transfers** or native token requirements
-   **Badge ownership** requirements or other side effects
-   **Advanced conditional logic**

Default balances are limited to basic approval structures without complex criteria.

## Related Concepts

-   **[Balance System](balance-system.md)** - How user balances work and inherit from defaults
-   **[Transferability & Approvals](transferability-approvals.md)** - User-level approval system
-   **[Permissions](permissions/)** - User permissions for updating their own approvals
-   **[Manager](manager.md)** - Collection-level controls that can override user defaults

Default balances provide a powerful way to establish baseline behavior for all users while maintaining the flexibility for users to customize their own approval settings after initialization.
