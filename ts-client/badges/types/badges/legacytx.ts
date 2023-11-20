/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { AddressMapping } from "./address_mappings";
import { Balance, UintRange } from "./balances";
import { CollectionPermissions, UserPermissions } from "./permissions";
import {
  BadgeMetadataTimeline,
  CollectionMetadataTimeline,
  ContractAddressTimeline,
  CustomDataTimeline,
  IsArchivedTimeline,
  ManagerTimeline,
  OffChainBalancesMetadataTimeline,
  StandardsTimeline,
} from "./timelines";
import { CollectionApproval, Transfer, UserIncomingApproval, UserOutgoingApproval } from "./transfers";

export const protobufPackage = "badges";

/** The types defined in these files are used to define the MsgServer types for all requests and responses for Msgs of the badges module. */
export interface UserOutgoingApprovalTimeline {
  outgoingApprovals: UserOutgoingApproval[];
  timelineTimes: UintRange[];
}

export interface UserIncomingApprovalTimeline {
  incomingApprovals: UserIncomingApproval[];
  timelineTimes: UintRange[];
}

export interface CollectionApprovalTimeline {
  collectionApprovals: CollectionApproval[];
  timelineTimes: UintRange[];
}

export interface MsgNewCollection {
  /** See collections.proto for more details about these MsgNewBadge fields. Defines the badge details. Leave unneeded fields empty. */
  creator: string;
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  customDataTimeline: CustomDataTimeline[];
  balancesType: string;
  collectionApprovals: CollectionApproval[];
  permissions: CollectionPermissions | undefined;
  standardsTimeline: StandardsTimeline[];
  /**
   * Badge supplys and amounts to create. For each idx, we create amounts[idx] badges each with a supply of supplys[idx].
   * If supply[idx] == 0, we assume default supply. amountsToCreate[idx] can't equal 0.
   */
  badgesToCreate: Balance[];
  transfers: Transfer[];
  contractAddressTimeline: ContractAddressTimeline[];
  addressMappings: AddressMapping[];
  /** The user's approved transfers for each badge ID. */
  defaultOutgoingApprovals: UserOutgoingApproval[];
  /** The user's approved incoming transfers for each badge ID. */
  defaultIncomingApprovals: UserIncomingApproval[];
  defaultDisapproveSelfInitiated: boolean;
}

export interface MsgNewCollectionResponse {
  /** ID of created badge collecon */
  collectionId: string;
}

/** This handles both minting more of existing badges and creating new badges. */
export interface MsgMintAndDistributeBadges {
  creator: string;
  collectionId: string;
  badgesToCreate: Balance[];
  transfers: Transfer[];
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  collectionApprovals: CollectionApproval[];
  addressMappings: AddressMapping[];
}

export interface MsgMintAndDistributeBadgesResponse {
}

export interface MsgUniversalUpdateCollectionApprovals {
  creator: string;
  collectionId: string;
  collectionApprovals: CollectionApproval[];
  addressMappings: AddressMapping[];
}

export interface MsgUniversalUpdateCollectionApprovalsResponse {
}

/** Update badge Uris with new URI object, if permitted. */
export interface MsgUpdateMetadata {
  creator: string;
  collectionId: string;
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  customDataTimeline: CustomDataTimeline[];
  contractAddressTimeline: ContractAddressTimeline[];
  standardsTimeline: StandardsTimeline[];
}

export interface MsgUpdateMetadataResponse {
}

/** Update badge permissions with new permissions, if permitted. */
export interface MsgUniversalUpdateCollectionPermissions {
  creator: string;
  collectionId: string;
  permissions: CollectionPermissions | undefined;
  addressMappings: AddressMapping[];
}

export interface MsgUniversalUpdateCollectionPermissionsResponse {
}

export interface MsgUpdateUserPermissions {
  creator: string;
  collectionId: string;
  permissions: UserPermissions | undefined;
  addressMappings: AddressMapping[];
}

export interface MsgUpdateUserPermissionsResponse {
}

/** Transfer manager to this address. Recipient must have made a request. */
export interface MsgUpdateManager {
  creator: string;
  collectionId: string;
  managerTimeline: ManagerTimeline[];
}

export interface MsgUpdateManagerResponse {
}

