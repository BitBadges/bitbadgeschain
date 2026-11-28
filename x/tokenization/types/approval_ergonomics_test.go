package types_test

import (
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestApprovalCriteriaCastReflectionRoundTrip guards against the C-1 class of bugs
// (e.g., omitted EvmQueryChallenges in approved_transfers_casts.go) by reflecting
// over EVERY exported field of ApprovalCriteria, OutgoingApprovalCriteria, and
// IncomingApprovalCriteria and asserting 100% preservation across all 4 cast functions.
func TestApprovalCriteriaCastReflectionRoundTrip(t *testing.T) {
	// Explicit allowlist of fields intentionally absent from Outgoing/Incoming structs.
	// If a new field is added to ApprovalCriteria, it MUST either be added to the user
	// structs + casts OR explicitly categorized here.
	collectionOnlyFields := map[string]bool{
		"OverridesFromOutgoingApprovals": true,
		"OverridesToIncomingApprovals":   true,
		"AllowBackedMinting":             true,
		"AllowSpecialWrapping":           true,
		"UserApprovalSettings":           true,
	}
	outgoingOnlyVsIncoming := map[string]bool{
		"RequireToEqualsInitiatedBy":       true,
		"RequireToDoesNotEqualInitiatedBy": true,
		"RecipientChecks":                  true,
	}
	incomingOnlyVsOutgoing := map[string]bool{
		"RequireFromEqualsInitiatedBy":       true,
		"RequireFromDoesNotEqualInitiatedBy": true,
		"SenderChecks":                       true,
	}

	acType := reflect.TypeOf(types.ApprovalCriteria{})
	outType := reflect.TypeOf(types.OutgoingApprovalCriteria{})
	inType := reflect.TypeOf(types.IncomingApprovalCriteria{})

	// Verify structural partition: every ApprovalCriteria field must appear in Outgoing,
	// Incoming, or collectionOnlyFields.
	for i := 0; i < acType.NumField(); i++ {
		f := acType.Field(i)
		if !f.IsExported() {
			continue
		}
		_, inOut := outType.FieldByName(f.Name)
		_, inIn := inType.FieldByName(f.Name)

		switch {
		case collectionOnlyFields[f.Name]:
			require.False(t, inOut, "collection-only field %s unexpectedly on OutgoingApprovalCriteria", f.Name)
			require.False(t, inIn, "collection-only field %s unexpectedly on IncomingApprovalCriteria", f.Name)
		case outgoingOnlyVsIncoming[f.Name]:
			require.True(t, inOut, "outgoing field %s missing from OutgoingApprovalCriteria", f.Name)
			require.False(t, inIn, "outgoing-only field %s unexpectedly on IncomingApprovalCriteria", f.Name)
		case incomingOnlyVsOutgoing[f.Name]:
			require.False(t, inOut, "incoming-only field %s unexpectedly on OutgoingApprovalCriteria", f.Name)
			require.True(t, inIn, "incoming field %s missing from IncomingApprovalCriteria", f.Name)
		default:
			require.True(t, inOut, "shared field %s missing from OutgoingApprovalCriteria", f.Name)
			require.True(t, inIn, "shared field %s missing from IncomingApprovalCriteria", f.Name)
		}
	}

	// Build a fully populated ApprovalCriteria where EVERY SINGLE exported field is non-zero.
	coin := sdk.NewInt64Coin("ubadge", 42)
	fullCollection := &types.ApprovalCriteria{
		MerkleChallenges: []*types.MerkleChallenge{{Root: "0xabc", ExpectedProofLength: sdkmath.NewUint(2), MaxUsesPerLeaf: sdkmath.NewUint(1), ChallengeTrackerId: "mc1"}},
		PredeterminedBalances: &types.PredeterminedBalances{
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
			ManualBalances:         []*types.ManualBalances{{Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: singleRange(1), OwnershipTimes: fullRanges()}}}},
		},
		ApprovalAmounts: &types.ApprovalAmounts{OverallApprovalAmount: sdkmath.NewUint(10), AmountTrackerId: "amt1"},
		MaxNumTransfers: &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(5), AmountTrackerId: "cnt1"},
		CoinTransfers:   []*types.CoinTransfer{{To: "bb1dest", Coins: []*sdk.Coin{&coin}}},
		RequireToEqualsInitiatedBy:         true,
		RequireFromEqualsInitiatedBy:       true,
		RequireToDoesNotEqualInitiatedBy:   true,
		RequireFromDoesNotEqualInitiatedBy: true,
		OverridesFromOutgoingApprovals:     true,
		OverridesToIncomingApprovals:       true,
		AutoDeletionOptions:                &types.AutoDeletionOptions{AfterOneUse: true},
		MustOwnTokens:                      []*types.MustOwnTokens{{CollectionId: sdkmath.NewUint(7), AmountRange: &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}},
		DynamicStoreChallenges:             []*types.DynamicStoreChallenge{{StoreId: sdkmath.NewUint(3)}},
		EthSignatureChallenges:             []*types.ETHSignatureChallenge{{Signer: "0x0000000000000000000000000000000000000001"}},
		SenderChecks:                       &types.AddressChecks{MustNotBeEvmContract: true},
		RecipientChecks:                    &types.AddressChecks{MustBeEvmContract: true},
		InitiatorChecks:                    &types.AddressChecks{MustNotBeLiquidityPool: true},
		AltTimeChecks:                      &types.AltTimeChecks{TimezoneOffsetMinutes: sdkmath.NewUint(60)},
		MustPrioritize:                     true,
		VotingChallenges:                   []*types.VotingChallenge{{ProposalId: "prop-1", QuorumThreshold: sdkmath.NewUint(51)}},
		AllowBackedMinting:                 true,
		AllowSpecialWrapping:               true,
		EvmQueryChallenges:                 []*types.EVMQueryChallenge{{ContractAddress: "0x0000000000000000000000000000000000000002", Calldata: "0x12345678", GasLimit: sdkmath.NewUint(100000)}},
		UserApprovalSettings:               &types.UserApprovalSettings{DisableUserCoinTransfers: true},
	}

	// Assert our test fixture actually populated EVERY exported field on ApprovalCriteria!
	fullVal := reflect.ValueOf(fullCollection).Elem()
	for i := 0; i < acType.NumField(); i++ {
		f := acType.Field(i)
		if !f.IsExported() {
			continue
		}
		require.False(t, fullVal.Field(i).IsZero(), "test fixture left ApprovalCriteria.%s zero", f.Name)
	}

	// Cast 1 & 2: Collection -> Outgoing -> Collection
	outCriteria := types.CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(fullCollection)
	require.NotNil(t, outCriteria)
	outVal := reflect.ValueOf(outCriteria).Elem()
	for i := 0; i < outType.NumField(); i++ {
		f := outType.Field(i)
		if !f.IsExported() {
			continue
		}
		got := outVal.Field(i).Interface()
		want := fullVal.FieldByName(f.Name).Interface()
		require.False(t, outVal.Field(i).IsZero(), "CastFromCollection->Outgoing dropped field %s", f.Name)
		require.Equal(t, want, got, "CastFromCollection->Outgoing mismatch on field %s", f.Name)
	}

	backFromOut := types.CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(outCriteria)
	require.NotNil(t, backFromOut)
	backOutVal := reflect.ValueOf(backFromOut).Elem()
	for i := 0; i < outType.NumField(); i++ {
		f := outType.Field(i)
		if !f.IsExported() {
			continue
		}
		got := backOutVal.FieldByName(f.Name).Interface()
		want := outVal.Field(i).Interface()
		require.False(t, backOutVal.FieldByName(f.Name).IsZero(), "CastOutgoing->Collection dropped field %s", f.Name)
		require.Equal(t, want, got, "CastOutgoing->Collection mismatch on field %s", f.Name)
	}

	// Cast 3 & 4: Collection -> Incoming -> Collection
	inCriteria := types.CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(fullCollection)
	require.NotNil(t, inCriteria)
	inVal := reflect.ValueOf(inCriteria).Elem()
	for i := 0; i < inType.NumField(); i++ {
		f := inType.Field(i)
		if !f.IsExported() {
			continue
		}
		got := inVal.Field(i).Interface()
		want := fullVal.FieldByName(f.Name).Interface()
		require.False(t, inVal.Field(i).IsZero(), "CastFromCollection->Incoming dropped field %s", f.Name)
		require.Equal(t, want, got, "CastFromCollection->Incoming mismatch on field %s", f.Name)
	}

	backFromIn := types.CastIncomingApprovalCriteriaToCollectionApprovalCriteria(inCriteria)
	require.NotNil(t, backFromIn)
	backInVal := reflect.ValueOf(backFromIn).Elem()
	for i := 0; i < inType.NumField(); i++ {
		f := inType.Field(i)
		if !f.IsExported() {
			continue
		}
		got := backInVal.FieldByName(f.Name).Interface()
		want := inVal.Field(i).Interface()
		require.False(t, backInVal.FieldByName(f.Name).IsZero(), "CastIncoming->Collection dropped field %s", f.Name)
		require.Equal(t, want, got, "CastIncoming->Collection mismatch on field %s", f.Name)
	}

	// Nil-safety checks for all transfer and criteria casts.
	require.Nil(t, types.CastOutgoingTransferToCollectionTransfer(nil, "bb1from"))
	require.Nil(t, types.CastFromCollectionTransferToOutgoingTransfer(nil))
	require.Nil(t, types.CastIncomingTransferToCollectionTransfer(nil, "bb1to"))
	require.Nil(t, types.CastFromCollectionTransferToIncomingTransfer(nil))
	require.Empty(t, types.CastOutgoingTransfersToCollectionTransfers([]*types.UserOutgoingApproval{nil}, "bb1from"))
	require.Empty(t, types.CastIncomingTransfersToCollectionTransfers([]*types.UserIncomingApproval{nil}, "bb1to"))
}

