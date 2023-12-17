package chain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
)

type HDKeyStore struct {
	keyDirPath string       // 文件所在路径
	scryptN    int          // 生成加密文件的参数N
	scryptP    int          // 生成加密文件的参数p
	Key        keystore.Key // HDKeyStore对应的key
}

func NewHDKeystore(path string, privateKey *ecdsa.PrivateKey) *HDKeyStore {
	// 获取UUID
	uuid1 := NewRandom()
	key := keystore.Key{
		Id:         uuid.UUID(uuid1),
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
	//keystore.StandardScryptN, keystore.StandardScryptP
	return &HDKeyStore{
		keyDirPath: path,
		scryptN:    keystore.StandardScryptN, // 这两个会影响加密时间，从而影响生成钱包数据的时间
		scryptP:    keystore.StandardScryptP,
		Key:        key,
	}
}

type UUID [16]byte

// 全局加密随机阅读器

// 生成UUID

func NewRandom() UUID {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		log.Fatal("new random err", err)
	}
	// 版本4规范处理与变形
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return [16]byte(uuid)
}

func (ks *HDKeyStore) StoreKey(filename string, key *keystore.Key, auth string) error {
	keyJson, err := keystore.EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		return err
	}
	return WriteKeyFile(filename, keyJson)
}

func WriteKeyFile(filename string, context []byte) error {
	const dirPerm = 0700
	_, err := os.Stat(filepath.Dir(filename))
	if os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(filename), dirPerm); err != nil {
			return err
		}
	}

	f, err := os.CreateTemp(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}

	if _, err := f.Write(context); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return err
	}
	_ = f.Close()
	return os.Rename(f.Name(), filename)
}

func (ks *HDKeyStore) JoinPath(filename string) string {
	//如果filename是绝对路径，则直接返回
	if filepath.IsAbs(filename) {
		return filename
	}

	// 将路径与文件拼接
	return filepath.Join(ks.keyDirPath, filename)
}

func (ks *HDKeyStore) GetKey(addr common.Address, filename, auth string) (*keystore.Key, error) {
	// 读取文件内容
	keyJoin, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 利用以太坊DecryptKey解码json文件
	key, err := keystore.DecryptKey(keyJoin, auth)
	if err != nil {
		return nil, err
	}

	// 如果地址不同代表解析失败
	if key.Address != addr {
		return nil, fmt.Errorf("key context mismath: have account%x, want %x", key.Address, addr)
	}
	ks.Key = *key
	return key, nil
}
