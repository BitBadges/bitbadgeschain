package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type ApprovalEventsTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestApprovalEventsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(ApprovalEventsTestSuite))
}

func (suite *ApprovalEventsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// filterApprovalChangeEvents returns only events with type "approvalChange"
func filterApprovalChangeEvents(events sdk.Events) []sdk.Event {
	var result []sdk.Event
	for _, e := range events {
		if e.Type == "approvalChange" {
			result = append(result, e)
		}
	}
	return result
}

// getEventAttribute returns the value of an attribute by key from an event
func getEventAttribute(event sdk.Event, key string) string {
	for _, attr := range event.Attributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}

// TestCreateApprovals_EmitsCreatedEvents verifies that adding new approvals emits "created" events
func (suite *ApprovalEventsTestSuite) TestCreateApprovals_EmitsCreatedEvents() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval-a", "AllWithoutMint", "All"),
		testutil.GenerateCollectionApproval("approval-b", "AllWithoutMint", "All"),
	}

	// Reset event manager to isolate our events
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Check response has approval changes
	suite.Require().Len(resp.ApprovalChanges, 2, "should report 2 approval changes")
	for _, change := range resp.ApprovalChanges {
		suite.Require().Equal("created", change.Action)
		suite.Require().Equal("collection", change.ApprovalLevel)
	}

	// Check events
	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	suite.Require().Len(events, 2, "should emit 2 approvalChange events")

	for _, event := range events {
		suite.Require().Equal(suite.CollectionId.String(), getEventAttribute(event, "collectionId"))
		suite.Require().Equal("created", getEventAttribute(event, "action"))
		suite.Require().Equal("collection", getEventAttribute(event, "approvalLevel"))
		suite.Require().Equal("", getEventAttribute(event, "approverAddress"), "collection-level should have empty approverAddress")
		suite.Require().NotEmpty(getEventAttribute(event, "approval"), "should include approval JSON")
	}

	// Verify specific approval IDs are in events
	eventIds := map[string]bool{}
	for _, event := range events {
		eventIds[getEventAttribute(event, "approvalId")] = true
	}
	suite.Require().True(eventIds["approval-a"])
	suite.Require().True(eventIds["approval-b"])
}

// TestEditApproval_EmitsEditedEvent verifies that modifying an existing approval emits an "edited" event
func (suite *ApprovalEventsTestSuite) TestEditApproval_EmitsEditedEvent() {
	// First, create an approval
	original := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("my-approval", "AllWithoutMint", "All"),
	}

	msg1 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          original,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Now edit it — change the toListId
	edited := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("my-approval", "AllWithoutMint", "AllWithoutMint"),
	}

	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          edited,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Len(resp.ApprovalChanges, 1)
	suite.Require().Equal("edited", resp.ApprovalChanges[0].Action)
	suite.Require().Equal("my-approval", resp.ApprovalChanges[0].ApprovalId)

	// Check events
	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	suite.Require().Len(events, 1, "should emit 1 approvalChange event for edit")
	suite.Require().Equal("edited", getEventAttribute(events[0], "action"))
	suite.Require().Equal("my-approval", getEventAttribute(events[0], "approvalId"))
}

// TestDeleteApproval_EmitsDeletedEvent verifies that removing an approval emits a "deleted" event
func (suite *ApprovalEventsTestSuite) TestDeleteApproval_EmitsDeletedEvent() {
	// Create two approvals
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("keep-me", "AllWithoutMint", "All"),
		testutil.GenerateCollectionApproval("delete-me", "AllWithoutMint", "All"),
	}

	msg1 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Now update with only one — the other is implicitly deleted
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	msg2 := &types.MsgSetCollectionApprovals{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CollectionApprovals: []*types.CollectionApproval{testutil.GenerateCollectionApproval("keep-me", "AllWithoutMint", "All")},
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Find the deleted change in response
	var deletedChange *types.ApprovalChange
	for _, c := range resp.ApprovalChanges {
		if c.Action == "deleted" {
			deletedChange = c
		}
	}
	suite.Require().NotNil(deletedChange, "should have a deleted change")
	suite.Require().Equal("delete-me", deletedChange.ApprovalId)

	// Check events
	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	deletedEvents := []sdk.Event{}
	for _, e := range events {
		if getEventAttribute(e, "action") == "deleted" {
			deletedEvents = append(deletedEvents, e)
		}
	}
	suite.Require().Len(deletedEvents, 1)
	suite.Require().Equal("delete-me", getEventAttribute(deletedEvents[0], "approvalId"))
	suite.Require().NotEmpty(getEventAttribute(deletedEvents[0], "approval"), "deleted event should include old approval JSON")
}

