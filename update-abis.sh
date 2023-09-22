# enter in stader-lib/contracts and throw and error if it fails
echo "Updating ABIs..."

cd contracts;

abigen --abi ./../abis/atomicSwap.json --pkg contracts --type AtomicSwap --out atomic-swap.go;
abigen --abi ./../abis/dai.json --pkg contracts --type Dai --out dai.go;
abigen --abi ./../abis/erc20.json --pkg contracts --type Erc20 --out erc-20.go;
abigen --abi ./../abis/UniswapFactory.json --pkg contracts --type UniswapFactory --out uniswap-factory.go;
abigen --abi ./../abis/UniswapPair.json --pkg contracts --type UniswapPair --out uniswap-pair.go;
abigen --abi ./../abis/weth.json --pkg contracts --type Weth --out weth.go;
abigen --abi ./../abis/UniswapV3SwapRouter.json --pkg contracts --type UniswapV3SwapRouter --out uniswap-v3-swap-router.go;
abigen --abi ./../abis/UniswapV3QuoterV2.json --pkg contracts --type UniswapV3QuoterV2 --out uniswap-v3-quoter-v2.go;
abigen --abi ./../abis/UniswapV3Quoter.json --pkg contracts --type UniswapV3Quoter --out uniswap-v3-quoter.go;


echo "Done updating ABIs."