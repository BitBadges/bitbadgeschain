package types

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// Prototype errors for BB-379 CoinPullAllowance recurring subscriptions.
var (
	ErrAllowanceNotFound       = errors.New("coin pull allowance not found")
	ErrAllowanceAlreadyExists  = errors.New("coin pull allowance already exists")
	ErrInvalidAllowanceConfig  = errors.New("invalid coin pull allowance configuration")
	ErrAllowanceRevoked        = errors.New("coin pull allowance has been revoked")
	ErrAllowanceExpired        = errors.New("coin pull allowance has expired")
	ErrAllowanceNotYetActive   = errors.New("coin pull allowance start time not reached")
	ErrTermsMutated            = errors.New("subscription terms or approval version mutated after subscriber consent")
	ErrPeriodCapExceeded       = errors.New("requested pull exceeds per-period allowance cap")
	ErrExecutorFeeExceeded     = errors.New("executor fee exceeds subscriber maxExecutorFee")
	ErrUnauthorizedRevocation  = errors.New("only the subscriber may revoke a coin pull allowance")
	ErrDenomMismatch           = errors.New("requested pull denom does not match allowance denom")
)

// PullStatus describes the outcome of a permissionless pull attempt.
type PullStatus string

const (
	PullStatusSucceeded           PullStatus = "SUCCEEDED"
	PullStatusInsufficientBalance PullStatus = "INSUFFICIENT_BALANCE_GRACE"
)

// SubscriptionTerms captures the exact immutable merchant billing terms bound
// cryptographically into a subscriber's CoinPullAllowance at consent time.
type SubscriptionTerms struct {
	CollectionId    uint64       `json:"collectionId"`
	ApprovalId      string       `json:"approvalId"`
	ApprovalVersion uint64       `json:"approvalVersion"`
	PriceDenom      string       `json:"priceDenom"`
	PriceAmount     sdkmath.Uint `json:"priceAmount"`
	PayoutAddress   string       `json:"payoutAddress"`
	PeriodLengthMs  uint64       `json:"periodLengthMs"`
}

