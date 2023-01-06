/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { WhitelistMintInfo } from "./balances";
import { IdRange } from "./ranges";
import { UriObject } from "./uris";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

export interface MsgNewBadge {
  /** See badges.proto for more details about these MsgNewBadge fields. Defines the badge details. Leave unneeded fields empty. */
  creator: string;
  uri: UriObject | undefined;
  permissions: number;
  arbitraryBytes: string;
  defaultSubassetSupply: number;
  freezeAddressRanges: IdRange[];
  standard: number;
  /**
   * Subasset supplys and amounts to create must be same length. For each idx, we create amounts[idx] subbadges each with a supply of supplys[idx].
   * If supply[idx] == 0, we assume default supply. amountsToCreate[idx] can't equal 0.
   * TODO: convert this into one struct?
   */
  subassetSupplys: number[];
  subassetAmountsToCreate: number[];
  whitelistedRecipients: WhitelistMintInfo[];
}

export interface MsgNewBadgeResponse {
  /** ID of created badge */
  id: number;
}

export interface MsgNewSubBadge {
  creator: string;
  badgeId: number;
  /**
   * Subasset supplys and amounts to create must be same length. For each idx, we create amounts[idx] subbadges each with a supply of supplys[idx].
   * If supply[idx] == 0, we assume default supply. amountsToCreate[idx] can't equal 0.
   */
  supplys: number[];
  amountsToCreate: number[];
}

export interface MsgNewSubBadgeResponse {
  /** ID of next subbadgeId after creating all subbadges. */
  nextSubassetId: number;
}

/** For each amount, for each toAddress, we will attempt to transfer all the subbadgeIds for the badge with ID badgeId. */
export interface MsgTransferBadge {
  creator: string;
  from: number;
  toAddresses: number[];
  amounts: number[];
  badgeId: number;
  subbadgeRanges: IdRange[];
  /** If 0, never expires and assumed to be the max possible time. */
  expirationTime: number;
  /** If 0, always cancellable. Must be <= expirationTime. */
  cantCancelBeforeTime: number;
}

export interface MsgTransferBadgeResponse {
}

/** For each amount, for each toAddress, we will request a transfer all the subbadgeIds for the badge with ID badgeId. Other party must approve / reject the transfer request. */
export interface MsgRequestTransferBadge {
  creator: string;
  from: number;
  amount: number;
  badgeId: number;
  subbadgeRanges: IdRange[];
  /** If 0, never expires and assumed to be the max possible time. */
  expirationTime: number;
  /** If 0, always cancellable. Must be <= expirationTime. */
  cantCancelBeforeTime: number;
}

export interface MsgRequestTransferBadgeResponse {
}

/** For all pending transfers of the badge where ThisPendingNonce is within some nonceRange in nonceRanges, we accept or deny the pending transfer. */
export interface MsgHandlePendingTransfer {
  creator: string;
  accept: boolean;
  badgeId: number;
  nonceRanges: IdRange[];
  /** Forceful accept is an option to accept the transfer forcefully instead of just marking it as approved. */
  forcefulAccept: boolean;
}

export interface MsgHandlePendingTransferResponse {
}

/** Sets an approval (no add or remove), just set it for an address. */
export interface MsgSetApproval {
  creator: string;
  amount: number;
  address: number;
  badgeId: number;
  subbadgeRanges: IdRange[];
}

export interface MsgSetApprovalResponse {
}

/** For each address and for each amount, revoke badge. */
export interface MsgRevokeBadge {
  creator: string;
  addresses: number[];
  amounts: number[];
  badgeId: number;
  subbadgeRanges: IdRange[];
}

export interface MsgRevokeBadgeResponse {
}

/** Add or remove addreses from the freeze address range */
export interface MsgFreezeAddress {
  creator: string;
  addressRanges: IdRange[];
  badgeId: number;
  add: boolean;
}

export interface MsgFreezeAddressResponse {
}

/** Update badge Uris with new URI object, if permitted. */
export interface MsgUpdateUris {
  creator: string;
  badgeId: number;
  uri: UriObject | undefined;
}

export interface MsgUpdateUrisResponse {
}

/** Update badge permissions with new permissions, if permitted. */
export interface MsgUpdatePermissions {
  creator: string;
  badgeId: number;
  permissions: number;
}

export interface MsgUpdatePermissionsResponse {
}

/** Transfer manager to this address. Recipient must have made a request. */
export interface MsgTransferManager {
  creator: string;
  badgeId: number;
  address: number;
}

export interface MsgTransferManagerResponse {
}

/** Add / remove request for manager to be transferred. */
export interface MsgRequestTransferManager {
  creator: string;
  badgeId: number;
  add: boolean;
}

export interface MsgRequestTransferManagerResponse {
}

/** Self destructs the badge, if permitted. */
export interface MsgSelfDestructBadge {
  creator: string;
  badgeId: number;
}

export interface MsgSelfDestructBadgeResponse {
}

/** Prunes balances of self destructed badges. Can be called by anyone */
export interface MsgPruneBalances {
  creator: string;
  badgeIds: number[];
  addresses: number[];
}

export interface MsgPruneBalancesResponse {
}

/** Update badge bytes, if permitted */
export interface MsgUpdateBytes {
  creator: string;
  badgeId: number;
  newBytes: string;
}

export interface MsgUpdateBytesResponse {
}

export interface MsgRegisterAddresses {
  creator: string;
  addressesToRegister: string[];
}

export interface MsgRegisterAddressesResponse {
  registeredAddressNumbers: IdRange | undefined;
}

function createBaseMsgNewBadge(): MsgNewBadge {
  return {
    creator: "",
    uri: undefined,
    permissions: 0,
    arbitraryBytes: "",
    defaultSubassetSupply: 0,
    freezeAddressRanges: [],
    standard: 0,
    subassetSupplys: [],
    subassetAmountsToCreate: [],
    whitelistedRecipients: [],
  };
}

