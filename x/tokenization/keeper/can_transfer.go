package keeper

import (
	"fmt"
	"math"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CanTransferRequest specifies a read-only transfer policy introspection query (BB-378).
type CanTransferRequest struct {
	CollectionId sdkmath.Uint    `json:"collectionId"`
	Transfer     *types.Transfer `json:"transfer"`
	InitiatedBy  string          `json:"initiatedBy"`
	// Optional block-time override in unix milliseconds (0 uses ctx.BlockTime()).
	AtTimestampMs uint64 `json:"atTimestampMs,omitempty"`
}

// FailedCriterion describes a specific structured reason why a candidate approval
// or transfer precondition rejected the transfer.
type FailedCriterion struct {
	ApprovalId    string `json:"approvalId,omitempty"`
	ApprovalLevel string `json:"approvalLevel"` // "collection", "outgoing", "incoming", "validation", "balance"
	Criterion     string `json:"criterion"`
	Expected      string `json:"expected,omitempty"`
	Actual        string `json:"actual,omitempty"`
	Message       string `json:"message"`
}

// CanTransferResult is the structured, side-effect-free response of EvaluateCanTransfer.
type CanTransferResult struct {
	Allowed            bool              `json:"allowed"`
	MatchedApprovalIds []string          `json:"matchedApprovalIds"`
	FailedCriteria     []FailedCriterion `json:"failedCriteria,omitempty"`
	RequiredCoins      sdk.Coins         `json:"requiredCoins"`
}

// EvaluateCanTransfer performs a strictly read-only evaluation of whether req.Transfer
// would be permitted under current on-chain state (optionally evaluated at req.AtTimestampMs).
//
// All execution runs inside discarded ctx.CacheContext() forks so it is impossible
// for EvaluateCanTransfer to mutate store state, bank balances, approval trackers,
// or caller-owned structs.
func (k Keeper) EvaluateCanTransfer(ctx sdk.Context, req CanTransferRequest) CanTransferResult {
	res := CanTransferResult{
		Allowed:            false,
		MatchedApprovalIds: []string{},
		FailedCriteria:     []FailedCriterion{},
		RequiredCoins:      sdk.NewCoins(),
	}

	if req.Transfer == nil {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "transfer",
			Message:       "transfer cannot be nil",
		})
		return res
	}
	if len(req.Transfer.ToAddresses) == 0 {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "toAddresses",
			Message:       "transfer must specify at least one recipient address",
		})
		return res
	}
	if req.InitiatedBy == "" {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "initiatedBy",
			Message:       "initiatedBy address cannot be empty",
		})
		return res
	}
	if req.CollectionId.IsNil() || req.CollectionId.IsZero() {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "collectionId",
			Message:       "collectionId must be initialized and greater than zero",
		})
		return res
	}
	if req.AtTimestampMs > uint64(math.MaxInt64) {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "atTimestampMs",
			Expected:      fmt.Sprintf("<= %d", int64(math.MaxInt64)),
			Actual:        fmt.Sprintf("%d", req.AtTimestampMs),
			Message:       "atTimestampMs exceeds int64 range",
		})
		return res
	}

	// Fork an isolated cache context and discard its write-back function so zero
	// state mutations or events can ever leak to the caller's ctx.
	evalCtx, _ := ctx.CacheContext()
	if req.AtTimestampMs > 0 {
		evalCtx = evalCtx.WithBlockTime(time.UnixMilli(int64(req.AtTimestampMs)))
	}

	collection, found := k.GetCollectionFromStore(evalCtx, req.CollectionId)
	if !found || collection == nil {
		res.FailedCriteria = append(res.FailedCriteria, FailedCriterion{
			ApprovalLevel: "validation",
			Criterion:     "collectionNotFound",
			Expected:      "existing collection",
			Actual:        req.CollectionId.String(),
			Message:       fmt.Sprintf("collection %s not found", req.CollectionId.String()),
		})
		return res
	}

	// Deep-copy transfer so HandleTransfersEnriched cannot mutate caller memory.
	transferCopy := cloneTransfer(req.Transfer)

	tracking, err := k.HandleTransfersEnriched(evalCtx, collection, []*types.Transfer{transferCopy}, req.InitiatedBy)
	if err == nil {
		res.Allowed = true
		res.MatchedApprovalIds = uniqueApprovalIds(tracking.AllApprovalsUsed)
		res.RequiredCoins = aggregateCoinTransfers(tracking.AllCoinTransfers)
		return res
	}

	// Transfer failed — run structured diagnostics in a fresh discarded CacheContext.
	diagCtx, _ := ctx.CacheContext()
	if req.AtTimestampMs > 0 {
		diagCtx = diagCtx.WithBlockTime(time.UnixMilli(int64(req.AtTimestampMs)))
	}
	diagCol, _ := k.GetCollectionFromStore(diagCtx, req.CollectionId)
	diagTransfer := cloneTransfer(req.Transfer)

	failed, reqCoins := k.diagnoseTransferFailure(diagCtx, diagCol, diagTransfer, req.InitiatedBy, err)
	res.FailedCriteria = failed
	res.RequiredCoins = reqCoins
	return res
}

