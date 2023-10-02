package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// BlockchainIterator 区块链迭代对象
type BlockchainIterator struct {
	currentHash []byte   // 当前区块的hash
	DB          *bolt.DB // 数据库
}

// Next 获取上一个区块
func (blockchainIterator *BlockchainIterator) Next() *Block {
	block := &Block{}
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockTableName))
		if bucket == nil {
			return nil
		}
		// 获取当前区块，通过当前区块获取上一个区块哈说
		blockBytes := bucket.Get(blockchainIterator.currentHash)
		block = DeSerializeBlock(blockBytes)
		// 更新当前hash
		blockchainIterator.currentHash = block.PreBlockHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
