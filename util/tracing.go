package util

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"pepc-tester/constants"
	types2 "pepc-tester/types"

	"github.com/ethereum/go-ethereum/core/types"
)

func GetCallTraces(tx *types.Transaction) (*types2.CallTraceResponse, error) {
	ecRpcUrl := fmt.Sprintf("http://%s", constants.BuilderUrl)

	signer := types.NewCancunSigner(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("tx bytes are %v\n", hex.EncodeToString(tx.Data()))
	// Create the JSON-RPC request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "debug_traceCall",
		"params": []interface{}{map[string]interface{}{
			"from":  sender,
			"to":    tx.To(),
			"gas":   fmt.Sprintf("0x%x", tx.Gas()),
			"data":  fmt.Sprintf("0x%s", hex.EncodeToString(tx.Data())),
			"value": fmt.Sprintf("0x%x", tx.Value()),
		}, "latest", map[string]interface{}{"tracer": "callTracer", "disableStorage": false, "disableMemory": false}},
		"id": 1,
	}

	// Serialize the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to serialize JSON request: %v\n", err)
		return nil, err
	}

	// Send the HTTP POST request to the Ethereum client
	resp, err := http.Post(ecRpcUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Failed to send HTTP request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	//// Read the response body
	//responseBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("Failed to read response body: %v\n", err)
	//	return nil, err
	//}
	//
	//fmt.Printf("tracer response is")
	//fmt.Println(string(responseBody))
	//
	//return nil, nil
	// Parse the JSON response
	var jsonResponse types2.CallTraceResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		fmt.Printf("Failed to parse JSON response: %v\n", err)
		return nil, err
	}

	// Print the response
	return &jsonResponse, nil
}
