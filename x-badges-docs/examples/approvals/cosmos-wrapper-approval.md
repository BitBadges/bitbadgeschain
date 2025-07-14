# Cosmos Wrapper Approval

This example demonstrates how to create an approval that allows badges to be sent to a Cosmos coin wrapper address, enabling conversion to native Cosmos SDK coins.

You pretty much: 1) figure out your address and 2) figure out a path that users can send to this address without needing the address to control its approvals. Oftentimes, you may not even need to forcefully override the incoming approvals because you default allow all incoming transfers which also applies to the wrapper address automatically.

Full example: [Cosmos Coin Wrapper Example](../cosmos-coin-wrapper-example.md)

## Code Example

```typescript
export const wrapperApproval = ({
    specialAddress,
    badgeIds,
    ownershipTimes,
    approvalId,
}: {
    specialAddress: string;
    badgeIds: iUintRange<bigint>[];
    ownershipTimes: iUintRange<bigint>[];
    approvalId: string;
}): RequiredApprovalProps => {
    const id = approvalId;
    const toSet: RequiredApprovalProps = {
        version: 0n,
        toListId: specialAddress,
        toList: AddressList.getReservedAddressList(specialAddress),
        fromListId: 'AllWithoutMint',
        fromList: AddressList.getReservedAddressList('AllWithoutMint'),
        initiatedByListId: 'All',
        initiatedByList: AddressList.AllAddresses(),
        transferTimes: UintRangeArray.FullRanges(),
        badgeIds: badgeIds,
        ownershipTimes: ownershipTimes,
        approvalId: id,
        approvalCriteria: {
            ...EmptyApprovalCriteria,
            overridesToIncomingApprovals: true,
        },
    };

    return toSet;
};
```

## Related Concepts

-   [Cosmos Wrapper Paths](../../concepts/cosmos-wrapper-paths.md)
-   [Transferability / Approvals](../../concepts/transferability-approvals.md)
-   [Address Lists](../../concepts/address-lists.md)
