// Quint spec for BitBadges approval replay protection
// Tests that incrementing approval versions prevents replay attacks
module replay_protection {
  // === TYPES ===
  type Address = str
  type ApprovalId = str

  // === CONSTANTS ===
  pure val ADDRESSES = Set("alice", "bob", "carol")
  pure val APPROVAL_IDS = Set("approval1", "approval2")

  // === STATE ===

  // Balances
  var balances: Address -> int

  // Approval definitions: (approvalId, from, to) -> version
  var approvalVersion: (ApprovalId, Address, Address) -> int

  // User's consumed versions: (user, approvalId, from, to) -> lastConsumedVersion
  var consumedVersion: (Address, ApprovalId, Address, Address) -> int

  // Track replay attempts
  var replayAttempts: int
  var successfulTransfers: int

  // === INITIALIZATION ===
  action init = all {
    balances' = ADDRESSES.mapBy(_ => 100),
    approvalVersion' = tuples(APPROVAL_IDS, ADDRESSES, ADDRESSES).mapBy(_ => 1),
    consumedVersion' = tuples(ADDRESSES, APPROVAL_IDS, ADDRESSES, ADDRESSES).mapBy(_ => 0),
    replayAttempts' = 0,
    successfulTransfers' = 0
  }

  // === ACTIONS ===

  // Increment approval version (revoke and re-grant)
  action incrementVersion(approvalId: ApprovalId, fromAddr: Address, toAddr: Address): bool =
    val currentVersion = approvalVersion.get((approvalId, fromAddr, toAddr))
    all {
      currentVersion < 100,  // Prevent overflow
      approvalVersion' = approvalVersion.set((approvalId, fromAddr, toAddr), currentVersion + 1),
      balances' = balances,
      consumedVersion' = consumedVersion,
      replayAttempts' = replayAttempts,
      successfulTransfers' = successfulTransfers
    }

  // Transfer using approval (with version check)
  action transfer(initiator: Address, fromAddr: Address, toAddr: Address, amount: int, approvalId: ApprovalId): bool =
    val currentApprovalVersion = approvalVersion.get((approvalId, fromAddr, toAddr))
    val userConsumedVersion = consumedVersion.get((initiator, approvalId, fromAddr, toAddr))

    // Version check: user must not have already used this version
    val versionValid = userConsumedVersion < currentApprovalVersion

    all {
      amount > 0,
      fromAddr != toAddr,
      balances.get(fromAddr) >= amount,

      // Version must be valid (not a replay)
      versionValid,

      // Execute transfer
      balances' = balances
        .set(fromAddr, balances.get(fromAddr) - amount)
        .set(toAddr, balances.get(toAddr) + amount),

      // Mark this version as consumed for this user
      consumedVersion' = consumedVersion.set(
        (initiator, approvalId, fromAddr, toAddr),
        currentApprovalVersion
      ),

      approvalVersion' = approvalVersion,
      replayAttempts' = replayAttempts,
      successfulTransfers' = successfulTransfers + 1
    }

  // Attempt replay (should always fail, tracked for verification)
  action attemptReplay(initiator: Address, fromAddr: Address, toAddr: Address, amount: int, approvalId: ApprovalId): bool =
    val currentApprovalVersion = approvalVersion.get((approvalId, fromAddr, toAddr))
    val userConsumedVersion = consumedVersion.get((initiator, approvalId, fromAddr, toAddr))

    // This is a replay attempt when version is already consumed
    val isReplayAttempt = userConsumedVersion >= currentApprovalVersion

    all {
      amount > 0,
      fromAddr != toAddr,
      isReplayAttempt,  // This IS a replay attempt

      // Replay should NOT change balances
      balances' = balances,
      consumedVersion' = consumedVersion,
      approvalVersion' = approvalVersion,
      replayAttempts' = replayAttempts + 1,
      successfulTransfers' = successfulTransfers
    }

  // === STEP ===
  action step = any {
    // Increment version
    nondet aid = APPROVAL_IDS.oneOf()
    nondet f1 = ADDRESSES.oneOf()
    nondet t1 = ADDRESSES.oneOf()
    incrementVersion(aid, f1, t1),

    // Valid transfer
    nondet initiator = ADDRESSES.oneOf()
    nondet f2 = ADDRESSES.oneOf()
    nondet t2 = ADDRESSES.oneOf()
    nondet amt = 1.to(50).oneOf()
    nondet aid2 = APPROVAL_IDS.oneOf()
    transfer(initiator, f2, t2, amt, aid2),

    // Attempt replay
    nondet attacker = ADDRESSES.oneOf()
    nondet f3 = ADDRESSES.oneOf()
    nondet t3 = ADDRESSES.oneOf()
    nondet amt2 = 1.to(50).oneOf()
    nondet aid3 = APPROVAL_IDS.oneOf()
    attemptReplay(attacker, f3, t3, amt2, aid3)
  }

  // === INVARIANTS ===

  // No negative balances
  val inv_no_negative_balances: bool =
    ADDRESSES.forall(addr => balances.get(addr) >= 0)

  // Balance conservation
  val inv_balance_conservation: bool =
    ADDRESSES.fold(0, (sum, addr) => sum + balances.get(addr)) == ADDRESSES.size() * 100

  // CRITICAL: Consumed version never exceeds approval version
  // This ensures users can't "time travel" to future versions
  val inv_consumed_never_exceeds_approval: bool =
    tuples(ADDRESSES, APPROVAL_IDS, ADDRESSES, ADDRESSES).forall(key =>
      val consumed = consumedVersion.get(key)
      val approval = approvalVersion.get((key._2, key._3, key._4))
      consumed <= approval
    )

  // CRITICAL: Replay attempts don't change balances
  // (Implicitly tested by balance conservation + no transfers on replay)
  val inv_replay_doesnt_transfer: bool =
    // If only replays happened (no successful transfers), balances unchanged
    successfulTransfers == 0 implies
      ADDRESSES.forall(addr => balances.get(addr) == 100)

  // Combined invariant
  val inv_all: bool =
    inv_no_negative_balances and
    inv_balance_conservation and
    inv_consumed_never_exceeds_approval
}
