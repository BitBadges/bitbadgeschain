# Upgrade rehearsal (`scripts/upgrade`)

Tracked, parametrised, non-interactive tooling for every chain upgrade. It
replaces the legacy `scripts-dev/` scripts (`handle_upgrade_logic.sh`,
`upgrade-helper.sh`, `propose-version.sh`, `upgrade-rehearsal/`), which were
hardcoded to one version pair, needed two terminals and `ignite`, and used
Linux-only paths. Everything here runs unattended on macOS (Apple silicon or
Intel) or Linux with Docker; the exit code is meaningful (0 = every check in
every stage passed).

```
scripts/upgrade/
  new-version.sh       start vNN: scaffold + wire the handler, bump Makefile VERSION
  rehearse.sh          Docker rehearsal between two git refs (host driver)
  propose.sh           submit + vote a MsgSoftwareUpgrade (container or real node)
  checks/<vNN>.sh      per-version post-upgrade assertions, run by rehearse.sh
  container/           scripts that run inside the image (build, upgrade, ...)
  cosmovisor/          wrapper module that builds cosmovisor on the chain's toolchain
  lib/, templates/     helpers and scaffolding templates
  test/run.sh          self-tests + shellcheck  (make upgrade-tooling-test)
```

## Release flow in three commands

```sh
make new-version V=v36                       # 1. start the version
#    ... implement app/upgrades/v36/upgrades.go, commit ...
make upgrade-rehearsal FROM=v35 TO=HEAD      # 2. rehearse v35 -> HEAD
scripts/upgrade/propose.sh --name v36 ...    # 3. propose on testnet/mainnet
```

## 1. `new-version.sh vNN`

```
scripts/upgrade/new-version.sh v36 [--dry-run] [--snapshot-proto [--proto-ref REF] [--no-proto-gen]]
```

Idempotent; running it twice is a no-op. It:

- bumps `VERSION := vNN` in the Makefile;
- scaffolds `app/upgrades/vNN/upgrades.go` from `templates/upgrades.go.tmpl`
  **only if absent**;
- adds the import, the `SetUpgradeHandler` registration and a
  `StoreUpgrades` switch case to `app/upgrades.go` **only if vNN is not
  already registered there**. A hand-wired version is left alone entirely,
  because whoever wired it already decided whether it needs a store case.

Unlike `handle_upgrade_logic.sh` it never renames or deletes
`app/upgrades/v<N-1>`; every handler stays in history.

`--snapshot-proto` ports `sync-old-proto.sh`: it copies
`proto/tokenization/*.proto` from git ref `v<N-1>` (or `--proto-ref`) into
`proto/tokenization/v<N-1>/`, rewrites `package`, imports and `go_package`
to the versioned names, drops `query.proto`/`legacytx.proto`, writes
`x/tokenization/types/v<N-1>/math.go`, repoints the `oldtypes` import in
`x/tokenization/keeper/migrations.go`, and runs `ignite generate proto-go -y`
when `ignite` is installed. **Only use it when the tokenization proto types
changed since the previous version** and the new handler needs to decode the
old encoding for a state migration. Bumping `ConsensusVersion` in
`x/tokenization/module/module.go` and registering the migration stay manual:
they depend on whether there is a state migration at all.

## 2. `rehearse.sh`

```
scripts/upgrade/rehearse.sh --from v34 --to HEAD [--name v35]
    [--rollback] [--evmcheck] [--multivalidator] [--all]
    [--workdir DIR] [--skip-build] [--dry-run]
make upgrade-rehearsal FROM=v34 TO=HEAD [NAME=v35] [REHEARSAL_FLAGS="--all"]
```

Both binaries and cosmovisor are built **inside** the image from `git
archive` of each ref, never from a working tree (uncommitted changes are not
rehearsed; commit first). The upgrade name defaults to the highest
`app/upgrades/vNN` directory on the `--to` ref. `--dry-run` resolves refs,
toolchains and the name and prints the plan without touching Docker.
`--skip-build` reuses the binaries already in the workdir. Logs and per-stage
transcripts land in `$TMPDIR/bb-upgrade-rehearsal/<from>-to-<to>/logs`.

### Stages and what they assert

**build** (always) — FROM, TO and cosmovisor compile on the toolchain each
`go.mod` requires.

**upgrade** (always) — single-validator cosmovisor rehearsal. After
cosmovisor starts nothing touches the node; the only inputs are the proposal
and the vote, via `propose.sh`. Asserts:

1. the chain advanced past the upgrade height;
2. the FROM binary halted with `UPGRADE "<name>" NEEDED`;
3. cosmovisor ran the `pre-upgrade` hook;
4. the TO binary logged `applying upgrade "<name>"`;
5. `cosmovisor/current` resolves into `upgrades/<name>/`;
6. blocks are still being produced afterwards (2 + 6 together prove the
   TO binary, not FROM, produced them);
