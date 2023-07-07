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
import { UserApprovedIncomingTransferTimeline, UserApprovedOutgoingTransferTimeline } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

export interface Transfer {
  from: string;
  toAddresses: string[];
  balances: Balance[];
  /**
   * Note here we remain optimistic that the solutions will apply to all potential challenges.
   * It is the Tx Sender's responsibility to ensure that the solutions are valid for all potential challenges.
   * If you are attempting to claim badges with different sets of challenges, you will need to make multiple transfers.
   */
  solutions: ChallengeSolution[];
}

export interface MsgNewCollection {
  /** See collections.proto for more details about these MsgNewBadge fields. Defines the badge details. Leave unneeded fields empty. */
  creator: string;
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  customDataTimeline: CustomDataTimeline[];
  balancesType: string;
  inheritedBalancesTimeline: InheritedBalancesTimeline[];
  collectionApprovedTransfersTimeline: CollectionApprovedTransferTimeline[];
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
  defaultApprovedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  /** The user's approved incoming transfers for each badge ID. */
  defaultApprovedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
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
  inheritedBalancesTimeline: InheritedBalancesTimeline[];
  collectionMetadataTimeline: CollectionMetadataTimeline[];
  badgeMetadataTimeline: BadgeMetadataTimeline[];
  offChainBalancesMetadataTimeline: OffChainBalancesMetadataTimeline[];
  collectionApprovedTransfersTimeline: CollectionApprovedTransferTimeline[];
  addressMappings: AddressMapping[];
}

export interface MsgMintAndDistributeBadgesResponse {
}

/** For each amount, for each toAddress, we will attempt to transfer all the badgeIds for the badge with ID badgeId. */
export interface MsgTransferBadges {
  creator: string;
  collectionId: string;
  transfers: Transfer[];
}

export interface MsgTransferBadgesResponse {
}

export interface MsgUpdateCollectionApprovedTransfers {
  creator: string;
  collectionId: string;
  collectionApprovedTransfersTimeline: CollectionApprovedTransferTimeline[];
  addressMappings: AddressMapping[];
}

export interface MsgUpdateCollectionApprovedTransfersResponse {
}

export interface MsgUpdateUserApprovedTransfers {
  creator: string;
  collectionId: string;
  approvedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  approvedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
  addressMappings: AddressMapping[];
}

