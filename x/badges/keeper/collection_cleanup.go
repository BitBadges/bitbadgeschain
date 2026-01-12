package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
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
	balancePrefix := append(UserBalanceKey, []byte(collectionIdStr+BalanceKeyDelimiter)...)
	balanceIterator := storetypes.KVStorePrefixIterator(store, balancePrefix)
	balanceKeysToDelete := []string{}
	for ; balanceIterator.Valid(); balanceIterator.Next() {
		// Extract balance key (remove the prefix byte)
		balanceKey := string(balanceIterator.Key()[len(UserBalanceKey):])
		balanceKeysToDelete = append(balanceKeysToDelete, balanceKey)
	}
	balanceIterator.Close()

	// Delete all balance keys
	for _, balanceKey := range balanceKeysToDelete {
		k.DeleteUserBalanceFromStore(ctx, balanceKey)
	}

	// 2. Purge all challenge trackers for this collection
	// Challenge tracker keys format: collectionId-address-approvalLevel-approvalId-challengeId-leafIndex
	challengePrefix := append(UsedClaimChallengeKey, []byte(collectionIdStr+BalanceKeyDelimiter)...)
	challengeIterator := storetypes.KVStorePrefixIterator(store, challengePrefix)
	challengeKeysToDelete := []string{}
	for ; challengeIterator.Valid(); challengeIterator.Next() {
		// Extract challenge key (remove the prefix byte)
		challengeKey := string(challengeIterator.Key()[len(UsedClaimChallengeKey):])
		challengeKeysToDelete = append(challengeKeysToDelete, challengeKey)
		store.Delete(challengeIterator.Key())
	}
	challengeIterator.Close()

	// 3. Purge all approval trackers for this collection
	// Approval tracker keys format: collectionId-addressForApproval-approvalId-amountTrackerId-level-trackerType-address
	approvalPrefix := append(ApprovalTrackerKey, []byte(collectionIdStr+BalanceKeyDelimiter)...)
	approvalIterator := storetypes.KVStorePrefixIterator(store, approvalPrefix)
	approvalKeysToDelete := []string{}
	for ; approvalIterator.Valid(); approvalIterator.Next() {
		// Extract approval tracker key (remove the prefix byte)
		approvalKey := string(approvalIterator.Key()[len(ApprovalTrackerKey):])
		approvalKeysToDelete = append(approvalKeysToDelete, approvalKey)
		store.Delete(approvalIterator.Key())
	}
	approvalIterator.Close()

	// 4. Purge all approval versions for this collection
	// Approval version keys format: collectionId-approvalLevel-approverAddress-approvalId
	versionPrefix := append(ApprovalVersionKey, []byte(collectionIdStr+BalanceKeyDelimiter)...)
	versionIterator := storetypes.KVStorePrefixIterator(store, versionPrefix)
	versionKeysToDelete := []string{}
	for ; versionIterator.Valid(); versionIterator.Next() {
		// Extract approval version key (remove the prefix byte)
		versionKey := string(versionIterator.Key()[len(ApprovalVersionKey):])
		versionKeysToDelete = append(versionKeysToDelete, versionKey)
		store.Delete(versionIterator.Key())
	}
	versionIterator.Close()

	// 5. Purge all ETH signature trackers for this collection
	// ETH signature tracker keys format: collectionId-addressForChallenge-approvalLevel-approvalId-challengeId-signature
	ethPrefix := append(ETHSignatureTrackerKey, []byte(collectionIdStr+BalanceKeyDelimiter)...)
	ethIterator := storetypes.KVStorePrefixIterator(store, ethPrefix)
	ethKeysToDelete := []string{}
	for ; ethIterator.Valid(); ethIterator.Next() {
		// Extract ETH signature tracker key (remove the prefix)
		ethKey := string(ethIterator.Key()[len(ETHSignatureTrackerKey):])
		ethKeysToDelete = append(ethKeysToDelete, ethKey)
		store.Delete(ethIterator.Key())
	}
	ethIterator.Close()

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
