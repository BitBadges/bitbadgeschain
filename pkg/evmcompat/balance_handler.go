package evmcompat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/ethereum/go-ethereum/common"
)

// NewBalanceHandlerFactory reconciles only accounts that have an EVM address.
func NewBalanceHandlerFactory(bankKeeper cmn.BankKeeper) *cmn.BalanceHandlerFactory {
	return cmn.NewBalanceHandlerFactory(evmBalanceAccounts{BankKeeper: bankKeeper})
}

type evmBalanceAccounts struct {
	cmn.BankKeeper
}

func (k evmBalanceAccounts) BlockedAddr(addr sdk.AccAddress) bool {
	// Other address lengths retain their native store balances without EVM conversion.
	return len(addr) != common.AddressLength || k.BankKeeper.BlockedAddr(addr)
}
