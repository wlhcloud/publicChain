package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 1、区块高度
	Height int64
	// 2、上一个区块Hash
	PreBlockHash []byte
	// 3、交易数据
	Data []byte
	// 4、 时间戳
	Timestamp int64
	// 5、Hash
	Hash []byte
	// 6、Nonce
	Nonce int64
}

// Serialize 序列化区块
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeSerializeBlock 反序列化
func DeSerializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// NewBlock 创建新的区块
func NewBlock(data string, height int64, preBlockHash []byte) *Block {
	// 创建区块
	block := Block{
		height, preBlockHash, []byte(data), time.Now().Unix(), nil, 0,
	}
	// 调用工作量证明的方法并返回有效hash 和nonce
	pow := NewProofOfWork(&block)
	// 挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// CreateGenesisBlock 创建一个创世区块
func CreateGenesisBlock(data string) *Block {
	preBlockHash := make([]byte, 32)
	return NewBlock(data, 1, preBlockHash)
}
