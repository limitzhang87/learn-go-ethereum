package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 监听新区块
	client, err := ethclient.Dial("ws://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	headers := make(chan *types.Header) // 使用一个通过接受新的区块头信息

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		fmt.Println("SUBSCRIBE NEW HEAD ERR")
		log.Fatal(err)
	}

	fmt.Println("START SUBSCRIBE")
	for {
		select {
		case header := <-headers:
			fmt.Printf("RECEIVE HEAD: %s\n", header.Hash().Hex())
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64())   // 3477413
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		case err := <-sub.Err():
			fmt.Println("SUB HEADER ERR ")
			log.Fatal(err)
		}
	}
}
