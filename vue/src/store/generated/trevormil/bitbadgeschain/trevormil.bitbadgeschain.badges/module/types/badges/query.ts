/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { Params } from "../badges/params";
import { BitBadge } from "../badges/badges";
import { BadgeBalanceInfo } from "../badges/balances";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}

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
  subbadgeId: number;
  address: string;
}

export interface QueryGetBalanceResponse {
  balanceInfo: BadgeBalanceInfo | undefined;
  message: string;
}

const baseQueryParamsRequest: object = {};

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
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
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = { ...baseQueryParamsRequest } as QueryParamsRequest;
    return message;
  },
};

const baseQueryParamsResponse: object = {};

export const QueryParamsResponse = {
  encode(
    message: QueryParamsResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
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
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = { ...baseQueryParamsResponse } as QueryParamsResponse;
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    return message;
  },
};

const baseQueryGetBadgeRequest: object = { id: 0 };

export const QueryGetBadgeRequest = {
  encode(
    message: QueryGetBadgeRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetBadgeRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetBadgeRequest } as QueryGetBadgeRequest;
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
    const message = { ...baseQueryGetBadgeRequest } as QueryGetBadgeRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    return message;
  },

  toJSON(message: QueryGetBadgeRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<QueryGetBadgeRequest>): QueryGetBadgeRequest {
    const message = { ...baseQueryGetBadgeRequest } as QueryGetBadgeRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    return message;
  },
};

const baseQueryGetBadgeResponse: object = {};

export const QueryGetBadgeResponse = {
  encode(
    message: QueryGetBadgeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.badge !== undefined) {
      BitBadge.encode(message.badge, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetBadgeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetBadgeResponse } as QueryGetBadgeResponse;
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
    const message = { ...baseQueryGetBadgeResponse } as QueryGetBadgeResponse;
    if (object.badge !== undefined && object.badge !== null) {
      message.badge = BitBadge.fromJSON(object.badge);
    } else {
      message.badge = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetBadgeResponse): unknown {
    const obj: any = {};
    message.badge !== undefined &&
      (obj.badge = message.badge ? BitBadge.toJSON(message.badge) : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetBadgeResponse>
  ): QueryGetBadgeResponse {
    const message = { ...baseQueryGetBadgeResponse } as QueryGetBadgeResponse;
    if (object.badge !== undefined && object.badge !== null) {
      message.badge = BitBadge.fromPartial(object.badge);
    } else {
      message.badge = undefined;
    }
    return message;
  },
};

const baseQueryGetBalanceRequest: object = {
  badgeId: 0,
  subbadgeId: 0,
  address: "",
};

export const QueryGetBalanceRequest = {
  encode(
    message: QueryGetBalanceRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.badgeId !== 0) {
      writer.uint32(8).uint64(message.badgeId);
    }
    if (message.subbadgeId !== 0) {
      writer.uint32(16).uint64(message.subbadgeId);
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetBalanceRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseQueryGetBalanceRequest } as QueryGetBalanceRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.subbadgeId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBalanceRequest {
    const message = { ...baseQueryGetBalanceRequest } as QueryGetBalanceRequest;
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = Number(object.badgeId);
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = Number(object.subbadgeId);
    } else {
      message.subbadgeId = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = String(object.address);
    } else {
      message.address = "";
    }
    return message;
  },

  toJSON(message: QueryGetBalanceRequest): unknown {
    const obj: any = {};
    message.badgeId !== undefined && (obj.badgeId = message.badgeId);
    message.subbadgeId !== undefined && (obj.subbadgeId = message.subbadgeId);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetBalanceRequest>
  ): QueryGetBalanceRequest {
    const message = { ...baseQueryGetBalanceRequest } as QueryGetBalanceRequest;
    if (object.badgeId !== undefined && object.badgeId !== null) {
      message.badgeId = object.badgeId;
    } else {
      message.badgeId = 0;
    }
    if (object.subbadgeId !== undefined && object.subbadgeId !== null) {
      message.subbadgeId = object.subbadgeId;
    } else {
      message.subbadgeId = 0;
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    } else {
      message.address = "";
    }
    return message;
  },
};

const baseQueryGetBalanceResponse: object = { message: "" };

export const QueryGetBalanceResponse = {
  encode(
    message: QueryGetBalanceResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.balanceInfo !== undefined) {
      BadgeBalanceInfo.encode(
        message.balanceInfo,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): QueryGetBalanceResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetBalanceResponse,
    } as QueryGetBalanceResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balanceInfo = BadgeBalanceInfo.decode(
            reader,
            reader.uint32()
          );
          break;
        case 2:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBalanceResponse {
    const message = {
      ...baseQueryGetBalanceResponse,
    } as QueryGetBalanceResponse;
    if (object.balanceInfo !== undefined && object.balanceInfo !== null) {
      message.balanceInfo = BadgeBalanceInfo.fromJSON(object.balanceInfo);
    } else {
      message.balanceInfo = undefined;
    }
    if (object.message !== undefined && object.message !== null) {
      message.message = String(object.message);
    } else {
      message.message = "";
    }
    return message;
  },

  toJSON(message: QueryGetBalanceResponse): unknown {
    const obj: any = {};
    message.balanceInfo !== undefined &&
      (obj.balanceInfo = message.balanceInfo
        ? BadgeBalanceInfo.toJSON(message.balanceInfo)
        : undefined);
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetBalanceResponse>
  ): QueryGetBalanceResponse {
    const message = {
      ...baseQueryGetBalanceResponse,
    } as QueryGetBalanceResponse;
    if (object.balanceInfo !== undefined && object.balanceInfo !== null) {
      message.balanceInfo = BadgeBalanceInfo.fromPartial(object.balanceInfo);
    } else {
      message.balanceInfo = undefined;
    }
    if (object.message !== undefined && object.message !== null) {
      message.message = object.message;
    } else {
      message.message = "";
    }
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
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Query",
      "Params",
      data
    );
    return promise.then((data) => QueryParamsResponse.decode(new Reader(data)));
  }

  GetBadge(request: QueryGetBadgeRequest): Promise<QueryGetBadgeResponse> {
    const data = QueryGetBadgeRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Query",
      "GetBadge",
      data
    );
    return promise.then((data) =>
      QueryGetBadgeResponse.decode(new Reader(data))
    );
  }

  GetBalance(
    request: QueryGetBalanceRequest
  ): Promise<QueryGetBalanceResponse> {
    const data = QueryGetBalanceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "trevormil.bitbadgeschain.badges.Query",
      "GetBalance",
      data
    );
    return promise.then((data) =>
      QueryGetBalanceResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

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
