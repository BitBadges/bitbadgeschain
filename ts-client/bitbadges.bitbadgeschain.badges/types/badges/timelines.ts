/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { InheritedBalance, UintRange } from "./balances";
import { BadgeMetadata, CollectionMetadata, OffChainBalancesMetadata } from "./metadata";
import { CollectionApprovedTransfer } from "./transfers";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

export interface CollectionMetadataTimeline {
  collectionMetadata: CollectionMetadata | undefined;
  timelineTimes: UintRange[];
}

export interface BadgeMetadataTimeline {
  badgeMetadata: BadgeMetadata[];
  timelineTimes: UintRange[];
}

export interface OffChainBalancesMetadataTimeline {
  offChainBalancesMetadata: OffChainBalancesMetadata | undefined;
  timelineTimes: UintRange[];
}

export interface InheritedBalancesTimeline {
  inheritedBalances: InheritedBalance[];
  timelineTimes: UintRange[];
}

export interface CustomDataTimeline {
  customData: string;
  timelineTimes: UintRange[];
}

export interface ManagerTimeline {
  manager: string;
  timelineTimes: UintRange[];
}

export interface CollectionApprovedTransferTimeline {
  approvedTransfers: CollectionApprovedTransfer[];
  timelineTimes: UintRange[];
}

export interface IsArchivedTimeline {
  isArchived: boolean;
  timelineTimes: UintRange[];
}

export interface ContractAddressTimeline {
  contractAddress: string;
  timelineTimes: UintRange[];
}

export interface StandardsTimeline {
  standards: string[];
  timelineTimes: UintRange[];
}

function createBaseCollectionMetadataTimeline(): CollectionMetadataTimeline {
  return { collectionMetadata: undefined, timelineTimes: [] };
}

