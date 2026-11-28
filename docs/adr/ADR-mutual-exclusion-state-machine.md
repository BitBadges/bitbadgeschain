# ADR: Native Approval Mutual-Exclusion Groups & State-Machine Transitions

- **Ticket:** BB-376
- **Status:** Scoping ADR (Implementation deferred to dedicated chain upgrade PR)
- **Scope:** `bitbadgeschain` (Go), parity analysis for `bitbadges-solana` (Rust Anchor), `crates/bb-standards`, and `BitBadgesSpec` (`BB-360`)

## 1. Context & Problem Statement

Multi-outcome standards on BitBadges — including **Bounties** (fulfill vs. refund), **Crowdfunds** (release goal vs. refund backers), **Refundable Invoices / PaymentRequestV2** (pay vs. cancel vs. refund), **Auctions** (accept bid vs. cancel), and **Prediction Markets** (resolve YES vs. NO vs. invalid refund) — all require **provably race-free mutual exclusion**: once one outcome branch executes, competing branches must be atomically and permanently disabled.

Today, builders simulate mutual exclusion by sharing a single `maxNumTransfers.amountTrackerId` with `overallMaxNumTransfers: 1` across two approvals. This workaround has severe limitations:
1. **Fragile convention:** Nothing stops an accidental mismatch in `amountTrackerId`, or a cap > 1 on one branch, or a missing `autoDeletionOptions`.
2. **No multi-phase lifecycles:** A shared counter cannot model multi-stage transitions such as `INITIAL -> PAID -> REFUNDED` vs `INITIAL -> CANCELLED` (e.g., a refundable invoice or crowdfund with distinct contribution and settlement phases).
3. **Difficult to audit:** Both humans and deterministic verifiers (`BB-374`) must reverse-engineer cross-approval tracker ID graphs rather than inspecting an explicit state transition rule.

## 2. Options Analysis

### Option A: `exclusivityGroupId` (Atomic Latch / Shared Counter)

Add an optional field on `ApprovalCriteria`:
```protobuf
message ExclusivityGroup {
  string groupId = 1;       // Scoped to (collectionId, approvalLevel, approverAddress)
  uint64 maxExecutions = 2; // Default 1 (strict one-of-N latch)
}
```
- **Semantics:** All approvals referencing `groupId` share a single atomic execution counter stored at `ExclusivityPrefix | collectionId | level | approver | groupId`. When any approval in the group executes, the counter increments atomically inside `HandleTransfer`'s `AtomicContext`; once `count >= maxExecutions`, every approval in the group is rejected.
- **Pros:** Minimal schema (2 fields), tiny storage footprint (~8 bytes per group), trivial to understand for 1-of-N races (Bounty fulfill vs. refund).
- **Cons:** Cannot express multi-phase lifecycles (e.g., `INITIAL -> PAID -> REFUNDED`, where `REFUNDED` is ONLY reachable *after* `PAID`, never directly from `INITIAL`, and `CANCELLED` is ONLY reachable from `INITIAL`).

### Option B: Lightweight `StateTransition` (Scoped Per-Key State Machine)

Add an optional field on `ApprovalCriteria`:
```protobuf
message StateTransition {
  string stateKey = 1;                // e.g. "bounty:1", "invoice:42", "crowdfund"
  repeated string requiredStates = 2; // Canonical initial state is "INITIAL" (unset/empty "" normalizes to "INITIAL")
  string nextState = 3;               // Target state written atomically on commit
  bool terminal = 4;                  // If true, no subsequent transition can leave nextState
}
```
- **Semantics:** Each `(collectionId, stateKey)` stores a lightweight state string (or compact enum byte in Anchor/Solana), whose canonical default initial state when unset is `"INITIAL"` (empty string `""` normalizes to `"INITIAL"`). When an approval executes:
  1. Assert `currentState` is NOT marked terminal.
  2. Assert `currentState` matches one of `requiredStates` (treating an empty `requiredStates` slice as `["INITIAL"]`).
  3. Atomically set `currentState = nextState` (and `terminal = rule.terminal`) inside `AtomicContext`. If any later transfer step fails, `AtomicContext.Rollback()` reverts the state!
- **Pros:**
  - Strictly subsumes Option A: a 1-of-N `exclusivityGroupId` is simply `{ stateKey: groupId, requiredStates: ["INITIAL"], nextState: "EXECUTED", terminal: true }`.
  - Cleanly expresses multi-phase standards (`INITIAL -> PAID -> REFUNDED` vs `INITIAL -> CANCELLED`, Crowdfund `INITIAL -> RELEASED` vs `INITIAL -> REFUNDING`).
  - Trivial to audit & verify deterministically in `BB-374` (`crates/bb-standards` and the dedicated chain verifier) by inspecting the finite state graph directly.
- **Cons:** Slightly richer struct (4 fields vs 2), requires bounding `stateKey` and state name lengths (e.g. `<= 64` bytes ASCII alphanumeric/`_-:`).

