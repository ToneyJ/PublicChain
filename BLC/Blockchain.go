package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

//数据库名
const dbName = "blockchain.db"

//表名
const blockTableName = "blocks"

//链
type Blockchain struct {
	Tlp []byte //最新区块hash
	DB  *bolt.DB
}

//迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tlp, blockchain.DB}
}

//新增区块
func (blockchain *Blockchain) AddBlock(Txs []*Transation) {
	//打开更新数据库
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取最新区块hash
			blockBytes := b.Get(blockchain.Tlp)
			//反序列化
			block := DeserialzeBlock(blockBytes)

			//将区块序列化并且存储到数据库中
			newBlock := NewBlock(Txs, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			//更新blockchain的Tlp
			blockchain.Tlp = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

//创建一个带有创世区块节点的区块
func NewBlockhain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("创世区块已经存在......")
		os.Exit(1)
	}

	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var genisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		//校验
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			//创建数据库表
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		if b != nil {
			//创建创世区块
			txCoinbase := NewCoinbaseTransation(address)
			genisBlock := NewGeneisBlock([]*Transation{txCoinbase})
			//将创世区块存储到表中
			err := b.Put(genisBlock.Hash, genisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新区块的hash
			err = b.Put([]byte("l"), genisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genisHash = genisBlock.Hash
		}

		return nil
	})
	return &Blockchain{genisHash, db}
}

//遍历区块
func (blockChain *Blockchain) Print() {

	blockchainIterator := blockChain.Iterator()

	for {

		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("PrevBlockHash：%x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)
		fmt.Println("TXS:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("TxHash:%x\n", in.TxHash)
				fmt.Printf("vout:%d\n", in.Vout)
				fmt.Printf("from:%s\n", in.PublicKey)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("money:%d\n", out.Value)
				fmt.Printf("to:%s\n", out.Pubkey)
			}
		}

		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//如果一个地址的对应的TXOutput未花费，那么这个Transaction就应该添加到数组中并返回
func (blockchain *Blockchain) Utxos(address string, txs []*Transation) []*UTXO {
	iterator := blockchain.Iterator()

	var unUtxos []*UTXO
	spentTxoutputs := make(map[string][]int)

	for _, tx := range txs {

		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				pubKey := Base58Decode([]byte(address))
				hash160 := pubKey[1 : len(pubKey)-addressChecksumLen]
				if in.UnLockRipmed(hash160) {

					key := hex.EncodeToString(in.TxHash)

					spentTxoutputs[key] = append(spentTxoutputs[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.Vouts {

			if out.UnLockWithAddress(address) {
				if len(spentTxoutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUtxos = append(unUtxos, utxo)
				} else {
					for hash, indexArray := range spentTxoutputs {

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _, outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUtxos = append(unUtxos, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUtxos = append(unUtxos, utxo)
						}
					}
				}

			}

		}

	}

	for {
		block := iterator.Next()

		//txHash
		//Vins
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					pubKey := Base58Decode([]byte(address))
					hash160 := pubKey[1 : len(pubKey)-addressChecksumLen]
					if in.UnLockRipmed(hash160) {
						key := hex.EncodeToString(in.TxHash)
						spentTxoutputs[key] = append(spentTxoutputs[key], in.Vout)
					}
				}
			}

		work:
			//Vout
			for index, out := range tx.Vouts {
				if out.UnLockWithAddress(address) {
					if spentTxoutputs != nil {
						if len(spentTxoutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTxoutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUtxos = append(unUtxos, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUtxos = append(unUtxos, utxo)
						}

					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUtxos
}

//挖掘新的区块
func (blockchain *Blockchain) MineNewBlcok(from, to, amount []string) {
	var txs []*Transation

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	//奖励
	tx := NewCoinbaseTransation(from[0])
	txs = append(txs, tx)
	//1.通过相关算法建立transaction数组

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserialzeBlock(blockBytes)
		}

		return nil
	})

	//建立新区块之前进行验证
	for _, tx := range txs {
		if blockchain.VerifyTransaction(tx) == false {
			log.Panic("用户签名不正确...")
		}
	}

	//2.建立新区块
	block = NewBlock(txs, block.Height+1, block.Hash)
	//将新区块存储到数据库中
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockTableName))
		if bucket != nil {
			bucket.Put(block.Hash, block.Serialize())
			bucket.Put([]byte("l"), block.Hash)
			blockchain.Tlp = block.Hash
		}

		return nil
	})
}

//查找可用的utxo
func (blockchain *Blockchain) FindUTXOs(from string, amount int, txs []*Transation) (int, map[string][]int) {
	//获取所有的UTXO
	utxos := blockchain.Utxos(from, txs)

	var value int
	spendAbleUTXO := make(map[string][]int)
	//遍历utxos
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		TxHash := hex.EncodeToString(utxo.TxHash)
		spendAbleUTXO[TxHash] = append(spendAbleUTXO[TxHash], utxo.Index)
		if value >= amount {
			break
		}
	}

	if value < amount {
		fmt.Printf("%s fund is not enough\n", from)
		os.Exit(1)
	}

	return value, spendAbleUTXO
}

//查询余额
func (blockchain *Blockchain) GetBalance(address string) int {
	utxos := blockchain.Utxos(address, []*Transation{})

	var amount int

	for _, utxo := range utxos {
		amount = amount + utxo.Output.Value
	}

	return amount
}

//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

//返回blockchain对象
func BlockchainObject() *Blockchain {

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))

		}

		return nil
	})

	return &Blockchain{tip, db}
}

//签名
func (blockchain *Blockchain) SignTransaction(tx *Transation, private ecdsa.PrivateKey) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transation)

	for _, vin := range tx.Vins {
		prevTX, err := blockchain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(private, prevTXs)
}

//验证
func (blockchian *Blockchain) VerifyTransaction(tx *Transation) bool {
	prevTxs := make(map[string]Transation)

	for _, vin := range tx.Vins {
		prevTx, err := blockchian.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}

	return tx.Verify(prevTxs)
}

//
func (blockchain *Blockchain) FindTransaction(Txhash []byte) (Transation, error) {
	iterator := blockchain.Iterator()

	for {
		block := iterator.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, Txhash) == 0 {
				return *tx, nil
			}

			var hashInt big.Int
			hashInt.SetBytes(block.PrevBlockHash)

			if big.NewInt(0).Cmp(&hashInt) == 0 {
				break
			}
		}
	}

	return Transation{}, nil
}
