/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface BadgesActionCombination {
  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesActionDefaultValues {
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* ActionPermission defines the permissions for performing an action.

This is simple and straightforward as the only thing we need to check is the permitted/forbidden times.
*/
export interface BadgesActionPermission {
  defaultValues?: BadgesActionDefaultValues;
  combinations?: BadgesActionCombination[];
}

/**
* An AddressMapping is a permanent list of addresses that are referenced by a mapping ID.
The mapping may include only the specified addresses, or it may include all addresses but
the specified addresses (depending on if includeAddresses is true or false).

AddressMappings are used for things like whitelists, blacklists, approvals, etc.
*/
export interface BadgesAddressMapping {
  mappingId?: string;
  addresses?: string[];
  includeAddresses?: boolean;
  uri?: string;
  customData?: string;
}

export interface BadgesApprovalsTracker {
  numTransfers?: string;
  amounts?: BadgesBalance[];
}

/**
* A BadgeCollection is the top level object for a collection of badges. 
It defines everything about the collection, such as the manager, metadata, etc.

All collections are identified by a collectionId assigned by the blockchain, which is a uint64 that increments (i.e. first collection has ID 1).

All collections also have a manager who is responsible for managing the collection. 
They can be granted certain permissions, such as the ability to mint new badges.

Certain fields are timeline-based, which means they may have different values at different block heights. 
We fetch the value according to the current time.
For example, we may set the manager to be Alice from Time1 to Time2, and then set the manager to be Bob from Time2 to Time3.

Collections may have different balance types: standard vs off-chain vs inherited. See documentation for differences.
*/
export interface BadgesBadgeCollection {
  /** The collectionId is the unique identifier for this collection. */
  collectionId?: string;

  /** The collection metadata is the metadata for the collection itself. */
  collectionMetadataTimeline?: BadgesCollectionMetadataTimeline[];

  /** The badge metadata is the metadata for each badge in the collection. */
  badgeMetadataTimeline?: BadgesBadgeMetadataTimeline[];

  /** The balancesType is the type of balances this collection uses (standard, off-chain, or inherited). */
  balancesType?: string;

  /** The off-chain balances metadata defines where to fetch the balances for collections with off-chain balances. */
  offChainBalancesMetadataTimeline?: BadgesOffChainBalancesMetadataTimeline[];

  /** The inherited balances metadata defines the parent balances for collections with inherited balances. */
  inheritedBalancesTimeline?: BadgesInheritedBalancesTimeline[];

  /** The custom data field is an arbitrary field that can be used to store any data. */
  customDataTimeline?: BadgesCustomDataTimeline[];

  /** The manager is the address of the manager of this collection. */
  managerTimeline?: BadgesManagerTimeline[];

  /** The permissions define what the manager of the collection can do or not do. */
  permissions?: BadgesCollectionPermissions;

  /**
   * The approved transfers timeline defines the transferability of the collection for collections with standard balances.
   * This defines it on a collection-level. All transfers must be explicitly allowed on the collection-level, or else, they will fail.
   *
   * Collection approved transfers can optionally specify to override the user approvals for a transfer (e.g. forcefully revoke a badge).
   * If user approvals are not overriden, then a transfer must also satisfy the From user's approved outgoing transfers and the To user's approved incoming transfers.
   */
  collectionApprovedTransfersTimeline?: BadgesCollectionApprovedTransferTimeline[];

  /** Standards allow us to define a standard for the collection. This lets others know how to interpret the fields of the collection. */
  standardsTimeline?: BadgesStandardsTimeline[];

  /**
   * The isArchivedTimeline defines whether the collection is archived or not.
   * When a collection is archived, it is read-only and no transactions can be processed.
   */
  isArchivedTimeline?: BadgesIsArchivedTimeline[];

  /** The contractAddressTimeline defines the contract address for the collection (if it has a corresponding contract). */
  contractAddressTimeline?: BadgesContractAddressTimeline[];

  /**
   * The defaultUserApprovedOutgoingTransfersTimeline defines the default user approved outgoing transfers for an uninitialized user balance.
   * The user can change this value at any time.
   */
  defaultUserApprovedOutgoingTransfersTimeline?: BadgesUserApprovedOutgoingTransferTimeline[];

  /**
   * The defaultUserApprovedIncomingTransfersTimeline defines the default user approved incoming transfers for an uninitialized user balance.
   * The user can change this value at any time.
   *
   * Ex: Set this to disallow all incoming transfers by default, making the user have to opt-in to receiving the badge.
   */
  defaultUserApprovedIncomingTransfersTimeline?: BadgesUserApprovedIncomingTransferTimeline[];
}

/**
* This defines the metadata for specific badge IDs.
This should be interpreted according to the collection standard.
*/
export interface BadgesBadgeMetadata {
  uri?: string;
  customData?: string;
  badgeIds?: BadgesUintRange[];
}

export interface BadgesBadgeMetadataTimeline {
  badgeMetadata?: BadgesBadgeMetadata[];
  timelineTimes?: BadgesUintRange[];
}

/**
* Balance represents the balance of a badge for a specific user.
The user amounts xAmount of a badge for the badgeID specified for the time ranges specified.

Ex: User A owns x10 of badge IDs 1-10 from 1/1/2020 to 1/1/2021.

If times or badgeIDs have len > 1, then the user owns all badge IDs specified for all time ranges specified.
*/
export interface BadgesBalance {
  amount?: string;
  ownedTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
}

export interface BadgesBalancesActionCombination {
  /** ValueOptions defines how we manipulate the default values. */
  badgeIdsOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  ownedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesBalancesActionDefaultValues {
  badgeIds?: BadgesUintRange[];
  ownedTimes?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* BalancesActionPermission defines the permissions for updating a timeline-based field for specific badges and specific badge ownership times.
Currently, this is only used for creating new badges.

Ex: If you want to lock the ability to create new badges for badgeIds [1,2] at ownedTimes 1/1/2020 - 1/1/2021, 
you could set the combination (badgeIds: [1,2], ownershipTimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
*/
export interface BadgesBalancesActionPermission {
  defaultValues?: BadgesBalancesActionDefaultValues;
  combinations?: BadgesBalancesActionCombination[];
}

/**
* Challenges define the rules for the approval.
If all challenge are not met with valid solutions, then the transfer is not approved.

Currently, we only support Merkle tree challenges where the Merkle path must be to the provided root
and be the expected length.

We also support the following options:
-useCreatorAddressAsLeaf: If true, then the leaf will be set to the creator address. Used for whitelist trees.
-maxOneUsePerLeaf: If true, then each leaf can only be used once. If false, then the leaf can be used multiple times.
This is very important to be set to true if you want to prevent replay attacks.
-useLeafIndexForDistributionOrder: If true, we will use the leafIndex to determine the order of the distribution of badges.
leafIndex 0 will be the leftmost leaf of the expectedProofLength layer

IMPORTANT: We track the number of uses per leaf according to a challenge ID.
Please use unique challenge IDs for different challenges of the same timeline.
If you update the challenge ID, then the used leaves tracker will reset and start a new tally.
It is highly recommended to avoid updating a challenge without resetting the tally via a new challenge ID.
*/
export interface BadgesChallenge {
  root?: string;
  expectedProofLength?: string;
  useCreatorAddressAsLeaf?: boolean;
  maxOneUsePerLeaf?: boolean;
  useLeafIndexForDistributionOrder?: boolean;
  challengeId?: string;
}

export interface BadgesMerkleProof {
  proof?: BadgesClaimProof;
}

export interface BadgesClaimProof {
  leaf?: string;
  aunts?: BadgesMerklePathItem[];
}

export interface BadgesMerklePathItem {
  aunt?: string;
  onRight?: boolean;
}

export interface BadgesCollectionApprovedTransfer {
  fromMappingId?: string;
  toMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  allowedCombinations?: BadgesIsCollectionTransferAllowed[];
  challenges?: BadgesChallenge[];
  approvalId?: string;
  incrementBadgeIdsBy?: string;
  incrementOwnedTimesBy?: string;
  overallApprovals?: BadgesApprovalsTracker;

  /** PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address. */
  perAddressApprovals?: BadgesPerAddressApprovals;
  overridesFromApprovedOutgoingTransfers?: boolean;
  overridesToApprovedIncomingTransfers?: boolean;
  requireToEqualsInitiatedBy?: boolean;
  requireFromEqualsInitiatedBy?: boolean;
  requireToDoesNotEqualInitiatedBy?: boolean;
  requireFromDoesNotEqualInitiatedBy?: boolean;
  uri?: string;
  customData?: string;
}

export interface BadgesCollectionApprovedTransferCombination {
  /** ValueOptions defines how we manipulate the default values. */
  timelineTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  fromMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  toMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  initiatedByMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  transferTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  badgeIdsOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesCollectionApprovedTransferDefaultValues {
  timelineTimes?: BadgesUintRange[];
  fromMappingId?: string;
  toMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* CollectionApprovedTransferPermission defines what collection approved transfers can be updated vs are locked.

Each transfer is broken down to a (from, to, initiatedBy, transferTime, badgeId) tuple.
For a transfer to match, we need to match ALL of the fields in the combination. 
These are detemined by the fromMappingId, toMappingId, initiatedByMappingId, transferTimes, badgeIds fields.
AddressMappings are used for (from, to, initiatedBy) which are a permanent list of addresses identified by an ID (see AddressMappings). 

TimelineTimes: which timeline times of the collection's approvedTransfersTimeline field can be updated or not?
permitted/forbidden TimelineTimes: when can the manager execute this permission?

Ex: Let's say we are updating the transferability for timelineTime 1 and the transfer tuple ("All", "All", "All", 10, 1000).
We would check to find the FIRST CollectionApprovedTransferPermission that matches this combination.
If we find a match, we would check the permitted/forbidden times to see if we can execute this permission (default is ALLOWED).

Ex: So if you wanted to freeze the transferability to enforce that badge ID 1 will always be transferable, you could set
the combination ("All", "All", "All", "All Transfer Times", 1) to always be forbidden at all timelineTimes.
*/
export interface BadgesCollectionApprovedTransferPermission {
  defaultValues?: BadgesCollectionApprovedTransferDefaultValues;
  combinations?: BadgesCollectionApprovedTransferCombination[];
}

export interface BadgesCollectionApprovedTransferTimeline {
  approvedTransfers?: BadgesCollectionApprovedTransfer[];
  timelineTimes?: BadgesUintRange[];
}

/**
* This defines the metadata for the collection.
This should be interpreted according to the collection standard.
*/
export interface BadgesCollectionMetadata {
  uri?: string;
  customData?: string;
}

export interface BadgesCollectionMetadataTimeline {
  /**
   * This defines the metadata for the collection.
   * This should be interpreted according to the collection standard.
   */
  collectionMetadata?: BadgesCollectionMetadata;
  timelineTimes?: BadgesUintRange[];
}

/**
* CollectionPermissions defines the permissions for the collection (i.e. what the manager can and cannot do).

There are five types of permissions for a collection: ActionPermission, TimedUpdatePermission, TimedUpdateWithBadgeIdsPermission, BalancesActionPermission, and CollectionApprovedTransferPermission.

The permission type allows fine-grained access control for each action.
ActionPermission: defines when the manager can perform an action.
TimedUpdatePermission: defines when the manager can update a timeline-based field and what times of the timeline can be updated.
TimedUpdateWithBadgeIdsPermission: defines when the manager can update a timeline-based field for specific badges and what times of the timeline can be updated.
BalancesActionPermission: defines when the manager can perform an action for specific badges and specific badge ownership times.
CollectionApprovedTransferPermission: defines when the manager can update the transferability of the collection and what transfers can be updated vs locked

Note there are a few different times here which could get confusing:
- timelineTimes: the times when a timeline-based field is a specific value
- permitted/forbiddenTimes - the times that a permission can be performed
- transferTimes - the times that a transfer occurs
- ownedTimes - the times when a badge is owned by a user

The permitted/forbiddenTimes are used to determine when a permission can be executed.
Once a time is set to be permitted or forbidden, it is PERMANENT and cannot be changed.
If a time is not set to be permitted or forbidden, it is considered NEUTRAL and can be updated but is ALLOWED by default.

Each permission type has a defaultValues field and a combinations field.
The defaultValues field defines the default values for the permission which can be manipulated by the combinations field (to avoid unnecessary repetition).
Ex: We can have default value badgeIds = [1,2] and combinations = [{invertDefault: true, isAllowed: false}, {isAllowed: true}].
This would mean that badgeIds [1,2] are allowed but everything else is not allowed.

IMPORTANT: For all permissions, we ONLY take the first combination that matches. Any subsequent combinations are ignored. 
Ex: If we have defaultValues = {badgeIds: [1,2]} and combinations = [{isAllowed: true}, {isAllowed: false}].
This would mean that badgeIds [1,2] are allowed and the second combination is ignored.
*/
export interface BadgesCollectionPermissions {
  canDeleteCollection?: BadgesActionPermission[];
  canArchive?: BadgesTimedUpdatePermission[];
  canUpdateContractAddress?: BadgesTimedUpdatePermission[];
  canUpdateOffChainBalancesMetadata?: BadgesTimedUpdatePermission[];
  canUpdateStandards?: BadgesTimedUpdatePermission[];
  canUpdateCustomData?: BadgesTimedUpdatePermission[];
  canUpdateManager?: BadgesTimedUpdatePermission[];
  canUpdateCollectionMetadata?: BadgesTimedUpdatePermission[];
  canCreateMoreBadges?: BadgesBalancesActionPermission[];
  canUpdateBadgeMetadata?: BadgesTimedUpdateWithBadgeIdsPermission[];
  canUpdateInheritedBalances?: BadgesTimedUpdateWithBadgeIdsPermission[];
  canUpdateCollectionApprovedTransfers?: BadgesCollectionApprovedTransferPermission[];
}

export interface BadgesContractAddressTimeline {
  contractAddress?: string;
  timelineTimes?: BadgesUintRange[];
}

export interface BadgesCustomDataTimeline {
  customData?: string;
  timelineTimes?: BadgesUintRange[];
}

/**
* InheritedBalances are a powerful feature of the BitBadges module.
They allow a colllection to inherit the balances from another collection.
Ex: Badges from Collection A inherits the balances from badges from Collection B.

The badgeIds specified will inherit the balances from the parent collection and badges specified.
If the total number of parent badges == 1, then all the badgeIds will inherit the balance from that parent badge.
Otherwise, the total number of parent badges must equal the total number of badgeIds specified.
By total number, we mean the sum of the number of badgeIds in each UintRange.
*/
export interface BadgesInheritedBalance {
  badgeIds?: BadgesUintRange[];
  parentCollectionId?: string;
  parentBadgeIds?: BadgesUintRange[];
}

export interface BadgesInheritedBalancesTimeline {
  inheritedBalances?: BadgesInheritedBalance[];
  timelineTimes?: BadgesUintRange[];
}

export interface BadgesIsArchivedTimeline {
  isArchived?: boolean;
  timelineTimes?: BadgesUintRange[];
}

export interface BadgesIsCollectionTransferAllowed {
  invertFrom?: boolean;
  invertTo?: boolean;
  invertInitiatedBy?: boolean;
  invertTransferTimes?: boolean;
  invertBadgeIds?: boolean;
  isAllowed?: boolean;
}

export interface BadgesIsUserIncomingTransferAllowed {
  invertFrom?: boolean;
  invertInitiatedBy?: boolean;
  invertTransferTimes?: boolean;
  invertBadgeIds?: boolean;
  isAllowed?: boolean;
}

export interface BadgesIsUserOutgoingTransferAllowed {
  invertTo?: boolean;
  invertInitiatedBy?: boolean;
  invertTransferTimes?: boolean;
  invertBadgeIds?: boolean;
  isAllowed?: boolean;
}

export interface BadgesManagerTimeline {
  manager?: string;
  timelineTimes?: BadgesUintRange[];
}

export type BadgesMsgArchiveCollectionResponse = object;

export type BadgesMsgDeleteCollectionResponse = object;

export type BadgesMsgMintAndDistributeBadgesResponse = object;

export interface BadgesMsgNewCollectionResponse {
  /** ID of created badge collecon */
  collectionId?: string;
}

export type BadgesMsgTransferBadgesResponse = object;

export type BadgesMsgUpdateCollectionApprovedTransfersResponse = object;

export type BadgesMsgUpdateCollectionPermissionsResponse = object;

export type BadgesMsgUpdateManagerResponse = object;

export type BadgesMsgUpdateMetadataResponse = object;

export type BadgesMsgUpdateUserApprovedTransfersResponse = object;

export type BadgesMsgUpdateUserPermissionsResponse = object;

/**
* This defines the metadata for the off-chain balances (if using this balances type).
This should be interpreted according to the collection standard.
*/
export interface BadgesOffChainBalancesMetadata {
  uri?: string;
  customData?: string;
}

export interface BadgesOffChainBalancesMetadataTimeline {
  /**
   * This defines the metadata for the off-chain balances (if using this balances type).
   * This should be interpreted according to the collection standard.
   */
  offChainBalancesMetadata?: BadgesOffChainBalancesMetadata;
  timelineTimes?: BadgesUintRange[];
}

/**
 * Params defines the parameters for the module.
 */
export type BadgesParams = object;

/**
 * PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address.
 */
export interface BadgesPerAddressApprovals {
  approvalsPerFromAddress?: BadgesApprovalsTracker;
  approvalsPerToAddress?: BadgesApprovalsTracker;
  approvalsPerInitiatedByAddress?: BadgesApprovalsTracker;
}

export interface BadgesQueryGetAddressMappingResponse {
  /**
   * An AddressMapping is a permanent list of addresses that are referenced by a mapping ID.
   * The mapping may include only the specified addresses, or it may include all addresses but
   * the specified addresses (depending on if includeAddresses is true or false).
   *
   * AddressMappings are used for things like whitelists, blacklists, approvals, etc.
   */
  mapping?: BadgesAddressMapping;
}

export interface BadgesQueryGetApprovalsTrackerResponse {
  tracker?: BadgesApprovalsTracker;
}

export interface BadgesQueryGetBalanceResponse {
  /**
   * UserBalanceStore is the store for the user balances
   * It consists of a list of balances, a list of approved outgoing transfers, and a list of approved incoming transfers,
   * and the permissions for updating the approved incoming/outgoing transfers.
   *
   * The default approved outgoing / incoming transfers are defined by the collection.
   * The outgoing transfers can be used to allow / disallow transfers which are sent from this user.
   * If a transfer has no match, then it is disallowed by default, unless from == initiatedBy (i.e. initiated by this user).
   * The incoming transfers can be used to allow / disallow transfers which are sent to this user.
   * If a transfer has no match, then it is disallowed by default, unless to == initiatedBy (i.e. initiated by this user).
   * Note that the user approved transfers are only checked if the collection approved transfers do not specify to override
   * the user approved transfers.
   */
  balance?: BadgesUserBalanceStore;
}

export interface BadgesQueryGetCollectionResponse {
  /**
   * A BadgeCollection is the top level object for a collection of badges.
   * It defines everything about the collection, such as the manager, metadata, etc.
   *
   * All collections are identified by a collectionId assigned by the blockchain, which is a uint64 that increments (i.e. first collection has ID 1).
   * All collections also have a manager who is responsible for managing the collection.
   * They can be granted certain permissions, such as the ability to mint new badges.
   * Certain fields are timeline-based, which means they may have different values at different block heights.
   * We fetch the value according to the current time.
   * For example, we may set the manager to be Alice from Time1 to Time2, and then set the manager to be Bob from Time2 to Time3.
   * Collections may have different balance types: standard vs off-chain vs inherited. See documentation for differences.
   */
  collection?: BadgesBadgeCollection;
}

export interface BadgesQueryGetNumUsedForMerkleChallengeResponse {
  numUsed?: string;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface BadgesQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: BadgesParams;
}

export interface BadgesStandardsTimeline {
  standards?: string[];
  timelineTimes?: BadgesUintRange[];
}

export interface BadgesTimedUpdateCombination {
  /** ValueOptions defines how we manipulate the default values. */
  timelineTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesTimedUpdateDefaultValues {
  timelineTimes?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* TimedUpdatePermission defines the permissions for updating a timeline-based field.

Ex: If you want to lock the ability to update the collection's metadata for timelineTimes 1/1/2020 - 1/1/2021,
you could set the combination (TimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
*/
export interface BadgesTimedUpdatePermission {
  defaultValues?: BadgesTimedUpdateDefaultValues;
  combinations?: BadgesTimedUpdateCombination[];
}

export interface BadgesTimedUpdateWithBadgeIdsCombination {
  /** ValueOptions defines how we manipulate the default values. */
  timelineTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  badgeIdsOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesTimedUpdateWithBadgeIdsDefaultValues {
  badgeIds?: BadgesUintRange[];
  timelineTimes?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* TimedUpdateWithBadgeIdsPermission defines the permissions for updating a timeline-based field for specific badges.

Ex: If you want to lock the ability to update the metadata for badgeIds [1,2] for timelineTimes 1/1/2020 - 1/1/2021,
you could set the combination (badgeIds: [1,2], TimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
*/
export interface BadgesTimedUpdateWithBadgeIdsPermission {
  defaultValues?: BadgesTimedUpdateWithBadgeIdsDefaultValues;
  combinations?: BadgesTimedUpdateWithBadgeIdsCombination[];
}

export interface BadgesTransfer {
  from?: string;
  toAddresses?: string[];
  balances?: BadgesBalance[];

  /**
   * Note here we remain optimistic that the solutions will apply to all potential challenges.
   * It is the Tx Sender's responsibility to ensure that the solutions are valid for all potential challenges.
   * If you are attempting to claim badges with different sets of challenges, you will need to make multiple transfers.
   */
  solutions?: BadgesMerkleProof[];
}

/**
* uintRange is a range of IDs from some start to some end (inclusive).

uintRanges are one of the core types used in the BitBadgesChain module.
They are used for evrything from badge IDs to time ranges to min / max balance amounts.
*/
export interface BadgesUintRange {
  start?: string;
  end?: string;
}

/**
* UserApprovedIncomingTransfer defines the rules for the approval of an incoming transfer to a user.
See CollectionApprovedTransfer for more details. This is the same minus a few fields.
*/
export interface BadgesUserApprovedIncomingTransfer {
  fromMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  allowedCombinations?: BadgesIsUserIncomingTransferAllowed[];
  challenges?: BadgesChallenge[];
  approvalId?: string;
  incrementBadgeIdsBy?: string;
  incrementOwnedTimesBy?: string;

  /** PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address. */
  perAddressApprovals?: BadgesPerAddressApprovals;
  uri?: string;
  customData?: string;
  requireFromEqualsInitiatedBy?: boolean;
  requireFromDoesNotEqualInitiatedBy?: boolean;
}

export interface BadgesUserApprovedIncomingTransferCombination {
  /** ValueOptions defines how we manipulate the default values. */
  timelineTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  fromMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  initiatedByMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  transferTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  badgeIdsOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesUserApprovedIncomingTransferDefaultValues {
  timelineTimes?: BadgesUintRange[];
  fromMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* UserApprovedIncomingTransferPermission defines the permissions for updating the user's approved incoming transfers.
See CollectionApprovedTransferPermission for more details. This is equivalent without the toMappingId field because that is always the user.
*/
export interface BadgesUserApprovedIncomingTransferPermission {
  defaultValues?: BadgesUserApprovedIncomingTransferDefaultValues;
  combinations?: BadgesUserApprovedIncomingTransferCombination[];
}

export interface BadgesUserApprovedIncomingTransferTimeline {
  approvedIncomingTransfers?: BadgesUserApprovedIncomingTransfer[];
  timelineTimes?: BadgesUintRange[];
}

/**
* UserApprovedOutgoingTransfer defines the rules for the approval of an outgoing transfer from a user.
See CollectionApprovedTransfer for more details. This is the same minus a few fields.
*/
export interface BadgesUserApprovedOutgoingTransfer {
  toMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  allowedCombinations?: BadgesIsUserOutgoingTransferAllowed[];
  challenges?: BadgesChallenge[];
  approvalId?: string;
  incrementBadgeIdsBy?: string;
  incrementOwnedTimesBy?: string;

  /** PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address. */
  perAddressApprovals?: BadgesPerAddressApprovals;
  uri?: string;
  customData?: string;
  requireToEqualsInitiatedBy?: boolean;
  requireToDoesNotEqualInitiatedBy?: boolean;
}

export interface BadgesUserApprovedOutgoingTransferCombination {
  /** ValueOptions defines how we manipulate the default values. */
  timelineTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  toMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  initiatedByMappingOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  transferTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  badgeIdsOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  permittedTimesOptions?: BadgesValueOptions;

  /** ValueOptions defines how we manipulate the default values. */
  forbiddenTimesOptions?: BadgesValueOptions;
}

export interface BadgesUserApprovedOutgoingTransferDefaultValues {
  timelineTimes?: BadgesUintRange[];
  toMappingId?: string;
  initiatedByMappingId?: string;
  transferTimes?: BadgesUintRange[];
  badgeIds?: BadgesUintRange[];
  permittedTimes?: BadgesUintRange[];
  forbiddenTimes?: BadgesUintRange[];
}

/**
* UserApprovedOutgoingTransferPermission defines the permissions for updating the user's approved outgoing transfers.
See CollectionApprovedTransferPermission for more details. This is equivalent without the fromMappingId field because that is always the user.
*/
export interface BadgesUserApprovedOutgoingTransferPermission {
  defaultValues?: BadgesUserApprovedOutgoingTransferDefaultValues;
  combinations?: BadgesUserApprovedOutgoingTransferCombination[];
}

export interface BadgesUserApprovedOutgoingTransferTimeline {
  approvedOutgoingTransfers?: BadgesUserApprovedOutgoingTransfer[];
  timelineTimes?: BadgesUintRange[];
}

/**
* UserBalanceStore is the store for the user balances
It consists of a list of balances, a list of approved outgoing transfers, and a list of approved incoming transfers,
and the permissions for updating the approved incoming/outgoing transfers.

The default approved outgoing / incoming transfers are defined by the collection.

The outgoing transfers can be used to allow / disallow transfers which are sent from this user.
If a transfer has no match, then it is disallowed by default, unless from == initiatedBy (i.e. initiated by this user).

The incoming transfers can be used to allow / disallow transfers which are sent to this user.
If a transfer has no match, then it is disallowed by default, unless to == initiatedBy (i.e. initiated by this user).

Note that the user approved transfers are only checked if the collection approved transfers do not specify to override
the user approved transfers.
*/
export interface BadgesUserBalanceStore {
  balances?: BadgesBalance[];
  approvedOutgoingTransfersTimeline?: BadgesUserApprovedOutgoingTransferTimeline[];
  approvedIncomingTransfersTimeline?: BadgesUserApprovedIncomingTransferTimeline[];

  /**
   * UserPermissions defines the permissions for the user (i.e. what the user can and cannot do).
   *
   * See CollectionPermissions for more details on the different types of permissions.
   * The UserApprovedOutgoing and UserApprovedIncoming permissions are the same as the CollectionApprovedTransferPermission,
   * but certain fields are removed because they are not relevant to the user.
   */
  permissions?: BadgesUserPermissions;
}

/**
* UserPermissions defines the permissions for the user (i.e. what the user can and cannot do).

See CollectionPermissions for more details on the different types of permissions.
The UserApprovedOutgoing and UserApprovedIncoming permissions are the same as the CollectionApprovedTransferPermission,
but certain fields are removed because they are not relevant to the user.
*/
export interface BadgesUserPermissions {
  canUpdateApprovedOutgoingTransfers?: BadgesUserApprovedOutgoingTransferPermission[];
  canUpdateApprovedIncomingTransfers?: BadgesUserApprovedIncomingTransferPermission[];
}

/**
 * ValueOptions defines how we manipulate the default values.
 */
export interface BadgesValueOptions {
  invertDefault?: boolean;

  /** Override default values with all possible values */
  allValues?: boolean;

  /** Override default values with no values */
  noValues?: boolean;
}

/**
* `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

 Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }
*/
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
            ? JSON.stringify(property)
            : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title badges/address_mappings.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryGetAddressMapping
   * @request GET:/bitbadges/bitbadgeschain/badges/get_address_mapping/{mappingId}
   */
  queryGetAddressMapping = (mappingId: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetAddressMappingResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_address_mapping/${mappingId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGetApprovalsTracker
   * @request GET:/bitbadges/bitbadgeschain/badges/get_approvals_tracker/{approvalId}/{level}
   */
  queryGetApprovalsTracker = (
    approvalId: string,
    level: string,
    query?: { depth?: string; address?: string; collectionId?: string },
    params: RequestParams = {},
  ) =>
    this.request<BadgesQueryGetApprovalsTrackerResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_approvals_tracker/${approvalId}/${level}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGetBalance
   * @summary Queries an addresses balance for a badge collection, specified by its ID.
   * @request GET:/bitbadges/bitbadgeschain/badges/get_balance/{collectionId}/{address}
   */
  queryGetBalance = (collectionId: string, address: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetBalanceResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_balance/${collectionId}/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGetCollection
   * @summary Queries a badge collection by ID.
   * @request GET:/bitbadges/bitbadgeschain/badges/get_collection/{collectionId}
   */
  queryGetCollection = (collectionId: string, params: RequestParams = {}) =>
    this.request<BadgesQueryGetCollectionResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_collection/${collectionId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGetNumUsedForMerkleChallenge
   * @request GET:/bitbadges/bitbadgeschain/badges/get_num_used_for_challenge/{challengeId}/{level}/{leafIndex}
   */
  queryGetNumUsedForMerkleChallenge = (
    challengeId: string,
    level: string,
    leafIndex: string,
    query?: { collectionId?: string },
    params: RequestParams = {},
  ) =>
    this.request<BadgesQueryGetNumUsedForMerkleChallengeResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/get_num_used_for_challenge/${challengeId}/${level}/${leafIndex}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/bitbadges/bitbadgeschain/badges/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<BadgesQueryParamsResponse, RpcStatus>({
      path: `/bitbadges/bitbadgeschain/badges/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
