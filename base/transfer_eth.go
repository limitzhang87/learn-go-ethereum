package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	// crypto加载密钥时，不能有Ox,
	privateKey, err := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// 获取转账的随机数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("get nonce err")
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("get gas price err")
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xE52C8A1A9aA9b6a3e28F278Af3F8A84a8a4f738e")

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	})

	// 获取最新chain id
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("get current chain id err")
		log.Fatal(err)
	}

	// 使用交易信息，chainId, 密钥，生成签名后的交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		fmt.Println("send transaction err")
		log.Fatal(err)
	}

	fmt.Printf("tx send %s\n", signedTx.Hash().Hex()) // 0xa2b946d4befddece638bab2a1c3ae5d1b90d0b482253c627cd53fd8b339c3c22
}
