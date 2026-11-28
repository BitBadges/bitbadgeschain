package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
)

// Canonical BitBadges standard identifiers (BB-374).
// All existing standards are intentionally preserved with NO consolidation.
const (
	StandardNonTransferable  = "Non-Transferable"
	StandardSubscriptions    = "Subscriptions"
	StandardNFTs             = "NFTs"
	StandardFungibleTokens   = "Fungible Tokens"
	StandardSmartToken       = "Smart Token"
	StandardAddressList      = "Address List"
	StandardVault            = "Vault"
	StandardAgentWallet      = "Agent Wallet"
	StandardPaymentRequestV2 = "PaymentRequestV2"
	StandardBounty           = "Bounty"
	StandardCrowdfund        = "Crowdfund"
	StandardSpendableCredit  = "Spendable Credit"
)

// StandardViolation captures a deterministic conformance failure for a claimed standard.
type StandardViolation struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// StandardVerificationResult separates untrusted creator-declared ClaimedStandards
// from deterministically verified VerifiedStandards and structured FailedStandards.
type StandardVerificationResult struct {
	ClaimedStandards  []string                       `json:"claimedStandards"`
	VerifiedStandards []string                       `json:"verifiedStandards"`
	FailedStandards   map[string][]StandardViolation `json:"failedStandards"`
}

// KnownVerifiableStandards returns the canonical ordered list of all standards
// supported by the deterministic on-chain verifier.
func KnownVerifiableStandards() []string {
	return []string{
		StandardNonTransferable,
		StandardSubscriptions,
		StandardNFTs,
		StandardFungibleTokens,
		StandardSmartToken,
		StandardAddressList,
		StandardVault,
		StandardAgentWallet,
		StandardPaymentRequestV2,
		StandardBounty,
		StandardCrowdfund,
		StandardSpendableCredit,
	}
}

// VerifyCollectionStandards deterministically verifies every standard claimed in
// collection.Standards without mutating the collection or reading external store state.
func VerifyCollectionStandards(collection *TokenCollection) StandardVerificationResult {
	res := StandardVerificationResult{
		ClaimedStandards:  []string{},
		VerifiedStandards: []string{},
		FailedStandards:   make(map[string][]StandardViolation),
	}
	if collection == nil {
		return res
	}

	seen := make(map[string]bool, len(collection.Standards))
	for _, raw := range collection.Standards {
		std := strings.TrimSpace(raw)
		if std == "" || seen[std] {
			continue
		}
		seen[std] = true
		res.ClaimedStandards = append(res.ClaimedStandards, std)

		violations := VerifyStandard(collection, std)
		if len(violations) == 0 {
			res.VerifiedStandards = append(res.VerifiedStandards, std)
		} else {
			res.FailedStandards[std] = violations
		}
	}

	return res
}

// VerifyStandard checks whether collection satisfies the deterministic invariants
// for a single named standard. Returns nil/empty slice on conformance.
func VerifyStandard(collection *TokenCollection, standard string) []StandardViolation {
	if collection == nil {
		return []StandardViolation{{
			Code:    "NIL_COLLECTION",
			Message: "collection cannot be nil",
		}}
	}

	switch standard {
	case StandardNonTransferable:
		return verifyNonTransferable(collection)
	case StandardSubscriptions:
		return verifySubscriptions(collection)
	case StandardNFTs:
		return verifyNFTs(collection)
	case StandardFungibleTokens:
		return verifyFungibleTokens(collection)
	case StandardSmartToken:
		return verifySmartToken(collection)
	case StandardAddressList:
		return verifyAddressList(collection)
	case StandardVault:
		return verifyVault(collection)
	case StandardAgentWallet:
		return verifyAgentWallet(collection)
	case StandardPaymentRequestV2:
		return verifyPaymentRequestV2(collection)
	case StandardBounty:
		return verifyBounty(collection)
	case StandardCrowdfund:
		return verifyCrowdfund(collection)
	case StandardSpendableCredit:
		return verifySpendableCredit(collection)
	default:
		return []StandardViolation{{
			Code:    "UNKNOWN_STANDARD",
			Message: fmt.Sprintf("standard %q has no deterministic verifier", standard),
		}}
	}
}

