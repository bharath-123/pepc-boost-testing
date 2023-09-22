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

	//// print out the resp body in string
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//fmt.Printf("resp body is %v\n", buf.String())
	//
	//return nil, fmt.Errorf("Printing as string")

	//// Parse the JSON response
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
func IsTxUniv3EthUsdcSwap(trace *types2.CallTrace) (bool, error) {
	stack := []types2.CallTrace{*trace}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		res, err := IsTraceUniV3EthUsdcSwap(current)
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

func GetInterfaceMethodArgs(data []byte, method string, contractAbi *abi.ABI) ([]interface{}, error) {
	abiMethod := contractAbi.Methods[method]
	res, err := abiMethod.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, err
	}

	return res, nil
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

func GetMethodArgs(data []byte, method string, contractAbi *abi.ABI) (interface{}, error) {
	abiMethod := contractAbi.Methods[method]
	res, err := abiMethod.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetOutputMethodArgs(data []byte, method string, contractAbi *abi.ABI) (interface{}, error) {
	abiMethod := contractAbi.Methods[method]
	res, err := abiMethod.Outputs.Unpack(data[4:])
	if err != nil {
		return nil, err
	}

	return res, nil
}

// This will change based on the state interference check
func IsTraceUniV3EthUsdcSwap(callTrace types2.CallTrace) (bool, error) {
	if callTrace.To == nil {
		return false, nil
	}
	if callTrace.Type == "STATICCALL" {
		return false, nil
	}
	if len(callTrace.Input) < 4 {
		return false, nil
	}

	if *callTrace.To != constants.UniV3SwapRouter {
		return false, nil
	}

	uniV3SwapRouterAbi, err := contracts.UniswapV3SwapRouterMetaData.GetAbi()
	if err != nil {
		return false, err
	}
	exactInputSingleId := uniV3SwapRouterAbi.Methods["exactInputSingle"].ID
	if !bytes.Equal(callTrace.Input[:4], exactInputSingleId) {
		return false, nil
	}

	// unpack the args
	args, err := GetInterfaceMethodArgs(callTrace.Input, "exactInputSingle", uniV3SwapRouterAbi)
	if err != nil {
		return false, err
	}
	firstEle := args[0]
	fmt.Printf("args are %v\n", args[0])
	firstEleBytes, err := json.Marshal(firstEle)
	fmt.Printf("firstEleBytes are %v\n", string(firstEleBytes))
	exactInputSingleParams := new(contracts.ISwapRouterExactInputSingleParams)
	err = json.Unmarshal(firstEleBytes, exactInputSingleParams)
	if err != nil {
		return false, err
	}
	fmt.Printf("exactInputSingleParams are %v\n", exactInputSingleParams.AmountIn)
	//swapRouterParams, ok := firstEle.(contracts.ISwapRouterExactInputSingleParams)
	//if !ok {
	//	return false, fmt.Errorf("failed to parse args of swapRouter tx to ISwapRouterExactInputSingleParams")
	//}
	//fmt.Printf("swapRouterParams are %v\n", swapRouterParams)
	return false, nil
	//if swapRouterParams.TokenIn != constants.WethGoerliAddress || swapRouterParams.TokenOut != constants.UsdcAddress {
	//	return false, nil
	//}
	//if swapRouterParams.TokenOut != constants.UsdcAddress || swapRouterParams.TokenOut != constants.WethGoerliAddress {
	//	return false, nil
	//}

	return true, nil
}
