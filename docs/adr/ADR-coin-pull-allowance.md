# ADR: Native Scoped Pull Allowances (`CoinPullAllowance`) for Recurring Subscriptions & Agent Budgets

- **Ticket:** BB-379
- **Aligned With:** `BB-146` (Subscriptions Redesign — Plan Changes, Proration, Renewal Consent Separation)
- **Status:** Scoping ADR (Implementation deferred to dedicated chain upgrade PR)
- **Scope:** `bitbadgeschain` (`x/tokenization` + `x/bank`), parity analysis for `bitbadges-solana` (SPL Token Delegate / PDA)

## 1. Context & Alignment with Subscriptions Redesign (`BB-146`)

In `BB-146`, recurring subscriptions on BitBadges established a strict separation between **paid access** (time-bounded token ownership ranges) and **renewal consent** (user-authorized recurring charges). Today, recurring auto-debit relies on `UserIncomingApproval` with `coinTransfers` (`overrideFromWithApproverAddress: true`) + `recurringOwnershipTimes`. While functional for static plans, this approach has sharp ergonomic and operational edges:
1. **Brittle external consent management:** Subscribers and operators must carefully match incoming approvals, versions, and charge periods; legacy consents cannot easily express an explicit periodic spending cap independent of token transfer matching.
2. **Plan upgrades & proration (`BB-146`):** Immediate prorated upgrades require atomically retiring old renewal consent and binding the new tier's exact terms without allowing the merchant to silently raise the recurring price.
3. **Agent Vault / AI Agent budgets:** Autonomous agents need bounded recurring spending allowances (e.g., *"up to 50 USDC per 30 days until Dec 31, revocable in 1 tx"*) without handing out unbounded `x/authz` `SendAuthorization` grants.

### Why Standard Cosmos `x/authz` `SendAuthorization` Is Insufficient
- `x/authz` `SendAuthorization` only supports a single **total** `SpendLimit` — it has **no per-period reset interval** (`resetTimeInterval`), **no binding to collection approval version / immutable terms hash**, and **no atomic coupling** to minting a subscription entitlement token!

## 2. Proposed Primitive: `CoinPullAllowance`

A subscriber signs **one** consent transaction (`MsgGrantCoinPullAllowance`) creating a scoped pull allowance stored in `x/tokenization`:

```protobuf
message CoinPullAllowance {
  string allowanceId = 1;
  string subscriber = 2;          // Payer whose bank coins are debited (or Agent Vault address)
  string beneficiary = 3;         // Recipient of the minted subscription token (defaults to subscriber; supports BB-146 sponsored beneficiaries)
  uint64 collectionId = 4;
  string approvalId = 5;
  uint64 approvalVersion = 6;     // Pinned collection approval version
  bytes termsHash = 7;            // SHA-256 over canonical deterministic protobuf bytes of PullAllowanceTerms (NOT `|` string concatenation)
  string denom = 8;
  string maxAmountPerPeriod = 9;  // Hard cap per billing period (never overcharged)
  uint64 startTimeMs = 10;
  uint64 periodLengthMs = 11;     // Fixed interval (e.g., 30d = 2_592_000_000 ms; not calendar billing per BB-146)
  uint64 validUntilMs = 12;       // Hard expiry timestamp
  string maxExecutorFee = 13;     // Optional bounded fee paid to permissionless keeper/executor
  bool revoked = 14;
}

message PullAllowanceTerms {
  uint64 collectionId = 1;
  string approvalId = 2;
  uint64 approvalVersion = 3;
  string priceDenom = 4;
  string priceAmount = 5;
  string payoutAddress = 6;
  uint64 periodLengthMs = 7;
}
```

### Core Safety Guarantees

1. **Hard Per-Period Cap & Executor Fee Validation (Never Overcharged):**
   - For any pull at timestamp `nowMs`, the current period index is computed deterministically: `periodIdx = (nowMs - startTimeMs) / periodLengthMs`.
   - `spentInPeriod + pullAmount + executorFee <= maxAmountPerPeriod` and `executorFee <= maxExecutorFee` are enforced unconditionally.
   - **Non-empty `executorAddress` requirement:** Whenever `executorFee > 0`, `MsgExecutePullRenewal` MUST require a valid, non-empty `executorAddress` (rejecting empty or malformed addresses before any debit occurs) so executor fees can never be burned, lost, or misrouted.
   - Period quota resets cleanly when `periodIdx` advances; unused quota does **not** roll over (preventing burst drains after dormancy).
2. **Immutable Terms Binding via Canonical Protobuf Hashing (Merchant Price-Mutation Rejection):**
   - Because `CoinPullAllowance` pins BOTH `approvalVersion` AND `termsHash = SHA256(proto.MarshalDeterministic(PullAllowanceTerms{...}))` using canonical deterministic protobuf serialization (explicitly avoiding delimiter-based `|` string concatenation, which is vulnerable to delimiter-collision ambiguity across string fields such as `approvalId`, `priceDenom`, or `payoutAddress`), a merchant **cannot** raise the subscription price, shorten the period, or redirect payout funds after consent.
   - Any modification to the collection approval bumps `approvalVersion` and changes `termsHash` -> pulls are immediately rejected (`ErrTermsMutated`) until the subscriber explicitly signs a new consent!
   - Aligns directly with `BB-146`'s rule: *"existing incoming consent contains explicit values and is not rewritten by changing collection terms."*
