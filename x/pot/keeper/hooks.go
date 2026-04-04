package keeper

// This file previously contained the transient-queue and hook-based approach
// (EnqueueAddress, DrainQueue, OnCredentialTransfer).
//
// That approach was removed because:
// 1. Time-dependent credential expiry is missed when no transfer event fires.
// 2. Returning ValidatorUpdates from two modules (x/staking + x/pot) causes
//    an instant chain halt ("validator EndBlock updates already set").
//
// The new approach iterates all bonded validators in EndBlocker and uses
// staking Jail/Unjail. See abci.go.
