package util

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

func GetCurrentSlot(mevRelayUrl string) (*big.Int, error) {
	fmt.Printf("Mev url is %s\n", mevRelayUrl)
	// Create a new HTTP POST request with the JSON data as the request body
	req, err := http.NewRequest("GET", mevRelayUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	// Send the request
	fmt.Printf("Sending request")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}

	// read big.Int type from response body
	var currentSlot big.Int
	err = json.NewDecoder(resp.Body).Decode(&currentSlot)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	return &currentSlot, nil
}

func GetCurrentBlock(mevRelayUrl string, slot uint64) (string, error) {
	finalUrl := fmt.Sprintf("%s/%d", mevRelayUrl, slot)
	// Create a new HTTP POST request with the JSON data as the request body
	req, err := http.NewRequest("GET", finalUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	// Send the request
	fmt.Print("Sending request")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}

	// print body as string
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body:", err)
	//	return "", err
	//}
	//fmt.Println(string(body))
	//
	//return "nil", fmt.Errorf("not implemented")

	var parentHash string
	err = json.NewDecoder(resp.Body).Decode(&parentHash)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}

	return parentHash, nil
}

func GetCurrentProposerFeeRecipient(mevRelayUrl string, slot uint64) (string, error) {
	finalUrl := fmt.Sprintf("%s/%d", mevRelayUrl, slot)
	// Create a new HTTP POST request with the JSON data as the request body
	req, err := http.NewRequest("GET", finalUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	// Send the request
	fmt.Print("Sending request")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Request failed with status code:", resp.StatusCode)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body:", err)
	//	return "", err
	//}
	//fmt.Println(string(body))
	//return "", err

	var feeRecipient string
	err = json.NewDecoder(resp.Body).Decode(&feeRecipient)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}

	return feeRecipient, nil
}

func GetTobTxs(mevRelayerUrl string, slot uint64, parentHash string) ([]*types.Transaction, error) {
	trimmedParentHash := strings.Trim(parentHash, "\"")
	trimmedParentHash = strings.Trim(trimmedParentHash, "\n")
	finalUrl := fmt.Sprintf("%s/%d/%s", mevRelayerUrl, slot, trimmedParentHash)
	// Create a new HTTP POST request with the JSON data as the request body
	req, err := http.NewRequest("GET", finalUrl, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	// Send the request
	fmt.Print("Sending request")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		return nil, fmt.Errorf("Request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	txBytes := string(body)
	txs := strings.Split(txBytes, ",")
	finalTxs := []*types.Transaction{}

	for _, tx := range txs {
		tobTx, err := hex.DecodeString(tx)
		if err != nil {
			return nil, err
		}
		decodedTx := new(types.Transaction)
		err = decodedTx.UnmarshalBinary(tobTx)
		if err != nil {
			return nil, err
		}
		finalTxs = append(finalTxs, decodedTx)
	}

	return finalTxs, nil
}
