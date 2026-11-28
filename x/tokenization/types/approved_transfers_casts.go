package types

func CastOutgoingTransfersToCollectionTransfers(transfers []*UserOutgoingApproval, fromAddress string) []*CollectionApproval {
	collectionTransfers := make([]*CollectionApproval, 0, len(transfers))
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastOutgoingTransferToCollectionTransfer(transfer, fromAddress))
	}

	return collectionTransfers
}

func CastIncomingTransfersToCollectionTransfers(transfers []*UserIncomingApproval, toAddress string) []*CollectionApproval {
	collectionTransfers := make([]*CollectionApproval, 0, len(transfers))
	for _, transfer := range transfers {
		collectionTransfers = append(collectionTransfers, CastIncomingTransferToCollectionTransfer(transfer, toAddress))
	}

	return collectionTransfers
}

func CastOutgoingTransferToCollectionTransfer(transfer *UserOutgoingApproval, fromAddress string) *CollectionApproval {
	if transfer == nil {
		return nil
	}
	approvalCriteria := CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(transfer.ApprovalCriteria)
	return &CollectionApproval{
		ToListId:          transfer.ToListId,
		FromListId:        fromAddress,
		InitiatedByListId: transfer.InitiatedByListId,
		TransferTimes:     transfer.TransferTimes,
		TokenIds:          transfer.TokenIds,
		OwnershipTimes:    transfer.OwnershipTimes,
		ApprovalCriteria:  approvalCriteria,
		ApprovalId:        transfer.ApprovalId,
		Uri:               transfer.Uri,
		CustomData:        transfer.CustomData,
		Version:           transfer.Version,
	}
}

func CastFromCollectionTransferToOutgoingTransfer(transfer *CollectionApproval) *UserOutgoingApproval {
	if transfer == nil {
		return nil
	}
	approvalCriteria := CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserOutgoingApproval{
		ToListId:          transfer.ToListId,
		InitiatedByListId: transfer.InitiatedByListId,
		TransferTimes:     transfer.TransferTimes,
		TokenIds:          transfer.TokenIds,
		OwnershipTimes:    transfer.OwnershipTimes,
		ApprovalCriteria:  approvalCriteria,
		ApprovalId:        transfer.ApprovalId,
		Uri:               transfer.Uri,
		CustomData:        transfer.CustomData,
		Version:           transfer.Version,
	}
}

func CastIncomingTransferToCollectionTransfer(transfer *UserIncomingApproval, toAddress string) *CollectionApproval {
	if transfer == nil {
		return nil
	}
	approvalCriteria := CastIncomingApprovalCriteriaToCollectionApprovalCriteria(transfer.ApprovalCriteria)

	return &CollectionApproval{
		ToListId:          toAddress,
		FromListId:        transfer.FromListId,
		InitiatedByListId: transfer.InitiatedByListId,
		TransferTimes:     transfer.TransferTimes,
		TokenIds:          transfer.TokenIds,
		OwnershipTimes:    transfer.OwnershipTimes,
		ApprovalCriteria:  approvalCriteria,
		ApprovalId:        transfer.ApprovalId,
		Uri:               transfer.Uri,
		CustomData:        transfer.CustomData,
		Version:           transfer.Version,
	}
}

func CastFromCollectionTransferToIncomingTransfer(transfer *CollectionApproval) *UserIncomingApproval {
	if transfer == nil {
		return nil
	}
	approvalCriteria := CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(transfer.ApprovalCriteria)

	return &UserIncomingApproval{
		FromListId:        transfer.FromListId,
		InitiatedByListId: transfer.InitiatedByListId,
		TransferTimes:     transfer.TransferTimes,
		TokenIds:          transfer.TokenIds,
		OwnershipTimes:    transfer.OwnershipTimes,
		ApprovalCriteria:  approvalCriteria,
		ApprovalId:        transfer.ApprovalId,
		Uri:               transfer.Uri,
		CustomData:        transfer.CustomData,
		Version:           transfer.Version,
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
		MerkleChallenges:                   approvalCriteria.MerkleChallenges,
		CoinTransfers:                      approvalCriteria.CoinTransfers,
		AutoDeletionOptions:                approvalCriteria.AutoDeletionOptions,
		MustOwnTokens:                      approvalCriteria.MustOwnTokens,
		DynamicStoreChallenges:             approvalCriteria.DynamicStoreChallenges,
		EthSignatureChallenges:             approvalCriteria.EthSignatureChallenges,
		VotingChallenges:                   approvalCriteria.VotingChallenges,
		SenderChecks:                       approvalCriteria.SenderChecks,
		InitiatorChecks:                    approvalCriteria.InitiatorChecks,
		AltTimeChecks:                      approvalCriteria.AltTimeChecks,
		MustPrioritize:                     approvalCriteria.MustPrioritize,
		EvmQueryChallenges:                 approvalCriteria.EvmQueryChallenges,
		AllowBackedMinting:                 false, // Not applicable to incoming approvals
		AllowSpecialWrapping:               false, // Not applicable to incoming approvals
		// Note: IncomingApprovalCriteria doesn't have recipientChecks
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
		MerkleChallenges:                 approvalCriteria.MerkleChallenges,
		CoinTransfers:                    approvalCriteria.CoinTransfers,
		AutoDeletionOptions:              approvalCriteria.AutoDeletionOptions,
		MustOwnTokens:                    approvalCriteria.MustOwnTokens,
		DynamicStoreChallenges:           approvalCriteria.DynamicStoreChallenges,
		EthSignatureChallenges:           approvalCriteria.EthSignatureChallenges,
		VotingChallenges:                 approvalCriteria.VotingChallenges,
		RecipientChecks:                  approvalCriteria.RecipientChecks,
		InitiatorChecks:                  approvalCriteria.InitiatorChecks,
		AltTimeChecks:                    approvalCriteria.AltTimeChecks,
		MustPrioritize:                   approvalCriteria.MustPrioritize,
		EvmQueryChallenges:               approvalCriteria.EvmQueryChallenges,
		AllowBackedMinting:               false, // Not applicable to outgoing approvals
		AllowSpecialWrapping:             false, // Not applicable to outgoing approvals
		// Note: OutgoingApprovalCriteria doesn't have senderChecks
	}
}

