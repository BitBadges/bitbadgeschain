package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	bitcoin "github.com/bitbadges/bitbadgeschain/chain-handlers/bitcoin/utils"
	ethereumcodec "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/codec"
	ethereum "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/utils"
	solana "github.com/bitbadges/bitbadgeschain/chain-handlers/solana/utils"

	// EVM module types - required for JSON-RPC tx decoding
	evmcryptocodec "github.com/cosmos/evm/crypto/codec"
	evmethsecp256k1 "github.com/cosmos/evm/crypto/ethsecp256k1"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
)

/**
  The legacy `/ethereum.PubKey` types below are this chain's pre-cosmos/evm EVM
  key. They stay registered so historical accounts and txs still decode.

  As of v34 the accounts themselves are migrated to cosmos/evm's key type by
  app/upgrades/v34/pubkeys.go, so nothing new is written with the legacy type.
  Keep these registrations regardless: dropping them would break decoding of
  every historical tx and block that carried one.
*/

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	ethereumcodec.RegisterCrypto(cdc)

	// cosmos/evm's own ethsecp256k1 key. See RegisterInterfaces below for why
	// this was missing and what it broke.
	//
	// Register the concretes directly rather than calling
	// evmcryptocodec.RegisterCrypto: that helper also calls
	// keyring.RegisterLegacyAminoCodec and cryptocodec.RegisterCrypto, both of
	// which ethereumcodec.RegisterCrypto above has already done, and amino
	// panics on a duplicate interface registration.
	cdc.RegisterConcrete(&evmethsecp256k1.PubKey{}, evmethsecp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(&evmethsecp256k1.PrivKey{}, evmethsecp256k1.PrivKeyName, nil)
	// Note: evmtypes amino registration is handled by the EVM module itself
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	ethereumcodec.RegisterInterfaces(interfaceRegistry)

	// cosmos/evm's ethsecp256k1 key type. This registration was missing, and
	// the omission is the root cause of BB-12.
	//
	// Only the legacy `/ethereum.PubKey` above was registered, so
	// `/cosmos.evm.crypto.v1.ethsecp256k1.PubKey` — the type cosmos/evm's ante
	// handler actually recognises, and the type an EVM account should carry —
	// could not be decoded at all:
	//
	//	no concrete type registered for type URL
	//	/cosmos.evm.crypto.v1.ethsecp256k1.PubKey against interface *types.PubKey
	//
	// EVM users appeared fine only because the common path is MsgEthereumTx via
	// the precompile, where the signature is recovered from the Ethereum
	// transaction and no Cosmos pubkey is involved. Signing a *Cosmos* message
	// from an EVM wallet does need one, which is why a gamm swap simulate was
	// where this finally surfaced.
	evmcryptocodec.RegisterInterfaces(interfaceRegistry)
	ethereum.RegisterInterfaces(interfaceRegistry)
	solana.RegisterInterfaces(interfaceRegistry)
	bitcoin.RegisterInterfaces(interfaceRegistry)

	// EVM module types - required for JSON-RPC tx decoding (MsgEthereumTx, etc.)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	erc20types.RegisterInterfaces(interfaceRegistry)
	feemarkettypes.RegisterInterfaces(interfaceRegistry)
}
