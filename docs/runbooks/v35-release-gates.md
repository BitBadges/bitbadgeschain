---
title: v35 security remediation and release gates
last-verified: 2026-09-06
---

## [1] Before setting an upgrade height

The security regression suite and captured-state replay do not establish a
validator's upgrade-block capacity. Complete the in-scope checks below before setting a height:

1. Capture the migration prefixes at a recent fixed mainnet height using
   `python3 scripts/upgrade/scan-v35-state.py --height HEIGHT --out /tmp/v35-HEIGHT`.
   The RPC must retain that height and expose read-only ABCI subspace queries.
   Preserve the raw responses and counts. Recheck immediately before the proposal
   because pre-v35 state can still be inflated cheaply.
2. Full-mainnet-database rehearsal is excluded from this release scope by the
   release owner's explicit instruction on 2026-09-06. Use the populated-state
   replay and the fresh local upgrade/rollback/EVM/multivalidator rehearsal.
   These do not establish full-mainnet IAVL timing; record that limitation.
3. Measure a block filled with the most expensive allowed JSON/precompile and
   ETH-recovery work on validator-class hardware, including nested alias routing.
   The 100,000,000 block gas limit is a ceiling, not a demonstrated throughput
   budget. Lower it before release if measured time/RSS exceeds consensus margins.
4. Review IBC configuration order and coverage. Matching is first-match; put
   concrete channel rules ahead of wildcard fallbacks. Use wildcard supply-shift
   limits for denoms that must remain bounded on newly opened channels. Inbound
   address/unique-sender limits identify counterparty-provided strings, not
   independently authenticated identities; they are not a substitute for a
   supply-shift cap.
5. Recheck existing collections for empty approval token-ID ranges, missing
   predetermined-balance order methods, nonzero maximum supply with a backed
   path, ERC-3643 unlimited default allocations, and override wrapper paths
   without `{id}`. Do not silently rewrite arbitrary default allocations or
   reinterpret existing backing: any newly discovered affected collection needs
   an explicit compatible migration before release.

Captured height 11,979,578 contained 998 tokenization rows, including 170 collection
records representing 97 distinct collection IDs (both historical and current key
encodings exist), 203 balance rows, no uppercase address spellings, no ETH usage
trackers, no manager splitters, and no pair fee overrides. The exact legacy-shape
checks above found no affected record. Two historical/current records for
collection 50 contain an empty ETH challenge placeholder; neither configures a
usable signer. This is a dated observation, not a guarantee about upgrade state.

## [2] ETH claim voucher compatibility

v35 accepts only the **v2** EIP-191 signed message. Existing pre-v35 vouchers must
be reissued; there is no legacy verification fallback. This also prevents a
previously consumed legacy signature from being redeemed once more under the
nonce-based tracker scheme. Inventory external issuers before rollout even if
on-chain signer/tracker scans are empty.

The exact message is `BitBadges ETH Signature Challenge v2\n` followed by these
fields in order: chain ID, nonce, initiator, collection ID, approver address,
approval level, approval ID, challenge tracker ID. Each field is encoded as its
UTF-8 byte length in decimal, `:`, and the field bytes, with no extra separator.
Use the SDK's `getETHSignatureChallengeMessage` helper (0.45.0 release candidate)
and sign its returned string with EIP-191 `signMessage`; do not recreate the
format with delimiter concatenation. The chain's `ETHSignatureChallengeMessage`
and SDK helper share golden fixtures.

Use canonical addresses and the exact execution chain ID. Nonces contain 1–256
ASCII letters, digits, underscores or hyphens. Nonces must be unique within a
collection/approver/level/approval/challenge scope **across recipients**. The
tracker query's legacy-named `signature` parameter carries the nonce, not the
signature bytes. A transaction accepts at most 100 ETH proofs; nil proofs and
oversized signatures/nonces are rejected, and every attempted recovery is metered.

This voucher change is separate from EIP-712 transaction signing and transaction
session sign-in. Consumers must use the updated SDK before activating v35 flows.

## [3] Other behavior changes

- Deleting a collection fails while its aliases have bank supply or a registered
  GAMM pool, or its mint escrow retains bank coins. Existing wrapper/backing
  protections remain. An orphan alias denom with no registered router uses bank
  routing, so a prefix alone cannot strand an ordinary bank coin.
- Internal user transfers now enforce bank blocked-address and send-enabled
  policies. Explicit module-to-module APIs retain their protocol semantics. The
  pinned IBC transfer module already rejects blocked receiving module accounts.
