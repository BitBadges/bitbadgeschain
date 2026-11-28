# ADR: Core Approval Ergonomics (2.1–2.4)

- **Ticket:** BB-375
- **Status:** Proposed / Phased Implementation
- **Scope:** `x/tokenization/types` (`approved_transfers_casts.go`, `approval_ergonomics.go`, `validate_basic.go`), `proto/tokenization/approval_criteria.proto` (additive evolution), **`proto/tokenization/permissions.proto` EXPLICITLY UNCHANGED**

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
- **Phase 1 (Zero-Breaking Immediate Hardening — Implemented Now):**
  - Keep existing wire protobuf messages intact so zero on-chain state migration is required.
  - Add nil-safety guards to all 4 transfer cast functions in `x/tokenization/types/approved_transfers_casts.go`.
  - Enforce a **Go reflection round-trip invariant test** (`TestApprovalCriteriaCastReflectionRoundTrip`) that enumerates EVERY exported field of `ApprovalCriteria`, `OutgoingApprovalCriteria`, and `IncomingApprovalCriteria` via `reflect.TypeOf`, verifies the exact partition of shared vs tier-exclusive fields, populates every field with a non-zero sentinel value, and asserts full round-trip fidelity across all 4 casts. Any future proto field added without updating all 4 casts fails `go test` immediately.
- **Phase 2 (Next Major Wire Upgrade — Unified Protobuf Criteria):**
  - Introduce a single canonical `ApprovalCriteria` message across `CollectionApproval`, `UserOutgoingApproval`, and `UserIncomingApproval` (or a `v2` message family sharing tag numbers `1..26`) paired with a single `ValidateCriteriaForTier(tier ApprovalTier, c *ApprovalCriteria)` validator:
    - `TierOutgoing`: rejects `SenderChecks`, `RequireFrom*`, `Overrides*`, `AllowBackedMinting`, `AllowSpecialWrapping`, `UserApprovalSettings`.
    - `TierIncoming`: rejects `RecipientChecks`, `RequireTo*`, `Overrides*`, `AllowBackedMinting`, `AllowSpecialWrapping`, `UserApprovalSettings`.
    - `TierCollection`: permits all fields.
  - Eliminates manual 4-way cast functions entirely and shrinks LLM/MCP tool schema definitions by ~3x.

### 2.2 Typed Mutual-Exclusivity (`oneof`) & Structured Address-List Predicates

#### Mutual-Exclusivity Rules (Enforced Immediately in `approval_ergonomics.go` + Future Proto `oneof`)
1. **Quantity / Balance Mode (`PredeterminedBalances` XOR `ApprovalAmounts`):**
   - `PredeterminedBalances` dictates exact per-transfer balances (incremented, scaled, or manual), whereas `ApprovalAmounts` tracks arbitrary flexible transfer amounts up to a quota. Combining both produces confusing silent rejections when predetermined outputs conflict with cumulative caps.
   - Strict validation helper `ValidateApprovalCriteriaErgonomics` rejects any criteria where both `PredeterminedBalances` (non-nil/non-empty) and `ApprovalAmounts` (non-nil/non-empty) are simultaneously active.
   - Future v2 proto models this as `oneof balance_constraint { PredeterminedBalances predetermined = 2; ApprovalAmounts amounts = 3; }`.
2. **Ownership-Time Mode (`DurationFromTimestamp` XOR `RecurringOwnershipTimes` XOR `IncrementOwnershipTimesBy`):**
   - Already mutually exclusive at runtime; enforced explicitly in `ValidateApprovalCriteriaErgonomics` (and planned as `oneof ownership_time_mode` in v2 proto) so illegal combinations cannot pass pre-flight or builder validation.
3. **Party-Relation Constraints (Replacing 4 Overlapping Booleans):**
   - Setting both `RequireToEqualsInitiatedBy: true` and `RequireToDoesNotEqualInitiatedBy: true` (or the `From` equivalents) is unsatisfiable (always rejects every transfer).
   - `ValidateApprovalCriteriaErgonomics` rejects contradictory boolean pairs and contradictory `AddressChecks` (`MustBeEvmContract && MustNotBeEvmContract`, `MustBeLiquidityPool && MustNotBeLiquidityPool`).
   - Planned v2 proto replaces the 4 booleans with `PartyRelation to_initiator_relation` and `PartyRelation from_initiator_relation` (`UNSPECIFIED | MUST_EQUAL | MUST_NOT_EQUAL`).
