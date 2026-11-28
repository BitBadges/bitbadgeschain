package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func fullRanges() []*types.UintRange {
	return []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(types.MaxUint64Value)}}
}

func singleRange(id uint64) []*types.UintRange {
	return []*types.UintRange{{Start: sdkmath.NewUint(id), End: sdkmath.NewUint(id)}}
}

func frozenAllPermissions() *types.CollectionPermissions {
	return &types.CollectionPermissions{
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
			{
				FromListId:                "All",
				ToListId:                  "All",
				InitiatedByListId:         "All",
				TransferTimes:             fullRanges(),
				TokenIds:                  fullRanges(),
				OwnershipTimes:            fullRanges(),
				ApprovalId:                "All",
				PermanentlyForbiddenTimes: fullRanges(),
			},
		},
	}
}

func frozenPostMintPermissions() *types.CollectionPermissions {
	return &types.CollectionPermissions{
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
			{
				FromListId:                "AllWithoutMint",
				ToListId:                  "AllWithoutMint",
				InitiatedByListId:         "AllWithoutMint",
				TransferTimes:             fullRanges(),
				TokenIds:                  fullRanges(),
				OwnershipTimes:            fullRanges(),
				ApprovalId:                "All",
				PermanentlyForbiddenTimes: fullRanges(),
			},
		},
	}
}

func sampleBackedPath(addr string) *types.CosmosCoinBackedPath {
	return &types.CosmosCoinBackedPath{
		Address: addr,
		Conversion: &types.Conversion{
			SideA: &types.ConversionSideAWithDenom{
				Amount: sdkmath.NewUint(1),
				Denom:  "uusdc",
			},
			SideB: []*types.Balance{{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       singleRange(1),
				OwnershipTimes: fullRanges(),
			}},
		},
	}
}

func sampleCoin(denom string, amt int64) *sdk.Coin {
	c := sdk.NewInt64Coin(denom, amt)
	return &c
}

