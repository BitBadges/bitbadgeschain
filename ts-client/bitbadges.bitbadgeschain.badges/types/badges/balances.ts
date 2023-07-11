/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/**
 * uintRange is a range of IDs from some start to some end (inclusive).
 *
 * uintRanges are one of the core types used in the BitBadgesChain module.
 * They are used for evrything from badge IDs to time ranges to min / max balance amounts.
 */
export interface UintRange {
  start: string;
  end: string;
}

/**
 * Balance represents the balance of a badge for a specific user.
 * The user amounts xAmount of a badge for the badgeID specified for the time ranges specified.
 *
 * Ex: User A owns x10 of badge IDs 1-10 from 1/1/2020 to 1/1/2021.
 *
 * If times or badgeIDs have len > 1, then the user owns all badge IDs specified for all time ranges specified.
 */
export interface Balance {
  amount: string;
  ownedTimes: UintRange[];
  badgeIds: UintRange[];
}

export interface MustOwnBadges {
  collectionId: string;
  amountRange: UintRange | undefined;
  ownedTimes: UintRange[];
  badgeIds: UintRange[];
}

/**
 * InheritedBalances are a powerful feature of the BitBadges module.
 * They allow a colllection to inherit the balances from another collection.
 * Ex: Badges from Collection A inherits the balances from badges from Collection B.
 *
 * The badgeIds specified will inherit the balances from the parent collection and badges specified.
 * If the total number of parent badges == 1, then all the badgeIds will inherit the balance from that parent badge.
 * Otherwise, the total number of parent badges must equal the total number of badgeIds specified.
 * By total number, we mean the sum of the number of badgeIds in each UintRange.
 */
export interface InheritedBalance {
  badgeIds: UintRange[];
  parentCollectionId: string;
  parentBadgeIds: UintRange[];
}

function createBaseUintRange(): UintRange {
  return { start: "", end: "" };
}

export const UintRange = {
  encode(message: UintRange, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== "") {
      writer.uint32(10).string(message.start);
    }
    if (message.end !== "") {
      writer.uint32(18).string(message.end);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UintRange {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUintRange();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = reader.string();
          break;
        case 2:
          message.end = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UintRange {
    return { start: isSet(object.start) ? String(object.start) : "", end: isSet(object.end) ? String(object.end) : "" };
  },

  toJSON(message: UintRange): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = message.start);
    message.end !== undefined && (obj.end = message.end);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UintRange>, I>>(object: I): UintRange {
    const message = createBaseUintRange();
    message.start = object.start ?? "";
    message.end = object.end ?? "";
    return message;
  },
};

function createBaseBalance(): Balance {
  return { amount: "", ownedTimes: [], badgeIds: [] };
}

export const Balance = {
  encode(message: Balance, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.amount !== "") {
      writer.uint32(10).string(message.amount);
    }
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Balance {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBalance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount = reader.string();
          break;
        case 2:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
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

  fromJSON(object: any): Balance {
    return {
      amount: isSet(object.amount) ? String(object.amount) : "",
      ownedTimes: Array.isArray(object?.ownedTimes) ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e)) : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
    };
  },

  toJSON(message: Balance): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = message.amount);
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Balance>, I>>(object: I): Balance {
    const message = createBaseBalance();
    message.amount = object.amount ?? "";
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMustOwnBadges(): MustOwnBadges {
  return { collectionId: "", amountRange: undefined, ownedTimes: [], badgeIds: [] };
}

export const MustOwnBadges = {
  encode(message: MustOwnBadges, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionId !== "") {
      writer.uint32(10).string(message.collectionId);
    }
    if (message.amountRange !== undefined) {
      UintRange.encode(message.amountRange, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MustOwnBadges {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMustOwnBadges();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionId = reader.string();
          break;
        case 2:
          message.amountRange = UintRange.decode(reader, reader.uint32());
          break;
        case 3:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MustOwnBadges {
    return {
      collectionId: isSet(object.collectionId) ? String(object.collectionId) : "",
      amountRange: isSet(object.amountRange) ? UintRange.fromJSON(object.amountRange) : undefined,
      ownedTimes: Array.isArray(object?.ownedTimes) ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e)) : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
    };
  },

  toJSON(message: MustOwnBadges): unknown {
    const obj: any = {};
    message.collectionId !== undefined && (obj.collectionId = message.collectionId);
    message.amountRange !== undefined
      && (obj.amountRange = message.amountRange ? UintRange.toJSON(message.amountRange) : undefined);
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MustOwnBadges>, I>>(object: I): MustOwnBadges {
    const message = createBaseMustOwnBadges();
    message.collectionId = object.collectionId ?? "";
    message.amountRange = (object.amountRange !== undefined && object.amountRange !== null)
      ? UintRange.fromPartial(object.amountRange)
      : undefined;
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseInheritedBalance(): InheritedBalance {
  return { badgeIds: [], parentCollectionId: "", parentBadgeIds: [] };
}

export const InheritedBalance = {
  encode(message: InheritedBalance, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.parentCollectionId !== "") {
      writer.uint32(18).string(message.parentCollectionId);
    }
    for (const v of message.parentBadgeIds) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InheritedBalance {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInheritedBalance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 2:
          message.parentCollectionId = reader.string();
          break;
        case 3:
          message.parentBadgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InheritedBalance {
    return {
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      parentCollectionId: isSet(object.parentCollectionId) ? String(object.parentCollectionId) : "",
      parentBadgeIds: Array.isArray(object?.parentBadgeIds)
        ? object.parentBadgeIds.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: InheritedBalance): unknown {
    const obj: any = {};
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    message.parentCollectionId !== undefined && (obj.parentCollectionId = message.parentCollectionId);
    if (message.parentBadgeIds) {
      obj.parentBadgeIds = message.parentBadgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.parentBadgeIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<InheritedBalance>, I>>(object: I): InheritedBalance {
    const message = createBaseInheritedBalance();
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.parentCollectionId = object.parentCollectionId ?? "";
    message.parentBadgeIds = object.parentBadgeIds?.map((e) => UintRange.fromPartial(e)) || [];
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