export const MsgNewBadge = {
  encode(message: MsgNewBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.uri !== undefined) {
      UriObject.encode(message.uri, writer.uint32(18).fork()).ldelim();
    }
    if (message.permissions !== 0) {
      writer.uint32(32).uint64(message.permissions);
    }
    if (message.arbitraryBytes !== "") {
      writer.uint32(42).string(message.arbitraryBytes);
    }
    if (message.defaultSubassetSupply !== 0) {
      writer.uint32(48).uint64(message.defaultSubassetSupply);
    }
    for (const v of message.freezeAddressRanges) {
      IdRange.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    if (message.standard !== 0) {
      writer.uint32(80).uint64(message.standard);
    }
    writer.uint32(58).fork();
    for (const v of message.subassetSupplys) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(66).fork();
    for (const v of message.subassetAmountsToCreate) {
      writer.uint64(v);
    }
    writer.ldelim();
    for (const v of message.whitelistedRecipients) {
      WhitelistMintInfo.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.uri = UriObject.decode(reader, reader.uint32());
          break;
        case 4:
          message.permissions = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.arbitraryBytes = reader.string();
          break;
        case 6:
          message.defaultSubassetSupply = longToNumber(reader.uint64() as Long);
          break;
        case 9:
          message.freezeAddressRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 10:
          message.standard = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.subassetSupplys.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.subassetSupplys.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 8:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.subassetAmountsToCreate.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.subassetAmountsToCreate.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 11:
          message.whitelistedRecipients.push(WhitelistMintInfo.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      uri: isSet(object.uri) ? UriObject.fromJSON(object.uri) : undefined,
      permissions: isSet(object.permissions) ? Number(object.permissions) : 0,
      arbitraryBytes: isSet(object.arbitraryBytes) ? String(object.arbitraryBytes) : "",
      defaultSubassetSupply: isSet(object.defaultSubassetSupply) ? Number(object.defaultSubassetSupply) : 0,
      freezeAddressRanges: Array.isArray(object?.freezeAddressRanges)
        ? object.freezeAddressRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
      standard: isSet(object.standard) ? Number(object.standard) : 0,
      subassetSupplys: Array.isArray(object?.subassetSupplys) ? object.subassetSupplys.map((e: any) => Number(e)) : [],
      subassetAmountsToCreate: Array.isArray(object?.subassetAmountsToCreate)
        ? object.subassetAmountsToCreate.map((e: any) => Number(e))
        : [],
      whitelistedRecipients: Array.isArray(object?.whitelistedRecipients)
        ? object.whitelistedRecipients.map((e: any) => WhitelistMintInfo.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgNewBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.uri !== undefined && (obj.uri = message.uri ? UriObject.toJSON(message.uri) : undefined);
    message.permissions !== undefined && (obj.permissions = Math.round(message.permissions));
    message.arbitraryBytes !== undefined && (obj.arbitraryBytes = message.arbitraryBytes);
    message.defaultSubassetSupply !== undefined
      && (obj.defaultSubassetSupply = Math.round(message.defaultSubassetSupply));
    if (message.freezeAddressRanges) {
      obj.freezeAddressRanges = message.freezeAddressRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.freezeAddressRanges = [];
    }
    message.standard !== undefined && (obj.standard = Math.round(message.standard));
    if (message.subassetSupplys) {
      obj.subassetSupplys = message.subassetSupplys.map((e) => Math.round(e));
    } else {
      obj.subassetSupplys = [];
    }
    if (message.subassetAmountsToCreate) {
      obj.subassetAmountsToCreate = message.subassetAmountsToCreate.map((e) => Math.round(e));
    } else {
      obj.subassetAmountsToCreate = [];
    }
    if (message.whitelistedRecipients) {
      obj.whitelistedRecipients = message.whitelistedRecipients.map((e) => e ? WhitelistMintInfo.toJSON(e) : undefined);
    } else {
      obj.whitelistedRecipients = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewBadge>, I>>(object: I): MsgNewBadge {
    const message = createBaseMsgNewBadge();
    message.creator = object.creator ?? "";
    message.uri = (object.uri !== undefined && object.uri !== null) ? UriObject.fromPartial(object.uri) : undefined;
    message.permissions = object.permissions ?? 0;
    message.arbitraryBytes = object.arbitraryBytes ?? "";
    message.defaultSubassetSupply = object.defaultSubassetSupply ?? 0;
    message.freezeAddressRanges = object.freezeAddressRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.standard = object.standard ?? 0;
    message.subassetSupplys = object.subassetSupplys?.map((e) => e) || [];
    message.subassetAmountsToCreate = object.subassetAmountsToCreate?.map((e) => e) || [];
    message.whitelistedRecipients = object.whitelistedRecipients?.map((e) => WhitelistMintInfo.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgNewBadgeResponse(): MsgNewBadgeResponse {
  return { id: 0 };
}

export const MsgNewBadgeResponse = {
  encode(message: MsgNewBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewBadgeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewBadgeResponse {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: MsgNewBadgeResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewBadgeResponse>, I>>(object: I): MsgNewBadgeResponse {
    const message = createBaseMsgNewBadgeResponse();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseMsgNewSubBadge(): MsgNewSubBadge {
  return { creator: "", badgeId: 0, supplys: [], amountsToCreate: [] };
}

export const MsgNewSubBadge = {
  encode(message: MsgNewSubBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    writer.uint32(26).fork();
    for (const v of message.supplys) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(34).fork();
    for (const v of message.amountsToCreate) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewSubBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewSubBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.supplys.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.supplys.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 4:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.amountsToCreate.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.amountsToCreate.push(longToNumber(reader.uint64() as Long));
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewSubBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      supplys: Array.isArray(object?.supplys) ? object.supplys.map((e: any) => Number(e)) : [],
      amountsToCreate: Array.isArray(object?.amountsToCreate) ? object.amountsToCreate.map((e: any) => Number(e)) : [],
    };
  },

  toJSON(message: MsgNewSubBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.supplys) {
      obj.supplys = message.supplys.map((e) => Math.round(e));
    } else {
      obj.supplys = [];
    }
    if (message.amountsToCreate) {
      obj.amountsToCreate = message.amountsToCreate.map((e) => Math.round(e));
    } else {
      obj.amountsToCreate = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewSubBadge>, I>>(object: I): MsgNewSubBadge {
    const message = createBaseMsgNewSubBadge();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.supplys = object.supplys?.map((e) => e) || [];
    message.amountsToCreate = object.amountsToCreate?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgNewSubBadgeResponse(): MsgNewSubBadgeResponse {
  return { nextSubassetId: 0 };
}

export const MsgNewSubBadgeResponse = {
  encode(message: MsgNewSubBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextSubassetId !== 0) {
      writer.uint32(8).uint64(message.nextSubassetId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgNewSubBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgNewSubBadgeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nextSubassetId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewSubBadgeResponse {
    return { nextSubassetId: isSet(object.nextSubassetId) ? Number(object.nextSubassetId) : 0 };
  },

  toJSON(message: MsgNewSubBadgeResponse): unknown {
    const obj: any = {};
    message.nextSubassetId !== undefined && (obj.nextSubassetId = Math.round(message.nextSubassetId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgNewSubBadgeResponse>, I>>(object: I): MsgNewSubBadgeResponse {
    const message = createBaseMsgNewSubBadgeResponse();
    message.nextSubassetId = object.nextSubassetId ?? 0;
    return message;
  },
};

function createBaseMsgTransferBadge(): MsgTransferBadge {
  return {
    creator: "",
    from: 0,
    toAddresses: [],
    amounts: [],
    badgeId: 0,
    subbadgeRanges: [],
    expirationTime: 0,
    cantCancelBeforeTime: 0,
  };
}

export const MsgTransferBadge = {
  encode(message: MsgTransferBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.from !== 0) {
      writer.uint32(16).uint64(message.from);
    }
    writer.uint32(26).fork();
    for (const v of message.toAddresses) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(34).fork();
    for (const v of message.amounts) {
      writer.uint64(v);
    }
    writer.ldelim();
    if (message.badgeId !== 0) {
      writer.uint32(40).uint64(message.badgeId);
    }
    for (const v of message.subbadgeRanges) {
      IdRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.expirationTime !== 0) {
      writer.uint32(56).uint64(message.expirationTime);
    }
    if (message.cantCancelBeforeTime !== 0) {
      writer.uint32(64).uint64(message.cantCancelBeforeTime);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.from = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.toAddresses.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.toAddresses.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 4:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.amounts.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.amounts.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 5:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.subbadgeRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.expirationTime = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.cantCancelBeforeTime = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgTransferBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      from: isSet(object.from) ? Number(object.from) : 0,
      toAddresses: Array.isArray(object?.toAddresses) ? object.toAddresses.map((e: any) => Number(e)) : [],
      amounts: Array.isArray(object?.amounts) ? object.amounts.map((e: any) => Number(e)) : [],
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      subbadgeRanges: Array.isArray(object?.subbadgeRanges)
        ? object.subbadgeRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
      expirationTime: isSet(object.expirationTime) ? Number(object.expirationTime) : 0,
      cantCancelBeforeTime: isSet(object.cantCancelBeforeTime) ? Number(object.cantCancelBeforeTime) : 0,
    };
  },

  toJSON(message: MsgTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = Math.round(message.from));
    if (message.toAddresses) {
      obj.toAddresses = message.toAddresses.map((e) => Math.round(e));
    } else {
      obj.toAddresses = [];
    }
    if (message.amounts) {
      obj.amounts = message.amounts.map((e) => Math.round(e));
    } else {
      obj.amounts = [];
    }
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.subbadgeRanges = [];
    }
    message.expirationTime !== undefined && (obj.expirationTime = Math.round(message.expirationTime));
    message.cantCancelBeforeTime !== undefined && (obj.cantCancelBeforeTime = Math.round(message.cantCancelBeforeTime));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferBadge>, I>>(object: I): MsgTransferBadge {
    const message = createBaseMsgTransferBadge();
    message.creator = object.creator ?? "";
    message.from = object.from ?? 0;
    message.toAddresses = object.toAddresses?.map((e) => e) || [];
    message.amounts = object.amounts?.map((e) => e) || [];
    message.badgeId = object.badgeId ?? 0;
    message.subbadgeRanges = object.subbadgeRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.expirationTime = object.expirationTime ?? 0;
    message.cantCancelBeforeTime = object.cantCancelBeforeTime ?? 0;
    return message;
  },
};

function createBaseMsgTransferBadgeResponse(): MsgTransferBadgeResponse {
  return {};
}

export const MsgTransferBadgeResponse = {
  encode(_: MsgTransferBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferBadgeResponse();
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

  fromJSON(_: any): MsgTransferBadgeResponse {
    return {};
  },

  toJSON(_: MsgTransferBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferBadgeResponse>, I>>(_: I): MsgTransferBadgeResponse {
    const message = createBaseMsgTransferBadgeResponse();
    return message;
  },
};

function createBaseMsgRequestTransferBadge(): MsgRequestTransferBadge {
  return {
    creator: "",
    from: 0,
    amount: 0,
    badgeId: 0,
    subbadgeRanges: [],
    expirationTime: 0,
    cantCancelBeforeTime: 0,
  };
}

export const MsgRequestTransferBadge = {
  encode(message: MsgRequestTransferBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.from !== 0) {
      writer.uint32(16).uint64(message.from);
    }
    if (message.amount !== 0) {
      writer.uint32(32).uint64(message.amount);
    }
    if (message.badgeId !== 0) {
      writer.uint32(40).uint64(message.badgeId);
    }
    for (const v of message.subbadgeRanges) {
      IdRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.expirationTime !== 0) {
      writer.uint32(56).uint64(message.expirationTime);
    }
    if (message.cantCancelBeforeTime !== 0) {
      writer.uint32(64).uint64(message.cantCancelBeforeTime);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestTransferBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestTransferBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.from = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.subbadgeRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.expirationTime = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.cantCancelBeforeTime = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRequestTransferBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      from: isSet(object.from) ? Number(object.from) : 0,
      amount: isSet(object.amount) ? Number(object.amount) : 0,
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      subbadgeRanges: Array.isArray(object?.subbadgeRanges)
        ? object.subbadgeRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
      expirationTime: isSet(object.expirationTime) ? Number(object.expirationTime) : 0,
      cantCancelBeforeTime: isSet(object.cantCancelBeforeTime) ? Number(object.cantCancelBeforeTime) : 0,
    };
  },

  toJSON(message: MsgRequestTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = Math.round(message.from));
    message.amount !== undefined && (obj.amount = Math.round(message.amount));
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.subbadgeRanges = [];
    }
    message.expirationTime !== undefined && (obj.expirationTime = Math.round(message.expirationTime));
    message.cantCancelBeforeTime !== undefined && (obj.cantCancelBeforeTime = Math.round(message.cantCancelBeforeTime));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRequestTransferBadge>, I>>(object: I): MsgRequestTransferBadge {
    const message = createBaseMsgRequestTransferBadge();
    message.creator = object.creator ?? "";
    message.from = object.from ?? 0;
    message.amount = object.amount ?? 0;
    message.badgeId = object.badgeId ?? 0;
    message.subbadgeRanges = object.subbadgeRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.expirationTime = object.expirationTime ?? 0;
    message.cantCancelBeforeTime = object.cantCancelBeforeTime ?? 0;
    return message;
  },
};

function createBaseMsgRequestTransferBadgeResponse(): MsgRequestTransferBadgeResponse {
  return {};
}

export const MsgRequestTransferBadgeResponse = {
  encode(_: MsgRequestTransferBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestTransferBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestTransferBadgeResponse();
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

  fromJSON(_: any): MsgRequestTransferBadgeResponse {
    return {};
  },

  toJSON(_: MsgRequestTransferBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRequestTransferBadgeResponse>, I>>(_: I): MsgRequestTransferBadgeResponse {
    const message = createBaseMsgRequestTransferBadgeResponse();
    return message;
  },
};

function createBaseMsgHandlePendingTransfer(): MsgHandlePendingTransfer {
  return { creator: "", accept: false, badgeId: 0, nonceRanges: [], forcefulAccept: false };
}

export const MsgHandlePendingTransfer = {
  encode(message: MsgHandlePendingTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.accept === true) {
      writer.uint32(16).bool(message.accept);
    }
    if (message.badgeId !== 0) {
      writer.uint32(24).uint64(message.badgeId);
    }
    for (const v of message.nonceRanges) {
      IdRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.forcefulAccept === true) {
      writer.uint32(40).bool(message.forcefulAccept);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgHandlePendingTransfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgHandlePendingTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.accept = reader.bool();
          break;
        case 3:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.nonceRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.forcefulAccept = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgHandlePendingTransfer {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      accept: isSet(object.accept) ? Boolean(object.accept) : false,
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      nonceRanges: Array.isArray(object?.nonceRanges) ? object.nonceRanges.map((e: any) => IdRange.fromJSON(e)) : [],
      forcefulAccept: isSet(object.forcefulAccept) ? Boolean(object.forcefulAccept) : false,
    };
  },

  toJSON(message: MsgHandlePendingTransfer): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.accept !== undefined && (obj.accept = message.accept);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.nonceRanges) {
      obj.nonceRanges = message.nonceRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.nonceRanges = [];
    }
    message.forcefulAccept !== undefined && (obj.forcefulAccept = message.forcefulAccept);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgHandlePendingTransfer>, I>>(object: I): MsgHandlePendingTransfer {
    const message = createBaseMsgHandlePendingTransfer();
    message.creator = object.creator ?? "";
    message.accept = object.accept ?? false;
    message.badgeId = object.badgeId ?? 0;
    message.nonceRanges = object.nonceRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.forcefulAccept = object.forcefulAccept ?? false;
    return message;
  },
};

function createBaseMsgHandlePendingTransferResponse(): MsgHandlePendingTransferResponse {
  return {};
}

export const MsgHandlePendingTransferResponse = {
  encode(_: MsgHandlePendingTransferResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgHandlePendingTransferResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgHandlePendingTransferResponse();
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

  fromJSON(_: any): MsgHandlePendingTransferResponse {
    return {};
  },

  toJSON(_: MsgHandlePendingTransferResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgHandlePendingTransferResponse>, I>>(
    _: I,
  ): MsgHandlePendingTransferResponse {
    const message = createBaseMsgHandlePendingTransferResponse();
    return message;
  },
};

function createBaseMsgSetApproval(): MsgSetApproval {
  return { creator: "", amount: 0, address: 0, badgeId: 0, subbadgeRanges: [] };
}

export const MsgSetApproval = {
  encode(message: MsgSetApproval, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.amount !== 0) {
      writer.uint32(16).uint64(message.amount);
    }
    if (message.address !== 0) {
      writer.uint32(24).uint64(message.address);
    }
    if (message.badgeId !== 0) {
      writer.uint32(32).uint64(message.badgeId);
    }
    for (const v of message.subbadgeRanges) {
      IdRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetApproval {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetApproval();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.subbadgeRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetApproval {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      amount: isSet(object.amount) ? Number(object.amount) : 0,
      address: isSet(object.address) ? Number(object.address) : 0,
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      subbadgeRanges: Array.isArray(object?.subbadgeRanges)
        ? object.subbadgeRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgSetApproval): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.amount !== undefined && (obj.amount = Math.round(message.amount));
    message.address !== undefined && (obj.address = Math.round(message.address));
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetApproval>, I>>(object: I): MsgSetApproval {
    const message = createBaseMsgSetApproval();
    message.creator = object.creator ?? "";
    message.amount = object.amount ?? 0;
    message.address = object.address ?? 0;
    message.badgeId = object.badgeId ?? 0;
    message.subbadgeRanges = object.subbadgeRanges?.map((e) => IdRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgSetApprovalResponse(): MsgSetApprovalResponse {
  return {};
}

export const MsgSetApprovalResponse = {
  encode(_: MsgSetApprovalResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetApprovalResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetApprovalResponse();
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

  fromJSON(_: any): MsgSetApprovalResponse {
    return {};
  },

  toJSON(_: MsgSetApprovalResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSetApprovalResponse>, I>>(_: I): MsgSetApprovalResponse {
    const message = createBaseMsgSetApprovalResponse();
    return message;
  },
};

function createBaseMsgRevokeBadge(): MsgRevokeBadge {
  return { creator: "", addresses: [], amounts: [], badgeId: 0, subbadgeRanges: [] };
}

export const MsgRevokeBadge = {
  encode(message: MsgRevokeBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    writer.uint32(18).fork();
    for (const v of message.addresses) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(26).fork();
    for (const v of message.amounts) {
      writer.uint64(v);
    }
    writer.ldelim();
    if (message.badgeId !== 0) {
      writer.uint32(32).uint64(message.badgeId);
    }
    for (const v of message.subbadgeRanges) {
      IdRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRevokeBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRevokeBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.addresses.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.addresses.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.amounts.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.amounts.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 4:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.subbadgeRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRevokeBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      addresses: Array.isArray(object?.addresses) ? object.addresses.map((e: any) => Number(e)) : [],
      amounts: Array.isArray(object?.amounts) ? object.amounts.map((e: any) => Number(e)) : [],
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      subbadgeRanges: Array.isArray(object?.subbadgeRanges)
        ? object.subbadgeRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MsgRevokeBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => Math.round(e));
    } else {
      obj.addresses = [];
    }
    if (message.amounts) {
      obj.amounts = message.amounts.map((e) => Math.round(e));
    } else {
      obj.amounts = [];
    }
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRevokeBadge>, I>>(object: I): MsgRevokeBadge {
    const message = createBaseMsgRevokeBadge();
    message.creator = object.creator ?? "";
    message.addresses = object.addresses?.map((e) => e) || [];
    message.amounts = object.amounts?.map((e) => e) || [];
    message.badgeId = object.badgeId ?? 0;
    message.subbadgeRanges = object.subbadgeRanges?.map((e) => IdRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMsgRevokeBadgeResponse(): MsgRevokeBadgeResponse {
  return {};
}

export const MsgRevokeBadgeResponse = {
  encode(_: MsgRevokeBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRevokeBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRevokeBadgeResponse();
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

  fromJSON(_: any): MsgRevokeBadgeResponse {
    return {};
  },

  toJSON(_: MsgRevokeBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRevokeBadgeResponse>, I>>(_: I): MsgRevokeBadgeResponse {
    const message = createBaseMsgRevokeBadgeResponse();
    return message;
  },
};

function createBaseMsgFreezeAddress(): MsgFreezeAddress {
  return { creator: "", addressRanges: [], badgeId: 0, add: false };
}

export const MsgFreezeAddress = {
  encode(message: MsgFreezeAddress, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.addressRanges) {
      IdRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.badgeId !== 0) {
      writer.uint32(24).uint64(message.badgeId);
    }
    if (message.add === true) {
      writer.uint32(32).bool(message.add);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFreezeAddress {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFreezeAddress();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.addressRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 3:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.add = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgFreezeAddress {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      addressRanges: Array.isArray(object?.addressRanges)
        ? object.addressRanges.map((e: any) => IdRange.fromJSON(e))
        : [],
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      add: isSet(object.add) ? Boolean(object.add) : false,
    };
  },

  toJSON(message: MsgFreezeAddress): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addressRanges) {
      obj.addressRanges = message.addressRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.addressRanges = [];
    }
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.add !== undefined && (obj.add = message.add);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFreezeAddress>, I>>(object: I): MsgFreezeAddress {
    const message = createBaseMsgFreezeAddress();
    message.creator = object.creator ?? "";
    message.addressRanges = object.addressRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.badgeId = object.badgeId ?? 0;
    message.add = object.add ?? false;
    return message;
  },
};

function createBaseMsgFreezeAddressResponse(): MsgFreezeAddressResponse {
  return {};
}

export const MsgFreezeAddressResponse = {
  encode(_: MsgFreezeAddressResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFreezeAddressResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFreezeAddressResponse();
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

  fromJSON(_: any): MsgFreezeAddressResponse {
    return {};
  },

  toJSON(_: MsgFreezeAddressResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFreezeAddressResponse>, I>>(_: I): MsgFreezeAddressResponse {
    const message = createBaseMsgFreezeAddressResponse();
    return message;
  },
};

function createBaseMsgUpdateUris(): MsgUpdateUris {
  return { creator: "", badgeId: 0, uri: undefined };
}

export const MsgUpdateUris = {
  encode(message: MsgUpdateUris, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.uri !== undefined) {
      UriObject.encode(message.uri, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUris {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUris();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.uri = UriObject.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateUris {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      uri: isSet(object.uri) ? UriObject.fromJSON(object.uri) : undefined,
    };
  },

  toJSON(message: MsgUpdateUris): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.uri !== undefined && (obj.uri = message.uri ? UriObject.toJSON(message.uri) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUris>, I>>(object: I): MsgUpdateUris {
    const message = createBaseMsgUpdateUris();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.uri = (object.uri !== undefined && object.uri !== null) ? UriObject.fromPartial(object.uri) : undefined;
    return message;
  },
};

function createBaseMsgUpdateUrisResponse(): MsgUpdateUrisResponse {
  return {};
}

export const MsgUpdateUrisResponse = {
  encode(_: MsgUpdateUrisResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUrisResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUrisResponse();
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

  fromJSON(_: any): MsgUpdateUrisResponse {
    return {};
  },

  toJSON(_: MsgUpdateUrisResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateUrisResponse>, I>>(_: I): MsgUpdateUrisResponse {
    const message = createBaseMsgUpdateUrisResponse();
    return message;
  },
};

function createBaseMsgUpdatePermissions(): MsgUpdatePermissions {
  return { creator: "", badgeId: 0, permissions: 0 };
}

export const MsgUpdatePermissions = {
  encode(message: MsgUpdatePermissions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.permissions !== 0) {
      writer.uint32(24).uint64(message.permissions);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePermissions {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePermissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.permissions = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdatePermissions {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      permissions: isSet(object.permissions) ? Number(object.permissions) : 0,
    };
  },

  toJSON(message: MsgUpdatePermissions): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.permissions !== undefined && (obj.permissions = Math.round(message.permissions));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdatePermissions>, I>>(object: I): MsgUpdatePermissions {
    const message = createBaseMsgUpdatePermissions();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.permissions = object.permissions ?? 0;
    return message;
  },
};

function createBaseMsgUpdatePermissionsResponse(): MsgUpdatePermissionsResponse {
  return {};
}

export const MsgUpdatePermissionsResponse = {
  encode(_: MsgUpdatePermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePermissionsResponse();
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

  fromJSON(_: any): MsgUpdatePermissionsResponse {
    return {};
  },

  toJSON(_: MsgUpdatePermissionsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdatePermissionsResponse>, I>>(_: I): MsgUpdatePermissionsResponse {
    const message = createBaseMsgUpdatePermissionsResponse();
    return message;
  },
};

function createBaseMsgTransferManager(): MsgTransferManager {
  return { creator: "", badgeId: 0, address: 0 };
}

export const MsgTransferManager = {
  encode(message: MsgTransferManager, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.address !== 0) {
      writer.uint32(24).uint64(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferManager {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferManager();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgTransferManager {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      address: isSet(object.address) ? Number(object.address) : 0,
    };
  },

  toJSON(message: MsgTransferManager): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.address !== undefined && (obj.address = Math.round(message.address));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferManager>, I>>(object: I): MsgTransferManager {
    const message = createBaseMsgTransferManager();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.address = object.address ?? 0;
    return message;
  },
};

function createBaseMsgTransferManagerResponse(): MsgTransferManagerResponse {
  return {};
}

export const MsgTransferManagerResponse = {
  encode(_: MsgTransferManagerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgTransferManagerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgTransferManagerResponse();
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

  fromJSON(_: any): MsgTransferManagerResponse {
    return {};
  },

  toJSON(_: MsgTransferManagerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgTransferManagerResponse>, I>>(_: I): MsgTransferManagerResponse {
    const message = createBaseMsgTransferManagerResponse();
    return message;
  },
};

function createBaseMsgRequestTransferManager(): MsgRequestTransferManager {
  return { creator: "", badgeId: 0, add: false };
}

export const MsgRequestTransferManager = {
  encode(message: MsgRequestTransferManager, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.add === true) {
      writer.uint32(24).bool(message.add);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestTransferManager {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestTransferManager();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.add = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRequestTransferManager {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      add: isSet(object.add) ? Boolean(object.add) : false,
    };
  },

  toJSON(message: MsgRequestTransferManager): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.add !== undefined && (obj.add = message.add);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRequestTransferManager>, I>>(object: I): MsgRequestTransferManager {
    const message = createBaseMsgRequestTransferManager();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.add = object.add ?? false;
    return message;
  },
};

function createBaseMsgRequestTransferManagerResponse(): MsgRequestTransferManagerResponse {
  return {};
}

export const MsgRequestTransferManagerResponse = {
  encode(_: MsgRequestTransferManagerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRequestTransferManagerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRequestTransferManagerResponse();
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

  fromJSON(_: any): MsgRequestTransferManagerResponse {
    return {};
  },

  toJSON(_: MsgRequestTransferManagerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRequestTransferManagerResponse>, I>>(
    _: I,
  ): MsgRequestTransferManagerResponse {
    const message = createBaseMsgRequestTransferManagerResponse();
    return message;
  },
};

function createBaseMsgSelfDestructBadge(): MsgSelfDestructBadge {
  return { creator: "", badgeId: 0 };
}

export const MsgSelfDestructBadge = {
  encode(message: MsgSelfDestructBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSelfDestructBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSelfDestructBadge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSelfDestructBadge {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
    };
  },

  toJSON(message: MsgSelfDestructBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSelfDestructBadge>, I>>(object: I): MsgSelfDestructBadge {
    const message = createBaseMsgSelfDestructBadge();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    return message;
  },
};

function createBaseMsgSelfDestructBadgeResponse(): MsgSelfDestructBadgeResponse {
  return {};
}

export const MsgSelfDestructBadgeResponse = {
  encode(_: MsgSelfDestructBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSelfDestructBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSelfDestructBadgeResponse();
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

  fromJSON(_: any): MsgSelfDestructBadgeResponse {
    return {};
  },

  toJSON(_: MsgSelfDestructBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgSelfDestructBadgeResponse>, I>>(_: I): MsgSelfDestructBadgeResponse {
    const message = createBaseMsgSelfDestructBadgeResponse();
    return message;
  },
};

function createBaseMsgPruneBalances(): MsgPruneBalances {
  return { creator: "", badgeIds: [], addresses: [] };
}

export const MsgPruneBalances = {
  encode(message: MsgPruneBalances, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    writer.uint32(18).fork();
    for (const v of message.badgeIds) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(26).fork();
    for (const v of message.addresses) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgPruneBalances {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgPruneBalances();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.badgeIds.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.badgeIds.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.addresses.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.addresses.push(longToNumber(reader.uint64() as Long));
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgPruneBalances {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => Number(e)) : [],
      addresses: Array.isArray(object?.addresses) ? object.addresses.map((e: any) => Number(e)) : [],
    };
  },

  toJSON(message: MsgPruneBalances): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => Math.round(e));
    } else {
      obj.badgeIds = [];
    }
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => Math.round(e));
    } else {
      obj.addresses = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgPruneBalances>, I>>(object: I): MsgPruneBalances {
    const message = createBaseMsgPruneBalances();
    message.creator = object.creator ?? "";
    message.badgeIds = object.badgeIds?.map((e) => e) || [];
    message.addresses = object.addresses?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgPruneBalancesResponse(): MsgPruneBalancesResponse {
  return {};
}

export const MsgPruneBalancesResponse = {
  encode(_: MsgPruneBalancesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgPruneBalancesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgPruneBalancesResponse();
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

  fromJSON(_: any): MsgPruneBalancesResponse {
    return {};
  },

  toJSON(_: MsgPruneBalancesResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgPruneBalancesResponse>, I>>(_: I): MsgPruneBalancesResponse {
    const message = createBaseMsgPruneBalancesResponse();
    return message;
  },
};

function createBaseMsgUpdateBytes(): MsgUpdateBytes {
  return { creator: "", badgeId: 0, newBytes: "" };
}

export const MsgUpdateBytes = {
  encode(message: MsgUpdateBytes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.newBytes !== "") {
      writer.uint32(26).string(message.newBytes);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateBytes {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateBytes();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.newBytes = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateBytes {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      newBytes: isSet(object.newBytes) ? String(object.newBytes) : "",
    };
  },

  toJSON(message: MsgUpdateBytes): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.newBytes !== undefined && (obj.newBytes = message.newBytes);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateBytes>, I>>(object: I): MsgUpdateBytes {
    const message = createBaseMsgUpdateBytes();
    message.creator = object.creator ?? "";
    message.badgeId = object.badgeId ?? 0;
    message.newBytes = object.newBytes ?? "";
    return message;
  },
};

function createBaseMsgUpdateBytesResponse(): MsgUpdateBytesResponse {
  return {};
}

export const MsgUpdateBytesResponse = {
  encode(_: MsgUpdateBytesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateBytesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateBytesResponse();
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

  fromJSON(_: any): MsgUpdateBytesResponse {
    return {};
  },

  toJSON(_: MsgUpdateBytesResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateBytesResponse>, I>>(_: I): MsgUpdateBytesResponse {
    const message = createBaseMsgUpdateBytesResponse();
    return message;
  },
};

function createBaseMsgRegisterAddresses(): MsgRegisterAddresses {
  return { creator: "", addressesToRegister: [] };
}

export const MsgRegisterAddresses = {
  encode(message: MsgRegisterAddresses, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.addressesToRegister) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRegisterAddresses {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRegisterAddresses();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.addressesToRegister.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRegisterAddresses {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      addressesToRegister: Array.isArray(object?.addressesToRegister)
        ? object.addressesToRegister.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: MsgRegisterAddresses): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addressesToRegister) {
      obj.addressesToRegister = message.addressesToRegister.map((e) => e);
    } else {
      obj.addressesToRegister = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRegisterAddresses>, I>>(object: I): MsgRegisterAddresses {
    const message = createBaseMsgRegisterAddresses();
    message.creator = object.creator ?? "";
    message.addressesToRegister = object.addressesToRegister?.map((e) => e) || [];
    return message;
  },
};

function createBaseMsgRegisterAddressesResponse(): MsgRegisterAddressesResponse {
  return { registeredAddressNumbers: undefined };
}

export const MsgRegisterAddressesResponse = {
  encode(message: MsgRegisterAddressesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.registeredAddressNumbers !== undefined) {
      IdRange.encode(message.registeredAddressNumbers, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRegisterAddressesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRegisterAddressesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.registeredAddressNumbers = IdRange.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRegisterAddressesResponse {
    return {
      registeredAddressNumbers: isSet(object.registeredAddressNumbers)
        ? IdRange.fromJSON(object.registeredAddressNumbers)
        : undefined,
    };
  },

  toJSON(message: MsgRegisterAddressesResponse): unknown {
    const obj: any = {};
    message.registeredAddressNumbers !== undefined && (obj.registeredAddressNumbers = message.registeredAddressNumbers
      ? IdRange.toJSON(message.registeredAddressNumbers)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgRegisterAddressesResponse>, I>>(object: I): MsgRegisterAddressesResponse {
    const message = createBaseMsgRegisterAddressesResponse();
    message.registeredAddressNumbers =
      (object.registeredAddressNumbers !== undefined && object.registeredAddressNumbers !== null)
        ? IdRange.fromPartial(object.registeredAddressNumbers)
        : undefined;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  NewBadge(request: MsgNewBadge): Promise<MsgNewBadgeResponse>;
  NewSubBadge(request: MsgNewSubBadge): Promise<MsgNewSubBadgeResponse>;
  TransferBadge(request: MsgTransferBadge): Promise<MsgTransferBadgeResponse>;
  RequestTransferBadge(request: MsgRequestTransferBadge): Promise<MsgRequestTransferBadgeResponse>;
  HandlePendingTransfer(request: MsgHandlePendingTransfer): Promise<MsgHandlePendingTransferResponse>;
  SetApproval(request: MsgSetApproval): Promise<MsgSetApprovalResponse>;
  RevokeBadge(request: MsgRevokeBadge): Promise<MsgRevokeBadgeResponse>;
  FreezeAddress(request: MsgFreezeAddress): Promise<MsgFreezeAddressResponse>;
  UpdateUris(request: MsgUpdateUris): Promise<MsgUpdateUrisResponse>;
  UpdatePermissions(request: MsgUpdatePermissions): Promise<MsgUpdatePermissionsResponse>;
  TransferManager(request: MsgTransferManager): Promise<MsgTransferManagerResponse>;
  RequestTransferManager(request: MsgRequestTransferManager): Promise<MsgRequestTransferManagerResponse>;
  SelfDestructBadge(request: MsgSelfDestructBadge): Promise<MsgSelfDestructBadgeResponse>;
  PruneBalances(request: MsgPruneBalances): Promise<MsgPruneBalancesResponse>;
  UpdateBytes(request: MsgUpdateBytes): Promise<MsgUpdateBytesResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  RegisterAddresses(request: MsgRegisterAddresses): Promise<MsgRegisterAddressesResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.NewBadge = this.NewBadge.bind(this);
    this.NewSubBadge = this.NewSubBadge.bind(this);
    this.TransferBadge = this.TransferBadge.bind(this);
    this.RequestTransferBadge = this.RequestTransferBadge.bind(this);
    this.HandlePendingTransfer = this.HandlePendingTransfer.bind(this);
    this.SetApproval = this.SetApproval.bind(this);
    this.RevokeBadge = this.RevokeBadge.bind(this);
    this.FreezeAddress = this.FreezeAddress.bind(this);
    this.UpdateUris = this.UpdateUris.bind(this);
    this.UpdatePermissions = this.UpdatePermissions.bind(this);
    this.TransferManager = this.TransferManager.bind(this);
    this.RequestTransferManager = this.RequestTransferManager.bind(this);
    this.SelfDestructBadge = this.SelfDestructBadge.bind(this);
    this.PruneBalances = this.PruneBalances.bind(this);
    this.UpdateBytes = this.UpdateBytes.bind(this);
    this.RegisterAddresses = this.RegisterAddresses.bind(this);
  }
  NewBadge(request: MsgNewBadge): Promise<MsgNewBadgeResponse> {
    const data = MsgNewBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "NewBadge", data);
    return promise.then((data) => MsgNewBadgeResponse.decode(new _m0.Reader(data)));
  }

  NewSubBadge(request: MsgNewSubBadge): Promise<MsgNewSubBadgeResponse> {
    const data = MsgNewSubBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "NewSubBadge", data);
    return promise.then((data) => MsgNewSubBadgeResponse.decode(new _m0.Reader(data)));
  }

  TransferBadge(request: MsgTransferBadge): Promise<MsgTransferBadgeResponse> {
    const data = MsgTransferBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "TransferBadge", data);
    return promise.then((data) => MsgTransferBadgeResponse.decode(new _m0.Reader(data)));
  }

  RequestTransferBadge(request: MsgRequestTransferBadge): Promise<MsgRequestTransferBadgeResponse> {
    const data = MsgRequestTransferBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "RequestTransferBadge", data);
    return promise.then((data) => MsgRequestTransferBadgeResponse.decode(new _m0.Reader(data)));
  }

  HandlePendingTransfer(request: MsgHandlePendingTransfer): Promise<MsgHandlePendingTransferResponse> {
    const data = MsgHandlePendingTransfer.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "HandlePendingTransfer", data);
    return promise.then((data) => MsgHandlePendingTransferResponse.decode(new _m0.Reader(data)));
  }

  SetApproval(request: MsgSetApproval): Promise<MsgSetApprovalResponse> {
    const data = MsgSetApproval.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "SetApproval", data);
    return promise.then((data) => MsgSetApprovalResponse.decode(new _m0.Reader(data)));
  }

  RevokeBadge(request: MsgRevokeBadge): Promise<MsgRevokeBadgeResponse> {
    const data = MsgRevokeBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "RevokeBadge", data);
    return promise.then((data) => MsgRevokeBadgeResponse.decode(new _m0.Reader(data)));
  }

  FreezeAddress(request: MsgFreezeAddress): Promise<MsgFreezeAddressResponse> {
    const data = MsgFreezeAddress.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "FreezeAddress", data);
    return promise.then((data) => MsgFreezeAddressResponse.decode(new _m0.Reader(data)));
  }

  UpdateUris(request: MsgUpdateUris): Promise<MsgUpdateUrisResponse> {
    const data = MsgUpdateUris.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateUris", data);
    return promise.then((data) => MsgUpdateUrisResponse.decode(new _m0.Reader(data)));
  }

  UpdatePermissions(request: MsgUpdatePermissions): Promise<MsgUpdatePermissionsResponse> {
    const data = MsgUpdatePermissions.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdatePermissions", data);
    return promise.then((data) => MsgUpdatePermissionsResponse.decode(new _m0.Reader(data)));
  }

  TransferManager(request: MsgTransferManager): Promise<MsgTransferManagerResponse> {
    const data = MsgTransferManager.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "TransferManager", data);
    return promise.then((data) => MsgTransferManagerResponse.decode(new _m0.Reader(data)));
  }

  RequestTransferManager(request: MsgRequestTransferManager): Promise<MsgRequestTransferManagerResponse> {
    const data = MsgRequestTransferManager.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "RequestTransferManager", data);
    return promise.then((data) => MsgRequestTransferManagerResponse.decode(new _m0.Reader(data)));
  }

  SelfDestructBadge(request: MsgSelfDestructBadge): Promise<MsgSelfDestructBadgeResponse> {
    const data = MsgSelfDestructBadge.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "SelfDestructBadge", data);
    return promise.then((data) => MsgSelfDestructBadgeResponse.decode(new _m0.Reader(data)));
  }

  PruneBalances(request: MsgPruneBalances): Promise<MsgPruneBalancesResponse> {
    const data = MsgPruneBalances.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "PruneBalances", data);
    return promise.then((data) => MsgPruneBalancesResponse.decode(new _m0.Reader(data)));
  }

  UpdateBytes(request: MsgUpdateBytes): Promise<MsgUpdateBytesResponse> {
    const data = MsgUpdateBytes.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "UpdateBytes", data);
    return promise.then((data) => MsgUpdateBytesResponse.decode(new _m0.Reader(data)));
  }

  RegisterAddresses(request: MsgRegisterAddresses): Promise<MsgRegisterAddressesResponse> {
    const data = MsgRegisterAddresses.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Msg", "RegisterAddresses", data);
    return promise.then((data) => MsgRegisterAddressesResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
