# Quest Token Collection Example

This example demonstrates creating a quest collection.

## Transaction Structure

```json
[
    {
        "creator": "bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl",
        "collectionId": "0",
        "balancesType": "Standard",
        "defaultBalances": {
            "balances": [],
            "outgoingApprovals": [],
            "incomingApprovals": [],
            "autoApproveSelfInitiatedOutgoingTransfers": true,
            "autoApproveSelfInitiatedIncomingTransfers": true,
            "autoApproveAllIncomingTransfers": true,
            "userPermissions": {
                "canUpdateOutgoingApprovals": [],
                "canUpdateIncomingApprovals": [],
                "canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [],
                "canUpdateAutoApproveSelfInitiatedIncomingTransfers": [],
                "canUpdateAutoApproveAllIncomingTransfers": []
            }
        },
        "validBadgeIds": [
            {
                "start": "1",
                "end": "1"
            }
        ],
        "collectionPermissions": {
            "canDeleteCollection": [],
            "canArchiveCollection": [],
            "canUpdateOffChainBalancesMetadata": [],
            "canUpdateStandards": [],
            "canUpdateCustomData": [],
            "canUpdateManager": [],
            "canUpdateCollectionMetadata": [],
            "canUpdateValidBadgeIds": [],
            "canUpdateBadgeMetadata": [],
            "canUpdateCollectionApprovals": []
        },
        "managerTimeline": [
            {
                "manager": "bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl",
                "timelineTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ]
            }
        ],
        "collectionMetadataTimeline": [
            {
                "collectionMetadata": {
                    "uri": "ipfs://QmRbRYYyphz73apphqP3QQmkeZxbtMWmAxasGfhcw1RApD",
                    "customData": ""
                },
                "timelineTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ]
            }
        ],
        "badgeMetadataTimeline": [
            {
                "badgeMetadata": [
                    {
                        "uri": "ipfs://QmRbRYYyphz73apphqP3QQmkeZxbtMWmAxasGfhcw1RApD",
                        "customData": "",
                        "badgeIds": [
                            {
                                "start": "1",
                                "end": "18446744073709551615"
                            }
                        ]
                    }
                ],
                "timelineTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ]
            }
        ],
        "offChainBalancesMetadataTimeline": [],
        "customDataTimeline": [],
        "collectionApprovals": [
            {
                "fromListId": "Mint",
                "toListId": "All",
                "initiatedByListId": "All",
                "transferTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ],
                "badgeIds": [
                    {
                        "start": "1",
                        "end": "1"
                    }
                ],
                "ownershipTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ],
                "uri": "ipfs://QmPUAjAPMDQMJZV8mpnaYbhBu2BUS4u449c7KsRNZip9uf",
                "customData": "",
                "approvalId": "quests-approval",
                "approvalCriteria": {
                    "merkleChallenges": [
                        {
                            "root": "5958c51f7c54d8e27ac42a9a2f03069c1412071abb87bf0e7be0dde790a82dbb",
                            "expectedProofLength": "0",
                            "useCreatorAddressAsLeaf": false,
                            "maxUsesPerLeaf": "1",
                            "uri": "ipfs://QmRsSK3Fw63bcJPuiYutNfBK3TYdnB8X5QG8W6ksVMuNcH",
                            "customData": "",
                            "challengeTrackerId": "1c5b9f3c390d26981996a6b593fe42300023b0e43534954a73075b912d9ca2e6",
                            "leafSigner": "0xa612B14Ff99DAe9FBC9613bF4553781086c5F887"
                        }
                    ],
                    "ethSignatureChallenges": [],
                    "predeterminedBalances": {
                        "manualBalances": [],
                        "incrementedBalances": {
                            "startBalances": [
                                {
                                    "amount": "1",
                                    "ownershipTimes": [
                                        {
                                            "start": "1",
                                            "end": "18446744073709551615"
                                        }
                                    ],
                                    "badgeIds": [
                                        {
                                            "start": "1",
                                            "end": "1"
                                        }
                                    ]
                                }
                            ],
                            "incrementBadgeIdsBy": "0",
                            "incrementOwnershipTimesBy": "0",
                            "durationFromTimestamp": "0",
                            "allowOverrideTimestamp": false,
                            "recurringOwnershipTimes": {
                                "startTime": "0",
                                "intervalLength": "0",
                                "chargePeriodLength": "0"
                            },
                            "allowOverrideWithAnyValidBadge": false
                        },
                        "orderCalculationMethod": {
                            "useOverallNumTransfers": true,
                            "usePerToAddressNumTransfers": false,
                            "usePerFromAddressNumTransfers": false,
                            "usePerInitiatedByAddressNumTransfers": false,
                            "useMerkleChallengeLeafIndex": false,
                            "challengeTrackerId": ""
                        }
                    },
                    "approvalAmounts": {
                        "overallApprovalAmount": "0",
                        "perToAddressApprovalAmount": "0",
                        "perFromAddressApprovalAmount": "0",
                        "perInitiatedByAddressApprovalAmount": "0",
                        "amountTrackerId": "quests-approval",
                        "resetTimeIntervals": {
                            "startTime": "0",
                            "intervalLength": "0"
                        }
                    },
                    "maxNumTransfers": {
                        "overallMaxNumTransfers": "1",
                        "perToAddressMaxNumTransfers": "0",
                        "perFromAddressMaxNumTransfers": "0",
                        "perInitiatedByAddressMaxNumTransfers": "0",
                        "amountTrackerId": "quests-approval",
                        "resetTimeIntervals": {
                            "startTime": "0",
                            "intervalLength": "0"
                        }
                    },
                    "coinTransfers": [
                        {
                            "to": "",
                            "coins": [
                                {
                                    "denom": "ubadge",
                                    "amount": "5000000000"
                                }
                            ],
                            "overrideFromWithApproverAddress": true,
                            "overrideToWithInitiator": true
                        }
                    ],
                    "requireToEqualsInitiatedBy": false,
                    "requireFromEqualsInitiatedBy": false,
                    "requireToDoesNotEqualInitiatedBy": false,
                    "requireFromDoesNotEqualInitiatedBy": false,
                    "overridesFromOutgoingApprovals": true,
                    "overridesToIncomingApprovals": false,
                    "autoDeletionOptions": {
                        "afterOneUse": false,
                        "afterOverallMaxNumTransfers": false
                    },
                    "userRoyalties": {
                        "percentage": "0",
                        "payoutAddress": ""
                    },
                    "mustOwnBadges": []
                },
                "version": "0"
            }
        ],
        "standardsTimeline": [
            {
                "standards": ["Quests"],
                "timelineTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ]
            }
        ],
        "isArchivedTimeline": [],
        "mintEscrowCoinsToTransfer": [
            {
                "denom": "ubadge",
                "amount": "5000000000"
            }
        ],
        "cosmosCoinWrapperPathsToAdd": []
    }
]
```

