package util

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"pepc-tester/constants"
	"pepc-tester/contracts"
	types2 "pepc-tester/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func GetCallTraces(tx *types.Transaction) (*types2.CallTraceResponse, error) {
	ecRpcUrl := fmt.Sprintf("http://%s", constants.EcUrl)

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

// just check if it goes to the DaiWethPair with a swap tx
func IsTxWEthDaiSwap(trace *types2.CallTrace) (bool, error) {
	stack := []types2.CallTrace{*trace}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		res, err := IsTraceToWEthDaiPair(current)
		if err != nil {
			return false, err
		}
		// we found a weth/dai swap i.e the tx contains a weth/dai swap
		if res {
			return true, nil
		}

		for _, call := range current.Calls {
			stack = append(stack, call)
		}
	}

	return false, nil
}

func GetMethodArgs(tx *types.Transaction, method string, contractAbi *abi.ABI) (map[string]interface{}, error) {
	res := contractAbi.Methods[method]
	argsMap := make(map[string]interface{})
	err := res.Inputs.UnpackIntoMap(argsMap, tx.Data()[4:])
	if err != nil {
		return nil, err
	}

	return argsMap, nil
}

// This will change based on the state interference check
func IsTraceToWEthDaiPair(callTrace types2.CallTrace) (bool, error) {
	if callTrace.To == nil {
		return false, nil
	}
	if callTrace.Type == "STATICCALL" {
		return false, nil
	}

	uniswapDaiWethAddress1 := common.HexToAddress("0x0D6b80a9Cefc2C58308F0Adc26586E550E4422ef")
	uniswapDaiWethAddress2 := common.HexToAddress("0x2ed2B47342450C006F83913a422F7C2BDAB8377a")
	if *callTrace.To != uniswapDaiWethAddress1 && *callTrace.To != uniswapDaiWethAddress2 {
		return false, nil
	}

	if len(callTrace.Input) < 4 {
		return false, nil
	}

	// this will be the same across all environments
	uniswapPairAbi, err := contracts.UniswapPairMetaData.GetAbi()
	if err != nil {
		return false, err
	}
	swapId := uniswapPairAbi.Methods["swap"].ID
	if !bytes.Equal(callTrace.Input[:4], swapId) {
		return false, nil
	}

	return true, nil
}
