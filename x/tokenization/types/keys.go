package types

const (
	// ModuleName defines the module name
	ModuleName = "tokenization"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tokenization"

	// Version defines the current version the IBC module supports
	Version = "tokenization-1"

	// PortID is the default port id that module binds to
	PortID = "tokenization"

	// ICQ (Interchain Query) constants for cross-chain token ownership verification
	// ICQVersion defines the version for ICQ channels
	ICQVersion = "icq-tokenization-1"

	// ICQMaxProofSize is the maximum size of balance proofs in bytes
	ICQMaxProofSize = 10240 // 10KB
)

var ParamsKey = []byte("p_tokenization")

// PortKey defines the key to store the port ID in store
var PortKey = KeyPrefix("tokenization-port-")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
