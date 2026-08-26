package app

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// buildAndDropApps builds n test apps, drops each one, and reports the change
// in live goroutines and live heap across the whole batch.
//
// The GC is run twice because the first cycle can queue finalizers whose work
// only lands on the second, and a leak measured after one cycle reads high for
// reasons that are not leaks.
func buildAndDropApps(t *testing.T, n int) (goroutineGrowth int, heapGrowthBytes int64) {
	t.Helper()

	var ms runtime.MemStats
	settle := func() (int, int64) {
		runtime.GC()
		runtime.GC()
		runtime.ReadMemStats(&ms)
		return runtime.NumGoroutine(), int64(ms.HeapAlloc)
	}

	baseGoroutines, baseHeap := settle()

	for i := 0; i < n; i++ {
		app := Setup(false)
		// Touch the app the way a test would, so this measures the state a real
		// test leaves behind rather than a bare constructor.
		_ = app.BaseApp.NewContext(false)
	}

	afterGoroutines, afterHeap := settle()
	return afterGoroutines - baseGoroutines, afterHeap - baseHeap
}

// TestTestAppIsNotRetainedAfterUse pins that building a test app and dropping it
// actually releases it.
//
// It did not. Every Setup() leaked 36 goroutines - 30 iavl nodeDB pruners, one
// per store key, plus 6 belonging to the cosmos/evm mempool - and each of those
// goroutines pinned its app, the keepers behind it, and the merged protobuf file
// descriptor set built during construction. x/tokenization/keeper builds one app
// per test case in SetupTest, so its test binary retained all 320 of them and
// peaked at 12.88 GB RSS on a 16 GB CI runner. The job died with a bare exit 143
// and no message.
//
// Nothing in the suite could see this. Every test passed at every point; the
// only symptom was the runner being killed, which reads as CI flake rather than
// as a defect in the code. Hence this test: the property has no other observable
// consequence until the day it takes a runner down.
//
// The thresholds are deliberately loose. This test is not measuring an exact
// footprint, it is separating "released" from "retained per app", and those two
// are two orders of magnitude apart.
func TestTestAppIsNotRetainedAfterUse(t *testing.T) {
	const apps = 10

	goroutineGrowth, heapGrowth := buildAndDropApps(t, apps)

	// Pre-fix: +360. A few goroutines of slack absorbs lazily-started package
	// internals (the desertbit timer routine, for one) that start once and are
	// not per-app.
	require.LessOrEqual(t, goroutineGrowth, 5,
		"building %d test apps leaked %d goroutines; a goroutine that outlives its app "+
			"pins the app, its keepers and its descriptor set, and the leak is linear in "+
			"the number of tests in the package", apps, goroutineGrowth)

	// Pre-fix: ~6.8 MB an app, so ~68 MB here. Steady state after the fix is
	// under 10 MB for the whole batch.
	const maxHeapGrowth = 20 << 20
	require.Less(t, heapGrowth, int64(maxHeapGrowth),
		"building %d test apps retained %.1f MB of heap; test apps must be collectable "+
			"once dropped", apps, float64(heapGrowth)/(1<<20))
}
