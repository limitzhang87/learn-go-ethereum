package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  // 时间戳
	Data          []byte // 数据域
	PrevBlockHash []byte //前块hash
	Hash          []byte // 当前快hash
	Nonce         int64  // 随机数
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers) // [32]byte
	b.Hash = hash[:]               // []byte
}

// Serialize 序列化区块
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	// 编码器
	encoder := gob.NewEncoder(&result)
	// 编码
	_ = encoder.Encode(b)
	return result.Bytes()
}

// DeserializeBlock 反序列化，区块数据还原为block
func DeserializeBlock(d []byte) *Block {
	var block Block
	// 创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(d))
	//解析区块数据
	_ = decoder.Decode(&block) // TODO
	return &block
}

func NewBlock(preHash []byte, data []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: preHash,
		Hash:          nil,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("Genesis Block"))
}
