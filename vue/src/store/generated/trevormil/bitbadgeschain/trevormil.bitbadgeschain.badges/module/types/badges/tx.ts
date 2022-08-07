/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { NumberRange } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

export interface MsgNewBadge {
  creator: string;
  uri: string;
  subassetUris: string;
  permissions: number;
  metadataHash: string;
  defaultSubassetSupply: number;
}

export interface MsgNewBadgeResponse {
  id: number;
}

export interface MsgNewSubBadge {
  creator: string;
  id: number;
  supplys: number[];
  amountsToCreate: number[];
}

export interface MsgNewSubBadgeResponse {
  subassetId: number;
}

export interface MsgTransferBadge {
  creator: string;
  from: number;
  toAddresses: number[];
  amounts: number[];
  badgeId: number;
  subbadgeRanges: NumberRange[];
}

export interface MsgTransferBadgeResponse {}

export interface MsgRequestTransferBadge {
  creator: string;
  from: number;
  amount: number;
  badgeId: number;
  subbadgeRanges: NumberRange[];
}

export interface MsgRequestTransferBadgeResponse {}

export interface MsgHandlePendingTransfer {
  creator: string;
  accept: boolean;
  badgeId: number;
  nonceRanges: NumberRange[];
  forcefulAccept: boolean;
}

export interface MsgHandlePendingTransferResponse {}

export interface MsgSetApproval {
  creator: string;
  amount: number;
  address: number;
  badgeId: number;
  subbadgeRanges: NumberRange[];
}

export interface MsgSetApprovalResponse {}

export interface MsgRevokeBadge {
  creator: string;
  addresses: number[];
  amounts: number[];
  badgeId: number;
  subbadgeRanges: NumberRange[];
}

export interface MsgRevokeBadgeResponse {}

export interface MsgFreezeAddress {
  creator: string;
  addressRanges: NumberRange[];
  badgeId: number;
  add: boolean;
}

export interface MsgFreezeAddressResponse {}

export interface MsgUpdateUris {
  creator: string;
  badgeId: number;
  uri: string;
  subassetUri: string;
}

export interface MsgUpdateUrisResponse {}

export interface MsgUpdatePermissions {
  creator: string;
  badgeId: number;
  permissions: number;
}

export interface MsgUpdatePermissionsResponse {}

export interface MsgTransferManager {
  creator: string;
  badgeId: number;
  address: number;
}

export interface MsgTransferManagerResponse {}

export interface MsgRequestTransferManager {
  creator: string;
  badgeId: number;
  add: boolean;
}

export interface MsgRequestTransferManagerResponse {}

export interface MsgSelfDestructBadge {
  creator: string;
  badgeId: number;
}

export interface MsgSelfDestructBadgeResponse {}

export interface MsgPruneBalances {
  creator: string;
  badgeIds: number;
  addresses: number;
}

export interface MsgPruneBalancesResponse {}

const baseMsgNewBadge: object = {
  creator: "",
  uri: "",
  subassetUris: "",
  permissions: 0,
  metadataHash: "",
  defaultSubassetSupply: 0,
};

export const MsgNewBadge = {
  encode(message: MsgNewBadge, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.uri !== "") {
      writer.uint32(18).string(message.uri);
    }
    if (message.subassetUris !== "") {
      writer.uint32(26).string(message.subassetUris);
    }
    if (message.permissions !== 0) {
      writer.uint32(32).uint64(message.permissions);
    }
    if (message.metadataHash !== "") {
      writer.uint32(42).string(message.metadataHash);
    }
    if (message.defaultSubassetSupply !== 0) {
      writer.uint32(48).uint64(message.defaultSubassetSupply);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgNewBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgNewBadge } as MsgNewBadge;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.uri = reader.string();
          break;
        case 3:
          message.subassetUris = reader.string();
          break;
        case 4:
          message.permissions = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.metadataHash = reader.string();
          break;
        case 6:
          message.defaultSubassetSupply = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewBadge {
    const message = { ...baseMsgNewBadge } as MsgNewBadge;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = String(object.uri);
    } else {
      message.uri = "";
    }
    if (object.subassetUris !== undefined && object.subassetUris !== null) {
      message.subassetUris = String(object.subassetUris);
    } else {
      message.subassetUris = "";
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = Number(object.permissions);
    } else {
      message.permissions = 0;
    }
    if (object.metadataHash !== undefined && object.metadataHash !== null) {
      message.metadataHash = String(object.metadataHash);
    } else {
      message.metadataHash = "";
    }
    if (
      object.defaultSubassetSupply !== undefined &&
      object.defaultSubassetSupply !== null
    ) {
      message.defaultSubassetSupply = Number(object.defaultSubassetSupply);
    } else {
      message.defaultSubassetSupply = 0;
    }
    return message;
  },

  toJSON(message: MsgNewBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.uri !== undefined && (obj.uri = message.uri);
    message.subassetUris !== undefined &&
      (obj.subassetUris = message.subassetUris);
    message.permissions !== undefined &&
      (obj.permissions = message.permissions);
    message.metadataHash !== undefined &&
      (obj.metadataHash = message.metadataHash);
    message.defaultSubassetSupply !== undefined &&
      (obj.defaultSubassetSupply = message.defaultSubassetSupply);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgNewBadge>): MsgNewBadge {
    const message = { ...baseMsgNewBadge } as MsgNewBadge;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = object.uri;
    } else {
      message.uri = "";
    }
    if (object.subassetUris !== undefined && object.subassetUris !== null) {
      message.subassetUris = object.subassetUris;
    } else {
      message.subassetUris = "";
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = object.permissions;
    } else {
      message.permissions = 0;
    }
    if (object.metadataHash !== undefined && object.metadataHash !== null) {
      message.metadataHash = object.metadataHash;
    } else {
      message.metadataHash = "";
    }
    if (
      object.defaultSubassetSupply !== undefined &&
      object.defaultSubassetSupply !== null
    ) {
      message.defaultSubassetSupply = object.defaultSubassetSupply;
    } else {
      message.defaultSubassetSupply = 0;
    }
    return message;
  },
};

