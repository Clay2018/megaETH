pragma solidity >=0.4.0 <0.8.0;

library Stack {

    uint256 constant MAX_LENGTH = 20;

    struct uint256Stack {
        uint256[MAX_LENGTH] buffer;
        uint256 index;
    }

    function push(uint256Stack memory s, uint256 addr) internal pure {
        s.buffer[s.index] = addr;
        s.index += 1;
    }

    function pop(uint256Stack memory s) internal pure returns(uint256) {
        require(s.index != 0);
        s.index -= 1;
        return s.buffer[s.index];
    }

    function peek(uint256Stack memory s) internal pure returns(uint256, uint256) {
        //need justify s.index
        return (s.buffer[(s.index-1)], s.index-1);
    }
}