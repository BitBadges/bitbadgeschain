# ADR: Deterministic Standard Verification (`VerifiedStandards` vs `ClaimedStandards`)

- **Ticket:** BB-374
- **Status:** Scoping ADR (Implementation tracked under dedicated Standards Verification epic `BB-374`)
- **Scope:** `BitBadges/bitbadgeschain` (`x/tokenization/types`), parity with `bitbadges-solana`, `services/indexer-rs/crates/bb-standards`, and `bitbadgesjs`

## 1. Context & Problem Statement

Today, `TokenCollection.Standards` (`repeated string standards = 8` in `proto/tokenization/collections.proto`) is an unvalidated, creator-declared array of strings. Any collection creator can declare `standards: ["Non-Transferable"]` while leaving `collectionPermissions.canUpdateCollectionApprovals` unfrozen (allowing them to enable transfers later), or tag `"Subscriptions"` without any `durationFromTimestamp` / `recurringOwnershipTimes` criteria, or tag `"Bounty"` without mutual exclusion between payout and refund approvals.

Wallets, indexers, frontends, and AI agents currently have no protocol-level way to distinguish between **self-declared intent** (`ClaimedStandards`) and **cryptographically / structurally verified conformance** (`VerifiedStandards`).

### Explicit Non-Goal: No Standard Consolidation

Per BB-374 scope, this ADR **does NOT consolidate, merge, rename, or retire any existing standards**. The full existing standard catalog (`Non-Transferable`, `Subscriptions`, `NFTs`, `Fungible Tokens`, `Smart Token`, `Address List`, `Vault`, `Agent Wallet`, `PaymentRequestV2`, `Bounty`, `Crowdfund`, `Spendable Credit`) remains intact. Consolidation may be evaluated in a separate future ADR.

## 2. Decision

1. **Separate `ClaimedStandards` from `VerifiedStandards` + `FailedStandards`:**
   - `ClaimedStandards`: The raw, untrusted `collection.Standards` tags declared on-chain.
   - `VerifiedStandards`: The deterministic subset of `ClaimedStandards` whose on-chain structural predicates pass with zero violations.
   - `FailedStandards`: A map of `standard -> []StandardViolation{Code, Message}` explaining every exact invariant that failed (including `UNKNOWN_STANDARD` for unrecognized tags).
2. **Proposed pure, deterministic, zero-store-read verification engine (planned for `x/tokenization/types/verify_standards.go` in the dedicated `BB-374` PR):**
   - Will evaluate a `*TokenCollection` purely from its in-memory state (`CollectionApprovals`, `CollectionPermissions`, `Invariants`, `ValidTokenIds`, `DefaultBalances`, `AliasPaths`, `CosmosCoinWrapperPaths`).
   - Can be called identically on-chain (e.g., query responses, optional strict validation on `MsgCreateCollection` / `MsgUpdateCollection`) and mirrored 1:1 in Rust (`crates/bb-standards`) and TypeScript (`bitbadgesjs/src/core/verify-standards.ts`).

## 3. On-Chain vs Indexer Verification Boundaries

| Dimension | On-Chain (`x/tokenization/types`) | Indexer (`bb-standards` / `bb-api` / `bb-mcp`) |
| --- | --- | --- |
| **Inputs** | `TokenCollection` struct (approvals, permissions, invariants, validTokenIds) | `TokenCollection` + resolved `AddressList` store contents + off-chain IPFS/HTTP metadata JSON + cross-collection state |
| **Guarantees** | Consensus-deterministic, trustless, verifiable by light clients & EVM precompiles | Enriched verification (can check off-chain metadata schema, external list membership, historical funding) |
| **Failure Mode** | Never depends on network/IPFS; fails closed if an invariant cannot be proven from on-chain fields | Marks metadata-dependent sub-checks separately from structural on-chain verification |
| **Gas / Latency** | O(A + P) struct traversal (~1–5 µs; < 5,000 gas equivalent) | Computed at index time and cached in `CollectionIndexRecord.verifiedStandards` |

