package chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

// TXInput 交易输入结构
type TXInput struct {
	TxId       []byte // 应用交易ID
	VOutIdx    int    // 应用的交易输出编号
	PubKeyHash []byte // 公钥hash值
	Signature  []byte // 签名信息
}

// TXOutput 交易输出结构
type TXOutput struct {
	Value      int    // 输出金额
	PubKeyHash []byte // 公钥hash值
}

// CanUnlockOutputWith 判断该输入是否可以被某个账户使用
func (in *TXInput) CanUnlockOutputWith(addr []byte) bool {
	return bytes.Equal(in.PubKeyHash, addr)
}

// CanBeUnLockWith 判断某输出是否可以被账户使用
func (out *TXOutput) CanBeUnLockWith(addr []byte) bool {
	return bytes.Equal(out.PubKeyHash, addr)
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

// TrimmedCopy 修剪交易信息，将输入项的签名和公钥位置空出来
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	// 将原交易内的签名和公钥都置空
	for _, vIn := range tx.VIn {
		inputs = append(inputs, TXInput{TxId: vIn.TxId, VOutIdx: vIn.VOutIdx, PubKeyHash: nil, Signature: nil})
	}

	for _, vOut := range tx.VOut {
		outputs = append(outputs, TXOutput{vOut.Value, vOut.PubKeyHash})
	}

	// 复制一份交易
	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Sign 交易签名
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 1. CoinBase交易无须签名
	if tx.IsCoinBase() {
		return
	}

	// 2. 修剪交易
	txCopy := tx.TrimmedCopy()
	// 3. 循环向输入项签名
	for inID, vin := range txCopy.VIn {
		// 找到输入项引用的交易
		prevTx := prevTXs[hex.EncodeToString(vin.TxId)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKeyHash = prevTx.VOut[vin.VOutIdx].PubKeyHash
		txCopy.SetId()
		// txID生成后吧PubKey置空
		txCopy.VIn[inID].PubKeyHash = nil
		// 使用ecsda签名获得r和s
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		// 签名签名数据
		signature := append(r.Bytes(), s.Bytes()...)
		tx.VIn[inID].Signature = signature
	}
}

// Subsidy 定义奖励数量
const Subsidy = 10

// NewCoinBaseTX 创建CoinBase交易，CoinBase就是挖矿奖励
func NewCoinBaseTX(to []byte, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 创建一个输入项
	txIn := TXInput{TxId: []byte{}, VOutIdx: -1}
	// 创建输出项
	txOut := TXOutput{Value: Subsidy, PubKeyHash: to}
	tx := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
	tx.SetId()
	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].VOutIdx == -1
}

// NewUTXOTransaction 新建交易
func NewUTXOTransaction(from, to []byte, amount int, bc *Blockchain) (*Transaction, error) {
	// 1. 获取未花费的交易输出
	//acc, validOutput := bc.FindUnSpendableOutputs(from, amount)
	unspentTxs := bc.FindUnspentTransactions(from)
	acc, validOutput := bc.FindUnSpendableOutputsWithTX(from, amount, unspentTxs)

	if acc < amount {
		return nil, errors.New("no enough funds")
	}

	input := make([]TXInput, 0, len(validOutput))
	output := make([]TXOutput, 0, 2)
	// 将前面的未花费交易输出变为新交易的交易输入
	for txId, ids := range validOutput {
		txIdByte, err := hex.DecodeString(txId)
		if err != nil {
			fmt.Println("decode tx id err")
			return nil, err
		}
		for _, idx := range ids {
			input = append(input, TXInput{
				TxId:       txIdByte,
				VOutIdx:    idx,
				PubKeyHash: from,
			})
		}
	}

	output = append(output, TXOutput{
		Value:      amount,
		PubKeyHash: to,
	})
	if acc > amount {
		// 剩余的币要还给发送帐号
		output = append(output, TXOutput{
			Value:      acc - amount,
			PubKeyHash: from,
		})
	}

	tx := &Transaction{
		ID:   nil,
		VIn:  input,
		VOut: output,
	}
	tx.SetId()

	unspentTxMap := make(map[string]*Transaction, len(unspentTxs))
	for _, unspentTx := range unspentTxs {
		txId := hex.EncodeToString(unspentTx.ID)
		unspentTxMap[txId] = unspentTx
	}

	wallet, err := NewWallet("./keystore")
	if err != nil {
		return nil, err
	}
	tx.Sign(wallet.HDKeyStore.Key.PrivateKey, unspentTxMap)
	return tx, nil
}
