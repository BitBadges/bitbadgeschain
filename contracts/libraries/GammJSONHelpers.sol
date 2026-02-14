// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/GammTypes.sol";

/**
 * @title GammJSONHelpers
 * @notice Helper library for constructing JSON strings for the gamm precompile
 * @dev All methods return JSON strings that match the protobuf JSON format
 * 
 * Usage example:
 * ```solidity
 * string memory coinsJson = GammJSONHelpers.coinsToJson(tokenInMaxs);
 * string memory json = GammJSONHelpers.joinPoolJSON(poolId, shareOutAmount, coinsJson);
 * (uint256 shares, Coin[] memory tokens) = precompile.joinPool(json);
 * ```
 */
library GammJSONHelpers {
    // ============ Transaction JSON Constructors ============

    /**
     * @notice Construct JSON for joinPool
     */
    function joinPoolJSON(
        uint64 poolId,
        uint256 shareOutAmount,
        string memory tokenInMaxsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId),
            '","shareOutAmount":"', _uintToString(shareOutAmount),
            '","tokenInMaxs":', tokenInMaxsJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for exitPool
     */
    function exitPoolJSON(
        uint64 poolId,
        uint256 shareInAmount,
        string memory tokenOutMinsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId),
            '","shareInAmount":"', _uintToString(shareInAmount),
            '","tokenOutMins":', tokenOutMinsJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for swapExactAmountIn
     */
    function swapExactAmountInJSON(
        string memory routesJson,
        string memory tokenInJson,
        uint256 tokenOutMinAmount,
        string memory affiliatesJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"routes":', routesJson,
            ',"tokenIn":', tokenInJson,
            ',"tokenOutMinAmount":"', _uintToString(tokenOutMinAmount),
            '","affiliates":', affiliatesJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for swapExactAmountInWithIBCTransfer
     */
    function swapExactAmountInWithIBCTransferJSON(
        string memory routesJson,
        string memory tokenInJson,
        uint256 tokenOutMinAmount,
        string memory ibcTransferInfoJson,
        string memory affiliatesJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"routes":', routesJson,
            ',"tokenIn":', tokenInJson,
            ',"tokenOutMinAmount":"', _uintToString(tokenOutMinAmount),
            '","ibcTransferInfo":', ibcTransferInfoJson,
            ',"affiliates":', affiliatesJson,
            '}'
        ));
    }

    // ============ Query JSON Constructors ============

    /**
     * @notice Construct JSON for getPool
     */
    function getPoolJSON(
        uint64 poolId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getPools
     * @param paginationJson JSON string for pagination (can be empty string for default)
     */
    function getPoolsJSON(
        string memory paginationJson
    ) internal pure returns (string memory) {
        if (bytes(paginationJson).length == 0) {
            return "{}";
        }
        return string(abi.encodePacked(
            '{"pagination":', paginationJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for getPoolType
     */
    function getPoolTypeJSON(
        uint64 poolId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for calcJoinPoolNoSwapShares
     */
    function calcJoinPoolNoSwapSharesJSON(
        uint64 poolId,
        string memory tokenInMaxsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId),
            '","tokenInMaxs":', tokenInMaxsJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for calcExitPoolCoinsFromShares
     */
    function calcExitPoolCoinsFromSharesJSON(
        uint64 poolId,
        uint256 shareInAmount
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId),
            '","shareInAmount":"', _uintToString(shareInAmount), '"}'
        ));
    }

    /**
     * @notice Construct JSON for calcJoinPoolShares
     */
    function calcJoinPoolSharesJSON(
        uint64 poolId,
        string memory tokenInMaxsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId),
            '","tokenInMaxs":', tokenInMaxsJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for getPoolParams
     */
    function getPoolParamsJSON(
        uint64 poolId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getTotalShares
     */
    function getTotalSharesJSON(
        uint64 poolId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getTotalLiquidity
     */
    function getTotalLiquidityJSON(
        uint64 poolId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(poolId), '"}'
        ));
    }

    // ============ Type to JSON Converters ============

    /**
     * @notice Convert Coin to JSON
     */
    function coinToJson(GammTypes.Coin memory coin) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"denom":"', _escapeJsonString(coin.denom),
            '","amount":"', _uintToString(coin.amount), '"}'
        ));
    }

    /**
     * @notice Convert Coin array to JSON array
     */
    function coinsToJson(GammTypes.Coin[] memory coins) internal pure returns (string memory) {
        if (coins.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < coins.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, coinToJson(coins[i])));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Convert SwapAmountInRoute to JSON
     */
    function swapRouteToJson(GammTypes.SwapAmountInRoute memory route) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"poolId":"', _uint64ToString(route.poolId),
            '","tokenOutDenom":"', _escapeJsonString(route.tokenOutDenom), '"}'
        ));
    }

    /**
     * @notice Convert SwapAmountInRoute array to JSON array
     */
    function swapRoutesToJson(GammTypes.SwapAmountInRoute[] memory routes) internal pure returns (string memory) {
        if (routes.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < routes.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, swapRouteToJson(routes[i])));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Convert Affiliate to JSON
     * @dev Note: address is converted to hex string, basisPointsFee is converted to string
     */
    function affiliateToJson(GammTypes.Affiliate memory affiliate) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"address":"', _addressToString(affiliate.address_),
            '","basisPointsFee":"', _uintToString(affiliate.basisPointsFee), '"}'
        ));
    }

    /**
     * @notice Convert Affiliate array to JSON array
     */
    function affiliatesToJson(GammTypes.Affiliate[] memory affiliates) internal pure returns (string memory) {
        if (affiliates.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < affiliates.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, affiliateToJson(affiliates[i])));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Convert IBCTransferInfo to JSON
     */
    function ibcTransferInfoToJson(GammTypes.IBCTransferInfo memory info) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"sourceChannel":"', _escapeJsonString(info.sourceChannel),
            '","receiver":"', _escapeJsonString(info.receiver),
            '","memo":"', _escapeJsonString(info.memo),
            '","timeoutTimestamp":"', _uint64ToString(info.timeoutTimestamp), '"}'
        ));
    }

    /**
     * @notice Convert pagination to JSON
     * @param key Pagination key (can be empty)
     * @param offset Offset (0 if not used)
     * @param limit Limit (0 if not used)
     * @param countTotal Whether to count total (false if not used)
     */
    function paginationToJson(
        string memory key,
        uint64 offset,
        uint64 limit,
        bool countTotal
    ) internal pure returns (string memory) {
        string memory result = "{";
        bool hasFields = false;

        if (bytes(key).length > 0) {
            result = string(abi.encodePacked(result, '"key":"', _escapeJsonString(key), '"'));
            hasFields = true;
        }

        if (offset > 0) {
            if (hasFields) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"offset":"', _uint64ToString(offset), '"'));
            hasFields = true;
        }

        if (limit > 0) {
            if (hasFields) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"limit":"', _uint64ToString(limit), '"'));
            hasFields = true;
        }

        if (countTotal) {
            if (hasFields) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"countTotal":true'));
        }

        result = string(abi.encodePacked(result, "}"));
        return result;
    }

    // ============ Internal Helpers ============

    /**
     * @notice Convert uint256 to string
     */
    function _uintToString(uint256 value) private pure returns (string memory) {
        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }

    /**
     * @notice Convert uint64 to string
     */
    function _uint64ToString(uint64 value) private pure returns (string memory) {
        return _uintToString(uint256(value));
    }

    /**
     * @notice Convert address to string (hex format)
     */
    function _addressToString(address addr) private pure returns (string memory) {
        bytes memory data = abi.encodePacked(addr);
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < data.length; i++) {
            str[2 + i * 2] = alphabet[uint8(data[i] >> 4)];
            str[3 + i * 2] = alphabet[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }

    /**
     * @notice Escape JSON string (escape quotes, backslashes, newlines, etc.)
     * @dev Basic escaping for JSON strings - escapes quotes, backslashes, and control characters
     */
    function _escapeJsonString(string memory str) private pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        bytes memory result = new bytes(strBytes.length * 2); // Worst case: all chars need escaping
        uint256 resultIndex = 0;
        
        for (uint256 i = 0; i < strBytes.length; i++) {
            bytes1 char = strBytes[i];
            if (char == 0x22) { // "
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x22; // "
            } else if (char == 0x5C) { // \
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x5C; // \
            } else if (char == 0x0A) { // \n
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x6E; // n
            } else if (char == 0x0D) { // \r
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x72; // r
            } else if (char == 0x09) { // \t
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x74; // t
            } else if (char >= 0x20) { // Printable ASCII
                result[resultIndex++] = char;
            }
            // Control characters < 0x20 are skipped (except \n, \r, \t which are handled above)
        }
        
        // Resize result to actual length
        bytes memory finalResult = new bytes(resultIndex);
        for (uint256 i = 0; i < resultIndex; i++) {
            finalResult[i] = result[i];
        }
        return string(finalResult);
    }
}

