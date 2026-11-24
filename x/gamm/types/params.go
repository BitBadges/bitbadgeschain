package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Parameter store keys.
var (
	KeyPoolCreationFee = []byte("PoolCreationFee")
)

func NewParams() Params {
	return Params{}
}

// default gamm module parameters.
func DefaultParams() Params {
	return Params{}
}

// validate params.
func (p Params) Validate() error {
	return nil
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid pool creation fee: %+v", i)
	}

	return nil
}