func cloneTransfer(t *types.Transfer) *types.Transfer {
	if t == nil {
		return nil
	}
	copied := *t
	if t.ToAddresses != nil {
		copied.ToAddresses = append([]string(nil), t.ToAddresses...)
	}
	copied.Balances = types.DeepCopyBalances(t.Balances)
	if t.PrioritizedApprovals != nil {
		copied.PrioritizedApprovals = make([]*types.ApprovalIdentifierDetails, len(t.PrioritizedApprovals))
		for i, pa := range t.PrioritizedApprovals {
			if pa != nil {
				c := *pa
				copied.PrioritizedApprovals[i] = &c
			}
		}
	}
	if t.PrecalculateBalancesFromApproval != nil {
		p := *t.PrecalculateBalancesFromApproval
		if p.PrecalculationOptions != nil {
			opts := *p.PrecalculationOptions
			p.PrecalculationOptions = &opts
		}
		copied.PrecalculateBalancesFromApproval = &p
	}
	if t.MerkleProofs != nil {
		copied.MerkleProofs = make([]*types.MerkleProof, len(t.MerkleProofs))
		for i, mp := range t.MerkleProofs {
			if mp != nil {
				c := *mp
				copied.MerkleProofs[i] = &c
			}
		}
	}
	if t.EthSignatureProofs != nil {
		copied.EthSignatureProofs = make([]*types.ETHSignatureProof, len(t.EthSignatureProofs))
		for i, ep := range t.EthSignatureProofs {
			if ep != nil {
				c := *ep
				copied.EthSignatureProofs[i] = &c
			}
		}
	}
	return &copied
}

