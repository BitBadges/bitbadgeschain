# Badge Ownership

Require specific badge holdings from the initiator. Checked on-chain.

## Interface

```typescript
interface MustOwnBadges<T extends NumberType> {
    collectionId: T;
    amountRange: UintRange<T>; // Min/max amount expected
    ownershipTimes: UintRange<T>[];
    badgeIds: UintRange<T>[];

    overrideWithCurrentTime: boolean; // Use current block time. Overrides ownershipTimes with [{ start: currentTime, end: currentTime }]
    mustSatisfyForAllAssets: boolean; // All vs one badge requirement
}
```
