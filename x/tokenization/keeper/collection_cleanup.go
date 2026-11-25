package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PurgeCollectionState removes all state associated with a collection
// This includes balances, approval trackers, challenge trackers, approval versions, and ETH signature trackers
func (k Keeper) PurgeCollectionState(ctx sdk.Context, collectionId sdkmath.Uint) error {
	collectionIdStr := collectionId.String()
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	// 1. Purge all balances for this collection
	// Balance keys format: collectionId-address
	balancePrefix := append(append([]byte(nil), UserBalanceKey...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
	balanceIterator := storetypes.KVStorePrefixIterator(store, balancePrefix)
	balanceKeysToDelete := []string{}
	for ; balanceIterator.Valid(); balanceIterator.Next() {
		// Extract balance key (remove the prefix byte)
		balanceKey := string(balanceIterator.Key()[len(UserBalanceKey):])
		balanceKeysToDelete = append(balanceKeysToDelete, balanceKey)
	}
	if err := balanceIterator.Close(); err != nil {
		k.logger.Error("failed to close balance iterator", "error", err)
	}

	// Delete all balance keys
	for _, balanceKey := range balanceKeysToDelete {
		k.DeleteUserBalanceFromStore(ctx, balanceKey)
	}

	// 2. Purge all challenge trackers for this collection
	// Challenge tracker keys format: collectionId-address-approvalLevel-approvalId-challengeId-leafIndex
	challengePrefix := append(append([]byte(nil), UsedClaimChallengeKey...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
	challengeIterator := storetypes.KVStorePrefixIterator(store, challengePrefix)
	challengeKeysToDelete := [][]byte{}
	for ; challengeIterator.Valid(); challengeIterator.Next() {
		challengeKeysToDelete = append(challengeKeysToDelete, append([]byte(nil), challengeIterator.Key()...))
	}
	if err := challengeIterator.Close(); err != nil {
		k.logger.Error("failed to close challenge iterator", "error", err)
	}
	for _, key := range challengeKeysToDelete {
		store.Delete(key)
	}

	// 3. Purge all approval trackers for this collection
	// Approval tracker keys format: collectionId-addressForApproval-approvalId-amountTrackerId-level-trackerType-address
	approvalPrefix := append(append([]byte(nil), ApprovalTrackerKey...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
	approvalIterator := storetypes.KVStorePrefixIterator(store, approvalPrefix)
	approvalKeysToDelete := [][]byte{}
	for ; approvalIterator.Valid(); approvalIterator.Next() {
		approvalKeysToDelete = append(approvalKeysToDelete, append([]byte(nil), approvalIterator.Key()...))
	}
	if err := approvalIterator.Close(); err != nil {
		k.logger.Error("failed to close approval iterator", "error", err)
	}
	for _, key := range approvalKeysToDelete {
		store.Delete(key)
	}

	// 4. Purge all approval versions for this collection
	// Approval version keys format: collectionId-approvalLevel-approverAddress-approvalId
	versionPrefix := append(append([]byte(nil), ApprovalVersionKey...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
	versionIterator := storetypes.KVStorePrefixIterator(store, versionPrefix)
	versionKeysToDelete := [][]byte{}
	for ; versionIterator.Valid(); versionIterator.Next() {
		versionKeysToDelete = append(versionKeysToDelete, append([]byte(nil), versionIterator.Key()...))
	}
	if err := versionIterator.Close(); err != nil {
		k.logger.Error("failed to close version iterator", "error", err)
	}
	for _, key := range versionKeysToDelete {
		store.Delete(key)
	}

	// 5. Purge all ETH signature trackers for this collection
	// ETH signature tracker keys format: collectionId-addressForChallenge-approvalLevel-approvalId-challengeId-signature
	ethPrefix := append(append([]byte(nil), ETHSignatureTrackerKey...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
	ethIterator := storetypes.KVStorePrefixIterator(store, ethPrefix)
	ethKeysToDelete := [][]byte{}
	for ; ethIterator.Valid(); ethIterator.Next() {
		ethKeysToDelete = append(ethKeysToDelete, append([]byte(nil), ethIterator.Key()...))
	}
	if err := ethIterator.Close(); err != nil {
		k.logger.Error("failed to close eth iterator", "error", err)
	}
	for _, key := range ethKeysToDelete {
		store.Delete(key)
	}

	// 6. Purge votes and voting challenge trackers for this collection
	// Keys format: collectionId-approverAddress-approvalLevel-approvalId-proposalId[-voter]
	for _, votePrefix := range [][]byte{VotingTrackerKey, VotingChallengeTrackerKey} {
		fullPrefix := append(append([]byte(nil), votePrefix...), []byte(collectionIdStr+BalanceKeyDelimiter)...)
		voteIterator := storetypes.KVStorePrefixIterator(store, fullPrefix)
		voteKeysToDelete := [][]byte{}
		for ; voteIterator.Valid(); voteIterator.Next() {
			voteKeysToDelete = append(voteKeysToDelete, append([]byte(nil), voteIterator.Key()...))
		}
		if err := voteIterator.Close(); err != nil {
			k.logger.Error("failed to close vote iterator", "error", err)
		}
		for _, key := range voteKeysToDelete {
			store.Delete(key)
		}
	}

	// 7. Purge collection stats
	if !collectionId.IsNil() && collectionId.BigInt().IsUint64() {
		store.Delete(collectionStatsStoreKey(collectionId))
	}

	// Log the cleanup for monitoring
	k.logger.Info("purged collection state",
		"collection_id", collectionIdStr,
		"balances_deleted", len(balanceKeysToDelete),
		"challenge_trackers_deleted", len(challengeKeysToDelete),
		"approval_trackers_deleted", len(approvalKeysToDelete),
		"approval_versions_deleted", len(versionKeysToDelete),
		"eth_signature_trackers_deleted", len(ethKeysToDelete),
	)

	return nil
}
