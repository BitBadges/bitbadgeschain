/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { UintRange } from "./balances";

export const protobufPackage = "badges";

/**
 * CollectionPermissions defines the permissions for the collection (i.e. what the manager can and cannot do).
 *
 * There are five types of permissions for a collection: ActionPermission, TimedUpdatePermission, TimedUpdateWithBadgeIdsPermission, BalancesActionPermission, and CollectionApprovalPermission.
 *
 * The permission type allows fine-grained access control for each action.
 * ActionPermission: defines when the manager can perform an action.
 * TimedUpdatePermission: defines when the manager can update a timeline-based field and what times of the timeline can be updated.
 * TimedUpdateWithBadgeIdsPermission: defines when the manager can update a timeline-based field for specific badges and what times of the timeline can be updated.
 * BalancesActionPermission: defines when the manager can perform an action for specific badges and specific badge ownership times.
 * CollectionApprovalPermission: defines when the manager can update the transferability of the collection and what transfers can be updated vs locked
 *
 * Note there are a few different times here which could get confusing:
 * - timelineTimes: the times when a timeline-based field is a specific value
 * - permitted/forbiddenTimes - the times that a permission can be performed
 * - transferTimes - the times that a transfer occurs
 * - ownershipTimes - the times when a badge is owned by a user
 *
 * The permitted/forbiddenTimes are used to determine when a permission can be executed.
 * Once a time is set to be permitted or forbidden, it is PERMANENT and cannot be changed.
 * If a time is not set to be permitted or forbidden, it is considered NEUTRAL and can be updated but is ALLOWED by default.
 *
 * Each permission type has a defaultValues field and a combinations field.
 * The defaultValues field defines the default values for the permission which can be manipulated by the combinations field (to avoid unnecessary repetition).
 * Ex: We can have default value badgeIds = [1,2] and combinations = [{invertDefault: true, isApproved: false}, {isApproved: true}].
 * This would mean that badgeIds [1,2] are allowed but everything else is not allowed.
 *
 * IMPORTANT: For all permissions, we ONLY take the first combination that matches. Any subsequent combinations are ignored.
 * Ex: If we have defaultValues = {badgeIds: [1,2]} and combinations = [{isApproved: true}, {isApproved: false}].
 * This would mean that badgeIds [1,2] are allowed and the second combination is ignored.
 */
export interface CollectionPermissions {
  canDeleteCollection: ActionPermission[];
  canArchiveCollection: TimedUpdatePermission[];
  canUpdateOffChainBalancesMetadata: TimedUpdatePermission[];
  canUpdateStandards: TimedUpdatePermission[];
  canUpdateCustomData: TimedUpdatePermission[];
  canUpdateManager: TimedUpdatePermission[];
  canUpdateCollectionMetadata: TimedUpdatePermission[];
  canCreateMoreBadges: BalancesActionPermission[];
  canUpdateBadgeMetadata: TimedUpdateWithBadgeIdsPermission[];
  canUpdateCollectionApprovals: CollectionApprovalPermission[];
}

/**
 * UserPermissions defines the permissions for the user (i.e. what the user can and cannot do).
 *
 * See CollectionPermissions for more details on the different types of permissions.
 * The UserApprovedOutgoing and UserApprovedIncoming permissions are the same as the CollectionApprovalPermission,
 * but certain fields are removed because they are not relevant to the user.
 */
export interface UserPermissions {
  canUpdateOutgoingApprovals: UserOutgoingApprovalPermission[];
  canUpdateIncomingApprovals: UserIncomingApprovalPermission[];
  canUpdateAutoApproveSelfInitiatedOutgoingTransfers: ActionPermission[];
  canUpdateAutoApproveSelfInitiatedIncomingTransfers: ActionPermission[];
}

