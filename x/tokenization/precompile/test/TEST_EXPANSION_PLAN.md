# Test Expansion Plan

## Overview

This document outlines a comprehensive plan to expand test coverage for the tokenization precompile, including full EVM keeper integration tests and Solidity contract tests.

## Current Test Coverage Analysis

### Existing Test Files ✅
- `precompile_test.go` - Basic precompile structure tests
- `handlers_test.go` - Handler method unit tests  
- `conversions_test.go` - Type conversion tests
- `validation_test.go` - Input validation tests
- `error_test.go` - Error handling tests
- `security_test.go` - Security pattern tests
- `query_test.go` - Query method tests
- `query_methods_test.go` - Additional query tests
- `e2e_test.go` - Basic E2E tests (direct precompile calls)
- `integration_e2e_test.go` - Integration E2E tests
- `edge_cases_test.go` - Edge case tests
- `gas_test.go` - Gas calculation tests
- `gas_benchmark_test.go` - Gas benchmarks
- `return_types_test.go` - Return type conversion tests
- `erc3643_compliance_test.go` - ERC3643 compliance tests
- `evm_integration_test.go` - ⚠️ Placeholder (needs full implementation)

### Coverage Gaps Identified

1. **Full EVM Keeper Integration** - Tests using actual EVM keeper (not just direct calls)
2. **Solidity Contract Tests** - Deploy and test actual Solidity contracts
3. **Comprehensive ApprovalCriteria Tests** - All approval criteria features
4. **Concurrency/Race Condition Tests** - Parallel execution safety
5. **Reentrancy Tests** - Reentrancy attack scenarios
6. **Event Verification Through EVM** - Event emission verification
7. **Gas Accuracy Verification** - Verify estimates match actual consumption
8. **Stress Tests** - Large-scale operations
9. **Multi-User Scenarios** - Complex multi-user workflows
10. **Error Recovery Tests** - Error handling and recovery paths

---

## Phase 1: Full EVM Keeper Integration Tests

### Goal
Create tests that use the actual EVM keeper to execute precompile calls, verifying the full integration path from EVM transaction → precompile → keeper → state.

### Implementation Plan

#### 1.1 Create EVM Keeper Test Suite
**File:** `evm_keeper_integration_test.go`

**Structure:**
```go
type EVMKeeperIntegrationTestSuite struct {
    suite.Suite
    App            *app.App
    Ctx            sdk.Context
    EVMKeeper      *evmkeeper.Keeper
    Precompile     *Precompile
    
    // Test accounts with private keys
    AliceKey       *ecdsa.PrivateKey
    BobKey         *ecdsa.PrivateKey
    AliceEVM       common.Address
    BobEVM         common.Address
    Alice          sdk.AccAddress
    Bob            sdk.AccAddress
    
    CollectionId   sdkmath.Uint
}
```

**Setup:**
- Use `app.Setup()` to create full app instance
- Initialize EVM keeper with precompile registered (see `app/evm.go:137`)
- Create EVM instance for transaction execution
- Fund test accounts with native tokens

#### 1.2 Test Cases

**TestEVMKeeper_PrecompileRegistration**
```go
// Verify precompile is registered in EVM keeper
// Check precompile address is correct
// Verify precompile can be called via EVM
```

**TestEVMKeeper_TransferTokens_ThroughEVM**
```go
// Create EVM transaction calling precompile.transferTokens()
// Execute through EVM keeper
// Verify state changes in tokenization keeper
// Check events emitted through EVM
```

**TestEVMKeeper_QueryMethods_ThroughEVM**
```go
// Call query methods via EVM (static calls)
// Verify return values
// Check gas consumption
```

**TestEVMKeeper_GasAccounting**
```go
// Verify gas is properly deducted
// Check gas refunds (if applicable)
// Verify gas limits are enforced
```

