package keeper

var (
	// ManagerSplitterKey is the prefix for manager splitter storage
	ManagerSplitterKey = []byte{0x01}
	// ManagerSplitterCountKey is the prefix for the next manager splitter ID
	ManagerSplitterCountKey = []byte{0x02}
)

// managerSplitterStoreKey returns the key for a manager splitter by address
func managerSplitterStoreKey(addr string) []byte {
	key := make([]byte, len(ManagerSplitterKey)+len(addr))
	copy(key, ManagerSplitterKey)
	copy(key[len(ManagerSplitterKey):], []byte(addr))
	return key
}

