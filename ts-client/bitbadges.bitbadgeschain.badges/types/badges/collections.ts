/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { CollectionPermissions, UserPermissions } from "./permissions";
import {
  BadgeMetadataTimeline,
  CollectionApprovedTransferTimeline,
  CollectionMetadataTimeline,
  ContractAddressTimeline,
  CustomDataTimeline,
  InheritedBalancesTimeline,
  IsArchivedTimeline,
  ManagerTimeline,
  OffChainBalancesMetadataTimeline,
  StandardsTimeline,
} from "./timelines";
import { UserApprovedIncomingTransferTimeline, UserApprovedOutgoingTransferTimeline } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

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
  /** The inherited balances metadata defines the parent balances for collections with inherited balances. */
  inheritedBalancesTimeline: InheritedBalancesTimeline[];
  /** The custom data field is an arbitrary field that can be used to store any data. */
  customDataTimeline: CustomDataTimeline[];
  /** The manager is the address of the manager of this collection. */
  managerTimeline: ManagerTimeline[];
  /** The permissions define what the manager of the collection can do or not do. */
  collectionPermissions:
    | CollectionPermissions
    | undefined;
  /**
   * The approved transfers timeline defines the transferability of the collection for collections with standard balances.
   * This defines it on a collection-level. All transfers must be explicitly allowed on the collection-level, or else, they will fail.
   *
   * Collection approved transfers can optionally specify to override the user approvals for a transfer (e.g. forcefully revoke a badge).
   * If user approvals are not overriden, then a transfer must also satisfy the From user's approved outgoing transfers and the To user's approved incoming transfers.
   */
  collectionApprovedTransfersTimeline: CollectionApprovedTransferTimeline[];
  /** Standards allow us to define a standard for the collection. This lets others know how to interpret the fields of the collection. */
  standardsTimeline: StandardsTimeline[];
  /**
   * The isArchivedTimeline defines whether the collection is archived or not.
   * When a collection is archived, it is read-only and no transactions can be processed.
   */
  isArchivedTimeline: IsArchivedTimeline[];
  /** The contractAddressTimeline defines the contract address for the collection (if it has a corresponding contract). */
  contractAddressTimeline: ContractAddressTimeline[];
  /**
   * The defaultUserApprovedOutgoingTransfersTimeline defines the default user approved outgoing transfers for an uninitialized user balance.
   * The user can change this value at any time.
   */
  defaultUserApprovedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  /**
   * The defaultUserApprovedIncomingTransfersTimeline defines the default user approved incoming transfers for an uninitialized user balance.
   * The user can change this value at any time.
   *
   * Ex: Set this to disallow all incoming transfers by default, making the user have to opt-in to receiving the badge.
   */
  defaultUserApprovedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
  defaultUserPermissions: UserPermissions | undefined;
  createdBy: string;
}

