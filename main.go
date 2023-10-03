package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"pepc-tester/constants"
	"pepc-tester/util"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	for {
		// sleep for 1 slot
		time.Sleep(12 * time.Second)

		fmt.Println("2. Creating private key")
		privateKey, err := crypto.HexToECDSA("<your private key>")
		if err != nil {
			log.Print(err)
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
			log.Print(err)
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
			log.Print(err)
			continue
		}

		fmt.Println("6. Creating nonce")
		nonce, err := client.NonceAt(context.Background(), sender, nil)
		if err != nil {
			log.Print(err)
			continue
		}

		transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			log.Print(err)
			continue
		}
		transactor.Context = context.Background()
		transactor.NoSend = false

		// Set the gas price and gas limit
		fmt.Println("7. Setting gas price and gas limit")
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Print(err)
			continue
		}

		gasLimit := uint64(21000) // You may need to adjust this depending on the type of transaction

		// Specify the amount to send (in Wei)
		value := big.NewInt(1000000000000000000) // 0.1 ETH

		// Create a new Ethereum transaction
		fmt.Printf("8. Creating the txs")
		//tobTx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
		builderPayout := types.NewTransaction(nonce, common.HexToAddress("<Builder coinbase address>"), value, gasLimit, gasPrice, nil)

		fmt.Println("9. get local chain id")
		localChainId, err := client.ChainID(context.Background())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("ChainID: %d\n", localChainId)

		// Sign the transaction with the sender's private key
		signedBuilderPayout, err := types.SignTx(builderPayout, types.NewCancunSigner(localChainId), privateKey)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("Sending tx with hash: %s\n", signedBuilderPayout.Hash())
		err = client.SendTransaction(context.Background(), signedBuilderPayout)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("Sent tx with hash: %s\n", signedBuilderPayout.Hash())
	}
}
