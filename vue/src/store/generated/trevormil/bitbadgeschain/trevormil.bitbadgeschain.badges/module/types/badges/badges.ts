/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** BitBadge defines a badge type. Think of this like the smart contract definition */
export interface BitBadge {
  /**
   * id defines the unique identifier of the Badge classification, similar to the contract address of ERC721
   * starts at 0 and increments by 1 each badge
   */
  id: number;
  /** uri for the class metadata stored off chain. must match a valid metadata standard (bitbadge, collection, etc) */
  uri: string;
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
  permission_flags: number;
  /**
   * if frozen_by_default is true, this is an accumulator of unfrozen addresses; and vice versa for false
   * big.Int will always only be 32 uint64s long
   */
  freeze_addresses: number[];
  /**
   * uri for the subassets metadata stored off chain; include {id} in the string, it will be replaced with the subasset id
   * if not specified, uses a default Class (ID # 1) like metadata
   */
  subasset_uri_format: string;
  /** starts at 0; each subasset created will incrementally have an increasing ID # */
  next_subasset_id: number;
  /** only store if not 1 (default); will be sorted in order of subsasset ids; (maybe add defaut option in future) */
  subassets_total_supply: Subasset[];
}

/** Only will be created if supply >= 0 */
export interface Subasset {
  startId: number;
  endId: number;
  supply: number;
}

const baseBitBadge: object = {
  id: 0,
  uri: "",
  manager: 0,
  permission_flags: 0,
  freeze_addresses: 0,
  subasset_uri_format: "",
  next_subasset_id: 0,
};

export const BitBadge = {
  encode(message: BitBadge, writer: Writer = Writer.create()): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.uri !== "") {
      writer.uint32(18).string(message.uri);
    }
    if (message.manager !== 0) {
      writer.uint32(32).uint64(message.manager);
    }
    if (message.permission_flags !== 0) {
      writer.uint32(40).uint64(message.permission_flags);
    }
    writer.uint32(82).fork();
    for (const v of message.freeze_addresses) {
      writer.uint64(v);
    }
    writer.ldelim();
    if (message.subasset_uri_format !== "") {
      writer.uint32(90).string(message.subasset_uri_format);
    }
    if (message.next_subasset_id !== 0) {
      writer.uint32(96).uint64(message.next_subasset_id);
    }
    for (const v of message.subassets_total_supply) {
      Subasset.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BitBadge {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBitBadge } as BitBadge;
    message.freeze_addresses = [];
    message.subassets_total_supply = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.uri = reader.string();
          break;
        case 4:
          message.manager = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.permission_flags = longToNumber(reader.uint64() as Long);
          break;
        case 10:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.freeze_addresses.push(
                longToNumber(reader.uint64() as Long)
              );
            }
          } else {
            message.freeze_addresses.push(
              longToNumber(reader.uint64() as Long)
            );
          }
          break;
        case 11:
          message.subasset_uri_format = reader.string();
          break;
        case 12:
          message.next_subasset_id = longToNumber(reader.uint64() as Long);
          break;
        case 13:
          message.subassets_total_supply.push(
            Subasset.decode(reader, reader.uint32())
          );
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
    message.freeze_addresses = [];
    message.subassets_total_supply = [];
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
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = Number(object.manager);
    } else {
      message.manager = 0;
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
      object.freeze_addresses !== undefined &&
      object.freeze_addresses !== null
    ) {
      for (const e of object.freeze_addresses) {
        message.freeze_addresses.push(Number(e));
      }
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
      for (const e of object.subassets_total_supply) {
        message.subassets_total_supply.push(Subasset.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: BitBadge): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.uri !== undefined && (obj.uri = message.uri);
    message.manager !== undefined && (obj.manager = message.manager);
    message.permission_flags !== undefined &&
      (obj.permission_flags = message.permission_flags);
    if (message.freeze_addresses) {
      obj.freeze_addresses = message.freeze_addresses.map((e) => e);
    } else {
      obj.freeze_addresses = [];
    }
    message.subasset_uri_format !== undefined &&
      (obj.subasset_uri_format = message.subasset_uri_format);
    message.next_subasset_id !== undefined &&
      (obj.next_subasset_id = message.next_subasset_id);
    if (message.subassets_total_supply) {
      obj.subassets_total_supply = message.subassets_total_supply.map((e) =>
        e ? Subasset.toJSON(e) : undefined
      );
    } else {
      obj.subassets_total_supply = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BitBadge>): BitBadge {
    const message = { ...baseBitBadge } as BitBadge;
    message.freeze_addresses = [];
    message.subassets_total_supply = [];
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
    if (object.manager !== undefined && object.manager !== null) {
      message.manager = object.manager;
    } else {
      message.manager = 0;
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
      object.freeze_addresses !== undefined &&
      object.freeze_addresses !== null
    ) {
      for (const e of object.freeze_addresses) {
        message.freeze_addresses.push(e);
      }
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
      for (const e of object.subassets_total_supply) {
        message.subassets_total_supply.push(Subasset.fromPartial(e));
      }
    }
    return message;
  },
};

const baseSubasset: object = { startId: 0, endId: 0, supply: 0 };

export const Subasset = {
  encode(message: Subasset, writer: Writer = Writer.create()): Writer {
    if (message.startId !== 0) {
      writer.uint32(8).uint64(message.startId);
    }
    if (message.endId !== 0) {
      writer.uint32(16).uint64(message.endId);
    }
    if (message.supply !== 0) {
      writer.uint32(24).uint64(message.supply);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Subasset {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSubasset } as Subasset;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.startId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.endId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.supply = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Subasset {
    const message = { ...baseSubasset } as Subasset;
    if (object.startId !== undefined && object.startId !== null) {
      message.startId = Number(object.startId);
    } else {
      message.startId = 0;
    }
    if (object.endId !== undefined && object.endId !== null) {
      message.endId = Number(object.endId);
    } else {
      message.endId = 0;
    }
    if (object.supply !== undefined && object.supply !== null) {
      message.supply = Number(object.supply);
    } else {
      message.supply = 0;
    }
    return message;
  },

  toJSON(message: Subasset): unknown {
    const obj: any = {};
    message.startId !== undefined && (obj.startId = message.startId);
    message.endId !== undefined && (obj.endId = message.endId);
    message.supply !== undefined && (obj.supply = message.supply);
    return obj;
  },

  fromPartial(object: DeepPartial<Subasset>): Subasset {
    const message = { ...baseSubasset } as Subasset;
    if (object.startId !== undefined && object.startId !== null) {
      message.startId = object.startId;
    } else {
      message.startId = 0;
    }
    if (object.endId !== undefined && object.endId !== null) {
      message.endId = object.endId;
    } else {
      message.endId = 0;
    }
    if (object.supply !== undefined && object.supply !== null) {
      message.supply = object.supply;
    } else {
      message.supply = 0;
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
