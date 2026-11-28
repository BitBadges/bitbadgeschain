package types_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
)

func TestExclusivityLatch_PreventsFulfillVsRefundRace(t *testing.T) {
	engine := types.NewExclusivityEngine()
	spec := types.ExclusivityGroupSpec{GroupId: "bounty-42", MaxExecutions: 1}

	var fulfillWins, refundWins atomic.Int32
	var wg sync.WaitGroup
	start := make(chan struct{})

	// Launch 20 competing goroutines (10 fulfill vs 10 refund) simultaneously.
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			if err := engine.ExecuteLatch(spec, "bounty-fulfill", nil); err == nil {
				fulfillWins.Add(1)
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			if err := engine.ExecuteLatch(spec, "bounty-refund", nil); err == nil {
				refundWins.Add(1)
			}
		}()
	}

	close(start)
	wg.Wait()

	totalWins := fulfillWins.Load() + refundWins.Load()
	require.Equal(t, int32(1), totalWins, "strictly one branch (either fulfill OR refund) may win")
	require.Equal(t, uint64(1), engine.LatchCount("bounty-42"))

	// Any subsequent attempt must fail with ErrExclusivityConflict.
	err := engine.ExecuteLatch(spec, "late-refund", nil)
	require.ErrorIs(t, err, types.ErrExclusivityConflict)
}

func TestExclusivityLatch_RollsBackOnCommitFailure(t *testing.T) {
	engine := types.NewExclusivityEngine()
	spec := types.ExclusivityGroupSpec{GroupId: "bounty-rollback", MaxExecutions: 1}

	boom := errors.New("insufficient escrow coins")
	err := engine.ExecuteLatch(spec, "bounty-fulfill", func() error { return boom })
	require.ErrorIs(t, err, boom)
	require.Equal(t, uint64(0), engine.LatchCount("bounty-rollback"))

	// After rollback, refund can still cleanly succeed.
	require.NoError(t, engine.ExecuteLatch(spec, "bounty-refund", nil))
	require.Equal(t, uint64(1), engine.LatchCount("bounty-rollback"))
}

func TestStateMachine_CrowdfundMultiBackerRefundsVsCreatorRelease(t *testing.T) {
	engine := types.NewExclusivityEngine()
	const key = "crowdfund:7"

	contribute := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"OPEN"}, NextState: "OPEN"}
	openRefunds := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"OPEN"}, NextState: "REFUNDING"}
	claimRefund := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"REFUNDING"}, NextState: "REFUNDING"}
	releaseToCreator := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"OPEN"}, NextState: "RELEASED", Terminal: true}

	// Multiple backers contribute while OPEN.
	require.NoError(t, engine.ExecuteTransition(contribute, "backer-1-contribute", nil))
	require.NoError(t, engine.ExecuteTransition(contribute, "backer-2-contribute", nil))
	require.Equal(t, "OPEN", engine.GetState(key).State)

	// Goal missed -> open refunds.
	require.NoError(t, engine.ExecuteTransition(openRefunds, "open-refunds", nil))
	require.Equal(t, "REFUNDING", engine.GetState(key).State)

	// Creator can NO LONGER drain escrow via releaseToCreator!
	err := engine.ExecuteTransition(releaseToCreator, "creator-drain-attempt", nil)
	require.ErrorIs(t, err, types.ErrInvalidStateTransition)

	// Multiple backers can each claim refunds safely in REFUNDING state.
	require.NoError(t, engine.ExecuteTransition(claimRefund, "backer-1-refund", nil))
	require.NoError(t, engine.ExecuteTransition(claimRefund, "backer-2-refund", nil))
	require.Len(t, engine.AuditLog(), 5)
}

func TestStateMachine_RefundableInvoiceLifecycleAndTerminalLock(t *testing.T) {
	engine := types.NewExclusivityEngine()
	const key = "invoice:99"

	pay := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"OPEN"}, NextState: "PAID"}
	cancel := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"OPEN"}, NextState: "CANCELLED", Terminal: true}
	refund := types.StateTransitionSpec{StateKey: key, RequiredStates: []string{"PAID"}, NextState: "REFUNDED", Terminal: true}

	// Cannot refund an unpaid invoice.
	require.ErrorIs(t, engine.ExecuteTransition(refund, "premature-refund", nil), types.ErrInvalidStateTransition)

	// Pay succeeds (OPEN -> PAID).
	require.NoError(t, engine.ExecuteTransition(pay, "pay-invoice", nil))
	require.Equal(t, "PAID", engine.GetState(key).State)

	// Cancel fails once paid.
	require.ErrorIs(t, engine.ExecuteTransition(cancel, "late-cancel", nil), types.ErrInvalidStateTransition)

	// Refund succeeds and locks terminal state (PAID -> REFUNDED).
	require.NoError(t, engine.ExecuteTransition(refund, "merchant-refund", nil))
	slot := engine.GetState(key)
	require.Equal(t, "REFUNDED", slot.State)
	require.True(t, slot.Terminal)

	// Second refund fails because REFUNDED is terminal.
	require.ErrorIs(t, engine.ExecuteTransition(refund, "double-refund", nil), types.ErrTerminalStateLocked)
}

func TestStateMachine_CompiledExclusivitySugarMatchesLatchSemantics(t *testing.T) {
	engine := types.NewExclusivityEngine()
	rule := types.CompileExclusivityGroupToStateTransition("bounty-sugar")

	require.NoError(t, engine.ExecuteTransition(rule, "fulfill", nil))
	err := engine.ExecuteTransition(rule, "refund", nil)
	require.ErrorIs(t, err, types.ErrTerminalStateLocked)
}
