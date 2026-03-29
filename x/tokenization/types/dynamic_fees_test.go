package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	// Ensure the SDK bech32 prefix is set to "bb" for address validation
	params.SetAddressPrefixes()
}

// testAddr is a valid bb-prefixed address for use in tests.
var testAddr = func() string {
	params.SetAddressPrefixes()
	return sdk.AccAddress(make([]byte, 20)).String()
}()

// --- Feature B: Dynamic Fee Schedule Tests ---

func TestComputeDynamicFee_SingleTier(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{
				MinAmount:   sdkmath.NewUint(0),
				MaxAmount:   sdkmath.NewUint(1000000),
				BasisPoints: 50, // 0.5%
			},
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	// 10000 tokens * 50bp / 10000 = 50
	fee, found := schedule.ComputeDynamicFee(sdkmath.NewUint(10000))
	if !found {
		t.Fatal("expected to find matching tier")
	}
	if !fee.Equal(sdkmath.NewInt(50)) {
		t.Fatalf("expected fee 50, got %s", fee.String())
	}
}

func TestComputeDynamicFee_TieredBrackets(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(999), BasisPoints: 100},     // 1%
			{MinAmount: sdkmath.NewUint(1000), MaxAmount: sdkmath.NewUint(9999), BasisPoints: 50},   // 0.5%
			{MinAmount: sdkmath.NewUint(10000), MaxAmount: sdkmath.NewUint(999999), BasisPoints: 25}, // 0.25%
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	tests := []struct {
		name           string
		amount         sdkmath.Uint
		expectedFee    sdkmath.Int
		expectedFound  bool
	}{
		{"small transfer 500 at 1%", sdkmath.NewUint(500), sdkmath.NewInt(5), true},
		{"medium transfer 5000 at 0.5%", sdkmath.NewUint(5000), sdkmath.NewInt(25), true},
		{"large transfer 50000 at 0.25%", sdkmath.NewUint(50000), sdkmath.NewInt(125), true},
		{"boundary lower 0 at 1%", sdkmath.NewUint(0), sdkmath.NewInt(0), true},
		{"boundary upper 999 at 1%", sdkmath.NewUint(999), sdkmath.NewInt(9), true},
		{"boundary lower 1000 at 0.5%", sdkmath.NewUint(1000), sdkmath.NewInt(5), true},
		{"no matching tier", sdkmath.NewUint(1000000), sdkmath.ZeroInt(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee, found := schedule.ComputeDynamicFee(tt.amount)
			if found != tt.expectedFound {
				t.Fatalf("expected found=%v, got %v", tt.expectedFound, found)
			}
			if found && !fee.Equal(tt.expectedFee) {
				t.Fatalf("expected fee %s, got %s", tt.expectedFee.String(), fee.String())
			}
		})
	}
}

func TestComputeDynamicFee_ZeroBasisPoints(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(1000000), BasisPoints: 0},
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	fee, found := schedule.ComputeDynamicFee(sdkmath.NewUint(10000))
	if !found {
		t.Fatal("expected to find matching tier")
	}
	if !fee.IsZero() {
		t.Fatalf("expected zero fee, got %s", fee.String())
	}
}

func TestComputeDynamicFee_MaxBasisPoints(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(1000000), BasisPoints: 10000}, // 100%
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	fee, found := schedule.ComputeDynamicFee(sdkmath.NewUint(5000))
	if !found {
		t.Fatal("expected to find matching tier")
	}
	if !fee.Equal(sdkmath.NewInt(5000)) {
		t.Fatalf("expected fee 5000 (100%%), got %s", fee.String())
	}
}

// --- Feature A: Metered Balances (CoinPerTokenMultiplier) Tests ---

func TestComputeMeteredAmount_SingleCoin(t *testing.T) {
	m := &CoinPerTokenMultiplier{
		CoinAmountPerToken: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(100)), // 100 ubadge per token
		},
	}

	result := m.ComputeMeteredAmount(sdkmath.NewUint(50)) // 50 tokens
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	if !result[0].Amount.Equal(sdkmath.NewInt(5000)) {
		t.Fatalf("expected 5000 ubadge, got %s", result[0].Amount.String())
	}
	if result[0].Denom != "ubadge" {
		t.Fatalf("expected denom ubadge, got %s", result[0].Denom)
	}
}

