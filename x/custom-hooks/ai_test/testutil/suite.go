package testutil

import (
	"testing"

	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
)

// AITestSuite provides a comprehensive test suite for AI-generated tests
// This complements the existing keeper_test.go which uses real app keepers
// Note: KeeperTestHelper already embeds suite.Suite, so we don't embed it again
type AITestSuite struct {
	apptesting.KeeperTestHelper

	Keeper keeper.Keeper

	// Test addresses
	Alice   string
	Bob     string
	Charlie string
}

// SetupTest initializes the test suite with a fresh keeper and context
func (suite *AITestSuite) SetupTest() {
	// Initialize app if not already done - Reset() will handle this, but ensure it's safe
	if suite.App == nil {
		suite.Setup()
	} else {
		suite.Reset()
	}

	// Create keeper with real app keepers (same as keeper_test.go)
	// IBC v10: ScopedIBCTransferKeeper removed - capabilities no longer used
	suite.Keeper = keeper.NewKeeper(
		suite.App.Logger(),
		&suite.App.GammKeeper,
		suite.App.BankKeeper,
		suite.App.TokenizationKeeper,
		&suite.App.SendmanagerKeeper,
		suite.App.TransferKeeper,
		suite.App.HooksICS4Wrapper,
		suite.App.IBCKeeper.ChannelKeeper,
	)

	// Initialize test addresses - TestAccs should be initialized by Reset()
	// Use fallback if not available (shouldn't happen, but safe)
	if len(suite.TestAccs) >= 3 {
		suite.Alice = suite.TestAccs[0].String()
		suite.Bob = suite.TestAccs[1].String()
		suite.Charlie = suite.TestAccs[2].String()
	} else {
		// Fallback addresses - these are valid Bech32 addresses
		suite.Alice = "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
		suite.Bob = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
		suite.Charlie = "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	}
}

// RunTestSuite runs a test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

// Helper function to get error from acknowledgement
func GetAckError(ack ibcexported.Acknowledgement) string {
	if !ack.Success() {
		// Extract error from acknowledgement
		if channelAck, ok := ack.(channeltypes.Acknowledgement); ok {
			if errResp, ok := channelAck.Response.(*channeltypes.Acknowledgement_Error); ok {
				return errResp.Error
			}
		}
	}
	return ""
}
