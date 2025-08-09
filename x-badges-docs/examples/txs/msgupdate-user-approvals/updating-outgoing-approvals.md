# Updating Outgoing Approvals

This example demonstrates how to update user outgoing approvals to control what tokens a user can send and to whom.

## Transaction Structure

```json
[
    {
        "creator": "bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl",
        "collectionId": "1",
        "updateOutgoingApprovals": true,
        "outgoingApprovals": [
            {
                "toListId": "All",
                "initiatedByListId": "bb18el5ug46umcws58m445ql5scgg2n3tzagfecvl",
                "transferTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ],
                "badgeIds": [
                    {
                        "start": "1",
                        "end": "20"
                    }
                ],
                "ownershipTimes": [
                    {
                        "start": "1",
                        "end": "18446744073709551615"
                    }
                ],
                "uri": "",
                "customData": "",
                "approvalId": "87bc6dd97492b913b3d2b6c91c71b7a2bc98d41a715e49285180e8db9f4ea0bb",
                "approvalCriteria": {
                    "merkleChallenges": [],
                    "ethSignatureChallenges": [],
                    "predeterminedBalances": {
                        "manualBalances": [],
                        "incrementedBalances": {
                            "startBalances": [],
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
                            "useOverallNumTransfers": false,
                            "usePerToAddressNumTransfers": false,
                            "usePerFromAddressNumTransfers": false,
                            "usePerInitiatedByAddressNumTransfers": false,
                            "useMerkleChallengeLeafIndex": false,
                            "challengeTrackerId": ""
                        }
                    },
                    "approvalAmounts": {
                        "overallApprovalAmount": "1",
                        "perToAddressApprovalAmount": "0",
                        "perFromAddressApprovalAmount": "0",
                        "perInitiatedByAddressApprovalAmount": "1",
                        "amountTrackerId": "87bc6dd97492b913b3d2b6c91c71b7a2bc98d41a715e49285180e8db9f4ea0bb",
                        "resetTimeIntervals": {
                            "startTime": "0",
                            "intervalLength": "0"
                        }
                    },
                    "maxNumTransfers": {
                        "overallMaxNumTransfers": "0",
                        "perToAddressMaxNumTransfers": "0",
                        "perFromAddressMaxNumTransfers": "0",
                        "perInitiatedByAddressMaxNumTransfers": "0",
                        "amountTrackerId": "fe1ffc5f6ff98f0e41b097f33623248868d367dc36dd7f22b2717b61b9d7c91c",
                        "resetTimeIntervals": {
                            "startTime": "0",
                            "intervalLength": "0"
                        }
                    },
                    "coinTransfers": [],
                    "requireToEqualsInitiatedBy": false,
                    "requireToDoesNotEqualInitiatedBy": false,
                    "autoDeletionOptions": {
                        "afterOneUse": false,
                        "afterOverallMaxNumTransfers": false
                    },
                    "mustOwnBadges": []
                },
                "version": "0"
            }
        ],

        // All other updates are false, so values do not ma
        "updateIncomingApprovals": false,
        "incomingApprovals": [],
        "updateAutoApproveSelfInitiatedOutgoingTransfers": true,
        "autoApproveSelfInitiatedOutgoingTransfers": true,
        "updateAutoApproveSelfInitiatedIncomingTransfers": false,
        "autoApproveSelfInitiatedIncomingTransfers": false,
        "updateAutoApproveAllIncomingTransfers": false,
        "autoApproveAllIncomingTransfers": false,
        "updateUserPermissions": false,
        "userPermissions": {
            "canUpdateOutgoingApprovals": [],
            "canUpdateIncomingApprovals": [],
            "canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [],
            "canUpdateAutoApproveSelfInitiatedIncomingTransfers": [],
            "canUpdateAutoApproveAllIncomingTransfers": []
        }
    }
]
```

## Related Examples

-   [Building User Approvals](../../building-user-approvals.md) - User approval patterns
-   [Empty Approval Criteria](../../empty-approval-criteria.md) - Template for unrestricted approvals

## Related Concepts

-   [Transferability / Approvals](../../../concepts/transferability-approvals.md) - Approval system overview
-   [Approval Criteria](../../../concepts/approval-criteria/README.md) - Criteria configuration
-   [Tallied Approval Amounts](../../../concepts/approval-criteria/tallied-approval-amounts.md) - Amount tracking mechanics
