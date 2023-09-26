# Checking tx data using geth style call traces

## Introduction

Let's say we have an ethereum tx and we want to figure out whether this tx performs a Weth/Dai swap on Uniswap V2. One way to achieve this is by
checking if the `To` address of the tx is to the Uniswap V2 Weth/Dai pool and then parsing the `Data` field of the tx to see if it matches the
signature of the `swap` function of the Uniswap V2 Weth/Dai pool contract. This works for simple txs made by EOA accounts, but what if the tx is
made by a contract? In this case, the `To` address of the tx will be to the contract address, and the `Data` field of the tx will be the ABI
encoded function call to the contract. We cannot parse the `Data` field of the tx to see if it matches the signature of the `swap` function of
the Uniswap V2 Weth/Dai pool contract. 

To acheive this, we can use the `debug_traceCall` RPC call of geth to get the call traces of the tx. The call traces will contain the multiple traces of the
tx, and we can parse the `input` field of the traces to see if it matches the signature of the `swap` function of the Uniswap V2 Weth/Dai pool contract.
We can think of the calltracer breaking the tx down into the individual sub-txs performed. One of the individual sub-txs will be the `swap` function
call to the Uniswap V2 Weth/Dai pool contract.

## Running this example

This example uses a Kurtosis devnet setup. You can refer to https://github.com/kurtosis-tech/eth2-package to setup kurtosis and the devnet.

Once you have the devnet setup, please enter the EC url in the `constants/addresses.go` file. You can then run the example using the following command:

```bash
go run main.go
```

The files in `contracts` directory are generated using abigen(https://geth.ethereum.org/docs/tools/abigen). The `update-abis.sh` script has more details on
how the files were generated.

The abis in `abis` directory are picked up from https://github.com/flashbots/mev-flood. You can compile the contracts in the `mev-flood` repo by following the README file

The bulk of the logic is in `util/tracing.go`. The logic includes making the `debug_traceCall` and parsing the traces to see if the tx performs a Weth/Dai swap on Uniswap V2.
We use the abi generated from the `mev-flood` repo to parse the `input` field of the traces.
