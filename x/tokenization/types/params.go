package types

import (
	sdkmath "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		AllowedDenoms: []string{
			"ubadge",
			"ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349",
			"ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701",
			"ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
		},
		AffiliatePercentage: sdkmath.NewUint(500),
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
