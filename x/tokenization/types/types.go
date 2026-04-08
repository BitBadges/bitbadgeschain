package types

import (
	"fmt"
	"strconv"
	"strings"
)

// Special addresses used throughout the tokens module
const (
	// MintAddress represents the special address for minting new tokens
	MintAddress = "Mint"

	// TotalAddress represents the special address for tracking total supply
	TotalAddress = "Total"

	// MintEscrowAlias is the alias for the collection's mint escrow bb1 address
	MintEscrowAlias = "MintEscrow"

	// CosmosWrapperPrefix is the prefix for cosmos wrapper path aliases (e.g., "CosmosWrapper/0")
	CosmosWrapperPrefix = "CosmosWrapper/"

	// IBCBackingAlias is the alias for the collection's IBC backing path address
	IBCBackingAlias = "IBCBacking"
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

// IsAddressAlias checks if the given string is a recognized address alias
func IsAddressAlias(address string) bool {
	if address == MintEscrowAlias || address == IBCBackingAlias {
		return true
	}
	if strings.HasPrefix(address, CosmosWrapperPrefix) {
		idx := address[len(CosmosWrapperPrefix):]
		if _, err := strconv.ParseUint(idx, 10, 64); err == nil {
			return true
		}
	}
	return false
}

// ValidateAddressAlias validates the syntax of an address alias (no collection context needed)
func ValidateAddressAlias(address string) error {
	if address == MintEscrowAlias || address == IBCBackingAlias {
		return nil
	}
	if strings.HasPrefix(address, CosmosWrapperPrefix) {
		idx := address[len(CosmosWrapperPrefix):]
		if idx == "" {
			return fmt.Errorf("CosmosWrapper alias must specify an index (e.g., CosmosWrapper/0)")
		}
		if _, err := strconv.ParseUint(idx, 10, 64); err != nil {
			return fmt.Errorf("CosmosWrapper alias index must be a non-negative integer: %s", idx)
		}
		return nil
	}
	return fmt.Errorf("unrecognized address alias: %s", address)
}

// IsReservedAliasListId checks if a list ID matches a reserved alias pattern.
// Used to prevent creation of stored address lists with names that would collide with aliases.
func IsReservedAliasListId(listId string) bool {
	return listId == MintEscrowAlias || listId == IBCBackingAlias || strings.HasPrefix(listId, CosmosWrapperPrefix)
}
