// Focused Quint spec for BitBadges approval hierarchy
// Tests the three-layer approval system and override semantics
module approval_hierarchy {
  // === TYPES ===
  type Address = str

  // === CONSTANTS ===
  pure val ADDRESSES = Set("alice", "bob", "carol")

  // === STATE ===

  // Balances: address -> amount
  var balances: Address -> int

  // Collection-level approval: (from, to) -> allowed
  var collectionApproval: (Address, Address) -> bool

  // Collection-level override flags
  var collectionOverridesOutgoing: (Address, Address) -> bool
  var collectionOverridesIncoming: (Address, Address) -> bool

  // User outgoing approvals: (owner, to) -> allowed
  var outgoingApproval: (Address, Address) -> bool

  // User incoming approvals: (recipient, from) -> allowed
  var incomingApproval: (Address, Address) -> bool

  // Track which approvals were used (for debugging)
  var lastUsedCollectionApproval: bool
  var lastUsedOutgoingApproval: bool
  var lastUsedIncomingApproval: bool
  var lastOverrodeOutgoing: bool
  var lastOverrodeIncoming: bool

  // === INITIALIZATION ===
  action init = all {
    balances' = ADDRESSES.mapBy(_ => 100),
    collectionApproval' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => false),
    collectionOverridesOutgoing' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => false),
    collectionOverridesIncoming' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => false),
    outgoingApproval' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => false),
    incomingApproval' = tuples(ADDRESSES, ADDRESSES).mapBy(_ => false),
    lastUsedCollectionApproval' = false,
    lastUsedOutgoingApproval' = false,
    lastUsedIncomingApproval' = false,
    lastOverrodeOutgoing' = false,
    lastOverrodeIncoming' = false
  }

  // === ACTIONS ===

  // Set collection approval
  action setCollectionApproval(fromAddr: Address, toAddr: Address, allowed: bool, overrideOut: bool, overrideIn: bool): bool = all {
    collectionApproval' = collectionApproval.set((fromAddr, toAddr), allowed),
    collectionOverridesOutgoing' = collectionOverridesOutgoing.set((fromAddr, toAddr), overrideOut),
    collectionOverridesIncoming' = collectionOverridesIncoming.set((fromAddr, toAddr), overrideIn),
    // Unchanged
    balances' = balances,
    outgoingApproval' = outgoingApproval,
    incomingApproval' = incomingApproval,
    lastUsedCollectionApproval' = lastUsedCollectionApproval,
    lastUsedOutgoingApproval' = lastUsedOutgoingApproval,
    lastUsedIncomingApproval' = lastUsedIncomingApproval,
    lastOverrodeOutgoing' = lastOverrodeOutgoing,
    lastOverrodeIncoming' = lastOverrodeIncoming
  }

  // Set user outgoing approval
  action setOutgoingApproval(owner: Address, toAddr: Address, allowed: bool): bool = all {
    outgoingApproval' = outgoingApproval.set((owner, toAddr), allowed),
    // Unchanged
    balances' = balances,
    collectionApproval' = collectionApproval,
    collectionOverridesOutgoing' = collectionOverridesOutgoing,
    collectionOverridesIncoming' = collectionOverridesIncoming,
    incomingApproval' = incomingApproval,
    lastUsedCollectionApproval' = lastUsedCollectionApproval,
    lastUsedOutgoingApproval' = lastUsedOutgoingApproval,
    lastUsedIncomingApproval' = lastUsedIncomingApproval,
    lastOverrodeOutgoing' = lastOverrodeOutgoing,
    lastOverrodeIncoming' = lastOverrodeIncoming
  }

  // Set user incoming approval
  action setIncomingApproval(recipient: Address, fromAddr: Address, allowed: bool): bool = all {
    incomingApproval' = incomingApproval.set((recipient, fromAddr), allowed),
    // Unchanged
    balances' = balances,
    collectionApproval' = collectionApproval,
    collectionOverridesOutgoing' = collectionOverridesOutgoing,
    collectionOverridesIncoming' = collectionOverridesIncoming,
    outgoingApproval' = outgoingApproval,
    lastUsedCollectionApproval' = lastUsedCollectionApproval,
    lastUsedOutgoingApproval' = lastUsedOutgoingApproval,
    lastUsedIncomingApproval' = lastUsedIncomingApproval,
    lastOverrodeOutgoing' = lastOverrodeOutgoing,
    lastOverrodeIncoming' = lastOverrodeIncoming
  }

