package app

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// appVersionGetterAdapter provides GetAppVersion for wasmd's IBCHandler
// wasmd expects an appVersionGetter interface to get the application version
// by extracting it from the channel version string (stripping middleware data)
type appVersionGetterAdapter struct {
	app *App
}

// GetAppVersion implements the appVersionGetter interface expected by wasmd
// It extracts the application-level version from the channel version string
// by stripping out middleware version data. This is used for IBC version negotiation.
func (a *appVersionGetterAdapter) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	// Get the channel to access its version string
	channel, found := a.app.IBCKeeper.ChannelKeeper.GetChannel(ctx, portID, channelID)
	if !found {
		return "", false
	}

	// Extract app version from channel version string
	// Channel version format in IBC v10: "{app_version}" or "{middleware1}/{middleware2}/{app_version}"
	// We need to extract the base app version by stripping middleware prefixes
	version := channel.Version

	// If version contains slashes, it means middleware is present
	// The app version is typically the last component after all middleware
	// However, for wasm channels, the version might be just the wasm version
	// We'll return the full version and let wasmd handle the parsing
	// as it knows how to extract its own version from the string

	// For wasm ports, the version should be the wasm module version
	// Check if this is a wasm port
	if strings.HasPrefix(portID, "wasm.") {
		// For wasm channels, return the channel version as-is
		// wasmd will parse it to extract the contract version
		return version, true
	}

	// For non-wasm channels, try to extract the base app version
	// This is a best-effort approach - the actual parsing depends on middleware stack
	parts := strings.Split(version, "/")
	if len(parts) > 0 {
		// Return the last component as it's typically the app version
		return parts[len(parts)-1], true
	}

	// Fallback: return the version as-is
	return version, true
}
