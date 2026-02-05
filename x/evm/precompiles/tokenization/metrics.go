package tokenization

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Metrics tracks usage statistics for the precompile
type Metrics struct {
	// Transaction counts
	TransferCount      uint64
	IncomingApprovalCount uint64
	OutgoingApprovalCount uint64

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

// IncrementTransfer increments the transfer counter
func (m *Metrics) IncrementTransfer() {
	m.TransferCount++
}

// IncrementIncomingApproval increments the incoming approval counter
func (m *Metrics) IncrementIncomingApproval() {
	m.IncomingApprovalCount++
}

// IncrementOutgoingApproval increments the outgoing approval counter
func (m *Metrics) IncrementOutgoingApproval() {
	m.OutgoingApprovalCount++
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

// LogPrecompileUsage logs precompile usage for monitoring
func LogPrecompileUsage(ctx sdk.Context, method string, success bool, gasUsed uint64, err error) {
	logger := ctx.Logger()
	
	if err != nil {
		// Extract error code if it's a PrecompileError
		if precompileErr, ok := err.(*PrecompileError); ok {
			logger.Error("precompile error",
				"method", method,
				"error_code", precompileErr.Code,
				"error_message", precompileErr.Message,
				"gas_used", gasUsed,
			)
		} else {
			logger.Error("precompile error",
				"method", method,
				"error", err.Error(),
				"gas_used", gasUsed,
			)
		}
	} else {
		logger.Info("precompile success",
			"method", method,
			"gas_used", gasUsed,
		)
	}
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

