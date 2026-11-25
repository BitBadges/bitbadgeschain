package types

import (
	"encoding/binary"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
)

const (
	ModuleName = "customhooks"
	RouterKey  = ModuleName
	StoreKey   = ModuleName
)

var (
	// TransientStoreKey is the transient store key for storing deterministic error messages
	// Transient stores are automatically cleared at the end of each transaction
	TransientStoreKey = storetypes.NewTransientStoreKey("customhooks_transient")

	// DeterministicErrorKey is the key used to store deterministic error messages in the transient store
	DeterministicErrorKey = []byte("deterministic_error")

	// PendingIBCRecoveryPrefix namespaces outbound IBC recovery records by (port, channel, sequence).
	PendingIBCRecoveryPrefix = []byte("customhooks_ibc_recovery/")
)

// PendingIBCRecoveryKey builds a deterministic store key for an outbound IBC packet recovery entry.
func PendingIBCRecoveryKey(portID, channelID string, sequence uint64) []byte {
	key := make([]byte, 0, len(PendingIBCRecoveryPrefix)+len(portID)+1+len(channelID)+1+8)
	key = append(key, PendingIBCRecoveryPrefix...)
	key = append(key, []byte(portID)...)
	key = append(key, '/')
	key = append(key, []byte(channelID)...)
	key = append(key, '/')
	seqBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBytes, sequence)
	key = append(key, seqBytes...)
	return key
}