const baseMsgNewBadgeResponse: object = { id: 0 };

export const MsgNewBadgeResponse = {
  encode(
    message: MsgNewBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgNewBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgNewBadgeResponse } as MsgNewBadgeResponse;
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
    const message = { ...baseMsgNewBadgeResponse } as MsgNewBadgeResponse;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    return message;
  },

  toJSON(message: MsgNewBadgeResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgNewBadgeResponse>): MsgNewBadgeResponse {
    const message = { ...baseMsgNewBadgeResponse } as MsgNewBadgeResponse;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    return message;
  },
};

const baseMsgNewSubBadge: object = {
  creator: "",
  id: 0,
  supplys: 0,
  amountsToCreate: 0,
};

export const MsgNewSubBadge = {
  encode(message: MsgNewSubBadge, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
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

  decode(input: Reader | Uint8Array, length?: number): MsgNewSubBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgNewSubBadge } as MsgNewSubBadge;
    message.supplys = [];
    message.amountsToCreate = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
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
              message.amountsToCreate.push(
                longToNumber(reader.uint64() as Long)
              );
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
    const message = { ...baseMsgNewSubBadge } as MsgNewSubBadge;
    message.supplys = [];
    message.amountsToCreate = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    if (object.supplys !== undefined && object.supplys !== null) {
      for (const e of object.supplys) {
        message.supplys.push(Number(e));
      }
    }
    if (
      object.amountsToCreate !== undefined &&
      object.amountsToCreate !== null
    ) {
      for (const e of object.amountsToCreate) {
        message.amountsToCreate.push(Number(e));
      }
    }
    return message;
  },

  toJSON(message: MsgNewSubBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = message.id);
    if (message.supplys) {
      obj.supplys = message.supplys.map((e) => e);
    } else {
      obj.supplys = [];
    }
    if (message.amountsToCreate) {
      obj.amountsToCreate = message.amountsToCreate.map((e) => e);
    } else {
      obj.amountsToCreate = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<MsgNewSubBadge>): MsgNewSubBadge {
    const message = { ...baseMsgNewSubBadge } as MsgNewSubBadge;
    message.supplys = [];
    message.amountsToCreate = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    if (object.supplys !== undefined && object.supplys !== null) {
      for (const e of object.supplys) {
        message.supplys.push(e);
      }
    }
    if (
      object.amountsToCreate !== undefined &&
      object.amountsToCreate !== null
    ) {
      for (const e of object.amountsToCreate) {
        message.amountsToCreate.push(e);
      }
    }
    return message;
  },
};

const baseMsgNewSubBadgeResponse: object = { subassetId: 0 };

