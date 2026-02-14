package gamm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Metrics tracks usage statistics for the precompile
type Metrics struct {
	// Transaction counts
	JoinPoolCount                      uint64
	ExitPoolCount                      uint64
	SwapExactAmountInCount            uint64
	SwapExactAmountInWithIBCTransferCount uint64

	// Query counts
	QueryCount uint64

	// Error counts by type
	ErrorCounts map[ErrorCode]uint64

	// Gas consumption
	TotalGasConsumed uint64
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		ErrorCounts: make(map[ErrorCode]uint64),
	}
}

// IncrementJoinPool increments the join pool counter
func (m *Metrics) IncrementJoinPool() {
	m.JoinPoolCount++
}

// IncrementExitPool increments the exit pool counter
func (m *Metrics) IncrementExitPool() {
	m.ExitPoolCount++
}

// IncrementSwap increments the swap counter
func (m *Metrics) IncrementSwap() {
	m.SwapExactAmountInCount++
}

// IncrementSwapWithIBC increments the swap with IBC transfer counter
func (m *Metrics) IncrementSwapWithIBC() {
	m.SwapExactAmountInWithIBCTransferCount++
}

// IncrementQuery increments the query counter
func (m *Metrics) IncrementQuery() {
	m.QueryCount++
}

// IncrementError increments the error counter for a specific error code
func (m *Metrics) IncrementError(code ErrorCode) {
	m.ErrorCounts[code]++
}

// AddGas adds gas consumption to the total
func (m *Metrics) AddGas(amount uint64) {
	m.TotalGasConsumed += amount
}

// GetMetrics returns current metrics (for monitoring/observability)
// In a production system, this would integrate with a metrics collection system
func GetMetrics(ctx sdk.Context) *Metrics {
	// For now, return a new instance
	// In production, this would retrieve from context or a metrics store
	return NewMetrics()
}

// EmitMetricsEvent emits a metrics event for observability
func EmitMetricsEvent(ctx sdk.Context, method string, success bool, gasUsed uint64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_metrics",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("method", method),
			sdk.NewAttribute("success", fmt.Sprintf("%t", success)),
			sdk.NewAttribute("gas_used", fmt.Sprintf("%d", gasUsed)),
		),
	)
}

