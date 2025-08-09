# Simple Token Transfer

This example demonstrates a basic token transfer from the mint to a specific address.

## Overview

This transfer creates token ID 1 from collection 20 and sends it to the creator address. The token has full ownership time range and uses collection-level approval.

## Transfer Details

-   **Collection ID**: 20
-   **Token ID**: 1
-   **Amount**: 1
-   **From**: Mint (new token creation)
-   **To**: Creator address
-   **Approval**: Collection-level approval (assumes user-level approvals successfully auto-scan)

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
                "balances": [
                    {
                        "amount": "1",
                        "ownershipTimes": [
                            {
                                "start": "1",
                                "end": "18446744073709551615"
                            }
                        ],
                        "tokenIds": [
                            {
                                "start": "1",
                                "end": "1"
                            }
                        ]
                    }
                ],
                "precalculateBalancesFromApproval": {
                    "approvalId": "",
                    "approvalLevel": "",
                    "approverAddress": "",
                    "version": "0"
                },
                "merkleProofs": [],
                "ethSignatureProofs": [],
                "memo": "",
                "prioritizedApprovals": [
                    {
                        "approvalId": "4a1ed47db7bc0f9f7174eab12aa9b8c9b9e4e37474ca2264668cf8e1b1598dde",
                        "approvalLevel": "collection",
                        "approverAddress": "",
                        "version": "0"
                    }
                ],
                "onlyCheckPrioritizedCollectionApprovals": true,
                "onlyCheckPrioritizedIncomingApprovals": false,
                "onlyCheckPrioritizedOutgoingApprovals": false,
                "affiliateAddress": "",
                "numAttempts": "1"
            }
        ]
    }
]
```

## Key Components Explained

### Transfer Source

-   `"from": "Mint"` - Indicates this is a new token creation from the mint

### Destination

-   `"toAddresses": ["bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl"]` - The recipient address

### Balance Specification

-   `"amount": "1"` - Transfer 1 token
-   `"ownershipTimes"` - Full ownership time range (1 to max uint64)
-   `"tokenIds"` - Specific token ID range (1 to 1)

### Approval Configuration

-   `"prioritizedApprovals"` - Uses collection-level approval
-   `"onlyCheckPrioritizedCollectionApprovals": true` - Only check collection approvals
-   `"approvalId"` - Specific approval identifier for the collection

### Additional Settings

-   `"merkleProofs": []` - No merkle proofs required for this simple transfer
-   `"ethSignatureProofs": []` - No ETH signature proofs required for this simple transfer
-   `"memo": ""` - No memo attached
-   `"numAttempts": "1"` - Single transfer attempt

## Usage

This example can be used as a template for basic token minting operations where you want to create a new token and transfer it to a specific address using collection-level approval.