func (k Keeper) diagnoseTransferFailure(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transfer *types.Transfer,
	initiatedBy string,
	execErr error,
) ([]FailedCriterion, sdk.Coins) {
	var failures []FailedCriterion
	requiredCoins := sdk.NewCoins()

	if collection == nil || transfer == nil {
		return []FailedCriterion{{
			ApprovalLevel: "validation",
			Criterion:     "transferExecution",
			Message:       execErr.Error(),
		}}, requiredCoins
	}

	// 1. Archive check
	if types.GetIsArchived(ctx, collection) {
		return []FailedCriterion{{
			ApprovalLevel: "validation",
			Criterion:     "archived",
			Expected:      "isArchived=false",
			Actual:        "isArchived=true",
			Message:       "collection is archived (read-only)",
		}}, requiredCoins
	}

	// 2. Basic & invariant validation
	if err := types.ValidateTransferWithInvariants(ctx, transfer, true, collection); err != nil {
		criterion := "validation"
		if strings.Contains(err.Error(), "noCustomOwnershipTimes") {
			criterion = "ownershipTimes"
		} else if strings.Contains(err.Error(), "cosmosCoinBackedPath") {
			criterion = "invariants"
		}
		return []FailedCriterion{{
			ApprovalLevel: "validation",
			Criterion:     criterion,
			Message:       err.Error(),
		}}, requiredCoins
	}

	// 3. ValidTokenIds check on Mint / backing path
	if types.IsMintAddress(transfer.From) || k.IsBackingPathAddress(ctx, collection, transfer.From) {
		if err := types.ValidateBalancesWithinValidTokenIds(transfer.Balances, collection.ValidTokenIds); err != nil {
			return []FailedCriterion{{
				ApprovalLevel: "validation",
				Criterion:     "validTokenIds",
				Message:       err.Error(),
			}}, requiredCoins
		}
	}

	// 4. Sender balance check (for non-Mint, non-backed addresses)
	fromBal, _, err := k.GetBalanceOrApplyDefaultReadOnly(ctx, collection, transfer.From)
	if err == nil && !types.IsMintAddress(transfer.From) && !k.IsSpecialBackedAddress(ctx, collection, transfer.From) {
		testFrom := types.DeepCopyBalances(fromBal.Balances)
		for range transfer.ToAddresses {
			for _, bal := range transfer.Balances {
				var subErr error
				testFrom, subErr = types.SubtractBalance(ctx, testFrom, bal, false)
				if subErr != nil {
					failures = append(failures, FailedCriterion{
						ApprovalLevel: "balance",
						Criterion:     "insufficientBalance",
						Message:       fmt.Sprintf("sender %s has insufficient token balance: %s", transfer.From, subErr.Error()),
					})
					break
				}
			}
		}
	}

	// 5. Inspect collection approvals per recipient
	currTimeMs := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	for _, to := range transfer.ToAddresses {
		toBal, _, _ := k.GetBalanceOrApplyDefaultReadOnly(ctx, collection, to)

		// First test whether collection approvals actually pass (meaning failure is user-level or post-transfer).
		colTestCtx, _ := ctx.CacheContext()
		colApprovalsUsed := []ApprovalsUsed{}
		colCoinTransfers := []CoinTransfers{}
		colEvents := &EventTracking{ApprovalsUsed: &colApprovalsUsed, CoinTransfers: &colCoinTransfers}
		meta := TransferMetadata{
			From:          transfer.From,
			To:            to,
			InitiatedBy:   initiatedBy,
			ApprovalLevel: "collection",
		}
		userChecks, colErr := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(colTestCtx, collection, transfer, meta, colEvents, "collection")
		if colErr == nil {
			requiredCoins = requiredCoins.Add(aggregateCoinTransfers(colCoinTransfers)...)
			// Collection approval passed; inspect user outgoing/incoming approvals.
			for _, uc := range userChecks {
				uTransfer := &types.Transfer{
					From:                                    transfer.From,
					ToAddresses:                             []string{to},
					Balances:                                uc.Balances,
					MerkleProofs:                            transfer.MerkleProofs,
					EthSignatureProofs:                      transfer.EthSignatureProofs,
					Memo:                                    transfer.Memo,
					PrioritizedApprovals:                    transfer.PrioritizedApprovals,
					OnlyCheckPrioritizedCollectionApprovals: transfer.OnlyCheckPrioritizedCollectionApprovals,
					OnlyCheckPrioritizedIncomingApprovals:   transfer.OnlyCheckPrioritizedIncomingApprovals,
					OnlyCheckPrioritizedOutgoingApprovals:   transfer.OnlyCheckPrioritizedOutgoingApprovals,
					PrecalculateBalancesFromApproval:        transfer.PrecalculateBalancesFromApproval,
				}
				uCtx, _ := ctx.CacheContext()
				uUsed := []ApprovalsUsed{}
				uCoins := []CoinTransfers{}
				uEvents := &EventTracking{ApprovalsUsed: &uUsed, CoinTransfers: &uCoins}
				if uc.Outgoing {
					if fromBal == nil {
						failures = append(failures, FailedCriterion{
							ApprovalLevel: "outgoing",
							Criterion:     "userOutgoingApproval",
							Message:       "sender balance store unavailable",
						})
						continue
					}
					uMeta := TransferMetadata{From: transfer.From, To: to, InitiatedBy: initiatedBy, ApproverAddress: transfer.From, ApprovalLevel: "outgoing"}
					if uErr := k.DeductUserOutgoingApprovals(uCtx, collection, transfer.Balances, uTransfer, uMeta, fromBal, uEvents, uc.UserApprovalSettings); uErr != nil {
						failures = append(failures, FailedCriterion{
							ApprovalLevel: "outgoing",
							Criterion:     "userOutgoingApproval",
							Message:       uErr.Error(),
						})
					}
				} else {
					if toBal == nil {
						failures = append(failures, FailedCriterion{
							ApprovalLevel: "incoming",
							Criterion:     "userIncomingApproval",
							Message:       "recipient balance store unavailable",
						})
						continue
					}
					uMeta := TransferMetadata{From: transfer.From, To: to, InitiatedBy: initiatedBy, ApproverAddress: to, ApprovalLevel: "incoming"}
					if uErr := k.DeductUserIncomingApprovals(uCtx, collection, transfer.Balances, uTransfer, uMeta, toBal, uEvents, uc.UserApprovalSettings); uErr != nil {
						failures = append(failures, FailedCriterion{
							ApprovalLevel: "incoming",
							Criterion:     "userIncomingApproval",
							Message:       uErr.Error(),
						})
					}
				}
			}
			continue
		}

		// Collection approval failed: diagnose candidate collection approvals.
		if len(collection.CollectionApprovals) == 0 {
			failures = append(failures, FailedCriterion{
				ApprovalLevel: "collection",
				Criterion:     "noCollectionApprovals",
				Message:       "collection has no collection approvals configured",
			})
			continue
		}

		candidates := collection.CollectionApprovals
		if transfer.OnlyCheckPrioritizedCollectionApprovals {
			prioritizedSet := make(map[string]bool)
			for _, pa := range transfer.PrioritizedApprovals {
				if pa != nil && pa.ApprovalLevel == "collection" && pa.ApproverAddress == "" {
					prioritizedSet[pa.ApprovalId] = true
				}
			}
			filtered := make([]*types.CollectionApproval, 0, len(prioritizedSet))
			for _, ap := range collection.CollectionApprovals {
				if ap != nil && prioritizedSet[ap.ApprovalId] {
					filtered = append(filtered, ap)
				}
			}
			if len(filtered) == 0 {
				failures = append(failures, FailedCriterion{
					ApprovalLevel: "collection",
					Criterion:     "prioritizedApprovalNotFound",
					Message:       "onlyCheckPrioritizedCollectionApprovals is true, but none of the specified prioritized approvals exist on the collection",
				})
				continue
			}
			candidates = filtered
		}

		for _, approval := range candidates {
			if approval == nil {
				continue
			}
			apFailures, apCoins := k.diagnoseSingleCollectionApproval(ctx, collection, approval, transfer, to, initiatedBy, currTimeMs)
			failures = append(failures, apFailures...)
			requiredCoins = requiredCoins.Add(apCoins...)
		}
	}

	if len(failures) == 0 {
		failures = append(failures, FailedCriterion{
			ApprovalLevel: "collection",
			Criterion:     "transferExecution",
			Message:       execErr.Error(),
		})
	}

	return failures, requiredCoins
}

