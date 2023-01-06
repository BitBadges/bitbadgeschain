/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { IdRange } from "./ranges";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** A URI object defines a uri and subasset uri for a badge and its subbadges. Designed to save storage and avoid reused text and common patterns. */
export interface UriObject {
  /** This will be == 0 represeting plaintext URLs for now, but in the future, if we want to add other decoding / encoding schemes like Base64, etc, we just define a new decoding scheme */
  decodeScheme: number;
  /** Helps to save space by not storing methods like https:// every time. If this == 0, we assume it's included in the uri itself. Else, we define a few predefined int -> scheme maps in uris.go that will prefix the uri bytes. */
  scheme: number;
  /** The uri bytes to store. Will be converted to a string. To be manipulated according to other properties of this URI object. */
  uri: string;
  /** The four fields below are used to convert the uri from above to the subasset URI. */
  idxRangeToRemove: IdRange | undefined;
  insertSubassetBytesIdx: number;
  /** After removing the above range, insert these bytes at insertSubassetBytesIdx. insertSubassetBytesIdx is the idx after the range removal, not before. */
  bytesToInsert: string;
  /** This is the idx where we insert the id number of the subbadge. For example, ex.com/ and insertIdIdx == 6, the subasset URI for ID 0 would be ex.com/0 */
  insertIdIdx: number;
}

function createBaseUriObject(): UriObject {
  return {
    decodeScheme: 0,
    scheme: 0,
    uri: "",
    idxRangeToRemove: undefined,
    insertSubassetBytesIdx: 0,
    bytesToInsert: "",
    insertIdIdx: 0,
  };
}

export const UriObject = {
  encode(message: UriObject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.decodeScheme !== 0) {
      writer.uint32(8).uint64(message.decodeScheme);
    }
    if (message.scheme !== 0) {
      writer.uint32(16).uint64(message.scheme);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.idxRangeToRemove !== undefined) {
      IdRange.encode(message.idxRangeToRemove, writer.uint32(34).fork()).ldelim();
    }
    if (message.insertSubassetBytesIdx !== 0) {
      writer.uint32(40).uint64(message.insertSubassetBytesIdx);
    }
    if (message.bytesToInsert !== "") {
      writer.uint32(50).string(message.bytesToInsert);
    }
    if (message.insertIdIdx !== 0) {
      writer.uint32(56).uint64(message.insertIdIdx);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UriObject {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUriObject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.decodeScheme = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.scheme = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.uri = reader.string();
          break;
        case 4:
          message.idxRangeToRemove = IdRange.decode(reader, reader.uint32());
          break;
        case 5:
          message.insertSubassetBytesIdx = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.bytesToInsert = reader.string();
          break;
        case 7:
          message.insertIdIdx = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UriObject {
    return {
      decodeScheme: isSet(object.decodeScheme) ? Number(object.decodeScheme) : 0,
      scheme: isSet(object.scheme) ? Number(object.scheme) : 0,
      uri: isSet(object.uri) ? String(object.uri) : "",
      idxRangeToRemove: isSet(object.idxRangeToRemove) ? IdRange.fromJSON(object.idxRangeToRemove) : undefined,
      insertSubassetBytesIdx: isSet(object.insertSubassetBytesIdx) ? Number(object.insertSubassetBytesIdx) : 0,
      bytesToInsert: isSet(object.bytesToInsert) ? String(object.bytesToInsert) : "",
      insertIdIdx: isSet(object.insertIdIdx) ? Number(object.insertIdIdx) : 0,
    };
  },

  toJSON(message: UriObject): unknown {
    const obj: any = {};
    message.decodeScheme !== undefined && (obj.decodeScheme = Math.round(message.decodeScheme));
    message.scheme !== undefined && (obj.scheme = Math.round(message.scheme));
    message.uri !== undefined && (obj.uri = message.uri);
    message.idxRangeToRemove !== undefined
      && (obj.idxRangeToRemove = message.idxRangeToRemove ? IdRange.toJSON(message.idxRangeToRemove) : undefined);
    message.insertSubassetBytesIdx !== undefined
      && (obj.insertSubassetBytesIdx = Math.round(message.insertSubassetBytesIdx));
    message.bytesToInsert !== undefined && (obj.bytesToInsert = message.bytesToInsert);
    message.insertIdIdx !== undefined && (obj.insertIdIdx = Math.round(message.insertIdIdx));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UriObject>, I>>(object: I): UriObject {
    const message = createBaseUriObject();
    message.decodeScheme = object.decodeScheme ?? 0;
    message.scheme = object.scheme ?? 0;
    message.uri = object.uri ?? "";
    message.idxRangeToRemove = (object.idxRangeToRemove !== undefined && object.idxRangeToRemove !== null)
      ? IdRange.fromPartial(object.idxRangeToRemove)
      : undefined;
    message.insertSubassetBytesIdx = object.insertSubassetBytesIdx ?? 0;
    message.bytesToInsert = object.bytesToInsert ?? "";
    message.insertIdIdx = object.insertIdIdx ?? 0;
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