/**
 * CollectionApprovalPermission defines what collection approved transfers can be updated vs are locked.
 *
 * Each transfer is broken down to a (from, to, initiatedBy, transferTime, badgeId) tuple.
 * For a transfer to match, we need to match ALL of the fields in the combination.
 * These are detemined by the fromMappingId, toMappingId, initiatedByMappingId, transferTimes, badgeIds fields.
 * AddressMappings are used for (from, to, initiatedBy) which are a permanent list of addresses identified by an ID (see AddressMappings).
 *
 * TimelineTimes: which timeline times of the collection's approvalsTimeline field can be updated or not?
 * permitted/forbidden TimelineTimes: when can the manager execute this permission?
 *
 * Ex: Let's say we are updating the transferability for timelineTime 1 and the transfer tuple ("AllWithoutMint", "AllWithoutMint", "AllWithoutMint", 10, 1000).
 * We would check to find the FIRST CollectionApprovalPermission that matches this combination.
 * If we find a match, we would check the permitted/forbidden times to see if we can execute this permission (default is ALLOWED).
 *
 * Ex: So if you wanted to freeze the transferability to enforce that badge ID 1 will always be transferable, you could set
 * the combination ("AllWithoutMint", "AllWithoutMint", "AllWithoutMint", "All Transfer Times", 1) to always be forbidden at all timelineTimes.
 */
