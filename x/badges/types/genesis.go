package types

import (
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		Params:      DefaultParams(),
		NextAssetId: 0,
		Badges:      []*BitBadge{},
		Balances:    []*BadgeBalanceInfo{},
		BalanceIds:  []string{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.

//TODO: validate badges and owners are formatted correctly
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
