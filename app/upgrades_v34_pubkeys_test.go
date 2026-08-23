package app

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmethsecp256k1 "github.com/cosmos/evm/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	legacyethsecp256k1 "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"
	v34 "github.com/bitbadges/bitbadgeschain/app/upgrades/v34"
)

// newEVMKeyPair returns a real ethsecp256k1 key in both representations: this
// chain's original /ethereum.PubKey and cosmos/evm's. They wrap identical key
// bytes, which is what makes the v34 migration lossless.
func newEVMKeyPair(t *testing.T) (*legacyethsecp256k1.PubKey, *evmethsecp256k1.PubKey) {
	t.Helper()

	priv, err := evmethsecp256k1.GenerateKey()
	require.NoError(t, err)

	keyBytes := priv.PubKey().Bytes()
	legacy := &legacyethsecp256k1.PubKey{Key: keyBytes}
	modern := &evmethsecp256k1.PubKey{Key: keyBytes}

	// The premise of the whole migration: same bytes, same address.
	require.Equal(t, modern.Address().Bytes(), legacy.Address().Bytes(),
		"legacy and cosmos/evm ethsecp256k1 must derive the same address")

	return legacy, modern
}

func TestV34MigrateLegacyEthereumPubKeys(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	ak := app.AccountKeeper

	// 1. A legacy EVM account — the case this migration exists for.
	legacyPub, modernPub := newEVMKeyPair(t)
	legacyAddr := sdk.AccAddress(legacyPub.Address())
	ak.SetAccount(ctx, authtypes.NewBaseAccount(legacyAddr, legacyPub, 41, 7))

	// 2. An account already on cosmos/evm's key — must not be touched.
	_, alreadyModern := newEVMKeyPair(t)
	modernAddr := sdk.AccAddress(alreadyModern.Address())
	ak.SetAccount(ctx, authtypes.NewBaseAccount(modernAddr, alreadyModern, 42, 3))

	// 3. A plain Cosmos secp256k1 account — must not be touched.
	cosmosPub := secp256k1.GenPrivKey().PubKey()
	cosmosAddr := sdk.AccAddress(cosmosPub.Address())
	ak.SetAccount(ctx, authtypes.NewBaseAccount(cosmosAddr, cosmosPub, 43, 5))

	// 4. An account that has never signed (nil pubkey) — must not be touched.
	neverSignedAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	ak.SetAccount(ctx, authtypes.NewBaseAccount(neverSignedAddr, nil, 44, 0))

	res, err := v34.MigrateLegacyEthereumPubKeys(ctx, ak)
	require.NoError(t, err)
	require.Equal(t, 1, res.Converted, "exactly the one legacy account should convert")
	require.Equal(t, 0, res.Cleared)

	// The legacy account now carries cosmos/evm's type, with the key bytes and
	// account metadata intact.
	got := ak.GetAccount(ctx, legacyAddr)
	require.NotNil(t, got)
	migrated, ok := got.GetPubKey().(*evmethsecp256k1.PubKey)
	require.True(t, ok, "expected cosmos/evm ethsecp256k1.PubKey, got %T", got.GetPubKey())
	require.Equal(t, modernPub.Key, migrated.Key, "key bytes must survive unchanged")
	require.Equal(t, legacyAddr, got.GetAddress(), "address must not move")
	require.Equal(t, uint64(41), got.GetAccountNumber(), "account number must survive")
	require.Equal(t, uint64(7), got.GetSequence(), "sequence must survive")

	// And it is now the type cosmos/evm's ante handler actually matches, which
	// is the entire point — the legacy type failed that switch.
	require.IsType(t, &evmethsecp256k1.PubKey{}, ak.GetAccount(ctx, legacyAddr).GetPubKey())

	// Untouched accounts stay exactly as they were.
	require.IsType(t, &evmethsecp256k1.PubKey{}, ak.GetAccount(ctx, modernAddr).GetPubKey())
	require.Equal(t, uint64(3), ak.GetAccount(ctx, modernAddr).GetSequence())
	require.IsType(t, &secp256k1.PubKey{}, ak.GetAccount(ctx, cosmosAddr).GetPubKey())
	require.Equal(t, uint64(5), ak.GetAccount(ctx, cosmosAddr).GetSequence())
	require.Nil(t, ak.GetAccount(ctx, neverSignedAddr).GetPubKey())
}

// A legacy key stored under an address it does not derive is already unusable.
// The migration must not mint a working-looking key for it; it clears the key so
// the ante handler's SetPubKeyDecorator can repopulate it correctly on the next
// signed tx.
func TestV34MigrateLegacyEthereumPubKeys_MismatchedAddressIsCleared(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	ak := app.AccountKeeper

	legacyPub, _ := newEVMKeyPair(t)
	wrongAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	require.NotEqual(t, wrongAddr, sdk.AccAddress(legacyPub.Address()))

	acc := authtypes.NewBaseAccount(wrongAddr, nil, 12, 4)
	require.NoError(t, acc.SetPubKey(legacyPub))
	ak.SetAccount(ctx, acc)

	res, err := v34.MigrateLegacyEthereumPubKeys(ctx, ak)
	require.NoError(t, err)
	require.Equal(t, 0, res.Converted)
	require.Equal(t, 1, res.Cleared)

	got := ak.GetAccount(ctx, wrongAddr)
	require.NotNil(t, got)
	require.Nil(t, got.GetPubKey(), "unusable legacy key must be cleared, not converted")
	require.Equal(t, uint64(12), got.GetAccountNumber())
	require.Equal(t, uint64(4), got.GetSequence())
}

// Running the migration twice must be a no-op the second time. Upgrade handlers
// get re-run in replay and recovery scenarios.
func TestV34MigrateLegacyEthereumPubKeys_Idempotent(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	ak := app.AccountKeeper

	legacyPub, _ := newEVMKeyPair(t)
	addr := sdk.AccAddress(legacyPub.Address())
	ak.SetAccount(ctx, authtypes.NewBaseAccount(addr, legacyPub, 51, 2))

	first, err := v34.MigrateLegacyEthereumPubKeys(ctx, ak)
	require.NoError(t, err)
	require.Equal(t, 1, first.Converted)

	second, err := v34.MigrateLegacyEthereumPubKeys(ctx, ak)
	require.NoError(t, err)
	require.Equal(t, 0, second.Converted, "second run must find nothing to do")
	require.Equal(t, 0, second.Cleared)

	require.IsType(t, &evmethsecp256k1.PubKey{}, ak.GetAccount(ctx, addr).GetPubKey())
	require.Equal(t, uint64(2), ak.GetAccount(ctx, addr).GetSequence())
}
