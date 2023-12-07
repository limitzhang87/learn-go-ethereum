package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

func main() {

	privateKey, err := crypto.GenerateKey() // 生成私钥
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥转为字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// 将私钥字节转为16进制字符串
	privateKeyHex := hexutil.Encode(privateKeyBytes)
	fmt.Printf("private key hex: %s\n", privateKeyHex[2:]) // 使用密钥需要删除前面的0x

	//密钥生成公钥
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	// 公钥转为字节
	publicKeyBytes := crypto.FromECDSAPub(publicKey)
	// 公钥字节转为16进制字符串, 前4位中的0x和04是无效饿
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	fmt.Printf("public key hex: %s len(%d)\n ", publicKeyHex, len(publicKeyHex)-4)

	// 公钥生成ECDSA地址
	address := crypto.PubkeyToAddress(*publicKey)
	fmt.Printf("public key address hex: %s\n", address.Hex())
	// 生成地址的逻辑是Keccak-256哈希，然后我们取最后40个字符（20个字节）并用“0x”作为前缀
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
