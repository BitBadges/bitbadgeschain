package evmcompat

import (
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// RunNativeAction delegates to cosmos/evm's RunNativeAction.
func RunNativeAction(p cmn.Precompile, evm *vm.EVM, contract *vm.Contract, action cmn.NativeAction) ([]byte, error) {
	return p.RunNativeAction(evm, contract, action)
}
