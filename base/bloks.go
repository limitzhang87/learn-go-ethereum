package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	header, err := client.HeaderByNumber(ctx, nil) // 获取区块头部信息，第二个参数传入空则获取最新的区块
	if err != nil {
		fmt.Println("get header err")
		log.Fatal(err)
	}

	fmt.Printf("block_id: %s\n", header.Number.String())

	blockId := header.Number

	// 根据区块ID查询区块全部信息
	block, err := client.BlockByNumber(ctx, blockId)
	if err != nil {
		fmt.Println("get block by number err")
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))

	// 获取一个区块的交易数目

	tsCount, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		fmt.Println("get block transaction count err")
		log.Fatal(tsCount)
	}
	fmt.Printf("transaction count: %d\n", tsCount)
}
