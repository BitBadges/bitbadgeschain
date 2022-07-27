/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** BitBadge defines a badge type. Think of this like the smart contract definition */
export interface BitBadge {
  /** id defines the unique identifier of the Badge classification, similar to the contract address of ERC721 */
  id: string;
  /** uri for the class metadata stored off chain. must match a valid metadata standard (bitbadge, collection, etc) */
  uri: string;
  /** inital creator address of the class */
  creator: string;
  /** manager addressof the class; can have special permissions; is used as the reserve address for the assets */
  manager: string;
  /**
   * Flag bits are in the following order from left to right; leading zeroes are applied and any future additions will be appended to the right
   *
   * can_update_uris: can the manager update the uris of the class and subassets; if false, locked forever
   * forceful_transfers: if true, one can send a badge to an account without pending approval; these badges should not by default be displayed on public profiles (can also use collections)
   * can_create: when true, manager can create more subassets of the class; once set to false, it is locked
   * can_revoke: when true, manager can revoke subassets of the class (including null address); once set to false, it is locked
   * can_freeze: when true, manager can freeze addresseses from transferring; once set to false, it is locked
   * frozen_by_default: when true, all addresses are considered frozen and must be unfrozen to transfer; when false, all addresses are considered unfrozen and must be frozen to freeze
   */
  permission_flags: number;
  /**
   * if frozen_by_default is true, this is a list of unfrozen addresses; and vice versa for false
   * TODO: make this a fixed length set efficient accumulator (no need to store a list of all addresses; just lookup membership)
   * TODO: set max length
   */
  frozen_or_unfrozen_addresses_digest: string;
  /**
   * uri for the subassets metadata stored off chain; include {id} in the string, it will be replaced with the subasset id
   * if not specified, uses a default Class (ID # 1) like metadata
   */
  subasset_uri_format: string;
  /** starts at 0; each subasset created will incrementally have an increasing ID # */
  next_subasset_id: number;
  /** subasset id => total supply map; only store if not 1 (default) */
  subassets_total_supply: { [key: number]: number };
}

export interface BitBadge_SubassetsTotalSupplyEntry {
  key: number;
  value: number;
}

const baseBitBadge: object = {
  id: "",
  uri: "",
  creator: "",
  manager: "",
  permission_flags: 0,
  frozen_or_unfrozen_addresses_digest: "",
  subasset_uri_format: "",
  next_subasset_id: 0,
};