export const MsgNewSubBadgeResponse = {
  encode(
    message: MsgNewSubBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.subassetId !== 0) {
      writer.uint32(8).uint64(message.subassetId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgNewSubBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgNewSubBadgeResponse } as MsgNewSubBadgeResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subassetId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgNewSubBadgeResponse {
    const message = { ...baseMsgNewSubBadgeResponse } as MsgNewSubBadgeResponse;
    if (object.subassetId !== undefined && object.subassetId !== null) {
      message.subassetId = Number(object.subassetId);
    } else {
      message.subassetId = 0;
    }
    return message;
  },

  toJSON(message: MsgNewSubBadgeResponse): unknown {
    const obj: any = {};
    message.subassetId !== undefined && (obj.subassetId = message.subassetId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgNewSubBadgeResponse>
  ): MsgNewSubBadgeResponse {
    const message = { ...baseMsgNewSubBadgeResponse } as MsgNewSubBadgeResponse;
    if (object.subassetId !== undefined && object.subassetId !== null) {
      message.subassetId = object.subassetId;
    } else {
      message.subassetId = 0;
    }
    return message;
  },
};

const baseMsgTransferBadge: object = {
  creator: "",
  from: 0,
  toAddresses: 0,
  amounts: 0,
  badgeId: 0,
};

export const MsgTransferBadge = {
  encode(message: MsgTransferBadge, writer: Writer = Writer.create()): Writer {
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
      NumberRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgTransferBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgTransferBadge } as MsgTransferBadge;
    message.toAddresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
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
          message.subbadgeRanges.push(
            NumberRange.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgTransferBadge {
    const message = { ...baseMsgTransferBadge } as MsgTransferBadge;
    message.toAddresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = Number(object.from);
    } else {
      message.from = 0;
    }
    if (object.toAddresses !== undefined && object.toAddresses !== null) {
      for (const e of object.toAddresses) {
        message.toAddresses.push(Number(e));
      }
    }
    if (object.amounts !== undefined && object.amounts !== null) {
      for (const e of object.amounts) {
        message.amounts.push(Number(e));
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: MsgTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = message.from);
    if (message.toAddresses) {
      obj.toAddresses = message.toAddresses.map((e) => e);
    } else {
      obj.toAddresses = [];
    }
    if (message.amounts) {
      obj.amounts = message.amounts.map((e) => e);
    } else {
      obj.amounts = [];
    }
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<MsgTransferBadge>): MsgTransferBadge {
    const message = { ...baseMsgTransferBadge } as MsgTransferBadge;
    message.toAddresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = object.from;
    } else {
      message.from = 0;
    }
    if (object.toAddresses !== undefined && object.toAddresses !== null) {
      for (const e of object.toAddresses) {
        message.toAddresses.push(e);
      }
    }
    if (object.amounts !== undefined && object.amounts !== null) {
      for (const e of object.amounts) {
        message.amounts.push(e);
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromPartial(e));
      }
    }
    return message;
  },
};

const baseMsgTransferBadgeResponse: object = {};

export const MsgTransferBadgeResponse = {
  encode(
    _: MsgTransferBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgTransferBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgTransferBadgeResponse,
    } as MsgTransferBadgeResponse;
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
    const message = {
      ...baseMsgTransferBadgeResponse,
    } as MsgTransferBadgeResponse;
    return message;
  },

  toJSON(_: MsgTransferBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgTransferBadgeResponse>
  ): MsgTransferBadgeResponse {
    const message = {
      ...baseMsgTransferBadgeResponse,
    } as MsgTransferBadgeResponse;
    return message;
  },
};

const baseMsgRequestTransferBadge: object = {
  creator: "",
  from: 0,
  amount: 0,
  badgeId: 0,
};

export const MsgRequestTransferBadge = {
  encode(
    message: MsgRequestTransferBadge,
    writer: Writer = Writer.create()
  ): Writer {
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
      NumberRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRequestTransferBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestTransferBadge,
    } as MsgRequestTransferBadge;
    message.subbadgeRanges = [];
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
          message.subbadgeRanges.push(
            NumberRange.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRequestTransferBadge {
    const message = {
      ...baseMsgRequestTransferBadge,
    } as MsgRequestTransferBadge;
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = Number(object.from);
    } else {
      message.from = 0;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: MsgRequestTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = message.from);
    message.amount !== undefined && (obj.amount = message.amount);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgRequestTransferBadge>
  ): MsgRequestTransferBadge {
    const message = {
      ...baseMsgRequestTransferBadge,
    } as MsgRequestTransferBadge;
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = object.from;
    } else {
      message.from = 0;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromPartial(e));
      }
    }
    return message;
  },
};

const baseMsgRequestTransferBadgeResponse: object = {};

export const MsgRequestTransferBadgeResponse = {
  encode(
    _: MsgRequestTransferBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgRequestTransferBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestTransferBadgeResponse,
    } as MsgRequestTransferBadgeResponse;
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
    const message = {
      ...baseMsgRequestTransferBadgeResponse,
    } as MsgRequestTransferBadgeResponse;
    return message;
  },

  toJSON(_: MsgRequestTransferBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgRequestTransferBadgeResponse>
  ): MsgRequestTransferBadgeResponse {
    const message = {
      ...baseMsgRequestTransferBadgeResponse,
    } as MsgRequestTransferBadgeResponse;
    return message;
  },
};

const baseMsgHandlePendingTransfer: object = {
  creator: "",
  accept: false,
  badgeId: 0,
  forcefulAccept: false,
};

