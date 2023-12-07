package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 循环遍历一个区块总的交易事务，代表一次次的交易
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().Uint64())
		fmt.Println(tx.Gas())
		fmt.Println(tx.GasPrice().Uint64())
		fmt.Println(tx.Nonce())
		fmt.Println(tx.Data())
		fmt.Println(tx.To().Hex())

		// 获取交易中的发送者
		signer := types.NewEIP155Signer(tx.ChainId())
		sender, err := types.Sender(signer, tx)
		if err != nil {
			fmt.Println("get sender from types err")
			log.Fatal(err)
		}
		fmt.Printf("sender: %s\n", sender.String())

		msg, err := core.TransactionToMessage(tx, signer, nil)
		if err != nil {
			fmt.Println("get tx message err")
			log.Fatal(err)
		}
		fmt.Printf("sender: %s\n", msg.From)

		// 获取交易的凭证，判断交易是否成功
		txReceipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println("get transaction receipt err")
			log.Fatal(err)
		}
		fmt.Printf("tx recepit: %d\n", txReceipt.Status)

		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}

	// 根据索引下表获取区块中的交易
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")

	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		fmt.Println("get block err")
		log.Fatal(err)
	}
	fmt.Printf("ts count %d\n", count)

	for i := uint(0); i < count; i++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			fmt.Println("get tx in block err")
			log.Fatal(err)
		}
		fmt.Println(tx.Hash().Hex())
	}

}
