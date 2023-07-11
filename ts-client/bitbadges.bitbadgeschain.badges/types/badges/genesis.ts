/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { AddressMapping } from "./address_mappings";
import { BadgeCollection } from "./collections";
import { Params } from "./params";
import { ApprovalsTracker, UserBalanceStore } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** GenesisState defines the badges module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  portId: string;
  collections: BadgeCollection[];
  nextCollectionId: string;
  balances: UserBalanceStore[];
  balanceStoreKeys: string[];
  numUsedForMerkleChallenges: string[];
  numUsedForMerkleChallengesStoreKeys: string[];
  addressMappings: AddressMapping[];
  approvalsTrackers: ApprovalsTracker[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  approvalsTrackerStoreKeys: string[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    portId: "",
    collections: [],
    nextCollectionId: "",
    balances: [],
    balanceStoreKeys: [],
    numUsedForMerkleChallenges: [],
    numUsedForMerkleChallengesStoreKeys: [],
    addressMappings: [],
    approvalsTrackers: [],
    approvalsTrackerStoreKeys: [],
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    if (message.portId !== "") {
      writer.uint32(18).string(message.portId);
    }
    for (const v of message.collections) {
      BadgeCollection.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.nextCollectionId !== "") {
      writer.uint32(34).string(message.nextCollectionId);
    }
    for (const v of message.balances) {
      UserBalanceStore.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.balanceStoreKeys) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.numUsedForMerkleChallenges) {
      writer.uint32(58).string(v!);
    }
    for (const v of message.numUsedForMerkleChallengesStoreKeys) {
      writer.uint32(66).string(v!);
    }
    for (const v of message.addressMappings) {
      AddressMapping.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.approvalsTrackers) {
      ApprovalsTracker.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.approvalsTrackerStoreKeys) {
      writer.uint32(90).string(v!);
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
          message.collections.push(BadgeCollection.decode(reader, reader.uint32()));
          break;
        case 4:
          message.nextCollectionId = reader.string();
          break;
        case 5:
          message.balances.push(UserBalanceStore.decode(reader, reader.uint32()));
          break;
        case 6:
          message.balanceStoreKeys.push(reader.string());
          break;
        case 7:
          message.numUsedForMerkleChallenges.push(reader.string());
          break;
        case 8:
          message.numUsedForMerkleChallengesStoreKeys.push(reader.string());
          break;
        case 9:
          message.addressMappings.push(AddressMapping.decode(reader, reader.uint32()));
          break;
        case 10:
          message.approvalsTrackers.push(ApprovalsTracker.decode(reader, reader.uint32()));
          break;
        case 11:
          message.approvalsTrackerStoreKeys.push(reader.string());
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
      collections: Array.isArray(object?.collections)
        ? object.collections.map((e: any) => BadgeCollection.fromJSON(e))
        : [],
      nextCollectionId: isSet(object.nextCollectionId) ? String(object.nextCollectionId) : "",
      balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => UserBalanceStore.fromJSON(e)) : [],
      balanceStoreKeys: Array.isArray(object?.balanceStoreKeys)
        ? object.balanceStoreKeys.map((e: any) => String(e))
        : [],
      numUsedForMerkleChallenges: Array.isArray(object?.numUsedForMerkleChallenges)
        ? object.numUsedForMerkleChallenges.map((e: any) => String(e))
        : [],
      numUsedForMerkleChallengesStoreKeys: Array.isArray(object?.numUsedForMerkleChallengesStoreKeys)
        ? object.numUsedForMerkleChallengesStoreKeys.map((e: any) => String(e))
        : [],
      addressMappings: Array.isArray(object?.addressMappings)
        ? object.addressMappings.map((e: any) => AddressMapping.fromJSON(e))
        : [],
      approvalsTrackers: Array.isArray(object?.approvalsTrackers)
        ? object.approvalsTrackers.map((e: any) => ApprovalsTracker.fromJSON(e))
        : [],
      approvalsTrackerStoreKeys: Array.isArray(object?.approvalsTrackerStoreKeys)
        ? object.approvalsTrackerStoreKeys.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    message.portId !== undefined && (obj.portId = message.portId);
    if (message.collections) {
      obj.collections = message.collections.map((e) => e ? BadgeCollection.toJSON(e) : undefined);
    } else {
      obj.collections = [];
    }
    message.nextCollectionId !== undefined && (obj.nextCollectionId = message.nextCollectionId);
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? UserBalanceStore.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    if (message.balanceStoreKeys) {
      obj.balanceStoreKeys = message.balanceStoreKeys.map((e) => e);
    } else {
      obj.balanceStoreKeys = [];
    }
    if (message.numUsedForMerkleChallenges) {
      obj.numUsedForMerkleChallenges = message.numUsedForMerkleChallenges.map((e) => e);
    } else {
      obj.numUsedForMerkleChallenges = [];
    }
    if (message.numUsedForMerkleChallengesStoreKeys) {
      obj.numUsedForMerkleChallengesStoreKeys = message.numUsedForMerkleChallengesStoreKeys.map((e) => e);
    } else {
      obj.numUsedForMerkleChallengesStoreKeys = [];
    }
    if (message.addressMappings) {
      obj.addressMappings = message.addressMappings.map((e) => e ? AddressMapping.toJSON(e) : undefined);
    } else {
      obj.addressMappings = [];
    }
    if (message.approvalsTrackers) {
      obj.approvalsTrackers = message.approvalsTrackers.map((e) => e ? ApprovalsTracker.toJSON(e) : undefined);
    } else {
      obj.approvalsTrackers = [];
    }
    if (message.approvalsTrackerStoreKeys) {
      obj.approvalsTrackerStoreKeys = message.approvalsTrackerStoreKeys.map((e) => e);
    } else {
      obj.approvalsTrackerStoreKeys = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.portId = object.portId ?? "";
    message.collections = object.collections?.map((e) => BadgeCollection.fromPartial(e)) || [];
    message.nextCollectionId = object.nextCollectionId ?? "";
    message.balances = object.balances?.map((e) => UserBalanceStore.fromPartial(e)) || [];
    message.balanceStoreKeys = object.balanceStoreKeys?.map((e) => e) || [];
    message.numUsedForMerkleChallenges = object.numUsedForMerkleChallenges?.map((e) => e) || [];
    message.numUsedForMerkleChallengesStoreKeys = object.numUsedForMerkleChallengesStoreKeys?.map((e) => e) || [];
    message.addressMappings = object.addressMappings?.map((e) => AddressMapping.fromPartial(e)) || [];
    message.approvalsTrackers = object.approvalsTrackers?.map((e) => ApprovalsTracker.fromPartial(e)) || [];
    message.approvalsTrackerStoreKeys = object.approvalsTrackerStoreKeys?.map((e) => e) || [];
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
