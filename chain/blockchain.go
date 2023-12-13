package chain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const (
	DbFile              = "./db/blockchain.db"
	BlocksBucket        = "blocks" // bucket名称，相当于库名
	Miner               = "limit"
	GenesisCoinbaseData = "The TIme 10/Dec/2023 Chancellor on brink onf second bailout for banks"
)

type Blockchain struct {
	// 用来记录最新区块的hash
	tip []byte
	db  *bolt.DB
}

//func dbExists() bool {
//	if _, err := os.Stat(DbFile); os.IsNotExist(err) {
//		return false
//	}
//	return true
//}

func NewBlockchain() *Blockchain {
	// 1. 创建数据库文件，只能第一次创建
	//if dbExists() {
	//	fmt.Println("Blockchain already exists.")
	//	os.Exit(1)
	//}
	var tip []byte
	// 1. 打开数据库文件
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Fatal("Open db fail", err)
	}
	// 3. 更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		// 3.1 获取bucket
		bucket := tx.Bucket([]byte(BlocksBucket))
		if bucket == nil {
			// 3.2.1 第一次使用，创建创世块
			fmt.Println("No existing blockchain found. Creating a new one...")
			cbTx := NewCoinBaseTX(Miner, GenesisCoinbaseData)
			genesis := NewGenesisBlock(cbTx)
			// 3.2.2 区块数据编码
			blockData := genesis.Serialize()
			//3.2.3 创建新bucket，存入区块信息
			var createErr error
			bucket, createErr = tx.CreateBucket([]byte(BlocksBucket))
			if createErr != nil {
				fmt.Println("Create new bucket err", createErr)
				return createErr
			}
			err := bucket.Put(genesis.Hash, blockData)
			if err != nil {
				fmt.Println("Put block data err", err)
				return err
			}
			err = bucket.Put([]byte("last"), genesis.Hash)
			if err != nil {
				fmt.Println("Put block data err", err)
				return err
			}
			tip = genesis.Hash

		} else {
			// 3.3 不是第一次使用
			tip = bucket.Get([]byte("last"))
		}
		return nil
	})
	if err != nil {
		log.Fatal("bucket op err", err)
	}

	return &Blockchain{tip, db}
}

// MinedBlock 挖矿
func (bc *Blockchain) MinedBlock(txs []*Transaction, data string) {
	var tip []byte
	// 1. 获取上一个区块的hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		// 不用判空了
		tip = bucket.Get([]byte("last"))
		return nil
	})

	// 交易信息增加挖矿奖励
	cbtx := NewCoinBaseTX(Miner, data)
	txs = append(txs, cbtx)

	// 利用前块生成新块
	newBlock := NewBlock(txs, tip)
	blockData := newBlock.Serialize()

	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		err := bucket.Put(newBlock.Hash, blockData)
		if err != nil {
			fmt.Println("Bucket put fail", err)
			return err
		}
		err = bucket.Put([]byte("last"), newBlock.Hash)
		if err != nil {
			fmt.Println("Bucket put fail", err)
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Fatal("Block db update err", err)
	}
}

// FindUnspentTransactionsOld 查找账户可以解锁的全部交易
// Deprecated: 不再使用
func (bc *Blockchain) FindUnspentTransactionsOld(address string) []*Transaction {
	var unspentTXs []*Transaction
	// 已经花出的UTXO， 构建tx -> VOutIdx的map
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		// 遍历所有区块
		block, next := bci.PreBlock()

		// 遍历每一个区块中的所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 遍历交易中的输出，
		Outputs:
			for outIdx, out := range tx.VOut {
				// 已经被花出去了，直接跳过此交易
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// 可以被address解锁，代表属于address的utxo在此交易中
				if out.CanBeUnLockWith(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}

			// 用来维护spentTXOs, 已经被引用过了，代表被使用（需要排除挖矿收入）
			if tx.IsCoinBase() == false {
				for _, input := range tx.VIn {
					if input.CanUnlockOutputWith(address) {
						inTxId := hex.EncodeToString(input.TxId)
						spentTXOs[inTxId] = append(spentTXOs[inTxId], input.VOutIdx)
					}
				}
			}
		}

		if !next {
			break
		}
	}
	return unspentTXs
}

