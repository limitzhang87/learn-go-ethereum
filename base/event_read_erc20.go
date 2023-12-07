package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/ierc20"
	"log"
	"math/big"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval ..
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func main() {
	// 查询日志，获取topic

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xAe0d9464095942f8428AeB9BEA293189B4856768")
	instance, err := ierc20.NewIerc20(contractAddress, client)
	if err != nil {
		fmt.Printf("LOAD CONTRACT ERR")
		log.Fatal(err)
	}

	filterer := instance.Ierc20Filterer

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(20),
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		fmt.Println("FILTER LOGS ERR")
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	for _, logV := range logs {
		fmt.Printf("Log Block Number: %d\n", logV.BlockNumber)
		fmt.Printf("Log Index: %d\n", logV.Index)
		fmt.Printf("Log Topic %v\n", logV.Topics)
		// 根据日志中的topic和函数的签名比较
		switch logV.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			// 解析Log
			transferLog, err := filterer.ParseTransfer(logV)
			if err != nil {
				fmt.Println("PARSE TRANSFER ERR")
				log.Fatal(err)
			}
			fmt.Printf(
				"TransferLog From: %s,  equal topic: %t\n",
				transferLog.From.Hex(),
				transferLog.From == common.BytesToAddress(logV.Topics[1].Bytes()),
			)

			fmt.Printf(
				"TransferLog To: %s, equal topic: %t\n",
				transferLog.To.Hex(),
				transferLog.To == common.BytesToAddress(logV.Topics[2].Bytes()), // TOPIC返回的值中，大小是64位的
			)

		case logApprovalSigHash.Hex():
			approvalLog, err := filterer.ParseApproval(logV)
			if err != nil {
				fmt.Println("PARSE APPROVAL ERR")
				log.Fatal(err)
			}
			fmt.Printf(
				"Approval Owner: %s, equal topic: %t\n",
				approvalLog.Owner.Hex(),
				approvalLog.Owner == common.BytesToAddress(logV.Topics[1].Bytes()),
			)

			fmt.Printf(
				"Approval Spender: %s, equal topic: %t\n",
				approvalLog.Spender.Hex(),
				approvalLog.Spender == common.BytesToAddress(logV.Topics[2].Bytes()),
			)
		}
	}
}
