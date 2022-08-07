/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

export interface NumberRange {
  start: number;
  end: number;
}

export interface RangesToAmounts {
  ranges: NumberRange[];
  amount: number;
}

const baseNumberRange: object = { start: 0, end: 0 };

export const NumberRange = {
  encode(message: NumberRange, writer: Writer = Writer.create()): Writer {
    if (message.start !== 0) {
      writer.uint32(8).uint64(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).uint64(message.end);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): NumberRange {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseNumberRange } as NumberRange;
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

  fromJSON(object: any): NumberRange {
    const message = { ...baseNumberRange } as NumberRange;
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

  toJSON(message: NumberRange): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = message.start);
    message.end !== undefined && (obj.end = message.end);
    return obj;
  },

  fromPartial(object: DeepPartial<NumberRange>): NumberRange {
    const message = { ...baseNumberRange } as NumberRange;
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

const baseRangesToAmounts: object = { amount: 0 };

export const RangesToAmounts = {
  encode(message: RangesToAmounts, writer: Writer = Writer.create()): Writer {
    for (const v of message.ranges) {
      NumberRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.amount !== 0) {
      writer.uint32(16).uint64(message.amount);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RangesToAmounts {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRangesToAmounts } as RangesToAmounts;
    message.ranges = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ranges.push(NumberRange.decode(reader, reader.uint32()));
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

  fromJSON(object: any): RangesToAmounts {
    const message = { ...baseRangesToAmounts } as RangesToAmounts;
    message.ranges = [];
    if (object.ranges !== undefined && object.ranges !== null) {
      for (const e of object.ranges) {
        message.ranges.push(NumberRange.fromJSON(e));
      }
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    return message;
  },

  toJSON(message: RangesToAmounts): unknown {
    const obj: any = {};
    if (message.ranges) {
      obj.ranges = message.ranges.map((e) =>
        e ? NumberRange.toJSON(e) : undefined
      );
    } else {
      obj.ranges = [];
    }
    message.amount !== undefined && (obj.amount = message.amount);
    return obj;
  },

  fromPartial(object: DeepPartial<RangesToAmounts>): RangesToAmounts {
    const message = { ...baseRangesToAmounts } as RangesToAmounts;
    message.ranges = [];
    if (object.ranges !== undefined && object.ranges !== null) {
      for (const e of object.ranges) {
        message.ranges.push(NumberRange.fromPartial(e));
      }
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
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
