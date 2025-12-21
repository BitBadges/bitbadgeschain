package app

import (
	"fmt"

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

// registerSendManagerRouters registers the badges module router with sendmanager
// This is deferred until after both keepers are created to avoid circular dependency
func (app *App) registerSendManagerRouters() error {
	// Register badges module routers for both prefixes
	badgesRouter := sendmanagerkeeper.NewBadgesAliasDenomRouter(app.BadgesKeeper)

	// Register badgeslp: prefix using the exported constant from badges keeper
	if err := app.SendmanagerKeeper.RegisterRouter(badgeskeeper.AliasDenomPrefix, badgesRouter); err != nil {
		return fmt.Errorf("failed to register badges alias denom prefix router: %w", err)
	}

	return nil
}