function createBaseBadgeCollection(): BadgeCollection {
  return {
    collectionId: "",
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    balancesType: "",
    offChainBalancesMetadataTimeline: [],
    inheritedBalancesTimeline: [],
    customDataTimeline: [],
    managerTimeline: [],
    collectionPermissions: undefined,
    collectionApprovedTransfersTimeline: [],
    standardsTimeline: [],
    isArchivedTimeline: [],
    contractAddressTimeline: [],
    defaultUserApprovedOutgoingTransfersTimeline: [],
    defaultUserApprovedIncomingTransfersTimeline: [],
    defaultUserPermissions: undefined,
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
    for (const v of message.inheritedBalancesTimeline) {
      InheritedBalancesTimeline.encode(v!, writer.uint32(50).fork()).ldelim();
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
    for (const v of message.collectionApprovedTransfersTimeline) {
      CollectionApprovedTransferTimeline.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.contractAddressTimeline) {
      ContractAddressTimeline.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.defaultUserApprovedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.defaultUserApprovedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    if (message.defaultUserPermissions !== undefined) {
      UserPermissions.encode(message.defaultUserPermissions, writer.uint32(130).fork()).ldelim();
    }
    if (message.createdBy !== "") {
      writer.uint32(138).string(message.createdBy);
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
        case 6:
          message.inheritedBalancesTimeline.push(InheritedBalancesTimeline.decode(reader, reader.uint32()));
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
          message.collectionApprovedTransfersTimeline.push(
            CollectionApprovedTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 11:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
          break;
        case 12:
          message.isArchivedTimeline.push(IsArchivedTimeline.decode(reader, reader.uint32()));
          break;
        case 13:
          message.contractAddressTimeline.push(ContractAddressTimeline.decode(reader, reader.uint32()));
          break;
        case 14:
          message.defaultUserApprovedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 15:
          message.defaultUserApprovedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 16:
          message.defaultUserPermissions = UserPermissions.decode(reader, reader.uint32());
          break;
        case 17:
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
      inheritedBalancesTimeline: Array.isArray(object?.inheritedBalancesTimeline)
        ? object.inheritedBalancesTimeline.map((e: any) => InheritedBalancesTimeline.fromJSON(e))
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
      collectionApprovedTransfersTimeline: Array.isArray(object?.collectionApprovedTransfersTimeline)
        ? object.collectionApprovedTransfersTimeline.map((e: any) => CollectionApprovedTransferTimeline.fromJSON(e))
        : [],
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
      contractAddressTimeline: Array.isArray(object?.contractAddressTimeline)
        ? object.contractAddressTimeline.map((e: any) => ContractAddressTimeline.fromJSON(e))
        : [],
      defaultUserApprovedOutgoingTransfersTimeline: Array.isArray(object?.defaultUserApprovedOutgoingTransfersTimeline)
        ? object.defaultUserApprovedOutgoingTransfersTimeline.map((e: any) =>
          UserApprovedOutgoingTransferTimeline.fromJSON(e)
        )
        : [],
      defaultUserApprovedIncomingTransfersTimeline: Array.isArray(object?.defaultUserApprovedIncomingTransfersTimeline)
        ? object.defaultUserApprovedIncomingTransfersTimeline.map((e: any) =>
          UserApprovedIncomingTransferTimeline.fromJSON(e)
        )
        : [],
      defaultUserPermissions: isSet(object.defaultUserPermissions)
        ? UserPermissions.fromJSON(object.defaultUserPermissions)
        : undefined,
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
    if (message.inheritedBalancesTimeline) {
      obj.inheritedBalancesTimeline = message.inheritedBalancesTimeline.map((e) =>
        e ? InheritedBalancesTimeline.toJSON(e) : undefined
      );
    } else {
      obj.inheritedBalancesTimeline = [];
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
    if (message.collectionApprovedTransfersTimeline) {
      obj.collectionApprovedTransfersTimeline = message.collectionApprovedTransfersTimeline.map((e) =>
        e ? CollectionApprovedTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionApprovedTransfersTimeline = [];
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
    if (message.contractAddressTimeline) {
      obj.contractAddressTimeline = message.contractAddressTimeline.map((e) =>
        e ? ContractAddressTimeline.toJSON(e) : undefined
      );
    } else {
      obj.contractAddressTimeline = [];
    }
    if (message.defaultUserApprovedOutgoingTransfersTimeline) {
      obj.defaultUserApprovedOutgoingTransfersTimeline = message.defaultUserApprovedOutgoingTransfersTimeline.map((e) =>
        e ? UserApprovedOutgoingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.defaultUserApprovedOutgoingTransfersTimeline = [];
    }
    if (message.defaultUserApprovedIncomingTransfersTimeline) {
      obj.defaultUserApprovedIncomingTransfersTimeline = message.defaultUserApprovedIncomingTransfersTimeline.map((e) =>
        e ? UserApprovedIncomingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.defaultUserApprovedIncomingTransfersTimeline = [];
    }
    message.defaultUserPermissions !== undefined && (obj.defaultUserPermissions = message.defaultUserPermissions
      ? UserPermissions.toJSON(message.defaultUserPermissions)
      : undefined);
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
    message.inheritedBalancesTimeline =
      object.inheritedBalancesTimeline?.map((e) => InheritedBalancesTimeline.fromPartial(e)) || [];
    message.customDataTimeline = object.customDataTimeline?.map((e) => CustomDataTimeline.fromPartial(e)) || [];
    message.managerTimeline = object.managerTimeline?.map((e) => ManagerTimeline.fromPartial(e)) || [];
    message.collectionPermissions =
      (object.collectionPermissions !== undefined && object.collectionPermissions !== null)
        ? CollectionPermissions.fromPartial(object.collectionPermissions)
        : undefined;
    message.collectionApprovedTransfersTimeline =
      object.collectionApprovedTransfersTimeline?.map((e) => CollectionApprovedTransferTimeline.fromPartial(e)) || [];
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    message.contractAddressTimeline = object.contractAddressTimeline?.map((e) => ContractAddressTimeline.fromPartial(e))
      || [];
    message.defaultUserApprovedOutgoingTransfersTimeline =
      object.defaultUserApprovedOutgoingTransfersTimeline?.map((e) =>
        UserApprovedOutgoingTransferTimeline.fromPartial(e)
      ) || [];
    message.defaultUserApprovedIncomingTransfersTimeline =
      object.defaultUserApprovedIncomingTransfersTimeline?.map((e) =>
        UserApprovedIncomingTransferTimeline.fromPartial(e)
      ) || [];
    message.defaultUserPermissions =
      (object.defaultUserPermissions !== undefined && object.defaultUserPermissions !== null)
        ? UserPermissions.fromPartial(object.defaultUserPermissions)
        : undefined;
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
