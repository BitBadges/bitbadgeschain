/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { IdRange, BalanceObject } from "../badges/ranges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** BitBadge defines a badge type. Think of this like the smart contract definition. */
export interface BitBadge {
  /**
   * id defines the unique identifier of the Badge classification, similar to the contract address of ERC721
   * starts at 0 and increments by 1 each badge
   */
  id: number;
  /** uri for the class metadata stored off chain. must match a valid metadata standard (bitbadge, collection, etc) */
  uri: string;
  /**
   * these bytes can be used to store anything on-chain about the badge. This can be updatable or not depending on the permissions set.
   * Max 256 bytes allowed
   */
  arbitraryBytes: Uint8Array;
  /** manager address of the class; can have special permissions; is used as the reserve address for the assets */
  manager: number;
  /**
   * Flag bits are in the following order from left to right; leading zeroes are applied and any future additions will be appended to the right
   *
   * can_manager_transfer: can the manager transfer managerial privileges to another address
   * can_update_uris: can the manager update the uris of the class and subassets; if false, locked forever
   * forceful_transfers: if true, one can send a badge to an account without pending approval; these badges should not by default be displayed on public profiles (can also use collections)
   * can_create: when true, manager can create more subassets of the class; once set to false, it is locked
   * can_revoke: when true, manager can revoke subassets of the class (including null address); once set to false, it is locked
   * can_freeze: when true, manager can freeze addresseses from transferring; once set to false, it is locked
   * frozen_by_default: when true, all addresses are considered frozen and must be unfrozen to transfer; when false, all addresses are considered unfrozen and must be frozen to freeze
   * manager is not frozen by default
   *
   * More permissions to be added
   */
  permissions: number;
  /**
   * if frozen_by_default is true, this is an accumulator of unfrozen addresses; and vice versa for false
   * big.Int will always only be 32 uint64s long
   */
  freezeRanges: IdRange[];
  /**
   * uri for the subassets metadata stored off chain; include {id} in the string, it will be replaced with the subasset id
   * if not specified, uses a default Class (ID # 1) like metadata
   */
  subassetUriFormat: string;
  /** starts at 0; each subasset created will incrementally have an increasing ID # */
  nextSubassetId: number;
  /** only store if not == default; will be sorted in order of subsasset ids; (maybe add defaut option in future) */
  subassetSupplys: BalanceObject[];
  defaultSubassetSupply: number;
}

const baseBitBadge: object = {
  id: 0,
  uri: "",
  manager: 0,
  permissions: 0,
  subassetUriFormat: "",
  nextSubassetId: 0,
  defaultSubassetSupply: 0,
};

export const BitBadge = {
  encode(message: BitBadge, writer: Writer = Writer.create()): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.uri !== "") {
      writer.uint32(18).string(message.uri);
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
    if (message.subassetUriFormat !== "") {
      writer.uint32(90).string(message.subassetUriFormat);
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
          message.uri = reader.string();
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
        case 11:
          message.subassetUriFormat = reader.string();
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
      message.uri = String(object.uri);
    } else {
      message.uri = "";
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
    if (
      object.subassetUriFormat !== undefined &&
      object.subassetUriFormat !== null
    ) {
      message.subassetUriFormat = String(object.subassetUriFormat);
    } else {
      message.subassetUriFormat = "";
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
    return message;
  },

  toJSON(message: BitBadge): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.uri !== undefined && (obj.uri = message.uri);
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
    message.subassetUriFormat !== undefined &&
      (obj.subassetUriFormat = message.subassetUriFormat);
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
      message.uri = object.uri;
    } else {
      message.uri = "";
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
    if (
      object.subassetUriFormat !== undefined &&
      object.subassetUriFormat !== null
    ) {
      message.subassetUriFormat = object.subassetUriFormat;
    } else {
      message.subassetUriFormat = "";
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
