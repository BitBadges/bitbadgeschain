/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { AddressMapping } from "./address_mappings";
import { BadgeCollection } from "./collections";
import { Params } from "./params";
import { ApprovalsTracker, UserBalanceStore } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetCollectionRequest {
  collectionId: string;
}

export interface QueryGetCollectionResponse {
  collection: BadgeCollection | undefined;
}

export interface QueryGetBalanceRequest {
  collectionId: string;
  address: string;
}

export interface QueryGetBalanceResponse {
  balance: UserBalanceStore | undefined;
}

export interface QueryGetAddressMappingRequest {
  mappingId: string;
}

export interface QueryGetAddressMappingResponse {
  mapping: AddressMapping | undefined;
}

export interface QueryGetApprovalsTrackerRequest {
  approvalId: string;
  level: string;
  depth: string;
  address: string;
  collectionId: string;
}

export interface QueryGetApprovalsTrackerResponse {
  tracker: ApprovalsTracker | undefined;
}

export interface QueryGetNumUsedForMerkleChallengeRequest {
  challengeId: string;
  level: string;
  leafIndex: string;
  collectionId: string;
}

export interface QueryGetNumUsedForMerkleChallengeResponse {
  numUsed: string;
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

function createBaseQueryGetCollectionRequest(): QueryGetCollectionRequest {
  return { collectionId: "" };
}

export const QueryGetCollectionRequest = {
  encode(message: QueryGetCollectionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetCollectionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCollectionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetCollectionRequest {
    return { collectionId: isSet(object.collectionId) ? String(object.collectionId) : "" };
  },

  toJSON(message: QueryGetCollectionRequest): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetCollectionRequest>, I>>(object: I): QueryGetCollectionRequest {
    const message = createBaseQueryGetCollectionRequest();
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseQueryGetCollectionResponse(): QueryGetCollectionResponse {
  return { collection: undefined };
}

export const QueryGetCollectionResponse = {
  encode(message: QueryGetCollectionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collection !== undefined) {
      BadgeCollection.encode(message.collection, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetCollectionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCollectionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collection = BadgeCollection.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetCollectionResponse {
    return { collection: isSet(object.collection) ? BadgeCollection.fromJSON(object.collection) : undefined };
  },

  toJSON(message: QueryGetCollectionResponse): unknown {
    const obj: any = {};
    message.collection !== undefined
      && (obj.collection = message.collection ? BadgeCollection.toJSON(message.collection) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetCollectionResponse>, I>>(object: I): QueryGetCollectionResponse {
    const message = createBaseQueryGetCollectionResponse();
    message.collection = (object.collection !== undefined && object.collection !== null)
      ? BadgeCollection.fromPartial(object.collection)
      : undefined;
    return message;
  },
};

function createBaseQueryGetBalanceRequest(): QueryGetBalanceRequest {
  return { collectionId: "", address: "" };
}

export const QueryGetBalanceRequest = {
  encode(message: QueryGetBalanceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
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
          message.collectionId = reader.string();
          break;
        case 2:
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
    return {
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      address: isSet(object.address) ? String(object.address) : "",
    };
  },

  toJSON(message: QueryGetBalanceRequest): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBalanceRequest>, I>>(object: I): QueryGetBalanceRequest {
    const message = createBaseQueryGetBalanceRequest();
    message.collectionId = object.collectionId ?? "";
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGetBalanceResponse(): QueryGetBalanceResponse {
  return { balance: undefined };
}

export const QueryGetBalanceResponse = {
  encode(message: QueryGetBalanceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.balance !== undefined) {
      UserBalanceStore.encode(message.balance, writer.uint32(10).fork()).ldelim();
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
          message.balance = UserBalanceStore.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetBalanceResponse {
    return { balance: isSet(object.balance) ? UserBalanceStore.fromJSON(object.balance) : undefined };
  },

  toJSON(message: QueryGetBalanceResponse): unknown {
    const obj: any = {};
    message.balance !== undefined
      && (obj.balance = message.balance ? UserBalanceStore.toJSON(message.balance) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetBalanceResponse>, I>>(object: I): QueryGetBalanceResponse {
    const message = createBaseQueryGetBalanceResponse();
    message.balance = (object.balance !== undefined && object.balance !== null)
      ? UserBalanceStore.fromPartial(object.balance)
      : undefined;
    return message;
  },
};

function createBaseQueryGetAddressMappingRequest(): QueryGetAddressMappingRequest {
  return { mappingId: "" };
}

export const QueryGetAddressMappingRequest = {
  encode(message: QueryGetAddressMappingRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mappingId !== "") {
      writer.uint32(10).string(message.mappingId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetAddressMappingRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetAddressMappingRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mappingId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetAddressMappingRequest {
    return { mappingId: isSet(object.mappingId) ? String(object.mappingId) : "" };
  },

  toJSON(message: QueryGetAddressMappingRequest): unknown {
    const obj: any = {};
    message.mappingId !== undefined && (obj.mappingId = message.mappingId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetAddressMappingRequest>, I>>(
    object: I,
  ): QueryGetAddressMappingRequest {
    const message = createBaseQueryGetAddressMappingRequest();
    message.mappingId = object.mappingId ?? "";
    return message;
  },
};

function createBaseQueryGetAddressMappingResponse(): QueryGetAddressMappingResponse {
  return { mapping: undefined };
}

export const QueryGetAddressMappingResponse = {
  encode(message: QueryGetAddressMappingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mapping !== undefined) {
      AddressMapping.encode(message.mapping, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetAddressMappingResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetAddressMappingResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mapping = AddressMapping.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetAddressMappingResponse {
    return { mapping: isSet(object.mapping) ? AddressMapping.fromJSON(object.mapping) : undefined };
  },

  toJSON(message: QueryGetAddressMappingResponse): unknown {
    const obj: any = {};
    message.mapping !== undefined
      && (obj.mapping = message.mapping ? AddressMapping.toJSON(message.mapping) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetAddressMappingResponse>, I>>(
    object: I,
  ): QueryGetAddressMappingResponse {
    const message = createBaseQueryGetAddressMappingResponse();
    message.mapping = (object.mapping !== undefined && object.mapping !== null)
      ? AddressMapping.fromPartial(object.mapping)
      : undefined;
    return message;
  },
};

function createBaseQueryGetApprovalsTrackerRequest(): QueryGetApprovalsTrackerRequest {
  return { approvalId: "", level: "", depth: "", address: "", collectionId: "" };
}

export const QueryGetApprovalsTrackerRequest = {
  encode(message: QueryGetApprovalsTrackerRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(10).string(message.approvalId);
    }
    if (message.level !== "") {
      writer.uint32(18).string(message.level);
    }
    if (message.depth !== "") {
      writer.uint32(26).string(message.depth);
    }
    if (message.address !== "") {
      writer.uint32(34).string(message.address);
    }
    if (message.collectionId !== "") {
      writer.uint32(42).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetApprovalsTrackerRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetApprovalsTrackerRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvalId = reader.string();
          break;
        case 2:
          message.level = reader.string();
          break;
        case 3:
          message.depth = reader.string();
          break;
        case 4:
          message.address = reader.string();
          break;
        case 5:
          message.collectionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetApprovalsTrackerRequest {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      level: isSet(object.level) ? String(object.level) : "",
      depth: isSet(object.depth) ? String(object.depth) : "",
      address: isSet(object.address) ? String(object.address) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
    };
  },

  toJSON(message: QueryGetApprovalsTrackerRequest): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.level !== undefined && (obj.level = message.level);
    message.depth !== undefined && (obj.depth = message.depth);
    message.address !== undefined && (obj.address = message.address);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetApprovalsTrackerRequest>, I>>(
    object: I,
  ): QueryGetApprovalsTrackerRequest {
    const message = createBaseQueryGetApprovalsTrackerRequest();
    message.approvalId = object.approvalId ?? "";
    message.level = object.level ?? "";
    message.depth = object.depth ?? "";
    message.address = object.address ?? "";
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseQueryGetApprovalsTrackerResponse(): QueryGetApprovalsTrackerResponse {
  return { tracker: undefined };
}

export const QueryGetApprovalsTrackerResponse = {
  encode(message: QueryGetApprovalsTrackerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tracker !== undefined) {
      ApprovalsTracker.encode(message.tracker, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetApprovalsTrackerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetApprovalsTrackerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tracker = ApprovalsTracker.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetApprovalsTrackerResponse {
    return { tracker: isSet(object.tracker) ? ApprovalsTracker.fromJSON(object.tracker) : undefined };
  },

  toJSON(message: QueryGetApprovalsTrackerResponse): unknown {
    const obj: any = {};
    message.tracker !== undefined
      && (obj.tracker = message.tracker ? ApprovalsTracker.toJSON(message.tracker) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetApprovalsTrackerResponse>, I>>(
    object: I,
  ): QueryGetApprovalsTrackerResponse {
    const message = createBaseQueryGetApprovalsTrackerResponse();
    message.tracker = (object.tracker !== undefined && object.tracker !== null)
      ? ApprovalsTracker.fromPartial(object.tracker)
      : undefined;
    return message;
  },
};

function createBaseQueryGetNumUsedForMerkleChallengeRequest(): QueryGetNumUsedForMerkleChallengeRequest {
  return { challengeId: "", level: "", leafIndex: "", collectionId: "" };
}

export const QueryGetNumUsedForMerkleChallengeRequest = {
  encode(message: QueryGetNumUsedForMerkleChallengeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.challengeId !== "") {
      writer.uint32(10).string(message.challengeId);
    }
    if (message.level !== "") {
      writer.uint32(18).string(message.level);
    }
    if (message.leafIndex !== "") {
      writer.uint32(26).string(message.leafIndex);
    }
    if (message.collectionId !== "") {
      writer.uint32(34).string(message.collectionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetNumUsedForMerkleChallengeRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetNumUsedForMerkleChallengeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.challengeId = reader.string();
          break;
        case 2:
          message.level = reader.string();
          break;
        case 3:
          message.leafIndex = reader.string();
          break;
        case 4:
          message.collectionId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetNumUsedForMerkleChallengeRequest {
    return {
      challengeId: isSet(object.challengeId) ? String(object.challengeId) : "",
      level: isSet(object.level) ? String(object.level) : "",
      leafIndex: isSet(object.leafIndex) ? String(object.leafIndex) : "",
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
    };
  },

  toJSON(message: QueryGetNumUsedForMerkleChallengeRequest): unknown {
    const obj: any = {};
    message.challengeId !== undefined && (obj.challengeId = message.challengeId);
    message.level !== undefined && (obj.level = message.level);
    message.leafIndex !== undefined && (obj.leafIndex = message.leafIndex);
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetNumUsedForMerkleChallengeRequest>, I>>(
    object: I,
  ): QueryGetNumUsedForMerkleChallengeRequest {
    const message = createBaseQueryGetNumUsedForMerkleChallengeRequest();
    message.challengeId = object.challengeId ?? "";
    message.level = object.level ?? "";
    message.leafIndex = object.leafIndex ?? "";
    message.collectionId = object.collectionId ?? "";
    return message;
  },
};

function createBaseQueryGetNumUsedForMerkleChallengeResponse(): QueryGetNumUsedForMerkleChallengeResponse {
  return { numUsed: "" };
}

export const QueryGetNumUsedForMerkleChallengeResponse = {
  encode(message: QueryGetNumUsedForMerkleChallengeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.numUsed !== "") {
      writer.uint32(10).string(message.numUsed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetNumUsedForMerkleChallengeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetNumUsedForMerkleChallengeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.numUsed = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetNumUsedForMerkleChallengeResponse {
    return { numUsed: isSet(object.numUsed) ? String(object.numUsed) : "" };
  },

  toJSON(message: QueryGetNumUsedForMerkleChallengeResponse): unknown {
    const obj: any = {};
    message.numUsed !== undefined && (obj.numUsed = message.numUsed);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetNumUsedForMerkleChallengeResponse>, I>>(
    object: I,
  ): QueryGetNumUsedForMerkleChallengeResponse {
    const message = createBaseQueryGetNumUsedForMerkleChallengeResponse();
    message.numUsed = object.numUsed ?? "";
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a badge collection by ID. */
  GetCollection(request: QueryGetCollectionRequest): Promise<QueryGetCollectionResponse>;
  GetAddressMapping(request: QueryGetAddressMappingRequest): Promise<QueryGetAddressMappingResponse>;
  GetApprovalsTracker(request: QueryGetApprovalsTrackerRequest): Promise<QueryGetApprovalsTrackerResponse>;
  GetNumUsedForMerkleChallenge(request: QueryGetNumUsedForMerkleChallengeRequest): Promise<QueryGetNumUsedForMerkleChallengeResponse>;
  /** Queries an addresses balance for a badge collection, specified by its ID. */
  GetBalance(request: QueryGetBalanceRequest): Promise<QueryGetBalanceResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.GetCollection = this.GetCollection.bind(this);
    this.GetAddressMapping = this.GetAddressMapping.bind(this);
    this.GetApprovalsTracker = this.GetApprovalsTracker.bind(this);
    this.GetNumUsedForMerkleChallenge = this.GetNumUsedForMerkleChallenge.bind(this);
    this.GetBalance = this.GetBalance.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  GetCollection(request: QueryGetCollectionRequest): Promise<QueryGetCollectionResponse> {
    const data = QueryGetCollectionRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetCollection", data);
    return promise.then((data) => QueryGetCollectionResponse.decode(new _m0.Reader(data)));
  }

  GetAddressMapping(request: QueryGetAddressMappingRequest): Promise<QueryGetAddressMappingResponse> {
    const data = QueryGetAddressMappingRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetAddressMapping", data);
    return promise.then((data) => QueryGetAddressMappingResponse.decode(new _m0.Reader(data)));
  }

  GetApprovalsTracker(request: QueryGetApprovalsTrackerRequest): Promise<QueryGetApprovalsTrackerResponse> {
    const data = QueryGetApprovalsTrackerRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetApprovalsTracker", data);
    return promise.then((data) => QueryGetApprovalsTrackerResponse.decode(new _m0.Reader(data)));
  }

  GetNumUsedForMerkleChallenge(request: QueryGetNumUsedForMerkleChallengeRequest): Promise<QueryGetNumUsedForMerkleChallengeResponse> {
    const data = QueryGetNumUsedForMerkleChallengeRequest.encode(request).finish();
    const promise = this.rpc.request("bitbadges.bitbadgeschain.badges.Query", "GetNumUsedForMerkleChallenge", data);
    return promise.then((data) => QueryGetNumUsedForMerkleChallengeResponse.decode(new _m0.Reader(data)));
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
