package types

func CastOutgoingTransfersToCollectionTransfers(transfers []*UserApprovedOutgoingTransfer, fromAddress string) []*CollectionApprovedTransfer {
	collectionTransfers := []*CollectionApprovedTransfer{}
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastOutgoingTransferToCollectionTransfer(transfer, fromAddress))
	}
	
	return collectionTransfers
}

func CastIncomingTransfersToCollectionTransfers(transfers []*UserApprovedIncomingTransfer, toAddress string) []*CollectionApprovedTransfer {
	collectionTransfers := []*CollectionApprovedTransfer{}
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastIncomingTransferToCollectionTransfer(transfer, toAddress))
	}
	
	return collectionTransfers
}

func CastOutgoingTransferToCollectionTransfer(transfer *UserApprovedOutgoingTransfer, fromAddress string) *CollectionApprovedTransfer {
	allowedCombinations := []*IsCollectionTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastOutgoingCombinationToCollectionCombination(combination))
	}
	
	return &CollectionApprovedTransfer{
		ToMappingId: transfer.ToMappingId,
		FromMappingId: fromAddress,
		InitiatedByMappingId: transfer.InitiatedByMappingId,
		TransferTimes: transfer.TransferTimes,
		BadgeIds: transfer.BadgeIds,
		AllowedCombinations: allowedCombinations,
		PerAddressApprovals: transfer.PerAddressApprovals,
		IncrementIdsBy: transfer.IncrementIdsBy,
		IncrementTimesBy: transfer.IncrementTimesBy,
		RequireToEqualsInitiatedBy: transfer.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: transfer.RequireToDoesNotEqualInitiatedBy,
		CustomData: transfer.CustomData,
		TrackerId: transfer.TrackerId,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedOutgoingTransfer {
	allowedCombinations := []*IsUserOutgoingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToOutgoingCombination(combination))
	}
	
	return &UserApprovedOutgoingTransfer{
		ToMappingId: transfer.ToMappingId,
		InitiatedByMappingId: transfer.InitiatedByMappingId,
		TransferTimes: transfer.TransferTimes,
		BadgeIds: transfer.BadgeIds,
		AllowedCombinations: allowedCombinations,
		PerAddressApprovals: transfer.PerAddressApprovals,
		IncrementIdsBy: transfer.IncrementIdsBy,
		IncrementTimesBy: transfer.IncrementTimesBy,
		RequireToEqualsInitiatedBy: transfer.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: transfer.RequireToDoesNotEqualInitiatedBy,
		CustomData: transfer.CustomData,
		TrackerId: transfer.TrackerId,
	}
}

func CastOutgoingCombinationToCollectionCombination(combination *IsUserOutgoingTransferAllowed) *IsCollectionTransferAllowed {
	return &IsCollectionTransferAllowed{
		IsAllowed: combination.IsAllowed,
		InvertBadgeIds: combination.InvertBadgeIds,
		InvertTransferTimes: combination.InvertTransferTimes,
		InvertTo: combination.InvertTo,
		InvertInitiatedBy: combination.InvertInitiatedBy,
	}
}

func CastFromCollectionCombinationToOutgoingCombination(combination *IsCollectionTransferAllowed) *IsUserOutgoingTransferAllowed {
	return &IsUserOutgoingTransferAllowed{
		IsAllowed: combination.IsAllowed,
		InvertBadgeIds: combination.InvertBadgeIds,
		InvertTransferTimes: combination.InvertTransferTimes,
		InvertTo: combination.InvertTo,
		InvertInitiatedBy: combination.InvertInitiatedBy,
	}
}

func CastIncomingTransferToCollectionTransfer(transfer *UserApprovedIncomingTransfer, toAddress string) *CollectionApprovedTransfer {
	allowedCombinations := []*IsCollectionTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastIncomingCombinationToCollectionCombination(combination))
	}
	
	return &CollectionApprovedTransfer{
		ToMappingId: toAddress,
		FromMappingId: transfer.FromMappingId,
		InitiatedByMappingId: transfer.InitiatedByMappingId,
		TransferTimes: transfer.TransferTimes,
		BadgeIds: transfer.BadgeIds,
		AllowedCombinations: allowedCombinations,
		PerAddressApprovals: transfer.PerAddressApprovals,
		IncrementIdsBy: transfer.IncrementIdsBy,
		IncrementTimesBy: transfer.IncrementTimesBy,
		RequireFromEqualsInitiatedBy: transfer.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: transfer.RequireFromDoesNotEqualInitiatedBy,
		CustomData: transfer.CustomData,
		TrackerId: transfer.TrackerId,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedIncomingTransfer {
	allowedCombinations := []*IsUserIncomingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToIncomingCombination(combination))
	}
	
	return &UserApprovedIncomingTransfer{
		FromMappingId: transfer.FromMappingId,
		InitiatedByMappingId: transfer.InitiatedByMappingId,
		TransferTimes: transfer.TransferTimes,
		BadgeIds: transfer.BadgeIds,
		AllowedCombinations: allowedCombinations,
		PerAddressApprovals: transfer.PerAddressApprovals,
		IncrementIdsBy: transfer.IncrementIdsBy,
		IncrementTimesBy: transfer.IncrementTimesBy,
		RequireFromEqualsInitiatedBy: transfer.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: transfer.RequireFromDoesNotEqualInitiatedBy,
		CustomData: transfer.CustomData,
		TrackerId: transfer.TrackerId,
	}
}

func CastIncomingCombinationToCollectionCombination(combination *IsUserIncomingTransferAllowed) *IsCollectionTransferAllowed {
	return &IsCollectionTransferAllowed{
		IsAllowed: combination.IsAllowed,
		InvertBadgeIds: combination.InvertBadgeIds,
		InvertTransferTimes: combination.InvertTransferTimes,
		InvertTo: combination.InvertFrom,
		InvertInitiatedBy: combination.InvertInitiatedBy,
	}
}

func CastFromCollectionCombinationToIncomingCombination(combination *IsCollectionTransferAllowed) *IsUserIncomingTransferAllowed {
	return &IsUserIncomingTransferAllowed{
		IsAllowed: combination.IsAllowed,
		InvertBadgeIds: combination.InvertBadgeIds,
		InvertTransferTimes: combination.InvertTransferTimes,
		InvertFrom: combination.InvertTo,
		InvertInitiatedBy: combination.InvertInitiatedBy,
	}
}
