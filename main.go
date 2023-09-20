package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"pepc-tester/constants"
	"pepc-tester/contracts"
	"pepc-tester/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// Test private key, no problem in pushing this to github
	privateKey, err := crypto.HexToECDSA("ef5177cd0b6b21c87db5a0bf35d4084a8a57a9d6a064f86d51ac85f2b873a4e2")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new instance of the Ethereum sender's wallet
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)

	ecRpcUrl := fmt.Sprintf("http://%s", constants.EcUrl)
	fmt.Println("Connecting to the execution client.....")
	client, err := ethclient.Dial(ecRpcUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to the execution client!")

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	atomicSwapContract, err := util.GetAtomicSwap(constants.AtomicSwap.String(), client)
	if err != nil {
		log.Fatal(err)
		return
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
		return
	}
	transactor.Context = context.Background()
	transactor.NoSend = true

	fmt.Printf("Creating swao tx.....\n")
	swapTx, err := atomicSwapContract.Swap(transactor, []common.Address{constants.WethAddress, constants.DaiAddress}, big.NewInt(100000000000000000), constants.UniswapFactoryA, sender, false)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Created swap tx!!\n", swapTx.Hash())

	fmt.Printf("Fetching the traces of the swap tx.....\n")
	traces, err := util.GetCallTraces(swapTx)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Fetched the traces of the swap tx!\n")

	fmt.Printf("Traces are %v\n", traces.Result)

	// check if tx contains a weth dai swap
	isTxWethDaiSwap, err := util.IsTxWEthDaiSwap(&traces.Result)
	if err != nil {
		log.Fatal(err)
		return
	}
	// if it does, then unpack the args
	if isTxWethDaiSwap {
		fmt.Printf("Tx is a weth/dai swap!!\n")
		fmt.Println("Unpacking args now!")
		atomicSwapAbi, err := contracts.AtomicSwapMetaData.GetAbi()
		if err != nil {
			log.Fatal(err)
			return
		}
		args, err := util.GetMethodArgs(swapTx, "swap", atomicSwapAbi)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Args of swap Tx are %v\n", args)
	} else {
		fmt.Printf("Tx is not a weth/dai swap!!\n")
	}
}
