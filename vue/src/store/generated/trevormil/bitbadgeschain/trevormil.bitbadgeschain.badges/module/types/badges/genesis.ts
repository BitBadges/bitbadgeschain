/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Params } from "../badges/params";
import { BitBadge } from "../badges/badges";
import { UserBalanceInfo } from "../badges/balances";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** GenesisState defines the badges module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  port_id: string;
  badges: BitBadge[];
  balances: UserBalanceInfo[];
  balance_ids: string[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  nextBadgeId: number;
}

const baseGenesisState: object = {
  port_id: "",
  balance_ids: "",
  nextBadgeId: 0,
};

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    if (message.port_id !== "") {
      writer.uint32(18).string(message.port_id);
    }
    for (const v of message.badges) {
      BitBadge.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.balances) {
      UserBalanceInfo.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.balance_ids) {
      writer.uint32(42).string(v!);
    }
    if (message.nextBadgeId !== 0) {
      writer.uint32(48).uint64(message.nextBadgeId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = [];
    message.balances = [];
    message.balance_ids = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.port_id = reader.string();
          break;
        case 3:
          message.badges.push(BitBadge.decode(reader, reader.uint32()));
          break;
        case 4:
          message.balances.push(
            UserBalanceInfo.decode(reader, reader.uint32())
          );
          break;
        case 5:
          message.balance_ids.push(reader.string());
          break;
        case 6:
          message.nextBadgeId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = [];
    message.balances = [];
    message.balance_ids = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    if (object.port_id !== undefined && object.port_id !== null) {
      message.port_id = String(object.port_id);
    } else {
      message.port_id = "";
    }
    if (object.badges !== undefined && object.badges !== null) {
      for (const e of object.badges) {
        message.badges.push(BitBadge.fromJSON(e));
      }
    }
    if (object.balances !== undefined && object.balances !== null) {
      for (const e of object.balances) {
        message.balances.push(UserBalanceInfo.fromJSON(e));
      }
    }
    if (object.balance_ids !== undefined && object.balance_ids !== null) {
      for (const e of object.balance_ids) {
        message.balance_ids.push(String(e));
      }
    }
    if (object.nextBadgeId !== undefined && object.nextBadgeId !== null) {
      message.nextBadgeId = Number(object.nextBadgeId);
    } else {
      message.nextBadgeId = 0;
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    message.port_id !== undefined && (obj.port_id = message.port_id);
    if (message.badges) {
      obj.badges = message.badges.map((e) =>
        e ? BitBadge.toJSON(e) : undefined
      );
    } else {
      obj.badges = [];
    }
    if (message.balances) {
      obj.balances = message.balances.map((e) =>
        e ? UserBalanceInfo.toJSON(e) : undefined
      );
    } else {
      obj.balances = [];
    }
    if (message.balance_ids) {
      obj.balance_ids = message.balance_ids.map((e) => e);
    } else {
      obj.balance_ids = [];
    }
    message.nextBadgeId !== undefined &&
      (obj.nextBadgeId = message.nextBadgeId);
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = [];
    message.balances = [];
    message.balance_ids = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    if (object.port_id !== undefined && object.port_id !== null) {
      message.port_id = object.port_id;
    } else {
      message.port_id = "";
    }
    if (object.badges !== undefined && object.badges !== null) {
      for (const e of object.badges) {
        message.badges.push(BitBadge.fromPartial(e));
      }
    }
    if (object.balances !== undefined && object.balances !== null) {
      for (const e of object.balances) {
        message.balances.push(UserBalanceInfo.fromPartial(e));
      }
    }
    if (object.balance_ids !== undefined && object.balance_ids !== null) {
      for (const e of object.balance_ids) {
        message.balance_ids.push(e);
      }
    }
    if (object.nextBadgeId !== undefined && object.nextBadgeId !== null) {
      message.nextBadgeId = object.nextBadgeId;
    } else {
      message.nextBadgeId = 0;
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
