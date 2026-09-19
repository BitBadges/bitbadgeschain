# Local token lifecycle execution

The version-1 runner executes tokenization messages and bank sends through real
app message handlers with isolated in-memory state, fresh temporary homes,
fixture actors, and an explicit millisecond clock. It cannot broadcast.

```sh
GOTOOLCHAIN=auto go build -tags=test \
  -ldflags "-X github.com/bitbadges/bitbadgeschain/lifecycle.ChainCommit=$(git rev-parse HEAD)" \
  -o /tmp/bitbadges-lifecycle ./cmd/bitbadges-lifecycle
/tmp/bitbadges-lifecycle --version
/tmp/bitbadges-lifecycle --schema
/tmp/bitbadges-lifecycle < lifecycle/testdata/scenarios/subscription-paid.json
GOTOOLCHAIN=auto go test -tags=test ./lifecycle ./cmd/bitbadges-lifecycle -count=1
```

`--schema` emits the input JSON Schema. The Go decoder also enforces valid
addresses, supported assertions, explicit outcomes, unique actor/step IDs,
positive time advances, and one message or time advance per step. The input is
limited to 8 MiB, 100 actors, and 1,000 steps. A step has a 50 million gas limit.
Each step needs a non-null expectation; an empty scenario never passes.

A message envelope is `{ "typeUrl": "/tokenization.MsgCreateCollection", "value":
{...} }`. The actor must match the signer extracted by the app's codec. No mint
approval is inserted, no artifact is repaired, no identifier is substituted.
Unknown protobuf fields are rejected. Identity binding and metadata resolution
must happen before execution. Each successful message commits its cache; a
failed message discards all message writes, including earlier transfers in a
batch. Expected rejection is a passing step only if the state assertions also
pass. `errorContains` optionally pins the rejection reason.

Time steps advance only the context clock. No expiry sweep is invented: balance
ownership ranges remain stored, and `balanceAt` checks the specified point or
`"now"`. Adjacent intervals compare by their amounts everywhere, not by array
length or serialization. Endpoints are inclusive. Overlaps between balances sum;
overlaps inside a balance's coordinate ranges are rejected as ambiguous.

## Assertions

All assertions carry `kind` and `expected`. Numeric values are decimal strings.

| Kind | Selectors | Expected value |
| --- | --- | --- |
| `coins` | address, denom | Exact bank balance |
| `coinDelta` | address, denom | Signed change during this step |
| `balances` | collectionId, address | Entire balance surface, including token/time ranges |
| `balanceAt` | collectionId, address, tokenId, ownershipTime | Amount at point; time may be `now` |
| `approvals` | collectionId | Complete stored collection approvals |
| `permissions` | collectionId | Complete stored collection permissions |
| `manager` | collectionId | Manager address |
| `tracker` | collectionId, approvalId, amountTrackerId, level, trackerType, addressForApproval, address | Complete stored approval tracker |

Approvals and permissions use the chain's stored JSON representation. The chain
materializes defaults, including `mustPrioritize` for payment approvals. Fixture
expectations explicitly include those defaults; the runner does not normalize
away differences. Missing collections/trackers are errors rather than an empty
matching state. Use an explicit empty balance array to assert zero ownership.

## Evidence and limits

Output contains the exact-input SHA-256, chain commit, per-step outcome, event
records, observed and expected state, time, and exclusions. Exit 0 means all
assertions/outcomes passed; 1 means a failed expectation; 2 means malformed or
unsupported input or an execution failure. Unknown build revision stays
`unknown`; never treat it as verified provenance. Pin the binary and retain its
hash alongside evidence. Repeated runs use fresh state and deterministic actor
observations; app-internal test genesis accounts are outside the scenario scope.

`execution: module` excludes signatures, transaction fees, sequence checks,
block hooks, IBC, and external services. Protocol-level payment fees **are**
executed: with current default params the 7 BADGE invoice sends 7 to the merchant
and charges the payer another 0.007 BADGE. This is separate from transaction gas.
The full-transaction conformance test executes real signed bank sends and a
subscription creation through BaseApp, verifies module state parity, emitted
events, fees, sequence increments, replay rejection, and tampered-signature
rejection. A module rejection after successful ante still charges the fee and
increments sequence; the conformance test pins that distinction from module-only
rollback. Transactions whose signers cannot fund their fees are rejected before that state change. It is local Cosmos direct signing at genesis height, not a claim of
coverage for every wallet signature format or consensus/block processing.

## Builder fixtures

`testdata/artifacts` preserves SDK builder outputs. The five preset inputs are
in `builder-inputs.json`; NFT/fungible examples use the SDK session APIs and the
published skill configuration. The generator removes only the session `_meta`
sidecar before producing a chain envelope. Scenario assembly binds initially
empty creator/manager fields to the fixture signer. All approval, payment,
permission, range, and invariant fields execute unchanged.

```sh
bun lifecycle/testdata/generate-artifacts.ts /path/to/sdk/package
UPDATE_LIFECYCLE_FIXTURES=1 GOTOOLCHAIN=auto go test -tags=test ./lifecycle -count=1
```

The generator records SDK revision, dirty status, complete dist hash, and
individual artifact hashes in `provenance.json`. Regeneration requires an
explicitly built SDK, including the session normalization fix for
`updateDefaultBalances` and `updateInvariants`. Review fixture changes as spec
changes. Tests read stored artifacts and do not need Bun or a sibling checkout.
They also replay checked-in JSON scenarios, so examples are executable contracts.

Subscription renewal here is a new paid, contiguous 30-day ownership period,
not calendar-month arithmetic or an automatic recurring charge. The base
subscription allows timestamp overrides; callers must choose and verify desired
renewal dates. Its initial non-transferability is not immutable: manager powers
may allow adding non-mint approvals. Invoice settlement is direct payment, not
escrow or a refund guarantee. Backed redemption asserts on-chain reserves and
balances; it does not assert off-chain backing. Burning in NFT/fungible examples
sends to the canonical sink and does not restore lifetime mint allowance.

Credit fixtures additionally cover scaled purchase ratios, holder-transfer/burn
rejection for cumulative credits, spendable pack caps, holder-only consumption,
and expiry of purchase/consumption permissions. Spendable balances remain stored
after expiry: the test does not confuse a transfer deadline with ownership expiry.
Successful consumption proves token transfer to the sink, not service delivery.
