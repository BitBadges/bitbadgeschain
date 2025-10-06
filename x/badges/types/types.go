package types

// Special addresses used throughout the badges module
const (
	// MintAddress represents the special address for minting new badges
	MintAddress = "Mint"

	// TotalAddress represents the special address for tracking total supply
	TotalAddress = "Total"
)

// IsSpecialAddress checks if the given address is a special system address
func IsSpecialAddress(address string) bool {
	return address == MintAddress || address == TotalAddress
}

// IsMintAddress checks if the given address is the mint address
func IsMintAddress(address string) bool {
	return address == MintAddress
}

// IsTotalAddress checks if the given address is the total address
func IsTotalAddress(address string) bool {
	return address == TotalAddress
}
