// Copyright 2023 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/bitbadges/bitbadgeschain/x/ethermint/blob/main/LICENSE
package eip712

import (
	"errors"

	"github.com/bitbadges/bitbadgeschain/app/params"

	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	protoCodec codec.ProtoCodecMarshaler
	aminoCodec *codec.LegacyAmino
)

// SetEncodingConfig set the encoding config to the singleton codecs (Amino and Protobuf).
// The process of unmarshaling SignDoc bytes into a SignDoc object requires having a codec
// populated with all relevant message types. As a result, we must call this method on app
// initialization with the app's encoding config.
func SetEncodingConfig(cfg params.EncodingConfig) {
	aminoCodec = cfg.Amino
	protoCodec = codec.NewProtoCodec(cfg.InterfaceRegistry)
}

// validateCodecInit ensures that both Amino and Protobuf encoding codecs have been set on app init,
// so the module does not panic if either codec is not found.
func validateCodecInit() error {
	if aminoCodec == nil || protoCodec == nil {
		return errors.New("missing codec: codecs have not been properly initialized using SetEncodingConfig")
	}

	return nil
}