# v35 — 9 → 18 decimal migration (exploration)

Branch `exp/v35-18-decimals`, stacked on `exp/v34-sdk-0.54.4`. **Exploratory, not
shippable.** The chain-side core works and is tested; the sections under "Not
done" are real gaps, not polish.

## Shape of the change

`ubadge` (9 decimals) → `abadge` (18 decimals), ×10⁹. Nobody's holdings change
value, only their representation.

Base denom becomes `abadge` rather than keeping `ubadge` and moving its exponent,
for two reasons. x/vm's `validateCoinInfo` rejects any 18-decimal config where
`Denom != ExtendedDenom`, so the base must equal the extended denom the EVM
already uses. And "micro" naming a 10⁻¹⁸ unit would be actively misleading —
`abadge` (atto) is what it is.

x/precisebank goes away. With base == extended there is no precision gap left to
bridge, which is the whole point: precisebank now lives in `contrib` as a
compatibility shim for chains upstream considers unsupported.

## Done and tested

| Piece | File |
|---|---|
| Bank balance + supply redenomination | `app/upgrades/v35/redenominate.go` |
| Denom metadata (drives x/vm's decimals) | same |
| Staking rescale — validators, delegations, UBDs, redelegations | `app/upgrades/v35/staking.go` |
| Module params — staking, mint, gov, evm, tokenization `AllowedDenoms` | `app/upgrades/v35/params.go` |
| EVM config at 18 decimals, precisebank unwired | `app/params/constants.go`, `app/evm*.go`, `app/app.go` |
| `PowerReduction` 10⁶ → 10¹⁵ | `app/params/constants.go` |

Tests in `app/upgrades_v35_redenominate_test.go` cover exact ×10⁹ conversion
including the 1-unit dust case, legacy denom fully retired, balances summing to
supply (the assertion that catches value being invented or destroyed), the
staging module account's own balance surviving, foreign/IBC denoms untouched, and
a refusal when supply disagrees with the sum of balances.

Balances are rewritten with `UncheckedSetBalance` rather than transfers: this is
a representation change, not a movement of value, and routing it through
`SendCoins` would trip blocked-address checks on module accounts and fire
transfer hooks for something that is not a transfer. Supply is moved once, via a
burn and a mint through a module account whose own balance is restored exactly.

## Three hazards this surfaced

**1. `PowerReduction` has to move with the decimals, and nobody would notice.**

The chain never overrides `sdk.DefaultPowerReduction`, so it is 10⁶. Consensus
power is `tokens / PowerReduction`. Multiplying every balance by 10⁹ without
touching it multiplies every validator's power by 10⁹. With supply on the order
of 10¹⁵ ubadge, post-migration power lands around 10¹⁸ — against CometBFT's
`MaxTotalVotingPower` of ~8.2×10¹⁸, with no headroom for supply growth.

Set to 10¹⁵ here, which leaves every validator's power numerically identical
across the upgrade. This is a consensus-affecting global, so it takes effect for
whichever binary is running — correct under cosmovisor, which replays each height
range with the binary that produced it, but it does mean a node cannot sync from
genesis with only the v35 binary.

**2. The denom's sort position changes, and sort order is load-bearing.**

`abadge` sorts *first* where `ubadge` sorted *last*. `sdk.Coins` is denom-sorted,
so every stored coin list reorders. This is not cosmetic:
`x/poolmanager/types.FormatDenomTradePairKey` embeds the denom strings directly
in the store key, so **every poolmanager trade-pair record keyed on `ubadge`
would be orphaned by the rename** and has to be re-keyed in the migration. That
is not implemented (see below). It showed up first as a test failure, which is
the only reason it was noticed at all.

**3. The test suite hides denom coupling.**

198 hardcoded `"ubadge"` literals across 23 test files, several of them in
position-sensitive fixtures — sorted coin strings and ordered denom slices that
pass only because `ubadge` happened to sort last. Mechanically fixed here, but
the count is the useful signal: nothing in the codebase forces test fixtures
through `appparams.BaseCoinUnit`, so the denom is coupled in far more places than
a grep for the constant would suggest.

## Not done — the real remaining work

**Chain**

- **poolmanager trade-pair re-keying.** Hazard 2. Records keyed
  `<prefix>|ubadge|X` must be rewritten to `<prefix>|abadge|X`, both directions,
  including taker-fee entries. Without this, existing routes silently disappear.
- **x/distribution.** Every `DecCoins` in the module scales: outstanding rewards,
  accumulated commission, current and historical rewards, `DelegatorStartingInfo.Stake`,
  and `FeePool.CommunityPool`. Skipping it mis-scales every pending reward.
- **x/gov deposits.** Params are handled; the `Deposit` records and
  `Proposal.TotalDeposit` on in-flight proposals are not.
- **x/gamm pool records.** Pool reserves are bank balances and convert for free,
  but the pool objects store their own `PoolAssets` and `TotalShares`. Scaling
  assets without shares reprices every LP position.
- **x/tokenization stored amounts.** `CoinTransfers` inside approval criteria
  carry `ubadge` amounts; `CosmosCoinWrapperPaths` carry denoms. Both are
  user-authored and persisted, so they need a walk of every collection.
- **IBC denom traces.** Escrow balances convert with bank, but `ubadge` sent to
  another chain returns as `ibc/<hash-of-ubadge>`. Those vouchers do not
  re-derive, and counterparty chains keep the old trace.

**Off-chain, in lockstep**

~1370 references across 261 files: bitbadgesjs SDK (107 files / 452 refs),
indexer (46 / 557), frontend (72 / 268), docs (36 / 94). The indexer is the
sharpest edge — it stores historical `ubadge` amounts, so it needs a backfill
decision, not just a constant change.

## Suggested sequencing

The chain migration is the easy half and it is already large. The honest split:

1. Land v34 first. It is independent and ready.
2. Build the poolmanager re-keying and the distribution/gov/gamm/tokenization
   walks, each with its own before/after invariant test.
3. Rehearse against a mainnet state snapshot, asserting total value conserved
   per module — not just that the upgrade does not panic.
4. Coordinate the client repos on one release.

The forcing function remains block-STM (BB-8): `EnableVirtualFeeCollection()`
panics below 18 decimals, so parallel execution's fee path is gated behind this
migration. Until that is wanted, 9 decimals keeps working — precisebank is
restored on v34 and upstream still ships it.