## Key Features

### Quest Token Collection

This example creates a quest collection with the following characteristics:

-   **Single Token**: Only token ID 1 is valid (`validBadgeIds: [{"start": "1", "end": "1"}]`)
-   **Quest Standard**: Uses the "Quests" standard for quest-related functionality
-   **Merkle Proof Verification**: Requires users to provide valid Merkle proofs (BitBadges claims) to claim tokens
-   **Coin Rewards**: Transfers 5000000000 ubadge coins to successful claimants and properly handles the Mint escrow coins

### Approval System

#### Collection Approval (Minting)

-   **From/To**: Mint → All users
-   **Merkle Challenge**: Root hash `5958c51f7c54d8e27ac42a9a2f03069c1412071abb87bf0e7be0dde790a82dbb`
-   **Leaf Signer**: `0xa612B14Ff99DAe9FBC9613bF4553781086c5F887` (Ethereum address)
-   **Max Uses**: 1 per leaf, overall max 1 transfer
-   **Coin Transfer**: 5000000000 ubadge coins per successful claim

#### Default User Approvals

-   **Incoming**: Allows all incoming transfers from any source
-   **Auto-approve**: Self-initiated transfers are automatically approved
-   **User Permissions**: No user permissions to update approvals (all soft-enabled)

### Manager and Permissions

-   **Manager**: `bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl` (one manager)
-   **Collection Permissions**: All permissions are empty (all soft-enabled)
-   **Timeline**: All configurations are forever (start: 1, end: max uint64)

### Metadata

-   **Collection URI**: `ipfs://QmRbRYYyphz73apphqP3QQmkeZxbtMWmAxasGfhcw1RApD`
-   **Token URI**: Same as collection URI for all tokens
-   **Approval URI**: `ipfs://QmPUAjAPMDQMJZV8mpnaYbhBu2BUS4u449c7KsRNZip9uf`
-   **Challenge URI**: `ipfs://QmRsSK3Fw63bcJPuiYutNfBK3TYdnB8X5QG8W6ksVMuNcH`

### Escrow and Funding

-   **Mint Escrow**: 5000000000 ubadge coins are escrowed to fund the coin transfers
-   **No Cosmos Coin Wrapper**: Empty `cosmosCoinWrapperPathsToAdd` array. No wrapping allowed.
