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
  IMPORTANT: Even though these are not technically supported anymore, we need to keep them for
	legacy purposes (some accounts still have etheruem.PubKey and other dependent types).

	To fully remove this, we need to handle migrations of these accounts.
*/

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	ethereumcodec.RegisterCrypto(cdc)

	// Register cosmos/evm's `ethsecp256k1.PubKey` / `PrivKey` on the legacy
	// amino codec under their canonical amino names. Required because
	// `ConsumeTxSizeGasDecorator` (and any other ante decorator that
	// invokes `LegacyAmino.MustMarshal` on a tx for size estimation)
	// fails with "Cannot encode unregistered concrete type
	// ethsecp256k1.PubKey" if the SignerInfo's pubkey type isn't
	// registered as an amino concrete. The proto interface registration
	// in `RegisterInterfaces` covers tx decode, but the gas-size path
	// goes through legacy amino specifically.
	//
	// We register inline rather than calling `evmcryptocodec.RegisterCrypto`
	// because that helper also re-runs `cryptocodec.RegisterCrypto(cdc)` —
	// already invoked by `ethereumcodec.RegisterCrypto` above — which
	// double-registers the cosmos-sdk standard crypto types and panics
	// with "TypeInfo already exists for types.PubKey".
	cdc.RegisterConcrete(&evmethsecp256k1.PubKey{}, "cosmos-evm/PubKeyEthSecp256k1", nil)
	cdc.RegisterConcrete(&evmethsecp256k1.PrivKey{}, "cosmos-evm/PrivKeyEthSecp256k1", nil)
	_ = evmcryptocodec.RegisterCrypto // imported above for grep-ability; kept registered manually to avoid duplicate cryptocodec call
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	ethereumcodec.RegisterInterfaces(interfaceRegistry)
	ethereum.RegisterInterfaces(interfaceRegistry)
	solana.RegisterInterfaces(interfaceRegistry)
	bitcoin.RegisterInterfaces(interfaceRegistry)

	// cosmos/evm crypto types — registers `cosmos.evm.crypto.v1.ethsecp256k1.PubKey`
	// on the InterfaceRegistry so the chain can decode EIP-712-signed txs whose
	// SignerInfo wraps the pubkey under that canonical type URL. The legacy
	// `ethereum.PubKey` registration above is kept only for already-existing
	// accounts; new EVM signing flows route through this path.
	evmcryptocodec.RegisterInterfaces(interfaceRegistry)

	// EVM module types - required for JSON-RPC tx decoding (MsgEthereumTx, etc.)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	erc20types.RegisterInterfaces(interfaceRegistry)
	feemarkettypes.RegisterInterfaces(interfaceRegistry)
}