export const MsgHandlePendingTransfer = {
  encode(
    message: MsgHandlePendingTransfer,
    writer: Writer = Writer.create()
  ): Writer {
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
      NumberRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.forcefulAccept === true) {
      writer.uint32(40).bool(message.forcefulAccept);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgHandlePendingTransfer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgHandlePendingTransfer,
    } as MsgHandlePendingTransfer;
    message.nonceRanges = [];
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
          message.nonceRanges.push(NumberRange.decode(reader, reader.uint32()));
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
    const message = {
      ...baseMsgHandlePendingTransfer,
    } as MsgHandlePendingTransfer;
    message.nonceRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.accept !== undefined && object.accept !== null) {
      message.accept = Boolean(object.accept);
    } else {
      message.accept = false;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.nonceRanges !== undefined && object.nonceRanges !== null) {
      for (const e of object.nonceRanges) {
        message.nonceRanges.push(NumberRange.fromJSON(e));
      }
    }
    if (object.forcefulAccept !== undefined && object.forcefulAccept !== null) {
      message.forcefulAccept = Boolean(object.forcefulAccept);
    } else {
      message.forcefulAccept = false;
    }
    return message;
  },

  toJSON(message: MsgHandlePendingTransfer): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.accept !== undefined && (obj.accept = message.accept);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    if (message.nonceRanges) {
      obj.nonceRanges = message.nonceRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.nonceRanges = [];
    }
    message.forcefulAccept !== undefined &&
      (obj.forcefulAccept = message.forcefulAccept);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgHandlePendingTransfer>
  ): MsgHandlePendingTransfer {
    const message = {
      ...baseMsgHandlePendingTransfer,
    } as MsgHandlePendingTransfer;
    message.nonceRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.accept !== undefined && object.accept !== null) {
      message.accept = object.accept;
    } else {
      message.accept = false;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.nonceRanges !== undefined && object.nonceRanges !== null) {
      for (const e of object.nonceRanges) {
        message.nonceRanges.push(NumberRange.fromPartial(e));
      }
    }
    if (object.forcefulAccept !== undefined && object.forcefulAccept !== null) {
      message.forcefulAccept = object.forcefulAccept;
    } else {
      message.forcefulAccept = false;
    }
    return message;
  },
};

const baseMsgHandlePendingTransferResponse: object = {};

export const MsgHandlePendingTransferResponse = {
  encode(
    _: MsgHandlePendingTransferResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgHandlePendingTransferResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgHandlePendingTransferResponse,
    } as MsgHandlePendingTransferResponse;
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
    const message = {
      ...baseMsgHandlePendingTransferResponse,
    } as MsgHandlePendingTransferResponse;
    return message;
  },

  toJSON(_: MsgHandlePendingTransferResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgHandlePendingTransferResponse>
  ): MsgHandlePendingTransferResponse {
    const message = {
      ...baseMsgHandlePendingTransferResponse,
    } as MsgHandlePendingTransferResponse;
    return message;
  },
};

const baseMsgSetApproval: object = {
  creator: "",
  amount: 0,
  address: 0,
  badgeId: 0,
};

export const MsgSetApproval = {
  encode(message: MsgSetApproval, writer: Writer = Writer.create()): Writer {
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
      NumberRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgSetApproval {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSetApproval } as MsgSetApproval;
    message.subbadgeRanges = [];
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
          message.subbadgeRanges.push(
            NumberRange.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgSetApproval {
    const message = { ...baseMsgSetApproval } as MsgSetApproval;
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = Number(object.address);
    } else {
      message.address = 0;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: MsgSetApproval): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.amount !== undefined && (obj.amount = message.amount);
    message.address !== undefined && (obj.address = message.address);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<MsgSetApproval>): MsgSetApproval {
    const message = { ...baseMsgSetApproval } as MsgSetApproval;
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = 0;
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromPartial(e));
      }
    }
    return message;
  },
};

const baseMsgSetApprovalResponse: object = {};

export const MsgSetApprovalResponse = {
  encode(_: MsgSetApprovalResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgSetApprovalResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSetApprovalResponse } as MsgSetApprovalResponse;
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
    const message = { ...baseMsgSetApprovalResponse } as MsgSetApprovalResponse;
    return message;
  },

  toJSON(_: MsgSetApprovalResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgSetApprovalResponse>): MsgSetApprovalResponse {
    const message = { ...baseMsgSetApprovalResponse } as MsgSetApprovalResponse;
    return message;
  },
};

const baseMsgRevokeBadge: object = {
  creator: "",
  addresses: 0,
  amounts: 0,
  badgeId: 0,
};

