package keeper_test

import (
	// "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	// "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestAccountsAreCreatedAndStoredCorrectly() {
	sampleAccNumber := suite.app.BadgesKeeper.GetOrCreateAccountNumberForAccAddressBech32(suite.ctx, sdk.MustAccAddressFromBech32(alice))
	suite.Equal(sampleAccNumber, uint64(7)) //see test suite setup for why it is 1009
}