package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
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
			End:   sdkmath.NewUint(r.Uint64()),
		})
	}
	return timelineTimes
}

func GetRandomBalances(r *rand.Rand, length int) []*types.Balance {
	randomSubassets := []*types.Balance{}
	for i := 0; i < r.Intn(length); i++ {
		randomSubassets = append(randomSubassets, &types.Balance{
			Amount:         sdkmath.NewUint(r.Uint64()),
			TokenIds:       GetTimelineTimes(r, 3),
			OwnershipTimes: GetTimelineTimes(r, 3),
		})
	}

	return randomSubassets
}

func GetRandomAddresses(r *rand.Rand, length int, accs []simtypes.Account) []string {
	randomAddresses := []string{}
	for i := 0; i < r.Intn(length)+1; i++ {
		acc, _ := simtypes.RandomAcc(r, accs)
		randomAddresses = append(randomAddresses, acc.Address.String())
	}

	return randomAddresses
}

func GetRandomTransfers(r *rand.Rand, length int, accs []simtypes.Account) []*types.Transfer {
	randomTransfers := []*types.Transfer{}

	randomTransfers = append(randomTransfers, &types.Transfer{
		From:        "Mint",
		ToAddresses: GetRandomAddresses(r, 3, accs),
		Balances:    GetRandomBalances(r, 3),
	})

	for i := 0; i < r.Intn(length-1)+1; i++ {
		randomTransfers = append(randomTransfers, &types.Transfer{
			From:        GetRandomAddresses(r, 3, accs)[0],
			ToAddresses: GetRandomAddresses(r, 3, accs),
			Balances:    GetRandomBalances(r, 3),
		})
	}

	return randomTransfers
}

func GetRandomCollectionPermissions(r *rand.Rand, accs []simtypes.Account) *types.CollectionPermissions {
	randomCollectionPermissions := &types.CollectionPermissions{
		CanDeleteCollection: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanArchiveCollection: []*types.ActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanUpdateStandards: []*types.ActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanUpdateCustomData: []*types.ActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanUpdateManager: []*types.ActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
			},
		},
		CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
				TokenIds:                  GetTimelineTimes(r, 3),
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
				TokenIds:                  GetTimelineTimes(r, 3),
			},
		},
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
			{

				PermanentlyPermittedTimes: GetTimelineTimes(r, 3),
				PermanentlyForbiddenTimes: GetTimelineTimes(r, 3),
				TransferTimes:             GetTimelineTimes(r, 3),
				OwnershipTimes:            GetTimelineTimes(r, 3),
				TokenIds:                  GetTimelineTimes(r, 3),
				ToListId:                  GetRandomAddresses(r, 3, accs)[0],
				FromListId:                GetRandomAddresses(r, 3, accs)[0],
				InitiatedByListId:         GetRandomAddresses(r, 3, accs)[0],
			},
		},
	}

	return randomCollectionPermissions
}
