package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const (
	MaxNonce   = math.MaxInt64
	TargetBits = 16 // 难度值
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	// target 为最终安度值
	target := big.NewInt(1)
	// target 为1想做移256-24（w挖矿难度）
	target.Lsh(target, uint(256-TargetBits))
	// 生成pow结构
	pow := &ProofOfWork{block: b, target: target}
	return pow
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0
	fmt.Printf("Mining the block containing %s, maxNonce=%d\n", pow.block.Data, MaxNonce)
	for nonce < MaxNonce {
		// 数据准备
		data := pow.prepareData(nonce)
		// 计算hash
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x %x", hash, pow.target)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		Int2Hex(pow.block.Timestamp),
		Int2Hex(int64(TargetBits)),
		Int2Hex(nonce),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Validate() bool {
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

func Int2Hex(num int64) []byte {
	buff := new(bytes.Buffer)
	// 大端法写入
	_ = binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
