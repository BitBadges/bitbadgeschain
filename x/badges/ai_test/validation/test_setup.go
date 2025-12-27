package validation

import (
	"sync"

	"github.com/bitbadges/bitbadgeschain/app/params"
)

var (
	configInitOnce sync.Once
)

// EnsureSDKConfig ensures the SDK config is initialized with "bb" prefix
// This should be called at the start of tests that need address validation
func EnsureSDKConfig() {
	configInitOnce.Do(func() {
		// Try to initialize SDK config with "bb" prefix
		config := params.InitSDKConfigWithoutSeal()
		
		// Force set the prefix if it's not already "bb"
		// This might panic if config is sealed, but we'll catch it
		defer func() {
			if r := recover(); r != nil {
				// Config might be sealed, that's ok - we'll use whatever prefix is set
			}
		}()
		
		currentPrefix := config.GetBech32AccountAddrPrefix()
		if currentPrefix != "bb" {
			// Try to set it - might panic if sealed
			config.SetBech32PrefixForAccount("bb", "bbpub")
			config.SetBech32PrefixForValidator("bbvaloper", "bbvaloperpub")
			config.SetBech32PrefixForConsensusNode("bbvalcons", "bbvalconspub")
		}
	})
}

// init initializes the SDK config with "bb" prefix for validation tests
func init() {
	// Ensure SDK config is initialized with "bb" prefix before any address validation
	// This must be called before ValidateAddress is used
	EnsureSDKConfig()
}

