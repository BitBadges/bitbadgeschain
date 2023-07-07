/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { UintRange } from "./balances";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/**
 * This defines the metadata for specific badge IDs.
 * This should be interpreted according to the collection standard.
 */
export interface BadgeMetadata {
  uri: string;
  customData: string;
  badgeIds: UintRange[];
}

/**
 * This defines the metadata for the collection.
 * This should be interpreted according to the collection standard.
 */
export interface CollectionMetadata {
  uri: string;
  customData: string;
}

/**
 * This defines the metadata for the off-chain balances (if using this balances type).
 * This should be interpreted according to the collection standard.
 */
export interface OffChainBalancesMetadata {
  uri: string;
  customData: string;
}

function createBaseBadgeMetadata(): BadgeMetadata {
  return { uri: "", customData: "", badgeIds: [] };
}

export const BadgeMetadata = {
  encode(message: BadgeMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.uri !== "") {
      writer.uint32(10).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(18).string(message.customData);
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BadgeMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBadgeMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uri = reader.string();
          break;
        case 2:
          message.customData = reader.string();
          break;
        case 3:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeMetadata {
    return {
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
    };
  },

  toJSON(message: BadgeMetadata): unknown {
    const obj: any = {};
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BadgeMetadata>, I>>(object: I): BadgeMetadata {
    const message = createBaseBadgeMetadata();
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCollectionMetadata(): CollectionMetadata {
  return { uri: "", customData: "" };
}

export const CollectionMetadata = {
  encode(message: CollectionMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.uri !== "") {
      writer.uint32(10).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(18).string(message.customData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uri = reader.string();
          break;
        case 2:
          message.customData = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionMetadata {
    return {
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
    };
  },

  toJSON(message: CollectionMetadata): unknown {
    const obj: any = {};
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionMetadata>, I>>(object: I): CollectionMetadata {
    const message = createBaseCollectionMetadata();
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    return message;
  },
};

function createBaseOffChainBalancesMetadata(): OffChainBalancesMetadata {
  return { uri: "", customData: "" };
}

export const OffChainBalancesMetadata = {
  encode(message: OffChainBalancesMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.uri !== "") {
      writer.uint32(10).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(18).string(message.customData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OffChainBalancesMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOffChainBalancesMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uri = reader.string();
          break;
        case 2:
          message.customData = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): OffChainBalancesMetadata {
    return {
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
    };
  },

  toJSON(message: OffChainBalancesMetadata): unknown {
    const obj: any = {};
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<OffChainBalancesMetadata>, I>>(object: I): OffChainBalancesMetadata {
    const message = createBaseOffChainBalancesMetadata();
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    return message;
  },
};

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
