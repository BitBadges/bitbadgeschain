package types

import (
	sdkmath "cosmossdk.io/math"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                DefaultParams(),
		ManagerSplitters:      []*ManagerSplitter{},
		NextManagerSplitterId: sdkmath.NewUint(1),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate all manager splitters
	for _, ms := range gs.ManagerSplitters {
		if ms.Admin == "" {
			return ErrInvalidAdmin
		}
		if ms.Address == "" {
			return ErrInvalidAddress
		}
	}

	return nil
}

