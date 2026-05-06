// Helper for distinguishing "field unset" from "field set to zero value"
// after the chain normalizes nil pointers on load.
//
// Background: the proto schema marks nested struct fields
// `(gogoproto.nullable) = true` so the canonical Amino JSON omits
// them when the SDK's `dropEmptyProtoSubMessages` strips empty
// sub-messages from the proto wire. This is required for EIP-712
// typed-data parity with what MetaMask signs.
//
// To prevent nil-deref panics across ~170 keeper access sites,
// `keeper.NormalizeNilPointers` fills nil pointer-to-struct fields
// with fresh zero-value instances after storage load. After
// normalization, `field != nil` no longer reliably distinguishes
// "user set this field" from "user did not set this field" — both
// produce a non-nil pointer (one points to a zero-value struct,
// the other to a populated one).
//
// `IsBasicallyEmpty` recovers that distinction by checking the
// proto-encoded size: a populated message has wire size > 0, a
// zero-value message has size 0. Use this at "is field set"
// semantic check sites; keep `!= nil` for plain nil-safety guards
// (those continue to work after normalize).

package types

import (
	"reflect"
	"strings"
)

// IsBasicallyEmpty reports whether `m` is conceptually unset.
//
// Returns true when:
//   - `m` is nil (interface or typed nil pointer)
//   - all of `m`'s fields are zero values
//   - all message-typed sub-fields are themselves IsBasicallyEmpty
//     (recursive — a non-nil pointer to a zero-value struct counts as
//     empty, including when its sub-fields are also non-nil pointers to
//     zero-value structs)
//
// Why not just `m.Size() == 0`: gogoproto's Size() includes the wire
// tag + length-0 marker (typically 2 bytes) for every non-nil
// embedded-message field, even when that embedded message is itself
// all-zero. After `keeper.NormalizeNilPointers` walks a fresh struct,
// every sub-field is a non-nil pointer to a zero struct, so Size
// returns nonzero — falsely reporting "set" — which over-fires
// security checks like the cosmosCoinBackedPath / Mint-approval
// rule.
//
// Reflection walk treats normalize-filled empty pointer chains as
// empty so callers get the original "user actually configured this"
// semantic back.
func IsBasicallyEmpty(m interface{}) bool {
	if m == nil {
		return true
	}
	rv := reflect.ValueOf(m)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		// scalar pointed-to value: zero <=> empty
		return rv.IsZero()
	}
	return isStructAllZero(rv, 0)
}

func isStructAllZero(rv reflect.Value, depth int) bool {
	if depth > 50 {
		// safety: assume non-empty rather than recurse infinitely
		return false
	}
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		name := t.Field(i).Name
		// Skip protobuf bookkeeping fields (XXX_unrecognized, XXX_sizecache, etc).
		if strings.HasPrefix(name, "XXX_") {
			continue
		}
		if !isFieldZero(f, depth+1) {
			return false
		}
	}
	return true
}

func isFieldZero(f reflect.Value, depth int) bool {
	switch f.Kind() {
	case reflect.Ptr:
		if f.IsNil() {
			return true
		}
		// Non-nil pointer to a struct: empty if the pointee is all-zero
		// (recursive). This is the load-bearing case — it makes
		// `&Conversion{}` (normalize-filled) report as empty even though
		// the parent `&CosmosCoinBackedPath{Conversion: &Conversion{}}`
		// has non-zero proto.Size().
		if f.Type().Elem().Kind() == reflect.Struct {
			return isStructAllZero(f.Elem(), depth)
		}
		return f.IsZero()
	case reflect.Struct:
		return isStructAllZero(f, depth)
	case reflect.Slice, reflect.Map:
		return f.Len() == 0
	default:
		return f.IsZero()
	}
}
