package keeper

import (
	"strings"

	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

var (
	ProposalKey       = []byte{0x01}
	NextProposalIdKey = []byte{0x02}
	Delimiter         = []byte{0xDD}
	Placeholder       = []byte{0xFF}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

func ConstructProposalKey(id sdkmath.Uint) string {
	keyParts := []string{
		id.String(),
	}
	return strings.Join(keyParts, BalanceKeyDelimiter)
}

func proposalKey(id sdkmath.Uint) []byte {
	key := make([]byte, len(ProposalKey)+IDLength)
	copy(key, ProposalKey)
	copy(key[len(ProposalKey):], []byte(id.String()))
	return key
}

func nextProposalIdKey() []byte {
	return NextProposalIdKey
}
