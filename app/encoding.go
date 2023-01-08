package app

import (
	"github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/encoding"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := encoding.MakeConfig(ModuleBasics)
	return params.EncodingConfig{
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Marshaler:         encodingConfig.Codec,
		TxConfig:          encodingConfig.TxConfig,
		Amino:             encodingConfig.Amino,
	}
}