func verifyNonTransferable(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	for _, approval := range c.CollectionApprovals {
		if approval == nil {
			continue
		}
		if approval.FromListId != MintAddress && approvalIsActive(approval) {
			violations = append(violations, StandardViolation{
				Code:    "NON_TRANSFERABLE_POST_MINT_APPROVAL",
				Message: fmt.Sprintf("approval %q allows post-mint transfers (fromListId=%q)", approval.ApprovalId, approval.FromListId),
			})
		}
	}

	if !arePostMintCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "NON_TRANSFERABLE_UNFROZEN_APPROVALS",
			Message: "collectionPermissions.canUpdateCollectionApprovals does not permanently forbid updating post-mint approvals",
		})
	}

	return violations
}

func verifySubscriptions(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	if len(c.ValidTokenIds) == 0 {
		violations = append(violations, StandardViolation{
			Code:    "SUBSCRIPTION_NO_VALID_TOKEN_IDS",
			Message: "subscription collection must define validTokenIds",
		})
	}
	if c.Invariants != nil && c.Invariants.NoCustomOwnershipTimes {
		violations = append(violations, StandardViolation{
			Code:    "SUBSCRIPTION_NO_CUSTOM_OWNERSHIP_TIMES",
			Message: "subscription collection cannot enable noCustomOwnershipTimes invariant",
		})
	}

	hasTimeBoundedApproval := false
	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		pb := approval.ApprovalCriteria.PredeterminedBalances
		if pb == nil || pb.IncrementedBalances == nil {
			continue
		}
		inc := pb.IncrementedBalances
		hasDuration := !inc.DurationFromTimestamp.IsNil() && !inc.DurationFromTimestamp.IsZero()
		hasRecurring := inc.RecurringOwnershipTimes != nil &&
			!inc.RecurringOwnershipTimes.IntervalLength.IsNil() &&
			!inc.RecurringOwnershipTimes.IntervalLength.IsZero()
		if hasDuration || hasRecurring {
			hasTimeBoundedApproval = true
			break
		}
	}

	if !hasTimeBoundedApproval {
		violations = append(violations, StandardViolation{
			Code:    "SUBSCRIPTION_MISSING_TIME_BOUNDED_APPROVAL",
			Message: "subscription collection must have at least one active approval with durationFromTimestamp > 0 or recurringOwnershipTimes",
		})
	}

	return violations
}

func verifyNFTs(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	if len(c.ValidTokenIds) == 0 {
		violations = append(violations, StandardViolation{
			Code:    "NFT_NO_VALID_TOKEN_IDS",
			Message: "NFT collection must define validTokenIds",
		})
	}
	if c.Invariants == nil || !c.Invariants.MaxSupplyPerId.Equal(sdkmath.NewUint(1)) {
		violations = append(violations, StandardViolation{
			Code:    "NFT_MAX_SUPPLY_NOT_ONE",
			Message: "NFT collection must enforce invariants.maxSupplyPerId == 1",
		})
	}
	if c.Invariants == nil || !c.Invariants.NoCustomOwnershipTimes {
		violations = append(violations, StandardViolation{
			Code:    "NFT_CUSTOM_OWNERSHIP_TIMES_ALLOWED",
			Message: "NFT collection must enforce invariants.noCustomOwnershipTimes == true",
		})
	}

	return violations
}

