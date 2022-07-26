/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { Params } from "../badges/params";
import { BitBadge } from "../badges/badges";

export const protobufPackage = "trevormil.bitbadgeschain.badges";

/** GenesisState defines the badges module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  port_id: string;
  /** badge id => BitBadge object */
  badges: { [key: string]: BitBadge };
  /** address => AddressOwnershipInfo object */
  addresses: { [key: string]: AddressOwnershipInfo };
}

export interface GenesisState_BadgesEntry {
  key: string;
  value: BitBadge | undefined;
}

export interface GenesisState_AddressesEntry {
  key: string;
  value: AddressOwnershipInfo | undefined;
}

/** AddressOwnershipInfo stores all badge data relating to a specific address */
export interface AddressOwnershipInfo {
  /** badges is a map of badges owned by the owner: indexed by (badge id + subasset_id) */
  badges_owned: { [key: string]: BadgeOwnershipInfo };
}

export interface AddressOwnershipInfo_BadgesOwnedEntry {
  key: string;
  value: BadgeOwnershipInfo | undefined;
}

export interface BadgeOwnershipInfo {
  balance: number;
  pending_nonce: number;
  /** indexed by (from's pending_nonce - to's pending_nonce) */
  pending: { [key: string]: PendingTransfers };
  /** address => amount */
  approvals: { [key: string]: number };
}

export interface BadgeOwnershipInfo_PendingEntry {
  key: string;
  value: PendingTransfers | undefined;
}

export interface BadgeOwnershipInfo_ApprovalsEntry {
  key: string;
  value: number;
}

/** Pending transfers will not be saved after accept / reject */
export interface PendingTransfers {
  amount: number;
  /** vs. receive request */
  send_request: boolean;
  to: string;
  from: string;
  memo: string;
}

const baseGenesisState: object = { port_id: "" };

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    if (message.port_id !== "") {
      writer.uint32(18).string(message.port_id);
    }
    Object.entries(message.badges).forEach(([key, value]) => {
      GenesisState_BadgesEntry.encode(
        { key: key as any, value },
        writer.uint32(26).fork()
      ).ldelim();
    });
    Object.entries(message.addresses).forEach(([key, value]) => {
      GenesisState_AddressesEntry.encode(
        { key: key as any, value },
        writer.uint32(34).fork()
      ).ldelim();
    });
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = {};
    message.addresses = {};
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.port_id = reader.string();
          break;
        case 3:
          const entry3 = GenesisState_BadgesEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry3.value !== undefined) {
            message.badges[entry3.key] = entry3.value;
          }
          break;
        case 4:
          const entry4 = GenesisState_AddressesEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry4.value !== undefined) {
            message.addresses[entry4.key] = entry4.value;
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = {};
    message.addresses = {};
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    if (object.port_id !== undefined && object.port_id !== null) {
      message.port_id = String(object.port_id);
    } else {
      message.port_id = "";
    }
    if (object.badges !== undefined && object.badges !== null) {
      Object.entries(object.badges).forEach(([key, value]) => {
        message.badges[key] = BitBadge.fromJSON(value);
      });
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      Object.entries(object.addresses).forEach(([key, value]) => {
        message.addresses[key] = AddressOwnershipInfo.fromJSON(value);
      });
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    message.port_id !== undefined && (obj.port_id = message.port_id);
    obj.badges = {};
    if (message.badges) {
      Object.entries(message.badges).forEach(([k, v]) => {
        obj.badges[k] = BitBadge.toJSON(v);
      });
    }
    obj.addresses = {};
    if (message.addresses) {
      Object.entries(message.addresses).forEach(([k, v]) => {
        obj.addresses[k] = AddressOwnershipInfo.toJSON(v);
      });
    }
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.badges = {};
    message.addresses = {};
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    if (object.port_id !== undefined && object.port_id !== null) {
      message.port_id = object.port_id;
    } else {
      message.port_id = "";
    }
    if (object.badges !== undefined && object.badges !== null) {
      Object.entries(object.badges).forEach(([key, value]) => {
        if (value !== undefined) {
          message.badges[key] = BitBadge.fromPartial(value);
        }
      });
    }
    if (object.addresses !== undefined && object.addresses !== null) {
      Object.entries(object.addresses).forEach(([key, value]) => {
        if (value !== undefined) {
          message.addresses[key] = AddressOwnershipInfo.fromPartial(value);
        }
      });
    }
    return message;
  },
};

const baseGenesisState_BadgesEntry: object = { key: "" };