func TestComputeMeteredAmount_MultipleCoins(t *testing.T) {
	m := &CoinPerTokenMultiplier{
		CoinAmountPerToken: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(100)),
			sdk.NewCoin("uatom", sdkmath.NewInt(50)),
		},
	}

	result := m.ComputeMeteredAmount(sdkmath.NewUint(10))
	if len(result) != 2 {
		t.Fatalf("expected 2 coins, got %d", len(result))
	}
	if !result[0].Amount.Equal(sdkmath.NewInt(1000)) {
		t.Fatalf("expected 1000 ubadge, got %s", result[0].Amount.String())
	}
	if !result[1].Amount.Equal(sdkmath.NewInt(500)) {
		t.Fatalf("expected 500 uatom, got %s", result[1].Amount.String())
	}
}

func TestComputeMeteredAmount_ZeroTokens(t *testing.T) {
	m := &CoinPerTokenMultiplier{
		CoinAmountPerToken: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(100)),
		},
	}

	result := m.ComputeMeteredAmount(sdkmath.NewUint(0))
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	if !result[0].Amount.IsZero() {
		t.Fatalf("expected 0, got %s", result[0].Amount.String())
	}
}

// --- Feature C: Time-Based Refund Tests ---

func TestComputeRefundAmount_FullRefund(t *testing.T) {
	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(30 * 24 * 60 * 60 * 1000), // 30 days in ms
		EndTime:       sdkmath.NewUint(1000000 + 30*24*60*60*1000),
	}

	// Current time = start time (endTime - totalDuration), full refund expected
	currentTime := sdkmath.NewUint(1000000)
	result := formula.ComputeRefundAmount(currentTime)
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	if !result[0].Amount.Equal(sdkmath.NewInt(10000)) {
		t.Fatalf("expected full refund 10000, got %s", result[0].Amount.String())
	}
}

func TestComputeRefundAmount_HalfRefund(t *testing.T) {
	totalDuration := uint64(30 * 24 * 60 * 60 * 1000) // 30 days in ms
	startTime := uint64(1000000)
	endTime := startTime + totalDuration

	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(totalDuration),
		EndTime:       sdkmath.NewUint(endTime),
	}

	// Current time = halfway through
	halfwayTime := sdkmath.NewUint(startTime + totalDuration/2)
	result := formula.ComputeRefundAmount(halfwayTime)
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	if !result[0].Amount.Equal(sdkmath.NewInt(5000)) {
		t.Fatalf("expected half refund 5000, got %s", result[0].Amount.String())
	}
}

func TestComputeRefundAmount_Expired(t *testing.T) {
	totalDuration := uint64(30 * 24 * 60 * 60 * 1000)
	startTime := uint64(1000000)
	endTime := startTime + totalDuration

	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(totalDuration),
		EndTime:       sdkmath.NewUint(endTime),
	}

	// Current time = after end time
	result := formula.ComputeRefundAmount(sdkmath.NewUint(endTime + 1000))
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	if !result[0].Amount.IsZero() {
		t.Fatalf("expected zero refund, got %s", result[0].Amount.String())
	}
}

func TestComputeRefundAmount_NearEnd(t *testing.T) {
	totalDuration := uint64(100000) // 100 seconds in ms
	startTime := uint64(1000000)
	endTime := startTime + totalDuration

	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(totalDuration),
		EndTime:       sdkmath.NewUint(endTime),
	}

	// 10% time remaining
	currentTime := sdkmath.NewUint(endTime - 10000)
	result := formula.ComputeRefundAmount(currentTime)
	if len(result) != 1 {
		t.Fatalf("expected 1 coin, got %d", len(result))
	}
	// 10000 * 10000 / 100000 = 1000
	if !result[0].Amount.Equal(sdkmath.NewInt(1000)) {
		t.Fatalf("expected refund 1000 (10%%), got %s", result[0].Amount.String())
	}
}

func TestComputeRefundAmount_MultipleCoins(t *testing.T) {
	totalDuration := uint64(100000)
	startTime := uint64(1000000)
	endTime := startTime + totalDuration

	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
			sdk.NewCoin("uatom", sdkmath.NewInt(5000)),
		},
		TotalDuration: sdkmath.NewUint(totalDuration),
		EndTime:       sdkmath.NewUint(endTime),
	}

	// 50% time remaining
	currentTime := sdkmath.NewUint(startTime + totalDuration/2)
	result := formula.ComputeRefundAmount(currentTime)
	if len(result) != 2 {
		t.Fatalf("expected 2 coins, got %d", len(result))
	}
	if !result[0].Amount.Equal(sdkmath.NewInt(5000)) {
		t.Fatalf("expected refund 5000 ubadge, got %s", result[0].Amount.String())
	}
	if !result[1].Amount.Equal(sdkmath.NewInt(2500)) {
		t.Fatalf("expected refund 2500 uatom, got %s", result[1].Amount.String())
	}
}

