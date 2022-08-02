/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

export interface MsgNewBadge {
  creator: string;
  uri: string;
  subassetUris: string;
  permissions: number;
}

export interface MsgNewBadgeResponse {
  id: number;
}

export interface MsgNewSubBadge {
  creator: string;
  id: number;
  supply: number;
}

export interface MsgNewSubBadgeResponse {
  subassetId: number;
}

export interface MsgTransferBadge {
  creator: string;
  from: number;
  to: number;
  amount: number;
  badgeId: number;
  subbadgeId: number;
}

export interface MsgTransferBadgeResponse {}

export interface MsgRequestTransferBadge {
  creator: string;
  from: number;
  amount: number;
  badgeId: number;
  subbadgeId: number;
}

export interface MsgRequestTransferBadgeResponse {}

export interface MsgHandlePendingTransfer {
  creator: string;
  accept: boolean;
  badgeId: number;
  subbadgeId: number;
  this_nonce: number;
}

export interface MsgHandlePendingTransferResponse {}

export interface MsgSetApproval {
  creator: string;
  amount: number;
  address: number;
  badgeId: number;
  subbadgeId: number;
}

export interface MsgSetApprovalResponse {}

export interface MsgRevokeBadge {
  creator: string;
  address: number;
  amount: number;
  badgeId: number;
  subbadgeId: number;
}

export interface MsgRevokeBadgeResponse {}

export interface MsgFreezeAddress {
  creator: string;
  address: string;
  badgeId: string;
  subbadgeId: string;
}

export interface MsgFreezeAddressResponse {}