// Hash returns a canonical SHA-256 hex digest over all merchant billing terms.
// Any price increase, payout redirect, period shortening, or version bump
// produces a different hash and invalidates prior subscriber consent.
func (t SubscriptionTerms) Hash() string {
	priceStr := "0"
	if !t.PriceAmount.IsNil() {
		priceStr = t.PriceAmount.String()
	}
	payload := fmt.Sprintf("v1|%d|%s|%d|%s|%s|%s|%d",
		t.CollectionId,
		t.ApprovalId,
		t.ApprovalVersion,
		t.PriceDenom,
		priceStr,
		t.PayoutAddress,
		t.PeriodLengthMs,
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

// CoinPullAllowance represents a subscriber's signed, scoped, revocable
// recurring auto-debit authorization aligned with BB-146 & BB-379.
type CoinPullAllowance struct {
	AllowanceId        string       `json:"allowanceId"`
	Subscriber         string       `json:"subscriber"`
	Beneficiary        string       `json:"beneficiary"`
	CollectionId       uint64       `json:"collectionId"`
	ApprovalId         string       `json:"approvalId"`
	ApprovalVersion    uint64       `json:"approvalVersion"`
	TermsHash          string       `json:"termsHash"`
	Denom              string       `json:"denom"`
	MaxAmountPerPeriod sdkmath.Uint `json:"maxAmountPerPeriod"`
	MaxExecutorFee     sdkmath.Uint `json:"maxExecutorFee"`
	StartTimeMs        uint64       `json:"startTimeMs"`
	PeriodLengthMs     uint64       `json:"periodLengthMs"`
	ValidUntilMs       uint64       `json:"validUntilMs"`

	// Runtime accounting state
	CurrentPeriodIndex   uint64       `json:"currentPeriodIndex"`
	SpentInCurrentPeriod sdkmath.Uint `json:"spentInCurrentPeriod"`
	TotalPulled          sdkmath.Uint `json:"totalPulled"`
	LastPulledAtMs       uint64       `json:"lastPulledAtMs"`
	FailedBalanceChecks  uint64       `json:"failedBalanceChecks"`
	Revoked              bool         `json:"revoked"`
	RevokedAtMs          uint64       `json:"revokedAtMs"`
}

// PullResult summarizes a permissionless renewal pull attempt.
type PullResult struct {
	Status                PullStatus   `json:"status"`
	PeriodIndex           uint64       `json:"periodIndex"`
	ChargedAmount         sdkmath.Uint `json:"chargedAmount"`
	ExecutorFeePaid       sdkmath.Uint `json:"executorFeePaid"`
	RemainingPayerBalance sdkmath.Uint `json:"remainingPayerBalance"`
	RemainingPeriodQuota  sdkmath.Uint `json:"remainingPeriodQuota"`
}

// PullAllowanceRegistry is an in-memory reference state machine for CoinPullAllowance.
type PullAllowanceRegistry struct {
	allowances map[string]*CoinPullAllowance
	balances   map[string]map[string]sdkmath.Uint // address -> denom -> balance
}

// NewPullAllowanceRegistry constructs an empty PullAllowanceRegistry.
func NewPullAllowanceRegistry() *PullAllowanceRegistry {
	return &PullAllowanceRegistry{
		allowances: make(map[string]*CoinPullAllowance),
		balances:   make(map[string]map[string]sdkmath.Uint),
	}
}

// SetBalance sets a simulated bank balance for an address and denom.
func (r *PullAllowanceRegistry) SetBalance(addr, denom string, amount sdkmath.Uint) {
	if r.balances[addr] == nil {
		r.balances[addr] = make(map[string]sdkmath.Uint)
	}
	r.balances[addr][denom] = amount
}

// GetBalance returns the simulated bank balance for an address and denom.
func (r *PullAllowanceRegistry) GetBalance(addr, denom string) sdkmath.Uint {
	if r.balances[addr] == nil {
		return sdkmath.ZeroUint()
	}
	bal, ok := r.balances[addr][denom]
	if !ok || bal.IsNil() {
		return sdkmath.ZeroUint()
	}
	return bal
}

// GrantAllowance records initial subscriber consent pinned to terms.Hash() and terms.ApprovalVersion.
func (r *PullAllowanceRegistry) GrantAllowance(
	allowanceId string,
	subscriber string,
	beneficiary string,
	terms SubscriptionTerms,
	maxAmountPerPeriod sdkmath.Uint,
	maxExecutorFee sdkmath.Uint,
	startTimeMs uint64,
	validUntilMs uint64,
) (*CoinPullAllowance, error) {
	if allowanceId == "" || subscriber == "" || terms.ApprovalId == "" || terms.PriceDenom == "" || terms.PayoutAddress == "" {
		return nil, fmt.Errorf("%w: missing required identifier or address", ErrInvalidAllowanceConfig)
	}
	if _, exists := r.allowances[allowanceId]; exists {
		return nil, fmt.Errorf("%w: %s", ErrAllowanceAlreadyExists, allowanceId)
	}
	if terms.PeriodLengthMs == 0 {
		return nil, fmt.Errorf("%w: periodLengthMs must be > 0", ErrInvalidAllowanceConfig)
	}
	if validUntilMs <= startTimeMs {
		return nil, fmt.Errorf("%w: validUntilMs must be greater than startTimeMs", ErrInvalidAllowanceConfig)
	}
	if terms.PriceAmount.IsNil() || terms.PriceAmount.IsZero() {
		return nil, fmt.Errorf("%w: terms priceAmount must be > 0", ErrInvalidAllowanceConfig)
	}
	if maxAmountPerPeriod.IsNil() || maxAmountPerPeriod.IsZero() {
		return nil, fmt.Errorf("%w: maxAmountPerPeriod must be > 0", ErrInvalidAllowanceConfig)
	}
	if maxAmountPerPeriod.LT(terms.PriceAmount) {
		return nil, fmt.Errorf("%w: maxAmountPerPeriod (%s) cannot be less than terms priceAmount (%s)",
			ErrInvalidAllowanceConfig, maxAmountPerPeriod.String(), terms.PriceAmount.String())
	}
	if maxExecutorFee.IsNil() {
		maxExecutorFee = sdkmath.ZeroUint()
	}
	if beneficiary == "" {
		beneficiary = subscriber
	}

	a := &CoinPullAllowance{
		AllowanceId:          allowanceId,
		Subscriber:           subscriber,
		Beneficiary:          beneficiary,
		CollectionId:         terms.CollectionId,
		ApprovalId:           terms.ApprovalId,
		ApprovalVersion:      terms.ApprovalVersion,
		TermsHash:            terms.Hash(),
		Denom:                terms.PriceDenom,
		MaxAmountPerPeriod:   maxAmountPerPeriod,
		MaxExecutorFee:       maxExecutorFee,
		StartTimeMs:          startTimeMs,
		PeriodLengthMs:       terms.PeriodLengthMs,
		ValidUntilMs:         validUntilMs,
		SpentInCurrentPeriod: sdkmath.ZeroUint(),
		TotalPulled:          sdkmath.ZeroUint(),
	}
	r.allowances[allowanceId] = a
	copied := *a
	return &copied, nil
}

// RevokeAllowance immediately revokes the allowance. Only the subscriber may revoke.
func (r *PullAllowanceRegistry) RevokeAllowance(allowanceId, caller string, nowMs uint64) error {
	a, ok := r.allowances[allowanceId]
	if !ok {
		return fmt.Errorf("%w: %s", ErrAllowanceNotFound, allowanceId)
	}
	if caller != a.Subscriber {
		return fmt.Errorf("%w: caller %s != subscriber %s", ErrUnauthorizedRevocation, caller, a.Subscriber)
	}
	a.Revoked = true
	a.RevokedAtMs = nowMs
	return nil
}

// GetAllowance returns a defensive copy of the allowance state.
func (r *PullAllowanceRegistry) GetAllowance(allowanceId string) (*CoinPullAllowance, bool) {
	a, ok := r.allowances[allowanceId]
	if !ok {
		return nil, false
	}
	copied := *a
	return &copied, true
}

// ExecutePull runs a permissionless renewal charge against allowanceId using currentTerms.
// Enforces revocation, expiry, immutable terms binding, per-period hard cap, and
// graceful insufficient-balance handling (never bricking state or burning quota).
func (r *PullAllowanceRegistry) ExecutePull(
	allowanceId string,
	currentTerms SubscriptionTerms,
	executorAddress string,
	nowMs uint64,
	executorFee sdkmath.Uint,
) (PullResult, error) {
	a, ok := r.allowances[allowanceId]
	if !ok {
		return PullResult{}, fmt.Errorf("%w: %s", ErrAllowanceNotFound, allowanceId)
	}
	if a.Revoked {
		return PullResult{}, fmt.Errorf("%w: revoked at %d", ErrAllowanceRevoked, a.RevokedAtMs)
	}
	if nowMs < a.StartTimeMs {
		return PullResult{}, fmt.Errorf("%w: now=%d < start=%d", ErrAllowanceNotYetActive, nowMs, a.StartTimeMs)
	}
	if nowMs > a.ValidUntilMs {
		return PullResult{}, fmt.Errorf("%w: now=%d > validUntil=%d", ErrAllowanceExpired, nowMs, a.ValidUntilMs)
	}
	if currentTerms.PriceAmount.IsNil() || currentTerms.PriceAmount.IsZero() {
		return PullResult{}, fmt.Errorf("%w: currentTerms priceAmount must be > 0", ErrInvalidAllowanceConfig)
	}
	if a.LastPulledAtMs > 0 && nowMs < a.LastPulledAtMs {
		return PullResult{}, fmt.Errorf("%w: timestamp %d precedes last pull %d", ErrInvalidAllowanceConfig, nowMs, a.LastPulledAtMs)
	}
	if currentTerms.PriceDenom != a.Denom {
		return PullResult{}, fmt.Errorf("%w: got %s, expected %s", ErrDenomMismatch, currentTerms.PriceDenom, a.Denom)
	}
	if currentTerms.ApprovalVersion != a.ApprovalVersion || currentTerms.Hash() != a.TermsHash {
		return PullResult{}, fmt.Errorf("%w: expected version=%d hash=%s, got version=%d hash=%s",
			ErrTermsMutated, a.ApprovalVersion, a.TermsHash, currentTerms.ApprovalVersion, currentTerms.Hash())
	}
	if executorFee.IsNil() {
		executorFee = sdkmath.ZeroUint()
	}
	if executorFee.GT(a.MaxExecutorFee) {
		return PullResult{}, fmt.Errorf("%w: requested %s > max %s", ErrExecutorFeeExceeded, executorFee.String(), a.MaxExecutorFee.String())
	}

	// Determine deterministic period index and roll over quota on new period.
	periodIdx := (nowMs - a.StartTimeMs) / a.PeriodLengthMs
	spentInPeriod := a.SpentInCurrentPeriod
	if periodIdx > a.CurrentPeriodIndex || a.LastPulledAtMs == 0 {
		spentInPeriod = sdkmath.ZeroUint()
	}

	totalDebit := currentTerms.PriceAmount.Add(executorFee)
	projectedSpend := spentInPeriod.Add(totalDebit)
	if projectedSpend.GT(a.MaxAmountPerPeriod) {
		return PullResult{}, fmt.Errorf("%w: spent=%s + debit=%s > cap=%s in period %d",
			ErrPeriodCapExceeded, spentInPeriod.String(), totalDebit.String(), a.MaxAmountPerPeriod.String(), periodIdx)
	}

	// Graceful insufficient-balance check: do NOT consume quota or brick allowance.
	payerBal := r.GetBalance(a.Subscriber, a.Denom)
	if payerBal.LT(totalDebit) {
		a.FailedBalanceChecks++
		return PullResult{
			Status:                PullStatusInsufficientBalance,
			PeriodIndex:           periodIdx,
			ChargedAmount:         sdkmath.ZeroUint(),
			ExecutorFeePaid:       sdkmath.ZeroUint(),
			RemainingPayerBalance: payerBal,
			RemainingPeriodQuota:  a.MaxAmountPerPeriod.Sub(spentInPeriod),
		}, nil
	}

	// Commit atomic balance transfers and period accounting.
	newPayerBal := payerBal.Sub(totalDebit)
	r.SetBalance(a.Subscriber, a.Denom, newPayerBal)

	payoutBal := r.GetBalance(currentTerms.PayoutAddress, a.Denom)
	r.SetBalance(currentTerms.PayoutAddress, a.Denom, payoutBal.Add(currentTerms.PriceAmount))

	if !executorFee.IsZero() && executorAddress != "" {
		execBal := r.GetBalance(executorAddress, a.Denom)
		r.SetBalance(executorAddress, a.Denom, execBal.Add(executorFee))
	}

	a.CurrentPeriodIndex = periodIdx
	a.SpentInCurrentPeriod = projectedSpend
	a.TotalPulled = a.TotalPulled.Add(totalDebit)
	a.LastPulledAtMs = nowMs

	return PullResult{
		Status:                PullStatusSucceeded,
		PeriodIndex:           periodIdx,
		ChargedAmount:         currentTerms.PriceAmount,
		ExecutorFeePaid:       executorFee,
		RemainingPayerBalance: newPayerBal,
		RemainingPeriodQuota:  a.MaxAmountPerPeriod.Sub(projectedSpend),
	}, nil
}
