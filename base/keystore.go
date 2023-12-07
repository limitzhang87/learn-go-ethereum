package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
)

func main() {
	// 使用keyStore初始化一个钱包
	//ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	//password := "limitzhang"
	//// 使用钱包生成一个账号， 内部还是使用的ecdsa.GenerateKey()
	//account, err := ks.NewAccount(password) // 同时会在wallets中生成账号的信息文件
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("new account private key %s\n", account.Address.Hex()) // 0x6D6afD9331ae6293A661E9803fa7F198299fba94

	// 可以通过账号文件信息导入

	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "limitzhang"
	file := "./wallets/UTC--2023-11-06T18-23-40.832674100Z--6d6afd9331ae6293a661e9803fa7f198299fba94"
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("read file err")
		log.Fatal(err)
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		fmt.Println("keystore import account err")
		log.Fatal(err)
	}

	fmt.Printf("import account address hex: %s\n", account.Address.Hex())
}