func verifyFungibleTokens(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	if len(c.ValidTokenIds) != 1 || c.ValidTokenIds[0] == nil ||
		c.ValidTokenIds[0].Start.IsNil() || c.ValidTokenIds[0].End.IsNil() ||
		!c.ValidTokenIds[0].Start.Equal(c.ValidTokenIds[0].End) {
		violations = append(violations, StandardViolation{
			Code:    "FUNGIBLE_REQUIRES_SINGLE_TOKEN_ID",
			Message: "Fungible Tokens collection must define exactly one fungible token ID (start == end)",
		})
	}
	if c.Invariants == nil || !c.Invariants.NoCustomOwnershipTimes {
		violations = append(violations, StandardViolation{
			Code:    "FUNGIBLE_CUSTOM_OWNERSHIP_TIMES_ALLOWED",
			Message: "Fungible Tokens collection must enforce invariants.noCustomOwnershipTimes == true",
		})
	}
	if c.Invariants != nil && c.Invariants.MaxSupplyPerId.Equal(sdkmath.NewUint(1)) {
		violations = append(violations, StandardViolation{
			Code:    "FUNGIBLE_MAX_SUPPLY_IS_ONE",
			Message: "Fungible Tokens collection cannot have maxSupplyPerId == 1",
		})
	}

	return violations
}

func verifySmartToken(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	backingAddr, ok := getValidBackingAddress(c)
	if !ok {
		violations = append(violations, StandardViolation{
			Code:    "SMART_TOKEN_MISSING_BACKED_PATH",
			Message: "Smart Token collection must define a valid invariants.cosmosCoinBackedPath with address and conversion",
		})
	}

	hasBackedApproval := false
	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) {
			continue
		}
		if approval.FromListId == MintAddress {
			violations = append(violations, StandardViolation{
				Code:    "SMART_TOKEN_UNBACKED_MINT_APPROVAL",
				Message: fmt.Sprintf("Smart Token approval %q cannot mint from Mint address", approval.ApprovalId),
			})
		}
		if ok && approval.ApprovalCriteria != nil &&
			approval.ApprovalCriteria.AllowBackedMinting &&
			approval.ApprovalCriteria.MustPrioritize &&
			(approval.FromListId == backingAddr || approval.ToListId == backingAddr) {
			hasBackedApproval = true
		}
	}

	if ok && !hasBackedApproval {
		violations = append(violations, StandardViolation{
			Code:    "SMART_TOKEN_MISSING_BACKED_APPROVAL",
			Message: "Smart Token collection must have at least one prioritized allowBackedMinting approval referencing the backing address",
		})
	}

	return violations
}

func verifyAddressList(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	if len(c.ValidTokenIds) == 0 {
		violations = append(violations, StandardViolation{
			Code:    "ADDRESS_LIST_NO_VALID_TOKEN_IDS",
			Message: "Address List collection must define validTokenIds",
		})
	}

	hasScopedList := false
	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) {
			continue
		}
		if isConcreteAddressListId(approval.FromListId) ||
			isConcreteAddressListId(approval.ToListId) ||
			isConcreteAddressListId(approval.InitiatedByListId) {
			hasScopedList = true
			break
		}
	}

	if !hasScopedList {
		violations = append(violations, StandardViolation{
			Code:    "ADDRESS_LIST_NO_SCOPED_PREDICATE",
			Message: "Address List collection must have at least one active approval scoped to a non-wildcard address list",
		})
	}

	return violations
}