export const MsgRevokeBadge = {
  encode(message: MsgRevokeBadge, writer: Writer = Writer.create()): Writer {
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
      NumberRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRevokeBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgRevokeBadge } as MsgRevokeBadge;
    message.addresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
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
          message.subbadgeRanges.push(
            NumberRange.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRevokeBadge {
    const message = { ...baseMsgRevokeBadge } as MsgRevokeBadge;
    message.addresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      for (const e of object.addresses) {
        message.addresses.push(Number(e));
      }
    }
    if (object.amounts !== undefined && object.amounts !== null) {
      for (const e of object.amounts) {
        message.amounts.push(Number(e));
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: MsgRevokeBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => e);
    } else {
      obj.addresses = [];
    }
    if (message.amounts) {
      obj.amounts = message.amounts.map((e) => e);
    } else {
      obj.amounts = [];
    }
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    if (message.subbadgeRanges) {
      obj.subbadgeRanges = message.subbadgeRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.subbadgeRanges = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<MsgRevokeBadge>): MsgRevokeBadge {
    const message = { ...baseMsgRevokeBadge } as MsgRevokeBadge;
    message.addresses = [];
    message.amounts = [];
    message.subbadgeRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      for (const e of object.addresses) {
        message.addresses.push(e);
      }
    }
    if (object.amounts !== undefined && object.amounts !== null) {
      for (const e of object.amounts) {
        message.amounts.push(e);
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeRanges !== undefined && object.subbadgeRanges !== null) {
      for (const e of object.subbadgeRanges) {
        message.subbadgeRanges.push(NumberRange.fromPartial(e));
      }
    }
    return message;
  },
};

const baseMsgRevokeBadgeResponse: object = {};

export const MsgRevokeBadgeResponse = {
  encode(_: MsgRevokeBadgeResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRevokeBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgRevokeBadgeResponse } as MsgRevokeBadgeResponse;
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
    const message = { ...baseMsgRevokeBadgeResponse } as MsgRevokeBadgeResponse;
    return message;
  },

  toJSON(_: MsgRevokeBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgRevokeBadgeResponse>): MsgRevokeBadgeResponse {
    const message = { ...baseMsgRevokeBadgeResponse } as MsgRevokeBadgeResponse;
    return message;
  },
};

const baseMsgFreezeAddress: object = { creator: "", badgeId: 0, add: false };

export const MsgFreezeAddress = {
  encode(message: MsgFreezeAddress, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.addressRanges) {
      NumberRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.badgeId !== 0) {
      writer.uint32(24).uint64(message.badgeId);
    }
    if (message.add === true) {
      writer.uint32(32).bool(message.add);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgFreezeAddress {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgFreezeAddress } as MsgFreezeAddress;
    message.addressRanges = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.addressRanges.push(
            NumberRange.decode(reader, reader.uint32())
          );
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
    const message = { ...baseMsgFreezeAddress } as MsgFreezeAddress;
    message.addressRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.addressRanges !== undefined && object.addressRanges !== null) {
      for (const e of object.addressRanges) {
        message.addressRanges.push(NumberRange.fromJSON(e));
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.add !== undefined && object.add !== null) {
      message.add = Boolean(object.add);
    } else {
      message.add = false;
    }
    return message;
  },

  toJSON(message: MsgFreezeAddress): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.addressRanges) {
      obj.addressRanges = message.addressRanges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.addressRanges = [];
    }
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.add !== undefined && (obj.add = message.add);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgFreezeAddress>): MsgFreezeAddress {
    const message = { ...baseMsgFreezeAddress } as MsgFreezeAddress;
    message.addressRanges = [];
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.addressRanges !== undefined && object.addressRanges !== null) {
      for (const e of object.addressRanges) {
        message.addressRanges.push(NumberRange.fromPartial(e));
      }
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.add !== undefined && object.add !== null) {
      message.add = object.add;
    } else {
      message.add = false;
    }
    return message;
  },
};

const baseMsgFreezeAddressResponse: object = {};

export const MsgFreezeAddressResponse = {
  encode(
    _: MsgFreezeAddressResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgFreezeAddressResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgFreezeAddressResponse,
    } as MsgFreezeAddressResponse;
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
    const message = {
      ...baseMsgFreezeAddressResponse,
    } as MsgFreezeAddressResponse;
    return message;
  },

  toJSON(_: MsgFreezeAddressResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgFreezeAddressResponse>
  ): MsgFreezeAddressResponse {
    const message = {
      ...baseMsgFreezeAddressResponse,
    } as MsgFreezeAddressResponse;
    return message;
  },
};

const baseMsgUpdateUris: object = {
  creator: "",
  badgeId: 0,
  uri: "",
  subassetUri: "",
};

export const MsgUpdateUris = {
  encode(message: MsgUpdateUris, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.subassetUri !== "") {
      writer.uint32(34).string(message.subassetUri);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgUpdateUris {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgUpdateUris } as MsgUpdateUris;
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
          message.uri = reader.string();
          break;
        case 4:
          message.subassetUri = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateUris {
    const message = { ...baseMsgUpdateUris } as MsgUpdateUris;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = String(object.uri);
    } else {
      message.uri = "";
    }
    if (object.subassetUri !== undefined && object.subassetUri !== null) {
      message.subassetUri = String(object.subassetUri);
    } else {
      message.subassetUri = "";
    }
    return message;
  },

  toJSON(message: MsgUpdateUris): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.subassetUri !== undefined &&
      (obj.subassetUri = message.subassetUri);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgUpdateUris>): MsgUpdateUris {
    const message = { ...baseMsgUpdateUris } as MsgUpdateUris;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = object.uri;
    } else {
      message.uri = "";
    }
    if (object.subassetUri !== undefined && object.subassetUri !== null) {
      message.subassetUri = object.subassetUri;
    } else {
      message.subassetUri = "";
    }
    return message;
  },
};

