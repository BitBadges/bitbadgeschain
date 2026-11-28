# ADR: Core Approval Ergonomics (2.1–2.4)

- **Ticket:** BB-375
- **Status:** Scoping ADR + Phase 1 Cast/Nil-Safety Hardening
- **Scope:** `x/tokenization/types` (`approved_transfers_casts.go`, `approved_transfers_casts_test.go`, `validate_basic.go`), `proto/tokenization/approval_criteria.proto` (future additive evolution), **`proto/tokenization/permissions.proto` EXPLICITLY UNCHANGED**

## 1. Context & Motivation

`x/tokenization` currently models transfer approval criteria across three nearly identical protobuf messages (`ApprovalCriteria`, `OutgoingApprovalCriteria`, `IncomingApprovalCriteria`), allows contradictory/mutually exclusive criteria fields to be set simultaneously, relies on raw string-DSL address list predicates, and requires developers and AI agents to remember non-obvious boolean traps (such as `overridesFromOutgoingApprovals: true` when `fromListId == "Mint"`).

These ergonomics gaps have directly caused:
- **Security regressions (`C-1` class):** Whenever a new criterion (such as `EvmQueryChallenges`) is added to the 3 structs, developers must manually update 4 separate cast functions in `x/tokenization/types/approved_transfers_casts.go`. Omitting a single line silently drops security checks on user-level approvals.
- **Delimiter injection (`H-14` class):** String-DSL list IDs (`!Mint:bb1a:bb1b`, `!(...)`, `AllWithout...`) are parsed with ad-hoc string slicing, inviting malformed inputs and divergence between chain, indexer, and SDK parsers.
- **Combinatorial footguns & LLM hallucinations:** Duplicating 17+ fields across 3 structs inflates JSON/MCP schemas ~3x and lets callers specify illegal combinations (e.g., `PredeterminedBalances` + `ApprovalAmounts`, `DurationFromTimestamp` + `RecurringOwnershipTimes`, or `RequireToEqualsInitiatedBy` + `RequireToDoesNotEqualInitiatedBy`).

## 2. Design Decisions

### 2.1 Unifying the 3 `ApprovalCriteria` Structs

#### Current State
- `ApprovalCriteria` (25 fields): 14 universal fields + 3 incoming-side fields (`RequireFromEqualsInitiatedBy`, `RequireFromDoesNotEqualInitiatedBy`, `SenderChecks`) + 3 outgoing-side fields (`RequireToEqualsInitiatedBy`, `RequireToDoesNotEqualInitiatedBy`, `RecipientChecks`) + 5 collection-only fields (`OverridesFromOutgoingApprovals`, `OverridesToIncomingApprovals`, `AllowBackedMinting`, `AllowSpecialWrapping`, `UserApprovalSettings`).
- `OutgoingApprovalCriteria` (17 fields): duplicates the 14 universal + 3 outgoing-side fields with different protobuf field numbers (e.g., `evmQueryChallenges = 17` vs `25`).
- `IncomingApprovalCriteria` (17 fields): duplicates the 14 universal + 3 incoming-side fields.

#### Target Architecture & Phased Migration
- **Phase 1 (Immediate Cast & Nil-Safety Hardening — Implemented Now):**
  - Keep existing wire protobuf messages intact so zero `.proto` schema changes or store migrations are introduced in this PR.
  - Add nil-safety guards to all 4 single-element transfer cast functions in `x/tokenization/types/approved_transfers_casts.go` (without silently dropping `nil` elements in slice casts) and reject `nil` approval entries explicitly in `ValidateUserOutgoingApprovals`, `ValidateUserIncomingApprovals`, and `ValidateCollectionApprovals` (`x/tokenization/types/validate_basic.go`).
  - Enforce a **Go reflection round-trip invariant test** (`TestApprovalCriteriaCastReflectionRoundTrip` in `x/tokenization/types/approved_transfers_casts_test.go`) that enumerates EVERY exported field of `ApprovalCriteria`, `OutgoingApprovalCriteria`, and `IncomingApprovalCriteria` via `reflect.TypeOf`, verifies the exact partition of shared vs tier-exclusive fields, populates every field with a non-zero sentinel value, and asserts full round-trip fidelity across all 4 casts (including `EvmQueryChallenges`). Any future proto field added without updating all 4 casts fails `go test` immediately.
