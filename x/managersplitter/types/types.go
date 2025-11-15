package types

const (
	// ModuleName defines the module name
	ModuleName = "managersplitter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_managersplitter"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte("p_managersplitter")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

