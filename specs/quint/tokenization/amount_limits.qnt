// Quint spec for BitBadges approval amount tracking
// Tests that transfer limits can't be exceeded
module amount_limits {
  // === TYPES ===
  type Address = str

  // === CONSTANTS ===
  pure val ADDRESSES = Set("alice", "bob", "carol")

  // === STATE ===

  // Balances
  var balances: Address -> int

  // Approval with amount limit: (from, to) -> (allowed, maxAmount, usedAmount)
  var approvalLimit: (Address, Address) -> { allowed: bool, maxAmount: int, usedAmount: int }

  // Track successful transfers for verification
  var totalTransferred: int

  // === INITIALIZATION ===
  action init = all {
    balances' = ADDRESSES.mapBy(_ => 100),
    approvalLimit' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => { allowed: false, maxAmount: 0, usedAmount: 0 }),
    totalTransferred' = 0
  }

  // === ACTIONS ===

  // Set approval with limit
  action setApprovalWithLimit(fromAddr: Address, toAddr: Address, allowed: bool, maxAmount: int): bool = all {
    maxAmount >= 0,
    approvalLimit' = approvalLimit.set((fromAddr, toAddr), {
      allowed: allowed,
      maxAmount: maxAmount,
      usedAmount: 0  // Reset on new approval
    }),
    balances' = balances,
    totalTransferred' = totalTransferred
  }

  // Transfer respecting amount limits
  action transfer(fromAddr: Address, toAddr: Address, amount: int): bool =
    val approval = approvalLimit.get((fromAddr, toAddr))
    val remainingAllowance = approval.maxAmount - approval.usedAmount
    all {
      amount > 0,
      fromAddr != toAddr,
      balances.get(fromAddr) >= amount,

      // Must be approved
      approval.allowed,

      // Must not exceed remaining allowance
      amount <= remainingAllowance,

      // Execute transfer
      balances' = balances
        .set(fromAddr, balances.get(fromAddr) - amount)
        .set(toAddr, balances.get(toAddr) + amount),

      // Update used amount
      approvalLimit' = approvalLimit.set((fromAddr, toAddr), {
        allowed: approval.allowed,
        maxAmount: approval.maxAmount,
        usedAmount: approval.usedAmount + amount
      }),

      totalTransferred' = totalTransferred + amount
    }

  // Reset approval usage (e.g., for time-based reset)
  action resetApprovalUsage(fromAddr: Address, toAddr: Address): bool =
    val approval = approvalLimit.get((fromAddr, toAddr))
    all {
      approval.allowed,  // Can only reset active approvals
      approvalLimit' = approvalLimit.set((fromAddr, toAddr), {
        allowed: approval.allowed,
        maxAmount: approval.maxAmount,
        usedAmount: 0
      }),
      balances' = balances,
      totalTransferred' = totalTransferred
    }

  // === STEP ===
  action step = any {
    // Set approval with random limit
    nondet f = ADDRESSES.oneOf()
    nondet t = ADDRESSES.oneOf()
    nondet allowed = Set(true, false).oneOf()
    nondet limit = 0.to(200).oneOf()
    setApprovalWithLimit(f, t, allowed, limit),

    // Transfer
    nondet txFrom = ADDRESSES.oneOf()
    nondet txTo = ADDRESSES.oneOf()
    nondet txAmt = 1.to(50).oneOf()
    transfer(txFrom, txTo, txAmt),

    // Reset usage
    nondet rFrom = ADDRESSES.oneOf()
    nondet rTo = ADDRESSES.oneOf()
    resetApprovalUsage(rFrom, rTo)
  }

  // === INVARIANTS ===

  // No negative balances
  val inv_no_negative_balances: bool =
    ADDRESSES.forall(addr => balances.get(addr) >= 0)

  // Balance conservation
  val inv_balance_conservation: bool =
    ADDRESSES.fold(0, (sum, addr) => sum + balances.get(addr)) == ADDRESSES.size() * 100

  // CRITICAL: Used amount never exceeds max amount
  val inv_used_never_exceeds_max: bool =
    tuples(ADDRESSES, ADDRESSES).forall(pair =>
      val approval = approvalLimit.get(pair)
      approval.usedAmount <= approval.maxAmount
    )

  // CRITICAL: Used amount is always non-negative
  val inv_used_non_negative: bool =
    tuples(ADDRESSES, ADDRESSES).forall(pair =>
      approvalLimit.get(pair).usedAmount >= 0
    )

  // Combined invariant
  val inv_all: bool =
    inv_no_negative_balances and
    inv_balance_conservation and
    inv_used_never_exceeds_max and
    inv_used_non_negative
}
