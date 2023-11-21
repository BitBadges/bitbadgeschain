package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermintcodec "github.com/bitbadges/bitbadgeschain/x/ethermint/crypto/codec"
	ethermint "github.com/bitbadges/bitbadgeschain/x/ethermint/utils"
	solana "github.com/bitbadges/bitbadgeschain/x/solana/utils"
)

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdk.RegisterLegacyAminoCodec(cdc)
	ethermintcodec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	std.RegisterInterfaces(interfaceRegistry)
	ethermintcodec.RegisterInterfaces(interfaceRegistry)
	ethermint.RegisterInterfaces(interfaceRegistry)
	solana.RegisterInterfaces(interfaceRegistry)
}
