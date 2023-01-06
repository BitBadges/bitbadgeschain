/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** Id ranges define a range of IDs from start to end. Can be used for subbadgeIds, nonces, addresses anything. If end == 0, we assume end == start. Start must be >= end. */
export interface IdRange {
  start: number;
  end: number;
}

/** Defines a balance object. The specified balance holds for all ids specified within the id ranges array. */
export interface BalanceObject {
  balance: number;
  idRanges: IdRange[];
}

function createBaseIdRange(): IdRange {
  return { start: 0, end: 0 };
}

export const IdRange = {
  encode(message: IdRange, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== 0) {
      writer.uint32(8).uint64(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).uint64(message.end);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdRange {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdRange();
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
    return { start: isSet(object.start) ? Number(object.start) : 0, end: isSet(object.end) ? Number(object.end) : 0 };
  },

  toJSON(message: IdRange): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.end !== undefined && (obj.end = Math.round(message.end));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IdRange>, I>>(object: I): IdRange {
    const message = createBaseIdRange();
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    return message;
  },
};

function createBaseBalanceObject(): BalanceObject {
  return { balance: 0, idRanges: [] };
}

export const BalanceObject = {
  encode(message: BalanceObject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.balance !== 0) {
      writer.uint32(8).uint64(message.balance);
    }
    for (const v of message.idRanges) {
      IdRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BalanceObject {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalanceObject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balance = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.idRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BalanceObject {
    return {
      balance: isSet(object.balance) ? Number(object.balance) : 0,
      idRanges: Array.isArray(object?.idRanges) ? object.idRanges.map((e: any) => IdRange.fromJSON(e)) : [],
    };
  },

  toJSON(message: BalanceObject): unknown {
    const obj: any = {};
    message.balance !== undefined && (obj.balance = Math.round(message.balance));
    if (message.idRanges) {
      obj.idRanges = message.idRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.idRanges = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BalanceObject>, I>>(object: I): BalanceObject {
    const message = createBaseBalanceObject();
    message.balance = object.balance ?? 0;
    message.idRanges = object.idRanges?.map((e) => IdRange.fromPartial(e)) || [];
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
