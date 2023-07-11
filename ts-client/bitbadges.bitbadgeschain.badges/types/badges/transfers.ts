/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Balance, MustOwnBadges, UintRange } from "./balances";
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
  userPermissions: UserPermissions | undefined;
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
export interface MerkleChallenge {
  root: string;
  expectedProofLength: string;
  useCreatorAddressAsLeaf: boolean;
  maxOneUsePerLeaf: boolean;
  useLeafIndexForTransferOrder: boolean;
  challengeId: string;
}

export interface IsUserOutgoingTransferAllowed {
  invertTo: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  invertOwnedTimes: boolean;
  isAllowed: boolean;
}

export interface IsUserIncomingTransferAllowed {
  invertFrom: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  invertOwnedTimes: boolean;
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
  ownedTimes: UintRange[];
  allowedCombinations: IsUserOutgoingTransferAllowed[];
  approvalDetails: OutgoingApprovalDetails[];
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
  ownedTimes: UintRange[];
  allowedCombinations: IsUserIncomingTransferAllowed[];
  approvalDetails: IncomingApprovalDetails[];
}

export interface IsCollectionTransferAllowed {
  invertFrom: boolean;
  invertTo: boolean;
  invertInitiatedBy: boolean;
  invertTransferTimes: boolean;
  invertBadgeIds: boolean;
  invertOwnedTimes: boolean;
  isAllowed: boolean;
}

export interface ManualBalances {
  balances: Balance[];
}

export interface IncrementedBalances {
  startBalances: Balance[];
  incrementBadgeIdsBy: string;
  incrementOwnedTimesBy: string;
}

export interface PredeterminedOrderCalculationMethod {
  useOverallNumTransfers: boolean;
  usePerToAddressNumTransfers: boolean;
  usePerFromAddressNumTransfers: boolean;
  usePerInitiatedByAddressNumTransfers: boolean;
  useMerkleChallengeLeafIndex: boolean;
}

export interface PredeterminedBalances {
  manualBalances: ManualBalances[];
  incrementedBalances: IncrementedBalances | undefined;
  orderCalculationMethod: PredeterminedOrderCalculationMethod | undefined;
}

/** PerAddressApprovals defines the approvals per unique from, to, and/or initiatedBy address. */
export interface ApprovalAmounts {
  overallApprovalAmount: string;
  perToAddressApprovalAmount: string;
  perFromAddressApprovalAmount: string;
  perInitiatedByAddressApprovalAmount: string;
}

export interface MaxNumTransfers {
  overallMaxNumTransfers: string;
  perToAddressMaxNumTransfers: string;
  perFromAddressMaxNumTransfers: string;
  perInitiatedByAddressMaxNumTransfers: string;
}

export interface ApprovalsTracker {
  numTransfers: string;
  amounts: Balance[];
}

export interface ApprovalDetails {
  approvalId: string;
  uri: string;
  customData: string;
  mustOwnBadges: MustOwnBadges[];
  merkleChallenges: MerkleChallenge[];
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireToEqualsInitiatedBy: boolean;
  requireFromEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
  overridesFromApprovedOutgoingTransfers: boolean;
  overridesToApprovedIncomingTransfers: boolean;
}

export interface OutgoingApprovalDetails {
  approvalId: string;
  uri: string;
  customData: string;
  mustOwnBadges: MustOwnBadges[];
  merkleChallenges: MerkleChallenge[];
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireToEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
}

export interface IncomingApprovalDetails {
  approvalId: string;
  uri: string;
  customData: string;
  mustOwnBadges: MustOwnBadges[];
  merkleChallenges: MerkleChallenge[];
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireFromEqualsInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
}

export interface CollectionApprovedTransfer {
  /** Match Criteria */
  fromMappingId: string;
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownedTimes: UintRange[];
  allowedCombinations: IsCollectionTransferAllowed[];
  /** Restrictions (Challenges, Amounts, MaxNumTransfers) */
  approvalDetails: ApprovalDetails[];
}

export interface ApprovalIdDetails {
  approvalId: string;
  /** "collection", "incoming", "outgoing" */
  approvalLevel: string;
  /** Leave blank if approvalLevel == "collection" */
  address: string;
}

export interface Transfer {
  from: string;
  toAddresses: string[];
  balances: Balance[];
  precalculateFromApproval: ApprovalIdDetails | undefined;
  merkleProofs: MerkleProof[];
  memo: string;
}

export interface MerklePathItem {
  aunt: string;
  onRight: boolean;
}

/** Consistent with tendermint/crypto merkle tree */
export interface MerkleProof {
  leaf: string;
  aunts: MerklePathItem[];
}

