package types

func CastOutgoingTransfersToCollectionTransfers(transfers []*UserOutgoingApproval, fromAddress string) []*CollectionApproval {
	collectionTransfers := []*CollectionApproval{}
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastOutgoingTransferToCollectionTransfer(transfer, fromAddress))
	}

	return collectionTransfers
}

func CastIncomingTransfersToCollectionTransfers(transfers []*UserIncomingApproval, toAddress string) []*CollectionApproval {
	collectionTransfers := []*CollectionApproval{}
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastIncomingTransferToCollectionTransfer(transfer, toAddress))
	}

	return collectionTransfers
}

func CastOutgoingTransferToCollectionTransfer(transfer *UserOutgoingApproval, fromAddress string) *CollectionApproval {

	approvalCriteria := CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(transfer.ApprovalCriteria)
	return &CollectionApproval{
		ToMappingId:                      transfer.ToMappingId,
		FromMappingId:                    fromAddress,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		ApprovalCriteria: 								approvalCriteria,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
		Uri: transfer.Uri,
		CustomData: transfer.CustomData,

		IsApproved:           transfer.IsApproved,
		BadgeIdsOptions:      transfer.BadgeIdsOptions,
		TransferTimesOptions: transfer.TransferTimesOptions,
		ToMappingOptions:            transfer.ToMappingOptions,
		InitiatedByMappingOptions:   transfer.InitiatedByMappingOptions,
		OwnershipTimesOptions: transfer.OwnershipTimesOptions,
		ApprovalTrackerIdOptions: transfer.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: transfer.ChallengeTrackerIdOptions,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApproval) *UserOutgoingApproval {

	approvalCriteria := CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserOutgoingApproval{
		ToMappingId:                      transfer.ToMappingId,
		InitiatedByMappingId:             transfer.InitiatedByMappingId,
		TransferTimes:                    transfer.TransferTimes,
		BadgeIds:                         transfer.BadgeIds,
		OwnershipTimes: 								  transfer.OwnershipTimes,
		ApprovalCriteria: 								approvalCriteria,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
		Uri: transfer.Uri,
		CustomData: transfer.CustomData,

		IsApproved:           transfer.IsApproved,
		BadgeIdsOptions:      transfer.BadgeIdsOptions,
		TransferTimesOptions: transfer.TransferTimesOptions,
		ToMappingOptions:            transfer.ToMappingOptions,
		InitiatedByMappingOptions:   transfer.InitiatedByMappingOptions,
		OwnershipTimesOptions: transfer.OwnershipTimesOptions,
		ApprovalTrackerIdOptions: transfer.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: transfer.ChallengeTrackerIdOptions,
	}
}


func CastIncomingTransferToCollectionTransfer(transfer *UserIncomingApproval, toAddress string) *CollectionApproval {

	approvalCriteria := CastIncomingApprovalCriteriaToCollectionApprovalCriteria(transfer.ApprovalCriteria)

	return &CollectionApproval{
		ToMappingId:                        toAddress,
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  	transfer.OwnershipTimes,
		ApprovalCriteria:                    approvalCriteria,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
		Uri: transfer.Uri,
		CustomData: transfer.CustomData,

		IsApproved:           transfer.IsApproved,
		BadgeIdsOptions:      transfer.BadgeIdsOptions,
		TransferTimesOptions: transfer.TransferTimesOptions,
		FromMappingOptions:            transfer.FromMappingOptions,
		InitiatedByMappingOptions:   transfer.InitiatedByMappingOptions,
		OwnershipTimesOptions: transfer.OwnershipTimesOptions,
		ApprovalTrackerIdOptions: transfer.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: transfer.ChallengeTrackerIdOptions,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApproval) *UserIncomingApproval {

	approvalCriteria := CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserIncomingApproval{
		FromMappingId:                      transfer.FromMappingId,
		InitiatedByMappingId:               transfer.InitiatedByMappingId,
		TransferTimes:                      transfer.TransferTimes,
		BadgeIds:                           transfer.BadgeIds,
		OwnershipTimes: 								  			transfer.OwnershipTimes,
		ApprovalCriteria:                    approvalCriteria,
		ApprovalId: transfer.ApprovalId,
		ApprovalTrackerId: transfer.ApprovalTrackerId,
		ChallengeTrackerId: transfer.ChallengeTrackerId,
		Uri: transfer.Uri,
		CustomData: transfer.CustomData,

		IsApproved:           transfer.IsApproved,
		BadgeIdsOptions:      transfer.BadgeIdsOptions,
		TransferTimesOptions: transfer.TransferTimesOptions,
		FromMappingOptions:          transfer.FromMappingOptions,
		InitiatedByMappingOptions:   transfer.InitiatedByMappingOptions,
		OwnershipTimesOptions: transfer.OwnershipTimesOptions,
		ApprovalTrackerIdOptions: transfer.ApprovalTrackerIdOptions,
		ChallengeTrackerIdOptions: transfer.ChallengeTrackerIdOptions,
	}
}

func CastIncomingApprovalCriteriaToCollectionApprovalCriteria(approvalCriteria *IncomingApprovalCriteria) *ApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &ApprovalCriteria{
		ApprovalAmounts: approvalCriteria.ApprovalAmounts,
		MaxNumTransfers: approvalCriteria.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalCriteria.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalCriteria.RequireFromDoesNotEqualInitiatedBy,
		PredeterminedBalances: approvalCriteria.PredeterminedBalances,
		MustOwnBadges: approvalCriteria.MustOwnBadges,
		MerkleChallenge: approvalCriteria.MerkleChallenge,
	}
}

func CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(approvalCriteria *OutgoingApprovalCriteria) *ApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &ApprovalCriteria{
		ApprovalAmounts: approvalCriteria.ApprovalAmounts,
		MaxNumTransfers: approvalCriteria.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalCriteria.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalCriteria.RequireToDoesNotEqualInitiatedBy,
		PredeterminedBalances: approvalCriteria.PredeterminedBalances,
		MustOwnBadges: approvalCriteria.MustOwnBadges,
		MerkleChallenge: approvalCriteria.MerkleChallenge,
	}
}

func CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(approvalCriteria *ApprovalCriteria) *IncomingApprovalCriteria {
	return &IncomingApprovalCriteria{
		ApprovalAmounts: approvalCriteria.ApprovalAmounts,
		MaxNumTransfers: approvalCriteria.MaxNumTransfers,
		RequireFromEqualsInitiatedBy: approvalCriteria.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalCriteria.RequireFromDoesNotEqualInitiatedBy,
		PredeterminedBalances: approvalCriteria.PredeterminedBalances,
		MustOwnBadges: approvalCriteria.MustOwnBadges,
		MerkleChallenge: approvalCriteria.MerkleChallenge,
	}
}

func CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(approvalCriteria *ApprovalCriteria) *OutgoingApprovalCriteria {
	return &OutgoingApprovalCriteria{
		ApprovalAmounts: approvalCriteria.ApprovalAmounts,
		MaxNumTransfers: approvalCriteria.MaxNumTransfers,
		RequireToEqualsInitiatedBy: approvalCriteria.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalCriteria.RequireToDoesNotEqualInitiatedBy,
		PredeterminedBalances: approvalCriteria.PredeterminedBalances,
		MustOwnBadges: approvalCriteria.MustOwnBadges,
		MerkleChallenge: approvalCriteria.MerkleChallenge,
	}
}