/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { UintRange } from "./balances";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/**
 * CollectionPermissions defines the permissions for the collection (i.e. what the manager can and cannot do).
 *
 * There are five types of permissions for a collection: ActionPermission, TimedUpdatePermission, TimedUpdateWithBadgeIdsPermission, BalancesActionPermission, and CollectionApprovedTransferPermission.
 *
 * The permission type allows fine-grained access control for each action.
 * ActionPermission: defines when the manager can perform an action.
 * TimedUpdatePermission: defines when the manager can update a timeline-based field and what times of the timeline can be updated.
 * TimedUpdateWithBadgeIdsPermission: defines when the manager can update a timeline-based field for specific badges and what times of the timeline can be updated.
 * BalancesActionPermission: defines when the manager can perform an action for specific badges and specific badge ownership times.
 * CollectionApprovedTransferPermission: defines when the manager can update the transferability of the collection and what transfers can be updated vs locked
 *
 * Note there are a few different times here which could get confusing:
 * - timelineTimes: the times when a timeline-based field is a specific value
 * - permitted/forbiddenTimes - the times that a permission can be performed
 * - transferTimes - the times that a transfer occurs
 * - ownedTimes - the times when a badge is owned by a user
 *
 * The permitted/forbiddenTimes are used to determine when a permission can be executed.
 * Once a time is set to be permitted or forbidden, it is PERMANENT and cannot be changed.
 * If a time is not set to be permitted or forbidden, it is considered NEUTRAL and can be updated but is ALLOWED by default.
 *
 * Each permission type has a defaultValues field and a combinations field.
 * The defaultValues field defines the default values for the permission which can be manipulated by the combinations field (to avoid unnecessary repetition).
 * Ex: We can have default value badgeIds = [1,2] and combinations = [{invertDefault: true, isAllowed: false}, {isAllowed: true}].
 * This would mean that badgeIds [1,2] are allowed but everything else is not allowed.
 *
 * IMPORTANT: For all permissions, we ONLY take the first combination that matches. Any subsequent combinations are ignored.
 * Ex: If we have defaultValues = {badgeIds: [1,2]} and combinations = [{isAllowed: true}, {isAllowed: false}].
 * This would mean that badgeIds [1,2] are allowed and the second combination is ignored.
 */
export interface CollectionPermissions {
  canDeleteCollection: ActionPermission[];
  canArchive: TimedUpdatePermission[];
  canUpdateContractAddress: TimedUpdatePermission[];
  canUpdateOffChainBalancesMetadata: TimedUpdatePermission[];
  canUpdateStandards: TimedUpdatePermission[];
  canUpdateCustomData: TimedUpdatePermission[];
  canUpdateManager: TimedUpdatePermission[];
  canUpdateCollectionMetadata: TimedUpdatePermission[];
  canCreateMoreBadges: BalancesActionPermission[];
  canUpdateBadgeMetadata: TimedUpdateWithBadgeIdsPermission[];
  canUpdateInheritedBalances: TimedUpdateWithBadgeIdsPermission[];
  canUpdateCollectionApprovedTransfers: CollectionApprovedTransferPermission[];
}

/**
 * UserPermissions defines the permissions for the user (i.e. what the user can and cannot do).
 *
 * See CollectionPermissions for more details on the different types of permissions.
 * The UserApprovedOutgoing and UserApprovedIncoming permissions are the same as the CollectionApprovedTransferPermission,
 * but certain fields are removed because they are not relevant to the user.
 */
export interface UserPermissions {
  canUpdateApprovedOutgoingTransfers: UserApprovedOutgoingTransferPermission[];
  canUpdateApprovedIncomingTransfers: UserApprovedIncomingTransferPermission[];
}

/** ValueOptions defines how we manipulate the default values. */
export interface ValueOptions {
  invertDefault: boolean;
  /** Override default values with all possible values */
  allValues: boolean;
  /** Override default values with no values */
  noValues: boolean;
}

