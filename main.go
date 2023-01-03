package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func prepareSwap(data []byte, blockNumber uint64, TxHash common.Hash) { fmt.Println("Swap") }
func prepareMint(data []byte, blockNumber uint64, TxHash common.Hash) { fmt.Println("Mint") }
func prepareBurn(data []byte, blockNumber uint64, TxHash common.Hash) { fmt.Println("Burn") }
func errHandler(err error) {
	fmt.Println("Error", err)
}

const ETH_USDT = "0x11b815efB8f581194ae79006d24E0d814B7697F6"
const ETH_USDC = "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"

var currentBlockNumber uint64

func main() {
	// createError()
	fmt.Println("Listener is running")
	client, err := ethclient.Dial(ALCHEMY_KEY)
	if err != nil {
		log.Fatal(err)
	}

	logSwapSig := []byte("Swap(address,address,int256,int256,uint160,uint128,int24)")
	logMintSig := []byte("Mint(address,address,int24,int24,uint128,uint256,uint256)")
	logBurnSig := []byte("Burn(address,int24,int24,uint128,uint256,uint256)")
	logSwapSigHash := crypto.Keccak256Hash(logSwapSig)
	logMintSigHash := crypto.Keccak256Hash(logMintSig)
	logBurnSigHash := crypto.Keccak256Hash(logBurnSig)

	contractAddress := common.HexToAddress(ETH_USDT)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logSwapSigHash.Hex():
			prepareSwap(vLog.Data, vLog.BlockNumber, vLog.TxHash)
		case logMintSigHash.Hex():
			prepareMint(vLog.Data, vLog.BlockNumber, vLog.TxHash)
		case logBurnSigHash.Hex():
			prepareBurn(vLog.Data, vLog.BlockNumber, vLog.TxHash)
		}
		currentBlockNumber = vLog.BlockNumber
	}
	for {
		select {
		case err := <-sub.Err():
			errHandler(err)
		case vLog := <-logs:
			// fmt.Println(vLog) // pointer to event log
			fmt.Println(vLog.BlockNumber)
		}
	}
}
