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
| Staking — validators, delegations, UBDs, redelegations | `app/upgrades/v35/staking.go` |
| Distribution — rewards, commission, historical ratios, community pool, starting info | `app/upgrades/v35/distribution.go` |
| Gov — in-flight proposal deposits and `TotalDeposit` | `app/upgrades/v35/gov.go` |
| GAMM — balancer `PoolAssets`, stableswap `PoolLiquidity` | `app/upgrades/v35/gamm.go` |
| Poolmanager — taker-fee key re-keying | `app/upgrades/v35/poolmanager.go` |
| Tokenization — `CoinTransfers` on collection and user approvals | `app/upgrades/v35/tokenization.go` |
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

### Verified on a live node

A fresh 18-decimal chain boots and produces blocks with no precisebank in the
graph, and the EVM and Cosmos views are now the *same number*:

```
cosmos  abadge : 999000000000000000000000
eth_getBalance : 0xd38be6051f27c2600000  = 999000000000000000000000   (1:1, = 999,000 BADGE)
```

That 1:1 is the point of the migration — at 9 decimals those two numbers differ
by 10^9 and only agree because precisebank sits between them.

Suite is green at 76 packages on the 18-decimal branch.

## Three hazards this surfaced

**1. `PowerReduction` has to move with the decimals, and nobody would notice.**

Confirmed the hard way: with `PowerReduction` left at 10^6 the chain refuses to
start — *"validator set is empty after InitGenesis, please ensure at least one
validator is initialized with a delegation greater than or equal to the
DefaultPowerReduction"*. A genesis validator staking what used to be a healthy
amount now falls below a single unit of consensus power.

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

## What is deliberately not converted

**IBC denom traces.** `ubadge` already sent to another chain came back as
`ibc/<hash-of-ubadge>`. That hash is derived from the denom string on the
*counterparty*, so nothing this chain does re-derives it, and the counterparty
keeps the old trace regardless. Escrowed balances convert with bank; the
outstanding vouchers do not and cannot. This needs a per-channel decision
(drain before the upgrade, or accept two denominations for the same asset), and
it is the single item that genuinely cannot be solved chain-side.

**LP share denoms** (`gamm/pool/N`) and **pool weights**, which are their own
denom and dimensionless respectively. See the note in `gamm.go` — scaling either
would reprice every LP position.

**Slash-event fractions** in x/distribution, which are dimensionless.

## Still to do before this could ship

- **Rehearsal against a mainnet state snapshot**, asserting value conserved
  per module rather than only that the upgrade does not panic. The unit tests
  cover the arithmetic; they cannot cover the shape of real state.
- **Client repos in lockstep** — ~1370 references across 261 files: bitbadgesjs
  SDK (107 files / 452 refs), indexer (46 / 557), frontend (72 / 268), docs
  (36 / 94). The indexer is the sharp edge; it stores historical `ubadge`
  amounts, so it needs a backfill decision rather than a constant change.
- **An ADR** recording the decision, per the original ticket.

## Suggested sequencing

The chain migration is the easy half and it is already large. The honest split:

1. Land v34 first. It is independent and ready.
2. ~~Build the poolmanager re-keying and the distribution/gov/gamm/tokenization
   walks~~ — done, with an end-to-end value-conservation test over the real
   handler.
3. Rehearse against a mainnet state snapshot, asserting total value conserved
   per module — not just that the upgrade does not panic.
4. Coordinate the client repos on one release.

The forcing function remains block-STM (BB-8): `EnableVirtualFeeCollection()`
panics below 18 decimals, so parallel execution's fee path is gated behind this
migration. Until that is wanted, 9 decimals keeps working — precisebank is
restored on v34 and upstream still ships it.
