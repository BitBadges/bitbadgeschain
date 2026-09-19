package lifecycle

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func validateAssertion(a Assertion) error {
	if len(a.Expected) == 0 || string(a.Expected) == "null" {
		return fmt.Errorf("assertion needs explicit non-null expected value")
	}
	switch a.Kind {
	case "coins", "coinDelta":
		if err := sdk.ValidateDenom(a.Denom); err != nil {
			return err
		}
	case "balances", "balanceAt", "approvals", "permissions", "tracker", "manager":
		if !validUint(a.CollectionID) || a.CollectionID == "0" {
			return fmt.Errorf("invalid collectionId")
		}
	default:
		return fmt.Errorf("unsupported assertion kind %q", a.Kind)
	}
	switch a.Kind {
	case "coins", "coinDelta", "balances", "balanceAt":
		if _, err := sdk.AccAddressFromBech32(a.Address); err != nil {
			return err
		}
	}
	if a.Kind == "balanceAt" {
		if !validUint(a.TokenID) || a.TokenID == "0" || (a.OwnershipTime != "now" && !validUint(a.OwnershipTime)) {
			return fmt.Errorf("invalid balance coordinate")
		}
	}
	if a.Kind == "tracker" && (a.ApprovalID == "" || a.AmountTrackerID == "" || a.Level == "" || a.TrackerType == "") {
		return fmt.Errorf("tracker needs explicit identifiers")
	}
	if a.Kind == "balances" {
		var balances []*types.Balance
		if err := json.Unmarshal(a.Expected, &balances); err != nil {
			return err
		}
		return validBalances(balances)
	}
	return nil
}
func observe(application *app.App, ctx sdk.Context, a Assertion, before sdkmath.Int) (json.RawMessage, error) {
	marshal := func(v any) (json.RawMessage, error) { b, e := json.Marshal(v); return b, e }
	if a.Kind == "coins" || a.Kind == "coinDelta" {
		addr, _ := sdk.AccAddressFromBech32(a.Address)
		n := application.BankKeeper.GetBalance(ctx, addr, a.Denom).Amount
		if a.Kind == "coinDelta" {
			n = n.Sub(before)
		}
		return marshal(n.String())
	}
	id := sdkmath.NewUintFromString(a.CollectionID)
	collection, found := application.TokenizationKeeper.GetCollectionFromStore(ctx, id)
	if !found {
		return nil, fmt.Errorf("collection %s not found", a.CollectionID)
	}
	switch a.Kind {
	case "manager":
		return marshal(collection.Manager)
	case "approvals":
		return marshal(collection.CollectionApprovals)
	case "permissions":
		return marshal(collection.CollectionPermissions)
	case "tracker":
		tracker, found := application.TokenizationKeeper.GetApprovalTrackerFromStore(ctx, id, a.AddressForApproval, a.ApprovalID, a.AmountTrackerID, a.Level, a.TrackerType, a.Address)
		if !found {
			return nil, fmt.Errorf("tracker not found")
		}
		return marshal(tracker)
	case "balances", "balanceAt":
		balances, _, err := application.TokenizationKeeper.GetBalanceOrApplyDefault(ctx, collection, a.Address)
		if err != nil {
			return nil, err
		}
		if a.Kind == "balances" {
			if len(balances.Balances) == 0 {
				return json.RawMessage("[]"), nil
			}
			return marshal(balances.Balances)
		}
		token, _ := new(big.Int).SetString(a.TokenID, 10)
		ownership := a.OwnershipTime
		if ownership == "now" {
			ownership = fmt.Sprint(ctx.BlockTime().UnixMilli())
		}
		tm, _ := new(big.Int).SetString(ownership, 10)
		return marshal(amountAt(balances.Balances, token, tm).String())
	}
	return nil, fmt.Errorf("unsupported assertion")
}
func compare(a Assertion, actual json.RawMessage) (bool, error) {
	if a.Kind == "balances" {
		var want, got []*types.Balance
		if err := json.Unmarshal(a.Expected, &want); err != nil {
			return false, err
		}
		if err := json.Unmarshal(actual, &got); err != nil {
			return false, err
		}
		return equalBalances(want, got)
	}
	var want, got any
	if err := json.Unmarshal(a.Expected, &want); err != nil {
		return false, err
	}
	if err := json.Unmarshal(actual, &got); err != nil {
		return false, err
	}
	return reflect.DeepEqual(want, got), nil
}
func validBalances(bs []*types.Balance) error {
	for _, b := range bs {
		if b == nil || b.Amount.IsNil() || len(b.TokenIds) == 0 || len(b.OwnershipTimes) == 0 {
			return fmt.Errorf("invalid balance")
		}
		for _, rs := range [][]*types.UintRange{b.TokenIds, b.OwnershipTimes} {
			for i, r := range rs {
				if r == nil || r.Start.IsNil() || r.End.IsNil() || r.Start.GT(r.End) {
					return fmt.Errorf("invalid range")
				}
				for _, prev := range rs[:i] {
					if !r.Start.GT(prev.End) && !prev.Start.GT(r.End) {
						return fmt.Errorf("overlapping ranges inside one balance")
					}
				}
			}
		}
	}
	return nil
}
func amountAt(bs []*types.Balance, id, tm *big.Int) *big.Int {
	n := new(big.Int)
	for _, b := range bs {
		if contains(b.TokenIds, id) && contains(b.OwnershipTimes, tm) {
			n.Add(n, b.Amount.BigInt())
		}
	}
	return n
}
func contains(rs []*types.UintRange, n *big.Int) bool {
	for _, r := range rs {
		if n.Cmp(r.Start.BigInt()) >= 0 && n.Cmp(r.End.BigInt()) <= 0 {
			return true
		}
	}
	return false
}
func equalBalances(a, b []*types.Balance) (bool, error) {
	if err := validBalances(a); err != nil {
		return false, err
	}
	if err := validBalances(b); err != nil {
		return false, err
	}
	ids, times := map[string]*big.Int{}, map[string]*big.Int{}
	for _, bs := range [][]*types.Balance{a, b} {
		for _, bal := range bs {
			for i, rs := range [][]*types.UintRange{bal.TokenIds, bal.OwnershipTimes} {
				target := ids
				if i == 1 {
					target = times
				}
				for _, r := range rs {
					start := r.Start.BigInt()
					end := new(big.Int).Add(r.End.BigInt(), big.NewInt(1))
					target[start.String()] = start
					target[end.String()] = end
				}
			}
		}
	}
	if len(ids)*len(times) > 1000000 {
		return false, fmt.Errorf("balance assertion grid exceeds limit")
	}
	keys := func(m map[string]*big.Int) []*big.Int {
		out := []*big.Int{}
		for _, v := range m {
			out = append(out, v)
		}
		sort.Slice(out, func(i, j int) bool { return out[i].Cmp(out[j]) < 0 })
		return out
	}
	for _, id := range keys(ids) {
		for _, tm := range keys(times) {
			if amountAt(a, id, tm).Cmp(amountAt(b, id, tm)) != 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

var uintPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)

func validUint(s string) bool {
	n, ok := new(big.Int).SetString(s, 10)
	return ok && uintPattern.MatchString(s) && n.BitLen() <= 256
}