func TestVerifyStandards_AllTwelvePassAndRejectSpoofs(t *testing.T) {
	const (
		backingAddr = "bb1backing000000000000000000000000000000000"
		escrowAddr  = "bb1escrow0000000000000000000000000000000000"
		merchant    = "bb1merchant00000000000000000000000000000000"
		worker      = "bb1worker0000000000000000000000000000000000"
	)

	cases := []struct {
		name           string
		standard       string
		validCol       *types.TokenCollection
		spoofedCol     *types.TokenCollection
		wantSpoofCodes []string
	}{
		{
			name:     "Non-Transferable",
			standard: types.StandardNonTransferable,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardNonTransferable},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenPostMintPermissions(),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "mint-only",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: merchant,
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardNonTransferable},
				ValidTokenIds: singleRange(1),
				// Unfrozen approvals AND a post-mint transferable approval
				CollectionPermissions: &types.CollectionPermissions{},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "p2p-transfer",
					FromListId:        "AllWithoutMint",
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			wantSpoofCodes: []string{"NON_TRANSFERABLE_POST_MINT_APPROVAL", "NON_TRANSFERABLE_UNFROZEN_APPROVALS"},
		},
		{
			name:     "Subscriptions",
			standard: types.StandardSubscriptions,
			validCol: &types.TokenCollection{
				Standards:     []string{types.StandardSubscriptions},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "sub-mint",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						PredeterminedBalances: &types.PredeterminedBalances{
							IncrementedBalances: &types.IncrementedBalances{
								DurationFromTimestamp: sdkmath.NewUint(2592000000),
							},
						},
					},
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardSubscriptions},
				ValidTokenIds: singleRange(1),
				Invariants:    &types.CollectionInvariants{NoCustomOwnershipTimes: true},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "plain-mint",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			wantSpoofCodes: []string{"SUBSCRIPTION_NO_CUSTOM_OWNERSHIP_TIMES", "SUBSCRIPTION_MISSING_TIME_BOUNDED_APPROVAL"},
		},
		{
			name:     "NFTs",
			standard: types.StandardNFTs,
			validCol: &types.TokenCollection{
				Standards:     []string{types.StandardNFTs},
				ValidTokenIds: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
				Invariants: &types.CollectionInvariants{
					MaxSupplyPerId:         sdkmath.NewUint(1),
					NoCustomOwnershipTimes: true,
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardNFTs},
				ValidTokenIds: singleRange(1),
				Invariants: &types.CollectionInvariants{
					MaxSupplyPerId:         sdkmath.NewUint(500),
					NoCustomOwnershipTimes: false,
				},
			},
			wantSpoofCodes: []string{"NFT_MAX_SUPPLY_NOT_ONE", "NFT_CUSTOM_OWNERSHIP_TIMES_ALLOWED"},
		},
		{
			name:     "Fungible Tokens",
			standard: types.StandardFungibleTokens,
			validCol: &types.TokenCollection{
				Standards:     []string{types.StandardFungibleTokens},
				ValidTokenIds: singleRange(1),
				Invariants: &types.CollectionInvariants{
					MaxSupplyPerId:         sdkmath.NewUint(1000000),
					NoCustomOwnershipTimes: true,
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardFungibleTokens},
				ValidTokenIds: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
				Invariants: &types.CollectionInvariants{
					MaxSupplyPerId:         sdkmath.NewUint(1),
					NoCustomOwnershipTimes: true,
				},
			},
			wantSpoofCodes: []string{"FUNGIBLE_REQUIRES_SINGLE_TOKEN_ID", "FUNGIBLE_MAX_SUPPLY_IS_ONE"},
		},
		{
			name:     "Smart Token",
			standard: types.StandardSmartToken,
			validCol: &types.TokenCollection{
				Standards:     []string{types.StandardSmartToken},
				ValidTokenIds: singleRange(1),
				Invariants:    &types.CollectionInvariants{CosmosCoinBackedPath: sampleBackedPath(backingAddr)},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "backed-deposit",
					FromListId:        backingAddr,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						AllowBackedMinting: true,
						MustPrioritize:     true,
					},
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardSmartToken},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "rogue-mint",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			wantSpoofCodes: []string{"SMART_TOKEN_MISSING_BACKED_PATH", "SMART_TOKEN_UNBACKED_MINT_APPROVAL"},
		},
		{
			name:     "Address List",
			standard: types.StandardAddressList,
			validCol: &types.TokenCollection{
				Standards:     []string{types.StandardAddressList},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "allowlist-mint",
					FromListId:        types.MintAddress,
					ToListId:          "vip_members_list",
					InitiatedByListId: merchant,
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardAddressList},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "open-mint",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			wantSpoofCodes: []string{"ADDRESS_LIST_NO_SCOPED_PREDICATE"},
		},
		{
			name:     "Vault",
			standard: types.StandardVault,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardVault},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenAllPermissions(),
				Invariants: &types.CollectionInvariants{
					CosmosCoinBackedPath:       sampleBackedPath(backingAddr),
					NoForcefulPostMintTransfers: true,
				},
				CollectionApprovals: []*types.CollectionApproval{
					{
						ApprovalId:        "deposit",
						FromListId:        backingAddr,
						ToListId:          "AllWithoutMint",
						InitiatedByListId: "AllWithoutMint",
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria:  &types.ApprovalCriteria{AllowBackedMinting: true, MustPrioritize: true},
					},
					{
						ApprovalId:        "withdraw",
						FromListId:        "AllWithoutMint",
						ToListId:          backingAddr,
						InitiatedByListId: "AllWithoutMint",
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria:  &types.ApprovalCriteria{AllowBackedMinting: true, MustPrioritize: true},
					},
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardVault},
				ValidTokenIds: singleRange(1),
				Invariants:    &types.CollectionInvariants{CosmosCoinBackedPath: sampleBackedPath(backingAddr)},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "deposit-only-honeypot",
					FromListId:        backingAddr,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria:  &types.ApprovalCriteria{AllowBackedMinting: true, MustPrioritize: true},
				}},
			},
			wantSpoofCodes: []string{"VAULT_MISSING_WITHDRAW_APPROVAL", "VAULT_FORCEFUL_TRANSFERS_ALLOWED", "VAULT_UNFROZEN_APPROVALS"},
		},
		{
			name:     "Agent Wallet",
			standard: types.StandardAgentWallet,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardAgentWallet},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenAllPermissions(),
				Invariants:            &types.CollectionInvariants{NoForcefulPostMintTransfers: true},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "agent-daily-spend",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: worker,
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{
							To:                              merchant,
							OverrideFromWithApproverAddress: true,
							Coins:                           []*sdk.Coin{sampleCoin("uusdc", 10_000_000)},
						}},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(5),
							AmountTrackerId:        "agent-cap",
						},
					},
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardAgentWallet},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "unbounded-drain",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: worker,
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{
							To:                              merchant,
							OverrideFromWithApproverAddress: true,
							Coins:                           []*sdk.Coin{sampleCoin("uusdc", 10_000_000)},
						}},
					},
				}},
			},
			wantSpoofCodes: []string{"AGENT_WALLET_UNBOUNDED_SPEND", "AGENT_WALLET_UNPROTECTED_CUSTODY"},
		},
		{
			name:     "PaymentRequestV2",
			standard: types.StandardPaymentRequestV2,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardPaymentRequestV2},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenAllPermissions(),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "pay-invoice",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{
							To:    merchant,
							Coins: []*sdk.Coin{sampleCoin("uusdc", 25_000_000)},
						}},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1),
							AmountTrackerId:        "inv-1",
						},
					},
				}},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardPaymentRequestV2},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "repeatable-invoice",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{
							To:    merchant,
							Coins: []*sdk.Coin{sampleCoin("uusdc", 25_000_000)},
						}},
					},
				}},
			},
			wantSpoofCodes: []string{"PAYMENT_REQUEST_MISSING_BOUNDED_PAYMENT", "PAYMENT_REQUEST_UNFROZEN_APPROVALS"},
		},
		{
			name:     "Bounty",
			standard: types.StandardBounty,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardBounty},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenAllPermissions(),
				CollectionApprovals: []*types.CollectionApproval{
					{
						ApprovalId:        "fulfill",
						FromListId:        types.MintAddress,
						ToListId:          "AllWithoutMint",
						InitiatedByListId: merchant,
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:                              worker,
								OverrideFromWithApproverAddress: true,
								Coins:                           []*sdk.Coin{sampleCoin("uusdc", 100_000_000)},
							}},
							MaxNumTransfers: &types.MaxNumTransfers{
								OverallMaxNumTransfers: sdkmath.NewUint(1),
								AmountTrackerId:        "bounty-latch",
							},
						},
					},
					{
						ApprovalId:        "refund",
						FromListId:        types.MintAddress,
						ToListId:          merchant,
						InitiatedByListId: merchant,
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:                              merchant,
								OverrideFromWithApproverAddress: true,
								Coins:                           []*sdk.Coin{sampleCoin("uusdc", 100_000_000)},
							}},
							MaxNumTransfers: &types.MaxNumTransfers{
								OverallMaxNumTransfers: sdkmath.NewUint(1),
								AmountTrackerId:        "bounty-latch",
							},
						},
					},
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardBounty},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{
					{
						ApprovalId:        "fulfill",
						FromListId:        types.MintAddress,
						ToListId:          "AllWithoutMint",
						InitiatedByListId: merchant,
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:                              worker,
								OverrideFromWithApproverAddress: true,
								Coins:                           []*sdk.Coin{sampleCoin("uusdc", 100_000_000)},
							}},
							MaxNumTransfers: &types.MaxNumTransfers{
								OverallMaxNumTransfers: sdkmath.NewUint(1),
								AmountTrackerId:        "tracker-a", // DIFFERENT tracker IDs -> no mutual exclusion!
							},
						},
					},
					{
						ApprovalId:        "refund",
						FromListId:        types.MintAddress,
						ToListId:          merchant,
						InitiatedByListId: merchant,
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							MaxNumTransfers: &types.MaxNumTransfers{
								OverallMaxNumTransfers: sdkmath.NewUint(1),
								AmountTrackerId:        "tracker-b",
							},
						},
					},
				},
			},
			wantSpoofCodes: []string{"BOUNTY_MISSING_MUTUAL_EXCLUSION", "BOUNTY_UNFROZEN_APPROVALS"},
		},
		{
			name:     "Crowdfund",
			standard: types.StandardCrowdfund,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardCrowdfund},
				ValidTokenIds:         singleRange(1),
				MintEscrowAddress:     escrowAddr,
				CollectionPermissions: frozenAllPermissions(),
				CollectionApprovals: []*types.CollectionApproval{
					{
						ApprovalId:        "contribute",
						FromListId:        types.MintAddress,
						ToListId:          "AllWithoutMint",
						InitiatedByListId: "AllWithoutMint",
						TokenIds:          singleRange(1),
						TransferTimes:     []*types.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)}},
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:    types.MintAddress,
								Coins: []*sdk.Coin{sampleCoin("uusdc", 10_000_000)},
							}},
						},
					},
					{
						ApprovalId:        "release",
						FromListId:        types.MintAddress,
						ToListId:          merchant,
						InitiatedByListId: merchant,
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:                              merchant,
								OverrideFromWithApproverAddress: true,
								Coins:                           []*sdk.Coin{sampleCoin("uusdc", 100_000_000)},
							}},
						},
					},
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardCrowdfund},
				ValidTokenIds: singleRange(1),
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "no-deadline-contribute",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(), // infinite deadline!
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{
							To:    types.MintAddress,
							Coins: []*sdk.Coin{sampleCoin("uusdc", 10_000_000)},
						}},
					},
				}},
			},
			wantSpoofCodes: []string{"CROWDFUND_MISSING_BOUNDED_CONTRIBUTION", "CROWDFUND_MISSING_ESCROW_SETTLEMENT", "CROWDFUND_UNFROZEN_APPROVALS"},
		},
		{
			name:     "Spendable Credit",
			standard: types.StandardSpendableCredit,
			validCol: &types.TokenCollection{
				Standards:             []string{types.StandardSpendableCredit},
				ValidTokenIds:         singleRange(1),
				CollectionPermissions: frozenAllPermissions(),
				Invariants: &types.CollectionInvariants{
					NoCustomOwnershipTimes:     true,
					NoForcefulPostMintTransfers: true,
				},
				CollectionApprovals: []*types.CollectionApproval{
					{
						ApprovalId:        "buy-credit",
						FromListId:        types.MintAddress,
						ToListId:          "AllWithoutMint",
						InitiatedByListId: "AllWithoutMint",
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
						ApprovalCriteria: &types.ApprovalCriteria{
							CoinTransfers: []*types.CoinTransfer{{
								To:    merchant,
								Coins: []*sdk.Coin{sampleCoin("uusdc", 1_000_000)},
							}},
						},
					},
					{
						ApprovalId:        "redeem-credit",
						FromListId:        "AllWithoutMint",
						ToListId:          merchant,
						InitiatedByListId: "AllWithoutMint",
						TokenIds:          singleRange(1),
						TransferTimes:     fullRanges(),
						OwnershipTimes:    fullRanges(),
					},
				},
			},
			spoofedCol: &types.TokenCollection{
				Standards:     []string{types.StandardSpendableCredit},
				ValidTokenIds: singleRange(1),
				Invariants: &types.CollectionInvariants{
					NoCustomOwnershipTimes:     true,
					NoForcefulPostMintTransfers: false, // can confiscate user credits!
				},
				CollectionApprovals: []*types.CollectionApproval{{
					ApprovalId:        "buy-only",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
				}},
			},
			wantSpoofCodes: []string{"SPENDABLE_CREDIT_FORCEFUL_TRANSFERS_ALLOWED", "SPENDABLE_CREDIT_MISSING_ISSUE_PATH", "SPENDABLE_CREDIT_MISSING_SPEND_PATH", "SPENDABLE_CREDIT_UNFROZEN_APPROVALS"},
		},
	}

	require.Len(t, cases, len(types.KnownVerifiableStandards()), "test matrix must cover every known standard")

	for _, tc := range cases {
		t.Run(tc.name+"/valid", func(t *testing.T) {
			res := types.VerifyCollectionStandards(tc.validCol)
			require.Equal(t, []string{tc.standard}, res.ClaimedStandards)
			require.Equal(t, []string{tc.standard}, res.VerifiedStandards, "unexpected failures: %+v", res.FailedStandards)
			require.Empty(t, res.FailedStandards)
		})

		t.Run(tc.name+"/spoofed_rejected", func(t *testing.T) {
			res := types.VerifyCollectionStandards(tc.spoofedCol)
			require.Equal(t, []string{tc.standard}, res.ClaimedStandards)
			require.Empty(t, res.VerifiedStandards, "spoofed %s must NOT appear in VerifiedStandards", tc.standard)
			require.Contains(t, res.FailedStandards, tc.standard)

			gotCodes := make([]string, 0, len(res.FailedStandards[tc.standard]))
			for _, v := range res.FailedStandards[tc.standard] {
				gotCodes = append(gotCodes, v.Code)
			}
			for _, wantCode := range tc.wantSpoofCodes {
				require.Contains(t, gotCodes, wantCode)
			}
		})
	}
}