**TestEVMKeeper_ErrorHandling**
```go
// Test error propagation through EVM
// Verify error codes are preserved
// Check error messages are sanitized
```

**TestEVMKeeper_AllTransactionMethods**
```go
// Test all 24+ transaction methods through EVM
// Verify each method works correctly
// Check state changes
```

#### 1.3 Helper Functions

```go
// setupAppWithEVM creates full app with EVM keeper and precompile registered
func setupAppWithEVM(t *testing.T) (*app.App, sdk.Context, *evmkeeper.Keeper, *Precompile)

// createEVMTransaction creates an EVM transaction calling the precompile
func (suite *EVMKeeperIntegrationTestSuite) createEVMTransaction(
    fromKey *ecdsa.PrivateKey,
    methodName string,
    args []interface{},
) (*types.Transaction, error)

// executeEVMTransaction executes a transaction through EVM keeper
func (suite *EVMKeeperIntegrationTestSuite) executeEVMTransaction(
    tx *types.Transaction,
) (*evmtypes.MsgEthereumTxResponse, error)

// verifyEvents checks that expected events were emitted
func (suite *EVMKeeperIntegrationTestSuite) verifyEvents(
    events []sdk.Event,
    expectedEvents []string,
) error

// fundAccount funds an account with native tokens for gas
func (suite *EVMKeeperIntegrationTestSuite) fundAccount(
    addr sdk.AccAddress,
    amount sdk.Coins,
) error
```

---

## Phase 2: Solidity Contract Tests

### Goal
Deploy and test actual Solidity contracts that interact with the precompile, simulating real-world usage.

### Implementation Plan

#### 2.1 Create Solidity Test Contract
**File:** `contracts/test/PrecompileTestContract.sol`

**Features:**
- Wrapper functions for all precompile methods
- Event emission for test verification
- Error handling and recovery
- Gas estimation helpers

**Example:**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";

contract PrecompileTestContract {
    ITokenizationPrecompile constant precompile = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    event TransferExecuted(uint256 collectionId, address recipient, bool success);
    
    function testTransfer(
        uint256 collectionId,
        address recipient,
        uint256 amount
    ) external returns (bool) {
        address[] memory recipients = new address[](1);
        recipients[0] = recipient;
        // ... setup ranges
        
        bool success = precompile.transferTokens(...);
        emit TransferExecuted(collectionId, recipient, success);
        return success;
    }
}
```

#### 2.2 Create Solidity Test Suite
**File:** `solidity_contract_test.go`

**Structure:**
```go
type SolidityContractTestSuite struct {
    suite.Suite
    App            *app.App
    Ctx            sdk.Context
    EVMKeeper      *evmkeeper.Keeper
    
    // Deployed contract address and ABI
    TestContractAddr common.Address
    TestContractABI  abi.ABI
    ContractBytecode []byte
    
    // Test accounts
    DeployerKey    *ecdsa.PrivateKey
    AliceKey       *ecdsa.PrivateKey
    BobKey         *ecdsa.PrivateKey
    
    CollectionId   sdkmath.Uint
}
```

#### 2.3 Test Cases

**TestSolidity_DeployContract**
```go
// Deploy test contract
// Verify deployment succeeds
// Check contract code is stored
// Verify contract can call precompile
```

**TestSolidity_TransferTokens_ThroughContract**
```go
// Call precompile through Solidity contract
// Verify transfer succeeds
// Check events from both contract and precompile
// Verify state changes
```

**TestSolidity_QueryMethods_ThroughContract**
```go
// Call query methods through contract (view functions)
// Verify return values are correct
// Check gas consumption
```

**TestSolidity_ErrorHandling**
```go
// Test error propagation from precompile to contract
// Verify contract handles errors correctly
// Check error messages
```

**TestSolidity_ReentrancyProtection**
```go
// Attempt reentrancy attack through contract
// Verify protection works
// Check state consistency
```

**TestSolidity_ComplexWorkflows**
```go
// Multi-step workflows through contract
// State changes across multiple calls
// Event verification
```

**TestSolidity_AllPrecompileMethods**
```go
// Test all precompile methods through contract
// Verify each method works
// Check return values
```

#### 2.4 Helper Functions

```go
// compileSolidityContract compiles a Solidity contract
func compileSolidityContract(contractPath string) ([]byte, *abi.ABI, error)

