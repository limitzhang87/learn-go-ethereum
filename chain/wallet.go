package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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
