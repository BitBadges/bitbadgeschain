/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { IdRange } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** A URI object defines a uri and subasset uri for a badge and its subbadges. Designed to save storage and avoid reused text and common patterns. */
export interface UriObject {
  /** This will be == 0 represeting plaintext URLs for now, but in the future, if we want to add other decoding / encoding schemes like Base64, etc, we just define a new decoding scheme */
  decodeScheme: number;
  /** Helps to save space by not storing methods like https:// every time. If this == 0, we assume it's included in the uri itself. Else, we define a few predefined int -> scheme maps in uris.go that will prefix the uri bytes. */
  scheme: number;
  /** The uri bytes to store. Will be converted to a string. To be manipulated according to other properties of this URI object. */
  uri: Uint8Array;
  /** The four fields below are used to convert the uri from above to the subasset URI. */
  idxRangeToRemove: IdRange | undefined;
  insertSubassetBytesIdx: number;
  /** After removing the above range, insert these bytes at insertSubassetBytesIdx. insertSubassetBytesIdx is the idx after the range removal, not before. */
  bytesToInsert: Uint8Array;
  /** This is the idx where we insert the id number of the subbadge. For example, ex.com/ and insertIdIdx == 6, the subasset URI for ID 0 would be ex.com/0 */
  insertIdIdx: number;
}

const baseUriObject: object = {
  decodeScheme: 0,
  scheme: 0,
  insertSubassetBytesIdx: 0,
  insertIdIdx: 0,
};

export const UriObject = {
  encode(message: UriObject, writer: Writer = Writer.create()): Writer {
    if (message.decodeScheme !== 0) {
      writer.uint32(8).uint64(message.decodeScheme);
    }
    if (message.scheme !== 0) {
      writer.uint32(16).uint64(message.scheme);
    }
    if (message.uri.length !== 0) {
      writer.uint32(26).bytes(message.uri);
    }
    if (message.idxRangeToRemove !== undefined) {
      IdRange.encode(
        message.idxRangeToRemove,
        writer.uint32(34).fork()
      ).ldelim();
    }
    if (message.insertSubassetBytesIdx !== 0) {
      writer.uint32(40).uint64(message.insertSubassetBytesIdx);
    }
    if (message.bytesToInsert.length !== 0) {
      writer.uint32(50).bytes(message.bytesToInsert);
    }
    if (message.insertIdIdx !== 0) {
      writer.uint32(56).uint64(message.insertIdIdx);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UriObject {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUriObject } as UriObject;
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
          message.uri = reader.bytes();
          break;
        case 4:
          message.idxRangeToRemove = IdRange.decode(reader, reader.uint32());
          break;
        case 5:
          message.insertSubassetBytesIdx = longToNumber(
            reader.uint64() as Long
          );
          break;
        case 6:
          message.bytesToInsert = reader.bytes();
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
    const message = { ...baseUriObject } as UriObject;
    if (object.decodeScheme !== undefined && object.decodeScheme !== null) {
      message.decodeScheme = Number(object.decodeScheme);
    } else {
      message.decodeScheme = 0;
    }
    if (object.scheme !== undefined && object.scheme !== null) {
      message.scheme = Number(object.scheme);
    } else {
      message.scheme = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = bytesFromBase64(object.uri);
    }
    if (
      object.idxRangeToRemove !== undefined &&
      object.idxRangeToRemove !== null
    ) {
      message.idxRangeToRemove = IdRange.fromJSON(object.idxRangeToRemove);
    } else {
      message.idxRangeToRemove = undefined;
    }
    if (
      object.insertSubassetBytesIdx !== undefined &&
      object.insertSubassetBytesIdx !== null
    ) {
      message.insertSubassetBytesIdx = Number(object.insertSubassetBytesIdx);
    } else {
      message.insertSubassetBytesIdx = 0;
    }
    if (object.bytesToInsert !== undefined && object.bytesToInsert !== null) {
      message.bytesToInsert = bytesFromBase64(object.bytesToInsert);
    }
    if (object.insertIdIdx !== undefined && object.insertIdIdx !== null) {
      message.insertIdIdx = Number(object.insertIdIdx);
    } else {
      message.insertIdIdx = 0;
    }
    return message;
  },

  toJSON(message: UriObject): unknown {
    const obj: any = {};
    message.decodeScheme !== undefined &&
      (obj.decodeScheme = message.decodeScheme);
    message.scheme !== undefined && (obj.scheme = message.scheme);
    message.uri !== undefined &&
      (obj.uri = base64FromBytes(
        message.uri !== undefined ? message.uri : new Uint8Array()
      ));
    message.idxRangeToRemove !== undefined &&
      (obj.idxRangeToRemove = message.idxRangeToRemove
        ? IdRange.toJSON(message.idxRangeToRemove)
        : undefined);
    message.insertSubassetBytesIdx !== undefined &&
      (obj.insertSubassetBytesIdx = message.insertSubassetBytesIdx);
    message.bytesToInsert !== undefined &&
      (obj.bytesToInsert = base64FromBytes(
        message.bytesToInsert !== undefined
          ? message.bytesToInsert
          : new Uint8Array()
      ));
    message.insertIdIdx !== undefined &&
      (obj.insertIdIdx = message.insertIdIdx);
    return obj;
  },

  fromPartial(object: DeepPartial<UriObject>): UriObject {
    const message = { ...baseUriObject } as UriObject;
    if (object.decodeScheme !== undefined && object.decodeScheme !== null) {
      message.decodeScheme = object.decodeScheme;
    } else {
      message.decodeScheme = 0;
    }
    if (object.scheme !== undefined && object.scheme !== null) {
      message.scheme = object.scheme;
    } else {
      message.scheme = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = object.uri;
    } else {
      message.uri = new Uint8Array();
    }
    if (
      object.idxRangeToRemove !== undefined &&
      object.idxRangeToRemove !== null
    ) {
      message.idxRangeToRemove = IdRange.fromPartial(object.idxRangeToRemove);
    } else {
      message.idxRangeToRemove = undefined;
    }
    if (
      object.insertSubassetBytesIdx !== undefined &&
      object.insertSubassetBytesIdx !== null
    ) {
      message.insertSubassetBytesIdx = object.insertSubassetBytesIdx;
    } else {
      message.insertSubassetBytesIdx = 0;
    }
    if (object.bytesToInsert !== undefined && object.bytesToInsert !== null) {
      message.bytesToInsert = object.bytesToInsert;
    } else {
      message.bytesToInsert = new Uint8Array();
    }
    if (object.insertIdIdx !== undefined && object.insertIdIdx !== null) {
      message.insertIdIdx = object.insertIdIdx;
    } else {
      message.insertIdIdx = 0;
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

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (let i = 0; i < arr.byteLength; ++i) {
    bin.push(String.fromCharCode(arr[i]));
  }
  return btoa(bin.join(""));
}

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
