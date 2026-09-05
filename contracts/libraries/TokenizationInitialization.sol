// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/ITokenizationPrecompile.sol";
import "./TokenizationJSONHelpers.sol";

library TokenizationInitialization {
    // Only used immediately after creating an approval-free collection.
    function mintInitialSupply(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address recipient,
        string memory balancesJson
    ) internal {
        string memory recipientString = _addressToBech32(recipient);
        string memory fullRange = '[{"start":"1","end":"18446744073709551615"}]';
        string memory approval = string(abi.encodePacked(
            '[{"approvalId":"initial-supply","fromListId":"Mint","toListId":"', recipientString,
            '","initiatedByListId":"', _addressToBech32(address(this)),
            '","transferTimes":', fullRange, ',"tokenIds":', fullRange, ',"ownershipTimes":', fullRange,
            ',"approvalCriteria":{"overridesFromOutgoingApprovals":true,"overridesToIncomingApprovals":true}}]'
        ));
        precompile.setCollectionApprovals(TokenizationJSONHelpers.setCollectionApprovalsJSON(collectionId, approval, "[]"));
        require(precompile.transferTokens(string(abi.encodePacked(
            '{"collectionId":"', TokenizationJSONHelpers.uintToString(collectionId),
            '","transfers":[{"from":"Mint","toAddresses":["', recipientString,
            '"],"balances":', balancesJson, '}]}'
        ))), "Initial supply transfer failed");
        precompile.setCollectionApprovals(TokenizationJSONHelpers.setCollectionApprovalsJSON(collectionId, "[]", "[]"));
    }

    // Approval list literals use the chain's bb address format.
    function _addressToBech32(address account) private pure returns (string memory) {
        bytes memory alphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l";
        bytes memory output = new bytes(41);
        output[0] = "b";
        output[1] = "b";
        output[2] = "1";
        bytes memory values = new bytes(43);
        values[0] = 0x03;
        values[1] = 0x03;
        values[3] = 0x02;
        values[4] = 0x02;
        for (uint256 i = 0; i < 32; i++) {
            uint8 value = uint8(uint160(account) >> (5 * (31 - i))) & 31;
            values[5 + i] = bytes1(value);
            output[3 + i] = alphabet[value];
        }
        uint32 checksum = 1;
        uint32[5] memory generators = [uint32(0x3b6a57b2), 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3];
        for (uint256 i = 0; i < values.length; i++) {
            uint32 high = checksum >> 25;
            checksum = ((checksum & 0x1ffffff) << 5) ^ uint8(values[i]);
            for (uint256 j = 0; j < 5; j++) {
                if ((high >> j) & 1 != 0) checksum ^= generators[j];
            }
        }
        checksum ^= 1;
        for (uint256 i = 0; i < 6; i++) output[35 + i] = alphabet[(checksum >> (5 * (5 - i))) & 31];
        return string(output);
    }
}
