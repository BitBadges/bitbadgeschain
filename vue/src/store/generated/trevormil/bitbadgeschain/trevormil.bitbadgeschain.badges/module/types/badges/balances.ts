/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { BalanceObject, IdRange } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** Defines a user balance object for a badge w/ the user's balances, nonce, pending transfers, and approvals. All subbadge IDs for a badge are handled within this object. */
export interface UserBalanceInfo {
  /** The user's balance for each subbadge. */
  balanceAmounts: BalanceObject[];
  /** Nonce for pending transfers. Increments by 1 each time. */
  pendingNonce: number;
  /** IDs will be sorted in order of this account's pending nonce. */
  pending: PendingTransfer[];
  /** Approvals are sorted in order of address number. */
  approvals: Approval[];
}

/** Defines an approval object for a specific address. */
export interface Approval {
  /** account number for the address */
  address: number;
  /** approval balances for every subbadgeId */
  approvalAmounts: BalanceObject[];
}

/** Defines a pending transfer object for two addresses. A pending transfer will be stored in both parties' balance objects. */
export interface PendingTransfer {
  subbadgeRange: IdRange | undefined;
  /** This pending nonce is the nonce of the account for which this transfer is stored. Other is the other party's. Will be swapped for the other party's stored pending transfer. */
  thisPendingNonce: number;
  otherPendingNonce: number;
  amount: number;
  /** Sent defines who initiated this pending transfer */
  sent: boolean;
  to: number;
  from: number;
  approvedBy: number;
  /** For non forceful accepts, this will be true if the other party has accepted but doesn't want to pay the gas fees. */
  markedAsAccepted: boolean;
  /** Can't be accepted after expiration time. If == 0, we assume it never expires. */
  expirationTime: number;
  /** Can't cancel before must be less than expiration time. If == 0, we assume can cancel at any time. */
  cantCancelBeforeTime: number;
}

const baseUserBalanceInfo: object = { pendingNonce: 0 };

