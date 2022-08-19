/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** Id ranges define a range of IDs from start to end. Can be used for subbadgeIds, nonces, addresses anything. If end == 0, we assume end == start. Start must be >= end. */
export interface IdRange {
  start: number;
  end: number;
}

/** Defines a balance object. The specified balance holds for all ids specified within the id ranges array. */
export interface BalanceObject {
  balance: number;
  id_ranges: IdRange[];
}

const baseIdRange: object = { start: 0, end: 0 };

export const IdRange = {
  encode(message: IdRange, writer: Writer = Writer.create()): Writer {
    if (message.start !== 0) {
      writer.uint32(8).uint64(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).uint64(message.end);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): IdRange {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseIdRange } as IdRange;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.end = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IdRange {
    const message = { ...baseIdRange } as IdRange;
    if (object.start !== undefined && object.start !== null) {
      message.start = Number(object.start);
    } else {
      message.start = 0;
    }
    if (object.end !== undefined && object.end !== null) {
      message.end = Number(object.end);
    } else {
      message.end = 0;
    }
    return message;
  },

  toJSON(message: IdRange): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = message.start);
    message.end !== undefined && (obj.end = message.end);
    return obj;
  },

  fromPartial(object: DeepPartial<IdRange>): IdRange {
    const message = { ...baseIdRange } as IdRange;
    if (object.start !== undefined && object.start !== null) {
      message.start = object.start;
    } else {
      message.start = 0;
    }
    if (object.end !== undefined && object.end !== null) {
      message.end = object.end;
    } else {
      message.end = 0;
    }
    return message;
  },
};

const baseBalanceObject: object = { balance: 0 };

export const BalanceObject = {
  encode(message: BalanceObject, writer: Writer = Writer.create()): Writer {
    if (message.balance !== 0) {
      writer.uint32(8).uint64(message.balance);
    }
    for (const v of message.id_ranges) {
      IdRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BalanceObject {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBalanceObject } as BalanceObject;
    message.id_ranges = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balance = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.id_ranges.push(IdRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BalanceObject {
    const message = { ...baseBalanceObject } as BalanceObject;
    message.id_ranges = [];
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Number(object.balance);
    } else {
      message.balance = 0;
    }
    if (object.id_ranges !== undefined && object.id_ranges !== null) {
      for (const e of object.id_ranges) {
        message.id_ranges.push(IdRange.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: BalanceObject): unknown {
    const obj: any = {};
    message.balance !== undefined && (obj.balance = message.balance);
    if (message.id_ranges) {
      obj.id_ranges = message.id_ranges.map((e) =>
        e ? IdRange.toJSON(e) : undefined
      );
    } else {
      obj.id_ranges = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BalanceObject>): BalanceObject {
    const message = { ...baseBalanceObject } as BalanceObject;
    message.id_ranges = [];
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    } else {
      message.balance = 0;
    }
    if (object.id_ranges !== undefined && object.id_ranges !== null) {
      for (const e of object.id_ranges) {
        message.id_ranges.push(IdRange.fromPartial(e));
      }
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
