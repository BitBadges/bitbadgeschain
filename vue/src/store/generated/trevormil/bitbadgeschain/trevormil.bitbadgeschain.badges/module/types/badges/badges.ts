/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { UriObject } from "../badges/uris";
import { IdRange, BalanceObject } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** BitBadge defines a badge type. Think of this like the smart contract definition. */
export interface BitBadge {
  /**
   * id defines the unique identifier of the Badge classification, similar to the contract address of ERC721
   * starts at 0 and increments by 1 each badge
   */
  id: number;
  /**
   * uri object for the badge uri and subasset uris stored off chain. Stored in a special UriObject that attemtps to save space and avoid reused plaintext storage such as http:// and duplicate text for uri and subasset uris
   * data returned should corresponds to the Badge standard defined.
   */
  uri: UriObject | undefined;
  /**
   * these bytes can be used to store anything on-chain about the badge. This can be updatable or not depending on the permissions set.
   * Max 256 bytes allowed
   */
  arbitraryBytes: Uint8Array;
  /** manager address of the class; can have special permissions; is used as the reserve address for the assets */
  manager: number;
  /** Store permissions packed in a uint where the bits correspond to permissions from left to right; leading zeroes are applied and any future additions will be appended to the right. See types/permissions.go */
  permissions: number;
  /** FreezeRanges defines what addresses are frozen or unfrozen. If permissions.FrozenByDefault is false, this is used for frozen addresses. If true, this is used for unfrozen addresses. */
  freezeRanges: IdRange[];
  /** Starts at 0. Each subasset created will incrementally have an increasing ID #. Can't overflow. */
  nextSubassetId: number;
  /** Subasset supplys are stored if the subasset supply != default. Balance => SubbadgeIdRange map */
  subassetSupplys: BalanceObject[];
  /** Default subasset supply. If == 0, we assume == 1. */
  defaultSubassetSupply: number;
  /** Defines what standard this badge should implement. Must obey the rules of that standard. */
  standard: number;
}

const baseBitBadge: object = {
  id: 0,
  manager: 0,
  permissions: 0,
  nextSubassetId: 0,
  defaultSubassetSupply: 0,
  standard: 0,
};