// deployContract deploys a Solidity contract
func (suite *SolidityContractTestSuite) deployContract(
    contractCode []byte,
    constructorArgs []interface{},
) (common.Address, error)

// callContract calls a contract method
func (suite *SolidityContractTestSuite) callContract(
    contractAddr common.Address,
    methodName string,
    args []interface{},
) ([]interface{}, error)

// verifyContractEvents checks contract events
func (suite *SolidityContractTestSuite) verifyContractEvents(
    events []sdk.Event,
    expectedEventName string,
) error
```

---

## Phase 3: Comprehensive ApprovalCriteria Tests

### Goal
Test all ApprovalCriteria features comprehensively, including all nested structures.

### Implementation Plan

#### 3.1 Merkle Challenge Tests
**File:** `approval_criteria_test.go`

**Test Cases:**
- `TestApprovalCriteria_MerkleChallenge_Valid` - Valid merkle proof
- `TestApprovalCriteria_MerkleChallenge_Invalid` - Invalid merkle proof
- `TestApprovalCriteria_MerkleChallenge_Multiple` - Multiple merkle challenges
- `TestApprovalCriteria_MerkleChallenge_Expired` - Expired merkle challenge
- `TestApprovalCriteria_MerkleChallenge_TransferThroughPrecompile` - Full E2E with merkle

#### 3.2 Predetermined Balance Tests
- `TestApprovalCriteria_PredeterminedBalance_Exact` - Exact balance match
- `TestApprovalCriteria_PredeterminedBalance_Range` - Balance within range
- `TestApprovalCriteria_PredeterminedBalance_Multiple` - Multiple balance requirements
- `TestApprovalCriteria_PredeterminedBalance_NotMet` - Balance requirement not met
- `TestApprovalCriteria_PredeterminedBalance_TransferThroughPrecompile` - Full E2E

#### 3.3 Voting Challenge Tests
- `TestApprovalCriteria_VotingChallenge_SingleProposal` - Single proposal voting
- `TestApprovalCriteria_VotingChallenge_MultipleProposals` - Multiple proposals
- `TestApprovalCriteria_VotingChallenge_Threshold` - Voting threshold requirements
- `TestApprovalCriteria_VotingChallenge_Expired` - Expired proposals
- `TestApprovalCriteria_VotingChallenge_CastVoteThroughPrecompile` - Vote via precompile

#### 3.4 ETH Signature Challenge Tests
- `TestApprovalCriteria_ETHSignature_Valid` - Valid EIP712 signature
- `TestApprovalCriteria_ETHSignature_Invalid` - Invalid signature
- `TestApprovalCriteria_ETHSignature_Expired` - Expired signature
- `TestApprovalCriteria_ETHSignature_Replay` - Replay attack prevention
- `TestApprovalCriteria_ETHSignature_TransferThroughPrecompile` - Full E2E

#### 3.5 Time-Based Challenge Tests
- `TestApprovalCriteria_OfflineHours` - Offline hours requirement
- `TestApprovalCriteria_OfflineDays` - Offline days requirement
- `TestApprovalCriteria_TimeWindows` - Time window restrictions
- `TestApprovalCriteria_TimeChecks_TransferThroughPrecompile` - Full E2E

#### 3.6 Complex ApprovalCriteria Combinations
- `TestApprovalCriteria_AllCriteria_AND` - All criteria must be met
- `TestApprovalCriteria_AllCriteria_OR` - Any criteria can be met
- `TestApprovalCriteria_NestedCriteria` - Nested approval structures
- `TestApprovalCriteria_ComplexWorkflow` - Real-world complex scenario

---

## Phase 4: Concurrency and Race Condition Tests

### Goal
Test parallel execution safety and race conditions.

### Implementation Plan

#### 4.1 Concurrent Transfer Tests
**File:** `concurrency_test.go`

**Test Cases:**
- `TestConcurrency_ParallelTransfers_SameCollection` - Multiple transfers to same collection
- `TestConcurrency_ParallelTransfers_DifferentCollections` - Transfers to different collections
- `TestConcurrency_ParallelApprovals` - Multiple approvals set simultaneously
- `TestConcurrency_ParallelQueries` - Multiple queries in parallel
- `TestConcurrency_ParallelCollectionUpdates` - Multiple collection updates

#### 4.2 Race Condition Tests
- `TestRaceCondition_SimultaneousUpdates` - Update same collection simultaneously
- `TestRaceCondition_UpdateDuringTransfer` - Update collection during transfer
- `TestRaceCondition_ApprovalDuringTransfer` - Set approval during transfer
- `TestRaceCondition_DeleteDuringTransfer` - Delete collection during transfer

#### 4.3 Lock Verification Tests
- `TestConcurrency_StateConsistency` - State remains consistent under concurrency
- `TestConcurrency_NoPartialUpdates` - No partial state updates
- `TestConcurrency_AtomicOperations` - Operations are atomic

**Implementation:**
```go
// Use goroutines to simulate concurrent access
func (suite *ConcurrencyTestSuite) TestConcurrency_ParallelTransfers() {
    var wg sync.WaitGroup
    errors := make(chan error, 10)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            // Execute transfer
            err := suite.executeTransfer(index)
            if err != nil {
                errors <- err
            }
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    // Verify no errors occurred
    suite.Len(errors, 0)
    // Verify final state is correct
    suite.verifyFinalState()
}
```

---

## Phase 5: Reentrancy Tests

### Goal
Test reentrancy protection and attack scenarios.

### Implementation Plan

#### 5.1 Reentrancy Attack Tests
**File:** `reentrancy_test.go`

**Test Cases:**
- `TestReentrancy_TransferReentrancy` - Attempt reentrancy during transfer
- `TestReentrancy_ApprovalReentrancy` - Attempt reentrancy during approval
- `TestReentrancy_QueryReentrancy` - Attempt reentrancy during query
- `TestReentrancy_MaliciousContract` - Malicious contract attempting reentrancy

#### 5.2 Call Stack Tests
- `TestReentrancy_CallStackDepth` - Verify call stack depth limits
- `TestReentrancy_NestedCalls` - Nested precompile calls
- `TestReentrancy_CrossPrecompileCalls` - Calls between different precompiles

**Malicious Contract Example:**
```solidity
contract MaliciousContract {
    ITokenizationPrecompile precompile = ...;
    bool attacking = false;
    
    function attack() external {
        attacking = true;
        // Try to reenter during transfer
        precompile.transferTokens(...);
    }
    
    // Fallback that attempts reentrancy
    receive() external payable {
        if (attacking) {
            attacking = false;
            // Attempt reentrancy
            precompile.transferTokens(...);
        }
    }
}
```

---

## Phase 6: Event Verification Tests

### Goal
Verify events are properly emitted through EVM and can be indexed.

### Implementation Plan

#### 6.1 Event Emission Tests
**File:** `event_verification_test.go`

**Test Cases:**
- `TestEvents_TransferTokens_Emitted` - Transfer event emitted
- `TestEvents_CollectionCreated_Emitted` - Collection creation event
- `TestEvents_ApprovalSet_Emitted` - Approval set event
- `TestEvents_AllTransactionMethods` - All transaction methods emit events

#### 6.2 Event Data Verification
- `TestEvents_TransferTokens_DataCorrect` - Event data matches transaction
- `TestEvents_IndexedFields` - Indexed fields are correct
- `TestEvents_NonIndexedFields` - Non-indexed fields are correct
- `TestEvents_ThroughEVM` - Events emitted through EVM (not just direct calls)

#### 6.3 Event Parsing Tests
- `TestEvents_ParseTransferEvent` - Parse transfer event from logs
- `TestEvents_ParseAllEventTypes` - Parse all event types
- `TestEvents_EventIndexing` - Events can be indexed by external systems

---

## Phase 7: Gas Accuracy Tests

### Goal
Verify gas estimates match actual consumption within acceptable tolerance.

### Implementation Plan

#### 7.1 Gas Estimation Tests
**File:** `gas_accuracy_test.go`

**Test Cases:**
- `TestGasAccuracy_TransferTokens_EstimateVsActual` - Transfer gas accuracy
- `TestGasAccuracy_AllMethods_WithinTolerance` - All methods within 10% tolerance
- `TestGasAccuracy_DynamicGasCalculation` - Dynamic gas calculation accuracy
- `TestGasAccuracy_LargeInputs` - Gas accuracy with large inputs

#### 7.2 Gas Limit Tests
- `TestGasLimits_Enforced` - Gas limits are enforced
- `TestGasLimits_Refunds` - Gas refunds work correctly
- `TestGasLimits_OutOfGas` - Out of gas errors handled correctly

**Implementation:**
```go
func (suite *GasAccuracyTestSuite) TestGasAccuracy_TransferTokens() {
    // Estimate gas
    estimatedGas := suite.estimateGas("transferTokens", args)
    
    // Execute transaction
    actualGas, err := suite.executeAndMeasureGas("transferTokens", args)
    suite.NoError(err)
    
    // Verify within tolerance (10%)
    tolerance := estimatedGas / 10
    suite.InDelta(estimatedGas, actualGas, float64(tolerance))
}
```

---

## Phase 8: Stress Tests

### Goal
Test system behavior under load and with maximum inputs.

### Implementation Plan

#### 8.1 Large-Scale Tests
**File:** `stress_test.go`

**Test Cases:**
- `TestStress_MaxRecipients` - Transfer to maximum recipients (100)
- `TestStress_MaxRanges` - Maximum token ID ranges (100)
- `TestStress_MaxApprovals` - Maximum approvals (100)
- `TestStress_ManyCollections` - Create and manage many collections

#### 8.2 Performance Tests
- `TestPerformance_TransferThroughput` - Transfers per second
- `TestPerformance_QueryThroughput` - Queries per second
- `TestPerformance_MemoryUsage` - Memory usage under load
- `TestPerformance_GasConsumption` - Gas consumption patterns

---

## Phase 9: Multi-User Workflow Tests

### Goal
Test complex multi-user scenarios and workflows.

### Implementation Plan

#### 9.1 Multi-User Scenarios
**File:** `multi_user_workflow_test.go`

**Test Cases:**
- `TestMultiUser_ComplexApprovalWorkflow` - Multi-user approval workflow
- `TestMultiUser_ConcurrentCollectionManagement` - Multiple users managing collections
- `TestMultiUser_VotingWorkflow` - Voting workflow with multiple voters
- `TestMultiUser_AddressListManagement` - Shared address list management

**Example Workflow:**
```go
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_ComplexApprovalWorkflow() {
    // Alice creates collection
    collectionId := suite.createCollection(suite.Alice)
    
    // Bob sets incoming approval
    suite.setIncomingApproval(suite.Bob, collectionId)
    
    // Charlie sets outgoing approval
    suite.setOutgoingApproval(suite.Charlie, collectionId)
    
    // Alice transfers to Bob (requires Bob's incoming approval)
    suite.transferTokens(suite.Alice, suite.Bob, collectionId)
    
    // Bob transfers to Charlie (requires Charlie's incoming approval)
    suite.transferTokens(suite.Bob, suite.Charlie, collectionId)
    
    // Verify final state
    suite.verifyFinalBalances()
}
```

---

## Phase 10: Error Recovery Tests

### Goal
Test error handling and recovery paths.

### Implementation Plan

#### 10.1 Error Recovery Tests
**File:** `error_recovery_test.go`

**Test Cases:**
- `TestErrorRecovery_PartialFailure` - Partial operation failure
- `TestErrorRecovery_StateRollback` - State rollback on error
- `TestErrorRecovery_RetryMechanisms` - Retry after error
- `TestErrorRecovery_ErrorPropagation` - Error propagation through layers

---

## Implementation Priority

### 🔴 High Priority (Must Have)
1. **Phase 1: Full EVM Keeper Integration** - Critical for production readiness
2. **Phase 2: Solidity Contract Tests** - Essential for real-world usage

### 🟡 Medium Priority (Should Have)
3. **Phase 3: ApprovalCriteria Tests** - Important for feature completeness
4. **Phase 4: Concurrency Tests** - Important for production safety
5. **Phase 5: Reentrancy Tests** - Security critical
6. **Phase 6: Event Verification** - Important for indexing and monitoring

### 🟢 Lower Priority (Nice to Have)
7. **Phase 7: Gas Accuracy Tests** - Can be done post-launch
8. **Phase 8: Stress Tests** - Performance optimization
9. **Phase 9: Multi-User Workflows** - Edge case coverage
10. **Phase 10: Error Recovery** - Edge case coverage

---

## Test Infrastructure Requirements

### 1. EVM Keeper Setup Helper
**File:** `test_helpers_evm.go`

```go
// setupAppWithEVM creates a full app instance with EVM keeper
func setupAppWithEVM(t *testing.T) (*app.App, sdk.Context, *evmkeeper.Keeper, *Precompile) {
    app := app.Setup(false)
    ctx := app.BaseApp.NewContext(false, tmproto.Header{})
    
    // EVM keeper is already initialized in app
    evmKeeper := app.EVMKeeper
    
    // Register precompile (already done in app/evm.go, but verify)
    tokenizationPrecompile := NewPrecompile(app.TokenizationKeeper)
    precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
    
    // Verify precompile is registered
    precompile := evmKeeper.GetPrecompile(precompileAddr)
    require.NotNil(t, precompile)
    
    return app, ctx, evmKeeper, tokenizationPrecompile
}

