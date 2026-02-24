// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MinBankBalanceChecker
 * @notice Reference implementation for EVM query challenges in approval criteria (v25)
 * @dev This contract demonstrates how to use EVM query challenges for pre-transfer validation
 *      as part of approval criteria. It checks that the sender has a minimum balance before
 *      allowing the transfer.
 *
 * ## Overview
 * EVM query challenges in approval criteria allow you to add custom validation logic
 * that must pass before a transfer is approved. This is different from invariants which
 * validate post-transfer state.
 *
 * ## How It Works
 * 1. When a transfer is initiated, the chain calls this contract with the sender's address
 * 2. This contract checks if the sender has sufficient balance (simulated for testing)
 * 3. It returns `bytes32(1)` to approve, `bytes32(0)` to reject
 *
 * ## Placeholders
 * In EVM query challenge calldata, you can use these placeholders:
 * - $sender - The sender (from) address
 * - $recipient - The recipient (to) address
 * - $initiator - The transaction initiator address
 * - $collectionId - The collection ID
 *
 * ## Usage in Approval Criteria
 * ```solidity
 * EVMQueryChallenge memory challenge = TokenizationHelpers.createEVMQueryChallenge(
 *     address(minBalanceChecker),
 *     abi.encodeWithSelector(
 *         MinBankBalanceChecker.checkMinBalance.selector,
 *         "$sender",       // Placeholder for sender address
 *         1000000          // Min 1M units required
 *     ),
 *     bytes32(uint256(1)), // Expected: pass (1)
 *     "eq",
 *     100000
 * );
 * ```
 *
 * ## Note on Production Use
 * This example uses simulated balances for testing. In production, you could:
 * - Query actual bank balances via bank precompile (once ERC20 pairs are registered)
 * - Query ERC20 token balances
 * - Check any other on-chain state
 *
 * ## See Also
 * - MaxUniqueHoldersChecker.sol: Example of invariant-level EVM query challenges
 * - TokenizationHelpers.sol: Helper functions for creating EVM query challenges
 */
contract MinBankBalanceChecker {
    // Simulated balances for testing (address => balance)
    mapping(address => uint256) public simulatedBalances;

    /// @notice Set a simulated balance for an account (for E2E testing).
    /// @param account Address to set balance for.
    /// @param amount Balance amount to set.
    function setSimulatedBalance(address account, uint256 amount) external {
        simulatedBalances[account] = amount;
    }

    /// @notice Check that account has at least minAmount of simulated balance.
    /// @param account Address to check (use $sender in approval calldata).
    /// @param minAmount Minimum required amount.
    /// @return Pass: bytes32(1) if balance >= minAmount, else bytes32(0).
    function checkMinBalance(address account, uint256 minAmount) external view returns (bytes32) {
        return simulatedBalances[account] >= minAmount ? bytes32(uint256(1)) : bytes32(uint256(0));
    }
}