export interface MsgArchiveCollection {
  creator: string;
  collectionId: string;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgArchiveCollectionResponse {
}

function createBaseUserOutgoingApprovalTimeline(): UserOutgoingApprovalTimeline {
  return { outgoingApprovals: [], timelineTimes: [] };
}

export const UserOutgoingApprovalTimeline = {
  encode(message: UserOutgoingApprovalTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.outgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserOutgoingApprovalTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserOutgoingApprovalTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.outgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserOutgoingApprovalTimeline {
    return {
      outgoingApprovals: Array.isArray(object?.outgoingApprovals)
        ? object.outgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserOutgoingApprovalTimeline): unknown {
    const obj: any = {};
    if (message.outgoingApprovals) {
      obj.outgoingApprovals = message.outgoingApprovals.map((e) => e ? UserOutgoingApproval.toJSON(e) : undefined);
    } else {
      obj.outgoingApprovals = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserOutgoingApprovalTimeline>, I>>(object: I): UserOutgoingApprovalTimeline {
    const message = createBaseUserOutgoingApprovalTimeline();
    message.outgoingApprovals = object.outgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserIncomingApprovalTimeline(): UserIncomingApprovalTimeline {
  return { incomingApprovals: [], timelineTimes: [] };
}

export const UserIncomingApprovalTimeline = {
  encode(message: UserIncomingApprovalTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.incomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserIncomingApprovalTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserIncomingApprovalTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.incomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserIncomingApprovalTimeline {
    return {
      incomingApprovals: Array.isArray(object?.incomingApprovals)
        ? object.incomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserIncomingApprovalTimeline): unknown {
    const obj: any = {};
    if (message.incomingApprovals) {
      obj.incomingApprovals = message.incomingApprovals.map((e) => e ? UserIncomingApproval.toJSON(e) : undefined);
    } else {
      obj.incomingApprovals = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserIncomingApprovalTimeline>, I>>(object: I): UserIncomingApprovalTimeline {
    const message = createBaseUserIncomingApprovalTimeline();
    message.incomingApprovals = object.incomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCollectionApprovalTimeline(): CollectionApprovalTimeline {
  return { collectionApprovals: [], timelineTimes: [] };
}

export const CollectionApprovalTimeline = {
  encode(message: CollectionApprovalTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovalTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovalTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovalTimeline {
    return {
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionApprovalTimeline): unknown {
    const obj: any = {};
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovalTimeline>, I>>(object: I): CollectionApprovalTimeline {
    const message = createBaseCollectionApprovalTimeline();
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgNewCollection(): MsgNewCollection {
  return {
    creator: "",
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    offChainBalancesMetadataTimeline: [],
    customDataTimeline: [],
    balancesType: "",
    collectionApprovals: [],
    permissions: undefined,
    standardsTimeline: [],
    badgesToCreate: [],
    transfers: [],
    contractAddressTimeline: [],
    addressMappings: [],
    defaultOutgoingApprovals: [],
    defaultIncomingApprovals: [],
    defaultDisapproveSelfInitiated: false,
  };
}

export const MsgNewCollection = {
  encode(message: MsgNewCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.collectionMetadataTimeline) {
      CollectionMetadataTimeline.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.badgeMetadataTimeline) {
      BadgeMetadataTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.offChainBalancesMetadataTimeline) {
      OffChainBalancesMetadataTimeline.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.customDataTimeline) {
      CustomDataTimeline.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.balancesType !== "") {
      writer.uint32(50).string(message.balancesType);
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    if (message.permissions !== undefined) {
      CollectionPermissions.encode(message.permissions, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.badgesToCreate) {
      Balance.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.transfers) {
      Transfer.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.contractAddressTimeline) {
      ContractAddressTimeline.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.defaultOutgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    for (const v of message.defaultIncomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    if (message.defaultDisapproveSelfInitiated === true) {
      writer.uint32(136).bool(message.defaultDisapproveSelfInitiated);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionMetadataTimeline.push(CollectionMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 3:
          message.badgeMetadataTimeline.push(BadgeMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 4:
          message.offChainBalancesMetadataTimeline.push(
            OffChainBalancesMetadataTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 5:
          message.customDataTimeline.push(CustomDataTimeline.decode(reader, reader.uint32()));
          break;
        case 6:
          message.balancesType = reader.string();
          break;
        case 8:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 9:
          message.permissions = CollectionPermissions.decode(reader, reader.uint32());
          break;
        case 10:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
          break;
        case 11:
          message.badgesToCreate.push(Balance.decode(reader, reader.uint32()));
          break;
        case 12:
          message.transfers.push(Transfer.decode(reader, reader.uint32()));
          break;
        case 13:
          message.contractAddressTimeline.push(ContractAddressTimeline.decode(reader, reader.uint32()));
          break;
        case 14:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        case 15:
          message.defaultOutgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 16:
          message.defaultIncomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 17:
          message.defaultDisapproveSelfInitiated = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      customDataTimeline: Array.isArray(object?.customDataTimeline)
        ? object.customDataTimeline.map((e: any) => CustomDataTimeline.fromJSON(e))
        : [],
      balancesType: isSet(object.balancesType) ? String(object.balancesType) : "",
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      permissions: isSet(object.permissions) ? CollectionPermissions.fromJSON(object.permissions) : undefined,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      badgesToCreate: Array.isArray(object?.badgesToCreate)
        ? object.badgesToCreate.map((e: any) => Balance.fromJSON(e))
        : [],
      transfers: Array.isArray(object?.transfers) ? object.transfers.map((e: any) => Transfer.fromJSON(e)) : [],
      contractAddressTimeline: Array.isArray(object?.contractAddressTimeline)
        ? object.contractAddressTimeline.map((e: any) => ContractAddressTimeline.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
      defaultOutgoingApprovals: Array.isArray(object?.defaultOutgoingApprovals)
        ? object.defaultOutgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      defaultIncomingApprovals: Array.isArray(object?.defaultIncomingApprovals)
        ? object.defaultIncomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      defaultDisapproveSelfInitiated: isSet(object.defaultDisapproveSelfInitiated)
        ? Boolean(object.defaultDisapproveSelfInitiated)
        : false,
    };
  },

  toJSON(message: MsgNewCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
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
    message.balancesType !== undefined && (obj.balancesType = message.balancesType);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    message.permissions !== undefined
      && (obj.permissions = message.permissions ? CollectionPermissions.toJSON(message.permissions) : undefined);
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    if (message.badgesToCreate) {
      obj.badgesToCreate = message.badgesToCreate.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.badgesToCreate = [];
    }
    if (message.transfers) {
      obj.transfers = message.transfers.map((e) => e ? Transfer.toJSON(e) : undefined);
    } else {
      obj.transfers = [];
    }
    if (message.contractAddressTimeline) {
      obj.contractAddressTimeline = message.contractAddressTimeline.map((e) =>
        e ? ContractAddressTimeline.toJSON(e) : undefined
      );
    } else {
      obj.contractAddressTimeline = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    if (message.defaultOutgoingApprovals) {
      obj.defaultOutgoingApprovals = message.defaultOutgoingApprovals.map((e) =>
        e ? UserOutgoingApproval.toJSON(e) : undefined
      );
    } else {
      obj.defaultOutgoingApprovals = [];
    }
    if (message.defaultIncomingApprovals) {
      obj.defaultIncomingApprovals = message.defaultIncomingApprovals.map((e) =>
        e ? UserIncomingApproval.toJSON(e) : undefined
      );
    } else {
      obj.defaultIncomingApprovals = [];
    }
    message.defaultDisapproveSelfInitiated !== undefined
      && (obj.defaultDisapproveSelfInitiated = message.defaultDisapproveSelfInitiated);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewCollection>, I>>(object: I): MsgNewCollection {
    const message = createBaseMsgNewCollection();
    message.creator = object.creator ?? "";
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.customDataTimeline = object.customDataTimeline?.map((e) => CustomDataTimeline.fromPartial(e)) || [];
    message.balancesType = object.balancesType ?? "";
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? CollectionPermissions.fromPartial(object.permissions)
      : undefined;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.badgesToCreate = object.badgesToCreate?.map((e) => Balance.fromPartial(e)) || [];
    message.transfers = object.transfers?.map((e) => Transfer.fromPartial(e)) || [];
    message.contractAddressTimeline = object.contractAddressTimeline?.map((e) => ContractAddressTimeline.fromPartial(e))
      || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    message.defaultOutgoingApprovals = object.defaultOutgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e))
      || [];
    message.defaultIncomingApprovals = object.defaultIncomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e))
      || [];
    message.defaultDisapproveSelfInitiated = object.defaultDisapproveSelfInitiated ?? false;
    return message;
  },
};

function createBaseMsgNewCollectionResponse(): MsgNewCollectionResponse {
  return { collectionId: "" };
}

export const MsgNewCollectionResponse = {
  encode(message: MsgNewCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewCollectionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewCollectionResponse {
    return { collectionId: isSet(object.collectionId) ? String(object.collectionId) : "" };
  },

  toJSON(message: MsgNewCollectionResponse): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewCollectionResponse>, I>>(object: I): MsgNewCollectionResponse {
    const message = createBaseMsgNewCollectionResponse();
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseMsgMintAndDistributeBadges(): MsgMintAndDistributeBadges {
  return {
    creator: "",
    collectionId: "",
    badgesToCreate: [],
    transfers: [],
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    offChainBalancesMetadataTimeline: [],
    collectionApprovals: [],
    addressMappings: [],
  };
}

export const MsgMintAndDistributeBadges = {
  encode(message: MsgMintAndDistributeBadges, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.badgesToCreate) {
      Balance.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.transfers) {
      Transfer.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.collectionMetadataTimeline) {
      CollectionMetadataTimeline.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.badgeMetadataTimeline) {
      BadgeMetadataTimeline.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.offChainBalancesMetadataTimeline) {
      OffChainBalancesMetadataTimeline.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMintAndDistributeBadges {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMintAndDistributeBadges();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.badgesToCreate.push(Balance.decode(reader, reader.uint32()));
          break;
        case 4:
          message.transfers.push(Transfer.decode(reader, reader.uint32()));
          break;
        case 6:
          message.collectionMetadataTimeline.push(CollectionMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 7:
          message.badgeMetadataTimeline.push(BadgeMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 8:
          message.offChainBalancesMetadataTimeline.push(
            OffChainBalancesMetadataTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 9:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 10:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgMintAndDistributeBadges {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      badgesToCreate: Array.isArray(object?.badgesToCreate)
        ? object.badgesToCreate.map((e: any) => Balance.fromJSON(e))
        : [],
      transfers: Array.isArray(object?.transfers) ? object.transfers.map((e: any) => Transfer.fromJSON(e)) : [],
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgMintAndDistributeBadges): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.badgesToCreate) {
      obj.badgesToCreate = message.badgesToCreate.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.badgesToCreate = [];
    }
    if (message.transfers) {
      obj.transfers = message.transfers.map((e) => e ? Transfer.toJSON(e) : undefined);
    } else {
      obj.transfers = [];
    }
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
    if (message.offChainBalancesMetadataTimeline) {
      obj.offChainBalancesMetadataTimeline = message.offChainBalancesMetadataTimeline.map((e) =>
        e ? OffChainBalancesMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.offChainBalancesMetadataTimeline = [];
    }
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMintAndDistributeBadges>, I>>(object: I): MsgMintAndDistributeBadges {
    const message = createBaseMsgMintAndDistributeBadges();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.badgesToCreate = object.badgesToCreate?.map((e) => Balance.fromPartial(e)) || [];
    message.transfers = object.transfers?.map((e) => Transfer.fromPartial(e)) || [];
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgMintAndDistributeBadgesResponse(): MsgMintAndDistributeBadgesResponse {
  return {};
}

export const MsgMintAndDistributeBadgesResponse = {
  encode(_: MsgMintAndDistributeBadgesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMintAndDistributeBadgesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMintAndDistributeBadgesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgMintAndDistributeBadgesResponse {
    return {};
  },

  toJSON(_: MsgMintAndDistributeBadgesResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMintAndDistributeBadgesResponse>, I>>(
    _: I,
  ): MsgMintAndDistributeBadgesResponse {
    const message = createBaseMsgMintAndDistributeBadgesResponse();
    return message;
  },
};

function createBaseMsgUniversalUpdateCollectionApprovals(): MsgUniversalUpdateCollectionApprovals {
  return { creator: "", collectionId: "", collectionApprovals: [], addressMappings: [] };
}

export const MsgUniversalUpdateCollectionApprovals = {
  encode(message: MsgUniversalUpdateCollectionApprovals, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollectionApprovals {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollectionApprovals();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 4:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUniversalUpdateCollectionApprovals {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUniversalUpdateCollectionApprovals): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollectionApprovals>, I>>(
    object: I,
  ): MsgUniversalUpdateCollectionApprovals {
    const message = createBaseMsgUniversalUpdateCollectionApprovals();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUniversalUpdateCollectionApprovalsResponse(): MsgUniversalUpdateCollectionApprovalsResponse {
  return {};
}

export const MsgUniversalUpdateCollectionApprovalsResponse = {
  encode(_: MsgUniversalUpdateCollectionApprovalsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollectionApprovalsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollectionApprovalsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUniversalUpdateCollectionApprovalsResponse {
    return {};
  },

  toJSON(_: MsgUniversalUpdateCollectionApprovalsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollectionApprovalsResponse>, I>>(
    _: I,
  ): MsgUniversalUpdateCollectionApprovalsResponse {
    const message = createBaseMsgUniversalUpdateCollectionApprovalsResponse();
    return message;
  },
};

function createBaseMsgUpdateMetadata(): MsgUpdateMetadata {
  return {
    creator: "",
    collectionId: "",
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    offChainBalancesMetadataTimeline: [],
    customDataTimeline: [],
    contractAddressTimeline: [],
    standardsTimeline: [],
  };
}

export const MsgUpdateMetadata = {
  encode(message: MsgUpdateMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.collectionMetadataTimeline) {
      CollectionMetadataTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeMetadataTimeline) {
      BadgeMetadataTimeline.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.offChainBalancesMetadataTimeline) {
      OffChainBalancesMetadataTimeline.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.customDataTimeline) {
      CustomDataTimeline.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.contractAddressTimeline) {
      ContractAddressTimeline.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.collectionMetadataTimeline.push(CollectionMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeMetadataTimeline.push(BadgeMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 5:
          message.offChainBalancesMetadataTimeline.push(
            OffChainBalancesMetadataTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 6:
          message.customDataTimeline.push(CustomDataTimeline.decode(reader, reader.uint32()));
          break;
        case 7:
          message.contractAddressTimeline.push(ContractAddressTimeline.decode(reader, reader.uint32()));
          break;
        case 8:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateMetadata {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      customDataTimeline: Array.isArray(object?.customDataTimeline)
        ? object.customDataTimeline.map((e: any) => CustomDataTimeline.fromJSON(e))
        : [],
      contractAddressTimeline: Array.isArray(object?.contractAddressTimeline)
        ? object.contractAddressTimeline.map((e: any) => ContractAddressTimeline.fromJSON(e))
        : [],
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateMetadata): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
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
    if (message.contractAddressTimeline) {
      obj.contractAddressTimeline = message.contractAddressTimeline.map((e) =>
        e ? ContractAddressTimeline.toJSON(e) : undefined
      );
    } else {
      obj.contractAddressTimeline = [];
    }
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateMetadata>, I>>(object: I): MsgUpdateMetadata {
    const message = createBaseMsgUpdateMetadata();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.customDataTimeline = object.customDataTimeline?.map((e) => CustomDataTimeline.fromPartial(e)) || [];
    message.contractAddressTimeline = object.contractAddressTimeline?.map((e) => ContractAddressTimeline.fromPartial(e))
      || [];
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateMetadataResponse(): MsgUpdateMetadataResponse {
  return {};
}

export const MsgUpdateMetadataResponse = {
  encode(_: MsgUpdateMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateMetadataResponse {
    return {};
  },

  toJSON(_: MsgUpdateMetadataResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateMetadataResponse>, I>>(_: I): MsgUpdateMetadataResponse {
    const message = createBaseMsgUpdateMetadataResponse();
    return message;
  },
};

function createBaseMsgUniversalUpdateCollectionPermissions(): MsgUniversalUpdateCollectionPermissions {
  return { creator: "", collectionId: "", permissions: undefined, addressMappings: [] };
}

export const MsgUniversalUpdateCollectionPermissions = {
  encode(message: MsgUniversalUpdateCollectionPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.permissions !== undefined) {
      CollectionPermissions.encode(message.permissions, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollectionPermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollectionPermissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.permissions = CollectionPermissions.decode(reader, reader.uint32());
          break;
        case 4:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUniversalUpdateCollectionPermissions {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      permissions: isSet(object.permissions) ? CollectionPermissions.fromJSON(object.permissions) : undefined,
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUniversalUpdateCollectionPermissions): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.permissions !== undefined
      && (obj.permissions = message.permissions ? CollectionPermissions.toJSON(message.permissions) : undefined);
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollectionPermissions>, I>>(
    object: I,
  ): MsgUniversalUpdateCollectionPermissions {
    const message = createBaseMsgUniversalUpdateCollectionPermissions();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? CollectionPermissions.fromPartial(object.permissions)
      : undefined;
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUniversalUpdateCollectionPermissionsResponse(): MsgUniversalUpdateCollectionPermissionsResponse {
  return {};
}

export const MsgUniversalUpdateCollectionPermissionsResponse = {
  encode(_: MsgUniversalUpdateCollectionPermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollectionPermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollectionPermissionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUniversalUpdateCollectionPermissionsResponse {
    return {};
  },

  toJSON(_: MsgUniversalUpdateCollectionPermissionsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollectionPermissionsResponse>, I>>(
    _: I,
  ): MsgUniversalUpdateCollectionPermissionsResponse {
    const message = createBaseMsgUniversalUpdateCollectionPermissionsResponse();
    return message;
  },
};

function createBaseMsgUpdateUserPermissions(): MsgUpdateUserPermissions {
  return { creator: "", collectionId: "", permissions: undefined, addressMappings: [] };
}

export const MsgUpdateUserPermissions = {
  encode(message: MsgUpdateUserPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.permissions !== undefined) {
      UserPermissions.encode(message.permissions, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserPermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserPermissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.permissions = UserPermissions.decode(reader, reader.uint32());
          break;
        case 4:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateUserPermissions {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      permissions: isSet(object.permissions) ? UserPermissions.fromJSON(object.permissions) : undefined,
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateUserPermissions): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.permissions !== undefined
      && (obj.permissions = message.permissions ? UserPermissions.toJSON(message.permissions) : undefined);
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserPermissions>, I>>(object: I): MsgUpdateUserPermissions {
    const message = createBaseMsgUpdateUserPermissions();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? UserPermissions.fromPartial(object.permissions)
      : undefined;
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateUserPermissionsResponse(): MsgUpdateUserPermissionsResponse {
  return {};
}

export const MsgUpdateUserPermissionsResponse = {
  encode(_: MsgUpdateUserPermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserPermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserPermissionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateUserPermissionsResponse {
    return {};
  },

  toJSON(_: MsgUpdateUserPermissionsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserPermissionsResponse>, I>>(
    _: I,
  ): MsgUpdateUserPermissionsResponse {
    const message = createBaseMsgUpdateUserPermissionsResponse();
    return message;
  },
};

function createBaseMsgUpdateManager(): MsgUpdateManager {
  return { creator: "", collectionId: "", managerTimeline: [] };
}

export const MsgUpdateManager = {
  encode(message: MsgUpdateManager, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.managerTimeline) {
      ManagerTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateManager {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateManager();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.managerTimeline.push(ManagerTimeline.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateManager {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      managerTimeline: Array.isArray(object?.managerTimeline)
        ? object.managerTimeline.map((e: any) => ManagerTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateManager): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.managerTimeline) {
      obj.managerTimeline = message.managerTimeline.map((e) => e ? ManagerTimeline.toJSON(e) : undefined);
    } else {
      obj.managerTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateManager>, I>>(object: I): MsgUpdateManager {
    const message = createBaseMsgUpdateManager();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.managerTimeline = object.managerTimeline?.map((e) => ManagerTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateManagerResponse(): MsgUpdateManagerResponse {
  return {};
}

export const MsgUpdateManagerResponse = {
  encode(_: MsgUpdateManagerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateManagerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateManagerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateManagerResponse {
    return {};
  },

  toJSON(_: MsgUpdateManagerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateManagerResponse>, I>>(_: I): MsgUpdateManagerResponse {
    const message = createBaseMsgUpdateManagerResponse();
    return message;
  },
};

function createBaseMsgArchiveCollection(): MsgArchiveCollection {
  return { creator: "", collectionId: "", isArchivedTimeline: [] };
}

export const MsgArchiveCollection = {
  encode(message: MsgArchiveCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgArchiveCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgArchiveCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        case 3:
          message.isArchivedTimeline.push(IsArchivedTimeline.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgArchiveCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgArchiveCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.isArchivedTimeline) {
      obj.isArchivedTimeline = message.isArchivedTimeline.map((e) => e ? IsArchivedTimeline.toJSON(e) : undefined);
    } else {
      obj.isArchivedTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgArchiveCollection>, I>>(object: I): MsgArchiveCollection {
    const message = createBaseMsgArchiveCollection();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgArchiveCollectionResponse(): MsgArchiveCollectionResponse {
  return {};
}

export const MsgArchiveCollectionResponse = {
  encode(_: MsgArchiveCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgArchiveCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgArchiveCollectionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgArchiveCollectionResponse {
    return {};
  },

  toJSON(_: MsgArchiveCollectionResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgArchiveCollectionResponse>, I>>(_: I): MsgArchiveCollectionResponse {
    const message = createBaseMsgArchiveCollectionResponse();
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