func TestValidateApprovalCriteriaErgonomics_MutualExclusivity(t *testing.T) {
	t.Run("rejects PredeterminedBalances + ApprovalAmounts", func(t *testing.T) {
		criteria := &types.ApprovalCriteria{
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					DurationFromTimestamp: sdkmath.NewUint(1000),
				},
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "t1",
			},
		}
		err := types.ValidateApprovalCriteriaErgonomics(criteria)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mutually exclusive")
	})

	t.Run("rejects DurationFromTimestamp + RecurringOwnershipTimes", func(t *testing.T) {
		criteria := &types.ApprovalCriteria{
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					DurationFromTimestamp: sdkmath.NewUint(3600000),
					RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
						StartTime:          sdkmath.NewUint(1000),
						IntervalLength:     sdkmath.NewUint(3600000),
						ChargePeriodLength: sdkmath.NewUint(60000),
					},
				},
			},
		}
		err := types.ValidateApprovalCriteriaErgonomics(criteria)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mutually exclusive")
	})

	t.Run("rejects contradictory party-relation booleans", func(t *testing.T) {
		err := types.ValidateApprovalCriteriaErgonomics(&types.ApprovalCriteria{
			RequireToEqualsInitiatedBy:       true,
			RequireToDoesNotEqualInitiatedBy: true,
		})
		require.Error(t, err)

		err = types.ValidateApprovalCriteriaErgonomics(&types.ApprovalCriteria{
			RequireFromEqualsInitiatedBy:       true,
			RequireFromDoesNotEqualInitiatedBy: true,
		})
		require.Error(t, err)
	})

	t.Run("rejects contradictory AddressChecks", func(t *testing.T) {
		err := types.ValidateApprovalCriteriaErgonomics(&types.ApprovalCriteria{
			SenderChecks: &types.AddressChecks{
				MustBeEvmContract:    true,
				MustNotBeEvmContract: true,
			},
		})
		require.Error(t, err)
	})

	t.Run("allows valid non-conflicting criteria", func(t *testing.T) {
		criteria := &types.ApprovalCriteria{
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					DurationFromTimestamp: sdkmath.NewUint(86400000),
				},
			},
			RequireToEqualsInitiatedBy: true,
		}
		require.NoError(t, types.ValidateApprovalCriteriaErgonomics(criteria))
		require.NoError(t, types.ValidateApprovalCriteriaErgonomics(nil))
	})

	t.Run("validates user outgoing and incoming tiers identically", func(t *testing.T) {
		outErr := types.ValidateOutgoingApprovalCriteriaErgonomics(&types.OutgoingApprovalCriteria{
			RequireToEqualsInitiatedBy:       true,
			RequireToDoesNotEqualInitiatedBy: true,
		})
		require.Error(t, outErr)

		inErr := types.ValidateIncomingApprovalCriteriaErgonomics(&types.IncomingApprovalCriteria{
			RequireFromEqualsInitiatedBy:       true,
			RequireFromDoesNotEqualInitiatedBy: true,
		})
		require.Error(t, inErr)
	})
}

