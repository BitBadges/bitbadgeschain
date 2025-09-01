package poolmanager

import (
	"sort"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

//
// Taker Fee Share Agreements
//

// getAllTakerFeeShareAgreementsMap creates the map used for the taker fee share agreements cache.
func (k Keeper) getAllTakerFeeShareAgreementsMap(ctx sdk.Context) (map[string]types.TakerFeeShareAgreement, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyTakerFeeShare)
	defer iterator.Close()

	takerFeeShareAgreementsMap := make(map[string]types.TakerFeeShareAgreement)
	for ; iterator.Valid(); iterator.Next() {
		takerFeeShareAgreement := types.TakerFeeShareAgreement{}
		if err := proto.Unmarshal(iterator.Value(), &takerFeeShareAgreement); err != nil {
			return nil, err
		}
		takerFeeShareAgreementsMap[takerFeeShareAgreement.Denom] = takerFeeShareAgreement
	}

	return takerFeeShareAgreementsMap, nil
}

// GetAllTakerFeesShareAgreements creates a slice of all taker fee share agreements.
// Used in the AllTakerFeeShareAgreementsRequest gRPC query.
func (k Keeper) GetAllTakerFeesShareAgreements(ctx sdk.Context) ([]types.TakerFeeShareAgreement, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyTakerFeeShare)
	defer iterator.Close()

	takerFeeShareAgreements := []types.TakerFeeShareAgreement{}
	for ; iterator.Valid(); iterator.Next() {
		takerFeeShareAgreement := types.TakerFeeShareAgreement{}
		if err := proto.Unmarshal(iterator.Value(), &takerFeeShareAgreement); err != nil {
			return nil, err
		}
		takerFeeShareAgreements = append(takerFeeShareAgreements, takerFeeShareAgreement)
	}

	return takerFeeShareAgreements, nil
}

// setTakerFeeShareAgreementsMapCached is used for initializing the cache for the taker fee share agreements.
func (k *Keeper) setTakerFeeShareAgreementsMapCached(ctx sdk.Context) error {
	takerFeeShareAgreement, err := k.getAllTakerFeeShareAgreementsMap(ctx)
	if err != nil {
		return err
	}
	k.cachedTakerFeeShareAgreementMap = takerFeeShareAgreement
	return nil
}

// getTakerFeeShareAgreementFromDenom retrieves a specific taker fee share agreement from the store.
func (k Keeper) getTakerFeeShareAgreementFromDenom(takerFeeShareDenom string) (types.TakerFeeShareAgreement, bool) {
	takerFeeShareAgreement, found := k.cachedTakerFeeShareAgreementMap[takerFeeShareDenom]
	return takerFeeShareAgreement, found
}

// GetTakerFeeShareAgreementFromDenomUNSAFE is used to expose an internal method to gRPC query. This method should not be used in other modules, since the cache is not populated in those keepers.
// Used in the TakerFeeShareAgreementFromDenomRequest gRPC query.
func (k Keeper) GetTakerFeeShareAgreementFromDenomUNSAFE(takerFeeShareDenom string) (types.TakerFeeShareAgreement, bool) {
	return k.getTakerFeeShareAgreementFromDenom(takerFeeShareDenom)
}

// GetTakerFeeShareAgreementFromDenom retrieves a specific taker fee share agreement from the store, bypassing cache.
func (k Keeper) GetTakerFeeShareAgreementFromDenomNoCache(ctx sdk.Context, takerFeeShareDenom string) (types.TakerFeeShareAgreement, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.FormatTakerFeeShareAgreementKey(takerFeeShareDenom)
	bz := store.Get(key)
	if bz == nil {
		return types.TakerFeeShareAgreement{}, false
	}

	var takerFeeShareAgreement types.TakerFeeShareAgreement
	if err := proto.Unmarshal(bz, &takerFeeShareAgreement); err != nil {
		return types.TakerFeeShareAgreement{}, false
	}

	return takerFeeShareAgreement, true
}

// SetTakerFeeShareAgreementForDenom is used for setting a specific taker fee share agreement in the store.
// Used in the MsgSetTakerFeeShareAgreementForDenom, by the gov address only.
func (k *Keeper) SetTakerFeeShareAgreementForDenom(ctx sdk.Context, takerFeeShare types.TakerFeeShareAgreement) error {
	store := ctx.KVStore(k.storeKey)
	key := types.FormatTakerFeeShareAgreementKey(takerFeeShare.Denom)
	bz, err := proto.Marshal(&takerFeeShare)
	if err != nil {
		return err
	}

	store.Set(key, bz)

	// Set cache value
	k.cachedTakerFeeShareAgreementMap[takerFeeShare.Denom] = takerFeeShare

	return nil
}

//
// Taker Fee Share Accumulators
//

