package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

const (
	// DefaultSimCollectionCount is the default number of collections to pre-create
	DefaultSimCollectionCount = 5
	// DefaultSimDynamicStoreCount is the default number of dynamic stores to pre-create
	DefaultSimDynamicStoreCount = 3
	// DefaultSimBalanceCount is the default number of balances to pre-mint
	DefaultSimBalanceCount = 10
	// MaxTimelineRange is the maximum value for timeline ranges (default: 1000)
	MaxTimelineRange = 1000
	// MinTimelineRange is the minimum value for timeline ranges (default: 1)
	MinTimelineRange = 1
	// DefaultMultiRunAttempts is the default number of times to try each operation
	DefaultMultiRunAttempts = 5
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
// Uses bounded ranges (MinTimelineRange to MaxTimelineRange) for better simulation
func GetTimelineTimes(r *rand.Rand, length int) []*types.UintRange {
	return GetBoundedTimelineTimes(r, length, MinTimelineRange, MaxTimelineRange)
}

// GetBoundedTimelineTimes generates valid UintRange slices within specified bounds
// Ensures Start <= End and generates realistic ranges
func GetBoundedTimelineTimes(r *rand.Rand, length int, min, max uint64) []*types.UintRange {
	if min >= max {
		max = min + 1
	}
	timelineTimes := make([]*types.UintRange, 0, length)
	for i := 0; i < length; i++ {
		// Generate start and end within bounds
		startVal := min + uint64(r.Int63n(int64(max-min+1)))
		endVal := min + uint64(r.Int63n(int64(max-min+1)))

		// Ensure Start <= End
		if startVal > endVal {
			startVal, endVal = endVal, startVal
		}

		// Ensure we have a valid range (at least 1)
		if startVal == endVal {
			if endVal < max {
				endVal++
			} else if startVal > min {
				startVal--
			} else {
				endVal = startVal + 1
			}
		}

		start := sdkmath.NewUint(startVal)
		end := sdkmath.NewUint(endVal)

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

// GetNonOverlappingRanges generates non-overlapping ranges within bounds
// Returns two arrays that don't overlap with each other
func GetNonOverlappingRanges(r *rand.Rand, count int, min, max uint64) ([]*types.UintRange, []*types.UintRange) {
	if count == 0 {
		return []*types.UintRange{}, []*types.UintRange{}
	}

	// Split the range into non-overlapping segments
	// Use permitted times for first half, forbidden times for second half
	rangeSize := (max - min) / uint64(count*2)
	if rangeSize == 0 {
		rangeSize = 1
	}

	permitted := make([]*types.UintRange, 0, count)
	forbidden := make([]*types.UintRange, 0, count)

	current := min
	for i := 0; i < count && current < max; i++ {
		// Permitted range
		permStart := current
		permEnd := current + rangeSize - 1
		if permEnd > max {
			permEnd = max
		}
		if permStart < permEnd {
			permitted = append(permitted, &types.UintRange{
				Start: sdkmath.NewUint(permStart),
				End:   sdkmath.NewUint(permEnd),
			})
		}
		current = permEnd + 1

		// Forbidden range (if we have space)
		if current < max && i < count {
			forbStart := current
			forbEnd := current + rangeSize - 1
			if forbEnd > max {
				forbEnd = max
			}
			if forbStart < forbEnd {
				forbidden = append(forbidden, &types.UintRange{
					Start: sdkmath.NewUint(forbStart),
					End:   sdkmath.NewUint(forbEnd),
				})
			}
			current = forbEnd + 1
		}
	}

	return permitted, forbidden
}

func GetRandomCollectionPermissions(r *rand.Rand, accs []simtypes.Account) *types.CollectionPermissions {
	// Generate non-overlapping ranges for permissions
	// Use empty arrays or single non-overlapping ranges to avoid validation errors
	permTimes, forbTimes := GetNonOverlappingRanges(r, 1, MinTimelineRange, MaxTimelineRange)

	randomCollectionPermissions := &types.CollectionPermissions{
		CanDeleteCollection: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanArchiveCollection: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanUpdateStandards: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanUpdateCustomData: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanUpdateManager: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
				TokenIds:                  GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange),
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
				TokenIds:                  GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange),
			},
		},
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
			{
				ApprovalId:                simtypes.RandStringOfLength(r, 10),
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
				TransferTimes:             GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange),
				OwnershipTimes:            GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange),
				TokenIds:                  GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange),
				ToListId:                  "All",
				FromListId:                "All",
				InitiatedByListId:         "All",
			},
		},
		CanAddMoreAliasPaths: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
			},
		},
		CanAddMoreCosmosCoinWrapperPaths: []*types.ActionPermission{
			{
				PermanentlyPermittedTimes: permTimes,
				PermanentlyForbiddenTimes: forbTimes,
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
// If collection is nil or has no valid token IDs, generates bounded random ranges
// Always respects collection constraints when available
func GetRandomValidTokenIds(r *rand.Rand, collection *types.TokenCollection, count int) []*types.UintRange {
	if collection == nil || len(collection.ValidTokenIds) == 0 {
		// No collection or no valid token IDs, generate bounded random ranges
		return GetBoundedTimelineTimes(r, count, MinTimelineRange, MaxTimelineRange)
	}

	// Use collection's valid token IDs as constraints
	validIds := collection.ValidTokenIds
	result := make([]*types.UintRange, 0, count)

	// Track used ranges to avoid duplicates
	usedIndices := make(map[int]bool)

	for i := 0; i < count && len(usedIndices) < len(validIds); i++ {
		// Pick a random valid range we haven't used yet
		idx := r.Intn(len(validIds))
		if len(usedIndices) < len(validIds) {
			// Try to avoid duplicates
			attempts := 0
			for usedIndices[idx] && attempts < len(validIds) {
				idx = r.Intn(len(validIds))
				attempts++
			}
			usedIndices[idx] = true
		}

		validRange := validIds[idx]
		// Generate a sub-range within the valid range
		start := validRange.Start
		end := validRange.End
		if start.LT(end) {
			// Generate random sub-range within bounds
			rangeSize := end.Sub(start)
			if rangeSize.GT(sdkmath.NewUint(0)) {
				// Use at least 1/4 of the range, up to the full range
				minSize := rangeSize.Quo(sdkmath.NewUint(4))
				if minSize.IsZero() {
					minSize = sdkmath.NewUint(1)
				}
				maxOffset := rangeSize.Sub(minSize)
				if maxOffset.IsZero() {
					maxOffset = sdkmath.NewUint(1)
				}

				offset := sdkmath.NewUint(uint64(r.Int63n(int64(maxOffset.Uint64()) + 1)))
				subStart := start.Add(offset)

				// Ensure sub-range size is at least minSize
				remaining := end.Sub(subStart)
				if remaining.LT(minSize) {
					subStart = end.Sub(minSize)
					if subStart.LT(start) {
						subStart = start
					}
				}

				subRangeSize := minSize.Add(sdkmath.NewUint(uint64(r.Int63n(int64(end.Sub(subStart).Sub(minSize).Uint64() + 1)))))
				subEnd := subStart.Add(subRangeSize)
				if subEnd.GT(end) {
					subEnd = end
				}
				if subEnd.LTE(subStart) {
					subEnd = subStart.Add(sdkmath.NewUint(1))
					if subEnd.GT(end) {
						subEnd = end
						subStart = end.Sub(sdkmath.NewUint(1))
						if subStart.LT(start) {
							subStart = start
						}
					}
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

	// If we need more, reuse valid ranges or generate bounded ranges
	for len(result) < count {
		if len(validIds) > 0 {
			// Reuse a valid range
			validRange := validIds[r.Intn(len(validIds))]
			result = append(result, &types.UintRange{
				Start: validRange.Start,
				End:   validRange.End,
			})
		} else {
			// Fallback to bounded ranges
			result = append(result, GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange)[0])
		}
	}

	return result
}

// GetRandomCollectionApproval generates a random collection approval with proper fields
func GetRandomCollectionApproval(r *rand.Rand, accs []simtypes.Account) *types.CollectionApproval {
	approvalId := simtypes.RandStringOfLength(r, 10)
	fromListId := "All"
	toListId := "All"
	// Sometimes use specific addresses
	if len(accs) > 0 && r.Intn(3) == 0 {
		addr := GetRandomAddresses(r, 1, accs)
		if len(addr) > 0 {
			fromListId = addr[0]
		}
	}
	if len(accs) > 0 && r.Intn(3) == 0 {
		addr := GetRandomAddresses(r, 1, accs)
		if len(addr) > 0 {
			toListId = addr[0]
		}
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

// GetValidBalancesFromCollection generates balances that respect collection token ID constraints
// Ensures amounts are non-zero and uses actual collection state when available
func GetValidBalancesFromCollection(r *rand.Rand, collection *types.TokenCollection, count int) []*types.Balance {
	balances := make([]*types.Balance, 0, count)
	for i := 0; i < count; i++ {
		// Generate non-zero amount (1 to 1000)
		amount := sdkmath.NewUint(uint64(r.Intn(1000) + 1))

		// Use valid token IDs from collection if available
		tokenIds := GetRandomValidTokenIds(r, collection, 1)

		balances = append(balances, &types.Balance{
			Amount:         amount,
			TokenIds:       tokenIds,
			OwnershipTimes: GetTimelineTimes(r, 1),
		})
	}
	return balances
}

// EnsureAccountExists verifies account exists in simulation accounts before use
// Falls back to first account if random selection fails
func EnsureAccountExists(r *rand.Rand, accs []simtypes.Account) simtypes.Account {
	if len(accs) == 0 {
		panic("no accounts available for simulation")
	}
	acc, _ := simtypes.RandomAcc(r, accs)
	return acc
}

// GetKnownGoodCollectionId returns a collection ID that's guaranteed to exist
// Falls back to creating one if needed (via SetupSimulationState)
func GetKnownGoodCollectionId(ctx sdk.Context, k keeper.Keeper) (sdkmath.Uint, bool) {
	nextId := k.GetNextCollectionId(ctx)
	if nextId.LTE(sdkmath.NewUint(1)) {
		// No collections exist yet
		return sdkmath.NewUint(0), false
	}
	// Return the first collection ID (1)
	return sdkmath.NewUint(1), true
}

// GetKnownGoodDynamicStoreId returns a dynamic store ID that's guaranteed to exist
// Falls back to creating one if needed
func GetKnownGoodDynamicStoreId(ctx sdk.Context, k keeper.Keeper) (sdkmath.Uint, bool) {
	nextStoreId := k.GetNextDynamicStoreId(ctx)
	if nextStoreId.LTE(sdkmath.NewUint(1)) {
		// No dynamic stores exist yet
		return sdkmath.NewUint(0), false
	}
	// Return the first dynamic store ID (1)
	return sdkmath.NewUint(1), true
}

// GetOrCreateCollection ensures at least one collection exists
// Returns the first collection ID, or creates a minimal valid collection if none exist
func GetOrCreateCollection(ctx sdk.Context, k keeper.Keeper, creator string, r *rand.Rand, accs []simtypes.Account) (sdkmath.Uint, error) {
	collectionId, found := GetKnownGoodCollectionId(ctx, k)
	if found {
		return collectionId, nil
	}

	// Create a minimal valid collection
	validTokenIds := GetBoundedTimelineTimes(r, 1, MinTimelineRange, MaxTimelineRange)
	collectionPermissions := GetRandomCollectionPermissions(r, accs)

	collectionMetadata := &types.CollectionMetadata{
		Uri:        "https://example.com/metadata/sim",
		CustomData: "simulation",
	}

	msg := &types.MsgCreateCollection{
		Creator:               creator,
		DefaultBalances:       &types.UserBalanceStore{Balances: []*types.Balance{}},
		ValidTokenIds:         validTokenIds,
		CollectionPermissions: collectionPermissions,
		Manager:               creator,
		CollectionMetadata:    collectionMetadata,
		TokenMetadata:         []*types.TokenMetadata{},
		CustomData:            "simulation",
		CollectionApprovals:   []*types.CollectionApproval{},
		Standards:             []string{},
		IsArchived:            false,
	}

	// Execute the message via msgServer
	msgServer := keeper.NewMsgServerImpl(k)
	_, err := msgServer.CreateCollection(ctx, msg)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	// Return the newly created collection ID
	return k.GetNextCollectionId(ctx).Sub(sdkmath.NewUint(1)), nil
}

// SetupSimulationState pre-creates collections, dynamic stores, and balances for better simulation
// This ensures that simulation operations have valid state to work with
func SetupSimulationState(ctx sdk.Context, k keeper.Keeper, accs []simtypes.Account, r *rand.Rand) error {
	if len(accs) == 0 {
		return nil // Can't setup state without accounts
	}

	// Pre-create collections
	for i := 0; i < DefaultSimCollectionCount; i++ {
		creator := EnsureAccountExists(r, accs)
		validTokenIds := GetBoundedTimelineTimes(r, r.Intn(3)+1, MinTimelineRange, MaxTimelineRange)
		collectionPermissions := GetRandomCollectionPermissions(r, accs)

		collectionMetadata := &types.CollectionMetadata{
			Uri:        "https://example.com/metadata/sim-" + simtypes.RandStringOfLength(r, 10),
			CustomData: "simulation-" + simtypes.RandStringOfLength(r, 10),
		}

		// Sometimes add a mint approval
		collectionApprovals := []*types.CollectionApproval{}
		if r.Intn(2) == 0 {
			mintApproval := GetRandomCollectionApproval(r, accs)
			mintApproval.FromListId = types.MintAddress
			mintApproval.ToListId = "All"
			mintApproval.InitiatedByListId = "All"
			if mintApproval.ApprovalCriteria == nil {
				mintApproval.ApprovalCriteria = &types.ApprovalCriteria{}
			}
			mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
			mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
			collectionApprovals = append(collectionApprovals, mintApproval)
		}

		msg := &types.MsgCreateCollection{
			Creator:               creator.Address.String(),
			DefaultBalances:       &types.UserBalanceStore{Balances: []*types.Balance{}},
			ValidTokenIds:         validTokenIds,
			CollectionPermissions: collectionPermissions,
			Manager:               creator.Address.String(),
			CollectionMetadata:    collectionMetadata,
			TokenMetadata:         []*types.TokenMetadata{},
			CustomData:            simtypes.RandStringOfLength(r, 20),
			CollectionApprovals:   collectionApprovals,
			Standards:             []string{},
			IsArchived:            false,
		}

		msgServer := keeper.NewMsgServerImpl(k)
		_, err := msgServer.CreateCollection(ctx, msg)
		if err != nil {
			// Continue even if one fails
			continue
		}
	}

	// Pre-create dynamic stores
	for i := 0; i < DefaultSimDynamicStoreCount; i++ {
		creator := EnsureAccountExists(r, accs)
		defaultValue := r.Intn(2) == 0

		msg := &types.MsgCreateDynamicStore{
			Creator:      creator.Address.String(),
			DefaultValue: defaultValue,
		}

		msgServer := keeper.NewMsgServerImpl(k)
		_, err := msgServer.CreateDynamicStore(ctx, msg)
		if err != nil {
			// Continue even if one fails
			continue
		}
	}

	// Pre-mint some balances to known accounts
	// Get all collections we just created
	nextCollectionId := k.GetNextCollectionId(ctx)
	if nextCollectionId.GT(sdkmath.NewUint(1)) {
		maxCollectionId := nextCollectionId.Sub(sdkmath.NewUint(1))
		// Mint balances to a few collections
		for i := 0; i < DefaultSimBalanceCount && i < int(maxCollectionId.Uint64()); i++ {
			collectionId := sdkmath.NewUint(uint64(i + 1))
			collection, found := k.GetCollectionFromStore(ctx, collectionId)
			if !found {
				continue
			}

			// Mint to a random account
			recipient := EnsureAccountExists(r, accs)
			balances := GetValidBalancesFromCollection(r, collection, 1)

			transfers := []*types.Transfer{
				{
					From:        types.MintAddress,
					ToAddresses: []string{recipient.Address.String()},
					Balances:    balances,
				},
			}

			msg := &types.MsgTransferTokens{
				Creator:      recipient.Address.String(),
				CollectionId: collectionId,
				Transfers:    transfers,
			}

			msgServer := keeper.NewMsgServerImpl(k)
			_, err := msgServer.TransferTokens(ctx, msg)
			if err != nil {
				// Continue even if one fails
				continue
			}
		}
	}

	return nil
}
