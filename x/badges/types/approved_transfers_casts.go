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

	approvalDetails := []*ApprovalDetails{}
	for _, approvalDetail := range transfer.ApprovalDetails {
		approvalDetails = append(approvalDetails, CastOutgoingApprovalDetailsToCollectionApprovalDetails(approvalDetail))
	}

	return &CollectionApprovedTransfer{
		ToMappingId:                      transfer.ToMappingId,
		FromMappingId:                    fromAddress,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		AllowedCombinations:              allowedCombinations,
		ApprovalDetails: 								approvalDetails,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedOutgoingTransfer {
	allowedCombinations := []*IsUserOutgoingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToOutgoingCombination(combination))
	}

	approvalDetails := []*OutgoingApprovalDetails{}
	for _, approvalDetail := range transfer.ApprovalDetails {
		approvalDetails = append(approvalDetails, CastFromCollectionApprovalDetailsToOutgoingApprovalDetails(approvalDetail))
	}

	return &UserApprovedOutgoingTransfer{
		ToMappingId:                      transfer.ToMappingId,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		AllowedCombinations:              allowedCombinations,
		ApprovalDetails: 								approvalDetails,
	}
}

func CastOutgoingCombinationToCollectionCombination(combination *IsUserOutgoingTransferAllowed) *IsCollectionTransferAllowed {
	return &IsCollectionTransferAllowed{
		IsApproved:           combination.IsApproved,
		BadgeIdsOptions:      combination.BadgeIdsOptions,
		TransferTimesOptions: combination.TransferTimesOptions,
		ToMappingOptions:            combination.ToMappingOptions,
		InitiatedByMappingOptions:   combination.InitiatedByMappingOptions,
		OwnershipTimesOptions: combination.OwnershipTimesOptions,
	}
}

func CastFromCollectionCombinationToOutgoingCombination(combination *IsCollectionTransferAllowed) *IsUserOutgoingTransferAllowed {
	return &IsUserOutgoingTransferAllowed{
		IsApproved:           combination.IsApproved,
		BadgeIdsOptions:      combination.BadgeIdsOptions,
		TransferTimesOptions: combination.TransferTimesOptions,
		ToMappingOptions:            combination.ToMappingOptions,
		InitiatedByMappingOptions:   combination.InitiatedByMappingOptions,
		OwnershipTimesOptions: combination.OwnershipTimesOptions,
	}
}

func CastIncomingTransferToCollectionTransfer(transfer *UserApprovedIncomingTransfer, toAddress string) *CollectionApprovedTransfer {
	allowedCombinations := []*IsCollectionTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastIncomingCombinationToCollectionCombination(combination))
	}

	approvalDetails := []*ApprovalDetails{}
	for _, approvalDetail := range transfer.ApprovalDetails {
		approvalDetails = append(approvalDetails, CastIncomingApprovalDetailsToCollectionApprovalDetails(approvalDetail))
	}

	return &CollectionApprovedTransfer{
		ToMappingId:                        toAddress,
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  	transfer.OwnershipTimes,
		AllowedCombinations:                allowedCombinations,
		ApprovalDetails:                    approvalDetails,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedIncomingTransfer {
	allowedCombinations := []*IsUserIncomingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToIncomingCombination(combination))
	}

	approvalDetails := []*IncomingApprovalDetails{}
	for _, approvalDetail := range transfer.ApprovalDetails {
		approvalDetails = append(approvalDetails, CastFromCollectionApprovalDetailsToIncomingApprovalDetails(approvalDetail))
	}

	return &UserApprovedIncomingTransfer{
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  			transfer.OwnershipTimes,
		AllowedCombinations:                allowedCombinations,
		ApprovalDetails:                    approvalDetails,
	}
}

func CastIncomingCombinationToCollectionCombination(combination *IsUserIncomingTransferAllowed) *IsCollectionTransferAllowed {
	return &IsCollectionTransferAllowed{
		IsApproved:           combination.IsApproved,
		BadgeIdsOptions:      combination.BadgeIdsOptions,
		TransferTimesOptions: combination.TransferTimesOptions,
		FromMappingOptions:            combination.FromMappingOptions,
		InitiatedByMappingOptions:   combination.InitiatedByMappingOptions,
		OwnershipTimesOptions: combination.OwnershipTimesOptions,
	}
}

func CastFromCollectionCombinationToIncomingCombination(combination *IsCollectionTransferAllowed) *IsUserIncomingTransferAllowed {
	return &IsUserIncomingTransferAllowed{
		IsApproved:           combination.IsApproved,
		BadgeIdsOptions:      combination.BadgeIdsOptions,
		TransferTimesOptions: combination.TransferTimesOptions,
		FromMappingOptions:          combination.FromMappingOptions,
		InitiatedByMappingOptions:   combination.InitiatedByMappingOptions,
		OwnershipTimesOptions: combination.OwnershipTimesOptions,
	}
}


func CastIncomingApprovalDetailsToCollectionApprovalDetails(approvalDetails *IncomingApprovalDetails) *ApprovalDetails {
	return &ApprovalDetails{
		ApprovalId: approvalDetails.ApprovalId,
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalDetails.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalDetails.RequireFromDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenges: approvalDetails.MerkleChallenges,
	}
}

func CastOutgoingApprovalDetailsToCollectionApprovalDetails(approvalDetails *OutgoingApprovalDetails) *ApprovalDetails {
	return &ApprovalDetails{
		ApprovalId: approvalDetails.ApprovalId,
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalDetails.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalDetails.RequireToDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenges: approvalDetails.MerkleChallenges,
	}
}

func CastFromCollectionApprovalDetailsToIncomingApprovalDetails(approvalDetails *ApprovalDetails) *IncomingApprovalDetails {
	return &IncomingApprovalDetails{
		ApprovalId: approvalDetails.ApprovalId,
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalDetails.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalDetails.RequireFromDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenges: approvalDetails.MerkleChallenges,
	}
}

func CastFromCollectionApprovalDetailsToOutgoingApprovalDetails(approvalDetails *ApprovalDetails) *OutgoingApprovalDetails {
	return &OutgoingApprovalDetails{
		ApprovalId: approvalDetails.ApprovalId,
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalDetails.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalDetails.RequireToDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenges: approvalDetails.MerkleChallenges,
	}
}