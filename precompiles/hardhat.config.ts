const config = {
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./contracts",
  },
  solidity: {
    version: "0.8.19",
    settings: {
      optimizer: {
        enabled: false,
      },
    },
  },
};

export default config;
