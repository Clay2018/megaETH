const { expect } = require("chai");
const hre = require("hardhat");
const { time, loadFixture } = require("@nomicfoundation/hardhat-network-helpers")
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");
const BigNumber = require('bignumber.js');
const { ethers } = require("hardhat");


describe("dijsktras-test", function () {
    async function _deploy() {

        const [owner, GovernorAddr, userAddr2] = await hre.ethers.getSigners();

        const YourToken = await hre.ethers.getContractFactory("DijsktrasVariant");
        const yourToken = await YourToken.deploy();
        await yourToken.deployed()


        const addrs = [
            "0x0000000000000000000000000000000000000001",
            "0x0000000000000000000000000000000000000002",
            "0x0000000000000000000000000000000000000003",
            "0x0000000000000000000000000000000000000004",
            "0x0000000000000000000000000000000000000005",
            "0x0000000000000000000000000000000000000006",
            "0x0000000000000000000000000000000000000007",
            "0x0000000000000000000000000000000000000008",
            "0x0000000000000000000000000000000000000009",
        ]


        await yourToken.addPoolRelation(addrs[0], addrs[1]);
        await yourToken.addPoolRelation(addrs[0], addrs[2]);
        await yourToken.addPoolRelation(addrs[1], addrs[3]);
        await yourToken.addPoolRelation(addrs[1], addrs[0]);
        await yourToken.addPoolRelation(addrs[1], addrs[4]);
        await yourToken.addPoolRelation(addrs[2], addrs[0]);
        await yourToken.addPoolRelation(addrs[2], addrs[5]);
        await yourToken.addPoolRelation(addrs[2], addrs[6]);
        await yourToken.addPoolRelation(addrs[3], addrs[1]);
        await yourToken.addPoolRelation(addrs[3], addrs[7]);
        await yourToken.addPoolRelation(addrs[4], addrs[7]);
        await yourToken.addPoolRelation(addrs[4], addrs[1]);
        await yourToken.addPoolRelation(addrs[4], addrs[5]);
        await yourToken.addPoolRelation(addrs[5], addrs[2]);
        await yourToken.addPoolRelation(addrs[5], addrs[4]);
        await yourToken.addPoolRelation(addrs[5], addrs[6]);
        await yourToken.addPoolRelation(addrs[6], addrs[2]);
        await yourToken.addPoolRelation(addrs[6], addrs[5]);
        await yourToken.addPoolRelation(addrs[6], addrs[8]);
        await yourToken.addPoolRelation(addrs[7], addrs[3]);
        await yourToken.addPoolRelation(addrs[7], addrs[4]);
        await yourToken.addPoolRelation(addrs[8], addrs[6]);

        //await yourToken.debugModifyGraph(3, [1, 7, 6]);

        ret = await yourToken.dijsktras_variant(addrs[0], addrs[3]);
        console.log(ret);

        let transReceipt = await hre.ethers.provider.getTransactionReceipt(ret.hash)
        console.log(transReceipt);




        //paths = ret[0];
        //pathIndex = ret[1];
        //console.log("total path:", pathIndex);
        //for(i = 0; i < pathIndex; i++) {
        //    console.log("path: ", paths[i]);
        //}

        //ret = await yourToken.graph(3, 0);
        //console.log("ret:", ret);
        //ret = await yourToken.graph(3, 1);
        //console.log("ret:", ret);
        //ret = await yourToken.graph(3, 2);
        //console.log("ret:", ret);

        return {yourToken, owner};
    }

    describe("YourToken", function() {
        it("checkowner", async function() {
            await _deploy()
        })
    })
})
