package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	ETH_USDT_Address = common.HexToAddress("0x11b815efB8f581194ae79006d24E0d814B7697F6")
	ETH_USDC_Address = common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640")
	logSwapSigHash   = crypto.Keccak256Hash([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)"))
	logMintSigHash   = crypto.Keccak256Hash([]byte("Mint(address,address,int24,int24,uint128,uint256,uint256)"))
	logBurnSigHash   = crypto.Keccak256Hash([]byte("Burn(address,int24,int24,uint128,uint256,uint256)"))

	currentBlockNumber uint64
	startBlockNumber   uint64
	AFTER_RECONNECT    bool
	client             *ethclient.Client
	sub                ethereum.Subscription
	logs               (chan types.Log)
	err                error
	ALCHEMY_KEY        string
)

func createSubscription(client *ethclient.Client) (ethereum.Subscription, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ETH_USDT_Address, ETH_USDC_Address},
		Topics:    [][]common.Hash{{logSwapSigHash, logMintSigHash, logBurnSigHash}},
	}
	logs = make(chan types.Log)
	return client.SubscribeFilterLogs(context.Background(), query, logs)
}

func init() {
	// Init dotenv environmnet
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ALCHEMY_KEY = os.Getenv("ALCHEMY_KEY")
}

func main() {

Start:
	// Create new ethclient.Client
	client, err = ethclient.Dial(ALCHEMY_KEY)
	if err != nil {
		fmt.Println("Client Error:\t", err)
		fmt.Println("Try reconnect Client")
		time.Sleep(2 * time.Second)
		goto Start
	}

	// Create Filters
	fmt.Println("Listener is running")

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
		}
		if AFTER_RECONNECT {
			currentBlockNumber = vLog.BlockNumber
			AFTER_RECONNECT = false
			recovery(startBlockNumber, currentBlockNumber, client)
		}
		timerFunc()
		// timer = timer.NewTimer(time.Seconds * time.Duration(60*1))
	}

	for {
		select {
		// case <-timer.C:
		// fmt.Println("Alcehmy lost subscribe")
		// fmt.Println("Try reconnect Subscription. Last block:\t", currentBlockNumber)
		// AFTER_RECONNECT = true
		// startBlockNumber = currentBlockNumber
		// time.Sleep(2 * time.Second)
		// goto Subscribe
		case err := <-sub.Err():
			if err != nil {
				fmt.Errorf("Subscribe Error:\t", err)
				fmt.Println("Try reconnect Subscription. Last block:\t", currentBlockNumber)
				AFTER_RECONNECT = true
				startBlockNumber = currentBlockNumber
				time.Sleep(2 * time.Second)
				goto Subscribe
			}
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}
}
