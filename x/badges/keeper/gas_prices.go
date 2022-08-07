package keeper

//These can always be optimized and we are open to suggestions for them
//We also want to eventually include a refund mechanism for self destructing badges

//We can experiment with this. It may be fine just to leave gas as is w/ default Cosmos settings

// CosmosSDK KVGasConfig returns a default gas config for KVStores.
// func KVGasConfig() GasConfig {
// 	return GasConfig{
// 		HasCost:          1000,
// 		DeleteCost:       1000,
// 		ReadCostFlat:     1000,
// 		ReadCostPerByte:  3,
// 		WriteCostFlat:    2000,
// 		WriteCostPerByte: 30,
// 		IterNextCostFlat: 30,
// 	}
// }

const (
	FixedCostPerMsg = 0

	BadgeCost                       = 0
	SubbadgeWithSupplyNotEqualToOne = 0

	FreezeOrUnfreezeAddress       = 0 //This is probably the most expensive operation that can involve the most storage
	SimpleAdjustBalanceOrApproval = 0
	AddOrRemovePending            = 0
	RequestTransferManagerCost    = 0
	TransferManagerCost           = 0
	BadgeUpdate                   = 0

	PruneBalanceRefundAmountPerAddress 	  = 750
	PruneBalanceRefundAmountPerBadge 	  = 750
)
