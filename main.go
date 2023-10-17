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

func TransferWeth() {
	ecRpcUrl := fmt.Sprintf("http://%s", constants.BuilderUrl)
	fmt.Println("1. Creating EC client")
	client, err := ethclient.Dial(ecRpcUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("2. Creating private key")
	privateKey, err := crypto.HexToECDSA("ef5177cd0b6b21c87db5a0bf35d4084a8a57a9d6a064f86d51ac85f2b873a4e2")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("3. Creating sender instance")

	fmt.Println("2. Creating private key")
	receiverPrivateKey, err := crypto.HexToECDSA("df9bb6de5d3dc59595bcaa676397d837ff49441d211878c024eabda2cd067c9f")
	if err != nil {
		log.Fatal(err)
		return
	}

	receiver := crypto.PubkeyToAddress(receiverPrivateKey.PublicKey)

	wethContract, err := util.GetWeth(constants.WethAddress.String(), client)
	if err != nil {
		log.Fatal(err)
		return
	}

	chainId, err := client.ChainID(context.Background())
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

	res, err := wethContract.Transfer(transactor, receiver, big.NewInt(1000000000000000000))
	if err != nil {
		log.Fatal(err)
		return
	}

	transactor, err = bind.NewKeyedTransactorWithChainID(receiverPrivateKey, chainId)

	fmt.Printf("Approving wEth")
	res, err = wethContract.Approve(transactor, constants.AtomicSwap, big.NewInt(1000000000000000000))
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Successfully Transfereed weth: %v\n", res.Hash().String())

}

func main() {
	//TransferWeth()

	for {
		// sleep for 1 slot

		//fmt.Printf("Fetching the current slot from the relayer\n")
		//currentSlot, err := util.GetCurrentSlot(fmt.Sprintf("http://%s/eth/v1/relay/get_head_slot", constants.MevRelayerUrl))
		//if err != nil {
		//	fmt.Printf("Error fetching head slot from relayer: %v", err)
		//	continue
		//}
		//fmt.Printf("Current slot: %v\n", currentSlot)
		//
		//fmt.Printf("Fetching the parent hash for the next slot from the relayer\n")
		//slotParentHash, err := util.GetCurrentBlock(fmt.Sprintf("http://%s/eth/v1/relay/get_parent_hash_for_slot", constants.MevRelayerUrl), currentSlot.Uint64()+1)
		//if err != nil {
		//	fmt.Printf("Error fetching parent hash for slot relayer: %v", err)
		//	continue
		//}
		//fmt.Println("Parent hash: ", slotParentHash)
		//
		//fmt.Printf("Fetching the fee recipient of the proposer the next slot from the relayer\n")
		//proposerFeeRecipient, err := util.GetCurrentProposerFeeRecipient(fmt.Sprintf("http://%s/eth/v1/relay/get_proposer_for_slot", constants.MevRelayerUrl), currentSlot.Uint64()+1)
		//if err != nil {
		//	fmt.Printf("Error fetching parent hash for slot relayer: %v", err)
		//	continue
		//}
		//fmt.Println("proposer fee recipient: ", proposerFeeRecipient)

		fmt.Println("2. Creating private key")
		privateKey, err := crypto.HexToECDSA("cc877501b98c7171da436b1f1b6081941495795af765b5341e32558acd34e722")
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Println("3. Creating sender instance")
		// Create a new instance of the Ethereum sender's wallet
		sender := crypto.PubkeyToAddress(privateKey.PublicKey)
		fmt.Printf("Sender address: %v\n", sender.String())

		ecRpcUrl := fmt.Sprintf("http://%s", constants.BuilderUrl)
		fmt.Println("1. Creating EC client")
		client, err := ethclient.Dial(ecRpcUrl)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("client: %v\n", client)

		senderBalance, err := util.GetBalance(client, sender.String())
		if err != nil {
			fmt.Printf("Error fetching balance for sender: %v", err)
			continue
		}
		fmt.Printf("Sender balance: %v\n", senderBalance)

		chainId, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal(err)
			continue
		}

		uniV3SwapRouter, err := util.GetUniV3SwapRouter(constants.UniV3SwapRouter.String(), client)
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Println("6. Creating nonce")
		nonce, err := client.NonceAt(context.Background(), sender, nil)
		if err != nil {
			log.Fatal(err)
			continue
		}

		transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			log.Fatal(err)
			continue
		}
		transactor.Context = context.Background()
		transactor.NoSend = false
		transactor.Nonce = big.NewInt(int64(nonce))

		fmt.Printf("Creating swap tx\n")
		transactor.NoSend = true
		fmt.Printf("Creating swap tx.....\n")
		swapTx1, err := uniV3SwapRouter.ExactInputSingle(transactor, contracts.ISwapRouterExactInputSingleParams{
			TokenIn:           constants.WethGoerliAddress,
			TokenOut:          constants.UsdcAddress,
			Fee:               big.NewInt(3000),
			Recipient:         sender,
			Deadline:          big.NewInt(time.Now().Unix() + 1200),
			AmountIn:          big.NewInt(10000),
			AmountOutMinimum:  big.NewInt(0),
			SqrtPriceLimitX96: big.NewInt(0),
		})
		//transactor.Nonce = big.NewInt(0).Add(big.NewInt(int64(nonce)), big.NewInt(1))
		//swapTx2, err := uniV3SwapRouter.ExactInputSingle(transactor, contracts.ISwapRouterExactInputSingleParams{
		//	TokenIn:           constants.WethGoerliAddress,
		//	TokenOut:          constants.WBtcGoerliAddress,
		//	Fee:               big.NewInt(3000),
		//	Recipient:         sender,
		//	Deadline:          big.NewInt(time.Now().Unix() + 1200),
		//	AmountIn:          big.NewInt(1),
		//	AmountOutMinimum:  big.NewInt(0),
		//	SqrtPriceLimitX96: big.NewInt(0),
		//})
		//
		//transactor.Nonce = big.NewInt(0).Add(big.NewInt(int64(nonce)), big.NewInt(2))
		//swapTx3, err := uniV3SwapRouter.ExactInputSingle(transactor, contracts.ISwapRouterExactInputSingleParams{
		//	TokenIn:           constants.WethGoerliAddress,
		//	TokenOut:          constants.DaiAddress,
		//	Fee:               big.NewInt(3000),
		//	Recipient:         sender,
		//	Deadline:          big.NewInt(time.Now().Unix() + 1200),
		//	AmountIn:          big.NewInt(1),
		//	AmountOutMinimum:  big.NewInt(0),
		//	SqrtPriceLimitX96: big.NewInt(0),
		//})

		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("Created Eth/Usdc tx!!\n", swapTx1.Hash())

		time.Sleep(12 * time.Second)

		//fmt.Printf("Created Eth/WBtc tx!!\n", swapTx2.Hash())
		//fmt.Printf("Created Eth/Dai tx!!\n", swapTx3.Hash())

		//// Replace with the recipient's address
		//fmt.Println("5. Creating relayer instance")
		//proposerFeeRecipientAddress := common.HexToAddress(proposerFeeRecipient)
		//
		//proposerFeeRecipientBalance, err := util.GetBalance(client, proposerFeeRecipientAddress.String())
		//if err != nil {
		//	fmt.Printf("Error fetching balance for proposer fee recipient: %v", err)
		//	continue
		//}
		//fmt.Printf("proposer fee recipient address balance: %v\n", proposerFeeRecipientBalance)
		//
		//fmt.Printf("Nonce: %d\n", nonce)
		//fmt.Printf("Eth/Usdc Swap tx nonce is: %d\n", swapTx1.Nonce())
		//fmt.Printf("Eth/WBtc Swap tx nonce is: %d\n", swapTx2.Nonce())
		//fmt.Printf("Eth/Dai Swap tx nonce is: %d\n", swapTx3.Nonce())
		//
		//// Set the gas price and gas limit
		//fmt.Println("7. Setting gas price and gas limit")
		//gasPrice, err := client.SuggestGasPrice(context.Background())
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//
		//gasLimit := uint64(21000) // You may need to adjust this depending on the type of transaction
		//
		//// Specify the amount to send (in Wei)
		//value := big.NewInt(100000000000000000) // 0.1 ETH
		//
		//// Create a new Ethereum transaction
		//fmt.Printf("8. Creating the txs")
		////tobTx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
		//proposerPayout := types.NewTransaction(swapTx3.Nonce()+1, proposerFeeRecipientAddress, value, gasLimit, gasPrice, nil)
		//
		//fmt.Println("9. get local chain id")
		//localChainId, err := client.ChainID(context.Background())
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//fmt.Printf("ChainID: %d\n", localChainId)
		//
		//// Sign the transaction with the sender's private key
		//signedProposerPayout, err := types.SignTx(proposerPayout, types.NewCancunSigner(localChainId), privateKey)
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//binaryTx1, err := swapTx1.MarshalBinary()
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//binaryTx2, err := swapTx2.MarshalBinary()
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//binaryTx3, err := swapTx3.MarshalBinary()
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//binaryTx4, err := signedProposerPayout.MarshalBinary()
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//
		//tobTxs := new(types2.TobTxsSubmitRequest)
		//tobTxs.TobTxs = utilbellatrix.ExecutionPayloadTransactions{
		//	Transactions: []bellatrix.Transaction{binaryTx1, binaryTx2, binaryTx3, binaryTx4},
		//}
		//tobTxs.Slot = currentSlot.Uint64() + 1
		//tobTxs.ParentHash = slotParentHash
		//
		//jsonTobTxs, err := tobTxs.MarshalJSON()
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//
		//fmt.Println("10. Sending transactions to mev relayer")
		//// Create a new HTTP POST request with the JSON data as the request body
		//req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/relay/v1/builder/tob_txs", constants.MevRelayerUrl), bytes.NewBuffer(jsonTobTxs))
		//if err != nil {
		//	fmt.Println("Error creating request:", err)
		//	continue
		//}
		//req.Header.Set("Content-Type", "application/json")
		//
		//httpClient := &http.Client{}
		//
		//// Send the request
		//resp, err := httpClient.Do(req)
		//if err != nil {
		//	fmt.Println("Error sending request:", err)
		//	continue
		//}
		//defer resp.Body.Close()
		//
		//// Check the response status code
		//if resp.StatusCode == http.StatusOK {
		//	fmt.Printf("Tob Txs successfully submitted for slot: %d and parentHash: %s\n", tobTxs.Slot, tobTxs.ParentHash)
		//	fmt.Printf("First ToB tx hash is %s\n", swapTx1.Hash().String())
		//	fmt.Printf("Second ToB tx hash is %s\n", swapTx2.Hash().String())
		//	fmt.Printf("Third ToB tx hash is %s\n", swapTx3.Hash().String())
		//	fmt.Printf("Relayer payout hash is %s\n", signedProposerPayout.Hash().String())
		//} else {
		//	fmt.Println("Request failed with status code:", resp.StatusCode)
		//	// read the body
		//	body, err := io.ReadAll(resp.Body)
		//	if err != nil {
		//		fmt.Println("Error reading response body:", err)
		//		continue
		//	}
		//	fmt.Printf("Request failed with error message: %s\n", string(body))
		//
		//}
	}
}