export const UserBalanceInfo = {
  encode(message: UserBalanceInfo, writer: Writer = Writer.create()): Writer {
    for (const v of message.balanceAmounts) {
      BalanceObject.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.pendingNonce !== 0) {
      writer.uint32(24).uint64(message.pendingNonce);
    }
    for (const v of message.pending) {
      PendingTransfer.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.approvals) {
      Approval.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UserBalanceInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUserBalanceInfo } as UserBalanceInfo;
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.balanceAmounts.push(
            BalanceObject.decode(reader, reader.uint32())
          );
          break;
        case 3:
          message.pendingNonce = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.pending.push(PendingTransfer.decode(reader, reader.uint32()));
          break;
        case 5:
          message.approvals.push(Approval.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserBalanceInfo {
    const message = { ...baseUserBalanceInfo } as UserBalanceInfo;
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    if (object.balanceAmounts !== undefined && object.balanceAmounts !== null) {
      for (const e of object.balanceAmounts) {
        message.balanceAmounts.push(BalanceObject.fromJSON(e));
      }
    }
    if (object.pendingNonce !== undefined && object.pendingNonce !== null) {
      message.pendingNonce = Number(object.pendingNonce);
    } else {
      message.pendingNonce = 0;
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

  toJSON(message: UserBalanceInfo): unknown {
    const obj: any = {};
    if (message.balanceAmounts) {
      obj.balanceAmounts = message.balanceAmounts.map((e) =>
        e ? BalanceObject.toJSON(e) : undefined
      );
    } else {
      obj.balanceAmounts = [];
    }
    message.pendingNonce !== undefined &&
      (obj.pendingNonce = message.pendingNonce);
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

  fromPartial(object: DeepPartial<UserBalanceInfo>): UserBalanceInfo {
    const message = { ...baseUserBalanceInfo } as UserBalanceInfo;
    message.balanceAmounts = [];
    message.pending = [];
    message.approvals = [];
    if (object.balanceAmounts !== undefined && object.balanceAmounts !== null) {
      for (const e of object.balanceAmounts) {
        message.balanceAmounts.push(BalanceObject.fromPartial(e));
      }
    }
    if (object.pendingNonce !== undefined && object.pendingNonce !== null) {
      message.pendingNonce = object.pendingNonce;
    } else {
      message.pendingNonce = 0;
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

const baseApproval: object = { address: 0 };

export const Approval = {
  encode(message: Approval, writer: Writer = Writer.create()): Writer {
    if (message.address !== 0) {
      writer.uint32(8).uint64(message.address);
    }
    for (const v of message.approvalAmounts) {
      BalanceObject.encode(v!, writer.uint32(18).fork()).ldelim();
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
            BalanceObject.decode(reader, reader.uint32())
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
        message.approvalAmounts.push(BalanceObject.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: Approval): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    if (message.approvalAmounts) {
      obj.approvalAmounts = message.approvalAmounts.map((e) =>
        e ? BalanceObject.toJSON(e) : undefined
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
        message.approvalAmounts.push(BalanceObject.fromPartial(e));
      }
    }
    return message;
  },
};

const basePendingTransfer: object = {
  thisPendingNonce: 0,
  otherPendingNonce: 0,
  amount: 0,
  sent: false,
  to: 0,
  from: 0,
  approvedBy: 0,
  markedAsAccepted: false,
  expirationTime: 0,
  cantCancelBeforeTime: 0,
};

export const PendingTransfer = {
  encode(message: PendingTransfer, writer: Writer = Writer.create()): Writer {
    if (message.subbadgeRange !== undefined) {
      IdRange.encode(message.subbadgeRange, writer.uint32(10).fork()).ldelim();
    }
    if (message.thisPendingNonce !== 0) {
      writer.uint32(16).uint64(message.thisPendingNonce);
    }
    if (message.otherPendingNonce !== 0) {
      writer.uint32(24).uint64(message.otherPendingNonce);
    }
    if (message.amount !== 0) {
      writer.uint32(32).uint64(message.amount);
    }
    if (message.sent === true) {
      writer.uint32(40).bool(message.sent);
    }
    if (message.to !== 0) {
      writer.uint32(48).uint64(message.to);
    }
    if (message.from !== 0) {
      writer.uint32(56).uint64(message.from);
    }
    if (message.approvedBy !== 0) {
      writer.uint32(72).uint64(message.approvedBy);
    }
    if (message.markedAsAccepted === true) {
      writer.uint32(80).bool(message.markedAsAccepted);
    }
    if (message.expirationTime !== 0) {
      writer.uint32(88).uint64(message.expirationTime);
    }
    if (message.cantCancelBeforeTime !== 0) {
      writer.uint32(96).uint64(message.cantCancelBeforeTime);
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
          message.subbadgeRange = IdRange.decode(reader, reader.uint32());
          break;
        case 2:
          message.thisPendingNonce = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.otherPendingNonce = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.sent = reader.bool();
          break;
        case 6:
          message.to = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.from = longToNumber(reader.uint64() as Long);
          break;
        case 9:
          message.approvedBy = longToNumber(reader.uint64() as Long);
          break;
        case 10:
          message.markedAsAccepted = reader.bool();
          break;
        case 11:
          message.expirationTime = longToNumber(reader.uint64() as Long);
          break;
        case 12:
          message.cantCancelBeforeTime = longToNumber(reader.uint64() as Long);
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
      message.subbadgeRange = IdRange.fromJSON(object.subbadgeRange);
    } else {
      message.subbadgeRange = undefined;
    }
    if (
      object.thisPendingNonce !== undefined &&
      object.thisPendingNonce !== null
    ) {
      message.thisPendingNonce = Number(object.thisPendingNonce);
    } else {
      message.thisPendingNonce = 0;
    }
    if (
      object.otherPendingNonce !== undefined &&
      object.otherPendingNonce !== null
    ) {
      message.otherPendingNonce = Number(object.otherPendingNonce);
    } else {
      message.otherPendingNonce = 0;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.sent !== undefined && object.sent !== null) {
      message.sent = Boolean(object.sent);
    } else {
      message.sent = false;
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
    if (object.approvedBy !== undefined && object.approvedBy !== null) {
      message.approvedBy = Number(object.approvedBy);
    } else {
      message.approvedBy = 0;
    }
    if (
      object.markedAsAccepted !== undefined &&
      object.markedAsAccepted !== null
    ) {
      message.markedAsAccepted = Boolean(object.markedAsAccepted);
    } else {
      message.markedAsAccepted = false;
    }
    if (object.expirationTime !== undefined && object.expirationTime !== null) {
      message.expirationTime = Number(object.expirationTime);
    } else {
      message.expirationTime = 0;
    }
    if (
      object.cantCancelBeforeTime !== undefined &&
      object.cantCancelBeforeTime !== null
    ) {
      message.cantCancelBeforeTime = Number(object.cantCancelBeforeTime);
    } else {
      message.cantCancelBeforeTime = 0;
    }
    return message;
  },

  toJSON(message: PendingTransfer): unknown {
    const obj: any = {};
    message.subbadgeRange !== undefined &&
      (obj.subbadgeRange = message.subbadgeRange
        ? IdRange.toJSON(message.subbadgeRange)
        : undefined);
    message.thisPendingNonce !== undefined &&
      (obj.thisPendingNonce = message.thisPendingNonce);
    message.otherPendingNonce !== undefined &&
      (obj.otherPendingNonce = message.otherPendingNonce);
    message.amount !== undefined && (obj.amount = message.amount);
    message.sent !== undefined && (obj.sent = message.sent);
    message.to !== undefined && (obj.to = message.to);
    message.from !== undefined && (obj.from = message.from);
    message.approvedBy !== undefined && (obj.approvedBy = message.approvedBy);
    message.markedAsAccepted !== undefined &&
      (obj.markedAsAccepted = message.markedAsAccepted);
    message.expirationTime !== undefined &&
      (obj.expirationTime = message.expirationTime);
    message.cantCancelBeforeTime !== undefined &&
      (obj.cantCancelBeforeTime = message.cantCancelBeforeTime);
    return obj;
  },

  fromPartial(object: DeepPartial<PendingTransfer>): PendingTransfer {
    const message = { ...basePendingTransfer } as PendingTransfer;
    if (object.subbadgeRange !== undefined && object.subbadgeRange !== null) {
      message.subbadgeRange = IdRange.fromPartial(object.subbadgeRange);
    } else {
      message.subbadgeRange = undefined;
    }
    if (
      object.thisPendingNonce !== undefined &&
      object.thisPendingNonce !== null
    ) {
      message.thisPendingNonce = object.thisPendingNonce;
    } else {
      message.thisPendingNonce = 0;
    }
    if (
      object.otherPendingNonce !== undefined &&
      object.otherPendingNonce !== null
    ) {
      message.otherPendingNonce = object.otherPendingNonce;
    } else {
      message.otherPendingNonce = 0;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.sent !== undefined && object.sent !== null) {
      message.sent = object.sent;
    } else {
      message.sent = false;
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
    if (object.approvedBy !== undefined && object.approvedBy !== null) {
      message.approvedBy = object.approvedBy;
    } else {
      message.approvedBy = 0;
    }
    if (
      object.markedAsAccepted !== undefined &&
      object.markedAsAccepted !== null
    ) {
      message.markedAsAccepted = object.markedAsAccepted;
    } else {
      message.markedAsAccepted = false;
    }
    if (object.expirationTime !== undefined && object.expirationTime !== null) {
      message.expirationTime = object.expirationTime;
    } else {
      message.expirationTime = 0;
    }
    if (
      object.cantCancelBeforeTime !== undefined &&
      object.cantCancelBeforeTime !== null
    ) {
      message.cantCancelBeforeTime = object.cantCancelBeforeTime;
    } else {
      message.cantCancelBeforeTime = 0;
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