// createEVMAccount creates an EVM account with private key
func createEVMAccount() (*ecdsa.PrivateKey, common.Address, sdk.AccAddress) {
    key, _ := crypto.GenerateKey()
    addr := crypto.PubkeyToAddress(key.PublicKey)
    cosmosAddr := sdk.AccAddress(addr.Bytes())
    return key, addr, cosmosAddr
}

// fundEVMAccount funds an EVM account with native tokens
func fundEVMAccount(ctx sdk.Context, bankKeeper bankkeeper.Keeper, addr sdk.AccAddress, amount sdk.Coins) error {
    return bankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, amount)
}
```

### 2. Solidity Contract Compiler Integration
**File:** `solidity_test_helpers.go`

```go
// compileSolidityContract compiles a Solidity contract
func compileSolidityContract(contractPath string) ([]byte, *abi.ABI, error) {
    // Use solc or solcjs to compile
    // Return bytecode and ABI
}

// embedSolidityContract embeds compiled contract in test
//go:embed contracts/test/PrecompileTestContract.bin
var testContractBytecode string

//go:embed contracts/test/PrecompileTestContract.abi
var testContractABI string
```

### 3. EVM Transaction Builder
**File:** `evm_transaction_helpers.go`

```go
// buildEVMTransaction builds an EVM transaction
func buildEVMTransaction(
    fromKey *ecdsa.PrivateKey,
    to common.Address,
    data []byte,
    value *big.Int,
    gasLimit uint64,
    gasPrice *big.Int,
    nonce uint64,
    chainID *big.Int,
) (*types.Transaction, error) {
    tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
    return types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
}

