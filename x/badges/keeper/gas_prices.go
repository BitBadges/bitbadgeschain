package keeper

//For the most part, we want to just use the default gas costs of Cosmos SDK keeper implementations.
//However, in some situations, we may want additional gas incentives or penalties.

const (
	FixedCostPerMsg                 = 0
	BadgeCost                       = 0
	SubbadgeWithSupplyNotEqualToOne = 0
	FreezeOrUnfreezeAddress         = 0
	SimpleAdjustBalanceOrApproval   = 0
	AddOrRemovePending              = 0
	RequestTransferManagerCost      = 0
	TransferManagerCost             = 0
	BadgeUpdate                     = 0

	// Incentivize users to prune expired / self-destructed badges.
	PruneBalanceRefundAmountPerAddress = 750
	PruneBalanceRefundAmountPerBadge   = 750
)
