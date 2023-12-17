package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"log"
)

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256() // 椭圆曲率
	// 生成私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 利用私钥推导出公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, pubKey
}

// Wallet 钱包结构
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey // 私钥
	PublicKey  []byte            // 公钥
}

// NewWallet 创建钱包
func NewWallet() *Wallet {
	// 随机生成密钥对
	privateKey, publicKey := newKeyPair()
	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	return wallet
}

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

const ChecksumLen = 4

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:ChecksumLen]
}

const version = byte(0x00)

func (w *Wallet) GetAddress() []byte {
	// 1. 计算公钥 hash
	publicKeyHash := HashPublicKey(w.PublicKey)
	// 2. 计算校验和
	versionPayload := append([]byte{version}, publicKeyHash...)
	checksum := checksum(versionPayload)
	// 3. 计算base58编码
	fullPayload := append(versionPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

// DeriveAddressFromMnemonic 根据助记词和密码反推私钥
func DeriveAddressFromMnemonic(mnemonic string, password string, index int) {
	// 1. 推导路径
	/*
		m/44'  提案编号，39,44
		/60' 币种 比特比(0), 以太坊(60)
		/0' 逻辑性亚账户
		/0 HD钱包两个压树， 一个用来接收地址，一个用来找零
		/1 地址编号
	*/

	path := fmt.Sprintf("m/44'/60'/0'/0/%d", index)
	dPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 获取种子
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	// 3. 获取主key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.Params{})
	if err != nil {
		fmt.Println("Failed to NewMaster", err)
		return
	}
	// 4. 推导私钥
	privateKey, err := DerivePrivateKey(dPath, masterKey)
	if err != nil {
		log.Fatal(err)
	}
	// 5. 推导公钥
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 6. 利用公钥推导私钥
	address := crypto.PubkeyToAddress(*publicKey)
	fmt.Println(address.Hex())
}

func DerivePrivateKey(path accounts.DerivationPath, masterKey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
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

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return publicKeyECDSA, nil
}
