package app

import (
	"fmt"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

// registerSendManagerRouters registers the tokenization module router with sendmanager
// This is deferred until after both keepers are created to avoid circular dependency
func (app *App) registerSendManagerRouters() error {
	// Register tokenization module routers for both prefixes
	tokenizationRouter := sendmanagerkeeper.NewTokenizationAliasDenomRouter(app.TokenizationKeeper)

	// Register tokenizationlp: prefix using the exported constant from tokenization keeper
	if err := app.SendmanagerKeeper.RegisterRouter(tokenizationkeeper.AliasDenomPrefix, tokenizationRouter); err != nil {
		return fmt.Errorf("failed to register tokenization alias denom prefix router: %w", err)
	}

	return nil
}
