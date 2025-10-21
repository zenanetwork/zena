require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    compilers: [
      {
        version: "0.8.18",
      },
      // This version is required to compile the werc9 contract.
      {
        version: "0.4.22",
      },
    ],
  },
  networks: {
    cosmos: {
      url: "http://127.0.0.1:8545",
      chainId: 262144,
      accounts: [
        "0x8BE1E5311E4CB31002C5C84CEA459B5E598592F1D00C796E3DE2880D55FE9990", // mykey (validator)
        "0xF91067EF80B57C9D04D8F6E45F458D81D6B65397EBECCE54693E398A6AF6D347", // dev0
        "0x8112578E13EFAF1A189FACB02D595EAC6AC95C4416F4D037FE3C7A98CE36E80A", // dev1
      ],
      timeout: 60000,              // Network timeout 60 seconds
      gasPrice: 10000000000        // 10 Gwei fixed (handles base fee increases)
    },
  },
  mocha: {
    timeout: 120000  // Global Mocha timeout 120 seconds
  },
};
