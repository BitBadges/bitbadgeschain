package v34

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	evmethsecp256k1 "github.com/cosmos/evm/crypto/ethsecp256k1"

	legacyethsecp256k1 "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"
)

// LegacyEthereumPubKeyTypeURL is the proto type URL of the pre-cosmos/evm
// EVM public key this chain shipped before adopting cosmos/evm.
const LegacyEthereumPubKeyTypeURL = "/ethereum.PubKey"

// PubKeyMigrationResult reports what MigrateLegacyEthereumPubKeys did, so the
// upgrade handler can log it and a test can assert on it.
type PubKeyMigrationResult struct {
	Scanned   int
	Converted int
	Cleared   int
}

// MigrateLegacyEthereumPubKeys rewrites accounts still holding this chain's
// original, pre-cosmos/evm EVM public key.
//
// Before adopting cosmos/evm, EVM accounts were stored with a home-grown
// `/ethereum.PubKey` (chain-handlers/ethereum/crypto/ethsecp256k1). cosmos/evm
// brought its own `/cosmos.evm.crypto.v1.ethsecp256k1.PubKey`, and new accounts
// have used it ever since — but the accounts created under the old scheme were
// never migrated. encoding/codec/codec.go has carried a note about exactly this
// for as long as it has existed:
//
//	some accounts still have etheruem.PubKey and other dependent types.
//	To fully remove this, we need to handle migrations of these accounts.
//
// The consequence is that those accounts cannot transact. cosmos/evm's ante
// handler picks the signature-verification path with a Go type switch:
//
//	case *ethsecp256k1.PubKey:   // cosmos/evm's type, not ours
//
// Our legacy key is a *different Go type that happens to share the name*, so it
// falls through to the default branch and the tx is rejected with
//
//	unrecognized/unsupported public key type: *ethsecp256k1.PubKey: invalid pubkey
//
// which reads as though the type were supported — the %T verb prints only the
// package-qualified name, and both are "ethsecp256k1.PubKey". That collision is
// why this went unexplained for so long (see BB-12).
//
// Both types store the same thing: a 33-byte compressed secp256k1 key. Both
// derive the address identically (DecompressPubkey -> PubkeyToAddress), so the
// conversion is lossless and the account keeps its address, number and
// sequence. Converting in place is preferred over simply clearing the key,
// because a converted account works immediately — clearing would leave it
// unusable until its owner happened to sign something, and leaves a null
// pub_key that explorers and indexers read as "never seen".
//
// Clearing is still the fallback for a key that cannot be decompressed or whose
// address does not match the account it is stored under. Such an account is
// already unusable; leaving the legacy key in place would keep it that way,
// whereas a nil pub_key is repopulated correctly by the ante handler's
// SetPubKeyDecorator on the owner's next signed transaction.
func MigrateLegacyEthereumPubKeys(ctx sdk.Context, ak authkeeper.AccountKeeper) (PubKeyMigrationResult, error) {
	var (
		res     PubKeyMigrationResult
		iterErr error
	)

	ak.IterateAccounts(ctx, func(acc sdk.AccountI) bool {
		res.Scanned++

		legacy, ok := acc.GetPubKey().(*legacyethsecp256k1.PubKey)
		if !ok {
			return false
		}

		converted := &evmethsecp256k1.PubKey{Key: legacy.Key}

		// Only convert when the key still resolves to the address it is stored
		// under. A mismatch means the record is already broken; do not mint a
		// working-looking key for it.
		if addr := converted.Address(); addr == nil || !sdk.AccAddress(addr).Equals(acc.GetAddress()) {
			if err := acc.SetPubKey(nil); err != nil {
				iterErr = fmt.Errorf("clearing unusable legacy pubkey on %s: %w", acc.GetAddress(), err)
				return true
			}
			ak.SetAccount(ctx, acc)
			res.Cleared++
			ctx.Logger().Info(
				"v34: cleared unusable legacy /ethereum.PubKey",
				"address", acc.GetAddress().String(),
			)
			return false
		}

		if err := acc.SetPubKey(converted); err != nil {
			iterErr = fmt.Errorf("converting legacy pubkey on %s: %w", acc.GetAddress(), err)
			return true
		}
		ak.SetAccount(ctx, acc)
		res.Converted++
		return false
	})

	if iterErr != nil {
		return res, iterErr
	}

	ctx.Logger().Info(
		"v34: migrated legacy /ethereum.PubKey accounts",
		"scanned", res.Scanned,
		"converted", res.Converted,
		"cleared", res.Cleared,
	)
	return res, nil
}