const baseMsgNewBadge: object = {
  creator: "",
  uri: "",
  subassetUris: "",
  permissions: 0,
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

const baseMsgNewSubBadge: object = { creator: "", id: 0, supply: 0 };

export const MsgNewSubBadge = {
  encode(message: MsgNewSubBadge, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    if (message.supply !== 0) {
      writer.uint32(24).uint64(message.supply);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgNewSubBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgNewSubBadge } as MsgNewSubBadge;
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
          message.supply = longToNumber(reader.uint64() as Long);
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
    if (object.supply !== undefined && object.supply !== null) {
      message.supply = Number(object.supply);
    } else {
      message.supply = 0;
    }
    return message;
  },

  toJSON(message: MsgNewSubBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = message.id);
    message.supply !== undefined && (obj.supply = message.supply);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgNewSubBadge>): MsgNewSubBadge {
    const message = { ...baseMsgNewSubBadge } as MsgNewSubBadge;
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
    if (object.supply !== undefined && object.supply !== null) {
      message.supply = object.supply;
    } else {
      message.supply = 0;
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
  to: 0,
  amount: 0,
  badgeId: 0,
  subbadgeId: 0,
};

export const MsgTransferBadge = {
  encode(message: MsgTransferBadge, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.from !== 0) {
      writer.uint32(16).uint64(message.from);
    }
    if (message.to !== 0) {
      writer.uint32(24).uint64(message.to);
    }
    if (message.amount !== 0) {
      writer.uint32(32).uint64(message.amount);
    }
    if (message.badgeId !== 0) {
      writer.uint32(40).uint64(message.badgeId);
    }
    if (message.subbadgeId !== 0) {
      writer.uint32(48).uint64(message.subbadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgTransferBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgTransferBadge } as MsgTransferBadge;
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
          message.to = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.subbadgeId = longToNumber(reader.uint64() as Long);
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
    if (object.to !== undefined && object.to !== null) {
      message.to = Number(object.to);
    } else {
      message.to = 0;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    return message;
  },

  toJSON(message: MsgTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = message.from);
    message.to !== undefined && (obj.to = message.to);
    message.amount !== undefined && (obj.amount = message.amount);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgTransferBadge>): MsgTransferBadge {
    const message = { ...baseMsgTransferBadge } as MsgTransferBadge;
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
    if (object.to !== undefined && object.to !== null) {
      message.to = object.to;
    } else {
      message.to = 0;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
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
  subbadgeId: 0,
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
    if (message.subbadgeId !== 0) {
      writer.uint32(48).uint64(message.subbadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRequestTransferBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRequestTransferBadge,
    } as MsgRequestTransferBadge;
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
          message.subbadgeId = longToNumber(reader.uint64() as Long);
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    return message;
  },

  toJSON(message: MsgRequestTransferBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.from !== undefined && (obj.from = message.from);
    message.amount !== undefined && (obj.amount = message.amount);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgRequestTransferBadge>
  ): MsgRequestTransferBadge {
    const message = {
      ...baseMsgRequestTransferBadge,
    } as MsgRequestTransferBadge;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
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
  subbadgeId: 0,
  this_nonce: 0,
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
    if (message.subbadgeId !== 0) {
      writer.uint32(32).uint64(message.subbadgeId);
    }
    if (message.this_nonce !== 0) {
      writer.uint32(40).uint64(message.this_nonce);
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
          message.subbadgeId = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.this_nonce = longToNumber(reader.uint64() as Long);
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    if (object.this_nonce !== undefined && object.this_nonce !== null) {
      message.this_nonce = Number(object.this_nonce);
    } else {
      message.this_nonce = 0;
    }
    return message;
  },

  toJSON(message: MsgHandlePendingTransfer): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.accept !== undefined && (obj.accept = message.accept);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    message.this_nonce !== undefined && (obj.this_nonce = message.this_nonce);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgHandlePendingTransfer>
  ): MsgHandlePendingTransfer {
    const message = {
      ...baseMsgHandlePendingTransfer,
    } as MsgHandlePendingTransfer;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
    }
    if (object.this_nonce !== undefined && object.this_nonce !== null) {
      message.this_nonce = object.this_nonce;
    } else {
      message.this_nonce = 0;
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
  subbadgeId: 0,
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
    if (message.subbadgeId !== 0) {
      writer.uint32(40).uint64(message.subbadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgSetApproval {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgSetApproval } as MsgSetApproval;
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
          message.subbadgeId = longToNumber(reader.uint64() as Long);
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    return message;
  },

  toJSON(message: MsgSetApproval): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.amount !== undefined && (obj.amount = message.amount);
    message.address !== undefined && (obj.address = message.address);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgSetApproval>): MsgSetApproval {
    const message = { ...baseMsgSetApproval } as MsgSetApproval;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
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
  address: 0,
  amount: 0,
  badgeId: 0,
  subbadgeId: 0,
};

export const MsgRevokeBadge = {
  encode(message: MsgRevokeBadge, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.address !== 0) {
      writer.uint32(16).uint64(message.address);
    }
    if (message.amount !== 0) {
      writer.uint32(24).uint64(message.amount);
    }
    if (message.badgeId !== 0) {
      writer.uint32(32).uint64(message.badgeId);
    }
    if (message.subbadgeId !== 0) {
      writer.uint32(40).uint64(message.subbadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRevokeBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgRevokeBadge } as MsgRevokeBadge;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.subbadgeId = longToNumber(reader.uint64() as Long);
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
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = Number(object.address);
    } else {
      message.address = 0;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    return message;
  },

  toJSON(message: MsgRevokeBadge): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.address !== undefined && (obj.address = message.address);
    message.amount !== undefined && (obj.amount = message.amount);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgRevokeBadge>): MsgRevokeBadge {
    const message = { ...baseMsgRevokeBadge } as MsgRevokeBadge;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = 0;
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
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
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

const baseMsgFreezeAddress: object = {
  creator: "",
  address: "",
  badgeId: "",
  subbadgeId: "",
};

export const MsgFreezeAddress = {
  encode(message: MsgFreezeAddress, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    if (message.badgeId !== "") {
      writer.uint32(26).string(message.badgeId);
    }
    if (message.subbadgeId !== "") {
      writer.uint32(34).string(message.subbadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgFreezeAddress {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgFreezeAddress } as MsgFreezeAddress;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.address = reader.string();
          break;
        case 3:
          message.badgeId = reader.string();
          break;
        case 4:
          message.subbadgeId = reader.string();
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
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = String(object.address);
    } else {
      message.address = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = String(object.badgeId);
    } else {
      message.badgeId = "";
    }
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = String(object.subbadgeId);
    } else {
      message.subbadgeId = "";
    }
    return message;
  },

  toJSON(message: MsgFreezeAddress): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.address !== undefined && (obj.address = message.address);
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgFreezeAddress>): MsgFreezeAddress {
    const message = { ...baseMsgFreezeAddress } as MsgFreezeAddress;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = "";
    }
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = "";
    }
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = "";
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
  /** this line is used by starport scaffolding # proto/tx/rpc */
  FreezeAddress(request: MsgFreezeAddress): Promise<MsgFreezeAddressResponse>;
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