// TestUnchangedApproval_NoEvent verifies that re-submitting an identical approval emits no event
func (suite *ApprovalEventsTestSuite) TestUnchangedApproval_NoEvent() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("stable", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Submit the same approval again
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// No changes should be reported
	suite.Require().Empty(resp.ApprovalChanges, "unchanged approval should produce no changes")

	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	suite.Require().Empty(events, "unchanged approval should emit no events")
}

// TestMixedChanges_CorrectEventActions verifies create + edit + delete in a single update
func (suite *ApprovalEventsTestSuite) TestMixedChanges_CorrectEventActions() {
	// Setup: create two approvals
	initial := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("will-edit", "AllWithoutMint", "All"),
		testutil.GenerateCollectionApproval("will-delete", "AllWithoutMint", "All"),
	}

	msg1 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          initial,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Now: edit one, delete one, create one
	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	updated := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("will-edit", "AllWithoutMint", "AllWithoutMint"), // edited (toListId changed)
		testutil.GenerateCollectionApproval("brand-new", "AllWithoutMint", "All"),             // created
		// "will-delete" is omitted → deleted
	}

	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          updated,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Build action map from response
	actionMap := map[string]string{}
	for _, c := range resp.ApprovalChanges {
		actionMap[c.ApprovalId] = c.Action
	}
	suite.Require().Equal("edited", actionMap["will-edit"])
	suite.Require().Equal("created", actionMap["brand-new"])
	suite.Require().Equal("deleted", actionMap["will-delete"])

	// Verify events match
	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	suite.Require().Len(events, 3, "should emit 3 events (1 edit + 1 create + 1 delete)")

	eventActionMap := map[string]string{}
	for _, e := range events {
		eventActionMap[getEventAttribute(e, "approvalId")] = getEventAttribute(e, "action")
	}
	suite.Require().Equal("edited", eventActionMap["will-edit"])
	suite.Require().Equal("created", eventActionMap["brand-new"])
	suite.Require().Equal("deleted", eventActionMap["will-delete"])
}

// TestReviewItems_PopulatedOnResponse verifies reviewItems are returned in responses
func (suite *ApprovalEventsTestSuite) TestReviewItems_PopulatedOnResponse() {
	// Create collection with no approvals → should get "non-transferable" warning
	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			testutil.GenerateUintRange(1, 100),
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata:    testutil.GenerateCollectionMetadata("https://example.com", ""),
		TokenMetadata:         []*types.TokenMetadata{},
		CollectionApprovals:   []*types.CollectionApproval{},
		Standards:             []string{},
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(resp.ReviewItems, "should have review items")

	hasNonTransferable := false
	for _, item := range resp.ReviewItems {
		if item == "No transfer approvals set — tokens will be non-transferable" {
			hasNonTransferable = true
		}
	}
	suite.Require().True(hasNonTransferable, "should warn about non-transferable tokens")
}

// TestReviewItems_ApprovalChangeSummary verifies the summary string format
func (suite *ApprovalEventsTestSuite) TestReviewItems_ApprovalChangeSummary() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("new-one", "AllWithoutMint", "All"),
	}

	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	resp, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	hasSummary := false
	for _, item := range resp.ReviewItems {
		if item == "1 approvals created, 0 edited, 0 deleted" {
			hasSummary = true
		}
	}
	suite.Require().True(hasSummary, "should include approval change summary in reviewItems, got: %v", resp.ReviewItems)
}

// TestApprovalChangeEvent_HasVersionAttribute verifies version is populated on events
func (suite *ApprovalEventsTestSuite) TestApprovalChangeEvent_HasVersionAttribute() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("versioned", "AllWithoutMint", "All"),
	}

	suite.Ctx = suite.Ctx.WithEventManager(sdk.NewEventManager())

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}
	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	events := filterApprovalChangeEvents(suite.Ctx.EventManager().Events())
	suite.Require().Len(events, 1)

	version := getEventAttribute(events[0], "version")
	suite.Require().NotEmpty(version, "version attribute should be set")
}