function createBaseUserBalanceStore(): UserBalanceStore {
  return {
    balances: [],
    approvedOutgoingTransfersTimeline: [],
    approvedIncomingTransfersTimeline: [],
    userPermissions: undefined,
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
    if (message.userPermissions !== undefined) {
      UserPermissions.encode(message.userPermissions, writer.uint32(34).fork()).ldelim();
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
          message.userPermissions = UserPermissions.decode(reader, reader.uint32());
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
      userPermissions: isSet(object.userPermissions) ? UserPermissions.fromJSON(object.userPermissions) : undefined,
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
    message.userPermissions !== undefined
      && (obj.userPermissions = message.userPermissions ? UserPermissions.toJSON(message.userPermissions) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserBalanceStore>, I>>(object: I): UserBalanceStore {
    const message = createBaseUserBalanceStore();
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.approvedOutgoingTransfersTimeline =
      object.approvedOutgoingTransfersTimeline?.map((e) => UserApprovedOutgoingTransferTimeline.fromPartial(e)) || [];
    message.approvedIncomingTransfersTimeline =
      object.approvedIncomingTransfersTimeline?.map((e) => UserApprovedIncomingTransferTimeline.fromPartial(e)) || [];
    message.userPermissions = (object.userPermissions !== undefined && object.userPermissions !== null)
      ? UserPermissions.fromPartial(object.userPermissions)
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

function createBaseMerkleChallenge(): MerkleChallenge {
  return {
    root: "",
    expectedProofLength: "",
    useCreatorAddressAsLeaf: false,
    maxOneUsePerLeaf: false,
    useLeafIndexForTransferOrder: false,
    challengeId: "",
  };
}

export const MerkleChallenge = {
  encode(message: MerkleChallenge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    if (message.useLeafIndexForTransferOrder === true) {
      writer.uint32(40).bool(message.useLeafIndexForTransferOrder);
    }
    if (message.challengeId !== "") {
      writer.uint32(50).string(message.challengeId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MerkleChallenge {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMerkleChallenge();
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
          message.useLeafIndexForTransferOrder = reader.bool();
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

  fromJSON(object: any): MerkleChallenge {
    return {
      root: isSet(object.root) ? String(object.root) : "",
      expectedProofLength: isSet(object.expectedProofLength) ? String(object.expectedProofLength) : "",
      useCreatorAddressAsLeaf: isSet(object.useCreatorAddressAsLeaf) ? Boolean(object.useCreatorAddressAsLeaf) : false,
      maxOneUsePerLeaf: isSet(object.maxOneUsePerLeaf) ? Boolean(object.maxOneUsePerLeaf) : false,
      useLeafIndexForTransferOrder: isSet(object.useLeafIndexForTransferOrder)
        ? Boolean(object.useLeafIndexForTransferOrder)
        : false,
      challengeId: isSet(object.challengeId) ? String(object.challengeId) : "",
    };
  },

  toJSON(message: MerkleChallenge): unknown {
    const obj: any = {};
    message.root !== undefined && (obj.root = message.root);
    message.expectedProofLength !== undefined && (obj.expectedProofLength = message.expectedProofLength);
    message.useCreatorAddressAsLeaf !== undefined && (obj.useCreatorAddressAsLeaf = message.useCreatorAddressAsLeaf);
    message.maxOneUsePerLeaf !== undefined && (obj.maxOneUsePerLeaf = message.maxOneUsePerLeaf);
    message.useLeafIndexForTransferOrder !== undefined
      && (obj.useLeafIndexForTransferOrder = message.useLeafIndexForTransferOrder);
    message.challengeId !== undefined && (obj.challengeId = message.challengeId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MerkleChallenge>, I>>(object: I): MerkleChallenge {
    const message = createBaseMerkleChallenge();
    message.root = object.root ?? "";
    message.expectedProofLength = object.expectedProofLength ?? "";
    message.useCreatorAddressAsLeaf = object.useCreatorAddressAsLeaf ?? false;
    message.maxOneUsePerLeaf = object.maxOneUsePerLeaf ?? false;
    message.useLeafIndexForTransferOrder = object.useLeafIndexForTransferOrder ?? false;
    message.challengeId = object.challengeId ?? "";
    return message;
  },
};

function createBaseIsUserOutgoingTransferAllowed(): IsUserOutgoingTransferAllowed {
  return {
    invertTo: false,
    invertInitiatedBy: false,
    invertTransferTimes: false,
    invertBadgeIds: false,
    invertOwnedTimes: false,
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
    if (message.invertOwnedTimes === true) {
      writer.uint32(48).bool(message.invertOwnedTimes);
    }
    if (message.isAllowed === true) {
      writer.uint32(56).bool(message.isAllowed);
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
          message.invertOwnedTimes = reader.bool();
          break;
        case 7:
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
      invertOwnedTimes: isSet(object.invertOwnedTimes) ? Boolean(object.invertOwnedTimes) : false,
      isAllowed: isSet(object.isAllowed) ? Boolean(object.isAllowed) : false,
    };
  },

  toJSON(message: IsUserOutgoingTransferAllowed): unknown {
    const obj: any = {};
    message.invertTo !== undefined && (obj.invertTo = message.invertTo);
    message.invertInitiatedBy !== undefined && (obj.invertInitiatedBy = message.invertInitiatedBy);
    message.invertTransferTimes !== undefined && (obj.invertTransferTimes = message.invertTransferTimes);
    message.invertBadgeIds !== undefined && (obj.invertBadgeIds = message.invertBadgeIds);
    message.invertOwnedTimes !== undefined && (obj.invertOwnedTimes = message.invertOwnedTimes);
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
    message.invertOwnedTimes = object.invertOwnedTimes ?? false;
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
    invertOwnedTimes: false,
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
    if (message.invertOwnedTimes === true) {
      writer.uint32(48).bool(message.invertOwnedTimes);
    }
    if (message.isAllowed === true) {
      writer.uint32(56).bool(message.isAllowed);
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
          message.invertOwnedTimes = reader.bool();
          break;
        case 7:
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
      invertOwnedTimes: isSet(object.invertOwnedTimes) ? Boolean(object.invertOwnedTimes) : false,
      isAllowed: isSet(object.isAllowed) ? Boolean(object.isAllowed) : false,
    };
  },

  toJSON(message: IsUserIncomingTransferAllowed): unknown {
    const obj: any = {};
    message.invertFrom !== undefined && (obj.invertFrom = message.invertFrom);
    message.invertInitiatedBy !== undefined && (obj.invertInitiatedBy = message.invertInitiatedBy);
    message.invertTransferTimes !== undefined && (obj.invertTransferTimes = message.invertTransferTimes);
    message.invertBadgeIds !== undefined && (obj.invertBadgeIds = message.invertBadgeIds);
    message.invertOwnedTimes !== undefined && (obj.invertOwnedTimes = message.invertOwnedTimes);
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
    message.invertOwnedTimes = object.invertOwnedTimes ?? false;
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
    ownedTimes: [],
    allowedCombinations: [],
    approvalDetails: [],
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
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsUserOutgoingTransferAllowed.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.approvalDetails) {
      OutgoingApprovalDetails.encode(v!, writer.uint32(50).fork()).ldelim();
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
        case 16:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.allowedCombinations.push(IsUserOutgoingTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 6:
          message.approvalDetails.push(OutgoingApprovalDetails.decode(reader, reader.uint32()));
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
      ownedTimes: Array.isArray(object?.ownedTimes) ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsUserOutgoingTransferAllowed.fromJSON(e))
        : [],
      approvalDetails: Array.isArray(object?.approvalDetails)
        ? object.approvalDetails.map((e: any) => OutgoingApprovalDetails.fromJSON(e))
        : [],
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
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsUserOutgoingTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.approvalDetails) {
      obj.approvalDetails = message.approvalDetails.map((e) => e ? OutgoingApprovalDetails.toJSON(e) : undefined);
    } else {
      obj.approvalDetails = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedOutgoingTransfer>, I>>(object: I): UserApprovedOutgoingTransfer {
    const message = createBaseUserApprovedOutgoingTransfer();
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsUserOutgoingTransferAllowed.fromPartial(e))
      || [];
    message.approvalDetails = object.approvalDetails?.map((e) => OutgoingApprovalDetails.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserApprovedIncomingTransfer(): UserApprovedIncomingTransfer {
  return {
    fromMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownedTimes: [],
    allowedCombinations: [],
    approvalDetails: [],
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
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsUserIncomingTransferAllowed.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.approvalDetails) {
      IncomingApprovalDetails.encode(v!, writer.uint32(90).fork()).ldelim();
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
        case 16:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 5:
          message.allowedCombinations.push(IsUserIncomingTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 11:
          message.approvalDetails.push(IncomingApprovalDetails.decode(reader, reader.uint32()));
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
      ownedTimes: Array.isArray(object?.ownedTimes) ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsUserIncomingTransferAllowed.fromJSON(e))
        : [],
      approvalDetails: Array.isArray(object?.approvalDetails)
        ? object.approvalDetails.map((e: any) => IncomingApprovalDetails.fromJSON(e))
        : [],
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
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsUserIncomingTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.approvalDetails) {
      obj.approvalDetails = message.approvalDetails.map((e) => e ? IncomingApprovalDetails.toJSON(e) : undefined);
    } else {
      obj.approvalDetails = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserApprovedIncomingTransfer>, I>>(object: I): UserApprovedIncomingTransfer {
    const message = createBaseUserApprovedIncomingTransfer();
    message.fromMappingId = object.fromMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsUserIncomingTransferAllowed.fromPartial(e))
      || [];
    message.approvalDetails = object.approvalDetails?.map((e) => IncomingApprovalDetails.fromPartial(e)) || [];
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
    invertOwnedTimes: false,
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
    if (message.invertOwnedTimes === true) {
      writer.uint32(48).bool(message.invertOwnedTimes);
    }
    if (message.isAllowed === true) {
      writer.uint32(56).bool(message.isAllowed);
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
          message.invertOwnedTimes = reader.bool();
          break;
        case 7:
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
      invertOwnedTimes: isSet(object.invertOwnedTimes) ? Boolean(object.invertOwnedTimes) : false,
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
    message.invertOwnedTimes !== undefined && (obj.invertOwnedTimes = message.invertOwnedTimes);
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
    message.invertOwnedTimes = object.invertOwnedTimes ?? false;
    message.isAllowed = object.isAllowed ?? false;
    return message;
  },
};

function createBaseManualBalances(): ManualBalances {
  return { balances: [] };
}

export const ManualBalances = {
  encode(message: ManualBalances, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.balances) {
      Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ManualBalances {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseManualBalances();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balances.push(Balance.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ManualBalances {
    return { balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => Balance.fromJSON(e)) : [] };
  },

  toJSON(message: ManualBalances): unknown {
    const obj: any = {};
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ManualBalances>, I>>(object: I): ManualBalances {
    const message = createBaseManualBalances();
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    return message;
  },
};

function createBaseIncrementedBalances(): IncrementedBalances {
  return { startBalances: [], incrementBadgeIdsBy: "", incrementOwnedTimesBy: "" };
}

export const IncrementedBalances = {
  encode(message: IncrementedBalances, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.startBalances) {
      Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.incrementBadgeIdsBy !== "") {
      writer.uint32(18).string(message.incrementBadgeIdsBy);
    }
    if (message.incrementOwnedTimesBy !== "") {
      writer.uint32(26).string(message.incrementOwnedTimesBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IncrementedBalances {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIncrementedBalances();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.startBalances.push(Balance.decode(reader, reader.uint32()));
          break;
        case 2:
          message.incrementBadgeIdsBy = reader.string();
          break;
        case 3:
          message.incrementOwnedTimesBy = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IncrementedBalances {
    return {
      startBalances: Array.isArray(object?.startBalances)
        ? object.startBalances.map((e: any) => Balance.fromJSON(e))
        : [],
      incrementBadgeIdsBy: isSet(object.incrementBadgeIdsBy) ? String(object.incrementBadgeIdsBy) : "",
      incrementOwnedTimesBy: isSet(object.incrementOwnedTimesBy) ? String(object.incrementOwnedTimesBy) : "",
    };
  },

  toJSON(message: IncrementedBalances): unknown {
    const obj: any = {};
    if (message.startBalances) {
      obj.startBalances = message.startBalances.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.startBalances = [];
    }
    message.incrementBadgeIdsBy !== undefined && (obj.incrementBadgeIdsBy = message.incrementBadgeIdsBy);
    message.incrementOwnedTimesBy !== undefined && (obj.incrementOwnedTimesBy = message.incrementOwnedTimesBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IncrementedBalances>, I>>(object: I): IncrementedBalances {
    const message = createBaseIncrementedBalances();
    message.startBalances = object.startBalances?.map((e) => Balance.fromPartial(e)) || [];
    message.incrementBadgeIdsBy = object.incrementBadgeIdsBy ?? "";
    message.incrementOwnedTimesBy = object.incrementOwnedTimesBy ?? "";
    return message;
  },
};

function createBasePredeterminedOrderCalculationMethod(): PredeterminedOrderCalculationMethod {
  return {
    useOverallNumTransfers: false,
    usePerToAddressNumTransfers: false,
    usePerFromAddressNumTransfers: false,
    usePerInitiatedByAddressNumTransfers: false,
    useMerkleChallengeLeafIndex: false,
  };
}

export const PredeterminedOrderCalculationMethod = {
  encode(message: PredeterminedOrderCalculationMethod, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.useOverallNumTransfers === true) {
      writer.uint32(8).bool(message.useOverallNumTransfers);
    }
    if (message.usePerToAddressNumTransfers === true) {
      writer.uint32(16).bool(message.usePerToAddressNumTransfers);
    }
    if (message.usePerFromAddressNumTransfers === true) {
      writer.uint32(24).bool(message.usePerFromAddressNumTransfers);
    }
    if (message.usePerInitiatedByAddressNumTransfers === true) {
      writer.uint32(32).bool(message.usePerInitiatedByAddressNumTransfers);
    }
    if (message.useMerkleChallengeLeafIndex === true) {
      writer.uint32(40).bool(message.useMerkleChallengeLeafIndex);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PredeterminedOrderCalculationMethod {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePredeterminedOrderCalculationMethod();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.useOverallNumTransfers = reader.bool();
          break;
        case 2:
          message.usePerToAddressNumTransfers = reader.bool();
          break;
        case 3:
          message.usePerFromAddressNumTransfers = reader.bool();
          break;
        case 4:
          message.usePerInitiatedByAddressNumTransfers = reader.bool();
          break;
        case 5:
          message.useMerkleChallengeLeafIndex = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PredeterminedOrderCalculationMethod {
    return {
      useOverallNumTransfers: isSet(object.useOverallNumTransfers) ? Boolean(object.useOverallNumTransfers) : false,
      usePerToAddressNumTransfers: isSet(object.usePerToAddressNumTransfers)
        ? Boolean(object.usePerToAddressNumTransfers)
        : false,
      usePerFromAddressNumTransfers: isSet(object.usePerFromAddressNumTransfers)
        ? Boolean(object.usePerFromAddressNumTransfers)
        : false,
      usePerInitiatedByAddressNumTransfers: isSet(object.usePerInitiatedByAddressNumTransfers)
        ? Boolean(object.usePerInitiatedByAddressNumTransfers)
        : false,
      useMerkleChallengeLeafIndex: isSet(object.useMerkleChallengeLeafIndex)
        ? Boolean(object.useMerkleChallengeLeafIndex)
        : false,
    };
  },

  toJSON(message: PredeterminedOrderCalculationMethod): unknown {
    const obj: any = {};
    message.useOverallNumTransfers !== undefined && (obj.useOverallNumTransfers = message.useOverallNumTransfers);
    message.usePerToAddressNumTransfers !== undefined
      && (obj.usePerToAddressNumTransfers = message.usePerToAddressNumTransfers);
    message.usePerFromAddressNumTransfers !== undefined
      && (obj.usePerFromAddressNumTransfers = message.usePerFromAddressNumTransfers);
    message.usePerInitiatedByAddressNumTransfers !== undefined
      && (obj.usePerInitiatedByAddressNumTransfers = message.usePerInitiatedByAddressNumTransfers);
    message.useMerkleChallengeLeafIndex !== undefined
      && (obj.useMerkleChallengeLeafIndex = message.useMerkleChallengeLeafIndex);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PredeterminedOrderCalculationMethod>, I>>(
    object: I,
  ): PredeterminedOrderCalculationMethod {
    const message = createBasePredeterminedOrderCalculationMethod();
    message.useOverallNumTransfers = object.useOverallNumTransfers ?? false;
    message.usePerToAddressNumTransfers = object.usePerToAddressNumTransfers ?? false;
    message.usePerFromAddressNumTransfers = object.usePerFromAddressNumTransfers ?? false;
    message.usePerInitiatedByAddressNumTransfers = object.usePerInitiatedByAddressNumTransfers ?? false;
    message.useMerkleChallengeLeafIndex = object.useMerkleChallengeLeafIndex ?? false;
    return message;
  },
};

function createBasePredeterminedBalances(): PredeterminedBalances {
  return { manualBalances: [], incrementedBalances: undefined, orderCalculationMethod: undefined };
}

export const PredeterminedBalances = {
  encode(message: PredeterminedBalances, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.manualBalances) {
      ManualBalances.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.incrementedBalances !== undefined) {
      IncrementedBalances.encode(message.incrementedBalances, writer.uint32(18).fork()).ldelim();
    }
    if (message.orderCalculationMethod !== undefined) {
      PredeterminedOrderCalculationMethod.encode(message.orderCalculationMethod, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PredeterminedBalances {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePredeterminedBalances();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.manualBalances.push(ManualBalances.decode(reader, reader.uint32()));
          break;
        case 2:
          message.incrementedBalances = IncrementedBalances.decode(reader, reader.uint32());
          break;
        case 3:
          message.orderCalculationMethod = PredeterminedOrderCalculationMethod.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PredeterminedBalances {
    return {
      manualBalances: Array.isArray(object?.manualBalances)
        ? object.manualBalances.map((e: any) => ManualBalances.fromJSON(e))
        : [],
      incrementedBalances: isSet(object.incrementedBalances)
        ? IncrementedBalances.fromJSON(object.incrementedBalances)
        : undefined,
      orderCalculationMethod: isSet(object.orderCalculationMethod)
        ? PredeterminedOrderCalculationMethod.fromJSON(object.orderCalculationMethod)
        : undefined,
    };
  },

  toJSON(message: PredeterminedBalances): unknown {
    const obj: any = {};
    if (message.manualBalances) {
      obj.manualBalances = message.manualBalances.map((e) => e ? ManualBalances.toJSON(e) : undefined);
    } else {
      obj.manualBalances = [];
    }
    message.incrementedBalances !== undefined && (obj.incrementedBalances = message.incrementedBalances
      ? IncrementedBalances.toJSON(message.incrementedBalances)
      : undefined);
    message.orderCalculationMethod !== undefined && (obj.orderCalculationMethod = message.orderCalculationMethod
      ? PredeterminedOrderCalculationMethod.toJSON(message.orderCalculationMethod)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PredeterminedBalances>, I>>(object: I): PredeterminedBalances {
    const message = createBasePredeterminedBalances();
    message.manualBalances = object.manualBalances?.map((e) => ManualBalances.fromPartial(e)) || [];
    message.incrementedBalances = (object.incrementedBalances !== undefined && object.incrementedBalances !== null)
      ? IncrementedBalances.fromPartial(object.incrementedBalances)
      : undefined;
    message.orderCalculationMethod =
      (object.orderCalculationMethod !== undefined && object.orderCalculationMethod !== null)
        ? PredeterminedOrderCalculationMethod.fromPartial(object.orderCalculationMethod)
        : undefined;
    return message;
  },
};

function createBaseApprovalAmounts(): ApprovalAmounts {
  return {
    overallApprovalAmount: "",
    perToAddressApprovalAmount: "",
    perFromAddressApprovalAmount: "",
    perInitiatedByAddressApprovalAmount: "",
  };
}

export const ApprovalAmounts = {
  encode(message: ApprovalAmounts, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.overallApprovalAmount !== "") {
      writer.uint32(10).string(message.overallApprovalAmount);
    }
    if (message.perToAddressApprovalAmount !== "") {
      writer.uint32(18).string(message.perToAddressApprovalAmount);
    }
    if (message.perFromAddressApprovalAmount !== "") {
      writer.uint32(26).string(message.perFromAddressApprovalAmount);
    }
    if (message.perInitiatedByAddressApprovalAmount !== "") {
      writer.uint32(34).string(message.perInitiatedByAddressApprovalAmount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalAmounts {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalAmounts();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.overallApprovalAmount = reader.string();
          break;
        case 2:
          message.perToAddressApprovalAmount = reader.string();
          break;
        case 3:
          message.perFromAddressApprovalAmount = reader.string();
          break;
        case 4:
          message.perInitiatedByAddressApprovalAmount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalAmounts {
    return {
      overallApprovalAmount: isSet(object.overallApprovalAmount) ? String(object.overallApprovalAmount) : "",
      perToAddressApprovalAmount: isSet(object.perToAddressApprovalAmount)
        ? String(object.perToAddressApprovalAmount)
        : "",
      perFromAddressApprovalAmount: isSet(object.perFromAddressApprovalAmount)
        ? String(object.perFromAddressApprovalAmount)
        : "",
      perInitiatedByAddressApprovalAmount: isSet(object.perInitiatedByAddressApprovalAmount)
        ? String(object.perInitiatedByAddressApprovalAmount)
        : "",
    };
  },

  toJSON(message: ApprovalAmounts): unknown {
    const obj: any = {};
    message.overallApprovalAmount !== undefined && (obj.overallApprovalAmount = message.overallApprovalAmount);
    message.perToAddressApprovalAmount !== undefined
      && (obj.perToAddressApprovalAmount = message.perToAddressApprovalAmount);
    message.perFromAddressApprovalAmount !== undefined
      && (obj.perFromAddressApprovalAmount = message.perFromAddressApprovalAmount);
    message.perInitiatedByAddressApprovalAmount !== undefined
      && (obj.perInitiatedByAddressApprovalAmount = message.perInitiatedByAddressApprovalAmount);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalAmounts>, I>>(object: I): ApprovalAmounts {
    const message = createBaseApprovalAmounts();
    message.overallApprovalAmount = object.overallApprovalAmount ?? "";
    message.perToAddressApprovalAmount = object.perToAddressApprovalAmount ?? "";
    message.perFromAddressApprovalAmount = object.perFromAddressApprovalAmount ?? "";
    message.perInitiatedByAddressApprovalAmount = object.perInitiatedByAddressApprovalAmount ?? "";
    return message;
  },
};

function createBaseMaxNumTransfers(): MaxNumTransfers {
  return {
    overallMaxNumTransfers: "",
    perToAddressMaxNumTransfers: "",
    perFromAddressMaxNumTransfers: "",
    perInitiatedByAddressMaxNumTransfers: "",
  };
}

export const MaxNumTransfers = {
  encode(message: MaxNumTransfers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.overallMaxNumTransfers !== "") {
      writer.uint32(10).string(message.overallMaxNumTransfers);
    }
    if (message.perToAddressMaxNumTransfers !== "") {
      writer.uint32(18).string(message.perToAddressMaxNumTransfers);
    }
    if (message.perFromAddressMaxNumTransfers !== "") {
      writer.uint32(26).string(message.perFromAddressMaxNumTransfers);
    }
    if (message.perInitiatedByAddressMaxNumTransfers !== "") {
      writer.uint32(34).string(message.perInitiatedByAddressMaxNumTransfers);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MaxNumTransfers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMaxNumTransfers();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.overallMaxNumTransfers = reader.string();
          break;
        case 2:
          message.perToAddressMaxNumTransfers = reader.string();
          break;
        case 3:
          message.perFromAddressMaxNumTransfers = reader.string();
          break;
        case 4:
          message.perInitiatedByAddressMaxNumTransfers = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MaxNumTransfers {
    return {
      overallMaxNumTransfers: isSet(object.overallMaxNumTransfers) ? String(object.overallMaxNumTransfers) : "",
      perToAddressMaxNumTransfers: isSet(object.perToAddressMaxNumTransfers)
        ? String(object.perToAddressMaxNumTransfers)
        : "",
      perFromAddressMaxNumTransfers: isSet(object.perFromAddressMaxNumTransfers)
        ? String(object.perFromAddressMaxNumTransfers)
        : "",
      perInitiatedByAddressMaxNumTransfers: isSet(object.perInitiatedByAddressMaxNumTransfers)
        ? String(object.perInitiatedByAddressMaxNumTransfers)
        : "",
    };
  },

  toJSON(message: MaxNumTransfers): unknown {
    const obj: any = {};
    message.overallMaxNumTransfers !== undefined && (obj.overallMaxNumTransfers = message.overallMaxNumTransfers);
    message.perToAddressMaxNumTransfers !== undefined
      && (obj.perToAddressMaxNumTransfers = message.perToAddressMaxNumTransfers);
    message.perFromAddressMaxNumTransfers !== undefined
      && (obj.perFromAddressMaxNumTransfers = message.perFromAddressMaxNumTransfers);
    message.perInitiatedByAddressMaxNumTransfers !== undefined
      && (obj.perInitiatedByAddressMaxNumTransfers = message.perInitiatedByAddressMaxNumTransfers);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MaxNumTransfers>, I>>(object: I): MaxNumTransfers {
    const message = createBaseMaxNumTransfers();
    message.overallMaxNumTransfers = object.overallMaxNumTransfers ?? "";
    message.perToAddressMaxNumTransfers = object.perToAddressMaxNumTransfers ?? "";
    message.perFromAddressMaxNumTransfers = object.perFromAddressMaxNumTransfers ?? "";
    message.perInitiatedByAddressMaxNumTransfers = object.perInitiatedByAddressMaxNumTransfers ?? "";
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

function createBaseApprovalDetails(): ApprovalDetails {
  return {
    approvalId: "",
    uri: "",
    customData: "",
    mustOwnBadges: [],
    merkleChallenges: [],
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireToEqualsInitiatedBy: false,
    requireFromEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
    overridesFromApprovedOutgoingTransfers: false,
    overridesToApprovedIncomingTransfers: false,
  };
}

export const ApprovalDetails = {
  encode(message: ApprovalDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(50).string(message.approvalId);
    }
    if (message.uri !== "") {
      writer.uint32(58).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(66).string(message.customData);
    }
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.merkleChallenges) {
      MerkleChallenge.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.predeterminedBalances !== undefined) {
      PredeterminedBalances.encode(message.predeterminedBalances, writer.uint32(26).fork()).ldelim();
    }
    if (message.approvalAmounts !== undefined) {
      ApprovalAmounts.encode(message.approvalAmounts, writer.uint32(34).fork()).ldelim();
    }
    if (message.maxNumTransfers !== undefined) {
      MaxNumTransfers.encode(message.maxNumTransfers, writer.uint32(42).fork()).ldelim();
    }
    if (message.requireToEqualsInitiatedBy === true) {
      writer.uint32(72).bool(message.requireToEqualsInitiatedBy);
    }
    if (message.requireFromEqualsInitiatedBy === true) {
      writer.uint32(80).bool(message.requireFromEqualsInitiatedBy);
    }
    if (message.requireToDoesNotEqualInitiatedBy === true) {
      writer.uint32(88).bool(message.requireToDoesNotEqualInitiatedBy);
    }
    if (message.requireFromDoesNotEqualInitiatedBy === true) {
      writer.uint32(96).bool(message.requireFromDoesNotEqualInitiatedBy);
    }
    if (message.overridesFromApprovedOutgoingTransfers === true) {
      writer.uint32(104).bool(message.overridesFromApprovedOutgoingTransfers);
    }
    if (message.overridesToApprovedIncomingTransfers === true) {
      writer.uint32(112).bool(message.overridesToApprovedIncomingTransfers);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 6:
          message.approvalId = reader.string();
          break;
        case 7:
          message.uri = reader.string();
          break;
        case 8:
          message.customData = reader.string();
          break;
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenges.push(MerkleChallenge.decode(reader, reader.uint32()));
          break;
        case 3:
          message.predeterminedBalances = PredeterminedBalances.decode(reader, reader.uint32());
          break;
        case 4:
          message.approvalAmounts = ApprovalAmounts.decode(reader, reader.uint32());
          break;
        case 5:
          message.maxNumTransfers = MaxNumTransfers.decode(reader, reader.uint32());
          break;
        case 9:
          message.requireToEqualsInitiatedBy = reader.bool();
          break;
        case 10:
          message.requireFromEqualsInitiatedBy = reader.bool();
          break;
        case 11:
          message.requireToDoesNotEqualInitiatedBy = reader.bool();
          break;
        case 12:
          message.requireFromDoesNotEqualInitiatedBy = reader.bool();
          break;
        case 13:
          message.overridesFromApprovedOutgoingTransfers = reader.bool();
          break;
        case 14:
          message.overridesToApprovedIncomingTransfers = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalDetails {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenges: Array.isArray(object?.merkleChallenges)
        ? object.merkleChallenges.map((e: any) => MerkleChallenge.fromJSON(e))
        : [],
      predeterminedBalances: isSet(object.predeterminedBalances)
        ? PredeterminedBalances.fromJSON(object.predeterminedBalances)
        : undefined,
      approvalAmounts: isSet(object.approvalAmounts) ? ApprovalAmounts.fromJSON(object.approvalAmounts) : undefined,
      maxNumTransfers: isSet(object.maxNumTransfers) ? MaxNumTransfers.fromJSON(object.maxNumTransfers) : undefined,
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
      overridesFromApprovedOutgoingTransfers: isSet(object.overridesFromApprovedOutgoingTransfers)
        ? Boolean(object.overridesFromApprovedOutgoingTransfers)
        : false,
      overridesToApprovedIncomingTransfers: isSet(object.overridesToApprovedIncomingTransfers)
        ? Boolean(object.overridesToApprovedIncomingTransfers)
        : false,
    };
  },

  toJSON(message: ApprovalDetails): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    if (message.merkleChallenges) {
      obj.merkleChallenges = message.merkleChallenges.map((e) => e ? MerkleChallenge.toJSON(e) : undefined);
    } else {
      obj.merkleChallenges = [];
    }
    message.predeterminedBalances !== undefined && (obj.predeterminedBalances = message.predeterminedBalances
      ? PredeterminedBalances.toJSON(message.predeterminedBalances)
      : undefined);
    message.approvalAmounts !== undefined
      && (obj.approvalAmounts = message.approvalAmounts ? ApprovalAmounts.toJSON(message.approvalAmounts) : undefined);
    message.maxNumTransfers !== undefined
      && (obj.maxNumTransfers = message.maxNumTransfers ? MaxNumTransfers.toJSON(message.maxNumTransfers) : undefined);
    message.requireToEqualsInitiatedBy !== undefined
      && (obj.requireToEqualsInitiatedBy = message.requireToEqualsInitiatedBy);
    message.requireFromEqualsInitiatedBy !== undefined
      && (obj.requireFromEqualsInitiatedBy = message.requireFromEqualsInitiatedBy);
    message.requireToDoesNotEqualInitiatedBy !== undefined
      && (obj.requireToDoesNotEqualInitiatedBy = message.requireToDoesNotEqualInitiatedBy);
    message.requireFromDoesNotEqualInitiatedBy !== undefined
      && (obj.requireFromDoesNotEqualInitiatedBy = message.requireFromDoesNotEqualInitiatedBy);
    message.overridesFromApprovedOutgoingTransfers !== undefined
      && (obj.overridesFromApprovedOutgoingTransfers = message.overridesFromApprovedOutgoingTransfers);
    message.overridesToApprovedIncomingTransfers !== undefined
      && (obj.overridesToApprovedIncomingTransfers = message.overridesToApprovedIncomingTransfers);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalDetails>, I>>(object: I): ApprovalDetails {
    const message = createBaseApprovalDetails();
    message.approvalId = object.approvalId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenges = object.merkleChallenges?.map((e) => MerkleChallenge.fromPartial(e)) || [];
    message.predeterminedBalances =
      (object.predeterminedBalances !== undefined && object.predeterminedBalances !== null)
        ? PredeterminedBalances.fromPartial(object.predeterminedBalances)
        : undefined;
    message.approvalAmounts = (object.approvalAmounts !== undefined && object.approvalAmounts !== null)
      ? ApprovalAmounts.fromPartial(object.approvalAmounts)
      : undefined;
    message.maxNumTransfers = (object.maxNumTransfers !== undefined && object.maxNumTransfers !== null)
      ? MaxNumTransfers.fromPartial(object.maxNumTransfers)
      : undefined;
    message.requireToEqualsInitiatedBy = object.requireToEqualsInitiatedBy ?? false;
    message.requireFromEqualsInitiatedBy = object.requireFromEqualsInitiatedBy ?? false;
    message.requireToDoesNotEqualInitiatedBy = object.requireToDoesNotEqualInitiatedBy ?? false;
    message.requireFromDoesNotEqualInitiatedBy = object.requireFromDoesNotEqualInitiatedBy ?? false;
    message.overridesFromApprovedOutgoingTransfers = object.overridesFromApprovedOutgoingTransfers ?? false;
    message.overridesToApprovedIncomingTransfers = object.overridesToApprovedIncomingTransfers ?? false;
    return message;
  },
};

function createBaseOutgoingApprovalDetails(): OutgoingApprovalDetails {
  return {
    approvalId: "",
    uri: "",
    customData: "",
    mustOwnBadges: [],
    merkleChallenges: [],
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireToEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
  };
}

export const OutgoingApprovalDetails = {
  encode(message: OutgoingApprovalDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(50).string(message.approvalId);
    }
    if (message.uri !== "") {
      writer.uint32(58).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(66).string(message.customData);
    }
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.merkleChallenges) {
      MerkleChallenge.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.predeterminedBalances !== undefined) {
      PredeterminedBalances.encode(message.predeterminedBalances, writer.uint32(26).fork()).ldelim();
    }
    if (message.approvalAmounts !== undefined) {
      ApprovalAmounts.encode(message.approvalAmounts, writer.uint32(34).fork()).ldelim();
    }
    if (message.maxNumTransfers !== undefined) {
      MaxNumTransfers.encode(message.maxNumTransfers, writer.uint32(42).fork()).ldelim();
    }
    if (message.requireToEqualsInitiatedBy === true) {
      writer.uint32(72).bool(message.requireToEqualsInitiatedBy);
    }
    if (message.requireToDoesNotEqualInitiatedBy === true) {
      writer.uint32(88).bool(message.requireToDoesNotEqualInitiatedBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OutgoingApprovalDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutgoingApprovalDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 6:
          message.approvalId = reader.string();
          break;
        case 7:
          message.uri = reader.string();
          break;
        case 8:
          message.customData = reader.string();
          break;
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenges.push(MerkleChallenge.decode(reader, reader.uint32()));
          break;
        case 3:
          message.predeterminedBalances = PredeterminedBalances.decode(reader, reader.uint32());
          break;
        case 4:
          message.approvalAmounts = ApprovalAmounts.decode(reader, reader.uint32());
          break;
        case 5:
          message.maxNumTransfers = MaxNumTransfers.decode(reader, reader.uint32());
          break;
        case 9:
          message.requireToEqualsInitiatedBy = reader.bool();
          break;
        case 11:
          message.requireToDoesNotEqualInitiatedBy = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): OutgoingApprovalDetails {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenges: Array.isArray(object?.merkleChallenges)
        ? object.merkleChallenges.map((e: any) => MerkleChallenge.fromJSON(e))
        : [],
      predeterminedBalances: isSet(object.predeterminedBalances)
        ? PredeterminedBalances.fromJSON(object.predeterminedBalances)
        : undefined,
      approvalAmounts: isSet(object.approvalAmounts) ? ApprovalAmounts.fromJSON(object.approvalAmounts) : undefined,
      maxNumTransfers: isSet(object.maxNumTransfers) ? MaxNumTransfers.fromJSON(object.maxNumTransfers) : undefined,
      requireToEqualsInitiatedBy: isSet(object.requireToEqualsInitiatedBy)
        ? Boolean(object.requireToEqualsInitiatedBy)
        : false,
      requireToDoesNotEqualInitiatedBy: isSet(object.requireToDoesNotEqualInitiatedBy)
        ? Boolean(object.requireToDoesNotEqualInitiatedBy)
        : false,
    };
  },

  toJSON(message: OutgoingApprovalDetails): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    if (message.merkleChallenges) {
      obj.merkleChallenges = message.merkleChallenges.map((e) => e ? MerkleChallenge.toJSON(e) : undefined);
    } else {
      obj.merkleChallenges = [];
    }
    message.predeterminedBalances !== undefined && (obj.predeterminedBalances = message.predeterminedBalances
      ? PredeterminedBalances.toJSON(message.predeterminedBalances)
      : undefined);
    message.approvalAmounts !== undefined
      && (obj.approvalAmounts = message.approvalAmounts ? ApprovalAmounts.toJSON(message.approvalAmounts) : undefined);
    message.maxNumTransfers !== undefined
      && (obj.maxNumTransfers = message.maxNumTransfers ? MaxNumTransfers.toJSON(message.maxNumTransfers) : undefined);
    message.requireToEqualsInitiatedBy !== undefined
      && (obj.requireToEqualsInitiatedBy = message.requireToEqualsInitiatedBy);
    message.requireToDoesNotEqualInitiatedBy !== undefined
      && (obj.requireToDoesNotEqualInitiatedBy = message.requireToDoesNotEqualInitiatedBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<OutgoingApprovalDetails>, I>>(object: I): OutgoingApprovalDetails {
    const message = createBaseOutgoingApprovalDetails();
    message.approvalId = object.approvalId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenges = object.merkleChallenges?.map((e) => MerkleChallenge.fromPartial(e)) || [];
    message.predeterminedBalances =
      (object.predeterminedBalances !== undefined && object.predeterminedBalances !== null)
        ? PredeterminedBalances.fromPartial(object.predeterminedBalances)
        : undefined;
    message.approvalAmounts = (object.approvalAmounts !== undefined && object.approvalAmounts !== null)
      ? ApprovalAmounts.fromPartial(object.approvalAmounts)
      : undefined;
    message.maxNumTransfers = (object.maxNumTransfers !== undefined && object.maxNumTransfers !== null)
      ? MaxNumTransfers.fromPartial(object.maxNumTransfers)
      : undefined;
    message.requireToEqualsInitiatedBy = object.requireToEqualsInitiatedBy ?? false;
    message.requireToDoesNotEqualInitiatedBy = object.requireToDoesNotEqualInitiatedBy ?? false;
    return message;
  },
};

function createBaseIncomingApprovalDetails(): IncomingApprovalDetails {
  return {
    approvalId: "",
    uri: "",
    customData: "",
    mustOwnBadges: [],
    merkleChallenges: [],
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireFromEqualsInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
  };
}

export const IncomingApprovalDetails = {
  encode(message: IncomingApprovalDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(50).string(message.approvalId);
    }
    if (message.uri !== "") {
      writer.uint32(58).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(66).string(message.customData);
    }
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.merkleChallenges) {
      MerkleChallenge.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.predeterminedBalances !== undefined) {
      PredeterminedBalances.encode(message.predeterminedBalances, writer.uint32(26).fork()).ldelim();
    }
    if (message.approvalAmounts !== undefined) {
      ApprovalAmounts.encode(message.approvalAmounts, writer.uint32(34).fork()).ldelim();
    }
    if (message.maxNumTransfers !== undefined) {
      MaxNumTransfers.encode(message.maxNumTransfers, writer.uint32(42).fork()).ldelim();
    }
    if (message.requireFromEqualsInitiatedBy === true) {
      writer.uint32(80).bool(message.requireFromEqualsInitiatedBy);
    }
    if (message.requireFromDoesNotEqualInitiatedBy === true) {
      writer.uint32(96).bool(message.requireFromDoesNotEqualInitiatedBy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IncomingApprovalDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIncomingApprovalDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 6:
          message.approvalId = reader.string();
          break;
        case 7:
          message.uri = reader.string();
          break;
        case 8:
          message.customData = reader.string();
          break;
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenges.push(MerkleChallenge.decode(reader, reader.uint32()));
          break;
        case 3:
          message.predeterminedBalances = PredeterminedBalances.decode(reader, reader.uint32());
          break;
        case 4:
          message.approvalAmounts = ApprovalAmounts.decode(reader, reader.uint32());
          break;
        case 5:
          message.maxNumTransfers = MaxNumTransfers.decode(reader, reader.uint32());
          break;
        case 10:
          message.requireFromEqualsInitiatedBy = reader.bool();
          break;
        case 12:
          message.requireFromDoesNotEqualInitiatedBy = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IncomingApprovalDetails {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenges: Array.isArray(object?.merkleChallenges)
        ? object.merkleChallenges.map((e: any) => MerkleChallenge.fromJSON(e))
        : [],
      predeterminedBalances: isSet(object.predeterminedBalances)
        ? PredeterminedBalances.fromJSON(object.predeterminedBalances)
        : undefined,
      approvalAmounts: isSet(object.approvalAmounts) ? ApprovalAmounts.fromJSON(object.approvalAmounts) : undefined,
      maxNumTransfers: isSet(object.maxNumTransfers) ? MaxNumTransfers.fromJSON(object.maxNumTransfers) : undefined,
      requireFromEqualsInitiatedBy: isSet(object.requireFromEqualsInitiatedBy)
        ? Boolean(object.requireFromEqualsInitiatedBy)
        : false,
      requireFromDoesNotEqualInitiatedBy: isSet(object.requireFromDoesNotEqualInitiatedBy)
        ? Boolean(object.requireFromDoesNotEqualInitiatedBy)
        : false,
    };
  },

  toJSON(message: IncomingApprovalDetails): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    if (message.merkleChallenges) {
      obj.merkleChallenges = message.merkleChallenges.map((e) => e ? MerkleChallenge.toJSON(e) : undefined);
    } else {
      obj.merkleChallenges = [];
    }
    message.predeterminedBalances !== undefined && (obj.predeterminedBalances = message.predeterminedBalances
      ? PredeterminedBalances.toJSON(message.predeterminedBalances)
      : undefined);
    message.approvalAmounts !== undefined
      && (obj.approvalAmounts = message.approvalAmounts ? ApprovalAmounts.toJSON(message.approvalAmounts) : undefined);
    message.maxNumTransfers !== undefined
      && (obj.maxNumTransfers = message.maxNumTransfers ? MaxNumTransfers.toJSON(message.maxNumTransfers) : undefined);
    message.requireFromEqualsInitiatedBy !== undefined
      && (obj.requireFromEqualsInitiatedBy = message.requireFromEqualsInitiatedBy);
    message.requireFromDoesNotEqualInitiatedBy !== undefined
      && (obj.requireFromDoesNotEqualInitiatedBy = message.requireFromDoesNotEqualInitiatedBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IncomingApprovalDetails>, I>>(object: I): IncomingApprovalDetails {
    const message = createBaseIncomingApprovalDetails();
    message.approvalId = object.approvalId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenges = object.merkleChallenges?.map((e) => MerkleChallenge.fromPartial(e)) || [];
    message.predeterminedBalances =
      (object.predeterminedBalances !== undefined && object.predeterminedBalances !== null)
        ? PredeterminedBalances.fromPartial(object.predeterminedBalances)
        : undefined;
    message.approvalAmounts = (object.approvalAmounts !== undefined && object.approvalAmounts !== null)
      ? ApprovalAmounts.fromPartial(object.approvalAmounts)
      : undefined;
    message.maxNumTransfers = (object.maxNumTransfers !== undefined && object.maxNumTransfers !== null)
      ? MaxNumTransfers.fromPartial(object.maxNumTransfers)
      : undefined;
    message.requireFromEqualsInitiatedBy = object.requireFromEqualsInitiatedBy ?? false;
    message.requireFromDoesNotEqualInitiatedBy = object.requireFromDoesNotEqualInitiatedBy ?? false;
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
    ownedTimes: [],
    allowedCombinations: [],
    approvalDetails: [],
  };
}

export const CollectionApprovedTransfer = {
  encode(message: CollectionApprovedTransfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fromMappingId !== "") {
      writer.uint32(34).string(message.fromMappingId);
    }
    if (message.toMappingId !== "") {
      writer.uint32(42).string(message.toMappingId);
    }
    if (message.initiatedByMappingId !== "") {
      writer.uint32(50).string(message.initiatedByMappingId);
    }
    for (const v of message.transferTimes) {
      UintRange.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.badgeIds) {
      UintRange.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.ownedTimes) {
      UintRange.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.allowedCombinations) {
      IsCollectionTransferAllowed.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.approvalDetails) {
      ApprovalDetails.encode(v!, writer.uint32(90).fork()).ldelim();
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
        case 4:
          message.fromMappingId = reader.string();
          break;
        case 5:
          message.toMappingId = reader.string();
          break;
        case 6:
          message.initiatedByMappingId = reader.string();
          break;
        case 7:
          message.transferTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 8:
          message.badgeIds.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 9:
          message.ownedTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 10:
          message.allowedCombinations.push(IsCollectionTransferAllowed.decode(reader, reader.uint32()));
          break;
        case 11:
          message.approvalDetails.push(ApprovalDetails.decode(reader, reader.uint32()));
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
      ownedTimes: Array.isArray(object?.ownedTimes) ? object.ownedTimes.map((e: any) => UintRange.fromJSON(e)) : [],
      allowedCombinations: Array.isArray(object?.allowedCombinations)
        ? object.allowedCombinations.map((e: any) => IsCollectionTransferAllowed.fromJSON(e))
        : [],
      approvalDetails: Array.isArray(object?.approvalDetails)
        ? object.approvalDetails.map((e: any) => ApprovalDetails.fromJSON(e))
        : [],
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
    if (message.ownedTimes) {
      obj.ownedTimes = message.ownedTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownedTimes = [];
    }
    if (message.allowedCombinations) {
      obj.allowedCombinations = message.allowedCombinations.map((e) =>
        e ? IsCollectionTransferAllowed.toJSON(e) : undefined
      );
    } else {
      obj.allowedCombinations = [];
    }
    if (message.approvalDetails) {
      obj.approvalDetails = message.approvalDetails.map((e) => e ? ApprovalDetails.toJSON(e) : undefined);
    } else {
      obj.approvalDetails = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApprovedTransfer>, I>>(object: I): CollectionApprovedTransfer {
    const message = createBaseCollectionApprovedTransfer();
    message.fromMappingId = object.fromMappingId ?? "";
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownedTimes = object.ownedTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.allowedCombinations = object.allowedCombinations?.map((e) => IsCollectionTransferAllowed.fromPartial(e))
      || [];
    message.approvalDetails = object.approvalDetails?.map((e) => ApprovalDetails.fromPartial(e)) || [];
    return message;
  },
};

function createBaseApprovalIdDetails(): ApprovalIdDetails {
  return { approvalId: "", approvalLevel: "", address: "" };
}

export const ApprovalIdDetails = {
  encode(message: ApprovalIdDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(10).string(message.approvalId);
    }
    if (message.approvalLevel !== "") {
      writer.uint32(18).string(message.approvalLevel);
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalIdDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalIdDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.approvalId = reader.string();
          break;
        case 2:
          message.approvalLevel = reader.string();
          break;
        case 3:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalIdDetails {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      approvalLevel: isSet(object.approvalLevel) ? String(object.approvalLevel) : "",
      address: isSet(object.address) ? String(object.address) : "",
    };
  },

  toJSON(message: ApprovalIdDetails): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.approvalLevel !== undefined && (obj.approvalLevel = message.approvalLevel);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalIdDetails>, I>>(object: I): ApprovalIdDetails {
    const message = createBaseApprovalIdDetails();
    message.approvalId = object.approvalId ?? "";
    message.approvalLevel = object.approvalLevel ?? "";
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseTransfer(): Transfer {
  return { from: "", toAddresses: [], balances: [], precalculateFromApproval: undefined, merkleProofs: [], memo: "" };
}

export const Transfer = {
  encode(message: Transfer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.from !== "") {
      writer.uint32(10).string(message.from);
    }
    for (const v of message.toAddresses) {
      writer.uint32(18).string(v!);
    }
    for (const v of message.balances) {
      Balance.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.precalculateFromApproval !== undefined) {
      ApprovalIdDetails.encode(message.precalculateFromApproval, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.merkleProofs) {
      MerkleProof.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.memo !== "") {
      writer.uint32(50).string(message.memo);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Transfer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.from = reader.string();
          break;
        case 2:
          message.toAddresses.push(reader.string());
          break;
        case 3:
          message.balances.push(Balance.decode(reader, reader.uint32()));
          break;
        case 4:
          message.precalculateFromApproval = ApprovalIdDetails.decode(reader, reader.uint32());
          break;
        case 5:
          message.merkleProofs.push(MerkleProof.decode(reader, reader.uint32()));
          break;
        case 6:
          message.memo = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Transfer {
    return {
      from: isSet(object.from) ? String(object.from) : "",
      toAddresses: Array.isArray(object?.toAddresses) ? object.toAddresses.map((e: any) => String(e)) : [],
      balances: Array.isArray(object?.balances) ? object.balances.map((e: any) => Balance.fromJSON(e)) : [],
      precalculateFromApproval: isSet(object.precalculateFromApproval)
        ? ApprovalIdDetails.fromJSON(object.precalculateFromApproval)
        : undefined,
      merkleProofs: Array.isArray(object?.merkleProofs)
        ? object.merkleProofs.map((e: any) => MerkleProof.fromJSON(e))
        : [],
      memo: isSet(object.memo) ? String(object.memo) : "",
    };
  },

  toJSON(message: Transfer): unknown {
    const obj: any = {};
    message.from !== undefined && (obj.from = message.from);
    if (message.toAddresses) {
      obj.toAddresses = message.toAddresses.map((e) => e);
    } else {
      obj.toAddresses = [];
    }
    if (message.balances) {
      obj.balances = message.balances.map((e) => e ? Balance.toJSON(e) : undefined);
    } else {
      obj.balances = [];
    }
    message.precalculateFromApproval !== undefined && (obj.precalculateFromApproval = message.precalculateFromApproval
      ? ApprovalIdDetails.toJSON(message.precalculateFromApproval)
      : undefined);
    if (message.merkleProofs) {
      obj.merkleProofs = message.merkleProofs.map((e) => e ? MerkleProof.toJSON(e) : undefined);
    } else {
      obj.merkleProofs = [];
    }
    message.memo !== undefined && (obj.memo = message.memo);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Transfer>, I>>(object: I): Transfer {
    const message = createBaseTransfer();
    message.from = object.from ?? "";
    message.toAddresses = object.toAddresses?.map((e) => e) || [];
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.precalculateFromApproval =
      (object.precalculateFromApproval !== undefined && object.precalculateFromApproval !== null)
        ? ApprovalIdDetails.fromPartial(object.precalculateFromApproval)
        : undefined;
    message.merkleProofs = object.merkleProofs?.map((e) => MerkleProof.fromPartial(e)) || [];
    message.memo = object.memo ?? "";
    return message;
  },
};

function createBaseMerklePathItem(): MerklePathItem {
  return { aunt: "", onRight: false };
}

export const MerklePathItem = {
  encode(message: MerklePathItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.aunt !== "") {
      writer.uint32(10).string(message.aunt);
    }
    if (message.onRight === true) {
      writer.uint32(16).bool(message.onRight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MerklePathItem {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMerklePathItem();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.aunt = reader.string();
          break;
        case 2:
          message.onRight = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MerklePathItem {
    return {
      aunt: isSet(object.aunt) ? String(object.aunt) : "",
      onRight: isSet(object.onRight) ? Boolean(object.onRight) : false,
    };
  },

  toJSON(message: MerklePathItem): unknown {
    const obj: any = {};
    message.aunt !== undefined && (obj.aunt = message.aunt);
    message.onRight !== undefined && (obj.onRight = message.onRight);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MerklePathItem>, I>>(object: I): MerklePathItem {
    const message = createBaseMerklePathItem();
    message.aunt = object.aunt ?? "";
    message.onRight = object.onRight ?? false;
    return message;
  },
};

function createBaseMerkleProof(): MerkleProof {
  return { leaf: "", aunts: [] };
}

export const MerkleProof = {
  encode(message: MerkleProof, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.leaf !== "") {
      writer.uint32(10).string(message.leaf);
    }
    for (const v of message.aunts) {
      MerklePathItem.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MerkleProof {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMerkleProof();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.leaf = reader.string();
          break;
        case 2:
          message.aunts.push(MerklePathItem.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MerkleProof {
    return {
      leaf: isSet(object.leaf) ? String(object.leaf) : "",
      aunts: Array.isArray(object?.aunts) ? object.aunts.map((e: any) => MerklePathItem.fromJSON(e)) : [],
    };
  },

  toJSON(message: MerkleProof): unknown {
    const obj: any = {};
    message.leaf !== undefined && (obj.leaf = message.leaf);
    if (message.aunts) {
      obj.aunts = message.aunts.map((e) => e ? MerklePathItem.toJSON(e) : undefined);
    } else {
      obj.aunts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MerkleProof>, I>>(object: I): MerkleProof {
    const message = createBaseMerkleProof();
    message.leaf = object.leaf ?? "";
    message.aunts = object.aunts?.map((e) => MerklePathItem.fromPartial(e)) || [];
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
