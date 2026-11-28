# ADR: Read-Only Policy Introspection (`EvaluateCanTransfer` / `QueryCanTransfer`) vs Native Cosmos `Simulate`

- **Ticket:** BB-378
- **Status:** Scoping ADR (Implementation deferred to dedicated gRPC `Query/CanTransfer` PR)
- **Scope:** `bitbadgeschain` (`x/tokenization/keeper`, gRPC `Query/CanTransfer`, EVM precompile `0x1001`), `bitbadges-solana`, `bb-api`, `bb-mcp`

## 1. Context & Problem Statement

Wallets, frontends, EVM contracts, and AI agents frequently need to answer a policy question **before** constructing and signing a transaction:

> *"Can address `A` transfer tokens `X` to `B` at time `T` — and if not, which exact approval criterion failed, and what coins are required?"*

Today, callers must either:
1. Build a dummy `MsgTransferTokens` transaction, fabricate signatures/account sequence/gas, call Cosmos gRPC `tx.Service/Simulate`, and regex-parse opaque error strings; OR
2. Re-implement BitBadges' approval-matching engine off-chain (causing drift and false positives/negatives).

## 2. Tradeoff Analysis: `EvaluateCanTransfer` (`QueryCanTransfer`) vs Native Cosmos `Simulate`

| Dimension | `EvaluateCanTransfer` / `QueryCanTransfer` | Native Cosmos `tx.Service/Simulate` |
| --- | --- | --- |
| **Output Structure** | Structured `{ allowed, matchedApprovalIds, failedCriteria: [{ approvalId, approvalLevel, criterion, expected, actual, message }], requiredCoins }` | Binary pass/fail + `GasInfo` + single unstructured error string |
| **Boilerplate Required** | None — only `(collectionId, transfer, initiatedBy, atTimestampMs?)` | Full `TxRaw`, `AuthInfo`, account sequence lookup, dummy public key/signature, fee/gas envelope |
| **Time-Travel (`atTimestamp`)** | **Native** — overrides `cacheCtx.WithBlockTime(atTime)` cleanly without mutating node clock | **Not supported** — always runs at current block time |
| **EVM Precompile (`0x1001`)** | **Callable via `STATICCALL`** from Solidity contracts (pure read-only view) | **Impossible** from EVM contracts |
| **State Safety** | Runs inside `ctx.CacheContext()` and **discards** commit handle (zero store writes, zero events emitted to parent ctx) | Runs in `BaseApp` simulate cache (read-only) |
| **AnteHandler / Gas / Fees** | **Not covered** — does not estimate tx gas, check signature validity, or verify tx fee payer balance | **Full coverage** — executes full `AnteHandler`, gas metering, fee deduction, and module hooks |
| **Cross-Module Effects** | Evaluates `x/tokenization` + its bank/evm/gamm keeper calls inside `CacheContext` | Evaluates full multi-message tx including other Cosmos/EVM messages |
| **Latency (in-process / gRPC)** | Fast (~0.2–1.5 ms; skips tx decode, protobuf `Any` unpacking, AnteHandler, signature verification) | Moderate (~3–15 ms; full tx decode + AnteHandler + BaseApp routing + gas meter) |

### Where `EvaluateCanTransfer` Wins
1. **Actionable per-criterion diagnostics:** When a transfer fails, `EvaluateCanTransfer` inspects each candidate approval inside isolated `CacheContext` branches and reports structured `FailedCriteria` (e.g., `criterion: "transferTime"`, `criterion: "coinTransfer"`, `criterion: "mustOwnTokens"`, `criterion: "maxNumTransfers"`, `criterion: "merkleChallenge"`, `criterion: "userOutgoingApproval"`, `criterion: "userIncomingApproval"`, `criterion: "insufficientBalance"`) plus `RequiredCoins` even on failure.
2. **Future/past time-window preview (`AtTimestampMs`):** Enables *"Will this subscription/claim still be valid tomorrow?"* or *"When does the refund window open?"* without advancing chain height.
3. **EVM smart contract composability:** Exposing `EvaluateCanTransfer` on the tokenization precompile (`0x1001`) lets Solidity dApps query transfer eligibility via `STATICCALL`.
4. **Zero account/signature friction:** Works for addresses that have never transacted on-chain (no on-chain account number / sequence yet).

