package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//迭代
type BlockchainIterator struct {
	CurrentHash []byte //当前区块hash
	DB          *bolt.DB
}

//获取下一个区块
func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block
	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			//获取迭代器所对应的区块
			block = DeserialzeBlock(currentBlockBytes)
			//更新迭代器的currenthash
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return block
}

