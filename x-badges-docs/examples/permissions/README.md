# Permission Examples

This directory contains practical examples of different permission configurations for badge collections. Each example demonstrates specific patterns and use cases for controlling collection management.

## Contents

-   [Freezing Mint Transferability](freezing-mint-transferability.md) - Permanently freeze minting capabilities
-   [Locking Specific Approval ID](locking-specific-approval-id.md) - Lock specific approval IDs with granular control
-   [Locking Specific Badge IDs](locking-specific-badge-ids.md) - Lock approvals for specific badge ID ranges
-   [Locking Valid Badge IDs](locking-valid-badge-ids.md) - Control valid badge ID range updates
-   [Locked Collection](locked-collection.md) - Collection with permanently locked permissions
-   [Temporary Permissions](temporary-permissions.md) - Time-limited management permissions
-   [Badge Specific Permissions](badge-specific-permissions.md) - Permissions that apply to specific badge IDs
-   [Community Controlled](community-controlled.md) - Permissions for community-managed collections

## Permission System Overview

BitBadges permissions follow a timeline-based system where:

1. **Permanently Permitted Times** - Permission is always allowed
2. **Permanently Forbidden Times** - Permission is always denied
3. **Default (Empty)** - Permission is soft-enabled (manager can change)

## Common Patterns

-   **No Manager** - Set manager to empty string to disable all management
-   **Complete Control** - Empty permission arrays for full soft-enabled control
-   **Locked Forever** - Use `permanentlyForbiddenTimes: FullTimeRanges`
-   **Time-Limited** - Use specific time ranges for temporary control

## Related Concepts

-   [Permission System](../../concepts/permissions/permission-system.md)
-   [Manager](../../concepts/manager.md)
-   [Timeline System](../../concepts/timeline-system.md)