- LP joins/swaps no longer silently enable a user's outgoing self-approval flag.
  Users whose outgoing approval does not permit the transfer must authorize it.
- Existing validators must set `minimum-gas-prices = "10ubadge"` in `app.toml`.
  The generated-config default changes only newly initialized homes.
- Wildcard flow changes to an aggregate channel key cause a one-time quota reset
  at upgrade. Known unsuffixed legacy windows remain stored but inactive; current
  timeframe windows initialize under their new keys. Review this reset operationally.
- Canonical balance collisions sum balances and retain canonical permissions;
  duplicate holders decrement the holder count. Boolean dynamic-store collisions
  use AND, preserving a revocation. Colliding voting scopes preserve votes but
  reset the quorum-delay timestamp; a pre-positioned eligible uppercase scope can
  therefore restart the delay. Isolated moved scopes retain their timestamp.
- All precompile message dispatch cases are metered. JSON nesting is capped at
  32, object fields at 100, arrays at 100 (address lists at 1,000), and strings at
  the existing 10,000-byte metadata limit; specialized tighter limits still apply.
- Deleting an in-use manager splitter is rejected. Delegate permissions and
  admin/address spellings are canonicalized at upgrade.

Collection purge remains atomic and gas-metered. Very large collections can
exceed the block limit and become undeletable; this fails closed without deleting
escrow protection. Solving that availability limitation requires a separately
designed staged purge/tombstone protocol. Do not remove metering or leave partially
deleted live collections as a shortcut.

## [4] Reproducing regressions

Run `GOTOOLCHAIN=auto go test ./... -tags=test -count=1` and
`scripts/upgrade/test/run.sh`. `TestV35MigrationsAgainstCapturedMainnetStores`
replays the captured nonempty prefixes and checks idempotence. Synthetic tests
exercise uppercase collisions, revocations, iterator lifetime, and authorization.
Batch migration walks retain at most 1,000 copied keys (about 1 MiB key budget,
plus one key), read values at visit time, and log progress. Total migration time
is still proportional to the state processed.

Four example contracts are compiled and exercised through the EVM keeper and
actual tokenization precompile. Regenerate the checked-in ABI/bytecode with
`python3 scripts/upgrade/test/compile_examples.py`; its compiler image is pinned
by digest. A source-hash regression rejects stale fixtures after Solidity edits.
The examples restore wrapper approvals after minting, debit actual holders,
use millisecond ownership times, and exercise credential/carbon issuance.


## [5] Release-candidate evidence (2026-09-06)

Consensus/application source tested: `815e88068cd5968c534863935557fcd0645ba35e`.
Chain and SDK CI passed. The complete local `--all` rehearsal passed build,
upgrade, rollback, EVM checks and four-validator app-hash agreement, including
transaction-bearing blocks before and after upgrade. Fresh scan at height
11,980,782 returned byte-identical audited module rows to height 11,979,578.
Mainnet's enabled precompile list includes 0x1001, 0x1002 and 0x1003.

On a local Apple M5 Pro Linux/arm64 container capped at **2 CPUs / 4 GiB**, two
sampled near-cap transactions executed successfully with the reviewed binary:

| Workload | Receipt gas | Finalize through commit (log timestamps) |
| --- | --- | --- |
| 28,000 valid-curve EVM signature recoveries | 95,410,717 | 1.088 s |
| 400 rejected tokenization calls with a 10KB string, caught by Solidity | 99,545,863 | 0.140 s |

Container peak memory during the probe was 331,165,696 bytes (~316 MiB).
The precompile probe explicitly enables the custom precompiles; a call to an
inactive precompile is not valid performance or functional evidence. Receipts
were successful while every nested tokenization call failed, confirming the
caller pays gas despite catching errors. These samples support retaining 100M
for the candidate; they are **not exhaustive worst-case coverage** of state-heavy
alias routing or a benchmark of the actual production validators.

Packed SDK 0.45.0 consumer checks: frontend 276 unit tests and typecheck; indexer
10 v35 tests and 32 EIP-712/auth tests. The HTTP smoke passed unsigned simulation,
foreign-session/tampering/concurrent replay rejection, cookie rotation, private
claim gating, and execution of the same EIP-712 signature on the reviewed binary.
SDK publication, registry-backed consumer lockfiles/CI, and the final browser
wallet check remain separate release steps.