### Gas Cost Analysis

Each proposed on-chain verifier in `verify_standards.go` will perform a bounded linear scan over `collection.CollectionApprovals` (length `A`) and `collection.CollectionPermissions.CanUpdateCollectionApprovals` (length `P`), with zero IAVL store reads:

- **Time complexity:** `O(S * (A + P))` where `S = len(collection.Standards)` (typically `1..4`), `A` is bounded by collection approval limits, and `P` is bounded by permission list size.
- **Estimated gas cost:** ~500 base gas + ~150 gas per approval/permission inspected (~2,000–5,000 gas typical), negligible relative to `MsgCreateCollection` (~150,000–300,000 gas).
- **Read-path cost:** Free when served via gRPC query / `bb-api` / `bb-mcp` off-chain execution.

## 4. Backwards Compatibility for Pre-Existing Mainnet Collections

- **Zero state migration required:** `proto/tokenization/collections.proto` field `standards` is untouched. Existing mainnet collections remain 100% readable and valid.
- **Non-breaking rollout:**
  1. **Phase 1 (Planned — Dedicated `BB-374` PR):** Expose `VerifyCollectionStandards(collection)` as a pure deterministic function used by queries, indexers (`bb-standards`), SDK (`bitbadgesjs`), and AI tooling (`bb-mcp`) to surface `{ claimedStandards, verifiedStandards, failedStandards }`. Pre-existing collections with legacy/incomplete configs simply report unverified tags in `failedStandards` without breaking transfers or updates.
  2. **Phase 2 (Opt-in / Governance-gated):** Add an optional `strictStandardVerification` flag on `MsgCreateCollection` / `MsgUpdateCollection` (or governance param) that rejects transactions claiming a standard they do not satisfy. Pre-existing collections are grandfathered unless they opt in.

## 5. Canonical Verification Matrix (12 Standards)

All 12 standards are checked deterministically with zero consolidation:

| Standard | Required On-Chain Invariants | Rejected Spoof Examples |
| --- | --- | --- |
| `Non-Transferable` | No post-mint peer-to-peer collection approval (`FromListId != "Mint"` with active ranges), except optional legitimate burn paths (e.g., holder-initiated transfers exclusively to a designated unspendable burn address); `canUpdateCollectionApprovals` permanently forbids updating post-mint approvals across `[1, MaxUint64]` | Unfrozen `canUpdateCollectionApprovals`; active peer-to-peer approval |
| `Subscriptions` | `ValidTokenIds` non-empty; `NoCustomOwnershipTimes == false`; at least one active approval with `PredeterminedBalances.IncrementedBalances` setting `DurationFromTimestamp > 0` or active `RecurringOwnershipTimes` (`IntervalLength > 0`) | Missing `durationFromTimestamp` and `recurringOwnershipTimes`; `noCustomOwnershipTimes: true` |
| `NFTs` | `ValidTokenIds` non-empty; `Invariants.MaxSupplyPerId == 1`; `Invariants.NoCustomOwnershipTimes == true` | `maxSupplyPerId == 0` or `> 1`; custom ownership times enabled |
| `Fungible Tokens` | Single fungible token ID (`len(ValidTokenIds) == 1 && Start == End`); `Invariants.NoCustomOwnershipTimes == true`; `MaxSupplyPerId == 0 || MaxSupplyPerId > 1` | `maxSupplyPerId == 1` (NFT); multi-ID non-fungible range |
| `Smart Token` | Valid `Invariants.CosmosCoinBackedPath` (non-empty address + conversion); at least one approval with `AllowBackedMinting && MustPrioritize` matching backing address; no unbacked `FromListId == "Mint"` approval | Missing `cosmosCoinBackedPath`; unbacked Mint approval |
| `Address List` | `ValidTokenIds` non-empty; at least one approval scoped to a concrete non-wildcard address list ID (`FromListId`, `ToListId`, or `InitiatedByListId` not in `All`, `AllWithMint`, `AllWithoutMint`, `None`, `Mint`) | Open collection using only `AllWithoutMint` wildcards |
| `Vault` | Valid `CosmosCoinBackedPath`; both deposit (`FromListId == backingAddr`) and withdrawal (`ToListId == backingAddr`) approvals with `AllowBackedMinting && MustPrioritize`; `NoForcefulPostMintTransfers == true`; frozen `canUpdateCollectionApprovals` | Missing withdrawal path (honeypot vault); unfrozen approvals; forceful post-mint transfers allowed |
| `Agent Wallet` | Has `CosmosCoinBackedPath` or escrow-funded `CoinTransfers` (`OverrideFromWithApproverAddress`); every spending/backed approval has a bounded cap (`MaxNumTransfers`, `ApprovalAmounts`, or finite `TransferTimes`); `NoForcefulPostMintTransfers` or frozen approvals | Unbounded agent spending approval with unlimited transfers and infinite window |
| `PaymentRequestV2` | At least one approval with non-empty `CoinTransfers`; either a single-use bound (`MaxNumTransfers.OverallMaxNumTransfers == 1` with `ResetTimeIntervals` unset/zero, or `AutoDeletionOptions.AfterOneUse`) OR a bounded cumulative partial-payment cap (`ApprovalAmounts` with `ResetTimeIntervals` unset/zero, matching `ADR-approval-ergonomics.md`); frozen `canUpdateCollectionApprovals` | Unlimited repeatable invoice; resetting `ResetTimeIntervals` on payment latch/cap; unfrozen invoice terms |
| `Bounty` | At least two outcome approvals (e.g., fulfill + refund) sharing a terminal mutual-exclusion latch (`MaxNumTransfers.OverallMaxNumTransfers == 1`, identical non-empty `AmountTrackerId`, and `ResetTimeIntervals` unset/zero); escrow payout via `CoinTransfers` (`OverrideFromWithApproverAddress`); frozen `canUpdateCollectionApprovals` | Missing mutual-exclusion shared tracker (fulfill+refund double-spend); resetting `ResetTimeIntervals` on latch; unfrozen approvals |
| `Crowdfund` | Contribution approval (`FromListId == "Mint"`, `CoinTransfers` to escrow, finite deadline `TransferTimes.End < MaxUint64`); escrow release/refund approval (`OverrideFromWithApproverAddress`); frozen `canUpdateCollectionApprovals` | Infinite contribution window; missing escrow funding or release/refund branch |
| `Spendable Credit` | `NoCustomOwnershipTimes == true`; issue/mint approval (`FromListId == "Mint"` with `CoinTransfers`, `AllowAmountScaling`, or `ApprovalAmounts`) and redemption/spend approval (`FromListId != "Mint"` or backed path); `NoForcefulPostMintTransfers == true`; frozen `canUpdateCollectionApprovals` | Forceful post-mint confiscation enabled; missing redemption path; unfrozen terms |

## 6. Cross-Repo Parity Plan

- **Go (`bitbadgeschain`):** Dedicated `BB-374` Standards Verification PR introducing `x/tokenization/types/verify_standards.go` wired into query/validation surfaces.
- **Rust (`bitbadges-monorepo` / `bitbadges-solana`):** `crates/bb-standards` mirrors the exact violation codes (`NON_TRANSFERABLE_UNFROZEN_APPROVALS`, `SUBSCRIPTION_MISSING_DURATION`, `BOUNTY_MISSING_MUTUAL_EXCLUSION`, etc.).
- **TypeScript (`bitbadgesjs`):** `src/core/verify-standards.ts` shares golden pass/fail fixtures with Go and Rust to permanently eliminate `H-13` verification drift.
