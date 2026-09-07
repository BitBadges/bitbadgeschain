# EIP-712 GAMM and sendmanager coverage

GAMM already registers its Amino types through AppModuleBasic in v35. The ten
routed GAMM messages (nine core messages plus balancer pool creation) work with
EIP-712. Stableswap message services remain disabled; these changes do not enable
a new pool model or bypass any message authorization.

Sendmanager previously had no LegacyAmino registrations. The app now registers
MsgSendWithAliasRouting and MsgUpdateParams using their exact proto amino.name
values. The SDK companion adds their generated proto classes, Amino converters,
and decoding registry entries. Parameter updates still require governance authority.

`app/testdata/eip712-module-messages.json` is copied from the SDK's
`src/eip712/module-messages.fixture.json`. Its test creates frontend transaction
payloads, signs all twelve messages, checks indexer-side signature verification,
and rejects a changed memo. The Go test independently hashes each sign document,
decodes the actual SDK-signed transaction bytes, accepts the original signature,
and rejects a changed memo through the app's ante handler. This verifies signing
and decoding, not successful execution of every pool/IBC operation against live state.

Regenerate fixtures in the SDK with:

```bash
UPDATE_EIP712_MODULE_FIXTURES=1 bun run test --runInBand src/eip712/module-messages.spec.ts
```

Copy the fixture to the chain testdata path above, then run:

```bash
GOTOOLCHAIN=auto go test ./app -tags=test -run TestEIP712 -count=1
```

Balancer fee fields in proto messages contain scaled integer strings (for example
3000000000000000), as supplied by the frontend. Decimal text such as 0.003 is not
valid protobuf wire input for the Go custom decimal type.

v35 has not rolled out. After human merge, retag v35 at the merged chain commit
and regenerate its release binaries and checksums before rollout. Merge alone
does not activate this registration. Publish/adopt the companion SDK in
coordination with the updated v35 binary.