func (k Keeper) diagnoseSingleCollectionApproval(
	ctx sdk.Context,
	collection *types.TokenCollection,
	approval *types.CollectionApproval,
	transfer *types.Transfer,
	to string,
	initiatedBy string,
	currTimeMs sdkmath.Uint,
) ([]FailedCriterion, sdk.Coins) {
	var failures []FailedCriterion
	reqCoins := sdk.NewCoins()

	// Collect any coin transfers declared on this approval for caller ergonomics.
	if approval.ApprovalCriteria != nil {
		for _, ct := range approval.ApprovalCriteria.CoinTransfers {
			if ct == nil {
				continue
			}
			for _, c := range ct.Coins {
				if c != nil && !c.Amount.IsNil() && c.Amount.IsPositive() {
					reqCoins = reqCoins.Add(*c)
				}
			}
		}
	}

	// 1. Address list match
	if !k.CheckIfAddressesMatchCollectionListIds(ctx, approval, transfer.From, to, initiatedBy) {
		return []FailedCriterion{{
			ApprovalId:    approval.ApprovalId,
			ApprovalLevel: "collection",
			Criterion:     "addressMatch",
			Expected:      fmt.Sprintf("from=%s to=%s initiatedBy=%s", approval.FromListId, approval.ToListId, approval.InitiatedByListId),
			Actual:        fmt.Sprintf("from=%s to=%s initiatedBy=%s", transfer.From, to, initiatedBy),
			Message:       "transfer addresses do not match approval list IDs",
		}}, sdk.NewCoins()
	}

	// 2. Transfer time window
	inTime, timeErr := types.SearchUintRangesForUint(currTimeMs, approval.TransferTimes)
	if !inTime || timeErr != nil {
		return []FailedCriterion{{
			ApprovalId:    approval.ApprovalId,
			ApprovalLevel: "collection",
			Criterion:     "transferTime",
			Actual:        currTimeMs.String(),
			Message:       fmt.Sprintf("transfer timestamp %s ms is outside approval transferTimes", currTimeMs.String()),
		}}, reqCoins
	}

	// 3. Token IDs & ownership times overlap
	matchedBals, balErr := types.GetBalancesForIds(ctx, approval.TokenIds, approval.OwnershipTimes, transfer.Balances)
	matchedBals = types.FilterZeroBalances(matchedBals)
	if balErr != nil || len(matchedBals) == 0 {
		return []FailedCriterion{{
			ApprovalId:    approval.ApprovalId,
			ApprovalLevel: "collection",
			Criterion:     "ownershipTimes",
			Message:       "requested tokenIds or ownershipTimes do not overlap with approval",
		}}, reqCoins
	}

	// 4. Prioritization requirement
	isPrioritized := false
	for _, pa := range transfer.PrioritizedApprovals {
		if pa != nil && pa.ApprovalId == approval.ApprovalId && pa.ApprovalLevel == "collection" && pa.ApproverAddress == "" {
			if !pa.Version.IsNil() && !pa.Version.Equal(approval.Version) {
				failures = append(failures, FailedCriterion{
					ApprovalId:    approval.ApprovalId,
					ApprovalLevel: "collection",
					Criterion:     "approvalVersion",
					Expected:      approval.Version.String(),
					Actual:        pa.Version.String(),
					Message:       "prioritized approval version mismatch",
				})
				return failures, reqCoins
			}
			isPrioritized = true
			break
		}
	}
	if !types.CollectionApprovalIsAutoScannable(approval.ApprovalCriteria) && !isPrioritized {
		failures = append(failures, FailedCriterion{
			ApprovalId:    approval.ApprovalId,
			ApprovalLevel: "collection",
			Criterion:     "mustPrioritize",
			Expected:      fmt.Sprintf("prioritizedApprovals containing %q (version %s)", approval.ApprovalId, approval.Version.String()),
			Actual:        "not prioritized",
			Message:       "approval is non-auto-scannable or has mustPrioritize=true and must be listed in prioritizedApprovals",
		})
	}

	ac := approval.ApprovalCriteria
	if ac == nil {
		return failures, reqCoins
	}

	// 5. Criteria checkers (MustOwnTokens, AddressChecks, AltTimeChecks, DynamicStore, Voting, EVMQuery, party equality)
	meta := TransferMetadata{
		From:          transfer.From,
		To:            to,
		InitiatedBy:   initiatedBy,
		ApprovalLevel: "collection",
	}
	for _, checker := range (&k).GetApprovalCriteriaCheckers(approval) {
		cCtx, _ := ctx.CacheContext()
		detMsg, err := checker.Check(cCtx, approval, collection, to, transfer.From, initiatedBy, "collection", "", transfer.MerkleProofs, transfer.EthSignatureProofs, transfer.Memo, isPrioritized)
		if err != nil {
			msg := detMsg
			if msg == "" {
				msg = err.Error()
			}
			failures = append(failures, FailedCriterion{
				ApprovalId:    approval.ApprovalId,
				ApprovalLevel: "collection",
				Criterion:     normalizeCheckerCriterion(checker.Name()),
				Message:       msg,
			})
		}
	}

	// 6. CoinTransfers check
	if len(ac.CoinTransfers) > 0 {
		cCtx, _ := ctx.CacheContext()
		coinLog := []CoinTransfers{}
		detMsg, err := k.ExecuteCoinTransfers(cCtx, ac.CoinTransfers, meta, &coinLog, collection, nil)
		if err != nil {
			msg := detMsg
			if msg == "" {
				msg = err.Error()
			}
			failures = append(failures, FailedCriterion{
				ApprovalId:    approval.ApprovalId,
				ApprovalLevel: "collection",
				Criterion:     "coinTransfer",
				Expected:      reqCoins.String(),
				Message:       msg,
			})
		}
	}

	// 7. MerkleChallenges check
	challengeIncrements := sdkmath.ZeroUint()
	if len(ac.MerkleChallenges) > 0 {
		mCtx, _ := ctx.CacheContext()
		detMsg, inc, err := k.HandleMerkleChallenges(mCtx, collection.CollectionId, transfer, approval, meta, false)
		if err != nil {
			msg := detMsg
			if msg == "" {
				msg = err.Error()
			}
			failures = append(failures, FailedCriterion{
				ApprovalId:    approval.ApprovalId,
				ApprovalLevel: "collection",
				Criterion:     "merkleChallenge",
				Message:       msg,
			})
		} else {
			challengeIncrements = inc
		}
	}

	// 8. ETHSignatureChallenges check
	if len(ac.EthSignatureChallenges) > 0 {
		eCtx, _ := ctx.CacheContext()
		detMsg, err := k.HandleETHSignatureChallenges(eCtx, collection.CollectionId, transfer, approval, meta)
		if err != nil {
			msg := detMsg
			if msg == "" {
				msg = err.Error()
			}
			failures = append(failures, FailedCriterion{
				ApprovalId:    approval.ApprovalId,
				ApprovalLevel: "collection",
				Criterion:     "ethSignatureChallenge",
				Message:       msg,
			})
		}
	}

	// 9. Trackers: MaxNumTransfers, ApprovalAmounts, PredeterminedBalances
	tCtx, _ := ctx.CacheContext()
	trackerTypes := []string{"overall", "to", "from", "initiatedBy"}
	approvedAmounts := []sdkmath.Uint{sdkmath.ZeroUint(), sdkmath.ZeroUint(), sdkmath.ZeroUint(), sdkmath.ZeroUint()}
	if ac.ApprovalAmounts != nil {
		approvedAmounts = []sdkmath.Uint{
			ac.ApprovalAmounts.OverallApprovalAmount,
			ac.ApprovalAmounts.PerToAddressApprovalAmount,
			ac.ApprovalAmounts.PerFromAddressApprovalAmount,
			ac.ApprovalAmounts.PerInitiatedByAddressApprovalAmount,
		}
	}
	maxTransfers := []sdkmath.Uint{sdkmath.ZeroUint(), sdkmath.ZeroUint(), sdkmath.ZeroUint(), sdkmath.ZeroUint()}
	if ac.MaxNumTransfers != nil {
		maxTransfers = []sdkmath.Uint{
			ac.MaxNumTransfers.OverallMaxNumTransfers,
			ac.MaxNumTransfers.PerToAddressMaxNumTransfers,
			ac.MaxNumTransfers.PerFromAddressMaxNumTransfers,
			ac.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers,
		}
	}
	addresses := []string{"", to, transfer.From, initiatedBy}

	var precalcOpts *types.PrecalculationOptions
	if transfer.PrecalculateBalancesFromApproval != nil {
		precalcOpts = transfer.PrecalculateBalancesFromApproval.PrecalculationOptions
	}

	for i, trackerType := range trackerTypes {
		detMsg, err := k.IncrementApprovalsAndAssertWithinThreshold(
			tCtx, collection, approval, transfer.Balances,
			approvedAmounts[i], maxTransfers[i], matchedBals,
			challengeIncrements, meta, trackerType, addresses[i], precalcOpts,
		)
		if err != nil {
			msg := detMsg
			if msg == "" {
				msg = err.Error()
			}
			criterion := "maxNumTransfers"
			lower := strings.ToLower(msg + " " + err.Error())
			switch {
			case strings.Contains(lower, "predetermined") || strings.Contains(lower, "scaling"):
				criterion = "predeterminedBalances"
			case strings.Contains(lower, "amount"):
				criterion = "approvalAmounts"
			}
			failures = append(failures, FailedCriterion{
				ApprovalId:    approval.ApprovalId,
				ApprovalLevel: "collection",
				Criterion:     criterion,
				Message:       msg,
			})
			break
		}
	}

	return failures, reqCoins
}

