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
		ToListId:           transfer.ToListId,
		FromListId:         fromAddress,
		InitiatedByListId:  transfer.InitiatedByListId,
		TransferTimes:      transfer.TransferTimes,
		BadgeIds:           transfer.BadgeIds,
		OwnershipTimes:     transfer.OwnershipTimes,
		ApprovalCriteria:   approvalCriteria,
		ApprovalId:         transfer.ApprovalId,
		Uri:                transfer.Uri,
		CustomData:         transfer.CustomData,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApproval) *UserOutgoingApproval {

	approvalCriteria := CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserOutgoingApproval{
		ToListId:           transfer.ToListId,
		InitiatedByListId:  transfer.InitiatedByListId,
		TransferTimes:      transfer.TransferTimes,
		BadgeIds:           transfer.BadgeIds,
		OwnershipTimes:     transfer.OwnershipTimes,
		ApprovalCriteria:   approvalCriteria,
		ApprovalId:         transfer.ApprovalId,
		Uri:                transfer.Uri,
		CustomData:         transfer.CustomData,
	}
}

func CastIncomingTransferToCollectionTransfer(transfer *UserIncomingApproval, toAddress string) *CollectionApproval {

	approvalCriteria := CastIncomingApprovalCriteriaToCollectionApprovalCriteria(transfer.ApprovalCriteria)

	return &CollectionApproval{
		ToListId:           toAddress,
		FromListId:         transfer.FromListId,
		InitiatedByListId:  transfer.InitiatedByListId,
		TransferTimes:      transfer.TransferTimes,
		BadgeIds:           transfer.BadgeIds,
		OwnershipTimes:     transfer.OwnershipTimes,
		ApprovalCriteria:   approvalCriteria,
		ApprovalId:         transfer.ApprovalId,
		Uri:                transfer.Uri,
		CustomData:         transfer.CustomData,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApproval) *UserIncomingApproval {

	approvalCriteria := CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserIncomingApproval{
		FromListId:         transfer.FromListId,
		InitiatedByListId:  transfer.InitiatedByListId,
		TransferTimes:      transfer.TransferTimes,
		BadgeIds:           transfer.BadgeIds,
		OwnershipTimes:     transfer.OwnershipTimes,
		ApprovalCriteria:   approvalCriteria,
		ApprovalId:         transfer.ApprovalId,
		Uri:                transfer.Uri,
		CustomData:         transfer.CustomData,
	}
}

func CastIncomingApprovalCriteriaToCollectionApprovalCriteria(approvalCriteria *IncomingApprovalCriteria) *ApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &ApprovalCriteria{
		ApprovalAmounts:                    approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                    approvalCriteria.MaxNumTransfers,
		RequireFromEqualsInitiatedBy:       approvalCriteria.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalCriteria.RequireFromDoesNotEqualInitiatedBy,
		PredeterminedBalances:              approvalCriteria.PredeterminedBalances,
		MustOwnBadges:                      approvalCriteria.MustOwnBadges,
		MerkleChallenges:                    approvalCriteria.MerkleChallenges,
	}
}

func CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(approvalCriteria *OutgoingApprovalCriteria) *ApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &ApprovalCriteria{
		ApprovalAmounts:                  approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                  approvalCriteria.MaxNumTransfers,
		RequireToEqualsInitiatedBy:       approvalCriteria.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalCriteria.RequireToDoesNotEqualInitiatedBy,
		PredeterminedBalances:            approvalCriteria.PredeterminedBalances,
		MustOwnBadges:                    approvalCriteria.MustOwnBadges,
		MerkleChallenges:                  approvalCriteria.MerkleChallenges,
	}
}

func CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(approvalCriteria *ApprovalCriteria) *IncomingApprovalCriteria {
	return &IncomingApprovalCriteria{
		ApprovalAmounts:                    approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                    approvalCriteria.MaxNumTransfers,
		RequireFromEqualsInitiatedBy:       approvalCriteria.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalCriteria.RequireFromDoesNotEqualInitiatedBy,
		PredeterminedBalances:              approvalCriteria.PredeterminedBalances,
		MustOwnBadges:                      approvalCriteria.MustOwnBadges,
		MerkleChallenges:                    approvalCriteria.MerkleChallenges,
	}
}

func CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(approvalCriteria *ApprovalCriteria) *OutgoingApprovalCriteria {
	return &OutgoingApprovalCriteria{
		ApprovalAmounts:                  approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                  approvalCriteria.MaxNumTransfers,
		RequireToEqualsInitiatedBy:       approvalCriteria.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalCriteria.RequireToDoesNotEqualInitiatedBy,
		PredeterminedBalances:            approvalCriteria.PredeterminedBalances,
		MustOwnBadges:                    approvalCriteria.MustOwnBadges,
		MerkleChallenges:                  approvalCriteria.MerkleChallenges,
	}
}