  // CORE TRANSFER LOGIC - Three-layer approval check with overrides
  action transfer(fromAddr: Address, toAddr: Address, amount: int): bool =
    // Pre-compute all approval checks
    val hasCollApproval = collectionApproval.get((fromAddr, toAddr))
    val overridesOut = collectionOverridesOutgoing.get((fromAddr, toAddr))
    val overridesIn = collectionOverridesIncoming.get((fromAddr, toAddr))
    val hasOutApproval = outgoingApproval.get((fromAddr, toAddr))
    val hasInApproval = incomingApproval.get((toAddr, fromAddr))

    // Outgoing check: needed unless collection approval overrides it
    val outgoingPasses = overridesOut or hasOutApproval

    // Incoming check: needed unless collection approval overrides it
    val incomingPasses = overridesIn or hasInApproval

    all {
      // Basic validation
      amount > 0,
      fromAddr != toAddr,
      balances.get(fromAddr) >= amount,

      // Layer 1: Must have collection approval
      hasCollApproval,

      // Layer 2: Outgoing must pass (or be overridden)
      outgoingPasses,

      // Layer 3: Incoming must pass (or be overridden)
      incomingPasses,

      // Execute transfer
      balances' = balances
        .set(fromAddr, balances.get(fromAddr) - amount)
        .set(toAddr, balances.get(toAddr) + amount),

      // Track what was used for debugging
      lastUsedCollectionApproval' = hasCollApproval,
      lastUsedOutgoingApproval' = hasOutApproval,
      lastUsedIncomingApproval' = hasInApproval,
      lastOverrodeOutgoing' = overridesOut,
      lastOverrodeIncoming' = overridesIn,

      // Unchanged
      collectionApproval' = collectionApproval,
      collectionOverridesOutgoing' = collectionOverridesOutgoing,
      collectionOverridesIncoming' = collectionOverridesIncoming,
      outgoingApproval' = outgoingApproval,
      incomingApproval' = incomingApproval
    }

  // === STEP ===
  action step = any {
    // Set collection approval with random override flags
    nondet f1 = ADDRESSES.oneOf()
    nondet t1 = ADDRESSES.oneOf()
    nondet allowed1 = Set(true, false).oneOf()
    nondet overOut = Set(true, false).oneOf()
    nondet overIn = Set(true, false).oneOf()
    setCollectionApproval(f1, t1, allowed1, overOut, overIn),

    // Set outgoing approval
    nondet owner = ADDRESSES.oneOf()
    nondet t2 = ADDRESSES.oneOf()
    nondet allowed2 = Set(true, false).oneOf()
    setOutgoingApproval(owner, t2, allowed2),

    // Set incoming approval
    nondet recip = ADDRESSES.oneOf()
    nondet f3 = ADDRESSES.oneOf()
    nondet allowed3 = Set(true, false).oneOf()
    setIncomingApproval(recip, f3, allowed3),

    // Transfer
    nondet txFrom = ADDRESSES.oneOf()
    nondet txTo = ADDRESSES.oneOf()
    nondet txAmt = 1.to(50).oneOf()
    transfer(txFrom, txTo, txAmt)
  }

  // === INVARIANTS ===

  // No negative balances
  val inv_no_negative_balances: bool =
    ADDRESSES.forall(addr => balances.get(addr) >= 0)

  // Conservation: total supply constant
  val inv_balance_conservation: bool =
    ADDRESSES.fold(0, (sum, addr) => sum + balances.get(addr)) == ADDRESSES.size() * 100

  // CRITICAL: If transfer happened, collection approval MUST have been true
  // (This validates that transfers can't bypass collection-level check)
  val inv_collection_approval_required: bool =
    // If a transfer changed balances, collection approval was used
    lastUsedCollectionApproval or ADDRESSES.forall(addr => balances.get(addr) == 100)

  // CRITICAL: Override semantics are respected
  // If outgoing was overridden, outgoing approval wasn't required
  // If not overridden, outgoing approval must have been present
  val inv_override_outgoing_correct: bool =
    // After any transfer, either:
    // 1. No transfer happened (all balances = 100), OR
    // 2. Outgoing was overridden, OR
    // 3. Outgoing approval was explicitly granted
    ADDRESSES.forall(addr => balances.get(addr) == 100) or
    lastOverrodeOutgoing or
    lastUsedOutgoingApproval

  val inv_override_incoming_correct: bool =
    ADDRESSES.forall(addr => balances.get(addr) == 100) or
    lastOverrodeIncoming or
    lastUsedIncomingApproval

  // Combined invariant
  val inv_all: bool =
    inv_no_negative_balances and
    inv_balance_conservation and
    inv_collection_approval_required and
    inv_override_outgoing_correct and
    inv_override_incoming_correct
}