// GetTakerFeeShareDenomsToAccruedValue retrieves the accrued value for a specific taker fee share denomination and taker fee charged denomination from the store.
// Used in the TakerFeeShareDenomsToAccruedValueRequest gRPC query.
func (k Keeper) GetTakerFeeShareDenomsToAccruedValue(ctx sdk.Context, takerFeeShareDenom string, takerFeeChargedDenom string) (osmomath.Int, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyTakerFeeShareDenomAccrualForTakerFeeChargedDenom(takerFeeShareDenom, takerFeeChargedDenom)
	accruedValue := sdk.IntProto{}
	found, err := osmoutils.Get(store, key, &accruedValue)
	if err != nil {
		return osmomath.Int{}, err
	}
	if !found {
		return osmomath.Int{}, types.NoAccruedValueError{TakerFeeShareDenom: takerFeeShareDenom, TakerFeeChargedDenom: takerFeeChargedDenom}
	}
	return accruedValue.Int, nil
}

// SetTakerFeeShareDenomsToAccruedValue sets the accrued value for a specific taker fee share denomination and taker fee charged denomination in the store.
func (k Keeper) SetTakerFeeShareDenomsToAccruedValue(ctx sdk.Context, takerFeeShareDenom string, takerFeeChargedDenom string, accruedValue osmomath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyTakerFeeShareDenomAccrualForTakerFeeChargedDenom(takerFeeShareDenom, takerFeeChargedDenom)
	accruedValueProto := sdk.IntProto{Int: accruedValue}
	bz, err := proto.Marshal(&accruedValueProto)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// increaseTakerFeeShareDenomsToAccruedValue increases (adds to, not replace) the accrued value for a specific taker fee share denomination and taker fee charged denomination in the store.
func (k Keeper) increaseTakerFeeShareDenomsToAccruedValue(ctx sdk.Context, takerFeeShareDenom string, takerFeeChargedDenom string, additiveValue osmomath.Int) error {
	accruedValueBefore, err := k.GetTakerFeeShareDenomsToAccruedValue(ctx, takerFeeShareDenom, takerFeeChargedDenom)
	if err != nil {
		if _, ok := err.(types.NoAccruedValueError); ok {
			accruedValueBefore = osmomath.ZeroInt()
		} else {
			return err
		}
	}

	accruedValueAfter := accruedValueBefore.Add(additiveValue)
	return k.SetTakerFeeShareDenomsToAccruedValue(ctx, takerFeeShareDenom, takerFeeChargedDenom, accruedValueAfter)
}

// GetAllTakerFeeShareAccumulators creates a slice of all taker fee share accumulators.
// Used in the AllTakerFeeShareAccumulatorsRequest gRPC query.
func (k Keeper) GetAllTakerFeeShareAccumulators(ctx sdk.Context) ([]types.TakerFeeSkimAccumulator, error) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.TakerFeeSkimAccrualPrefix)
	defer iter.Close()

	takerFeeAgreementDenomToCoins := make(map[string]sdk.Coins)
	var denoms []string // Slice to keep track of the keys and ensure deterministic ordering

	for ; iter.Valid(); iter.Next() {
		accruedValue := sdk.IntProto{}
		if err := proto.Unmarshal(iter.Value(), &accruedValue); err != nil {
			return nil, err
		}
		keyParts := strings.Split(string(iter.Key()), types.KeySeparator)
		tierDenom := keyParts[1]
		takerFeeDenom := keyParts[2]
		accruedValueInt := accruedValue.Int
		currentCoins := takerFeeAgreementDenomToCoins[tierDenom]

		// Add the denom to the slice if it's not already there
		if _, exists := takerFeeAgreementDenomToCoins[tierDenom]; !exists {
			denoms = append(denoms, tierDenom)
		}

		takerFeeAgreementDenomToCoins[tierDenom] = currentCoins.Add(sdk.NewCoin(takerFeeDenom, accruedValueInt))
	}

	takerFeeSkimAccumulators := []types.TakerFeeSkimAccumulator{}
	for _, denom := range denoms {
		takerFeeSkimAccumulators = append(takerFeeSkimAccumulators, types.TakerFeeSkimAccumulator{
			Denom:            denom,
			SkimmedTakerFees: takerFeeAgreementDenomToCoins[denom],
		})
	}

	return takerFeeSkimAccumulators, nil
}