export const BitBadge = {
  encode(message: BitBadge, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.uri !== "") {
      writer.uint32(18).string(message.uri);
    }
    if (message.creator !== "") {
      writer.uint32(26).string(message.creator);
    }
    if (message.manager !== "") {
      writer.uint32(34).string(message.manager);
    }
    if (message.permission_flags !== 0) {
      writer.uint32(40).uint64(message.permission_flags);
    }
    if (message.frozen_or_unfrozen_addresses_digest !== "") {
      writer.uint32(82).string(message.frozen_or_unfrozen_addresses_digest);
    }
    if (message.subasset_uri_format !== "") {
      writer.uint32(90).string(message.subasset_uri_format);
    }
    if (message.next_subasset_id !== 0) {
      writer.uint32(96).uint64(message.next_subasset_id);
    }
    Object.entries(message.subassets_total_supply).forEach(([key, value]) => {
      BitBadge_SubassetsTotalSupplyEntry.encode(
        { key: key as any, value },
        writer.uint32(106).fork()
      ).ldelim();
    });
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BitBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBitBadge } as BitBadge;
    message.subassets_total_supply = {};
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.uri = reader.string();
          break;
        case 3:
          message.creator = reader.string();
          break;
        case 4:
          message.manager = reader.string();
          break;
        case 5:
          message.permission_flags = longToNumber(reader.uint64() as Long);
          break;
        case 10:
          message.frozen_or_unfrozen_addresses_digest = reader.string();
          break;
        case 11:
          message.subasset_uri_format = reader.string();
          break;
        case 12:
          message.next_subasset_id = longToNumber(reader.uint64() as Long);
          break;
        case 13:
          const entry13 = BitBadge_SubassetsTotalSupplyEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry13.value !== undefined) {
            message.subassets_total_supply[entry13.key] = entry13.value;
          }
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
    message.subassets_total_supply = {};
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = String(object.uri);
    } else {
      message.uri = "";
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = String(object.manager);
    } else {
      message.manager = "";
    }
    if (
      object.permission_flags !== undefined &&
      object.permission_flags !== null
    ) {
      message.permission_flags = Number(object.permission_flags);
    } else {
      message.permission_flags = 0;
    }
    if (
      object.frozen_or_unfrozen_addresses_digest !== undefined &&
      object.frozen_or_unfrozen_addresses_digest !== null
    ) {
      message.frozen_or_unfrozen_addresses_digest = String(
        object.frozen_or_unfrozen_addresses_digest
      );
    } else {
      message.frozen_or_unfrozen_addresses_digest = "";
    }
    if (
      object.subasset_uri_format !== undefined &&
      object.subasset_uri_format !== null
    ) {
      message.subasset_uri_format = String(object.subasset_uri_format);
    } else {
      message.subasset_uri_format = "";
    }
    if (
      object.next_subasset_id !== undefined &&
      object.next_subasset_id !== null
    ) {
      message.next_subasset_id = Number(object.next_subasset_id);
    } else {
      message.next_subasset_id = 0;
    }
    if (
      object.subassets_total_supply !== undefined &&
      object.subassets_total_supply !== null
    ) {
      Object.entries(object.subassets_total_supply).forEach(([key, value]) => {
        message.subassets_total_supply[Number(key)] = Number(value);
      });
    }
    return message;
  },

  toJSON(message: BitBadge): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.uri !== undefined && (obj.uri = message.uri);
    message.creator !== undefined && (obj.creator = message.creator);
    message.manager !== undefined && (obj.manager = message.manager);
    message.permission_flags !== undefined &&
      (obj.permission_flags = message.permission_flags);
    message.frozen_or_unfrozen_addresses_digest !== undefined &&
      (obj.frozen_or_unfrozen_addresses_digest =
        message.frozen_or_unfrozen_addresses_digest);
    message.subasset_uri_format !== undefined &&
      (obj.subasset_uri_format = message.subasset_uri_format);
    message.next_subasset_id !== undefined &&
      (obj.next_subasset_id = message.next_subasset_id);
    obj.subassets_total_supply = {};
    if (message.subassets_total_supply) {
      Object.entries(message.subassets_total_supply).forEach(([k, v]) => {
        obj.subassets_total_supply[k] = v;
      });
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BitBadge>): BitBadge {
    const message = { ...baseBitBadge } as BitBadge;
    message.subassets_total_supply = {};
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.uri !== undefined && object.uri !== null) {
      message.uri = object.uri;
    } else {
      message.uri = "";
    }
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = object.manager;
    } else {
      message.manager = "";
    }
    if (
      object.permission_flags !== undefined &&
      object.permission_flags !== null
    ) {
      message.permission_flags = object.permission_flags;
    } else {
      message.permission_flags = 0;
    }
    if (
      object.frozen_or_unfrozen_addresses_digest !== undefined &&
      object.frozen_or_unfrozen_addresses_digest !== null
    ) {
      message.frozen_or_unfrozen_addresses_digest =
        object.frozen_or_unfrozen_addresses_digest;
    } else {
      message.frozen_or_unfrozen_addresses_digest = "";
    }
    if (
      object.subasset_uri_format !== undefined &&
      object.subasset_uri_format !== null
    ) {
      message.subasset_uri_format = object.subasset_uri_format;
    } else {
      message.subasset_uri_format = "";
    }
    if (
      object.next_subasset_id !== undefined &&
      object.next_subasset_id !== null
    ) {
      message.next_subasset_id = object.next_subasset_id;
    } else {
      message.next_subasset_id = 0;
    }
    if (
      object.subassets_total_supply !== undefined &&
      object.subassets_total_supply !== null
    ) {
      Object.entries(object.subassets_total_supply).forEach(([key, value]) => {
        if (value !== undefined) {
          message.subassets_total_supply[Number(key)] = Number(value);
        }
      });
    }
    return message;
  },
};

const baseBitBadge_SubassetsTotalSupplyEntry: object = { key: 0, value: 0 };

export const BitBadge_SubassetsTotalSupplyEntry = {
  encode(
    message: BitBadge_SubassetsTotalSupplyEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint64(message.key);
    }
    if (message.value !== 0) {
      writer.uint32(16).uint64(message.value);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): BitBadge_SubassetsTotalSupplyEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBitBadge_SubassetsTotalSupplyEntry,
    } as BitBadge_SubassetsTotalSupplyEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.value = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BitBadge_SubassetsTotalSupplyEntry {
    const message = {
      ...baseBitBadge_SubassetsTotalSupplyEntry,
    } as BitBadge_SubassetsTotalSupplyEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = Number(object.key);
    } else {
      message.key = 0;
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = Number(object.value);
    } else {
      message.value = 0;
    }
    return message;
  },

  toJSON(message: BitBadge_SubassetsTotalSupplyEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial(
    object: DeepPartial<BitBadge_SubassetsTotalSupplyEntry>
  ): BitBadge_SubassetsTotalSupplyEntry {
    const message = {
      ...baseBitBadge_SubassetsTotalSupplyEntry,
    } as BitBadge_SubassetsTotalSupplyEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = 0;
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = object.value;
    } else {
      message.value = 0;
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