const baseMsgUpdateUrisResponse: object = {};

export const MsgUpdateUrisResponse = {
  encode(_: MsgUpdateUrisResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgUpdateUrisResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgUpdateUrisResponse } as MsgUpdateUrisResponse;
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
    const message = { ...baseMsgUpdateUrisResponse } as MsgUpdateUrisResponse;
    return message;
  },

  toJSON(_: MsgUpdateUrisResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgUpdateUrisResponse>): MsgUpdateUrisResponse {
    const message = { ...baseMsgUpdateUrisResponse } as MsgUpdateUrisResponse;
    return message;
  },
};

const baseMsgUpdatePermissions: object = {
  creator: "",
  badgeId: 0,
  permissions: 0,
};

export const MsgUpdatePermissions = {
  encode(
    message: MsgUpdatePermissions,
    writer: Writer = Writer.create()
  ): Writer {
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

  decode(input: Reader | Uint8Array, length?: number): MsgUpdatePermissions {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgUpdatePermissions } as MsgUpdatePermissions;
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
    const message = { ...baseMsgUpdatePermissions } as MsgUpdatePermissions;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = Number(object.permissions);
    } else {
      message.permissions = 0;
    }
    return message;
  },

  toJSON(message: MsgUpdatePermissions): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.permissions !== undefined &&
      (obj.permissions = message.permissions);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgUpdatePermissions>): MsgUpdatePermissions {
    const message = { ...baseMsgUpdatePermissions } as MsgUpdatePermissions;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = object.permissions;
    } else {
      message.permissions = 0;
    }
    return message;
  },
};

const baseMsgUpdatePermissionsResponse: object = {};

export const MsgUpdatePermissionsResponse = {
  encode(
    _: MsgUpdatePermissionsResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgUpdatePermissionsResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgUpdatePermissionsResponse,
    } as MsgUpdatePermissionsResponse;
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
    const message = {
      ...baseMsgUpdatePermissionsResponse,
    } as MsgUpdatePermissionsResponse;
    return message;
  },

  toJSON(_: MsgUpdatePermissionsResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgUpdatePermissionsResponse>
  ): MsgUpdatePermissionsResponse {
    const message = {
      ...baseMsgUpdatePermissionsResponse,
    } as MsgUpdatePermissionsResponse;
    return message;
  },
};

const baseMsgTransferManager: object = { creator: "", badgeId: 0, address: 0 };

export const MsgTransferManager = {
  encode(
    message: MsgTransferManager,
    writer: Writer = Writer.create()
  ): Writer {
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

  decode(input: Reader | Uint8Array, length?: number): MsgTransferManager {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgTransferManager } as MsgTransferManager;
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
    const message = { ...baseMsgTransferManager } as MsgTransferManager;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = Number(object.address);
    } else {
      message.address = 0;
    }
    return message;
  },

  toJSON(message: MsgTransferManager): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgTransferManager>): MsgTransferManager {
    const message = { ...baseMsgTransferManager } as MsgTransferManager;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = 0;
    }
    return message;
  },
};

const baseMsgTransferManagerResponse: object = {};

export const MsgTransferManagerResponse = {
  encode(
    _: MsgTransferManagerResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgTransferManagerResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgTransferManagerResponse,
    } as MsgTransferManagerResponse;
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
    const message = {
      ...baseMsgTransferManagerResponse,
    } as MsgTransferManagerResponse;
    return message;
  },

  toJSON(_: MsgTransferManagerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgTransferManagerResponse>
  ): MsgTransferManagerResponse {
    const message = {
      ...baseMsgTransferManagerResponse,
    } as MsgTransferManagerResponse;
    return message;
  },
};

const baseMsgRequestTransferManager: object = {
  creator: "",
  badgeId: 0,
  add: false,
};

export const MsgRequestTransferManager = {
  encode(
    message: MsgRequestTransferManager,
    writer: Writer = Writer.create()
  ): Writer {
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

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgRequestTransferManager {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestTransferManager,
    } as MsgRequestTransferManager;
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
    const message = {
      ...baseMsgRequestTransferManager,
    } as MsgRequestTransferManager;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.add !== undefined && object.add !== null) {
      message.add = Boolean(object.add);
    } else {
      message.add = false;
    }
    return message;
  },

  toJSON(message: MsgRequestTransferManager): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.add !== undefined && (obj.add = message.add);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgRequestTransferManager>
  ): MsgRequestTransferManager {
    const message = {
      ...baseMsgRequestTransferManager,
    } as MsgRequestTransferManager;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.add !== undefined && object.add !== null) {
      message.add = object.add;
    } else {
      message.add = false;
    }
    return message;
  },
};

const baseMsgRequestTransferManagerResponse: object = {};

