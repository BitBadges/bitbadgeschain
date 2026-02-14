package tokenization_test

import (
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
)

// IndexerEventsTestSuite tests that indexer events are properly emitted during EVM transactions
type IndexerEventsTestSuite struct {
	EVMKeeperIntegrationTestSuite
}

func TestIndexerEventsTestSuite(t *testing.T) {
	suite.Run(t, new(IndexerEventsTestSuite))
}

// SetupTest sets up the test suite
func (suite *IndexerEventsTestSuite) SetupTest() {
	suite.EVMKeeperIntegrationTestSuite.SetupTest()
}

// TestIndexerEvents_UniversalUpdateCollection verifies that indexer events are emitted
// when universalUpdateCollection is called through the EVM precompile
// Note: We use setManager as a simpler test case since it also calls UniversalUpdateCollection internally
func (suite *IndexerEventsTestSuite) TestIndexerEvents_UniversalUpdateCollection() {
	suite.SetupTest()

	// Clear events before transaction
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	// Build setManager transaction (which internally calls UniversalUpdateCollection)
	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	method := suite.Precompile.ABI.Methods["setManager"]
	suite.Require().NotNil(method, "setManager method should exist")

	// Build JSON for setManager
	jsonMsg, err := helpers.BuildSetManagerJSON(suite.Alice.String(), suite.CollectionId.BigInt(), suite.Bob.String())
	suite.Require().NoError(err, "Failed to build JSON")

	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.Require().NoError(err, "Failed to pack arguments")

	nonce := suite.getNonce(suite.AliceEVM)
	tx, err := helpers.BuildEVMTransaction(
		suite.AliceKey,
		&precompileAddr,
		input,
		big.NewInt(0),
		500000,
		big.NewInt(0),
		nonce,
		chainID,
	)
	suite.Require().NoError(err, "Failed to build EVM transaction")

	// Capture events before execution
	eventsBefore := suite.Ctx.EventManager().Events()
	eventCountBefore := len(eventsBefore)
	suite.T().Logf("Events before execution: %d", eventCountBefore)

	// Execute transaction
	response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)

	// Handle snapshot errors gracefully
	if err != nil && strings.Contains(err.Error(), "snapshot index") {
		suite.T().Skip("Skipping test due to snapshot error (known upstream bug)")
		return
	}
	if response != nil && strings.Contains(response.VmError, "snapshot revert error") {
		suite.T().Skip("Skipping test due to snapshot error (known upstream bug)")
		return
	}

	suite.Require().NoError(err, "Transaction should execute successfully")
	suite.Require().NotNil(response, "Response should not be nil")
	suite.Require().Empty(response.VmError, "Transaction should not have VM error: %s", response.VmError)

	// Get events from context after execution
	eventsAfter := suite.Ctx.EventManager().Events()
	eventCountAfter := len(eventsAfter)
	suite.T().Logf("Events after execution: %d (added: %d)", eventCountAfter, eventCountAfter-eventCountBefore)

	// Log all events for debugging
	suite.T().Logf("\n=== ALL EVENTS AFTER EXECUTION ===")
	for i, event := range eventsAfter {
		suite.T().Logf("Event[%d]: Type=%s, Attributes=%d", i, event.Type, len(event.Attributes))
		for j, attr := range event.Attributes {
			suite.T().Logf("  Attr[%d]: %s=%s", j, attr.Key, attr.Value)
		}
	}
	suite.T().Logf("=== END EVENTS ===\n")

	// Verify that indexer events were emitted
	indexerEvents := []sdk.Event{}
	for _, event := range eventsAfter {
		if event.Type == "indexer" {
			indexerEvents = append(indexerEvents, event)
		}
	}

	suite.T().Logf("Found %d indexer events", len(indexerEvents))
	suite.Require().NotEmpty(indexerEvents, "At least one indexer event should be emitted")

	// Verify that we have an indexer event for universal_update_collection
	// (setManager internally calls UniversalUpdateCollection)
	foundUniversalUpdateIndexerEvent := false
	for _, event := range indexerEvents {
		for _, attr := range event.Attributes {
			if attr.Key == "msg_type" && attr.Value == "universal_update_collection" {
				foundUniversalUpdateIndexerEvent = true
				suite.T().Logf("Found indexer event for universal_update_collection")
				// Log all attributes of this event
				suite.T().Logf("Indexer event attributes:")
				for _, a := range event.Attributes {
					suite.T().Logf("  %s=%s", a.Key, a.Value)
				}
				break
			}
		}
		if foundUniversalUpdateIndexerEvent {
			break
		}
	}

	suite.Require().True(foundUniversalUpdateIndexerEvent, "Indexer event for universal_update_collection should be emitted when setManager is called")

	// Also verify that the message event was emitted (from EmitMessageAndIndexerEvents)
	foundMessageEvent := false
	for _, event := range eventsAfter {
		if event.Type == sdk.EventTypeMessage {
			for _, attr := range event.Attributes {
				if attr.Key == "msg_type" && attr.Value == "universal_update_collection" {
					foundMessageEvent = true
					suite.T().Logf("Found message event for universal_update_collection")
					break
				}
			}
		}
		if foundMessageEvent {
			break
		}
	}

	suite.Require().True(foundMessageEvent, "Message event for universal_update_collection should be emitted")
}