func verifyVault(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	backingAddr, ok := getValidBackingAddress(c)
	if !ok {
		violations = append(violations, StandardViolation{
			Code:    "VAULT_MISSING_BACKED_PATH",
			Message: "Vault collection must define a valid invariants.cosmosCoinBackedPath",
		})
	}

	hasDeposit := false
	hasWithdraw := false
	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		if ok && approval.ApprovalCriteria.AllowBackedMinting && approval.ApprovalCriteria.MustPrioritize {
			if approval.FromListId == backingAddr {
				hasDeposit = true
			}
			if approval.ToListId == backingAddr {
				hasWithdraw = true
			}
		}
	}

	if ok && !hasDeposit {
		violations = append(violations, StandardViolation{
			Code:    "VAULT_MISSING_DEPOSIT_APPROVAL",
			Message: "Vault collection must define a prioritized allowBackedMinting deposit approval (fromListId == backingAddr)",
		})
	}
	if ok && !hasWithdraw {
		violations = append(violations, StandardViolation{
			Code:    "VAULT_MISSING_WITHDRAW_APPROVAL",
			Message: "Vault collection must define a prioritized allowBackedMinting withdrawal approval (toListId == backingAddr)",
		})
	}
	if c.Invariants == nil || !c.Invariants.NoForcefulPostMintTransfers {
		violations = append(violations, StandardViolation{
			Code:    "VAULT_FORCEFUL_TRANSFERS_ALLOWED",
			Message: "Vault collection must enforce invariants.noForcefulPostMintTransfers == true",
		})
	}
	if !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "VAULT_UNFROZEN_APPROVALS",
			Message: "Vault collection must permanently freeze collectionPermissions.canUpdateCollectionApprovals",
		})
	}

	return violations
}

func verifyAgentWallet(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	_, hasBacked := getValidBackingAddress(c)
	hasSpendApproval := false

	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		ac := approval.ApprovalCriteria
		isEscrowSpend := hasEscrowCoinTransfer(ac.CoinTransfers)
		isSpendPath := isEscrowSpend || ac.AllowBackedMinting
		if !isSpendPath {
			continue
		}
		hasSpendApproval = true

		if !approvalHasSpendingBound(approval) {
			violations = append(violations, StandardViolation{
				Code:    "AGENT_WALLET_UNBOUNDED_SPEND",
				Message: fmt.Sprintf("Agent Wallet spending approval %q must enforce maxNumTransfers, approvalAmounts, or a finite transferTimes window", approval.ApprovalId),
			})
		}
	}

	if !hasBacked && !hasSpendApproval {
		violations = append(violations, StandardViolation{
			Code:    "AGENT_WALLET_MISSING_SPEND_CAPABILITY",
			Message: "Agent Wallet must configure a cosmosCoinBackedPath or bounded escrow coinTransfer approval",
		})
	}

	hasForcefulGuard := c.Invariants != nil && c.Invariants.NoForcefulPostMintTransfers
	if !hasForcefulGuard && !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "AGENT_WALLET_UNPROTECTED_CUSTODY",
			Message: "Agent Wallet must enforce noForcefulPostMintTransfers or permanently freeze canUpdateCollectionApprovals",
		})
	}

	return violations
}

func verifyPaymentRequestV2(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	hasSingleUsePayment := false
	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		ac := approval.ApprovalCriteria
		if len(ac.CoinTransfers) == 0 {
			continue
		}
		singleUse := (ac.MaxNumTransfers != nil && ac.MaxNumTransfers.OverallMaxNumTransfers.Equal(sdkmath.NewUint(1))) ||
			(ac.AutoDeletionOptions != nil && ac.AutoDeletionOptions.AfterOneUse)
		if singleUse {
			hasSingleUsePayment = true
			break
		}
	}

	if !hasSingleUsePayment {
		violations = append(violations, StandardViolation{
			Code:    "PAYMENT_REQUEST_MISSING_BOUNDED_PAYMENT",
			Message: "PaymentRequestV2 must have an active single-use approval (overallMaxNumTransfers == 1 or afterOneUse) with coinTransfers",
		})
	}
	if !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "PAYMENT_REQUEST_UNFROZEN_APPROVALS",
			Message: "PaymentRequestV2 must permanently freeze canUpdateCollectionApprovals so invoice terms cannot be mutated",
		})
	}

	return violations
}