export const MsgRequestTransferManagerResponse = {
  encode(
    _: MsgRequestTransferManagerResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgRequestTransferManagerResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestTransferManagerResponse,
    } as MsgRequestTransferManagerResponse;
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
    const message = {
      ...baseMsgRequestTransferManagerResponse,
    } as MsgRequestTransferManagerResponse;
    return message;
  },

  toJSON(_: MsgRequestTransferManagerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgRequestTransferManagerResponse>
  ): MsgRequestTransferManagerResponse {
    const message = {
      ...baseMsgRequestTransferManagerResponse,
    } as MsgRequestTransferManagerResponse;
    return message;
  },
};

const baseMsgSelfDestructBadge: object = { creator: "", badgeId: 0 };

export const MsgSelfDestructBadge = {
  encode(
    message: MsgSelfDestructBadge,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeId !== 0) {
      writer.uint32(16).uint64(message.badgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgSelfDestructBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSelfDestructBadge } as MsgSelfDestructBadge;
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
    const message = { ...baseMsgSelfDestructBadge } as MsgSelfDestructBadge;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    return message;
  },

  toJSON(message: MsgSelfDestructBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgSelfDestructBadge>): MsgSelfDestructBadge {
    const message = { ...baseMsgSelfDestructBadge } as MsgSelfDestructBadge;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    return message;
  },
};

const baseMsgSelfDestructBadgeResponse: object = {};

export const MsgSelfDestructBadgeResponse = {
  encode(
    _: MsgSelfDestructBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgSelfDestructBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgSelfDestructBadgeResponse,
    } as MsgSelfDestructBadgeResponse;
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
    const message = {
      ...baseMsgSelfDestructBadgeResponse,
    } as MsgSelfDestructBadgeResponse;
    return message;
  },

  toJSON(_: MsgSelfDestructBadgeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgSelfDestructBadgeResponse>
  ): MsgSelfDestructBadgeResponse {
    const message = {
      ...baseMsgSelfDestructBadgeResponse,
    } as MsgSelfDestructBadgeResponse;
    return message;
  },
};

const baseMsgPruneBalances: object = { creator: "", badgeIds: 0, addresses: 0 };

export const MsgPruneBalances = {
  encode(message: MsgPruneBalances, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.badgeIds !== 0) {
      writer.uint32(16).uint64(message.badgeIds);
    }
    if (message.addresses !== 0) {
      writer.uint32(24).uint64(message.addresses);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgPruneBalances {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgPruneBalances } as MsgPruneBalances;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.badgeIds = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.addresses = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgPruneBalances {
    const message = { ...baseMsgPruneBalances } as MsgPruneBalances;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.badgeIds !== undefined && object.badgeIds !== null) {
      message.badgeIds = Number(object.badgeIds);
    } else {
      message.badgeIds = 0;
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      message.addresses = Number(object.addresses);
    } else {
      message.addresses = 0;
    }
    return message;
  },

  toJSON(message: MsgPruneBalances): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.badgeIds !== undefined && (obj.badgeIds = message.badgeIds);
    message.addresses !== undefined && (obj.addresses = message.addresses);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgPruneBalances>): MsgPruneBalances {
    const message = { ...baseMsgPruneBalances } as MsgPruneBalances;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.badgeIds !== undefined && object.badgeIds !== null) {
      message.badgeIds = object.badgeIds;
    } else {
      message.badgeIds = 0;
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      message.addresses = object.addresses;
    } else {
      message.addresses = 0;
    }
    return message;
  },
};

const baseMsgPruneBalancesResponse: object = {};

