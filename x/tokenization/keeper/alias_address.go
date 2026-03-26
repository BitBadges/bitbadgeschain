package keeper

import "fmt"

// GenerateWrapperPathAddress derives the deterministic wrapper path address from a denom.
// Exposed for tests and any logic that needs the address without storing it on-chain.
func GenerateWrapperPathAddress(denom string) (string, error) {
	accountAddr, err := generatePathAddress(denom, WrapperPathGenerationPrefix)
	if err != nil {
		return "", err
	}
	return accountAddr.String(), nil
}

// MustGenerateWrapperPathAddress is a convenience helper that panics on error (suitable for tests).
func MustGenerateWrapperPathAddress(denom string) string {
	addr, err := GenerateWrapperPathAddress(denom)
	if err != nil {
		panic(fmt.Errorf("failed to generate wrapper path address: %w", err))
	}
	return addr
}
