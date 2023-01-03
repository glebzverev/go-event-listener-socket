package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
