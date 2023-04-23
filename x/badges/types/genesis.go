package types

import (
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		Params:           DefaultParams(),
		NextCollectionId: 1,
		NextClaimId:      1,
		Collections:      []*BadgeCollection{},
		Balances:         []*UserBalance{},
		BalanceIds:       []string{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.

// IMPORTANT: We assume badges are well-formed and validated here
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
