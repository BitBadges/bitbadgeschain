package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"
)

func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ethsecp256k1.PubKey{}, "ethereum/PubKey", nil)
	cdc.RegisterConcrete(&ethsecp256k1.PrivKey{}, "ethereum/PrivKey", nil)

	keyring.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)

	// HACK: Important for the SDK to work with amino (specifically an antehandler still uses "legacy.Cdc").
	legacy.Cdc = cdc
}
