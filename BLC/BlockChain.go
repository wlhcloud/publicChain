package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"time"
)

const dbName = "blockchain.db"  // 数据库名称
const blockTableName = "blocks" //表名称

// Blockchain 区块链
type Blockchain struct {
	Tip []byte   // 最新区块的hash
	DB  *bolt.DB // 数据库
}

// Iterator 迭代所有区块
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

// PrintChain 打印所有区块
func (blockchain *Blockchain) PrintChain() {
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("PreBlockHash：%x\n", block.PreBlockHash)
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		fmt.Println()

		// 判断是否是第一个区块
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// AddBlockToBlockchain 添加区块
func (blockchain *Blockchain) AddBlockToBlockchain(data string) {

	// 添加到区块链
	// 创建表
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		// 往表中存储数据
		if b != nil {
			// 获取最新区块
			lastBlock := b.Get(blockchain.Tip)
			serializeBlock := DeSerializeBlock(lastBlock)

			// 创建新区块
			block := NewBlock(data, serializeBlock.Height+1, blockchain.Tip)
			// 新区块存储到表中
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic("新区块存储到表中失败")
			}
			// 存储最新的hash
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic("存储最新的hash失败")
			}

			blockchain.Tip = block.Hash
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// CreateBlockchainWithGenesis 创建带有创世区块的区块链
func CreateBlockchainWithGenesis() *Blockchain {
	var blockHash []byte // 存储最新区块hash

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 创建表
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		// 往表中存储数据
		if b != nil {
			// 创建创世区块
			genesisBlock := CreateGenesisBlock("Genesis Block")

			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(" 将创世区块存储到表中存储失败")
			}
			// 存储最新的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic("存储最新的hash失败")
			}

			blockHash = genesisBlock.Hash
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})

	// 返回区块链对象
	return &Blockchain{blockHash, db}
}