- **Phase 2 (Next Major Wire Upgrade — Unified Protobuf Criteria):**
  - Introduce a single canonical `ApprovalCriteria` message across `CollectionApproval`, `UserOutgoingApproval`, and `UserIncomingApproval` (or a `v2` message family sharing tag numbers `1..26`) paired with a single `ValidateCriteriaForTier(tier ApprovalTier, c *ApprovalCriteria)` validator:
    - `TierOutgoing`: rejects `SenderChecks`, `RequireFrom*`, `Overrides*`, `AllowBackedMinting`, `AllowSpecialWrapping`, `UserApprovalSettings`.
    - `TierIncoming`: rejects `RecipientChecks`, `RequireTo*`, `Overrides*`, `AllowBackedMinting`, `AllowSpecialWrapping`, `UserApprovalSettings`.
    - `TierCollection`: permits all fields.
  - Eliminates manual 4-way cast functions entirely and shrinks LLM/MCP tool schema definitions by ~3x.

### 2.2 Typed Mutual-Exclusivity (`oneof`) & Structured Address-List Predicates

#### Proposed Mutual-Exclusivity & Satisfiability Rules (Future Consensus Upgrade / Builder Validation)
1. **Quantity / Balance Mode (`PredeterminedBalances` + `ApprovalAmounts` Partial-Payment Caveat):**
   - `PredeterminedBalances` dictates exact per-transfer balances (incremented, scaled, or manual), whereas `ApprovalAmounts` tracks cumulative transfer amounts up to a quota.
   - **Important compatibility constraint (`PaymentRequestV2` partial payments):** `PaymentRequestV2` partial-payment approvals **legitimately combine** `PredeterminedBalances` (e.g., scaled/incremented payment receipt balances) with `ApprovalAmounts` (capping cumulative paid amount across multiple partial-payment transfers). Therefore, `PredeterminedBalances` and `ApprovalAmounts` MUST **NOT** be made strictly mutually exclusive (nor collapsed into a naive `oneof`) without first providing a dedicated partial-payment/total-cap replacement or explicitly allowing the partial-payment combination.
   - Any future v2 proto or strict validator MUST preserve support for `PaymentRequestV2` partial payments (e.g., either modeling a separate cumulative cap alongside `PredeterminedBalances` or gating strict exclusivity to non-scaled/fixed modes).
2. **Ownership-Time Mode (`DurationFromTimestamp` XOR `RecurringOwnershipTimes` XOR `IncrementOwnershipTimesBy`):**
   - Already mutually exclusive at runtime; should be enforced explicitly in future validation (and modeled as `oneof ownership_time_mode` in v2 proto) so illegal combinations cannot pass pre-flight or builder validation.
3. **Party-Relation Constraints (Replacing 4 Overlapping Booleans):**
   - Setting both `RequireToEqualsInitiatedBy: true` and `RequireToDoesNotEqualInitiatedBy: true` (or the `From` equivalents) is unsatisfiable (always rejects every transfer).
   - Future consensus validation and builders should reject contradictory boolean pairs and contradictory `AddressChecks` (`MustBeEvmContract && MustNotBeEvmContract`, `MustBeLiquidityPool && MustNotBeLiquidityPool`).
   - Planned v2 proto replaces the 4 booleans with `PartyRelation to_initiator_relation` and `PartyRelation from_initiator_relation` (`UNSPECIFIED | MUST_EQUAL | MUST_NOT_EQUAL`).
4. **Structured Address-List Predicates (Eliminating `H-14` Delimiter Injection):**
   - Future validation should validate list-ID syntax strictly: rejecting empty tokens (`bb1a::bb1b`), leading/trailing colons, double negations (`!!Mint`), unmatched parentheses (`!(Mint`), and whitespace/control characters.
   - Planned v2 proto adds a typed `AddressPredicate` (`oneof { BuiltinSet builtin; string stored_list_id; AddressSet inline_addresses; AddressPredicate not; }`) while maintaining lossless conversion from legacy `listId` strings for backwards compatibility.

### 2.3 Sane Protocol Defaults (Eliminating Mandatory "Gotcha" Flags)

Future SDK builders, CLI compilers (`BB-360`), and optional upgrade-gated normalizers should eliminate the three most common setup failures:

