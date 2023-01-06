/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { BitBadge } from "./badges";
import { UserBalanceInfo } from "./balances";
import { Params } from "./params";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetBadgeRequest {
  id: number;
}

export interface QueryGetBadgeResponse {
  badge: BitBadge | undefined;
}

export interface QueryGetBalanceRequest {
  badgeId: number;
  address: number;
}

export interface QueryGetBalanceResponse {
  balanceInfo: UserBalanceInfo | undefined;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryGetBadgeRequest(): QueryGetBadgeRequest {
  return { id: 0 };
}

export const QueryGetBadgeRequest = {
  encode(message: QueryGetBadgeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetBadgeRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBadgeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBadgeRequest {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: QueryGetBadgeRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBadgeRequest>, I>>(object: I): QueryGetBadgeRequest {
    const message = createBaseQueryGetBadgeRequest();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseQueryGetBadgeResponse(): QueryGetBadgeResponse {
  return { badge: undefined };
}

export const QueryGetBadgeResponse = {
  encode(message: QueryGetBadgeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.badge !== undefined) {
      BitBadge.encode(message.badge, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetBadgeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBadgeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badge = BitBadge.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBadgeResponse {
    return { badge: isSet(object.badge) ? BitBadge.fromJSON(object.badge) : undefined };
  },

  toJSON(message: QueryGetBadgeResponse): unknown {
    const obj: any = {};
    message.badge !== undefined && (obj.badge = message.badge ? BitBadge.toJSON(message.badge) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBadgeResponse>, I>>(object: I): QueryGetBadgeResponse {
    const message = createBaseQueryGetBadgeResponse();
    message.badge = (object.badge !== undefined && object.badge !== null)
      ? BitBadge.fromPartial(object.badge)
      : undefined;
    return message;
  },
};

function createBaseQueryGetBalanceRequest(): QueryGetBalanceRequest {
  return { badgeId: 0, address: 0 };
}

export const QueryGetBalanceRequest = {
  encode(message: QueryGetBalanceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.badgeId !== 0) {
      writer.uint32(8).uint64(message.badgeId);
    }
    if (message.address !== 0) {
      writer.uint32(16).uint64(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetBalanceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBalanceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.address = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBalanceRequest {
    return {
      badgeId: isSet(object.badgeId) ? Number(object.badgeId) : 0,
      address: isSet(object.address) ? Number(object.address) : 0,
    };
  },

  toJSON(message: QueryGetBalanceRequest): unknown {
    const obj: any = {};
    message.badgeId !== undefined && (obj.badgeId = Math.round(message.badgeId));
    message.address !== undefined && (obj.address = Math.round(message.address));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBalanceRequest>, I>>(object: I): QueryGetBalanceRequest {
    const message = createBaseQueryGetBalanceRequest();
    message.badgeId = object.badgeId ?? 0;
    message.address = object.address ?? 0;
    return message;
  },
};

function createBaseQueryGetBalanceResponse(): QueryGetBalanceResponse {
  return { balanceInfo: undefined };
}

export const QueryGetBalanceResponse = {
  encode(message: QueryGetBalanceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.balanceInfo !== undefined) {
      UserBalanceInfo.encode(message.balanceInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetBalanceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBalanceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balanceInfo = UserBalanceInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBalanceResponse {
    return { balanceInfo: isSet(object.balanceInfo) ? UserBalanceInfo.fromJSON(object.balanceInfo) : undefined };
  },

  toJSON(message: QueryGetBalanceResponse): unknown {
    const obj: any = {};
    message.balanceInfo !== undefined
      && (obj.balanceInfo = message.balanceInfo ? UserBalanceInfo.toJSON(message.balanceInfo) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBalanceResponse>, I>>(object: I): QueryGetBalanceResponse {
    const message = createBaseQueryGetBalanceResponse();
    message.balanceInfo = (object.balanceInfo !== undefined && object.balanceInfo !== null)
      ? UserBalanceInfo.fromPartial(object.balanceInfo)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a list of GetBadge items. */
  GetBadge(request: QueryGetBadgeRequest): Promise<QueryGetBadgeResponse>;
  /** Queries a list of GetBalance items. */
  GetBalance(request: QueryGetBalanceRequest): Promise<QueryGetBalanceResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.GetBadge = this.GetBadge.bind(this);
    this.GetBalance = this.GetBalance.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  GetBadge(request: QueryGetBadgeRequest): Promise<QueryGetBadgeResponse> {
    const data = QueryGetBadgeRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetBadge", data);
    return promise.then((data) => QueryGetBadgeResponse.decode(new _m0.Reader(data)));
  }

  GetBalance(request: QueryGetBalanceRequest): Promise<QueryGetBalanceResponse> {
    const data = QueryGetBalanceRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetBalance", data);
    return promise.then((data) => QueryGetBalanceResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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