// executeEVMTransaction executes a transaction through EVM keeper
func executeEVMTransaction(
    ctx sdk.Context,
    evmKeeper *evmkeeper.Keeper,
    tx *types.Transaction,
) (*evmtypes.MsgEthereumTxResponse, error) {
    // Convert transaction to MsgEthereumTx
    msg := &evmtypes.MsgEthereumTx{
        Data: &evmtypes.AnyMsg{
            Value: tx,
        },
    }
    
    // Execute through keeper
    return evmKeeper.EthereumTx(ctx, msg)
}
```

### 4. Event Parser
**File:** `event_parser.go`

```go
// ParsedEvent represents a parsed event
type ParsedEvent struct {
    Name      string
    Indexed   map[string]interface{}
    Data      map[string]interface{}
    Address   common.Address
}

// parseEvents parses events from transaction response
func parseEvents(events []sdk.Event, precompileABI abi.ABI) ([]ParsedEvent, error) {
    // Parse events from SDK events
    // Match against precompile ABI events
    // Return parsed events
}

// verifyEvent verifies an event matches expected values
func verifyEvent(event ParsedEvent, expectedName string, expectedData map[string]interface{}) error {
    // Verify event name
    // Verify event data matches expected values
}
```

---

## Testing Tools and Libraries

### Required
- `github.com/ethereum/go-ethereum` - EVM integration, transaction building
- `github.com/stretchr/testify/suite` - Test suites
- `github.com/cosmos/evm` - Cosmos EVM module
- `github.com/ethereum/go-ethereum/crypto` - Cryptographic functions

### Optional
- `github.com/onsi/ginkgo` - BDD testing (if preferred)
- `github.com/onsi/gomega` - Matchers (if using Ginkgo)
- Solidity compiler (`solc` or `solcjs`) - For compiling test contracts

---

## Success Criteria

### Phase 1-2 (Critical) ✅
- [ ] All transaction methods tested through EVM keeper
- [ ] All query methods tested through EVM keeper
- [ ] At least one Solidity contract deployed and tested
- [ ] Events verified through EVM
- [ ] Gas accounting verified
- [ ] Error handling verified through EVM

### Phase 3-6 (Important) ✅
- [ ] All ApprovalCriteria features tested
- [ ] Concurrency safety verified
- [ ] Reentrancy protection verified
- [ ] Event emission verified

### Phase 7-10 (Nice to Have) ✅
- [ ] Gas accuracy within 10% tolerance
- [ ] Stress tests pass
- [ ] Multi-user workflows tested
- [ ] Error recovery verified

---

## Estimated Effort

| Phase | Estimated Time | Priority | Dependencies |
|-------|---------------|----------|--------------|
| Phase 1: EVM Keeper Integration | 2-3 days | 🔴 High | app.Setup() |
| Phase 2: Solidity Contract Tests | 2-3 days | 🔴 High | Phase 1 |
| Phase 3: ApprovalCriteria Tests | 3-4 days | 🟡 Medium | Phase 1 |
| Phase 4: Concurrency Tests | 2-3 days | 🟡 Medium | Phase 1 |
| Phase 5: Reentrancy Tests | 1-2 days | 🟡 Medium | Phase 2 |
| Phase 6: Event Verification | 1-2 days | 🟡 Medium | Phase 1 |
| Phase 7: Gas Accuracy | 1-2 days | 🟢 Low | Phase 1 |
| Phase 8: Stress Tests | 2-3 days | 🟢 Low | Phase 1 |
| Phase 9: Multi-User Workflows | 2-3 days | 🟢 Low | Phase 1 |
| Phase 10: Error Recovery | 1-2 days | 🟢 Low | Phase 1 |
| **Total** | **17-27 days** | | |

---

## Next Steps

1. **Start with Phase 1**: Implement full EVM keeper integration tests
   - Create `evm_keeper_integration_test.go`
   - Set up test infrastructure
   - Implement basic test cases

2. **Then Phase 2**: Add Solidity contract tests
   - Create test contract
   - Set up compilation infrastructure
   - Implement contract interaction tests

3. **Iterate**: Add remaining phases based on priority and time available

---

## Example: Phase 1 Implementation Skeleton

```go
package tokenization

