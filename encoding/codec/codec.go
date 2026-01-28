package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	bitcoin "github.com/bitbadges/bitbadgeschain/chain-handlers/bitcoin/utils"
	ethereumcodec "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/codec"
	ethereum "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/utils"
	solana "github.com/bitbadges/bitbadgeschain/chain-handlers/solana/utils"
)

/**
  IMPORTANT: Even though these are not technically supported anymore, we need to keep them for
	legacy purposes (some accounts still have etheruem.PubKey and other dependent types).

	To fully remove this, we need to handle migrations of these accounts.
*/

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
