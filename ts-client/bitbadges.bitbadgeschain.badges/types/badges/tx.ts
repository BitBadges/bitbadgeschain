/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { AddressMapping } from "./address_mappings";
import { Balance } from "./balances";
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
import { Transfer, UserApprovedIncomingTransferTimeline, UserApprovedOutgoingTransferTimeline } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

export interface MsgUpdateCollection {
  creator: string;
  /** 0 for new collection */
  collectionId: string;
  /** The following section of fields are only allowed to be set upon creation of a new collection. */
  balancesType: string;
  /** The user's approved transfers for each badge ID. */
  defaultApprovedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  /** The user's approved incoming transfers for each badge ID. */
  defaultApprovedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
  defaultUserPermissions:
    | UserPermissions
    | undefined;
  /** The rest of the fields are allowed to be set on creation or update. */
  badgesToCreate: Balance[];
  updateCollectionPermissions: boolean;
  collectionPermissions: CollectionPermissions | undefined;
  updateManagerTimeline: boolean;
  managerTimeline: ManagerTimeline[];
  updateCollectionMetadataTimeline: boolean;
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  updateBadgeMetadataTimeline: boolean;
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  updateOffChainBalancesMetadataTimeline: boolean;
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  updateCustomDataTimeline: boolean;
  customDataTimeline: CustomDataTimeline[];
  updateInheritedBalancesTimeline: boolean;
  inheritedBalancesTimeline: InheritedBalancesTimeline[];
  updateCollectionApprovedTransfersTimeline: boolean;
  collectionApprovedTransfersTimeline: CollectionApprovedTransferTimeline[];
  updateStandardsTimeline: boolean;
  standardsTimeline: StandardsTimeline[];
  updateContractAddressTimeline: boolean;
  contractAddressTimeline: ContractAddressTimeline[];
  updateIsArchivedTimeline: boolean;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgUpdateCollectionResponse {
  /** ID of badge collection */
  collectionId: string;
}

export interface MsgCreateAddressMappings {
  creator: string;
  addressMappings: AddressMapping[];
}

export interface MsgCreateAddressMappingsResponse {
}

/** For each amount, for each toAddress, we will attempt to transfer all the badgeIds for the badge with ID badgeId. */
export interface MsgTransferBadges {
  creator: string;
  collectionId: string;
  transfers: Transfer[];
}

export interface MsgTransferBadgesResponse {
}

export interface MsgDeleteCollection {
  creator: string;
  collectionId: string;
}

export interface MsgDeleteCollectionResponse {
}

export interface MsgUpdateUserApprovedTransfers {
  creator: string;
  collectionId: string;
  updateApprovedOutgoingTransfersTimeline: boolean;
  approvedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  updateApprovedIncomingTransfersTimeline: boolean;
  approvedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
  updateUserPermissions: boolean;
  userPermissions: UserPermissions | undefined;
}

export interface MsgUpdateUserApprovedTransfersResponse {
}

function createBaseMsgUpdateCollection(): MsgUpdateCollection {
  return {
    creator: "",
    collectionId: "",
    balancesType: "",
    defaultApprovedOutgoingTransfersTimeline: [],
    defaultApprovedIncomingTransfersTimeline: [],
    defaultUserPermissions: undefined,
    badgesToCreate: [],
    updateCollectionPermissions: false,
    collectionPermissions: undefined,
    updateManagerTimeline: false,
    managerTimeline: [],
    updateCollectionMetadataTimeline: false,
    collectionMetadataTimeline: [],
    updateBadgeMetadataTimeline: false,
    badgeMetadataTimeline: [],
    updateOffChainBalancesMetadataTimeline: false,
    offChainBalancesMetadataTimeline: [],
    updateCustomDataTimeline: false,
    customDataTimeline: [],
    updateInheritedBalancesTimeline: false,
    inheritedBalancesTimeline: [],
    updateCollectionApprovedTransfersTimeline: false,
    collectionApprovedTransfersTimeline: [],
    updateStandardsTimeline: false,
    standardsTimeline: [],
    updateContractAddressTimeline: false,
    contractAddressTimeline: [],
    updateIsArchivedTimeline: false,
    isArchivedTimeline: [],
  };
}

export const MsgUpdateCollection = {
  encode(message: MsgUpdateCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.balancesType !== "") {
      writer.uint32(26).string(message.balancesType);
    }
    for (const v of message.defaultApprovedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.defaultApprovedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.defaultUserPermissions !== undefined) {
      UserPermissions.encode(message.defaultUserPermissions, writer.uint32(234).fork()).ldelim();
    }
    for (const v of message.badgesToCreate) {
      Balance.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.updateCollectionPermissions === true) {
      writer.uint32(56).bool(message.updateCollectionPermissions);
    }
    if (message.collectionPermissions !== undefined) {
      CollectionPermissions.encode(message.collectionPermissions, writer.uint32(66).fork()).ldelim();
    }
    if (message.updateManagerTimeline === true) {
      writer.uint32(72).bool(message.updateManagerTimeline);
    }
    for (const v of message.managerTimeline) {
      ManagerTimeline.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    if (message.updateCollectionMetadataTimeline === true) {
      writer.uint32(88).bool(message.updateCollectionMetadataTimeline);
    }
    for (const v of message.collectionMetadataTimeline) {
      CollectionMetadataTimeline.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    if (message.updateBadgeMetadataTimeline === true) {
      writer.uint32(104).bool(message.updateBadgeMetadataTimeline);
    }
    for (const v of message.badgeMetadataTimeline) {
      BadgeMetadataTimeline.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    if (message.updateOffChainBalancesMetadataTimeline === true) {
      writer.uint32(120).bool(message.updateOffChainBalancesMetadataTimeline);
    }
    for (const v of message.offChainBalancesMetadataTimeline) {
      OffChainBalancesMetadataTimeline.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    if (message.updateCustomDataTimeline === true) {
      writer.uint32(136).bool(message.updateCustomDataTimeline);
    }
    for (const v of message.customDataTimeline) {
      CustomDataTimeline.encode(v!, writer.uint32(146).fork()).ldelim();
    }
    if (message.updateInheritedBalancesTimeline === true) {
      writer.uint32(152).bool(message.updateInheritedBalancesTimeline);
    }
    for (const v of message.inheritedBalancesTimeline) {
      InheritedBalancesTimeline.encode(v!, writer.uint32(162).fork()).ldelim();
    }
    if (message.updateCollectionApprovedTransfersTimeline === true) {
      writer.uint32(168).bool(message.updateCollectionApprovedTransfersTimeline);
    }
    for (const v of message.collectionApprovedTransfersTimeline) {
      CollectionApprovedTransferTimeline.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    if (message.updateStandardsTimeline === true) {
      writer.uint32(184).bool(message.updateStandardsTimeline);
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(194).fork()).ldelim();
    }
    if (message.updateContractAddressTimeline === true) {
      writer.uint32(200).bool(message.updateContractAddressTimeline);
    }
    for (const v of message.contractAddressTimeline) {
      ContractAddressTimeline.encode(v!, writer.uint32(210).fork()).ldelim();
    }
    if (message.updateIsArchivedTimeline === true) {
      writer.uint32(216).bool(message.updateIsArchivedTimeline);
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(226).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollection();
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
          message.balancesType = reader.string();
          break;
        case 4:
          message.defaultApprovedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 5:
          message.defaultApprovedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 29:
          message.defaultUserPermissions = UserPermissions.decode(reader, reader.uint32());
          break;
        case 6:
          message.badgesToCreate.push(Balance.decode(reader, reader.uint32()));
          break;
        case 7:
          message.updateCollectionPermissions = reader.bool();
          break;
        case 8:
          message.collectionPermissions = CollectionPermissions.decode(reader, reader.uint32());
          break;
        case 9:
          message.updateManagerTimeline = reader.bool();
          break;
        case 10:
          message.managerTimeline.push(ManagerTimeline.decode(reader, reader.uint32()));
          break;
        case 11:
          message.updateCollectionMetadataTimeline = reader.bool();
          break;
        case 12:
          message.collectionMetadataTimeline.push(CollectionMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 13:
          message.updateBadgeMetadataTimeline = reader.bool();
          break;
        case 14:
          message.badgeMetadataTimeline.push(BadgeMetadataTimeline.decode(reader, reader.uint32()));
          break;
        case 15:
          message.updateOffChainBalancesMetadataTimeline = reader.bool();
          break;
        case 16:
          message.offChainBalancesMetadataTimeline.push(
            OffChainBalancesMetadataTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 17:
          message.updateCustomDataTimeline = reader.bool();
          break;
        case 18:
          message.customDataTimeline.push(CustomDataTimeline.decode(reader, reader.uint32()));
          break;
        case 19:
          message.updateInheritedBalancesTimeline = reader.bool();
          break;
        case 20:
          message.inheritedBalancesTimeline.push(InheritedBalancesTimeline.decode(reader, reader.uint32()));
          break;
        case 21:
          message.updateCollectionApprovedTransfersTimeline = reader.bool();
          break;
        case 22:
          message.collectionApprovedTransfersTimeline.push(
            CollectionApprovedTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 23:
          message.updateStandardsTimeline = reader.bool();
          break;
        case 24:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
          break;
        case 25:
          message.updateContractAddressTimeline = reader.bool();
          break;
        case 26:
          message.contractAddressTimeline.push(ContractAddressTimeline.decode(reader, reader.uint32()));
          break;
        case 27:
          message.updateIsArchivedTimeline = reader.bool();
          break;
        case 28:
          message.isArchivedTimeline.push(IsArchivedTimeline.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      balancesType: isSet(object.balancesType) ? String(object.balancesType) : "",
      defaultApprovedOutgoingTransfersTimeline: Array.isArray(object?.defaultApprovedOutgoingTransfersTimeline)
        ? object.defaultApprovedOutgoingTransfersTimeline.map((e: any) =>
          UserApprovedOutgoingTransferTimeline.fromJSON(e)
        )
        : [],
      defaultApprovedIncomingTransfersTimeline: Array.isArray(object?.defaultApprovedIncomingTransfersTimeline)
        ? object.defaultApprovedIncomingTransfersTimeline.map((e: any) =>
          UserApprovedIncomingTransferTimeline.fromJSON(e)
        )
        : [],
      defaultUserPermissions: isSet(object.defaultUserPermissions)
        ? UserPermissions.fromJSON(object.defaultUserPermissions)
        : undefined,
      badgesToCreate: Array.isArray(object?.badgesToCreate)
        ? object.badgesToCreate.map((e: any) => Balance.fromJSON(e))
        : [],
      updateCollectionPermissions: isSet(object.updateCollectionPermissions)
        ? Boolean(object.updateCollectionPermissions)
        : false,
      collectionPermissions: isSet(object.collectionPermissions)
        ? CollectionPermissions.fromJSON(object.collectionPermissions)
        : undefined,
      updateManagerTimeline: isSet(object.updateManagerTimeline) ? Boolean(object.updateManagerTimeline) : false,
      managerTimeline: Array.isArray(object?.managerTimeline)
        ? object.managerTimeline.map((e: any) => ManagerTimeline.fromJSON(e))
        : [],
      updateCollectionMetadataTimeline: isSet(object.updateCollectionMetadataTimeline)
        ? Boolean(object.updateCollectionMetadataTimeline)
        : false,
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      updateBadgeMetadataTimeline: isSet(object.updateBadgeMetadataTimeline)
        ? Boolean(object.updateBadgeMetadataTimeline)
        : false,
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      updateOffChainBalancesMetadataTimeline: isSet(object.updateOffChainBalancesMetadataTimeline)
        ? Boolean(object.updateOffChainBalancesMetadataTimeline)
        : false,
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      updateCustomDataTimeline: isSet(object.updateCustomDataTimeline)
        ? Boolean(object.updateCustomDataTimeline)
        : false,
      customDataTimeline: Array.isArray(object?.customDataTimeline)
        ? object.customDataTimeline.map((e: any) => CustomDataTimeline.fromJSON(e))
        : [],
      updateInheritedBalancesTimeline: isSet(object.updateInheritedBalancesTimeline)
        ? Boolean(object.updateInheritedBalancesTimeline)
        : false,
      inheritedBalancesTimeline: Array.isArray(object?.inheritedBalancesTimeline)
        ? object.inheritedBalancesTimeline.map((e: any) => InheritedBalancesTimeline.fromJSON(e))
        : [],
      updateCollectionApprovedTransfersTimeline: isSet(object.updateCollectionApprovedTransfersTimeline)
        ? Boolean(object.updateCollectionApprovedTransfersTimeline)
        : false,
      collectionApprovedTransfersTimeline: Array.isArray(object?.collectionApprovedTransfersTimeline)
        ? object.collectionApprovedTransfersTimeline.map((e: any) => CollectionApprovedTransferTimeline.fromJSON(e))
        : [],
      updateStandardsTimeline: isSet(object.updateStandardsTimeline) ? Boolean(object.updateStandardsTimeline) : false,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      updateContractAddressTimeline: isSet(object.updateContractAddressTimeline)
        ? Boolean(object.updateContractAddressTimeline)
        : false,
      contractAddressTimeline: Array.isArray(object?.contractAddressTimeline)
        ? object.contractAddressTimeline.map((e: any) => ContractAddressTimeline.fromJSON(e))
        : [],
      updateIsArchivedTimeline: isSet(object.updateIsArchivedTimeline)
        ? Boolean(object.updateIsArchivedTimeline)
        : false,
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.balancesType !== undefined && (obj.balancesType = message.balancesType);
    if (message.defaultApprovedOutgoingTransfersTimeline) {
      obj.defaultApprovedOutgoingTransfersTimeline = message.defaultApprovedOutgoingTransfersTimeline.map((e) =>
        e ? UserApprovedOutgoingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.defaultApprovedOutgoingTransfersTimeline = [];
    }
    if (message.defaultApprovedIncomingTransfersTimeline) {
      obj.defaultApprovedIncomingTransfersTimeline = message.defaultApprovedIncomingTransfersTimeline.map((e) =>
        e ? UserApprovedIncomingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.defaultApprovedIncomingTransfersTimeline = [];
    }
    message.defaultUserPermissions !== undefined && (obj.defaultUserPermissions = message.defaultUserPermissions
      ? UserPermissions.toJSON(message.defaultUserPermissions)
      : undefined);
    if (message.badgesToCreate) {
      obj.badgesToCreate = message.badgesToCreate.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.badgesToCreate = [];
    }
    message.updateCollectionPermissions !== undefined
      && (obj.updateCollectionPermissions = message.updateCollectionPermissions);
    message.collectionPermissions !== undefined && (obj.collectionPermissions = message.collectionPermissions
      ? CollectionPermissions.toJSON(message.collectionPermissions)
      : undefined);
    message.updateManagerTimeline !== undefined && (obj.updateManagerTimeline = message.updateManagerTimeline);
    if (message.managerTimeline) {
      obj.managerTimeline = message.managerTimeline.map((e) => e ? ManagerTimeline.toJSON(e) : undefined);
    } else {
      obj.managerTimeline = [];
    }
    message.updateCollectionMetadataTimeline !== undefined
      && (obj.updateCollectionMetadataTimeline = message.updateCollectionMetadataTimeline);
    if (message.collectionMetadataTimeline) {
      obj.collectionMetadataTimeline = message.collectionMetadataTimeline.map((e) =>
        e ? CollectionMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionMetadataTimeline = [];
    }
    message.updateBadgeMetadataTimeline !== undefined
      && (obj.updateBadgeMetadataTimeline = message.updateBadgeMetadataTimeline);
    if (message.badgeMetadataTimeline) {
      obj.badgeMetadataTimeline = message.badgeMetadataTimeline.map((e) =>
        e ? BadgeMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.badgeMetadataTimeline = [];
    }
    message.updateOffChainBalancesMetadataTimeline !== undefined
      && (obj.updateOffChainBalancesMetadataTimeline = message.updateOffChainBalancesMetadataTimeline);
    if (message.offChainBalancesMetadataTimeline) {
      obj.offChainBalancesMetadataTimeline = message.offChainBalancesMetadataTimeline.map((e) =>
        e ? OffChainBalancesMetadataTimeline.toJSON(e) : undefined
      );
    } else {
      obj.offChainBalancesMetadataTimeline = [];
    }
    message.updateCustomDataTimeline !== undefined && (obj.updateCustomDataTimeline = message.updateCustomDataTimeline);
    if (message.customDataTimeline) {
      obj.customDataTimeline = message.customDataTimeline.map((e) => e ? CustomDataTimeline.toJSON(e) : undefined);
    } else {
      obj.customDataTimeline = [];
    }
    message.updateInheritedBalancesTimeline !== undefined
      && (obj.updateInheritedBalancesTimeline = message.updateInheritedBalancesTimeline);
    if (message.inheritedBalancesTimeline) {
      obj.inheritedBalancesTimeline = message.inheritedBalancesTimeline.map((e) =>
        e ? InheritedBalancesTimeline.toJSON(e) : undefined
      );
    } else {
      obj.inheritedBalancesTimeline = [];
    }
    message.updateCollectionApprovedTransfersTimeline !== undefined
      && (obj.updateCollectionApprovedTransfersTimeline = message.updateCollectionApprovedTransfersTimeline);
    if (message.collectionApprovedTransfersTimeline) {
      obj.collectionApprovedTransfersTimeline = message.collectionApprovedTransfersTimeline.map((e) =>
        e ? CollectionApprovedTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionApprovedTransfersTimeline = [];
    }
    message.updateStandardsTimeline !== undefined && (obj.updateStandardsTimeline = message.updateStandardsTimeline);
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    message.updateContractAddressTimeline !== undefined
      && (obj.updateContractAddressTimeline = message.updateContractAddressTimeline);
    if (message.contractAddressTimeline) {
      obj.contractAddressTimeline = message.contractAddressTimeline.map((e) =>
        e ? ContractAddressTimeline.toJSON(e) : undefined
      );
    } else {
      obj.contractAddressTimeline = [];
    }
    message.updateIsArchivedTimeline !== undefined && (obj.updateIsArchivedTimeline = message.updateIsArchivedTimeline);
    if (message.isArchivedTimeline) {
      obj.isArchivedTimeline = message.isArchivedTimeline.map((e) => e ? IsArchivedTimeline.toJSON(e) : undefined);
    } else {
      obj.isArchivedTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollection>, I>>(object: I): MsgUpdateCollection {
    const message = createBaseMsgUpdateCollection();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.balancesType = object.balancesType ?? "";
    message.defaultApprovedOutgoingTransfersTimeline =
      object.defaultApprovedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e))
      || [];
    message.defaultApprovedIncomingTransfersTimeline =
      object.defaultApprovedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e))
      || [];
    message.defaultUserPermissions =
      (object.defaultUserPermissions !== undefined && object.defaultUserPermissions !== null)
        ? UserPermissions.fromPartial(object.defaultUserPermissions)
        : undefined;
    message.badgesToCreate = object.badgesToCreate?.map((e) => Balance.fromPartial(e)) || [];
    message.updateCollectionPermissions = object.updateCollectionPermissions ?? false;
    message.collectionPermissions =
      (object.collectionPermissions !== undefined && object.collectionPermissions !== null)
        ? CollectionPermissions.fromPartial(object.collectionPermissions)
        : undefined;
    message.updateManagerTimeline = object.updateManagerTimeline ?? false;
    message.managerTimeline = object.managerTimeline?.map((e) => ManagerTimeline.fromPartial(e)) || [];
    message.updateCollectionMetadataTimeline = object.updateCollectionMetadataTimeline ?? false;
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.updateBadgeMetadataTimeline = object.updateBadgeMetadataTimeline ?? false;
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.updateOffChainBalancesMetadataTimeline = object.updateOffChainBalancesMetadataTimeline ?? false;
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.updateCustomDataTimeline = object.updateCustomDataTimeline ?? false;
    message.customDataTimeline = object.customDataTimeline?.map((e) => CustomDataTimeline.fromPartial(e)) || [];
    message.updateInheritedBalancesTimeline = object.updateInheritedBalancesTimeline ?? false;
    message.inheritedBalancesTimeline =
      object.inheritedBalancesTimeline?.map((e) => InheritedBalancesTimeline.fromPartial(e)) || [];
    message.updateCollectionApprovedTransfersTimeline = object.updateCollectionApprovedTransfersTimeline ?? false;
    message.collectionApprovedTransfersTimeline =
      object.collectionApprovedTransfersTimeline?.map((e) => CollectionApprovedTransferTimeline.fromPartial(e)) || [];
    message.updateStandardsTimeline = object.updateStandardsTimeline ?? false;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.updateContractAddressTimeline = object.updateContractAddressTimeline ?? false;
    message.contractAddressTimeline = object.contractAddressTimeline?.map((e) => ContractAddressTimeline.fromPartial(e))
      || [];
    message.updateIsArchivedTimeline = object.updateIsArchivedTimeline ?? false;
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateCollectionResponse(): MsgUpdateCollectionResponse {
  return { collectionId: "" };
}

export const MsgUpdateCollectionResponse = {
  encode(message: MsgUpdateCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollectionResponse();
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

  fromJSON(object: any): MsgUpdateCollectionResponse {
    return { collectionId: isSet(object.collectionId) ? String(object.collectionId) : "" };
  },

  toJSON(message: MsgUpdateCollectionResponse): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollectionResponse>, I>>(object: I): MsgUpdateCollectionResponse {
    const message = createBaseMsgUpdateCollectionResponse();
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseMsgCreateAddressMappings(): MsgCreateAddressMappings {
  return { creator: "", addressMappings: [] };
}

export const MsgCreateAddressMappings = {
  encode(message: MsgCreateAddressMappings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateAddressMappings {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateAddressMappings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateAddressMappings {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgCreateAddressMappings): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateAddressMappings>, I>>(object: I): MsgCreateAddressMappings {
    const message = createBaseMsgCreateAddressMappings();
    message.creator = object.creator ?? "";
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgCreateAddressMappingsResponse(): MsgCreateAddressMappingsResponse {
  return {};
}

export const MsgCreateAddressMappingsResponse = {
  encode(_: MsgCreateAddressMappingsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateAddressMappingsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateAddressMappingsResponse();
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

  fromJSON(_: any): MsgCreateAddressMappingsResponse {
    return {};
  },

  toJSON(_: MsgCreateAddressMappingsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateAddressMappingsResponse>, I>>(
    _: I,
  ): MsgCreateAddressMappingsResponse {
    const message = createBaseMsgCreateAddressMappingsResponse();
    return message;
  },
};

function createBaseMsgTransferBadges(): MsgTransferBadges {
  return { creator: "", collectionId: "", transfers: [] };
}

export const MsgTransferBadges = {
  encode(message: MsgTransferBadges, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.transfers) {
      Transfer.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferBadges {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferBadges();
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
          message.transfers.push(Transfer.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgTransferBadges {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      transfers: Array.isArray(object?.transfers) ? object.transfers.map((e: any) => Transfer.fromJSON(e)) : [],
    };
  },

  toJSON(message: MsgTransferBadges): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.transfers) {
      obj.transfers = message.transfers.map((e) => e ? Transfer.toJSON(e) : undefined);
    } else {
      obj.transfers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferBadges>, I>>(object: I): MsgTransferBadges {
    const message = createBaseMsgTransferBadges();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.transfers = object.transfers?.map((e) => Transfer.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgTransferBadgesResponse(): MsgTransferBadgesResponse {
  return {};
}

export const MsgTransferBadgesResponse = {
  encode(_: MsgTransferBadgesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferBadgesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferBadgesResponse();
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

  fromJSON(_: any): MsgTransferBadgesResponse {
    return {};
  },

  toJSON(_: MsgTransferBadgesResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferBadgesResponse>, I>>(_: I): MsgTransferBadgesResponse {
    const message = createBaseMsgTransferBadgesResponse();
    return message;
  },
};

function createBaseMsgDeleteCollection(): MsgDeleteCollection {
  return { creator: "", collectionId: "" };
}

export const MsgDeleteCollection = {
  encode(message: MsgDeleteCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.collectionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDeleteCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
    };
  },

  toJSON(message: MsgDeleteCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteCollection>, I>>(object: I): MsgDeleteCollection {
    const message = createBaseMsgDeleteCollection();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseMsgDeleteCollectionResponse(): MsgDeleteCollectionResponse {
  return {};
}

export const MsgDeleteCollectionResponse = {
  encode(_: MsgDeleteCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteCollectionResponse();
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

  fromJSON(_: any): MsgDeleteCollectionResponse {
    return {};
  },

  toJSON(_: MsgDeleteCollectionResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteCollectionResponse>, I>>(_: I): MsgDeleteCollectionResponse {
    const message = createBaseMsgDeleteCollectionResponse();
    return message;
  },
};

function createBaseMsgUpdateUserApprovedTransfers(): MsgUpdateUserApprovedTransfers {
  return {
    creator: "",
    collectionId: "",
    updateApprovedOutgoingTransfersTimeline: false,
    approvedOutgoingTransfersTimeline: [],
    updateApprovedIncomingTransfersTimeline: false,
    approvedIncomingTransfersTimeline: [],
    updateUserPermissions: false,
    userPermissions: undefined,
  };
}

export const MsgUpdateUserApprovedTransfers = {
  encode(message: MsgUpdateUserApprovedTransfers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.updateApprovedOutgoingTransfersTimeline === true) {
      writer.uint32(24).bool(message.updateApprovedOutgoingTransfersTimeline);
    }
    for (const v of message.approvedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.updateApprovedIncomingTransfersTimeline === true) {
      writer.uint32(40).bool(message.updateApprovedIncomingTransfersTimeline);
    }
    for (const v of message.approvedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.updateUserPermissions === true) {
      writer.uint32(56).bool(message.updateUserPermissions);
    }
    if (message.userPermissions !== undefined) {
      UserPermissions.encode(message.userPermissions, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserApprovedTransfers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserApprovedTransfers();
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
          message.updateApprovedOutgoingTransfersTimeline = reader.bool();
          break;
        case 4:
          message.approvedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 5:
          message.updateApprovedIncomingTransfersTimeline = reader.bool();
          break;
        case 6:
          message.approvedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 7:
          message.updateUserPermissions = reader.bool();
          break;
        case 8:
          message.userPermissions = UserPermissions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateUserApprovedTransfers {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      updateApprovedOutgoingTransfersTimeline: isSet(object.updateApprovedOutgoingTransfersTimeline)
        ? Boolean(object.updateApprovedOutgoingTransfersTimeline)
        : false,
      approvedOutgoingTransfersTimeline: Array.isArray(object?.approvedOutgoingTransfersTimeline)
        ? object.approvedOutgoingTransfersTimeline.map((e: any) => UserApprovedOutgoingTransferTimeline.fromJSON(e))
        : [],
      updateApprovedIncomingTransfersTimeline: isSet(object.updateApprovedIncomingTransfersTimeline)
        ? Boolean(object.updateApprovedIncomingTransfersTimeline)
        : false,
      approvedIncomingTransfersTimeline: Array.isArray(object?.approvedIncomingTransfersTimeline)
        ? object.approvedIncomingTransfersTimeline.map((e: any) => UserApprovedIncomingTransferTimeline.fromJSON(e))
        : [],
      updateUserPermissions: isSet(object.updateUserPermissions) ? Boolean(object.updateUserPermissions) : false,
      userPermissions: isSet(object.userPermissions) ? UserPermissions.fromJSON(object.userPermissions) : undefined,
    };
  },

  toJSON(message: MsgUpdateUserApprovedTransfers): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.updateApprovedOutgoingTransfersTimeline !== undefined
      && (obj.updateApprovedOutgoingTransfersTimeline = message.updateApprovedOutgoingTransfersTimeline);
    if (message.approvedOutgoingTransfersTimeline) {
      obj.approvedOutgoingTransfersTimeline = message.approvedOutgoingTransfersTimeline.map((e) =>
        e ? UserApprovedOutgoingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedOutgoingTransfersTimeline = [];
    }
    message.updateApprovedIncomingTransfersTimeline !== undefined
      && (obj.updateApprovedIncomingTransfersTimeline = message.updateApprovedIncomingTransfersTimeline);
    if (message.approvedIncomingTransfersTimeline) {
      obj.approvedIncomingTransfersTimeline = message.approvedIncomingTransfersTimeline.map((e) =>
        e ? UserApprovedIncomingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedIncomingTransfersTimeline = [];
    }
    message.updateUserPermissions !== undefined && (obj.updateUserPermissions = message.updateUserPermissions);
    message.userPermissions !== undefined
      && (obj.userPermissions = message.userPermissions ? UserPermissions.toJSON(message.userPermissions) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserApprovedTransfers>, I>>(
    object: I,
  ): MsgUpdateUserApprovedTransfers {
    const message = createBaseMsgUpdateUserApprovedTransfers();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.updateApprovedOutgoingTransfersTimeline = object.updateApprovedOutgoingTransfersTimeline ?? false;
    message.approvedOutgoingTransfersTimeline =
      object.approvedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e)) || [];
    message.updateApprovedIncomingTransfersTimeline = object.updateApprovedIncomingTransfersTimeline ?? false;
    message.approvedIncomingTransfersTimeline =
      object.approvedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e)) || [];
    message.updateUserPermissions = object.updateUserPermissions ?? false;
    message.userPermissions = (object.userPermissions !== undefined && object.userPermissions !== null)
      ? UserPermissions.fromPartial(object.userPermissions)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateUserApprovedTransfersResponse(): MsgUpdateUserApprovedTransfersResponse {
  return {};
}

export const MsgUpdateUserApprovedTransfersResponse = {
  encode(_: MsgUpdateUserApprovedTransfersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserApprovedTransfersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserApprovedTransfersResponse();
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

  fromJSON(_: any): MsgUpdateUserApprovedTransfersResponse {
    return {};
  },

  toJSON(_: MsgUpdateUserApprovedTransfersResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserApprovedTransfersResponse>, I>>(
    _: I,
  ): MsgUpdateUserApprovedTransfersResponse {
    const message = createBaseMsgUpdateUserApprovedTransfersResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  UpdateCollection(request: MsgUpdateCollection): Promise<MsgUpdateCollectionResponse>;
  CreateAddressMappings(request: MsgCreateAddressMappings): Promise<MsgCreateAddressMappingsResponse>;
  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse>;
  UpdateUserApprovedTransfers(request: MsgUpdateUserApprovedTransfers): Promise<MsgUpdateUserApprovedTransfersResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.UpdateCollection = this.UpdateCollection.bind(this);
    this.CreateAddressMappings = this.CreateAddressMappings.bind(this);
    this.TransferBadges = this.TransferBadges.bind(this);
    this.UpdateUserApprovedTransfers = this.UpdateUserApprovedTransfers.bind(this);
    this.DeleteCollection = this.DeleteCollection.bind(this);
  }
  UpdateCollection(request: MsgUpdateCollection): Promise<MsgUpdateCollectionResponse> {
    const data = MsgUpdateCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateCollection", data);
    return promise.then((data) => MsgUpdateCollectionResponse.decode(new _m0.Reader(data)));
  }

  CreateAddressMappings(request: MsgCreateAddressMappings): Promise<MsgCreateAddressMappingsResponse> {
    const data = MsgCreateAddressMappings.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "CreateAddressMappings", data);
    return promise.then((data) => MsgCreateAddressMappingsResponse.decode(new _m0.Reader(data)));
  }

  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse> {
    const data = MsgTransferBadges.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "TransferBadges", data);
    return promise.then((data) => MsgTransferBadgesResponse.decode(new _m0.Reader(data)));
  }

  UpdateUserApprovedTransfers(
    request: MsgUpdateUserApprovedTransfers,
  ): Promise<MsgUpdateUserApprovedTransfersResponse> {
    const data = MsgUpdateUserApprovedTransfers.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateUserApprovedTransfers", data);
    return promise.then((data) => MsgUpdateUserApprovedTransfersResponse.decode(new _m0.Reader(data)));
  }

  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse> {
    const data = MsgDeleteCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "DeleteCollection", data);
    return promise.then((data) => MsgDeleteCollectionResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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
