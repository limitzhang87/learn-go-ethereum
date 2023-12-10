package chain

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
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

func dbExists() bool {
	if _, err := os.Stat(DbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

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

// AddBlock 添加区块
func (bc *Blockchain) AddBlock(txs []*Transaction) {
	var tip []byte
	// 1. 获取上一个区块的hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		// 不用判空了
		tip = bucket.Get([]byte("last"))
		return nil
	})

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