export interface CollectionApprovalPermission {
  fromMappingId: string;
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  amountTrackerId: string;
  challengeTrackerId: string;
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * UserOutgoingApprovalPermission defines the permissions for updating the user's approved outgoing transfers.
 * See CollectionApprovalPermission for more details. This is equivalent without the fromMappingId field because that is always the user.
 */
export interface UserOutgoingApprovalPermission {
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  amountTrackerId: string;
  challengeTrackerId: string;
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * UserIncomingApprovalPermission defines the permissions for updating the user's approved incoming transfers.
 * See CollectionApprovalPermission for more details. This is equivalent without the toMappingId field because that is always the user.
 */
export interface UserIncomingApprovalPermission {
  fromMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  amountTrackerId: string;
  challengeTrackerId: string;
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * BalancesActionPermission defines the permissions for updating a timeline-based field for specific badges and specific badge ownership times.
 * Currently, this is only used for creating new badges.
 *
 * Ex: If you want to lock the ability to create new badges for badgeIds [1,2] at ownershipTimes 1/1/2020 - 1/1/2021,
 * you could set the combination (badgeIds: [1,2], ownershipTimelineTimes: [1/1/2020 - 1/1/2021]) to always be forbidden.
 */
export interface BalancesActionPermission {
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

/**
 * ActionPermission defines the permissions for performing an action.
 *
 * This is simple and straightforward as the only thing we need to check is the permitted/forbidden times.
 */
export interface ActionPermission {
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
  badgeIds: UintRange[];
  timelineTimes: UintRange[];
  permittedTimes: UintRange[];
  forbiddenTimes: UintRange[];
}

function createBaseCollectionPermissions(): CollectionPermissions {
  return {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canCreateMoreBadges: [],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [],
  };
}

export const CollectionPermissions = {
  encode(message: CollectionPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.canDeleteCollection) {
      ActionPermission.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.canArchiveCollection) {
      TimedUpdatePermission.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.canUpdateOffChainBalancesMetadata) {
      TimedUpdatePermission.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.canUpdateStandards) {
      TimedUpdatePermission.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.canUpdateCustomData) {
      TimedUpdatePermission.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.canUpdateManager) {
      TimedUpdatePermission.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.canUpdateCollectionMetadata) {
      TimedUpdatePermission.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.canCreateMoreBadges) {
      BalancesActionPermission.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.canUpdateBadgeMetadata) {
      TimedUpdateWithBadgeIdsPermission.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.canUpdateCollectionApprovals) {
      CollectionApprovalPermission.encode(v!, writer.uint32(82).fork()).ldelim();
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
          message.canArchiveCollection.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 3:
          message.canUpdateOffChainBalancesMetadata.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 4:
          message.canUpdateStandards.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 5:
          message.canUpdateCustomData.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 6:
          message.canUpdateManager.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 7:
          message.canUpdateCollectionMetadata.push(TimedUpdatePermission.decode(reader, reader.uint32()));
          break;
        case 8:
          message.canCreateMoreBadges.push(BalancesActionPermission.decode(reader, reader.uint32()));
          break;
        case 9:
          message.canUpdateBadgeMetadata.push(TimedUpdateWithBadgeIdsPermission.decode(reader, reader.uint32()));
          break;
        case 10:
          message.canUpdateCollectionApprovals.push(CollectionApprovalPermission.decode(reader, reader.uint32()));
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
      canArchiveCollection: Array.isArray(object?.canArchiveCollection)
        ? object.canArchiveCollection.map((e: any) => TimedUpdatePermission.fromJSON(e))
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
      canUpdateCollectionApprovals: Array.isArray(object?.canUpdateCollectionApprovals)
        ? object.canUpdateCollectionApprovals.map((e: any) => CollectionApprovalPermission.fromJSON(e))
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
    if (message.canArchiveCollection) {
      obj.canArchiveCollection = message.canArchiveCollection.map((e) =>
        e ? TimedUpdatePermission.toJSON(e) : undefined
      );
    } else {
      obj.canArchiveCollection = [];
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
    if (message.canUpdateCollectionApprovals) {
      obj.canUpdateCollectionApprovals = message.canUpdateCollectionApprovals.map((e) =>
        e ? CollectionApprovalPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateCollectionApprovals = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionPermissions>, I>>(object: I): CollectionPermissions {
    const message = createBaseCollectionPermissions();
    message.canDeleteCollection = object.canDeleteCollection?.map((e) => ActionPermission.fromPartial(e)) || [];
    message.canArchiveCollection = object.canArchiveCollection?.map((e) => TimedUpdatePermission.fromPartial(e)) || [];
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
    message.canUpdateCollectionApprovals =
      object.canUpdateCollectionApprovals?.map((e) => CollectionApprovalPermission.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserPermissions(): UserPermissions {
  return {
    canUpdateOutgoingApprovals: [],
    canUpdateIncomingApprovals: [],
    canUpdateAutoApproveSelfInitiatedOutgoingTransfers: [],
    canUpdateAutoApproveSelfInitiatedIncomingTransfers: [],
  };
}

export const UserPermissions = {
  encode(message: UserPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.canUpdateOutgoingApprovals) {
      UserOutgoingApprovalPermission.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.canUpdateIncomingApprovals) {
      UserIncomingApprovalPermission.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.canUpdateAutoApproveSelfInitiatedOutgoingTransfers) {
      ActionPermission.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.canUpdateAutoApproveSelfInitiatedIncomingTransfers) {
      ActionPermission.encode(v!, writer.uint32(34).fork()).ldelim();
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
          message.canUpdateOutgoingApprovals.push(UserOutgoingApprovalPermission.decode(reader, reader.uint32()));
          break;
        case 2:
          message.canUpdateIncomingApprovals.push(UserIncomingApprovalPermission.decode(reader, reader.uint32()));
          break;
        case 3:
          message.canUpdateAutoApproveSelfInitiatedOutgoingTransfers.push(
            ActionPermission.decode(reader, reader.uint32()),
          );
          break;
        case 4:
          message.canUpdateAutoApproveSelfInitiatedIncomingTransfers.push(
            ActionPermission.decode(reader, reader.uint32()),
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
      canUpdateOutgoingApprovals: Array.isArray(object?.canUpdateOutgoingApprovals)
        ? object.canUpdateOutgoingApprovals.map((e: any) => UserOutgoingApprovalPermission.fromJSON(e))
        : [],
      canUpdateIncomingApprovals: Array.isArray(object?.canUpdateIncomingApprovals)
        ? object.canUpdateIncomingApprovals.map((e: any) => UserIncomingApprovalPermission.fromJSON(e))
        : [],
      canUpdateAutoApproveSelfInitiatedOutgoingTransfers:
        Array.isArray(object?.canUpdateAutoApproveSelfInitiatedOutgoingTransfers)
          ? object.canUpdateAutoApproveSelfInitiatedOutgoingTransfers.map((e: any) => ActionPermission.fromJSON(e))
          : [],
      canUpdateAutoApproveSelfInitiatedIncomingTransfers:
        Array.isArray(object?.canUpdateAutoApproveSelfInitiatedIncomingTransfers)
          ? object.canUpdateAutoApproveSelfInitiatedIncomingTransfers.map((e: any) => ActionPermission.fromJSON(e))
          : [],
    };
  },

  toJSON(message: UserPermissions): unknown {
    const obj: any = {};
    if (message.canUpdateOutgoingApprovals) {
      obj.canUpdateOutgoingApprovals = message.canUpdateOutgoingApprovals.map((e) =>
        e ? UserOutgoingApprovalPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateOutgoingApprovals = [];
    }
    if (message.canUpdateIncomingApprovals) {
      obj.canUpdateIncomingApprovals = message.canUpdateIncomingApprovals.map((e) =>
        e ? UserIncomingApprovalPermission.toJSON(e) : undefined
      );
    } else {
      obj.canUpdateIncomingApprovals = [];
    }
    if (message.canUpdateAutoApproveSelfInitiatedOutgoingTransfers) {
      obj.canUpdateAutoApproveSelfInitiatedOutgoingTransfers = message
        .canUpdateAutoApproveSelfInitiatedOutgoingTransfers.map((e) => e ? ActionPermission.toJSON(e) : undefined);
    } else {
      obj.canUpdateAutoApproveSelfInitiatedOutgoingTransfers = [];
    }
    if (message.canUpdateAutoApproveSelfInitiatedIncomingTransfers) {
      obj.canUpdateAutoApproveSelfInitiatedIncomingTransfers = message
        .canUpdateAutoApproveSelfInitiatedIncomingTransfers.map((e) => e ? ActionPermission.toJSON(e) : undefined);
    } else {
      obj.canUpdateAutoApproveSelfInitiatedIncomingTransfers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserPermissions>, I>>(object: I): UserPermissions {
    const message = createBaseUserPermissions();
    message.canUpdateOutgoingApprovals =
      object.canUpdateOutgoingApprovals?.map((e) => UserOutgoingApprovalPermission.fromPartial(e)) || [];
    message.canUpdateIncomingApprovals =
      object.canUpdateIncomingApprovals?.map((e) => UserIncomingApprovalPermission.fromPartial(e)) || [];
    message.canUpdateAutoApproveSelfInitiatedOutgoingTransfers =
      object.canUpdateAutoApproveSelfInitiatedOutgoingTransfers?.map((e) => ActionPermission.fromPartial(e)) || [];
    message.canUpdateAutoApproveSelfInitiatedIncomingTransfers =
      object.canUpdateAutoApproveSelfInitiatedIncomingTransfers?.map((e) => ActionPermission.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCollectionApprovalPermission(): CollectionApprovalPermission {
  return {
    fromMappingId: "",
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const CollectionApprovalPermission = {
  encode(message: CollectionApprovalPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromMappingId !== "") {
      writer.uint32(10).string(message.fromMappingId);
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
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(58).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(66).string(message.challengeTrackerId);
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovalPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovalPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromMappingId = reader.string();
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
        case 6:
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.amountTrackerId = reader.string();
          break;
        case 8:
          message.challengeTrackerId = reader.string();
          break;
        case 9:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 10:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovalPermission {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionApprovalPermission): unknown {
    const obj: any = {};
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
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

  fromPartial<I extends Exact<DeepPartial<CollectionApprovalPermission>, I>>(object: I): CollectionApprovalPermission {
    const message = createBaseCollectionApprovalPermission();
    message.fromMappingId = object.fromMappingId ?? "";
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserOutgoingApprovalPermission(): UserOutgoingApprovalPermission {
  return {
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const UserOutgoingApprovalPermission = {
  encode(message: UserOutgoingApprovalPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.toMappingId !== "") {
      writer.uint32(10).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(18).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(50).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(58).string(message.challengeTrackerId);
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserOutgoingApprovalPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserOutgoingApprovalPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.toMappingId = reader.string();
          break;
        case 2:
          message.initiatedByMappingId = reader.string();
          break;
        case 3:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.amountTrackerId = reader.string();
          break;
        case 7:
          message.challengeTrackerId = reader.string();
          break;
        case 8:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 9:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserOutgoingApprovalPermission {
    return {
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserOutgoingApprovalPermission): unknown {
    const obj: any = {};
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
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

  fromPartial<I extends Exact<DeepPartial<UserOutgoingApprovalPermission>, I>>(
    object: I,
  ): UserOutgoingApprovalPermission {
    const message = createBaseUserOutgoingApprovalPermission();
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserIncomingApprovalPermission(): UserIncomingApprovalPermission {
  return {
    fromMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    permittedTimes: [],
    forbiddenTimes: [],
  };
}

export const UserIncomingApprovalPermission = {
  encode(message: UserIncomingApprovalPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromMappingId !== "") {
      writer.uint32(10).string(message.fromMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(18).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(50).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(58).string(message.challengeTrackerId);
    }
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserIncomingApprovalPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserIncomingApprovalPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromMappingId = reader.string();
          break;
        case 2:
          message.initiatedByMappingId = reader.string();
          break;
        case 3:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.amountTrackerId = reader.string();
          break;
        case 7:
          message.challengeTrackerId = reader.string();
          break;
        case 8:
          message.permittedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 9:
          message.forbiddenTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserIncomingApprovalPermission {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserIncomingApprovalPermission): unknown {
    const obj: any = {};
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
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

  fromPartial<I extends Exact<DeepPartial<UserIncomingApprovalPermission>, I>>(
    object: I,
  ): UserIncomingApprovalPermission {
    const message = createBaseUserIncomingApprovalPermission();
    message.fromMappingId = object.fromMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseBalancesActionPermission(): BalancesActionPermission {
  return { badgeIds: [], ownershipTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const BalancesActionPermission = {
  encode(message: BalancesActionPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.ownershipTimes) {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): BalancesActionPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalancesActionPermission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
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

  fromJSON(object: any): BalancesActionPermission {
    return {
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BalancesActionPermission): unknown {
    const obj: any = {};
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
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

  fromPartial<I extends Exact<DeepPartial<BalancesActionPermission>, I>>(object: I): BalancesActionPermission {
    const message = createBaseBalancesActionPermission();
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseActionPermission(): ActionPermission {
  return { permittedTimes: [], forbiddenTimes: [] };
}

export const ActionPermission = {
  encode(message: ActionPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.permittedTimes) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.forbiddenTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
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

  fromJSON(object: any): ActionPermission {
    return {
      permittedTimes: Array.isArray(object?.permittedTimes)
        ? object.permittedTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      forbiddenTimes: Array.isArray(object?.forbiddenTimes)
        ? object.forbiddenTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ActionPermission): unknown {
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

  fromPartial<I extends Exact<DeepPartial<ActionPermission>, I>>(object: I): ActionPermission {
    const message = createBaseActionPermission();
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdatePermission(): TimedUpdatePermission {
  return { timelineTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const TimedUpdatePermission = {
  encode(message: TimedUpdatePermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdatePermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdatePermission();
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

  fromJSON(object: any): TimedUpdatePermission {
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

  toJSON(message: TimedUpdatePermission): unknown {
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

  fromPartial<I extends Exact<DeepPartial<TimedUpdatePermission>, I>>(object: I): TimedUpdatePermission {
    const message = createBaseTimedUpdatePermission();
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTimedUpdateWithBadgeIdsPermission(): TimedUpdateWithBadgeIdsPermission {
  return { badgeIds: [], timelineTimes: [], permittedTimes: [], forbiddenTimes: [] };
}

export const TimedUpdateWithBadgeIdsPermission = {
  encode(message: TimedUpdateWithBadgeIdsPermission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): TimedUpdateWithBadgeIdsPermission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTimedUpdateWithBadgeIdsPermission();
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

  fromJSON(object: any): TimedUpdateWithBadgeIdsPermission {
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

  toJSON(message: TimedUpdateWithBadgeIdsPermission): unknown {
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

  fromPartial<I extends Exact<DeepPartial<TimedUpdateWithBadgeIdsPermission>, I>>(
    object: I,
  ): TimedUpdateWithBadgeIdsPermission {
    const message = createBaseTimedUpdateWithBadgeIdsPermission();
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.permittedTimes = object.permittedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.forbiddenTimes = object.forbiddenTimes?.map((e) => UintRange.fromPartial(e)) || [];
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