func normalizeCheckerCriterion(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "mustowntoken") || strings.Contains(lower, "must_own"):
		return "mustOwnTokens"
	case strings.Contains(lower, "alttime") || strings.Contains(lower, "alt_time"):
		return "altTimeChecks"
	case strings.Contains(lower, "dynamicstore") || strings.Contains(lower, "dynamic_store"):
		return "dynamicStoreChallenge"
	case strings.Contains(lower, "voting"):
		return "votingChallenge"
	case strings.Contains(lower, "evm"):
		return "evmQueryChallenge"
	case strings.Contains(lower, "sender"):
		return "senderChecks"
	case strings.Contains(lower, "recipient"):
		return "recipientChecks"
	case strings.Contains(lower, "initiator"):
		return "initiatorChecks"
	default:
		if name == "" {
			return "criteriaChecker"
		}
		return name
	}
}

func uniqueApprovalIds(used []ApprovalsUsed) []string {
	seen := make(map[string]bool, len(used))
	out := make([]string, 0, len(used))
	for _, u := range used {
		if u.ApprovalId == "" || seen[u.ApprovalId] {
			continue
		}
		seen[u.ApprovalId] = true
		out = append(out, u.ApprovalId)
	}
	return out
}

func aggregateCoinTransfers(transfers []CoinTransfers) sdk.Coins {
	coins := sdk.NewCoins()
	for _, ct := range transfers {
		if ct.Denom == "" || ct.Amount == "" {
			continue
		}
		amt, ok := sdkmath.NewIntFromString(ct.Amount)
		if !ok || !amt.IsPositive() {
			continue
		}
		coins = coins.Add(sdk.NewCoin(ct.Denom, amt))
	}
	return coins
}
