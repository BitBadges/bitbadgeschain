// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct UintRange {
    uint256 start;
    uint256 end;
}

interface IBadgesPrecompile {
    function transferTokens(
        uint256 collectionId,
        address[] calldata toAddresses,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool);
}

