// Nil-pointer normalization for proto messages loaded from storage.
//
// Background: nullable=true on nested struct fields is required so the
// chain's `MarshalAminoJSON` omits empty sub-structs (otherwise empty
// types in the typed-data tree trigger go-ethereum's encodeType bug
// and break EIP-712 verification — see codec.go and the proto changes).
// But it means stored messages may have `nil` pointers where the older
// schema always had non-nil zero-value structs.
//
// Existing keeper code (~170+ sites) reads chains like
// `approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId`
// without nil-guarding every level — those checks were unnecessary
// before nullable=true. Rather than thread nil checks through every
// access site, this normalizer walks a freshly-unmarshaled message and
// replaces any nil pointer-to-struct field with a fresh empty instance.
// Downstream code keeps reading direct field paths and gets the
// zero-value semantics it would have gotten under the old non-nullable
// schema.
//
// Apply at storage-load boundaries (GetCollectionFromStore, etc.) and
// at msg-handler entry points. The empty struct values default to
// fields that mean "no constraint" / "default behavior" — same as the
// pre-nullable=true world, so semantics are preserved.

package keeper

import (
	"reflect"
	"strings"

	sdkmath "cosmossdk.io/math"
)

// NormalizeNilPointers recursively walks `v` (must be a pointer) and
// initializes any uninitialized `sdkmath.Uint` / `sdkmath.Int` field
// with `NewUint(0)` / `NewInt(0)`. The Go zero value of those types
// holds a nil internal `*big.Int`, so any arithmetic op (.GT, .Cmp,
// .Uint64, etc.) panics. Proto fields tagged `customtype = "Uint",
// nullable = false` deserialize to the broken zero value when the
// wire bytes don't include the field — common after the SDK side
// strips empty Uint customtype strings before broadcasting.
//
// This intentionally DOES NOT fill nil pointer-to-struct fields with
// fresh empty structs. Earlier versions did, but that broke the
// nil-vs-empty distinction many code paths rely on (equality checks,
// "user explicitly set this" semantic checks, change-detection in
// approval-update tracking, etc.). Callers that need direct field
// access on potentially-nil pointers must use explicit nil guards or
// the proto-generated `GetX()` accessors.
func NormalizeNilPointers(v interface{}) {
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	walkStructForNormalize(rv.Elem(), 0)
}

// uintZero / intZero are pre-built non-nil zero values used to replace
// uninitialized sdkmath.Uint{} / sdkmath.Int{} fields. Those types wrap
// `*big.Int`; their Go zero value has a nil internal pointer, so any
// `.GT() / .Cmp() / .Uint64()` etc. panics. Proto fields with
// `customtype = "Uint"` and `nullable = false` deserialize to the zero
// value when the wire bytes don't include the field — typical when the
// SDK side dropped the empty string. Normalize-time replacement
// guarantees downstream math operations never see nil internals.
var (
	uintZero = sdkmath.NewUint(0)
	intZero  = sdkmath.NewInt(0)

	uintType = reflect.TypeOf(uintZero)
	intType  = reflect.TypeOf(intZero)
)

func walkStructForNormalize(v reflect.Value, depth int) {
	if depth > 50 {
		return // safety: bounded recursion
	}
	if v.Kind() != reflect.Struct {
		return
	}

	// sdkmath.Uint / sdkmath.Int are struct values (non-pointer), but
	// their Go zero state is *invalid* — internal *big.Int is nil and any
	// arithmetic operation panics. Detect and replace with a properly-
	// initialized zero. Must run BEFORE the field walk because the
	// recursive call would skip these (Int/Uint internals aren't fields
	// we want to recurse into).
	if v.Type() == uintType {
		if (sdkmath.Uint{}) == v.Interface().(sdkmath.Uint) {
			if v.CanSet() {
				v.Set(reflect.ValueOf(uintZero))
			}
		}
		return
	}
	if v.Type() == intType {
		if (sdkmath.Int{}) == v.Interface().(sdkmath.Int) {
			if v.CanSet() {
				v.Set(reflect.ValueOf(intZero))
			}
		}
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ft := t.Field(i)
		if !field.CanSet() {
			continue
		}
		// Skip protobuf bookkeeping fields
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			// Only recurse into ALREADY-non-nil pointers — don't fill
			// nil pointers (would break nil-vs-empty equality semantics).
			if !field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
				walkStructForNormalize(field.Elem(), depth+1)
			}
		case reflect.Struct:
			walkStructForNormalize(field, depth+1)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				switch elem.Kind() {
				case reflect.Ptr:
					if !elem.IsNil() && elem.Type().Elem().Kind() == reflect.Struct {
						walkStructForNormalize(elem.Elem(), depth+1)
					}
				case reflect.Struct:
					walkStructForNormalize(elem, depth+1)
				}
			}
		}
	}
}
