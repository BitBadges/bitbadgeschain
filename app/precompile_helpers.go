package app

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
)

// RegisterAndEnableAllPrecompiles is a helper function for tests that registers and enables
// all precompiles (both default Cosmos and custom BitBadges precompiles).
// This ensures consistent precompile setup across all test files.
func RegisterAndEnableAllPrecompiles(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	tokenizationKeeper tokenizationkeeper.Keeper,
	gammKeeper gammkeeper.Keeper,
	sendManagerKeeper sendmanagerkeeper.Keeper,
) error {
	// Register custom BitBadges precompiles
	// Note: Default Cosmos precompiles are already registered via WithStaticPrecompiles
	// during EVM keeper initialization, so we only need to register custom ones here

	// Register tokenization precompile
	tokenizationPrecompile := tokenizationprecompile.NewPrecompile(tokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	evmKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)

	// Register gamm precompile
	gammPrecompile := gammprecompile.NewPrecompile(gammKeeper)
	gammPrecompileAddr := common.HexToAddress(gammprecompile.GammPrecompileAddress)
	evmKeeper.RegisterStaticPrecompile(gammPrecompileAddr, gammPrecompile)

	// Register sendmanager precompile
	sendManagerPrecompile := sendmanagerprecompile.NewPrecompile(sendManagerKeeper)
	sendManagerPrecompileAddr := common.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress)
	evmKeeper.RegisterStaticPrecompile(sendManagerPrecompileAddr, sendManagerPrecompile)

	// Enable all precompiles (both default Cosmos and custom BitBadges)
	allAddresses := []common.Address{}

	// Add default Cosmos precompiles (0x0001-0x0009)
	allAddresses = append(allAddresses, GetDefaultCosmosPrecompileAddresses()...)

	// Add custom BitBadges precompiles (0x1001+)
	allAddresses = append(allAddresses, GetAllCustomPrecompileAddresses()...)

	for _, addr := range allAddresses {
		if err := evmKeeper.EnableStaticPrecompiles(ctx, addr); err != nil {
			// Log error but don't fail if precompile is already enabled (idempotent)
			ctx.Logger().Info("Precompile enable attempt in test", "error", err, "address", addr.Hex())
		}
	}

	return nil
}

// GetAllCustomPrecompileAddresses returns all custom BitBadges precompile addresses
// This is useful for validation and testing
func GetAllCustomPrecompileAddresses() []common.Address {
	return []common.Address{
		common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress), // 0x1001 - Tokenization
		common.HexToAddress(gammprecompile.GammPrecompileAddress),                 // 0x1002 - Gamm
		common.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress),   // 0x1003 - SendManager
		// Next available address: 0x1004
	}
}

// GetAllDefaultCosmosPrecompileAddresses returns all default Cosmos precompile addresses
// This is useful for validation and testing
func GetAllDefaultCosmosPrecompileAddresses() []common.Address {
	return GetDefaultCosmosPrecompileAddresses()
}

// ValidateNoAddressCollisions checks that all precompile addresses are unique
// This helps prevent address collision bugs between default Cosmos and custom BitBadges precompiles
func ValidateNoAddressCollisions() error {
	// Check custom precompiles for duplicates
	customAddresses := GetAllCustomPrecompileAddresses()
	seen := make(map[common.Address]bool)

	for _, addr := range customAddresses {
		if seen[addr] {
			return fmt.Errorf("duplicate custom precompile address detected: %s", addr.Hex())
		}
		seen[addr] = true
	}

	// Check for collisions between default Cosmos and custom precompiles
	defaultAddresses := GetDefaultCosmosPrecompileAddresses()
	for _, defaultAddr := range defaultAddresses {
		if seen[defaultAddr] {
			return fmt.Errorf("custom precompile address collides with default Cosmos precompile: %s", defaultAddr.Hex())
		}
		seen[defaultAddr] = true
	}

	return nil
}
