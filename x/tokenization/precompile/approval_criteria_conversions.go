package tokenization

import (
	"math/big"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// autoDeletionOptionsToSolidity converts AutoDeletionOptions to Solidity struct format.
func autoDeletionOptionsToSolidity(o *tokenizationtypes.AutoDeletionOptions) []interface{} {
	if o == nil {
		return []interface{}{false, false, false, false}
	}
	return []interface{}{
		o.AfterOneUse,
		o.AfterOverallMaxNumTransfers,
		o.AllowCounterpartyPurge,
		o.AllowPurgeIfExpired,
	}
}

// resetTimeIntervalsToSolidity converts ResetTimeIntervals to Solidity struct format.
func resetTimeIntervalsToSolidity(r *tokenizationtypes.ResetTimeIntervals) []interface{} {
	if r == nil {
		return []interface{}{big.NewInt(0), big.NewInt(0)}
	}
	start := big.NewInt(0)
	interval := big.NewInt(0)
	if !r.StartTime.IsNil() {
		start = r.StartTime.BigInt()
	}
	if !r.IntervalLength.IsNil() {
		interval = r.IntervalLength.BigInt()
	}
	return []interface{}{start, interval}
}

// approvalAmountsToSolidity converts ApprovalAmounts to Solidity struct format.
func approvalAmountsToSolidity(a *tokenizationtypes.ApprovalAmounts) []interface{} {
	if a == nil {
		return []interface{}{
			big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0),
			"", resetTimeIntervalsToSolidity(nil),
		}
	}
	return []interface{}{
		bigIntFromUint(a.OverallApprovalAmount),
		bigIntFromUint(a.PerToAddressApprovalAmount),
		bigIntFromUint(a.PerFromAddressApprovalAmount),
		bigIntFromUint(a.PerInitiatedByAddressApprovalAmount),
		a.AmountTrackerId,
		resetTimeIntervalsToSolidity(a.ResetTimeIntervals),
	}
}

// maxNumTransfersToSolidity converts MaxNumTransfers to Solidity struct format.
func maxNumTransfersToSolidity(m *tokenizationtypes.MaxNumTransfers) []interface{} {
	if m == nil {
		return []interface{}{
			big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0),
			"", resetTimeIntervalsToSolidity(nil),
		}
	}
	return []interface{}{
		bigIntFromUint(m.OverallMaxNumTransfers),
		bigIntFromUint(m.PerToAddressMaxNumTransfers),
		bigIntFromUint(m.PerFromAddressMaxNumTransfers),
		bigIntFromUint(m.PerInitiatedByAddressMaxNumTransfers),
		m.AmountTrackerId,
		resetTimeIntervalsToSolidity(m.ResetTimeIntervals),
	}
}

func bigIntFromUint(u tokenizationtypes.Uint) *big.Int {
	if u.IsNil() {
		return big.NewInt(0)
	}
	return u.BigInt()
}

// coinTransferToSolidity converts CoinTransfer to Solidity struct format.
// Coins (sdk.Coin) are flattened to coinDenoms []string and coinAmounts []*big.Int.
func coinTransferToSolidity(c *tokenizationtypes.CoinTransfer) []interface{} {
	if c == nil {
		return []interface{}{"", make([]interface{}, 0), make([]interface{}, 0), false, false}
	}
	var denoms []interface{}
	var amounts []interface{}
	for _, coin := range c.Coins {
		if coin == nil {
			continue
		}
		denoms = append(denoms, coin.Denom)
		amt := big.NewInt(0)
		if !coin.Amount.IsNil() {
			amt = coin.Amount.BigInt()
		}
		amounts = append(amounts, amt)
	}
	return []interface{}{
		c.To,
		denoms,
		amounts,
		c.OverrideFromWithApproverAddress,
		c.OverrideToWithInitiator,
	}
}

