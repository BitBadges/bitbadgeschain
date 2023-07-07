/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { GenesisState } from "./genesis";
import { Params } from "./wasmx";

export const protobufPackage = "bitbadges.bitbadgeschain.wasmx";

/** QueryWasmxParamsRequest is the request type for the Query/WasmxParams RPC method. */
export interface QueryWasmxParamsRequest {
}

/** QueryWasmxParamsRequest is the response type for the Query/WasmxParams RPC method. */
export interface QueryWasmxParamsResponse {
  params: Params | undefined;
}

/** QueryModuleStateRequest is the request type for the Query/WasmxModuleState RPC method. */
export interface QueryModuleStateRequest {
}

/** QueryModuleStateResponse is the response type for the Query/WasmxModuleState RPC method. */
export interface QueryModuleStateResponse {
  state: GenesisState | undefined;
}

function createBaseQueryWasmxParamsRequest(): QueryWasmxParamsRequest {
  return {};
}

export const QueryWasmxParamsRequest = {
  encode(_: QueryWasmxParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryWasmxParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryWasmxParamsRequest();
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

  fromJSON(_: any): QueryWasmxParamsRequest {
    return {};
  },

  toJSON(_: QueryWasmxParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryWasmxParamsRequest>, I>>(_: I): QueryWasmxParamsRequest {
    const message = createBaseQueryWasmxParamsRequest();
    return message;
  },
};

function createBaseQueryWasmxParamsResponse(): QueryWasmxParamsResponse {
  return { params: undefined };
}

export const QueryWasmxParamsResponse = {
  encode(message: QueryWasmxParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryWasmxParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryWasmxParamsResponse();
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

  fromJSON(object: any): QueryWasmxParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryWasmxParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryWasmxParamsResponse>, I>>(object: I): QueryWasmxParamsResponse {
    const message = createBaseQueryWasmxParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryModuleStateRequest(): QueryModuleStateRequest {
  return {};
}

export const QueryModuleStateRequest = {
  encode(_: QueryModuleStateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryModuleStateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryModuleStateRequest();
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

  fromJSON(_: any): QueryModuleStateRequest {
    return {};
  },

  toJSON(_: QueryModuleStateRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryModuleStateRequest>, I>>(_: I): QueryModuleStateRequest {
    const message = createBaseQueryModuleStateRequest();
    return message;
  },
};

function createBaseQueryModuleStateResponse(): QueryModuleStateResponse {
  return { state: undefined };
}

export const QueryModuleStateResponse = {
  encode(message: QueryModuleStateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.state !== undefined) {
      GenesisState.encode(message.state, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryModuleStateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryModuleStateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.state = GenesisState.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryModuleStateResponse {
    return { state: isSet(object.state) ? GenesisState.fromJSON(object.state) : undefined };
  },

  toJSON(message: QueryModuleStateResponse): unknown {
    const obj: any = {};
    message.state !== undefined && (obj.state = message.state ? GenesisState.toJSON(message.state) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryModuleStateResponse>, I>>(object: I): QueryModuleStateResponse {
    const message = createBaseQueryModuleStateResponse();
    message.state = (object.state !== undefined && object.state !== null)
      ? GenesisState.fromPartial(object.state)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Retrieves wasmx params */
  WasmxParams(request: QueryWasmxParamsRequest): Promise<QueryWasmxParamsResponse>;
  /** Retrieves the entire wasmx module's state */
  WasmxModuleState(request: QueryModuleStateRequest): Promise<QueryModuleStateResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.WasmxParams = this.WasmxParams.bind(this);
    this.WasmxModuleState = this.WasmxModuleState.bind(this);
  }
  WasmxParams(request: QueryWasmxParamsRequest): Promise<QueryWasmxParamsResponse> {
    const data = QueryWasmxParamsRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.wasmx.Query", "WasmxParams", data);
    return promise.then((data) => QueryWasmxParamsResponse.decode(new _m0.Reader(data)));
  }

  WasmxModuleState(request: QueryModuleStateRequest): Promise<QueryModuleStateResponse> {
    const data = QueryModuleStateRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.wasmx.Query", "WasmxModuleState", data);
    return promise.then((data) => QueryModuleStateResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
