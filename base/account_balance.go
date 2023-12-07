package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

func main() {
	// 读取账户余额
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}
	adr := "0x2c2059b05CfF9Fde6027298b2cABE365BcF74DE3"
	ctx := context.Background()
	address := common.HexToAddress(adr)
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		_ = fmt.Errorf("get balance err: %s\n", err.Error())
		log.Fatal(err)
	}
	fmt.Printf("account:%s, balance %d\n", adr, balance)

	// 读取余额是，传入区块号，返回读取该区块时的账户余额
	blockId, err := client.BlockNumber(ctx)
	if err != nil {
		log.Printf("get last block number err")
		log.Fatal(err)
	}
	fmt.Printf("last block id: %d\n", blockId)
	balance, err = client.BalanceAt(ctx, address, big.NewInt(int64(blockId)))
	if err != nil {
		_ = fmt.Errorf("get balance(by block) err: %s\n", err.Error())
		log.Fatal(err)
	}
	fmt.Printf("account:%s, balance %d\n", adr, balance)

	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("balance of eth  %.18f\n", ethValue)

	// 带确认余额，在提交或等待交易确认
	pendingBalance, err := client.PendingBalanceAt(ctx, address)
	if err != nil {
		fmt.Println("get pending balance err")
		log.Fatal(err)
	}
	fmt.Printf("pending balance %d\n", pendingBalance)
}
