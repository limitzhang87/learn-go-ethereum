package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("2ecd5a140398e2df5e0d7b928aee5fdd24b1e8c9a5fe739331423409c50256af")
	if err != nil {
		fmt.Println("get private key error")
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	tokenAddress := common.HexToAddress("0xAe0d9464095942f8428AeB9BEA293189B4856768")
	//toAddress := common.HexToAddress("0xE52C8A1A9aA9b6a3e28F278Af3F8A84a8a4f738e")

	//callMint(client, privateKey, tokenAddress)
	//err = callTransfer(client, privateKey, tokenAddress, toAddress)
	err = callBalance(client, privateKey, tokenAddress, fromAddress)

	//if err != nil {
	//	fmt.Println("call transfer err")
	//	log.Fatal(err)
	//}
}

func getFnSignature(method string) []byte {
	hash := sha3.NewLegacyKeccak256()
	mintFnSignature := []byte(method)
	hash.Reset()
	hash.Write(mintFnSignature)
	return hash.Sum(nil)[:4] // 取前四位代表函数的签名

	//return crypto.Keccak256([]byte(method))[:4]  等价于上面
}

func callMint(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress common.Address) {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	value := big.NewInt(0)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("get pending nonce err")
		log.Fatal(err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("get work id(chain id) err")
		log.Fatal(err)
	}

	mintSg := getFnSignature("mint(uint256)")

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	padAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, mintSg...)
	data = append(data, padAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		fmt.Println("get gas limit err")
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("get gas price err")
		log.Fatal(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &tokenAddress,
		Value:    value,
		Data:     data,
		V:        nil,
		R:        nil,
		S:        nil,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		fmt.Println("signed tx err")
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("send tx err")
		log.Fatal(err)
	}
}

func callTransfer(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, toAddress common.Address) error {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	value := big.NewInt(0)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("get pending nonce err")
		return err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("get chain id err")
		return err
	}

	var data []byte

	sgTx := getFnSignature("transfer(address,uint256)")
	padToAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount := new(big.Int)
	amount.SetString("100000000000000000000", 10) // 100 tokens
	padAmount := common.LeftPadBytes(amount.Bytes(), 32)

	data = append(data, sgTx...)
	data = append(data, padToAddress...)
	data = append(data, padAmount...)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("get gas price err")
		return err
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &tokenAddress,
		Value:    value,
		Data:     data,
		V:        nil,
		R:        nil,
		S:        nil,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		fmt.Println("sign tx err")
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	return err
}

func callBalance(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, userAddress common.Address) error {
	ctx := context.Background()
	var err error

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	sgBalance := getFnSignature("balanceOf(address)")
	padAddress := common.LeftPadBytes(userAddress.Bytes(), 32)

	var data []byte
	data = append(data, sgBalance...)
	data = append(data, padAddress...)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		fmt.Println("GET PENDING NONCE ERR")
		return err
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("GET PAS PRICE ERR")
		return err
	}
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &tokenAddress,
		Value:    big.NewInt(0),
		Data:     data,
		V:        nil,
		R:        nil,
		S:        nil,
	})

	chainId, err := client.ChainID(ctx)
	if err != nil {
		fmt.Println("GET CHAIN ID ERR")
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		fmt.Println("sign tx err")
		return err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		fmt.Println("send transaction err")
		return err
	}

	return nil
}