export const MsgPruneBalancesResponse = {
  encode(
    _: MsgPruneBalancesResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgPruneBalancesResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgPruneBalancesResponse,
    } as MsgPruneBalancesResponse;
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
    const message = {
      ...baseMsgPruneBalancesResponse,
    } as MsgPruneBalancesResponse;
    return message;
  },

  toJSON(_: MsgPruneBalancesResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgPruneBalancesResponse>
  ): MsgPruneBalancesResponse {
    const message = {
      ...baseMsgPruneBalancesResponse,
    } as MsgPruneBalancesResponse;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  NewBadge(request: MsgNewBadge): Promise<MsgNewBadgeResponse>;
  NewSubBadge(request: MsgNewSubBadge): Promise<MsgNewSubBadgeResponse>;
  TransferBadge(request: MsgTransferBadge): Promise<MsgTransferBadgeResponse>;
  RequestTransferBadge(
    request: MsgRequestTransferBadge
  ): Promise<MsgRequestTransferBadgeResponse>;
  HandlePendingTransfer(
    request: MsgHandlePendingTransfer
  ): Promise<MsgHandlePendingTransferResponse>;
  SetApproval(request: MsgSetApproval): Promise<MsgSetApprovalResponse>;
  RevokeBadge(request: MsgRevokeBadge): Promise<MsgRevokeBadgeResponse>;
  FreezeAddress(request: MsgFreezeAddress): Promise<MsgFreezeAddressResponse>;
  UpdateUris(request: MsgUpdateUris): Promise<MsgUpdateUrisResponse>;
  UpdatePermissions(
    request: MsgUpdatePermissions
  ): Promise<MsgUpdatePermissionsResponse>;
  TransferManager(
    request: MsgTransferManager
  ): Promise<MsgTransferManagerResponse>;
  RequestTransferManager(
    request: MsgRequestTransferManager
  ): Promise<MsgRequestTransferManagerResponse>;
  SelfDestructBadge(
    request: MsgSelfDestructBadge
  ): Promise<MsgSelfDestructBadgeResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  PruneBalances(request: MsgPruneBalances): Promise<MsgPruneBalancesResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  NewBadge(request: MsgNewBadge): Promise<MsgNewBadgeResponse> {
    const data = MsgNewBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "NewBadge",
      data
    );
    return promise.then((data) => MsgNewBadgeResponse.decode(new Reader(data)));
  }

  NewSubBadge(request: MsgNewSubBadge): Promise<MsgNewSubBadgeResponse> {
    const data = MsgNewSubBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "NewSubBadge",
      data
    );
    return promise.then((data) =>
      MsgNewSubBadgeResponse.decode(new Reader(data))
    );
  }

  TransferBadge(request: MsgTransferBadge): Promise<MsgTransferBadgeResponse> {
    const data = MsgTransferBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "TransferBadge",
      data
    );
    return promise.then((data) =>
      MsgTransferBadgeResponse.decode(new Reader(data))
    );
  }

  RequestTransferBadge(
    request: MsgRequestTransferBadge
  ): Promise<MsgRequestTransferBadgeResponse> {
    const data = MsgRequestTransferBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "RequestTransferBadge",
      data
    );
    return promise.then((data) =>
      MsgRequestTransferBadgeResponse.decode(new Reader(data))
    );
  }

  HandlePendingTransfer(
    request: MsgHandlePendingTransfer
  ): Promise<MsgHandlePendingTransferResponse> {
    const data = MsgHandlePendingTransfer.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "HandlePendingTransfer",
      data
    );
    return promise.then((data) =>
      MsgHandlePendingTransferResponse.decode(new Reader(data))
    );
  }

  SetApproval(request: MsgSetApproval): Promise<MsgSetApprovalResponse> {
    const data = MsgSetApproval.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "SetApproval",
      data
    );
    return promise.then((data) =>
      MsgSetApprovalResponse.decode(new Reader(data))
    );
  }

  RevokeBadge(request: MsgRevokeBadge): Promise<MsgRevokeBadgeResponse> {
    const data = MsgRevokeBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "RevokeBadge",
      data
    );
    return promise.then((data) =>
      MsgRevokeBadgeResponse.decode(new Reader(data))
    );
  }

  FreezeAddress(request: MsgFreezeAddress): Promise<MsgFreezeAddressResponse> {
    const data = MsgFreezeAddress.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "FreezeAddress",
      data
    );
    return promise.then((data) =>
      MsgFreezeAddressResponse.decode(new Reader(data))
    );
  }

  UpdateUris(request: MsgUpdateUris): Promise<MsgUpdateUrisResponse> {
    const data = MsgUpdateUris.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "UpdateUris",
      data
    );
    return promise.then((data) =>
      MsgUpdateUrisResponse.decode(new Reader(data))
    );
  }

  UpdatePermissions(
    request: MsgUpdatePermissions
  ): Promise<MsgUpdatePermissionsResponse> {
    const data = MsgUpdatePermissions.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "UpdatePermissions",
      data
    );
    return promise.then((data) =>
      MsgUpdatePermissionsResponse.decode(new Reader(data))
    );
  }

  TransferManager(
    request: MsgTransferManager
  ): Promise<MsgTransferManagerResponse> {
    const data = MsgTransferManager.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "TransferManager",
      data
    );
    return promise.then((data) =>
      MsgTransferManagerResponse.decode(new Reader(data))
    );
  }

  RequestTransferManager(
    request: MsgRequestTransferManager
  ): Promise<MsgRequestTransferManagerResponse> {
    const data = MsgRequestTransferManager.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "RequestTransferManager",
      data
    );
    return promise.then((data) =>
      MsgRequestTransferManagerResponse.decode(new Reader(data))
    );
  }

  SelfDestructBadge(
    request: MsgSelfDestructBadge
  ): Promise<MsgSelfDestructBadgeResponse> {
    const data = MsgSelfDestructBadge.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "SelfDestructBadge",
      data
    );
    return promise.then((data) =>
      MsgSelfDestructBadgeResponse.decode(new Reader(data))
    );
  }

  PruneBalances(request: MsgPruneBalances): Promise<MsgPruneBalancesResponse> {
    const data = MsgPruneBalances.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Msg",
      "PruneBalances",
      data
    );
    return promise.then((data) =>
      MsgPruneBalancesResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
