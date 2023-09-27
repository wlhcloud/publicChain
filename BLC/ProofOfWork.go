package BLC

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

// 256位Hash里面前面至少有16个01
const targetBit = 16

type ProofOfWork struct {
	block  *Block   // 当前要验证的区块
	target *big.Int // 大数据存储
}

// IsValid 验证Hash是否有效
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.block.Hash)
	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

// prepareData 返回区块字节数组
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.block.PreBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBit)),
		IntToHex(nonce),
		IntToHex(pow.block.Height),
	}, []byte{})
	return data
}

// NewProofOfWork 创建新的工作证明
func NewProofOfWork(block *Block) *ProofOfWork {
	// 1、创建一个初始值为1的target
	target := big.NewInt(1)
	// 2、左位移256 - targetBit
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}

// Run 挖矿
func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	nonce := 0
	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte
	for {
		// 1、将block属性拼接成为字节数组
		data := proofOfWork.prepareData(int64(nonce))
		// 2、生成Hash
		hash = sha256.Sum256(data)
		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])
		// 判断hashInt是否小雨Block里面的target
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce++
	}
	return hash[:], int64(nonce)
}