7. `config.toml`/`app.toml`: any difference from the pre-start snapshot is
   printed (only the pre-upgrade hook could have written it), and a rewritten
   file must have a matching `<file>.bak`. This is deliberately not a
   hardcoded "mempool flood -> app" check: each version's hook may or may not
   rewrite config, and the diff is the report;
8. `checks/<name>.sh`, if present, runs with `BIN` (TO), `HOME_DIR`,
   `UPGRADE_HEIGHT`, `LOG` and `API_URL` in the environment and its exit code is the
   number of failed checks. `checks/v35.sh` asserts
   `consensus block.max_gas == 100000000` and
   `feemarket min_gas_price == 10`, then runs the EVM value-transfer
   assertions from `lib/evm_transfer.sh` (below). Add one file per version.

**EVM value transfers** (`lib/evm_transfer.sh`, run by `checks/v35.sh` in
the upgrade stage and by the evmcheck stage) — creates two `eth_secp256k1`
keys, funds the sender with 2 BADGE by bank send at `--gas-prices 10ubadge`,
exports its hex key and signs two legacy EIP-155 transfers with the tracked
`scripts/upgrade/evmtx` helper (built in the container from the TO tree; gas
price = max(`eth_gasPrice`, 10 gwei)). Asserts, for a 0.5 BADGE transfer:
receipt status `0x1`, recipient `eth_getBalance == value`, recipient bank
`ubadge` up by value/1e9, sender bank `ubadge` down by (value+fee)/1e9 with
the fee a whole number of ubadge; and for a second 1e9+5 wei transfer:
status `0x1`, recipient EVM balance exactly v1+v2 and bank balance holding
the whole part. After both, `/cosmos/bank/v1beta1/denom_owners/abadge` has
zero owners and `supply/by_denom?denom=abadge` is 0 — the 18-decimal denom
never reaches x/bank.

**rollback** (`--rollback`) — runs the TO `pre-upgrade` hook the way
cosmovisor does (no `--home`, only `DAEMON_HOME`) on a stock FROM home and
asserts it exits 0 and is idempotent. If it rewrote config, `<file>.bak`
must match the pristine file and the FROM binary must produce blocks again
once the config is restored; whether FROM also starts against the rewritten
config is *reported* (that is what tells the release notes whether the
rollback order matters). If it rewrote nothing, the rollback is binary-only
and FROM must still start. Note this covers configuration only: state
written by the TO binary after the upgrade height cannot be rolled back.

**evmcheck** (`--evmcheck`) — fresh chain from a plain `init` on the TO
binary with JSON-RPC enabled: `eth_getBalance` must equal the bank `ubadge`
balance x 10^9, and a bank send must succeed. Run it every time: the upgrade
stage starts from a config the FROM binary wrote, so it cannot see whether a
config `init` produces is startable (the v34 defect this caught).

**multivalidator** (`--multivalidator`) — four independent validators under
cosmovisor. Asserts all four agree on the app hash before the upgrade, in a
block that provably contains a transaction, at the upgrade height, after it,
and in a post-upgrade transaction block; all four apply the handler and run
the hook. This is the only stage that can see nondeterminism: with one
validator, whatever it computes is consensus.

### Genesis and config deviations

Genesis: only governance timing: `voting_period` 15s (20s multivalidator),
`expedited_voting_period` 10s, `max_deposit_period` 30s, deposits 1/2 of
whatever denom genesis already uses. Denoms are read from genesis
(`staking.params.bond_denom`, `gov.params.min_deposit[0].denom`), not
assumed.

Config: the upgrade stage enables the REST API server (`[api] enable = true`
in `app.toml`, before the snapshot the hook diff is taken against) so
`checks/<name>.sh` can read module params that have no CLI query, and the
EVM JSON-RPC server (`[json-rpc] enable = true`, 127.0.0.1:8545) for the EVM
transfer assertions. The chain
binary registers no `query feemarket`/`query vm` commands, so `checks/v35.sh`
reads `$API_URL/cosmos/evm/feemarket/v1/params`.

### Toolchain and module notes

- **One Go per ref.** `rehearse.sh` reads the `go`/`toolchain` line of each
  ref's `go.mod`, installs every distinct version into the image
  (`/usr/local/go-<v>`) and builds each ref with its own. The image tag
  carries the set, so a new version pair rebuilds the image once.
- **cosmovisor on Go 1.26.** cosmovisor v1.7.1's own `go.mod` resolves
  `bytedance/sonic` to a release that does not link on Go 1.26 (the v33->v34
  rehearsal had to keep a go1.25.8 toolchain around just for it). It is now
  built through `scripts/upgrade/cosmovisor/`, a wrapper module that raises
  sonic to v1.15.1, so it builds on the same toolchain as the chain. To bump
  cosmovisor: edit `require cosmossdk.io/tools/cosmovisor` there and run
  `go mod tidy` in that directory.
