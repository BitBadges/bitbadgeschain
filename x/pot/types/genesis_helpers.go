package types

// DefaultGenesis returns the default genesis state for x/pot.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic validation of the genesis state.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
