package keeper

import (
	"context"
	"fmt"
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//TODO: Clean and DRY

func GetCurrentUserApprovedTransfers(ctx sdk.Context, userBalance *types.UserBalanceStore) []*types.UserApprovedTransfer {
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	approvedTransfersTimeline := userBalance.ApprovedTransfersTimeline
	for _, approvedTransfersTimelineVal := range approvedTransfersTimeline {
		_, found := types.SearchIdRangesForId(blockTime, approvedTransfersTimelineVal.Times)
		if found {
			return approvedTransfersTimelineVal.ApprovedTransfers
		}
	}

	return []*types.UserApprovedTransfer{}
}

func GetCurrentManager(ctx sdk.Context, collection types.BadgeCollection) string {
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	managerTimeline := collection.ManagerTimeline
	for _, managerTimelineVal := range managerTimeline {
		_, found := types.SearchIdRangesForId(blockTime, managerTimelineVal.Times)
		if found {
			return managerTimelineVal.Manager
		}
	}

	return ""
}

func GetIsArchived(ctx sdk.Context, collection types.BadgeCollection) bool {
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	isArchivedTimeline := collection.IsArchivedTimeline
	for _, isArchivedTimelineVal := range isArchivedTimeline {
		_, found := types.SearchIdRangesForId(blockTime, isArchivedTimelineVal.Times)
		if found {
			return isArchivedTimelineVal.IsArchived
		}
	}

	return false
}

func GetIsArchivedTimesAndValues(isArchivedTimeline []*types.IsArchivedTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range isArchivedTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.IsArchived)
	}
	return times, values
}

func GetCollectionApprovedTransferTimesAndValues(approvedTransfers []*types.CollectionApprovedTransferTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range approvedTransfers {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.ApprovedTransfers)
	}
	return times, values
}

func GetUserApprovedTransferTimesAndValues(approvedTransfers []*types.UserApprovedTransferTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range approvedTransfers {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.ApprovedTransfers)
	}
	return times, values
}


func GetInheritedBalancesTimesAndValues(inheritedBalances []*types.InheritedBalancesTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range inheritedBalances {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.InheritedBalances)
	}

	return times, values
}

func GetOffChainBalancesMetadataTimesAndValues(inheritedBalancesMetadata []*types.OffChainBalancesMetadataTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range inheritedBalancesMetadata {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.OffChainBalancesMetadata)
	}
	return times, values
}

func GetCollectionMetadataTimesAndValues(timeline []*types.CollectionMetadataTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.CollectionMetadata)
	}
	return times, values
}

func GetBadgeMetadataTimesAndValues(timeline []*types.BadgeMetadataTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.BadgeMetadata)
	}
	return times, values
}

func GetManagerTimesAndValues(managerTimeline []*types.ManagerTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range managerTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.Manager)
	}
	return times, values
}

func GetContractAddressTimesAndValues(contractAddressTimeline []*types.ContractAddressTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range contractAddressTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.ContractAddress)
	}
	return times, values
}

func GetCustomDataTimesAndValues(customDataTimeline []*types.CustomDataTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range customDataTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.CustomData)
	}
	return times, values
}

func GetStandardsTimesAndValues(standardsTimeline []*types.StandardTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range standardsTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.Standards)
	}
	return times, values
}

