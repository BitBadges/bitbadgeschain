package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
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

// GetTimelineTimes generates valid UintRange slices where Start <= End
func GetTimelineTimes(r *rand.Rand, length int) []*types.UintRange {
	timelineTimes := make([]*types.UintRange, 0, length)
	for i := 0; i < length; i++ {
		start := sdkmath.NewUint(r.Uint64())
		end := sdkmath.NewUint(r.Uint64())
		// Ensure Start <= End
		if start.GT(end) {
			start, end = end, start
		}
		// Ensure End is at least Start (handle case where both are 0)
		if end.IsZero() && start.IsZero() {
			end = sdkmath.NewUint(1)
		}
		timelineTimes = append(timelineTimes, &types.UintRange{
			Start: start,
			End:   end,
		})
	}
	return timelineTimes
}

// GetRandomBalances generates random balances (may include zero amounts)
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

// GetRandomValidBalances generates random balances with amounts > 0
func GetRandomValidBalances(r *rand.Rand, length int) []*types.Balance {
	randomBalances := []*types.Balance{}
	count := r.Intn(length) + 1
	for i := 0; i < count; i++ {
		amount := sdkmath.NewUint(r.Uint64())
		// Ensure amount > 0
		if amount.IsZero() {
			amount = sdkmath.NewUint(1)
		}
		randomBalances = append(randomBalances, &types.Balance{
			Amount:         amount,
			TokenIds:       GetTimelineTimes(r, 1),
			OwnershipTimes: GetTimelineTimes(r, 1),
		})
	}
	return randomBalances
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

// GetRandomCollectionId returns a random existing collection ID, or zero if none exist
func GetRandomCollectionId(r *rand.Rand, ctx sdk.Context, k keeper.Keeper) (sdkmath.Uint, bool) {
	nextId := k.GetNextCollectionId(ctx)
	if nextId.LTE(sdkmath.NewUint(1)) {
		// No collections exist yet (nextId is 1 means no collections created)
		return sdkmath.NewUint(0), false
	}
	// Collections exist from 1 to (nextId - 1)
	maxId := nextId.Sub(sdkmath.NewUint(1))
	if maxId.IsZero() {
		return sdkmath.NewUint(0), false
	}
	// Return random ID between 1 and maxId
	randomId := sdkmath.NewUint(uint64(r.Int63n(int64(maxId.Uint64()))) + 1)
	return randomId, true
}

// GetRandomValidTokenIds generates token IDs that are valid for a collection
// If collection is nil or has no valid token IDs, generates random valid ranges
func GetRandomValidTokenIds(r *rand.Rand, collection *types.TokenCollection, count int) []*types.UintRange {
	if collection == nil || len(collection.ValidTokenIds) == 0 {
		// No collection or no valid token IDs, generate random ranges
		return GetTimelineTimes(r, count)
	}
	
	// Use collection's valid token IDs as constraints
	validIds := collection.ValidTokenIds
	result := make([]*types.UintRange, 0, count)
	
	for i := 0; i < count && i < len(validIds); i++ {
		validRange := validIds[r.Intn(len(validIds))]
		// Generate a sub-range within the valid range
		start := validRange.Start
		end := validRange.End
		if start.LT(end) {
			// Generate random sub-range
			rangeSize := end.Sub(start)
			if rangeSize.GT(sdkmath.NewUint(0)) {
				offset := sdkmath.NewUint(uint64(r.Int63n(int64(rangeSize.Uint64()))))
				subStart := start.Add(offset)
				subEnd := subStart.Add(sdkmath.NewUint(uint64(r.Int63n(int64(end.Sub(subStart).Uint64()) + 1))))
				if subEnd.GT(end) {
					subEnd = end
				}
				result = append(result, &types.UintRange{
					Start: subStart,
					End:   subEnd,
				})
			} else {
				result = append(result, validRange)
			}
		} else {
			result = append(result, validRange)
		}
	}
	
	// If we need more, fill with random ranges
	for len(result) < count {
		result = append(result, GetTimelineTimes(r, 1)[0])
	}
	
	return result
}

// GetRandomCollectionApproval generates a random collection approval with proper fields
func GetRandomCollectionApproval(r *rand.Rand, accs []simtypes.Account) *types.CollectionApproval {
	approvalId := simtypes.RandStringOfLength(r, 10)
	fromListId := "All"
	toListId := "All"
	// Sometimes use specific addresses
	if r.Intn(3) == 0 {
		fromListId = GetRandomAddresses(r, 1, accs)[0]
	}
	if r.Intn(3) == 0 {
		toListId = GetRandomAddresses(r, 1, accs)[0]
	}
	
	return &types.CollectionApproval{
		ApprovalId:        approvalId,
		FromListId:        fromListId,
		ToListId:          toListId,
		InitiatedByListId: "All", // Default to "All" for initiated by list
		TransferTimes:     GetTimelineTimes(r, 1),
		TokenIds:          GetTimelineTimes(r, 1),
		OwnershipTimes:    GetTimelineTimes(r, 1),
		ApprovalCriteria:  &types.ApprovalCriteria{},
		Version:           sdkmath.NewUint(0),
	}
}
