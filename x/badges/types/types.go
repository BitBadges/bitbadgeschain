package types

// Special addresses used throughout the tokens module
const (
	// MintAddress represents the special address for minting new tokens
	MintAddress = "Mint"

	// TotalAddress represents the special address for tracking total supply
	TotalAddress = "Total"
)

// IsMintOrTotalAddress checks if the given address is a mint or total address
func IsMintOrTotalAddress(address string) bool {
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