3. **Instant 1-Tx Revocation:**
   - Subscriber calls `MsgRevokeCoinPullAllowance(allowanceId)` -> sets `revoked = true` immediately.
   - Only `subscriber` (never merchant or executor) can revoke.
   - Revoking renewal consent **never revokes already-paid ownership intervals** (preserving `BB-146` paid-access vs. renewal-consent separation).
4. **Graceful Insufficient-Balance Handling:**
   - If the subscriber's bank balance is `< pullAmount + executorFee`, the renewal executor receives a structured `PullStatusInsufficientBalance` outcome.
   - The allowance is **NOT bricked**, its period cap is **NOT consumed**, and no partial debit occurs.
   - If the subscriber tops up before expiry/grace ends, a subsequent retry succeeds; otherwise the subscription token's time-bounded ownership naturally expires at the end of the last paid interval without corrupting state.

## 3. Execution Model: Permissionless Keeper vs Native `x/tokenization` Trigger

| Model | Mechanism | Pros | Cons |
| --- | --- | --- | --- |
| **Permissionless Keeper (`MsgExecutePullRenewal`)** | Any relayer/merchant/bot submits `MsgExecutePullRenewal(allowanceId)` during the charge window and receives up to `maxExecutorFee` (bounded inside the allowance) | Zero `EndBlocker` consensus bloat; scales to millions of subscriptions; predictable gas; works identically on Solana (crank) | Requires off-chain crank / merchant worker (already exists in BitBadges payment worker) |
| **Native `EndBlocker` Auto-Trigger** | Chain automatically iterates due allowances every block | No external relayer | Unbounded block-time DoS vector; hard to meter gas fairly; impossible on Solana |

**Decision:** Use **Permissionless Keeper Execution (`MsgExecutePullRenewal`)** with an optional, subscriber-capped `maxExecutorFee` (default `0`, where merchants/operators absorb gas as a cost of billing).

## 4. Alignment with `BB-146` Plan Changes & Agent Vaults

- **Immediate Prorated Upgrade (`BB-146`):** An atomic upgrade transaction (1) revokes the old tier's `CoinPullAllowance`, (2) settles the prorated price difference for the remaining period, (3) swaps the old-tier entitlement for the new-tier entitlement over the remaining period boundary, and (4) grants a new `CoinPullAllowance` pinned to the new tier's `termsHash`.
- **Scheduled Downgrade (`BB-146`):** Revoke the old tier's allowance (or set `validUntilMs = currentPeriodEnd`) and grant a new allowance with `startTimeMs = currentPeriodEnd` pinned to the lower tier's `termsHash`.
- **Sponsored Beneficiaries (`BB-146`):** `subscriber` (payer) and `beneficiary` (token recipient) are distinct fields, allowing an employer or Agent Vault to fund a user's subscription cleanly.
- **Agent Vault Expiry:** If `subscriber` is an Agent Vault, once the vault expires or is swept, insufficient balance gracefully stops renewals without trapping funds.

## 5. Cross-Chain Feasibility (Cosmos + Solana)

- **Cosmos (`bitbadgeschain`):** Implemented natively inside `x/tokenization` using a dedicated module-scoped debit path via `bankKeeper.SendCoins` (gated strictly by `CoinPullAllowance` verification, never general arbitrary debit) inside `AtomicContext` so payment + token mint succeed or roll back atomically.
- **Solana (`bitbadges-solana`):** Implemented via an Anchor `PullAllowance` PDA + SPL Token `ApproveChecked` delegation to the program's vault authority PDA. The program enforces `terms_hash`, `period_length`, `max_amount_per_period`, `valid_until`, and `revoked` before invoking `token::transfer_checked` and minting the time-bounded subscription entitlement in the same instruction.

## 6. Planned Implementation & Verification Criteria

Implementation is deferred to a dedicated chain upgrade PR introducing `CoinPullAllowance`, `PullAllowanceTerms`, `MsgGrantCoinPullAllowance`, `MsgRevokeCoinPullAllowance`, and `MsgExecutePullRenewal`. Acceptance tests must cover:
1. Initial subscriber consent grant & first charge.
2. Periodic capped pulls (rejects over-cap pulls within the same period, rolls over cleanly on next period boundary).
3. Non-empty, valid `executorAddress` enforcement whenever `executorFee > 0` (and `executorFee <= maxExecutorFee`).
4. Insufficient-balance grace handling (zero corruption, zero quota loss, succeeds after top-up).
5. Instant mid-period revocation by subscriber (and rejection of unauthorized revocation attempts).
6. Merchant price / payout / version mutation rejection via canonical protobuf `TermsHash` & `ApprovalVersion` verification.
