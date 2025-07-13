# Badge Ownership

Must own badges are another unique feature that is very powerful. This allows you to specify certain badges and amounts of badges of a collection (typically a different collection) that must be owned in order to be approved. This is checked on-chain.

For example, you may implement a badge collection where only holders of a verified badge are approved to send and receive badges. Or, you may implement what you must NOT own (own x0) a scammer badge in order to interact.

Note that alternatively, you may choose to check / enforce this off-chain as well via BitBadges claims.

```typescript
export interface MustOwnBadges<T extends NumberType> {
  collectionId: T;

  amountRange: UintRange<T>; //min/max amount expected to be owned
  ownershipTimes: UintRange<T>[];
  badgeIds: UintRange<T>[];
  
  //override ownershipTimes with the exact block millisecond at execution
  //Ex: [{start: 12345, end: 12345}]
  overrideWithCurrentTime: boolean;
  
  //if true, must meet ownership requirements for ALL badges
  //if false, must meet ownership requirements for ONE badge
  mustSatisfyForAllAssets: boolean; 
}
```
