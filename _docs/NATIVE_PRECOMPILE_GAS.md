---
title: Native precompile gas settlement
date: 2026-09-05
---

The tokenization, GAMM, and sendmanager `Run` entry points use
`pkg/evmcompat.RunNativeAction`. It delegates snapshots, balance reconciliation,
successful-action billing, error conversion, and SDK out-of-gas recovery to the
pinned Cosmos EVM wrapper. On an ordinary action error only, it charges the SDK
gas consumed by that action before returning the error. An EVM revert preserves
remaining gas, so rolling back state must not also discard this charge.

The shim deliberately does not bill successful actions or recover panics itself.
The pinned wrapper already handles those paths; billing them again would double
charge. Its pre-action balance hook only records an event offset and does not
consume SDK gas. The shim measures from callback entry, excluding previously
consumed context gas. Static `RequiredGas` charges remain separate and unchanged.

This is a chain-owned integration change in v35, not a dependency update or a
module-cache patch. It applies to the three custom precompiles above, not direct
calls to the dependency's base wrapper or other upstream precompiles. Revisit the
shim when upgrading Cosmos EVM: if its wrapper begins billing ordinary action
errors, remove the extra settlement here and retain the regression assertions.

`app/native_action_gas_test.go` checks exact static-plus-dynamic charges, successful
actions, ordinary errors, SDK out-of-gas, state rollback, a real EVM parent catching
two child reverts, and failing calls through all three custom entry points.

```sh
go test -mod=readonly ./app -tags=test -run 'TestNativeAction|TestCustomPrecompileFailures' -count=1
```

Full-suite and upgrade-rehearsal checks remain necessary for release. No circuit
breaker state is changed by this integration; EVM activation is a separate
operator decision.