export const GenesisState_BadgesEntry = {
  encode(
    message: GenesisState_BadgesEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      BitBadge.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): GenesisState_BadgesEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseGenesisState_BadgesEntry,
    } as GenesisState_BadgesEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = BitBadge.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState_BadgesEntry {
    const message = {
      ...baseGenesisState_BadgesEntry,
    } as GenesisState_BadgesEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = String(object.key);
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = BitBadge.fromJSON(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },

  toJSON(message: GenesisState_BadgesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value ? BitBadge.toJSON(message.value) : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<GenesisState_BadgesEntry>
  ): GenesisState_BadgesEntry {
    const message = {
      ...baseGenesisState_BadgesEntry,
    } as GenesisState_BadgesEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = BitBadge.fromPartial(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },
};

const baseGenesisState_AddressesEntry: object = { key: "" };

export const GenesisState_AddressesEntry = {
  encode(
    message: GenesisState_AddressesEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      AddressOwnershipInfo.encode(
        message.value,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): GenesisState_AddressesEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseGenesisState_AddressesEntry,
    } as GenesisState_AddressesEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = AddressOwnershipInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState_AddressesEntry {
    const message = {
      ...baseGenesisState_AddressesEntry,
    } as GenesisState_AddressesEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = String(object.key);
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = AddressOwnershipInfo.fromJSON(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },

  toJSON(message: GenesisState_AddressesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value
        ? AddressOwnershipInfo.toJSON(message.value)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<GenesisState_AddressesEntry>
  ): GenesisState_AddressesEntry {
    const message = {
      ...baseGenesisState_AddressesEntry,
    } as GenesisState_AddressesEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = AddressOwnershipInfo.fromPartial(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },
};

const baseAddressOwnershipInfo: object = {};

export const AddressOwnershipInfo = {
  encode(
    message: AddressOwnershipInfo,
    writer: Writer = Writer.create()
  ): Writer {
    Object.entries(message.badges_owned).forEach(([key, value]) => {
      AddressOwnershipInfo_BadgesOwnedEntry.encode(
        { key: key as any, value },
        writer.uint32(10).fork()
      ).ldelim();
    });
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AddressOwnershipInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddressOwnershipInfo } as AddressOwnershipInfo;
    message.badges_owned = {};
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          const entry1 = AddressOwnershipInfo_BadgesOwnedEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry1.value !== undefined) {
            message.badges_owned[entry1.key] = entry1.value;
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddressOwnershipInfo {
    const message = { ...baseAddressOwnershipInfo } as AddressOwnershipInfo;
    message.badges_owned = {};
    if (object.badges_owned !== undefined && object.badges_owned !== null) {
      Object.entries(object.badges_owned).forEach(([key, value]) => {
        message.badges_owned[key] = BadgeOwnershipInfo.fromJSON(value);
      });
    }
    return message;
  },

  toJSON(message: AddressOwnershipInfo): unknown {
    const obj: any = {};
    obj.badges_owned = {};
    if (message.badges_owned) {
      Object.entries(message.badges_owned).forEach(([k, v]) => {
        obj.badges_owned[k] = BadgeOwnershipInfo.toJSON(v);
      });
    }
    return obj;
  },

  fromPartial(object: DeepPartial<AddressOwnershipInfo>): AddressOwnershipInfo {
    const message = { ...baseAddressOwnershipInfo } as AddressOwnershipInfo;
    message.badges_owned = {};
    if (object.badges_owned !== undefined && object.badges_owned !== null) {
      Object.entries(object.badges_owned).forEach(([key, value]) => {
        if (value !== undefined) {
          message.badges_owned[key] = BadgeOwnershipInfo.fromPartial(value);
        }
      });
    }
    return message;
  },
};

const baseAddressOwnershipInfo_BadgesOwnedEntry: object = { key: "" };

export const AddressOwnershipInfo_BadgesOwnedEntry = {
  encode(
    message: AddressOwnershipInfo_BadgesOwnedEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      BadgeOwnershipInfo.encode(
        message.value,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): AddressOwnershipInfo_BadgesOwnedEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseAddressOwnershipInfo_BadgesOwnedEntry,
    } as AddressOwnershipInfo_BadgesOwnedEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = BadgeOwnershipInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddressOwnershipInfo_BadgesOwnedEntry {
    const message = {
      ...baseAddressOwnershipInfo_BadgesOwnedEntry,
    } as AddressOwnershipInfo_BadgesOwnedEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = String(object.key);
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = BadgeOwnershipInfo.fromJSON(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },

  toJSON(message: AddressOwnershipInfo_BadgesOwnedEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value
        ? BadgeOwnershipInfo.toJSON(message.value)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<AddressOwnershipInfo_BadgesOwnedEntry>
  ): AddressOwnershipInfo_BadgesOwnedEntry {
    const message = {
      ...baseAddressOwnershipInfo_BadgesOwnedEntry,
    } as AddressOwnershipInfo_BadgesOwnedEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = BadgeOwnershipInfo.fromPartial(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },
};

const baseBadgeOwnershipInfo: object = { balance: 0, pending_nonce: 0 };

export const BadgeOwnershipInfo = {
  encode(
    message: BadgeOwnershipInfo,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.balance !== 0) {
      writer.uint32(8).uint64(message.balance);
    }
    if (message.pending_nonce !== 0) {
      writer.uint32(16).uint64(message.pending_nonce);
    }
    Object.entries(message.pending).forEach(([key, value]) => {
      BadgeOwnershipInfo_PendingEntry.encode(
        { key: key as any, value },
        writer.uint32(26).fork()
      ).ldelim();
    });
    Object.entries(message.approvals).forEach(([key, value]) => {
      BadgeOwnershipInfo_ApprovalsEntry.encode(
        { key: key as any, value },
        writer.uint32(34).fork()
      ).ldelim();
    });
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): BadgeOwnershipInfo {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseBadgeOwnershipInfo } as BadgeOwnershipInfo;
    message.pending = {};
    message.approvals = {};
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balance = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.pending_nonce = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          const entry3 = BadgeOwnershipInfo_PendingEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry3.value !== undefined) {
            message.pending[entry3.key] = entry3.value;
          }
          break;
        case 4:
          const entry4 = BadgeOwnershipInfo_ApprovalsEntry.decode(
            reader,
            reader.uint32()
          );
          if (entry4.value !== undefined) {
            message.approvals[entry4.key] = entry4.value;
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeOwnershipInfo {
    const message = { ...baseBadgeOwnershipInfo } as BadgeOwnershipInfo;
    message.pending = {};
    message.approvals = {};
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = Number(object.balance);
    } else {
      message.balance = 0;
    }
    if (object.pending_nonce !== undefined && object.pending_nonce !== null) {
      message.pending_nonce = Number(object.pending_nonce);
    } else {
      message.pending_nonce = 0;
    }
    if (object.pending !== undefined && object.pending !== null) {
      Object.entries(object.pending).forEach(([key, value]) => {
        message.pending[key] = PendingTransfers.fromJSON(value);
      });
    }
    if (object.approvals !== undefined && object.approvals !== null) {
      Object.entries(object.approvals).forEach(([key, value]) => {
        message.approvals[key] = Number(value);
      });
    }
    return message;
  },

  toJSON(message: BadgeOwnershipInfo): unknown {
    const obj: any = {};
    message.balance !== undefined && (obj.balance = message.balance);
    message.pending_nonce !== undefined &&
      (obj.pending_nonce = message.pending_nonce);
    obj.pending = {};
    if (message.pending) {
      Object.entries(message.pending).forEach(([k, v]) => {
        obj.pending[k] = PendingTransfers.toJSON(v);
      });
    }
    obj.approvals = {};
    if (message.approvals) {
      Object.entries(message.approvals).forEach(([k, v]) => {
        obj.approvals[k] = v;
      });
    }
    return obj;
  },

  fromPartial(object: DeepPartial<BadgeOwnershipInfo>): BadgeOwnershipInfo {
    const message = { ...baseBadgeOwnershipInfo } as BadgeOwnershipInfo;
    message.pending = {};
    message.approvals = {};
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    } else {
      message.balance = 0;
    }
    if (object.pending_nonce !== undefined && object.pending_nonce !== null) {
      message.pending_nonce = object.pending_nonce;
    } else {
      message.pending_nonce = 0;
    }
    if (object.pending !== undefined && object.pending !== null) {
      Object.entries(object.pending).forEach(([key, value]) => {
        if (value !== undefined) {
          message.pending[key] = PendingTransfers.fromPartial(value);
        }
      });
    }
    if (object.approvals !== undefined && object.approvals !== null) {
      Object.entries(object.approvals).forEach(([key, value]) => {
        if (value !== undefined) {
          message.approvals[key] = Number(value);
        }
      });
    }
    return message;
  },
};

const baseBadgeOwnershipInfo_PendingEntry: object = { key: "" };

export const BadgeOwnershipInfo_PendingEntry = {
  encode(
    message: BadgeOwnershipInfo_PendingEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      PendingTransfers.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): BadgeOwnershipInfo_PendingEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBadgeOwnershipInfo_PendingEntry,
    } as BadgeOwnershipInfo_PendingEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = PendingTransfers.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeOwnershipInfo_PendingEntry {
    const message = {
      ...baseBadgeOwnershipInfo_PendingEntry,
    } as BadgeOwnershipInfo_PendingEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = String(object.key);
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = PendingTransfers.fromJSON(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },

  toJSON(message: BadgeOwnershipInfo_PendingEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value
        ? PendingTransfers.toJSON(message.value)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<BadgeOwnershipInfo_PendingEntry>
  ): BadgeOwnershipInfo_PendingEntry {
    const message = {
      ...baseBadgeOwnershipInfo_PendingEntry,
    } as BadgeOwnershipInfo_PendingEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = PendingTransfers.fromPartial(object.value);
    } else {
      message.value = undefined;
    }
    return message;
  },
};

const baseBadgeOwnershipInfo_ApprovalsEntry: object = { key: "", value: 0 };

export const BadgeOwnershipInfo_ApprovalsEntry = {
  encode(
    message: BadgeOwnershipInfo_ApprovalsEntry,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== 0) {
      writer.uint32(16).uint64(message.value);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): BadgeOwnershipInfo_ApprovalsEntry {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseBadgeOwnershipInfo_ApprovalsEntry,
    } as BadgeOwnershipInfo_ApprovalsEntry;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BadgeOwnershipInfo_ApprovalsEntry {
    const message = {
      ...baseBadgeOwnershipInfo_ApprovalsEntry,
    } as BadgeOwnershipInfo_ApprovalsEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = String(object.key);
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = Number(object.value);
    } else {
      message.value = 0;
    }
    return message;
  },

  toJSON(message: BadgeOwnershipInfo_ApprovalsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial(
    object: DeepPartial<BadgeOwnershipInfo_ApprovalsEntry>
  ): BadgeOwnershipInfo_ApprovalsEntry {
    const message = {
      ...baseBadgeOwnershipInfo_ApprovalsEntry,
    } as BadgeOwnershipInfo_ApprovalsEntry;
    if (object.key !== undefined && object.key !== null) {
      message.key = object.key;
    } else {
      message.key = "";
    }
    if (object.value !== undefined && object.value !== null) {
      message.value = object.value;
    } else {
      message.value = 0;
    }
    return message;
  },
};

const basePendingTransfers: object = {
  amount: 0,
  send_request: false,
  to: "",
  from: "",
  memo: "",
};

export const PendingTransfers = {
  encode(message: PendingTransfers, writer: Writer = Writer.create()): Writer {
    if (message.amount !== 0) {
      writer.uint32(8).uint64(message.amount);
    }
    if (message.send_request === true) {
      writer.uint32(16).bool(message.send_request);
    }
    if (message.to !== "") {
      writer.uint32(26).string(message.to);
    }
    if (message.from !== "") {
      writer.uint32(34).string(message.from);
    }
    if (message.memo !== "") {
      writer.uint32(42).string(message.memo);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PendingTransfers {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePendingTransfers } as PendingTransfers;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.send_request = reader.bool();
          break;
        case 3:
          message.to = reader.string();
          break;
        case 4:
          message.from = reader.string();
          break;
        case 5:
          message.memo = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PendingTransfers {
    const message = { ...basePendingTransfers } as PendingTransfers;
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.send_request !== undefined && object.send_request !== null) {
      message.send_request = Boolean(object.send_request);
    } else {
      message.send_request = false;
    }
    if (object.to !== undefined && object.to !== null) {
      message.to = String(object.to);
    } else {
      message.to = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = String(object.from);
    } else {
      message.from = "";
    }
    if (object.memo !== undefined && object.memo !== null) {
      message.memo = String(object.memo);
    } else {
      message.memo = "";
    }
    return message;
  },

  toJSON(message: PendingTransfers): unknown {
    const obj: any = {};
    message.amount !== undefined && (obj.amount = message.amount);
    message.send_request !== undefined &&
      (obj.send_request = message.send_request);
    message.to !== undefined && (obj.to = message.to);
    message.from !== undefined && (obj.from = message.from);
    message.memo !== undefined && (obj.memo = message.memo);
    return obj;
  },

  fromPartial(object: DeepPartial<PendingTransfers>): PendingTransfers {
    const message = { ...basePendingTransfers } as PendingTransfers;
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.send_request !== undefined && object.send_request !== null) {
      message.send_request = object.send_request;
    } else {
      message.send_request = false;
    }
    if (object.to !== undefined && object.to !== null) {
      message.to = object.to;
    } else {
      message.to = "";
    }
    if (object.from !== undefined && object.from !== null) {
      message.from = object.from;
    } else {
      message.from = "";
    }
    if (object.memo !== undefined && object.memo !== null) {
      message.memo = object.memo;
    } else {
      message.memo = "";
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