// DeleteAllTakerFeeShareAccumulatorsForTakerFeeShareDenom clears the TakerFeeShareAccumulator records for a specific taker fee share denom.
// Is specifically used after the distributions have been completed after epoch for each denom.
func (k Keeper) DeleteAllTakerFeeShareAccumulatorsForTakerFeeShareDenom(ctx sdk.Context, takerFeeShareDenom string) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyTakerFeeShareDenomAccrualForAllDenoms(takerFeeShareDenom))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// getRegisteredAlloyedPoolFromDenom retrieves a specific registered alloyed pool from the store via the alloyed denom.
func (k Keeper) getRegisteredAlloyedPoolFromDenom(alloyedDenom string) (types.AlloyContractTakerFeeShareState, bool) {
	registeredAlloyedPool, found := k.cachedRegisteredAlloyPoolByAlloyDenomMap[alloyedDenom]
	if !found {
		return types.AlloyContractTakerFeeShareState{}, false
	}
	return registeredAlloyedPool, true
}

// GetRegisteredAlloyedPoolFromDenomUNSAFE is used to expose an internal method to gRPC query. This method should not be used in other modules, since the cache is not populated in those keepers.
// Used in the RegisteredAlloyedPoolFromDenomRequest gRPC query.
func (k Keeper) GetRegisteredAlloyedPoolFromDenomUNSAFE(alloyedDenom string) (types.AlloyContractTakerFeeShareState, bool) {
	return k.getRegisteredAlloyedPoolFromDenom(alloyedDenom)
}

// getRegisteredAlloyedPoolFromPoolId retrieves a specific registered alloyed pool from the store via the pool id.
func (k Keeper) getRegisteredAlloyedPoolFromPoolId(ctx sdk.Context, poolId uint64) (types.AlloyContractTakerFeeShareState, error) {
	alloyedDenom, err := k.getAlloyedDenomFromPoolId(ctx, poolId)
	if err != nil {
		return types.AlloyContractTakerFeeShareState{}, err
	}
	registeredAlloyedPool, found := k.getRegisteredAlloyedPoolFromDenom(alloyedDenom)
	if !found {
		return types.AlloyContractTakerFeeShareState{}, types.NoRegisteredAlloyedPoolError{PoolId: poolId}
	}
	return registeredAlloyedPool, nil
}

// GetRegisteredAlloyedPoolFromPoolIdUNSAFE is used to expose an internal method to gRPC query. This method should not be used in other modules, since the cache is not populated in those keepers.
// Used in the RegisteredAlloyedPoolFromPoolIdRequest gRPC query.
func (k Keeper) GetRegisteredAlloyedPoolFromPoolIdUNSAFE(ctx sdk.Context, poolId uint64) (types.AlloyContractTakerFeeShareState, error) {
	return k.getRegisteredAlloyedPoolFromPoolId(ctx, poolId)
}

// GetAllRegisteredAlloyedPools creates a slice of all registered alloyed pools.
// Used in the AllRegisteredAlloyedPoolsRequest gRPC query.
func (k Keeper) GetAllRegisteredAlloyedPools(ctx sdk.Context) ([]types.AlloyContractTakerFeeShareState, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyRegisteredAlloyPool)
	defer iterator.Close()

	var registeredAlloyedPools []types.AlloyContractTakerFeeShareState
	for ; iterator.Valid(); iterator.Next() {
		registeredAlloyedPool := types.AlloyContractTakerFeeShareState{}
		err := proto.Unmarshal(iterator.Value(), &registeredAlloyedPool)
		if err != nil {
			return nil, err
		}

		registeredAlloyedPools = append(registeredAlloyedPools, registeredAlloyedPool)
	}

	return registeredAlloyedPools, nil
}

// GetAllRegisteredAlloyedPoolsByDenomMap creates the map used for the registered alloyed pools cache.
func (k Keeper) getAllRegisteredAlloyedPoolsByDenomMap(ctx sdk.Context) (map[string]types.AlloyContractTakerFeeShareState, error) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.KeyRegisteredAlloyPool)
	defer iter.Close()

	registeredAlloyedPoolsMap := make(map[string]types.AlloyContractTakerFeeShareState)
	for ; iter.Valid(); iter.Next() {
		registeredAlloyedPool := types.AlloyContractTakerFeeShareState{}
		if err := proto.Unmarshal(iter.Value(), &registeredAlloyedPool); err != nil {
			return nil, err
		}

		key := string(iter.Key())
		parts := strings.Split(key, types.KeySeparator)
		if len(parts) < 3 {
			return nil, types.ErrInvalidKeyFormat
		}
		alloyedDenom := parts[len(parts)-1]
		registeredAlloyedPoolsMap[alloyedDenom] = registeredAlloyedPool
	}

	return registeredAlloyedPoolsMap, nil
}

// setAllRegisteredAlloyedPoolsByDenomCached initializes the cache for the registered alloyed pools.
func (k *Keeper) setAllRegisteredAlloyedPoolsByDenomCached(ctx sdk.Context) error {
	registeredAlloyPools, err := k.getAllRegisteredAlloyedPoolsByDenomMap(ctx)
	if err != nil {
		return err
	}
	k.cachedRegisteredAlloyPoolByAlloyDenomMap = registeredAlloyPools
	return nil
}

