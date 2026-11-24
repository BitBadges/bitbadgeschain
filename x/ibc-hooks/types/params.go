package types

// DefaultParams returns default ibc-hooks module parameters.
func DefaultParams() Params {
	return Params{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return nil
}
