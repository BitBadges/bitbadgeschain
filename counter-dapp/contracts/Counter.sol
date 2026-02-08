// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;



contract Counter {
    uint256 public count;
    
    event CountIncremented(uint256 newCount);
    
    constructor(uint256 _initialCount) {
        count = _initialCount;
    }
    
    function increment() public {
        count += 1;
        emit CountIncremented(count);
    }
    
    function getCount() public view returns (uint256) {
        return count;
    }
}