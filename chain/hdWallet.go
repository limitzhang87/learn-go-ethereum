package chain

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
)

const DefaultDerivationPath = "m/40'/60'/0'/0/1"

func createMnemonic() (string, error) {
	// Entropy 生成， 注意传入值y=32*x, 并且 128<=y<=256
	b, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	// 生成助记词
	return bip39.NewMnemonic(b)
}

// HDWallet 钱包结构体
type HDWallet struct {
	Address   common.Address
	HDKeytory *HDKeyStore
}

// NewHDWallet 钱包构造函数
func NewHDWallet(keyPath string) (*HDWallet, error) {
	// 1. 创建助记词
	mn, err := createMnemonic()
	if err != nil {
		fmt.Println("Failed to NewHDWallet", err)
		return nil, err
	}
	fmt.Println(mn)
	// 2. 推导私钥
	privateKey, err := NewKeyFromMnemonic(mn)
	if err != nil {
		fmt.Println("failed to NewKeyFromMnemonic", err)
		return nil, err
	}
	// 3. 获取地址
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		fmt.Println("failed to DerivePublicKey", err)
		return nil, err
	}
	// 利用公钥推导私钥
	address := crypto.PubkeyToAddress(*publicKey)
	// 4. 创建keystore
	hdKs := NewHDKeystore(keyPath, privateKey)
	// 5. 创建钱包
	return &HDWallet{address, hdKs}, nil
}

func (w *HDWallet) StoreKey(pass string) error {
	// 账户即文件名
	filename := w.HDKeytory.JoinPath(w.Address.Hex())
	return w.HDKeytory.StoreKey(filename, &w.HDKeytory.Key, pass)
}

func NewKeyFromMnemonic(mn string) (*ecdsa.PrivateKey, error) {
	//1. 推导目录
	path, err := accounts.ParseDerivationPath(DefaultDerivationPath)
	if err != nil {
		log.Panic("Failed to ParseDerivationPath ", err)
	}
	//2. 通过助记词生成种子
	//NewSeedWithErrorChecking(mnemonic string, password string) ([]byte, error)
	seed, err := bip39.NewSeedWithErrorChecking(mn, "")
	if err != nil {
		log.Panic("Failed to NewSeedWithErrorChecking ", err)
	}
	//3. 获得主key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Panic("Failed to NewMaster", err)
	}
	//4. 推导私钥
	return DerivePrivateKey(path, masterKey)
}