//
// Registered Alloyed Pool Ids
//

// getAllRegisteredAlloyedPoolsIdArray creates an array of all registered alloyed pools IDs.
func (k Keeper) getAllRegisteredAlloyedPoolsIdArray(ctx sdk.Context) ([]uint64, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyRegisteredAlloyPool)
	defer iterator.Close()

	registeredAlloyedPoolsIdArray := []uint64{}
	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		parts := strings.Split(key, types.KeySeparator)
		if len(parts) < 3 {
			return nil, types.ErrInvalidKeyFormat
		}
		alloyedIdStr := parts[len(parts)-2]
		// Convert the string to uint64
		alloyedId, err := strconv.ParseUint(alloyedIdStr, 10, 64)
		if err != nil {
			return nil, types.InvalidAlloyedPoolIDError{AlloyedIDStr: alloyedIdStr, Err: err}
		}
		registeredAlloyedPoolsIdArray = append(registeredAlloyedPoolsIdArray, alloyedId)
	}
	sort.Slice(registeredAlloyedPoolsIdArray, func(i, j int) bool { return registeredAlloyedPoolsIdArray[i] < registeredAlloyedPoolsIdArray[j] })

	return registeredAlloyedPoolsIdArray, nil
}

// calculateTakerFeeShareAgreements calculates the taker fee share agreements based on the total pool liquidity
// and normalization factors. It iterates through the pool liquidity, normalizes the amounts, and calculates
// the scaled skim percentages for each asset with a share agreement. Returns a slice of TakerFeeShareAgreement
// objects if successful, otherwise returns an error.
func (k Keeper) calculateTakerFeeShareAgreements(totalPoolLiquidity []sdk.Coin, normalizationFactors map[string]osmomath.Dec) ([]types.TakerFeeShareAgreement, error) {
	totalAlloyedLiquidity := types.ZeroDec
	var assetsWithShareAgreement []sdk.Coin
	var takerFeeShareAgreements []types.TakerFeeShareAgreement
	var skimAddresses []string
	var skimPercents []osmomath.Dec

	for _, coin := range totalPoolLiquidity {
		normalizationFactor := normalizationFactors[coin.Denom]
		normalizedAmount := coin.Amount.ToLegacyDec().Quo(normalizationFactor)
		totalAlloyedLiquidity = totalAlloyedLiquidity.Add(normalizedAmount)

		takerFeeShareAgreement, found := k.getTakerFeeShareAgreementFromDenom(coin.Denom)
		if !found {
			continue
		}
		assetsWithShareAgreement = append(assetsWithShareAgreement, coin)
		skimAddresses = append(skimAddresses, takerFeeShareAgreement.SkimAddress)
		skimPercents = append(skimPercents, takerFeeShareAgreement.SkimPercent)
	}

	if totalAlloyedLiquidity.IsZero() {
		return nil, types.ErrTotalAlloyedLiquidityIsZero
	}

	for i, coin := range assetsWithShareAgreement {
		normalizationFactor := normalizationFactors[coin.Denom]
		normalizedAmount := coin.Amount.ToLegacyDec().Quo(normalizationFactor)
		scaledSkim := normalizedAmount.Quo(totalAlloyedLiquidity).Mul(skimPercents[i])
		takerFeeShareAgreements = append(takerFeeShareAgreements, types.TakerFeeShareAgreement{
			Denom:       coin.Denom,
			SkimPercent: scaledSkim,
			SkimAddress: skimAddresses[i],
		})
	}

	return takerFeeShareAgreements, nil
}

// getAlloyedDenomFromPoolId retrieves the alloyed denomination associated with a given pool ID from the store.
// It iterates through the registered alloyed pools and matches the pool ID to find the corresponding alloyed denomination.
// Returns the alloyed denomination if found, otherwise returns an error.
func (k Keeper) getAlloyedDenomFromPoolId(ctx sdk.Context, poolId uint64) (string, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyRegisteredAlloyPool)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		parts := strings.Split(key, types.KeySeparator)
		if len(parts) < 3 {
			return "", types.ErrInvalidKeyFormat
		}
		alloyedIdStr := parts[len(parts)-2]
		// Convert the string to uint64
		alloyedId, err := strconv.ParseUint(alloyedIdStr, 10, 64)
		if err != nil {
			return "", types.InvalidAlloyedPoolIDError{AlloyedIDStr: alloyedIdStr, Err: err}
		}
		if alloyedId == poolId {
			alloyedDenom := parts[len(parts)-1]
			return alloyedDenom, nil
		}
	}
	return "", types.NoRegisteredAlloyedPoolError{PoolId: poolId}
}