// mustOwnTokensToSolidity converts MustOwnTokens to Solidity struct format.
func mustOwnTokensToSolidity(m *tokenizationtypes.MustOwnTokens) []interface{} {
	if m == nil {
		return []interface{}{
			big.NewInt(0), []interface{}{big.NewInt(0), big.NewInt(0)}, uintRangesToSolidity(nil), uintRangesToSolidity(nil),
			false, false, "",
		}
	}
	amountRange := []interface{}{big.NewInt(0), big.NewInt(0)}
	if m.AmountRange != nil {
		amountRange = []interface{}{bigIntFromUint(m.AmountRange.Start), bigIntFromUint(m.AmountRange.End)}
	}
	return []interface{}{
		bigIntFromUint(m.CollectionId),
		amountRange,
		uintRangesToSolidity(m.OwnershipTimes),
		uintRangesToSolidity(m.TokenIds),
		m.OverrideWithCurrentTime,
		m.MustSatisfyForAllAssets,
		m.OwnershipCheckParty,
	}
}

// dynamicStoreChallengeToSolidity converts DynamicStoreChallenge to Solidity struct format.
func dynamicStoreChallengeToSolidity(d *tokenizationtypes.DynamicStoreChallenge) []interface{} {
	if d == nil {
		return []interface{}{big.NewInt(0), ""}
	}
	return []interface{}{
		bigIntFromUint(d.StoreId),
		d.OwnershipCheckParty,
	}
}

// addressChecksToSolidity converts AddressChecks to Solidity struct format.
func addressChecksToSolidity(a *tokenizationtypes.AddressChecks) []interface{} {
	if a == nil {
		return []interface{}{false, false, false, false}
	}
	return []interface{}{
		a.MustBeEvmContract,
		a.MustNotBeEvmContract,
		a.MustBeLiquidityPool,
		a.MustNotBeLiquidityPool,
	}
}

// altTimeChecksToSolidity converts AltTimeChecks to Solidity struct format.
func altTimeChecksToSolidity(a *tokenizationtypes.AltTimeChecks) []interface{} {
	if a == nil {
		return []interface{}{uintRangesToSolidity(nil), uintRangesToSolidity(nil)}
	}
	return []interface{}{
		uintRangesToSolidity(a.OfflineHours),
		uintRangesToSolidity(a.OfflineDays),
	}
}

// userRoyaltiesToSolidity converts UserRoyalties to Solidity struct format.
func userRoyaltiesToSolidity(u *tokenizationtypes.UserRoyalties) []interface{} {
	if u == nil {
		return []interface{}{big.NewInt(0), ""}
	}
	return []interface{}{
		bigIntFromUint(u.Percentage),
		u.PayoutAddress,
	}
}

// merkleChallengeToSolidity converts MerkleChallenge to Solidity struct format.
func merkleChallengeToSolidity(m *tokenizationtypes.MerkleChallenge) []interface{} {
	if m == nil {
		return []interface{}{
			"", big.NewInt(0), false, big.NewInt(0), "", "", "", "",
		}
	}
	return []interface{}{
		m.Root,
		bigIntFromUint(m.ExpectedProofLength),
		m.UseCreatorAddressAsLeaf,
		bigIntFromUint(m.MaxUsesPerLeaf),
		m.Uri,
		m.CustomData,
		m.ChallengeTrackerId,
		m.LeafSigner,
	}
}

// ethSignatureChallengeToSolidity converts ETHSignatureChallenge to Solidity struct format.
func ethSignatureChallengeToSolidity(e *tokenizationtypes.ETHSignatureChallenge) []interface{} {
	if e == nil {
		return []interface{}{"", "", "", ""}
	}
	return []interface{}{e.Signer, e.ChallengeTrackerId, e.Uri, e.CustomData}
}

// voterToSolidity converts Voter to Solidity struct format.
func voterToSolidity(v *tokenizationtypes.Voter) []interface{} {
	if v == nil {
		return []interface{}{"", big.NewInt(0)}
	}
	return []interface{}{v.Address, bigIntFromUint(v.Weight)}
}