func TestVerifyCollectionStandards_MixedClaimedUnknownAndNil(t *testing.T) {
	// Nil collection must never panic.
	nilRes := types.VerifyCollectionStandards(nil)
	require.Empty(t, nilRes.ClaimedStandards)
	require.Empty(t, nilRes.VerifiedStandards)
	require.Empty(t, nilRes.FailedStandards)

	// Collection claiming one valid, one spoofed, one unknown, and a duplicate.
	col := &types.TokenCollection{
		Standards: []string{
			types.StandardNFTs,
			types.StandardNonTransferable, // spoofed (unfrozen)
			"FakeCustomTag",
			types.StandardNFTs, // duplicate
		},
		ValidTokenIds: singleRange(1),
		Invariants: &types.CollectionInvariants{
			MaxSupplyPerId:         sdkmath.NewUint(1),
			NoCustomOwnershipTimes: true,
		},
	}

	res := types.VerifyCollectionStandards(col)
	require.Equal(t, []string{types.StandardNFTs, types.StandardNonTransferable, "FakeCustomTag"}, res.ClaimedStandards)
	require.Equal(t, []string{types.StandardNFTs}, res.VerifiedStandards)
	require.Contains(t, res.FailedStandards, types.StandardNonTransferable)
	require.Contains(t, res.FailedStandards, "FakeCustomTag")
	require.Equal(t, "UNKNOWN_STANDARD", res.FailedStandards["FakeCustomTag"][0].Code)
}