func verifyBounty(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	hasEscrowPayout := false
	hasUnlatchedEscrowPayout := false
	trackerBranches := make(map[string]int)
	escrowTrackerIds := []string{}

	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		ac := approval.ApprovalCriteria
		isLatched := ac.MaxNumTransfers != nil &&
			ac.MaxNumTransfers.OverallMaxNumTransfers.Equal(sdkmath.NewUint(1)) &&
			ac.MaxNumTransfers.AmountTrackerId != ""
		if isLatched {
			trackerBranches[ac.MaxNumTransfers.AmountTrackerId]++
		}
		if hasEscrowCoinTransfer(ac.CoinTransfers) {
			hasEscrowPayout = true
			if !isLatched {
				hasUnlatchedEscrowPayout = true
			} else {
				escrowTrackerIds = append(escrowTrackerIds, ac.MaxNumTransfers.AmountTrackerId)
			}
		}
	}

	hasMutualExclusion := len(escrowTrackerIds) > 0 && !hasUnlatchedEscrowPayout
	for _, tid := range escrowTrackerIds {
		if trackerBranches[tid] < 2 {
			hasMutualExclusion = false
			break
		}
	}

	if !hasEscrowPayout {
		violations = append(violations, StandardViolation{
			Code:    "BOUNTY_MISSING_ESCROW_PAYOUT",
			Message: "Bounty must define an approval with escrow-funded coinTransfers (overrideFromWithApproverAddress == true)",
		})
	}
	if !hasMutualExclusion {
		violations = append(violations, StandardViolation{
			Code:    "BOUNTY_MISSING_MUTUAL_EXCLUSION",
			Message: "Bounty must enforce terminal-state mutual exclusion on every escrow payout approval (shared overallMaxNumTransfers == 1 and identical amountTrackerId across >= 2 branches)",
		})
	}
	if !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "BOUNTY_UNFROZEN_APPROVALS",
			Message: "Bounty must permanently freeze canUpdateCollectionApprovals so payout/refund exclusion cannot be bypassed",
		})
	}

	return violations
}

func verifyCrowdfund(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	hasBoundedContribution := false
	hasEscrowSettlement := false

	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) || approval.ApprovalCriteria == nil {
			continue
		}
		ac := approval.ApprovalCriteria

		if approval.FromListId == MintAddress &&
			hasDepositToEscrowCoinTransfer(ac.CoinTransfers, c.MintEscrowAddress) &&
			hasFiniteTransferWindow(approval.TransferTimes) {
			hasBoundedContribution = true
		}

		if hasEscrowCoinTransfer(ac.CoinTransfers) {
			hasEscrowSettlement = true
		}
	}

	if !hasBoundedContribution {
		violations = append(violations, StandardViolation{
			Code:    "CROWDFUND_MISSING_BOUNDED_CONTRIBUTION",
			Message: "Crowdfund must define a Mint contribution approval with escrow coinTransfers and a finite deadline (transferTimes.end < MaxUint64)",
		})
	}
	if !hasEscrowSettlement {
		violations = append(violations, StandardViolation{
			Code:    "CROWDFUND_MISSING_ESCROW_SETTLEMENT",
			Message: "Crowdfund must define an escrow release/refund approval (overrideFromWithApproverAddress == true)",
		})
	}
	if !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "CROWDFUND_UNFROZEN_APPROVALS",
			Message: "Crowdfund must permanently freeze canUpdateCollectionApprovals",
		})
	}

	return violations
}

