// Package evmcompat provides utilities for writing code that works correctly
// in both normal Cosmos SDK contexts and EVM precompile contexts.
//
// # The Problem
//
// When running inside an EVM precompile (via RunNativeAction), the context's
// MultiStore is a snapshotmulti.Store that manages EVM state snapshots.
// Calling ctx.CacheContext() in this context is problematic because:
//
//  1. snapshotmulti.Store.CacheMultiStore() returns ITSELF (not a new cache)
//  2. The returned writeCache() function calls Write() which clears the
//     entire snapshot stack (sets cacheStores = nil)
//  3. The EVM's journal still holds references to snapshot indices
//  4. When EVM tries to revert, it panics: "snapshot index X out of bound [0..0)"
//
// # The Solution
//
// This package provides AtomicContext which:
//   - In EVM context: Uses the native Snapshot()/RevertToSnapshot() mechanism
//   - In Cosmos context: Uses the standard CacheContext() pattern
//
// # Usage
//
//	atomic := evmcompat.NewAtomicContext(ctx)
//	defer atomic.Rollback() // Safe no-op if Commit() was called
//
//	// Do operations on atomic.Ctx()
//	if err := doSomething(atomic.Ctx()); err != nil {
//	    return err // Rollback happens via defer
//	}
//
//	atomic.Commit() // Commit changes (no-op in EVM context, writeCache in Cosmos)
//	return nil
package evmcompat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmstoretypes "github.com/cosmos/evm/x/vm/store/types"
)

// AtomicContext provides atomic operations with automatic rollback support
// that works correctly in both EVM precompile and normal Cosmos contexts.
type AtomicContext struct {
	ctx         sdk.Context
	commit      func()
	rollback    func()
	committed   bool
	isEvmContext bool
}

// NewAtomicContext creates a new atomic context for operations that need
// rollback capability. In EVM context, it uses the native snapshot mechanism.
// In normal Cosmos context, it uses CacheContext.
//
// Usage:
//
//	atomic := evmcompat.NewAtomicContext(ctx)
//	// Do operations on atomic.Ctx()
//	if err != nil {
//	    atomic.Rollback()
//	    return err
//	}
//	atomic.Commit()
func NewAtomicContext(ctx sdk.Context) *AtomicContext {
	if snapshotter := tryGetEvmSnapshotter(ctx); snapshotter != nil {
		idx := snapshotter.Snapshot()
		return &AtomicContext{
			ctx:          ctx, // Use original ctx directly in EVM context
			commit:       func() {}, // No-op: changes already on ctx
			rollback:     func() { snapshotter.RevertToSnapshot(idx) },
			isEvmContext: true,
		}
	}

	// Normal Cosmos context: use CacheContext
	cachedCtx, writeCache := ctx.CacheContext()
	return &AtomicContext{
		ctx:          cachedCtx,
		commit:       writeCache,
		rollback:     func() {}, // No-op: just don't call commit
		isEvmContext: false,
	}
}

// Ctx returns the context to use for operations.
// In EVM context, this is the original context.
// In Cosmos context, this is the cached context.
func (a *AtomicContext) Ctx() sdk.Context {
	return a.ctx
}

// Commit finalizes the atomic operations.
// In EVM context, this is a no-op (changes are already on the context).
// In Cosmos context, this writes the cache to the parent context.
// After Commit(), Rollback() becomes a no-op.
func (a *AtomicContext) Commit() {
	if !a.committed {
		a.commit()
		a.committed = true
	}
}

// Rollback discards any uncommitted changes.
// In EVM context, this reverts to the snapshot taken at creation.
// In Cosmos context, this is a no-op (cache is simply discarded).
// Safe to call multiple times or after Commit() (becomes no-op).
func (a *AtomicContext) Rollback() {
	if !a.committed {
		a.rollback()
		a.committed = true // Prevent double rollback
	}
}

// IsEvmContext returns true if running inside an EVM precompile context.
func (a *AtomicContext) IsEvmContext() bool {
	return a.isEvmContext
}

// tryGetEvmSnapshotter checks if the context's MultiStore is an EVM
// snapshotmulti.Store and returns its Snapshotter interface.
// Returns nil if not in EVM context.
func tryGetEvmSnapshotter(ctx sdk.Context) evmstoretypes.Snapshotter {
	if ctx.MultiStore() == nil {
		return nil
	}
	snapshotter, ok := ctx.MultiStore().(evmstoretypes.Snapshotter)
	if !ok {
		return nil
	}
	return snapshotter
}

// IsEvmSnapshotContext returns true if the given context is running inside
// an EVM precompile with snapshotmulti.Store. This can be used to check
// context type without creating an AtomicContext.
func IsEvmSnapshotContext(ctx sdk.Context) bool {
	return tryGetEvmSnapshotter(ctx) != nil
}

// WithAtomicRollback executes a function with automatic rollback on error.
// This is a convenience wrapper around AtomicContext for simple use cases.
//
// Usage:
//
//	err := evmcompat.WithAtomicRollback(ctx, func(ctx sdk.Context) error {
//	    // Do operations that might fail
//	    return possiblyFailingOperation(ctx)
//	})
func WithAtomicRollback(ctx sdk.Context, fn func(ctx sdk.Context) error) error {
	atomic := NewAtomicContext(ctx)
	if err := fn(atomic.Ctx()); err != nil {
		atomic.Rollback()
		return err
	}
	atomic.Commit()
	return nil
}