// TestIndexerEvents_TransferTokens verifies that indexer events are emitted
// when transferTokens is called through the EVM precompile
func (suite *IndexerEventsTestSuite) TestIndexerEvents_TransferTokens() {
	suite.SetupTest()

	// Clear events before transaction
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	// Build transferTokens transaction
	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	toAddressesStr := []string{suite.Bob.String()}
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		suite.CollectionId.BigInt(),
		suite.Alice.String(),
		toAddressesStr,
		big.NewInt(10),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
	)
	suite.Require().NoError(err)

	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.Require().NoError(err)

	nonce := suite.getNonce(suite.AliceEVM)
	tx, err := helpers.BuildEVMTransaction(
		suite.AliceKey,
		&precompileAddr,
		input,
		big.NewInt(0),
		500000,
		big.NewInt(0),
		nonce,
		chainID,
	)
	suite.Require().NoError(err)

	// Capture events before execution
	eventsBefore := suite.Ctx.EventManager().Events()
	eventCountBefore := len(eventsBefore)

	// Execute transaction
	response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)

	// Handle snapshot errors gracefully
	if err != nil && strings.Contains(err.Error(), "snapshot index") {
		suite.T().Skip("Skipping test due to snapshot error (known upstream bug)")
		return
	}
	if response != nil && strings.Contains(response.VmError, "snapshot revert error") {
		suite.T().Skip("Skipping test due to snapshot error (known upstream bug)")
		return
	}

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Get events from context after execution
	eventsAfter := suite.Ctx.EventManager().Events()
	eventCountAfter := len(eventsAfter)

	suite.T().Logf("Events before: %d, after: %d, added: %d", eventCountBefore, eventCountAfter, eventCountAfter-eventCountBefore)

	// Log all events for debugging
	suite.T().Logf("\n=== ALL EVENTS AFTER EXECUTION ===")
	for i, event := range eventsAfter {
		suite.T().Logf("Event[%d]: Type=%s", i, event.Type)
	}
	suite.T().Logf("=== END EVENTS ===\n")

	// Verify that indexer events were emitted
	indexerEvents := []sdk.Event{}
	for _, event := range eventsAfter {
		if event.Type == "indexer" {
			indexerEvents = append(indexerEvents, event)
		}
	}

	suite.T().Logf("Found %d indexer events", len(indexerEvents))
	suite.Require().NotEmpty(indexerEvents, "At least one indexer event should be emitted")

	// Verify we have indexer events for the transfer_tokens message
	// Precompile-specific events are no longer emitted (events come from underlying message handlers)
	foundTransferIndexerEvent := false

	for _, event := range indexerEvents {
		for _, attr := range event.Attributes {
			if attr.Key == "msg_type" && attr.Value == "transfer_tokens" {
				foundTransferIndexerEvent = true
				suite.T().Logf("Found indexer event for transfer_tokens message")
				break
			}
		}
	}

	suite.Require().True(foundTransferIndexerEvent, "Indexer event for transfer_tokens message should be emitted")
}
