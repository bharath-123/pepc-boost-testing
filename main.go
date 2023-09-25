package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"pepc-tester/constants"
	"pepc-tester/util"

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
	//
	//fmt.Printf("Fetching the current slot from the relayer\n")
	//currentSlot, err := util.GetCurrentSlot(fmt.Sprintf("http://%s/eth/v1/relay/get_head_slot", constants.MevRelayerUrl))
	//if err != nil {
	//	fmt.Printf("Error fetching head slot from relayer: %v", err)
	//	return
	//}
	//fmt.Printf("Current slot: %v\n", currentSlot)
	//
	//fmt.Printf("Fetching the parent hash for the next slot from the relayer\n")
	//slotParentHash, err := util.GetCurrentBlock(fmt.Sprintf("http://%s/eth/v1/relay/get_parent_hash_for_slot", constants.MevRelayerUrl), currentSlot.Uint64()+1)
	//if err != nil {
	//	fmt.Printf("Error fetching parent hash for slot relayer: %v", err)
	//	return
	//}
	//fmt.Println("Parent hash: ", slotParentHash)
	//
	//fmt.Printf("Fetching the parent hash for the next slot from the relayer\n")
	//proposerFeeRecipient, err := util.GetCurrentProposerFeeRecipient(fmt.Sprintf("http://%s/eth/v1/relay/get_proposer_for_slot", constants.MevRelayerUrl), currentSlot.Uint64()+1)
	//if err != nil {
	//	fmt.Printf("Error fetching parent hash for slot relayer: %v", err)
	//	return
	//}
	//fmt.Println("proposer fee recipient: ", proposerFeeRecipient)
	//
	//fmt.Println("2. Creating private key")
	//privateKey, err := crypto.HexToECDSA("df9bb6de5d3dc59595bcaa676397d837ff49441d211878c024eabda2cd067c9f")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//fmt.Println("3. Creating sender instance")
	//// Create a new instance of the Ethereum sender's wallet
	//sender := crypto.PubkeyToAddress(privateKey.PublicKey)
	//
	//ecRpcUrl := fmt.Sprintf("http://%s", constants.BuilderUrl)
	//fmt.Println("1. Creating EC client")
	//client, err := ethclient.Dial(ecRpcUrl)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("client: %v\n", client)
	//
	//senderBalance, err := util.GetBalance(client, sender.String())
	//if err != nil {
	//	fmt.Printf("Error fetching balance for sender: %v", err)
	//	return
	//}
	//fmt.Printf("Sender balance: %v\n", senderBalance)
	//
	//chainId, err := client.ChainID(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//atomicSwapContract, err := util.GetAtomicSwap(constants.AtomicSwap.String(), client)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//fmt.Println("6. Creating nonce")
	//nonce, err := client.NonceAt(context.Background(), sender, nil)
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
	//transactor.Nonce = big.NewInt(int64(nonce))
	//
	//fmt.Printf("Creating swap tx\n")
	//swapTx, err := atomicSwapContract.Swap(transactor, []common.Address{constants.WethAddress, constants.DaiAddress}, big.NewInt(100000000000000000), constants.UniswapFactoryA, sender, false)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("Created tx!!\n", swapTx.Hash())
	//
	//// Replace with the recipient's address
	//fmt.Println("5. Creating relayer instance")
	//proposerFeeRecipientAddress := common.HexToAddress(proposerFeeRecipient)
	//
	//proposerFeeRecipientBalance, err := util.GetBalance(client, proposerFeeRecipientAddress.String())
	//if err != nil {
	//	fmt.Printf("Error fetching balance for proposer fee recipient: %v", err)
	//	return
	//}
	//fmt.Printf("proposer fee recipient address balance: %v\n", proposerFeeRecipientBalance)
	//
	//fmt.Printf("Nonce: %d\n", nonce)
	//fmt.Printf("Swap tx nonce is: %d\n", swapTx.Nonce())
	//
	//// Set the gas price and gas limit
	//fmt.Println("7. Setting gas price and gas limit")
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//gasLimit := uint64(21000) // You may need to adjust this depending on the type of transaction
	//
	//// Specify the amount to send (in Wei)
	//value := big.NewInt(1000000000000000000) // 0.1 ETH
	//
	//// Create a new Ethereum transaction
	//fmt.Printf("8. Creating the txs")
	////tobTx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	//proposerPayout := types.NewTransaction(swapTx.Nonce()+1, proposerFeeRecipientAddress, value, gasLimit, gasPrice, nil)
	//
	//fmt.Println("9. get local chain id")
	//localChainId, err := client.ChainID(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("ChainID: %d\n", localChainId)
	//
	//// Sign the transaction with the sender's private key
	//signedProposerPayout, err := types.SignTx(proposerPayout, types.NewCancunSigner(localChainId), privateKey)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//binaryTx1, err := swapTx.MarshalBinary()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//binaryTx2, err := signedProposerPayout.MarshalBinary()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//tobTxs := new(types2.TobTxsSubmitRequest)
	//tobTxs.TobTxs = utilbellatrix.ExecutionPayloadTransactions{
	//	Transactions: []bellatrix.Transaction{binaryTx1, binaryTx2},
	//}
	//tobTxs.Slot = currentSlot.Uint64() + 1
	//tobTxs.ParentHash = slotParentHash
	//
	//jsonTobTxs, err := tobTxs.MarshalJSON()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("10. Sending transactions to mev relayer")
	//// Create a new HTTP POST request with the JSON data as the request body
	//req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/relay/v1/builder/tob_txs", constants.MevRelayerUrl), bytes.NewBuffer(jsonTobTxs))
	//if err != nil {
	//	fmt.Println("Error creating request:", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//httpClient := &http.Client{}
	//
	//// Send the request
	//resp, err := httpClient.Do(req)
	//if err != nil {
	//	fmt.Println("Error sending request:", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//// Check the response status code
	//if resp.StatusCode == http.StatusOK {
	//	fmt.Printf("Tob Txs successfully submitted for slot: %d and parentHash: %s\n", tobTxs.Slot, tobTxs.ParentHash)
	//	fmt.Printf("ToB tx hash is %s\n", swapTx.Hash().String())
	//	fmt.Printf("Relayer payout hash is %s\n", signedProposerPayout.Hash().String())
	//} else {
	//	fmt.Println("Request failed with status code:", resp.StatusCode)
	//	fmt.Printf("")
	//	return
	//}
	//
	//fmt.Printf("Getting ToB txs!\n")
	//res, err := util.GetTobTxs(fmt.Sprintf("http://%s/eth/v1/relay/get_tob_txs", constants.MevRelayerUrl), currentSlot.Uint64()+1, slotParentHash)
	//for _, r := range res {
	//	fmt.Printf("tx hash stored in relayer is %s\n", r.Hash().String())
	//}
	//fmt.Printf("Successfully got ToB txs\n")

	//======
	//Block inspection
	//======
	ecRpcUrl := fmt.Sprintf("http://%s", constants.BuilderUrl)
	fmt.Println("1. Creating EC client")
	client, err := ethclient.Dial(ecRpcUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("client: %v\n", client)

	blockDetails, err := util.GetBlockByBlockHash(client, "0x75cd835787486af61df75d96fe82fba7a2c4a47d03cc822b02171e6152d7fd8c")
	if err != nil {
		fmt.Printf("Error fetching block by hash: %v", err)
		return
	}
	fmt.Printf("Block number: %d\n", blockDetails.Number().Uint64())
	fmt.Printf("Parent hash: %v\n", blockDetails.ParentHash())
	i := 0
	fmt.Printf("Checking txs\n")
	for _, tx := range blockDetails.Transactions() {
		fmt.Printf("=======\n")
		fmt.Printf("To: %v\n", tx.To())
		fmt.Printf("Value: %v\n", tx.Value())
		fmt.Printf("Cost: %v\n", tx.Cost())
		fmt.Printf("Gas: %v\n", tx.Gas())
		fmt.Printf("Nonce: %v\n", tx.Nonce())
		fmt.Printf("Hash is %s\n", tx.Hash().String())
		receipt, err := util.GetTxReceipt(client, tx.Hash().String())
		if err != nil {
			return
		}
		fmt.Printf("Tx status is %d\n", receipt.Status)
		//receipt.
		//fmt.Printf("Tx status is %d\n", receipt.Status)
		//receipt.
		i++
		if i > 2 {
			break
		}
	}

}

//res, err := util.GetCallTraces(swapTx)
//if err != nil {
//	log.Fatal(err)
//	return
//}
//
//fmt.Printf("Res is %v\n", res.Result)

//
//atomicSwapMetadata, err := contracts.AtomicSwapMetaData.GetAbi()
//if err != nil {
//	log.Fatal(err)
//	return
//}
//swapId := atomicSwapMetadata.Methods["swap"].ID
//fmt.Printf("id: %x\n", swapId)
//
//res := atomicSwapMetadata.Methods["swap"]
//argsMap := make(map[string]interface{})
//err = res.Inputs.UnpackIntoMap(argsMap, swapTx.Data())
//if err != nil {
//	log.Fatal(err)
//	return
//}
//fmt.Printf("res is %v\n", argsMap)
