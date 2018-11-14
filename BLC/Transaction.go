package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
)

type Transation struct {
	TxHash []byte //交易hash
	Vins   []*TXInput
	Vouts  []*TxOutput
}

//判断当前交易是否是Coinbase交易
func (tx *Transation) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

//创建创世区块时的transation
func NewCoinbaseTransation(address string) *Transation {
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	//txOutput := &TxOutput{10, address}
	txOutput := NewTxoutput(10, address)
	txCoinbase := &Transation{[]byte{}, []*TXInput{txInput}, []*TxOutput{txOutput}}
	txCoinbase.HashTransaction()

	return txCoinbase
}

//转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transation) *Transation {
	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from]

	money, spendableUTXODic := blockchain.FindUTXOs(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TxOutput

	//消费者
	for txhash, indexArry := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txhash)
		for _, index := range indexArry {
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}

	//转账
	//txOutput := &TxOutput{amount, []byte(to)}
	txOutput := NewTxoutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput = NewTxoutput(money-amount, from)
	//txOutput = &TxOutput{money - amount, []byte(from)}
	txOutputs = append(txOutputs, txOutput)

	//交易hash
	tx := &Transation{[]byte{}, txInputs, txOutputs}
	tx.HashTransaction()

	//数字签名
	blockchain.SignTransaction(tx, wallet.PrivateKey)

	return tx
}

//转化hsahTransation
func (tx *Transation) HashTransaction() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}

//进行签名
func (tx *Transation) Sign(privkey ecdsa.PrivateKey, prevTxs map[string]Transation) {
	if tx.IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROE:Previoux transaction is not found!")
		}
	}
	txCopy := tx.TimemedCopy()

	for inHash, vin := range txCopy.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inHash].Signature = nil
		txCopy.Vins[inHash].PublicKey = prevTx.Vouts[vin.Vout].Pubkey
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inHash].PublicKey = nil

		//签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privkey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[inHash].Signature = signature
	}
}

//验证签名
func (tx *Transation) Verify(prevTxs map[string]Transation) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROE:Previoux transaction is not found!")
		}
	}

	txCopy := tx.TimemedCopy()
	curve := elliptic.P256()

	for inHash, vin := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inHash].Signature = nil
		txCopy.Vins[inHash].PublicKey = prevTx.Vouts[vin.Vout].Pubkey
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inHash].PublicKey = nil

		//私钥ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}

//备份transaction
func (tx *Transation) TimemedCopy() Transation {
	var inputs []*TXInput
	var outputs []*TxOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Pubkey})
	}

	txCopy := Transation{tx.TxHash, inputs, outputs}

	return txCopy
}

//Hash
func (tx *Transation) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//序列化
func (tx *Transation) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}