export const BitBadge = {
  encode(message: BitBadge, writer: Writer = Writer.create()): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.uri !== undefined) {
      UriObject.encode(message.uri, writer.uint32(18).fork()).ldelim();
    }
    if (message.arbitraryBytes.length !== 0) {
      writer.uint32(26).bytes(message.arbitraryBytes);
    }
    if (message.manager !== 0) {
      writer.uint32(32).uint64(message.manager);
    }
    if (message.permissions !== 0) {
      writer.uint32(40).uint64(message.permissions);
    }
    for (const v of message.freezeRanges) {
      IdRange.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    if (message.nextSubassetId !== 0) {
      writer.uint32(96).uint64(message.nextSubassetId);
    }
    for (const v of message.subassetSupplys) {
      BalanceObject.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    if (message.defaultSubassetSupply !== 0) {
      writer.uint32(112).uint64(message.defaultSubassetSupply);
    }
    if (message.standard !== 0) {
      writer.uint32(120).uint64(message.standard);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BitBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBitBadge } as BitBadge;
    message.freezeRanges = [];
    message.subassetSupplys = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.uri = UriObject.decode(reader, reader.uint32());
          break;
        case 3:
          message.arbitraryBytes = reader.bytes();
          break;
        case 4:
          message.manager = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.permissions = longToNumber(reader.uint64() as Long);
          break;
        case 10:
          message.freezeRanges.push(IdRange.decode(reader, reader.uint32()));
          break;
        case 12:
          message.nextSubassetId = longToNumber(reader.uint64() as Long);
          break;
        case 13:
          message.subassetSupplys.push(
            BalanceObject.decode(reader, reader.uint32())
          );
          break;
        case 14:
          message.defaultSubassetSupply = longToNumber(reader.uint64() as Long);
          break;
        case 15:
          message.standard = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BitBadge {
    const message = { ...baseBitBadge } as BitBadge;
    message.freezeRanges = [];
    message.subassetSupplys = [];
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = UriObject.fromJSON(object.uri);
    } else {
      message.uri = undefined;
    }
    if (object.arbitraryBytes !== undefined && object.arbitraryBytes !== null) {
      message.arbitraryBytes = bytesFromBase64(object.arbitraryBytes);
    }
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = Number(object.manager);
    } else {
      message.manager = 0;
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = Number(object.permissions);
    } else {
      message.permissions = 0;
    }
    if (object.freezeRanges !== undefined && object.freezeRanges !== null) {
      for (const e of object.freezeRanges) {
        message.freezeRanges.push(IdRange.fromJSON(e));
      }
    }
    if (object.nextSubassetId !== undefined && object.nextSubassetId !== null) {
      message.nextSubassetId = Number(object.nextSubassetId);
    } else {
      message.nextSubassetId = 0;
    }
    if (
      object.subassetSupplys !== undefined &&
      object.subassetSupplys !== null
    ) {
      for (const e of object.subassetSupplys) {
        message.subassetSupplys.push(BalanceObject.fromJSON(e));
      }
    }
    if (
      object.defaultSubassetSupply !== undefined &&
      object.defaultSubassetSupply !== null
    ) {
      message.defaultSubassetSupply = Number(object.defaultSubassetSupply);
    } else {
      message.defaultSubassetSupply = 0;
    }
    if (object.standard !== undefined && object.standard !== null) {
      message.standard = Number(object.standard);
    } else {
      message.standard = 0;
    }
    return message;
  },

  toJSON(message: BitBadge): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.uri !== undefined &&
      (obj.uri = message.uri ? UriObject.toJSON(message.uri) : undefined);
    message.arbitraryBytes !== undefined &&
      (obj.arbitraryBytes = base64FromBytes(
        message.arbitraryBytes !== undefined
          ? message.arbitraryBytes
          : new Uint8Array()
      ));
    message.manager !== undefined && (obj.manager = message.manager);
    message.permissions !== undefined &&
      (obj.permissions = message.permissions);
    if (message.freezeRanges) {
      obj.freezeRanges = message.freezeRanges.map((e) =>
        e ? IdRange.toJSON(e) : undefined
      );
    } else {
      obj.freezeRanges = [];
    }
    message.nextSubassetId !== undefined &&
      (obj.nextSubassetId = message.nextSubassetId);
    if (message.subassetSupplys) {
      obj.subassetSupplys = message.subassetSupplys.map((e) =>
        e ? BalanceObject.toJSON(e) : undefined
      );
    } else {
      obj.subassetSupplys = [];
    }
    message.defaultSubassetSupply !== undefined &&
      (obj.defaultSubassetSupply = message.defaultSubassetSupply);
    message.standard !== undefined && (obj.standard = message.standard);
    return obj;
  },

  fromPartial(object: DeepPartial<BitBadge>): BitBadge {
    const message = { ...baseBitBadge } as BitBadge;
    message.freezeRanges = [];
    message.subassetSupplys = [];
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = UriObject.fromPartial(object.uri);
    } else {
      message.uri = undefined;
    }
    if (object.arbitraryBytes !== undefined && object.arbitraryBytes !== null) {
      message.arbitraryBytes = object.arbitraryBytes;
    } else {
      message.arbitraryBytes = new Uint8Array();
    }
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = object.manager;
    } else {
      message.manager = 0;
    }
    if (object.permissions !== undefined && object.permissions !== null) {
      message.permissions = object.permissions;
    } else {
      message.permissions = 0;
    }
    if (object.freezeRanges !== undefined && object.freezeRanges !== null) {
      for (const e of object.freezeRanges) {
        message.freezeRanges.push(IdRange.fromPartial(e));
      }
    }
    if (object.nextSubassetId !== undefined && object.nextSubassetId !== null) {
      message.nextSubassetId = object.nextSubassetId;
    } else {
      message.nextSubassetId = 0;
    }
    if (
      object.subassetSupplys !== undefined &&
      object.subassetSupplys !== null
    ) {
      for (const e of object.subassetSupplys) {
        message.subassetSupplys.push(BalanceObject.fromPartial(e));
      }
    }
    if (
      object.defaultSubassetSupply !== undefined &&
      object.defaultSubassetSupply !== null
    ) {
      message.defaultSubassetSupply = object.defaultSubassetSupply;
    } else {
      message.defaultSubassetSupply = 0;
    }
    if (object.standard !== undefined && object.standard !== null) {
      message.standard = object.standard;
    } else {
      message.standard = 0;
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
