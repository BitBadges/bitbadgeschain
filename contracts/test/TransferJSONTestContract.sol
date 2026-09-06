// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../libraries/TokenizationJSONHelpers.sol";

contract TransferJSONTestContract {
    function defaultSender(
        uint256 collectionId, address[] memory recipients, uint256 amount,
        string memory tokenIds, string memory ownershipTimes
    ) external pure returns (string memory) {
        return TokenizationJSONHelpers.transferTokensJSON(collectionId, recipients, amount, tokenIds, ownershipTimes);
    }

    function explicitSender(
        uint256 collectionId, address from, address[] memory recipients, uint256 amount,
        string memory tokenIds, string memory ownershipTimes
    ) external pure returns (string memory) {
        return TokenizationJSONHelpers.transferTokensJSON(collectionId, from, recipients, amount, tokenIds, ownershipTimes);
    }
}
