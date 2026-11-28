package types

import (
	"errors"
	"fmt"
	"sync"
)

// Prototype errors for BB-376 native approval mutual-exclusion and state-machine transitions.
var (
	ErrExclusivityConflict    = errors.New("exclusivity group already exhausted")
	ErrTerminalStateLocked    = errors.New("state machine key is in a terminal state")
	ErrInvalidStateTransition = errors.New("current state does not satisfy required transition states")
	ErrInvalidTransitionRule  = errors.New("invalid state transition rule")
)

// DefaultInitialState is the implicit state of any newly created StateMachine key.
const DefaultInitialState = "OPEN"

// ExclusivityGroupSpec models Option A from ADR-mutual-exclusion-state-machine.md:
// multiple approvals sharing GroupId share an atomic execution counter capped at MaxExecutions.
type ExclusivityGroupSpec struct {
	GroupId       string `json:"groupId"`
	MaxExecutions uint64 `json:"maxExecutions"` // defaults to 1 if 0
}

// StateTransitionSpec models Option B from ADR-mutual-exclusion-state-machine.md:
// an approval asserts that StateKey is in one of RequiredStates and atomically
// advances it to NextState (optionally locking it permanently when Terminal is true).
type StateTransitionSpec struct {
	StateKey       string   `json:"stateKey"`
	RequiredStates []string `json:"requiredStates"` // empty => []string{DefaultInitialState}
	NextState      string   `json:"nextState"`
	Terminal       bool     `json:"terminal"`
}

// StateSlot holds the current state for a single (collection, stateKey) slot.
type StateSlot struct {
	State      string `json:"state"`
	Terminal   bool   `json:"terminal"`
	LastBranch string `json:"lastBranch"`
	Version    uint64 `json:"version"`
}

// TransitionAuditEntry records an executed transition for deterministic inspection.
type TransitionAuditEntry struct {
	StateKey  string `json:"stateKey"`
	BranchId  string `json:"branchId"`
	FromState string `json:"fromState"`
	ToState   string `json:"toState"`
	Terminal  bool   `json:"terminal"`
}

// ExclusivityEngine is a thread-safe prototype implementing both Option A
// (exclusivityGroupId atomic latches) and Option B (lightweight StateTransition
// state machines) with atomic rollback on transfer failure.
type ExclusivityEngine struct {
	mu      sync.Mutex
	latches map[string]uint64
	winners map[string][]string
	states  map[string]StateSlot
	audit   []TransitionAuditEntry
}

// NewExclusivityEngine creates an empty ExclusivityEngine.
func NewExclusivityEngine() *ExclusivityEngine {
	return &ExclusivityEngine{
		latches: make(map[string]uint64),
		winners: make(map[string][]string),
		states:  make(map[string]StateSlot),
	}
}

// ExecuteLatch atomically claims a slot in spec.GroupId and invokes onCommit.
// If onCommit returns an error, the latch mutation is rolled back atomically
// (mirroring Cosmos evmcompat.AtomicContext).
func (e *ExclusivityEngine) ExecuteLatch(spec ExclusivityGroupSpec, branchId string, onCommit func() error) error {
	if spec.GroupId == "" {
		return fmt.Errorf("%w: groupId cannot be empty", ErrInvalidTransitionRule)
	}
	maxExec := spec.MaxExecutions
	if maxExec == 0 {
		maxExec = 1
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	used := e.latches[spec.GroupId]
	if used >= maxExec {
		return fmt.Errorf("%w: group %q reached limit %d (winners=%v)", ErrExclusivityConflict, spec.GroupId, maxExec, e.winners[spec.GroupId])
	}

	if onCommit != nil {
		if err := onCommit(); err != nil {
			return err
		}
	}

	e.latches[spec.GroupId] = used + 1
	e.winners[spec.GroupId] = append(e.winners[spec.GroupId], branchId)
	return nil
}

// LatchCount returns how many executions have succeeded for groupId.
func (e *ExclusivityEngine) LatchCount(groupId string) uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.latches[groupId]
}

// ExecuteTransition atomically checks RequiredStates, invokes onCommit, and
// transitions StateKey to NextState. If onCommit returns an error, state is untouched.
func (e *ExclusivityEngine) ExecuteTransition(spec StateTransitionSpec, branchId string, onCommit func() error) error {
	if spec.StateKey == "" || spec.NextState == "" {
		return fmt.Errorf("%w: stateKey and nextState are required", ErrInvalidTransitionRule)
	}
	for _, req := range spec.RequiredStates {
		if req == "" {
			return fmt.Errorf("%w: requiredStates cannot contain empty state name", ErrInvalidTransitionRule)
		}
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	slot, exists := e.states[spec.StateKey]
	if !exists {
		slot = StateSlot{State: DefaultInitialState}
	}

	if slot.Terminal {
		return fmt.Errorf("%w: key %q locked in terminal state %q", ErrTerminalStateLocked, spec.StateKey, slot.State)
	}

	required := spec.RequiredStates
	if len(required) == 0 {
		required = []string{DefaultInitialState}
	}

	matched := false
	for _, req := range required {
		if slot.State == req {
			matched = true
			break
		}
	}
	if !matched {
		return fmt.Errorf("%w: key %q is in state %q, expected one of %v", ErrInvalidStateTransition, spec.StateKey, slot.State, required)
	}

	if onCommit != nil {
		if err := onCommit(); err != nil {
			return err
		}
	}

	prevState := slot.State
	slot.State = spec.NextState
	slot.Terminal = spec.Terminal
	slot.LastBranch = branchId
	slot.Version++
	e.states[spec.StateKey] = slot

	e.audit = append(e.audit, TransitionAuditEntry{
		StateKey:  spec.StateKey,
		BranchId:  branchId,
		FromState: prevState,
		ToState:   spec.NextState,
		Terminal:  spec.Terminal,
	})
	return nil
}

// GetState returns the current StateSlot for stateKey (defaulting to OPEN).
func (e *ExclusivityEngine) GetState(stateKey string) StateSlot {
	e.mu.Lock()
	defer e.mu.Unlock()
	slot, ok := e.states[stateKey]
	if !ok {
		return StateSlot{State: DefaultInitialState}
	}
	return slot
}

// AuditLog returns a copy of all committed transitions in execution order.
func (e *ExclusivityEngine) AuditLog() []TransitionAuditEntry {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]TransitionAuditEntry, len(e.audit))
	copy(out, e.audit)
	return out
}

// CompileExclusivityGroupToStateTransition demonstrates how Option A sugar
// compiles losslessly onto Option B's StateTransition primitive.
func CompileExclusivityGroupToStateTransition(groupId string) StateTransitionSpec {
	return StateTransitionSpec{
		StateKey:       groupId,
		RequiredStates: []string{DefaultInitialState},
		NextState:      "CLAIMED",
		Terminal:       true,
	}
}