### Where Native Cosmos `Simulate` Wins
1. **Accurate gas estimation & fee calculation:** Required before broadcasting any real Cosmos transaction.
2. **Full `AnteHandler` & multi-message atomicity:** Necessary when a transaction bundles multiple messages (e.g., `MsgBankSend` + `MsgTransferTokens` or `MsgUpdateUserApprovals` + `MsgTransferTokens`).
3. **Full state-diff simulation (`BB-362`):** Required when showing exact final token + bank balance diffs before user signature.

## 3. Hybrid Backend Strategy (`bb-api` & `bb-mcp`)

Neither mechanism replaces the other; they complement each other in a two-tier pipeline:

```text
[1. Eligibility / UI / Agent Exploration]
   -> Call QueryCanTransfer / EvaluateCanTransfer
   -> Fast (~1ms), no dummy tx, supports atTimestamp, returns structured FailedCriteria + RequiredCoins
   -> Powers UI eligibility badges, "Why can't I claim?" hints, agent pre-checks, EVM STATICCALL

[2. Pre-Broadcast Verification & Signing]
   -> Build real tx -> Call Cosmos Simulate (enriched via BB-362 balance diffs)
   -> Verifies AnteHandler, accurate gas limit, fee balance, and multi-msg atomicity
```

## 4. Planned Implementation Architecture (Dedicated gRPC `Query/CanTransfer` PR)

1. **Strict Read-Only Isolation:**
   - Entry point `Keeper.EvaluateCanTransfer(ctx sdk.Context, req CanTransferRequest) CanTransferResult` immediately creates `cacheCtx, _ := ctx.CacheContext()`.
   - The commit callback is discarded (`_`), guaranteeing **zero state mutations** (balances, approval trackers, auto-deletions, bank coin transfers, events) can ever leak back into `ctx`.
   - Caller-supplied `*types.Transfer` is deep-copied before evaluation so in-place normalization inside `HandleTransfersEnriched` never mutates caller memory.
2. **Optional `AtTimestampMs` Override:**
   - When `req.AtTimestampMs > 0` (and within `math.MaxInt64`), `cacheCtx = cacheCtx.WithBlockTime(time.UnixMilli(int64(req.AtTimestampMs)))`.
3. **Exact Execution Parity & Accurate `RequiredCoins` Accounting:**
   - Runs the real `k.HandleTransfersEnriched(cacheCtx, collection, []*types.Transfer{transferCopy}, req.InitiatedBy)` path so `Allowed == true` matches real consensus execution 1:1 (zero logic drift!).
   - **`RequiredCoins` semantics (Critical Design Requirement):**
     - `RequiredCoins` MUST strictly distinguish **initiator-paid coins** (coins debited from `initiatedBy`) from **escrow / approver-paid payouts** (`OverrideFromWithApproverAddress: true`, where coins are debited from the collection/approver escrow rather than from the initiator — e.g. bounty payouts, crowdfund withdrawals, or refunds MUST NOT be misreported as coins the initiator must pay!).
     - `RequiredCoins` MUST also include **backed / vault bank deposits** required when minting or wrapping via `CosmosCoinBackedPath` (e.g., `AllowBackedMinting: true` vault/smart-token deposits), not just `ApprovalCriteria.CoinTransfers`.
4. **Structured Diagnostic Pass on Failure:**
   - When execution fails, `diagnoseTransferFailure` inspects validation, archive status, sender token balances, every candidate collection approval (respecting `OnlyCheckPrioritizedCollectionApprovals`), and user outgoing/incoming approvals — each inside its own discarded `CacheContext`.

## 5. Wire Surface & Rollout (Dedicated Follow-Up PR)

- Implement `EvaluateCanTransfer` and wire it directly to a new additive gRPC query endpoint `Query/CanTransfer` in `proto/tokenization/query.proto` (plus EVM precompile `0x1001` `canTransfer(...) external view returns (...)`) in a single cohesive PR rather than shipping unwired keeper code.
- Mirror the same response shape (including initiator-paid vs escrow-paid coin distinction and backed/vault deposits) in `bitbadges-solana` view planner and `bb-api` / `bb-mcp`.