### Comparison Matrix

| Criterion | Option A (`exclusivityGroupId`) | Option B (`StateTransition`) | Recommendation |
| --- | --- | --- | --- |
| 1-of-N Mutual Exclusion (Bounty fulfill vs refund) | Native | Native (special case) | Tie |
| Multi-phase lifecycles (Paid -> Refunded, Crowdfund) | Not supported | Native | **Option B** |
| Storage per key | ~16 bytes | ~32–64 bytes | Negligible diff |
| Gas overhead | 1 KV read + 1 KV write | 1 KV read + 1 KV write | Tie |
| Deterministic Standard Verification (`BB-374`) | Checks shared `groupId` | Checks reachability & terminal states | **Option B** |
| Solana Anchor parity | PDA `ExclusivityLatch { count: u64 }` | PDA `StateSlot { state: u8, terminal: bool }` | **Option B** (compact u8 enum) |

### Decision: Recommend Option B (`StateTransition`) + Option A Sugar

Adopt **Option B (`StateTransition`)** as the single on-chain state-machine primitive with `"INITIAL"` as the single canonical default initial state across all standards, and provide **Option A (`exclusivityGroupId`)** as syntactic sugar in `BitBadgesSpec` (`BB-360`) and SDK builders (compiling `exclusivityGroupId: "X"` to `{ stateKey: "X", requiredStates: ["INITIAL"], nextState: "CLAIMED", terminal: true }`).

## 3. Concrete Before / After JSON Comparisons

### 3.1 Bounty (Fulfill vs. Refund)

**Before (brittle shared `amountTrackerId` convention):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "bounty-fulfill",
      "fromListId": "Mint",
      "toListId": "AllWithoutMint",
      "initiatedByListId": "bb1issuer...",
      "approvalCriteria": {
        "overridesFromOutgoingApprovals": true,
        "coinTransfers": [{ "to": "bb1worker...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "500000000" }] }],
        "maxNumTransfers": { "overallMaxNumTransfers": "1", "amountTrackerId": "bounty-shared-latch-DO-NOT-RENAME" },
        "autoDeletionOptions": { "afterOverallMaxNumTransfers": true }
      }
    },
    {
      "approvalId": "bounty-refund",
      "fromListId": "Mint",
      "toListId": "bb1issuer...",
      "initiatedByListId": "bb1issuer...",
      "transferTimes": [{ "start": "1735689600000", "end": "18446744073709551615" }],
      "approvalCriteria": {
        "overridesFromOutgoingApprovals": true,
        "coinTransfers": [{ "to": "bb1issuer...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "500000000" }] }],
        "maxNumTransfers": { "overallMaxNumTransfers": "1", "amountTrackerId": "bounty-shared-latch-DO-NOT-RENAME" },
        "autoDeletionOptions": { "afterOverallMaxNumTransfers": true }
      }
    }
  ]
}
```

**After (native `stateTransition` / `exclusivityGroupId`):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "bounty-fulfill",
      "fromListId": "Mint",
      "toListId": "AllWithoutMint",
      "initiatedByListId": "bb1issuer...",
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1worker...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "500000000" }] }],
        "stateTransition": { "stateKey": "bounty", "requiredStates": ["INITIAL"], "nextState": "FULFILLED", "terminal": true }
      }
    },
    {
      "approvalId": "bounty-refund",
      "fromListId": "Mint",
      "toListId": "bb1issuer...",
      "initiatedByListId": "bb1issuer...",
      "transferTimes": [{ "start": "1735689600000", "end": "18446744073709551615" }],
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1issuer...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "500000000" }] }],
        "stateTransition": { "stateKey": "bounty", "requiredStates": ["INITIAL"], "nextState": "REFUNDED", "terminal": true }
      }
    }
  ]
}
```

### 3.2 Crowdfund (Contribute -> Release OR Refund)

**Before (cannot cleanly gate multi-backer refunds vs creator withdrawal without race windows):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "contribute",
      "transferTimes": [{ "start": "1700000000000", "end": "1702592000000" }],
      "approvalCriteria": {
        "coinTransfers": [{ "to": "Mint", "coins": [{ "denom": "uusdc", "amount": "10000000" }] }]
      }
    },
    {
      "approvalId": "withdraw-funds",
      "transferTimes": [{ "start": "1702592000001", "end": "18446744073709551615" }],
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1creator...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "100000000" }] }],
        "maxNumTransfers": { "overallMaxNumTransfers": "1", "amountTrackerId": "crowdfund-settlement-fragile" }
      }
    }
  ]
}
```
*(Notice: individual backers cannot each trigger a refund while sharing a single `overallMaxNumTransfers: 1` tracker with `withdraw-funds`!)*

**After (state-machine separates initial funding `INITIAL`, terminal `RELEASED`, and multi-backer `REFUNDING`):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "contribute",
      "transferTimes": [{ "start": "1700000000000", "end": "1702592000000" }],
      "approvalCriteria": {
        "coinTransfers": [{ "to": "Mint", "coins": [{ "denom": "uusdc", "amount": "10000000" }] }],
        "stateTransition": { "stateKey": "campaign", "requiredStates": ["INITIAL"], "nextState": "INITIAL" }
      }
    },
    {
      "approvalId": "release-to-creator",
      "transferTimes": [{ "start": "1702592000001", "end": "18446744073709551615" }],
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1creator...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "100000000" }] }],
        "stateTransition": { "stateKey": "campaign", "requiredStates": ["INITIAL"], "nextState": "RELEASED", "terminal": true }
      }
    },
    {
      "approvalId": "open-refunds",
      "transferTimes": [{ "start": "1702592000001", "end": "18446744073709551615" }],
      "approvalCriteria": {
        "stateTransition": { "stateKey": "campaign", "requiredStates": ["INITIAL"], "nextState": "REFUNDING" }
      }
    },
    {
      "approvalId": "claim-backer-refund",
      "approvalCriteria": {
        "coinTransfers": [{ "overrideToWithInitiator": true, "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "10000000" }] }],
        "stateTransition": { "stateKey": "campaign", "requiredStates": ["REFUNDING"], "nextState": "REFUNDING" }
      }
    }
  ]
}
```
*(Once `open-refunds` flips `INITIAL -> REFUNDING`, `release-to-creator` is permanently impossible while N backers can each claim their individual refund!)*

