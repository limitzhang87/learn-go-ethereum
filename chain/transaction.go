package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// TXInput 交易输入结构
type TXInput struct {
	TxId     []byte // 应用交易ID
	VOutIdx  int    // 应用的交易输出编号
	FromAddr string // 输入方验签， 类似于上一个交易中TxOutput的ToAddr
}

// TXOutput 交易输出结构
type TXOutput struct {
	Value  int    // 输出金额
	ToAddr string // 输出方验签
}

// CanUnlockOutputWith 判断该输入是否可以被某个账户使用
func (in *TXInput) CanUnlockOutputWith(addr string) bool {
	return in.FromAddr == addr
}

// CanBeUnLockWith 判断某输出是否可以被账户使用
func (out *TXOutput) CanBeUnLockWith(addr string) bool {
	return out.ToAddr == addr
}

type Transaction struct {
	ID   []byte     // 交易ID
	VIn  []TXInput  // 交易输入项
	VOut []TXOutput // 交易输出项
}

// SetId 将交易信息转为hash， 并设为ID
func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	_ = enc.Encode(tx)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Subsidy 定义奖励数量
const Subsidy = 10

// NewCoinBaseTX 创建CoinBase交易，CoinBase就是挖矿奖励
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 创建一个输入项
	txIn := TXInput{TxId: []byte{}, VOutIdx: -1, FromAddr: data}
	// 创建输出项
	txOut := TXOutput{Value: Subsidy, ToAddr: to}
	tx := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
	tx.SetId()
	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].VOutIdx == -1
}
