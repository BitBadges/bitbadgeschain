# 🧠 Concepts

This directory contains detailed explanations of the core concepts and data structures that form the foundation of the BitBadges module.

-   **[UintRange](uintrange.md)** - Representing ranges of unsigned integers
-   **[Timeline System](timeline-system.md)** - How properties change over time with immutable historical records
-   **[Badge Collections](badge-collections.md)** - The primary entity that defines groups of related badges
-   **[Balance System](balance-system.md)** - How badge ownership is tracked with precise control over quantities and time
-   **[Balances Type](balances-type.md)** - How badge ownership and transfers are managed within collections
-   **[Off-Chain Balances](off-chain-balances.md)** - Legacy off-chain balance types (deprecated - use Standard balances)
-   **[Default Balances](default-balances.md)** - Predefined balance stores assigned to new users upon genesis creation
-   **[Valid Badge IDs](valid-badge-ids.md)** - Defining the range of badge identifiers that can exist within a collection
-   **[Total Supply](total-supply.md)** - Maximum number of badges that can exist for each badge ID
-   **[Mint Escrow Address](mint-escrow-address.md)** - Reserved address for holding native funds on behalf of the Mint address
-   **[Archived Collections](archived-collections.md)** - Temporarily or permanently disabling collection transactions
-   **[Metadata](metadata.md)** - Rich, dynamic content for collections and badges with timeline support
-   **[Manager](manager.md)** - Central authority for collection administration with timeline-based control
-   **[Transferability & Approvals](transferability-approvals.md)** - Overview of the three-tier approval system for badge transfers
-   **[Approval Criteria](approval-criteria/)** - Detailed approval criteria and conditions for badge transfers
-   **[Address Lists](address-lists.md)** - Reusable collections of addresses for approval configurations
-   **[Permissions](permissions/README.md)** - Granular control over collection management operations
-   **[Standards](standards.md)** - Generic framework for defining collection behavior and interpretation guidelines
-   **[Time Fields](time-fields.md)** - Understanding the different time-related fields used throughout BitBadges
-   **[Custom Data](custom-data.md)** - Generic string fields for storing arbitrary application-specific data
-   **[Cosmos Wrapper Paths](cosmos-wrapper-paths.md)** - 1:1 wrapping between badges and native Cosmos SDK coins for IBC compatibility
-   **[Collection Invariants](collection-invariants.md)** - Immutable rules that enforce fundamental constraints on collection behavior
-   **[Protocols](protocols/)** - Standardized implementation patterns for badge collections
