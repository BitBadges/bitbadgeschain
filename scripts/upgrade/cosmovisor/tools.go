//go:build tools

// Package tools pins cosmovisor as a build dependency of this wrapper module.
//
// cosmovisor v1.7.1's own go.mod resolves github.com/bytedance/sonic to a
// version that does not link on Go 1.26 ("invalid reference to
// encoding/json.unquoteBytes"). Building it through this module lets go.mod
// raise sonic to a Go-1.26-compatible release, so the chain binaries and
// cosmovisor can all be built with the toolchain the chain's go.mod requires.
package tools

import _ "cosmossdk.io/tools/cosmovisor/cmd/cosmovisor"
