/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Balance, UintRange } from "./balances";
import { UserPermissions } from "./permissions";

export const protobufPackage = "bitbadges.bitbadgeschain.badges";

/**
 * UserBalanceStore is the store for the user balances
 * It consists of a list of balances, a list of approved outgoing transfers, and a list of approved incoming transfers,
 * and the permissions for updating the approved incoming/outgoing transfers.
 *
 * The default approved outgoing / incoming transfers are defined by the collection.
 *
 * The outgoing transfers can be used to allow / disallow transfers which are sent from this user.
 * If a transfer has no match, then it is disallowed by default, unless from == initiatedBy (i.e. initiated by this user).
 *
 * The incoming transfers can be used to allow / disallow transfers which are sent to this user.
 * If a transfer has no match, then it is disallowed by default, unless to == initiatedBy (i.e. initiated by this user).
 *
 * Note that the user approved transfers are only checked if the collection approved transfers do not specify to override
 * the user approved transfers.
 */
export interface UserBalanceStore {
  balances: Balance[];
  approvedOutgoingTransfersTimeline: UserApprovedOutgoingTransferTimeline[];
  approvedIncomingTransfersTimeline: UserApprovedIncomingTransferTimeline[];
  permissions: UserPermissions | undefined;
}

export interface UserApprovedOutgoingTransferTimeline {
  approvedOutgoingTransfers: UserApprovedOutgoingTransfer[];
  timelineTimes: UintRange[];
}

export interface UserApprovedIncomingTransferTimeline {
  approvedIncomingTransfers: UserApprovedIncomingTransfer[];
  timelineTimes: UintRange[];
}

/**
 * Challenges define the rules for the approval.
 * If all challenge are not met with valid solutions, then the transfer is not approved.
 *
 * Currently, we only support Merkle tree challenges where the Merkle path must be to the provided root
 * and be the expected length.
 *
 * We also support the following options:
 * -useCreatorAddressAsLeaf: If true, then the leaf will be set to the creator address. Used for whitelist trees.
 * -maxOneUsePerLeaf: If true, then each leaf can only be used once. If false, then the leaf can be used multiple times.
 * This is very important to be set to true if you want to prevent replay attacks.
 * -useLeafIndexForDistributionOrder: If true, we will use the leafIndex to determine the order of the distribution of badges.
 * leafIndex 0 will be the leftmost leaf of the expectedProofLength layer
 *
 * IMPORTANT: We track the number of uses per leaf according to a challenge ID.
 * Please use unique challenge IDs for different challenges of the same timeline.
 * If you update the challenge ID, then the used leaves tracker will reset and start a new tally.
 * It is highly recommended to avoid updating a challenge without resetting the tally via a new challenge ID.
 */
export interface Challenge {
  root: string;
  expectedProofLength: string;
  useCreatorAddressAsLeaf: boolean;
  maxOneUsePerLeaf: boolean;
  useLeafIndexForDistributionOrder: boolean;
  challengeId: string;
}

/** PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address. */
export interface PerAddressApprovals {
  approvalsPerFromAddress: ApprovalsTracker | undefined;
  approvalsPerToAddress: ApprovalsTracker | undefined;
  approvalsPerInitiatedByAddress: ApprovalsTracker | undefined;
}

export interface IsUserOutgoingTransferAllowed {
  invertTo: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  isAllowed: boolean;
}

export interface IsUserIncomingTransferAllowed {
  invertFrom: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  isAllowed: boolean;
}

/**
 * UserApprovedOutgoingTransfer defines the rules for the approval of an outgoing transfer from a user.
 * See CollectionApprovedTransfer for more details. This is the same minus a few fields.
 */
export interface UserApprovedOutgoingTransfer {
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  allowedCombinations: IsUserOutgoingTransferAllowed[];
  challenges: Challenge[];
  trackerId: string;
  incrementBadgeIdsBy: string;
  incrementOwnershipTimesBy: string;
  perAddressApprovals: PerAddressApprovals | undefined;
  uri: string;
  customData: string;
  requireToEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
}

/**
 * UserApprovedIncomingTransfer defines the rules for the approval of an incoming transfer to a user.
 * See CollectionApprovedTransfer for more details. This is the same minus a few fields.
 */
export interface UserApprovedIncomingTransfer {
  fromMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  allowedCombinations: IsUserIncomingTransferAllowed[];
  challenges: Challenge[];
  trackerId: string;
  incrementBadgeIdsBy: string;
  incrementOwnershipTimesBy: string;
  perAddressApprovals: PerAddressApprovals | undefined;
  uri: string;
  customData: string;
  requireFromEqualsInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
}

export interface IsCollectionTransferAllowed {
  invertFrom: boolean;
  invertTo: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  isAllowed: boolean;
}

export interface CollectionApprovedTransfer {
  fromMappingId: string;
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  allowedCombinations: IsCollectionTransferAllowed[];
  challenges: Challenge[];
  trackerId: string;
  incrementBadgeIdsBy: string;
  incrementOwnershipTimesBy: string;
  overallApprovals: ApprovalsTracker | undefined;
  perAddressApprovals: PerAddressApprovals | undefined;
  overridesFromApprovedOutgoingTransfers: boolean;
  overridesToApprovedIncomingTransfers: boolean;
  requireToEqualsInitiatedBy: boolean;
  requireFromEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
  uri: string;
  customData: string;
}

export interface ApprovalsTracker {
  numTransfers: string;
  amounts: Balance[];
}

function createBaseUserBalanceStore(): UserBalanceStore {
  return {
    balances: [],
    approvedOutgoingTransfersTimeline: [],
    approvedIncomingTransfersTimeline: [],
    permissions: undefined,
  };
}