func TestVerifyNonTransferable_RejectsEarlierPermittedPermission(t *testing.T) {
	// First-match semantics: an earlier PermanentlyPermittedTimes entry overrides
	// a later PermanentlyForbiddenTimes entry, so this MUST fail verification.
	col := &types.TokenCollection{
		Standards:     []string{types.StandardNonTransferable},
		ValidTokenIds: singleRange(1),
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
				{
					FromListId:                "AllWithoutMint",
					ToListId:                  "AllWithoutMint",
					InitiatedByListId:         "AllWithoutMint",
					TransferTimes:             fullRanges(),
					TokenIds:                  fullRanges(),
					OwnershipTimes:            fullRanges(),
					ApprovalId:                "All",
					PermanentlyPermittedTimes: fullRanges(),
				},
				{
					FromListId:                "AllWithoutMint",
					ToListId:                  "AllWithoutMint",
					InitiatedByListId:         "AllWithoutMint",
					TransferTimes:             fullRanges(),
					TokenIds:                  fullRanges(),
					OwnershipTimes:            fullRanges(),
					ApprovalId:                "All",
					PermanentlyForbiddenTimes: fullRanges(),
				},
			},
		},
	}

	res := types.VerifyCollectionStandards(col)
	require.Empty(t, res.VerifiedStandards)
	require.Contains(t, res.FailedStandards, types.StandardNonTransferable)
}

