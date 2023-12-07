package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/solc/build/ierc20"
	"log"
)

func main() {
	// 根据IERC20接口生成的go文件，访问链上的ERC20合约

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, _ := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	tokenAddress := common.HexToAddress("0xAe0d9464095942f8428AeB9BEA293189B4856768")

	instance, err := ierc20.NewIerc20(tokenAddress, client)
	if err != nil {
		fmt.Println("NEW ERC20 CONTRACT ERR")
		log.Fatal(err)
	}

	balance, err := instance.BalanceOf(&bind.CallOpts{}, fromAddress)
	if err != nil {
		return
	}
	fmt.Println(balance.String())
}
