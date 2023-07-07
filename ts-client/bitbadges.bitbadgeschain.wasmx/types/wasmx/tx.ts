/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bitbadges.bitbadgeschain.wasmx";

/** MsgExecuteContractCompat submits the given message data to a smart contract, compatible with EIP712 */
export interface MsgExecuteContractCompat {
  /** Sender is the that actor that signed the messages */
  sender: string;
  /** Contract is the address of the smart contract */
  contract: string;
  /** Msg json encoded message to be passed to the contract */
  msg: string;
  /** Funds coins that are transferred to the contract on execution */
  funds: string;
}

/** MsgExecuteContractCompatResponse returns execution result data. */
export interface MsgExecuteContractCompatResponse {
  /** Data contains bytes to returned from the contract */
  data: Uint8Array;
}

function createBaseMsgExecuteContractCompat(): MsgExecuteContractCompat {
  return { sender: "", contract: "", msg: "", funds: "" };
}

export const MsgExecuteContractCompat = {
  encode(message: MsgExecuteContractCompat, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sender !== "") {
      writer.uint32(10).string(message.sender);
    }
    if (message.contract !== "") {
      writer.uint32(18).string(message.contract);
    }
    if (message.msg !== "") {
      writer.uint32(26).string(message.msg);
    }
    if (message.funds !== "") {
      writer.uint32(34).string(message.funds);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgExecuteContractCompat {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgExecuteContractCompat();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string();
          break;
        case 2:
          message.contract = reader.string();
          break;
        case 3:
          message.msg = reader.string();
          break;
        case 4:
          message.funds = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgExecuteContractCompat {
    return {
      sender: isSet(object.sender) ? String(object.sender) : "",
      contract: isSet(object.contract) ? String(object.contract) : "",
      msg: isSet(object.msg) ? String(object.msg) : "",
      funds: isSet(object.funds) ? String(object.funds) : "",
    };
  },

  toJSON(message: MsgExecuteContractCompat): unknown {
    const obj: any = {};
    message.sender !== undefined && (obj.sender = message.sender);
    message.contract !== undefined && (obj.contract = message.contract);
    message.msg !== undefined && (obj.msg = message.msg);
    message.funds !== undefined && (obj.funds = message.funds);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgExecuteContractCompat>, I>>(object: I): MsgExecuteContractCompat {
    const message = createBaseMsgExecuteContractCompat();
    message.sender = object.sender ?? "";
    message.contract = object.contract ?? "";
    message.msg = object.msg ?? "";
    message.funds = object.funds ?? "";
    return message;
  },
};

function createBaseMsgExecuteContractCompatResponse(): MsgExecuteContractCompatResponse {
  return { data: new Uint8Array() };
}

export const MsgExecuteContractCompatResponse = {
  encode(message: MsgExecuteContractCompatResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgExecuteContractCompatResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgExecuteContractCompatResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgExecuteContractCompatResponse {
    return { data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array() };
  },

  toJSON(message: MsgExecuteContractCompatResponse): unknown {
    const obj: any = {};
    message.data !== undefined
      && (obj.data = base64FromBytes(message.data !== undefined ? message.data : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgExecuteContractCompatResponse>, I>>(
    object: I,
  ): MsgExecuteContractCompatResponse {
    const message = createBaseMsgExecuteContractCompatResponse();
    message.data = object.data ?? new Uint8Array();
    return message;
  },
};

/** Msg defines the wasmx Msg service. */
export interface Msg {
  ExecuteContractCompat(request: MsgExecuteContractCompat): Promise<MsgExecuteContractCompatResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ExecuteContractCompat = this.ExecuteContractCompat.bind(this);
  }
  ExecuteContractCompat(request: MsgExecuteContractCompat): Promise<MsgExecuteContractCompatResponse> {
    const data = MsgExecuteContractCompat.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.wasmx.Msg", "ExecuteContractCompat", data);
    return promise.then((data) => MsgExecuteContractCompatResponse.decode(new _m0.Reader(data)));
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

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
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
