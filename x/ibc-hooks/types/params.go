package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ paramtypes.ParamSet = &Params{}
)

// ParamKeyTable returns the param key table for the ibc-hooks module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default ibc-hooks module parameters.
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return nil
}
