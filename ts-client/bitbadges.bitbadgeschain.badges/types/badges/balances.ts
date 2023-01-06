/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { BalanceObject, IdRange } from "./ranges";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

export interface WhitelistMintInfo {
  addresses: number[];
  balanceAmounts: BalanceObject[];
}

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
  subbadgeRange:
    | IdRange
    | undefined;
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

function createBaseWhitelistMintInfo(): WhitelistMintInfo {
  return { addresses: [], balanceAmounts: [] };
}

export const WhitelistMintInfo = {
  encode(message: WhitelistMintInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();
    for (const v of message.addresses) {
      writer.uint64(v);
    }
    writer.ldelim();
    for (const v of message.balanceAmounts) {
      BalanceObject.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WhitelistMintInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWhitelistMintInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.addresses.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.addresses.push(longToNumber(reader.uint64() as Long));
          }
          break;
        case 2:
          message.balanceAmounts.push(BalanceObject.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): WhitelistMintInfo {
    return {
      addresses: Array.isArray(object?.addresses) ? object.addresses.map((e: any) => Number(e)) : [],
      balanceAmounts: Array.isArray(object?.balanceAmounts)
        ? object.balanceAmounts.map((e: any) => BalanceObject.fromJSON(e))
        : [],
    };
  },

  toJSON(message: WhitelistMintInfo): unknown {
    const obj: any = {};
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => Math.round(e));
    } else {
      obj.addresses = [];
    }
    if (message.balanceAmounts) {
      obj.balanceAmounts = message.balanceAmounts.map((e) => e ? BalanceObject.toJSON(e) : undefined);
    } else {
      obj.balanceAmounts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<WhitelistMintInfo>, I>>(object: I): WhitelistMintInfo {
    const message = createBaseWhitelistMintInfo();
    message.addresses = object.addresses?.map((e) => e) || [];
    message.balanceAmounts = object.balanceAmounts?.map((e) => BalanceObject.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserBalanceInfo(): UserBalanceInfo {
  return { balanceAmounts: [], pendingNonce: 0, pending: [], approvals: [] };
}

export const UserBalanceInfo = {
  encode(message: UserBalanceInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): UserBalanceInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserBalanceInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.balanceAmounts.push(BalanceObject.decode(reader, reader.uint32()));
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
    return {
      balanceAmounts: Array.isArray(object?.balanceAmounts)
        ? object.balanceAmounts.map((e: any) => BalanceObject.fromJSON(e))
        : [],
      pendingNonce: isSet(object.pendingNonce) ? Number(object.pendingNonce) : 0,
      pending: Array.isArray(object?.pending) ? object.pending.map((e: any) => PendingTransfer.fromJSON(e)) : [],
      approvals: Array.isArray(object?.approvals) ? object.approvals.map((e: any) => Approval.fromJSON(e)) : [],
    };
  },

  toJSON(message: UserBalanceInfo): unknown {
    const obj: any = {};
    if (message.balanceAmounts) {
      obj.balanceAmounts = message.balanceAmounts.map((e) => e ? BalanceObject.toJSON(e) : undefined);
    } else {
      obj.balanceAmounts = [];
    }
    message.pendingNonce !== undefined && (obj.pendingNonce = Math.round(message.pendingNonce));
    if (message.pending) {
      obj.pending = message.pending.map((e) => e ? PendingTransfer.toJSON(e) : undefined);
    } else {
      obj.pending = [];
    }
    if (message.approvals) {
      obj.approvals = message.approvals.map((e) => e ? Approval.toJSON(e) : undefined);
    } else {
      obj.approvals = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserBalanceInfo>, I>>(object: I): UserBalanceInfo {
    const message = createBaseUserBalanceInfo();
    message.balanceAmounts = object.balanceAmounts?.map((e) => BalanceObject.fromPartial(e)) || [];
    message.pendingNonce = object.pendingNonce ?? 0;
    message.pending = object.pending?.map((e) => PendingTransfer.fromPartial(e)) || [];
    message.approvals = object.approvals?.map((e) => Approval.fromPartial(e)) || [];
    return message;
  },
};

function createBaseApproval(): Approval {
  return { address: 0, approvalAmounts: [] };
}

export const Approval = {
  encode(message: Approval, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== 0) {
      writer.uint32(8).uint64(message.address);
    }
    for (const v of message.approvalAmounts) {
      BalanceObject.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Approval {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApproval();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.approvalAmounts.push(BalanceObject.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Approval {
    return {
      address: isSet(object.address) ? Number(object.address) : 0,
      approvalAmounts: Array.isArray(object?.approvalAmounts)
        ? object.approvalAmounts.map((e: any) => BalanceObject.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Approval): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = Math.round(message.address));
    if (message.approvalAmounts) {
      obj.approvalAmounts = message.approvalAmounts.map((e) => e ? BalanceObject.toJSON(e) : undefined);
    } else {
      obj.approvalAmounts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Approval>, I>>(object: I): Approval {
    const message = createBaseApproval();
    message.address = object.address ?? 0;
    message.approvalAmounts = object.approvalAmounts?.map((e) => BalanceObject.fromPartial(e)) || [];
    return message;
  },
};

function createBasePendingTransfer(): PendingTransfer {
  return {
    subbadgeRange: undefined,
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
}

export const PendingTransfer = {
  encode(message: PendingTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): PendingTransfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePendingTransfer();
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
    return {
      subbadgeRange: isSet(object.subbadgeRange) ? IdRange.fromJSON(object.subbadgeRange) : undefined,
      thisPendingNonce: isSet(object.thisPendingNonce) ? Number(object.thisPendingNonce) : 0,
      otherPendingNonce: isSet(object.otherPendingNonce) ? Number(object.otherPendingNonce) : 0,
      amount: isSet(object.amount) ? Number(object.amount) : 0,
      sent: isSet(object.sent) ? Boolean(object.sent) : false,
      to: isSet(object.to) ? Number(object.to) : 0,
      from: isSet(object.from) ? Number(object.from) : 0,
      approvedBy: isSet(object.approvedBy) ? Number(object.approvedBy) : 0,
      markedAsAccepted: isSet(object.markedAsAccepted) ? Boolean(object.markedAsAccepted) : false,
      expirationTime: isSet(object.expirationTime) ? Number(object.expirationTime) : 0,
      cantCancelBeforeTime: isSet(object.cantCancelBeforeTime) ? Number(object.cantCancelBeforeTime) : 0,
    };
  },

  toJSON(message: PendingTransfer): unknown {
    const obj: any = {};
    message.subbadgeRange !== undefined
      && (obj.subbadgeRange = message.subbadgeRange ? IdRange.toJSON(message.subbadgeRange) : undefined);
    message.thisPendingNonce !== undefined && (obj.thisPendingNonce = Math.round(message.thisPendingNonce));
    message.otherPendingNonce !== undefined && (obj.otherPendingNonce = Math.round(message.otherPendingNonce));
    message.amount !== undefined && (obj.amount = Math.round(message.amount));
    message.sent !== undefined && (obj.sent = message.sent);
    message.to !== undefined && (obj.to = Math.round(message.to));
    message.from !== undefined && (obj.from = Math.round(message.from));
    message.approvedBy !== undefined && (obj.approvedBy = Math.round(message.approvedBy));
    message.markedAsAccepted !== undefined && (obj.markedAsAccepted = message.markedAsAccepted);
    message.expirationTime !== undefined && (obj.expirationTime = Math.round(message.expirationTime));
    message.cantCancelBeforeTime !== undefined && (obj.cantCancelBeforeTime = Math.round(message.cantCancelBeforeTime));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PendingTransfer>, I>>(object: I): PendingTransfer {
    const message = createBasePendingTransfer();
    message.subbadgeRange = (object.subbadgeRange !== undefined && object.subbadgeRange !== null)
      ? IdRange.fromPartial(object.subbadgeRange)
      : undefined;
    message.thisPendingNonce = object.thisPendingNonce ?? 0;
    message.otherPendingNonce = object.otherPendingNonce ?? 0;
    message.amount = object.amount ?? 0;
    message.sent = object.sent ?? false;
    message.to = object.to ?? 0;
    message.from = object.from ?? 0;
    message.approvedBy = object.approvedBy ?? 0;
    message.markedAsAccepted = object.markedAsAccepted ?? false;
    message.expirationTime = object.expirationTime ?? 0;
    message.cantCancelBeforeTime = object.cantCancelBeforeTime ?? 0;
    return message;
  },
};

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
