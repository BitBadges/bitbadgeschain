package types

func IdRangeEquals(a, b []*IdRange) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, aIdRange := range a {
		bIdRange := b[idx]
		if !aIdRange.Start.Equal(bIdRange.Start) {
			return false
		}

		if !aIdRange.End.Equal(bIdRange.End) {
			return false
		}
	}

	return true
}

func BalancesEqual(a, b []*Balance) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, aBalance := range a {
		bBalance := b[idx]
		if !aBalance.Amount.Equal(bBalance.Amount) {
			return false
		}

		if !IdRangeEquals(aBalance.BadgeIds, bBalance.BadgeIds) {
			return false
		}

		if !IdRangeEquals(aBalance.Times, bBalance.Times) {
			return false
		}
	}

	return true

}
