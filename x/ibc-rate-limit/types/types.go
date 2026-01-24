package types

import (
	"fmt"

	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		RateLimits: []RateLimitConfig{}, // Empty by default (no rate limits)
	}
}

// Validate validates params
func (p Params) Validate() error {
	for i, config := range p.RateLimits {
		if err := config.Validate(); err != nil {
			return fmt.Errorf("invalid rate limit config at index %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates a rate limit config
func (c RateLimitConfig) Validate() error {
	// Denom must always be specified (no empty denoms allowed)
	if c.Denom == "" {
		return fmt.Errorf("denom must be specified (empty denoms are not allowed)")
	}

	// Validate supply shift limits
	for i, limit := range c.SupplyShiftLimits {
		if err := limit.Validate(); err != nil {
			return fmt.Errorf("invalid supply_shift_limits[%d]: %w", i, err)
		}
	}

	// Validate unique sender limits
	for i, limit := range c.UniqueSenderLimits {
		if err := limit.Validate(); err != nil {
			return fmt.Errorf("invalid unique_sender_limits[%d]: %w", i, err)
		}
	}

	// Validate address limits
	for i, limit := range c.AddressLimits {
		if err := limit.Validate(); err != nil {
			return fmt.Errorf("invalid address_limits[%d]: %w", i, err)
		}
	}

	return nil
}

// Validate validates a timeframe limit
func (t TimeframeLimit) Validate() error {
	if t.MaxAmount.IsNegative() {
		return fmt.Errorf("max_amount cannot be negative: %s", t.MaxAmount)
	}
	if t.TimeframeType == TimeframeType_TIMEFRAME_TYPE_UNSPECIFIED {
		return fmt.Errorf("timeframe_type must be specified")
	}
	if t.TimeframeDuration <= 0 {
		return fmt.Errorf("timeframe_duration must be positive: %d", t.TimeframeDuration)
	}
	return nil
}

// Validate validates a unique sender limit
func (u UniqueSenderLimit) Validate() error {
	if u.MaxUniqueSenders < 0 {
		return fmt.Errorf("max_unique_senders cannot be negative: %d", u.MaxUniqueSenders)
	}
	if u.TimeframeType == TimeframeType_TIMEFRAME_TYPE_UNSPECIFIED {
		return fmt.Errorf("timeframe_type must be specified")
	}
	if u.TimeframeDuration <= 0 {
		return fmt.Errorf("timeframe_duration must be positive: %d", u.TimeframeDuration)
	}
	return nil
}

// Validate validates an address limit
func (a AddressLimit) Validate() error {
	if a.MaxTransfers < 0 {
		return fmt.Errorf("max_transfers cannot be negative: %d", a.MaxTransfers)
	}
	if a.MaxAmount.IsNegative() {
		return fmt.Errorf("max_amount cannot be negative: %s", a.MaxAmount)
	}
	if a.TimeframeType == TimeframeType_TIMEFRAME_TYPE_UNSPECIFIED {
		return fmt.Errorf("timeframe_type must be specified")
	}
	if a.TimeframeDuration <= 0 {
		return fmt.Errorf("timeframe_duration must be positive: %d", a.TimeframeDuration)
	}
	return nil
}

// TimeframeDurationInBlocks converts a timeframe duration to blocks
// blockTime is the average block time in seconds (e.g., 6 seconds)
func TimeframeDurationInBlocks(timeframeType TimeframeType, timeframeDuration int64, blockTimeSeconds int64) int64 {
	switch timeframeType {
	case TimeframeType_TIMEFRAME_TYPE_BLOCK:
		return timeframeDuration
	case TimeframeType_TIMEFRAME_TYPE_HOUR:
		// Convert hours to seconds, then to blocks
		seconds := timeframeDuration * 3600
		return seconds / blockTimeSeconds
	case TimeframeType_TIMEFRAME_TYPE_DAY:
		// Convert days to seconds, then to blocks
		seconds := timeframeDuration * 24 * 3600
		return seconds / blockTimeSeconds
	default:
		return timeframeDuration // Fallback to treating as blocks
	}
}

// FindMatchingConfig finds the first rate limit config that matches the given channel and denom
// Returns nil if no config matches
// Matching rules:
// - Empty channel_id matches all channels
// - Denom must exactly match (empty denoms are not allowed in configs)
// - Both must match (channel can be empty to match all, but denom must match exactly)
func (p Params) FindMatchingConfig(channelID, denom string) *RateLimitConfig {
	for i := range p.RateLimits {
		config := &p.RateLimits[i]

		// Check channel match (empty channel_id matches all channels)
		channelMatch := config.ChannelId == "" || config.ChannelId == channelID

		// Check denom match (must match exactly, empty denoms are not allowed)
		denomMatch := config.Denom == denom

		if channelMatch && denomMatch {
			return config
		}
	}
	return nil
}

// NewCustomErrorAcknowledgement creates a custom error acknowledgement with a deterministic error string
// IMPORTANT: The error string must be deterministic (no traces, logs, or non-deterministic values)
// This is used instead of channeltypes.NewErrorAcknowledgement to provide more friendly error messages
func NewCustomErrorAcknowledgement(errorMsg string) ibcexported.Acknowledgement {
	return channeltypes.Acknowledgement{
		Response: &channeltypes.Acknowledgement_Error{
			Error: fmt.Sprintf("ibc-rate-limit: %s", errorMsg),
		},
	}
}

// NewSuccessAcknowledgement creates a success acknowledgement
// This is used internally to indicate successful execution
func NewSuccessAcknowledgement() ibcexported.Acknowledgement {
	return channeltypes.NewResultAcknowledgement([]byte("success"))
}

// IsSuccessAcknowledgement checks if an acknowledgement indicates success
func IsSuccessAcknowledgement(ack ibcexported.Acknowledgement) bool {
	return ack.Success()
}