func CastFromCollectionApprovalCriteriaToIncomingApprovalCriteria(approvalCriteria *ApprovalCriteria) *IncomingApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &IncomingApprovalCriteria{
		ApprovalAmounts:                    approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                    approvalCriteria.MaxNumTransfers,
		RequireFromEqualsInitiatedBy:       approvalCriteria.RequireFromEqualsInitiatedBy,
		RequireFromDoesNotEqualInitiatedBy: approvalCriteria.RequireFromDoesNotEqualInitiatedBy,
		PredeterminedBalances:              approvalCriteria.PredeterminedBalances,
		MerkleChallenges:                   approvalCriteria.MerkleChallenges,
		CoinTransfers:                      approvalCriteria.CoinTransfers,
		AutoDeletionOptions:                approvalCriteria.AutoDeletionOptions,
		MustOwnTokens:                      approvalCriteria.MustOwnTokens,
		DynamicStoreChallenges:             approvalCriteria.DynamicStoreChallenges,
		EthSignatureChallenges:             approvalCriteria.EthSignatureChallenges,
		VotingChallenges:                   approvalCriteria.VotingChallenges,
		SenderChecks:                       approvalCriteria.SenderChecks,
		InitiatorChecks:                    approvalCriteria.InitiatorChecks,
		AltTimeChecks:                      approvalCriteria.AltTimeChecks,
		MustPrioritize:                     approvalCriteria.MustPrioritize,
		EvmQueryChallenges:                 approvalCriteria.EvmQueryChallenges,
		// Note: IncomingApprovalCriteria doesn't have recipientChecks
	}
}

func CastFromCollectionApprovalCriteriaToOutgoingApprovalCriteria(approvalCriteria *ApprovalCriteria) *OutgoingApprovalCriteria {
	if approvalCriteria == nil {
		return nil
	}

	return &OutgoingApprovalCriteria{
		ApprovalAmounts:                  approvalCriteria.ApprovalAmounts,
		MaxNumTransfers:                  approvalCriteria.MaxNumTransfers,
		RequireToEqualsInitiatedBy:       approvalCriteria.RequireToEqualsInitiatedBy,
		RequireToDoesNotEqualInitiatedBy: approvalCriteria.RequireToDoesNotEqualInitiatedBy,
		PredeterminedBalances:            approvalCriteria.PredeterminedBalances,
		MerkleChallenges:                 approvalCriteria.MerkleChallenges,
		CoinTransfers:                    approvalCriteria.CoinTransfers,
		AutoDeletionOptions:              approvalCriteria.AutoDeletionOptions,
		MustOwnTokens:                    approvalCriteria.MustOwnTokens,
		DynamicStoreChallenges:           approvalCriteria.DynamicStoreChallenges,
		EthSignatureChallenges:           approvalCriteria.EthSignatureChallenges,
		VotingChallenges:                 approvalCriteria.VotingChallenges,
		RecipientChecks:                  approvalCriteria.RecipientChecks,
		InitiatorChecks:                  approvalCriteria.InitiatorChecks,
		AltTimeChecks:                    approvalCriteria.AltTimeChecks,
		MustPrioritize:                   approvalCriteria.MustPrioritize,
		EvmQueryChallenges:               approvalCriteria.EvmQueryChallenges,
		// Note: OutgoingApprovalCriteria doesn't have senderChecks
	}
}