1. **Implicit Mint Outgoing Override (`fromListId == "Mint"`):**
   - The `"Mint"` pseudo-address has no private key and cannot sign `MsgUpdateUserApprovals`. Forgetting `overridesFromOutgoingApprovals: true` on a mint approval silently causes every mint transfer to fail on outgoing user-approval checks.
   - Builders/compilers should automatically ensure `ApprovalCriteria != nil` and set `OverridesFromOutgoingApprovals = true` whenever `FromListId == "Mint"`.
2. **Automatic `MustPrioritize` for Backed Minting & Special Wrapping:**
   - Existing keeper guardrails (`validate_backed_minting_guardrails_test.go`) reject approvals with `AllowBackedMinting: true` when `MustPrioritize: false`.
   - Builders/compilers should automatically set `MustPrioritize = true` whenever `AllowBackedMinting == true` or `AllowSpecialWrapping == true`, alongside `EnforceMustPrioritizeForNonAutoScannable`.
3. **Safe `DefaultBalances.AutoApproveAllIncomingTransfers` Default (Opt-Out):**
   - Leaving `defaultBalances` nil or zero-initialized defaults `AutoApproveAllIncomingTransfers = false`, causing airdrops and third-party mints to fail unless the recipient first submits a tx or collection approvals override incoming approvals.
   - Builders/compilers should default `AutoApproveAllIncomingTransfers = true` unless explicitly opted out, plus `AutoApproveSelfInitiatedOutgoingTransfers = true` and `AutoApproveSelfInitiatedIncomingTransfers = true`.

### 2.4 Explicit Constraint: Keeping `permissions.proto` Unchanged

**`proto/tokenization/permissions.proto` MUST remain completely unchanged.**

#### Rationale
- `CollectionPermissions` and `UserPermissions` (including `ActionPermission`, `TokenIdsActionPermission`, `CollectionApprovalPermission`, and the neutral / `permanentlyPermittedTimes` / `permanentlyForbiddenTimes` first-match semantics) are persisted across all existing mainnet collections and user balance stores.
- Altering `permissions.proto` would require a consensus-breaking store migration with high risk of inadvertently unlocking permanently frozen mainnet guarantees.
- Instead, neutral-vs-permanently-permitted UX ergonomics are handled strictly **above consensus** in:
  - The `BitBadgesSpec` declarative compiler (`BB-360`), which defaults permissions to safe explicit freeze profiles (`immutable`, `mutable-metadata-only`, `managed`).
  - The deterministic verifier (tracked under the dedicated Standards Verification epic `BB-374`), which audits `CollectionPermissions` for actual permanence.
  - UI/MCP guidance and `bb review` diagnostics.

## 3. Wire Compatibility, Consensus Impact & Mainnet Upgrade Plan

1. **No `.proto` schema changes in this PR:** `proto/tokenization/permissions.proto` and all other `.proto` files remain completely untouched; existing store keys and wire encodings remain intact.
2. **Consensus-affecting nature of `EvmQueryChallenges` cast preservation & nil-approval rejection:**
   - Preserving `EvmQueryChallenges` in the 4 user-approval cast functions (`approved_transfers_casts.go`) **IS a consensus-affecting bug fix**, because previously any `UserOutgoingApproval` or `UserIncomingApproval` carrying `EvmQueryChallenges` had that field silently dropped during validation and transfer evaluation. Deploying this fix onto an existing chain requires a **coordinated chain upgrade** and **historical replay verification** (confirming no pre-upgrade block contained user approvals with `EvmQueryChallenges` whose evaluation outcome would diverge on replay).
   - Likewise, rejecting `nil` approval entries explicitly in `ValidateUserOutgoingApprovals` / `ValidateUserIncomingApprovals` / `ValidateCollectionApprovals` turns a prior nil-pointer dereference panic into a deterministic `ErrInvalidRequest` error, and should ship alongside coordinated upgrade verification.
3. **Future ergonomics validators & v2 proto upgrade (if adopted):** Any consensus-enforced ergonomics rules (such as §2.2/§2.3 rules, carefully preserving `PaymentRequestV2` partial payments) or v2 unified criteria proto MUST be gated behind a coordinated `x/upgrade` height/handler with store migration or a translation shim for v1 clients during a deprecation window.
