package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"pepc-tester/constants"
	"pepc-tester/contracts"
	"pepc-tester/util"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// Test private key, no problem in pushing this to github
	privateKey, err := crypto.HexToECDSA("cc877501b98c7171da436b1f1b6081941495795af765b5341e32558acd34e722")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new instance of the Ethereum sender's wallet
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("Sender is %v\n", sender)

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

	fmt.Printf("ChainId is %v\n", chainId)

	uniV3SwapRouter, err := util.GetUniV3SwapRouter(constants.UniV3SwapRouter.String(), client)
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
	transactor.NoSend = false
	transactor.GasLimit = 2300000
	transactor.GasFeeCap = big.NewInt(180)
	transactor.GasTipCap = big.NewInt(180)

	fmt.Printf("Approving weth to uniswap v3 router.....\n")
	wethToken, err := util.GetWeth(constants.WethGoerliAddress.String(), client)
	if err != nil {
		log.Fatal(err)
		return
	}
	approveTx, err := wethToken.Approve(transactor, constants.UniV3SwapRouter, big.NewInt(1000000000000000000))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Approved weth to uniswap v3 router!: %s\n", approveTx.Hash().String())

	//uniV3Quoter, err := util.GetUniV3Quoter(constants.UniV3Quoter.String(), client)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Quoting exact input single.....\n")
	//res, err := uniV3Quoter.QuoteExactInputSingle(transactor,
	//	constants.WethGoerliAddress,
	//	constants.UsdcAddress,
	//	big.NewInt(1000000000000000000),
	//	big.NewInt(100),
	//	big.NewInt(0),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Quoted exact input single!: %v\n", res)
	//trace, err := util.GetCallTraces(res)
	//if err != nil {
	//	fmt.Printf("Failed to get call traces: %v\n", err)
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Trace is %v\n", trace.Result)
	//
	//uniV3QuoterAbi, err := contracts.UniswapV3QuoterMetaData.GetAbi()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	////callTraces, err := util.GetCallTraces(res)
	//outputArgs, err := util.GetOutputMethodArgs(trace.Result.Output, "quoteExactInputSingle", uniV3QuoterAbi)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Output args are %v\n", outputArgs)

	//fmt.Printf("Quoted exact input single!: %v\n", res)

	transactor.NoSend = true
	fmt.Printf("Creating swap tx.....\n")
	swapTx, err := uniV3SwapRouter.ExactInputSingle(transactor, contracts.ISwapRouterExactInputSingleParams{
		TokenIn:           constants.WethGoerliAddress,
		TokenOut:          constants.UsdcAddress,
		Fee:               big.NewInt(3000),
		Recipient:         sender,
		Deadline:          big.NewInt(time.Now().Unix() + 1200),
		AmountIn:          big.NewInt(1),
		AmountOutMinimum:  big.NewInt(0),
		SqrtPriceLimitX96: big.NewInt(0),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	trace, err := util.GetCallTraces(swapTx)
	if err != nil {
		fmt.Printf("Failed to get call traces: %v\n", err)
		log.Fatal(err)
		return
	}
	res, err := util.IsTxUniv3EthUsdcSwap(&trace.Result)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("IsTxUniv3EthUsdcSwap is %v\n", res)
	//swapTx, err := uniV3SwapRouter.RefundETH(transactor)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//jsonTx, err := swapTx.MarshalJSON()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//// write jsonTx to a file tx1.json
	//err = os.WriteFile("invalid_tx1.json", jsonTx, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//fmt.Printf("Created swap tx!!\n", swapTx.Hash())

	//fmt.Printf("Trace is %v\n", trace.Result)

	//atomicSwapContract, err := util.GetAtomicSwap(constants.AtomicSwap.String(), client)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//transactor.Context = context.Background()
	//transactor.NoSend = true
	//
	//fmt.Printf("Creating swao tx.....\n")
	//swapTx, err := atomicSwapContract.Swap(transactor, []common.Address{constants.WethAddress, constants.DaiAddress}, big.NewInt(100000000000000000), constants.UniswapFactoryA, sender, false)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Created swap tx!!\n", swapTx.Hash())
	//
	//fmt.Printf("Fetching the traces of the swap tx.....\n")
	//traces, err := util.GetCallTraces(swapTx)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Fetched the traces of the swap tx!\n")
	//
	//fmt.Printf("Traces are %v\n", traces.Result)
	//
	//// check if tx contains a weth dai swap
	//isTxWethDaiSwap, err := util.IsTxWEthDaiSwap(&traces.Result)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//// if it does, then unpack the args
	//if isTxWethDaiSwap {
	//	fmt.Printf("Tx is a weth/dai swap!!\n")
	//	fmt.Println("Unpacking args now!")
	//	atomicSwapAbi, err := contracts.AtomicSwapMetaData.GetAbi()
	//	if err != nil {
	//		log.Fatal(err)
	//		return
	//	}
	//	args, err := util.GetMethodArgs(swapTx, "swap", atomicSwapAbi)
	//	if err != nil {
	//		log.Fatal(err)
	//		return
	//	}
	//	fmt.Printf("Args of swap Tx are %v\n", args)
	//} else {
	//	fmt.Printf("Tx is not a weth/dai swap!!\n")
	//}
}