import (
    "crypto/ecdsa"
    "math/big"
    "testing"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/stretchr/testify/suite"
    
    evmkeeper "github.com/cosmos/evm/x/vm/keeper"
    evmtypes "github.com/cosmos/evm/x/vm/types"
    
    "github.com/bitbadges/bitbadgeschain/app"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type EVMKeeperIntegrationTestSuite struct {
    suite.Suite
    App            *app.App
    Ctx            sdk.Context
    EVMKeeper      *evmkeeper.Keeper
    Precompile     *Precompile
    
    AliceKey       *ecdsa.PrivateKey
    BobKey         *ecdsa.PrivateKey
    AliceEVM       common.Address
    BobEVM         common.Address
    Alice          sdk.AccAddress
    Bob            sdk.AccAddress
    
    CollectionId   sdkmath.Uint
}

func TestEVMKeeperIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(EVMKeeperIntegrationTestSuite))
}

func (suite *EVMKeeperIntegrationTestSuite) SetupTest() {
    // Set up app with EVM keeper
    suite.App = app.Setup(false)
    suite.Ctx = suite.App.BaseApp.NewContext(false, tmproto.Header{})
    suite.EVMKeeper = suite.App.EVMKeeper
    
    // Get precompile (should be registered in app/evm.go)
    precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
    precompile := suite.EVMKeeper.GetPrecompile(precompileAddr)
    suite.Require().NotNil(precompile)
    suite.Precompile = precompile.(*Precompile)
    
    // Create test accounts
    suite.AliceKey, suite.AliceEVM, suite.Alice = createEVMAccount()
    suite.BobKey, suite.BobEVM, suite.Bob = createEVMAccount()
    
    // Fund accounts
    suite.fundAccount(suite.Alice, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(1000000))))
    suite.fundAccount(suite.Bob, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(1000000))))
    
    // Create test collection
    suite.CollectionId = suite.createTestCollection()
}

