/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { CollectionPermissions, UserPermissions } from "./permissions";
import {
  BadgeMetadataTimeline,
  CollectionMetadataTimeline,
  CustomDataTimeline,
  IsArchivedTimeline,
  ManagerTimeline,
  OffChainBalancesMetadataTimeline,
  StandardsTimeline,
} from "./timelines";
import { CollectionApproval, UserIncomingApproval, UserOutgoingApproval } from "./transfers";

export const protobufPackage = "badges";

/**
 * A BadgeCollection is the top level object for a collection of badges.
 * It defines everything about the collection, such as the manager, metadata, etc.
 *
 * All collections are identified by a collectionId assigned by the blockchain, which is a uint64 that increments (i.e. first collection has ID 1).
 *
 * All collections also have a manager who is responsible for managing the collection.
 * They can be granted certain permissions, such as the ability to mint new badges.
 *
 * Certain fields are timeline-based, which means they may have different values at different block heights.
 * We fetch the value according to the current time.
 * For example, we may set the manager to be Alice from Time1 to Time2, and then set the manager to be Bob from Time2 to Time3.
 *
 * Collections may have different balance types: standard vs off-chain vs inherited. See documentation for differences.
 */
export interface BadgeCollection {
  /** The collectionId is the unique identifier for this collection. */
  collectionId: string;
  /** The collection metadata is the metadata for the collection itself. */
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  /** The badge metadata is the metadata for each badge in the collection. */
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  /** The balancesType is the type of balances this collection uses (standard, off-chain, or inherited). */
  balancesType: string;
  /** The off-chain balances metadata defines where to fetch the balances for collections with off-chain balances. */
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  /** The custom data field is an arbitrary field that can be used to store any data. */
  customDataTimeline: CustomDataTimeline[];
  /** The manager is the address of the manager of this collection. */
  managerTimeline: ManagerTimeline[];
  /** The permissions define what the manager of the collection can do or not do. */
  collectionPermissions:
    | CollectionPermissions
    | undefined;
  /**
   * The approved transfers defines the transferability of the collection for collections with standard balances.
   * This defines it on a collection-level. All transfers must be explicitly allowed on the collection-level, or else, they will fail.
   *
   * Collection approved transfers can optionally specify to override the user approvals for a transfer (e.g. forcefully revoke a badge).
   * If user approvals are not overriden, then a transfer must also satisfy the From user's approved outgoing transfers and the To user's approved incoming transfers.
   */
  collectionApprovals: CollectionApproval[];
  /** Standards allow us to define a standard for the collection. This lets others know how to interpret the fields of the collection. */
  standardsTimeline: StandardsTimeline[];
  /**
   * The isArchivedTimeline defines whether the collection is archived or not.
   * When a collection is archived, it is read-only and no transactions can be processed.
   */
  isArchivedTimeline: IsArchivedTimeline[];
  /**
   * The defaultUserOutgoingApprovals defines the default user approved outgoing transfers for an uninitialized user balance.
   * The user can change this value at any time.
   */
  defaultUserOutgoingApprovals: UserOutgoingApproval[];
  /**
   * The defaultUserIncomingApprovals defines the default user approved incoming transfers for an uninitialized user balance.
   * The user can change this value at any time.
   *
   * Ex: Set this to disallow all incoming transfers by default, making the user have to opt-in to receiving the badge.
   */
  defaultUserIncomingApprovals: UserIncomingApproval[];
  defaultUserPermissions: UserPermissions | undefined;
  defaultAutoApproveSelfInitiatedOutgoingTransfers: boolean;
  defaultAutoApproveSelfInitiatedIncomingTransfers: boolean;
  createdBy: string;
}

