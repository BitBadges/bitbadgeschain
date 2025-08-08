# Empty Approval Criteria Template

When creating collection approvals with empty approval criteria, you can use this template for "no additional restrictions". We reference this for simplicity in other examples.

## Template

```typescript
const EmptyApprovalCriteria = {
    approvalCriteria: {
        // No challenges to be completed
        merkleChallenges: [],
        // No specific balances to check
        predeterminedBalances: {
            manualBalances: [],
            incrementedBalances: {
                startBalances: [],
                incrementBadgeIdsBy: '0',
                incrementOwnershipTimesBy: '0',
                durationFromTimestamp: '0',
                allowOverrideTimestamp: false,
                recurringOwnershipTimes: {
                    startTime: '0',
                    intervalLength: '0',
                    chargePeriodLength: '0',
                },
                allowOverrideWithAnyValidBadge: false,
            },
            orderCalculationMethod: {
                useOverallNumTransfers: false,
                usePerToAddressNumTransfers: false,
                usePerFromAddressNumTransfers: false,
                usePerInitiatedByAddressNumTransfers: false,
                useMerkleChallengeLeafIndex: false,
                challengeTrackerId: '',
            },
        },
        // No approval amounts to check (0 = unlimited)
        approvalAmounts: {
            overallApprovalAmount: '0',
            perToAddressApprovalAmount: '0',
            perFromAddressApprovalAmount: '0',
            perInitiatedByAddressApprovalAmount: '0',
            amountTrackerId:
                'a4ab9bc5e8752842a35a79238de4f627677ceae1d8fa9de44b52416e085f7f11',
            resetTimeIntervals: {
                startTime: '0',
                intervalLength: '0',
            },
        },
        // No max number of transfers to check (0 = unlimited)
        maxNumTransfers: {
            overallMaxNumTransfers: '0',
            perToAddressMaxNumTransfers: '0',
            perFromAddressMaxNumTransfers: '0',
            perInitiatedByAddressMaxNumTransfers: '0',
            amountTrackerId:
                'd711e23dbe57b786dfb2d86d4a6792fb8c9951a18223065ea0c07d424225a738',
            resetTimeIntervals: {
                startTime: '0',
                intervalLength: '0',
            },
        },
        // No coin transfers to execute
        coinTransfers: [],

        // No ETH signature challenges to be completed
        ethSignatureChallenges: [],
        // No dynamic store challenges to be completed
        dynamicStoreChallenges: [],

        // No address matching requirements
        requireToEqualsInitiatedBy: false,
        requireFromEqualsInitiatedBy: false,
        requireToDoesNotEqualInitiatedBy: false,
        requireFromDoesNotEqualInitiatedBy: false,
        // No overrides from outgoing approvals
        overridesFromOutgoingApprovals: false,
        // No overrides to incoming approvals
        overridesToIncomingApprovals: false,
        // No auto deletion options
        autoDeletionOptions: {
            afterOneUse: false,
            afterOverallMaxNumTransfers: false,
        },
        // No user royalties
        userRoyalties: {
            percentage: '0',
            payoutAddress: '',
        },
        // No tokens to check ownership of
        mustOwnBadges: [],
    },
};
```

## Related Documentation

-   [Approval Criteria Overview](../concepts/approval-criteria/README.md)
-   [Building Collection Approvals](./building-collection-approvals.md)
-   [Transferability / Approvals](../concepts/transferability-approvals.md)