4. **Structured Address-List Predicates (Eliminating `H-14` Delimiter Injection):**
   - `ValidateAddressListPredicate` validates list-ID syntax strictly: rejects empty tokens (`bb1a::bb1b`), leading/trailing colons, double negations (`!!Mint`), unmatched parentheses (`!(Mint`), and whitespace/control characters.
   - Planned v2 proto adds a typed `AddressPredicate` (`oneof { BuiltinSet builtin; string stored_list_id; AddressSet inline_addresses; AddressPredicate not; }`) while maintaining lossless conversion from legacy `listId` strings for backwards compatibility.

### 2.3 Sane Protocol Defaults (Eliminating Mandatory "Gotcha" Flags)

`x/tokenization/types/approval_ergonomics.go` provides deterministic normalization helpers (`NormalizeCollectionApprovalDefaults`, `NormalizeDefaultBalances`) that eliminate the three most common setup failures:

1. **Implicit Mint Outgoing Override (`fromListId == "Mint"`):**
   - The `"Mint"` pseudo-address has no private key and cannot sign `MsgUpdateUserApprovals`. Forgetting `overridesFromOutgoingApprovals: true` on a mint approval silently causes every mint transfer to fail on outgoing user-approval checks.
   - `NormalizeCollectionApprovalDefaults` automatically ensures `ApprovalCriteria != nil` and sets `OverridesFromOutgoingApprovals = true` whenever `FromListId == "Mint"`.
2. **Automatic `MustPrioritize` for Backed Minting & Special Wrapping:**
   - Existing keeper guardrails (`validate_backed_minting_guardrails_test.go`) reject approvals with `AllowBackedMinting: true` when `MustPrioritize: false`.
   - `NormalizeCollectionApprovalDefaults` automatically sets `MustPrioritize = true` whenever `AllowBackedMinting == true` or `AllowSpecialWrapping == true`, and also applies `EnforceMustPrioritizeForNonAutoScannable`.
3. **Safe `DefaultBalances.AutoApproveAllIncomingTransfers` Default (Opt-Out):**
   - Leaving `defaultBalances` nil or zero-initialized defaults `AutoApproveAllIncomingTransfers = false`, causing airdrops and third-party mints to fail unless the recipient first submits a tx or collection approvals override incoming approvals.
   - `NormalizeDefaultBalances(defaults, optOutAutoApproveIncoming)` initializes `AutoApproveAllIncomingTransfers = !optOutAutoApproveIncoming` (defaulting to `true` unless explicitly opted out), plus `AutoApproveSelfInitiatedOutgoingTransfers = true` and `AutoApproveSelfInitiatedIncomingTransfers = true`.

### 2.4 Explicit Constraint: Keeping `permissions.proto` Unchanged

**`proto/tokenization/permissions.proto` MUST remain completely unchanged.**

#### Rationale
- `CollectionPermissions` and `UserPermissions` (including `ActionPermission`, `TokenIdsActionPermission`, `CollectionApprovalPermission`, and the neutral / `permanentlyPermittedTimes` / `permanentlyForbiddenTimes` first-match semantics) are persisted across all existing mainnet collections and user balance stores.
- Altering `permissions.proto` would require a consensus-breaking store migration with high risk of inadvertently unlocking permanently frozen mainnet guarantees.
- Instead, neutral-vs-permanently-permitted UX ergonomics are handled strictly **above consensus** in:
  - The `BitBadgesSpec` declarative compiler (`BB-360`), which defaults permissions to safe explicit freeze profiles (`immutable`, `mutable-metadata-only`, `managed`).
  - The deterministic verifier (`BB-374` `verify_standards.go`), which audits `CollectionPermissions` for actual permanence.
  - UI/MCP guidance and `bb review` diagnostics.

## 3. Wire Compatibility & Mainnet Upgrade Plan

1. **Zero breaking changes in this PR:** No `.proto` files are modified; all existing mainnet collections, store keys, and JSON payloads remain 100% compatible.
2. **Additive validation & normalization:** New ergonomics helpers live in `x/tokenization/types/approval_ergonomics.go` for opt-in use by builders, CLI/MCP, and future upgrade gates (leaving existing historical mainnet collections unaffected).
3. **Future v2 proto upgrade (if adopted):** Store migration upgrades existing collections lazily or via an `x/upgrade` handler while maintaining a JSON/proto translation shim for v1 clients during a deprecation window.
