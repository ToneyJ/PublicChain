package BLC

import (
	"bottos/common"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	//定义最大值
	maxNonce = math.MaxInt64
)

const targetBits = 16 //控制区块难度

type ProofOfWork struct {
	block  *Block   //当前需要验证区块
	target *big.Int //大数据存储,用作为挖矿难度值
}
//数据拼接，返回数组
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.SerialzeTransation(),
		IntToHex(int64(pow.block.Timestamp)),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
		IntToHex(pow.block.Height)},
		[]byte{})

	return data
}

//运行工作量证明机制
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < maxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		//转化为big.int类型
		hashInt.SetBytes(hash[:])
		//pow.target 和 hash 比较
		//hashInt<pow.target -1
		// == 0
		//> 1
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()

	return nonce, hash[:]
}

//创建新的工作量证明机制
func NewProofWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block, target}

	return pow
}

//验证当前工作量证明的有效性
func (pow *ProofOfWork) Validata() bool {
	var hashInt big.Int
	data := pow.PrepareData(pow.block.Nonce)
	hash := common.Sha256(data)
	hashInt.SetBytes(hash[:])

	isVaild := hashInt.Cmp(pow.target) == -1

	return isVaild
}
