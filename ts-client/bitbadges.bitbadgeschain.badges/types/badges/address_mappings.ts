/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/**
 * An AddressMapping is a permanent list of addresses that are referenced by a mapping ID.
 * The mapping may include only the specified addresses, or it may include all addresses but
 * the specified addresses (depending on if includeAddresses is true or false).
 *
 * AddressMappings are used for things like whitelists, blacklists, approvals, etc.
 */
export interface AddressMapping {
  mappingId: string;
  addresses: string[];
  includeAddresses: boolean;
  uri: string;
  customData: string;
}

function createBaseAddressMapping(): AddressMapping {
  return { mappingId: "", addresses: [], includeAddresses: false, uri: "", customData: "" };
}

export const AddressMapping = {
  encode(message: AddressMapping, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mappingId !== "") {
      writer.uint32(10).string(message.mappingId);
    }
    for (const v of message.addresses) {
      writer.uint32(18).string(v!);
    }
    if (message.includeAddresses === true) {
      writer.uint32(24).bool(message.includeAddresses);
    }
    if (message.uri !== "") {
      writer.uint32(34).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(42).string(message.customData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddressMapping {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddressMapping();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mappingId = reader.string();
          break;
        case 2:
          message.addresses.push(reader.string());
          break;
        case 3:
          message.includeAddresses = reader.bool();
          break;
        case 4:
          message.uri = reader.string();
          break;
        case 5:
          message.customData = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddressMapping {
    return {
      mappingId: isSet(object.mappingId) ? String(object.mappingId) : "",
      addresses: Array.isArray(object?.addresses) ? object.addresses.map((e: any) => String(e)) : [],
      includeAddresses: isSet(object.includeAddresses) ? Boolean(object.includeAddresses) : false,
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
    };
  },

  toJSON(message: AddressMapping): unknown {
    const obj: any = {};
    message.mappingId !== undefined && (obj.mappingId = message.mappingId);
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => e);
    } else {
      obj.addresses = [];
    }
    message.includeAddresses !== undefined && (obj.includeAddresses = message.includeAddresses);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AddressMapping>, I>>(object: I): AddressMapping {
    const message = createBaseAddressMapping();
    message.mappingId = object.mappingId ?? "";
    message.addresses = object.addresses?.map((e) => e) || [];
    message.includeAddresses = object.includeAddresses ?? false;
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
