/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { BalanceObject, IdRange } from "./ranges";
import { UriObject } from "./uris";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

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
  uri:
    | UriObject
    | undefined;
  /**
   * these bytes can be used to store anything on-chain about the badge. This can be updatable or not depending on the permissions set.
   * Max 256 bytes allowed
   */
  arbitraryBytes: string;
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

function createBaseBitBadge(): BitBadge {
  return {
    id: 0,
    uri: undefined,
    arbitraryBytes: "",
    manager: 0,
    permissions: 0,
    freezeRanges: [],
    nextSubassetId: 0,
    subassetSupplys: [],
    defaultSubassetSupply: 0,
    standard: 0,
  };
}

export const BitBadge = {
  encode(message: BitBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.uri !== undefined) {
      UriObject.encode(message.uri, writer.uint32(18).fork()).ldelim();
    }
    if (message.arbitraryBytes !== "") {
      writer.uint32(26).string(message.arbitraryBytes);
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

  decode(input: _m0.Reader | Uint8Array, length?: number): BitBadge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBitBadge();
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
          message.arbitraryBytes = reader.string();
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
          message.subassetSupplys.push(BalanceObject.decode(reader, reader.uint32()));
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
    return {
      id: isSet(object.id) ? Number(object.id) : 0,
      uri: isSet(object.uri) ? UriObject.fromJSON(object.uri) : undefined,
      arbitraryBytes: isSet(object.arbitraryBytes) ? String(object.arbitraryBytes) : "",
      manager: isSet(object.manager) ? Number(object.manager) : 0,
      permissions: isSet(object.permissions) ? Number(object.permissions) : 0,
      freezeRanges: Array.isArray(object?.freezeRanges) ? object.freezeRanges.map((e: any) => IdRange.fromJSON(e)) : [],
      nextSubassetId: isSet(object.nextSubassetId) ? Number(object.nextSubassetId) : 0,
      subassetSupplys: Array.isArray(object?.subassetSupplys)
        ? object.subassetSupplys.map((e: any) => BalanceObject.fromJSON(e))
        : [],
      defaultSubassetSupply: isSet(object.defaultSubassetSupply) ? Number(object.defaultSubassetSupply) : 0,
      standard: isSet(object.standard) ? Number(object.standard) : 0,
    };
  },

  toJSON(message: BitBadge): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    message.uri !== undefined && (obj.uri = message.uri ? UriObject.toJSON(message.uri) : undefined);
    message.arbitraryBytes !== undefined && (obj.arbitraryBytes = message.arbitraryBytes);
    message.manager !== undefined && (obj.manager = Math.round(message.manager));
    message.permissions !== undefined && (obj.permissions = Math.round(message.permissions));
    if (message.freezeRanges) {
      obj.freezeRanges = message.freezeRanges.map((e) => e ? IdRange.toJSON(e) : undefined);
    } else {
      obj.freezeRanges = [];
    }
    message.nextSubassetId !== undefined && (obj.nextSubassetId = Math.round(message.nextSubassetId));
    if (message.subassetSupplys) {
      obj.subassetSupplys = message.subassetSupplys.map((e) => e ? BalanceObject.toJSON(e) : undefined);
    } else {
      obj.subassetSupplys = [];
    }
    message.defaultSubassetSupply !== undefined
      && (obj.defaultSubassetSupply = Math.round(message.defaultSubassetSupply));
    message.standard !== undefined && (obj.standard = Math.round(message.standard));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BitBadge>, I>>(object: I): BitBadge {
    const message = createBaseBitBadge();
    message.id = object.id ?? 0;
    message.uri = (object.uri !== undefined && object.uri !== null) ? UriObject.fromPartial(object.uri) : undefined;
    message.arbitraryBytes = object.arbitraryBytes ?? "";
    message.manager = object.manager ?? 0;
    message.permissions = object.permissions ?? 0;
    message.freezeRanges = object.freezeRanges?.map((e) => IdRange.fromPartial(e)) || [];
    message.nextSubassetId = object.nextSubassetId ?? 0;
    message.subassetSupplys = object.subassetSupplys?.map((e) => BalanceObject.fromPartial(e)) || [];
    message.defaultSubassetSupply = object.defaultSubassetSupply ?? 0;
    message.standard = object.standard ?? 0;
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