function createBaseBadgeCollection(): BadgeCollection {
  return {
    collectionId: "",
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    balancesType: "",
    offChainBalancesMetadataTimeline: [],
    customDataTimeline: [],
    managerTimeline: [],
    collectionPermissions: undefined,
    collectionApprovals: [],
    standardsTimeline: [],
    isArchivedTimeline: [],
    defaultUserOutgoingApprovals: [],
    defaultUserIncomingApprovals: [],
    defaultUserPermissions: undefined,
    defaultAutoApproveSelfInitiatedOutgoingTransfers: false,
    defaultAutoApproveSelfInitiatedIncomingTransfers: false,
    createdBy: "",
  };
}

export const BadgeCollection = {
  encode(message: BadgeCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    for (const v of message.collectionMetadataTimeline) {
      CollectionMetadataTimeline.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.badgeMetadataTimeline) {
      BadgeMetadataTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.balancesType !== "") {
      writer.uint32(34).string(message.balancesType);
    }
    for (const v of message.offChainBalancesMetadataTimeline) {
      OffChainBalancesMetadataTimeline.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.customDataTimeline) {
      CustomDataTimeline.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.managerTimeline) {
      ManagerTimeline.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    if (message.collectionPermissions !== undefined) {
      CollectionPermissions.encode(message.collectionPermissions, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.defaultUserOutgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.defaultUserIncomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    if (message.defaultUserPermissions !== undefined) {
      UserPermissions.encode(message.defaultUserPermissions, writer.uint32(130).fork()).ldelim();
    }
    if (message.defaultAutoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(136).bool(message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.defaultAutoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(144).bool(message.defaultAutoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.createdBy !== "") {
      writer.uint32(154).string(message.createdBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BadgeCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBadgeCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionId = reader.string();
          break;
        case 2:
          message.collectionMetadataTimeline.push(CollectionMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 3:
          message.badgeMetadataTimeline.push(BadgeMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 4:
          message.balancesType = reader.string();
          break;
        case 5:
          message.offChainBalancesMetadataTimeline.push(
            OffChainBalancesMetadataTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 7:
          message.customDataTimeline.push(CustomDataTimeline.decode(reader, reader.uint32()));
          break;
        case 8:
          message.managerTimeline.push(ManagerTimeline.decode(reader, reader.uint32()));
          break;
        case 9:
          message.collectionPermissions = CollectionPermissions.decode(reader, reader.uint32());
          break;
        case 10:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 11:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
          break;
        case 12:
          message.isArchivedTimeline.push(IsArchivedTimeline.decode(reader, reader.uint32()));
          break;
        case 14:
          message.defaultUserOutgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 15:
          message.defaultUserIncomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 16:
          message.defaultUserPermissions = UserPermissions.decode(reader, reader.uint32());
          break;
        case 17:
          message.defaultAutoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 18:
          message.defaultAutoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 19:
          message.createdBy = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeCollection {
    return {
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      balancesType: isSet(object.balancesType) ? String(object.balancesType) : "",
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      customDataTimeline: Array.isArray(object?.customDataTimeline)
        ? object.customDataTimeline.map((e: any) => CustomDataTimeline.fromJSON(e))
        : [],
      managerTimeline: Array.isArray(object?.managerTimeline)
        ? object.managerTimeline.map((e: any) => ManagerTimeline.fromJSON(e))
        : [],
      collectionPermissions: isSet(object.collectionPermissions)
        ? CollectionPermissions.fromJSON(object.collectionPermissions)
        : undefined,
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
      defaultUserOutgoingApprovals: Array.isArray(object?.defaultUserOutgoingApprovals)
        ? object.defaultUserOutgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      defaultUserIncomingApprovals: Array.isArray(object?.defaultUserIncomingApprovals)
        ? object.defaultUserIncomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      defaultUserPermissions: isSet(object.defaultUserPermissions)
        ? UserPermissions.fromJSON(object.defaultUserPermissions)
        : undefined,
      defaultAutoApproveSelfInitiatedOutgoingTransfers: isSet(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        : false,
      defaultAutoApproveSelfInitiatedIncomingTransfers: isSet(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        : false,
      createdBy: isSet(object.createdBy) ? String(object.createdBy) : "",
    };
  },

  toJSON(message: BadgeCollection): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.collectionMetadataTimeline) {
      obj.collectionMetadataTimeline = message.collectionMetadataTimeline.map((e) =>
        e ? CollectionMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionMetadataTimeline = [];
    }
    if (message.badgeMetadataTimeline) {
      obj.badgeMetadataTimeline = message.badgeMetadataTimeline.map((e) =>
        e ? BadgeMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.badgeMetadataTimeline = [];
    }
    message.balancesType !== undefined && (obj.balancesType = message.balancesType);
    if (message.offChainBalancesMetadataTimeline) {
      obj.offChainBalancesMetadataTimeline = message.offChainBalancesMetadataTimeline.map((e) =>
        e ? OffChainBalancesMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.offChainBalancesMetadataTimeline = [];
    }
    if (message.customDataTimeline) {
      obj.customDataTimeline = message.customDataTimeline.map((e) => e ? CustomDataTimeline.toJSON(e) : undefined);
    } else {
      obj.customDataTimeline = [];
    }
    if (message.managerTimeline) {
      obj.managerTimeline = message.managerTimeline.map((e) => e ? ManagerTimeline.toJSON(e) : undefined);
    } else {
      obj.managerTimeline = [];
    }
    message.collectionPermissions !== undefined && (obj.collectionPermissions = message.collectionPermissions
      ? CollectionPermissions.toJSON(message.collectionPermissions)
      : undefined);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    if (message.isArchivedTimeline) {
      obj.isArchivedTimeline = message.isArchivedTimeline.map((e) => e ? IsArchivedTimeline.toJSON(e) : undefined);
    } else {
      obj.isArchivedTimeline = [];
    }
    if (message.defaultUserOutgoingApprovals) {
      obj.defaultUserOutgoingApprovals = message.defaultUserOutgoingApprovals.map((e) =>
        e ? UserOutgoingApproval.toJSON(e) : undefined
      );
    } else {
      obj.defaultUserOutgoingApprovals = [];
    }
    if (message.defaultUserIncomingApprovals) {
      obj.defaultUserIncomingApprovals = message.defaultUserIncomingApprovals.map((e) =>
        e ? UserIncomingApproval.toJSON(e) : undefined
      );
    } else {
      obj.defaultUserIncomingApprovals = [];
    }
    message.defaultUserPermissions !== undefined && (obj.defaultUserPermissions = message.defaultUserPermissions
      ? UserPermissions.toJSON(message.defaultUserPermissions)
      : undefined);
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedOutgoingTransfers =
        message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    message.defaultAutoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedIncomingTransfers =
        message.defaultAutoApproveSelfInitiatedIncomingTransfers);
    message.createdBy !== undefined && (obj.createdBy = message.createdBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BadgeCollection>, I>>(object: I): BadgeCollection {
    const message = createBaseBadgeCollection();
    message.collectionId = object.collectionId ?? "";
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.balancesType = object.balancesType ?? "";
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.customDataTimeline = object.customDataTimeline?.map((e) => CustomDataTimeline.fromPartial(e)) || [];
    message.managerTimeline = object.managerTimeline?.map((e) => ManagerTimeline.fromPartial(e)) || [];
    message.collectionPermissions =
      (object.collectionPermissions !== undefined && object.collectionPermissions !== null)
        ? CollectionPermissions.fromPartial(object.collectionPermissions)
        : undefined;
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    message.defaultUserOutgoingApprovals =
      object.defaultUserOutgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e)) || [];
    message.defaultUserIncomingApprovals =
      object.defaultUserIncomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e)) || [];
    message.defaultUserPermissions =
      (object.defaultUserPermissions !== undefined && object.defaultUserPermissions !== null)
        ? UserPermissions.fromPartial(object.defaultUserPermissions)
        : undefined;
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers = object.defaultAutoApproveSelfInitiatedOutgoingTransfers
      ?? false;
    message.defaultAutoApproveSelfInitiatedIncomingTransfers = object.defaultAutoApproveSelfInitiatedIncomingTransfers
      ?? false;
    message.createdBy = object.createdBy ?? "";
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
