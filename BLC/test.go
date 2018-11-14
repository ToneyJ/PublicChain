package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type Iterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (it *Iterator) Next() *Block {
	var block *Block
	err := it.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockbytes := b.Get(it.CurrentHash)
			block = DeserialzeBlock(blockbytes)
			it.CurrentHash = block.PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}

type CLi struct {
}
