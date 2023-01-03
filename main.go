package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func prepareSwap(data []byte, blockNumber uint64, TxHash common.Hash) {
	fmt.Println("Swap", blockNumber, TxHash)
}
func prepareMint(data []byte, blockNumber uint64, TxHash common.Hash) {
	fmt.Println("Mint", blockNumber, TxHash)
}
func prepareBurn(data []byte, blockNumber uint64, TxHash common.Hash) {
	fmt.Println("Burn", blockNumber, TxHash)
}

func createSubscription(client *ethclient.Client) (ethereum.Subscription, error) {
	ETH_USDT_Address := common.HexToAddress(ETH_USDT)
	ETH_USDC_Address := common.HexToAddress(ETH_USDC)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ETH_USDT_Address, ETH_USDC_Address},
	}
	logs = make(chan types.Log)
	return client.SubscribeFilterLogs(context.Background(), query, logs)
}

// func errClientHandler(err error, sub ethereum.Subscription, client *ethclient.Client) {
// 	fmt.Println("ErrHandler Error:\t", err)
// 	fmt.Println("Try reconnect. Last block:\t", currentBlockNumber)
// 	time.Sleep(2 * time.Second)
// 	main()
// }

const ETH_USDT = "0x11b815efB8f581194ae79006d24E0d814B7697F6"
const ETH_USDC = "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"

var currentBlockNumber uint64
var client *ethclient.Client
var sub ethereum.Subscription
var logs (chan types.Log)
var err error

func init() {

}

func main() {
Start:
	ALCHEMY_KEY := "wss://eth-mainnet.g.alchemy.com/v2/dt-X9e68ahXP9Sl1bDn89XXHVT3vBohA"
	client, err = ethclient.Dial(ALCHEMY_KEY)
	if err != nil {
		fmt.Println("Client Error:\t", err)
		fmt.Println("Try reconnect Client")
		time.Sleep(2 * time.Second)
		goto Start
	}

	fmt.Println("Listener is running")
	logSwapSig := []byte("Swap(address,address,int256,int256,uint160,uint128,int24)")
	logMintSig := []byte("Mint(address,address,int24,int24,uint128,uint256,uint256)")
	logBurnSig := []byte("Burn(address,int24,int24,uint128,uint256,uint256)")
	logSwapSigHash := crypto.Keccak256Hash(logSwapSig)
	logMintSigHash := crypto.Keccak256Hash(logMintSig)
	logBurnSigHash := crypto.Keccak256Hash(logBurnSig)

	logs = make(chan types.Log)
Subscribe:
	sub, err = createSubscription(client)
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
		default:
			sub.Unsubscribe()
			close(logs)
		}
		currentBlockNumber = vLog.BlockNumber
	}

	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				fmt.Errorf("Subscribe Error:\t", err)
				fmt.Println("Try reconnect Subscription. Last block:\t", currentBlockNumber)
				time.Sleep(2 * time.Second)
				goto Subscribe
			}
		case vLog := <-logs:
			fmt.Println(vLog)
		default:
			checkTime()
		}
	}
}