export const UserBalanceStore = {
  encode(message: UserBalanceStore, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.balances) {
      Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.approvedOutgoingTransfersTimeline) {
      UserApprovedOutgoingTransferTimeline.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.approvedIncomingTransfersTimeline) {
      UserApprovedIncomingTransferTimeline.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.permissions !== undefined) {
      UserPermissions.encode(message.permissions, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserBalanceStore {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserBalanceStore();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balances.push(Balance.decode(reader, reader.uint32()));
          break;
        case 2:
          message.approvedOutgoingTransfersTimeline.push(
            UserApprovedOutgoingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 3:
          message.approvedIncomingTransfersTimeline.push(
            UserApprovedIncomingTransferTimeline.decode(reader, reader.uint32()),
          );
          break;
        case 4:
          message.permissions = UserPermissions.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserBalanceStore {
    return {
      balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => Balance.fromJSON(e)) : [],
      approvedOutgoingTransfersTimeline: Array.isArray(object?.approvedOutgoingTransfersTimeline)
        ? object.approvedOutgoingTransfersTimeline.map((e: any) => UserApprovedOutgoingTransferTimeline.fromJSON(e))
        : [],
      approvedIncomingTransfersTimeline: Array.isArray(object?.approvedIncomingTransfersTimeline)
        ? object.approvedIncomingTransfersTimeline.map((e: any) => UserApprovedIncomingTransferTimeline.fromJSON(e))
        : [],
      permissions: isSet(object.permissions) ? UserPermissions.fromJSON(object.permissions) : undefined,
    };
  },

  toJSON(message: UserBalanceStore): unknown {
    const obj: any = {};
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    if (message.approvedOutgoingTransfersTimeline) {
      obj.approvedOutgoingTransfersTimeline = message.approvedOutgoingTransfersTimeline.map((e) =>
        e ? UserApprovedOutgoingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedOutgoingTransfersTimeline = [];
    }
    if (message.approvedIncomingTransfersTimeline) {
      obj.approvedIncomingTransfersTimeline = message.approvedIncomingTransfersTimeline.map((e) =>
        e ? UserApprovedIncomingTransferTimeline.toJSON(e) : undefined
      );
    } else {
      obj.approvedIncomingTransfersTimeline = [];
    }
    message.permissions !== undefined
      && (obj.permissions = message.permissions ? UserPermissions.toJSON(message.permissions) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserBalanceStore>, I>>(object: I): UserBalanceStore {
    const message = createBaseUserBalanceStore();
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.approvedOutgoingTransfersTimeline =
      object.approvedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e)) || [];
    message.approvedIncomingTransfersTimeline =
      object.approvedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e)) || [];
    message.permissions = (object.permissions !== undefined && object.permissions !== null)
      ? UserPermissions.fromPartial(object.permissions)
      : undefined;
    return message;
  },
};

function createBaseUserApprovedOutgoingTransferTimeline(): UserApprovedOutgoingTransferTimeline {
  return { approvedOutgoingTransfers: [], timelineTimes: [] };
}

export const UserApprovedOutgoingTransferTimeline = {
  encode(message: UserApprovedOutgoingTransferTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.approvedOutgoingTransfers) {
      UserApprovedOutgoingTransfer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedOutgoingTransferTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedOutgoingTransferTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvedOutgoingTransfers.push(UserApprovedOutgoingTransfer.decode(reader, reader.uint32()));
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

  fromJSON(object: any): UserApprovedOutgoingTransferTimeline {
    return {
      approvedOutgoingTransfers: Array.isArray(object?.approvedOutgoingTransfers)
        ? object.approvedOutgoingTransfers.map((e: any) => UserApprovedOutgoingTransfer.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedOutgoingTransferTimeline): unknown {
    const obj: any = {};
    if (message.approvedOutgoingTransfers) {
      obj.approvedOutgoingTransfers = message.approvedOutgoingTransfers.map((e) =>
        e ? UserApprovedOutgoingTransfer.toJSON(e) : undefined
      );
    } else {
      obj.approvedOutgoingTransfers = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransferTimeline>, I>>(
    object: I,
  ): UserApprovedOutgoingTransferTimeline {
    const message = createBaseUserApprovedOutgoingTransferTimeline();
    message.approvedOutgoingTransfers =
      object.approvedOutgoingTransfers?.map((e) => UserApprovedOutgoingTransfer.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserApprovedIncomingTransferTimeline(): UserApprovedIncomingTransferTimeline {
  return { approvedIncomingTransfers: [], timelineTimes: [] };
}

export const UserApprovedIncomingTransferTimeline = {
  encode(message: UserApprovedIncomingTransferTimeline, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.approvedIncomingTransfers) {
      UserApprovedIncomingTransfer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.timelineTimes) {
      UintRange.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedIncomingTransferTimeline {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedIncomingTransferTimeline();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvedIncomingTransfers.push(UserApprovedIncomingTransfer.decode(reader, reader.uint32()));
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

  fromJSON(object: any): UserApprovedIncomingTransferTimeline {
    return {
      approvedIncomingTransfers: Array.isArray(object?.approvedIncomingTransfers)
        ? object.approvedIncomingTransfers.map((e: any) => UserApprovedIncomingTransfer.fromJSON(e))
        : [],
      timelineTimes: Array.isArray(object?.timelineTimes)
        ? object.timelineTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
    };
  },

  toJSON(message: UserApprovedIncomingTransferTimeline): unknown {
    const obj: any = {};
    if (message.approvedIncomingTransfers) {
      obj.approvedIncomingTransfers = message.approvedIncomingTransfers.map((e) =>
        e ? UserApprovedIncomingTransfer.toJSON(e) : undefined
      );
    } else {
      obj.approvedIncomingTransfers = [];
    }
    if (message.timelineTimes) {
      obj.timelineTimes = message.timelineTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.timelineTimes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransferTimeline>, I>>(
    object: I,
  ): UserApprovedIncomingTransferTimeline {
    const message = createBaseUserApprovedIncomingTransferTimeline();
    message.approvedIncomingTransfers =
      object.approvedIncomingTransfers?.map((e) => UserApprovedIncomingTransfer.fromPartial(e)) || [];
    message.timelineTimes = object.timelineTimes?.map((e) => UintRange.fromPartial(e)) || [];
    return message;
  },
};

function createBaseChallenge(): Challenge {
  return {
    root: "",
    expectedProofLength: "",
    useCreatorAddressAsLeaf: false,
    maxOneUsePerLeaf: false,
    useLeafIndexForDistributionOrder: false,
    challengeId: "",
  };
}

export const Challenge = {
  encode(message: Challenge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.root !== "") {
      writer.uint32(10).string(message.root);
    }
    if (message.expectedProofLength !== "") {
      writer.uint32(18).string(message.expectedProofLength);
    }
    if (message.useCreatorAddressAsLeaf === true) {
      writer.uint32(24).bool(message.useCreatorAddressAsLeaf);
    }
    if (message.maxOneUsePerLeaf === true) {
      writer.uint32(32).bool(message.maxOneUsePerLeaf);
    }
    if (message.useLeafIndexForDistributionOrder === true) {
      writer.uint32(40).bool(message.useLeafIndexForDistributionOrder);
    }
    if (message.challengeId !== "") {
      writer.uint32(50).string(message.challengeId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Challenge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChallenge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.root = reader.string();
          break;
        case 2:
          message.expectedProofLength = reader.string();
          break;
        case 3:
          message.useCreatorAddressAsLeaf = reader.bool();
          break;
        case 4:
          message.maxOneUsePerLeaf = reader.bool();
          break;
        case 5:
          message.useLeafIndexForDistributionOrder = reader.bool();
          break;
        case 6:
          message.challengeId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Challenge {
    return {
      root: isSet(object.root) ? String(object.root) : "",
      expectedProofLength: isSet(object.expectedProofLength) ? String(object.expectedProofLength) : "",
      useCreatorAddressAsLeaf: isSet(object.useCreatorAddressAsLeaf) ? Boolean(object.useCreatorAddressAsLeaf) : false,
      maxOneUsePerLeaf: isSet(object.maxOneUsePerLeaf) ? Boolean(object.maxOneUsePerLeaf) : false,
      useLeafIndexForDistributionOrder: isSet(object.useLeafIndexForDistributionOrder)
        ? Boolean(object.useLeafIndexForDistributionOrder)
        : false,
      challengeId: isSet(object.challengeId) ? String(object.challengeId) : "",
    };
  },

  toJSON(message: Challenge): unknown {
    const obj: any = {};
    message.root !== undefined && (obj.root = message.root);
    message.expectedProofLength !== undefined && (obj.expectedProofLength = message.expectedProofLength);
    message.useCreatorAddressAsLeaf !== undefined && (obj.useCreatorAddressAsLeaf = message.useCreatorAddressAsLeaf);
    message.maxOneUsePerLeaf !== undefined && (obj.maxOneUsePerLeaf = message.maxOneUsePerLeaf);
    message.useLeafIndexForDistributionOrder !== undefined
      && (obj.useLeafIndexForDistributionOrder = message.useLeafIndexForDistributionOrder);
    message.challengeId !== undefined && (obj.challengeId = message.challengeId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Challenge>, I>>(object: I): Challenge {
    const message = createBaseChallenge();
    message.root = object.root ?? "";
    message.expectedProofLength = object.expectedProofLength ?? "";
    message.useCreatorAddressAsLeaf = object.useCreatorAddressAsLeaf ?? false;
    message.maxOneUsePerLeaf = object.maxOneUsePerLeaf ?? false;
    message.useLeafIndexForDistributionOrder = object.useLeafIndexForDistributionOrder ?? false;
    message.challengeId = object.challengeId ?? "";
    return message;
  },
};

function createBasePerAddressApprovals(): PerAddressApprovals {
  return {
    approvalsPerFromAddress: undefined,
    approvalsPerToAddress: undefined,
    approvalsPerInitiatedByAddress: undefined,
  };
}

export const PerAddressApprovals = {
  encode(message: PerAddressApprovals, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalsPerFromAddress !== undefined) {
      ApprovalsTracker.encode(message.approvalsPerFromAddress, writer.uint32(10).fork()).ldelim();
    }
    if (message.approvalsPerToAddress !== undefined) {
      ApprovalsTracker.encode(message.approvalsPerToAddress, writer.uint32(18).fork()).ldelim();
    }
    if (message.approvalsPerInitiatedByAddress !== undefined) {
      ApprovalsTracker.encode(message.approvalsPerInitiatedByAddress, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerAddressApprovals {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerAddressApprovals();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvalsPerFromAddress = ApprovalsTracker.decode(reader, reader.uint32());
          break;
        case 2:
          message.approvalsPerToAddress = ApprovalsTracker.decode(reader, reader.uint32());
          break;
        case 3:
          message.approvalsPerInitiatedByAddress = ApprovalsTracker.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PerAddressApprovals {
    return {
      approvalsPerFromAddress: isSet(object.approvalsPerFromAddress)
        ? ApprovalsTracker.fromJSON(object.approvalsPerFromAddress)
        : undefined,
      approvalsPerToAddress: isSet(object.approvalsPerToAddress)
        ? ApprovalsTracker.fromJSON(object.approvalsPerToAddress)
        : undefined,
      approvalsPerInitiatedByAddress: isSet(object.approvalsPerInitiatedByAddress)
        ? ApprovalsTracker.fromJSON(object.approvalsPerInitiatedByAddress)
        : undefined,
    };
  },

  toJSON(message: PerAddressApprovals): unknown {
    const obj: any = {};
    message.approvalsPerFromAddress !== undefined && (obj.approvalsPerFromAddress = message.approvalsPerFromAddress
      ? ApprovalsTracker.toJSON(message.approvalsPerFromAddress)
      : undefined);
    message.approvalsPerToAddress !== undefined && (obj.approvalsPerToAddress = message.approvalsPerToAddress
      ? ApprovalsTracker.toJSON(message.approvalsPerToAddress)
      : undefined);
    message.approvalsPerInitiatedByAddress !== undefined
      && (obj.approvalsPerInitiatedByAddress = message.approvalsPerInitiatedByAddress
        ? ApprovalsTracker.toJSON(message.approvalsPerInitiatedByAddress)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PerAddressApprovals>, I>>(object: I): PerAddressApprovals {
    const message = createBasePerAddressApprovals();
    message.approvalsPerFromAddress =
      (object.approvalsPerFromAddress !== undefined && object.approvalsPerFromAddress !== null)
        ? ApprovalsTracker.fromPartial(object.approvalsPerFromAddress)
        : undefined;
    message.approvalsPerToAddress =
      (object.approvalsPerToAddress !== undefined && object.approvalsPerToAddress !== null)
        ? ApprovalsTracker.fromPartial(object.approvalsPerToAddress)
        : undefined;
    message.approvalsPerInitiatedByAddress =
      (object.approvalsPerInitiatedByAddress !== undefined && object.approvalsPerInitiatedByAddress !== null)
        ? ApprovalsTracker.fromPartial(object.approvalsPerInitiatedByAddress)
        : undefined;
    return message;
  },
};

function createBaseIsUserOutgoingTransferAllowed(): IsUserOutgoingTransferAllowed {
  return {
    invertTo: false,
    invertInitiatedBy: false,
    invertTransferTimes: false,
    invertBadgeIds: false,
    isAllowed: false,
  };
}

export const IsUserOutgoingTransferAllowed = {
  encode(message: IsUserOutgoingTransferAllowed, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.invertTo === true) {
      writer.uint32(16).bool(message.invertTo);
    }
    if (message.invertInitiatedBy === true) {
      writer.uint32(24).bool(message.invertInitiatedBy);
    }
    if (message.invertTransferTimes === true) {
      writer.uint32(32).bool(message.invertTransferTimes);
    }
    if (message.invertBadgeIds === true) {
      writer.uint32(40).bool(message.invertBadgeIds);
    }
    if (message.isAllowed === true) {
      writer.uint32(48).bool(message.isAllowed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsUserOutgoingTransferAllowed {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsUserOutgoingTransferAllowed();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.invertTo = reader.bool();
          break;
        case 3:
          message.invertInitiatedBy = reader.bool();
          break;
        case 4:
          message.invertTransferTimes = reader.bool();
          break;
        case 5:
          message.invertBadgeIds = reader.bool();
          break;
        case 6:
          message.isAllowed = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsUserOutgoingTransferAllowed {
    return {
      invertTo: isSet(object.invertTo) ? Boolean(object.invertTo) : false,
      invertInitiatedBy: isSet(object.invertInitiatedBy) ? Boolean(object.invertInitiatedBy) : false,
      invertTransferTimes: isSet(object.invertTransferTimes) ? Boolean(object.invertTransferTimes) : false,
      invertBadgeIds: isSet(object.invertBadgeIds) ? Boolean(object.invertBadgeIds) : false,
      isAllowed: isSet(object.isAllowed) ? Boolean(object.isAllowed) : false,
    };
  },

  toJSON(message: IsUserOutgoingTransferAllowed): unknown {
    const obj: any = {};
    message.invertTo !== undefined && (obj.invertTo = message.invertTo);
    message.invertInitiatedBy !== undefined && (obj.invertInitiatedBy = message.invertInitiatedBy);
    message.invertTransferTimes !== undefined && (obj.invertTransferTimes = message.invertTransferTimes);
    message.invertBadgeIds !== undefined && (obj.invertBadgeIds = message.invertBadgeIds);
    message.isAllowed !== undefined && (obj.isAllowed = message.isAllowed);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsUserOutgoingTransferAllowed>, I>>(
    object: I,
  ): IsUserOutgoingTransferAllowed {
    const message = createBaseIsUserOutgoingTransferAllowed();
    message.invertTo = object.invertTo ?? false;
    message.invertInitiatedBy = object.invertInitiatedBy ?? false;
    message.invertTransferTimes = object.invertTransferTimes ?? false;
    message.invertBadgeIds = object.invertBadgeIds ?? false;
    message.isAllowed = object.isAllowed ?? false;
    return message;
  },
};

function createBaseIsUserIncomingTransferAllowed(): IsUserIncomingTransferAllowed {
  return {
    invertFrom: false,
    invertInitiatedBy: false,
    invertTransferTimes: false,
    invertBadgeIds: false,
    isAllowed: false,
  };
}

export const IsUserIncomingTransferAllowed = {
  encode(message: IsUserIncomingTransferAllowed, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.invertFrom === true) {
      writer.uint32(16).bool(message.invertFrom);
    }
    if (message.invertInitiatedBy === true) {
      writer.uint32(24).bool(message.invertInitiatedBy);
    }
    if (message.invertTransferTimes === true) {
      writer.uint32(32).bool(message.invertTransferTimes);
    }
    if (message.invertBadgeIds === true) {
      writer.uint32(40).bool(message.invertBadgeIds);
    }
    if (message.isAllowed === true) {
      writer.uint32(48).bool(message.isAllowed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsUserIncomingTransferAllowed {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsUserIncomingTransferAllowed();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.invertFrom = reader.bool();
          break;
        case 3:
          message.invertInitiatedBy = reader.bool();
          break;
        case 4:
          message.invertTransferTimes = reader.bool();
          break;
        case 5:
          message.invertBadgeIds = reader.bool();
          break;
        case 6:
          message.isAllowed = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsUserIncomingTransferAllowed {
    return {
      invertFrom: isSet(object.invertFrom) ? Boolean(object.invertFrom) : false,
      invertInitiatedBy: isSet(object.invertInitiatedBy) ? Boolean(object.invertInitiatedBy) : false,
      invertTransferTimes: isSet(object.invertTransferTimes) ? Boolean(object.invertTransferTimes) : false,
      invertBadgeIds: isSet(object.invertBadgeIds) ? Boolean(object.invertBadgeIds) : false,
      isAllowed: isSet(object.isAllowed) ? Boolean(object.isAllowed) : false,
    };
  },

  toJSON(message: IsUserIncomingTransferAllowed): unknown {
    const obj: any = {};
    message.invertFrom !== undefined && (obj.invertFrom = message.invertFrom);
    message.invertInitiatedBy !== undefined && (obj.invertInitiatedBy = message.invertInitiatedBy);
    message.invertTransferTimes !== undefined && (obj.invertTransferTimes = message.invertTransferTimes);
    message.invertBadgeIds !== undefined && (obj.invertBadgeIds = message.invertBadgeIds);
    message.isAllowed !== undefined && (obj.isAllowed = message.isAllowed);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsUserIncomingTransferAllowed>, I>>(
    object: I,
  ): IsUserIncomingTransferAllowed {
    const message = createBaseIsUserIncomingTransferAllowed();
    message.invertFrom = object.invertFrom ?? false;
    message.invertInitiatedBy = object.invertInitiatedBy ?? false;
    message.invertTransferTimes = object.invertTransferTimes ?? false;
    message.invertBadgeIds = object.invertBadgeIds ?? false;
    message.isAllowed = object.isAllowed ?? false;
    return message;
  },
};

function createBaseUserApprovedOutgoingTransfer(): UserApprovedOutgoingTransfer {
  return {
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    allowedCombinations: [],
    challenges: [],
    trackerId: "",
    incrementBadgeIdsBy: "",
    incrementOwnershipTimesBy: "",
    perAddressApprovals: undefined,
    uri: "",
    customData: "",
    requireToEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
  };
}

export const UserApprovedOutgoingTransfer = {
  encode(message: UserApprovedOutgoingTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.toMappingId !== "") {
      writer.uint32(10).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(18).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsUserOutgoingTransferAllowed.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.challenges) {
      Challenge.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.trackerId !== "") {
      writer.uint32(58).string(message.trackerId);
    }
    if (message.incrementBadgeIdsBy !== "") {
      writer.uint32(66).string(message.incrementBadgeIdsBy);
    }
    if (message.incrementOwnershipTimesBy !== "") {
      writer.uint32(74).string(message.incrementOwnershipTimesBy);
    }
    if (message.perAddressApprovals !== undefined) {
      PerAddressApprovals.encode(message.perAddressApprovals, writer.uint32(82).fork()).ldelim();
    }
    if (message.uri !== "") {
      writer.uint32(98).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(106).string(message.customData);
    }
    if (message.requireToEqualsInitiatedBy === true) {
      writer.uint32(112).bool(message.requireToEqualsInitiatedBy);
    }
    if (message.requireToDoesNotEqualInitiatedBy === true) {
      writer.uint32(120).bool(message.requireToDoesNotEqualInitiatedBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedOutgoingTransfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedOutgoingTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.toMappingId = reader.string();
          break;
        case 2:
          message.initiatedByMappingId = reader.string();
          break;
        case 3:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.allowedCombinations.push(IsUserOutgoingTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 6:
          message.challenges.push(Challenge.decode(reader, reader.uint32()));
          break;
        case 7:
          message.trackerId = reader.string();
          break;
        case 8:
          message.incrementBadgeIdsBy = reader.string();
          break;
        case 9:
          message.incrementOwnershipTimesBy = reader.string();
          break;
        case 10:
          message.perAddressApprovals = PerAddressApprovals.decode(reader, reader.uint32());
          break;
        case 12:
          message.uri = reader.string();
          break;
        case 13:
          message.customData = reader.string();
          break;
        case 14:
          message.requireToEqualsInitiatedBy = reader.bool();
          break;
        case 15:
          message.requireToDoesNotEqualInitiatedBy = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedOutgoingTransfer {
    return {
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsUserOutgoingTransferAllowed.fromJSON(e))
        : [],
      challenges: Array.isArray(object?.challenges) ? object.challenges.map((e: any) => Challenge.fromJSON(e)) : [],
      trackerId: isSet(object.trackerId) ? String(object.trackerId) : "",
      incrementBadgeIdsBy: isSet(object.incrementBadgeIdsBy) ? String(object.incrementBadgeIdsBy) : "",
      incrementOwnershipTimesBy: isSet(object.incrementOwnershipTimesBy)
        ? String(object.incrementOwnershipTimesBy)
        : "",
      perAddressApprovals: isSet(object.perAddressApprovals)
        ? PerAddressApprovals.fromJSON(object.perAddressApprovals)
        : undefined,
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      requireToEqualsInitiatedBy: isSet(object.requireToEqualsInitiatedBy)
        ? Boolean(object.requireToEqualsInitiatedBy)
        : false,
      requireToDoesNotEqualInitiatedBy: isSet(object.requireToDoesNotEqualInitiatedBy)
        ? Boolean(object.requireToDoesNotEqualInitiatedBy)
        : false,
    };
  },

  toJSON(message: UserApprovedOutgoingTransfer): unknown {
    const obj: any = {};
    message.toMappingId !== undefined && (obj.toMappingId = message.toMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsUserOutgoingTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.challenges) {
      obj.challenges = message.challenges.map((e) => e ? Challenge.toJSON(e) : undefined);
    } else {
      obj.challenges = [];
    }
    message.trackerId !== undefined && (obj.trackerId = message.trackerId);
    message.incrementBadgeIdsBy !== undefined && (obj.incrementBadgeIdsBy = message.incrementBadgeIdsBy);
    message.incrementOwnershipTimesBy !== undefined
      && (obj.incrementOwnershipTimesBy = message.incrementOwnershipTimesBy);
    message.perAddressApprovals !== undefined && (obj.perAddressApprovals = message.perAddressApprovals
      ? PerAddressApprovals.toJSON(message.perAddressApprovals)
      : undefined);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    message.requireToEqualsInitiatedBy !== undefined
      && (obj.requireToEqualsInitiatedBy = message.requireToEqualsInitiatedBy);
    message.requireToDoesNotEqualInitiatedBy !== undefined
      && (obj.requireToDoesNotEqualInitiatedBy = message.requireToDoesNotEqualInitiatedBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransfer>, I>>(object: I): UserApprovedOutgoingTransfer {
    const message = createBaseUserApprovedOutgoingTransfer();
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsUserOutgoingTransferAllowed.fromPartial(e))
      || [];
    message.challenges = object.challenges?.map((e) => Challenge.fromPartial(e)) || [];
    message.trackerId = object.trackerId ?? "";
    message.incrementBadgeIdsBy = object.incrementBadgeIdsBy ?? "";
    message.incrementOwnershipTimesBy = object.incrementOwnershipTimesBy ?? "";
    message.perAddressApprovals = (object.perAddressApprovals !== undefined && object.perAddressApprovals !== null)
      ? PerAddressApprovals.fromPartial(object.perAddressApprovals)
      : undefined;
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.requireToEqualsInitiatedBy = object.requireToEqualsInitiatedBy ?? false;
    message.requireToDoesNotEqualInitiatedBy = object.requireToDoesNotEqualInitiatedBy ?? false;
    return message;
  },
};

function createBaseUserApprovedIncomingTransfer(): UserApprovedIncomingTransfer {
  return {
    fromMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    allowedCombinations: [],
    challenges: [],
    trackerId: "",
    incrementBadgeIdsBy: "",
    incrementOwnershipTimesBy: "",
    perAddressApprovals: undefined,
    uri: "",
    customData: "",
    requireFromEqualsInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
  };
}

export const UserApprovedIncomingTransfer = {
  encode(message: UserApprovedIncomingTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromMappingId !== "") {
      writer.uint32(10).string(message.fromMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(18).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsUserIncomingTransferAllowed.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.challenges) {
      Challenge.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.trackerId !== "") {
      writer.uint32(58).string(message.trackerId);
    }
    if (message.incrementBadgeIdsBy !== "") {
      writer.uint32(66).string(message.incrementBadgeIdsBy);
    }
    if (message.incrementOwnershipTimesBy !== "") {
      writer.uint32(74).string(message.incrementOwnershipTimesBy);
    }
    if (message.perAddressApprovals !== undefined) {
      PerAddressApprovals.encode(message.perAddressApprovals, writer.uint32(90).fork()).ldelim();
    }
    if (message.uri !== "") {
      writer.uint32(98).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(106).string(message.customData);
    }
    if (message.requireFromEqualsInitiatedBy === true) {
      writer.uint32(112).bool(message.requireFromEqualsInitiatedBy);
    }
    if (message.requireFromDoesNotEqualInitiatedBy === true) {
      writer.uint32(120).bool(message.requireFromDoesNotEqualInitiatedBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserApprovedIncomingTransfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserApprovedIncomingTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromMappingId = reader.string();
          break;
        case 2:
          message.initiatedByMappingId = reader.string();
          break;
        case 3:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 4:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.allowedCombinations.push(IsUserIncomingTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 6:
          message.challenges.push(Challenge.decode(reader, reader.uint32()));
          break;
        case 7:
          message.trackerId = reader.string();
          break;
        case 8:
          message.incrementBadgeIdsBy = reader.string();
          break;
        case 9:
          message.incrementOwnershipTimesBy = reader.string();
          break;
        case 11:
          message.perAddressApprovals = PerAddressApprovals.decode(reader, reader.uint32());
          break;
        case 12:
          message.uri = reader.string();
          break;
        case 13:
          message.customData = reader.string();
          break;
        case 14:
          message.requireFromEqualsInitiatedBy = reader.bool();
          break;
        case 15:
          message.requireFromDoesNotEqualInitiatedBy = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserApprovedIncomingTransfer {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsUserIncomingTransferAllowed.fromJSON(e))
        : [],
      challenges: Array.isArray(object?.challenges) ? object.challenges.map((e: any) => Challenge.fromJSON(e)) : [],
      trackerId: isSet(object.trackerId) ? String(object.trackerId) : "",
      incrementBadgeIdsBy: isSet(object.incrementBadgeIdsBy) ? String(object.incrementBadgeIdsBy) : "",
      incrementOwnershipTimesBy: isSet(object.incrementOwnershipTimesBy)
        ? String(object.incrementOwnershipTimesBy)
        : "",
      perAddressApprovals: isSet(object.perAddressApprovals)
        ? PerAddressApprovals.fromJSON(object.perAddressApprovals)
        : undefined,
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      requireFromEqualsInitiatedBy: isSet(object.requireFromEqualsInitiatedBy)
        ? Boolean(object.requireFromEqualsInitiatedBy)
        : false,
      requireFromDoesNotEqualInitiatedBy: isSet(object.requireFromDoesNotEqualInitiatedBy)
        ? Boolean(object.requireFromDoesNotEqualInitiatedBy)
        : false,
    };
  },

  toJSON(message: UserApprovedIncomingTransfer): unknown {
    const obj: any = {};
    message.fromMappingId !== undefined && (obj.fromMappingId = message.fromMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsUserIncomingTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.challenges) {
      obj.challenges = message.challenges.map((e) => e ? Challenge.toJSON(e) : undefined);
    } else {
      obj.challenges = [];
    }
    message.trackerId !== undefined && (obj.trackerId = message.trackerId);
    message.incrementBadgeIdsBy !== undefined && (obj.incrementBadgeIdsBy = message.incrementBadgeIdsBy);
    message.incrementOwnershipTimesBy !== undefined
      && (obj.incrementOwnershipTimesBy = message.incrementOwnershipTimesBy);
    message.perAddressApprovals !== undefined && (obj.perAddressApprovals = message.perAddressApprovals
      ? PerAddressApprovals.toJSON(message.perAddressApprovals)
      : undefined);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    message.requireFromEqualsInitiatedBy !== undefined
      && (obj.requireFromEqualsInitiatedBy = message.requireFromEqualsInitiatedBy);
    message.requireFromDoesNotEqualInitiatedBy !== undefined
      && (obj.requireFromDoesNotEqualInitiatedBy = message.requireFromDoesNotEqualInitiatedBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransfer>, I>>(object: I): UserApprovedIncomingTransfer {
    const message = createBaseUserApprovedIncomingTransfer();
    message.fromMappingId = object.fromMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsUserIncomingTransferAllowed.fromPartial(e))
      || [];
    message.challenges = object.challenges?.map((e) => Challenge.fromPartial(e)) || [];
    message.trackerId = object.trackerId ?? "";
    message.incrementBadgeIdsBy = object.incrementBadgeIdsBy ?? "";
    message.incrementOwnershipTimesBy = object.incrementOwnershipTimesBy ?? "";
    message.perAddressApprovals = (object.perAddressApprovals !== undefined && object.perAddressApprovals !== null)
      ? PerAddressApprovals.fromPartial(object.perAddressApprovals)
      : undefined;
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.requireFromEqualsInitiatedBy = object.requireFromEqualsInitiatedBy ?? false;
    message.requireFromDoesNotEqualInitiatedBy = object.requireFromDoesNotEqualInitiatedBy ?? false;
    return message;
  },
};

function createBaseIsCollectionTransferAllowed(): IsCollectionTransferAllowed {
  return {
    invertFrom: false,
    invertTo: false,
    invertInitiatedBy: false,
    invertTransferTimes: false,
    invertBadgeIds: false,
    isAllowed: false,
  };
}

export const IsCollectionTransferAllowed = {
  encode(message: IsCollectionTransferAllowed, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.invertFrom === true) {
      writer.uint32(8).bool(message.invertFrom);
    }
    if (message.invertTo === true) {
      writer.uint32(16).bool(message.invertTo);
    }
    if (message.invertInitiatedBy === true) {
      writer.uint32(24).bool(message.invertInitiatedBy);
    }
    if (message.invertTransferTimes === true) {
      writer.uint32(32).bool(message.invertTransferTimes);
    }
    if (message.invertBadgeIds === true) {
      writer.uint32(40).bool(message.invertBadgeIds);
    }
    if (message.isAllowed === true) {
      writer.uint32(48).bool(message.isAllowed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsCollectionTransferAllowed {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsCollectionTransferAllowed();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.invertFrom = reader.bool();
          break;
        case 2:
          message.invertTo = reader.bool();
          break;
        case 3:
          message.invertInitiatedBy = reader.bool();
          break;
        case 4:
          message.invertTransferTimes = reader.bool();
          break;
        case 5:
          message.invertBadgeIds = reader.bool();
          break;
        case 6:
          message.isAllowed = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsCollectionTransferAllowed {
    return {
      invertFrom: isSet(object.invertFrom) ? Boolean(object.invertFrom) : false,
      invertTo: isSet(object.invertTo) ? Boolean(object.invertTo) : false,
      invertInitiatedBy: isSet(object.invertInitiatedBy) ? Boolean(object.invertInitiatedBy) : false,
      invertTransferTimes: isSet(object.invertTransferTimes) ? Boolean(object.invertTransferTimes) : false,
      invertBadgeIds: isSet(object.invertBadgeIds) ? Boolean(object.invertBadgeIds) : false,
      isAllowed: isSet(object.isAllowed) ? Boolean(object.isAllowed) : false,
    };
  },

  toJSON(message: IsCollectionTransferAllowed): unknown {
    const obj: any = {};
    message.invertFrom !== undefined && (obj.invertFrom = message.invertFrom);
    message.invertTo !== undefined && (obj.invertTo = message.invertTo);
    message.invertInitiatedBy !== undefined && (obj.invertInitiatedBy = message.invertInitiatedBy);
    message.invertTransferTimes !== undefined && (obj.invertTransferTimes = message.invertTransferTimes);
    message.invertBadgeIds !== undefined && (obj.invertBadgeIds = message.invertBadgeIds);
    message.isAllowed !== undefined && (obj.isAllowed = message.isAllowed);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsCollectionTransferAllowed>, I>>(object: I): IsCollectionTransferAllowed {
    const message = createBaseIsCollectionTransferAllowed();
    message.invertFrom = object.invertFrom ?? false;
    message.invertTo = object.invertTo ?? false;
    message.invertInitiatedBy = object.invertInitiatedBy ?? false;
    message.invertTransferTimes = object.invertTransferTimes ?? false;
    message.invertBadgeIds = object.invertBadgeIds ?? false;
    message.isAllowed = object.isAllowed ?? false;
    return message;
  },
};

function createBaseCollectionApprovedTransfer(): CollectionApprovedTransfer {
  return {
    fromMappingId: "",
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    allowedCombinations: [],
    challenges: [],
    trackerId: "",
    incrementBadgeIdsBy: "",
    incrementOwnershipTimesBy: "",
    overallApprovals: undefined,
    perAddressApprovals: undefined,
    overridesFromApprovedOutgoingTransfers: false,
    overridesToApprovedIncomingTransfers: false,
    requireToEqualsInitiatedBy: false,
    requireFromEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
    uri: "",
    customData: "",
  };
}

export const CollectionApprovedTransfer = {
  encode(message: CollectionApprovedTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromMappingId !== "") {
      writer.uint32(10).string(message.fromMappingId);
    }
    if (message.toMappingId !== "") {
      writer.uint32(18).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(26).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsCollectionTransferAllowed.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.challenges) {
      Challenge.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    if (message.trackerId !== "") {
      writer.uint32(66).string(message.trackerId);
    }
    if (message.incrementBadgeIdsBy !== "") {
      writer.uint32(74).string(message.incrementBadgeIdsBy);
    }
    if (message.incrementOwnershipTimesBy !== "") {
      writer.uint32(82).string(message.incrementOwnershipTimesBy);
    }
    if (message.overallApprovals !== undefined) {
      ApprovalsTracker.encode(message.overallApprovals, writer.uint32(90).fork()).ldelim();
    }
    if (message.perAddressApprovals !== undefined) {
      PerAddressApprovals.encode(message.perAddressApprovals, writer.uint32(98).fork()).ldelim();
    }
    if (message.overridesFromApprovedOutgoingTransfers === true) {
      writer.uint32(120).bool(message.overridesFromApprovedOutgoingTransfers);
    }
    if (message.overridesToApprovedIncomingTransfers === true) {
      writer.uint32(128).bool(message.overridesToApprovedIncomingTransfers);
    }
    if (message.requireToEqualsInitiatedBy === true) {
      writer.uint32(136).bool(message.requireToEqualsInitiatedBy);
    }
    if (message.requireFromEqualsInitiatedBy === true) {
      writer.uint32(144).bool(message.requireFromEqualsInitiatedBy);
    }
    if (message.requireToDoesNotEqualInitiatedBy === true) {
      writer.uint32(152).bool(message.requireToDoesNotEqualInitiatedBy);
    }
    if (message.requireFromDoesNotEqualInitiatedBy === true) {
      writer.uint32(160).bool(message.requireFromDoesNotEqualInitiatedBy);
    }
    if (message.uri !== "") {
      writer.uint32(170).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(178).string(message.customData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApprovedTransfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApprovedTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fromMappingId = reader.string();
          break;
        case 2:
          message.toMappingId = reader.string();
          break;
        case 3:
          message.initiatedByMappingId = reader.string();
          break;
        case 4:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.allowedCombinations.push(IsCollectionTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 7:
          message.challenges.push(Challenge.decode(reader, reader.uint32()));
          break;
        case 8:
          message.trackerId = reader.string();
          break;
        case 9:
          message.incrementBadgeIdsBy = reader.string();
          break;
        case 10:
          message.incrementOwnershipTimesBy = reader.string();
          break;
        case 11:
          message.overallApprovals = ApprovalsTracker.decode(reader, reader.uint32());
          break;
        case 12:
          message.perAddressApprovals = PerAddressApprovals.decode(reader, reader.uint32());
          break;
        case 15:
          message.overridesFromApprovedOutgoingTransfers = reader.bool();
          break;
        case 16:
          message.overridesToApprovedIncomingTransfers = reader.bool();
          break;
        case 17:
          message.requireToEqualsInitiatedBy = reader.bool();
          break;
        case 18:
          message.requireFromEqualsInitiatedBy = reader.bool();
          break;
        case 19:
          message.requireToDoesNotEqualInitiatedBy = reader.bool();
          break;
        case 20:
          message.requireFromDoesNotEqualInitiatedBy = reader.bool();
          break;
        case 21:
          message.uri = reader.string();
          break;
        case 22:
          message.customData = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApprovedTransfer {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsCollectionTransferAllowed.fromJSON(e))
        : [],
      challenges: Array.isArray(object?.challenges) ? object.challenges.map((e: any) => Challenge.fromJSON(e)) : [],
      trackerId: isSet(object.trackerId) ? String(object.trackerId) : "",
      incrementBadgeIdsBy: isSet(object.incrementBadgeIdsBy) ? String(object.incrementBadgeIdsBy) : "",
      incrementOwnershipTimesBy: isSet(object.incrementOwnershipTimesBy)
        ? String(object.incrementOwnershipTimesBy)
        : "",
      overallApprovals: isSet(object.overallApprovals) ? ApprovalsTracker.fromJSON(object.overallApprovals) : undefined,
      perAddressApprovals: isSet(object.perAddressApprovals)
        ? PerAddressApprovals.fromJSON(object.perAddressApprovals)
        : undefined,
      overridesFromApprovedOutgoingTransfers: isSet(object.overridesFromApprovedOutgoingTransfers)
        ? Boolean(object.overridesFromApprovedOutgoingTransfers)
        : false,
      overridesToApprovedIncomingTransfers: isSet(object.overridesToApprovedIncomingTransfers)
        ? Boolean(object.overridesToApprovedIncomingTransfers)
        : false,
      requireToEqualsInitiatedBy: isSet(object.requireToEqualsInitiatedBy)
        ? Boolean(object.requireToEqualsInitiatedBy)
        : false,
      requireFromEqualsInitiatedBy: isSet(object.requireFromEqualsInitiatedBy)
        ? Boolean(object.requireFromEqualsInitiatedBy)
        : false,
      requireToDoesNotEqualInitiatedBy: isSet(object.requireToDoesNotEqualInitiatedBy)
        ? Boolean(object.requireToDoesNotEqualInitiatedBy)
        : false,
      requireFromDoesNotEqualInitiatedBy: isSet(object.requireFromDoesNotEqualInitiatedBy)
        ? Boolean(object.requireFromDoesNotEqualInitiatedBy)
        : false,
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
    };
  },

  toJSON(message: CollectionApprovedTransfer): unknown {
    const obj: any = {};
    message.fromMappingId !== undefined && (obj.fromMappingId = message.fromMappingId);
    message.toMappingId !== undefined && (obj.toMappingId = message.toMappingId);
    message.initiatedByMappingId !== undefined && (obj.initiatedByMappingId = message.initiatedByMappingId);
    if (message.transferTimes) {
      obj.transferTimes = message.transferTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.transferTimes = [];
    }
    if (message.badgeIds) {
      obj.badgeIds = message.badgeIds.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.badgeIds = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsCollectionTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.challenges) {
      obj.challenges = message.challenges.map((e) => e ? Challenge.toJSON(e) : undefined);
    } else {
      obj.challenges = [];
    }
    message.trackerId !== undefined && (obj.trackerId = message.trackerId);
    message.incrementBadgeIdsBy !== undefined && (obj.incrementBadgeIdsBy = message.incrementBadgeIdsBy);
    message.incrementOwnershipTimesBy !== undefined
      && (obj.incrementOwnershipTimesBy = message.incrementOwnershipTimesBy);
    message.overallApprovals !== undefined && (obj.overallApprovals = message.overallApprovals
      ? ApprovalsTracker.toJSON(message.overallApprovals)
      : undefined);
    message.perAddressApprovals !== undefined && (obj.perAddressApprovals = message.perAddressApprovals
      ? PerAddressApprovals.toJSON(message.perAddressApprovals)
      : undefined);
    message.overridesFromApprovedOutgoingTransfers !== undefined
      && (obj.overridesFromApprovedOutgoingTransfers = message.overridesFromApprovedOutgoingTransfers);
    message.overridesToApprovedIncomingTransfers !== undefined
      && (obj.overridesToApprovedIncomingTransfers = message.overridesToApprovedIncomingTransfers);
    message.requireToEqualsInitiatedBy !== undefined
      && (obj.requireToEqualsInitiatedBy = message.requireToEqualsInitiatedBy);
    message.requireFromEqualsInitiatedBy !== undefined
      && (obj.requireFromEqualsInitiatedBy = message.requireFromEqualsInitiatedBy);
    message.requireToDoesNotEqualInitiatedBy !== undefined
      && (obj.requireToDoesNotEqualInitiatedBy = message.requireToDoesNotEqualInitiatedBy);
    message.requireFromDoesNotEqualInitiatedBy !== undefined
      && (obj.requireFromDoesNotEqualInitiatedBy = message.requireFromDoesNotEqualInitiatedBy);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransfer>, I>>(object: I): CollectionApprovedTransfer {
    const message = createBaseCollectionApprovedTransfer();
    message.fromMappingId = object.fromMappingId ?? "";
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsCollectionTransferAllowed.fromPartial(e))
      || [];
    message.challenges = object.challenges?.map((e) => Challenge.fromPartial(e)) || [];
    message.trackerId = object.trackerId ?? "";
    message.incrementBadgeIdsBy = object.incrementBadgeIdsBy ?? "";
    message.incrementOwnershipTimesBy = object.incrementOwnershipTimesBy ?? "";
    message.overallApprovals = (object.overallApprovals !== undefined && object.overallApprovals !== null)
      ? ApprovalsTracker.fromPartial(object.overallApprovals)
      : undefined;
    message.perAddressApprovals = (object.perAddressApprovals !== undefined && object.perAddressApprovals !== null)
      ? PerAddressApprovals.fromPartial(object.perAddressApprovals)
      : undefined;
    message.overridesFromApprovedOutgoingTransfers = object.overridesFromApprovedOutgoingTransfers ?? false;
    message.overridesToApprovedIncomingTransfers = object.overridesToApprovedIncomingTransfers ?? false;
    message.requireToEqualsInitiatedBy = object.requireToEqualsInitiatedBy ?? false;
    message.requireFromEqualsInitiatedBy = object.requireFromEqualsInitiatedBy ?? false;
    message.requireToDoesNotEqualInitiatedBy = object.requireToDoesNotEqualInitiatedBy ?? false;
    message.requireFromDoesNotEqualInitiatedBy = object.requireFromDoesNotEqualInitiatedBy ?? false;
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    return message;
  },
};

function createBaseApprovalsTracker(): ApprovalsTracker {
  return { numTransfers: "", amounts: [] };
}

export const ApprovalsTracker = {
  encode(message: ApprovalsTracker, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.numTransfers !== "") {
      writer.uint32(10).string(message.numTransfers);
    }
    for (const v of message.amounts) {
      Balance.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalsTracker {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalsTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.numTransfers = reader.string();
          break;
        case 2:
          message.amounts.push(Balance.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalsTracker {
    return {
      numTransfers: isSet(object.numTransfers) ? String(object.numTransfers) : "",
      amounts: Array.isArray(object?.amounts) ? object.amounts.map((e: any) => Balance.fromJSON(e)) : [],
    };
  },

  toJSON(message: ApprovalsTracker): unknown {
    const obj: any = {};
    message.numTransfers !== undefined && (obj.numTransfers = message.numTransfers);
    if (message.amounts) {
      obj.amounts = message.amounts.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.amounts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalsTracker>, I>>(object: I): ApprovalsTracker {
    const message = createBaseApprovalsTracker();
    message.numTransfers = object.numTransfers ?? "";
    message.amounts = object.amounts?.map((e) => Balance.fromPartial(e)) || [];
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
