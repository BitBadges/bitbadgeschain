# ERC3643 template transfer policy

`initializeCollection` mints the initial supply to the admin using a temporary
Mint approval, then replaces that approval with `wrapper-transfers`. Zero-supply
initialization installs the same transfer policy without minting.

The ongoing approval applies only to token ID 1, excludes Mint, and permits only
this wrapper contract as the transfer initiator. It overrides user-level incoming
and outgoing approvals because the wrapper checks KYC and freeze status and binds
the debit address to `msg.sender`. Holders cannot bypass these checks by calling
the tokenization precompile directly. No allowance or `transferFrom` API is
provided by this template.

This policy is installed only for collections created by the template. Attaching
an existing collection with `setCollectionId` does not rewrite that collection's
approvals; its manager must configure an appropriate policy separately.

Regression coverage runs against the tracked Solidity bytecode:

```sh
GOTOOLCHAIN=go1.26.6 go test ./x/tokenization/precompile/test/integration -tags=test -mod=readonly -run TestERC3643TokenizationTestSuite -count=1
```

Regenerate the template ABI and bytecode with the existing Solidity 0.8.28
compiler, optimizer enabled, and `contracts` as the base path. Copy only the
template output into `contracts/test`; imported library artifacts are not used
by these tests.