// votingChallengeToSolidity converts VotingChallenge to Solidity struct format.
func votingChallengeToSolidity(v *tokenizationtypes.VotingChallenge) []interface{} {
	if v == nil {
		return []interface{}{"", big.NewInt(0), make([]interface{}, 0), "", ""}
	}
	voters := make([]interface{}, 0, len(v.Voters))
	for _, x := range v.Voters {
		voters = append(voters, voterToSolidity(x))
	}
	return []interface{}{
		v.ProposalId,
		bigIntFromUint(v.QuorumThreshold),
		voters,
		v.Uri,
		v.CustomData,
	}
}

// manualBalancesToSolidity converts ManualBalances to Solidity struct format.
func manualBalancesToSolidity(m *tokenizationtypes.ManualBalances) []interface{} {
	if m == nil {
		return []interface{}{make([]interface{}, 0)}
	}
	balances := make([]interface{}, 0, len(m.Balances))
	for _, b := range m.Balances {
		if b == nil {
			continue
		}
		conv, err := ConvertBalanceToSolidityStruct(b)
		if err != nil {
			continue
		}
		balances = append(balances, conv)
	}
	return []interface{}{balances}
}

// recurringOwnershipTimesToSolidity converts RecurringOwnershipTimes to Solidity struct format.
func recurringOwnershipTimesToSolidity(r *tokenizationtypes.RecurringOwnershipTimes) []interface{} {
	if r == nil {
		return []interface{}{big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	}
	return []interface{}{
		bigIntFromUint(r.StartTime),
		bigIntFromUint(r.IntervalLength),
		bigIntFromUint(r.ChargePeriodLength),
	}
}

// incrementedBalancesToSolidity converts IncrementedBalances to Solidity struct format.
func incrementedBalancesToSolidity(i *tokenizationtypes.IncrementedBalances) []interface{} {
	if i == nil {
		return []interface{}{
			make([]interface{}, 0), big.NewInt(0), big.NewInt(0), big.NewInt(0),
			false, recurringOwnershipTimesToSolidity(nil), false,
		}
	}
	startBalances := make([]interface{}, 0, len(i.StartBalances))
	for _, b := range i.StartBalances {
		if b == nil {
			continue
		}
		conv, err := ConvertBalanceToSolidityStruct(b)
		if err != nil {
			continue
		}
		startBalances = append(startBalances, conv)
	}
	return []interface{}{
		startBalances,
		bigIntFromUint(i.IncrementTokenIdsBy),
		bigIntFromUint(i.IncrementOwnershipTimesBy),
		bigIntFromUint(i.DurationFromTimestamp),
		i.AllowOverrideTimestamp,
		recurringOwnershipTimesToSolidity(i.RecurringOwnershipTimes),
		i.AllowOverrideWithAnyValidToken,
	}
}

// predeterminedOrderCalculationMethodToSolidity converts PredeterminedOrderCalculationMethod to Solidity struct format.
func predeterminedOrderCalculationMethodToSolidity(p *tokenizationtypes.PredeterminedOrderCalculationMethod) []interface{} {
	if p == nil {
		return []interface{}{false, false, false, false, false, ""}
	}
	return []interface{}{
		p.UseOverallNumTransfers,
		p.UsePerToAddressNumTransfers,
		p.UsePerFromAddressNumTransfers,
		p.UsePerInitiatedByAddressNumTransfers,
		p.UseMerkleChallengeLeafIndex,
		p.ChallengeTrackerId,
	}
}

// predeterminedBalancesToSolidity converts PredeterminedBalances to Solidity struct format.
func predeterminedBalancesToSolidity(p *tokenizationtypes.PredeterminedBalances) []interface{} {
	if p == nil {
		return []interface{}{
			make([]interface{}, 0),
			incrementedBalancesToSolidity(nil),
			predeterminedOrderCalculationMethodToSolidity(nil),
		}
	}
	manual := make([]interface{}, 0, len(p.ManualBalances))
	for _, m := range p.ManualBalances {
		manual = append(manual, manualBalancesToSolidity(m))
	}
	return []interface{}{
		manual,
		incrementedBalancesToSolidity(p.IncrementedBalances),
		predeterminedOrderCalculationMethodToSolidity(p.OrderCalculationMethod),
	}
}

// approvalCriteriaToSolidity converts ApprovalCriteria to Solidity struct format (full struct).
func approvalCriteriaToSolidity(a *tokenizationtypes.ApprovalCriteria) []interface{} {
	if a == nil {
		return approvalCriteriaEmptySolidity()
	}
	merkle := make([]interface{}, 0, len(a.MerkleChallenges))
	for _, m := range a.MerkleChallenges {
		merkle = append(merkle, merkleChallengeToSolidity(m))
	}
	coinTransfers := make([]interface{}, 0, len(a.CoinTransfers))
	for _, c := range a.CoinTransfers {
		coinTransfers = append(coinTransfers, coinTransferToSolidity(c))
	}
	mustOwn := make([]interface{}, 0, len(a.MustOwnTokens))
	for _, m := range a.MustOwnTokens {
		mustOwn = append(mustOwn, mustOwnTokensToSolidity(m))
	}
	dynStore := make([]interface{}, 0, len(a.DynamicStoreChallenges))
	for _, d := range a.DynamicStoreChallenges {
		dynStore = append(dynStore, dynamicStoreChallengeToSolidity(d))
	}
	ethSig := make([]interface{}, 0, len(a.EthSignatureChallenges))
	for _, e := range a.EthSignatureChallenges {
		ethSig = append(ethSig, ethSignatureChallengeToSolidity(e))
	}
	voting := make([]interface{}, 0, len(a.VotingChallenges))
	for _, v := range a.VotingChallenges {
		voting = append(voting, votingChallengeToSolidity(v))
	}
	// Note: Solidity ApprovalCriteria struct does not include evmQueryChallenges; only CollectionInvariants does.

	return []interface{}{
		merkle,
		predeterminedBalancesToSolidity(a.PredeterminedBalances),
		approvalAmountsToSolidity(a.ApprovalAmounts),
		maxNumTransfersToSolidity(a.MaxNumTransfers),
		coinTransfers,
		a.RequireToEqualsInitiatedBy,
		a.RequireFromEqualsInitiatedBy,
		a.RequireToDoesNotEqualInitiatedBy,
		a.RequireFromDoesNotEqualInitiatedBy,
		a.OverridesFromOutgoingApprovals,
		a.OverridesToIncomingApprovals,
		autoDeletionOptionsToSolidity(a.AutoDeletionOptions),
		userRoyaltiesToSolidity(a.UserRoyalties),
		mustOwn,
		dynStore,
		ethSig,
		addressChecksToSolidity(a.SenderChecks),
		addressChecksToSolidity(a.RecipientChecks),
		addressChecksToSolidity(a.InitiatorChecks),
		altTimeChecksToSolidity(a.AltTimeChecks),
		a.MustPrioritize,
		voting,
		a.AllowBackedMinting,
		a.AllowSpecialWrapping,
	}
}

func approvalCriteriaEmptySolidity() []interface{} {
	empty := make([]interface{}, 0)
	return []interface{}{
		empty,
		predeterminedBalancesToSolidity(nil),
		approvalAmountsToSolidity(nil),
		maxNumTransfersToSolidity(nil),
		empty,
		false, false, false, false,
		false, false,
		autoDeletionOptionsToSolidity(nil),
		userRoyaltiesToSolidity(nil),
		empty,
		empty,
		empty,
		addressChecksToSolidity(nil),
		addressChecksToSolidity(nil),
		addressChecksToSolidity(nil),
		altTimeChecksToSolidity(nil),
		false,
		empty,
		false,
		false,
	}
}

// outgoingApprovalCriteriaToSolidity converts OutgoingApprovalCriteria to Solidity struct format.
func outgoingApprovalCriteriaToSolidity(a *tokenizationtypes.OutgoingApprovalCriteria) []interface{} {
	if a == nil {
		empty := make([]interface{}, 0)
		return []interface{}{
			empty, predeterminedBalancesToSolidity(nil), approvalAmountsToSolidity(nil), maxNumTransfersToSolidity(nil),
			empty, false, false, autoDeletionOptionsToSolidity(nil),
			empty, empty, empty,
			addressChecksToSolidity(nil), addressChecksToSolidity(nil),
			altTimeChecksToSolidity(nil), false, empty,
		}
	}
	merkle := make([]interface{}, 0, len(a.MerkleChallenges))
	for _, m := range a.MerkleChallenges {
		merkle = append(merkle, merkleChallengeToSolidity(m))
	}
	coinTransfers := make([]interface{}, 0, len(a.CoinTransfers))
	for _, c := range a.CoinTransfers {
		coinTransfers = append(coinTransfers, coinTransferToSolidity(c))
	}
	mustOwn := make([]interface{}, 0, len(a.MustOwnTokens))
	for _, m := range a.MustOwnTokens {
		mustOwn = append(mustOwn, mustOwnTokensToSolidity(m))
	}
	dynStore := make([]interface{}, 0, len(a.DynamicStoreChallenges))
	for _, d := range a.DynamicStoreChallenges {
		dynStore = append(dynStore, dynamicStoreChallengeToSolidity(d))
	}
	ethSig := make([]interface{}, 0, len(a.EthSignatureChallenges))
	for _, e := range a.EthSignatureChallenges {
		ethSig = append(ethSig, ethSignatureChallengeToSolidity(e))
	}
	voting := make([]interface{}, 0, len(a.VotingChallenges))
	for _, v := range a.VotingChallenges {
		voting = append(voting, votingChallengeToSolidity(v))
	}
	return []interface{}{
		merkle,
		predeterminedBalancesToSolidity(a.PredeterminedBalances),
		approvalAmountsToSolidity(a.ApprovalAmounts),
		maxNumTransfersToSolidity(a.MaxNumTransfers),
		coinTransfers,
		a.RequireToEqualsInitiatedBy,
		a.RequireToDoesNotEqualInitiatedBy,
		autoDeletionOptionsToSolidity(a.AutoDeletionOptions),
		mustOwn,
		dynStore,
		ethSig,
		addressChecksToSolidity(a.RecipientChecks),
		addressChecksToSolidity(a.InitiatorChecks),
		altTimeChecksToSolidity(a.AltTimeChecks),
		a.MustPrioritize,
		voting,
	}
}

// incomingApprovalCriteriaToSolidity converts IncomingApprovalCriteria to Solidity struct format.
func incomingApprovalCriteriaToSolidity(a *tokenizationtypes.IncomingApprovalCriteria) []interface{} {
	if a == nil {
		empty := make([]interface{}, 0)
		return []interface{}{
			empty, predeterminedBalancesToSolidity(nil), approvalAmountsToSolidity(nil), maxNumTransfersToSolidity(nil),
			empty, false, false, autoDeletionOptionsToSolidity(nil),
			empty, empty, empty,
			addressChecksToSolidity(nil), addressChecksToSolidity(nil),
			altTimeChecksToSolidity(nil), false, empty,
		}
	}
	merkle := make([]interface{}, 0, len(a.MerkleChallenges))
	for _, m := range a.MerkleChallenges {
		merkle = append(merkle, merkleChallengeToSolidity(m))
	}
	coinTransfers := make([]interface{}, 0, len(a.CoinTransfers))
	for _, c := range a.CoinTransfers {
		coinTransfers = append(coinTransfers, coinTransferToSolidity(c))
	}
	mustOwn := make([]interface{}, 0, len(a.MustOwnTokens))
	for _, m := range a.MustOwnTokens {
		mustOwn = append(mustOwn, mustOwnTokensToSolidity(m))
	}
	dynStore := make([]interface{}, 0, len(a.DynamicStoreChallenges))
	for _, d := range a.DynamicStoreChallenges {
		dynStore = append(dynStore, dynamicStoreChallengeToSolidity(d))
	}
	ethSig := make([]interface{}, 0, len(a.EthSignatureChallenges))
	for _, e := range a.EthSignatureChallenges {
		ethSig = append(ethSig, ethSignatureChallengeToSolidity(e))
	}
	voting := make([]interface{}, 0, len(a.VotingChallenges))
	for _, v := range a.VotingChallenges {
		voting = append(voting, votingChallengeToSolidity(v))
	}
	return []interface{}{
		merkle,
		predeterminedBalancesToSolidity(a.PredeterminedBalances),
		approvalAmountsToSolidity(a.ApprovalAmounts),
		maxNumTransfersToSolidity(a.MaxNumTransfers),
		coinTransfers,
		a.RequireFromEqualsInitiatedBy,
		a.RequireFromDoesNotEqualInitiatedBy,
		autoDeletionOptionsToSolidity(a.AutoDeletionOptions),
		mustOwn,
		dynStore,
		ethSig,
		addressChecksToSolidity(a.SenderChecks),
		addressChecksToSolidity(a.InitiatorChecks),
		altTimeChecksToSolidity(a.AltTimeChecks),
		a.MustPrioritize,
		voting,
	}
}

// userOutgoingApprovalsToSolidity converts []*UserOutgoingApproval to Solidity UserOutgoingApproval[] (each tuple has 10 elements).
func userOutgoingApprovalsToSolidity(approvals []*tokenizationtypes.UserOutgoingApproval) []interface{} {
	out := make([]interface{}, 0, len(approvals))
	for _, app := range approvals {
		if app == nil {
			continue
		}
		transferTimes := uintRangesToSolidity(app.TransferTimes)
		tokenIds := uintRangesToSolidity(app.TokenIds)
		ownershipTimes := uintRangesToSolidity(app.OwnershipTimes)
		version := big.NewInt(0)
		if !app.Version.IsNil() {
			version = app.Version.BigInt()
		}
		out = append(out, []interface{}{
			app.ApprovalId,
			app.ToListId,
			app.InitiatedByListId,
			transferTimes,
			tokenIds,
			ownershipTimes,
			app.Uri,
			app.CustomData,
			outgoingApprovalCriteriaToSolidity(app.ApprovalCriteria),
			version,
		})
	}
	return out
}

// userIncomingApprovalsToSolidity converts []*UserIncomingApproval to Solidity UserIncomingApproval[] (each tuple has 10 elements).
func userIncomingApprovalsToSolidity(approvals []*tokenizationtypes.UserIncomingApproval) []interface{} {
	out := make([]interface{}, 0, len(approvals))
	for _, app := range approvals {
		if app == nil {
			continue
		}
		transferTimes := uintRangesToSolidity(app.TransferTimes)
		tokenIds := uintRangesToSolidity(app.TokenIds)
		ownershipTimes := uintRangesToSolidity(app.OwnershipTimes)
		version := big.NewInt(0)
		if !app.Version.IsNil() {
			version = app.Version.BigInt()
		}
		out = append(out, []interface{}{
			app.ApprovalId,
			app.FromListId,
			app.InitiatedByListId,
			transferTimes,
			tokenIds,
			ownershipTimes,
			app.Uri,
			app.CustomData,
			incomingApprovalCriteriaToSolidity(app.ApprovalCriteria),
			version,
		})
	}
	return out
}
