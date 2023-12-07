package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/store"
	"log"
	"sync"
)

func main() {
	// 监听合约的日志

	client, err := ethclient.Dial("ws://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x665Ab94CaB1c9fD296aAdCbfD9F012C9663F0594")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		oldSubLog(client, contractAddress)
	}()

	go func() {
		defer wg.Done()
		newSubLog(client, contractAddress)
	}()
	wg.Wait()
}

func oldSubLog(client *ethclient.Client, contractAddress common.Address) {
	ctx := context.Background()
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress}, // 合约地址数组
	}

	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	// 加载ABI，为了解析监听返回的ABI日志数据
	contractAbi, err := store.StoreMetaData.GetAbi()
	if err != nil {
		fmt.Println("GET CONTRACT ABI ERR")
		log.Fatal(err)
	}

	fmt.Println("START SUBSCRIBE")
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("SUB EVENT LOG ERR")
			log.Fatal(err)
		case logV := <-logs:
			//fmt.Println(logV) // pointer to event log
			d, _ := contractAbi.Unpack("ItemSet", logV.Data)
			key := d[0].([32]byte)
			value := d[1].([32]byte)
			fmt.Println(string(key[:]), string(value[:]))
		}
	}
}

func newSubLog(client *ethclient.Client, contractAddress common.Address) {
	instance, err := store.NewStore(contractAddress, client)
	if err != nil {
		fmt.Println("LOAD CONTRACT ERR")
		log.Fatal(err)
	}

	filter := instance.StoreFilterer

	logs := make(chan *store.StoreItemSet)
	sub, err := filter.WatchItemSet(&bind.WatchOpts{}, logs)

	for {
		select {
		case err := <-sub.Err():
			fmt.Println("SUB EVENT LOGS ERR")
			log.Fatal(err)
		case logV := <-logs:
			fmt.Printf("EVENT: %s, %s\n", string(logV.Key[:]), string(logV.Value[:]))
		}
	}
}
