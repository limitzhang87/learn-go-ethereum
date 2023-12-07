package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/store"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	if err != nil {
		fmt.Println("GET PRIVATE KEY ERR")
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	ctAddress := common.HexToAddress("0x665Ab94CaB1c9fD296aAdCbfD9F012C9663F0594")

	instance, err := store.NewStore(ctAddress, client)
	if err != nil {
		fmt.Println("NEW CONTRACT ERR")
		log.Fatal(err)
	}

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
		fmt.Println("NEW TX OPS ERR")
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(50000)

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], "foo")
	copy(value[:], "bar")
	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		fmt.Println("CONTRAST CALL ERR")
		log.Fatal(err)
	}

	fmt.Printf("tx hash: %s\n", tx.Hash().Hex())
}
