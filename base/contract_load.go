package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/store"
	"log"
)

func main() {
	/*
		加载智能合约
		可以通过地址加载智能合约
	*/

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x665Ab94CaB1c9fD296aAdCbfD9F012C9663F0594")

	instance, err := store.NewStore(address, client) // 加载节点上的智能合约
	if err != nil {
		fmt.Println("LOAD CONTRACT ERR")
	}

	_ = instance
}
