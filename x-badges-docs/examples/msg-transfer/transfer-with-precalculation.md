# Transfer with Precalculation

This example demonstrates a badge transfer that uses precalculation from approval criteria instead of manually specifying balances.

## Overview

This transfer creates badges from collection 20 and sends them to the creator address. Instead of manually specifying the balance amounts, it uses precalculation from the approval criteria to determine what badges to transfer.

## Transfer Details

-   **Collection ID**: 20
-   **From**: Mint (new badge creation)
-   **To**: Creator address
-   **Approval**: Collection-level approval with precalculation
-   **Precalculation**: Enabled with specific approval ID

## JSON Structure

```json
[
    {
        "creator": "bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl",
        "collectionId": "20",
        "transfers": [
            {
                "from": "Mint",
                "toAddresses": ["bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl"],
                "balances": [],
                "precalculateBalancesFromApproval": {
                    "approvalId": "fd1cef5941fb08487ecc1038af09fb29a6d7d40a89d8e4889c9c954978aa7e41",
                    "approvalLevel": "collection",
                    "approverAddress": "",
                    "version": "0"
                },
                "merkleProofs": [],
                "memo": "",
                "prioritizedApprovals": [
                    {
                        "approvalId": "fd1cef5941fb08487ecc1038af09fb29a6d7d40a89d8e4889c9c954978aa7e41",
                        "approvalLevel": "collection",
                        "approverAddress": "",
                        "version": "0"
                    }
                ],
                "onlyCheckPrioritizedCollectionApprovals": true,
                "onlyCheckPrioritizedIncomingApprovals": false,
                "onlyCheckPrioritizedOutgoingApprovals": false,
                "precalculationOptions": {
                    "overrideTimestamp": "0",
                    "badgeIdsOverride": []
                },
                "affiliateAddress": "",
                "numAttempts": "1"
            }
        ]
    }
]
```

## Key Components Explained

### Precalculation Configuration

-   `"balances": []` - Empty balances array since amounts are calculated from approval
-   `"precalculateBalancesFromApproval"` - Specifies which approval to use for calculation
-   `"approvalId"` - The specific approval ID that defines the transfer criteria

### Prioritized Approvals

-   `"prioritizedApprovals"` - Uses the same approval ID for both precalculation and transfer
-   `"onlyCheckPrioritizedCollectionApprovals": true` - Only check collection-level approvals
-   `"onlyCheckPrioritizedIncomingApprovals": false` - Skip incoming approval checks
-   `"onlyCheckPrioritizedOutgoingApprovals": false` - Skip outgoing approval checks

### Precalculation Options

-   `"overrideTimestamp": "0"` - Use current timestamp for calculations
-   `"badgeIdsOverride": []` - No badge ID overrides, use approval criteria

### Non-Auto-Scan Behavior

This example demonstrates "prioritized non-auto-scan" behavior where:

-   Only the specified approval is checked (no automatic scanning of other approvals)
-   The system doesn't automatically look for other valid approvals
-   Transfer is limited to what the specified approval allows
-   Can use approvals with side effects and custom criteria like merkle challenges
-   Shows proper versioning of approvals

## Usage

This example is useful when:

-   You want to transfer badges based on approval criteria rather than manual specification
-   You need precise control over which approval is used
-   You want to avoid automatic approval scanning
-   The approval criteria dynamically determine badge amounts and IDs

## Differences from Simple Transfer

| Feature               | Simple Transfer      | Precalculation Transfer   |
| --------------------- | -------------------- | ------------------------- |
| Balance Specification | Manual amounts       | Calculated from approval  |
| Approval Scanning     | Auto-scan enabled    | Only specified approval   |
| Flexibility           | Fixed amounts        | Dynamic based on criteria |
| Control               | Direct specification | Approval-driven           |
