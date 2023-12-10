package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64          // 时间戳
	Transactions  []*Transaction // 交易信息
	PrevBlockHash []byte         //前块hash
	Hash          []byte         // 当前快hash
	Nonce         int64          // 随机数
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers) // [32]byte
	b.Hash = hash[:]               // []byte
}

func (b *Block) HashTransactions() []byte {
	hashes := make([][]byte, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.ID)
	}
	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
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

func NewBlock(tx []*Transaction, preHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  tx,
		PrevBlockHash: preHash,
		Hash:          nil,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return block
}

func NewGenesisBlock(coinBase *Transaction) *Block {
	return NewBlock([]*Transaction{coinBase}, []byte{})
}
