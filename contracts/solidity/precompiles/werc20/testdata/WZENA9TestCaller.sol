// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.17;

import "../IWERC20.sol";

contract WZENA9TestCaller {
    address payable public immutable WZENA;
    uint256 public counter;

    constructor(address payable _wrappedTokenAddress) {
        WZENA = _wrappedTokenAddress;
        counter = 0;
    }

    event Log(string message);

    function depositWithRevert(bool before, bool aft) public payable {
        counter++;

        uint amountIn = msg.value;
        IWERC20(WZENA).deposit{value: amountIn}();

        if (before) {
            require(false, "revert here");
        }

        counter--;

        if (aft) {
            require(false, "revert here");
        }
        return;
    }
}