func verifySpendableCredit(c *TokenCollection) []StandardViolation {
	var violations []StandardViolation

	if c.Invariants == nil || !c.Invariants.NoCustomOwnershipTimes {
		violations = append(violations, StandardViolation{
			Code:    "SPENDABLE_CREDIT_CUSTOM_OWNERSHIP_TIMES_ALLOWED",
			Message: "Spendable Credit must enforce invariants.noCustomOwnershipTimes == true",
		})
	}
	if c.Invariants == nil || !c.Invariants.NoForcefulPostMintTransfers {
		violations = append(violations, StandardViolation{
			Code:    "SPENDABLE_CREDIT_FORCEFUL_TRANSFERS_ALLOWED",
			Message: "Spendable Credit must enforce invariants.noForcefulPostMintTransfers == true",
		})
	}

	_, hasBacked := getValidBackingAddress(c)
	hasIssueApproval := false
	hasSpendApproval := false

	for _, approval := range c.CollectionApprovals {
		if approval == nil || !approvalIsActive(approval) {
			continue
		}
		ac := approval.ApprovalCriteria
		if approval.FromListId == MintAddress {
			if ac != nil && (len(ac.CoinTransfers) > 0 ||
				!ApprovalAmountsIsBasicallyNil(ac.ApprovalAmounts) ||
				(ac.PredeterminedBalances != nil && ac.PredeterminedBalances.IncrementedBalances != nil && ac.PredeterminedBalances.IncrementedBalances.AllowAmountScaling)) {
				hasIssueApproval = true
			}
		} else {
			hasSpendApproval = true
		}
	}

	if !hasIssueApproval && !hasBacked {
		violations = append(violations, StandardViolation{
			Code:    "SPENDABLE_CREDIT_MISSING_ISSUE_PATH",
			Message: "Spendable Credit must define a Mint issuance approval (with payment, scaling, or amount cap) or backed path",
		})
	}
	if !hasSpendApproval {
		violations = append(violations, StandardViolation{
			Code:    "SPENDABLE_CREDIT_MISSING_SPEND_PATH",
			Message: "Spendable Credit must define a post-mint spend/redemption approval",
		})
	}
	if !areAllCollectionApprovalsFrozen(c.CollectionPermissions) {
		violations = append(violations, StandardViolation{
			Code:    "SPENDABLE_CREDIT_UNFROZEN_APPROVALS",
			Message: "Spendable Credit must permanently freeze canUpdateCollectionApprovals",
		})
	}

	return violations
}

// --- Internal deterministic helpers ---

func approvalIsActive(a *CollectionApproval) bool {
	if a == nil {
		return false
	}
	if len(a.TokenIds) == 0 || len(a.TransferTimes) == 0 || len(a.OwnershipTimes) == 0 {
		return false
	}
	if a.FromListId == "None" || a.ToListId == "None" || a.InitiatedByListId == "None" {
		return false
	}
	return true
}

func isFullRange(ranges []*UintRange) bool {
	if len(ranges) == 0 {
		return false
	}
	maxVal := sdkmath.NewUint(MaxUint64Value)
	for _, r := range ranges {
		if r == nil || r.Start.IsNil() || r.End.IsNil() {
			continue
		}
		if r.Start.LTE(sdkmath.NewUint(1)) && r.End.Equal(maxVal) {
			return true
		}
	}
	return false
}

func coversAllWildcard(listId string) bool {
	return listId == "All" || listId == "AllWithMint"
}

func coversAllNonMint(listId string) bool {
	return listId == "All" || listId == "AllWithMint" || listId == "AllWithoutMint" || listId == "!Mint"
}

// arePostMintCollectionApprovalsFrozen checks first-match permission semantics
// in CollectionPermissions.CanUpdateCollectionApprovals to verify that updating
// non-Mint collection approvals is permanently forbidden for [1, MaxUint64].
func arePostMintCollectionApprovalsFrozen(perms *CollectionPermissions) bool {
	if perms == nil || len(perms.CanUpdateCollectionApprovals) == 0 {
		return false
	}
	for _, p := range perms.CanUpdateCollectionApprovals {
		if p == nil {
			continue
		}
		matchesPostMint := coversAllNonMint(p.FromListId) &&
			coversAllNonMint(p.ToListId) &&
			coversAllNonMint(p.InitiatedByListId) &&
			p.ApprovalId == "All" &&
			isFullRange(p.TransferTimes) &&
			isFullRange(p.TokenIds) &&
			isFullRange(p.OwnershipTimes)

		// First-match semantics: ANY earlier rule permitting updates for a non-Mint
		// sender (including a targeted address carve-out) breaks the post-mint freeze.
		if len(p.PermanentlyPermittedTimes) > 0 && p.FromListId != MintAddress {
			return false
		}
		if matchesPostMint && isFullRange(p.PermanentlyForbiddenTimes) {
			return true
		}
	}
	return false
}