- **Offline module resolution.** The image sets
  `GOPROXY=file:///hostmod/cache/download` with the host `GOMODCACHE` mounted
  read-only. Everything a ref needs must already be in the host cache
  (`go mod download` on the host first), including any private module the
  ref's `go.mod` replaces `cosmos/evm` with — no token ever enters the
  container. Do not set `GOPRIVATE`/`GONOPROXY` when running: they mean
  "skip the proxy" and send Go to the network.
- **CGO.** go-ethereum's secp256k1 has no pure-Go fallback, so binaries are
  Linux builds made in the container, matching what validators run.
- **arm64.** Runs natively on Apple silicon. Release binaries are amd64, but
  the upgrade logic under test is architecture-independent.
- **libwasmvm** switching from `upgrade-helper.sh` is not reproduced: the
  chain no longer depends on wasmvm.

## 3. `propose.sh`

```
scripts/upgrade/propose.sh --name v35 --home ~/.bitbadgeschain --from <key>
    [--deposit 10000000ubadge] [--height +N | --height H] [--expedited]
    [--chain-id X] [--node tcp://...] [--bin bitbadgeschaind]
    [--keyring-backend os] [--fees 0ubadge | --gas-prices 10ubadge]
    [--info '<binaries json>'] [--voters k1,k2] [--no-vote] [--dry-run]
```

Non-interactive replacement for `propose-version.sh`: no `.env` reads, no
`ustake` (the deposit denom defaults to `ubadge` and `ustake` is rejected),
no bootstrapped-collection transactions. It queries the gov module account
and current height from the node, writes `<home>/proposal-<name>.json`,
submits, reads the proposal id back from the tx events, votes yes from
`--from` (and `--voters`), and prints `PROPOSAL_ID=` / `UPGRADE_HEIGHT=`.
The rehearsal container uses exactly this script. Mainnet example:

```sh
scripts/upgrade/propose.sh --name v35 --home ~/.bitbadgeschain --from faucet \
  --height 12345678 --expedited --deposit 20000000000ubadge --gas-prices 10ubadge \
  --info '{"binaries":{"linux/amd64":"https://github.com/BitBadges/bitbadgeschain/releases/download/v35/bitbadgeschain-linux-amd64","linux/arm64":"https://github.com/BitBadges/bitbadgeschain/releases/download/v35/bitbadgeschain-linux-arm64"}}'
```

## What the v34 -> v35 rehearsal established

Recorded here because each item shapes how the next rehearsal is read.

- **Hands-off.** The v35 handler applies under cosmovisor with no operator
  action; four validators agree on every app hash through the upgrade.
- **Fee floor semantics.** `feemarket.min_gas_price = 10` is enforced on
  Cosmos transactions as exactly 10 ubadge per gas, in `ubadge` only: a
  200k-gas transaction needs `--gas-prices 10ubadge` (2000000ubadge), 9.99
  is rejected, and a fee paid in `stake` is rejected outright. The EIP-1559
  `base_fee` (18-decimal EVM units, starts at 1e9 and decays on empty
  blocks) does not gate Cosmos transactions. After v35, `propose.sh` and any
  other CLI use on mainnet must pass `--gas-prices 10ubadge` (or higher).
- **Fresh-chain gotcha.** A genesis with `min_gas_price > 0` fails at
  `InitGenesis` because the gentx carries no fee ("minimum global fee for
  this tx is: 2000000ubadge"). Local/test chains that want the floor must
  either raise it after genesis or create the gentx with `--fees`.
- **EVM value transfers work through precisebank.** Legacy EIP-155
  transfers of 0.5 BADGE and of 1e9+5 wei both land with status `0x1`;
  `eth_getBalance`, bank `ubadge` and the fee reconcile exactly, and `abadge`
  never appears in x/bank (zero denom owners, zero supply).
- **Young-chain gas price.** The feemarket `base_fee` starts at its genesis
  default of 1e9 and is applied as ubadge per gas, so `eth_gasPrice` on a
  fresh chain is ~6e17 wei/gas and a 21000-gas transfer costs ~12,000 BADGE;
  it decays 12.5% per empty block (about 140 blocks to reach the 10 ubadge
  floor). Mainnet's base fee decayed long ago and sits at the floor, so this
  only bites fresh testnets/local chains. The rehearsal funds its EVM sender
  from the live `eth_gasPrice` for that reason.
- **Hook is a no-op for v35.** `pre-upgrade` finds `mempool.type` already
  `"app"` on a v34 home, rewrites nothing, and rollback is binary-only.
- **No feemarket CLI.** The chain binary registers no `query feemarket`
  command; read `/cosmos/evm/feemarket/v1/params` over REST.

## Self-tests

`make upgrade-tooling-test` runs `scripts/upgrade/test/*.sh` (new-version
fixtures including the proto snapshot, `propose.sh --dry-run`, `rehearse.sh
--dry-run`, a fake-docker test that a failing stage fails the run, and
`go vet`/`go build` plus flag validation of the `evmtx` helper) and
shellcheck (local binary, else `koalaman/shellcheck` via Docker). No chain
binary is built by the tests.
