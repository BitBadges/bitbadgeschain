/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Balance, MustOwnBadges, UintRange } from "./balances";
import { UserPermissions } from "./permissions";

export const protobufPackage = "badges";

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
  outgoingApprovals: UserOutgoingApproval[];
  incomingApprovals: UserIncomingApproval[];
  autoApproveSelfInitiatedOutgoingTransfers: boolean;
  autoApproveSelfInitiatedIncomingTransfers: boolean;
  userPermissions: UserPermissions | undefined;
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
  maxUsesPerLeaf: string;
  uri: string;
  customData: string;
}

/**
 * UserOutgoingApproval defines the rules for the approval of an outgoing transfer from a user.
 * See CollectionApproval for more details. This is the same minus a few fields.
 */
export interface UserOutgoingApproval {
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  amountTrackerId: string;
  challengeTrackerId: string;
  /** if approved, we use these. if not, these are ignored */
  uri: string;
  customData: string;
  approvalId: string;
  approvalCriteria: OutgoingApprovalCriteria | undefined;
}

/**
 * UserIncomingApproval defines the rules for the approval of an incoming transfer to a user.
 * See CollectionApproval for more details. This is the same minus a few fields.
 */
export interface UserIncomingApproval {
  fromMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  /** if applicable */
  amountTrackerId: string;
  /** if applicable */
  challengeTrackerId: string;
  uri: string;
  customData: string;
  /** if applicable */
  approvalId: string;
  approvalCriteria: IncomingApprovalCriteria | undefined;
}

export interface ManualBalances {
  balances: Balance[];
}

export interface IncrementedBalances {
  startBalances: Balance[];
  incrementBadgeIdsBy: string;
  incrementOwnershipTimesBy: string;
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

export interface ApprovalCriteria {
  mustOwnBadges: MustOwnBadges[];
  merkleChallenge: MerkleChallenge | undefined;
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireToEqualsInitiatedBy: boolean;
  requireFromEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
  overridesFromOutgoingApprovals: boolean;
  overridesToIncomingApprovals: boolean;
}

export interface OutgoingApprovalCriteria {
  mustOwnBadges: MustOwnBadges[];
  merkleChallenge: MerkleChallenge | undefined;
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireToEqualsInitiatedBy: boolean;
  requireToDoesNotEqualInitiatedBy: boolean;
}

export interface IncomingApprovalCriteria {
  mustOwnBadges: MustOwnBadges[];
  merkleChallenge: MerkleChallenge | undefined;
  predeterminedBalances: PredeterminedBalances | undefined;
  approvalAmounts: ApprovalAmounts | undefined;
  maxNumTransfers: MaxNumTransfers | undefined;
  requireFromEqualsInitiatedBy: boolean;
  requireFromDoesNotEqualInitiatedBy: boolean;
}

export interface CollectionApproval {
  /** Match Criteria */
  fromMappingId: string;
  toMappingId: string;
  initiatedByMappingId: string;
  transferTimes: UintRange[];
  badgeIds: UintRange[];
  ownershipTimes: UintRange[];
  /** if applicable */
  amountTrackerId: string;
  /** if applicable */
  challengeTrackerId: string;
  uri: string;
  customData: string;
  /** if applicable */
  approvalId: string;
  approvalCriteria: ApprovalCriteria | undefined;
}

export interface ApprovalIdentifierDetails {
  approvalId: string;
  /** "collection", "incoming", "outgoing" */
  approvalLevel: string;
  /** Leave blank if approvalLevel == "collection" */
  approverAddress: string;
}

export interface Transfer {
  from: string;
  toAddresses: string[];
  balances: Balance[];
  precalculateBalancesFromApproval: ApprovalIdentifierDetails | undefined;
  merkleProofs: MerkleProof[];
  memo: string;
  prioritizedApprovals: ApprovalIdentifierDetails[];
  onlyCheckPrioritizedApprovals: boolean;
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
    outgoingApprovals: [],
    incomingApprovals: [],
    autoApproveSelfInitiatedOutgoingTransfers: false,
    autoApproveSelfInitiatedIncomingTransfers: false,
    userPermissions: undefined,
  };
}

