package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/howeyc/gopass"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
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

// Wallet 钱包结构体
type Wallet struct {
	Address    common.Address
	HDKeyStore *HDKeyStore
}

// NewWallet 钱包构造函数
func NewWallet(keyPath string) (*Wallet, error) {
	// 1. 创建助记词
	mn, err := createMnemonic()
	if err != nil {
		fmt.Println("Failed to NewWallet", err)
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
	publicKey, err := GetPublicKey(privateKey)
	if err != nil {
		fmt.Println("failed to GetPublicKey", err)
		return nil, err
	}
	// 利用公钥推导私钥
	address := crypto.PubkeyToAddress(*publicKey)
	// 4. 创建keystore
	hdKs := NewHDKeyStore(keyPath, privateKey)
	// 5. 创建钱包
	return &Wallet{address, hdKs}, nil
}

func LoadWallet(filename, dataDir string) (*Wallet, error) {
	hdKs := NewHDKeyStoreNoKey(dataDir)
	// 解决密码问题
	fmt.Println("Please input password for:", filename)
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println("get pass err", err)
		return nil, err
	}
	address := common.HexToAddress(filename) // 文件名就是账户地址
	_, err = hdKs.GetKey(address, hdKs.JoinPath(filename), string(pass))
	if err != nil {
		fmt.Println("failed to get key", err)
		return nil, err
	}
	//hdKs.Key = key 在GetKey中已经执行过这句了
	return &Wallet{
		Address:    address,
		HDKeyStore: hdKs,
	}, nil

}

func (w *Wallet) StoreKey(pass string) error {
	// 账户即文件名
	filename := w.HDKeyStore.JoinPath(w.Address.Hex())
	return w.HDKeyStore.StoreKey(filename, w.HDKeyStore.Key, pass)
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
	return GetPrivateKey(path, masterKey)
}

func GetPrivateKey(path accounts.DerivationPath, masterKey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	var err error
	key := masterKey
	for _, n := range path {
		// 按照路径跌倒获得最终key
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}
	}
	// 将key转换为ecdsa私钥
	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return privateKeyECDSA, nil
}

func GetPublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return publicKeyECDSA, nil
}

// =====================================

const ChecksumLen = 4

const version = byte(0x00)

// HashPublicKey 计算公钥 hash
func HashPublicKey(publicKey []byte) []byte {
	// 1. 先hash一次
	publicSHA256 := sha3.Sum256(publicKey)
	// 2. 计算 ripemed160
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func GetAddress(publicKey []byte) []byte {
	// 1. 计算公钥 hash
	publicKeyHash := HashPublicKey(publicKey)
	// 2. 计算校验和
	versionPayload := append([]byte{version}, publicKeyHash...)
	checksum := checksum(versionPayload)
	// 3. 计算base58编码
	fullPayload := append(versionPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:ChecksumLen]
}