export const CollectionMetadataTimeline = {
  encode(message: CollectionMetadataTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.collectionMetadata !== undefined) {
      CollectionMetadata.encode(message.collectionMetadata, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionMetadataTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionMetadataTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.collectionMetadata = CollectionMetadata.decode(reader, reader.uint32());
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionMetadataTimeline {
    return {
      collectionMetadata: isSet(object.collectionMetadata)
        ? CollectionMetadata.fromJSON(object.collectionMetadata)
        : undefined,
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionMetadataTimeline): unknown {
    const obj: any = {};
    message.collectionMetadata !== undefined && (obj.collectionMetadata = message.collectionMetadata
      ? CollectionMetadata.toJSON(message.collectionMetadata)
      : undefined);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionMetadataTimeline>, I>>(object: I): CollectionMetadataTimeline {
    const message = createBaseCollectionMetadataTimeline();
    message.collectionMetadata = (object.collectionMetadata !== undefined && object.collectionMetadata !== null)
      ? CollectionMetadata.fromPartial(object.collectionMetadata)
      : undefined;
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseBadgeMetadataTimeline(): BadgeMetadataTimeline {
  return { badgeMetadata: [], timelineTimes: [] };
}

export const BadgeMetadataTimeline = {
  encode(message: BadgeMetadataTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.badgeMetadata) {
      BadgeMetadata.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BadgeMetadataTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBadgeMetadataTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.badgeMetadata.push(BadgeMetadata.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeMetadataTimeline {
    return {
      badgeMetadata: Array.isArray(object?.badgeMetadata)
        ? object.badgeMetadata.map((e: any) => BadgeMetadata.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BadgeMetadataTimeline): unknown {
    const obj: any = {};
    if (message.badgeMetadata) {
      obj.badgeMetadata = message.badgeMetadata.map((e) => e ? BadgeMetadata.toJSON(e) : undefined);
    } else {
      obj.badgeMetadata = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<BadgeMetadataTimeline>, I>>(object: I): BadgeMetadataTimeline {
    const message = createBaseBadgeMetadataTimeline();
    message.badgeMetadata = object.badgeMetadata?.map((e) => BadgeMetadata.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseOffChainBalancesMetadataTimeline(): OffChainBalancesMetadataTimeline {
  return { offChainBalancesMetadata: undefined, timelineTimes: [] };
}

export const OffChainBalancesMetadataTimeline = {
  encode(message: OffChainBalancesMetadataTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.offChainBalancesMetadata !== undefined) {
      OffChainBalancesMetadata.encode(message.offChainBalancesMetadata, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OffChainBalancesMetadataTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOffChainBalancesMetadataTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.offChainBalancesMetadata = OffChainBalancesMetadata.decode(reader, reader.uint32());
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): OffChainBalancesMetadataTimeline {
    return {
      offChainBalancesMetadata: isSet(object.offChainBalancesMetadata)
        ? OffChainBalancesMetadata.fromJSON(object.offChainBalancesMetadata)
        : undefined,
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: OffChainBalancesMetadataTimeline): unknown {
    const obj: any = {};
    message.offChainBalancesMetadata !== undefined && (obj.offChainBalancesMetadata = message.offChainBalancesMetadata
      ? OffChainBalancesMetadata.toJSON(message.offChainBalancesMetadata)
      : undefined);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<OffChainBalancesMetadataTimeline>, I>>(
    object: I,
  ): OffChainBalancesMetadataTimeline {
    const message = createBaseOffChainBalancesMetadataTimeline();
    message.offChainBalancesMetadata =
      (object.offChainBalancesMetadata !== undefined && object.offChainBalancesMetadata !== null)
        ? OffChainBalancesMetadata.fromPartial(object.offChainBalancesMetadata)
        : undefined;
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseInheritedBalancesTimeline(): InheritedBalancesTimeline {
  return { inheritedBalances: [], timelineTimes: [] };
}

export const InheritedBalancesTimeline = {
  encode(message: InheritedBalancesTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.inheritedBalances) {
      InheritedBalance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InheritedBalancesTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInheritedBalancesTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inheritedBalances.push(InheritedBalance.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InheritedBalancesTimeline {
    return {
      inheritedBalances: Array.isArray(object?.inheritedBalances)
        ? object.inheritedBalances.map((e: any) => InheritedBalance.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: InheritedBalancesTimeline): unknown {
    const obj: any = {};
    if (message.inheritedBalances) {
      obj.inheritedBalances = message.inheritedBalances.map((e) => e ? InheritedBalance.toJSON(e) : undefined);
    } else {
      obj.inheritedBalances = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<InheritedBalancesTimeline>, I>>(object: I): InheritedBalancesTimeline {
    const message = createBaseInheritedBalancesTimeline();
    message.inheritedBalances = object.inheritedBalances?.map((e) => InheritedBalance.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCustomDataTimeline(): CustomDataTimeline {
  return { customData: "", timelineTimes: [] };
}

export const CustomDataTimeline = {
  encode(message: CustomDataTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customData !== "") {
      writer.uint32(10).string(message.customData);
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CustomDataTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCustomDataTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customData = reader.string();
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CustomDataTimeline {
    return {
      customData: isSet(object.customData) ? String(object.customData) : "",
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CustomDataTimeline): unknown {
    const obj: any = {};
    message.customData !== undefined && (obj.customData = message.customData);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CustomDataTimeline>, I>>(object: I): CustomDataTimeline {
    const message = createBaseCustomDataTimeline();
    message.customData = object.customData ?? "";
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseManagerTimeline(): ManagerTimeline {
  return { manager: "", timelineTimes: [] };
}

export const ManagerTimeline = {
  encode(message: ManagerTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.manager !== "") {
      writer.uint32(10).string(message.manager);
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ManagerTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseManagerTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.manager = reader.string();
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ManagerTimeline {
    return {
      manager: isSet(object.manager) ? String(object.manager) : "",
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ManagerTimeline): unknown {
    const obj: any = {};
    message.manager !== undefined && (obj.manager = message.manager);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ManagerTimeline>, I>>(object: I): ManagerTimeline {
    const message = createBaseManagerTimeline();
    message.manager = object.manager ?? "";
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCollectionApprovedTransferTimeline(): CollectionApprovedTransferTimeline {
  return { approvedTransfers: [], timelineTimes: [] };
}

export const CollectionApprovedTransferTimeline = {
  encode(message: CollectionApprovedTransferTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.approvedTransfers) {
      CollectionApprovedTransfer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovedTransferTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovedTransferTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvedTransfers.push(CollectionApprovedTransfer.decode(reader, reader.uint32()));
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovedTransferTimeline {
    return {
      approvedTransfers: Array.isArray(object?.approvedTransfers)
        ? object.approvedTransfers.map((e: any) => CollectionApprovedTransfer.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CollectionApprovedTransferTimeline): unknown {
    const obj: any = {};
    if (message.approvedTransfers) {
      obj.approvedTransfers = message.approvedTransfers.map((e) =>
        e ? CollectionApprovedTransfer.toJSON(e) : undefined
      );
    } else {
      obj.approvedTransfers = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransferTimeline>, I>>(
    object: I,
  ): CollectionApprovedTransferTimeline {
    const message = createBaseCollectionApprovedTransferTimeline();
    message.approvedTransfers = object.approvedTransfers?.map((e) => CollectionApprovedTransfer.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseIsArchivedTimeline(): IsArchivedTimeline {
  return { isArchived: false, timelineTimes: [] };
}

export const IsArchivedTimeline = {
  encode(message: IsArchivedTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isArchived === true) {
      writer.uint32(8).bool(message.isArchived);
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsArchivedTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsArchivedTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isArchived = reader.bool();
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsArchivedTimeline {
    return {
      isArchived: isSet(object.isArchived) ? Boolean(object.isArchived) : false,
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: IsArchivedTimeline): unknown {
    const obj: any = {};
    message.isArchived !== undefined && (obj.isArchived = message.isArchived);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsArchivedTimeline>, I>>(object: I): IsArchivedTimeline {
    const message = createBaseIsArchivedTimeline();
    message.isArchived = object.isArchived ?? false;
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseContractAddressTimeline(): ContractAddressTimeline {
  return { contractAddress: "", timelineTimes: [] };
}

export const ContractAddressTimeline = {
  encode(message: ContractAddressTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contractAddress !== "") {
      writer.uint32(10).string(message.contractAddress);
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContractAddressTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContractAddressTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractAddress = reader.string();
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractAddressTimeline {
    return {
      contractAddress: isSet(object.contractAddress) ? String(object.contractAddress) : "",
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ContractAddressTimeline): unknown {
    const obj: any = {};
    message.contractAddress !== undefined && (obj.contractAddress = message.contractAddress);
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ContractAddressTimeline>, I>>(object: I): ContractAddressTimeline {
    const message = createBaseContractAddressTimeline();
    message.contractAddress = object.contractAddress ?? "";
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseStandardsTimeline(): StandardsTimeline {
  return { standards: [], timelineTimes: [] };
}

export const StandardsTimeline = {
  encode(message: StandardsTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.standards) {
      writer.uint32(10).string(v!);
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StandardsTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStandardsTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.standards.push(reader.string());
          break;
        case 2:
          message.timelineTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StandardsTimeline {
    return {
      standards: Array.isArray(object?.standards) ? object.standards.map((e: any) => String(e)) : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: StandardsTimeline): unknown {
    const obj: any = {};
    if (message.standards) {
      obj.standards = message.standards.map((e) => e);
    } else {
      obj.standards = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StandardsTimeline>, I>>(object: I): StandardsTimeline {
    const message = createBaseStandardsTimeline();
    message.standards = object.standards?.map((e) => e) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
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
