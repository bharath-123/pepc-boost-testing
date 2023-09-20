package util

import (
	"pepc-tester/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetAtomicSwap(address string, client *ethclient.Client) (*contracts.AtomicSwap, error) {
	return contracts.NewAtomicSwap(common.HexToAddress(address), client)
}

func GetDai(address string, client *ethclient.Client) (*contracts.Dai, error) {
	return contracts.NewDai(common.HexToAddress(address), client)
}

func GetErc20(address string, client *ethclient.Client) (*contracts.Erc20, error) {
	return contracts.NewErc20(common.HexToAddress(address), client)
}

func GetUniswapFactory(address string, client *ethclient.Client) (*contracts.UniswapFactory, error) {
	return contracts.NewUniswapFactory(common.HexToAddress(address), client)
}

func GetWeth(address string, client *ethclient.Client) (*contracts.Weth, error) {
	return contracts.NewWeth(common.HexToAddress(address), client)
}

func GetUniswapSwap(address string, client *ethclient.Client) (*contracts.UniswapPair, error) {
	return contracts.NewUniswapPair(common.HexToAddress(address), client)
}