export interface MsgUpdateUserApprovedTransfersResponse {
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
export interface MsgUpdateCollectionPermissions {
  creator: string;
  collectionId: string;
  permissions: CollectionPermissions | undefined;
  addressMappings: AddressMapping[];
}

export interface MsgUpdateCollectionPermissionsResponse {
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

export interface ClaimProofItem {
  aunt: string;
  onRight: boolean;
}

/** Consistent with tendermint/crypto merkle tree */
export interface ClaimProof {
  leaf: string;
  aunts: ClaimProofItem[];
}

export interface ChallengeSolution {
  proof: ClaimProof | undefined;
}

export interface MsgDeleteCollection {
  creator: string;
  collectionId: string;
}

export interface MsgDeleteCollectionResponse {
}

export interface MsgArchiveCollection {
  creator: string;
  collectionId: string;
  isArchivedTimeline: IsArchivedTimeline[];
}

export interface MsgArchiveCollectionResponse {
}

function createBaseTransfer(): Transfer {
  return { from: "", toAddresses: [], balances: [], solutions: [] };
}

export const Transfer = {
  encode(message: Transfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.from !== "") {
      writer.uint32(10).string(message.from);
    }
    for (const v of message.toAddresses) {
      writer.uint32(18).string(v!);
    }
    for (const v of message.balances) {
      Balance.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.solutions) {
      ChallengeSolution.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Transfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.from = reader.string();
          break;
        case 2:
          message.toAddresses.push(reader.string());
          break;
        case 3:
          message.balances.push(Balance.decode(reader, reader.uint32()));
          break;
        case 4:
          message.solutions.push(ChallengeSolution.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Transfer {
    return {
      from: isSet(object.from) ? String(object.from) : "",
      toAddresses: Array.isArray(object?.toAddresses) ? object.toAddresses.map((e: any) => String(e)) : [],
      balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => Balance.fromJSON(e)) : [],
      solutions: Array.isArray(object?.solutions)
        ? object.solutions.map((e: any) => ChallengeSolution.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Transfer): unknown {
    const obj: any = {};
    message.from !== undefined && (obj.from = message.from);
    if (message.toAddresses) {
      obj.toAddresses = message.toAddresses.map((e) => e);
    } else {
      obj.toAddresses = [];
    }
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    if (message.solutions) {
      obj.solutions = message.solutions.map((e) => e ? ChallengeSolution.toJSON(e) : undefined);
    } else {
      obj.solutions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Transfer>, I>>(object: I): Transfer {
    const message = createBaseTransfer();
    message.from = object.from ?? "";
    message.toAddresses = object.toAddresses?.map((e) => e) || [];
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.solutions = object.solutions?.map((e) => ChallengeSolution.fromPartial(e)) || [];
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
    inheritedBalancesTimeline: [],
    collectionApprovedTransfersTimeline: [],
    permissions: undefined,
    standardsTimeline: [],
    badgesToCreate: [],
    transfers: [],
    contractAddressTimeline: [],
    addressMappings: [],
    defaultApprovedOutgoingTransfersTimeline: [],
    defaultApprovedIncomingTransfersTimeline: [],
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
    for (const v of message.inheritedBalancesTimeline) {
      InheritedBalancesTimeline.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.collectionApprovedTransfersTimeline) {
      CollectionApprovedTransferTimeline.encode(v!, writer.uint32(66).fork()).ldelim();
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
    for (const v of message.defaultApprovedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    for (const v of message.defaultApprovedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(130).fork()).ldelim();
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
        case 7:
          message.inheritedBalancesTimeline.push(InheritedBalancesTimeline.decode(reader, reader.uint32()));
          break;
        case 8:
          message.collectionApprovedTransfersTimeline.push(
            CollectionApprovedTransferTimeline.decode(reader, reader.uint32()),
          );
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
          message.defaultApprovedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 16:
          message.defaultApprovedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
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
      inheritedBalancesTimeline: Array.isArray(object?.inheritedBalancesTimeline)
        ? object.inheritedBalancesTimeline.map((e: any) => InheritedBalancesTimeline.fromJSON(e))
        : [],
      collectionApprovedTransfersTimeline: Array.isArray(object?.collectionApprovedTransfersTimeline)
        ? object.collectionApprovedTransfersTimeline.map((e: any) => CollectionApprovedTransferTimeline.fromJSON(e))
        : [],
      permissions: isSet(object.permissions) ? CollectionPermissions.fromJSON(object.permissions) : undefined,
      standardsTimeline: Array.isArray(object?.standardsTimeline)
        ? object.standardsTimeline.map((e: any) => StandardsTimeline.fromJSON(e))
        : [],
      badgesToCreate: Array.isArray(object?.badgesToCreate)
        ? object.badgesToCreate.map((e: any) => Balance.fromJSON(e))
        : [],
      transfers: Array.isArray(object?.transfers)
        ? object.transfers.map((e: any) => Transfer.fromJSON(e))
        : [],
      contractAddressTimeline: Array.isArray(object?.contractAddressTimeline)
        ? object.contractAddressTimeline.map((e: any) => ContractAddressTimeline.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
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
    if (message.inheritedBalancesTimeline) {
      obj.inheritedBalancesTimeline = message.inheritedBalancesTimeline.map((e) =>
        e ? InheritedBalancesTimeline.toJSON(e) : undefined
      );
    } else {
      obj.inheritedBalancesTimeline = [];
    }
    if (message.collectionApprovedTransfersTimeline) {
      obj.collectionApprovedTransfersTimeline = message.collectionApprovedTransfersTimeline.map((e) =>
        e ? CollectionApprovedTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionApprovedTransfersTimeline = [];
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
    message.inheritedBalancesTimeline =
      object.inheritedBalancesTimeline?.map((e) => InheritedBalancesTimeline.fromPartial(e)) || [];
    message.collectionApprovedTransfersTimeline =
      object.collectionApprovedTransfersTimeline?.map((e) => CollectionApprovedTransferTimeline.fromPartial(e)) || [];
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? CollectionPermissions.fromPartial(object.permissions)
      : undefined;
    message.standardsTimeline = object.standardsTimeline?.map((e) => StandardsTimeline.fromPartial(e)) || [];
    message.badgesToCreate = object.badgesToCreate?.map((e) => Balance.fromPartial(e)) || [];
    message.transfers = object.transfers?.map((e) => Transfer.fromPartial(e)) || [];
    message.contractAddressTimeline = object.contractAddressTimeline?.map((e) => ContractAddressTimeline.fromPartial(e))
      || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    message.defaultApprovedOutgoingTransfersTimeline =
      object.defaultApprovedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e))
      || [];
    message.defaultApprovedIncomingTransfersTimeline =
      object.defaultApprovedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e))
      || [];
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
    inheritedBalancesTimeline: [],
    collectionMetadataTimeline: [],
    badgeMetadataTimeline: [],
    offChainBalancesMetadataTimeline: [],
    collectionApprovedTransfersTimeline: [],
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
    for (const v of message.inheritedBalancesTimeline) {
      InheritedBalancesTimeline.encode(v!, writer.uint32(42).fork()).ldelim();
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
    for (const v of message.collectionApprovedTransfersTimeline) {
      CollectionApprovedTransferTimeline.encode(v!, writer.uint32(74).fork()).ldelim();
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
        case 5:
          message.inheritedBalancesTimeline.push(InheritedBalancesTimeline.decode(reader, reader.uint32()));
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
          message.collectionApprovedTransfersTimeline.push(
            CollectionApprovedTransferTimeline.decode(reader, reader.uint32()),
          );
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
      inheritedBalancesTimeline: Array.isArray(object?.inheritedBalancesTimeline)
        ? object.inheritedBalancesTimeline.map((e: any) => InheritedBalancesTimeline.fromJSON(e))
        : [],
      collectionMetadataTimeline: Array.isArray(object?.collectionMetadataTimeline)
        ? object.collectionMetadataTimeline.map((e: any) => CollectionMetadataTimeline.fromJSON(e))
        : [],
      badgeMetadataTimeline: Array.isArray(object?.badgeMetadataTimeline)
        ? object.badgeMetadataTimeline.map((e: any) => BadgeMetadataTimeline.fromJSON(e))
        : [],
      offChainBalancesMetadataTimeline: Array.isArray(object?.offChainBalancesMetadataTimeline)
        ? object.offChainBalancesMetadataTimeline.map((e: any) => OffChainBalancesMetadataTimeline.fromJSON(e))
        : [],
      collectionApprovedTransfersTimeline: Array.isArray(object?.collectionApprovedTransfersTimeline)
        ? object.collectionApprovedTransfersTimeline.map((e: any) => CollectionApprovedTransferTimeline.fromJSON(e))
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
    if (message.inheritedBalancesTimeline) {
      obj.inheritedBalancesTimeline = message.inheritedBalancesTimeline.map((e) =>
        e ? InheritedBalancesTimeline.toJSON(e) : undefined
      );
    } else {
      obj.inheritedBalancesTimeline = [];
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
    if (message.collectionApprovedTransfersTimeline) {
      obj.collectionApprovedTransfersTimeline = message.collectionApprovedTransfersTimeline.map((e) =>
        e ? CollectionApprovedTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionApprovedTransfersTimeline = [];
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
    message.inheritedBalancesTimeline =
      object.inheritedBalancesTimeline?.map((e) => InheritedBalancesTimeline.fromPartial(e)) || [];
    message.collectionMetadataTimeline =
      object.collectionMetadataTimeline?.map((e) => CollectionMetadataTimeline.fromPartial(e)) || [];
    message.badgeMetadataTimeline = object.badgeMetadataTimeline?.map((e) => BadgeMetadataTimeline.fromPartial(e))
      || [];
    message.offChainBalancesMetadataTimeline =
      object.offChainBalancesMetadataTimeline?.map((e) => OffChainBalancesMetadataTimeline.fromPartial(e)) || [];
    message.collectionApprovedTransfersTimeline =
      object.collectionApprovedTransfersTimeline?.map((e) => CollectionApprovedTransferTimeline.fromPartial(e)) || [];
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

function createBaseMsgUpdateCollectionApprovedTransfers(): MsgUpdateCollectionApprovedTransfers {
  return { creator: "", collectionId: "", collectionApprovedTransfersTimeline: [], addressMappings: [] };
}

export const MsgUpdateCollectionApprovedTransfers = {
  encode(message: MsgUpdateCollectionApprovedTransfers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.collectionId !== "") {
      writer.uint32(18).string(message.collectionId);
    }
    for (const v of message.collectionApprovedTransfersTimeline) {
      CollectionApprovedTransferTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollectionApprovedTransfers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollectionApprovedTransfers();
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
          message.collectionApprovedTransfersTimeline.push(
            CollectionApprovedTransferTimeline.decode(reader, reader.uint32()),
          );
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

  fromJSON(object: any): MsgUpdateCollectionApprovedTransfers {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      collectionApprovedTransfersTimeline: Array.isArray(object?.collectionApprovedTransfersTimeline)
        ? object.collectionApprovedTransfersTimeline.map((e: any) => CollectionApprovedTransferTimeline.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateCollectionApprovedTransfers): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.collectionApprovedTransfersTimeline) {
      obj.collectionApprovedTransfersTimeline = message.collectionApprovedTransfersTimeline.map((e) =>
        e ? CollectionApprovedTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.collectionApprovedTransfersTimeline = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollectionApprovedTransfers>, I>>(
    object: I,
  ): MsgUpdateCollectionApprovedTransfers {
    const message = createBaseMsgUpdateCollectionApprovedTransfers();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.collectionApprovedTransfersTimeline =
      object.collectionApprovedTransfersTimeline?.map((e) => CollectionApprovedTransferTimeline.fromPartial(e)) || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateCollectionApprovedTransfersResponse(): MsgUpdateCollectionApprovedTransfersResponse {
  return {};
}

export const MsgUpdateCollectionApprovedTransfersResponse = {
  encode(_: MsgUpdateCollectionApprovedTransfersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollectionApprovedTransfersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollectionApprovedTransfersResponse();
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

  fromJSON(_: any): MsgUpdateCollectionApprovedTransfersResponse {
    return {};
  },

  toJSON(_: MsgUpdateCollectionApprovedTransfersResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollectionApprovedTransfersResponse>, I>>(
    _: I,
  ): MsgUpdateCollectionApprovedTransfersResponse {
    const message = createBaseMsgUpdateCollectionApprovedTransfersResponse();
    return message;
  },
};

function createBaseMsgUpdateUserApprovedTransfers(): MsgUpdateUserApprovedTransfers {
  return {
    creator: "",
    collectionId: "",
    approvedOutgoingTransfersTimeline: [],
    approvedIncomingTransfersTimeline: [],
    addressMappings: [],
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
    for (const v of message.approvedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.approvedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(42).fork()).ldelim();
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
          message.approvedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 4:
          message.approvedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 5:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
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
      approvedOutgoingTransfersTimeline: Array.isArray(object?.approvedOutgoingTransfersTimeline)
        ? object.approvedOutgoingTransfersTimeline.map((e: any) => UserApprovedOutgoingTransferTimeline.fromJSON(e))
        : [],
      approvedIncomingTransfersTimeline: Array.isArray(object?.approvedIncomingTransfersTimeline)
        ? object.approvedIncomingTransfersTimeline.map((e: any) => UserApprovedIncomingTransferTimeline.fromJSON(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateUserApprovedTransfers): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    if (message.approvedOutgoingTransfersTimeline) {
      obj.approvedOutgoingTransfersTimeline = message.approvedOutgoingTransfersTimeline.map((e) =>
        e ? UserApprovedOutgoingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedOutgoingTransfersTimeline = [];
    }
    if (message.approvedIncomingTransfersTimeline) {
      obj.approvedIncomingTransfersTimeline = message.approvedIncomingTransfersTimeline.map((e) =>
        e ? UserApprovedIncomingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedIncomingTransfersTimeline = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUserApprovedTransfers>, I>>(
    object: I,
  ): MsgUpdateUserApprovedTransfers {
    const message = createBaseMsgUpdateUserApprovedTransfers();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.approvedOutgoingTransfersTimeline =
      object.approvedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e)) || [];
    message.approvedIncomingTransfersTimeline =
      object.approvedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e)) || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
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

function createBaseMsgUpdateCollectionPermissions(): MsgUpdateCollectionPermissions {
  return { creator: "", collectionId: "", permissions: undefined, addressMappings: [] };
}

export const MsgUpdateCollectionPermissions = {
  encode(message: MsgUpdateCollectionPermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollectionPermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollectionPermissions();
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

  fromJSON(object: any): MsgUpdateCollectionPermissions {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      permissions: isSet(object.permissions) ? CollectionPermissions.fromJSON(object.permissions) : undefined,
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgUpdateCollectionPermissions): unknown {
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

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollectionPermissions>, I>>(
    object: I,
  ): MsgUpdateCollectionPermissions {
    const message = createBaseMsgUpdateCollectionPermissions();
    message.creator = object.creator ?? "";
    message.collectionId = object.collectionId ?? "";
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? CollectionPermissions.fromPartial(object.permissions)
      : undefined;
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgUpdateCollectionPermissionsResponse(): MsgUpdateCollectionPermissionsResponse {
  return {};
}

export const MsgUpdateCollectionPermissionsResponse = {
  encode(_: MsgUpdateCollectionPermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateCollectionPermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCollectionPermissionsResponse();
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

  fromJSON(_: any): MsgUpdateCollectionPermissionsResponse {
    return {};
  },

  toJSON(_: MsgUpdateCollectionPermissionsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateCollectionPermissionsResponse>, I>>(
    _: I,
  ): MsgUpdateCollectionPermissionsResponse {
    const message = createBaseMsgUpdateCollectionPermissionsResponse();
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

function createBaseClaimProofItem(): ClaimProofItem {
  return { aunt: "", onRight: false };
}

export const ClaimProofItem = {
  encode(message: ClaimProofItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.aunt !== "") {
      writer.uint32(10).string(message.aunt);
    }
    if (message.onRight === true) {
      writer.uint32(16).bool(message.onRight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClaimProofItem {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClaimProofItem();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.aunt = reader.string();
          break;
        case 2:
          message.onRight = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClaimProofItem {
    return {
      aunt: isSet(object.aunt) ? String(object.aunt) : "",
      onRight: isSet(object.onRight) ? Boolean(object.onRight) : false,
    };
  },

  toJSON(message: ClaimProofItem): unknown {
    const obj: any = {};
    message.aunt !== undefined && (obj.aunt = message.aunt);
    message.onRight !== undefined && (obj.onRight = message.onRight);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClaimProofItem>, I>>(object: I): ClaimProofItem {
    const message = createBaseClaimProofItem();
    message.aunt = object.aunt ?? "";
    message.onRight = object.onRight ?? false;
    return message;
  },
};

function createBaseClaimProof(): ClaimProof {
  return { leaf: "", aunts: [] };
}

export const ClaimProof = {
  encode(message: ClaimProof, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.leaf !== "") {
      writer.uint32(10).string(message.leaf);
    }
    for (const v of message.aunts) {
      ClaimProofItem.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClaimProof {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClaimProof();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.leaf = reader.string();
          break;
        case 2:
          message.aunts.push(ClaimProofItem.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClaimProof {
    return {
      leaf: isSet(object.leaf) ? String(object.leaf) : "",
      aunts: Array.isArray(object?.aunts) ? object.aunts.map((e: any) => ClaimProofItem.fromJSON(e)) : [],
    };
  },

  toJSON(message: ClaimProof): unknown {
    const obj: any = {};
    message.leaf !== undefined && (obj.leaf = message.leaf);
    if (message.aunts) {
      obj.aunts = message.aunts.map((e) => e ? ClaimProofItem.toJSON(e) : undefined);
    } else {
      obj.aunts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClaimProof>, I>>(object: I): ClaimProof {
    const message = createBaseClaimProof();
    message.leaf = object.leaf ?? "";
    message.aunts = object.aunts?.map((e) => ClaimProofItem.fromPartial(e)) || [];
    return message;
  },
};

function createBaseChallengeSolution(): ChallengeSolution {
  return { proof: undefined };
}

export const ChallengeSolution = {
  encode(message: ChallengeSolution, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.proof !== undefined) {
      ClaimProof.encode(message.proof, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChallengeSolution {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChallengeSolution();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.proof = ClaimProof.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ChallengeSolution {
    return { proof: isSet(object.proof) ? ClaimProof.fromJSON(object.proof) : undefined };
  },

  toJSON(message: ChallengeSolution): unknown {
    const obj: any = {};
    message.proof !== undefined && (obj.proof = message.proof ? ClaimProof.toJSON(message.proof) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ChallengeSolution>, I>>(object: I): ChallengeSolution {
    const message = createBaseChallengeSolution();
    message.proof = (object.proof !== undefined && object.proof !== null)
      ? ClaimProof.fromPartial(object.proof)
      : undefined;
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

/** Msg defines the Msg service. */
export interface Msg {
  NewCollection(request: MsgNewCollection): Promise<MsgNewCollectionResponse>;
  MintAndDistributeBadges(request: MsgMintAndDistributeBadges): Promise<MsgMintAndDistributeBadgesResponse>;
  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse>;
  UpdateCollectionApprovedTransfers(
    request: MsgUpdateCollectionApprovedTransfers,
  ): Promise<MsgUpdateCollectionApprovedTransfersResponse>;
  UpdateUserApprovedTransfers(request: MsgUpdateUserApprovedTransfers): Promise<MsgUpdateUserApprovedTransfersResponse>;
  UpdateMetadata(request: MsgUpdateMetadata): Promise<MsgUpdateMetadataResponse>;
  UpdateCollectionPermissions(request: MsgUpdateCollectionPermissions): Promise<MsgUpdateCollectionPermissionsResponse>;
  UpdateUserPermissions(request: MsgUpdateUserPermissions): Promise<MsgUpdateUserPermissionsResponse>;
  UpdateManager(request: MsgUpdateManager): Promise<MsgUpdateManagerResponse>;
  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  ArchiveCollection(request: MsgArchiveCollection): Promise<MsgArchiveCollectionResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.NewCollection = this.NewCollection.bind(this);
    this.MintAndDistributeBadges = this.MintAndDistributeBadges.bind(this);
    this.TransferBadges = this.TransferBadges.bind(this);
    this.UpdateCollectionApprovedTransfers = this.UpdateCollectionApprovedTransfers.bind(this);
    this.UpdateUserApprovedTransfers = this.UpdateUserApprovedTransfers.bind(this);
    this.UpdateMetadata = this.UpdateMetadata.bind(this);
    this.UpdateCollectionPermissions = this.UpdateCollectionPermissions.bind(this);
    this.UpdateUserPermissions = this.UpdateUserPermissions.bind(this);
    this.UpdateManager = this.UpdateManager.bind(this);
    this.DeleteCollection = this.DeleteCollection.bind(this);
    this.ArchiveCollection = this.ArchiveCollection.bind(this);
  }
  NewCollection(request: MsgNewCollection): Promise<MsgNewCollectionResponse> {
    const data = MsgNewCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "NewCollection", data);
    return promise.then((data) => MsgNewCollectionResponse.decode(new _m0.Reader(data)));
  }

  MintAndDistributeBadges(request: MsgMintAndDistributeBadges): Promise<MsgMintAndDistributeBadgesResponse> {
    const data = MsgMintAndDistributeBadges.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "MintAndDistributeBadges", data);
    return promise.then((data) => MsgMintAndDistributeBadgesResponse.decode(new _m0.Reader(data)));
  }

  TransferBadges(request: MsgTransferBadges): Promise<MsgTransferBadgesResponse> {
    const data = MsgTransferBadges.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "TransferBadges", data);
    return promise.then((data) => MsgTransferBadgesResponse.decode(new _m0.Reader(data)));
  }

  UpdateCollectionApprovedTransfers(
    request: MsgUpdateCollectionApprovedTransfers,
  ): Promise<MsgUpdateCollectionApprovedTransfersResponse> {
    const data = MsgUpdateCollectionApprovedTransfers.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateCollectionApprovedTransfers", data);
    return promise.then((data) => MsgUpdateCollectionApprovedTransfersResponse.decode(new _m0.Reader(data)));
  }

  UpdateUserApprovedTransfers(
    request: MsgUpdateUserApprovedTransfers,
  ): Promise<MsgUpdateUserApprovedTransfersResponse> {
    const data = MsgUpdateUserApprovedTransfers.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateUserApprovedTransfers", data);
    return promise.then((data) => MsgUpdateUserApprovedTransfersResponse.decode(new _m0.Reader(data)));
  }

  UpdateMetadata(request: MsgUpdateMetadata): Promise<MsgUpdateMetadataResponse> {
    const data = MsgUpdateMetadata.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateMetadata", data);
    return promise.then((data) => MsgUpdateMetadataResponse.decode(new _m0.Reader(data)));
  }

  UpdateCollectionPermissions(
    request: MsgUpdateCollectionPermissions,
  ): Promise<MsgUpdateCollectionPermissionsResponse> {
    const data = MsgUpdateCollectionPermissions.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateCollectionPermissions", data);
    return promise.then((data) => MsgUpdateCollectionPermissionsResponse.decode(new _m0.Reader(data)));
  }

  UpdateUserPermissions(request: MsgUpdateUserPermissions): Promise<MsgUpdateUserPermissionsResponse> {
    const data = MsgUpdateUserPermissions.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateUserPermissions", data);
    return promise.then((data) => MsgUpdateUserPermissionsResponse.decode(new _m0.Reader(data)));
  }

  UpdateManager(request: MsgUpdateManager): Promise<MsgUpdateManagerResponse> {
    const data = MsgUpdateManager.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateManager", data);
    return promise.then((data) => MsgUpdateManagerResponse.decode(new _m0.Reader(data)));
  }

  DeleteCollection(request: MsgDeleteCollection): Promise<MsgDeleteCollectionResponse> {
    const data = MsgDeleteCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "DeleteCollection", data);
    return promise.then((data) => MsgDeleteCollectionResponse.decode(new _m0.Reader(data)));
  }

  ArchiveCollection(request: MsgArchiveCollection): Promise<MsgArchiveCollectionResponse> {
    const data = MsgArchiveCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "ArchiveCollection", data);
    return promise.then((data) => MsgArchiveCollectionResponse.decode(new _m0.Reader(data)));
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