export interface CollectionApprovedTransferCombination {
  timelineTimesOptions: ValueOptions | undefined;
  fromMappingOptions: ValueOptions | undefined;
  toMappingOptions: ValueOptions | undefined;
  initiatedByMappingOptions: ValueOptions | undefined;
  transferTimesOptions: ValueOptions | undefined;
  badgeIdsOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface CollectionApprovedTransferDefaultValues {
  timelineTimes: UintRange[];
  fromMappingId: string;
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * CollectionApprovedTransferPermission defines what collection approved transfers can be updated vs are locked.
 *
 * Each transfer is broken down to a (from, to, initiatedBy, transferTime, badgeId) tuple.
 * For a transfer to match, we need to match ALL of the fields in the combination.
 * These are detemined by the fromMappingId, toMappingId, initiatedByMappingId, transferTimes, badgeIds fields.
 * AddressMappings are used for (from, to, initiatedBy) which are a permanent list of addresses identified by an ID (see AddressMappings).
 *
 * TimelineTimes: which timeline times of the collection's approvedTransfersTimeline field can be updated or not?
 * permitted/forbidden TimelineTimes: when can the manager execute this permission?
 *
 * Ex: Let's say we are updating the transferability for timelineTime 1 and the transfer tuple ("All", "All", "All", 10, 1000).
 * We would check to find the FIRST CollectionApprovedTransferPermission that matches this combination.
 * If we find a match, we would check the permitted/forbidden times to see if we can execute this permission (default is ALLOWED).
 *
 * Ex: So if you wanted to freeze the transferability to enforce that badge ID 1 will always be transferable, you could set
 * the combination ("All", "All", "All", "All Transfer Times", 1) to always be forbidden at all timelineTimes.
 */
export interface CollectionApprovedTransferPermission {
  defaultValues: CollectionApprovedTransferDefaultValues | undefined;
  combinations: CollectionApprovedTransferCombination[];
}

export interface UserApprovedOutgoingTransferCombination {
  timelineTimesOptions: ValueOptions | undefined;
  toMappingOptions: ValueOptions | undefined;
  initiatedByMappingOptions: ValueOptions | undefined;
  transferTimesOptions: ValueOptions | undefined;
  badgeIdsOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface UserApprovedOutgoingTransferDefaultValues {
  timelineTimes: UintRange[];
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * UserApprovedOutgoingTransferPermission defines the permissions for updating the user's approved outgoing transfers.
 * See CollectionApprovedTransferPermission for more details. This is equivalent without the fromMappingId field because that is always the user.
 */
export interface UserApprovedOutgoingTransferPermission {
  defaultValues: UserApprovedOutgoingTransferDefaultValues | undefined;
  combinations: UserApprovedOutgoingTransferCombination[];
}

export interface UserApprovedIncomingTransferCombination {
  timelineTimesOptions: ValueOptions | undefined;
  fromMappingOptions: ValueOptions | undefined;
  initiatedByMappingOptions: ValueOptions | undefined;
  transferTimesOptions: ValueOptions | undefined;
  badgeIdsOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface UserApprovedIncomingTransferDefaultValues {
  timelineTimes: UintRange[];
  fromMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * UserApprovedIncomingTransferPermission defines the permissions for updating the user's approved incoming transfers.
 * See CollectionApprovedTransferPermission for more details. This is equivalent without the toMappingId field because that is always the user.
 */
export interface UserApprovedIncomingTransferPermission {
  defaultValues: UserApprovedIncomingTransferDefaultValues | undefined;
  combinations: UserApprovedIncomingTransferCombination[];
}

export interface BalancesActionCombination {
  badgeIdsOptions: ValueOptions | undefined;
  ownedTimesOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface BalancesActionDefaultValues {
  badgeIds: UintRange[];
  ownedTimes: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * BalancesActionPermission defines the permissions for updating a timeline-based field for specific badges and specific badge ownership times.
 * Currently, this is only used for creating new badges.
 *
 * Ex: If you want to lock the ability to create new badges for badgeIds [1,2] at ownedTimes 1/1/2020 - 1/1/2021,
 * you could set the combination (badgeIds: [1,2], ownershipTimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
 */
export interface BalancesActionPermission {
  defaultValues: BalancesActionDefaultValues | undefined;
  combinations: BalancesActionCombination[];
}

export interface ActionDefaultValues {
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

export interface ActionCombination {
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

/**
 * ActionPermission defines the permissions for performing an action.
 *
 * This is simple and straightforward as the only thing we need to check is the permitted/forbidden times.
 */
export interface ActionPermission {
  defaultValues: ActionDefaultValues | undefined;
  combinations: ActionCombination[];
}

export interface TimedUpdateCombination {
  timelineTimesOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface TimedUpdateDefaultValues {
  timelineTimes: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * TimedUpdatePermission defines the permissions for updating a timeline-based field.
 *
 * Ex: If you want to lock the ability to update the collection's metadata for timelineTimes 1/1/2020 - 1/1/2021,
 * you could set the combination (TimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
 */
export interface TimedUpdatePermission {
  defaultValues: TimedUpdateDefaultValues | undefined;
  combinations: TimedUpdateCombination[];
}

export interface TimedUpdateWithBadgeIdsCombination {
  timelineTimesOptions: ValueOptions | undefined;
  badgeIdsOptions: ValueOptions | undefined;
  permittedTimesOptions: ValueOptions | undefined;
  forbiddenTimesOptions: ValueOptions | undefined;
}

export interface TimedUpdateWithBadgeIdsDefaultValues {
  badgeIds: UintRange[];
  timelineTimes: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * TimedUpdateWithBadgeIdsPermission defines the permissions for updating a timeline-based field for specific badges.
 *
 * Ex: If you want to lock the ability to update the metadata for badgeIds [1,2] for timelineTimes 1/1/2020 - 1/1/2021,
 * you could set the combination (badgeIds: [1,2], TimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
 */
export interface TimedUpdateWithBadgeIdsPermission {
  defaultValues: TimedUpdateWithBadgeIdsDefaultValues | undefined;
  combinations: TimedUpdateWithBadgeIdsCombination[];
}

function createBaseCollectionPermissions(): CollectionPermissions {
  return {
    canDeleteCollection: [],
    canArchive: [],
    canUpdateContractAddress: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canCreateMoreBadges: [],
    canUpdateBadgeMetadata: [],
    canUpdateInheritedBalances: [],
    canUpdateCollectionApprovedTransfers: [],
  };
}

export const CollectionPermissions = {
  encode(message: CollectionPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.canDeleteCollection) {
      ActionPermission.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.canArchive) {
      TimedUpdatePermission.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.canUpdateContractAddress) {
      TimedUpdatePermission.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.canUpdateOffChainBalancesMetadata) {
      TimedUpdatePermission.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.canUpdateStandards) {
      TimedUpdatePermission.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.canUpdateCustomData) {
      TimedUpdatePermission.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.canUpdateManager) {
      TimedUpdatePermission.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.canUpdateCollectionMetadata) {
      TimedUpdatePermission.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.canCreateMoreBadges) {
      BalancesActionPermission.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.canUpdateBadgeMetadata) {
      TimedUpdateWithBadgeIdsPermission.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.canUpdateInheritedBalances) {
      TimedUpdateWithBadgeIdsPermission.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.canUpdateCollectionApprovedTransfers) {
      CollectionApprovedTransferPermission.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionPermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionPermissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.canDeleteCollection.push(ActionPermission.decode(reader, reader.uint32()));
          break;
        case 2:
          message.canArchive.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 3:
          message.canUpdateContractAddress.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 4:
          message.canUpdateOffChainBalancesMetadata.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 5:
          message.canUpdateStandards.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 6:
          message.canUpdateCustomData.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 7:
          message.canUpdateManager.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 8:
          message.canUpdateCollectionMetadata.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 9:
          message.canCreateMoreBadges.push(BalancesActionPermission.decode(reader, reader.uint32()));
          break;
        case 10:
          message.canUpdateBadgeMetadata.push(TimedUpdateWithBadgeIdsPermission.decode(reader, reader.uint32()));
          break;
        case 11:
          message.canUpdateInheritedBalances.push(TimedUpdateWithBadgeIdsPermission.decode(reader, reader.uint32()));
          break;
        case 12:
          message.canUpdateCollectionApprovedTransfers.push(
            CollectionApprovedTransferPermission.decode(reader, reader.uint32()),
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionPermissions {
    return {
      canDeleteCollection: Array.isArray(object?.canDeleteCollection)
        ? object.canDeleteCollection.map((e: any) => ActionPermission.fromJSON(e))
        : [],
      canArchive: Array.isArray(object?.canArchive)
        ? object.canArchive.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateContractAddress: Array.isArray(object?.canUpdateContractAddress)
        ? object.canUpdateContractAddress.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateOffChainBalancesMetadata: Array.isArray(object?.canUpdateOffChainBalancesMetadata)
        ? object.canUpdateOffChainBalancesMetadata.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateStandards: Array.isArray(object?.canUpdateStandards)
        ? object.canUpdateStandards.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateCustomData: Array.isArray(object?.canUpdateCustomData)
        ? object.canUpdateCustomData.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateManager: Array.isArray(object?.canUpdateManager)
        ? object.canUpdateManager.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canUpdateCollectionMetadata: Array.isArray(object?.canUpdateCollectionMetadata)
        ? object.canUpdateCollectionMetadata.map((e: any) => TimedUpdatePermission.fromJSON(e))
        : [],
      canCreateMoreBadges: Array.isArray(object?.canCreateMoreBadges)
        ? object.canCreateMoreBadges.map((e: any) => BalancesActionPermission.fromJSON(e))
        : [],
      canUpdateBadgeMetadata: Array.isArray(object?.canUpdateBadgeMetadata)
        ? object.canUpdateBadgeMetadata.map((e: any) => TimedUpdateWithBadgeIdsPermission.fromJSON(e))
        : [],
      canUpdateInheritedBalances: Array.isArray(object?.canUpdateInheritedBalances)
        ? object.canUpdateInheritedBalances.map((e: any) => TimedUpdateWithBadgeIdsPermission.fromJSON(e))
        : [],
      canUpdateCollectionApprovedTransfers: Array.isArray(object?.canUpdateCollectionApprovedTransfers)
        ? object.canUpdateCollectionApprovedTransfers.map((e: any) => CollectionApprovedTransferPermission.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionPermissions): unknown {
    const obj: any = {};
    if (message.canDeleteCollection) {
      obj.canDeleteCollection = message.canDeleteCollection.map((e) => e ? ActionPermission.toJSON(e) : undefined);
    } else {
      obj.canDeleteCollection = [];
    }
    if (message.canArchive) {
      obj.canArchive = message.canArchive.map((e) => e ? TimedUpdatePermission.toJSON(e) : undefined);
    } else {
      obj.canArchive = [];
    }
    if (message.canUpdateContractAddress) {
      obj.canUpdateContractAddress = message.canUpdateContractAddress.map((e) =>
        e ? TimedUpdatePermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateContractAddress = [];
    }
    if (message.canUpdateOffChainBalancesMetadata) {
      obj.canUpdateOffChainBalancesMetadata = message.canUpdateOffChainBalancesMetadata.map((e) =>
        e ? TimedUpdatePermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateOffChainBalancesMetadata = [];
    }
    if (message.canUpdateStandards) {
      obj.canUpdateStandards = message.canUpdateStandards.map((e) => e ? TimedUpdatePermission.toJSON(e) : undefined);
    } else {
      obj.canUpdateStandards = [];
    }
    if (message.canUpdateCustomData) {
      obj.canUpdateCustomData = message.canUpdateCustomData.map((e) => e ? TimedUpdatePermission.toJSON(e) : undefined);
    } else {
      obj.canUpdateCustomData = [];
    }
    if (message.canUpdateManager) {
      obj.canUpdateManager = message.canUpdateManager.map((e) => e ? TimedUpdatePermission.toJSON(e) : undefined);
    } else {
      obj.canUpdateManager = [];
    }
    if (message.canUpdateCollectionMetadata) {
      obj.canUpdateCollectionMetadata = message.canUpdateCollectionMetadata.map((e) =>
        e ? TimedUpdatePermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateCollectionMetadata = [];
    }
    if (message.canCreateMoreBadges) {
      obj.canCreateMoreBadges = message.canCreateMoreBadges.map((e) =>
        e ? BalancesActionPermission.toJSON(e) : undefined
      );
    } else {
      obj.canCreateMoreBadges = [];
    }
    if (message.canUpdateBadgeMetadata) {
      obj.canUpdateBadgeMetadata = message.canUpdateBadgeMetadata.map((e) =>
        e ? TimedUpdateWithBadgeIdsPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateBadgeMetadata = [];
    }
    if (message.canUpdateInheritedBalances) {
      obj.canUpdateInheritedBalances = message.canUpdateInheritedBalances.map((e) =>
        e ? TimedUpdateWithBadgeIdsPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateInheritedBalances = [];
    }
    if (message.canUpdateCollectionApprovedTransfers) {
      obj.canUpdateCollectionApprovedTransfers = message.canUpdateCollectionApprovedTransfers.map((e) =>
        e ? CollectionApprovedTransferPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateCollectionApprovedTransfers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionPermissions>, I>>(object: I): CollectionPermissions {
    const message = createBaseCollectionPermissions();
    message.canDeleteCollection = object.canDeleteCollection?.map((e) => ActionPermission.fromPartial(e)) || [];
    message.canArchive = object.canArchive?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canUpdateContractAddress = object.canUpdateContractAddress?.map((e) => TimedUpdatePermission.fromPartial(e))
      || [];
    message.canUpdateOffChainBalancesMetadata =
      object.canUpdateOffChainBalancesMetadata?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canUpdateStandards = object.canUpdateStandards?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canUpdateCustomData = object.canUpdateCustomData?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canUpdateManager = object.canUpdateManager?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canUpdateCollectionMetadata =
      object.canUpdateCollectionMetadata?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
    message.canCreateMoreBadges = object.canCreateMoreBadges?.map((e) => BalancesActionPermission.fromPartial(e)) || [];
    message.canUpdateBadgeMetadata =
      object.canUpdateBadgeMetadata?.map((e) => TimedUpdateWithBadgeIdsPermission.fromPartial(e)) || [];
    message.canUpdateInheritedBalances =
      object.canUpdateInheritedBalances?.map((e) => TimedUpdateWithBadgeIdsPermission.fromPartial(e)) || [];
    message.canUpdateCollectionApprovedTransfers =
      object.canUpdateCollectionApprovedTransfers?.map((e) => CollectionApprovedTransferPermission.fromPartial(e))
      || [];
    return message;
  },
};

function createBaseUserPermissions(): UserPermissions {
  return { canUpdateApprovedOutgoingTransfers: [], canUpdateApprovedIncomingTransfers: [] };
}

export const UserPermissions = {
  encode(message: UserPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.canUpdateApprovedOutgoingTransfers) {
      UserApprovedOutgoingTransferPermission.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.canUpdateApprovedIncomingTransfers) {
      UserApprovedIncomingTransferPermission.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserPermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserPermissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.canUpdateApprovedOutgoingTransfers.push(
            UserApprovedOutgoingTransferPermission.decode(reader, reader.uint32()),
          );
          break;
        case 2:
          message.canUpdateApprovedIncomingTransfers.push(
            UserApprovedIncomingTransferPermission.decode(reader, reader.uint32()),
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserPermissions {
    return {
      canUpdateApprovedOutgoingTransfers: Array.isArray(object?.canUpdateApprovedOutgoingTransfers)
        ? object.canUpdateApprovedOutgoingTransfers.map((e: any) => UserApprovedOutgoingTransferPermission.fromJSON(e))
        : [],
      canUpdateApprovedIncomingTransfers: Array.isArray(object?.canUpdateApprovedIncomingTransfers)
        ? object.canUpdateApprovedIncomingTransfers.map((e: any) => UserApprovedIncomingTransferPermission.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserPermissions): unknown {
    const obj: any = {};
    if (message.canUpdateApprovedOutgoingTransfers) {
      obj.canUpdateApprovedOutgoingTransfers = message.canUpdateApprovedOutgoingTransfers.map((e) =>
        e ? UserApprovedOutgoingTransferPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateApprovedOutgoingTransfers = [];
    }
    if (message.canUpdateApprovedIncomingTransfers) {
      obj.canUpdateApprovedIncomingTransfers = message.canUpdateApprovedIncomingTransfers.map((e) =>
        e ? UserApprovedIncomingTransferPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateApprovedIncomingTransfers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserPermissions>, I>>(object: I): UserPermissions {
    const message = createBaseUserPermissions();
    message.canUpdateApprovedOutgoingTransfers =
      object.canUpdateApprovedOutgoingTransfers?.map((e) => UserApprovedOutgoingTransferPermission.fromPartial(e))
      || [];
    message.canUpdateApprovedIncomingTransfers =
      object.canUpdateApprovedIncomingTransfers?.map((e) => UserApprovedIncomingTransferPermission.fromPartial(e))
      || [];
    return message;
  },
};

function createBaseValueOptions(): ValueOptions {
  return { invertDefault: false, allValues: false, noValues: false };
}

export const ValueOptions = {
  encode(message: ValueOptions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.invertDefault === true) {
      writer.uint32(8).bool(message.invertDefault);
    }
    if (message.allValues === true) {
      writer.uint32(16).bool(message.allValues);
    }
    if (message.noValues === true) {
      writer.uint32(24).bool(message.noValues);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValueOptions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValueOptions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invertDefault = reader.bool();
          break;
        case 2:
          message.allValues = reader.bool();
          break;
        case 3:
          message.noValues = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ValueOptions {
    return {
      invertDefault: isSet(object.invertDefault) ? Boolean(object.invertDefault) : false,
      allValues: isSet(object.allValues) ? Boolean(object.allValues) : false,
      noValues: isSet(object.noValues) ? Boolean(object.noValues) : false,
    };
  },

  toJSON(message: ValueOptions): unknown {
    const obj: any = {};
    message.invertDefault !== undefined && (obj.invertDefault = message.invertDefault);
    message.allValues !== undefined && (obj.allValues = message.allValues);
    message.noValues !== undefined && (obj.noValues = message.noValues);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ValueOptions>, I>>(object: I): ValueOptions {
    const message = createBaseValueOptions();
    message.invertDefault = object.invertDefault ?? false;
    message.allValues = object.allValues ?? false;
    message.noValues = object.noValues ?? false;
    return message;
  },
};

function createBaseCollectionApprovedTransferCombination(): CollectionApprovedTransferCombination {
  return {
    timelineTimesOptions: undefined,
    fromMappingOptions: undefined,
    toMappingOptions: undefined,
    initiatedByMappingOptions: undefined,
    transferTimesOptions: undefined,
    badgeIdsOptions: undefined,
    permittedTimesOptions: undefined,
    forbiddenTimesOptions: undefined,
  };
}

export const CollectionApprovedTransferCombination = {
  encode(message: CollectionApprovedTransferCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timelineTimesOptions !== undefined) {
      ValueOptions.encode(message.timelineTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.fromMappingOptions !== undefined) {
      ValueOptions.encode(message.fromMappingOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.toMappingOptions !== undefined) {
      ValueOptions.encode(message.toMappingOptions, writer.uint32(26).fork()).ldelim();
    }
    if (message.initiatedByMappingOptions !== undefined) {
      ValueOptions.encode(message.initiatedByMappingOptions, writer.uint32(34).fork()).ldelim();
    }
    if (message.transferTimesOptions !== undefined) {
      ValueOptions.encode(message.transferTimesOptions, writer.uint32(42).fork()).ldelim();
    }
    if (message.badgeIdsOptions !== undefined) {
      ValueOptions.encode(message.badgeIdsOptions, writer.uint32(50).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(58).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovedTransferCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovedTransferCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.fromMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.toMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 4:
          message.initiatedByMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 5:
          message.transferTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 6:
          message.badgeIdsOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 7:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 8:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovedTransferCombination {
    return {
      timelineTimesOptions: isSet(object.timelineTimesOptions)
        ? ValueOptions.fromJSON(object.timelineTimesOptions)
        : undefined,
      fromMappingOptions: isSet(object.fromMappingOptions)
        ? ValueOptions.fromJSON(object.fromMappingOptions)
        : undefined,
      toMappingOptions: isSet(object.toMappingOptions) ? ValueOptions.fromJSON(object.toMappingOptions) : undefined,
      initiatedByMappingOptions: isSet(object.initiatedByMappingOptions)
        ? ValueOptions.fromJSON(object.initiatedByMappingOptions)
        : undefined,
      transferTimesOptions: isSet(object.transferTimesOptions)
        ? ValueOptions.fromJSON(object.transferTimesOptions)
        : undefined,
      badgeIdsOptions: isSet(object.badgeIdsOptions) ? ValueOptions.fromJSON(object.badgeIdsOptions) : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: CollectionApprovedTransferCombination): unknown {
    const obj: any = {};
    message.timelineTimesOptions !== undefined && (obj.timelineTimesOptions = message.timelineTimesOptions
      ? ValueOptions.toJSON(message.timelineTimesOptions)
      : undefined);
    message.fromMappingOptions !== undefined && (obj.fromMappingOptions = message.fromMappingOptions
      ? ValueOptions.toJSON(message.fromMappingOptions)
      : undefined);
    message.toMappingOptions !== undefined
      && (obj.toMappingOptions = message.toMappingOptions ? ValueOptions.toJSON(message.toMappingOptions) : undefined);
    message.initiatedByMappingOptions !== undefined
      && (obj.initiatedByMappingOptions = message.initiatedByMappingOptions
        ? ValueOptions.toJSON(message.initiatedByMappingOptions)
        : undefined);
    message.transferTimesOptions !== undefined && (obj.transferTimesOptions = message.transferTimesOptions
      ? ValueOptions.toJSON(message.transferTimesOptions)
      : undefined);
    message.badgeIdsOptions !== undefined
      && (obj.badgeIdsOptions = message.badgeIdsOptions ? ValueOptions.toJSON(message.badgeIdsOptions) : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransferCombination>, I>>(
    object: I,
  ): CollectionApprovedTransferCombination {
    const message = createBaseCollectionApprovedTransferCombination();
    message.timelineTimesOptions = (object.timelineTimesOptions !== undefined && object.timelineTimesOptions !== null)
      ? ValueOptions.fromPartial(object.timelineTimesOptions)
      : undefined;
    message.fromMappingOptions = (object.fromMappingOptions !== undefined && object.fromMappingOptions !== null)
      ? ValueOptions.fromPartial(object.fromMappingOptions)
      : undefined;
    message.toMappingOptions = (object.toMappingOptions !== undefined && object.toMappingOptions !== null)
      ? ValueOptions.fromPartial(object.toMappingOptions)
      : undefined;
    message.initiatedByMappingOptions =
      (object.initiatedByMappingOptions !== undefined && object.initiatedByMappingOptions !== null)
        ? ValueOptions.fromPartial(object.initiatedByMappingOptions)
        : undefined;
    message.transferTimesOptions = (object.transferTimesOptions !== undefined && object.transferTimesOptions !== null)
      ? ValueOptions.fromPartial(object.transferTimesOptions)
      : undefined;
    message.badgeIdsOptions = (object.badgeIdsOptions !== undefined && object.badgeIdsOptions !== null)
      ? ValueOptions.fromPartial(object.badgeIdsOptions)
      : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseCollectionApprovedTransferDefaultValues(): CollectionApprovedTransferDefaultValues {
  return {
    timelineTimes: [],
    fromMappingId: "",
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const CollectionApprovedTransferDefaultValues = {
  encode(message: CollectionApprovedTransferDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.fromMappingId !== "") {
      writer.uint32(18).string(message.fromMappingId);
    }
    if (message.toMappingId !== "") {
      writer.uint32(26).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(34).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovedTransferDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovedTransferDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.fromMappingId = reader.string();
          break;
        case 3:
          message.toMappingId = reader.string();
          break;
        case 4:
          message.initiatedByMappingId = reader.string();
          break;
        case 5:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 8:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovedTransferDefaultValues {
    return {
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionApprovedTransferDefaultValues): unknown {
    const obj: any = {};
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    message.fromMappingId !== undefined && (obj.fromMappingId = message.fromMappingId);
    message.toMappingId !== undefined && (obj.toMappingId = message.toMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransferDefaultValues>, I>>(
    object: I,
  ): CollectionApprovedTransferDefaultValues {
    const message = createBaseCollectionApprovedTransferDefaultValues();
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.fromMappingId = object.fromMappingId ?? "";
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCollectionApprovedTransferPermission(): CollectionApprovedTransferPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const CollectionApprovedTransferPermission = {
  encode(message: CollectionApprovedTransferPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      CollectionApprovedTransferDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      CollectionApprovedTransferCombination.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovedTransferPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovedTransferPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = CollectionApprovedTransferDefaultValues.decode(reader, reader.uint32());
          break;
        case 7:
          message.combinations.push(CollectionApprovedTransferCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovedTransferPermission {
    return {
      defaultValues: isSet(object.defaultValues)
        ? CollectionApprovedTransferDefaultValues.fromJSON(object.defaultValues)
        : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => CollectionApprovedTransferCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionApprovedTransferPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined && (obj.defaultValues = message.defaultValues
      ? CollectionApprovedTransferDefaultValues.toJSON(message.defaultValues)
      : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) =>
        e ? CollectionApprovedTransferCombination.toJSON(e) : undefined
      );
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransferPermission>, I>>(
    object: I,
  ): CollectionApprovedTransferPermission {
    const message = createBaseCollectionApprovedTransferPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? CollectionApprovedTransferDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => CollectionApprovedTransferCombination.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserApprovedOutgoingTransferCombination(): UserApprovedOutgoingTransferCombination {
  return {
    timelineTimesOptions: undefined,
    toMappingOptions: undefined,
    initiatedByMappingOptions: undefined,
    transferTimesOptions: undefined,
    badgeIdsOptions: undefined,
    permittedTimesOptions: undefined,
    forbiddenTimesOptions: undefined,
  };
}

export const UserApprovedOutgoingTransferCombination = {
  encode(message: UserApprovedOutgoingTransferCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timelineTimesOptions !== undefined) {
      ValueOptions.encode(message.timelineTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.toMappingOptions !== undefined) {
      ValueOptions.encode(message.toMappingOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.initiatedByMappingOptions !== undefined) {
      ValueOptions.encode(message.initiatedByMappingOptions, writer.uint32(26).fork()).ldelim();
    }
    if (message.transferTimesOptions !== undefined) {
      ValueOptions.encode(message.transferTimesOptions, writer.uint32(34).fork()).ldelim();
    }
    if (message.badgeIdsOptions !== undefined) {
      ValueOptions.encode(message.badgeIdsOptions, writer.uint32(42).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(50).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedOutgoingTransferCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedOutgoingTransferCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.toMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.initiatedByMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 4:
          message.transferTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 5:
          message.badgeIdsOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 6:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 7:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedOutgoingTransferCombination {
    return {
      timelineTimesOptions: isSet(object.timelineTimesOptions)
        ? ValueOptions.fromJSON(object.timelineTimesOptions)
        : undefined,
      toMappingOptions: isSet(object.toMappingOptions) ? ValueOptions.fromJSON(object.toMappingOptions) : undefined,
      initiatedByMappingOptions: isSet(object.initiatedByMappingOptions)
        ? ValueOptions.fromJSON(object.initiatedByMappingOptions)
        : undefined,
      transferTimesOptions: isSet(object.transferTimesOptions)
        ? ValueOptions.fromJSON(object.transferTimesOptions)
        : undefined,
      badgeIdsOptions: isSet(object.badgeIdsOptions) ? ValueOptions.fromJSON(object.badgeIdsOptions) : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: UserApprovedOutgoingTransferCombination): unknown {
    const obj: any = {};
    message.timelineTimesOptions !== undefined && (obj.timelineTimesOptions = message.timelineTimesOptions
      ? ValueOptions.toJSON(message.timelineTimesOptions)
      : undefined);
    message.toMappingOptions !== undefined
      && (obj.toMappingOptions = message.toMappingOptions ? ValueOptions.toJSON(message.toMappingOptions) : undefined);
    message.initiatedByMappingOptions !== undefined
      && (obj.initiatedByMappingOptions = message.initiatedByMappingOptions
        ? ValueOptions.toJSON(message.initiatedByMappingOptions)
        : undefined);
    message.transferTimesOptions !== undefined && (obj.transferTimesOptions = message.transferTimesOptions
      ? ValueOptions.toJSON(message.transferTimesOptions)
      : undefined);
    message.badgeIdsOptions !== undefined
      && (obj.badgeIdsOptions = message.badgeIdsOptions ? ValueOptions.toJSON(message.badgeIdsOptions) : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransferCombination>, I>>(
    object: I,
  ): UserApprovedOutgoingTransferCombination {
    const message = createBaseUserApprovedOutgoingTransferCombination();
    message.timelineTimesOptions = (object.timelineTimesOptions !== undefined && object.timelineTimesOptions !== null)
      ? ValueOptions.fromPartial(object.timelineTimesOptions)
      : undefined;
    message.toMappingOptions = (object.toMappingOptions !== undefined && object.toMappingOptions !== null)
      ? ValueOptions.fromPartial(object.toMappingOptions)
      : undefined;
    message.initiatedByMappingOptions =
      (object.initiatedByMappingOptions !== undefined && object.initiatedByMappingOptions !== null)
        ? ValueOptions.fromPartial(object.initiatedByMappingOptions)
        : undefined;
    message.transferTimesOptions = (object.transferTimesOptions !== undefined && object.transferTimesOptions !== null)
      ? ValueOptions.fromPartial(object.transferTimesOptions)
      : undefined;
    message.badgeIdsOptions = (object.badgeIdsOptions !== undefined && object.badgeIdsOptions !== null)
      ? ValueOptions.fromPartial(object.badgeIdsOptions)
      : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseUserApprovedOutgoingTransferDefaultValues(): UserApprovedOutgoingTransferDefaultValues {
  return {
    timelineTimes: [],
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const UserApprovedOutgoingTransferDefaultValues = {
  encode(message: UserApprovedOutgoingTransferDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.toMappingId !== "") {
      writer.uint32(18).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(26).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedOutgoingTransferDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedOutgoingTransferDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.toMappingId = reader.string();
          break;
        case 3:
          message.initiatedByMappingId = reader.string();
          break;
        case 4:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 8:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedOutgoingTransferDefaultValues {
    return {
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedOutgoingTransferDefaultValues): unknown {
    const obj: any = {};
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    message.toMappingId !== undefined && (obj.toMappingId = message.toMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransferDefaultValues>, I>>(
    object: I,
  ): UserApprovedOutgoingTransferDefaultValues {
    const message = createBaseUserApprovedOutgoingTransferDefaultValues();
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserApprovedOutgoingTransferPermission(): UserApprovedOutgoingTransferPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const UserApprovedOutgoingTransferPermission = {
  encode(message: UserApprovedOutgoingTransferPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      UserApprovedOutgoingTransferDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      UserApprovedOutgoingTransferCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedOutgoingTransferPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedOutgoingTransferPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = UserApprovedOutgoingTransferDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(UserApprovedOutgoingTransferCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedOutgoingTransferPermission {
    return {
      defaultValues: isSet(object.defaultValues)
        ? UserApprovedOutgoingTransferDefaultValues.fromJSON(object.defaultValues)
        : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => UserApprovedOutgoingTransferCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedOutgoingTransferPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined && (obj.defaultValues = message.defaultValues
      ? UserApprovedOutgoingTransferDefaultValues.toJSON(message.defaultValues)
      : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) =>
        e ? UserApprovedOutgoingTransferCombination.toJSON(e) : undefined
      );
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransferPermission>, I>>(
    object: I,
  ): UserApprovedOutgoingTransferPermission {
    const message = createBaseUserApprovedOutgoingTransferPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? UserApprovedOutgoingTransferDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => UserApprovedOutgoingTransferCombination.fromPartial(e))
      || [];
    return message;
  },
};

function createBaseUserApprovedIncomingTransferCombination(): UserApprovedIncomingTransferCombination {
  return {
    timelineTimesOptions: undefined,
    fromMappingOptions: undefined,
    initiatedByMappingOptions: undefined,
    transferTimesOptions: undefined,
    badgeIdsOptions: undefined,
    permittedTimesOptions: undefined,
    forbiddenTimesOptions: undefined,
  };
}

export const UserApprovedIncomingTransferCombination = {
  encode(message: UserApprovedIncomingTransferCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timelineTimesOptions !== undefined) {
      ValueOptions.encode(message.timelineTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.fromMappingOptions !== undefined) {
      ValueOptions.encode(message.fromMappingOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.initiatedByMappingOptions !== undefined) {
      ValueOptions.encode(message.initiatedByMappingOptions, writer.uint32(26).fork()).ldelim();
    }
    if (message.transferTimesOptions !== undefined) {
      ValueOptions.encode(message.transferTimesOptions, writer.uint32(34).fork()).ldelim();
    }
    if (message.badgeIdsOptions !== undefined) {
      ValueOptions.encode(message.badgeIdsOptions, writer.uint32(42).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(50).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedIncomingTransferCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedIncomingTransferCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.fromMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.initiatedByMappingOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 4:
          message.transferTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 5:
          message.badgeIdsOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 6:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 7:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedIncomingTransferCombination {
    return {
      timelineTimesOptions: isSet(object.timelineTimesOptions)
        ? ValueOptions.fromJSON(object.timelineTimesOptions)
        : undefined,
      fromMappingOptions: isSet(object.fromMappingOptions)
        ? ValueOptions.fromJSON(object.fromMappingOptions)
        : undefined,
      initiatedByMappingOptions: isSet(object.initiatedByMappingOptions)
        ? ValueOptions.fromJSON(object.initiatedByMappingOptions)
        : undefined,
      transferTimesOptions: isSet(object.transferTimesOptions)
        ? ValueOptions.fromJSON(object.transferTimesOptions)
        : undefined,
      badgeIdsOptions: isSet(object.badgeIdsOptions) ? ValueOptions.fromJSON(object.badgeIdsOptions) : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: UserApprovedIncomingTransferCombination): unknown {
    const obj: any = {};
    message.timelineTimesOptions !== undefined && (obj.timelineTimesOptions = message.timelineTimesOptions
      ? ValueOptions.toJSON(message.timelineTimesOptions)
      : undefined);
    message.fromMappingOptions !== undefined && (obj.fromMappingOptions = message.fromMappingOptions
      ? ValueOptions.toJSON(message.fromMappingOptions)
      : undefined);
    message.initiatedByMappingOptions !== undefined
      && (obj.initiatedByMappingOptions = message.initiatedByMappingOptions
        ? ValueOptions.toJSON(message.initiatedByMappingOptions)
        : undefined);
    message.transferTimesOptions !== undefined && (obj.transferTimesOptions = message.transferTimesOptions
      ? ValueOptions.toJSON(message.transferTimesOptions)
      : undefined);
    message.badgeIdsOptions !== undefined
      && (obj.badgeIdsOptions = message.badgeIdsOptions ? ValueOptions.toJSON(message.badgeIdsOptions) : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransferCombination>, I>>(
    object: I,
  ): UserApprovedIncomingTransferCombination {
    const message = createBaseUserApprovedIncomingTransferCombination();
    message.timelineTimesOptions = (object.timelineTimesOptions !== undefined && object.timelineTimesOptions !== null)
      ? ValueOptions.fromPartial(object.timelineTimesOptions)
      : undefined;
    message.fromMappingOptions = (object.fromMappingOptions !== undefined && object.fromMappingOptions !== null)
      ? ValueOptions.fromPartial(object.fromMappingOptions)
      : undefined;
    message.initiatedByMappingOptions =
      (object.initiatedByMappingOptions !== undefined && object.initiatedByMappingOptions !== null)
        ? ValueOptions.fromPartial(object.initiatedByMappingOptions)
        : undefined;
    message.transferTimesOptions = (object.transferTimesOptions !== undefined && object.transferTimesOptions !== null)
      ? ValueOptions.fromPartial(object.transferTimesOptions)
      : undefined;
    message.badgeIdsOptions = (object.badgeIdsOptions !== undefined && object.badgeIdsOptions !== null)
      ? ValueOptions.fromPartial(object.badgeIdsOptions)
      : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseUserApprovedIncomingTransferDefaultValues(): UserApprovedIncomingTransferDefaultValues {
  return {
    timelineTimes: [],
    fromMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const UserApprovedIncomingTransferDefaultValues = {
  encode(message: UserApprovedIncomingTransferDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.fromMappingId !== "") {
      writer.uint32(18).string(message.fromMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(26).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedIncomingTransferDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedIncomingTransferDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.fromMappingId = reader.string();
          break;
        case 3:
          message.initiatedByMappingId = reader.string();
          break;
        case 4:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 8:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedIncomingTransferDefaultValues {
    return {
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedIncomingTransferDefaultValues): unknown {
    const obj: any = {};
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    message.fromMappingId !== undefined && (obj.fromMappingId = message.fromMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransferDefaultValues>, I>>(
    object: I,
  ): UserApprovedIncomingTransferDefaultValues {
    const message = createBaseUserApprovedIncomingTransferDefaultValues();
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.fromMappingId = object.fromMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserApprovedIncomingTransferPermission(): UserApprovedIncomingTransferPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const UserApprovedIncomingTransferPermission = {
  encode(message: UserApprovedIncomingTransferPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      UserApprovedIncomingTransferDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      UserApprovedIncomingTransferCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedIncomingTransferPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedIncomingTransferPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = UserApprovedIncomingTransferDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(UserApprovedIncomingTransferCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedIncomingTransferPermission {
    return {
      defaultValues: isSet(object.defaultValues)
        ? UserApprovedIncomingTransferDefaultValues.fromJSON(object.defaultValues)
        : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => UserApprovedIncomingTransferCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedIncomingTransferPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined && (obj.defaultValues = message.defaultValues
      ? UserApprovedIncomingTransferDefaultValues.toJSON(message.defaultValues)
      : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) =>
        e ? UserApprovedIncomingTransferCombination.toJSON(e) : undefined
      );
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransferPermission>, I>>(
    object: I,
  ): UserApprovedIncomingTransferPermission {
    const message = createBaseUserApprovedIncomingTransferPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? UserApprovedIncomingTransferDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => UserApprovedIncomingTransferCombination.fromPartial(e))
      || [];
    return message;
  },
};

function createBaseBalancesActionCombination(): BalancesActionCombination {
  return {
    badgeIdsOptions: undefined,
    ownedTimesOptions: undefined,
    permittedTimesOptions: undefined,
    forbiddenTimesOptions: undefined,
  };
}

export const BalancesActionCombination = {
  encode(message: BalancesActionCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.badgeIdsOptions !== undefined) {
      ValueOptions.encode(message.badgeIdsOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.ownedTimesOptions !== undefined) {
      ValueOptions.encode(message.ownedTimesOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(26).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BalancesActionCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalancesActionCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeIdsOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.ownedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 4:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BalancesActionCombination {
    return {
      badgeIdsOptions: isSet(object.badgeIdsOptions) ? ValueOptions.fromJSON(object.badgeIdsOptions) : undefined,
      ownedTimesOptions: isSet(object.ownedTimesOptions)
        ? ValueOptions.fromJSON(object.ownedTimesOptions)
        : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: BalancesActionCombination): unknown {
    const obj: any = {};
    message.badgeIdsOptions !== undefined
      && (obj.badgeIdsOptions = message.badgeIdsOptions ? ValueOptions.toJSON(message.badgeIdsOptions) : undefined);
    message.ownedTimesOptions !== undefined && (obj.ownedTimesOptions = message.ownedTimesOptions
      ? ValueOptions.toJSON(message.ownedTimesOptions)
      : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BalancesActionCombination>, I>>(object: I): BalancesActionCombination {
    const message = createBaseBalancesActionCombination();
    message.badgeIdsOptions = (object.badgeIdsOptions !== undefined && object.badgeIdsOptions !== null)
      ? ValueOptions.fromPartial(object.badgeIdsOptions)
      : undefined;
    message.ownedTimesOptions =
      (object.ownedTimesOptions !== undefined && object.ownedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.ownedTimesOptions)
        : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseBalancesActionDefaultValues(): BalancesActionDefaultValues {
  return { badgeIds: [], ownedTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const BalancesActionDefaultValues = {
  encode(message: BalancesActionDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BalancesActionDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalancesActionDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 3:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BalancesActionDefaultValues {
    return {
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownedTimes: Array.isArray(object?.ownedTimes)
        ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BalancesActionDefaultValues): unknown {
    const obj: any = {};
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BalancesActionDefaultValues>, I>>(object: I): BalancesActionDefaultValues {
    const message = createBaseBalancesActionDefaultValues();
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseBalancesActionPermission(): BalancesActionPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const BalancesActionPermission = {
  encode(message: BalancesActionPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      BalancesActionDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      BalancesActionCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BalancesActionPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalancesActionPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = BalancesActionDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(BalancesActionCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BalancesActionPermission {
    return {
      defaultValues: isSet(object.defaultValues)
        ? BalancesActionDefaultValues.fromJSON(object.defaultValues)
        : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => BalancesActionCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BalancesActionPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined && (obj.defaultValues = message.defaultValues
      ? BalancesActionDefaultValues.toJSON(message.defaultValues)
      : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) => e ? BalancesActionCombination.toJSON(e) : undefined);
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BalancesActionPermission>, I>>(object: I): BalancesActionPermission {
    const message = createBaseBalancesActionPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? BalancesActionDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => BalancesActionCombination.fromPartial(e)) || [];
    return message;
  },
};

function createBaseActionDefaultValues(): ActionDefaultValues {
  return { permittedTimes: [], forbiddenTimes: [] };
}

export const ActionDefaultValues = {
  encode(message: ActionDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActionDefaultValues {
    return {
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ActionDefaultValues): unknown {
    const obj: any = {};
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ActionDefaultValues>, I>>(object: I): ActionDefaultValues {
    const message = createBaseActionDefaultValues();
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseActionCombination(): ActionCombination {
  return { permittedTimesOptions: undefined, forbiddenTimesOptions: undefined };
}

export const ActionCombination = {
  encode(message: ActionCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActionCombination {
    return {
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: ActionCombination): unknown {
    const obj: any = {};
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ActionCombination>, I>>(object: I): ActionCombination {
    const message = createBaseActionCombination();
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseActionPermission(): ActionPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const ActionPermission = {
  encode(message: ActionPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      ActionDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      ActionCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = ActionDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(ActionCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActionPermission {
    return {
      defaultValues: isSet(object.defaultValues) ? ActionDefaultValues.fromJSON(object.defaultValues) : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => ActionCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ActionPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined
      && (obj.defaultValues = message.defaultValues ? ActionDefaultValues.toJSON(message.defaultValues) : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) => e ? ActionCombination.toJSON(e) : undefined);
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ActionPermission>, I>>(object: I): ActionPermission {
    const message = createBaseActionPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? ActionDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => ActionCombination.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdateCombination(): TimedUpdateCombination {
  return { timelineTimesOptions: undefined, permittedTimesOptions: undefined, forbiddenTimesOptions: undefined };
}

export const TimedUpdateCombination = {
  encode(message: TimedUpdateCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timelineTimesOptions !== undefined) {
      ValueOptions.encode(message.timelineTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdateCombination {
    return {
      timelineTimesOptions: isSet(object.timelineTimesOptions)
        ? ValueOptions.fromJSON(object.timelineTimesOptions)
        : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: TimedUpdateCombination): unknown {
    const obj: any = {};
    message.timelineTimesOptions !== undefined && (obj.timelineTimesOptions = message.timelineTimesOptions
      ? ValueOptions.toJSON(message.timelineTimesOptions)
      : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdateCombination>, I>>(object: I): TimedUpdateCombination {
    const message = createBaseTimedUpdateCombination();
    message.timelineTimesOptions = (object.timelineTimesOptions !== undefined && object.timelineTimesOptions !== null)
      ? ValueOptions.fromPartial(object.timelineTimesOptions)
      : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseTimedUpdateDefaultValues(): TimedUpdateDefaultValues {
  return { timelineTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const TimedUpdateDefaultValues = {
  encode(message: TimedUpdateDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 3:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdateDefaultValues {
    return {
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TimedUpdateDefaultValues): unknown {
    const obj: any = {};
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdateDefaultValues>, I>>(object: I): TimedUpdateDefaultValues {
    const message = createBaseTimedUpdateDefaultValues();
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdatePermission(): TimedUpdatePermission {
  return { defaultValues: undefined, combinations: [] };
}

export const TimedUpdatePermission = {
  encode(message: TimedUpdatePermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      TimedUpdateDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      TimedUpdateCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdatePermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdatePermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = TimedUpdateDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(TimedUpdateCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdatePermission {
    return {
      defaultValues: isSet(object.defaultValues) ? TimedUpdateDefaultValues.fromJSON(object.defaultValues) : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => TimedUpdateCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TimedUpdatePermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined
      && (obj.defaultValues = message.defaultValues
        ? TimedUpdateDefaultValues.toJSON(message.defaultValues)
        : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) => e ? TimedUpdateCombination.toJSON(e) : undefined);
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdatePermission>, I>>(object: I): TimedUpdatePermission {
    const message = createBaseTimedUpdatePermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? TimedUpdateDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => TimedUpdateCombination.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdateWithBadgeIdsCombination(): TimedUpdateWithBadgeIdsCombination {
  return {
    timelineTimesOptions: undefined,
    badgeIdsOptions: undefined,
    permittedTimesOptions: undefined,
    forbiddenTimesOptions: undefined,
  };
}

export const TimedUpdateWithBadgeIdsCombination = {
  encode(message: TimedUpdateWithBadgeIdsCombination, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timelineTimesOptions !== undefined) {
      ValueOptions.encode(message.timelineTimesOptions, writer.uint32(10).fork()).ldelim();
    }
    if (message.badgeIdsOptions !== undefined) {
      ValueOptions.encode(message.badgeIdsOptions, writer.uint32(18).fork()).ldelim();
    }
    if (message.permittedTimesOptions !== undefined) {
      ValueOptions.encode(message.permittedTimesOptions, writer.uint32(26).fork()).ldelim();
    }
    if (message.forbiddenTimesOptions !== undefined) {
      ValueOptions.encode(message.forbiddenTimesOptions, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateWithBadgeIdsCombination {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateWithBadgeIdsCombination();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.timelineTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 2:
          message.badgeIdsOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 3:
          message.permittedTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        case 4:
          message.forbiddenTimesOptions = ValueOptions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdateWithBadgeIdsCombination {
    return {
      timelineTimesOptions: isSet(object.timelineTimesOptions)
        ? ValueOptions.fromJSON(object.timelineTimesOptions)
        : undefined,
      badgeIdsOptions: isSet(object.badgeIdsOptions) ? ValueOptions.fromJSON(object.badgeIdsOptions) : undefined,
      permittedTimesOptions: isSet(object.permittedTimesOptions)
        ? ValueOptions.fromJSON(object.permittedTimesOptions)
        : undefined,
      forbiddenTimesOptions: isSet(object.forbiddenTimesOptions)
        ? ValueOptions.fromJSON(object.forbiddenTimesOptions)
        : undefined,
    };
  },

  toJSON(message: TimedUpdateWithBadgeIdsCombination): unknown {
    const obj: any = {};
    message.timelineTimesOptions !== undefined && (obj.timelineTimesOptions = message.timelineTimesOptions
      ? ValueOptions.toJSON(message.timelineTimesOptions)
      : undefined);
    message.badgeIdsOptions !== undefined
      && (obj.badgeIdsOptions = message.badgeIdsOptions ? ValueOptions.toJSON(message.badgeIdsOptions) : undefined);
    message.permittedTimesOptions !== undefined && (obj.permittedTimesOptions = message.permittedTimesOptions
      ? ValueOptions.toJSON(message.permittedTimesOptions)
      : undefined);
    message.forbiddenTimesOptions !== undefined && (obj.forbiddenTimesOptions = message.forbiddenTimesOptions
      ? ValueOptions.toJSON(message.forbiddenTimesOptions)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdateWithBadgeIdsCombination>, I>>(
    object: I,
  ): TimedUpdateWithBadgeIdsCombination {
    const message = createBaseTimedUpdateWithBadgeIdsCombination();
    message.timelineTimesOptions = (object.timelineTimesOptions !== undefined && object.timelineTimesOptions !== null)
      ? ValueOptions.fromPartial(object.timelineTimesOptions)
      : undefined;
    message.badgeIdsOptions = (object.badgeIdsOptions !== undefined && object.badgeIdsOptions !== null)
      ? ValueOptions.fromPartial(object.badgeIdsOptions)
      : undefined;
    message.permittedTimesOptions =
      (object.permittedTimesOptions !== undefined && object.permittedTimesOptions !== null)
        ? ValueOptions.fromPartial(object.permittedTimesOptions)
        : undefined;
    message.forbiddenTimesOptions =
      (object.forbiddenTimesOptions !== undefined && object.forbiddenTimesOptions !== null)
        ? ValueOptions.fromPartial(object.forbiddenTimesOptions)
        : undefined;
    return message;
  },
};

function createBaseTimedUpdateWithBadgeIdsDefaultValues(): TimedUpdateWithBadgeIdsDefaultValues {
  return { badgeIds: [], timelineTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const TimedUpdateWithBadgeIdsDefaultValues = {
  encode(message: TimedUpdateWithBadgeIdsDefaultValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateWithBadgeIdsDefaultValues {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateWithBadgeIdsDefaultValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 3:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdateWithBadgeIdsDefaultValues {
    return {
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TimedUpdateWithBadgeIdsDefaultValues): unknown {
    const obj: any = {};
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    if (message.permittedTimes) {
      obj.permittedTimes = message.permittedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.permittedTimes = [];
    }
    if (message.forbiddenTimes) {
      obj.forbiddenTimes = message.forbiddenTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.forbiddenTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdateWithBadgeIdsDefaultValues>, I>>(
    object: I,
  ): TimedUpdateWithBadgeIdsDefaultValues {
    const message = createBaseTimedUpdateWithBadgeIdsDefaultValues();
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdateWithBadgeIdsPermission(): TimedUpdateWithBadgeIdsPermission {
  return { defaultValues: undefined, combinations: [] };
}

export const TimedUpdateWithBadgeIdsPermission = {
  encode(message: TimedUpdateWithBadgeIdsPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.defaultValues !== undefined) {
      TimedUpdateWithBadgeIdsDefaultValues.encode(message.defaultValues, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.combinations) {
      TimedUpdateWithBadgeIdsCombination.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateWithBadgeIdsPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateWithBadgeIdsPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.defaultValues = TimedUpdateWithBadgeIdsDefaultValues.decode(reader, reader.uint32());
          break;
        case 2:
          message.combinations.push(TimedUpdateWithBadgeIdsCombination.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TimedUpdateWithBadgeIdsPermission {
    return {
      defaultValues: isSet(object.defaultValues)
        ? TimedUpdateWithBadgeIdsDefaultValues.fromJSON(object.defaultValues)
        : undefined,
      combinations: Array.isArray(object?.combinations)
        ? object.combinations.map((e: any) => TimedUpdateWithBadgeIdsCombination.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TimedUpdateWithBadgeIdsPermission): unknown {
    const obj: any = {};
    message.defaultValues !== undefined && (obj.defaultValues = message.defaultValues
      ? TimedUpdateWithBadgeIdsDefaultValues.toJSON(message.defaultValues)
      : undefined);
    if (message.combinations) {
      obj.combinations = message.combinations.map((e) => e ? TimedUpdateWithBadgeIdsCombination.toJSON(e) : undefined);
    } else {
      obj.combinations = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TimedUpdateWithBadgeIdsPermission>, I>>(
    object: I,
  ): TimedUpdateWithBadgeIdsPermission {
    const message = createBaseTimedUpdateWithBadgeIdsPermission();
    message.defaultValues = (object.defaultValues !== undefined && object.defaultValues !== null)
      ? TimedUpdateWithBadgeIdsDefaultValues.fromPartial(object.defaultValues)
      : undefined;
    message.combinations = object.combinations?.map((e) => TimedUpdateWithBadgeIdsCombination.fromPartial(e)) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