### 3.3 Refundable Invoice (`PaymentRequestV2`)

**Before (pay only; cancellation and post-payment refund require ad-hoc out-of-band workarounds):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "pay-invoice",
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1merchant...", "coins": [{ "denom": "uusdc", "amount": "25000000" }] }],
        "maxNumTransfers": { "overallMaxNumTransfers": "1", "amountTrackerId": "invoice-1" },
        "autoDeletionOptions": { "afterOneUse": true }
      }
    }
  ]
}
```

**After (provably race-free `INITIAL -> PAID | CANCELLED` and `PAID -> REFUNDED`):**
```json
{
  "collectionApprovals": [
    {
      "approvalId": "pay-invoice",
      "approvalCriteria": {
        "coinTransfers": [{ "to": "Mint", "coins": [{ "denom": "uusdc", "amount": "25000000" }] }],
        "stateTransition": { "stateKey": "inv:1", "requiredStates": ["INITIAL"], "nextState": "PAID" }
      }
    },
    {
      "approvalId": "cancel-invoice",
      "initiatedByListId": "bb1merchant...",
      "approvalCriteria": {
        "stateTransition": { "stateKey": "inv:1", "requiredStates": ["INITIAL"], "nextState": "CANCELLED", "terminal": true }
      }
    },
    {
      "approvalId": "refund-invoice",
      "initiatedByListId": "bb1merchant...",
      "approvalCriteria": {
        "coinTransfers": [{ "to": "bb1payer...", "overrideFromWithApproverAddress": true, "coins": [{ "denom": "uusdc", "amount": "25000000" }] }],
        "stateTransition": { "stateKey": "inv:1", "requiredStates": ["PAID"], "nextState": "REFUNDED", "terminal": true }
      }
    }
  ]
}
```

## 4. Impact on Additional Standards & Cross-Chain Parity

- **Auctions:** `INITIAL -> SETTLED | CANCELLED` prevents accepting a bid after auction cancellation.
- **Prediction Markets:** `INITIAL -> RESOLVED_YES | RESOLVED_NO | VOID` prevents dual-outcome settlement while allowing unlimited redemptions in the winning state.
- **Cross-Chain Parity:**
  - **Cosmos (`bitbadgeschain`):** Key `StateMachinePrefix | collectionId | stateKey` stored in `x/tokenization` KVStore; mutated inside `evmcompat.AtomicContext` for atomic rollback with transfers & bank sends.
  - **Solana (`bitbadges-solana`):** PDA seeded by `["state", collection_id, state_key]` holding `{ state: u8, terminal: bool }` checked via Anchor constraint.
  - **Indexer (`crates/bb-standards`) & `BitBadgesSpec` (`BB-360`):** Verifies transition graph reachability and compiles high-level `exclusivityGroup` declarations directly to `StateTransition`.

## 5. Planned Implementation & Verification Criteria

Implementation is deferred to a dedicated chain upgrade PR introducing the `StateTransition` proto field, KVStore persistence, and keeper execution hook inside `AtomicContext`. Acceptance tests must cover:
- Atomic 1-of-N mutual exclusion (`fulfill` vs. `refund`) with zero double-execution.
- Multi-phase transitions (Bounty, Crowdfund with multi-backer refunds, Refundable Invoice), terminal-state immutability, and full rollback on transfer failure.
