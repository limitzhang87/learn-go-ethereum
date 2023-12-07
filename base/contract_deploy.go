package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/store"
	"log"
	"math/big"
)

func main() {
	// 使用go部署智能合约
	/*
			solc --abi --bin -o ./solc/source ./solc/contract/Store.sol // 编译合约文件
			abigen --bin=./solc/source/Store.bin --abi=./solc/source/Store.abi --pkg=store --out=./solc/build/store/Store.go // 生成go文件


		问题：
			1. 使用solidity 0.8.22, 会报错 VM Exception while processing transaction: invalid opcode
				这是因为0.8.22之后 820的目标vm是shanghai以后的evm，增加了push0这个操作码。我看你贴的图用的是ganache cli，它应该还是shanghai之前的evm，不支持push0，所以出错日志中有invalid opercode。push0操作码在eip3855中规定，是0x5f，它之前是invalid opercode。
			简单说，820编译成的目标代码会出现0x5f，如果它跑在shanghai之前的evm中就会出你所遇见的错误
			2. Runtime error: code size to deposit exceeds maximum code size
				这是因为gaslimit不足导致的
	*/
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	if err != nil {
		fmt.Println("LOAD PRIVATE KEY ERR")
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		fmt.Println("GET PENDING NONCE ERR")
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("GET GAS PRICE ERR")
		log.Fatal(err)
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		fmt.Println("GET CHAIN ID ERR")
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		fmt.Println("NEW KEY TX ERR")
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(1000000)

	input := "1.0"                                                       // 合约构造参数
	address, tx, instance, err := store.DeployStore(auth, client, input) // 每一次执行都会重复部署
	if err != nil {
		fmt.Println("DEPLOY CONTRACT ERR")
		log.Fatal(err)
	}
	fmt.Println(address.Hex())   // 0x665Ab94CaB1c9fD296aAdCbfD9F012C9663F0594
	fmt.Println(tx.Hash().Hex()) // 0x425126a7223f2f0415c5055cd6334e5bf46d6a02eb0d6bdf51b18bf6f06225bb
	_ = instance
}
