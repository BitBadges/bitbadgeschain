/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bitbadgeschain.badges";

export interface MsgMsgUpdateCollection {
  creator: string;
}

export interface MsgMsgUpdateCollectionResponse {
}

function createBaseMsgMsgUpdateCollection(): MsgMsgUpdateCollection {
  return { creator: "" };
}

export const MsgMsgUpdateCollection = {
  encode(message: MsgMsgUpdateCollection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMsgUpdateCollection {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMsgUpdateCollection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgMsgUpdateCollection {
    return { creator: isSet(object.creator) ? String(object.creator) : "" };
  },

  toJSON(message: MsgMsgUpdateCollection): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMsgUpdateCollection>, I>>(object: I): MsgMsgUpdateCollection {
    const message = createBaseMsgMsgUpdateCollection();
    message.creator = object.creator ?? "";
    return message;
  },
};

function createBaseMsgMsgUpdateCollectionResponse(): MsgMsgUpdateCollectionResponse {
  return {};
}

export const MsgMsgUpdateCollectionResponse = {
  encode(_: MsgMsgUpdateCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgMsgUpdateCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMsgUpdateCollectionResponse();
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

  fromJSON(_: any): MsgMsgUpdateCollectionResponse {
    return {};
  },

  toJSON(_: MsgMsgUpdateCollectionResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgMsgUpdateCollectionResponse>, I>>(_: I): MsgMsgUpdateCollectionResponse {
    const message = createBaseMsgMsgUpdateCollectionResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  MsgUpdateCollection(request: MsgMsgUpdateCollection): Promise<MsgMsgUpdateCollectionResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.MsgUpdateCollection = this.MsgUpdateCollection.bind(this);
  }
  MsgUpdateCollection(request: MsgMsgUpdateCollection): Promise<MsgMsgUpdateCollectionResponse> {
    const data = MsgMsgUpdateCollection.encode(request).finish();
    const promise = this.rpc.request("bitbadgeschain.badges.Msg", "MsgUpdateCollection", data);
    return promise.then((data) => MsgMsgUpdateCollectionResponse.decode(new _m0.Reader(data)));
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
