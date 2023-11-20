/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { AddressMapping } from "./address_mappings";
import { Balance } from "./balances";
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
import { CollectionApproval, Transfer, UserIncomingApproval, UserOutgoingApproval } from "./transfers";

export const protobufPackage = "badges";

/** Used for WASM bindings and JSON parsing */
export interface BadgeCustomMsgType {
  createAddressMappingsMsg: MsgCreateAddressMappings | undefined;
  universalUpdateCollectionMsg: MsgUniversalUpdateCollection | undefined;
  deleteCollectionMsg: MsgDeleteCollection | undefined;
  transferBadgesMsg: MsgTransferBadges | undefined;
  updateUserApprovalsMsg: MsgUpdateUserApprovals | undefined;
  updateCollectionMsg: MsgUpdateCollection | undefined;
  createCollectionMsg: MsgCreateCollection | undefined;
}

/** The types defined in these files are used to define the MsgServer types for all requests and responses for Msgs of the badges module. */
export interface MsgUniversalUpdateCollection {
  creator: string;
  /** 0 for new collection */
  collectionId: string;
  /** The following section of fields are only allowed to be set upon creation of a new collection. */
  balancesType: string;
  /** Default balance options for newly initiated accounts. */
  defaultOutgoingApprovals: UserOutgoingApproval[];
  /** The user's approved incoming transfers for each badge ID. */
  defaultIncomingApprovals: UserIncomingApproval[];
  defaultAutoApproveSelfInitiatedOutgoingTransfers: boolean;
  defaultAutoApproveSelfInitiatedIncomingTransfers: boolean;
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
  updateCollectionApprovals: boolean;
  collectionApprovals: CollectionApproval[];
  updateStandardsTimeline: boolean;
  standardsTimeline: StandardsTimeline[];
  updateIsArchivedTimeline: boolean;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgUniversalUpdateCollectionResponse {
  /** ID of badge collection */
  collectionId: string;
}

export interface MsgUpdateCollection {
  creator: string;
  /** 0 for new collection */
  collectionId: string;
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
  updateCollectionApprovals: boolean;
  collectionApprovals: CollectionApproval[];
  updateStandardsTimeline: boolean;
  standardsTimeline: StandardsTimeline[];
  updateIsArchivedTimeline: boolean;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgUpdateCollectionResponse {
  /** ID of badge collection */
  collectionId: string;
}

export interface MsgCreateCollection {
  creator: string;
  /** The following section of fields are only allowed to be set upon creation of a new collection. */
  balancesType: string;
  /** Default balance options for newly initiated accounts. */
  defaultOutgoingApprovals: UserOutgoingApproval[];
  /** The user's approved incoming transfers for each badge ID. */
  defaultIncomingApprovals: UserIncomingApproval[];
  defaultAutoApproveSelfInitiatedOutgoingTransfers: boolean;
  defaultAutoApproveSelfInitiatedIncomingTransfers: boolean;
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
  updateCollectionApprovals: boolean;
  collectionApprovals: CollectionApproval[];
  updateStandardsTimeline: boolean;
  standardsTimeline: StandardsTimeline[];
  updateIsArchivedTimeline: boolean;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgCreateCollectionResponse {
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

export interface MsgUpdateUserApprovals {
  creator: string;
  collectionId: string;
  updateOutgoingApprovals: boolean;
  outgoingApprovals: UserOutgoingApproval[];
  updateIncomingApprovals: boolean;
  incomingApprovals: UserIncomingApproval[];
  updateAutoApproveSelfInitiatedOutgoingTransfers: boolean;
  autoApproveSelfInitiatedOutgoingTransfers: boolean;
  updateAutoApproveSelfInitiatedIncomingTransfers: boolean;
  autoApproveSelfInitiatedIncomingTransfers: boolean;
  updateUserPermissions: boolean;
  userPermissions: UserPermissions | undefined;
}

export interface MsgUpdateUserApprovalsResponse {
}

function createBaseBadgeCustomMsgType(): BadgeCustomMsgType {
  return {
    createAddressMappingsMsg: undefined,
    universalUpdateCollectionMsg: undefined,
    deleteCollectionMsg: undefined,
    transferBadgesMsg: undefined,
    updateUserApprovalsMsg: undefined,
    updateCollectionMsg: undefined,
    createCollectionMsg: undefined,
  };
}

export const BadgeCustomMsgType = {
  encode(message: BadgeCustomMsgType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.createAddressMappingsMsg !== undefined) {
      MsgCreateAddressMappings.encode(message.createAddressMappingsMsg, writer.uint32(10).fork()).ldelim();
    }
    if (message.universalUpdateCollectionMsg !== undefined) {
      MsgUniversalUpdateCollection.encode(message.universalUpdateCollectionMsg, writer.uint32(18).fork()).ldelim();
    }
    if (message.deleteCollectionMsg !== undefined) {
      MsgDeleteCollection.encode(message.deleteCollectionMsg, writer.uint32(26).fork()).ldelim();
    }
    if (message.transferBadgesMsg !== undefined) {
      MsgTransferBadges.encode(message.transferBadgesMsg, writer.uint32(34).fork()).ldelim();
    }
    if (message.updateUserApprovalsMsg !== undefined) {
      MsgUpdateUserApprovals.encode(message.updateUserApprovalsMsg, writer.uint32(42).fork()).ldelim();
    }
    if (message.updateCollectionMsg !== undefined) {
      MsgUpdateCollection.encode(message.updateCollectionMsg, writer.uint32(50).fork()).ldelim();
    }
    if (message.createCollectionMsg !== undefined) {
      MsgCreateCollection.encode(message.createCollectionMsg, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BadgeCustomMsgType {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBadgeCustomMsgType();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.createAddressMappingsMsg = MsgCreateAddressMappings.decode(reader, reader.uint32());
          break;
        case 2:
          message.universalUpdateCollectionMsg = MsgUniversalUpdateCollection.decode(reader, reader.uint32());
          break;
        case 3:
          message.deleteCollectionMsg = MsgDeleteCollection.decode(reader, reader.uint32());
          break;
        case 4:
          message.transferBadgesMsg = MsgTransferBadges.decode(reader, reader.uint32());
          break;
        case 5:
          message.updateUserApprovalsMsg = MsgUpdateUserApprovals.decode(reader, reader.uint32());
          break;
        case 6:
          message.updateCollectionMsg = MsgUpdateCollection.decode(reader, reader.uint32());
          break;
        case 7:
          message.createCollectionMsg = MsgCreateCollection.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeCustomMsgType {
    return {
      createAddressMappingsMsg: isSet(object.createAddressMappingsMsg)
        ? MsgCreateAddressMappings.fromJSON(object.createAddressMappingsMsg)
        : undefined,
      universalUpdateCollectionMsg: isSet(object.universalUpdateCollectionMsg)
        ? MsgUniversalUpdateCollection.fromJSON(object.universalUpdateCollectionMsg)
        : undefined,
      deleteCollectionMsg: isSet(object.deleteCollectionMsg)
        ? MsgDeleteCollection.fromJSON(object.deleteCollectionMsg)
        : undefined,
      transferBadgesMsg: isSet(object.transferBadgesMsg)
        ? MsgTransferBadges.fromJSON(object.transferBadgesMsg)
        : undefined,
      updateUserApprovalsMsg: isSet(object.updateUserApprovalsMsg)
        ? MsgUpdateUserApprovals.fromJSON(object.updateUserApprovalsMsg)
        : undefined,
      updateCollectionMsg: isSet(object.updateCollectionMsg)
        ? MsgUpdateCollection.fromJSON(object.updateCollectionMsg)
        : undefined,
      createCollectionMsg: isSet(object.createCollectionMsg)
        ? MsgCreateCollection.fromJSON(object.createCollectionMsg)
        : undefined,
    };
  },

  toJSON(message: BadgeCustomMsgType): unknown {
    const obj: any = {};
    message.createAddressMappingsMsg !== undefined && (obj.createAddressMappingsMsg = message.createAddressMappingsMsg
      ? MsgCreateAddressMappings.toJSON(message.createAddressMappingsMsg)
      : undefined);
    message.universalUpdateCollectionMsg !== undefined
      && (obj.universalUpdateCollectionMsg = message.universalUpdateCollectionMsg
        ? MsgUniversalUpdateCollection.toJSON(message.universalUpdateCollectionMsg)
        : undefined);
    message.deleteCollectionMsg !== undefined && (obj.deleteCollectionMsg = message.deleteCollectionMsg
      ? MsgDeleteCollection.toJSON(message.deleteCollectionMsg)
      : undefined);
    message.transferBadgesMsg !== undefined && (obj.transferBadgesMsg = message.transferBadgesMsg
      ? MsgTransferBadges.toJSON(message.transferBadgesMsg)
      : undefined);
    message.updateUserApprovalsMsg !== undefined && (obj.updateUserApprovalsMsg = message.updateUserApprovalsMsg
      ? MsgUpdateUserApprovals.toJSON(message.updateUserApprovalsMsg)
      : undefined);
    message.updateCollectionMsg !== undefined && (obj.updateCollectionMsg = message.updateCollectionMsg
      ? MsgUpdateCollection.toJSON(message.updateCollectionMsg)
      : undefined);
    message.createCollectionMsg !== undefined && (obj.createCollectionMsg = message.createCollectionMsg
      ? MsgCreateCollection.toJSON(message.createCollectionMsg)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BadgeCustomMsgType>, I>>(object: I): BadgeCustomMsgType {
    const message = createBaseBadgeCustomMsgType();
    message.createAddressMappingsMsg =
      (object.createAddressMappingsMsg !== undefined && object.createAddressMappingsMsg !== null)
        ? MsgCreateAddressMappings.fromPartial(object.createAddressMappingsMsg)
        : undefined;
    message.universalUpdateCollectionMsg =
      (object.universalUpdateCollectionMsg !== undefined && object.universalUpdateCollectionMsg !== null)
        ? MsgUniversalUpdateCollection.fromPartial(object.universalUpdateCollectionMsg)
        : undefined;
    message.deleteCollectionMsg = (object.deleteCollectionMsg !== undefined && object.deleteCollectionMsg !== null)
      ? MsgDeleteCollection.fromPartial(object.deleteCollectionMsg)
      : undefined;
    message.transferBadgesMsg = (object.transferBadgesMsg !== undefined && object.transferBadgesMsg !== null)
      ? MsgTransferBadges.fromPartial(object.transferBadgesMsg)
      : undefined;
    message.updateUserApprovalsMsg =
      (object.updateUserApprovalsMsg !== undefined && object.updateUserApprovalsMsg !== null)
        ? MsgUpdateUserApprovals.fromPartial(object.updateUserApprovalsMsg)
        : undefined;
    message.updateCollectionMsg = (object.updateCollectionMsg !== undefined && object.updateCollectionMsg !== null)
      ? MsgUpdateCollection.fromPartial(object.updateCollectionMsg)
      : undefined;
    message.createCollectionMsg = (object.createCollectionMsg !== undefined && object.createCollectionMsg !== null)
      ? MsgCreateCollection.fromPartial(object.createCollectionMsg)
      : undefined;
    return message;
  },
};

function createBaseMsgUniversalUpdateCollection(): MsgUniversalUpdateCollection {
  return {
    creator: "",
    collectionId: "",
    balancesType: "",
    defaultOutgoingApprovals: [],
    defaultIncomingApprovals: [],
    defaultAutoApproveSelfInitiatedOutgoingTransfers: false,
    defaultAutoApproveSelfInitiatedIncomingTransfers: false,
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
    updateCollectionApprovals: false,
    collectionApprovals: [],
    updateStandardsTimeline: false,
    standardsTimeline: [],
    updateIsArchivedTimeline: false,
    isArchivedTimeline: [],
  };
}

export const MsgUniversalUpdateCollection = {
  encode(message: MsgUniversalUpdateCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.balancesType !== "") {
      writer.uint32(26).string(message.balancesType);
    }
    for (const v of message.defaultOutgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.defaultIncomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.defaultAutoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(232).bool(message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.defaultAutoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(240).bool(message.defaultAutoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.defaultUserPermissions !== undefined) {
      UserPermissions.encode(message.defaultUserPermissions, writer.uint32(250).fork()).ldelim();
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
    if (message.updateCollectionApprovals === true) {
      writer.uint32(168).bool(message.updateCollectionApprovals);
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    if (message.updateStandardsTimeline === true) {
      writer.uint32(184).bool(message.updateStandardsTimeline);
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(194).fork()).ldelim();
    }
    if (message.updateIsArchivedTimeline === true) {
      writer.uint32(216).bool(message.updateIsArchivedTimeline);
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(226).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollection();
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
          message.defaultOutgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 5:
          message.defaultIncomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 29:
          message.defaultAutoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 30:
          message.defaultAutoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 31:
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
        case 21:
          message.updateCollectionApprovals = reader.bool();
          break;
        case 22:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 23:
          message.updateStandardsTimeline = reader.bool();
          break;
        case 24:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
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

  fromJSON(object: any): MsgUniversalUpdateCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      balancesType: isSet(object.balancesType) ? String(object.balancesType) : "",
      defaultOutgoingApprovals: Array.isArray(object?.defaultOutgoingApprovals)
        ? object.defaultOutgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      defaultIncomingApprovals: Array.isArray(object?.defaultIncomingApprovals)
        ? object.defaultIncomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      defaultAutoApproveSelfInitiatedOutgoingTransfers: isSet(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        : false,
      defaultAutoApproveSelfInitiatedIncomingTransfers: isSet(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        : false,
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
      updateCollectionApprovals: isSet(object.updateCollectionApprovals)
        ? Boolean(object.updateCollectionApprovals)
        : false,
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      updateStandardsTimeline: isSet(object.updateStandardsTimeline) ? Boolean(object.updateStandardsTimeline) : false,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      updateIsArchivedTimeline: isSet(object.updateIsArchivedTimeline)
        ? Boolean(object.updateIsArchivedTimeline)
        : false,
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUniversalUpdateCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.balancesType !== undefined && (obj.balancesType = message.balancesType);
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
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedOutgoingTransfers =
        message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    message.defaultAutoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedIncomingTransfers =
        message.defaultAutoApproveSelfInitiatedIncomingTransfers);
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
    message.updateCollectionApprovals !== undefined
      && (obj.updateCollectionApprovals = message.updateCollectionApprovals);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    message.updateStandardsTimeline !== undefined && (obj.updateStandardsTimeline = message.updateStandardsTimeline);
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    message.updateIsArchivedTimeline !== undefined && (obj.updateIsArchivedTimeline = message.updateIsArchivedTimeline);
    if (message.isArchivedTimeline) {
      obj.isArchivedTimeline = message.isArchivedTimeline.map((e) => e ? IsArchivedTimeline.toJSON(e) : undefined);
    } else {
      obj.isArchivedTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollection>, I>>(object: I): MsgUniversalUpdateCollection {
    const message = createBaseMsgUniversalUpdateCollection();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.balancesType = object.balancesType ?? "";
    message.defaultOutgoingApprovals = object.defaultOutgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e))
      || [];
    message.defaultIncomingApprovals = object.defaultIncomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e))
      || [];
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers = object.defaultAutoApproveSelfInitiatedOutgoingTransfers
      ?? false;
    message.defaultAutoApproveSelfInitiatedIncomingTransfers = object.defaultAutoApproveSelfInitiatedIncomingTransfers
      ?? false;
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
    message.updateCollectionApprovals = object.updateCollectionApprovals ?? false;
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.updateStandardsTimeline = object.updateStandardsTimeline ?? false;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.updateIsArchivedTimeline = object.updateIsArchivedTimeline ?? false;
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUniversalUpdateCollectionResponse(): MsgUniversalUpdateCollectionResponse {
  return { collectionId: "" };
}

export const MsgUniversalUpdateCollectionResponse = {
  encode(message: MsgUniversalUpdateCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUniversalUpdateCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUniversalUpdateCollectionResponse();
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

  fromJSON(object: any): MsgUniversalUpdateCollectionResponse {
    return { collectionId: isSet(object.collectionId) ? String(object.collectionId) : "" };
  },

  toJSON(message: MsgUniversalUpdateCollectionResponse): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUniversalUpdateCollectionResponse>, I>>(
    object: I,
  ): MsgUniversalUpdateCollectionResponse {
    const message = createBaseMsgUniversalUpdateCollectionResponse();
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseMsgUpdateCollection(): MsgUpdateCollection {
  return {
    creator: "",
    collectionId: "",
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
    updateCollectionApprovals: false,
    collectionApprovals: [],
    updateStandardsTimeline: false,
    standardsTimeline: [],
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
    if (message.updateCollectionApprovals === true) {
      writer.uint32(168).bool(message.updateCollectionApprovals);
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    if (message.updateStandardsTimeline === true) {
      writer.uint32(184).bool(message.updateStandardsTimeline);
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(194).fork()).ldelim();
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
        case 21:
          message.updateCollectionApprovals = reader.bool();
          break;
        case 22:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 23:
          message.updateStandardsTimeline = reader.bool();
          break;
        case 24:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
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
      updateCollectionApprovals: isSet(object.updateCollectionApprovals)
        ? Boolean(object.updateCollectionApprovals)
        : false,
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      updateStandardsTimeline: isSet(object.updateStandardsTimeline) ? Boolean(object.updateStandardsTimeline) : false,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
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
    message.updateCollectionApprovals !== undefined
      && (obj.updateCollectionApprovals = message.updateCollectionApprovals);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    message.updateStandardsTimeline !== undefined && (obj.updateStandardsTimeline = message.updateStandardsTimeline);
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
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
    message.updateCollectionApprovals = object.updateCollectionApprovals ?? false;
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.updateStandardsTimeline = object.updateStandardsTimeline ?? false;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
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

function createBaseMsgCreateCollection(): MsgCreateCollection {
  return {
    creator: "",
    balancesType: "",
    defaultOutgoingApprovals: [],
    defaultIncomingApprovals: [],
    defaultAutoApproveSelfInitiatedOutgoingTransfers: false,
    defaultAutoApproveSelfInitiatedIncomingTransfers: false,
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
    updateCollectionApprovals: false,
    collectionApprovals: [],
    updateStandardsTimeline: false,
    standardsTimeline: [],
    updateIsArchivedTimeline: false,
    isArchivedTimeline: [],
  };
}

export const MsgCreateCollection = {
  encode(message: MsgCreateCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.balancesType !== "") {
      writer.uint32(26).string(message.balancesType);
    }
    for (const v of message.defaultOutgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.defaultIncomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.defaultAutoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(232).bool(message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.defaultAutoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(240).bool(message.defaultAutoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.defaultUserPermissions !== undefined) {
      UserPermissions.encode(message.defaultUserPermissions, writer.uint32(250).fork()).ldelim();
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
    if (message.updateCollectionApprovals === true) {
      writer.uint32(168).bool(message.updateCollectionApprovals);
    }
    for (const v of message.collectionApprovals) {
      CollectionApproval.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    if (message.updateStandardsTimeline === true) {
      writer.uint32(184).bool(message.updateStandardsTimeline);
    }
    for (const v of message.standardsTimeline) {
      StandardsTimeline.encode(v!, writer.uint32(194).fork()).ldelim();
    }
    if (message.updateIsArchivedTimeline === true) {
      writer.uint32(216).bool(message.updateIsArchivedTimeline);
    }
    for (const v of message.isArchivedTimeline) {
      IsArchivedTimeline.encode(v!, writer.uint32(226).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 3:
          message.balancesType = reader.string();
          break;
        case 4:
          message.defaultOutgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 5:
          message.defaultIncomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 29:
          message.defaultAutoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 30:
          message.defaultAutoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 31:
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
        case 21:
          message.updateCollectionApprovals = reader.bool();
          break;
        case 22:
          message.collectionApprovals.push(CollectionApproval.decode(reader, reader.uint32()));
          break;
        case 23:
          message.updateStandardsTimeline = reader.bool();
          break;
        case 24:
          message.standardsTimeline.push(StandardsTimeline.decode(reader, reader.uint32()));
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

  fromJSON(object: any): MsgCreateCollection {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      balancesType: isSet(object.balancesType) ? String(object.balancesType) : "",
      defaultOutgoingApprovals: Array.isArray(object?.defaultOutgoingApprovals)
        ? object.defaultOutgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      defaultIncomingApprovals: Array.isArray(object?.defaultIncomingApprovals)
        ? object.defaultIncomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      defaultAutoApproveSelfInitiatedOutgoingTransfers: isSet(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedOutgoingTransfers)
        : false,
      defaultAutoApproveSelfInitiatedIncomingTransfers: isSet(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.defaultAutoApproveSelfInitiatedIncomingTransfers)
        : false,
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
      updateCollectionApprovals: isSet(object.updateCollectionApprovals)
        ? Boolean(object.updateCollectionApprovals)
        : false,
      collectionApprovals: Array.isArray(object?.collectionApprovals)
        ? object.collectionApprovals.map((e: any) => CollectionApproval.fromJSON(e))
        : [],
      updateStandardsTimeline: isSet(object.updateStandardsTimeline) ? Boolean(object.updateStandardsTimeline) : false,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      updateIsArchivedTimeline: isSet(object.updateIsArchivedTimeline)
        ? Boolean(object.updateIsArchivedTimeline)
        : false,
      isArchivedTimeline: Array.isArray(object?.isArchivedTimeline)
        ? object.isArchivedTimeline.map((e: any) => IsArchivedTimeline.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgCreateCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.balancesType !== undefined && (obj.balancesType = message.balancesType);
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
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedOutgoingTransfers =
        message.defaultAutoApproveSelfInitiatedOutgoingTransfers);
    message.defaultAutoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.defaultAutoApproveSelfInitiatedIncomingTransfers =
        message.defaultAutoApproveSelfInitiatedIncomingTransfers);
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
    message.updateCollectionApprovals !== undefined
      && (obj.updateCollectionApprovals = message.updateCollectionApprovals);
    if (message.collectionApprovals) {
      obj.collectionApprovals = message.collectionApprovals.map((e) => e ? CollectionApproval.toJSON(e) : undefined);
    } else {
      obj.collectionApprovals = [];
    }
    message.updateStandardsTimeline !== undefined && (obj.updateStandardsTimeline = message.updateStandardsTimeline);
    if (message.standardsTimeline) {
      obj.standardsTimeline = message.standardsTimeline.map((e) => e ? StandardsTimeline.toJSON(e) : undefined);
    } else {
      obj.standardsTimeline = [];
    }
    message.updateIsArchivedTimeline !== undefined && (obj.updateIsArchivedTimeline = message.updateIsArchivedTimeline);
    if (message.isArchivedTimeline) {
      obj.isArchivedTimeline = message.isArchivedTimeline.map((e) => e ? IsArchivedTimeline.toJSON(e) : undefined);
    } else {
      obj.isArchivedTimeline = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateCollection>, I>>(object: I): MsgCreateCollection {
    const message = createBaseMsgCreateCollection();
    message.creator = object.creator ?? "";
    message.balancesType = object.balancesType ?? "";
    message.defaultOutgoingApprovals = object.defaultOutgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e))
      || [];
    message.defaultIncomingApprovals = object.defaultIncomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e))
      || [];
    message.defaultAutoApproveSelfInitiatedOutgoingTransfers = object.defaultAutoApproveSelfInitiatedOutgoingTransfers
      ?? false;
    message.defaultAutoApproveSelfInitiatedIncomingTransfers = object.defaultAutoApproveSelfInitiatedIncomingTransfers
      ?? false;
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
    message.updateCollectionApprovals = object.updateCollectionApprovals ?? false;
    message.collectionApprovals = object.collectionApprovals?.map((e) => CollectionApproval.fromPartial(e)) || [];
    message.updateStandardsTimeline = object.updateStandardsTimeline ?? false;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.updateIsArchivedTimeline = object.updateIsArchivedTimeline ?? false;
    message.isArchivedTimeline = object.isArchivedTimeline?.map((e) => IsArchivedTimeline.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgCreateCollectionResponse(): MsgCreateCollectionResponse {
  return { collectionId: "" };
}

export const MsgCreateCollectionResponse = {
  encode(message: MsgCreateCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateCollectionResponse();
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

  fromJSON(object: any): MsgCreateCollectionResponse {
    return { collectionId: isSet(object.collectionId) ? String(object.collectionId) : "" };
  },

  toJSON(message: MsgCreateCollectionResponse): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateCollectionResponse>, I>>(object: I): MsgCreateCollectionResponse {
    const message = createBaseMsgCreateCollectionResponse();
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

function createBaseMsgUpdateUserApprovals(): MsgUpdateUserApprovals {
  return {
    creator: "",
    collectionId: "",
    updateOutgoingApprovals: false,
    outgoingApprovals: [],
    updateIncomingApprovals: false,
    incomingApprovals: [],
    updateAutoApproveSelfInitiatedOutgoingTransfers: false,
    autoApproveSelfInitiatedOutgoingTransfers: false,
    updateAutoApproveSelfInitiatedIncomingTransfers: false,
    autoApproveSelfInitiatedIncomingTransfers: false,
    updateUserPermissions: false,
    userPermissions: undefined,
  };
}

export const MsgUpdateUserApprovals = {
  encode(message: MsgUpdateUserApprovals, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    if (message.updateOutgoingApprovals === true) {
      writer.uint32(24).bool(message.updateOutgoingApprovals);
    }
    for (const v of message.outgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.updateIncomingApprovals === true) {
      writer.uint32(40).bool(message.updateIncomingApprovals);
    }
    for (const v of message.incomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.updateAutoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(56).bool(message.updateAutoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.autoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(64).bool(message.autoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.updateAutoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(72).bool(message.updateAutoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.autoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(80).bool(message.autoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.updateUserPermissions === true) {
      writer.uint32(88).bool(message.updateUserPermissions);
    }
    if (message.userPermissions !== undefined) {
      UserPermissions.encode(message.userPermissions, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserApprovals {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserApprovals();
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
          message.updateOutgoingApprovals = reader.bool();
          break;
        case 4:
          message.outgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 5:
          message.updateIncomingApprovals = reader.bool();
          break;
        case 6:
          message.incomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 7:
          message.updateAutoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 8:
          message.autoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 9:
          message.updateAutoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 10:
          message.autoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 11:
          message.updateUserPermissions = reader.bool();
          break;
        case 12:
          message.userPermissions = UserPermissions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateUserApprovals {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      updateOutgoingApprovals: isSet(object.updateOutgoingApprovals) ? Boolean(object.updateOutgoingApprovals) : false,
      outgoingApprovals: Array.isArray(object?.outgoingApprovals)
        ? object.outgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      updateIncomingApprovals: isSet(object.updateIncomingApprovals) ? Boolean(object.updateIncomingApprovals) : false,
      incomingApprovals: Array.isArray(object?.incomingApprovals)
        ? object.incomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      updateAutoApproveSelfInitiatedOutgoingTransfers: isSet(object.updateAutoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.updateAutoApproveSelfInitiatedOutgoingTransfers)
        : false,
      autoApproveSelfInitiatedOutgoingTransfers: isSet(object.autoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.autoApproveSelfInitiatedOutgoingTransfers)
        : false,
      updateAutoApproveSelfInitiatedIncomingTransfers: isSet(object.updateAutoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.updateAutoApproveSelfInitiatedIncomingTransfers)
        : false,
      autoApproveSelfInitiatedIncomingTransfers: isSet(object.autoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.autoApproveSelfInitiatedIncomingTransfers)
        : false,
      updateUserPermissions: isSet(object.updateUserPermissions) ? Boolean(object.updateUserPermissions) : false,
      userPermissions: isSet(object.userPermissions) ? UserPermissions.fromJSON(object.userPermissions) : undefined,
    };
  },

  toJSON(message: MsgUpdateUserApprovals): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.updateOutgoingApprovals !== undefined && (obj.updateOutgoingApprovals = message.updateOutgoingApprovals);
    if (message.outgoingApprovals) {
      obj.outgoingApprovals = message.outgoingApprovals.map((e) => e ? UserOutgoingApproval.toJSON(e) : undefined);
    } else {
      obj.outgoingApprovals = [];
    }
    message.updateIncomingApprovals !== undefined && (obj.updateIncomingApprovals = message.updateIncomingApprovals);
    if (message.incomingApprovals) {
      obj.incomingApprovals = message.incomingApprovals.map((e) => e ? UserIncomingApproval.toJSON(e) : undefined);
    } else {
      obj.incomingApprovals = [];
    }
    message.updateAutoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.updateAutoApproveSelfInitiatedOutgoingTransfers =
        message.updateAutoApproveSelfInitiatedOutgoingTransfers);
    message.autoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.autoApproveSelfInitiatedOutgoingTransfers = message.autoApproveSelfInitiatedOutgoingTransfers);
    message.updateAutoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.updateAutoApproveSelfInitiatedIncomingTransfers =
        message.updateAutoApproveSelfInitiatedIncomingTransfers);
    message.autoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.autoApproveSelfInitiatedIncomingTransfers = message.autoApproveSelfInitiatedIncomingTransfers);
    message.updateUserPermissions !== undefined && (obj.updateUserPermissions = message.updateUserPermissions);
    message.userPermissions !== undefined
      && (obj.userPermissions = message.userPermissions ? UserPermissions.toJSON(message.userPermissions) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserApprovals>, I>>(object: I): MsgUpdateUserApprovals {
    const message = createBaseMsgUpdateUserApprovals();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.updateOutgoingApprovals = object.updateOutgoingApprovals ?? false;
    message.outgoingApprovals = object.outgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e)) || [];
    message.updateIncomingApprovals = object.updateIncomingApprovals ?? false;
    message.incomingApprovals = object.incomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e)) || [];
    message.updateAutoApproveSelfInitiatedOutgoingTransfers = object.updateAutoApproveSelfInitiatedOutgoingTransfers
      ?? false;
    message.autoApproveSelfInitiatedOutgoingTransfers = object.autoApproveSelfInitiatedOutgoingTransfers ?? false;
    message.updateAutoApproveSelfInitiatedIncomingTransfers = object.updateAutoApproveSelfInitiatedIncomingTransfers
      ?? false;
    message.autoApproveSelfInitiatedIncomingTransfers = object.autoApproveSelfInitiatedIncomingTransfers ?? false;
    message.updateUserPermissions = object.updateUserPermissions ?? false;
    message.userPermissions = (object.userPermissions !== undefined && object.userPermissions !== null)
      ? UserPermissions.fromPartial(object.userPermissions)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateUserApprovalsResponse(): MsgUpdateUserApprovalsResponse {
  return {};
}

export const MsgUpdateUserApprovalsResponse = {
  encode(_: MsgUpdateUserApprovalsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUserApprovalsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUserApprovalsResponse();
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

  fromJSON(_: any): MsgUpdateUserApprovalsResponse {
    return {};
  },

  toJSON(_: MsgUpdateUserApprovalsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserApprovalsResponse>, I>>(_: I): MsgUpdateUserApprovalsResponse {
    const message = createBaseMsgUpdateUserApprovalsResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  UniversalUpdateCollection(request: MsgUniversalUpdateCollection): Promise<MsgUniversalUpdateCollectionResponse>;
  CreateAddressMappings(request: MsgCreateAddressMappings): Promise<MsgCreateAddressMappingsResponse>;
  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse>;
  UpdateUserApprovals(request: MsgUpdateUserApprovals): Promise<MsgUpdateUserApprovalsResponse>;
  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse>;
  UpdateCollection(request: MsgUpdateCollection): Promise<MsgUpdateCollectionResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  CreateCollection(request: MsgCreateCollection): Promise<MsgCreateCollectionResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.UniversalUpdateCollection = this.UniversalUpdateCollection.bind(this);
    this.CreateAddressMappings = this.CreateAddressMappings.bind(this);
    this.TransferBadges = this.TransferBadges.bind(this);
    this.UpdateUserApprovals = this.UpdateUserApprovals.bind(this);
    this.DeleteCollection = this.DeleteCollection.bind(this);
    this.UpdateCollection = this.UpdateCollection.bind(this);
    this.CreateCollection = this.CreateCollection.bind(this);
  }
  UniversalUpdateCollection(request: MsgUniversalUpdateCollection): Promise<MsgUniversalUpdateCollectionResponse> {
    const data = MsgUniversalUpdateCollection.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "UniversalUpdateCollection", data);
    return promise.then((data) => MsgUniversalUpdateCollectionResponse.decode(new _m0.Reader(data)));
  }

  CreateAddressMappings(request: MsgCreateAddressMappings): Promise<MsgCreateAddressMappingsResponse> {
    const data = MsgCreateAddressMappings.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "CreateAddressMappings", data);
    return promise.then((data) => MsgCreateAddressMappingsResponse.decode(new _m0.Reader(data)));
  }

  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse> {
    const data = MsgTransferBadges.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "TransferBadges", data);
    return promise.then((data) => MsgTransferBadgesResponse.decode(new _m0.Reader(data)));
  }

  UpdateUserApprovals(request: MsgUpdateUserApprovals): Promise<MsgUpdateUserApprovalsResponse> {
    const data = MsgUpdateUserApprovals.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "UpdateUserApprovals", data);
    return promise.then((data) => MsgUpdateUserApprovalsResponse.decode(new _m0.Reader(data)));
  }

  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse> {
    const data = MsgDeleteCollection.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "DeleteCollection", data);
    return promise.then((data) => MsgDeleteCollectionResponse.decode(new _m0.Reader(data)));
  }

  UpdateCollection(request: MsgUpdateCollection): Promise<MsgUpdateCollectionResponse> {
    const data = MsgUpdateCollection.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "UpdateCollection", data);
    return promise.then((data) => MsgUpdateCollectionResponse.decode(new _m0.Reader(data)));
  }

  CreateCollection(request: MsgCreateCollection): Promise<MsgCreateCollectionResponse> {
    const data = MsgCreateCollection.encode(request).finish();
    const promise = this.rpc.request("badges.Msg", "CreateCollection", data);
    return promise.then((data) => MsgCreateCollectionResponse.decode(new _m0.Reader(data)));
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
