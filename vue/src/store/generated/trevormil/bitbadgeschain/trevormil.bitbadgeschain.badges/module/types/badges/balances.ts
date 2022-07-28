/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** indexed by badgeid-subassetid-uniqueaccountnumber (26 bytes) */
export interface BadgeBalanceInfo {
  balance: number;
  pending_nonce: number;
  /** IDs will be sorted in order of pending_nonce */
  pending: PendingTransfer[];
  approvals: Approval[];
}

export interface Approval {
  address: string;
  amount: number;
}

/** Pending transfers will not be saved after accept / reject */
export interface PendingTransfer {
  /** indexed by (this account's pending_nonce - other account's pending_nonce) */
  id: string;
  amount: number;
  /** vs. receive request */
  send_request: boolean;
  to: number;
  from: number;
  memo: string;
}

const baseBadgeBalanceInfo: object = { balance: 0, pending_nonce: 0 };

export const BadgeBalanceInfo = {
  encode(message: BadgeBalanceInfo, writer: Writer = Writer.create()): Writer {
    if (message.balance !== 0) {
      writer.uint32(8).uint64(message.balance);
    }
    if (message.pending_nonce !== 0) {
      writer.uint32(16).uint64(message.pending_nonce);
    }
    for (const v of message.pending) {
      PendingTransfer.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.approvals) {
      Approval.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BadgeBalanceInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBadgeBalanceInfo } as BadgeBalanceInfo;
    message.pending = [];
    message.approvals = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balance = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.pending_nonce = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.pending.push(PendingTransfer.decode(reader, reader.uint32()));
          break;
        case 4:
          message.approvals.push(Approval.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeBalanceInfo {
    const message = { ...baseBadgeBalanceInfo } as BadgeBalanceInfo;
    message.pending = [];
    message.approvals = [];
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Number(object.balance);
    } else {
      message.balance = 0;
    }
    if (object.pending_nonce !== undefined && object.pending_nonce !== null) {
      message.pending_nonce = Number(object.pending_nonce);
    } else {
      message.pending_nonce = 0;
    }
    if (object.pending !== undefined && object.pending !== null) {
      for (const e of object.pending) {
        message.pending.push(PendingTransfer.fromJSON(e));
      }
    }
    if (object.approvals !== undefined && object.approvals !== null) {
      for (const e of object.approvals) {
        message.approvals.push(Approval.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: BadgeBalanceInfo): unknown {
    const obj: any = {};
    message.balance !== undefined && (obj.balance = message.balance);
    message.pending_nonce !== undefined &&
      (obj.pending_nonce = message.pending_nonce);
    if (message.pending) {
      obj.pending = message.pending.map((e) =>
        e ? PendingTransfer.toJSON(e) : undefined
      );
    } else {
      obj.pending = [];
    }
    if (message.approvals) {
      obj.approvals = message.approvals.map((e) =>
        e ? Approval.toJSON(e) : undefined
      );
    } else {
      obj.approvals = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BadgeBalanceInfo>): BadgeBalanceInfo {
    const message = { ...baseBadgeBalanceInfo } as BadgeBalanceInfo;
    message.pending = [];
    message.approvals = [];
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    } else {
      message.balance = 0;
    }
    if (object.pending_nonce !== undefined && object.pending_nonce !== null) {
      message.pending_nonce = object.pending_nonce;
    } else {
      message.pending_nonce = 0;
    }
    if (object.pending !== undefined && object.pending !== null) {
      for (const e of object.pending) {
        message.pending.push(PendingTransfer.fromPartial(e));
      }
    }
    if (object.approvals !== undefined && object.approvals !== null) {
      for (const e of object.approvals) {
        message.approvals.push(Approval.fromPartial(e));
      }
    }
    return message;
  },
};

const baseApproval: object = { address: "", amount: 0 };

export const Approval = {
  encode(message: Approval, writer: Writer = Writer.create()): Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.amount !== 0) {
      writer.uint32(16).uint64(message.amount);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Approval {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseApproval } as Approval;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Approval {
    const message = { ...baseApproval } as Approval;
    if (object.address !== undefined && object.address !== null) {
      message.address = String(object.address);
    } else {
      message.address = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    return message;
  },

  toJSON(message: Approval): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.amount !== undefined && (obj.amount = message.amount);
    return obj;
  },

  fromPartial(object: DeepPartial<Approval>): Approval {
    const message = { ...baseApproval } as Approval;
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    return message;
  },
};

const basePendingTransfer: object = {
  id: "",
  amount: 0,
  send_request: false,
  to: 0,
  from: 0,
  memo: "",
};

export const PendingTransfer = {
  encode(message: PendingTransfer, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.amount !== 0) {
      writer.uint32(16).uint64(message.amount);
    }
    if (message.send_request === true) {
      writer.uint32(24).bool(message.send_request);
    }
    if (message.to !== 0) {
      writer.uint32(32).uint64(message.to);
    }
    if (message.from !== 0) {
      writer.uint32(40).uint64(message.from);
    }
    if (message.memo !== "") {
      writer.uint32(50).string(message.memo);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PendingTransfer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePendingTransfer } as PendingTransfer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.send_request = reader.bool();
          break;
        case 4:
          message.to = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.from = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.memo = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PendingTransfer {
    const message = { ...basePendingTransfer } as PendingTransfer;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.send_request !== undefined && object.send_request !== null) {
      message.send_request = Boolean(object.send_request);
    } else {
      message.send_request = false;
    }
    if (object.to !== undefined && object.to !== null) {
      message.to = Number(object.to);
    } else {
      message.to = 0;
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = Number(object.from);
    } else {
      message.from = 0;
    }
    if (object.memo !== undefined && object.memo !== null) {
      message.memo = String(object.memo);
    } else {
      message.memo = "";
    }
    return message;
  },

  toJSON(message: PendingTransfer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.amount !== undefined && (obj.amount = message.amount);
    message.send_request !== undefined &&
      (obj.send_request = message.send_request);
    message.to !== undefined && (obj.to = message.to);
    message.from !== undefined && (obj.from = message.from);
    message.memo !== undefined && (obj.memo = message.memo);
    return obj;
  },

  fromPartial(object: DeepPartial<PendingTransfer>): PendingTransfer {
    const message = { ...basePendingTransfer } as PendingTransfer;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.send_request !== undefined && object.send_request !== null) {
      message.send_request = object.send_request;
    } else {
      message.send_request = false;
    }
    if (object.to !== undefined && object.to !== null) {
      message.to = object.to;
    } else {
      message.to = 0;
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = object.from;
    } else {
      message.from = 0;
    }
    if (object.memo !== undefined && object.memo !== null) {
      message.memo = object.memo;
    } else {
      message.memo = "";
    }
    return message;
  },
};

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
