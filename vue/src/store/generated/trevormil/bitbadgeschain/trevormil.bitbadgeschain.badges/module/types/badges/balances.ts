/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { RangesToAmounts, NumberRange } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** indexed by badgeid-subassetid-uniqueaccountnumber (26 bytes) */
export interface BadgeBalanceInfo {
  balanceAmounts: RangesToAmounts[];
  pending_nonce: number;
  /** IDs will be sorted in order of pending_nonce */
  pending: PendingTransfer[];
  approvals: Approval[];
  /** TODO: for (hidden on profile, pinned, etc) */
  user_flags: number;
}

export interface Approval {
  address: number;
  approvalAmounts: RangesToAmounts[];
}

export interface PendingTransfer {
  subbadgeRange: NumberRange | undefined;
  this_pending_nonce: number;
  other_pending_nonce: number;
  amount: number;
  /** vs. receive request */
  send_request: boolean;
  to: number;
  from: number;
  approved_by: number;
  markedAsApproved: boolean;
}

const baseBadgeBalanceInfo: object = { pending_nonce: 0, user_flags: 0 };

export const BadgeBalanceInfo = {
  encode(message: BadgeBalanceInfo, writer: Writer = Writer.create()): Writer {
    for (const v of message.balanceAmounts) {
      RangesToAmounts.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.pending_nonce !== 0) {
      writer.uint32(24).uint64(message.pending_nonce);
    }
    for (const v of message.pending) {
      PendingTransfer.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.approvals) {
      Approval.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.user_flags !== 0) {
      writer.uint32(48).uint64(message.user_flags);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BadgeBalanceInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBadgeBalanceInfo } as BadgeBalanceInfo;
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.balanceAmounts.push(
            RangesToAmounts.decode(reader, reader.uint32())
          );
          break;
        case 3:
          message.pending_nonce = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.pending.push(PendingTransfer.decode(reader, reader.uint32()));
          break;
        case 5:
          message.approvals.push(Approval.decode(reader, reader.uint32()));
          break;
        case 6:
          message.user_flags = longToNumber(reader.uint64() as Long);
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
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    if (object.balanceAmounts !== undefined && object.balanceAmounts !== null) {
      for (const e of object.balanceAmounts) {
        message.balanceAmounts.push(RangesToAmounts.fromJSON(e));
      }
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
    if (object.user_flags !== undefined && object.user_flags !== null) {
      message.user_flags = Number(object.user_flags);
    } else {
      message.user_flags = 0;
    }
    return message;
  },

  toJSON(message: BadgeBalanceInfo): unknown {
    const obj: any = {};
    if (message.balanceAmounts) {
      obj.balanceAmounts = message.balanceAmounts.map((e) =>
        e ? RangesToAmounts.toJSON(e) : undefined
      );
    } else {
      obj.balanceAmounts = [];
    }
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
    message.user_flags !== undefined && (obj.user_flags = message.user_flags);
    return obj;
  },

  fromPartial(object: DeepPartial<BadgeBalanceInfo>): BadgeBalanceInfo {
    const message = { ...baseBadgeBalanceInfo } as BadgeBalanceInfo;
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    if (object.balanceAmounts !== undefined && object.balanceAmounts !== null) {
      for (const e of object.balanceAmounts) {
        message.balanceAmounts.push(RangesToAmounts.fromPartial(e));
      }
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
    if (object.user_flags !== undefined && object.user_flags !== null) {
      message.user_flags = object.user_flags;
    } else {
      message.user_flags = 0;
    }
    return message;
  },
};

const baseApproval: object = { address: 0 };

export const Approval = {
  encode(message: Approval, writer: Writer = Writer.create()): Writer {
    if (message.address !== 0) {
      writer.uint32(8).uint64(message.address);
    }
    for (const v of message.approvalAmounts) {
      RangesToAmounts.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Approval {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseApproval } as Approval;
    message.approvalAmounts = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.approvalAmounts.push(
            RangesToAmounts.decode(reader, reader.uint32())
          );
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
    message.approvalAmounts = [];
    if (object.address !== undefined && object.address !== null) {
      message.address = Number(object.address);
    } else {
      message.address = 0;
    }
    if (
      object.approvalAmounts !== undefined &&
      object.approvalAmounts !== null
    ) {
      for (const e of object.approvalAmounts) {
        message.approvalAmounts.push(RangesToAmounts.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: Approval): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    if (message.approvalAmounts) {
      obj.approvalAmounts = message.approvalAmounts.map((e) =>
        e ? RangesToAmounts.toJSON(e) : undefined
      );
    } else {
      obj.approvalAmounts = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<Approval>): Approval {
    const message = { ...baseApproval } as Approval;
    message.approvalAmounts = [];
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = 0;
    }
    if (
      object.approvalAmounts !== undefined &&
      object.approvalAmounts !== null
    ) {
      for (const e of object.approvalAmounts) {
        message.approvalAmounts.push(RangesToAmounts.fromPartial(e));
      }
    }
    return message;
  },
};

const basePendingTransfer: object = {
  this_pending_nonce: 0,
  other_pending_nonce: 0,
  amount: 0,
  send_request: false,
  to: 0,
  from: 0,
  approved_by: 0,
  markedAsApproved: false,
};

export const PendingTransfer = {
  encode(message: PendingTransfer, writer: Writer = Writer.create()): Writer {
    if (message.subbadgeRange !== undefined) {
      NumberRange.encode(
        message.subbadgeRange,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.this_pending_nonce !== 0) {
      writer.uint32(16).uint64(message.this_pending_nonce);
    }
    if (message.other_pending_nonce !== 0) {
      writer.uint32(24).uint64(message.other_pending_nonce);
    }
    if (message.amount !== 0) {
      writer.uint32(32).uint64(message.amount);
    }
    if (message.send_request === true) {
      writer.uint32(40).bool(message.send_request);
    }
    if (message.to !== 0) {
      writer.uint32(48).uint64(message.to);
    }
    if (message.from !== 0) {
      writer.uint32(56).uint64(message.from);
    }
    if (message.approved_by !== 0) {
      writer.uint32(72).uint64(message.approved_by);
    }
    if (message.markedAsApproved === true) {
      writer.uint32(80).bool(message.markedAsApproved);
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
          message.subbadgeRange = NumberRange.decode(reader, reader.uint32());
          break;
        case 2:
          message.this_pending_nonce = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.other_pending_nonce = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.send_request = reader.bool();
          break;
        case 6:
          message.to = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.from = longToNumber(reader.uint64() as Long);
          break;
        case 9:
          message.approved_by = longToNumber(reader.uint64() as Long);
          break;
        case 10:
          message.markedAsApproved = reader.bool();
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
    if (object.subbadgeRange !== undefined && object.subbadgeRange !== null) {
      message.subbadgeRange = NumberRange.fromJSON(object.subbadgeRange);
    } else {
      message.subbadgeRange = undefined;
    }
    if (
      object.this_pending_nonce !== undefined &&
      object.this_pending_nonce !== null
    ) {
      message.this_pending_nonce = Number(object.this_pending_nonce);
    } else {
      message.this_pending_nonce = 0;
    }
    if (
      object.other_pending_nonce !== undefined &&
      object.other_pending_nonce !== null
    ) {
      message.other_pending_nonce = Number(object.other_pending_nonce);
    } else {
      message.other_pending_nonce = 0;
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
    if (object.approved_by !== undefined && object.approved_by !== null) {
      message.approved_by = Number(object.approved_by);
    } else {
      message.approved_by = 0;
    }
    if (
      object.markedAsApproved !== undefined &&
      object.markedAsApproved !== null
    ) {
      message.markedAsApproved = Boolean(object.markedAsApproved);
    } else {
      message.markedAsApproved = false;
    }
    return message;
  },

  toJSON(message: PendingTransfer): unknown {
    const obj: any = {};
    message.subbadgeRange !== undefined &&
      (obj.subbadgeRange = message.subbadgeRange
        ? NumberRange.toJSON(message.subbadgeRange)
        : undefined);
    message.this_pending_nonce !== undefined &&
      (obj.this_pending_nonce = message.this_pending_nonce);
    message.other_pending_nonce !== undefined &&
      (obj.other_pending_nonce = message.other_pending_nonce);
    message.amount !== undefined && (obj.amount = message.amount);
    message.send_request !== undefined &&
      (obj.send_request = message.send_request);
    message.to !== undefined && (obj.to = message.to);
    message.from !== undefined && (obj.from = message.from);
    message.approved_by !== undefined &&
      (obj.approved_by = message.approved_by);
    message.markedAsApproved !== undefined &&
      (obj.markedAsApproved = message.markedAsApproved);
    return obj;
  },

  fromPartial(object: DeepPartial<PendingTransfer>): PendingTransfer {
    const message = { ...basePendingTransfer } as PendingTransfer;
    if (object.subbadgeRange !== undefined && object.subbadgeRange !== null) {
      message.subbadgeRange = NumberRange.fromPartial(object.subbadgeRange);
    } else {
      message.subbadgeRange = undefined;
    }
    if (
      object.this_pending_nonce !== undefined &&
      object.this_pending_nonce !== null
    ) {
      message.this_pending_nonce = object.this_pending_nonce;
    } else {
      message.this_pending_nonce = 0;
    }
    if (
      object.other_pending_nonce !== undefined &&
      object.other_pending_nonce !== null
    ) {
      message.other_pending_nonce = object.other_pending_nonce;
    } else {
      message.other_pending_nonce = 0;
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
    if (object.approved_by !== undefined && object.approved_by !== null) {
      message.approved_by = object.approved_by;
    } else {
      message.approved_by = 0;
    }
    if (
      object.markedAsApproved !== undefined &&
      object.markedAsApproved !== null
    ) {
      message.markedAsApproved = object.markedAsApproved;
    } else {
      message.markedAsApproved = false;
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
