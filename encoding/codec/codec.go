package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	bitcoin "github.com/bitbadges/bitbadgeschain/chain-handlers/bitcoin/utils"
	ethereumcodec "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/codec"
	ethereum "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/utils"
	solana "github.com/bitbadges/bitbadgeschain/chain-handlers/solana/utils"
)

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	ethereumcodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	ethereumcodec.RegisterInterfaces(interfaceRegistry)
	ethereum.RegisterInterfaces(interfaceRegistry)
	solana.RegisterInterfaces(interfaceRegistry)
	bitcoin.RegisterInterfaces(interfaceRegistry)
}