// FindUnspentTransactions 查找账户可以解锁的全部交易
func (bc *Blockchain) FindUnspentTransactions(address string) []*Transaction {
	var unspentTXs []*Transaction
	// 已经花出的UTXO， 构建tx -> VOutIdx的map
	spentTXOs := make(map[string]map[int]struct{})
	bci := bc.Iterator()
	for {
		// 遍历所有区块
		block, next := bci.PreBlock()

		// 遍历每一个区块中的所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 遍历交易中的输出，找到没有被输入引用的交易
			for outIdx, out := range tx.VOut {
				// 不是当前发给当前地址的，直接跳过
				if out.CanBeUnLockWith(address) == false {
					continue
				}
				// 已经被花出去了(存在交易输入中)，直接跳过此交易输出
				if _, ok := spentTXOs[txID][outIdx]; ok {
					continue
				}
				// 添加到未花费输出中
				unspentTXs = append(unspentTXs, tx)
			}

			// 挖矿交易，直接跳过输入部分
			if tx.IsCoinBase() {
				continue
			}
			// 用来维护spentTXOs, 已经被引用过了，代表被使用
			for _, input := range tx.VIn {
				if input.CanUnlockOutputWith(address) == false {
					continue
				}
				inTxId := hex.EncodeToString(input.TxId)
				if spentTXOs[inTxId] == nil {
					spentTXOs[inTxId] = make(map[int]struct{})
				}
				spentTXOs[inTxId][input.VOutIdx] = struct{}{}
			}
		}

		if !next {
			break
		}
	}
	return unspentTXs
}

// FindUTXO 查找账户的全部UTXO
func (bc *Blockchain) FindUTXO(address string) []*TXOutput {
	unspentTXs := bc.FindUnspentTransactions(address)
	UTXO := make([]*TXOutput, 0, len(unspentTXs))
	for _, tx := range unspentTXs {
		for _, out := range tx.VOut {
			if out.CanBeUnLockWith(address) {
				tmp := out
				UTXO = append(UTXO, &tmp)
			}
		}
	}
	return UTXO
}

func (bc *Blockchain) FindUnSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentTXs := bc.FindUnspentTransactions(address)
	unspentTXOs := make(map[string][]int)
	value := 0
Loop:
	for _, tx := range unspentTXs {
		txId := hex.EncodeToString(tx.ID)
		for i, output := range tx.VOut {
			if output.CanBeUnLockWith(address) && value < amount {
				unspentTXOs[txId] = append(unspentTXOs[txId], i)
				value += output.Value
				if value > amount {
					break Loop
				}
			}
		}
	}
	return value, unspentTXOs
}

func (bc *Blockchain) GetBalance(address string) int {
	UTXO := bc.FindUTXO(address)
	value := 0
	for _, out := range UTXO {
		value += out.Value
	}
	return value
}

type BlockchainIterator struct {
	currentHash []byte   // 当前区块hash
	db          *bolt.DB // 已打开单位数据库
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockchainIterator) PreBlock() (*Block, bool) {
	if len(i.currentHash) <= 0 {
		return nil, false
	}
	var block *Block
	// 根据hash获取区块
	err := i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		if bucket == nil {
			fmt.Println("Get Bucket nil")
			return errors.New("get bucket nil")
		}
		blockData := bucket.Get(i.currentHash)
		if len(blockData) <= 0 {
			fmt.Println("block data nil")
			return errors.New("block data nil")
		}
		block = DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		log.Fatal("Get Block fail", err)
	}
	// 当前hash变更为前块hash
	i.currentHash = block.PrevBlockHash
	// 返回当前区块
	return block, len(block.PrevBlockHash) > 0
}
