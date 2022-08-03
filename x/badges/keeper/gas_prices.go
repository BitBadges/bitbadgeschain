package keeper

//These can always be optimized and we are open to suggestions for them
//We also want to eventually include a refund mechanism for self destructing badges
const (
	FixedCostPerMsg = 100
	
	BadgeCost = 1000 
	SubbadgeWithSupplyNotEqualToOne = 10 

	FreezeOrUnfreezeAddress = 50 //This is probably the most expensive operation that can involve the most storage
	SimpleAdjustBalanceOrApproval = 10
	AddOrRemovePending = 10
	RequestTransferManagerCost = 10
	TransferManagerCost = 10
	BadgeUpdate = 10
)