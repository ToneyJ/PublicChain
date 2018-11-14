package BLC

import (
	"fmt"
	"time"
)

type Block struct {
	//区块高度
	Height int64
	//时间戳
	Timestamp int64
	//上一个区块hash
	PrevBlockHash []byte
	//交易数据
	Txs []*Transation
	//当前区块的hash
	Hash []byte
	//随机数Nonce
	Nonce int
}

/*func (block *Block) SetHash() {
	//1.将时间戳转换为字节数组
	//1)将时间戳转化为字符串
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timeStamp := []byte(timeString)
	//2.将除了hash以外的其他属性，以字节数组的形式拼接起来
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timeStamp}, []byte{})
	//3.将拼接起来的数据进行256 hash
	hash := sha256.Sum256(headers)
	//4.将hash赋值给hash小户型
	block.Hash = hash[:]
}*/

//工厂方法
func NewBlock(Tx []*Transation, height int64,prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{height,time.Now().Unix(), prevBlockHash, Tx, []byte{}, 0}

	//将block作为参数，创建一个pow对象
	pow := NewProofWork(block)
	nonce, hash := pow.Run()

	//设置当前区块hash值
	block.Hash = hash[:]

	//设置Nonce
	block.Nonce = nonce

	//校验pow
	isValid := pow.Validata()
	if !pow.Validata() {
		fmt.Println("pow is", isValid)
	}

	//返回区块
	return block
}

//创建创世区块
func NewGeneisBlock(Txs []*Transation) *Block {
	return NewBlock(Txs,1, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
