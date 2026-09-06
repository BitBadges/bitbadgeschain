---
title: v35 voting identity migration
last-verified: 2026-09-05
---

This policy accompanies [chain PR #119](https://github.com/BitBadges/bitbadgeschain/pull/119).
For coordinated application releases, use the [v35 downstream rollout runbook](https://github.com/BitBadges/bitbadgesjs/blob/fix/v35-gas-and-proof-query/docs/runbooks/v35-downstream-rollout.md).

## Existing voters and votes

Voting configurations retain each voter slot's original address spelling and weight.
Existing votes retain their percentage and timestamp. Uppercase and lowercase slots
belonging to the same account are not merged: quorum still rounds each slot's
weighted contribution separately.

Messages must still use canonical lowercase addresses. `MsgCastVote` authenticates
the account identity against the configured slots, then overwrites every slot for
that account with the requested percentage. The account can therefore change or
revoke legacy votes without submitting a noncanonical message. Reset-after-execution
continues deleting the configured slots.

## User-level approver scopes

The migration normalizes the approver-address component of incoming/outgoing voting
keys, not the voter-address component. A vote already stored at the canonical key
wins a collision, even if the moved record is newer. Noncolliding records move
without changing their encoded vote data.

When voting scopes combine, their quorum timestamp is cleared to zero. The same
reset applies when uppercase and canonical user-balance namespaces collide: the
canonical approval configuration survives, so a delay established under the
discarded configuration must not carry over. Collision detection occurs before
balance namespaces are merged.

**A configured delay restarts only after a fresh authorized vote establishes
quorum.** Existing votes remain stored; resetting the timestamp does not erase them.
Clients should prompt an eligible voter to recast when the chain reports
`vote again to initialize`, then wait the full configured delay.

An isolated scope without either kind of collision retains its timestamp. Repeating
the migration does not reset an already re-established delay or change migrated
records.

## Verification

```sh
GOTOOLCHAIN=go1.26.6 go test -mod=readonly ./x/tokenization/keeper -tags=test \
  -run 'TestTokenizationKeeperTestSuite/(TestV35(LegacyVoter|DualSpelling|VotingApprover|VotingBalance)|TestMigrateV35PreservesVoter)' -count=1
```

The regressions cover canonical message validation, vote changes/revocations,
separate-weight rounding, canonical-record precedence, balance-namespace collisions,
delay reinitialization, and idempotence. These fixture-based checks do not establish
how many live-chain configurations need a recast; verify that during the upgrade
rehearsal.
