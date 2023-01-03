package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func recovery(startBlock uint64, currentBlock uint64, client *ethclient.Client) {
	fmt.Println("Recovery module:\t", startBlock, currentBlock)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{ETH_USDT_Address, ETH_USDC_Address},
		Topics:    [][]common.Hash{{logSwapSigHash, logMintSigHash, logBurnSigHash}},
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(currentBlock)),
	}

	recoveryLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for vLog := range recoveryLogs {
		vLogPrepare(recoveryLogs[vLog])
	}
}