func TestVerifyAdversarialSpoofs_TargetedCarveOutAndUnlatchedBountyDrain(t *testing.T) {
	t.Run("Non-Transferable rejects targeted address carve-out before freeze rule", func(t *testing.T) {
		col := &types.TokenCollection{
			Standards:     []string{types.StandardNonTransferable},
			ValidTokenIds: singleRange(1),
			CollectionPermissions: &types.CollectionPermissions{
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
					{
						FromListId:                "bb1attacker00000000000000000000000000000000",
						ToListId:                  "All",
						InitiatedByListId:         "All",
						TransferTimes:             fullRanges(),
						TokenIds:                  fullRanges(),
						OwnershipTimes:            fullRanges(),
						ApprovalId:                "All",
						PermanentlyPermittedTimes: fullRanges(),
					},
					{
						FromListId:                "AllWithoutMint",
						ToListId:                  "AllWithoutMint",
						InitiatedByListId:         "AllWithoutMint",
						TransferTimes:             fullRanges(),
						TokenIds:                  fullRanges(),
						OwnershipTimes:            fullRanges(),
						ApprovalId:                "All",
						PermanentlyForbiddenTimes: fullRanges(),
					},
				},
			},
		}
		res := types.VerifyCollectionStandards(col)
		require.Empty(t, res.VerifiedStandards)
		require.Contains(t, res.FailedStandards, types.StandardNonTransferable)
	})

	t.Run("Bounty rejects third unlatched escrow drain approval alongside valid pair", func(t *testing.T) {
		col := &types.TokenCollection{
			Standards:             []string{types.StandardBounty},
			ValidTokenIds:         singleRange(1),
			CollectionPermissions: frozenAllPermissions(),
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "fulfill",
					FromListId:        types.MintAddress,
					ToListId:          "AllWithoutMint",
					InitiatedByListId: "bb1merchant",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers:   []*types.CoinTransfer{{To: "bb1worker", OverrideFromWithApproverAddress: true, Coins: []*sdk.Coin{sampleCoin("uusdc", 100)}}},
						MaxNumTransfers: &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1), AmountTrackerId: "latch"},
					},
				},
				{
					ApprovalId:        "refund",
					FromListId:        types.MintAddress,
					ToListId:          "bb1merchant",
					InitiatedByListId: "bb1merchant",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers:   []*types.CoinTransfer{{To: "bb1merchant", OverrideFromWithApproverAddress: true, Coins: []*sdk.Coin{sampleCoin("uusdc", 100)}}},
						MaxNumTransfers: &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1), AmountTrackerId: "latch"},
					},
				},
				{
					ApprovalId:        "sneaky-backdoor-drain",
					FromListId:        types.MintAddress,
					ToListId:          "bb1attacker",
					InitiatedByListId: "bb1attacker",
					TokenIds:          singleRange(1),
					TransferTimes:     fullRanges(),
					OwnershipTimes:    fullRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						CoinTransfers: []*types.CoinTransfer{{To: "bb1attacker", OverrideFromWithApproverAddress: true, Coins: []*sdk.Coin{sampleCoin("uusdc", 100)}}},
					},
				},
			},
		}
		res := types.VerifyCollectionStandards(col)
		require.Empty(t, res.VerifiedStandards)
		require.Contains(t, res.FailedStandards, types.StandardBounty)
	})
}
