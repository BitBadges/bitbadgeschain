package simulation

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// FindAccount find a specific address from an account list
func FindAccount(accs []simtypes.Account, address string) (simtypes.Account, bool) {
	creator, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return simtypes.FindAccount(accs, creator)
}

func GetTimelineTimes(r *rand.Rand, length int) []*types.UintRange {
	timelineTimes := make([]*types.UintRange, 0, 1)
	for i := 0; i < length; i++ {
		timelineTimes = append(timelineTimes, &types.UintRange{
			Start: sdkmath.NewUint(r.Uint64()),
			End:  sdkmath.NewUint(r.Uint64()),
		})
	}
	return timelineTimes
}

func GetRandomBalances(r *rand.Rand, length int) []*types.Balance {
	randomSubassets := []*types.Balance{}
	for i := 0; i < r.Intn(length); i++ {
		randomSubassets = append(randomSubassets, &types.Balance{
			Amount: sdkmath.NewUint(r.Uint64()),
			BadgeIds: GetTimelineTimes(r, 3),
			OwnershipTimes: GetTimelineTimes(r, 3),
		})
	}

	return randomSubassets
}

func GetRandomAddresses(r *rand.Rand, length int, accs []simtypes.Account) []string {
	randomAddresses := []string{}
	for i := 0; i < r.Intn(length) + 1; i++ {
		acc, _ := simtypes.RandomAcc(r, accs)
		randomAddresses = append(randomAddresses, acc.Address.String())
	}

	return randomAddresses
}

func GetRandomTransfers(r *rand.Rand, length int, accs []simtypes.Account) []*types.Transfer {
	randomTransfers := []*types.Transfer{}


	randomTransfers = append(randomTransfers, &types.Transfer{
		From: "Mint",
		ToAddresses: GetRandomAddresses(r, 3, accs),
		Balances: GetRandomBalances(r, 3),
	})

	for i := 0; i < r.Intn(length - 1) + 1; i++ {
		randomTransfers = append(randomTransfers, &types.Transfer{
			From: GetRandomAddresses(r, 3, accs)[0],
			ToAddresses: GetRandomAddresses(r, 3, accs),
			Balances: GetRandomBalances(r, 3),
		})
	}

	return randomTransfers
}

func GetRandomValueOptions(r *rand.Rand) *types.ValueOptions {
	val := r.Int63n(4)

	return &types.ValueOptions{
		InvertDefault: val == 0,
		AllValues: val == 1,
		NoValues: val == 2,
	}
}

func GetRandomCollectionPermissions(r *rand.Rand, accs []simtypes.Account) *types.CollectionPermissions {
	randomCollectionPermissions := &types.CollectionPermissions{
		CanDeleteCollection: []*types.ActionPermission{
			{
				DefaultValues: &types.ActionDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.ActionCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanArchive: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateContractAddress: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateOffChainBalancesMetadata: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateStandards: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateCustomData: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateManager: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateCollectionMetadata: []*types.TimedUpdatePermission{
			{
				DefaultValues: &types.TimedUpdateDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanCreateMoreBadges: []*types.BalancesActionPermission{
			{
				DefaultValues: &types.BalancesActionDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					OwnershipTimes: GetTimelineTimes(r, 3),
					BadgeIds: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.BalancesActionCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					BadgeIdsOptions: GetRandomValueOptions(r),
					OwnershipTimesOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{
			{
				DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
					BadgeIds: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
					BadgeIdsOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateInheritedBalances: []*types.TimedUpdateWithBadgeIdsPermission{
			{
				DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
					BadgeIds: GetTimelineTimes(r, 3),
				},
				Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
					BadgeIdsOptions: GetRandomValueOptions(r),
				}},
			},
		},
		CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{
			{
				DefaultValues: &types.CollectionApprovedTransferDefaultValues{
					PermittedTimes: GetTimelineTimes(r, 3),
					ForbiddenTimes: GetTimelineTimes(r, 3),
					TransferTimes: GetTimelineTimes(r, 3),
					BadgeIds: GetTimelineTimes(r, 3),
					TimelineTimes: GetTimelineTimes(r, 3),
					ToMappingId: GetRandomAddresses(r, 3, accs)[0],
					FromMappingId: GetRandomAddresses(r, 3, accs)[0],
					InitiatedByMappingId: GetRandomAddresses(r, 3, accs)[0],
				},
				Combinations: []*types.CollectionApprovedTransferCombination{{
					PermittedTimesOptions: GetRandomValueOptions(r),
					ForbiddenTimesOptions: GetRandomValueOptions(r),
					TransferTimesOptions: GetRandomValueOptions(r),
					BadgeIdsOptions: GetRandomValueOptions(r),
					TimelineTimesOptions: GetRandomValueOptions(r),
					ToMappingOptions: GetRandomValueOptions(r),
					FromMappingOptions: GetRandomValueOptions(r),
					InitiatedByMappingOptions: GetRandomValueOptions(r),
				}},
			},
		},
	}

	return randomCollectionPermissions
}