func (suite *EVMKeeperIntegrationTestSuite) TestEVMKeeper_TransferTokens_ThroughEVM() {
    // Create EVM transaction
    method := suite.Precompile.ABI.Methods["transferTokens"]
    args := []interface{}{
        suite.CollectionId.BigInt(),
        []common.Address{suite.BobEVM},
        big.NewInt(10),
        []struct{Start, End *big.Int}{{Start: big.NewInt(1), End: big.NewInt(10)}},
        []struct{Start, End *big.Int}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
    }
    
    packed, err := method.Inputs.Pack(args...)
    suite.Require().NoError(err)
    
    input := append(method.ID, packed...)
    
    // Build transaction
    tx, err := buildEVMTransaction(
        suite.AliceKey,
        common.HexToAddress(TokenizationPrecompileAddress),
        input,
        big.NewInt(0),
        100000,
        big.NewInt(1000000000),
        0,
        big.NewInt(1),
    )
    suite.Require().NoError(err)
    
    // Execute transaction
    response, err := executeEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)
    suite.Require().NoError(err)
    suite.Require().NotNil(response)
    
    // Verify state changes
    // ... check balances changed
    
    // Verify events
    // ... check events emitted
}
```

---

## Resources

- [Cosmos EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview)
- [Ethereum Go Client Documentation](https://geth.ethereum.org/docs/)
- [Testify Suite Documentation](https://github.com/stretchr/testify#suite-package)
- [Solidity Compiler Documentation](https://docs.soliditylang.org/)