export const UserBalanceStore = {
  encode(message: UserBalanceStore, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.balances) {
      Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.outgoingApprovals) {
      UserOutgoingApproval.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.incomingApprovals) {
      UserIncomingApproval.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.autoApproveSelfInitiatedOutgoingTransfers === true) {
      writer.uint32(32).bool(message.autoApproveSelfInitiatedOutgoingTransfers);
    }
    if (message.autoApproveSelfInitiatedIncomingTransfers === true) {
      writer.uint32(40).bool(message.autoApproveSelfInitiatedIncomingTransfers);
    }
    if (message.userPermissions !== undefined) {
      UserPermissions.encode(message.userPermissions, writer.uint32(50).fork()).ldelim();
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
          message.outgoingApprovals.push(UserOutgoingApproval.decode(reader, reader.uint32()));
          break;
        case 3:
          message.incomingApprovals.push(UserIncomingApproval.decode(reader, reader.uint32()));
          break;
        case 4:
          message.autoApproveSelfInitiatedOutgoingTransfers = reader.bool();
          break;
        case 5:
          message.autoApproveSelfInitiatedIncomingTransfers = reader.bool();
          break;
        case 6:
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
      outgoingApprovals: Array.isArray(object?.outgoingApprovals)
        ? object.outgoingApprovals.map((e: any) => UserOutgoingApproval.fromJSON(e))
        : [],
      incomingApprovals: Array.isArray(object?.incomingApprovals)
        ? object.incomingApprovals.map((e: any) => UserIncomingApproval.fromJSON(e))
        : [],
      autoApproveSelfInitiatedOutgoingTransfers: isSet(object.autoApproveSelfInitiatedOutgoingTransfers)
        ? Boolean(object.autoApproveSelfInitiatedOutgoingTransfers)
        : false,
      autoApproveSelfInitiatedIncomingTransfers: isSet(object.autoApproveSelfInitiatedIncomingTransfers)
        ? Boolean(object.autoApproveSelfInitiatedIncomingTransfers)
        : false,
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
    if (message.outgoingApprovals) {
      obj.outgoingApprovals = message.outgoingApprovals.map((e) => e ? UserOutgoingApproval.toJSON(e) : undefined);
    } else {
      obj.outgoingApprovals = [];
    }
    if (message.incomingApprovals) {
      obj.incomingApprovals = message.incomingApprovals.map((e) => e ? UserIncomingApproval.toJSON(e) : undefined);
    } else {
      obj.incomingApprovals = [];
    }
    message.autoApproveSelfInitiatedOutgoingTransfers !== undefined
      && (obj.autoApproveSelfInitiatedOutgoingTransfers = message.autoApproveSelfInitiatedOutgoingTransfers);
    message.autoApproveSelfInitiatedIncomingTransfers !== undefined
      && (obj.autoApproveSelfInitiatedIncomingTransfers = message.autoApproveSelfInitiatedIncomingTransfers);
    message.userPermissions !== undefined
      && (obj.userPermissions = message.userPermissions ? UserPermissions.toJSON(message.userPermissions) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserBalanceStore>, I>>(object: I): UserBalanceStore {
    const message = createBaseUserBalanceStore();
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.outgoingApprovals = object.outgoingApprovals?.map((e) => UserOutgoingApproval.fromPartial(e)) || [];
    message.incomingApprovals = object.incomingApprovals?.map((e) => UserIncomingApproval.fromPartial(e)) || [];
    message.autoApproveSelfInitiatedOutgoingTransfers = object.autoApproveSelfInitiatedOutgoingTransfers ?? false;
    message.autoApproveSelfInitiatedIncomingTransfers = object.autoApproveSelfInitiatedIncomingTransfers ?? false;
    message.userPermissions = (object.userPermissions !== undefined && object.userPermissions !== null)
      ? UserPermissions.fromPartial(object.userPermissions)
      : undefined;
    return message;
  },
};

function createBaseMerkleChallenge(): MerkleChallenge {
  return {
    root: "",
    expectedProofLength: "",
    useCreatorAddressAsLeaf: false,
    maxUsesPerLeaf: "",
    uri: "",
    customData: "",
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
    if (message.maxUsesPerLeaf !== "") {
      writer.uint32(34).string(message.maxUsesPerLeaf);
    }
    if (message.uri !== "") {
      writer.uint32(50).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(58).string(message.customData);
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
          message.maxUsesPerLeaf = reader.string();
          break;
        case 6:
          message.uri = reader.string();
          break;
        case 7:
          message.customData = reader.string();
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
      maxUsesPerLeaf: isSet(object.maxUsesPerLeaf) ? String(object.maxUsesPerLeaf) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
    };
  },

  toJSON(message: MerkleChallenge): unknown {
    const obj: any = {};
    message.root !== undefined && (obj.root = message.root);
    message.expectedProofLength !== undefined && (obj.expectedProofLength = message.expectedProofLength);
    message.useCreatorAddressAsLeaf !== undefined && (obj.useCreatorAddressAsLeaf = message.useCreatorAddressAsLeaf);
    message.maxUsesPerLeaf !== undefined && (obj.maxUsesPerLeaf = message.maxUsesPerLeaf);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MerkleChallenge>, I>>(object: I): MerkleChallenge {
    const message = createBaseMerkleChallenge();
    message.root = object.root ?? "";
    message.expectedProofLength = object.expectedProofLength ?? "";
    message.useCreatorAddressAsLeaf = object.useCreatorAddressAsLeaf ?? false;
    message.maxUsesPerLeaf = object.maxUsesPerLeaf ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    return message;
  },
};

function createBaseUserOutgoingApproval(): UserOutgoingApproval {
  return {
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    uri: "",
    customData: "",
    approvalId: "",
    approvalCriteria: undefined,
  };
}

export const UserOutgoingApproval = {
  encode(message: UserOutgoingApproval, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(50).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(58).string(message.challengeTrackerId);
    }
    if (message.uri !== "") {
      writer.uint32(66).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(74).string(message.customData);
    }
    if (message.approvalId !== "") {
      writer.uint32(82).string(message.approvalId);
    }
    if (message.approvalCriteria !== undefined) {
      OutgoingApprovalCriteria.encode(message.approvalCriteria, writer.uint32(90).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserOutgoingApproval {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserOutgoingApproval();
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
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.amountTrackerId = reader.string();
          break;
        case 7:
          message.challengeTrackerId = reader.string();
          break;
        case 8:
          message.uri = reader.string();
          break;
        case 9:
          message.customData = reader.string();
          break;
        case 10:
          message.approvalId = reader.string();
          break;
        case 11:
          message.approvalCriteria = OutgoingApprovalCriteria.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserOutgoingApproval {
    return {
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      approvalCriteria: isSet(object.approvalCriteria)
        ? OutgoingApprovalCriteria.fromJSON(object.approvalCriteria)
        : undefined,
    };
  },

  toJSON(message: UserOutgoingApproval): unknown {
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.approvalCriteria !== undefined && (obj.approvalCriteria = message.approvalCriteria
      ? OutgoingApprovalCriteria.toJSON(message.approvalCriteria)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserOutgoingApproval>, I>>(object: I): UserOutgoingApproval {
    const message = createBaseUserOutgoingApproval();
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.approvalId = object.approvalId ?? "";
    message.approvalCriteria = (object.approvalCriteria !== undefined && object.approvalCriteria !== null)
      ? OutgoingApprovalCriteria.fromPartial(object.approvalCriteria)
      : undefined;
    return message;
  },
};

function createBaseUserIncomingApproval(): UserIncomingApproval {
  return {
    fromMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    uri: "",
    customData: "",
    approvalId: "",
    approvalCriteria: undefined,
  };
}

export const UserIncomingApproval = {
  encode(message: UserIncomingApproval, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(50).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(58).string(message.challengeTrackerId);
    }
    if (message.uri !== "") {
      writer.uint32(66).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(74).string(message.customData);
    }
    if (message.approvalId !== "") {
      writer.uint32(82).string(message.approvalId);
    }
    if (message.approvalCriteria !== undefined) {
      IncomingApprovalCriteria.encode(message.approvalCriteria, writer.uint32(90).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserIncomingApproval {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserIncomingApproval();
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
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 6:
          message.amountTrackerId = reader.string();
          break;
        case 7:
          message.challengeTrackerId = reader.string();
          break;
        case 8:
          message.uri = reader.string();
          break;
        case 9:
          message.customData = reader.string();
          break;
        case 10:
          message.approvalId = reader.string();
          break;
        case 11:
          message.approvalCriteria = IncomingApprovalCriteria.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserIncomingApproval {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      approvalCriteria: isSet(object.approvalCriteria)
        ? IncomingApprovalCriteria.fromJSON(object.approvalCriteria)
        : undefined,
    };
  },

  toJSON(message: UserIncomingApproval): unknown {
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.approvalCriteria !== undefined && (obj.approvalCriteria = message.approvalCriteria
      ? IncomingApprovalCriteria.toJSON(message.approvalCriteria)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserIncomingApproval>, I>>(object: I): UserIncomingApproval {
    const message = createBaseUserIncomingApproval();
    message.fromMappingId = object.fromMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.approvalId = object.approvalId ?? "";
    message.approvalCriteria = (object.approvalCriteria !== undefined && object.approvalCriteria !== null)
      ? IncomingApprovalCriteria.fromPartial(object.approvalCriteria)
      : undefined;
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
  return { startBalances: [], incrementBadgeIdsBy: "", incrementOwnershipTimesBy: "" };
}

export const IncrementedBalances = {
  encode(message: IncrementedBalances, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.startBalances) {
      Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.incrementBadgeIdsBy !== "") {
      writer.uint32(18).string(message.incrementBadgeIdsBy);
    }
    if (message.incrementOwnershipTimesBy !== "") {
      writer.uint32(26).string(message.incrementOwnershipTimesBy);
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
          message.incrementOwnershipTimesBy = reader.string();
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
      incrementOwnershipTimesBy: isSet(object.incrementOwnershipTimesBy)
        ? String(object.incrementOwnershipTimesBy)
        : "",
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
    message.incrementOwnershipTimesBy !== undefined
      && (obj.incrementOwnershipTimesBy = message.incrementOwnershipTimesBy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IncrementedBalances>, I>>(object: I): IncrementedBalances {
    const message = createBaseIncrementedBalances();
    message.startBalances = object.startBalances?.map((e) => Balance.fromPartial(e)) || [];
    message.incrementBadgeIdsBy = object.incrementBadgeIdsBy ?? "";
    message.incrementOwnershipTimesBy = object.incrementOwnershipTimesBy ?? "";
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

function createBaseApprovalCriteria(): ApprovalCriteria {
  return {
    mustOwnBadges: [],
    merkleChallenge: undefined,
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireToEqualsInitiatedBy: false,
    requireFromEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
    overridesFromOutgoingApprovals: false,
    overridesToIncomingApprovals: false,
  };
}

export const ApprovalCriteria = {
  encode(message: ApprovalCriteria, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.merkleChallenge !== undefined) {
      MerkleChallenge.encode(message.merkleChallenge, writer.uint32(18).fork()).ldelim();
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
    if (message.overridesFromOutgoingApprovals === true) {
      writer.uint32(104).bool(message.overridesFromOutgoingApprovals);
    }
    if (message.overridesToIncomingApprovals === true) {
      writer.uint32(112).bool(message.overridesToIncomingApprovals);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalCriteria {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalCriteria();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenge = MerkleChallenge.decode(reader, reader.uint32());
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
          message.overridesFromOutgoingApprovals = reader.bool();
          break;
        case 14:
          message.overridesToIncomingApprovals = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalCriteria {
    return {
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenge: isSet(object.merkleChallenge) ? MerkleChallenge.fromJSON(object.merkleChallenge) : undefined,
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
      overridesFromOutgoingApprovals: isSet(object.overridesFromOutgoingApprovals)
        ? Boolean(object.overridesFromOutgoingApprovals)
        : false,
      overridesToIncomingApprovals: isSet(object.overridesToIncomingApprovals)
        ? Boolean(object.overridesToIncomingApprovals)
        : false,
    };
  },

  toJSON(message: ApprovalCriteria): unknown {
    const obj: any = {};
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    message.merkleChallenge !== undefined
      && (obj.merkleChallenge = message.merkleChallenge ? MerkleChallenge.toJSON(message.merkleChallenge) : undefined);
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
    message.overridesFromOutgoingApprovals !== undefined
      && (obj.overridesFromOutgoingApprovals = message.overridesFromOutgoingApprovals);
    message.overridesToIncomingApprovals !== undefined
      && (obj.overridesToIncomingApprovals = message.overridesToIncomingApprovals);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalCriteria>, I>>(object: I): ApprovalCriteria {
    const message = createBaseApprovalCriteria();
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenge = (object.merkleChallenge !== undefined && object.merkleChallenge !== null)
      ? MerkleChallenge.fromPartial(object.merkleChallenge)
      : undefined;
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
    message.overridesFromOutgoingApprovals = object.overridesFromOutgoingApprovals ?? false;
    message.overridesToIncomingApprovals = object.overridesToIncomingApprovals ?? false;
    return message;
  },
};

function createBaseOutgoingApprovalCriteria(): OutgoingApprovalCriteria {
  return {
    mustOwnBadges: [],
    merkleChallenge: undefined,
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireToEqualsInitiatedBy: false,
    requireToDoesNotEqualInitiatedBy: false,
  };
}

export const OutgoingApprovalCriteria = {
  encode(message: OutgoingApprovalCriteria, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.merkleChallenge !== undefined) {
      MerkleChallenge.encode(message.merkleChallenge, writer.uint32(18).fork()).ldelim();
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

  decode(input: _m0.Reader | Uint8Array, length?: number): OutgoingApprovalCriteria {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutgoingApprovalCriteria();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenge = MerkleChallenge.decode(reader, reader.uint32());
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

  fromJSON(object: any): OutgoingApprovalCriteria {
    return {
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenge: isSet(object.merkleChallenge) ? MerkleChallenge.fromJSON(object.merkleChallenge) : undefined,
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

  toJSON(message: OutgoingApprovalCriteria): unknown {
    const obj: any = {};
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    message.merkleChallenge !== undefined
      && (obj.merkleChallenge = message.merkleChallenge ? MerkleChallenge.toJSON(message.merkleChallenge) : undefined);
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

  fromPartial<I extends Exact<DeepPartial<OutgoingApprovalCriteria>, I>>(object: I): OutgoingApprovalCriteria {
    const message = createBaseOutgoingApprovalCriteria();
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenge = (object.merkleChallenge !== undefined && object.merkleChallenge !== null)
      ? MerkleChallenge.fromPartial(object.merkleChallenge)
      : undefined;
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

function createBaseIncomingApprovalCriteria(): IncomingApprovalCriteria {
  return {
    mustOwnBadges: [],
    merkleChallenge: undefined,
    predeterminedBalances: undefined,
    approvalAmounts: undefined,
    maxNumTransfers: undefined,
    requireFromEqualsInitiatedBy: false,
    requireFromDoesNotEqualInitiatedBy: false,
  };
}

export const IncomingApprovalCriteria = {
  encode(message: IncomingApprovalCriteria, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.mustOwnBadges) {
      MustOwnBadges.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.merkleChallenge !== undefined) {
      MerkleChallenge.encode(message.merkleChallenge, writer.uint32(18).fork()).ldelim();
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

  decode(input: _m0.Reader | Uint8Array, length?: number): IncomingApprovalCriteria {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIncomingApprovalCriteria();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mustOwnBadges.push(MustOwnBadges.decode(reader, reader.uint32()));
          break;
        case 2:
          message.merkleChallenge = MerkleChallenge.decode(reader, reader.uint32());
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

  fromJSON(object: any): IncomingApprovalCriteria {
    return {
      mustOwnBadges: Array.isArray(object?.mustOwnBadges)
        ? object.mustOwnBadges.map((e: any) => MustOwnBadges.fromJSON(e))
        : [],
      merkleChallenge: isSet(object.merkleChallenge) ? MerkleChallenge.fromJSON(object.merkleChallenge) : undefined,
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

  toJSON(message: IncomingApprovalCriteria): unknown {
    const obj: any = {};
    if (message.mustOwnBadges) {
      obj.mustOwnBadges = message.mustOwnBadges.map((e) => e ? MustOwnBadges.toJSON(e) : undefined);
    } else {
      obj.mustOwnBadges = [];
    }
    message.merkleChallenge !== undefined
      && (obj.merkleChallenge = message.merkleChallenge ? MerkleChallenge.toJSON(message.merkleChallenge) : undefined);
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

  fromPartial<I extends Exact<DeepPartial<IncomingApprovalCriteria>, I>>(object: I): IncomingApprovalCriteria {
    const message = createBaseIncomingApprovalCriteria();
    message.mustOwnBadges = object.mustOwnBadges?.map((e) => MustOwnBadges.fromPartial(e)) || [];
    message.merkleChallenge = (object.merkleChallenge !== undefined && object.merkleChallenge !== null)
      ? MerkleChallenge.fromPartial(object.merkleChallenge)
      : undefined;
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

function createBaseCollectionApproval(): CollectionApproval {
  return {
    fromMappingId: "",
    toMappingId: "",
    initiatedByMappingId: "",
    transferTimes: [],
    badgeIds: [],
    ownershipTimes: [],
    amountTrackerId: "",
    challengeTrackerId: "",
    uri: "",
    customData: "",
    approvalId: "",
    approvalCriteria: undefined,
  };
}

export const CollectionApproval = {
  encode(message: CollectionApproval, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    for (const v of message.ownershipTimes) {
      UintRange.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.amountTrackerId !== "") {
      writer.uint32(58).string(message.amountTrackerId);
    }
    if (message.challengeTrackerId !== "") {
      writer.uint32(66).string(message.challengeTrackerId);
    }
    if (message.uri !== "") {
      writer.uint32(74).string(message.uri);
    }
    if (message.customData !== "") {
      writer.uint32(82).string(message.customData);
    }
    if (message.approvalId !== "") {
      writer.uint32(90).string(message.approvalId);
    }
    if (message.approvalCriteria !== undefined) {
      ApprovalCriteria.encode(message.approvalCriteria, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CollectionApproval {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCollectionApproval();
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
          message.ownershipTimes.push(UintRange.decode(reader, reader.uint32()));
          break;
        case 7:
          message.amountTrackerId = reader.string();
          break;
        case 8:
          message.challengeTrackerId = reader.string();
          break;
        case 9:
          message.uri = reader.string();
          break;
        case 10:
          message.customData = reader.string();
          break;
        case 11:
          message.approvalId = reader.string();
          break;
        case 12:
          message.approvalCriteria = ApprovalCriteria.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CollectionApproval {
    return {
      fromMappingId: isSet(object.fromMappingId) ? String(object.fromMappingId) : "",
      toMappingId: isSet(object.toMappingId) ? String(object.toMappingId) : "",
      initiatedByMappingId: isSet(object.initiatedByMappingId) ? String(object.initiatedByMappingId) : "",
      transferTimes: Array.isArray(object?.transferTimes)
        ? object.transferTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      badgeIds: Array.isArray(object?.badgeIds) ? object.badgeIds.map((e: any) => UintRange.fromJSON(e)) : [],
      ownershipTimes: Array.isArray(object?.ownershipTimes)
        ? object.ownershipTimes.map((e: any) => UintRange.fromJSON(e))
        : [],
      amountTrackerId: isSet(object.amountTrackerId) ? String(object.amountTrackerId) : "",
      challengeTrackerId: isSet(object.challengeTrackerId) ? String(object.challengeTrackerId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      customData: isSet(object.customData) ? String(object.customData) : "",
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      approvalCriteria: isSet(object.approvalCriteria) ? ApprovalCriteria.fromJSON(object.approvalCriteria) : undefined,
    };
  },

  toJSON(message: CollectionApproval): unknown {
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
    if (message.ownershipTimes) {
      obj.ownershipTimes = message.ownershipTimes.map((e) => e ? UintRange.toJSON(e) : undefined);
    } else {
      obj.ownershipTimes = [];
    }
    message.amountTrackerId !== undefined && (obj.amountTrackerId = message.amountTrackerId);
    message.challengeTrackerId !== undefined && (obj.challengeTrackerId = message.challengeTrackerId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.customData !== undefined && (obj.customData = message.customData);
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.approvalCriteria !== undefined && (obj.approvalCriteria = message.approvalCriteria
      ? ApprovalCriteria.toJSON(message.approvalCriteria)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CollectionApproval>, I>>(object: I): CollectionApproval {
    const message = createBaseCollectionApproval();
    message.fromMappingId = object.fromMappingId ?? "";
    message.toMappingId = object.toMappingId ?? "";
    message.initiatedByMappingId = object.initiatedByMappingId ?? "";
    message.transferTimes = object.transferTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.badgeIds = object.badgeIds?.map((e) => UintRange.fromPartial(e)) || [];
    message.ownershipTimes = object.ownershipTimes?.map((e) => UintRange.fromPartial(e)) || [];
    message.amountTrackerId = object.amountTrackerId ?? "";
    message.challengeTrackerId = object.challengeTrackerId ?? "";
    message.uri = object.uri ?? "";
    message.customData = object.customData ?? "";
    message.approvalId = object.approvalId ?? "";
    message.approvalCriteria = (object.approvalCriteria !== undefined && object.approvalCriteria !== null)
      ? ApprovalCriteria.fromPartial(object.approvalCriteria)
      : undefined;
    return message;
  },
};

function createBaseApprovalIdentifierDetails(): ApprovalIdentifierDetails {
  return { approvalId: "", approvalLevel: "", approverAddress: "" };
}

export const ApprovalIdentifierDetails = {
  encode(message: ApprovalIdentifierDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.approvalId !== "") {
      writer.uint32(10).string(message.approvalId);
    }
    if (message.approvalLevel !== "") {
      writer.uint32(18).string(message.approvalLevel);
    }
    if (message.approverAddress !== "") {
      writer.uint32(26).string(message.approverAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApprovalIdentifierDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApprovalIdentifierDetails();
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
          message.approverAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ApprovalIdentifierDetails {
    return {
      approvalId: isSet(object.approvalId) ? String(object.approvalId) : "",
      approvalLevel: isSet(object.approvalLevel) ? String(object.approvalLevel) : "",
      approverAddress: isSet(object.approverAddress) ? String(object.approverAddress) : "",
    };
  },

  toJSON(message: ApprovalIdentifierDetails): unknown {
    const obj: any = {};
    message.approvalId !== undefined && (obj.approvalId = message.approvalId);
    message.approvalLevel !== undefined && (obj.approvalLevel = message.approvalLevel);
    message.approverAddress !== undefined && (obj.approverAddress = message.approverAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApprovalIdentifierDetails>, I>>(object: I): ApprovalIdentifierDetails {
    const message = createBaseApprovalIdentifierDetails();
    message.approvalId = object.approvalId ?? "";
    message.approvalLevel = object.approvalLevel ?? "";
    message.approverAddress = object.approverAddress ?? "";
    return message;
  },
};

function createBaseTransfer(): Transfer {
  return {
    from: "",
    toAddresses: [],
    balances: [],
    precalculateBalancesFromApproval: undefined,
    merkleProofs: [],
    memo: "",
    prioritizedApprovals: [],
    onlyCheckPrioritizedApprovals: false,
  };
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
    if (message.precalculateBalancesFromApproval !== undefined) {
      ApprovalIdentifierDetails.encode(message.precalculateBalancesFromApproval, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.merkleProofs) {
      MerkleProof.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    if (message.memo !== "") {
      writer.uint32(50).string(message.memo);
    }
    for (const v of message.prioritizedApprovals) {
      ApprovalIdentifierDetails.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    if (message.onlyCheckPrioritizedApprovals === true) {
      writer.uint32(64).bool(message.onlyCheckPrioritizedApprovals);
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
          message.precalculateBalancesFromApproval = ApprovalIdentifierDetails.decode(reader, reader.uint32());
          break;
        case 5:
          message.merkleProofs.push(MerkleProof.decode(reader, reader.uint32()));
          break;
        case 6:
          message.memo = reader.string();
          break;
        case 7:
          message.prioritizedApprovals.push(ApprovalIdentifierDetails.decode(reader, reader.uint32()));
          break;
        case 8:
          message.onlyCheckPrioritizedApprovals = reader.bool();
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
      precalculateBalancesFromApproval: isSet(object.precalculateBalancesFromApproval)
        ? ApprovalIdentifierDetails.fromJSON(object.precalculateBalancesFromApproval)
        : undefined,
      merkleProofs: Array.isArray(object?.merkleProofs)
        ? object.merkleProofs.map((e: any) => MerkleProof.fromJSON(e))
        : [],
      memo: isSet(object.memo) ? String(object.memo) : "",
      prioritizedApprovals: Array.isArray(object?.prioritizedApprovals)
        ? object.prioritizedApprovals.map((e: any) => ApprovalIdentifierDetails.fromJSON(e))
        : [],
      onlyCheckPrioritizedApprovals: isSet(object.onlyCheckPrioritizedApprovals)
        ? Boolean(object.onlyCheckPrioritizedApprovals)
        : false,
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
    message.precalculateBalancesFromApproval !== undefined
      && (obj.precalculateBalancesFromApproval = message.precalculateBalancesFromApproval
        ? ApprovalIdentifierDetails.toJSON(message.precalculateBalancesFromApproval)
        : undefined);
    if (message.merkleProofs) {
      obj.merkleProofs = message.merkleProofs.map((e) => e ? MerkleProof.toJSON(e) : undefined);
    } else {
      obj.merkleProofs = [];
    }
    message.memo !== undefined && (obj.memo = message.memo);
    if (message.prioritizedApprovals) {
      obj.prioritizedApprovals = message.prioritizedApprovals.map((e) =>
        e ? ApprovalIdentifierDetails.toJSON(e) : undefined
      );
    } else {
      obj.prioritizedApprovals = [];
    }
    message.onlyCheckPrioritizedApprovals !== undefined
      && (obj.onlyCheckPrioritizedApprovals = message.onlyCheckPrioritizedApprovals);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Transfer>, I>>(object: I): Transfer {
    const message = createBaseTransfer();
    message.from = object.from ?? "";
    message.toAddresses = object.toAddresses?.map((e) => e) || [];
    message.balances = object.balances?.map((e) => Balance.fromPartial(e)) || [];
    message.precalculateBalancesFromApproval =
      (object.precalculateBalancesFromApproval !== undefined && object.precalculateBalancesFromApproval !== null)
        ? ApprovalIdentifierDetails.fromPartial(object.precalculateBalancesFromApproval)
        : undefined;
    message.merkleProofs = object.merkleProofs?.map((e) => MerkleProof.fromPartial(e)) || [];
    message.memo = object.memo ?? "";
    message.prioritizedApprovals = object.prioritizedApprovals?.map((e) => ApprovalIdentifierDetails.fromPartial(e))
      || [];
    message.onlyCheckPrioritizedApprovals = object.onlyCheckPrioritizedApprovals ?? false;
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
