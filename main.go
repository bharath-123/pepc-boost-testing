package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	utilbellatrix "github.com/attestantio/go-eth2-client/util/bellatrix"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TobTxsSubmitRequest struct {
	TobTxs     utilbellatrix.ExecutionPayloadTransactions
	Slot       uint64
	ParentHash string
}

func (t *TobTxsSubmitRequest) MarshalJSON() ([]byte, error) {
	txBytes, err := t.TobTxs.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		TobTxs     []byte `json:"tobTxs"`
		Slot       uint64 `json:"slot"`
		ParentHash string `json:"parentHash"`
	}{
		TobTxs:     txBytes,
		Slot:       t.Slot,
		ParentHash: t.ParentHash,
	})
}

func (t *TobTxsSubmitRequest) UnmarshalJSON(data []byte) error {
	var intermediateJson struct {
		TobTxs     []byte `json:"tobTxs"`
		Slot       uint64 `json:"slot"`
		ParentHash string `json:"parentHash"`
	}
	err := json.Unmarshal(data, &intermediateJson)
	if err != nil {
		return err
	}

	err = t.TobTxs.UnmarshalSSZ(intermediateJson.TobTxs)
	if err != nil {
		return err
	}
	t.Slot = intermediateJson.Slot
	t.ParentHash = intermediateJson.ParentHash

	return nil
}

func main() {
	const rpcUrl = ""
	const mevRelayerApi = ""
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with your sender's private key
	privateKey, err := crypto.HexToECDSA("0xb3c409b6b0b3aa5e65ab2dc1930534608239a478106acf6f3d9178e9f9b00b35")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the Ethereum sender's wallet
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Replace with the recipient's address
	toAddress := common.HexToAddress("0xB9D7a3554F221B34f49d7d3C61375E603aFb699e")
	relayerAddress := common.HexToAddress("0x4E9A3d9D1cd2A2b2371b8b3F489aE72259886f1A")

	// Create a new nonce for the sender
	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}

	// Set the gas price and gas limit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000) // You may need to adjust this depending on the type of transaction

	// Specify the amount to send (in Wei)
	value := big.NewInt(100000000000000000) // 0.1 ETH

	// Create a new Ethereum transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	relayerPayout := types.NewTransaction(nonce+1, relayerAddress, value, gasLimit, gasPrice, nil)
	// Sign the transaction with the sender's private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signedRelayerPayout, err := types.SignTx(relayerPayout, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	binaryTx1, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	binaryTx2, err := signedRelayerPayout.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	tobTxs := new(TobTxsSubmitRequest)
	tobTxs.TobTxs = utilbellatrix.ExecutionPayloadTransactions{
		Transactions: []bellatrix.Transaction{binaryTx1, binaryTx2},
	}
	tobTxs.Slot = 100
	tobTxs.ParentHash = "0x0000000"

	jsonTobTxs, err := tobTxs.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP POST request with the JSON data as the request body
	req, err := http.NewRequest("POST", mevRelayerApi, bytes.NewBuffer(jsonTobTxs))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}
}