// areAllCollectionApprovalsFrozen verifies that ALL collection approvals (including Mint)
// are permanently frozen in CanUpdateCollectionApprovals.
func areAllCollectionApprovalsFrozen(perms *CollectionPermissions) bool {
	if perms == nil || len(perms.CanUpdateCollectionApprovals) == 0 {
		return false
	}
	for _, p := range perms.CanUpdateCollectionApprovals {
		if p == nil {
			continue
		}
		if len(p.PermanentlyPermittedTimes) > 0 {
			return false
		}
		matchesAll := coversAllWildcard(p.FromListId) &&
			coversAllWildcard(p.ToListId) &&
			coversAllWildcard(p.InitiatedByListId) &&
			p.ApprovalId == "All" &&
			isFullRange(p.TransferTimes) &&
			isFullRange(p.TokenIds) &&
			isFullRange(p.OwnershipTimes)
		if matchesAll && isFullRange(p.PermanentlyForbiddenTimes) {
			return true
		}
	}
	return false
}

func getValidBackingAddress(c *TokenCollection) (string, bool) {
	if c == nil || c.Invariants == nil || c.Invariants.CosmosCoinBackedPath == nil {
		return "", false
	}
	bp := c.Invariants.CosmosCoinBackedPath
	if bp.Address == "" || bp.Conversion == nil || bp.Conversion.SideA == nil {
		return "", false
	}
	if bp.Conversion.SideA.Denom == "" || bp.Conversion.SideA.Amount.IsNil() || bp.Conversion.SideA.Amount.IsZero() {
		return "", false
	}
	hasPositiveSideB := false
	for _, b := range bp.Conversion.SideB {
		if b != nil && !b.Amount.IsNil() && !b.Amount.IsZero() {
			hasPositiveSideB = true
			break
		}
	}
	if !hasPositiveSideB {
		return "", false
	}
	return bp.Address, true
}

func isConcreteAddressListId(listId string) bool {
	switch listId {
	case "", "All", "AllWithMint", "AllWithoutMint", "None", MintAddress:
		return false
	default:
		return true
	}
}

func hasEscrowCoinTransfer(transfers []*CoinTransfer) bool {
	for _, ct := range transfers {
		if ct != nil && ct.OverrideFromWithApproverAddress && len(ct.Coins) > 0 {
			return true
		}
	}
	return false
}

func hasDepositToEscrowCoinTransfer(transfers []*CoinTransfer, escrowAddr string) bool {
	for _, ct := range transfers {
		if ct == nil || len(ct.Coins) == 0 {
			continue
		}
		if ct.To == MintAddress || (escrowAddr != "" && ct.To == escrowAddr) {
			return true
		}
	}
	return false
}

func hasFiniteTransferWindow(ranges []*UintRange) bool {
	if len(ranges) == 0 {
		return false
	}
	maxVal := sdkmath.NewUint(MaxUint64Value)
	for _, r := range ranges {
		if r == nil || r.Start.IsNil() || r.End.IsNil() || r.Start.GT(r.End) || !r.End.LT(maxVal) {
			return false
		}
	}
	return true
}

func approvalHasSpendingBound(a *CollectionApproval) bool {
	if a == nil || a.ApprovalCriteria == nil {
		return false
	}
	if !MaxNumTransfersIsBasicallyNil(a.ApprovalCriteria.MaxNumTransfers) {
		return true
	}
	if !ApprovalAmountsIsBasicallyNil(a.ApprovalCriteria.ApprovalAmounts) {
		return true
	}
	if hasFiniteTransferWindow(a.TransferTimes) {
		return true
	}
	return false
}
