/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { BitBadge } from "./badges";
import { UserBalanceInfo } from "./balances";
import { Params } from "./params";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** GenesisState defines the badges module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  portId: string;
  badges: BitBadge[];
  balances: UserBalanceInfo[];
  balanceIds: string[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  nextBadgeId: number;
}

function createBaseGenesisState(): GenesisState {
  return { params: undefined, portId: "", badges: [], balances: [], balanceIds: [], nextBadgeId: 0 };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    if (message.portId !== "") {
      writer.uint32(18).string(message.portId);
    }
    for (const v of message.badges) {
      BitBadge.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.balances) {
      UserBalanceInfo.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.balanceIds) {
      writer.uint32(42).string(v!);
    }
    if (message.nextBadgeId !== 0) {
      writer.uint32(48).uint64(message.nextBadgeId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.portId = reader.string();
          break;
        case 3:
          message.badges.push(BitBadge.decode(reader, reader.uint32()));
          break;
        case 4:
          message.balances.push(UserBalanceInfo.decode(reader, reader.uint32()));
          break;
        case 5:
          message.balanceIds.push(reader.string());
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
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      portId: isSet(object.portId) ? String(object.portId) : "",
      badges: Array.isArray(object?.badges) ? object.badges.map((e: any) => BitBadge.fromJSON(e)) : [],
      balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => UserBalanceInfo.fromJSON(e)) : [],
      balanceIds: Array.isArray(object?.balanceIds) ? object.balanceIds.map((e: any) => String(e)) : [],
      nextBadgeId: isSet(object.nextBadgeId) ? Number(object.nextBadgeId) : 0,
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    message.portId !== undefined && (obj.portId = message.portId);
    if (message.badges) {
      obj.badges = message.badges.map((e) => e ? BitBadge.toJSON(e) : undefined);
    } else {
      obj.badges = [];
    }
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? UserBalanceInfo.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    if (message.balanceIds) {
      obj.balanceIds = message.balanceIds.map((e) => e);
    } else {
      obj.balanceIds = [];
    }
    message.nextBadgeId !== undefined && (obj.nextBadgeId = Math.round(message.nextBadgeId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.portId = object.portId ?? "";
    message.badges = object.badges?.map((e) => BitBadge.fromPartial(e)) || [];
    message.balances = object.balances?.map((e) => UserBalanceInfo.fromPartial(e)) || [];
    message.balanceIds = object.balanceIds?.map((e) => e) || [];
    message.nextBadgeId = object.nextBadgeId ?? 0;
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
