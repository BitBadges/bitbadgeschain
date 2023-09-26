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

	approvalDetails := CastOutgoingApprovalDetailsToCollectionApprovalDetails(transfer.ApprovalDetails)

	return &CollectionApprovedTransfer{
		ToMappingId:                      transfer.ToMappingId,
		FromMappingId:                    fromAddress,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		AllowedCombinations:              allowedCombinations,
		ApprovalDetails: 								approvalDetails,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedOutgoingTransfer {
	allowedCombinations := []*IsUserOutgoingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToOutgoingCombination(combination))
	}

	approvalDetails := CastFromCollectionApprovalDetailsToOutgoingApprovalDetails(transfer.ApprovalDetails)

	return &UserApprovedOutgoingTransfer{
		ToMappingId:                      transfer.ToMappingId,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		AllowedCombinations:              allowedCombinations,
		ApprovalDetails: 								approvalDetails,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
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
		ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
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
		ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
	}
}

func CastIncomingTransferToCollectionTransfer(transfer *UserApprovedIncomingTransfer, toAddress string) *CollectionApprovedTransfer {
	allowedCombinations := []*IsCollectionTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastIncomingCombinationToCollectionCombination(combination))
	}

	approvalDetails := CastIncomingApprovalDetailsToCollectionApprovalDetails(transfer.ApprovalDetails)

	return &CollectionApprovedTransfer{
		ToMappingId:                        toAddress,
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  	transfer.OwnershipTimes,
		AllowedCombinations:                allowedCombinations,
		ApprovalDetails:                    approvalDetails,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApprovedTransfer) *UserApprovedIncomingTransfer {
	allowedCombinations := []*IsUserIncomingTransferAllowed{}
	for _, combination := range transfer.AllowedCombinations {
		allowedCombinations = append(allowedCombinations, CastFromCollectionCombinationToIncomingCombination(combination))
	}

	approvalDetails := CastFromCollectionApprovalDetailsToIncomingApprovalDetails(transfer.ApprovalDetails)

	return &UserApprovedIncomingTransfer{
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  			transfer.OwnershipTimes,
		AllowedCombinations:                allowedCombinations,
		ApprovalDetails:                    approvalDetails,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
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
		ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
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
		ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
	}
}


func CastIncomingApprovalDetailsToCollectionApprovalDetails(approvalDetails *IncomingApprovalDetails) *ApprovalDetails {
	if approvalDetails == nil {
		return nil
	}

	return &ApprovalDetails{
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalDetails.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalDetails.RequireFromDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenge: approvalDetails.MerkleChallenge,
	}
}

func CastOutgoingApprovalDetailsToCollectionApprovalDetails(approvalDetails *OutgoingApprovalDetails) *ApprovalDetails {
	if approvalDetails == nil {
		return nil
	}

	return &ApprovalDetails{
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalDetails.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalDetails.RequireToDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenge: approvalDetails.MerkleChallenge,
	}
}

func CastFromCollectionApprovalDetailsToIncomingApprovalDetails(approvalDetails *ApprovalDetails) *IncomingApprovalDetails {
	return &IncomingApprovalDetails{
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalDetails.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalDetails.RequireFromDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenge: approvalDetails.MerkleChallenge,
	}
}

func CastFromCollectionApprovalDetailsToOutgoingApprovalDetails(approvalDetails *ApprovalDetails) *OutgoingApprovalDetails {
	return &OutgoingApprovalDetails{
		ApprovalAmounts: approvalDetails.ApprovalAmounts,
		MaxNumTransfers: approvalDetails.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalDetails.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalDetails.RequireToDoesNotEqualInitiatedBy,
		Uri: approvalDetails.Uri,
		CustomData: approvalDetails.CustomData,
		PredeterminedBalances: approvalDetails.PredeterminedBalances,
		MustOwnBadges: approvalDetails.MustOwnBadges,
		MerkleChallenge: approvalDetails.MerkleChallenge,
	}
}