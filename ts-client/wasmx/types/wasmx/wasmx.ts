/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "wasmx";

export interface Params {
  /** Set the status to active to indicate that contracts can be executed in begin blocker. */
  isExecutionEnabled: boolean;
  /** Maximum aggregate total gas to be used for the contract executions in the BeginBlocker. */
  maxBeginBlockTotalGas: number;
  /** the maximum gas limit each individual contract can consume in the BeginBlocker. */
  maxContractGasLimit: number;
  /** min_gas_price defines the minimum gas price the contracts must pay to be executed in the BeginBlocker. */
  minGasPrice: number;
}

function createBaseParams(): Params {
  return { isExecutionEnabled: false, maxBeginBlockTotalGas: 0, maxContractGasLimit: 0, minGasPrice: 0 };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isExecutionEnabled === true) {
      writer.uint32(8).bool(message.isExecutionEnabled);
    }
    if (message.maxBeginBlockTotalGas !== 0) {
      writer.uint32(16).uint64(message.maxBeginBlockTotalGas);
    }
    if (message.maxContractGasLimit !== 0) {
      writer.uint32(24).uint64(message.maxContractGasLimit);
    }
    if (message.minGasPrice !== 0) {
      writer.uint32(32).uint64(message.minGasPrice);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isExecutionEnabled = reader.bool();
          break;
        case 2:
          message.maxBeginBlockTotalGas = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.maxContractGasLimit = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.minGasPrice = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    return {
      isExecutionEnabled: isSet(object.isExecutionEnabled) ? Boolean(object.isExecutionEnabled) : false,
      maxBeginBlockTotalGas: isSet(object.maxBeginBlockTotalGas) ? Number(object.maxBeginBlockTotalGas) : 0,
      maxContractGasLimit: isSet(object.maxContractGasLimit) ? Number(object.maxContractGasLimit) : 0,
      minGasPrice: isSet(object.minGasPrice) ? Number(object.minGasPrice) : 0,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.isExecutionEnabled !== undefined && (obj.isExecutionEnabled = message.isExecutionEnabled);
    message.maxBeginBlockTotalGas !== undefined
      && (obj.maxBeginBlockTotalGas = Math.round(message.maxBeginBlockTotalGas));
    message.maxContractGasLimit !== undefined && (obj.maxContractGasLimit = Math.round(message.maxContractGasLimit));
    message.minGasPrice !== undefined && (obj.minGasPrice = Math.round(message.minGasPrice));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.isExecutionEnabled = object.isExecutionEnabled ?? false;
    message.maxBeginBlockTotalGas = object.maxBeginBlockTotalGas ?? 0;
    message.maxContractGasLimit = object.maxContractGasLimit ?? 0;
    message.minGasPrice = object.minGasPrice ?? 0;
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
