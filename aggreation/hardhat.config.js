/** @type import('hardhat/config').HardhatUserConfig */
require('@nomiclabs/hardhat-ethers');
require('@nomicfoundation/hardhat-chai-matchers');
const dotenv = require('dotenv');
dotenv.config()

module.exports = {
  solidity: "0.7.6",
  mocha: {
    timeout: 100000
  },
  //networks: {
  //  hardhat: {
  //    forking: {
  //      url: `https://goerli.infura.io/v3/${process.env.WEB3_INFURA_PROJECT_ID}`,
  //      blockNumber: 7601392
  //    },
  //    accounts: {
  //      count: 10,
  //    }
  //  }
  //}
};