// --- GetTotalTokenTransferAmount Tests ---

func TestGetTotalTokenTransferAmount_SingleBalance(t *testing.T) {
	balances := []*Balance{
		{
			Amount:         sdkmath.NewUint(100),
			TokenIds:       []*UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*UintRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(999999999999)}},
		},
	}

	total := GetTotalTokenTransferAmount(balances)
	if !total.Equal(sdkmath.NewUint(100)) {
		t.Fatalf("expected 100, got %s", total.String())
	}
}

func TestGetTotalTokenTransferAmount_MultipleBalances(t *testing.T) {
	balances := []*Balance{
		{Amount: sdkmath.NewUint(100)},
		{Amount: sdkmath.NewUint(200)},
		{Amount: sdkmath.NewUint(50)},
	}

	total := GetTotalTokenTransferAmount(balances)
	if !total.Equal(sdkmath.NewUint(350)) {
		t.Fatalf("expected 350, got %s", total.String())
	}
}

func TestGetTotalTokenTransferAmount_Empty(t *testing.T) {
	total := GetTotalTokenTransferAmount([]*Balance{})
	if !total.IsZero() {
		t.Fatalf("expected 0, got %s", total.String())
	}
}

func TestGetTotalTokenTransferAmount_NilBalances(t *testing.T) {
	balances := []*Balance{nil, nil}
	total := GetTotalTokenTransferAmount(balances)
	if !total.IsZero() {
		t.Fatalf("expected 0, got %s", total.String())
	}
}

// --- Validation Tests ---

func TestValidateDynamicFeeSchedule_Valid(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(999), BasisPoints: 100},
			{MinAmount: sdkmath.NewUint(1000), MaxAmount: sdkmath.NewUint(9999), BasisPoints: 50},
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	if err := ValidateDynamicFeeSchedule(schedule); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDynamicFeeSchedule_Nil(t *testing.T) {
	if err := ValidateDynamicFeeSchedule(nil); err != nil {
		t.Fatalf("nil schedule should be valid, got: %v", err)
	}
}

func TestValidateDynamicFeeSchedule_EmptyTiers(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers:        []FeeTier{},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	if err := ValidateDynamicFeeSchedule(schedule); err == nil {
		t.Fatal("expected error for empty tiers")
	}
}

func TestValidateDynamicFeeSchedule_OverlappingTiers(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(1000), BasisPoints: 100},
			{MinAmount: sdkmath.NewUint(500), MaxAmount: sdkmath.NewUint(2000), BasisPoints: 50}, // overlaps
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	if err := ValidateDynamicFeeSchedule(schedule); err == nil {
		t.Fatal("expected error for overlapping tiers")
	}
}

func TestValidateDynamicFeeSchedule_BasisPointsTooHigh(t *testing.T) {
	schedule := &DynamicFeeSchedule{
		Tiers: []FeeTier{
			{MinAmount: sdkmath.NewUint(0), MaxAmount: sdkmath.NewUint(1000), BasisPoints: 10001},
		},
		FeeRecipient: testAddr,
		FeeDenom:     "ubadge",
	}

	if err := ValidateDynamicFeeSchedule(schedule); err == nil {
		t.Fatal("expected error for basis points > 10000")
	}
}

func TestValidateCoinPerTokenMultiplier_Valid(t *testing.T) {
	m := &CoinPerTokenMultiplier{
		CoinAmountPerToken: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(100)),
		},
	}
	if err := ValidateCoinPerTokenMultiplier(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCoinPerTokenMultiplier_EmptyCoins(t *testing.T) {
	m := &CoinPerTokenMultiplier{
		CoinAmountPerToken: []sdk.Coin{},
	}
	if err := ValidateCoinPerTokenMultiplier(m); err == nil {
		t.Fatal("expected error for empty coins")
	}
}

func TestValidateTimeBasedRefundFormula_Valid(t *testing.T) {
	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(86400000),
		EndTime:       sdkmath.NewUint(1000000 + 86400000),
	}
	if err := ValidateTimeBasedRefundFormula(formula); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateTimeBasedRefundFormula_ZeroDuration(t *testing.T) {
	formula := &TimeBasedRefundFormula{
		BaseRefundAmount: []sdk.Coin{
			sdk.NewCoin("ubadge", sdkmath.NewInt(10000)),
		},
		TotalDuration: sdkmath.NewUint(0),
		EndTime:       sdkmath.NewUint(1000000),
	}
	if err := ValidateTimeBasedRefundFormula(formula); err == nil {
		t.Fatal("expected error for zero duration")
	}
}
