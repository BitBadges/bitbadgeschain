package keeper

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// ValidateEVMQueryChallengesAreContracts ensures every EVM query challenge's contract address
// has code on the EVM (is a contract). Call this when storing collections or approvals that
// include EVM query challenges. When evmKeeper is nil, validation is skipped (e.g. in tests
// or chains without EVM); execution will still fail at runtime if the address has no code.
func (k Keeper) ValidateEVMQueryChallengesAreContracts(ctx sdk.Context, challenges []*types.EVMQueryChallenge) error {
	if len(challenges) == 0 {
		return nil
	}
	if k.evmKeeper == nil {
		return nil // Skip when EVM not available; execution-time check will fail if needed
	}
	for i, challenge := range challenges {
		if challenge == nil || challenge.ContractAddress == "" {
			continue
		}
		var contractAddr common.Address
		addr := strings.TrimSpace(challenge.ContractAddress)
		if len(addr) >= 2 && strings.ToLower(addr[:2]) == "0x" {
			contractAddr = common.HexToAddress(addr)
		} else {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return fmt.Errorf("EVM query challenge %d: invalid contract address: %w", i, err)
			}
			contractAddr = common.BytesToAddress(accAddr.Bytes())
		}
		if !k.evmKeeper.IsContract(ctx, contractAddr) {
			return fmt.Errorf("EVM query challenge %d: address is not a contract (no code): %s", i, challenge.ContractAddress)
		}
	}
	return nil
}

// ExecuteEVMQuery performs a read-only EVM contract call (uses zero address as caller).
// Use ExecuteEVMQueryWithCaller when the caller must be an existing account (e.g. for approval/invariant flows).
func (k Keeper) ExecuteEVMQuery(ctx sdk.Context, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	return k.ExecuteEVMQueryWithCaller(ctx, "", contractAddress, calldata, gasLimit)
}

// ExecuteEVMQueryWithCaller performs a read-only EVM contract call with the given caller.
// If callerAddress is empty, the zero address is used.
// Note: cosmos/evm requires the caller account to exist in auth keeper for GetSequence().
// For read-only queries, we ensure the zero address account exists temporarily if needed.
func (k Keeper) ExecuteEVMQueryWithCaller(ctx sdk.Context, callerAddress string, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	if k.evmKeeper == nil {
		return nil, fmt.Errorf("EVM keeper not available")
	}

	// Convert contract address to common.Address
	var contractAddr common.Address
	if len(contractAddress) >= 2 && strings.ToLower(contractAddress[:2]) == "0x" {
		contractAddr = common.HexToAddress(contractAddress)
	} else {
		accAddr, err := sdk.AccAddressFromBech32(contractAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid contract address: %w", err)
		}
		contractAddr = common.BytesToAddress(accAddr.Bytes())
	}

	// Verify it's a contract
	if !k.evmKeeper.IsContract(ctx, contractAddr) {
		return nil, fmt.Errorf("address is not a contract: %s", contractAddress)
	}

	// Resolve caller: use provided address, otherwise zero address
	var callerAddr common.Address
	if callerAddress != "" {
		if len(callerAddress) >= 2 && strings.ToLower(callerAddress[:2]) == "0x" {
			callerAddr = common.HexToAddress(callerAddress)
		} else {
			accAddr, err := sdk.AccAddressFromBech32(callerAddress)
			if err != nil {
				return nil, fmt.Errorf("invalid caller address: %w", err)
			}
			callerAddr = common.BytesToAddress(accAddr.Bytes())
		}
	}
	// If callerAddress is empty, callerAddr will be the zero address (default value)

	// cosmos/evm's CallEVMWithData requires the caller account to exist for GetSequence().
	// For read-only queries we use the zero address as caller; ensure that account exists.
	// This write persists to chain state (once per chain). It is intentional so that staticcall-style
	// EVM queries work without requiring a real signer account.
	callerAccAddr := sdk.AccAddress(callerAddr.Bytes())
	if k.accountKeeper != nil && k.accountKeeper.GetAccount(ctx, callerAccAddr) == nil {
		acc := k.accountKeeper.NewAccountWithAddress(ctx, callerAccAddr)
		k.accountKeeper.SetAccount(ctx, acc)
	}

	gasCap := new(big.Int).SetUint64(gasLimit)
	response, err := k.evmKeeper.CallEVMWithData(ctx, callerAddr, &contractAddr, calldata, false, gasCap)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %w", err)
	}

	if response.VmError != "" {
		return nil, fmt.Errorf("EVM execution error: %s", response.VmError)
	}

	return response.Ret, nil
}
