package BLC

import (
	"bytes"
)

type TxOutput struct {
	Value  int    //交易金额
	Pubkey []byte //交易用户公钥
}

func (txOutput *TxOutput) Lock(address string) {
	pubKey := Base58Decode([]byte(address))
	txOutput.Pubkey = pubKey[1 : len(pubKey)-addressChecksumLen]
}

func NewTxoutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{value, []byte(address)}
	//设置Pubkey
	txOutput.Lock(address)

	return txOutput
}

//判断当前的utxo拥有者是谁
func (txOutput *TxOutput) UnLockWithAddress(address string) bool {
	pubKey := Base58Decode([]byte(address))
	hash160 := pubKey[1 : len(pubKey)-addressChecksumLen]

	return bytes.Compare(txOutput.Pubkey, hash160) == 0
}