func TestValidateAddressListPredicate(t *testing.T) {
	valid := []string{
		"All",
		"AllWithoutMint",
		"Mint",
		"!Mint",
		"!(bb1a:bb1b)",
		"bb1a:bb1b",
		"customList123",
	}
	for _, id := range valid {
		require.NoError(t, types.ValidateAddressListPredicate(id), "expected %q to be valid", id)
	}

	invalid := []string{
		"",
		"!",
		"!!Mint",
		"!()",
		"!(Mint",
		"bb1a::bb1b",
		":bb1a",
		"bb1a:",
		"bb1a:!bb1b",
		"All With Space",
		"AllWithout",
	}
	for _, id := range invalid {
		require.Error(t, types.ValidateAddressListPredicate(id), "expected %q to be rejected", id)
	}
}

func TestNormalizeCollectionApprovalDefaultsAndDefaultBalances(t *testing.T) {
	t.Run("Mint approval automatically sets OverridesFromOutgoingApprovals", func(t *testing.T) {
		mintApproval := &types.CollectionApproval{
			ApprovalId:        "mint",
			FromListId:        types.MintAddress,
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
		}
		types.NormalizeCollectionApprovalDefaults(mintApproval)
		require.NotNil(t, mintApproval.ApprovalCriteria)
		require.True(t, mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals)
	})

	t.Run("AllowBackedMinting automatically enforces MustPrioritize", func(t *testing.T) {
		backed := &types.CollectionApproval{
			ApprovalId: "deposit",
			FromListId: "bb1backing",
			ToListId:   "AllWithoutMint",
			ApprovalCriteria: &types.ApprovalCriteria{
				AllowBackedMinting: true,
				MustPrioritize:     false,
			},
		}
		types.NormalizeCollectionApprovalDefaults(backed)
		require.True(t, backed.ApprovalCriteria.MustPrioritize)
		require.False(t, backed.ApprovalCriteria.OverridesFromOutgoingApprovals)
	})

	t.Run("NormalizeDefaultBalances defaults AutoApproveAllIncomingTransfers to true unless opted out", func(t *testing.T) {
		def := types.NormalizeDefaultBalances(nil, false)
		require.True(t, def.AutoApproveAllIncomingTransfers)
		require.True(t, def.AutoApproveSelfInitiatedOutgoingTransfers)
		require.True(t, def.AutoApproveSelfInitiatedIncomingTransfers)

		optOut := types.NormalizeDefaultBalances(nil, true)
		require.False(t, optOut.AutoApproveAllIncomingTransfers)
		require.True(t, optOut.AutoApproveSelfInitiatedOutgoingTransfers)
		require.True(t, optOut.AutoApproveSelfInitiatedIncomingTransfers)
	})
}
