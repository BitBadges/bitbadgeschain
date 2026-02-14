package keeper

import "fmt"

// GenerateAliasPathAddress derives the deterministic alias path address from a denom.
// Exposed for tests and any logic that needs the address without storing it on-chain.
func GenerateAliasPathAddress(denom string) (string, error) {
	accountAddr, err := generatePathAddress(denom, WrapperPathGenerationPrefix)
	if err != nil {
		return "", err
	}
	return accountAddr.String(), nil
}

// MustGenerateAliasPathAddress is a convenience helper that panics on error (suitable for tests).
func MustGenerateAliasPathAddress(denom string) string {
	addr, err := GenerateAliasPathAddress(denom)
	if err != nil {
		panic(fmt.Errorf("failed to generate alias path address: %w", err))
	}
	return addr
}
