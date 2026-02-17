// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "../libraries/TokenizationJSONHelpers.sol";

/// @title PrecompileDynamicStoreTestContract
/// @notice Test contract for dynamic store precompile methods
/// @dev Split from PrecompileTestContract to stay under EVM size limits
contract PrecompileDynamicStoreTestContract {
    ITokenizationPrecompile constant precompile =
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    // ============ Events ============

    event DynamicStoreCreated(uint256 indexed storeId);

    event DynamicStoreUpdated(uint256 indexed storeId, bool success);

    event DynamicStoreDeleted(uint256 indexed storeId, bool success);

    event DynamicStoreValueSet(
        uint256 indexed storeId,
        address indexed address_,
        bool value,
        bool success
    );

    // ============ Dynamic Store Methods ============

    /// @notice Wrapper for createDynamicStore
    function testCreateDynamicStore(
        bool defaultValue,
        string calldata uri,
        string calldata customData
    ) external returns (uint256) {
        string memory createJson = TokenizationJSONHelpers.createDynamicStoreJSON(
            defaultValue,
            uri,
            customData
        );
        uint256 storeId = precompile.createDynamicStore(createJson);
        emit DynamicStoreCreated(storeId);
        return storeId;
    }

    /// @notice Wrapper for updateDynamicStore
    function testUpdateDynamicStore(
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string calldata uri,
        string calldata customData
    ) external returns (bool) {
        string memory updateJson = TokenizationJSONHelpers.updateDynamicStoreJSON(
            storeId,
            defaultValue,
            globalEnabled,
            uri,
            customData
        );
        bool success = precompile.updateDynamicStore(updateJson);
        emit DynamicStoreUpdated(storeId, success);
        return success;
    }

    /// @notice Wrapper for deleteDynamicStore
    function testDeleteDynamicStore(uint256 storeId) external returns (bool) {
        string memory deleteJson = TokenizationJSONHelpers.deleteDynamicStoreJSON(storeId);
        bool success = precompile.deleteDynamicStore(deleteJson);
        emit DynamicStoreDeleted(storeId, success);
        return success;
    }

    /// @notice Wrapper for setDynamicStoreValue
    function testSetDynamicStoreValue(
        uint256 storeId,
        address address_,
        bool value
    ) external returns (bool) {
        string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            storeId,
            address_,
            value
        );
        bool success = precompile.setDynamicStoreValue(setValueJson);
        emit DynamicStoreValueSet(storeId, address_, value, success);
        return success;
    }

    // ============ Query Methods (View) ============

    /// @notice Wrapper for getDynamicStore
    function testGetDynamicStore(uint256 storeId) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getDynamicStoreJSON(storeId);
        return precompile.getDynamicStore(queryJson);
    }

    /// @notice Wrapper for getDynamicStoreValue
    function testGetDynamicStoreValue(
        uint256 storeId,
        address address_
    ) external view returns (bytes memory) {
        string memory queryJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(
            storeId,
            address_
        );
        return precompile.getDynamicStoreValue(queryJson);
    }

    /// @notice Wrapper for params
    function testParams() external view returns (bytes memory) {
        return precompile.params("{}");
    }

    /// @notice Wrapper for isAddressReservedProtocol
    function testIsAddressReservedProtocol(address addr) external view returns (bool) {
        string memory queryJson = TokenizationJSONHelpers.isAddressReservedProtocolJSON(addr);
        return precompile.isAddressReservedProtocol(queryJson);
    }

    /// @notice Wrapper for getAllReservedProtocolAddresses
    function testGetAllReservedProtocolAddresses() external view returns (address[] memory) {
        return precompile.getAllReservedProtocolAddresses("{}");
    }
}
