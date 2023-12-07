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
)

func main() {

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x665Ab94CaB1c9fD296aAdCbfD9F012C9663F0594")

	instance, err := store.NewStore(address, client)
	if err != nil {
		fmt.Println("LOAD CONTRACT ERR")
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	if err != nil {
		fmt.Println("GET PRIVATE KEY ERR")
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	opts := &bind.CallOpts{
		From:    fromAddress,
		Context: context.Background(),
	}

	version, err := instance.Version(opts)
	fmt.Printf("VERSION: %s\n", version)
}
