package util

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBlockByBlockHash(client *ethclient.Client, blockHash string) (*types.Block, error) {
	hash := common.HexToHash(blockHash)
	return client.BlockByHash(context.Background(), hash)
}

func GetTransactionByHash(client *ethclient.Client, txHash string) (*types.Transaction, bool, error) {
	hash := common.BytesToHash([]byte(txHash))
	return client.TransactionByHash(context.Background(), hash)
}

func GetTxReceipt(client *ethclient.Client, txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	return client.TransactionReceipt(context.Background(), hash)
}

func GetBalance(client *ethclient.Client, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	return client.BalanceAt(context.Background(), account, nil)
}
